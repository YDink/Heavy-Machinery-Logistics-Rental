package domain

import (
	"time"

	"github.com/google/uuid"
)

// DamageClaimDocket represents a formal post-rental structural damage claim filed within the 24h statutory window.
type DamageClaimDocket struct {
	ID                    uuid.UUID
	BookingID             uuid.UUID
	ClaimantID            uuid.UUID
	ComponentCategory     string
	IncidentNarrative     string
	RepairQuoteURI        string
	ClaimedAmountMinor    int64
	AcceptedAmountMinor   int64
	Currency              string
	DamagePhotos          []string
	Status                string
	DisputeReason         string
	ContractorRespondedAt *time.Time
	ArbitrationNotes      string
	ArbitratedBy          *uuid.UUID
	ArbitratedAt          *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
