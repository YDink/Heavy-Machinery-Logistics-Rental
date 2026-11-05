package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	transport "github.com/dinkmachinery/dink-core/internal/transport/http"
)

func TestSwagger_Handlers(t *testing.T) {
	server := transport.NewServer(nil, nil, nil, nil, nil, nil, nil)
	mux := http.NewServeMux()
	server.RegisterRoutes(mux)

	t.Run("GET /swagger/ returns Swagger UI HTML dashboard", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200 OK, got %d", rec.Code)
		}
		contentType := rec.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("expected Content-Type text/html, got %s", contentType)
		}
		body := rec.Body.String()
		if !strings.Contains(body, "SwaggerUIBundle") {
			t.Errorf("expected body to contain SwaggerUIBundle, got %s", body)
		}
		if !strings.Contains(body, "/swagger/openapi.yaml") {
			t.Errorf("expected body to reference /swagger/openapi.yaml")
		}
	})

	t.Run("GET /swagger/index.html returns Swagger UI HTML dashboard", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200 OK, got %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "SwaggerUIBundle") {
			t.Errorf("expected body to contain SwaggerUIBundle")
		}
	})

	t.Run("GET /swagger redirects to /swagger/", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMovedPermanently {
			t.Fatalf("expected status 301 Moved Permanently, got %d", rec.Code)
		}
		loc := rec.Header().Get("Location")
		if loc != "/swagger/" {
			t.Errorf("expected redirect location /swagger/, got %s", loc)
		}
	})

	t.Run("GET /swagger/openapi.yaml returns canonical OpenAPI 3.1.0 contract", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/openapi.yaml", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200 OK, got %d", rec.Code)
		}
		contentType := rec.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/yaml") {
			t.Errorf("expected Content-Type application/yaml, got %s", contentType)
		}
		body := rec.Body.String()
		if !strings.Contains(body, "openapi: 3.1.0") {
			t.Errorf("expected YAML to specify openapi: 3.1.0")
		}
		if !strings.Contains(body, "Dink Machinery Marketplace API") {
			t.Errorf("expected YAML to contain Dink Machinery Marketplace API title")
		}
	})
}

func TestSwagger_Disabled(t *testing.T) {
	server := transport.NewServer(nil, nil, nil, nil, nil, nil, nil)
	server.SetSwaggerEnabled(false)

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)

	t.Run("GET /swagger/ returns 404 when disabled", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status 404 Not Found when disabled, got %d", rec.Code)
		}
	})

	t.Run("GET /swagger/openapi.yaml returns 404 when disabled", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/openapi.yaml", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status 404 Not Found when disabled, got %d", rec.Code)
		}
	})
}
