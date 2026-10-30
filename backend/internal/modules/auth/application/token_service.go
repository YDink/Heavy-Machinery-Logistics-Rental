package application

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
)

const (
	AccessTokenDuration  = 15 * time.Minute      // 900 seconds (AUT-08)
	RefreshTokenDuration = 30 * 24 * time.Hour   // 30 days (ADR-014)
)

type TokenClaims struct {
	jwt.RegisteredClaims
	Role         string `json:"role"`
	CountryCode  string `json:"country_code"`
	TokenVersion int32  `json:"token_version"`
}

type TokenService struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	tokenStore domain.TokenStore
}

func NewTokenService(privateKey *ecdsa.PrivateKey, tokenStore domain.TokenStore) (*TokenService, error) {
	if privateKey == nil {
		// Generate an ephemeral P-256 key if none provided (for dev/test environments)
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ECDSA key: %w", err)
		}
		privateKey = key
	}

	return &TokenService{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		tokenStore: tokenStore,
	}, nil
}

// MintTokenPair issues an ES256 Access Token and a rotated Refresh Token (ADR-014).
func (s *TokenService) MintTokenPair(
	ctx context.Context,
	user *domain.User,
	familyID string,
) (*domain.TokenPair, error) {
	if familyID == "" {
		familyID = uuid.NewString()
	}

	now := time.Now()
	tokenID := uuid.NewString()

	// 1. Create Access Token (ES256)
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			Issuer:    "https://api.dinkmachinery.com",
			Audience:  jwt.ClaimStrings{"dink-core-api"},
			ID:        tokenID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenDuration)),
		},
		Role:         string(user.PrimaryRole),
		CountryCode:  user.CountryCode,
		TokenVersion: user.TokenVersion,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	accessTokenString, err := jwtToken.SignedString(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// 2. Generate Refresh Token (256-bit cryptographically secure string)
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshTokenString := hex.EncodeToString(randomBytes)

	// 3. Persist Refresh Token metadata in TokenStore
	if err := s.tokenStore.SaveRefreshToken(ctx, refreshTokenString, familyID, user.ID, RefreshTokenDuration); err != nil {
		return nil, fmt.Errorf("failed to persist refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:      accessTokenString,
		RefreshToken:     refreshTokenString,
		TokenType:        "Bearer",
		ExpiresInSeconds: int(AccessTokenDuration.Seconds()),
	}, nil
}

// RotateRefreshToken validates the refresh token, detects replay attacks, and issues a fresh token pair (ADR-014).
func (s *TokenService) RotateRefreshToken(
	ctx context.Context,
	refreshToken string,
	user *domain.User,
) (*domain.TokenPair, error) {
	metadata, err := s.tokenStore.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// REPLAY ATTACK DETECTION (OWASP ASVS §3.6 / D-05-01-002)
	// If the token is already revoked or not active, an attacker or compromised client is reusing it
	if metadata.IsRevoked || metadata.ActiveToken != refreshToken {
		// Immediately revoke the entire token family
		_ = s.tokenStore.RevokeFamily(ctx, metadata.FamilyID)
		return nil, domain.ErrTokenFamilyCompromised
	}

	// Issue fresh token pair under the same family
	newTokenPair, err := s.MintTokenPair(ctx, user, metadata.FamilyID)
	if err != nil {
		return nil, err
	}

	// Rotate the refresh token in the store (marking the old one as used)
	if err := s.tokenStore.RotateRefreshToken(ctx, refreshToken, newTokenPair.RefreshToken, metadata.FamilyID, RefreshTokenDuration); err != nil {
		return nil, fmt.Errorf("failed to rotate refresh token: %w", err)
	}

	return newTokenPair, nil
}

// VerifyAccessToken parses and validates an ES256 access token.
func (s *TokenService) VerifyAccessToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Strict algorithm enforcement (reject none and HS256)
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("%w: unexpected signing method %v", domain.ErrInvalidToken, t.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	// Check if token ID is blacklisted
	blacklisted, err := s.tokenStore.IsBlacklisted(ctx, claims.ID)
	if err == nil && blacklisted {
		return nil, domain.ErrTokenRevoked
	}

	return claims, nil
}
