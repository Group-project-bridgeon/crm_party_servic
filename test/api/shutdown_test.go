package api_test

import (
	"context"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestGracefulShutdown_InFlightRequestsComplete(t *testing.T) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
	})

	// Add an endpoint simulating an in-flight operation
	app.Get("/slow", func(c *fiber.Ctx) error {
		time.Sleep(100 * time.Millisecond)
		return c.SendString("completed successfully")
	})

	// Bind to an ephemeral port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to bind ephemeral listener: %v", err)
	}
	addr := ln.Addr().String()

	// Start server in background
	serverStarted := make(chan struct{})
	go func() {
		close(serverStarted)
		_ = app.Listener(ln)
	}()

	<-serverStarted
	time.Sleep(20 * time.Millisecond)

	var wg sync.WaitGroup
	var respStatus int
	var clientErr error

	// Send an in-flight request
	wg.Add(1)
	go func() {
		defer wg.Done()
		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get("http://" + addr + "/slow")
		if err != nil {
			clientErr = err
			return
		}
		defer resp.Body.Close()
		respStatus = resp.StatusCode
	}()

	// Wait slightly so request is in-flight
	time.Sleep(30 * time.Millisecond)

	// Initiate graceful shutdown with 1 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	shutdownErr := app.ShutdownWithContext(ctx)
	if shutdownErr != nil {
		t.Fatalf("error during graceful shutdown: %v", shutdownErr)
	}

	wg.Wait()

	if clientErr != nil {
		t.Fatalf("in-flight request failed during shutdown: %v", clientErr)
	}
	if respStatus != http.StatusOK {
		t.Fatalf("expected in-flight request status 200, got %d", respStatus)
	}
}
