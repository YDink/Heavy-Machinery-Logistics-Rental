package http

import (
	"encoding/json"
	"net/http"

	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// RespondDamageClaim handles POST /api/v1/bookings/{id}/damage-claims/{claimId}/respond
func (s *Server) RespondDamageClaim(w http.ResponseWriter, r *http.Request) {
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

	claimID, err := uuid.Parse(r.PathValue("claimId"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_CLAIM_ID",
			Detail: "Invalid claim ID in path.",
		})
		return
	}

	var req generated.DamageClaimResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Malformed request payload.",
		})
		return
	}

	acceptedMinor := int64(0)
	if req.AcceptedAmountMinor != nil {
		acceptedMinor = *req.AcceptedAmountMinor
	}
	disputeReason := ""
	if req.DisputeReason != nil {
		disputeReason = *req.DisputeReason
	}

	cmd := rentalApp.RespondDamageClaimCommand{
		BookingID:           bookingID,
		ClaimID:             claimID,
		ActorID:             userUUID,
		Action:              string(req.Action),
		AcceptedAmountMinor: acceptedMinor,
		DisputeReason:       disputeReason,
	}

	docket, err := s.rentalService.RespondToDamageClaim(r.Context(), cmd)
	if err != nil {
		response.Error(w, r, http.StatusConflict, response.ProblemDetails{
			Title:  "Conflict",
			Code:   "ERR_DAMAGE_CLAIM_RESPONSE_FAILED",
			Detail: err.Error(),
		})
		return
	}

	resp := generated.DamageClaimResponseResultData{
		ClaimId:             docket.ID,
		BookingId:           docket.BookingID,
		AcceptedAmountMinor: docket.AcceptedAmountMinor,
		Status:              docket.Status,
		UpdatedAt:           docket.UpdatedAt,
	}

	response.JSON(w, r, http.StatusOK, resp)
}
