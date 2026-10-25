package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// registerHealthRoutes registers health and status monitoring endpoints.
func registerHealthRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "UP",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})
}
