package config_test

import (
	"os"
	"testing"

	"github.com/dinkmachinery/dink-core/internal/platform/config"
)

func TestConfig_DevelopmentDefaults(t *testing.T) {
	os.Clearenv()
	_ = os.Setenv("APP_ENV", "development")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected clean config load, got error: %v", err)
	}

	if !cfg.IsDevelopment() {
		t.Errorf("expected IsDevelopment to be true")
	}
	if !cfg.MockOTP {
		t.Errorf("expected MockOTP to be true in development")
	}
	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
}

func TestConfig_StagingEnvironment(t *testing.T) {
	os.Clearenv()
	_ = os.Setenv("APP_ENV", "staging")
	_ = os.Setenv("MOCK_OTP", "false")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected clean config load, got error: %v", err)
	}

	if !cfg.IsStaging() {
		t.Errorf("expected IsStaging to be true")
	}
	if cfg.MockOTP {
		t.Errorf("expected MockOTP to be false in staging")
	}
}

func TestConfig_ProductionRejectsMockOTP(t *testing.T) {
	os.Clearenv()
	_ = os.Setenv("APP_ENV", "production")
	_ = os.Setenv("MOCK_OTP", "true")
	_ = os.Setenv("DATABASE_URL", "postgres://user:pass@db:5432/dink?sslmode=require")

	_, err := config.Load()
	if err == nil {
		t.Fatalf("expected error when MOCK_OTP is enabled in production, got nil")
	}
}

func TestConfig_ProductionRejectsInsecureSSL(t *testing.T) {
	os.Clearenv()
	_ = os.Setenv("APP_ENV", "production")
	_ = os.Setenv("MOCK_OTP", "false")
	_ = os.Setenv("DATABASE_URL", "postgres://user:pass@db:5432/dink?sslmode=disable")

	_, err := config.Load()
	if err == nil {
		t.Fatalf("expected error when sslmode=disable in production, got nil")
	}
}
