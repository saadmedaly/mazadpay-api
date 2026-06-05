package main

import (
	"context"
	"log"

	"github.com/saadmedaly/mazadpay-api/internal/config"
	"github.com/saadmedaly/mazadpay-api/internal/db"
	apphttp "github.com/saadmedaly/mazadpay-api/internal/http"
)

func main() {
	cfg := config.Load()

	// Connect to database — optional in development.
	// Server starts normally even when DATABASE_URL is not set.
	pool, err := db.Connect(context.Background(), cfg)
	if err != nil {
		log.Printf("warning: database unavailable: %v", err)
	}
	if pool != nil {
		defer pool.Close()
		log.Println("database connected")
	} else {
		log.Println("database not configured — running without DB")
	}

	app := apphttp.NewApp(cfg, pool)

	log.Printf("mazadpay-api starting — env=%s version=%s port=%s\n",
		cfg.AppEnv, cfg.AppVersion, cfg.Port)

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
