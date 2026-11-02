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

func TestDiscoveryService_SearchFleets(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryMachineryRepository()
	leaseChecker := application.NewMemoryHoldLeaseChecker()
	service := application.NewDiscoveryService(repo, repo, leaseChecker)

	ownerID := uuid.New()
	subcatID := uuid.MustParse("019a8f00-0002-7000-8000-000000000001")

	// Machine 1: Bole, Addis Ababa (~0km from 9.0300, 38.7400)
	f1, err := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:             ownerID,
		SubcategoryID:       subcatID,
		Make:                "Caterpillar",
		Model:               "320D",
		ManufactureYear:     2024,
		PlateNumber:         "ET-3-11111",
		DailyRateMinor:      2250000,
		ExcessHourRateMinor: 280000,
		Status:              domain.FleetStatusActive,
		Latitude:            9.0300,
		Longitude:           38.7400,
		AddressText:         "Bole, Addis Ababa",
		FuelType:            domain.FuelTypeDiesel,
		SelfMoveCapable:     false,
		OperatorIncluded:    true,
	})
	require.NoError(t, err)

	// Machine 2: Adama (~75km southeast from Addis Ababa)
	f2, err := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:             ownerID,
		SubcategoryID:       subcatID,
		Make:                "Komatsu",
		Model:               "PC200-8",
		ManufactureYear:     2022,
		PlateNumber:         "ET-3-22222",
		DailyRateMinor:      2000000,
		ExcessHourRateMinor: 250000,
		Status:              domain.FleetStatusActive,
		Latitude:            8.5400,
		Longitude:           39.2700,
		AddressText:         "Adama Depot",
		FuelType:            domain.FuelTypeDiesel,
		SelfMoveCapable:     false,
		OperatorIncluded:    false,
	})
	require.NoError(t, err)

	// Machine 3: Inactive machine in Addis Ababa (Must be excluded)
	_, err = repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:             ownerID,
		SubcategoryID:       subcatID,
		Make:                "Hitachi",
		Model:               "ZX200",
		ManufactureYear:     2023,
		PlateNumber:         "ET-3-33333",
		DailyRateMinor:      1900000,
		ExcessHourRateMinor: 240000,
		Status:              domain.FleetStatusInactive,
		Latitude:            9.0200,
		Longitude:           38.7500,
		AddressText:         "Kaliti, Addis Ababa",
		FuelType:            domain.FuelTypeDiesel,
	})
	require.NoError(t, err)

	t.Run("SpatialRadius_ExcludesFarMachines", func(t *testing.T) {
		filter := domain.FleetSearchFilter{
			Latitude:  9.0300,
			Longitude: 38.7400,
			RadiusKm:  50.0, // Adama (~75km) must be excluded
			SortBy:    domain.SortByDistanceAsc,
		}

		res, err := service.SearchFleets(ctx, filter)
		require.NoError(t, err)
		assert.Equal(t, 1, res.TotalMatches)
		assert.Equal(t, f1.ID, res.Results[0].FleetID)
		assert.Less(t, res.Results[0].DistanceMeters, 50000.0)
	})

	t.Run("WiderRadius_IncludesBothActiveMachines", func(t *testing.T) {
		filter := domain.FleetSearchFilter{
			Latitude:  9.0300,
			Longitude: 38.7400,
			RadiusKm:  100.0, // Includes both Addis Ababa and Adama
			SortBy:    domain.SortByDistanceAsc,
		}

		res, err := service.SearchFleets(ctx, filter)
		require.NoError(t, err)
		assert.Equal(t, 2, res.TotalMatches)
		assert.Equal(t, f1.ID, res.Results[0].FleetID)
		assert.Equal(t, f2.ID, res.Results[1].FleetID)
	})

	t.Run("ActiveHoldLease_DemotesAvailabilityBadge", func(t *testing.T) {
		// Put f1 on a 60-minute Valkey lease
		leaseChecker.SetHold(f1.ID, 60*time.Minute)
		defer leaseChecker.ReleaseHold(f1.ID)

		filter := domain.FleetSearchFilter{
			Latitude:  9.0300,
			Longitude: 38.7400,
			RadiusKm:  50.0,
		}

		res, err := service.SearchFleets(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 1, len(res.Results))
		assert.Equal(t, domain.AvailabilityStatusHeldPendingWire, res.Results[0].AvailabilityStatus)
	})

	t.Run("TemporalBlackout_ExcludesBookedMachine", func(t *testing.T) {
		start := time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC)

		repo.SeedBlackout(f1.ID, start, end, "Scheduled Track Replacement")

		searchStart := time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC)
		searchEnd := time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC)

		filter := domain.FleetSearchFilter{
			Latitude:  9.0300,
			Longitude: 38.7400,
			RadiusKm:  50.0,
			StartDate: &searchStart,
			EndDate:   &searchEnd,
		}

		res, err := service.SearchFleets(ctx, filter)
		require.NoError(t, err)
		// f1 is in maintenance during search window, so 0 matches
		assert.Equal(t, 0, res.TotalMatches)
	})
}

func TestDiscoveryService_CompareFleets(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryMachineryRepository()
	service := application.NewDiscoveryService(repo, repo, nil)

	ownerID := uuid.New()
	subcat1 := uuid.MustParse("019a8f00-0002-7000-8000-000000000001") // Mini Excavator
	subcat2 := uuid.New()                                              // Lowboy

	f1, _ := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:        ownerID,
		SubcategoryID:  subcat1,
		Make:           "CAT",
		Model:          "305E2",
		DailyRateMinor: 1500000,
	})
	f2, _ := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:        ownerID,
		SubcategoryID:  subcat1,
		Make:           "Komatsu",
		Model:          "PC55MR",
		DailyRateMinor: 1400000,
	})
	fDiff, _ := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:        ownerID,
		SubcategoryID:  subcat2,
		Make:           "Lowboy",
		Model:          "Tri-Axle",
		DailyRateMinor: 3500000,
	})

	t.Run("ValidComparison_ReturnsDynamicMatrix", func(t *testing.T) {
		matrix, err := service.CompareFleets(ctx, []uuid.UUID{f1.ID, f2.ID})
		require.NoError(t, err)
		assert.Equal(t, subcat1, matrix.SubcategoryID)
		assert.Equal(t, 2, len(matrix.Fleets))
		assert.NotEmpty(t, matrix.ComparisonMatrix)
	})

	t.Run("CrossCategory_ReturnsErrNonComparableCategories", func(t *testing.T) {
		_, err := service.CompareFleets(ctx, []uuid.UUID{f1.ID, fDiff.ID})
		assert.ErrorIs(t, err, domain.ErrNonComparableCategories)
	})

	t.Run("InvalidCount_ReturnsErrInvalidComparisonCount", func(t *testing.T) {
		_, err := service.CompareFleets(ctx, []uuid.UUID{f1.ID})
		assert.ErrorIs(t, err, domain.ErrInvalidComparisonCount)
	})
}

func TestDiscoveryService_GetFleetAvailability(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryMachineryRepository()
	leaseChecker := application.NewMemoryHoldLeaseChecker()
	service := application.NewDiscoveryService(repo, repo, leaseChecker)

	f, _ := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:       uuid.New(),
		SubcategoryID: uuid.MustParse("019a8f00-0002-7000-8000-000000000001"),
		Make:          "CAT",
		Model:         "320D",
	})

	start := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)

	t.Run("FullyAvailableAsset_ReturnsTrue", func(t *testing.T) {
		avail, err := service.GetFleetAvailability(ctx, f.ID, start, end)
		require.NoError(t, err)
		assert.True(t, avail.IsAvailable)
		assert.Nil(t, avail.UnavailableReason)
	})

	t.Run("AssetWithHoldLease_ReturnsUnavailableWithTTL", func(t *testing.T) {
		leaseChecker.SetHold(f.ID, 30*time.Minute)
		defer leaseChecker.ReleaseHold(f.ID)

		avail, err := service.GetFleetAvailability(ctx, f.ID, start, end)
		require.NoError(t, err)
		assert.False(t, avail.IsAvailable)
		require.NotNil(t, avail.UnavailableReason)
		assert.Equal(t, "HELD_BY_ANOTHER_RENTER", *avail.UnavailableReason)
		require.NotNil(t, avail.HoldExpiresInSeconds)
		assert.Greater(t, *avail.HoldExpiresInSeconds, 0)
	})
}
