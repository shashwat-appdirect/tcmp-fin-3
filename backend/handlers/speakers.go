package handlers

import (
	"encoding/json"
	"net/http"

	"event-registration-backend/firestore"
	"event-registration-backend/models"
)

// GetSpeakers returns all speakers
func GetSpeakers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	speakersRef := firestore.GetSpeakersCollection()
	docs, err := speakersRef.Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, "Failed to get speakers", http.StatusInternalServerError)
		return
	}

	var speakers []models.Speaker
	for _, doc := range docs {
		var speaker models.Speaker
		if err := doc.DataTo(&speaker); err != nil {
			continue
		}
		speaker.ID = doc.GetID()
		speakers = append(speakers, speaker)
	}

	// Ensure we return an empty array, not null
	if speakers == nil {
		speakers = []models.Speaker{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(speakers)
}

// CreateOrUpdateSpeaker creates or updates a speaker (admin only)
func CreateOrUpdateSpeaker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SpeakerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	speakersRef := firestore.GetSpeakersCollection()

	speaker := models.Speaker{
		Name:     req.Name,
		Bio:      req.Bio,
		PhotoURL: req.PhotoURL,
	}

	if req.ID != "" {
		// Update existing speaker
		_, err := speakersRef.Doc(req.ID).Set(ctx, speaker)
		if err != nil {
			http.Error(w, "Failed to update speaker", http.StatusInternalServerError)
			return
		}
		speaker.ID = req.ID
	} else {
		// Create new speaker
		docRef, _, err := speakersRef.Add(ctx, speaker)
		if err != nil {
			http.Error(w, "Failed to create speaker", http.StatusInternalServerError)
			return
		}
		speaker.ID = docRef.GetID()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(speaker)
}

