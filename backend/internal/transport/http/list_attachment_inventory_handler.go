package http

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// ListAttachmentInventory handles GET /api/v1/attachments/inventory
func (s *Server) ListAttachmentInventory(w http.ResponseWriter, r *http.Request) {
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

	statusFilter := r.URL.Query().Get("status")

	items, err := s.attachmentService.ListOwnerInventory(r.Context(), userUUID, statusFilter)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to list attachment inventory.",
		})
		return
	}

	inventory := make([]generated.AttachmentInventoryItem, 0, len(items))
	for _, it := range items {
		var primaryPhoto *string
		if len(it.Photos) > 0 {
			primaryPhoto = &it.Photos[0].StorageURI
		}
		inventory = append(inventory, generated.AttachmentInventoryItem{
			AttachmentId:             openapi_types.UUID(it.ID),
			AttachmentType:           string(it.AttachmentType),
			Make:                     it.Make,
			Model:                    it.Model,
			SubTypeSize:              it.SubTypeSize,
			Status:                   string(it.Status),
			StandaloneDailyRateMinor: it.StandaloneDailyRateMinor,
			IsLinkedToFleet:          false,
			PrimaryPhotoUri:          primaryPhoto,
		})
	}

	response.JSON(w, r, http.StatusOK, inventory)
}
