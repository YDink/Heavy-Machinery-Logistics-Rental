package http

import (
	"errors"
	"net/http"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// GenerateHandoverChallenge handles POST /api/v1/bookings/{id}/handover-challenge
func (s *Server) GenerateHandoverChallenge(w http.ResponseWriter, r *http.Request) {
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

	challenge, err := s.rentalService.GenerateHandoverPIN(r.Context(), bookingID, userUUID)
	if err != nil {
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: "Only the fleet owner can generate a handover challenge PIN.",
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Challenge Generation Failed",
			Code:   "ERR_CHALLENGE_GENERATION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	attemptCount := challenge.AttemptCount
	data := generated.HandoverChallengeResponseData{
		AttemptCount: &attemptCount,
		BookingId:    challenge.BookingID,
		ChallengePin: challenge.ChallengePIN,
		ExpiresAt:    challenge.ExpiresAt,
		MaxAttempts:  challenge.MaxAttempts,
	}

	response.JSON(w, r, http.StatusOK, data)
}
