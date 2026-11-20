package http

import (
	"encoding/json"
	"net/http"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	logisticsApp "github.com/dinkmachinery/dink-core/internal/modules/logistics/application"
	logisticsDomain "github.com/dinkmachinery/dink-core/internal/modules/logistics/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// UpdateLogisticsDispatchStage handles POST /api/v1/bookings/{id}/logistics/dispatch
func (s *Server) UpdateLogisticsDispatchStage(w http.ResponseWriter, r *http.Request) {
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

	var req generated.LogisticsDispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Request Payload",
			Code:   "ERR_INVALID_JSON",
			Detail: "Request body could not be parsed.",
		})
		return
	}

	var photos []logisticsDomain.WinchingPhoto
	if req.WinchingPhotos != nil {
		for _, p := range *req.WinchingPhotos {
			photos = append(photos, logisticsDomain.WinchingPhoto{
				Angle:      logisticsDomain.WinchingPhotoAngle(p.Angle),
				StorageURI: p.StorageUri,
			})
		}
	}
	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	res, err := s.logisticsService.AdvanceDispatchStage(r.Context(), logisticsApp.AdvanceDispatchStageCommand{
		SubBookingID:   subBookingID,
		CallerUserID:   userUUID,
		TargetStage:    bookingDomain.BookingStatus(req.TargetStage),
		WinchingPhotos: photos,
		Notes:          notes,
	})
	if err != nil {
		writeLogisticsDomainError(w, r, err)
		return
	}

	winchCount := int32(res.WinchingPhotosCount)
	data := generated.LogisticsDispatchResponseData{
		SubBookingId:        res.SubBookingID,
		ParentBookingId:     res.ParentBookingID,
		SubBookingStatus:    string(res.SubBookingStatus),
		ParentBookingStatus: string(res.ParentBookingStatus),
		WinchingPhotosCount: &winchCount,
		UpdatedAt:           res.UpdatedAt,
	}

	response.JSON(w, r, http.StatusOK, data)
}
