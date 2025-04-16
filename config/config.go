package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppEnv     string
	AppPort    string
}

// NewConfig creates a new Config instance
func NewConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "borderless"),
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8006"),
	}
}

// loadEnvFile loads the correct .env file based on APP_ENV
func loadEnvFile() {
	appEnv := os.Getenv("APP_ENV")

	switch appEnv {
	case "production":
		_ = godotenv.Load("production.env")
	default:
		_ = godotenv.Load(".env") // default to local
	}
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
