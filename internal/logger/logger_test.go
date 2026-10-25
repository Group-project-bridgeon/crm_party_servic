package logger

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func TestInit_JSON_Production(t *testing.T) {
	var buf bytes.Buffer
	l := Init(Config{
		Environment: "production",
		Level:       "INFO",
		ServiceName: "test_service",
		Output:      &buf,
	})

	l.Info("test message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, `"msg":"test message"`) {
		t.Fatalf("expected json log output with message, got: %s", output)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("failed to parse json log: %v", err)
	}

	if parsed["service"] != "test_service" {
		t.Errorf("expected service=test_service, got: %v", parsed["service"])
	}
	if parsed["env"] != "production" {
		t.Errorf("expected env=production, got: %v", parsed["env"])
	}
}

func TestInit_Text_Development(t *testing.T) {
	var buf bytes.Buffer
	l := Init(Config{
		Environment: "development",
		Level:       "DEBUG",
		ServiceName: "test_dev_service",
		Output:      &buf,
	})

	l.Debug("debug message")

	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Fatalf("expected text log with debug message, got: %s", output)
	}
}

func TestFiberMiddleware_Logging(t *testing.T) {
	var buf bytes.Buffer
	Init(Config{
		Environment: "production",
		Level:       "DEBUG",
		ServiceName: "test_service",
		Output:      &buf,
	})

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(NewFiberMiddleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("healthy")
	})

	// Test normal endpoint
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	output := buf.String()
	if !strings.Contains(output, `http request completed`) {
		t.Fatalf("expected log entry for /test request, got: %s", output)
	}
	if !strings.Contains(output, `"path":"/test"`) {
		t.Fatalf("expected path in log entry, got: %s", output)
	}
}
