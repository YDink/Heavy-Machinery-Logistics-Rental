package http

import (
	"encoding/json"
	"net/http"
	"time"

	logisticsApp "github.com/dinkmachinery/dink-core/internal/modules/logistics/application"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// AcceptLogisticsTransshipment handles POST /api/v1/bookings/{id}/logistics/transshipment-accept
func (s *Server) AcceptLogisticsTransshipment(w http.ResponseWriter, r *http.Request) {
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

	subBookingID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BOOKING_ID",
			Detail: "Invalid booking ID in path.",
		})
		return
	}

	var req generated.AcceptTransshipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Request Payload",
			Code:   "ERR_INVALID_JSON",
			Detail: "Request body could not be parsed.",
		})
		return
	}

	trans, subStatus, err := s.logisticsService.AcceptTransshipment(r.Context(), logisticsApp.AcceptTransshipmentCommand{
		SubBookingID:       subBookingID,
		CallerUserID:       userUUID,
		ReplacementFleetID: req.ReplacementFleetId,
	})
	if err != nil {
		writeLogisticsDomainError(w, r, err)
		return
	}

	acceptedAt := time.Now().UTC()
	if trans.AcceptedAt != nil {
		acceptedAt = *trans.AcceptedAt
	}

	data := generated.AcceptTransshipmentResponseData{
		TransshipmentId:    trans.ID,
		SubBookingId:       trans.SubBookingID,
		ParentBookingId:    trans.ParentBookingID,
		Status:             string(trans.Status),
		SubBookingStatus:   string(subStatus),
		ReplacementFleetId: req.ReplacementFleetId,
		AcceptedAt:         acceptedAt,
	}

	response.JSON(w, r, http.StatusOK, data)
}
