package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// SearchFleets handles GET /api/v1/fleets/search (SRS §05.04, SR-01 to SR-05, SR-07)
func (s *Server) SearchFleets(w http.ResponseWriter, r *http.Request) {
	if s.discoveryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}

	q := r.URL.Query()
	lat, _ := strconv.ParseFloat(q.Get("latitude"), 64)
	lng, _ := strconv.ParseFloat(q.Get("longitude"), 64)

	filter := machineryDomain.FleetSearchFilter{
		Latitude:  lat,
		Longitude: lng,
	}

	if radStr := q.Get("radiusKm"); radStr != "" {
		if rad, err := strconv.ParseFloat(radStr, 64); err == nil {
			filter.RadiusKm = rad
		}
	}
	if catStr := q.Get("categoryId"); catStr != "" {
		if id, err := uuid.Parse(catStr); err == nil {
			filter.CategoryID = &id
		}
	}
	if subStr := q.Get("subcategoryId"); subStr != "" {
		if id, err := uuid.Parse(subStr); err == nil {
			filter.SubcategoryID = &id
		}
	}
	if startStr := q.Get("startDate"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			filter.StartDate = &t
		}
	}
	if endStr := q.Get("endDate"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			filter.EndDate = &t
		}
	}
	if minStr := q.Get("minRateMinor"); minStr != "" {
		if m, err := strconv.ParseInt(minStr, 10, 64); err == nil {
			filter.MinRateMinor = &m
		}
	}
	if maxStr := q.Get("maxRateMinor"); maxStr != "" {
		if m, err := strconv.ParseInt(maxStr, 10, 64); err == nil {
			filter.MaxRateMinor = &m
		}
	}
	if opStr := q.Get("operatorIncluded"); opStr != "" {
		if b, err := strconv.ParseBool(opStr); err == nil {
			filter.OperatorIncluded = &b
		}
	}
	if fuelStr := q.Get("fuelType"); fuelStr != "" {
		filter.FuelType = &fuelStr
	}
	if smStr := q.Get("selfMoveCapable"); smStr != "" {
		if b, err := strconv.ParseBool(smStr); err == nil {
			filter.SelfMoveCapable = &b
		}
	}
	if queryStr := q.Get("query"); queryStr != "" {
		filter.Query = &queryStr
	}
	if sortStr := q.Get("sortBy"); sortStr != "" {
		filter.SortBy = machineryDomain.SortBy(sortStr)
	}
	if cc := r.Header.Get("X-Country-Code"); cc != "" {
		filter.CountryCode = cc
	}
	if limStr := q.Get("limit"); limStr != "" {
		if l, err := strconv.Atoi(limStr); err == nil {
			filter.Limit = l
		}
	}

	out, err := s.discoveryService.SearchFleets(r.Context(), filter)
	if err != nil {
		if errors.Is(err, machineryDomain.ErrSpatialRadiusOutOfBounds) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Search Radius",
				Code:   "ERR_SPATIAL_RADIUS_OUT_OF_BOUNDS",
				Detail: "Search radius must be between 1.0 and 500.0 km.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrInvalidRateRange) || errors.Is(err, machineryDomain.ErrInvalidDateRange) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Filter Bounds",
				Code:   "ERR_INVALID_FILTER",
				Detail: err.Error(),
			})
			return
		}

		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to execute fleet search.",
		})
		return
	}

	items := make([]generated.FleetSearchResultItem, len(out.Results))
	for i, res := range out.Results {
		var specs generated.FleetSearchKeySpecs
		if res.KeySpecs.OperatingWeightKg != nil {
			specs.OperatingWeightKg = res.KeySpecs.OperatingWeightKg
		}
		if res.KeySpecs.PayloadCapacityKg != nil {
			specs.PayloadCapacityKg = res.KeySpecs.PayloadCapacityKg
		}
		if res.KeySpecs.GrossPowerHp != nil {
			specs.GrossPowerHp = res.KeySpecs.GrossPowerHp
		}

		subEn := res.SubcategoryNameEN
		subAm := res.SubcategoryNameAM
		excess := res.ExcessHourRateMinor

		subcatID := openapi_types.UUID(res.SubcategoryID)

		items[i] = generated.FleetSearchResultItem{
			FleetId:             openapi_types.UUID(res.FleetID),
			SubcategoryId:       &subcatID,
			SubcategoryNameEn:   &subEn,
			SubcategoryNameAm:   &subAm,
			Make:                res.Make,
			Model:               res.Model,
			ManufactureYear:     res.ManufactureYear,
			DailyRateMinor:      res.DailyRateMinor,
			ExcessHourRateMinor: &excess,
			DistanceMeters:      res.DistanceMeters,
			AvailabilityStatus:  generated.FleetSearchResultItemAvailabilityStatus(res.AvailabilityStatus),
			SelfMoveCapable:     res.SelfMoveCapable,
			OperatorIncluded:    res.OperatorIncluded,
			FuelType:            res.FuelType,
			PrimaryPhotoUri:     res.PrimaryPhotoURI,
			Owner: generated.FleetSearchOwnerSummary{
				OwnerId:               openapi_types.UUID(res.Owner.OwnerID),
				TradingName:           res.Owner.TradingName,
				AverageRating:         res.Owner.AverageRating,
				TotalCompletedRentals: res.Owner.TotalCompletedRentals,
				IsVerified:            res.Owner.IsVerified,
			},
			KeySpecs: &specs,
		}
	}

	respData := generated.FleetSearchData{
		TotalMatches: out.TotalMatches,
		SearchCenter: generated.FleetSearchCenter{
			Latitude:  out.Latitude,
			Longitude: out.Longitude,
			RadiusKm:  out.RadiusKm,
		},
		Results: items,
	}

	response.JSON(w, r, http.StatusOK, respData)
}
