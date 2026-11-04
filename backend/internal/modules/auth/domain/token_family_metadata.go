package domain

import (
	"time"

	"github.com/google/uuid"
)

// TokenFamilyMetadata tracks refresh token rotation status in Valkey (ADR-014).
type TokenFamilyMetadata struct {
	FamilyID    string
	UserID      uuid.UUID
	ActiveToken string
	IsRevoked   bool
	CreatedAt   time.Time
}
