package application_test

import (
	"context"
	"testing"
	"time"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	bookingMemory "github.com/dinkmachinery/dink-core/internal/modules/booking/infrastructure"
	"github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	"github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	rentalMemory "github.com/dinkmachinery/dink-core/internal/modules/rental/infrastructure"
	"github.com/google/uuid"
)

type testRig struct {
	rentalSvc   *application.RentalService
	rentalRepo  *rentalMemory.MemoryRentalRepository
	bookingRepo *bookingMemory.MemoryBookingRepository
	ctx         context.Context
	ownerID     uuid.UUID
	renterID    uuid.UUID
	operatorID  uuid.UUID
	fleetID     uuid.UUID
	booking     *bookingDomain.Booking
}

func setupTestRig(t *testing.T) *testRig {
	rentalRepo := rentalMemory.NewMemoryRentalRepository()
	bookingRepo := bookingMemory.NewMemoryBookingRepository()
	svc := application.NewRentalService(rentalRepo, bookingRepo)

	ownerID := uuid.New()
	renterID := uuid.New()
	operatorID := uuid.New()
	fleetID := uuid.New()

	rentalRepo.SetFleetBaseline(fleetID, 1250.50)

	booking := &bookingDomain.Booking{
		ID:                     uuid.New(),
		BookingReference:       "DNK-2026-TEST01",
		OwnerID:                ownerID,
		RenterID:               renterID,
		AssignedOperatorID:     &operatorID,
		FleetID:                fleetID,
		StartDate:              time.Now().UTC(),
		EndDate:                time.Now().UTC().Add(72 * time.Hour),
		DurationDays:           3,
		DailyRateAppliedMinor:  2000000, // 20,000.00 ETB
		GrossRentalAmountMinor: 6000000,
		PrepaymentAmountMinor:  3000000,
		Currency:               "ETB",
		Status:                 bookingDomain.StatusConfirmed,
		Version:                1,
	}

	ctx := context.Background()
	if err := bookingRepo.CreateBooking(ctx, booking); err != nil {
		t.Fatalf("failed to create test booking: %v", err)
	}

	return &testRig{
		rentalSvc:   svc,
		rentalRepo:  rentalRepo,
		bookingRepo: bookingRepo,
		ctx:         ctx,
		ownerID:     ownerID,
		renterID:    renterID,
		operatorID:  operatorID,
		fleetID:     fleetID,
		booking:     booking,
	}
}

func TestRentalService_HandoverInspection_MonotonicMeterFloor(t *testing.T) {
	rig := setupTestRig(t)

	// Case 1: Regressive meter hours (< 1250.50)
	_, err := rig.rentalSvc.InspectHandover(rig.ctx, application.HandoverInspectCommand{
		BookingID:         rig.booking.ID,
		InspectorID:       rig.ownerID,
		CheckInMeterHours: 1200.00, // Regression!
		FluidLevelsOk:     true,
		SafetyGlassIntact: true,
		Photos: []domain.InspectionPhoto{
			{Angle: "METER_CLOSEUP", StorageURI: "s3://inspections/meter.jpg"},
		},
	})
	if err == nil {
		t.Fatalf("expected error on regressive meter reading, got nil")
	}

	// Case 2: Invalid photo angle
	_, err = rig.rentalSvc.InspectHandover(rig.ctx, application.HandoverInspectCommand{
		BookingID:         rig.booking.ID,
		InspectorID:       rig.ownerID,
		CheckInMeterHours: 1250.50,
		FluidLevelsOk:     true,
		SafetyGlassIntact: true,
		Photos: []domain.InspectionPhoto{
			{Angle: "RANDOM_INVALID_ANGLE", StorageURI: "s3://inspections/invalid.jpg"},
		},
	})
	if err == nil {
		t.Fatalf("expected error on invalid photo angle, got nil")
	}

	// Case 3: Valid monotonic inspection
	insp, err := rig.rentalSvc.InspectHandover(rig.ctx, application.HandoverInspectCommand{
		BookingID:         rig.booking.ID,
		InspectorID:       rig.ownerID,
		CheckInMeterHours: 1252.30,
		FluidLevelsOk:     true,
		SafetyGlassIntact: true,
		Notes:             "Equipment in pristine operational order",
		Photos: []domain.InspectionPhoto{
			{Angle: "METER_CLOSEUP", StorageURI: "s3://inspections/meter.jpg"},
			{Angle: "FRONT_ELEVATION", StorageURI: "s3://inspections/front.jpg"},
		},
	})
	if err != nil {
		t.Fatalf("expected successful inspection, got: %v", err)
	}
	if insp.CheckInMeterHours != 1252.30 {
		t.Errorf("expected meter 1252.30, got %.2f", insp.CheckInMeterHours)
	}
}

func TestRentalService_DualHandoverPIN_And_ShiftInitialization(t *testing.T) {
	rig := setupTestRig(t)

	// Step 1: Pre-condition inspection
	_, err := rig.rentalSvc.InspectHandover(rig.ctx, application.HandoverInspectCommand{
		BookingID:         rig.booking.ID,
		InspectorID:       rig.ownerID,
		CheckInMeterHours: 1255.00,
		FluidLevelsOk:     true,
		SafetyGlassIntact: true,
		Photos: []domain.InspectionPhoto{
			{Angle: "METER_CLOSEUP", StorageURI: "s3://inspections/meter.jpg"},
		},
	})
	if err != nil {
		t.Fatalf("inspection failed: %v", err)
	}

	// Step 2: Owner generates 6-digit PIN
	challenge, err := rig.rentalSvc.GenerateHandoverPIN(rig.ctx, rig.booking.ID, rig.ownerID)
	if err != nil {
		t.Fatalf("failed to generate handover PIN: %v", err)
	}
	if len(challenge.ChallengePIN) != 6 {
		t.Errorf("expected 6-digit PIN, got %s", challenge.ChallengePIN)
	}

	// Step 3: Unauthorized user cannot confirm
	randomUser := uuid.New()
	_, err = rig.rentalSvc.ConfirmHandover(rig.ctx, application.HandoverConfirmCommand{
		BookingID:    rig.booking.ID,
		CallerUserID: randomUser,
		PIN:          challenge.ChallengePIN,
	})
	if err != bookingDomain.ErrUnauthorizedAction {
		t.Fatalf("expected ErrUnauthorizedAction, got %v", err)
	}

	// Step 4: Incorrect PIN increments failed attempts
	_, err = rig.rentalSvc.ConfirmHandover(rig.ctx, application.HandoverConfirmCommand{
		BookingID:    rig.booking.ID,
		CallerUserID: rig.renterID,
		PIN:          "000000", // Wrong PIN
	})
	if err != domain.ErrInvalidChallengePIN {
		t.Fatalf("expected ErrInvalidChallengePIN, got %v", err)
	}

	// Step 5: Successful confirmation with contractor touch signature & GPS
	lat := 9.0300
	lng := 38.7400
	sig := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
	b, err := rig.rentalSvc.ConfirmHandover(rig.ctx, application.HandoverConfirmCommand{
		BookingID:            rig.booking.ID,
		CallerUserID:         rig.renterID,
		PIN:                  challenge.ChallengePIN,
		TouchSignatureBase64: sig,
		ContractorLat:        &lat,
		ContractorLng:        &lng,
	})
	if err != nil {
		t.Fatalf("failed to confirm handover: %v", err)
	}
	if b.Status != bookingDomain.StatusActive {
		t.Errorf("expected booking status ACTIVE, got %s", b.Status)
	}
	if b.ActivatedAt == nil {
		t.Errorf("expected ActivatedAt to be populated")
	}

	// Step 6: Verify Day 1 RentalShift initialized with 8.0 hr quota (AR-03)
	shift, err := rig.rentalRepo.GetActiveRentalShift(rig.ctx, rig.booking.ID)
	if err != nil {
		t.Fatalf("failed to find initialized shift: %v", err)
	}
	if shift.AllocatedHours != 8.00 {
		t.Errorf("expected allocated hours 8.00, got %.2f", shift.AllocatedHours)
	}
	if shift.StartMeterHours != 1255.00 {
		t.Errorf("expected start meter 1255.00, got %.2f", shift.StartMeterHours)
	}

	// Step 7: Re-confirming with used challenge should fail
	_, err = rig.rentalSvc.ConfirmHandover(rig.ctx, application.HandoverConfirmCommand{
		BookingID:    rig.booking.ID,
		CallerUserID: rig.renterID,
		PIN:          challenge.ChallengePIN,
	})
	if err == nil {
		t.Fatalf("expected error when confirming already used challenge, got nil")
	}
}

func TestRentalService_ActiveRentalHUD_TelemetryCompilation(t *testing.T) {
	rig := setupTestRig(t)

	// Activate booking
	challenge, _ := rig.rentalSvc.GenerateHandoverPIN(rig.ctx, rig.booking.ID, rig.ownerID)
	_, err := rig.rentalSvc.ConfirmHandover(rig.ctx, application.HandoverConfirmCommand{
		BookingID:    rig.booking.ID,
		CallerUserID: rig.renterID,
		PIN:          challenge.ChallengePIN,
	})
	if err != nil {
		t.Fatalf("failed to activate booking: %v", err)
	}

	// Progress the active shift
	shift, _ := rig.rentalRepo.GetActiveRentalShift(rig.ctx, rig.booking.ID)
	meter := 1256.0
	_ = rig.rentalRepo.UpdateRentalShiftProgress(rig.ctx, shift.ID, 4.5, 0.0, &meter, false)

	// Fetch HUD
	hud, err := rig.rentalSvc.GetActiveRentalHUD(rig.ctx, rig.booking.ID, rig.renterID)
	if err != nil {
		t.Fatalf("failed to retrieve HUD: %v", err)
	}

	if hud.Status != "ACTIVE" {
		t.Errorf("expected HUD status ACTIVE, got %s", hud.Status)
	}
	if hud.AccumulatedHoursToday != 4.5 {
		t.Errorf("expected 4.5 hours accumulated today, got %.2f", hud.AccumulatedHoursToday)
	}
	if hud.AllowedHoursPerDay != 8.0 {
		t.Errorf("expected 8.0 allowed hours, got %.2f", hud.AllowedHoursPerDay)
	}
	if !hud.CanExtend {
		t.Errorf("expected CanExtend to be true")
	}
	if !hud.CanReportBreakdown {
		t.Errorf("expected CanReportBreakdown to be true")
	}
}

func TestRentalService_ForwardExtension_QuoteAndConflict(t *testing.T) {
	rig := setupTestRig(t)

	// Activate booking
	challenge, _ := rig.rentalSvc.GenerateHandoverPIN(rig.ctx, rig.booking.ID, rig.ownerID)
	_, _ = rig.rentalSvc.ConfirmHandover(rig.ctx, application.HandoverConfirmCommand{
		BookingID:    rig.booking.ID,
		CallerUserID: rig.renterID,
		PIN:          challenge.ChallengePIN,
	})

	// Case 1: Extension date not forward
	_, err := rig.rentalSvc.RequestExtension(rig.ctx, application.RentalExtensionCommand{
		BookingID:  rig.booking.ID,
		RenterID:   rig.renterID,
		NewEndDate: rig.booking.EndDate.Add(-24 * time.Hour), // Earlier than current end date!
	})
	if err == nil {
		t.Fatalf("expected error on non-forward extension date, got nil")
	}

	// Case 2: Equipment reserved on forward dates (conflict)
	rig.rentalRepo.SetConflict(rig.fleetID, true)
	_, err = rig.rentalSvc.RequestExtension(rig.ctx, application.RentalExtensionCommand{
		BookingID:  rig.booking.ID,
		RenterID:   rig.renterID,
		NewEndDate: rig.booking.EndDate.Add(48 * time.Hour),
	})
	if err != domain.ErrExtensionDateUnavailable {
		t.Fatalf("expected ErrExtensionDateUnavailable, got %v", err)
	}

	// Case 3: Successful 2-day extension quote with 50% prepayment (AR-04)
	rig.rentalRepo.SetConflict(rig.fleetID, false)
	newEndDate := rig.booking.EndDate.Add(48 * time.Hour)
	ext, err := rig.rentalSvc.RequestExtension(rig.ctx, application.RentalExtensionCommand{
		BookingID:  rig.booking.ID,
		RenterID:   rig.renterID,
		NewEndDate: newEndDate,
	})
	if err != nil {
		t.Fatalf("expected successful extension request, got: %v", err)
	}

	if ext.ExtensionDays != 2 {
		t.Errorf("expected 2 extension days, got %d", ext.ExtensionDays)
	}
	// 2 days * 20,000.00 = 40,000.00 ETB (4,000,000 minor)
	if ext.IncrementalGrossMinor != 4000000 {
		t.Errorf("expected 4000000 incremental gross, got %d", ext.IncrementalGrossMinor)
	}
	// 50% prepay = 20,000.00 ETB (2,000,000 minor)
	if ext.IncrementalPrepaymentMinor != 2000000 {
		t.Errorf("expected 2000000 incremental prepayment, got %d", ext.IncrementalPrepaymentMinor)
	}
}

func TestRentalService_MechanicalBreakdown_HaltAndResume(t *testing.T) {
	rig := setupTestRig(t)

	// Activate booking
	challenge, _ := rig.rentalSvc.GenerateHandoverPIN(rig.ctx, rig.booking.ID, rig.ownerID)
	_, _ = rig.rentalSvc.ConfirmHandover(rig.ctx, application.HandoverConfirmCommand{
		BookingID:    rig.booking.ID,
		CallerUserID: rig.renterID,
		PIN:          challenge.ChallengePIN,
	})

	// Step 1: Report breakdown
	breakdown, err := rig.rentalSvc.ReportMechanicalBreakdown(rig.ctx, application.BreakdownReportCommand{
		BookingID:              rig.booking.ID,
		ReporterID:             rig.renterID,
		FaultCategory:          "HYDRAULIC_CIRCUIT",
		FaultDescription:       "Main boom pressure loss and ruptured hose",
		DiagnosticTroubleCodes: []string{"HYD-042", "ERR-910"},
		CurrentHourMeter:       1262.40,
		IncidentPhotos:         []string{"s3://photos/ruptured_hose.jpg"},
	})
	if err != nil {
		t.Fatalf("failed to report breakdown: %v", err)
	}

	if breakdown.Status != "ACTIVE_BREAKDOWN_PAUSED" {
		t.Errorf("expected status ACTIVE_BREAKDOWN_PAUSED, got %s", breakdown.Status)
	}

	// Verify booking transitioned to ACTIVE_BREAKDOWN_PAUSED
	updatedBooking, _ := rig.bookingRepo.GetBookingByID(rig.ctx, rig.booking.ID)
	if updatedBooking.Status != bookingDomain.StatusActiveBreakdownPaused {
		t.Errorf("expected booking status ACTIVE_BREAKDOWN_PAUSED, got %s", updatedBooking.Status)
	}

	// Verify shift telemetry paused
	shift, _ := rig.rentalRepo.GetActiveRentalShift(rig.ctx, rig.booking.ID)
	if !shift.IsPaused {
		t.Errorf("expected active shift to be paused during breakdown")
	}

	// Verify HUD reflects shift timer paused
	hud, _ := rig.rentalSvc.GetActiveRentalHUD(rig.ctx, rig.booking.ID, rig.renterID)
	if !hud.ShiftTimerPaused {
		t.Errorf("expected HUD ShiftTimerPaused to be true")
	}

	// Step 2: Resolve breakdown
	resolved, err := rig.rentalSvc.ResolveMechanicalBreakdown(rig.ctx, application.BreakdownResolveCommand{
		BreakdownID:     breakdown.ID,
		ResolverID:      rig.ownerID,
		MechanicName:    "Dawit Alemayehu",
		MechanicPhone:   "+251911445566",
		ResolutionNotes: "Replaced blown high-pressure hose and flushed hydraulic circuit.",
	})
	if err != nil {
		t.Fatalf("failed to resolve breakdown: %v", err)
	}
	if resolved.Status != "RESOLVED" {
		t.Errorf("expected breakdown status RESOLVED, got %s", resolved.Status)
	}

	// Verify booking restored to ACTIVE
	restoredBooking, _ := rig.bookingRepo.GetBookingByID(rig.ctx, rig.booking.ID)
	if restoredBooking.Status != bookingDomain.StatusActive {
		t.Errorf("expected booking status ACTIVE after resolution, got %s", restoredBooking.Status)
	}

	// Verify shift unpaused
	resumedShift, _ := rig.rentalRepo.GetActiveRentalShift(rig.ctx, rig.booking.ID)
	if resumedShift.IsPaused {
		t.Errorf("expected active shift to be unpaused after resolution")
	}
}
