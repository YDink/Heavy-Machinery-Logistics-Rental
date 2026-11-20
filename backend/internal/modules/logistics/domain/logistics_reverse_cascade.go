package domain

import (
	"time"

	"github.com/google/uuid"
)

// LogisticsReverseCascade records a Force Majeure cancellation and mutual party exoneration ledger entry.
type LogisticsReverseCascade struct {
	ID                             uuid.UUID
	ParentBookingID                uuid.UUID
	SubBookingID                   uuid.UUID
	TriggeredBy                    uuid.UUID
	Resolution                     string
	Reason                         string
	RenterRefundDueMinor           int64
	OwnerPenaltyLeviedMinor        int64
	OwnerCommissionRestitutionMinor int64
	CarrierDefaultPenaltyMinor     int64
	CarrierRatingDeduction         float64
	Currency                       string
	CreatedAt                      time.Time
}
