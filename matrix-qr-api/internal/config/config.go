package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	Port          string
	JWTSecret     []byte
	TokenTTL      time.Duration
	StatsURL      string
	StatsTimeout  time.Duration
	AuthUsername  string
	AuthPassword  string
	AllowedOrigin string
}

func Load() (Config, error) {
	cfg := Config{
		Port:          env("PORT", "8080"),
		JWTSecret:     []byte(os.Getenv("JWT_SECRET")),
		StatsURL:      env("STATS_API_URL", "http://localhost:3000"),
		AuthUsername:  env("AUTH_USERNAME", "admin"),
		AuthPassword:  os.Getenv("AUTH_PASSWORD"),
		AllowedOrigin: env("ALLOWED_ORIGIN", "*"),
	}

	if len(cfg.JWTSecret) < 16 {
		return cfg, errors.New("JWT_SECRET must be set and have at least 16 characters")
	}
	if cfg.AuthPassword == "" {
		return cfg, errors.New("AUTH_PASSWORD must be set")
	}

	var err error
	if cfg.TokenTTL, err = time.ParseDuration(env("TOKEN_TTL", "1h")); err != nil {
		return cfg, errors.New("TOKEN_TTL must be a valid duration")
	}
	if cfg.StatsTimeout, err = time.ParseDuration(env("STATS_TIMEOUT", "5s")); err != nil {
		return cfg, errors.New("STATS_TIMEOUT must be a valid duration")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
