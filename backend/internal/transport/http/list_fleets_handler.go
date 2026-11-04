package http

import (
	"net/http"

	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// ListFleets handles GET /api/v1/fleets
func (s *Server) ListFleets(w http.ResponseWriter, r *http.Request) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}
	response.JSON(w, r, http.StatusOK, []generated.FleetAsset{})
}
