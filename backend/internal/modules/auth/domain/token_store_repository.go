package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TokenStore specifies operations for refresh token rotation and revocation blacklists.
type TokenStore interface {
	SaveRefreshToken(ctx context.Context, token string, familyID string, userID uuid.UUID, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, token string) (*TokenFamilyMetadata, error)
	RotateRefreshToken(ctx context.Context, oldToken string, newToken string, familyID string, ttl time.Duration) error
	RevokeFamily(ctx context.Context, familyID string) error
	IsBlacklisted(ctx context.Context, tokenID string) (bool, error)
	BlacklistToken(ctx context.Context, tokenID string, ttl time.Duration) error
}
