package http

import (
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/dinkmachinery/dink-core/internal/platform/i18n"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// ListCategories handles GET /api/v1/categories
func (s *Server) ListCategories(w http.ResponseWriter, r *http.Request) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-02")
		return
	}

	lang := i18n.GetLanguage(r.Context())
	categories, err := s.machineryService.ListCategories(r.Context())
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to list machinery categories.",
		})
		return
	}

	result := make([]generated.Category, len(categories))
	for i, c := range categories {
		name := c.NameEN
		if lang == "am" && c.NameAM != "" {
			name = c.NameAM
		}
		var parentID *openapi_types.UUID
		if c.ParentID != nil {
			pid := openapi_types.UUID(*c.ParentID)
			parentID = &pid
		}
		result[i] = generated.Category{
			Id:       openapi_types.UUID(c.ID),
			Code:     c.Code,
			Name:     name,
			IconUrl:  c.IconURL,
			ParentId: parentID,
		}
	}

	response.JSON(w, r, http.StatusOK, result)
}
