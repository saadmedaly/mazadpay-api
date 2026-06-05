package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/saadmedaly/mazadpay-api/internal/config"
	"github.com/saadmedaly/mazadpay-api/internal/http/handlers"
)

// NewApp builds and returns the configured Fiber application.
func NewApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "mazadpay-api",
		ErrorHandler: jsonErrorHandler,
	})

	registerMiddleware(app, cfg)
	registerRoutes(app, cfg)

	return app
}

// registerMiddleware attaches global middleware in the correct order.
func registerMiddleware(app *fiber.App, cfg *config.Config) {
	// Panic recovery — always first
	app.Use(recover.New())

	// Unique request ID for tracing
	app.Use(requestid.New())

	// Security headers
	app.Use(helmet.New())

	// Structured request logger (skip in test env to keep output clean)
	if cfg.AppEnv != "test" {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} ${method} ${path} ${latency} rid=${locals:requestid}\n",
		}))
	}

	// CORS — only allow configured origins, never wildcard in production
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.CORSAllowedOrigins, ","),
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Authorization,X-Request-ID",
	}))
}

// registerRoutes registers all application routes.
func registerRoutes(app *fiber.App, cfg *config.Config) {
	health := handlers.NewHealthHandler(cfg)

	// Root health check (used by Render, uptime monitors)
	app.Get("/health", health.Check)

	// Versioned API group
	v1 := app.Group("/api/v1")
	v1.Get("/health", health.Check)

	// 404 catch-all — must be last
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": "Route not found",
		})
	})
}

// jsonErrorHandler returns all Fiber errors as JSON instead of plain text.
func jsonErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error":   "internal_error",
		"message": err.Error(),
	})
}
