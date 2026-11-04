package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	bookingApp "github.com/dinkmachinery/dink-core/internal/modules/booking/application"
	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// CreateBooking handles POST /api/v1/bookings
func (s *Server) CreateBooking(w http.ResponseWriter, r *http.Request) {
	if s.bookingService == nil {
		s.notImplemented(w, r, "05-05")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing or malformed Authorization header.",
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

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_INVALID_SUBJECT",
			Detail: "Subject claim is not a valid UUID.",
		})
		return
	}

	var body generated.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	var attUUIDs []uuid.UUID
	if body.AttachmentIds != nil {
		for _, id := range *body.AttachmentIds {
			attUUIDs = append(attUUIDs, uuid.UUID(id))
		}
	}

	input := bookingApp.CreateBookingInput{
		RenterID:            userUUID,
		FleetID:             uuid.UUID(body.FleetId),
		StartDate:           body.StartDate.Time,
		EndDate:             body.EndDate.Time,
		SiteLocationAddress: body.SiteLocationAddress,
		AttachmentIDs:       attUUIDs,
	}

	b, err := s.bookingService.CreateBooking(r.Context(), input)
	if err != nil {
		if errors.Is(err, bookingDomain.ErrSelfBookingProhibited) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Self Booking Prohibited",
				Code:   "ERR_SELF_BOOKING_PROHIBITED",
				Detail: "Equipment fleet owners cannot book their own machinery.",
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrConcurrentReservationCollision) {
			response.Error(w, r, http.StatusConflict, response.ProblemDetails{
				Title:  "Reservation Collision",
				Code:   "ERR_CONCURRENT_RESERVATION_COLLISION",
				Detail: "This fleet asset or an attachment is currently reserved or on hold for requested dates.",
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrInvalidDateRange) || errors.Is(err, bookingDomain.ErrInvalidDuration) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Rental Duration",
				Code:   "ERR_INVALID_RENTAL_DURATION",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrFleetNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Fleet Not Found",
				Code:   "ERR_FLEET_NOT_FOUND",
				Detail: "The specified machinery fleet could not be found.",
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

	response.JSON(w, r, http.StatusCreated, toGeneratedBooking(b))
}
