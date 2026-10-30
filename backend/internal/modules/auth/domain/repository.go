package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserRepository specifies data persistence invariants for User entities.
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByPhoneNumber(ctx context.Context, phone string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	MarkPhoneVerified(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateAccountStatus(ctx context.Context, id uuid.UUID, status AccountStatus) (*User, error)
	UpdateKYCStatus(ctx context.Context, id uuid.UUID, status KYCStatus) (*User, error)
	IncrementTokenVersion(ctx context.Context, id uuid.UUID) (int32, error)
}

// OTPRepository defines ephemeral storage and rate limiting for OTP challenges in Valkey.
type OTPRepository interface {
	SaveChallenge(ctx context.Context, challenge *OTPChallenge, ttl time.Duration) error
	GetChallenge(ctx context.Context, sessionID string) (*OTPChallenge, error)
	IncrementAttempts(ctx context.Context, sessionID string) (int, error)
	DeleteChallenge(ctx context.Context, sessionID string) error
	CheckAndIncrementRateLimit(ctx context.Context, phone string, window time.Duration, maxRequests int) (bool, error)
}

// TokenFamilyMetadata tracks refresh token rotation status in Valkey (ADR-014).
type TokenFamilyMetadata struct {
	FamilyID    string
	UserID      uuid.UUID
	ActiveToken string
	IsRevoked   bool
	CreatedAt   time.Time
}

// TokenStore specifies operations for refresh token rotation and revocation blacklists.
type TokenStore interface {
	SaveRefreshToken(ctx context.Context, token string, familyID string, userID uuid.UUID, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, token string) (*TokenFamilyMetadata, error)
	RotateRefreshToken(ctx context.Context, oldToken string, newToken string, familyID string, ttl time.Duration) error
	RevokeFamily(ctx context.Context, familyID string) error
	IsBlacklisted(ctx context.Context, tokenID string) (bool, error)
	BlacklistToken(ctx context.Context, tokenID string, ttl time.Duration) error
}

// SMSGateway defines the messaging abstraction for domestic carriers.
type SMSGateway interface {
	SendOTP(ctx context.Context, phone PhoneNumber, code string) error
}
