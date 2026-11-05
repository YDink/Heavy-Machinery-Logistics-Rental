package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	rentalDomain "github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// InspectReturn handles POST /api/v1/bookings/{id}/return-inspect
func (s *Server) InspectReturn(w http.ResponseWriter, r *http.Request) {
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

	var req generated.ReturnInspectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Malformed request payload.",
		})
		return
	}

	fluidOk := true
	if req.FluidLevelsOk != nil {
		fluidOk = *req.FluidLevelsOk
	}
	glassOk := true
	if req.SafetyGlassIntact != nil {
		glassOk = *req.SafetyGlassIntact
	}
	steadyOk := true
	if req.AccelerometerSteady != nil {
		steadyOk = *req.AccelerometerSteady
	}
	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	domainPhotos := make([]rentalDomain.InspectionPhoto, len(req.InspectionPhotos))
	for i, p := range req.InspectionPhotos {
		domainPhotos[i] = rentalDomain.InspectionPhoto{
			Angle:      string(p.Angle),
			StorageURI: p.StorageUri,
		}
	}

	cmd := rentalApp.ReturnInspectCommand{
		BookingID:           bookingID,
		InspectorID:         userUUID,
		CheckOutMeterHours:  req.CheckOutMeterHours,
		FluidLevelsOk:       fluidOk,
		SafetyGlassIntact:   glassOk,
		AccelerometerSteady: steadyOk,
		Notes:               notes,
		Photos:              domainPhotos,
	}

	inspection, err := s.rentalService.InspectReturn(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, rentalDomain.ErrReturnMeterRegression) || errors.Is(err, rentalDomain.ErrMissingReturnPhotos) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Unprocessable Entity",
				Code:   "ERR_INSPECTION_VALIDATION_FAILED",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusConflict, response.ProblemDetails{
			Title:  "Conflict",
			Code:   "ERR_RETURN_INSPECTION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	resp := generated.ReturnInspectionResponseData{
		BookingId:                  inspection.BookingID,
		CheckOutMeterHours:         inspection.CheckOutMeterHours,
		DamageAuditWindowExpiresAt: inspection.CreatedAt.Add(24 * time.Hour),
		InspectedAt:                inspection.CreatedAt,
		PhotosIngestedCount:        len(inspection.Photos),
		Status:                     "DAMAGE_PENDING_REVIEW",
	}

	response.JSON(w, r, http.StatusOK, resp)
}
