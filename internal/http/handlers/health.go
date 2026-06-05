package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/saadmedaly/mazadpay-api/internal/config"
	"github.com/saadmedaly/mazadpay-api/internal/db"
)

// HealthHandler holds dependencies for health endpoints.
type HealthHandler struct {
	cfg  *config.Config
	pool *pgxpool.Pool // nil when DB is not configured
}

// NewHealthHandler creates a HealthHandler.
// pool may be nil — DB health check handles that gracefully.
func NewHealthHandler(cfg *config.Config, pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{cfg: cfg, pool: pool}
}

// HealthResponse is the JSON body for the main health endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Env     string `json:"env"`
	Version string `json:"version"`
}

// Check handles GET /health and GET /api/v1/health.
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(HealthResponse{
		Status:  "ok",
		Service: "mazadpay-api",
		Env:     h.cfg.AppEnv,
		Version: h.cfg.AppVersion,
	})
}

// DBCheck handles GET /api/v1/health/db.
//
// Behaviour:
//   - pool == nil (DATABASE_URL not set) → 200 {"status":"disabled","database":"not_configured"}
//   - pool OK                            → 200 {"status":"ok","database":"connected"}
//   - pool unreachable                   → 503 {"status":"unreachable","database":"unreachable"}
//
// We return 200 for "disabled" because it is an intentional configuration
// choice in development, not a server failure.
func (h *HealthHandler) DBCheck(c *fiber.Ctx) error {
	result := db.Ping(h.pool)

	status := fiber.StatusOK
	if result.Status == db.StatusUnreachable {
		status = fiber.StatusServiceUnavailable
	}

	return c.Status(status).JSON(result)
}
