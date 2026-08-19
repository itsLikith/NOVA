package middleware

import (
	"strings"

	"github.com/nova/auth/config"
	"github.com/nova/auth/internal/models"
	"github.com/nova/auth/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

const (
	UserIDKey = "userID"
	RoleKey   = "role"
)

func JWTAuth(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")

		var tokenString string
		if strings.TrimSpace(header) == "" {
			// try cookie fallback
			tokenString = c.Cookies("token")
			if strings.TrimSpace(tokenString) == "" {
				return fiber.NewError(
					fiber.StatusUnauthorized,
					"authentication required",
				)
			}
		} else {
			var err error
			tokenString, err = utils.ExtractBearerToken(header)
			if err != nil {
				return fiber.NewError(
					fiber.StatusUnauthorized,
					"invalid authorization header",
				)
			}
		}

		claims, err := utils.ParseToken(
			tokenString,
			cfg.JWTSecret,
			cfg.JWTIssuer,
		)

		if err != nil {
			return fiber.NewError(
				fiber.StatusUnauthorized,
				"invalid or expired token",
			)
		}

		if claims.Subject == "" {
			return fiber.NewError(
				fiber.StatusUnauthorized,
				"invalid token claims",
			)
		}

		if claims.Role != models.RoleAdmin &&
			claims.Role != models.RoleUser {
			return fiber.NewError(
				fiber.StatusUnauthorized,
				"invalid token role",
			)
		}

		c.Locals(UserIDKey, claims.Subject)
		c.Locals(RoleKey, claims.Role)

		return c.Next()
	}
}

func RequireAdmin() fiber.Handler {
	return func(c fiber.Ctx) error {
		role, ok := c.Locals(RoleKey).(string)

		if !ok || role != models.RoleAdmin {
			return fiber.NewError(
				fiber.StatusForbidden,
				"admin privileges required",
			)
		}

		return c.Next()
	}
}
