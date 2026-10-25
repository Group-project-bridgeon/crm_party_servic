package api_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestObservabilityEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// 1. Test /health endpoint
	healthReq := httptest.NewRequest("GET", "/health", nil)
	healthResp, err := app.Test(healthReq)
	if err != nil {
		t.Fatalf("failed to call /health: %v", err)
	}
	if healthResp.StatusCode != 200 {
		t.Fatalf("expected 200 from /health, got %d", healthResp.StatusCode)
	}

	healthBody, _ := io.ReadAll(healthResp.Body)
	if !strings.Contains(string(healthBody), `"status":"UP"`) {
		t.Fatalf("expected UP status in health response, got: %s", string(healthBody))
	}

	// 2. Test /metrics endpoint
	metricsReq := httptest.NewRequest("GET", "/metrics", nil)
	metricsResp, err := app.Test(metricsReq)
	if err != nil {
		t.Fatalf("failed to call /metrics: %v", err)
	}
	if metricsResp.StatusCode != 200 {
		t.Fatalf("expected 200 from /metrics, got %d", metricsResp.StatusCode)
	}

	metricsBody, _ := io.ReadAll(metricsResp.Body)
	metricsStr := string(metricsBody)

	// Check for standard Go runtime metrics and custom party metrics
	if !strings.Contains(metricsStr, "go_goroutines") {
		t.Errorf("expected go_goroutines in /metrics output")
	}
	if !strings.Contains(metricsStr, "go_memstats_alloc_bytes") {
		t.Errorf("expected go_memstats_alloc_bytes in /metrics output")
	}
}
