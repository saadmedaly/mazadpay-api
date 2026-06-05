package config

import (
	"os"
	"strings"
)

// Config holds all runtime configuration read from environment variables.
type Config struct {
	Port               string
	AppEnv             string
	AppVersion         string
	CORSAllowedOrigins []string
}

const (
	defaultPort    = "8080"
	defaultEnv     = "development"
	defaultVersion = "0.1.0"
	defaultOrigins = "http://localhost:5500,http://127.0.0.1:5500,https://mazadpay.com,https://www.mazadpay.com,https://admin.mazadpay.com"
)

// Load reads configuration from environment variables and returns a Config.
// Missing variables fall back to safe defaults.
func Load() *Config {
	origins := getEnv("CORS_ALLOWED_ORIGINS", defaultOrigins)
	parsed := []string{}
	for _, o := range strings.Split(origins, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			parsed = append(parsed, o)
		}
	}

	return &Config{
		Port:               getEnv("PORT", defaultPort),
		AppEnv:             getEnv("APP_ENV", defaultEnv),
		AppVersion:         getEnv("APP_VERSION", defaultVersion),
		CORSAllowedOrigins: parsed,
	}
}

// IsDevelopment returns true when running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
