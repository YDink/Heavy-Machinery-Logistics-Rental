package domain

import (
	"time"

	"github.com/google/uuid"
)

// RentalSettlement represents the deterministic financial statement reconciling contract liabilities.
type RentalSettlement struct {
	ID                         uuid.UUID
	BookingID                  uuid.UUID
	BaseMachineFeeMinor        int64
	BundledAttachmentsFeeMinor int64
	ContractualAllowedHours    float64
	ActualOperatingHours       float64
	CreditedBreakdownHours     float64
	BillableExcessHours        float64
	ExcessHourlyRateMinor      int64
	ExcessHoursSurchargeMinor  int64
	AcceptedDamageMinor        int64
	GrossLiabilitiesMinor      int64
	PrepaymentSettledMinor     int64
	NetBalanceMinor            int64
	Currency                   string
	Status                     string
	FinalSlipStorageURI        string
	FinalSlipFileHash          string
	SettledAt                  *time.Time
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}
