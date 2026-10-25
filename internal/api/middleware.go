package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/group-project/crm_party_service/internal/logger"
	"github.com/group-project/crm_party_service/internal/metrics"
)

// SetupMiddlewares configures standard HTTP middlewares (recover, requestid, structured slog logger, RED metrics, cors).
func SetupMiddlewares(app *fiber.App) {
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.NewFiberMiddleware())
	app.Use(metrics.Get().NewFiberMiddleware())
	app.Use(cors.New())
}
