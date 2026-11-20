package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	bookingInfra "github.com/dinkmachinery/dink-core/internal/modules/booking/infrastructure"
	"github.com/dinkmachinery/dink-core/internal/modules/logistics/application"
	"github.com/dinkmachinery/dink-core/internal/modules/logistics/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/logistics/infrastructure"
	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	rentalInfra "github.com/dinkmachinery/dink-core/internal/modules/rental/infrastructure"
	"github.com/google/uuid"
)

func setupLogisticsHarness(t *testing.T) (
	*application.LogisticsService,
	*rentalApp.RentalService,
	*bookingInfra.MemoryBookingRepository,
	*infrastructure.MemoryLogisticsRepository,
	uuid.UUID, // renterID
	uuid.UUID, // ownerID
	uuid.UUID, // carrier1ID
	uuid.UUID, // carrier2ID
	uuid.UUID, // machineFleetID
	uuid.UUID, // lowboy1FleetID
	uuid.UUID, // lowboy2FleetID
	uuid.UUID, // parentBookingID
) {
	t.Helper()
	ctx := context.Background()

	bookingRepo := bookingInfra.NewMemoryBookingRepository()
	logRepo := infrastructure.NewMemoryLogisticsRepository(bookingRepo)
	rentalRepo := rentalInfra.NewMemoryRentalRepository()

	logSvc := application.NewLogisticsService(logRepo, bookingRepo)
	rentalSvc := rentalApp.NewRentalService(rentalRepo, bookingRepo)
	rentalSvc.SetLogisticsChecker(logSvc)

	renterID := uuid.New()
	ownerID := uuid.New()
	carrier1ID := uuid.New()
	carrier2ID := uuid.New()
	machineFleetID := uuid.New()
	lowboy1FleetID := uuid.New()
	lowboy2FleetID := uuid.New()
	parentBookingID := uuid.New()

	logRepo.SeedFleetProfile(&domain.FleetLogisticsProfile{
		FleetID:          machineFleetID,
		OwnerID:          ownerID,
		Make:             "Caterpillar",
		Model:            "320D",
		RootCategoryCode: "EXCAVATION",
		SelfMoveCapable:  false,
		DailyRateMinor:   3500000,
		Status:           "ACTIVE",
		CountryCode:      "ET",
		Location:         domain.GeoPoint{Latitude: 8.9806, Longitude: 38.7578, LocationName: "Bole Depot"},
		OwnerFullName:    "Almaz Tesfaye",
		OwnerPhoneNumber: "+251922222222",
		MachineSpecs: domain.MachineLogisticsSpecs{
			OperatingWeightKg: 21800,
			TransportLengthMm: 9530,
			TransportWidthMm:  2980,
			TransportHeightMm: 3150,
		},
	})

	logRepo.SeedFleetProfile(&domain.FleetLogisticsProfile{
		FleetID:          lowboy1FleetID,
		OwnerID:          carrier1ID,
		Make:             "Sinotruk",
		Model:            "Howo Tri-Axle Lowbed",
		RootCategoryCode: "LOGISTICS",
		SelfMoveCapable:  true,
		DailyRateMinor:   1800000,
		Status:           "ACTIVE",
		CountryCode:      "ET",
		Location:         domain.GeoPoint{Latitude: 8.9124, Longitude: 38.7412, LocationName: "Kality Depot"},
		OwnerFullName:    "Ato Girma Wolde",
		OwnerPhoneNumber: "+251911449120",
		TrailerProfile: domain.MatchedCarrier{
			FleetID:               lowboy1FleetID,
			OwnerID:               carrier1ID,
			Make:                  "Sinotruk",
			Model:                 "Howo Tri-Axle Lowbed",
			OwnerTradingName:      "Girma Heavy Haulage PLC",
			OwnerPhoneNumber:      "+251911449120",
			PayloadCapacityKg:     34000,
			DeckWellLengthMm:      11000,
			DeckWidthMm:           3100,
			OutriggersExtensionMm: 250,
			DeckHeightGroundMm:    900,
			TareWeightKg:          12500,
			AxleCount:             3,
			DailyRateMinor:        1800000,
			Currency:              "ETB",
			AverageRating:         4.95,
		},
	})

	logRepo.SeedFleetProfile(&domain.FleetLogisticsProfile{
		FleetID:          lowboy2FleetID,
		OwnerID:          carrier2ID,
		Make:             "Mercedes-Benz",
		Model:            "Actros Lowbed",
		RootCategoryCode: "LOGISTICS",
		SelfMoveCapable:  true,
		DailyRateMinor:   1950000,
		Status:           "ACTIVE",
		CountryCode:      "ET",
		Location:         domain.GeoPoint{Latitude: 8.7900, Longitude: 38.9100, LocationName: "Bishoftu Yard"},
		OwnerFullName:    "Samuel Odhiambo",
		OwnerPhoneNumber: "+251911773300",
		TrailerProfile: domain.MatchedCarrier{
			FleetID:               lowboy2FleetID,
			OwnerID:               carrier2ID,
			Make:                  "Mercedes-Benz",
			Model:                 "Actros Lowbed",
			OwnerTradingName:      "Corridor Lowboy Logistics",
			OwnerPhoneNumber:      "+251911773300",
			PayloadCapacityKg:     36000,
			DeckWellLengthMm:      11500,
			DeckWidthMm:           3200,
			OutriggersExtensionMm: 300,
			DeckHeightGroundMm:    900,
			TareWeightKg:          13000,
			AxleCount:             3,
			DailyRateMinor:        1950000,
			Currency:              "ETB",
			AverageRating:         4.88,
		},
	})

	now := time.Now().UTC()
	err := bookingRepo.CreateBooking(ctx, &bookingDomain.Booking{
		ID:                        parentBookingID,
		BookingReference:          "BK-PARENT-001",
		RenterID:                  renterID,
		OwnerID:                   ownerID,
		FleetID:                   machineFleetID,
		BookingKind:               "MACHINERY",
		StartDate:                 now.AddDate(0, 0, 2),
		EndDate:                   now.AddDate(0, 0, 7),
		DurationDays:              5,
		DailyRateAppliedMinor:     3500000,
		BaseHireFeeMinor:          17500000,
		GrossRentalAmountMinor:    17500000,
		PrepaymentAmountMinor:     8750000,
		PlatformCommissionRateBps: 1000,
		PlatformCommissionMinor:   1750000,
		Currency:                  "ETB",
		SiteLocationAddress:       "Mojo Site B",
		Status:                    bookingDomain.StatusConfirmed,
		Version:                   1,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	})
	if err != nil {
		t.Fatalf("seed parent booking failed: %v", err)
	}

	return logSvc, rentalSvc, bookingRepo, logRepo, renterID, ownerID, carrier1ID, carrier2ID, machineFleetID, lowboy1FleetID, lowboy2FleetID, parentBookingID
}

// TestLogistics_ParentBookingBlockedUntilTransportDelivery verifies D-05-08-001.
func TestLogistics_ParentBookingBlockedUntilTransportDelivery(t *testing.T) {
	ctx := context.Background()
	logSvc, rentalSvc, _, _, renterID, ownerID, carrier1ID, _, _, lowboy1FleetID, _, parentID := setupLogisticsHarness(t)

	// 1. Match carriers
	matchRes, err := logSvc.MatchCarriers(ctx, application.MatchCarriersCommand{
		ParentBookingID: parentID,
		CallerUserID:    ownerID,
		SearchRadiusKm:  75,
		PickupDate:      time.Now().UTC().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("MatchCarriers failed: %v", err)
	}
	if len(matchRes.AvailableLogisticsFleets) != 2 {
		t.Fatalf("expected 2 compatible lowboy carriers, got %d", len(matchRes.AvailableLogisticsFleets))
	}

	// 2. Create Sub-Booking (2.5% commission discount)
	subRes, err := logSvc.CreateSubBooking(ctx, application.CreateSubBookingCommand{
		ParentBookingID:  parentID,
		CallerUserID:     ownerID,
		LogisticsFleetID: lowboy1FleetID,
		PickupDateTime:   time.Now().UTC().Add(24 * time.Hour),
		Origin:           domain.GeoPoint{Latitude: 8.9124, Longitude: 38.7412, LocationName: "Kality Depot"},
		Destination:      domain.GeoPoint{Latitude: 8.5412, Longitude: 39.2718, LocationName: "Mojo Site B"},
	})
	if err != nil {
		t.Fatalf("CreateSubBooking failed: %v", err)
	}
	if subRes.PlatformCommissionBps != 250 || subRes.CommissionOwedMinor != 90000 {
		t.Fatalf("expected 2.5%% freight commission (250 bps / 90000), got %d bps / %d", subRes.PlatformCommissionBps, subRes.CommissionOwedMinor)
	}

	// 3. Advance carrier dispatch through DISPATCHED -> LOADING_INSPECTION -> IN_TRANSIT
	_, err = logSvc.AdvanceDispatchStage(ctx, application.AdvanceDispatchStageCommand{
		SubBookingID: subRes.SubBookingID,
		CallerUserID: carrier1ID,
		TargetStage:  bookingDomain.StatusDispatched,
	})
	if err != nil {
		t.Fatalf("DISPATCHED failed: %v", err)
	}

	_, err = logSvc.AdvanceDispatchStage(ctx, application.AdvanceDispatchStageCommand{
		SubBookingID: subRes.SubBookingID,
		CallerUserID: carrier1ID,
		TargetStage:  bookingDomain.StatusLoadingInspection,
		WinchingPhotos: []domain.WinchingPhoto{
			{Angle: domain.WinchingAngleRigProfile, StorageURI: "https://cdn.dinkmachinery.com/1.webp"},
			{Angle: domain.WinchingAngleChainBinders, StorageURI: "https://cdn.dinkmachinery.com/2.webp"},
			{Angle: domain.WinchingAngleBoomPinned, StorageURI: "https://cdn.dinkmachinery.com/3.webp"},
			{Angle: domain.WinchingAngleRampsLocked, StorageURI: "https://cdn.dinkmachinery.com/4.webp"},
		},
	})
	if err != nil {
		t.Fatalf("LOADING_INSPECTION failed: %v", err)
	}

	_, err = logSvc.AdvanceDispatchStage(ctx, application.AdvanceDispatchStageCommand{
		SubBookingID: subRes.SubBookingID,
		CallerUserID: carrier1ID,
		TargetStage:  bookingDomain.StatusInTransit,
	})
	if err != nil {
		t.Fatalf("IN_TRANSIT failed: %v", err)
	}

	// 4. Generate handover PIN while transport is still IN_TRANSIT
	challenge, err := rentalSvc.GenerateHandoverPIN(ctx, parentID, ownerID)
	if err != nil {
		t.Fatalf("GenerateHandoverPIN failed: %v", err)
	}

	// 5. Assert ConfirmHandover is BLOCKED while transport sub-booking is IN_TRANSIT (D-05-08-001)
	_, err = rentalSvc.ConfirmHandover(ctx, rentalApp.HandoverConfirmCommand{
		BookingID:    parentID,
		CallerUserID: renterID,
		PIN:          challenge.ChallengePIN,
	})
	if !errors.Is(err, domain.ErrParentBookingBlockedUntilDelivery) {
		t.Fatalf("expected ErrParentBookingBlockedUntilDelivery while sub-booking is IN_TRANSIT, got %v", err)
	}

	// 6. Complete delivery to site (READY_FOR_PICKUP) and confirm handover succeeds
	_, err = logSvc.AdvanceDispatchStage(ctx, application.AdvanceDispatchStageCommand{
		SubBookingID: subRes.SubBookingID,
		CallerUserID: carrier1ID,
		TargetStage:  bookingDomain.StatusReadyForPickup,
	})
	if err != nil {
		t.Fatalf("READY_FOR_PICKUP failed: %v", err)
	}

	activated, err := rentalSvc.ConfirmHandover(ctx, rentalApp.HandoverConfirmCommand{
		BookingID:    parentID,
		CallerUserID: renterID,
		PIN:          challenge.ChallengePIN,
	})
	if err != nil {
		t.Fatalf("expected ConfirmHandover to succeed after transport arrival, got %v", err)
	}
	if activated.Status != bookingDomain.StatusActive {
		t.Fatalf("expected parent booking ACTIVE, got %s", activated.Status)
	}
}

func TestLogistics_ExternalTransportDeclaration(t *testing.T) {
	ctx := context.Background()
	logSvc, _, _, _, _, ownerID, _, _, _, _, _, parentID := setupLogisticsHarness(t)

	ext, status, err := logSvc.RegisterExternalTransport(ctx, application.RegisterExternalTransportCommand{
		ParentBookingID:       parentID,
		CallerUserID:          ownerID,
		CarrierName:           "Private In-House Fleet",
		DriverFullName:        "Tadesse Bekele",
		DriverPhoneNumber:     "+251911884912",
		PrimeMoverPlateNumber: "ET-3-44192",
		EstimatedArrivalAt:    time.Now().UTC().Add(4 * time.Hour),
	})
	if err != nil {
		t.Fatalf("RegisterExternalTransport failed: %v", err)
	}
	if status != bookingDomain.StatusReadyForPickup {
		t.Fatalf("expected parent status READY_FOR_PICKUP, got %s", status)
	}
	if ext.DriverFullName != "Tadesse Bekele" {
		t.Fatalf("unexpected driver name %s", ext.DriverFullName)
	}
}

func TestLogistics_TransshipmentAndReverseCascadeForceMajeure(t *testing.T) {
	ctx := context.Background()
	logSvc, _, _, logRepo, renterID, ownerID, carrier1ID, _, _, lowboy1FleetID, _, parentID := setupLogisticsHarness(t)

	subRes, err := logSvc.CreateSubBooking(ctx, application.CreateSubBookingCommand{
		ParentBookingID:  parentID,
		CallerUserID:     ownerID,
		LogisticsFleetID: lowboy1FleetID,
		PickupDateTime:   time.Now().UTC().Add(2 * time.Hour),
		Origin:           domain.GeoPoint{Latitude: 8.9124, Longitude: 38.7412, LocationName: "Kality"},
		Destination:      domain.GeoPoint{Latitude: 8.5412, Longitude: 39.2718, LocationName: "Mojo"},
	})
	if err != nil {
		t.Fatalf("CreateSubBooking failed: %v", err)
	}

	// Move to IN_TRANSIT
	_, _ = logSvc.AdvanceDispatchStage(ctx, application.AdvanceDispatchStageCommand{SubBookingID: subRes.SubBookingID, CallerUserID: carrier1ID, TargetStage: bookingDomain.StatusDispatched})
	_, _ = logSvc.AdvanceDispatchStage(ctx, application.AdvanceDispatchStageCommand{
		SubBookingID: subRes.SubBookingID,
		CallerUserID: carrier1ID,
		TargetStage:  bookingDomain.StatusLoadingInspection,
		WinchingPhotos: []domain.WinchingPhoto{
			{Angle: domain.WinchingAngleRigProfile, StorageURI: "https://cdn.dinkmachinery.com/1.webp"},
			{Angle: domain.WinchingAngleChainBinders, StorageURI: "https://cdn.dinkmachinery.com/2.webp"},
			{Angle: domain.WinchingAngleBoomPinned, StorageURI: "https://cdn.dinkmachinery.com/3.webp"},
			{Angle: domain.WinchingAngleRampsLocked, StorageURI: "https://cdn.dinkmachinery.com/4.webp"},
		},
	})
	_, _ = logSvc.AdvanceDispatchStage(ctx, application.AdvanceDispatchStageCommand{SubBookingID: subRes.SubBookingID, CallerUserID: carrier1ID, TargetStage: bookingDomain.StatusInTransit})

	// Report roadside breakdown near Bishoftu (within 50km of alternate carrier)
	alertRes, err := logSvc.ReportTransshipmentAlert(ctx, application.TransshipmentAlertCommand{
		SubBookingID:       subRes.SubBookingID,
		CallerUserID:       carrier1ID,
		BreakdownLocation:  domain.GeoPoint{Latitude: 8.7500, Longitude: 38.9000, LocationName: "Expressway KM 38"},
		FailureCategory:    domain.FailurePrimeMoverEngineBlown,
		FailureDescription: "Turbocharger shaft sheared on expressway shoulder.",
		Photos:             []string{"https://cdn.dinkmachinery.com/breakdown/engine.webp"},
	})
	if err != nil {
		t.Fatalf("ReportTransshipmentAlert failed: %v", err)
	}
	if alertRes.SubBookingStatus != bookingDomain.StatusTransshipmentRequired {
		t.Fatalf("expected TRANSSHIPMENT_REQUIRED, got %s", alertRes.SubBookingStatus)
	}
	if alertRes.ParentBookingStatus != bookingDomain.StatusAwaitingFleetReadiness {
		t.Fatalf("expected parent AWAITING_FLEET_READINESS, got %s", alertRes.ParentBookingStatus)
	}
	if alertRes.CandidateCarriersCount != 1 {
		t.Fatalf("expected 1 replacement carrier within 50km, got %d", alertRes.CandidateCarriersCount)
	}

	// Simulate transshipment window failure -> execute Force Majeure Reverse Cascade
	cascadeRes, err := logSvc.TriggerReverseCascade(ctx, application.ReverseCascadeCommand{
		ParentBookingID: parentID,
		CallerUserID:    renterID,
		Reason:          "No replacement lowboy arrived within 60m SLA",
	})
	if err != nil {
		t.Fatalf("TriggerReverseCascade failed: %v", err)
	}

	if cascadeRes.ParentStatus != bookingDomain.StatusCancelled || cascadeRes.SubStatus != bookingDomain.StatusCancelledCarrierDefault {
		t.Fatalf("unexpected cascade statuses: parent=%s sub=%s", cascadeRes.ParentStatus, cascadeRes.SubStatus)
	}
	if cascadeRes.Cascade.RenterRefundDueMinor != 8750000 {
		t.Fatalf("expected 100%% renter refund 8750000, got %d", cascadeRes.Cascade.RenterRefundDueMinor)
	}
	if cascadeRes.Cascade.OwnerPenaltyLeviedMinor != 0 {
		t.Fatalf("expected 0 owner penalty under Force Majeure, got %d", cascadeRes.Cascade.OwnerPenaltyLeviedMinor)
	}
	if cascadeRes.Cascade.OwnerCommissionRestitutionMinor != 1750000 {
		t.Fatalf("expected commission restitution 1750000, got %d", cascadeRes.Cascade.OwnerCommissionRestitutionMinor)
	}
	if !logRepo.IsAvailabilityReleased(parentID) || !logRepo.IsAvailabilityReleased(subRes.SubBookingID) {
		t.Fatalf("expected availability locks released for both parent and sub-booking")
	}
	if logRepo.CarrierRatingDeductionTotal(carrier1ID) != 0.50 {
		t.Fatalf("expected -0.50 rating deduction for defaulting carrier, got %.2f", logRepo.CarrierRatingDeductionTotal(carrier1ID))
	}
}
