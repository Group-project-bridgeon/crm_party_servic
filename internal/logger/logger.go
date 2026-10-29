package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

type contextKey string

const requestIDKey contextKey = "request_id"

var defaultLogger *slog.Logger

// Config holds configuration parameters for setting up the structured logger.
type Config struct {
	Environment string // "production", "staging", "development"
	Level       string // "DEBUG", "INFO", "WARN", "ERROR"
	ServiceName string // e.g. "crm_party_service"
	Output      io.Writer
}

// Init configures and sets the global default slog logger based on the environment.
func Init(cfg Config) *slog.Logger {
	out := cfg.Output
	if out == nil {
		out = os.Stdout
	}

	var lvl slog.Level
	switch strings.ToUpper(cfg.Level) {
	case "DEBUG":
		lvl = slog.LevelDebug
	case "WARN":
		lvl = slog.LevelWarn
	case "ERROR":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: strings.ToLower(cfg.Environment) != "production",
	}

	var handler slog.Handler
	if strings.ToLower(cfg.Environment) == "production" || strings.ToLower(cfg.Environment) == "staging" {
		handler = slog.NewJSONHandler(out, opts)
	} else {
		handler = slog.NewTextHandler(out, opts)
	}

	serviceName := cfg.ServiceName
	if serviceName == "" {
		serviceName = "crm_party_service"
	}

	// Always attach service and env metadata to all logs
	handler = handler.WithAttrs([]slog.Attr{
		slog.String("service", serviceName),
		slog.String("env", cfg.Environment),
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
	defaultLogger = logger

	return logger
}

// L returns the active default slog.Logger.
func L() *slog.Logger {
	if defaultLogger == nil {
		return slog.Default()
	}
	return defaultLogger
}

// WithRequestID returns a new context with the given request ID embedded.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// FromContext extracts the request ID from context if available.
func FromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}
