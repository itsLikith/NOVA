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

	// Log presence of Set-Cookie header for debugging in dev
	if sc := c.Response().Header.Peek("Set-Cookie"); len(sc) > 0 {
		logger.Info("Forwarded Set-Cookie from auth upstream: " + string(sc))
		// Also log request cookies for debugging (if any)
		if rc := c.Request().Header.Peek("Cookie"); len(rc) > 0 {
			logger.Info("Request Cookie header present when forwarding auth response")
		} else {
			logger.Info("No request Cookie header present when forwarding auth response")
		}
	} else {
		logger.Info("No Set-Cookie header present in upstream response")
	}

	return nil
}
