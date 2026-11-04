package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	machineryApp "github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// LinkFleetAttachment handles POST /api/v1/fleets/{id}/attachments/link
func (s *Server) LinkFleetAttachment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	fleetUUID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Fleet ID",
			Code:   "ERR_INVALID_UUID",
			Detail: "The provided fleet ID is not a valid UUID.",
		})
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing or invalid Bearer authorization token.",
		})
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), tokenStr)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Access token is expired or signature is invalid.",
		})
		return
	}
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Invalid token subject claim.",
		})
		return
	}

	var body generated.LinkFleetAttachmentJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	attUUID := uuid.UUID(body.AttachmentId)

	link, err := s.attachmentService.LinkAttachmentToFleet(r.Context(), userUUID, machineryApp.LinkAttachmentInput{
		FleetID:               fleetUUID,
		AttachmentID:          attUUID,
		BundledDailyRateMinor: body.BundledDailyRateMinor,
	})
	if err != nil {
		if errors.Is(err, machineryDomain.ErrBundledRateExceedsStandalone) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Bundled Attachment Rate",
				Code:   "ERR_BUNDLED_RATE_EXCEEDS_STANDALONE",
				Detail: "The configured bundled daily rate cannot exceed the standalone inventory daily rate.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrFleetNotFound) || errors.Is(err, machineryDomain.ErrAttachmentNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Not Found",
				Code:   "ERR_NOT_FOUND",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrUnauthorizedFleetOwner) || errors.Is(err, machineryDomain.ErrUnauthorizedAttachmentOwner) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: "You do not have permission to link this equipment.",
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Link Request",
			Code:   "ERR_INVALID_LINK",
			Detail: err.Error(),
		})
		return
	}

	discountPct := float32(0)
	if link.StandaloneDailyRate > 0 {
		discountPct = float32(float64(link.StandaloneDailyRate-link.BundledDailyRateMinor) / float64(link.StandaloneDailyRate) * 100.0)
	}

	data := generated.LinkFleetAttachmentResponseData{
		FleetId:               openapi_types.UUID(link.FleetID),
		AttachmentId:          openapi_types.UUID(link.AttachmentID),
		BundledDailyRateMinor: link.BundledDailyRateMinor,
		DiscountPercentage:    discountPct,
		Status:                "LINKED",
	}

	response.JSON(w, r, http.StatusOK, data)
}
