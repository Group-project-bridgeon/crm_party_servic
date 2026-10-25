package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("ENV")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("SHUTDOWN_TIMEOUT_SECONDS")

	cfg := Load()

	if cfg.Port != ":8080" {
		t.Fatalf("expected :8080, got %s", cfg.Port)
	}

	if cfg.Environment != "development" {
		t.Fatalf("expected development, got %s", cfg.Environment)
	}

	if cfg.LogLevel != "INFO" {
		t.Fatalf("expected INFO, got %s", cfg.LogLevel)
	}

	if cfg.DatabaseURL == "" {
		t.Fatal("expected default DatabaseURL to be non-empty")
	}

	if cfg.ShutdownTimeout != 10*time.Second {
		t.Fatalf("expected default 10s shutdown timeout, got %v", cfg.ShutdownTimeout)
	}
}

func TestLoad_CustomEnv(t *testing.T) {
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgres://custom:custom@localhost:5432/custom_db")
	os.Setenv("ENV", "production")
	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "25")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("ENV")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("SHUTDOWN_TIMEOUT_SECONDS")
	}()

	cfg := Load()

	if cfg.Port != ":3000" {
		t.Fatalf("expected :3000, got %s", cfg.Port)
	}

	if cfg.Environment != "production" {
		t.Fatalf("expected production, got %s", cfg.Environment)
	}

	if cfg.LogLevel != "DEBUG" {
		t.Fatalf("expected DEBUG, got %s", cfg.LogLevel)
	}

	if cfg.DatabaseURL != "postgres://custom:custom@localhost:5432/custom_db" {
		t.Fatalf("expected custom db url, got %s", cfg.DatabaseURL)
	}

	if cfg.ShutdownTimeout != 25*time.Second {
		t.Fatalf("expected 25s shutdown timeout, got %v", cfg.ShutdownTimeout)
	}
}
