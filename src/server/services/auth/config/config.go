package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nova/pkg/logger"
)

type Config struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		logger.Info("Error loading environment variables; using default variables")
	}
	return Config{
		Port: getenv("PORT", "8081"),

		DBHost:     getenv("DB_HOST", "localhost"),
		DBPort:     getenv("DB_PORT", "5432"),
		DBUser:     getenv("DB_USER", "postgres"),
		DBPassword: getenv("DB_PASSWORD", "postgres"),
		DBName:     getenv("DB_NAME", "auth"),
	}
}

func getenv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
