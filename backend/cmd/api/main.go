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
	machineryApp "github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	machineryInfra "github.com/dinkmachinery/dink-core/internal/modules/machinery/infrastructure"
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
	var categoryRepo machineryDomain.CategoryRepository
	var specDefRepo machineryDomain.SpecDefinitionRepository
	var fleetRepo machineryDomain.FleetRepository
	var mediaRepo machineryDomain.FleetMediaRepository
	var specRepo machineryDomain.FleetSpecRepository

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err == nil && dbPool.Ping(ctx) == nil {
		queries := db.New(dbPool)
		userRepo = infrastructure.NewPostgresUserRepository(queries)
		pgMachineryRepo := machineryInfra.NewPostgresMachineryRepository(queries)
		categoryRepo = pgMachineryRepo
		specDefRepo = pgMachineryRepo
		fleetRepo = pgMachineryRepo
		mediaRepo = pgMachineryRepo
		specRepo = pgMachineryRepo
		logger.Info("Connected to PostgreSQL database for repositories")
	} else {
		logger.Warn("PostgreSQL connection unavailable; using in-memory caches for bootstrapping")
		userRepo = infrastructure.NewMemoryUserRepository()
		memMachineryRepo := machineryInfra.NewMemoryMachineryRepository()
		categoryRepo = memMachineryRepo
		specDefRepo = memMachineryRepo
		fleetRepo = memMachineryRepo
		mediaRepo = memMachineryRepo
		specRepo = memMachineryRepo
	}

	authService := application.NewAuthService(userRepo, otpService, tokenService)
	machineryService := machineryApp.NewMachineryService(
		categoryRepo,
		specDefRepo,
		fleetRepo,
		mediaRepo,
		specRepo,
	)
	server := transport.NewServer(authService, tokenService, machineryService)

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

	logger.Info("Backend modular monolith initialized with Auth and Machinery partitions",
		zap.String("environment", string(cfg.Env)),
		zap.String("port", cfg.Port),
		zap.Bool("mockOTP", cfg.MockOTP),
	)

	serverAddr := ":" + cfg.Port
	_ = http.ListenAndServe(serverAddr, r)
}
