package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultPort         = "8080"
	defaultJWTExpiresIn = "24h"
	defaultCORSOrigins  = "*"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	JWTExpiresIn       time.Duration
	CORSAllowedOrigins string
}

func Load() (Config, error) {
	if err := loadDotEnv(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		Port:               getEnv("PORT", defaultPort),
		DatabaseURL:        strings.TrimSpace(os.Getenv("DATABASE_URL")),
		JWTSecret:          strings.TrimSpace(os.Getenv("JWT_SECRET")),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", defaultCORSOrigins),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	expiresIn := getEnv("JWT_EXPIRES_IN", defaultJWTExpiresIn)
	duration, err := time.ParseDuration(expiresIn)
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_EXPIRES_IN: %w", err)
	}

	cfg.JWTExpiresIn = duration

	return cfg, nil
}

func ParseAllowedOrigins(value string) []string {
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))

	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	if len(origins) == 0 {
		return []string{defaultCORSOrigins}
	}

	return origins
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}
