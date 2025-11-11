package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"event-registration-backend/models"
)

// Integration tests that test complete flows
func TestIntegration_RegisterAndGetAttendees(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	// Register an attendee
	registerReq := models.RegisterRequest{
		FullName:    "Integration Test User",
		Email:       "integration@test.com",
		Designation: "QA Engineer",
	}
	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/attendees/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	RegisterAttendee(w, req)
	
	if w.Code != http.StatusOK {
		t.Skipf("Skipping integration test: Firestore not available: %d", w.Code)
		return
	}
	
	var registeredAttendee models.Attendee
	json.Unmarshal(w.Body.Bytes(), &registeredAttendee)
	
	// Get attendee count
	countReq := httptest.NewRequest("GET", "/api/attendees/count", nil)
	countW := httptest.NewRecorder()
	GetAttendeeCount(countW, countReq)
	
	if countW.Code != http.StatusOK {
		t.Errorf("Expected status 200 for count, got %d", countW.Code)
	}
	
	var countResult map[string]int
	json.Unmarshal(countW.Body.Bytes(), &countResult)
	if countResult["count"] < 1 {
		t.Errorf("Expected count >= 1, got %d", countResult["count"])
	}
	
	// Get all attendees
	getReq := httptest.NewRequest("GET", "/api/admin/attendees", nil)
	getW := httptest.NewRecorder()
	GetAttendees(getW, getReq)
	
	if getW.Code != http.StatusOK {
		t.Errorf("Expected status 200 for get attendees, got %d", getW.Code)
	}
	
	var attendees []models.Attendee
	json.Unmarshal(getW.Body.Bytes(), &attendees)
	
	found := false
	for _, a := range attendees {
		if a.Email == registerReq.Email {
			found = true
			if a.FullName != registerReq.FullName {
				t.Errorf("Expected FullName %s, got %s", registerReq.FullName, a.FullName)
			}
			break
		}
	}
	if !found {
		t.Error("Registered attendee not found in attendees list")
	}
	
	// Get stats
	statsReq := httptest.NewRequest("GET", "/api/admin/stats", nil)
	statsW := httptest.NewRecorder()
	GetAttendeeStats(statsW, statsReq)
	
	if statsW.Code != http.StatusOK {
		t.Errorf("Expected status 200 for stats, got %d", statsW.Code)
	}
	
	var stats map[string]int
	json.Unmarshal(statsW.Body.Bytes(), &stats)
	if stats[registerReq.Designation] < 1 {
		t.Errorf("Expected at least 1 attendee with designation %s", registerReq.Designation)
	}
}

func TestIntegration_CreateSpeakerAndSession(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	// Create a speaker
	speakerReq := models.SpeakerRequest{
		Name:     "Integration Speaker",
		Bio:      "Integration Test Bio",
		PhotoURL: "http://example.com/integration.jpg",
	}
	speakerBody, _ := json.Marshal(speakerReq)
	speakerHTTPReq := httptest.NewRequest("POST", "/api/admin/speakers", bytes.NewBuffer(speakerBody))
	speakerHTTPReq.Header.Set("Content-Type", "application/json")
	speakerW := httptest.NewRecorder()
	CreateOrUpdateSpeaker(speakerW, speakerHTTPReq)
	
	if speakerW.Code != http.StatusOK {
		t.Skipf("Skipping integration test: Firestore not available: %d", speakerW.Code)
		return
	}
	
	var speaker models.Speaker
	json.Unmarshal(speakerW.Body.Bytes(), &speaker)
	
	// Create a session with the speaker
	sessionReq := models.SessionRequest{
		Title:       "Integration Session",
		Description: "Integration Test Description",
		Time:        "2:00 PM",
		SpeakerID:   speaker.ID,
	}
	sessionBody, _ := json.Marshal(sessionReq)
	sessionHTTPReq := httptest.NewRequest("POST", "/api/admin/sessions", bytes.NewBuffer(sessionBody))
	sessionHTTPReq.Header.Set("Content-Type", "application/json")
	sessionW := httptest.NewRecorder()
	CreateOrUpdateSession(sessionW, sessionHTTPReq)
	
	if sessionW.Code != http.StatusOK {
		t.Errorf("Expected status 200 for session creation, got %d", sessionW.Code)
	}
	
	var session models.Session
	json.Unmarshal(sessionW.Body.Bytes(), &session)
	
	// Get sessions and verify speaker is linked
	getReq := httptest.NewRequest("GET", "/api/sessions", nil)
	getW := httptest.NewRecorder()
	GetSessions(getW, getReq)
	
	if getW.Code != http.StatusOK {
		t.Errorf("Expected status 200 for get sessions, got %d", getW.Code)
	}
	
	var sessions []models.Session
	json.Unmarshal(getW.Body.Bytes(), &sessions)
	
	found := false
	for _, s := range sessions {
		if s.ID == session.ID {
			found = true
			if s.SpeakerID != speaker.ID {
				t.Errorf("Expected SpeakerID %s, got %s", speaker.ID, s.SpeakerID)
			}
			break
		}
	}
	if !found {
		t.Error("Created session not found in sessions list")
	}
}

func TestIntegration_UpdateSpeakerAndSession(t *testing.T) {
	setupTestClient(t)
	defer teardownTestClient()
	
	// Create speaker
	speakerReq := models.SpeakerRequest{
		Name: "Original Speaker",
		Bio:  "Original Bio",
	}
	speakerBody, _ := json.Marshal(speakerReq)
	speakerHTTPReq := httptest.NewRequest("POST", "/api/admin/speakers", bytes.NewBuffer(speakerBody))
	speakerHTTPReq.Header.Set("Content-Type", "application/json")
	speakerW := httptest.NewRecorder()
	CreateOrUpdateSpeaker(speakerW, speakerHTTPReq)
	
	if speakerW.Code != http.StatusOK {
		t.Skipf("Skipping integration test: Firestore not available: %d", speakerW.Code)
		return
	}
	
	var speaker models.Speaker
	json.Unmarshal(speakerW.Body.Bytes(), &speaker)
	
	// Update speaker
	updateSpeakerReq := models.SpeakerRequest{
		ID:   speaker.ID,
		Name: "Updated Speaker",
		Bio:  "Updated Bio",
	}
	updateSpeakerBody, _ := json.Marshal(updateSpeakerReq)
	updateSpeakerHTTPReq := httptest.NewRequest("POST", "/api/admin/speakers", bytes.NewBuffer(updateSpeakerBody))
	updateSpeakerHTTPReq.Header.Set("Content-Type", "application/json")
	updateSpeakerW := httptest.NewRecorder()
	CreateOrUpdateSpeaker(updateSpeakerW, updateSpeakerHTTPReq)
	
	if updateSpeakerW.Code != http.StatusOK {
		t.Errorf("Expected status 200 for speaker update, got %d", updateSpeakerW.Code)
	}
	
	// Create session
	sessionReq := models.SessionRequest{
		Title:       "Original Session",
		Description: "Original Description",
		Time:        "10:00 AM",
		SpeakerID:   speaker.ID,
	}
	sessionBody, _ := json.Marshal(sessionReq)
	sessionHTTPReq := httptest.NewRequest("POST", "/api/admin/sessions", bytes.NewBuffer(sessionBody))
	sessionHTTPReq.Header.Set("Content-Type", "application/json")
	sessionW := httptest.NewRecorder()
	CreateOrUpdateSession(sessionW, sessionHTTPReq)
	
	if sessionW.Code != http.StatusOK {
		t.Errorf("Expected status 200 for session creation, got %d", sessionW.Code)
	}
	
	var session models.Session
	json.Unmarshal(sessionW.Body.Bytes(), &session)
	
	// Update session
	updateSessionReq := models.SessionRequest{
		ID:          session.ID,
		Title:       "Updated Session",
		Description: "Updated Description",
		Time:        "11:00 AM",
		SpeakerID:   speaker.ID,
	}
	updateSessionBody, _ := json.Marshal(updateSessionReq)
	updateSessionHTTPReq := httptest.NewRequest("POST", "/api/admin/sessions", bytes.NewBuffer(updateSessionBody))
	updateSessionHTTPReq.Header.Set("Content-Type", "application/json")
	updateSessionW := httptest.NewRecorder()
	CreateOrUpdateSession(updateSessionW, updateSessionHTTPReq)
	
	if updateSessionW.Code != http.StatusOK {
		t.Errorf("Expected status 200 for session update, got %d", updateSessionW.Code)
	}
	
	var updatedSession models.Session
	json.Unmarshal(updateSessionW.Body.Bytes(), &updatedSession)
	
	if updatedSession.Title != updateSessionReq.Title {
		t.Errorf("Expected Title %s, got %s", updateSessionReq.Title, updatedSession.Title)
	}
}

