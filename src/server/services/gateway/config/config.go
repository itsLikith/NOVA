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
		logger.Warn("No .env file found; using environment variables and defaults")
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
