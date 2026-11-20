package http

import (
	"encoding/json"
	"io"
	"net/http"

	logisticsApp "github.com/dinkmachinery/dink-core/internal/modules/logistics/application"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// TriggerLogisticsReverseCascade handles POST /api/v1/bookings/{id}/logistics/reverse-cascade
func (s *Server) TriggerLogisticsReverseCascade(w http.ResponseWriter, r *http.Request) {
	if s.logisticsService == nil {
		s.notImplemented(w, r, "05-08")
		return
	}

	userUUID, err := s.authenticateUser(r)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: err.Error(),
		})
		return
	}

	parentBookingID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BOOKING_ID",
			Detail: "Invalid booking ID in path.",
		})
		return
	}

	var req generated.ReverseCascadeRequest
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil && err != io.EOF {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Request Payload",
				Code:   "ERR_INVALID_JSON",
				Detail: "Request body could not be parsed.",
			})
			return
		}
	}

	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	res, err := s.logisticsService.TriggerReverseCascade(r.Context(), logisticsApp.ReverseCascadeCommand{
		ParentBookingID: parentBookingID,
		CallerUserID:    userUUID,
		Reason:          reason,
	})
	if err != nil {
		writeLogisticsDomainError(w, r, err)
		return
	}

	data := generated.ReverseCascadeResponseData{
		ParentBookingId:                 res.Cascade.ParentBookingID,
		SubBookingId:                    res.Cascade.SubBookingID,
		Resolution:                      res.Cascade.Resolution,
		RenterRefundDueMinor:            res.Cascade.RenterRefundDueMinor,
		OwnerPenaltyLeviedMinor:         res.Cascade.OwnerPenaltyLeviedMinor,
		OwnerCommissionRestitutionMinor: res.Cascade.OwnerCommissionRestitutionMinor,
		CarrierDefaultPenaltyMinor:      res.Cascade.CarrierDefaultPenaltyMinor,
		Currency:                        res.Cascade.Currency,
		ParentStatus:                    string(res.ParentStatus),
		SubStatus:                       string(res.SubStatus),
	}

	response.JSON(w, r, http.StatusOK, data)
}
