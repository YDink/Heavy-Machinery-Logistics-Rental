package http

import (
	"net/http"
	"time"

	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// InitiateReturn handles POST /api/v1/bookings/{id}/return-initiate
func (s *Server) InitiateReturn(w http.ResponseWriter, r *http.Request) {
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

	booking, err := s.rentalService.InitiateReturn(r.Context(), bookingID, userUUID)
	if err != nil {
		response.Error(w, r, http.StatusConflict, response.ProblemDetails{
			Title:  "Conflict",
			Code:   "ERR_RETURN_INITIATION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	checkIn := 0.0
	if booking.CheckInMeterHours != nil {
		checkIn = *booking.CheckInMeterHours
	}

	resp := generated.ReturnInitiateResponseData{
		BookingId:         booking.ID,
		Status:            string(booking.Status),
		CheckInMeterHours: checkIn,
		AllowedShiftHours: float64(booking.DurationDays) * 8.0,
		InitiatedAt:       time.Now().UTC(),
	}

	response.JSON(w, r, http.StatusOK, resp)
}
