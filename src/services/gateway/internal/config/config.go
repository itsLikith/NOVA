package config

import "os"

type Config struct {
	Port string

	AuthServiceURL string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:8081"
	}

	return &Config{
		Port:           port,
		AuthServiceURL: authServiceURL,
	}
}