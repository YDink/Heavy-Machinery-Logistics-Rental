package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
)

const (
	DefaultOTPTTL             = 300 * time.Second // 5 minutes (AUT-02)
	MaxOTPVerificationAttempts = 3                // Max 3 attempts before challenge invalidation (AUT-03)
	RateLimitWindow           = 15 * time.Minute  // 15-minute sliding window (D-05-01-001)
	MaxDispatchesPerWindow    = 3                 // Max 3 OTP dispatches per window
)

type OTPServiceConfig struct {
	MockOTP bool // If true, always generates "123456" for dev/staging test accounts
}

type OTPService struct {
	otpRepo    domain.OTPRepository
	smsGateway domain.SMSGateway
	config     OTPServiceConfig
}

func NewOTPService(
	otpRepo domain.OTPRepository,
	smsGateway domain.SMSGateway,
	cfg OTPServiceConfig,
) *OTPService {
	return &OTPService{
		otpRepo:    otpRepo,
		smsGateway: smsGateway,
		config:     cfg,
	}
}

// RequestOTP dispatches an OTP challenge enforcing distributed rate limits (REG-01, AUT-01 to AUT-03).
func (s *OTPService) RequestOTP(
	ctx context.Context,
	phone domain.PhoneNumber,
	purpose domain.OTPPurpose,
) (sessionID string, timeoutSeconds int, err error) {
	// 1. Enforce distributed rate limiting
	allowed, err := s.otpRepo.CheckAndIncrementRateLimit(ctx, phone.String(), RateLimitWindow, MaxDispatchesPerWindow)
	if err != nil {
		return "", 0, fmt.Errorf("rate limit check failed: %w", err)
	}
	if !allowed {
		return "", 0, domain.ErrRateLimitExceeded
	}

	// 2. Generate secure 6-digit numeric code (CSPRNG)
	var code string
	if s.config.MockOTP {
		code = "123456"
	} else {
		n, err := rand.Int(rand.Reader, big.NewInt(900000))
		if err != nil {
			return "", 0, fmt.Errorf("failed to generate random OTP: %w", err)
		}
		code = fmt.Sprintf("%06d", n.Int64()+100000)
	}

	// 3. Create challenge session
	sessionID = uuid.NewString()
	now := time.Now()
	challenge := &domain.OTPChallenge{
		SessionID: sessionID,
		Phone:     phone.String(),
		Code:      code,
		Purpose:   purpose,
		Attempts:  0,
		CreatedAt: now,
		ExpiresAt: now.Add(DefaultOTPTTL),
	}

	// 4. Persist challenge in ephemeral repository
	if err := s.otpRepo.SaveChallenge(ctx, challenge, DefaultOTPTTL); err != nil {
		return "", 0, fmt.Errorf("failed to save OTP challenge: %w", err)
	}

	// 5. Dispatch via domestic telco gateway
	if err := s.smsGateway.SendOTP(ctx, phone, code); err != nil {
		return "", 0, fmt.Errorf("failed to dispatch SMS: %w", err)
	}

	return sessionID, int(DefaultOTPTTL.Seconds()), nil
}

// VerifyOTP validates the submitted code against the pending challenge (AUT-03).
func (s *OTPService) VerifyOTP(
	ctx context.Context,
	sessionID string,
	submittedCode string,
) (*domain.OTPChallenge, error) {
	challenge, err := s.otpRepo.GetChallenge(ctx, sessionID)
	if err != nil {
		return nil, domain.ErrOTPExpired
	}

	if time.Now().After(challenge.ExpiresAt) {
		_ = s.otpRepo.DeleteChallenge(ctx, sessionID)
		return nil, domain.ErrOTPExpired
	}

	// Assert code match
	if challenge.Code != submittedCode {
		attempts, err := s.otpRepo.IncrementAttempts(ctx, sessionID)
		if err != nil {
			return nil, err
		}

		// Reached maximum attempts: delete challenge immediately (AUT-03)
		if attempts >= MaxOTPVerificationAttempts {
			_ = s.otpRepo.DeleteChallenge(ctx, sessionID)
			return nil, domain.ErrOTPMaxAttemptsExceeded
		}

		return nil, domain.ErrOTPInvalid
	}

	// Verification successful: consume and invalidate challenge
	_ = s.otpRepo.DeleteChallenge(ctx, sessionID)
	return challenge, nil
}
