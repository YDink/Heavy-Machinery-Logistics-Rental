package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RentalRepository interface {
	// Handover & Inspection
	CreateHandoverInspection(ctx context.Context, inspection *HandoverInspection) error
	GetHandoverInspection(ctx context.Context, bookingID uuid.UUID) (*HandoverInspection, error)
	CreateHandoverChallenge(ctx context.Context, challenge *HandoverChallenge) error
	GetActiveHandoverChallenge(ctx context.Context, bookingID uuid.UUID) (*HandoverChallenge, error)
	IncrementHandoverChallengeAttempts(ctx context.Context, challengeID uuid.UUID) (*HandoverChallenge, error)
	VerifyHandoverChallenge(ctx context.Context, challengeID uuid.UUID, touchSignatureBase64 string, lat, lng *float64) error

	// Shift Telemetry
	CreateRentalShift(ctx context.Context, shift *RentalShift) error
	GetActiveRentalShift(ctx context.Context, bookingID uuid.UUID) (*RentalShift, error)
	UpdateRentalShiftProgress(ctx context.Context, shiftID uuid.UUID, accumulatedHours, excessHours float64, endMeter *float64, isPaused bool) error
	EndRentalShift(ctx context.Context, shiftID uuid.UUID, endMeter float64, accumulatedHours, excessHours float64) error
	ListRentalShifts(ctx context.Context, bookingID uuid.UUID) ([]RentalShift, error)

	// Extension Management
	CreateRentalExtension(ctx context.Context, ext *RentalExtension) error
	GetRentalExtension(ctx context.Context, extensionID uuid.UUID) (*RentalExtension, error)
	ListRentalExtensions(ctx context.Context, bookingID uuid.UUID) ([]RentalExtension, error)
	UpdateRentalExtensionStatus(ctx context.Context, extensionID uuid.UUID, status string, notes *string, slipURI, slipHash *string) error
	ExpandBookingHirePeriod(ctx context.Context, bookingID uuid.UUID, newEndDate time.Time, additionalDays int, additionalGrossMinor, additionalPrepayMinor int64, expectedVersion int) error

	// Mechanical Breakdown
	CreateMechanicalBreakdown(ctx context.Context, breakdown *MechanicalBreakdown) error
	GetMechanicalBreakdown(ctx context.Context, breakdownID uuid.UUID) (*MechanicalBreakdown, error)
	GetActiveMechanicalBreakdown(ctx context.Context, bookingID uuid.UUID) (*MechanicalBreakdown, error)
	ListMechanicalBreakdowns(ctx context.Context, bookingID uuid.UUID) ([]MechanicalBreakdown, error)
	ResolveMechanicalBreakdown(ctx context.Context, breakdownID uuid.UUID, mechanicName, mechanicPhone, resolutionNotes string) error

	// Machine & Booking Baseline Queries
	GetFleetBaselineMeter(ctx context.Context, fleetID uuid.UUID) (float64, error)
	CheckForwardBookingConflict(ctx context.Context, fleetID uuid.UUID, excludeBookingID uuid.UUID, currentEndDate, newEndDate time.Time) (bool, error)

	// Return Inspection & Demobilization (Partition 05-07)
	CreateReturnInspection(ctx context.Context, inspection *ReturnInspection) error
	GetReturnInspection(ctx context.Context, bookingID uuid.UUID) (*ReturnInspection, error)

	// Damage Claims & Arbitration (Partition 05-07)
	CreateDamageClaimDocket(ctx context.Context, claim *DamageClaimDocket) error
	GetDamageClaimDocket(ctx context.Context, claimID uuid.UUID) (*DamageClaimDocket, error)
	GetDamageClaimDocketByBooking(ctx context.Context, bookingID uuid.UUID) (*DamageClaimDocket, error)
	ListDamageClaimDockets(ctx context.Context, bookingID uuid.UUID) ([]DamageClaimDocket, error)
	RespondToDamageClaim(ctx context.Context, claimID uuid.UUID, status string, acceptedMinor int64, disputeReason string) error
	ArbitrateDamageClaim(ctx context.Context, claimID uuid.UUID, status string, acceptedMinor int64, arbitrationNotes string, arbitratedBy uuid.UUID) error
	WaiveDamageClaim(ctx context.Context, claimID uuid.UUID) error

	// Rental Settlement Reconciliation (Partition 05-07)
	UpsertRentalSettlement(ctx context.Context, settlement *RentalSettlement) error
	GetRentalSettlement(ctx context.Context, bookingID uuid.UUID) (*RentalSettlement, error)
	UpdateRentalSettlementPaymentEvidence(ctx context.Context, bookingID uuid.UUID, slipURI, slipHash string) error
	ConfirmRentalSettlement(ctx context.Context, bookingID uuid.UUID) error

	// Double-Blind Reputation Reviews (Partition 05-07)
	CreateBookingReview(ctx context.Context, review *BookingReview) error
	GetBookingReviews(ctx context.Context, bookingID uuid.UUID) ([]BookingReview, error)
	PublishBookingReviews(ctx context.Context, bookingID uuid.UUID) error
	ListPublishedReviewsForTarget(ctx context.Context, targetID uuid.UUID) ([]BookingReview, error)

	// Booking Lifecycle Completion Updates
	UpdateBookingCompletionSummary(ctx context.Context, bookingID uuid.UUID, status string, checkOutMeter float64, damageClaimSubmitted, damageClaimWaived bool, acceptedDamageMinor, excessSurchargeMinor, finalNetBalanceMinor int64, expectedVersion int) error
}
