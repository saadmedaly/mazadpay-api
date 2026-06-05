package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/saadmedaly/mazadpay-api/internal/config"
	"github.com/saadmedaly/mazadpay-api/internal/http/handlers"
)

// NewApp builds and returns the configured Fiber application.
// pool may be nil when DATABASE_URL is not set.
func NewApp(cfg *config.Config, pool *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "mazadpay-api",
		ErrorHandler: jsonErrorHandler,
	})

	registerMiddleware(app, cfg)
	registerRoutes(app, cfg, pool)

	return app
}

func registerMiddleware(app *fiber.App, cfg *config.Config) {
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(helmet.New())

	if cfg.AppEnv != "test" {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} ${method} ${path} ${latency} rid=${locals:requestid}\n",
		}))
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.CORSAllowedOrigins, ","),
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Authorization,X-Request-ID",
	}))
}

func registerRoutes(app *fiber.App, cfg *config.Config, pool *pgxpool.Pool) {
	health := handlers.NewHealthHandler(cfg, pool)

	// Root health check (Render uptime monitor)
	app.Get("/health", health.Check)

	// Versioned API group
	v1 := app.Group("/api/v1")
	v1.Get("/health", health.Check)
	v1.Get("/health/db", health.DBCheck)

	// 404 catch-all — must be last
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": "Route not found",
		})
	})
}

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
