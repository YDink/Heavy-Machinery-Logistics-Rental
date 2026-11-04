package domain

import (
	"context"
	"time"
)

// OTPRepository defines ephemeral storage and rate limiting for OTP challenges in Valkey.
type OTPRepository interface {
	SaveChallenge(ctx context.Context, challenge *OTPChallenge, ttl time.Duration) error
	GetChallenge(ctx context.Context, sessionID string) (*OTPChallenge, error)
	IncrementAttempts(ctx context.Context, sessionID string) (int, error)
	DeleteChallenge(ctx context.Context, sessionID string) error
	CheckAndIncrementRateLimit(ctx context.Context, phone string, window time.Duration, maxRequests int) (bool, error)
}
