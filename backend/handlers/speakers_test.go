package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"event-registration-backend/models"
)

func TestGetSpeakers_Empty(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("GET", "/api/speakers", nil)
	w := httptest.NewRecorder()
	
	GetSpeakers(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var speakers []models.Speaker
	if err := json.Unmarshal(w.Body.Bytes(), &speakers); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if speakers == nil {
		t.Error("Expected empty array, got nil")
	}
}

func TestCreateOrUpdateSpeaker_Create(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	reqBody := models.SpeakerRequest{
		Name:     "Test Speaker",
		Bio:      "Test Bio",
		PhotoURL: "http://example.com/photo.jpg",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/admin/speakers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	CreateOrUpdateSpeaker(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var speaker models.Speaker
	if err := json.Unmarshal(w.Body.Bytes(), &speaker); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if speaker.Name != reqBody.Name {
		t.Errorf("Expected Name %s, got %s", reqBody.Name, speaker.Name)
	}
	if speaker.ID == "" {
		t.Error("Expected speaker ID to be set")
	}
}

func TestCreateOrUpdateSpeaker_Update(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	// First create a speaker
	createReq := models.SpeakerRequest{
		Name:     "Original Speaker",
		Bio:      "Original Bio",
		PhotoURL: "http://example.com/original.jpg",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/admin/speakers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	CreateOrUpdateSpeaker(w, req)
	
	if w.Code != http.StatusOK {
		t.Skipf("Skipping update test: create failed (likely no Firestore): %d", w.Code)
		return
	}
	
	var createdSpeaker models.Speaker
	json.Unmarshal(w.Body.Bytes(), &createdSpeaker)
	
	// Update the speaker
	updateReq := models.SpeakerRequest{
		ID:       createdSpeaker.ID,
		Name:     "Updated Speaker",
		Bio:      "Updated Bio",
		PhotoURL: "http://example.com/updated.jpg",
	}
	body2, _ := json.Marshal(updateReq)
	req2 := httptest.NewRequest("POST", "/api/admin/speakers", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	CreateOrUpdateSpeaker(w2, req2)
	
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}
	
	var updatedSpeaker models.Speaker
	if err := json.Unmarshal(w2.Body.Bytes(), &updatedSpeaker); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if updatedSpeaker.Name != updateReq.Name {
		t.Errorf("Expected Name %s, got %s", updateReq.Name, updatedSpeaker.Name)
	}
	if updatedSpeaker.ID != createdSpeaker.ID {
		t.Errorf("Expected same ID %s, got %s", createdSpeaker.ID, updatedSpeaker.ID)
	}
}

func TestCreateOrUpdateSpeaker_MissingName(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	reqBody := models.SpeakerRequest{
		Bio:      "Test Bio",
		PhotoURL: "http://example.com/photo.jpg",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/admin/speakers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	CreateOrUpdateSpeaker(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateOrUpdateSpeaker_InvalidJSON(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("POST", "/api/admin/speakers", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	CreateOrUpdateSpeaker(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetSpeakers_WrongMethod(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("POST", "/api/speakers", nil)
	w := httptest.NewRecorder()
	
	GetSpeakers(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

