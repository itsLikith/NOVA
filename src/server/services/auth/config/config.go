package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	AppPort       string
	JWTSecret     string
	AdminUsername string
	AdminPassword string
}

func New() *Config {
	godotenv.Load()
	return &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		AppPort:       os.Getenv("APP_PORT"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		AdminUsername: os.Getenv("ADMIN_USERNAME"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.AdminUsername) == "" {
		return fmt.Errorf("ADMIN_USERNAME is required")
	}
	if strings.TrimSpace(c.AdminPassword) == "" {
		return fmt.Errorf("ADMIN_PASSWORD is required")
	}
	if strings.TrimSpace(c.JWTSecret) == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
	)
}

func (c *Config) SystemDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=UTC",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBPort,
	)
}