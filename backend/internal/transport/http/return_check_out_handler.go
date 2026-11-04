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

// ReturnCheckOut handles POST /api/v1/bookings/{id}/return-check-out
func (s *Server) ReturnCheckOut(w http.ResponseWriter, r *http.Request) {
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

	var body generated.ReturnCheckOutRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Malformed JSON body.",
		})
		return
	}

	notes := ""
	if body.DamageNotes != nil {
		notes = *body.DamageNotes
	}
	var uris []string
	if body.PhotoUris != nil {
		uris = *body.PhotoUris
	}

	b, err := s.bookingService.ReturnCheckOut(r.Context(), bookingApp.ReturnCheckOutInput{
		BookingID:           bookingUUID,
		ActorID:             userUUID,
		ExpectedVersion:     body.ExpectedVersion,
		OperatingHoursMeter: float64(body.OperatingHoursMeter),
		PhotoURIs:           uris,
		DamageNotes:         notes,
	})
	if err != nil {
		if errors.Is(err, bookingDomain.ErrInvalidMeterReading) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Meter Reading",
				Code:   "ERR_INVALID_METER_READING",
				Detail: "Return check-out hour meter cannot regress below the check-in reading.",
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
