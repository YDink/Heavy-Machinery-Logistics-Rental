package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateTransition_PermittedPaths(t *testing.T) {
	tests := []struct {
		name    string
		from    BookingStatus
		to      BookingStatus
		wantErr bool
	}{
		{"PendingPayment -> PaymentEvidenceUploaded", StatusPendingPayment, StatusPaymentEvidenceUploaded, false},
		{"PendingPayment -> Cancelled", StatusPendingPayment, StatusCancelled, false},
		{"PaymentEvidenceUploaded -> Confirmed", StatusPaymentEvidenceUploaded, StatusConfirmed, false},
		{"PaymentEvidenceUploaded -> PaymentDeclined", StatusPaymentEvidenceUploaded, StatusPaymentDeclined, false},
		{"PaymentDeclined -> DisputeOpened", StatusPaymentDeclined, StatusDisputeOpened, false},
		{"PaymentDeclined -> PaymentEvidenceUploaded", StatusPaymentDeclined, StatusPaymentEvidenceUploaded, false},
		{"DisputeOpened -> DisputeResolvedRenter", StatusDisputeOpened, StatusDisputeResolvedRenter, false},
		{"DisputeOpened -> DisputeResolvedOwner", StatusDisputeOpened, StatusDisputeResolvedOwner, false},
		{"DisputeResolvedRenter -> Confirmed", StatusDisputeResolvedRenter, StatusConfirmed, false},
		{"DisputeResolvedOwner -> Cancelled", StatusDisputeResolvedOwner, StatusCancelled, false},
		{"Confirmed -> ReadyForPickup", StatusConfirmed, StatusReadyForPickup, false},
		{"ReadyForPickup -> Active", StatusReadyForPickup, StatusActive, false},
		{"Active -> ReturnInitiated", StatusActive, StatusReturnInitiated, false},
		{"ReturnInitiated -> Completed", StatusReturnInitiated, StatusCompleted, false},
		{"Completed -> Rated", StatusCompleted, StatusRated, false},
		// Invalid transitions
		{"Active -> Cancelled (Illegal)", StatusActive, StatusCancelled, true},
		{"Cancelled -> Confirmed (Terminal)", StatusCancelled, StatusConfirmed, true},
		{"Rated -> Active (Terminal)", StatusRated, StatusActive, true},
		{"PendingPayment -> Confirmed (Skipping Evidence)", StatusPendingPayment, StatusConfirmed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransition(tt.from, tt.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTransition(%s, %s) error = %v, wantErr = %v", tt.from, tt.to, err, tt.wantErr)
			}
		})
	}
}

func TestCalculateCancellation_Rules(t *testing.T) {
	now := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	startDateFuture := time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC) // > 48h
	startDateClose := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)  // < 48h

	prepayment := int64(10000000) // 100,000 ETB minor

	// 1. Pre-confirmation: Free cancellation
	bPre := &Booking{
		Status:                StatusPendingPayment,
		PrepaymentAmountMinor: prepayment,
	}
	resPre, err := CalculateCancellation(bPre, true, now, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resPre.PenaltyMinor != 0 || resPre.RefundMinor != 0 {
		t.Errorf("expected 0 penalty and 0 refund before confirmation, got penalty=%d refund=%d", resPre.PenaltyMinor, resPre.RefundMinor)
	}

	// 2. Post-confirmation: Contractor > 48h before start (20% penalty)
	bPostFuture := &Booking{
		Status:                StatusConfirmed,
		StartDate:             startDateFuture,
		PrepaymentAmountMinor: prepayment,
	}
	resPostFuture, err := CalculateCancellation(bPostFuture, true, now, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resPostFuture.PenaltyMinor != 2000000 || resPostFuture.RefundMinor != 8000000 {
		t.Errorf("expected 20%% penalty (20,000) and 80%% refund (80,000), got penalty=%d refund=%d", resPostFuture.PenaltyMinor, resPostFuture.RefundMinor)
	}

	// 3. Post-confirmation: Contractor <= 48h before start (50% penalty)
	bPostClose := &Booking{
		Status:                StatusConfirmed,
		StartDate:             startDateClose,
		PrepaymentAmountMinor: prepayment,
	}
	resPostClose, err := CalculateCancellation(bPostClose, true, now, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resPostClose.PenaltyMinor != 5000000 || resPostClose.RefundMinor != 5000000 {
		t.Errorf("expected 50%% penalty (50,000) and 50%% refund (50,000), got penalty=%d refund=%d", resPostClose.PenaltyMinor, resPostClose.RefundMinor)
	}

	// 4. Post-confirmation: Owner default (100% refund, 20% owner debt)
	resOwner, err := CalculateCancellation(bPostFuture, false, now, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resOwner.IsOwnerDebt || resOwner.PenaltyMinor != 2000000 || resOwner.RefundMinor != prepayment {
		t.Errorf("expected owner debt of 20%% (20,000) and 100%% refund (100,000), got isDebt=%v penalty=%d refund=%d",
			resOwner.IsOwnerDebt, resOwner.PenaltyMinor, resOwner.RefundMinor)
	}

	// 5. Force majeure
	resFM, err := CalculateCancellation(bPostClose, true, now, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resFM.PenaltyMinor != 0 || resFM.RefundMinor != prepayment {
		t.Errorf("expected force majeure 0 penalty and 100%% refund, got penalty=%d refund=%d", resFM.PenaltyMinor, resFM.RefundMinor)
	}

	// 6. Active booking cannot be cancelled
	bActive := &Booking{
		ID:     uuid.New(),
		Status: StatusActive,
	}
	if _, err := CalculateCancellation(bActive, true, now, false); err == nil {
		t.Errorf("expected error cancelling active booking, got nil")
	}
}
