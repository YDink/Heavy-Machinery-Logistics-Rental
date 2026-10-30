package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/infrastructure"
	"github.com/dinkmachinery/dink-core/internal/platform/config"
	"github.com/dinkmachinery/dink-core/internal/platform/database/db"
	transport "github.com/dinkmachinery/dink-core/internal/transport/http"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

type HealthResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Error   interface{}            `json:"error"`
	Meta    map[string]interface{} `json:"meta"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load runtime configuration: " + err.Error())
	}

	var logger *zap.Logger
	if cfg.IsDevelopment() {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()

	ctx := context.Background()

	// Initialize Auth Subsystem Infrastructure
	cache := infrastructure.NewMemoryCache()
	smsGateway := infrastructure.NewMockSMSGateway(logger)

	otpService := application.NewOTPService(cache, smsGateway, application.OTPServiceConfig{
		MockOTP: cfg.MockOTP,
	})

	tokenService, err := application.NewTokenService(nil, cache)
	if err != nil {
		panic("failed to initialize token service: " + err.Error())
	}

	var userRepo domain.UserRepository
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err == nil && dbPool.Ping(ctx) == nil {
		queries := db.New(dbPool)
		userRepo = infrastructure.NewPostgresUserRepository(queries)
		logger.Info("Connected to PostgreSQL database for identity repository")
	} else {
		logger.Warn("PostgreSQL connection unavailable; using in-memory identity cache for bootstrapping")
		userRepo = infrastructure.NewMemoryUserRepository()
	}

	authService := application.NewAuthService(userRepo, otpService, tokenService)
	server := transport.NewServer(authService, tokenService)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{
			Success: true,
			Data: map[string]interface{}{
				"status":      "healthy",
				"version":     "1.0.0-phase2",
				"environment": string(cfg.Env),
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

	// Mount OpenAPI routes onto /api/v1
	generated.HandlerWithOptions(server, generated.ChiServerOptions{
		BaseURL:    "/api/v1",
		BaseRouter: r,
	})

	logger.Info("Backend modular monolith initialized with Auth partition",
		zap.String("environment", string(cfg.Env)),
		zap.String("port", cfg.Port),
		zap.Bool("mockOTP", cfg.MockOTP),
	)

	serverAddr := ":" + cfg.Port
	_ = http.ListenAndServe(serverAddr, r)
}
