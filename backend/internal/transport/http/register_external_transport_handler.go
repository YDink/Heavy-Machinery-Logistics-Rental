package http

import (
	"encoding/json"
	"net/http"

	logisticsApp "github.com/dinkmachinery/dink-core/internal/modules/logistics/application"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// RegisterExternalTransport handles POST /api/v1/bookings/{id}/logistics/external-transport
func (s *Server) RegisterExternalTransport(w http.ResponseWriter, r *http.Request) {
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

	var req generated.ExternalTransportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Request Payload",
			Code:   "ERR_INVALID_JSON",
			Detail: "Request body could not be parsed.",
		})
		return
	}

	if req.DriverFullName == "" || req.DriverPhoneNumber == "" || req.PrimeMoverPlateNumber == "" {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Validation Error",
			Code:   "ERR_VALIDATION_FAILED",
			Detail: "driverFullName, driverPhoneNumber, and primeMoverPlateNumber are required.",
		})
		return
	}

	ext, status, err := s.logisticsService.RegisterExternalTransport(r.Context(), logisticsApp.RegisterExternalTransportCommand{
		ParentBookingID:       bookingID,
		CallerUserID:          userUUID,
		CarrierName:           req.CarrierName,
		DriverFullName:        req.DriverFullName,
		DriverPhoneNumber:     req.DriverPhoneNumber,
		PrimeMoverPlateNumber: req.PrimeMoverPlateNumber,
		EstimatedArrivalAt:    req.EstimatedArrivalDateTime,
	})
	if err != nil {
		writeLogisticsDomainError(w, r, err)
		return
	}

	data := generated.ExternalTransportResponseData{
		ParentBookingId:             ext.ParentBookingID,
		Status:                      string(status),
		ExternalTransportRegistered: true,
		CarrierName:                 ext.CarrierName,
		DriverFullName:              ext.DriverFullName,
		DriverPhoneNumber:           ext.DriverPhoneNumber,
		PrimeMoverPlateNumber:       ext.PrimeMoverPlateNumber,
		EstimatedArrivalDateTime:    ext.EstimatedArrivalAt,
	}

	response.JSON(w, r, http.StatusOK, data)
}
