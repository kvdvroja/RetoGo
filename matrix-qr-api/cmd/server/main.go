package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"matrix-qr-api/internal/config"
	"matrix-qr-api/internal/handler"
	"matrix-qr-api/internal/stats"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app := handler.New(handler.Options{
		JWTSecret:     cfg.JWTSecret,
		TokenTTL:      cfg.TokenTTL,
		Username:      cfg.AuthUsername,
		Password:      cfg.AuthPassword,
		AllowedOrigin: cfg.AllowedOrigin,
		Stats:         stats.NewHTTPClient(cfg.StatsURL, cfg.StatsTimeout),
		Logging:       true,
	})

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		_ = app.Shutdown()
	}()

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
