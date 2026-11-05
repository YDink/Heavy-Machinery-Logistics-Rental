package http

import (
	"encoding/json"
	"net/http"

	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// LodgeDamageClaim handles POST /api/v1/bookings/{id}/damage-claims
func (s *Server) LodgeDamageClaim(w http.ResponseWriter, r *http.Request) {
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

	var req generated.DamageClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Malformed request payload.",
		})
		return
	}

	quoteURI := ""
	if req.RepairQuoteUri != nil {
		quoteURI = *req.RepairQuoteUri
	}
	photos := []string{}
	if req.DamagePhotos != nil {
		photos = *req.DamagePhotos
	}

	cmd := rentalApp.DamageClaimCommand{
		BookingID:          bookingID,
		ClaimantID:         userUUID,
		ComponentCategory:  string(req.ComponentCategory),
		IncidentNarrative:  req.IncidentNarrative,
		RepairQuoteURI:     quoteURI,
		ClaimedAmountMinor: req.ClaimedAmountMinor,
		Currency:           req.Currency,
		DamagePhotos:       photos,
	}

	docket, err := s.rentalService.LodgeDamageClaim(r.Context(), cmd)
	if err != nil {
		response.Error(w, r, http.StatusConflict, response.ProblemDetails{
			Title:  "Conflict",
			Code:   "ERR_LODGE_DAMAGE_CLAIM_FAILED",
			Detail: err.Error(),
		})
		return
	}

	resp := generated.DamageClaimResponseData{
		ClaimId:             docket.ID,
		BookingId:           docket.BookingID,
		ComponentCategory:   docket.ComponentCategory,
		ClaimedAmountMinor:  docket.ClaimedAmountMinor,
		AcceptedAmountMinor: &docket.AcceptedAmountMinor,
		Currency:            docket.Currency,
		Status:              docket.Status,
		CreatedAt:           docket.CreatedAt,
	}

	response.JSON(w, r, http.StatusCreated, resp)
}
