package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port        string
	APIUsername string
	APIPassword string
}

// LoadConfig loads configuration
func LoadConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	apiUsername := os.Getenv("API_USERNAME")
	if apiUsername == "" {
		apiUsername = "admin"
	}

	apiPassword := os.Getenv("API_PASSWORD")
	if apiPassword == "" {
		apiPassword = "admin123"
	}

	return &Config{
		Port:        port,
		APIUsername: apiUsername,
		APIPassword: apiPassword,
	}
}
