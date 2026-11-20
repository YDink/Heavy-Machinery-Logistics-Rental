package http

import (
	"encoding/json"
	"net/http"

	logisticsApp "github.com/dinkmachinery/dink-core/internal/modules/logistics/application"
	logisticsDomain "github.com/dinkmachinery/dink-core/internal/modules/logistics/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// ReportLogisticsTransshipmentAlert handles POST /api/v1/bookings/{id}/logistics/transshipment-alert
func (s *Server) ReportLogisticsTransshipmentAlert(w http.ResponseWriter, r *http.Request) {
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

	var req generated.TransshipmentAlertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Request Payload",
			Code:   "ERR_INVALID_JSON",
			Detail: "Request body could not be parsed.",
		})
		return
	}

	locName := "Highway Corridor"
	if req.CurrentLocation.LocationName != nil && *req.CurrentLocation.LocationName != "" {
		locName = *req.CurrentLocation.LocationName
	}

	res, err := s.logisticsService.ReportTransshipmentAlert(r.Context(), logisticsApp.TransshipmentAlertCommand{
		SubBookingID: subBookingID,
		CallerUserID: userUUID,
		BreakdownLocation: logisticsDomain.GeoPoint{
			Latitude:     req.CurrentLocation.Latitude,
			Longitude:    req.CurrentLocation.Longitude,
			LocationName: locName,
		},
		FailureCategory:    logisticsDomain.TransshipmentFailureCategory(req.FailureCategory),
		FailureDescription: req.FailureDescription,
		Photos:             req.Photos,
	})
	if err != nil {
		writeLogisticsDomainError(w, r, err)
		return
	}

	data := generated.TransshipmentAlertResponseData{
		TransshipmentId:          res.Transshipment.ID,
		SubBookingId:             res.Transshipment.SubBookingID,
		ParentBookingId:          res.Transshipment.ParentBookingID,
		SubBookingStatus:         string(res.SubBookingStatus),
		ParentBookingStatus:      string(res.ParentBookingStatus),
		HotSwapSlaDeadline:       res.Transshipment.HotSwapSLADeadline,
		TransshipmentBountyMinor: res.Transshipment.TransshipmentBountyMinor,
		Currency:                 res.Transshipment.Currency,
		CandidateCarriersCount:   int32(res.CandidateCarriersCount),
	}

	response.JSON(w, r, http.StatusOK, data)
}
