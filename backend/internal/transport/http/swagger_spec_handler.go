package http

import (
	_ "embed"
	"net/http"
	"os"
)

//go:embed openapi.bundled.yaml
var embeddedOpenAPISpec []byte

// SwaggerSpec serves the canonical OpenAPI 3.1.0 specification bundle.
// In development, it searches candidate filesystem paths first to enable instant hot-reload,
// falling back cleanly to the embedded distribution bundle in production or containerized environments.
func (s *Server) SwaggerSpec(w http.ResponseWriter, r *http.Request) {
	candidatePaths := []string{
		os.Getenv("OPENAPI_SPEC_PATH"),
		"contracts/openapi.bundled.yaml",
		"../contracts/openapi.bundled.yaml",
		"../../contracts/openapi.bundled.yaml",
		"../../../contracts/openapi.bundled.yaml",
	}

	for _, path := range candidatePaths {
		if path == "" {
			continue
		}
		if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
			w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data)
			return
		}
	}

	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(embeddedOpenAPISpec)
}
