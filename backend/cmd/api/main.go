package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type HealthResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Error   interface{}            `json:"error"`
	Meta    map[string]interface{} `json:"meta"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{
			Success: true,
			Data: map[string]interface{}{
				"status":  "healthy",
				"version": "1.0.0-phase1",
			},
			Error: nil,
			Meta: map[string]interface{}{
				"version":   "1.0",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})

	logger.Info("Backend modular monolith initialized", zap.String("port", "8080"))
	_ = http.ListenAndServe(":8080", r)
}
