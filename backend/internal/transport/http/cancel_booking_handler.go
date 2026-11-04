package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	bookingApp "github.com/dinkmachinery/dink-core/internal/modules/booking/application"
	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// CancelBooking handles POST /api/v1/bookings/{id}/cancel
func (s *Server) CancelBooking(w http.ResponseWriter, r *http.Request) {
	if s.bookingService == nil {
		s.notImplemented(w, r, "05-05")
		return
	}

	authHeader := r.Header.Get("Authorization")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Access token is invalid.",
		})
		return
	}

	userUUID, _ := uuid.Parse(claims.Subject)
	idStr := r.PathValue("id")
	bookingUUID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Booking ID",
			Code:   "ERR_INVALID_UUID",
			Detail: "Invalid booking UUID format.",
		})
		return
	}

	var body generated.CancelBookingJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Malformed JSON body.",
		})
		return
	}

	b, err := s.bookingService.CancelBooking(r.Context(), bookingApp.CancelBookingInput{
		BookingID:       bookingUUID,
		ActorID:         userUUID,
		ExpectedVersion: body.ExpectedVersion,
		Reason:          body.Reason,
		IsForceMajeure:  false,
	})
	if err != nil {
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_UNAUTHORIZED_ACTION",
				Detail: "You are not a recognized party to this booking.",
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrConcurrencyConflict) {
			response.Error(w, r, http.StatusConflict, response.ProblemDetails{
				Title:  "Concurrency Conflict",
				Code:   "BOOKING_CONCURRENCY_CONFLICT",
				Detail: "The booking was modified by another transaction.",
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrInvalidStateTransition) {
			response.Error(w, r, http.StatusConflict, response.ProblemDetails{
				Title:  "Invalid State Transition",
				Code:   "ERR_INVALID_STATE_TRANSITION",
				Detail: err.Error(),
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
