package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nova/gateway/internal/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(
	service *services.AuthService,
) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Forward(c fiber.Ctx) error {
	return h.service.Forward(c)
}
