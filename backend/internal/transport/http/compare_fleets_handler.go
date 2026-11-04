package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// CompareFleets handles POST /api/v1/fleets/compare (SRS §05.04, SR-08, FL-15)
func (s *Server) CompareFleets(w http.ResponseWriter, r *http.Request) {
	if s.discoveryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}

	var req generated.FleetCompareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	if len(req.FleetIds) < 2 || len(req.FleetIds) > 3 {
		response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
			Title:  "Invalid Comparison Fleet Count",
			Code:   "ERR_INVALID_COMPARISON_COUNT",
			Detail: "Comparison requires between 2 and 3 fleets.",
		})
		return
	}

	ids := make([]uuid.UUID, len(req.FleetIds))
	for i, id := range req.FleetIds {
		ids[i] = uuid.UUID(id)
	}

	matrix, err := s.discoveryService.CompareFleets(r.Context(), ids)
	if err != nil {
		if errors.Is(err, machineryDomain.ErrNonComparableCategories) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Taxonomic Category Mismatch in Comparison",
				Code:   "ERR_NON_COMPARABLE_CATEGORIES",
				Detail: "Comparison matrix mandates identical subcategories across all compared assets.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrInvalidComparisonCount) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Fleet Count",
				Code:   "ERR_INVALID_COMPARISON_COUNT",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrFleetNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Fleet Not Found",
				Code:   "ERR_FLEET_NOT_FOUND",
				Detail: "One or more fleet assets selected for comparison were not found.",
			})
			return
		}

		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to build comparison matrix.",
		})
		return
	}

	headers := make([]generated.FleetCompareItemHeader, len(matrix.Fleets))
	for i, f := range matrix.Fleets {
		excess := f.ExcessHourRateMinor
		fuel := f.FuelType
		headers[i] = generated.FleetCompareItemHeader{
			Id:                  openapi_types.UUID(f.ID),
			Title:               f.Title,
			DailyRateMinor:      f.DailyRateMinor,
			ExcessHourRateMinor: &excess,
			PrimaryPhotoUri:     f.PrimaryPhotoURI,
			SelfMoveCapable:     f.SelfMoveCapable,
			OperatorIncluded:    f.OperatorIncluded,
			FuelType:            &fuel,
		}
	}

	rows := make([]generated.FleetCompareSpecRow, len(matrix.ComparisonMatrix))
	for i, row := range matrix.ComparisonMatrix {
		fieldKey := row.FieldKey
		fieldType := row.FieldType
		rows[i] = generated.FleetCompareSpecRow{
			FieldKey:    fieldKey,
			LabelEn:     row.LabelEN,
			LabelAm:     row.LabelAM,
			UnitSymbol:  row.UnitSymbol,
			FieldType:   &fieldType,
			Values:      row.Values,
			IsDifferent: row.IsDifferent,
		}
	}

	data := generated.FleetCompareData{
		SubcategoryId:     openapi_types.UUID(matrix.SubcategoryID),
		SubcategoryNameEn: matrix.SubcategoryNameEN,
		SubcategoryNameAm: matrix.SubcategoryNameAM,
		Fleets:            headers,
		ComparisonMatrix:  rows,
	}

	response.JSON(w, r, http.StatusOK, data)
}
