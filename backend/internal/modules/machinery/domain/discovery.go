package domain

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
)

type AvailabilityStatus string

const (
	AvailabilityStatusAvailableNow     AvailabilityStatus = "AVAILABLE_NOW"
	AvailabilityStatusHeldPendingWire  AvailabilityStatus = "HELD_PENDING_WIRE"
	AvailabilityStatusUnavailable      AvailabilityStatus = "UNAVAILABLE"
)

type SortBy string

const (
	SortByRecommended SortBy = "RECOMMENDED"
	SortByPriceAsc    SortBy = "PRICE_ASC"
	SortByPriceDesc   SortBy = "PRICE_DESC"
	SortByDistanceAsc SortBy = "DISTANCE_ASC"
	SortByRatingDesc  SortBy = "RATING_DESC"
)

type FleetSearchFilter struct {
	Latitude         float64
	Longitude        float64
	RadiusKm         float64
	CategoryID       *uuid.UUID
	SubcategoryID    *uuid.UUID
	StartDate        *time.Time
	EndDate          *time.Time
	MinRateMinor     *int64
	MaxRateMinor     *int64
	OperatorIncluded *bool
	FuelType         *string
	SelfMoveCapable  *bool
	Query            *string
	SortBy           SortBy
	CountryCode      string
	Limit            int
	Offset           int
}

type FleetSearchOwner struct {
	OwnerID               uuid.UUID
	TradingName           string
	AverageRating         float64
	TotalCompletedRentals int
	IsVerified            bool
}

type FleetSearchKeySpecs struct {
	OperatingWeightKg *float64
	PayloadCapacityKg *float64
	GrossPowerHp      *float64
}

type FleetSearchResult struct {
	FleetID             uuid.UUID
	SubcategoryID       uuid.UUID
	SubcategoryNameEN   string
	SubcategoryNameAM   string
	Make                string
	Model               string
	ManufactureYear     int
	DailyRateMinor      int64
	ExcessHourRateMinor int64
	DistanceMeters      float64
	AvailabilityStatus  AvailabilityStatus
	SelfMoveCapable     bool
	OperatorIncluded    bool
	FuelType            string
	PrimaryPhotoURI     *string
	Owner               FleetSearchOwner
	KeySpecs            FleetSearchKeySpecs
	RecommendationScore float64
}

type FleetCompareHeader struct {
	ID                  uuid.UUID
	Title               string
	DailyRateMinor      int64
	ExcessHourRateMinor int64
	PrimaryPhotoURI     *string
	SelfMoveCapable     bool
	OperatorIncluded    bool
	FuelType            string
}

type FleetCompareSpecRow struct {
	FieldKey    string
	LabelEN     string
	LabelAM     string
	UnitSymbol  *string
	FieldType   string
	Values      map[string]string // map[fleetID.String()]valueString
	IsDifferent bool
}

type FleetCompareResult struct {
	SubcategoryID     uuid.UUID
	SubcategoryNameEN string
	SubcategoryNameAM string
	Fleets            []FleetCompareHeader
	ComparisonMatrix  []FleetCompareSpecRow
}

type DateRangeReason struct {
	StartDate time.Time
	EndDate   time.Time
	Reason    string
}

type DateRangeState struct {
	StartDate time.Time
	EndDate   time.Time
	State     string
}

type FleetAvailabilityInfo struct {
	FleetID              uuid.UUID
	IsAvailable          bool
	UnavailableReason    *string
	HoldExpiresInSeconds *int
	NextAvailableDate    *time.Time
	BlackoutRanges       []DateRangeReason
	ReservedRanges       []DateRangeState
}

type OwnerAggregate struct {
	OwnerID               uuid.UUID
	AverageRating         float64
	TotalRentalsCount     int
	TotalCompletedRentals int
	CompletionRate        float64
	IsVerified            bool
	UpdatedAt             time.Time
}

type FleetInsurance struct {
	ID             uuid.UUID
	FleetID        uuid.UUID
	PolicyNumber   string
	ProviderName   string
	CoverageType   string
	Status         string
	ExpiresAt      time.Time
	CertificateURI string
	UpdatedAt      time.Time
}

// ValidateSearchFilter enforces spatial bounds and logical filter consistency.
func ValidateSearchFilter(filter *FleetSearchFilter) error {
	if filter.Latitude < -90.0 || filter.Latitude > 90.0 {
		return fmt.Errorf("latitude must be between -90 and +90 degrees")
	}
	if filter.Longitude < -180.0 || filter.Longitude > 180.0 {
		return fmt.Errorf("longitude must be between -180 and +180 degrees")
	}
	if filter.RadiusKm <= 0.0 {
		filter.RadiusKm = 50.0 // Default 50km
	}
	if filter.RadiusKm < 1.0 || filter.RadiusKm > 500.0 {
		return ErrSpatialRadiusOutOfBounds
	}
	if filter.MinRateMinor != nil && filter.MaxRateMinor != nil {
		if *filter.MinRateMinor > *filter.MaxRateMinor {
			return ErrInvalidRateRange
		}
	}
	if filter.StartDate != nil && filter.EndDate != nil {
		if filter.EndDate.Before(*filter.StartDate) {
			return ErrInvalidDateRange
		}
	}
	if filter.CountryCode == "" {
		filter.CountryCode = "ET"
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	return nil
}

// ValidateComparisonSet asserts that between 2 and 3 fleets of the identical subcategory are compared (SR-08).
func ValidateComparisonSet(fleets []FleetAsset) error {
	if len(fleets) < 2 || len(fleets) > 3 {
		return ErrInvalidComparisonCount
	}

	subcat := fleets[0].SubcategoryID
	for i := 1; i < len(fleets); i++ {
		if fleets[i].SubcategoryID != subcat {
			return ErrNonComparableCategories
		}
	}
	return nil
}

// CalculateRecommendationScore implements the multi-factor scoring formula (SRS §3.4, SR-04):
// Score = (0.40 * R_owner) + (0.30 * P_proximity) + (0.20 * U_fleet) + (0.10 * C_completion)
func CalculateRecommendationScore(
	averageRating float64,
	distanceMeters float64,
	radiusMeters float64,
	manufactureYear int,
	completionRate float64,
) float64 {
	// Normalized owner review score (0.0 to 1.0)
	rOwner := math.Min(math.Max(averageRating/5.0, 0.0), 1.0)

	// Inverse distance scalar (1.0 - (Distance / Search Radius))
	pProximity := 0.0
	if radiusMeters > 0 {
		pProximity = math.Max(1.0-(distanceMeters/radiusMeters), 0.0)
	}

	// Asset freshness factor (brand new listings near 1.0, degrading over 20 years)
	currentYear := time.Now().Year()
	fleetAgeYears := currentYear - manufactureYear
	if fleetAgeYears < 0 {
		fleetAgeYears = 0
	}
	uFleet := math.Max(1.0-(float64(fleetAgeYears)/20.0), 0.1)

	// Historical booking fulfillment rate (0.0 to 1.0)
	cCompletion := math.Min(math.Max(completionRate, 0.0), 1.0)

	score := (0.40 * rOwner) + (0.30 * pProximity) + (0.20 * uFleet) + (0.10 * cCompletion)
	return math.Round(score*10000) / 10000
}

// BuildComparisonMatrix transposes dynamic EAV specifications for up to 3 machines (SR-08, FL-15).
func BuildComparisonMatrix(
	fleets []FleetAsset,
	subcat Category,
	specDefs []SpecFieldDefinition,
	specValues []FleetSpecValue,
) *FleetCompareResult {
	headers := make([]FleetCompareHeader, len(fleets))
	fleetIDMap := make(map[uuid.UUID]bool)
	for i, f := range fleets {
		title := fmt.Sprintf("%s %s", f.Make, f.Model)
		var photoURI *string
		if len(f.Photos) > 0 {
			photoURI = &f.Photos[0].StorageURI
		}
		headers[i] = FleetCompareHeader{
			ID:                  f.ID,
			Title:               title,
			DailyRateMinor:      f.DailyRateMinor,
			ExcessHourRateMinor: f.ExcessHourRateMinor,
			PrimaryPhotoURI:     photoURI,
			SelfMoveCapable:     f.SelfMoveCapable,
			OperatorIncluded:    f.OperatorIncluded,
			FuelType:            string(f.FuelType),
		}
		fleetIDMap[f.ID] = true
	}

	// Group values by Spec Definition ID -> Fleet ID -> String representation
	valMap := make(map[uuid.UUID]map[uuid.UUID]string)
	for _, v := range specValues {
		if !fleetIDMap[v.FleetID] {
			continue
		}
		if _, ok := valMap[v.SpecID]; !ok {
			valMap[v.SpecID] = make(map[uuid.UUID]string)
		}
		strVal := ""
		if v.ValueString != nil {
			strVal = *v.ValueString
		} else if v.ValueNumeric != nil {
			strVal = fmt.Sprintf("%v", *v.ValueNumeric)
		} else if v.ValueBoolean != nil {
			strVal = fmt.Sprintf("%v", *v.ValueBoolean)
		}
		valMap[v.SpecID][v.FleetID] = strVal
	}

	var rows []FleetCompareSpecRow
	for _, def := range specDefs {
		if !def.IsComparable {
			continue
		}
		values := make(map[string]string)
		firstVal := ""
		firstSet := false
		isDifferent := false

		for _, f := range fleets {
			v := ""
			if m, ok := valMap[def.ID]; ok {
				v = m[f.ID]
			}
			values[f.ID.String()] = v

			if !firstSet {
				firstVal = v
				firstSet = true
			} else if v != firstVal {
				isDifferent = true
			}
		}

		rows = append(rows, FleetCompareSpecRow{
			FieldKey:    def.FieldKey,
			LabelEN:     def.LabelEN,
			LabelAM:     def.LabelAM,
			UnitSymbol:  def.UnitSymbol,
			FieldType:   def.FieldType,
			Values:      values,
			IsDifferent: isDifferent,
		})
	}

	return &FleetCompareResult{
		SubcategoryID:     subcat.ID,
		SubcategoryNameEN: subcat.NameEN,
		SubcategoryNameAM: subcat.NameAM,
		Fleets:            headers,
		ComparisonMatrix:  rows,
	}
}

// SortSearchResults applies deterministic sorting over candidate discovery results (SR-04).
func SortSearchResults(results []FleetSearchResult, sortBy SortBy) {
	switch sortBy {
	case SortByPriceAsc:
		sort.Slice(results, func(i, j int) bool {
			if results[i].DailyRateMinor == results[j].DailyRateMinor {
				return results[i].DistanceMeters < results[j].DistanceMeters
			}
			return results[i].DailyRateMinor < results[j].DailyRateMinor
		})
	case SortByPriceDesc:
		sort.Slice(results, func(i, j int) bool {
			if results[i].DailyRateMinor == results[j].DailyRateMinor {
				return results[i].DistanceMeters < results[j].DistanceMeters
			}
			return results[i].DailyRateMinor > results[j].DailyRateMinor
		})
	case SortByDistanceAsc:
		sort.Slice(results, func(i, j int) bool {
			return results[i].DistanceMeters < results[j].DistanceMeters
		})
	case SortByRatingDesc:
		sort.Slice(results, func(i, j int) bool {
			if results[i].Owner.AverageRating == results[j].Owner.AverageRating {
				return results[i].Owner.TotalCompletedRentals > results[j].Owner.TotalCompletedRentals
			}
			return results[i].Owner.AverageRating > results[j].Owner.AverageRating
		})
	case SortByRecommended:
		fallthrough
	default:
		sort.Slice(results, func(i, j int) bool {
			return results[i].RecommendationScore > results[j].RecommendationScore
		})
	}
}
