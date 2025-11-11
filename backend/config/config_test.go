package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("ADMIN_PASSWORD")
	os.Unsetenv("PORT")
	os.Unsetenv("SERVICE_ACCOUNT_PATH")
	
	cfg := LoadConfig()
	
	if cfg.AdminPassword != "admin123" {
		t.Errorf("Expected default admin password 'admin123', got '%s'", cfg.AdminPassword)
	}
	if cfg.Port != "8080" {
		t.Errorf("Expected default port '8080', got '%s'", cfg.Port)
	}
	if cfg.ServiceAccountPath != "./service-account.json" {
		t.Errorf("Expected default service account path './service-account.json', got '%s'", cfg.ServiceAccountPath)
	}
}

func TestLoadConfig_FromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("ADMIN_PASSWORD", "testpassword")
	os.Setenv("PORT", "3000")
	os.Setenv("SERVICE_ACCOUNT_PATH", "/path/to/service.json")
	defer func() {
		os.Unsetenv("ADMIN_PASSWORD")
		os.Unsetenv("PORT")
		os.Unsetenv("SERVICE_ACCOUNT_PATH")
	}()
	
	cfg := LoadConfig()
	
	if cfg.AdminPassword != "testpassword" {
		t.Errorf("Expected admin password 'testpassword', got '%s'", cfg.AdminPassword)
	}
	if cfg.Port != "3000" {
		t.Errorf("Expected port '3000', got '%s'", cfg.Port)
	}
	if cfg.ServiceAccountPath != "/path/to/service.json" {
		t.Errorf("Expected service account path '/path/to/service.json', got '%s'", cfg.ServiceAccountPath)
	}
}

func TestLoadConfig_PartialEnv(t *testing.T) {
	// Set only PORT
	os.Setenv("PORT", "9000")
	os.Unsetenv("ADMIN_PASSWORD")
	os.Unsetenv("SERVICE_ACCOUNT_PATH")
	defer os.Unsetenv("PORT")
	
	cfg := LoadConfig()
	
	if cfg.Port != "9000" {
		t.Errorf("Expected port '9000', got '%s'", cfg.Port)
	}
	if cfg.AdminPassword != "admin123" {
		t.Errorf("Expected default admin password 'admin123', got '%s'", cfg.AdminPassword)
	}
}

