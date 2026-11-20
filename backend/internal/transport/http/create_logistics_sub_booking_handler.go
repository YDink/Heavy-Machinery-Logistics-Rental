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

// CreateLogisticsSubBooking handles POST /api/v1/bookings/{id}/logistics/sub-booking
func (s *Server) CreateLogisticsSubBooking(w http.ResponseWriter, r *http.Request) {
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

	bookingID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BOOKING_ID",
			Detail: "Invalid booking ID in path.",
		})
		return
	}

	var req generated.CreateLogisticsSubBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Request Payload",
			Code:   "ERR_INVALID_JSON",
			Detail: "Request body could not be parsed.",
		})
		return
	}

	originName := "Origin Depot"
	if req.OriginCoordinates.LocationName != nil && *req.OriginCoordinates.LocationName != "" {
		originName = *req.OriginCoordinates.LocationName
	}
	destName := "Project Site"
	if req.DestinationCoordinates.LocationName != nil && *req.DestinationCoordinates.LocationName != "" {
		destName = *req.DestinationCoordinates.LocationName
	}

	res, err := s.logisticsService.CreateSubBooking(r.Context(), logisticsApp.CreateSubBookingCommand{
		ParentBookingID:  bookingID,
		CallerUserID:     userUUID,
		LogisticsFleetID: req.LogisticsFleetId,
		PickupDateTime:   req.PickupDateTime,
		Origin: logisticsDomain.GeoPoint{
			Latitude:     req.OriginCoordinates.Latitude,
			Longitude:    req.OriginCoordinates.Longitude,
			LocationName: originName,
		},
		Destination: logisticsDomain.GeoPoint{
			Latitude:     req.DestinationCoordinates.Latitude,
			Longitude:    req.DestinationCoordinates.Longitude,
			LocationName: destName,
		},
	})
	if err != nil {
		writeLogisticsDomainError(w, r, err)
		return
	}

	escortFee := res.EscortFeeMinor
	data := generated.LogisticsSubBookingResponseData{
		SubBookingId:               res.SubBookingID,
		ParentBookingId:            res.ParentBookingID,
		Status:                     string(res.Status),
		FreightHireFeeMinor:        res.FreightHireFeeMinor,
		PlatformCommissionBps:      int32(res.PlatformCommissionBps),
		PlatformCommissionRate:     "0.025",
		CommissionOwedMinor:        res.CommissionOwedMinor,
		Currency:                   res.Currency,
		AbnormalLoadEscortRequired: res.AbnormalLoadEscortRequired,
		EscortFeeMinor:             &escortFee,
		CarrierContact: generated.CarrierContactInfo{
			FullName:    res.CarrierFullName,
			PhoneNumber: res.CarrierPhoneNumber,
		},
	}

	response.JSON(w, r, http.StatusCreated, data)
}
