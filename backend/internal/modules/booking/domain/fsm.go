package domain

import (
	"fmt"
	"time"
)

var permittedTransitions = map[BookingStatus][]BookingStatus{
	StatusDraft: {
		StatusPendingPayment,
		StatusCancelled,
	},
	StatusPendingPayment: {
		StatusPaymentEvidenceUploaded,
		StatusCancelled,
	},
	StatusPaymentEvidenceUploaded: {
		StatusPaymentDeclined,
		StatusConfirmed,
		StatusCancelled,
	},
	StatusPaymentDeclined: {
		StatusPaymentEvidenceUploaded,
		StatusDisputeOpened,
		StatusCancelled,
	},
	StatusDisputeOpened: {
		StatusDisputeResolvedRenter,
		StatusDisputeResolvedOwner,
	},
	StatusDisputeResolvedRenter: {
		StatusConfirmed,
	},
	StatusDisputeResolvedOwner: {
		StatusCancelled,
	},
	StatusConfirmed: {
		StatusAwaitingFleetReadiness,
		StatusReadyForPickup,
		StatusActive, // Direct site activation if already mobilized
		StatusCancelled,
	},
	StatusAwaitingFleetReadiness: {
		StatusPendingAdditionalPayment,
		StatusReadyForPickup,
		StatusCancelled,
	},
	StatusPendingAdditionalPayment: {
		StatusReadyForPickup,
		StatusCancelled,
	},
	StatusReadyForPickup: {
		StatusActive,
		StatusCancelled,
	},
	StatusActive: {
		StatusActiveBreakdownPaused,
		StatusActiveExtensionHold,
		StatusReturnInitiated,
	},
	StatusActiveBreakdownPaused: {
		StatusActive,
		StatusReturnInitiated,
		StatusCancelled,
	},
	StatusActiveExtensionHold: {
		StatusActive,
		StatusReturnInitiated,
		StatusCancelled,
	},
	StatusReturnInitiated: {
		StatusDamagePendingReview,
		StatusPendingFinalPayment,
		StatusCompleted,
	},
	StatusDamagePendingReview: {
		StatusDamageDisputed,
		StatusPendingFinalPayment,
		StatusCompleted,
	},
	StatusDamageDisputed: {
		StatusPendingFinalPayment,
		StatusCompleted,
	},
	StatusPendingFinalPayment: {
		StatusCompleted,
	},
	StatusCompleted: {
		StatusRated,
	},
	StatusRated:     {},
	StatusCancelled: {},
}

// ValidateTransition asserts whether a state transition from `from` to `to` is legally permitted.
func ValidateTransition(from, to BookingStatus) error {
	allowed, ok := permittedTransitions[from]
	if !ok {
		return fmt.Errorf("%w: unknown source state '%s'", ErrInvalidStateTransition, from)
	}

	for _, target := range allowed {
		if target == to {
			return nil
		}
	}

	return fmt.Errorf("%w: cannot transition from '%s' to '%s'", ErrInvalidStateTransition, from, to)
}

// CancellationResult holds the financial penalty and refund calculation for a booking cancellation.
type CancellationResult struct {
	PenaltyMinor int64
	RefundMinor  int64
	IsOwnerDebt  bool
	Reason       string
}

// CalculateCancellation evaluates contractual cancellation penalty rules based on actor, timing, and current status.
func CalculateCancellation(b *Booking, isRenter bool, now time.Time, isForceMajeure bool) (CancellationResult, error) {
	if b.Status == StatusCancelled || b.Status == StatusRated || b.Status == StatusCompleted {
		return CancellationResult{}, fmt.Errorf("%w: cannot cancel booking in terminal state '%s'", ErrInvalidStateTransition, b.Status)
	}

	if b.Status == StatusActive || b.Status == StatusReturnInitiated || b.Status == StatusPendingFinalPayment {
		return CancellationResult{}, fmt.Errorf("%w: active or returning rentals cannot be cancelled; perform checkout return reconciliation", ErrInvalidStateTransition)
	}

	// 1. Pre-confirmation: Free cancellation for either party
	if b.Status == StatusDraft || b.Status == StatusPendingPayment || b.Status == StatusPaymentEvidenceUploaded || b.Status == StatusPaymentDeclined {
		return CancellationResult{
			PenaltyMinor: 0,
			RefundMinor:  0,
			IsOwnerDebt:  false,
			Reason:       "Pre-confirmation free cancellation",
		}, nil
	}

	// 2. Force Majeure (Logistics failure, transshipment crash)
	if isForceMajeure {
		return CancellationResult{
			PenaltyMinor: 0,
			RefundMinor:  b.PrepaymentAmountMinor,
			IsOwnerDebt:  false,
			Reason:       "Force majeure mutual exoneration (100% refund)",
		}, nil
	}

	// 3. Post-confirmation cancellation rules
	if isRenter {
		// Contractor cancellation post-confirmation
		// Check timing relative to StartDate
		hoursUntilStart := b.StartDate.Sub(now).Hours()
		if hoursUntilStart > 48 {
			// > 48h before start: 20% forfeiture of prepayment, 80% wire refund
			penalty := int64(float64(b.PrepaymentAmountMinor) * 0.20)
			refund := b.PrepaymentAmountMinor - penalty
			return CancellationResult{
				PenaltyMinor: penalty,
				RefundMinor:  refund,
				IsOwnerDebt:  false,
				Reason:       "Contractor cancelled > 48h before rental start (20% penalty)",
			}, nil
		}

		// <= 48h before start: 50% forfeiture of prepayment, 50% wire refund
		penalty := int64(float64(b.PrepaymentAmountMinor) * 0.50)
		refund := b.PrepaymentAmountMinor - penalty
		return CancellationResult{
			PenaltyMinor: penalty,
			RefundMinor:  refund,
			IsOwnerDebt:  false,
			Reason:       "Contractor cancelled <= 48h before rental start (50% penalty)",
		}, nil
	}

	// Equipment Owner unilateral cancellation post-confirmation
	// Contractor receives 100% refund; owner is debited 20% breach penalty
	ownerDebt := int64(float64(b.PrepaymentAmountMinor) * 0.20)
	return CancellationResult{
		PenaltyMinor: ownerDebt,
		RefundMinor:  b.PrepaymentAmountMinor,
		IsOwnerDebt:  true,
		Reason:       "Fleet owner unilateral cancellation (20% default fee levied on owner wallet)",
	}, nil
}
