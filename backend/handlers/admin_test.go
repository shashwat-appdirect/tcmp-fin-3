package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAdminLogin_Success(t *testing.T) {
	setupTestClient(t)
	// Set admin password
	SetAdminPassword("testpassword123")
	
	reqBody := LoginRequest{
		Password: "testpassword123",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	AdminLogin(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Message != "Login successful" {
		t.Errorf("Expected message 'Login successful', got '%s'", response.Message)
	}
}

func TestAdminLogin_InvalidPassword(t *testing.T) {
	setupTestClient(t)
	SetAdminPassword("correctpassword")
	
	reqBody := LoginRequest{
		Password: "wrongpassword",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	AdminLogin(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
	
	var response LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if response.Success {
		t.Error("Expected success to be false")
	}
	if response.Message != "Invalid password" {
		t.Errorf("Expected message 'Invalid password', got '%s'", response.Message)
	}
}

func TestAdminLogin_FromEnv(t *testing.T) {
	setupTestClient(t)
	// Clear set password
	SetAdminPassword("")
	
	// Set environment variable
	os.Setenv("ADMIN_PASSWORD", "envpassword123")
	defer os.Unsetenv("ADMIN_PASSWORD")
	
	reqBody := LoginRequest{
		Password: "envpassword123",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	AdminLogin(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAdminLogin_DefaultPassword(t *testing.T) {
	setupTestClient(t)
	// Clear set password and env
	SetAdminPassword("")
	os.Unsetenv("ADMIN_PASSWORD")
	
	reqBody := LoginRequest{
		Password: "admin123",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	AdminLogin(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAdminLogin_InvalidJSON(t *testing.T) {
	setupTestClient(t)
	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	AdminLogin(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAdminLogin_WrongMethod(t *testing.T) {
	setupTestClient(t)
	req := httptest.NewRequest("GET", "/api/admin/login", nil)
	w := httptest.NewRecorder()
	
	AdminLogin(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

