package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	rentalDomain "github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// ResolveMechanicalBreakdown handles POST /api/v1/bookings/{id}/breakdown/resolve
func (s *Server) ResolveMechanicalBreakdown(w http.ResponseWriter, r *http.Request) {
	userUUID, err := s.authenticateUser(r)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: err.Error(),
		})
		return
	}

	bookingID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BOOKING_ID",
			Detail: "Invalid booking ID in path.",
		})
		return
	}

	var req generated.BreakdownResolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BODY",
			Detail: "Malformed request payload.",
		})
		return
	}

	breakdown, err := s.rentalService.ResolveMechanicalBreakdown(r.Context(), rentalApp.BreakdownResolveCommand{
		BookingID:       bookingID,
		ResolverID:      userUUID,
		MechanicName:    req.MechanicName,
		MechanicPhone:   req.MechanicPhone,
		ResolutionNotes: req.ResolutionNotes,
	})
	if err != nil {
		if errors.Is(err, rentalDomain.ErrBreakdownNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Breakdown Not Found",
				Code:   "ERR_BREAKDOWN_NOT_FOUND",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: "Only equipment owner can resolve breakdown.",
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Resolution Failed",
			Code:   "ERR_RESOLUTION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	resTime := time.Now().UTC()
	if breakdown.ResolvedAt != nil {
		resTime = *breakdown.ResolvedAt
	}

	data := generated.BreakdownResolutionResponseData{
		BookingId:  breakdown.BookingID,
		IncidentId: breakdown.ID,
		ResumedAt:  resTime,
		Status:     breakdown.Status,
	}

	response.JSON(w, r, http.StatusOK, data)
}
