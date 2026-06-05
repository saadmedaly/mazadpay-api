package main

import (
	"log"

	"github.com/saadmedaly/mazadpay-api/internal/config"
	apphttp "github.com/saadmedaly/mazadpay-api/internal/http"
)

func main() {
	cfg := config.Load()

	app := apphttp.NewApp(cfg)

	addr := ":" + cfg.Port
	log.Printf("mazadpay-api starting — env=%s version=%s port=%s\n",
		cfg.AppEnv, cfg.AppVersion, cfg.Port)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
