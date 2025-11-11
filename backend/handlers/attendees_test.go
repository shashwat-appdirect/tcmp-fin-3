package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"event-registration-backend/models"
)

func TestRegisterAttendee_Success(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	reqBody := models.RegisterRequest{
		FullName:    "John Doe",
		Email:       "john@example.com",
		Designation: "Developer",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/attendees/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	RegisterAttendee(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var attendee models.Attendee
	if err := json.Unmarshal(w.Body.Bytes(), &attendee); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if attendee.FullName != reqBody.FullName {
		t.Errorf("Expected FullName %s, got %s", reqBody.FullName, attendee.FullName)
	}
	if attendee.Email != reqBody.Email {
		t.Errorf("Expected Email %s, got %s", reqBody.Email, attendee.Email)
	}
	if attendee.ID == "" {
		t.Error("Expected attendee ID to be set")
	}
}

func TestRegisterAttendee_MissingFields(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	testCases := []struct {
		name string
		body models.RegisterRequest
	}{
		{"Missing FullName", models.RegisterRequest{Email: "test@example.com", Designation: "Dev"}},
		{"Missing Email", models.RegisterRequest{FullName: "Test User", Designation: "Dev"}},
		{"Missing Designation", models.RegisterRequest{FullName: "Test User", Email: "test@example.com"}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest("POST", "/api/attendees/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			RegisterAttendee(w, req)
			
			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestRegisterAttendee_DuplicateEmail(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	// Register first attendee
	reqBody := models.RegisterRequest{
		FullName:    "John Doe",
		Email:       "duplicate@example.com",
		Designation: "Developer",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/attendees/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	RegisterAttendee(w, req)
	
	if w.Code != http.StatusOK {
		t.Skipf("Skipping duplicate test: first registration failed (likely no Firestore): %d", w.Code)
		return
	}
	
	// Try to register same email again
	req2 := httptest.NewRequest("POST", "/api/attendees/register", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	RegisterAttendee(w2, req2)
	
	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w2.Code)
	}
}

func TestRegisterAttendee_InvalidJSON(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("POST", "/api/attendees/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	RegisterAttendee(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetAttendeeCount_Empty(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("GET", "/api/attendees/count", nil)
	w := httptest.NewRecorder()
	
	GetAttendeeCount(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var result map[string]int
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if result["count"] != 0 {
		t.Errorf("Expected count 0, got %d", result["count"])
	}
}

func TestGetAttendeeCount_WithData(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	// Register an attendee first
	reqBody := models.RegisterRequest{
		FullName:    "Test User",
		Email:       "counttest@example.com",
		Designation: "Tester",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/attendees/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	RegisterAttendee(w, req)
	
	if w.Code != http.StatusOK {
		t.Skipf("Skipping count test: registration failed (likely no Firestore): %d", w.Code)
		return
	}
	
	// Get count
	req2 := httptest.NewRequest("GET", "/api/attendees/count", nil)
	w2 := httptest.NewRecorder()
	GetAttendeeCount(w2, req2)
	
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}
	
	var result map[string]int
	if err := json.Unmarshal(w2.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if result["count"] < 1 {
		t.Errorf("Expected count >= 1, got %d", result["count"])
	}
}

func TestGetAttendees_Empty(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("GET", "/api/admin/attendees", nil)
	w := httptest.NewRecorder()
	
	GetAttendees(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var attendees []models.Attendee
	if err := json.Unmarshal(w.Body.Bytes(), &attendees); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if attendees == nil {
		t.Error("Expected empty array, got nil")
	}
	if len(attendees) != 0 {
		t.Errorf("Expected 0 attendees, got %d", len(attendees))
	}
}

func TestGetAttendeeStats(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("GET", "/api/admin/stats", nil)
	w := httptest.NewRecorder()
	
	GetAttendeeStats(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var stats map[string]int
	if err := json.Unmarshal(w.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if stats == nil {
		t.Error("Expected stats map, got nil")
	}
}

func TestRegisterAttendee_WrongMethod(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("GET", "/api/attendees/register", nil)
	w := httptest.NewRecorder()
	
	RegisterAttendee(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

