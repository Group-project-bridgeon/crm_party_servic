package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("ENV")

	cfg := Load()

	if cfg.Port != ":8080" {
		t.Fatalf("expected :8080, got %s", cfg.Port)
	}

	if cfg.Environment != "development" {
		t.Fatalf("expected development, got %s", cfg.Environment)
	}

	if cfg.DatabaseURL == "" {
		t.Fatal("expected default DatabaseURL to be non-empty")
	}
}

func TestLoad_CustomEnv(t *testing.T) {
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgres://custom:custom@localhost:5432/custom_db")
	os.Setenv("ENV", "production")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("ENV")
	}()

	cfg := Load()

	if cfg.Port != ":3000" {
		t.Fatalf("expected :3000, got %s", cfg.Port)
	}

	if cfg.Environment != "production" {
		t.Fatalf("expected production, got %s", cfg.Environment)
	}

	if cfg.DatabaseURL != "postgres://custom:custom@localhost:5432/custom_db" {
		t.Fatalf("expected custom db url, got %s", cfg.DatabaseURL)
	}
}
