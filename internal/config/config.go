package config

import (
	"os"
	"strings"
)

type Config struct {
	Port        string
	DatabaseURL string
	Environment string
}

func Load() *Config {
	port := getEnv("PORT", "8080")
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	databaseURL := getEnv(
		"DATABASE_URL",
		"postgres://postgres:postgres@localhost:5432/party?sslmode=disable",
	)

	environment := getEnv("ENV", "development")

	return &Config{
		Port:        port,
		DatabaseURL: databaseURL,
		Environment: environment,
	}
}

func getEnv(key, fallback string) string {
	val, exists := os.LookupEnv(key)
	if !exists || strings.TrimSpace(val) == "" {
		return fallback
	}
	return strings.TrimSpace(val)
}
