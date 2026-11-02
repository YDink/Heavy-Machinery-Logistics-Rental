package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
)

func TestValidateSearchFilter(t *testing.T) {
	t.Run("ValidFilter_Passes", func(t *testing.T) {
		start := time.Now()
		end := start.Add(48 * time.Hour)
		minRate := int64(100000)
		maxRate := int64(500000)

		filter := domain.FleetSearchFilter{
			Latitude:     9.0300,
			Longitude:    38.7400,
			RadiusKm:     50.0,
			StartDate:    &start,
			EndDate:      &end,
			MinRateMinor: &minRate,
			MaxRateMinor: &maxRate,
		}

		err := domain.ValidateSearchFilter(&filter)
		assert.NoError(t, err)
		assert.Equal(t, "ET", filter.CountryCode)
		assert.Equal(t, 20, filter.Limit)
	})

	t.Run("InvalidLatitude_ReturnsError", func(t *testing.T) {
		filter := domain.FleetSearchFilter{
			Latitude:  95.0,
			Longitude: 38.7400,
		}
		err := domain.ValidateSearchFilter(&filter)
		assert.Error(t, err)
	})

	t.Run("RadiusExceeds500Km_ReturnsErrSpatialRadiusOutOfBounds", func(t *testing.T) {
		filter := domain.FleetSearchFilter{
			Latitude:  9.0300,
			Longitude: 38.7400,
			RadiusKm:  600.0,
		}
		err := domain.ValidateSearchFilter(&filter)
		assert.ErrorIs(t, err, domain.ErrSpatialRadiusOutOfBounds)
	})

	t.Run("InvalidRateRange_ReturnsErrInvalidRateRange", func(t *testing.T) {
		minRate := int64(500000)
		maxRate := int64(100000)
		filter := domain.FleetSearchFilter{
			Latitude:     9.0300,
			Longitude:    38.7400,
			RadiusKm:     50.0,
			MinRateMinor: &minRate,
			MaxRateMinor: &maxRate,
		}
		err := domain.ValidateSearchFilter(&filter)
		assert.ErrorIs(t, err, domain.ErrInvalidRateRange)
	})

	t.Run("EndDateBeforeStartDate_ReturnsErrInvalidDateRange", func(t *testing.T) {
		start := time.Now()
		end := start.Add(-24 * time.Hour)
		filter := domain.FleetSearchFilter{
			Latitude:  9.0300,
			Longitude: 38.7400,
			RadiusKm:  50.0,
			StartDate: &start,
			EndDate:   &end,
		}
		err := domain.ValidateSearchFilter(&filter)
		assert.ErrorIs(t, err, domain.ErrInvalidDateRange)
	})
}

func TestValidateComparisonSet(t *testing.T) {
	subcat1 := uuid.New()
	subcat2 := uuid.New()

	fleet1 := domain.FleetAsset{ID: uuid.New(), SubcategoryID: subcat1, Make: "CAT", Model: "320D"}
	fleet2 := domain.FleetAsset{ID: uuid.New(), SubcategoryID: subcat1, Make: "Komatsu", Model: "PC200"}
	fleet3 := domain.FleetAsset{ID: uuid.New(), SubcategoryID: subcat1, Make: "Hitachi", Model: "ZX200"}
	fleetDiffSubcat := domain.FleetAsset{ID: uuid.New(), SubcategoryID: subcat2, Make: "Dynapac", Model: "CA250"}

	t.Run("ExactlyTwoOrThreeHomogeneousFleets_Passes", func(t *testing.T) {
		assert.NoError(t, domain.ValidateComparisonSet([]domain.FleetAsset{fleet1, fleet2}))
		assert.NoError(t, domain.ValidateComparisonSet([]domain.FleetAsset{fleet1, fleet2, fleet3}))
	})

	t.Run("SingleFleet_ReturnsErrInvalidComparisonCount", func(t *testing.T) {
		err := domain.ValidateComparisonSet([]domain.FleetAsset{fleet1})
		assert.ErrorIs(t, err, domain.ErrInvalidComparisonCount)
	})

	t.Run("FourFleets_ReturnsErrInvalidComparisonCount", func(t *testing.T) {
		fleet4 := domain.FleetAsset{ID: uuid.New(), SubcategoryID: subcat1}
		err := domain.ValidateComparisonSet([]domain.FleetAsset{fleet1, fleet2, fleet3, fleet4})
		assert.ErrorIs(t, err, domain.ErrInvalidComparisonCount)
	})

	t.Run("HeterogeneousCategories_ReturnsErrNonComparableCategories", func(t *testing.T) {
		err := domain.ValidateComparisonSet([]domain.FleetAsset{fleet1, fleetDiffSubcat})
		assert.ErrorIs(t, err, domain.ErrNonComparableCategories)
	})
}

func TestCalculateRecommendationScore(t *testing.T) {
	scoreHigh := domain.CalculateRecommendationScore(5.0, 5000, 50000, 2024, 1.0)
	scoreLow := domain.CalculateRecommendationScore(2.5, 45000, 50000, 2005, 0.6)

	assert.Greater(t, scoreHigh, scoreLow)
	assert.GreaterOrEqual(t, scoreHigh, 0.0)
	assert.LessOrEqual(t, scoreHigh, 1.0)
}

func TestBuildComparisonMatrix(t *testing.T) {
	subcatID := uuid.New()
	subcat := domain.Category{
		ID:     subcatID,
		NameEN: "Crawler Excavator",
		NameAM: "ክራውለር ኤክስካቫተር",
	}

	fleet1 := domain.FleetAsset{
		ID:             uuid.New(),
		SubcategoryID:  subcatID,
		Make:           "Caterpillar",
		Model:          "320D",
		DailyRateMinor: 2250000,
	}
	fleet2 := domain.FleetAsset{
		ID:             uuid.New(),
		SubcategoryID:  subcatID,
		Make:           "Komatsu",
		Model:          "PC200-8",
		DailyRateMinor: 2100000,
	}

	specWeightID := uuid.New()
	specPowerID := uuid.New()
	specInternalID := uuid.New()

	unitKg := "kg"
	unitHp := "hp"

	defs := []domain.SpecFieldDefinition{
		{
			ID:           specWeightID,
			FieldKey:     "operating_weight_kg",
			LabelEN:      "Operating Weight",
			LabelAM:      "የስራ ክብደት",
			UnitSymbol:   &unitKg,
			FieldType:    "DECIMAL",
			IsComparable: true,
			DisplayOrder: 1,
		},
		{
			ID:           specPowerID,
			FieldKey:     "gross_power_hp",
			LabelEN:      "Gross Engine Power",
			LabelAM:      "ጠቅላላ የሞተር ጉልበት",
			UnitSymbol:   &unitHp,
			FieldType:    "DECIMAL",
			IsComparable: true,
			DisplayOrder: 2,
		},
		{
			ID:           specInternalID,
			FieldKey:     "internal_serial_code",
			LabelEN:      "Internal Code",
			IsComparable: false, // Not comparable, should be omitted
			DisplayOrder: 3,
		},
	}

	w1 := 21800.0
	w2 := 20010.0
	p1 := 148.0
	p2 := 148.0

	values := []domain.FleetSpecValue{
		{FleetID: fleet1.ID, SpecID: specWeightID, ValueNumeric: &w1},
		{FleetID: fleet2.ID, SpecID: specWeightID, ValueNumeric: &w2},
		{FleetID: fleet1.ID, SpecID: specPowerID, ValueNumeric: &p1},
		{FleetID: fleet2.ID, SpecID: specPowerID, ValueNumeric: &p2},
	}

	result := domain.BuildComparisonMatrix([]domain.FleetAsset{fleet1, fleet2}, subcat, defs, values)
	require.NotNil(t, result)
	assert.Equal(t, 2, len(result.Fleets))
	assert.Equal(t, 2, len(result.ComparisonMatrix)) // Non-comparable omitted

	// Row 1: operating weight differs
	assert.Equal(t, "operating_weight_kg", result.ComparisonMatrix[0].FieldKey)
	assert.True(t, result.ComparisonMatrix[0].IsDifferent)
	assert.Equal(t, "21800", result.ComparisonMatrix[0].Values[fleet1.ID.String()])
	assert.Equal(t, "20010", result.ComparisonMatrix[0].Values[fleet2.ID.String()])

	// Row 2: power is identical
	assert.Equal(t, "gross_power_hp", result.ComparisonMatrix[1].FieldKey)
	assert.False(t, result.ComparisonMatrix[1].IsDifferent)
}

func TestSortSearchResults(t *testing.T) {
	item1 := domain.FleetSearchResult{
		FleetID:             uuid.New(),
		DailyRateMinor:      200000,
		DistanceMeters:      10000,
		RecommendationScore: 0.85,
		Owner:               domain.FleetSearchOwner{AverageRating: 4.8},
	}
	item2 := domain.FleetSearchResult{
		FleetID:             uuid.New(),
		DailyRateMinor:      150000,
		DistanceMeters:      25000,
		RecommendationScore: 0.90,
		Owner:               domain.FleetSearchOwner{AverageRating: 4.5},
	}

	results := []domain.FleetSearchResult{item1, item2}

	domain.SortSearchResults(results, domain.SortByPriceAsc)
	assert.Equal(t, item2.FleetID, results[0].FleetID)

	domain.SortSearchResults(results, domain.SortByDistanceAsc)
	assert.Equal(t, item1.FleetID, results[0].FleetID)

	domain.SortSearchResults(results, domain.SortByRecommended)
	assert.Equal(t, item2.FleetID, results[0].FleetID)

	domain.SortSearchResults(results, domain.SortByRatingDesc)
	assert.Equal(t, item1.FleetID, results[0].FleetID)
}
