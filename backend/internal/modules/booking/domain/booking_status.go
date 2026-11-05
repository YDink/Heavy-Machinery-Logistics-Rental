package domain

type BookingStatus string

const (
	StatusDraft                    BookingStatus = "DRAFT"
	StatusPendingPayment           BookingStatus = "PENDING_PAYMENT"
	StatusPaymentEvidenceUploaded  BookingStatus = "PAYMENT_EVIDENCE_UPLOADED"
	StatusPaymentDeclined          BookingStatus = "PAYMENT_DECLINED"
	StatusDisputeOpened            BookingStatus = "DISPUTE_OPENED"
	StatusDisputeResolvedRenter    BookingStatus = "DISPUTE_RESOLVED_RENTER"
	StatusDisputeResolvedOwner     BookingStatus = "DISPUTE_RESOLVED_OWNER"
	StatusConfirmed                BookingStatus = "CONFIRMED"
	StatusAwaitingFleetReadiness   BookingStatus = "AWAITING_FLEET_READINESS"
	StatusPendingAdditionalPayment BookingStatus = "PENDING_ADDITIONAL_PAYMENT"
	StatusReadyForPickup           BookingStatus = "READY_FOR_PICKUP"
	StatusActive                   BookingStatus = "ACTIVE"
	StatusActiveBreakdownPaused    BookingStatus = "ACTIVE_BREAKDOWN_PAUSED"
	StatusActiveExtensionHold      BookingStatus = "ACTIVE_EXTENSION_HOLD"
	StatusReturnInitiated          BookingStatus = "RETURN_INITIATED"
	StatusDamagePendingReview      BookingStatus = "DAMAGE_PENDING_REVIEW"
	StatusDamageDisputed           BookingStatus = "DAMAGE_DISPUTED"
	StatusPendingFinalPayment      BookingStatus = "PENDING_FINAL_PAYMENT"
	StatusCompleted                BookingStatus = "COMPLETED"
	StatusRated                    BookingStatus = "RATED"
	StatusCancelled                BookingStatus = "CANCELLED"
)
