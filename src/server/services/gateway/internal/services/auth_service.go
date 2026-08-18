package services

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"

	"github.com/nova/pkg/logger"
	"github.com/nova/pkg/response"
)

type AuthService struct {
	forward fiber.Handler
}

func NewAuthService(upstreamURL string) *AuthService {
	policy := proxy.DefaultSecurityPolicy()
	// The target is a static, service-owned Compose/internal URL rather than
	// user input, so private network addresses are required and safe here.
	policy.AllowPrivateIPs = true

	return &AuthService{
		forward: proxy.Balancer(proxy.Config{
			Servers:        []string{strings.TrimRight(upstreamURL, "/")},
			SecurityPolicy: &policy,
		}),
	}
}

func (s *AuthService) Forward(c fiber.Ctx) error {
	if err := s.forward(c); err != nil {
		logger.Error("Auth upstream request failed: " + err.Error())
		return c.Status(fiber.StatusBadGateway).JSON(
			response.SendErrorResponse(
				fiber.StatusBadGateway,
				"auth service unavailable",
				err.Error(),
			),
		)
	}

	return nil
}
