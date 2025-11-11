package handlers

import (
	"encoding/json"
	"net/http"

	"event-registration-backend/firestore"
	"event-registration-backend/models"
)

// GetSessions returns all sessions with speaker details
func GetSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	sessionsRef := firestore.GetSessionsCollection()
	sessionsDocs, err := sessionsRef.Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, "Failed to get sessions", http.StatusInternalServerError)
		return
	}

	speakersRef := firestore.GetSpeakersCollection()
	speakersDocs, err := speakersRef.Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, "Failed to get speakers", http.StatusInternalServerError)
		return
	}

	// Create speakers map
	speakersMap := make(map[string]*models.Speaker)
	for _, doc := range speakersDocs {
		var speaker models.Speaker
		if err := doc.DataTo(&speaker); err != nil {
			continue
		}
		speaker.ID = doc.GetID()
		speakersMap[speaker.ID] = &speaker
	}

	// Build sessions with speaker details
	var sessions []models.Session
	for _, doc := range sessionsDocs {
		var session models.Session
		if err := doc.DataTo(&session); err != nil {
			continue
		}
		session.ID = doc.GetID()
		if speaker, ok := speakersMap[session.SpeakerID]; ok {
			session.Speaker = speaker
		}
		sessions = append(sessions, session)
	}

	// Ensure we return an empty array, not null
	if sessions == nil {
		sessions = []models.Session{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// CreateOrUpdateSession creates or updates a session (admin only)
func CreateOrUpdateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Title == "" || req.Description == "" || req.Time == "" || req.SpeakerID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	sessionsRef := firestore.GetSessionsCollection()

	session := models.Session{
		Title:       req.Title,
		Description: req.Description,
		Time:        req.Time,
		SpeakerID:   req.SpeakerID,
	}

	if req.ID != "" {
		// Update existing session
		_, err := sessionsRef.Doc(req.ID).Set(ctx, session)
		if err != nil {
			http.Error(w, "Failed to update session", http.StatusInternalServerError)
			return
		}
		session.ID = req.ID
	} else {
		// Create new session
		docRef, _, err := sessionsRef.Add(ctx, session)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}
		session.ID = docRef.GetID()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

