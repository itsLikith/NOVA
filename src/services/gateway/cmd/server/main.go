package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/nova/gateway/internal/config"
	"github.com/nova/gateway/internal/proxy"
	"github.com/nova/gateway/internal/handlers"
)

func main() {
	cfg := config.Load()

	app := fiber.New()
	app.Use(logger.New())

	app.Get("/api/v1/health", handlers.Health)

	app.All(
		"/api/v1/auth/*",
		proxy.Forward(cfg.AuthServiceURL),
	)

	log.Fatal(
		app.Listen(":" + cfg.Port),
	)
}