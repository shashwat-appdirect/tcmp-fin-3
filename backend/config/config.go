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
	GCPProjectID       string
	ClientID           string
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

	gcpProjectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if gcpProjectID == "" {
		gcpProjectID = os.Getenv("GCP_PROJECT")
	}
	// Cloud Run automatically sets GOOGLE_CLOUD_PROJECT, but we can also use GCP_PROJECT

	clientID := os.Getenv("CLIENT_ID")
	// CLIENT_ID is required when using ADC (Cloud Run), optional when using service account file

	return &Config{
		AdminPassword:      adminPassword,
		Port:               port,
		ServiceAccountPath: serviceAccountPath,
		FrontendDir:        frontendDir,
		GCPProjectID:       gcpProjectID,
		ClientID:           clientID,
	}
}

