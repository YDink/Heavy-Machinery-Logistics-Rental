package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// GetFleetAvailability handles GET /api/v1/fleets/{id}/availability (SRS §05.04, SR-07)
func (s *Server) GetFleetAvailability(w http.ResponseWriter, r *http.Request) {
	if s.discoveryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}

	idStr := r.PathValue("id")
	fleetID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Fleet ID",
			Code:   "ERR_INVALID_UUID",
			Detail: "The provided fleet ID is not a valid UUID.",
		})
		return
	}

	q := r.URL.Query()
	start, err := time.Parse("2006-01-02", q.Get("startDate"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Start Date",
			Code:   "ERR_INVALID_DATE_FORMAT",
			Detail: "Query parameter 'startDate' must be in YYYY-MM-DD format.",
		})
		return
	}
	end, err := time.Parse("2006-01-02", q.Get("endDate"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid End Date",
			Code:   "ERR_INVALID_DATE_FORMAT",
			Detail: "Query parameter 'endDate' must be in YYYY-MM-DD format.",
		})
		return
	}

	info, err := s.discoveryService.GetFleetAvailability(r.Context(), fleetID, start, end)
	if err != nil {
		if errors.Is(err, machineryDomain.ErrFleetNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Fleet Not Found",
				Code:   "ERR_FLEET_NOT_FOUND",
				Detail: "The requested fleet asset was not found.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrInvalidDateRange) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Date Range",
				Code:   "ERR_INVALID_DATE_RANGE",
				Detail: "Start date must be before or equal to end date.",
			})
			return
		}

		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to evaluate fleet availability.",
		})
		return
	}

	var blackouts []generated.FleetAvailabilityDateRange
	if len(info.BlackoutRanges) > 0 {
		blackouts = make([]generated.FleetAvailabilityDateRange, len(info.BlackoutRanges))
		for i, b := range info.BlackoutRanges {
			reason := b.Reason
			blackouts[i] = generated.FleetAvailabilityDateRange{
				StartDate: openapi_types.Date{Time: b.StartDate},
				EndDate:   openapi_types.Date{Time: b.EndDate},
				Reason:    &reason,
			}
		}
	}

	var reserved []generated.FleetAvailabilityDateRange
	if len(info.ReservedRanges) > 0 {
		reserved = make([]generated.FleetAvailabilityDateRange, len(info.ReservedRanges))
		for i, res := range info.ReservedRanges {
			state := res.State
			reserved[i] = generated.FleetAvailabilityDateRange{
				StartDate: openapi_types.Date{Time: res.StartDate},
				EndDate:   openapi_types.Date{Time: res.EndDate},
				State:     &state,
			}
		}
	}

	data := generated.FleetAvailabilityData{
		FleetId:              openapi_types.UUID(info.FleetID),
		IsAvailable:          info.IsAvailable,
		UnavailableReason:    info.UnavailableReason,
		HoldExpiresInSeconds: info.HoldExpiresInSeconds,
		BlackoutRanges:       &blackouts,
		ReservedRanges:       &reserved,
	}

	response.JSON(w, r, http.StatusOK, data)
}
