package http

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
)

// GetBookingById handles GET /api/v1/bookings/{id}
func (s *Server) GetBookingById(w http.ResponseWriter, r *http.Request) {
	if s.bookingService == nil {
		s.notImplemented(w, r, "05-05")
		return
	}

	idStr := r.PathValue("id")
	bookingUUID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Booking ID",
			Code:   "ERR_INVALID_UUID",
			Detail: "The provided booking ID path parameter is not a valid UUID.",
		})
		return
	}

	b, err := s.bookingService.GetBooking(r.Context(), bookingUUID)
	if err != nil {
		if errors.Is(err, bookingDomain.ErrBookingNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Booking Not Found",
				Code:   "ERR_BOOKING_NOT_FOUND",
				Detail: "The requested booking reservation does not exist.",
			})
			return
		}
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: err.Error(),
		})
		return
	}

	response.JSON(w, r, http.StatusOK, toGeneratedBooking(b))
}
