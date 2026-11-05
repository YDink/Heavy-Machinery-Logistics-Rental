package http

import (
	"net/http"

	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// FinalReconcile handles POST /api/v1/bookings/{id}/final-reconcile
func (s *Server) FinalReconcile(w http.ResponseWriter, r *http.Request) {
	_, err := s.authenticateUser(r)
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

	settlement, err := s.rentalService.ReconcileFinalSettlement(r.Context(), bookingID)
	if err != nil {
		response.Error(w, r, http.StatusConflict, response.ProblemDetails{
			Title:  "Conflict",
			Code:   "ERR_FINAL_RECONCILE_FAILED",
			Detail: err.Error(),
		})
		return
	}

	resp := generated.FinalReconciliationResponseData{
		BookingId:                  settlement.BookingID,
		BaseMachineFeeMinor:        settlement.BaseMachineFeeMinor,
		BundledAttachmentsFeeMinor: settlement.BundledAttachmentsFeeMinor,
		ContractualAllowedHours:    settlement.ContractualAllowedHours,
		ActualOperatingHours:       settlement.ActualOperatingHours,
		CreditedBreakdownHours:     settlement.CreditedBreakdownHours,
		BillableExcessHours:        settlement.BillableExcessHours,
		ExcessHourlyRateMinor:      settlement.ExcessHourlyRateMinor,
		ExcessHoursSurchargeMinor:  settlement.ExcessHoursSurchargeMinor,
		AcceptedDamageMinor:        settlement.AcceptedDamageMinor,
		GrossLiabilitiesMinor:      settlement.GrossLiabilitiesMinor,
		PrepaymentSettledMinor:     settlement.PrepaymentSettledMinor,
		NetBalanceMinor:            settlement.NetBalanceMinor,
		Currency:                   settlement.Currency,
		Status:                     settlement.Status,
		ReconciledAt:               settlement.CreatedAt,
	}

	response.JSON(w, r, http.StatusOK, resp)
}
