package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port            string
	DatabaseURL     string
	Environment     string
	LogLevel        string
	ShutdownTimeout time.Duration
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
	logLevel := getEnv("LOG_LEVEL", "INFO")

	shutdownTimeoutSecondsStr := getEnv("SHUTDOWN_TIMEOUT_SECONDS", "10")
	timeoutSec, err := strconv.Atoi(shutdownTimeoutSecondsStr)
	if err != nil || timeoutSec <= 0 {
		timeoutSec = 10
	}

	return &Config{
		Port:            port,
		DatabaseURL:     databaseURL,
		Environment:     environment,
		LogLevel:        logLevel,
		ShutdownTimeout: time.Duration(timeoutSec) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	val, exists := os.LookupEnv(key)
	if !exists || strings.TrimSpace(val) == "" {
		return fallback
	}
	return strings.TrimSpace(val)
}
