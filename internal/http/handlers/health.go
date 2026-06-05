package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saadmedaly/mazadpay-api/internal/config"
)

// HealthHandler holds dependencies for health endpoints.
type HealthHandler struct {
	cfg *config.Config
}

// NewHealthHandler creates a HealthHandler with the provided config.
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

// HealthResponse is the JSON body returned by health endpoints.
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
