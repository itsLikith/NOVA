package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nova/pkg/logger"
)

type Config struct {
	Port string

	AuthServiceURL string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		logger.Info("Error loading environment variables; using default variables")
	}

	port := getenv("PORT", "8080")

	authServiceURL := getenv("AUTH_SERVICE_URL", "http://localhost:8081")

	return Config{
		Port:           port,
		AuthServiceURL: authServiceURL,
	}
}

func getenv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
