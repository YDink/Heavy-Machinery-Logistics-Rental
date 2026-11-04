package http

import (
	"errors"
	"net/http"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// GetActiveRentalHUD handles GET /api/v1/bookings/{id}/hud
func (s *Server) GetActiveRentalHUD(w http.ResponseWriter, r *http.Request) {
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

	hud, err := s.rentalService.GetActiveRentalHUD(r.Context(), bookingID, userUUID)
	if err != nil {
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
			Title:  "HUD Not Found",
			Code:   "ERR_NOT_FOUND",
			Detail: err.Error(),
		})
		return
	}

	data := generated.ActiveRentalHUDResponseData{
		BookingId: hud.BookingID,
		Fleet: struct {
			Id          openapi_types.UUID `json:"id"`
			PlateNumber string             `json:"plateNumber"`
			Title       string             `json:"title"`
		}{
			Id:          hud.FleetID,
			PlateNumber: hud.PlateNumber,
			Title:       hud.FleetTitle,
		},
		Status: hud.Status,
	}

	if hud.OperatorID != nil {
		data.Operator = &struct {
			FullName     *string             `json:"fullName,omitempty"`
			Id           *openapi_types.UUID `json:"id,omitempty"`
			LicenseClass *string             `json:"licenseClass,omitempty"`
			PhoneNumber  *string             `json:"phoneNumber,omitempty"`
		}{
			FullName:     &hud.OperatorFullName,
			Id:           hud.OperatorID,
			LicenseClass: &hud.OperatorLicenseClass,
			PhoneNumber:  &hud.OperatorPhoneNumber,
		}
	}

	data.QuickActionsEnabled.CanExtend = hud.CanExtend
	data.QuickActionsEnabled.CanInitiateReturn = hud.CanInitiateReturn
	data.QuickActionsEnabled.CanReportBreakdown = hud.CanReportBreakdown

	data.ShiftTelemetry.AccumulatedHoursToday = hud.AccumulatedHoursToday
	data.ShiftTelemetry.AllowedHoursPerDay = hud.AllowedHoursPerDay
	data.ShiftTelemetry.CurrentRentalDay = hud.CurrentRentalDay
	data.ShiftTelemetry.ExcessHoursTotal = hud.ExcessHoursTotal
	data.ShiftTelemetry.HandoverBaselineHours = hud.HandoverBaselineHours
	data.ShiftTelemetry.ShiftTimerPaused = hud.ShiftTimerPaused
	data.ShiftTelemetry.TotalRentalDays = hud.TotalRentalDays

	response.JSON(w, r, http.StatusOK, data)
}
