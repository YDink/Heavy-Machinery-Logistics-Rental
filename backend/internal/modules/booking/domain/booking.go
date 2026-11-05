package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID                     uuid.UUID
	BookingReference       string
	RenterID               uuid.UUID
	OwnerID                uuid.UUID
	FleetID                uuid.UUID
	StartDate              time.Time
	EndDate                time.Time
	DurationDays           int
	DailyRateAppliedMinor  int64
	BaseHireFeeMinor       int64
	AttachmentsFeeMinor    int64
	GrossRentalAmountMinor int64
	PrepaymentAmountMinor  int64
	Currency               string
	SiteLocationAddress    string
	Status                 BookingStatus
	RejectionReasonCode    *string
	RejectionNotes         *string
	CancellationReason     *string
	CancelledBy            *uuid.UUID
	PenaltyAmountMinor     int64
	RefundAmountMinor      int64
	ExpiresAt              *time.Time
	ConfirmedAt            *time.Time
	ActivatedAt            *time.Time
	CompletedAt            *time.Time
	CancelledAt            *time.Time
	HandoverPIN            *string
	AssignedOperatorID     *uuid.UUID
	CheckInMeterHours      *float64
	CheckOutMeterHours     *float64
	Attachments               []BookingAttachment
	DamageClaimSubmitted      bool
	DamageClaimWaived         bool
	AcceptedDamageMinor       int64
	ExcessHoursSurchargeMinor int64
	FinalNetBalanceMinor      int64
	Version                   int
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
