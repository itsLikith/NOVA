package services

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"

	"github.com/nova/pkg/response"
)

type AuthService struct {
	upstreamURL string
}

func NewAuthService(upstreamURL string) *AuthService {
	return &AuthService{
		upstreamURL: strings.TrimRight(upstreamURL, "/"),
	}
}

func (s *AuthService) Forward(c fiber.Ctx) error {

	// /api/v1/auth/login
	originalURL := c.OriginalURL()

	// Remove gateway prefix
	path := strings.TrimPrefix(
		originalURL,
		"/api/v1/auth",
	)

	// Auth service receives:
	// /login
	targetURL := s.upstreamURL + path

	if err := proxy.Do(c, targetURL); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			response.SendErrorResponse(
				404,
				"Auth service unavailabe",
				err.Error(),
			),
		)
	}

	return nil
}
