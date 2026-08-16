package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/nova/auth/internal/service"
	"github.com/nova/auth/pkg/response"
)

func AuthRequired(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := setClaims(c, secret); err != nil {
			return err
		}
		return c.Next()
	}
}

func AdminRequired(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := setClaims(c, secret); err != nil {
			return err
		}

		claims, ok := c.Locals("user").(*service.Claims)
		if !ok {
			return response.ErrorMessage(c, fiber.StatusUnauthorized, "authentication required", "user not found in token")
		}
		if claims.Role != "admin" {
			return response.ErrorMessage(c, fiber.StatusForbidden, "access denied", "admin access required")
		}

		return c.Next()
	}
}

func setClaims(c *fiber.Ctx, secret string) error {
	tokenString, err := service.ExtractBearerToken(c.Get("Authorization"))
	if err != nil {
		message := "invalid token"
		if errors.Is(err, service.ErrMissingToken) {
			message = "missing authorization header"
		}
		return response.ErrorMessage(c, fiber.StatusUnauthorized, "authentication required", message)
	}

	claims, err := service.ValidateJWT(secret, tokenString)
	if err != nil {
		return response.ErrorMessage(c, fiber.StatusUnauthorized, "authentication required", "invalid token")
	}

	c.Locals("user", claims)
	return nil
}
