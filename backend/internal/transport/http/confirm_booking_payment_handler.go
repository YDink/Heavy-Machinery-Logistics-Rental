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

// ConfirmBookingPayment handles POST /api/v1/bookings/{id}/confirm-payment
func (s *Server) ConfirmBookingPayment(w http.ResponseWriter, r *http.Request) {
	if s.bookingService == nil {
		s.notImplemented(w, r, "05-05")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing Authorization header.",
		})
		return
	}

	claims, err := s.tokenService.VerifyAccessToken(r.Context(), strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_INVALID_TOKEN",
			Detail: "Access token is invalid or expired.",
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

	var body generated.ConfirmBookingPaymentJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Malformed JSON body.",
		})
		return
	}

	b, err := s.bookingService.ConfirmPayment(r.Context(), bookingApp.ConfirmPaymentInput{
		BookingID:       bookingUUID,
		ActorID:         userUUID,
		ExpectedVersion: body.ExpectedVersion,
		PaymentRef:      body.PaymentReference,
		Notes:           "",
	})
	if err != nil {
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_UNAUTHORIZED_ACTION",
				Detail: "Only the fleet owner can confirm receipt of payment.",
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrConcurrencyConflict) {
			response.Error(w, r, http.StatusConflict, response.ProblemDetails{
				Title:  "Concurrency Conflict",
				Code:   "BOOKING_CONCURRENCY_CONFLICT",
				Detail: "The booking was modified by another transaction. Reload and retry.",
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
