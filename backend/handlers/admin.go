package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var AdminPassword string

// SetAdminPassword sets the admin password from config
func SetAdminPassword(password string) {
	AdminPassword = password
}

// AdminLogin handles admin authentication
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use AdminPassword from config, fallback to env if not set
	adminPassword := AdminPassword
	if adminPassword == "" {
		adminPassword = os.Getenv("ADMIN_PASSWORD")
	}
	if adminPassword == "" {
		adminPassword = "admin123" // Default for development
	}

	if req.Password != adminPassword {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid password",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Message: "Login successful",
	})
}

