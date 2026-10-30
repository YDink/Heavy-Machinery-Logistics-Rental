package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
)

type AuthService struct {
	userRepo     domain.UserRepository
	otpService   *OTPService
	tokenService *TokenService
}

func NewAuthService(
	userRepo domain.UserRepository,
	otpService *OTPService,
	tokenService *TokenService,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		otpService:   otpService,
		tokenService: tokenService,
	}
}

// RequestOTPInput holds parameters for initiating an OTP challenge.
type RequestOTPInput struct {
	PhoneNumber    string
	CountryCode    string
	Purpose        domain.OTPPurpose
}

type RequestOTPOutput struct {
	VerificationSessionID string
	TimeoutSeconds        int
}

// RequestOTP parses MSISDN and issues an OTP challenge (REG-01, AUT-01).
func (s *AuthService) RequestOTP(ctx context.Context, input RequestOTPInput) (*RequestOTPOutput, error) {
	phone, err := domain.ParsePhoneNumber(input.PhoneNumber, input.CountryCode)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidPhoneNumber, err)
	}

	sessionID, timeout, err := s.otpService.RequestOTP(ctx, phone, input.Purpose)
	if err != nil {
		return nil, err
	}

	return &RequestOTPOutput{
		VerificationSessionID: sessionID,
		TimeoutSeconds:        timeout,
	}, nil
}

// VerifyOTPInput holds parameters for verifying an OTP challenge.
type VerifyOTPInput struct {
	VerificationSessionID string
	Code                  string
}

// VerifyOTP verifies the challenge and issues tokens (AUT-02, AUT-03, AUT-06).
func (s *AuthService) VerifyOTP(ctx context.Context, input VerifyOTPInput) (*domain.TokenPair, *domain.User, error) {
	challenge, err := s.otpService.VerifyOTP(ctx, input.VerificationSessionID, input.Code)
	if err != nil {
		return nil, nil, err
	}

	// Lookup or create user
	user, err := s.userRepo.GetByPhoneNumber(ctx, challenge.Phone)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// Create new user with active phone verified status
			phone, parseErr := domain.ParsePhoneNumber(challenge.Phone, "")
			if parseErr != nil {
				return nil, nil, parseErr
			}

			now := time.Now()
			newUser := &domain.User{
				ID:              uuid.New(),
				PhoneNumber:     phone,
				CountryCode:     phone.CountryCode(),
				FullName:        "New User",
				PrimaryRole:     domain.RoleRenter,
				AccountStatus:   domain.AccountStatusActive,
				KYCStatus:       domain.KYCStatusUnverified,
				TokenVersion:    1,
				PhoneVerifiedAt: &now,
			}

			createdUser, createErr := s.userRepo.Create(ctx, newUser)
			if createErr != nil {
				return nil, nil, fmt.Errorf("failed to create user after OTP verification: %w", createErr)
			}
			user = createdUser
		} else {
			return nil, nil, err
		}
	} else {
		// Mark phone verified if not yet verified
		if user.PhoneVerifiedAt == nil || user.AccountStatus == domain.AccountStatusUnverified {
			updated, updateErr := s.userRepo.MarkPhoneVerified(ctx, user.ID)
			if updateErr == nil {
				user = updated
			}
		}
	}

	// Assert account is not frozen
	if user.AccountStatus == domain.AccountStatusSuspendedSecurityFreeze ||
		user.AccountStatus == domain.AccountStatusTerminated {
		return nil, nil, domain.ErrAccountSuspended
	}

	// Mint token pair
	tokenPair, err := s.tokenService.MintTokenPair(ctx, user, "")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to mint session tokens: %w", err)
	}

	return tokenPair, user, nil
}

// RefreshSession exchanges a refresh token for a fresh token pair (ADR-014).
func (s *AuthService) RefreshSession(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	// Look up token family metadata to identify user
	metadata, err := s.tokenService.tokenStore.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, metadata.UserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	if user.AccountStatus == domain.AccountStatusSuspendedSecurityFreeze ||
		user.AccountStatus == domain.AccountStatusTerminated {
		return nil, domain.ErrAccountSuspended
	}

	return s.tokenService.RotateRefreshToken(ctx, refreshToken, user)
}

// GetUserProfile returns the authenticated user's profile.
func (s *AuthService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
