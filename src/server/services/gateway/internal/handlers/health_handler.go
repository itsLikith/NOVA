package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/nova/gateway/internal/service"
)

func HealthCheck(c fiber.Ctx) error {
	status := service.HealthStatus()

	return c.Status(fiber.StatusOK).JSON(status)
}