package http

import (
	"encoding/json"
	"errors"
	"net/http"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	rentalDomain "github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// RequestRentalExtension handles POST /api/v1/bookings/{id}/extend
func (s *Server) RequestRentalExtension(w http.ResponseWriter, r *http.Request) {
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

	var req generated.RentalExtensionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BODY",
			Detail: "Malformed request payload.",
		})
		return
	}

	ext, err := s.rentalService.RequestExtension(r.Context(), rentalApp.RentalExtensionCommand{
		BookingID:  bookingID,
		RenterID:   userUUID,
		NewEndDate: req.NewEndDate.Time,
	})
	if err != nil {
		if errors.Is(err, rentalDomain.ErrExtensionDateUnavailable) {
			response.Error(w, r, http.StatusConflict, response.ProblemDetails{
				Title:  "Extension Dates Conflict",
				Code:   "ERR_EXTENSION_CONFLICT",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, rentalDomain.ErrInvalidExtensionDate) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Extension Date",
				Code:   "ERR_INVALID_EXTENSION_DATE",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: "Only renter can request an extension.",
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Extension Request Failed",
			Code:   "ERR_EXTENSION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	data := generated.RentalExtensionResponseData{
		AdditionalGrossMinor:               ext.IncrementalGrossMinor,
		BookingId:                          ext.BookingID,
		Currency:                           ext.Currency,
		ExtensionDays:                      ext.ExtensionDays,
		IncrementalPrepaymentRequiredMinor: ext.IncrementalPrepaymentMinor,
		OwnerApprovalRequired:              true,
		Status:                             ext.Status,
	}

	response.JSON(w, r, http.StatusOK, data)
}
