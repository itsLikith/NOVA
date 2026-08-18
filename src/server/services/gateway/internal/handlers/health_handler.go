package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nova/gateway/internal/services"
)

func HealthCheck(c fiber.Ctx) error {
	status := services.HealthStatus()

	return c.Status(fiber.StatusOK).JSON(status)
}
