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

	// validate endpoint uses JWTAuth and will return user information if the token is valid
	auth.Get("/validate", middleware.JWTAuth(cfg), handler.Validate)

	admin := auth.Group(
		"/users",
		middleware.JWTAuth(cfg),
		middleware.RequireAdmin(),
	)

	admin.Post("/", handler.CreateUser)
}
