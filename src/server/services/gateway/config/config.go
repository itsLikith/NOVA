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
	return Config{
		Port:           getenv("PORT", "8080"),
		AuthServiceURL: getenv("AUTH_SERVICE_URL", "http://localhost:8081"),
	}
}

func getenv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
