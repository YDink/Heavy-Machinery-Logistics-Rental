package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	bookingApp "github.com/dinkmachinery/dink-core/internal/modules/booking/application"
	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// SubmitBookingPaymentEvidence handles POST /api/v1/bookings/{id}/payment-evidence
func (s *Server) SubmitBookingPaymentEvidence(w http.ResponseWriter, r *http.Request) {
	if s.bookingService == nil {
		s.notImplemented(w, r, "05-05")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing Authorization header.",
		})
		return
	}

	claims, err := s.tokenService.VerifyAccessToken(r.Context(), strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_INVALID_TOKEN",
			Detail: "Access token is invalid or expired.",
		})
		return
	}

	userUUID, _ := uuid.Parse(claims.Subject)
	idStr := r.PathValue("id")
	bookingUUID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Booking ID",
			Code:   "ERR_INVALID_UUID",
			Detail: "Invalid booking UUID format.",
		})
		return
	}

	var body generated.SubmitPaymentSlipRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Malformed JSON body.",
		})
		return
	}

	txnRef := ""
	if body.TransactionReference != nil {
		txnRef = *body.TransactionReference
	}
	currency := "ETB"
	if body.Currency != nil {
		currency = *body.Currency
	}

	slip, err := s.bookingService.SubmitPaymentSlip(r.Context(), bookingApp.SubmitPaymentSlipInput{
		BookingID:            bookingUUID,
		UploaderID:           userUUID,
		StorageURI:           body.StorageUri,
		FileHashSHA256:       body.FileHash,
		TransactionReference: txnRef,
		AmountMinor:          body.AmountMinor,
		Currency:             currency,
	})
	if err != nil {
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_UNAUTHORIZED_ACTION",
				Detail: "Only the renting contractor can submit payment evidence.",
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrInvalidStateTransition) {
			response.Error(w, r, http.StatusConflict, response.ProblemDetails{
				Title:  "Invalid State Transition",
				Code:   "ERR_INVALID_STATE_TRANSITION",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: err.Error(),
		})
		return
	}

	slipDTO := generated.PaymentSlip{
		Id:                   openapi_types.UUID(slip.ID),
		BookingId:            openapi_types.UUID(slip.BookingID),
		UploaderId:           openapi_types.UUID(slip.UploaderID),
		StorageUri:           slip.StorageURI,
		FileHash:             slip.FileHashSHA256,
		TransactionReference: slip.TransactionReference,
		Amount: generated.MoneyAmount{
			AmountMinor: slip.AmountMinor,
			Currency:    slip.Currency,
		},
		Status:    generated.PaymentSlipStatus(slip.Status),
		CreatedAt: slip.CreatedAt,
	}

	response.JSON(w, r, http.StatusCreated, slipDTO)
}
