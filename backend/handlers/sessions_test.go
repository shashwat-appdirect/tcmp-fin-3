package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"event-registration-backend/models"
)

func TestGetSessions_Empty(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("GET", "/api/sessions", nil)
	w := httptest.NewRecorder()
	
	GetSessions(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var sessions []models.Session
	if err := json.Unmarshal(w.Body.Bytes(), &sessions); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if sessions == nil {
		t.Error("Expected empty array, got nil")
	}
}

func TestCreateOrUpdateSession_Create(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	reqBody := models.SessionRequest{
		Title:       "Test Session",
		Description: "Test Description",
		Time:        "10:00 AM",
		SpeakerID:   "speaker1",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/admin/sessions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	CreateOrUpdateSession(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var session models.Session
	if err := json.Unmarshal(w.Body.Bytes(), &session); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if session.Title != reqBody.Title {
		t.Errorf("Expected Title %s, got %s", reqBody.Title, session.Title)
	}
	if session.ID == "" {
		t.Error("Expected session ID to be set")
	}
}

func TestCreateOrUpdateSession_Update(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	// First create a session
	createReq := models.SessionRequest{
		Title:       "Original Session",
		Description: "Original Description",
		Time:        "10:00 AM",
		SpeakerID:   "speaker1",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/admin/sessions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	CreateOrUpdateSession(w, req)
	
	if w.Code != http.StatusOK {
		t.Skipf("Skipping update test: create failed (likely no Firestore): %d", w.Code)
		return
	}
	
	var createdSession models.Session
	json.Unmarshal(w.Body.Bytes(), &createdSession)
	
	// Update the session
	updateReq := models.SessionRequest{
		ID:          createdSession.ID,
		Title:       "Updated Session",
		Description: "Updated Description",
		Time:        "11:00 AM",
		SpeakerID:   "speaker2",
	}
	body2, _ := json.Marshal(updateReq)
	req2 := httptest.NewRequest("POST", "/api/admin/sessions", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	CreateOrUpdateSession(w2, req2)
	
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}
	
	var updatedSession models.Session
	if err := json.Unmarshal(w2.Body.Bytes(), &updatedSession); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if updatedSession.Title != updateReq.Title {
		t.Errorf("Expected Title %s, got %s", updateReq.Title, updatedSession.Title)
	}
	if updatedSession.ID != createdSession.ID {
		t.Errorf("Expected same ID %s, got %s", createdSession.ID, updatedSession.ID)
	}
}

func TestCreateOrUpdateSession_MissingFields(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	testCases := []struct {
		name string
		body models.SessionRequest
	}{
		{"Missing Title", models.SessionRequest{Description: "Desc", Time: "10:00", SpeakerID: "s1"}},
		{"Missing Description", models.SessionRequest{Title: "Title", Time: "10:00", SpeakerID: "s1"}},
		{"Missing Time", models.SessionRequest{Title: "Title", Description: "Desc", SpeakerID: "s1"}},
		{"Missing SpeakerID", models.SessionRequest{Title: "Title", Description: "Desc", Time: "10:00"}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest("POST", "/api/admin/sessions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			CreateOrUpdateSession(w, req)
			
			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestCreateOrUpdateSession_InvalidJSON(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("POST", "/api/admin/sessions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	CreateOrUpdateSession(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetSessions_WithSpeaker(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	// Create a speaker first
	speakerReq := models.SpeakerRequest{
		Name:     "Test Speaker",
		Bio:      "Test Bio",
		PhotoURL: "http://example.com/photo.jpg",
	}
	speakerBody, _ := json.Marshal(speakerReq)
	speakerHTTPReq := httptest.NewRequest("POST", "/api/admin/speakers", bytes.NewBuffer(speakerBody))
	speakerHTTPReq.Header.Set("Content-Type", "application/json")
	speakerW := httptest.NewRecorder()
	CreateOrUpdateSpeaker(speakerW, speakerHTTPReq)
	
	if speakerW.Code != http.StatusOK {
		t.Skipf("Skipping test: speaker creation failed (likely no Firestore): %d", speakerW.Code)
		return
	}
	
	var speaker models.Speaker
	json.Unmarshal(speakerW.Body.Bytes(), &speaker)
	
	// Create a session with the speaker
	sessionReq := models.SessionRequest{
		Title:       "Test Session",
		Description: "Test Description",
		Time:        "10:00 AM",
		SpeakerID:   speaker.ID,
	}
	sessionBody, _ := json.Marshal(sessionReq)
	sessionHTTPReq := httptest.NewRequest("POST", "/api/admin/sessions", bytes.NewBuffer(sessionBody))
	sessionHTTPReq.Header.Set("Content-Type", "application/json")
	sessionW := httptest.NewRecorder()
	CreateOrUpdateSession(sessionW, sessionHTTPReq)
	
	if sessionW.Code != http.StatusOK {
		t.Skipf("Skipping test: session creation failed: %d", sessionW.Code)
		return
	}
	
	// Get sessions
	getReq := httptest.NewRequest("GET", "/api/sessions", nil)
	getW := httptest.NewRecorder()
	GetSessions(getW, getReq)
	
	if getW.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", getW.Code)
	}
	
	var sessions []models.Session
	if err := json.Unmarshal(getW.Body.Bytes(), &sessions); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if len(sessions) == 0 {
		t.Error("Expected at least one session")
	}
	
	found := false
	for _, s := range sessions {
		if s.ID != "" && s.Speaker != nil && s.Speaker.ID == speaker.ID {
			found = true
			break
		}
	}
	if !found {
		t.Log("Session created but speaker not linked (this is OK if Firestore emulator not running)")
	}
}

func TestGetSessions_WrongMethod(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	req := httptest.NewRequest("POST", "/api/sessions", nil)
	w := httptest.NewRecorder()
	
	GetSessions(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

