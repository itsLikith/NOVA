package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppEnv  string

	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string

	JWTSecret         string
	JWTIssuer         string
	JWTExpirationHour int

	AdminUserID   string
	AdminEmail    string
	AdminPassword string
}

func Load() (*Config, error) {
	// .env is useful locally. In production, environment variables
	// can be injected directly by the deployment environment.
	_ = godotenv.Load()

	jwtExpiration := getEnvInt("JWT_EXPIRATION_HOURS", 8)

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8081"),
		AppEnv:  getEnv("APP_ENV", "development"),

		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseUser:     getEnv("DATABASE_USER", "postgres"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", ""),
		DatabaseName:     getEnv("DATABASE_NAME", "auth"),
		DatabaseSSLMode:  getEnv("DATABASE_SSLMODE", "disable"),

		JWTSecret:         getEnv("JWT_SECRET", ""),
		JWTIssuer:         getEnv("JWT_ISSUER", "auth-service"),
		JWTExpirationHour: jwtExpiration,

		AdminUserID:   getEnv("ADMIN_USERID", ""),
		AdminEmail:    getEnv("ADMIN_EMAIL", ""),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	if c.AdminUserID == "" {
		return fmt.Errorf("ADMIN_USERID is required")
	}

	if c.AdminEmail == "" {
		return fmt.Errorf("ADMIN_EMAIL is required")
	}

	if c.AdminPassword == "" {
		return fmt.Errorf("ADMIN_PASSWORD is required")
	}

	if len(c.AdminPassword) < 8 {
		return fmt.Errorf("ADMIN_PASSWORD must be at least 8 characters")
	}

	if c.JWTExpirationHour <= 0 {
		return fmt.Errorf("JWT_EXPIRATION_HOURS must be greater than zero")
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return result
}
