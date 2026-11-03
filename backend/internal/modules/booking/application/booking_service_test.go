package application

import (
	"context"
	"testing"
	"time"

	"github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/booking/infrastructure"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/google/uuid"
)

type mockFleetProvider struct {
	fleets map[uuid.UUID]*machineryDomain.FleetAsset
}

func (m *mockFleetProvider) GetFleetAsset(ctx context.Context, id uuid.UUID) (*machineryDomain.FleetAsset, error) {
	if f, ok := m.fleets[id]; ok {
		return f, nil
	}
	return nil, machineryDomain.ErrFleetNotFound
}

func setupTestBookingService() (*BookingService, uuid.UUID, uuid.UUID, uuid.UUID) {
	repo := infrastructure.NewMemoryBookingRepository()
	leaseManager := infrastructure.NewMemoryHoldLeaseManager()

	ownerID := uuid.New()
	renterID := uuid.New()
	fleetID := uuid.New()

	fleetProvider := &mockFleetProvider{
		fleets: map[uuid.UUID]*machineryDomain.FleetAsset{
			fleetID: {
				ID:                  fleetID,
				OwnerID:             ownerID,
				DailyRateMinor:      3500000, // 35,000 ETB
				MinRentalDays:       2,
				BaselineMeterHours:  1200.0,
				PaymentTimeoutHours: 24,
				Status:              machineryDomain.FleetStatusActive,
			},
		},
	}

	service := NewBookingService(repo, leaseManager, fleetProvider, nil)
	return service, ownerID, renterID, fleetID
}

func TestBookingService_CreateBooking_SuccessAndCollisions(t *testing.T) {
	svc, ownerID, renterID, fleetID := setupTestBookingService()
	ctx := context.Background()

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(3 * 24 * time.Hour) // 3 days

	// 1. Successful creation
	input := CreateBookingInput{
		RenterID:            renterID,
		FleetID:             fleetID,
		StartDate:           start,
		EndDate:             end,
		SiteLocationAddress: "Bole Lemi Industrial Zone Phase 2, Addis Ababa",
	}

	b, err := svc.CreateBooking(ctx, input)
	if err != nil {
		t.Fatalf("expected successful booking creation, got error: %v", err)
	}

	if b.Status != domain.StatusPendingPayment {
		t.Errorf("expected status PENDING_PAYMENT, got %s", b.Status)
	}
	if b.DurationDays != 3 {
		t.Errorf("expected duration 3 days, got %d", b.DurationDays)
	}
	if b.GrossRentalAmountMinor != 3*3500000 {
		t.Errorf("expected gross rental 10,500,000, got %d", b.GrossRentalAmountMinor)
	}
	if b.PrepaymentAmountMinor != 10500000/2 {
		t.Errorf("expected prepayment 5,250,000, got %d", b.PrepaymentAmountMinor)
	}
	if b.HandoverPIN == nil || len(*b.HandoverPIN) != 6 {
		t.Errorf("expected 6-digit handover PIN, got %v", b.HandoverPIN)
	}

	// 2. Self-booking rejection
	selfInput := input
	selfInput.RenterID = ownerID
	if _, err := svc.CreateBooking(ctx, selfInput); err != domain.ErrSelfBookingProhibited {
		t.Errorf("expected ErrSelfBookingProhibited for owner self-booking, got %v", err)
	}

	// 3. Concurrent reservation collision
	competingRenter := uuid.New()
	competingInput := input
	competingInput.RenterID = competingRenter
	if _, err := svc.CreateBooking(ctx, competingInput); err != domain.ErrConcurrentReservationCollision {
		t.Errorf("expected ErrConcurrentReservationCollision when fleet on hold, got %v", err)
	}
}

func TestBookingService_FullLifecycleFSM(t *testing.T) {
	svc, ownerID, renterID, fleetID := setupTestBookingService()
	ctx := context.Background()

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(3 * 24 * time.Hour)

	// 1. Create Booking
	b, err := svc.CreateBooking(ctx, CreateBookingInput{
		RenterID:            renterID,
		FleetID:             fleetID,
		StartDate:           start,
		EndDate:             end,
		SiteLocationAddress: "Kilinto Industrial Park",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// 2. Submit Bank Slip
	slipInput := SubmitPaymentSlipInput{
		BookingID:            b.ID,
		UploaderID:           renterID,
		StorageURI:           "s3://payment-vault/slip1.webp",
		FileHashSHA256:       "hash123",
		TransactionReference: "TXN-888",
		AmountMinor:          b.PrepaymentAmountMinor,
		Currency:             "ETB",
	}
	slip, err := svc.SubmitPaymentSlip(ctx, slipInput)
	if err != nil {
		t.Fatalf("submit slip failed: %v", err)
	}
	if slip.Status != domain.SlipStatusPendingAudit {
		t.Errorf("expected slip status PENDING_AUDIT, got %s", slip.Status)
	}

	b, _ = svc.GetBooking(ctx, b.ID)
	if b.Status != domain.StatusPaymentEvidenceUploaded {
		t.Errorf("expected booking status PAYMENT_EVIDENCE_UPLOADED, got %s", b.Status)
	}

	// 3. Confirm Payment (Owner action)
	confirmInput := ConfirmPaymentInput{
		BookingID:       b.ID,
		ActorID:         ownerID,
		ExpectedVersion: b.Version,
		PaymentRef:      "TXN-888-VERIFIED",
		Notes:           "CBE audited",
	}
	b, err = svc.ConfirmPayment(ctx, confirmInput)
	if err != nil {
		t.Fatalf("confirm payment failed: %v", err)
	}
	if b.Status != domain.StatusConfirmed {
		t.Errorf("expected booking status CONFIRMED, got %s", b.Status)
	}
	if b.ConfirmedAt == nil {
		t.Errorf("expected ConfirmedAt timestamp to be set")
	}

	// 4. Handover Check-In (Renter enters challenge PIN on site)
	pin := *b.HandoverPIN
	checkInInput := HandoverCheckInInput{
		BookingID:           b.ID,
		ActorID:             renterID,
		ExpectedVersion:     b.Version,
		HandoverPIN:         pin,
		OperatingHoursMeter: 1205.5, // Baseline is 1200.0
		PhotoURIs:           []string{"s3://photos/meter.webp"},
	}
	b, err = svc.HandoverCheckIn(ctx, checkInInput)
	if err != nil {
		t.Fatalf("handover check-in failed: %v", err)
	}
	if b.Status != domain.StatusActive {
		t.Errorf("expected status ACTIVE, got %s", b.Status)
	}
	if b.CheckInMeterHours == nil || *b.CheckInMeterHours != 1205.5 {
		t.Errorf("expected check-in meter 1205.5, got %v", b.CheckInMeterHours)
	}

	// 5. Return Check-Out (Owner records return meter)
	// Non-monotonic meter rejected
	badReturnInput := ReturnCheckOutInput{
		BookingID:           b.ID,
		ActorID:             ownerID,
		ExpectedVersion:     b.Version,
		OperatingHoursMeter: 1200.0, // Less than check-in 1205.5
	}
	if _, err := svc.ReturnCheckOut(ctx, badReturnInput); err != domain.ErrInvalidMeterReading {
		t.Errorf("expected ErrInvalidMeterReading for regressing return meter, got %v", err)
	}

	// Valid clean return
	validReturnInput := ReturnCheckOutInput{
		BookingID:           b.ID,
		ActorID:             ownerID,
		ExpectedVersion:     b.Version,
		OperatingHoursMeter: 1228.0,
		DamageNotes:         "",
	}
	b, err = svc.ReturnCheckOut(ctx, validReturnInput)
	if err != nil {
		t.Fatalf("return check-out failed: %v", err)
	}
	if b.Status != domain.StatusCompleted {
		t.Errorf("expected status COMPLETED, got %s", b.Status)
	}
	if b.CompletedAt == nil {
		t.Errorf("expected CompletedAt to be set")
	}
}

func TestBookingService_DeclineAndDisputeFlow(t *testing.T) {
	svc, ownerID, renterID, fleetID := setupTestBookingService()
	ctx := context.Background()

	b, _ := svc.CreateBooking(ctx, CreateBookingInput{
		RenterID:  renterID,
		FleetID:   fleetID,
		StartDate: time.Now().Add(48 * time.Hour),
		EndDate:   time.Now().Add(96 * time.Hour),
	})

	_, _ = svc.SubmitPaymentSlip(ctx, SubmitPaymentSlipInput{
		BookingID:   b.ID,
		UploaderID:  renterID,
		AmountMinor: b.PrepaymentAmountMinor,
	})

	b, _ = svc.GetBooking(ctx, b.ID)

	// 1. Owner declines payment
	declineInput := DeclinePaymentInput{
		BookingID:       b.ID,
		ActorID:         ownerID,
		ExpectedVersion: b.Version,
		ReasonCode:      "FUNDS_NOT_RECEIVED_IN_ACCOUNT",
		RejectionNotes:  "No statement credit found for amount",
	}
	b, err := svc.DeclinePayment(ctx, declineInput)
	if err != nil {
		t.Fatalf("decline payment failed: %v", err)
	}
	if b.Status != domain.StatusPaymentDeclined {
		t.Errorf("expected status PAYMENT_DECLINED, got %s", b.Status)
	}

	// 2. Renter disputes decline
	disputeInput := DisputeBookingInput{
		BookingID:       b.ID,
		ActorID:         renterID,
		ExpectedVersion: b.Version,
		DisputeReason:   "Funds debited from Awash Bank account reference AW-9023",
	}
	b, err = svc.DisputeBooking(ctx, disputeInput)
	if err != nil {
		t.Fatalf("dispute booking failed: %v", err)
	}
	if b.Status != domain.StatusDisputeOpened {
		t.Errorf("expected status DISPUTE_OPENED, got %s", b.Status)
	}
}

func TestBookingService_CancellationAndOCCConflict(t *testing.T) {
	svc, _, renterID, fleetID := setupTestBookingService()
	ctx := context.Background()

	b, _ := svc.CreateBooking(ctx, CreateBookingInput{
		RenterID:  renterID,
		FleetID:   fleetID,
		StartDate: time.Now().Add(48 * time.Hour),
		EndDate:   time.Now().Add(96 * time.Hour),
	})

	// 1. OCC Version Conflict
	cancelInputWithBadVersion := CancelBookingInput{
		BookingID:       b.ID,
		ActorID:         renterID,
		ExpectedVersion: 999, // Stale version
		Reason:          "Cancelled site operation",
	}
	if _, err := svc.CancelBooking(ctx, cancelInputWithBadVersion); err != domain.ErrConcurrencyConflict {
		t.Errorf("expected ErrConcurrencyConflict for mismatched version, got %v", err)
	}

	// 2. Free pre-confirmation cancellation
	cancelInputValid := CancelBookingInput{
		BookingID:       b.ID,
		ActorID:         renterID,
		ExpectedVersion: b.Version,
		Reason:          "Cancelled site operation",
	}
	cancelledBooking, err := svc.CancelBooking(ctx, cancelInputValid)
	if err != nil {
		t.Fatalf("cancel booking failed: %v", err)
	}
	if cancelledBooking.Status != domain.StatusCancelled {
		t.Errorf("expected status CANCELLED, got %s", cancelledBooking.Status)
	}
	if cancelledBooking.PenaltyAmountMinor != 0 {
		t.Errorf("expected 0 penalty for pre-confirmation cancellation, got %d", cancelledBooking.PenaltyAmountMinor)
	}

	// 3. Confirm lease was released: another user can now book
	newRenter := uuid.New()
	bNew, err := svc.CreateBooking(ctx, CreateBookingInput{
		RenterID:  newRenter,
		FleetID:   fleetID,
		StartDate: time.Now().Add(48 * time.Hour),
		EndDate:   time.Now().Add(96 * time.Hour),
	})
	if err != nil {
		t.Fatalf("expected new booking to succeed after release, got error: %v", err)
	}
	if bNew.Status != domain.StatusPendingPayment {
		t.Errorf("expected new booking status PENDING_PAYMENT, got %s", bNew.Status)
	}
}
