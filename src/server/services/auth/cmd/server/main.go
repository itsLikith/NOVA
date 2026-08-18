package main

import (
	"context"

	"github.com/nova/auth/config"
	"github.com/nova/auth/internal/handlers"
	"github.com/nova/auth/internal/repository"
	"github.com/nova/auth/internal/routes"
	"github.com/nova/auth/internal/services"
	"github.com/nova/auth/pkg/database"
	"github.com/nova/pkg/logger"
	"github.com/nova/pkg/response"

	"github.com/gofiber/fiber/v3"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Configuration error: " + err.Error())
	}

	db, err := database.Connect(cfg)
	if err != nil {
		logger.Fatal("Database error: " + err.Error())
	}

	if err := database.Migrate(db); err != nil {
		logger.Fatal("Migration error: " + err.Error())
	}

	authRepository := repository.NewAuthRepository(db)

	authService := services.NewAuthService(
		authRepository,
		cfg,
	)

	if err := authService.EnsureAdmin(context.Background()); err != nil {
		logger.Fatal("Admin bootstrap error: " + err.Error())
	}

	healthService := services.NewHealthService(db)

	authHandler := handlers.NewAuthHandler(
		authService,
	)

	healthHandler := handlers.NewHealthHandler(
		healthService,
	)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			message := "internal server error"

			if fiberErr, ok := err.(*fiber.Error); ok {
				status = fiberErr.Code
				message = fiberErr.Message
			}

			logger.Warn("Request failed: " + err.Error())
			return c.Status(status).JSON(response.SendErrorResponse(status, message, nil))
		},
	})

	app.Use(recoverer.New())

	routes.RegisterAuthRoutes(
		app,
		authHandler,
		cfg,
	)

	routes.RegisterHealthRoutes(
		app,
		healthHandler,
	)

	logger.Info("Auth service starting on port " + cfg.AppPort)

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		logger.Fatal("Server error: " + err.Error())
	}
}
