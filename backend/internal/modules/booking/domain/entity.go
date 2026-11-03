package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDateRange                 = errors.New("start date must be in the future and before end date")
	ErrInvalidDuration                  = errors.New("rental duration must meet fleet minimum days")
	ErrConcurrentReservationCollision   = errors.New("fleet or attachment is currently reserved or on hold for overlapping dates")
	ErrInvalidStateTransition           = errors.New("invalid booking state transition")
	ErrBookingNotFound                  = errors.New("booking not found")
	ErrUnauthorizedAction               = errors.New("user is not authorized to perform action on this booking")
	ErrConcurrencyConflict              = errors.New("booking was updated by another transaction (version conflict)")
	ErrInvalidMeterReading              = errors.New("engine meter hours cannot regress (must be >= baseline check-in)")
	ErrInvalidChallengePIN              = errors.New("invalid 6-digit handover challenge PIN")
	ErrSlipNotFound                     = errors.New("bank transfer payment slip not found")
	ErrInvalidSlipStatus                = errors.New("bank slip status is not eligible for audit")
	ErrEmptyReason                      = errors.New("reason code and explanation cannot be empty")
	ErrSelfBookingProhibited            = errors.New("fleet owner cannot book their own machinery")
)

type BookingStatus string

const (
	StatusDraft                     BookingStatus = "DRAFT"
	StatusPendingPayment            BookingStatus = "PENDING_PAYMENT"
	StatusPaymentEvidenceUploaded   BookingStatus = "PAYMENT_EVIDENCE_UPLOADED"
	StatusPaymentDeclined           BookingStatus = "PAYMENT_DECLINED"
	StatusDisputeOpened             BookingStatus = "DISPUTE_OPENED"
	StatusDisputeResolvedRenter     BookingStatus = "DISPUTE_RESOLVED_RENTER"
	StatusDisputeResolvedOwner      BookingStatus = "DISPUTE_RESOLVED_OWNER"
	StatusConfirmed                 BookingStatus = "CONFIRMED"
	StatusAwaitingFleetReadiness    BookingStatus = "AWAITING_FLEET_READINESS"
	StatusPendingAdditionalPayment  BookingStatus = "PENDING_ADDITIONAL_PAYMENT"
	StatusReadyForPickup            BookingStatus = "READY_FOR_PICKUP"
	StatusActive                    BookingStatus = "ACTIVE"
	StatusActiveBreakdownPaused     BookingStatus = "ACTIVE_BREAKDOWN_PAUSED"
	StatusActiveExtensionHold       BookingStatus = "ACTIVE_EXTENSION_HOLD"
	StatusReturnInitiated           BookingStatus = "RETURN_INITIATED"
	StatusPendingFinalPayment       BookingStatus = "PENDING_FINAL_PAYMENT"
	StatusCompleted                 BookingStatus = "COMPLETED"
	StatusRated                     BookingStatus = "RATED"
	StatusCancelled                 BookingStatus = "CANCELLED"
)

type BankSlipStatus string

const (
	SlipStatusPendingAudit BankSlipStatus = "PENDING_AUDIT"
	SlipStatusVerified     BankSlipStatus = "VERIFIED"
	SlipStatusRejected     BankSlipStatus = "REJECTED"
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
	Attachments            []BookingAttachment
	Version                int
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type BookingAttachment struct {
	ID                    uuid.UUID
	BookingID             uuid.UUID
	AttachmentID          uuid.UUID
	AttachmentType        string
	Make                  string
	Model                 string
	SerialNumber          string
	BundledDailyRateMinor int64
	CreatedAt             time.Time
}

type BankSlip struct {
	ID                   uuid.UUID
	BookingID            uuid.UUID
	UploaderID           uuid.UUID
	StorageURI           string
	FileHashSHA256       string
	TransactionReference *string
	AmountMinor          int64
	Currency             string
	Status               BankSlipStatus
	RejectionReasonCode  *string
	RejectionNotes       *string
	CreatedAt            time.Time
	AuditedAt            *time.Time
}

type AuditLog struct {
	ID         uuid.UUID
	BookingID  uuid.UUID
	ActorID    *uuid.UUID
	Action     string
	FromStatus BookingStatus
	ToStatus   BookingStatus
	Reason     string
	Payload    map[string]any
	CreatedAt  time.Time
}
