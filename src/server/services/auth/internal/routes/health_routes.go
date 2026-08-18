package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/nova/auth/internal/handlers"
)

func HealthRoutes(router fiber.Router) {
	router.Get("/health", handlers.HealthCheck)
}