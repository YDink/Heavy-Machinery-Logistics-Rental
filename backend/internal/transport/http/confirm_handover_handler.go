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

// ConfirmHandover handles POST /api/v1/bookings/{id}/handover-confirm
func (s *Server) ConfirmHandover(w http.ResponseWriter, r *http.Request) {
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

	var req generated.HandoverConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BODY",
			Detail: "Malformed request payload.",
		})
		return
	}

	var touchSig string
	if req.TouchSignatureBase64 != nil {
		touchSig = *req.TouchSignatureBase64
	}
	var lat, lng *float64
	if req.ContractorLocation != nil {
		lat = req.ContractorLocation.Latitude
		lng = req.ContractorLocation.Longitude
	}

	activatedBooking, err := s.rentalService.ConfirmHandover(r.Context(), rentalApp.HandoverConfirmCommand{
		BookingID:            bookingID,
		CallerUserID:         userUUID,
		PIN:                  req.ChallengePin,
		TouchSignatureBase64: touchSig,
		ContractorLat:        lat,
		ContractorLng:        lng,
	})
	if err != nil {
		if errors.Is(err, rentalDomain.ErrPINLockout) {
			response.Error(w, r, http.StatusLocked, response.ProblemDetails{
				Title:  "Handover Challenge Locked Out",
				Code:   "ERR_PIN_LOCKOUT",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, rentalDomain.ErrPINExpired) {
			response.Error(w, r, http.StatusGone, response.ProblemDetails{
				Title:  "Challenge PIN Expired",
				Code:   "ERR_PIN_EXPIRED",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, rentalDomain.ErrInvalidChallengePIN) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Challenge PIN",
				Code:   "ERR_INVALID_CHALLENGE_PIN",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: "Only renter or assigned operator may verify handover.",
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Handover Confirmation Failed",
			Code:   "ERR_CONFIRMATION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	actTime := time.Now().UTC()
	if activatedBooking.ActivatedAt != nil {
		actTime = *activatedBooking.ActivatedAt
	}

	startMeter := 0.0
	if activatedBooking.CheckInMeterHours != nil {
		startMeter = *activatedBooking.CheckInMeterHours
	}

	data := generated.HandoverConfirmResponseData{
		ActivatedAt:              actTime,
		AllocatedShiftHoursTotal: rentalDomain.StandardShiftHoursPerDay,
		BookingId:                activatedBooking.ID,
		ShiftBaselineMeterHours:  startMeter,
		Status:                   string(activatedBooking.Status),
	}

	response.JSON(w, r, http.StatusOK, data)
}
