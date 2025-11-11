package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AdminPassword      string
	Port               string
	ServiceAccountPath string
	FrontendDir        string
}

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123" // Default for development
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serviceAccountPath := os.Getenv("SERVICE_ACCOUNT_PATH")
	if serviceAccountPath == "" {
		serviceAccountPath = "./service-account.json"
	}

	frontendDir := os.Getenv("FRONTEND_DIR")
	// Default is empty, will use "./frontend/dist" in main.go if not set

	return &Config{
		AdminPassword:      adminPassword,
		Port:               port,
		ServiceAccountPath: serviceAccountPath,
		FrontendDir:        frontendDir,
	}
}

