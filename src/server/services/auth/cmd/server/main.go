package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nova/auth/config"
	"github.com/nova/auth/internal/handler"
	"github.com/nova/auth/internal/model"
	"github.com/nova/auth/internal/repository"
	"github.com/nova/auth/internal/service"
	"github.com/nova/auth/pkg/database"
)

func main() {
	cfg := config.New()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("database migration failed: ", err)
	}

	repo := repository.NewPostgreSQLRepository(db)
	authService := service.NewService(repo, cfg.JWTSecret, cfg.AdminUsername, cfg.AdminPassword)

	if err := authService.SeedAdmin(); err != nil {
		log.Fatal("admin seed failed: ", err)
	}

	app := fiber.New()
	h := handler.NewHandler(authService, cfg.JWTSecret)
	h.RegisterRoutes(app)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
