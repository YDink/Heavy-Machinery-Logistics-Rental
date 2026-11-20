package http

import (
	"errors"
	"net/http"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	logisticsDomain "github.com/dinkmachinery/dink-core/internal/modules/logistics/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
)

func writeLogisticsDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, logisticsDomain.ErrUnauthorizedLogisticsAction),
		errors.Is(err, bookingDomain.ErrUnauthorizedAction):
		response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
			Title:  "Forbidden",
			Code:   "ERR_FORBIDDEN",
			Detail: err.Error(),
		})
	case errors.Is(err, logisticsDomain.ErrLogisticsBookingNotFound),
		errors.Is(err, bookingDomain.ErrBookingNotFound),
		errors.Is(err, logisticsDomain.ErrTransshipmentNotOpen):
		response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
			Title:  "Not Found",
			Code:   "ERR_LOGISTICS_NOT_FOUND",
			Detail: err.Error(),
		})
	case errors.Is(err, logisticsDomain.ErrTrailerPayloadInsufficient):
		response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
			Type:   "https://api.dinkmachinery.com/errors/trailer-payload-exceeded",
			Title:  "Machine Weight Exceeds Lowboy Payload Rating",
			Code:   "ERR_TRAILER_PAYLOAD_INSUFFICIENT",
			Detail: err.Error(),
		})
	case errors.Is(err, logisticsDomain.ErrTrailerDeckTooShort):
		response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
			Type:   "https://api.dinkmachinery.com/errors/trailer-deck-too-short",
			Title:  "Machine Length Exceeds Lowboy Deck Well Length",
			Code:   "ERR_TRAILER_DECK_TOO_SHORT",
			Detail: err.Error(),
		})
	case errors.Is(err, logisticsDomain.ErrTrailerDeckTooNarrow):
		response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
			Type:   "https://api.dinkmachinery.com/errors/trailer-deck-too-narrow",
			Title:  "Machine Width Exceeds Lowboy Deck Width",
			Code:   "ERR_TRAILER_DECK_TOO_NARROW",
			Detail: err.Error(),
		})
	case errors.Is(err, logisticsDomain.ErrCorridorHeightExceeded):
		response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
			Type:   "https://api.dinkmachinery.com/errors/corridor-height-exceeded",
			Title:  "Combined Transport Height Exceeds Corridor Clearance",
			Code:   "ERR_CORRIDOR_HEIGHT_EXCEEDED",
			Detail: err.Error(),
		})
	case errors.Is(err, logisticsDomain.ErrGVMExceeded):
		response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
			Type:   "https://api.dinkmachinery.com/errors/gvm-limit-exceeded",
			Title:  "Gross Vehicle Mass Exceeds Regional Highway Cap",
			Code:   "ERR_GVM_LIMIT_EXCEEDED",
			Detail: err.Error(),
		})
	case errors.Is(err, logisticsDomain.ErrInvalidWinchingPhotos):
		response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
			Title:  "Incomplete Winching Inspection Reticle",
			Code:   "ERR_INVALID_WINCHING_PHOTOS",
			Detail: err.Error(),
		})
	case errors.Is(err, logisticsDomain.ErrTransshipmentWindowExpired):
		response.Error(w, r, http.StatusGone, response.ProblemDetails{
			Type:   "https://api.dinkmachinery.com/errors/transshipment-sla-expired",
			Title:  "Roadside Transshipment Window Expired",
			Code:   "ERR_TRANSSHIPMENT_TIMEOUT",
			Detail: err.Error(),
		})
	case errors.Is(err, logisticsDomain.ErrSelfMoveNotRequired),
		errors.Is(err, logisticsDomain.ErrSubBookingAlreadyExists),
		errors.Is(err, logisticsDomain.ErrExternalTransportBlocked),
		errors.Is(err, logisticsDomain.ErrReverseCascadeAlreadyExecuted),
		errors.Is(err, logisticsDomain.ErrInvalidDispatchStage),
		errors.Is(err, logisticsDomain.ErrInvalidBookingStateForLogistics),
		errors.Is(err, bookingDomain.ErrInvalidStateTransition),
		errors.Is(err, bookingDomain.ErrConcurrencyConflict):
		response.Error(w, r, http.StatusConflict, response.ProblemDetails{
			Title:  "Logistics State Conflict",
			Code:   "ERR_LOGISTICS_CONFLICT",
			Detail: err.Error(),
		})
	default:
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL_SERVER",
			Detail: err.Error(),
		})
	}
}
