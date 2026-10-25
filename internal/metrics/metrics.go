package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// Metrics holds Prometheus collectors for HTTP telemetry.
type Metrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestsInFlight prometheus.Gauge
}

var defaultMetrics *Metrics

// Init registers and returns default Prometheus metrics collectors.
func Init(reg prometheus.Registerer) *Metrics {
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}

	m := &Metrics{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "crm",
				Subsystem: "party_service",
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests processed, partitioned by method, route path, and status code.",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "crm",
				Subsystem: "party_service",
				Name:      "http_request_duration_seconds",
				Help:      "Duration of HTTP requests in seconds.",
				Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		),
		requestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "crm",
				Subsystem: "party_service",
				Name:      "http_requests_in_flight",
				Help:      "Current number of in-flight HTTP requests being handled.",
			},
		),
	}

	reg.MustRegister(m.requestsTotal, m.requestDuration, m.requestsInFlight)
	defaultMetrics = m
	return m
}

// Get returns the initialized metrics instance or creates a default one.
func Get() *Metrics {
	if defaultMetrics == nil {
		return Init(prometheus.DefaultRegisterer)
	}
	return defaultMetrics
}

// NewFiberMiddleware returns a Fiber middleware that measures RED metrics.
func (m *Metrics) NewFiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Skip recording metrics for internal scrape endpoint itself
		if path == "/metrics" {
			return c.Next()
		}

		m.requestsInFlight.Inc()
		defer m.requestsInFlight.Dec()

		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()

		// Use route template (/api/v1/parties/:id) instead of raw path to avoid cardinality explosion
		routePath := c.Route().Path
		if routePath == "" {
			routePath = path
		}

		status := strconv.Itoa(c.Response().StatusCode())
		method := c.Method()

		m.requestsTotal.WithLabelValues(method, routePath, status).Inc()
		m.requestDuration.WithLabelValues(method, routePath).Observe(duration)

		return err
	}
}

// Handler returns a Fiber handler that serves Prometheus metrics using the default gatherer.
func Handler() fiber.Handler {
	return HandlerFor(prometheus.DefaultGatherer)
}

// HandlerFor returns a Fiber handler for a specific Prometheus gatherer.
func HandlerFor(gatherer prometheus.Gatherer) fiber.Handler {
	promHandler := fasthttpadaptor.NewFastHTTPHandler(promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}))
	return func(c *fiber.Ctx) error {
		promHandler(c.Context())
		return nil
	}
}
