package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/infrastructure"
)

func TestTokenRotation_NormalRotation(t *testing.T) {
	ctx := context.Background()
	cache := infrastructure.NewMemoryCache()

	tokenService, err := application.NewTokenService(nil, cache)
	if err != nil {
		t.Fatalf("failed to create TokenService: %v", err)
	}

	phone, _ := domain.ParsePhoneNumber("+251911442910", "ET")
	user := &domain.User{
		ID:            uuid.New(),
		PhoneNumber:   phone,
		CountryCode:   "ET",
		FullName:      "Dawit Haile-Mariam",
		PrimaryRole:   domain.RoleRenter,
		AccountStatus: domain.AccountStatusActive,
		KYCStatus:     domain.KYCStatusVerified,
		TokenVersion:  1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 1. Initial login mints Token Pair 1
	tokenPair1, err := tokenService.MintTokenPair(ctx, user, "")
	if err != nil {
		t.Fatalf("failed to mint initial token pair: %v", err)
	}

	// Verify access token claims
	claims, err := tokenService.VerifyAccessToken(ctx, tokenPair1.AccessToken)
	if err != nil {
		t.Fatalf("access token verification failed: %v", err)
	}
	if claims.Subject != user.ID.String() {
		t.Errorf("expected subject %s, got %s", user.ID.String(), claims.Subject)
	}
	if claims.Role != string(domain.RoleRenter) {
		t.Errorf("expected role %s, got %s", domain.RoleRenter, claims.Role)
	}

	// 2. Client legitimately rotates refresh token
	tokenPair2, err := tokenService.RotateRefreshToken(ctx, tokenPair1.RefreshToken, user)
	if err != nil {
		t.Fatalf("failed to rotate refresh token: %v", err)
	}

	if tokenPair2.RefreshToken == tokenPair1.RefreshToken {
		t.Fatalf("expected new refresh token string after rotation")
	}

	// Verify new access token is valid
	claims2, err := tokenService.VerifyAccessToken(ctx, tokenPair2.AccessToken)
	if err != nil {
		t.Fatalf("rotated access token verification failed: %v", err)
	}
	if claims2.Subject != user.ID.String() {
		t.Errorf("expected subject %s, got %s", user.ID.String(), claims2.Subject)
	}
}

// TestTokenRotation_ReplayAttack asserts D-05-01-002:
// Presenting a previously used refresh token triggers immediate family revocation.
func TestTokenRotation_ReplayAttack(t *testing.T) {
	ctx := context.Background()
	cache := infrastructure.NewMemoryCache()

	tokenService, err := application.NewTokenService(nil, cache)
	if err != nil {
		t.Fatalf("failed to create TokenService: %v", err)
	}

	phone, _ := domain.ParsePhoneNumber("+251911442910", "ET")
	user := &domain.User{
		ID:            uuid.New(),
		PhoneNumber:   phone,
		CountryCode:   "ET",
		FullName:      "Dawit Haile-Mariam",
		PrimaryRole:   domain.RoleRenter,
		AccountStatus: domain.AccountStatusActive,
		KYCStatus:     domain.KYCStatusVerified,
		TokenVersion:  1,
	}

	// 1. Initial mint (Version 1)
	pair1, err := tokenService.MintTokenPair(ctx, user, "")
	if err != nil {
		t.Fatalf("failed to mint pair1: %v", err)
	}

	// 2. Legitimate rotation: pair1.RefreshToken -> pair2.RefreshToken (Version 2)
	pair2, err := tokenService.RotateRefreshToken(ctx, pair1.RefreshToken, user)
	if err != nil {
		t.Fatalf("failed to rotate to pair2: %v", err)
	}

	// 3. REPLAY ATTACK SIMULATION:
	// Attacker steals pair1.RefreshToken and attempts to reuse it
	_, err = tokenService.RotateRefreshToken(ctx, pair1.RefreshToken, user)
	if err != domain.ErrTokenFamilyCompromised {
		t.Fatalf("expected ErrTokenFamilyCompromised on replay, got: %v", err)
	}

	// 4. Invariant Assertion: The entire family is revoked.
	// Even the legitimate user's current token (pair2.RefreshToken) must now be rejected!
	_, err = tokenService.RotateRefreshToken(ctx, pair2.RefreshToken, user)
	if err != domain.ErrTokenFamilyCompromised {
		t.Fatalf("expected ErrTokenFamilyCompromised for legitimate token after family revocation, got: %v", err)
	}
}
