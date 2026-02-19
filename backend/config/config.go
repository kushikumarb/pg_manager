package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port        string
	DBUser      string
	DBPass      string
	DBName      string
	DBHost      string
	DBPort      string
	JWTSecret   string
	Environment string
	BaseURL     string

	// Twilio Credentials
	TwilioSID        string
	TwilioAuthToken  string
	TwilioFromNumber string
}

var App AppConfig

func LoadConfig() {
	// 1. Try to load .env for local development.
	// In Docker, this will fail because variables are injected via docker-compose.
	// That is expected and handled by the error check.
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file not found. Relying on system environment variables.")
	}

	App = AppConfig{
		// Format: getEnv("KEY", "DEFAULT_VALUE")
		Port:        getEnv("PORT", "8080"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPass:      getEnv("DB_PASS", ""),
		DBName:      getEnv("DB_NAME", "pg_management"),
		DBHost:      getEnv("DB_HOST", "db"),   // Default to "db" for Docker network
		DBPort:      getEnv("DB_PORT", "5432"), // Default Postgres port
		JWTSecret:   getEnv("JWT_SECRET", "placeholder_for_dev_only"),
		Environment: getEnv("APP_ENV", "development"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),

		TwilioSID:        getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: getEnv("TWILIO_FROM_NUMBER", ""),
	}

	// 2. Production Security Warnings
	if App.Environment == "production" {
		// Updated to match your new placeholder
		if App.JWTSecret == "placeholder_for_dev_only" || App.JWTSecret == "" {
			log.Println("🚨 CRITICAL: Using default or empty JWT secret in production! Update your .env immediately.")
		}
		// Updated to warn if the password is empty in production
		if App.DBPass == "" {
			log.Println("⚠️ WARNING: No database password set in production mode.")
		}
	}

	log.Printf("✅ Configuration loaded successfully for [%s] mode", App.Environment)
}

// getEnv checks if an environment variable exists, otherwise returns a fallback value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
