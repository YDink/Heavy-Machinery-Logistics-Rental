package http

import (
	"net/http"
)

// SwaggerRedirect redirects /swagger to /swagger/.
func (s *Server) SwaggerRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
}
