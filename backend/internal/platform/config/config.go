package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
)

// Config represents runtime configuration for the Dink Machinery backend modular monolith.
type Config struct {
	Env            Environment
	Port           string
	DatabaseURL    string
	ValkeyURL      string
	SeaweedFSURL   string
	ClamAVAddress  string
	LogLevel       string
	MockOTP        bool
	AllowedOrigins []string
}

// Load loads configuration from environment variables with deterministic defaults.
func Load() (*Config, error) {
	envStr := strings.ToLower(getEnv("APP_ENV", string(EnvDevelopment)))
	var env Environment
	switch envStr {
	case "prod", "production":
		env = EnvProduction
	case "stage", "staging":
		env = EnvStaging
	default:
		env = EnvDevelopment
	}

	cfg := &Config{
		Env:           env,
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/dink_machinery?sslmode=disable"),
		ValkeyURL:     getEnv("VALKEY_URL", "localhost:6379"),
		SeaweedFSURL:  getEnv("SEAWEEDFS_URL", "http://localhost:8333"),
		ClamAVAddress: getEnv("CLAMAV_ADDR", "localhost:3310"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		MockOTP:       getEnvBool("MOCK_OTP", env == EnvDevelopment),
	}

	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		cfg.AllowedOrigins = strings.Split(origins, ",")
	} else {
		if env == EnvProduction {
			cfg.AllowedOrigins = []string{"https://admin.dinkmachinery.com"}
		} else if env == EnvStaging {
			cfg.AllowedOrigins = []string{"https://staging-admin.dinkmachinery.com"}
		} else {
			cfg.AllowedOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
		}
	}

	// Production Invariant Safety Checks
	if cfg.Env == EnvProduction {
		if cfg.MockOTP {
			return nil, fmt.Errorf("FATAL: MOCK_OTP must NEVER be enabled in production environment")
		}
		if strings.Contains(cfg.DatabaseURL, "sslmode=disable") {
			return nil, fmt.Errorf("FATAL: database connection must require SSL/TLS in production environment")
		}
	}

	return cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.Env == EnvDevelopment
}

func (c *Config) IsStaging() bool {
	return c.Env == EnvStaging
}

func (c *Config) IsProduction() bool {
	return c.Env == EnvProduction
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		parsed, err := strconv.ParseBool(val)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
