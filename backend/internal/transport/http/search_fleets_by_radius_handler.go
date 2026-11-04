package http

import (
	"net/http"
)

// SearchFleetsByRadius handles legacy GET /api/v1/geo/fleets/radius-search
func (s *Server) SearchFleetsByRadius(w http.ResponseWriter, r *http.Request) {
	s.SearchFleets(w, r)
}
