package routes

import (
	"github.com/nova/auth/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

func RegisterHealthRoutes(
	app *fiber.App,
	handler *handlers.HealthHandler,
) {
	app.Get("/api/v1/auth/health", handler.Health)
}
