package main

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nova/gateway/config"
	"github.com/nova/gateway/internal/routes"
	"github.com/nova/pkg/logger"
	"github.com/nova/pkg/response"
)

func main() {
	cfg := config.Load()

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

	api := app.Group("/api/v1")

	routes.HealthRoutes(api)
	routes.AuthRoutes(api, cfg)

	logger.Info("Gateway service starting on port " + cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		logger.Fatal("Server error: " + err.Error())
	}
}
