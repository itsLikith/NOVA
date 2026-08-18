package handlers

import (
	"github.com/nova/auth/internal/services"
	"github.com/nova/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type HealthHandler struct {
	service services.HealthService
}

func NewHealthHandler(service services.HealthService) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

func (h *HealthHandler) Health(c fiber.Ctx) error {
	if err := h.service.Check(c.Context()); err != nil {
		return fiber.NewError(
			fiber.StatusServiceUnavailable,
			"service unavailable",
		)
	}

	return c.Status(fiber.StatusOK).JSON(response.SendSuccessResponse(
		fiber.StatusOK,
		"Auth service is UP",
		fiber.Map{"service": "auth"},
	))
}
