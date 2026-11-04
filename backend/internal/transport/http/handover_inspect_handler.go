package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	rentalDomain "github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// HandoverInspect handles POST /api/v1/bookings/{id}/handover-inspect
func (s *Server) HandoverInspect(w http.ResponseWriter, r *http.Request) {
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

	var req generated.HandoverInspectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BODY",
			Detail: "Malformed request payload.",
		})
		return
	}

	photos := make([]rentalDomain.InspectionPhoto, len(req.InspectionPhotos))
	for i, p := range req.InspectionPhotos {
		photos[i] = rentalDomain.InspectionPhoto{
			Angle:      string(p.Angle),
			StorageURI: p.StorageUri,
		}
	}

	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	fluidLevelsOk := true
	if req.FluidLevelsOk != nil {
		fluidLevelsOk = *req.FluidLevelsOk
	}
	safetyGlassIntact := true
	if req.SafetyGlassIntact != nil {
		safetyGlassIntact = *req.SafetyGlassIntact
	}

	inspection, err := s.rentalService.InspectHandover(r.Context(), rentalApp.HandoverInspectCommand{
		BookingID:         bookingID,
		InspectorID:       userUUID,
		CheckInMeterHours: req.CheckInMeterHours,
		FluidLevelsOk:     fluidLevelsOk,
		SafetyGlassIntact: safetyGlassIntact,
		Notes:             notes,
		Photos:            photos,
	})
	if err != nil {
		if errors.Is(err, rentalDomain.ErrMeterReadingRegression) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Meter Reading Regression",
				Code:   "ERR_METER_READING_REGRESSION",
				Detail: err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "photo") {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Inspection Photos",
				Code:   "ERR_INVALID_PHOTOS",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Inspection Failed",
			Code:   "ERR_INSPECTION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	data := generated.HandoverInspectionResponseData{
		BookingId:           inspection.BookingID,
		InspectedAt:         inspection.CreatedAt,
		MeterVerified:       true,
		PhotosIngestedCount: len(inspection.Photos),
		Status:              "INSPECTED",
	}

	response.JSON(w, r, http.StatusOK, data)
}
