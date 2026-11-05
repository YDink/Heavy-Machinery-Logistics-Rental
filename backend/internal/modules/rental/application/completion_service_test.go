package application_test

import (
	"context"
	"testing"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	rentalDomain "github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	rentalInfra "github.com/dinkmachinery/dink-core/internal/modules/rental/infrastructure"
	"github.com/google/uuid"
)

type mockBookingProvider struct {
	bookings map[uuid.UUID]*bookingDomain.Booking
}

func newMockBookingProvider() *mockBookingProvider {
	return &mockBookingProvider{bookings: make(map[uuid.UUID]*bookingDomain.Booking)}
}

func (m *mockBookingProvider) GetBookingByID(ctx context.Context, id uuid.UUID) (*bookingDomain.Booking, error) {
	b, ok := m.bookings[id]
	if !ok {
		return nil, bookingDomain.ErrBookingNotFound
	}
	return b, nil
}

func (m *mockBookingProvider) UpdateBookingStatusOCC(ctx context.Context, booking *bookingDomain.Booking, expectedVersion int) error {
	booking.Version++
	m.bookings[booking.ID] = booking
	return nil
}

func (m *mockBookingProvider) UpdateBookingHandover(ctx context.Context, bookingID uuid.UUID, checkInMeter float64, expectedVersion int) error {
	b, ok := m.bookings[bookingID]
	if !ok {
		return bookingDomain.ErrBookingNotFound
	}
	b.CheckInMeterHours = &checkInMeter
	b.Version++
	return nil
}

func setupCompletionTest() (*rentalApp.RentalService, *rentalInfra.MemoryRentalRepository, *mockBookingProvider) {
	rentalRepo := rentalInfra.NewMemoryRentalRepository()
	bookingRepo := newMockBookingProvider()
	service := rentalApp.NewRentalService(rentalRepo, bookingRepo)
	return service, rentalRepo, bookingRepo
}

func TestRentalService_InitiateReturn(t *testing.T) {
	service, _, bookingRepo := setupCompletionTest()
	ctx := context.Background()

	renterID := uuid.New()
	ownerID := uuid.New()
	bookingID := uuid.New()

	booking := &bookingDomain.Booking{
		ID:           bookingID,
		RenterID:     renterID,
		OwnerID:      ownerID,
		Status:       bookingDomain.StatusActive,
		DurationDays: 7,
		Version:      1,
	}
	bookingRepo.bookings[bookingID] = booking

	t.Run("transitions ACTIVE to RETURN_INITIATED", func(t *testing.T) {
		updated, err := service.InitiateReturn(ctx, bookingID, renterID)
		if err != nil {
			t.Fatalf("expected successful return initiation, got: %v", err)
		}
		if updated.Status != bookingDomain.StatusReturnInitiated {
			t.Errorf("expected status RETURN_INITIATED, got %s", updated.Status)
		}
	})

	t.Run("rejects return on non-active booking", func(t *testing.T) {
		booking.Status = bookingDomain.StatusCompleted
		_, err := service.InitiateReturn(ctx, bookingID, renterID)
		if err == nil {
			t.Fatalf("expected error when initiating return on COMPLETED booking, got nil")
		}
	})
}

func TestRentalService_InspectReturn(t *testing.T) {
	service, _, bookingRepo := setupCompletionTest()
	ctx := context.Background()

	bookingID := uuid.New()
	inspectorID := uuid.New()
	checkIn := 1250.50

	booking := &bookingDomain.Booking{
		ID:                bookingID,
		Status:            bookingDomain.StatusReturnInitiated,
		CheckInMeterHours: &checkIn,
		DurationDays:      7,
		Version:           1,
	}
	bookingRepo.bookings[bookingID] = booking

	validPhotos := []rentalDomain.InspectionPhoto{
		{Angle: "FRONT_IMPLEMENT", StorageURI: "https://cdn.example.com/front.jpg"},
		{Angle: "REAR_COUNTERWEIGHT_ENGINE", StorageURI: "https://cdn.example.com/rear.jpg"},
		{Angle: "LEFT_SIDE_TRACKS", StorageURI: "https://cdn.example.com/left.jpg"},
		{Angle: "RIGHT_SIDE_TRACKS", StorageURI: "https://cdn.example.com/right.jpg"},
		{Angle: "HOUR_METER", StorageURI: "https://cdn.example.com/meter.jpg"},
		{Angle: "CABIN_INTERIOR", StorageURI: "https://cdn.example.com/cabin.jpg"},
	}

	t.Run("rejects regressive check-out meter reading", func(t *testing.T) {
		cmd := rentalApp.ReturnInspectCommand{
			BookingID:          bookingID,
			InspectorID:        inspectorID,
			CheckOutMeterHours: 1240.00, // lower than 1250.50
			Photos:             validPhotos,
		}
		_, err := service.InspectReturn(ctx, cmd)
		if err != rentalDomain.ErrReturnMeterRegression {
			t.Fatalf("expected ErrReturnMeterRegression, got: %v", err)
		}
	})

	t.Run("rejects incomplete photos", func(t *testing.T) {
		cmd := rentalApp.ReturnInspectCommand{
			BookingID:          bookingID,
			InspectorID:        inspectorID,
			CheckOutMeterHours: 1300.00,
			Photos:             validPhotos[:4],
		}
		_, err := service.InspectReturn(ctx, cmd)
		if err != rentalDomain.ErrMissingReturnPhotos {
			t.Fatalf("expected ErrMissingReturnPhotos, got: %v", err)
		}
	})

	t.Run("valid return inspection succeeds and updates booking", func(t *testing.T) {
		cmd := rentalApp.ReturnInspectCommand{
			BookingID:          bookingID,
			InspectorID:        inspectorID,
			CheckOutMeterHours: 1312.20,
			FluidLevelsOk:      true,
			SafetyGlassIntact:  true,
			Notes:              "Field check complete",
			Photos:             validPhotos,
		}
		insp, err := service.InspectReturn(ctx, cmd)
		if err != nil {
			t.Fatalf("expected clean inspection, got error: %v", err)
		}
		if insp.CheckOutMeterHours != 1312.20 {
			t.Errorf("expected meter 1312.20, got %.2f", insp.CheckOutMeterHours)
		}
		if booking.Status != bookingDomain.StatusDamagePendingReview {
			t.Errorf("expected status DAMAGE_PENDING_REVIEW, got %s", booking.Status)
		}
	})
}

func TestRentalService_DamageClaimLifecycle(t *testing.T) {
	service, _, bookingRepo := setupCompletionTest()
	ctx := context.Background()

	bookingID := uuid.New()
	renterID := uuid.New()
	ownerID := uuid.New()

	booking := &bookingDomain.Booking{
		ID:       bookingID,
		RenterID: renterID,
		OwnerID:  ownerID,
		Status:   bookingDomain.StatusDamagePendingReview,
		Version:  1,
	}
	bookingRepo.bookings[bookingID] = booking

	var claimID uuid.UUID
	t.Run("owner lodges damage claim", func(t *testing.T) {
		cmd := rentalApp.DamageClaimCommand{
			BookingID:          bookingID,
			ClaimantID:         ownerID,
			ComponentCategory:  "OPERATOR_CABIN",
			IncidentNarrative:  "Side window spiderweb crack from blast rock",
			RepairQuoteURI:     "https://cdn.example.com/quote.pdf",
			ClaimedAmountMinor: 850000,
			Currency:           "ETB",
		}
		docket, err := service.LodgeDamageClaim(ctx, cmd)
		if err != nil {
			t.Fatalf("expected clean damage claim lodging, got: %v", err)
		}
		claimID = docket.ID
		if docket.Status != "PENDING_REVIEW" {
			t.Errorf("expected status PENDING_REVIEW, got %s", docket.Status)
		}
		if !booking.DamageClaimSubmitted {
			t.Errorf("expected booking DamageClaimSubmitted to be true")
		}
	})

	t.Run("contractor accepts damage claim", func(t *testing.T) {
		cmd := rentalApp.RespondDamageClaimCommand{
			BookingID:           bookingID,
			ClaimID:             claimID,
			ActorID:             renterID,
			Action:              "ACCEPT",
			AcceptedAmountMinor: 850000,
		}
		docket, err := service.RespondToDamageClaim(ctx, cmd)
		if err != nil {
			t.Fatalf("expected clean response, got: %v", err)
		}
		if docket.Status != "ACCEPTED" {
			t.Errorf("expected status ACCEPTED, got %s", docket.Status)
		}
		if booking.AcceptedDamageMinor != 850000 {
			t.Errorf("expected accepted damage minor 850000, got %d", booking.AcceptedDamageMinor)
		}
	})
}

func TestRentalService_ReconcileFinalSettlement(t *testing.T) {
	service, _, bookingRepo := setupCompletionTest()
	ctx := context.Background()

	bookingID := uuid.New()
	checkIn := 1450.50
	checkOut := 1512.20

	booking := &bookingDomain.Booking{
		ID:                    bookingID,
		Status:                bookingDomain.StatusReturnInitiated,
		CheckInMeterHours:     &checkIn,
		CheckOutMeterHours:    &checkOut,
		DurationDays:          7,
		DailyRateAppliedMinor: 22400000 / 8, // ~2,800.00 ETB/hr
		BaseHireFeeMinor:      15750000,     // 157,500.00 ETB
		AttachmentsFeeMinor:   5600000,      // 56,000.00 ETB
		AcceptedDamageMinor:   850000,       // 8,500.00 ETB
		PrepaymentAmountMinor: 10675000,     // 106,750.00 ETB
		Currency:              "ETB",
		Version:               1,
	}
	bookingRepo.bookings[bookingID] = booking

	settlement, err := service.ReconcileFinalSettlement(ctx, bookingID)
	if err != nil {
		t.Fatalf("expected successful final reconciliation, got: %v", err)
	}

	if settlement.NetBalanceMinor <= 0 {
		t.Errorf("expected positive net balance due, got %d", settlement.NetBalanceMinor)
	}
	if booking.Status != bookingDomain.StatusPendingFinalPayment {
		t.Errorf("expected status PENDING_FINAL_PAYMENT, got %s", booking.Status)
	}
}

func TestRentalService_DoubleBlindReviews(t *testing.T) {
	service, _, bookingRepo := setupCompletionTest()
	ctx := context.Background()

	bookingID := uuid.New()
	renterID := uuid.New()
	ownerID := uuid.New()

	booking := &bookingDomain.Booking{
		ID:       bookingID,
		RenterID: renterID,
		OwnerID:  ownerID,
		Status:   bookingDomain.StatusCompleted,
		Version:  1,
	}
	bookingRepo.bookings[bookingID] = booking

	t.Run("first review remains un-published (isolated)", func(t *testing.T) {
		cmd := rentalApp.SubmitReviewCommand{
			BookingID: bookingID,
			AuthorID:  renterID,
			TargetID:  ownerID,
			Role:      "RENTER",
			Rating:    4.8,
			Feedback:  "Excellent machinery condition",
		}
		review, err := service.SubmitBookingReview(ctx, cmd)
		if err != nil {
			t.Fatalf("expected review submission to succeed, got: %v", err)
		}
		if review.IsPublished {
			t.Errorf("expected first review to remain unpublished, got isPublished=true")
		}
		if booking.Status == bookingDomain.StatusRated {
			t.Errorf("booking should not transition to RATED until both parties submit")
		}
	})

	t.Run("bilateral second review unblinds reviews and transitions to RATED", func(t *testing.T) {
		cmd := rentalApp.SubmitReviewCommand{
			BookingID: bookingID,
			AuthorID:  ownerID,
			TargetID:  renterID,
			Role:      "OWNER",
			Rating:    5.0,
			Feedback:  "Punctual wire settlement and careful machine handling",
		}
		review, err := service.SubmitBookingReview(ctx, cmd)
		if err != nil {
			t.Fatalf("expected second review submission to succeed, got: %v", err)
		}
		if !review.IsPublished {
			t.Errorf("expected bilateral submission to unblind reviews, got isPublished=false")
		}
		if booking.Status != bookingDomain.StatusRated {
			t.Errorf("expected booking to transition to RATED, got %s", booking.Status)
		}
	})
}
