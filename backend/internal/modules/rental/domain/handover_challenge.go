package domain

import (
	"time"

	"github.com/google/uuid"
)

type HandoverChallenge struct {
	ID                   uuid.UUID
	BookingID            uuid.UUID
	ChallengePIN         string
	ExpiresAt            time.Time
	AttemptCount         int
	MaxAttempts          int
	IsVerified           bool
	VerifiedAt           *time.Time
	TouchSignatureBase64 string
	ContractorLat        *float64
	ContractorLng        *float64
	CreatedAt            time.Time
}
