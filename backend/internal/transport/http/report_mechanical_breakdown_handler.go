package http

import (
	"encoding/json"
	"errors"
	"net/http"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	rentalDomain "github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// ReportMechanicalBreakdown handles POST /api/v1/bookings/{id}/breakdown
func (s *Server) ReportMechanicalBreakdown(w http.ResponseWriter, r *http.Request) {
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

	var req generated.MechanicalBreakdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BODY",
			Detail: "Malformed request payload.",
		})
		return
	}

	dtcs := []string{}
	if req.DiagnosticTroubleCodes != nil {
		dtcs = *req.DiagnosticTroubleCodes
	}

	photos := []string{}
	if req.IncidentPhotos != nil {
		photos = *req.IncidentPhotos
	}

	breakdown, err := s.rentalService.ReportMechanicalBreakdown(r.Context(), rentalApp.BreakdownReportCommand{
		BookingID:              bookingID,
		ReporterID:             userUUID,
		FaultCategory:          string(req.FaultCategory),
		FaultDescription:       req.FaultDescription,
		DiagnosticTroubleCodes: dtcs,
		CurrentHourMeter:       req.CurrentHourMeter,
		IncidentPhotos:         photos,
	})
	if err != nil {
		if errors.Is(err, rentalDomain.ErrBookingAlreadyPaused) {
			response.Error(w, r, http.StatusConflict, response.ProblemDetails{
				Title:  "Booking Already Paused",
				Code:   "ERR_BOOKING_ALREADY_PAUSED",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, rentalDomain.ErrInvalidFaultCategory) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Fault Category",
				Code:   "ERR_INVALID_FAULT_CATEGORY",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Breakdown Submission Failed",
			Code:   "ERR_BREAKDOWN_FAILED",
			Detail: err.Error(),
		})
		return
	}

	data := generated.MechanicalBreakdownResponseData{
		BookingId:          breakdown.BookingID,
		IncidentId:         breakdown.ID,
		OwnerSlaDeadline:   breakdown.SLADeadline,
		ShiftTimerPausedAt: breakdown.CreatedAt,
		Status:             breakdown.Status,
	}

	response.JSON(w, r, http.StatusOK, data)
}
