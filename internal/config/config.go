package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime configuration read from environment variables.
type Config struct {
	Port               string
	AppEnv             string
	AppVersion         string
	CORSAllowedOrigins []string

	// Database — optional; server starts without it in development.
	DatabaseURL       string
	DBMaxConns        int32
	DBMinConns        int32
	DBMaxConnLifetime time.Duration
}

const (
	defaultPort    = "8080"
	defaultEnv     = "development"
	defaultVersion = "0.1.0"
	defaultOrigins = "http://localhost:5500,http://127.0.0.1:5500,https://mazadpay.com,https://www.mazadpay.com,https://admin.mazadpay.com"
)

// Load reads configuration from environment variables with safe defaults.
func Load() *Config {
	origins := getEnv("CORS_ALLOWED_ORIGINS", defaultOrigins)
	parsed := []string{}
	for _, o := range strings.Split(origins, ",") {
		if o = strings.TrimSpace(o); o != "" {
			parsed = append(parsed, o)
		}
	}

	return &Config{
		Port:               getEnv("PORT", defaultPort),
		AppEnv:             getEnv("APP_ENV", defaultEnv),
		AppVersion:         getEnv("APP_VERSION", defaultVersion),
		CORSAllowedOrigins: parsed,

		DatabaseURL:       os.Getenv("DATABASE_URL"), // empty = disabled
		DBMaxConns:        int32(getEnvInt("DB_MAX_CONNS", 10)),
		DBMinConns:        int32(getEnvInt("DB_MIN_CONNS", 1)),
		DBMaxConnLifetime: getEnvDuration("DB_MAX_CONN_LIFETIME", time.Hour),
	}
}

// IsDevelopment returns true when running in development mode.
func (c *Config) IsDevelopment() bool { return c.AppEnv == "development" }

// HasDatabase returns true when a DATABASE_URL is configured.
func (c *Config) HasDatabase() bool { return c.DatabaseURL != "" }

// ---- helpers ----------------------------------------------------------------

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
