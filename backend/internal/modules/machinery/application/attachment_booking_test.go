package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/infrastructure"
)

func TestAttachment_BundledLockPreventsStandaloneRental(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryMachineryRepository()
	service := application.NewAttachmentService(repo, repo)

	ownerID := uuid.New()

	// 1. Register a prime host fleet (Caterpillar 320D)
	fleet := &domain.FleetAsset{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		SubcategoryID:  uuid.MustParse("019a8f00-0002-7000-8000-000000000001"),
		Make:           "Caterpillar",
		Model:          "320D",
		DailyRateMinor: 2250000, // 22,500 ETB
		Status:         domain.FleetStatusActive,
	}
	_, err := repo.CreateFleet(ctx, fleet)
	require.NoError(t, err)

	// 2. Register an attachment (Soosan SB-81 Hydraulic Breaker)
	photos := []domain.AttachmentPhoto{
		{Angle: domain.AttachmentPhotoAngleOverall, StorageURI: "s3://media/overall.webp"},
		{Angle: domain.AttachmentPhotoAngleBracket, StorageURI: "s3://media/bracket.webp"},
		{Angle: domain.AttachmentPhotoAnglePins, StorageURI: "s3://media/pins.webp"},
		{Angle: domain.AttachmentPhotoAngleToolTip, StorageURI: "s3://media/tooltip.webp"},
	}

	att, err := service.RegisterAttachment(ctx, application.RegisterAttachmentInput{
		OwnerID:                  ownerID,
		AttachmentType:           domain.AttachmentTypeHydraulicBreaker,
		Make:                     "Soosan",
		Model:                    "SB-81",
		SubTypeSize:              "135mm Chisel",
		SerialNumber:             "SN-SB81-2026-99120",
		OperatingWeightKg:        1850.0,
		MinCarrierWeightTonnes:   18.0,
		MaxCarrierWeightTonnes:   26.0,
		RequiredFlowMinLmin:      120.0,
		RequiredFlowMaxLmin:      150.0,
		RequiredPressureBar:      180.0,
		PinDiameterMm:            80.0,
		PinSpreadMm:              450.0,
		CompatibleSubcategories:  []string{"EXCAVATOR_MINI"},
		StandaloneDailyRateMinor: 1200000, // 12,000 ETB
		SecurityDepositMinor:     2500000,
		MinRentalDays:            1,
		Photos:                   photos,
	})
	require.NoError(t, err)
	require.NotNil(t, att)

	// 3. Link breaker to Caterpillar 320D with bundled discount (8,000 ETB)
	link, err := service.LinkAttachmentToFleet(ctx, ownerID, application.LinkAttachmentInput{
		FleetID:               fleet.ID,
		AttachmentID:          att.ID,
		BundledDailyRateMinor: 800000,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(800000), link.BundledDailyRateMinor)

	// 4. Initial check: verify both excavator and breaker are available for Nov 5 to Nov 12
	startDate := time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC)

	availCheck, err := service.CheckBundledAvailability(ctx, application.CheckBundledAvailabilityInput{
		FleetID:       fleet.ID,
		AttachmentIDs: []uuid.UUID{att.ID},
		StartDate:     startDate,
		EndDate:       endDate,
	})
	require.NoError(t, err)
	assert.True(t, availCheck.FleetAvailable)
	assert.True(t, availCheck.AllAttachmentsAvailable)
	assert.Equal(t, int64(3050000), availCheck.TotalCombinedDailyRateMinor) // 22,500 + 8,000 = 30,500 ETB

	// 5. Booking reservation executed: Renter books the excavator with the bundled breaker
	bookingID := uuid.New()
	err = service.ReserveAttachment(ctx, att.ID, &bookingID, startDate, endDate, link.BundledDailyRateMinor)
	require.NoError(t, err)

	// 6. Concurrency Defense (D-05-03-001, FL-34):
	// A second contractor tries to rent the SAME breaker as a standalone rental for an overlapping window (Nov 8 to Nov 15)
	overlappingStart := time.Date(2026, 11, 8, 0, 0, 0, 0, time.UTC)
	overlappingEnd := time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC)

	secondBookingID := uuid.New()
	err = service.ReserveAttachment(ctx, att.ID, &secondBookingID, overlappingStart, overlappingEnd, att.StandaloneDailyRateMinor)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrAttachmentScheduleConflict,
		"bundled booking must atomically prevent concurrent standalone rental of the same attachment")

	// 7. Check availability for overlapping window now reports attachment as unavailable
	availCheckSecond, err := service.CheckBundledAvailability(ctx, application.CheckBundledAvailabilityInput{
		FleetID:       fleet.ID,
		AttachmentIDs: []uuid.UUID{att.ID},
		StartDate:     overlappingStart,
		EndDate:       overlappingEnd,
	})
	require.NoError(t, err)
	assert.True(t, availCheckSecond.FleetAvailable)
	assert.False(t, availCheckSecond.AllAttachmentsAvailable)
	assert.False(t, availCheckSecond.AttachmentStatuses[0].IsAvailable)

	// 8. Non-overlapping window (e.g. Dec 1 to Dec 5) succeeds
	futureStart := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	futureEnd := time.Date(2026, 12, 5, 0, 0, 0, 0, time.UTC)
	futureBookingID := uuid.New()
	err = service.ReserveAttachment(ctx, att.ID, &futureBookingID, futureStart, futureEnd, att.StandaloneDailyRateMinor)
	require.NoError(t, err)
}

func TestAttachment_LinkValidationRejectsExceedingStandaloneRate(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryMachineryRepository()
	service := application.NewAttachmentService(repo, repo)

	ownerID := uuid.New()
	fleet := &domain.FleetAsset{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		SubcategoryID:  uuid.MustParse("019a8f00-0002-7000-8000-000000000001"),
		DailyRateMinor: 2000000,
		Status:         domain.FleetStatusActive,
	}
	_, err := repo.CreateFleet(ctx, fleet)
	require.NoError(t, err)

	photos := []domain.AttachmentPhoto{
		{Angle: domain.AttachmentPhotoAngleOverall, StorageURI: "s3://media/overall.webp"},
		{Angle: domain.AttachmentPhotoAngleBracket, StorageURI: "s3://media/bracket.webp"},
		{Angle: domain.AttachmentPhotoAnglePins, StorageURI: "s3://media/pins.webp"},
		{Angle: domain.AttachmentPhotoAngleToolTip, StorageURI: "s3://media/tooltip.webp"},
	}

	att, err := service.RegisterAttachment(ctx, application.RegisterAttachmentInput{
		OwnerID:                  ownerID,
		AttachmentType:           domain.AttachmentTypeStandardBucket,
		Make:                     "CAT",
		Model:                    "Heavy Bucket",
		SerialNumber:             "SN-BUCKET-01",
		StandaloneDailyRateMinor: 500000, // 5,000 ETB
		Photos:                   photos,
	})
	require.NoError(t, err)

	// Attempting to link with 6,000 ETB (> 5,000 ETB) -> rejected with ErrBundledRateExceedsStandalone
	_, err = service.LinkAttachmentToFleet(ctx, ownerID, application.LinkAttachmentInput{
		FleetID:               fleet.ID,
		AttachmentID:          att.ID,
		BundledDailyRateMinor: 600000,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrBundledRateExceedsStandalone)
}
