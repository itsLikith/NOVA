package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/nova/auth/config"
	"github.com/nova/auth/internal/routes"
)

func main() {
	cfg := config.Load()

	app := fiber.New()

	api := app.Group("/api/v1")

	routes.HealthRoutes(api)
	// routes.AuthRoutes(api, cfg)

	log.Fatal(app.Listen(":" + cfg.Port))
}
