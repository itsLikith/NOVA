package routes

import (
	"github.com/nova/auth/config"
	"github.com/nova/auth/internal/handlers"
	"github.com/nova/auth/pkg/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(
	app *fiber.App,
	handler *handlers.AuthHandler,
	cfg *config.Config,
) {
	auth := app.Group("/api/v1/auth")

	auth.Post("/login", handler.Login)

	admin := auth.Group(
		"/users",
		middleware.JWTAuth(cfg),
		middleware.RequireAdmin(),
	)

	admin.Post("/", handler.CreateUser)
}
