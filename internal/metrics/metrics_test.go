package metrics

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMetrics_FiberMiddlewareAndHandler(t *testing.T) {
	reg := prometheus.NewRegistry()
	m := Init(reg)

	app := fiber.New()
	app.Use(m.NewFiberMiddleware())

	app.Get("/test-endpoint", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Get("/metrics", HandlerFor(reg))

	// Make a test request to trigger metric increment
	req := httptest.NewRequest("GET", "/test-endpoint", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}

	// Scrape /metrics endpoint
	metricsReq := httptest.NewRequest("GET", "/metrics", nil)
	metricsResp, err := app.Test(metricsReq)
	if err != nil {
		t.Fatalf("failed to scrape metrics: %v", err)
	}
	if metricsResp.StatusCode != 200 {
		t.Fatalf("expected status 200 for metrics, got: %d", metricsResp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(metricsResp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	body := string(bodyBytes)

	if !strings.Contains(body, "crm_party_service_http_requests_total") {
		t.Errorf("expected crm_party_service_http_requests_total in metrics output")
	}
	if !strings.Contains(body, `path="/test-endpoint"`) {
		t.Errorf("expected path label in metrics output")
	}
}
