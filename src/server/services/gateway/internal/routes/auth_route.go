package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nova/gateway/config"
	"github.com/nova/gateway/internal/handlers"
	"github.com/nova/gateway/internal/services"
)

func AuthRoutes(
	router fiber.Router,
	cfg config.Config,
) {
	authService := services.NewAuthService(
		cfg.AuthServiceURL,
	)

	authHandler := handlers.NewAuthHandler(
		authService,
	)

	auth := router.Group("/auth")

	auth.All("/*", authHandler.Forward)
}
