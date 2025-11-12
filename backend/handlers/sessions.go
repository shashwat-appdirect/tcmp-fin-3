package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	fs "github.com/tcmp-fin-3/backend/firestore"
	"github.com/tcmp-fin-3/backend/models"
)

func GetSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionsPath := fs.GetCollectionPath("sessions")
	speakersPath := fs.GetCollectionPath("speakers")

	// Get all sessions
	sessionDocs, err := fs.Client.Collection(sessionsPath).Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch sessions: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all speakers
	speakerDocs, err := fs.Client.Collection(speakersPath).Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch speakers: %v", err), http.StatusInternalServerError)
		return
	}

	// Build speaker map
	speakerMap := make(map[string]models.Speaker)
	for _, doc := range speakerDocs {
		var speaker models.Speaker
		if err := doc.DataTo(&speaker); err != nil {
			continue
		}
		speaker.ID = doc.Ref.ID
		speakerMap[speaker.ID] = speaker
	}

	// Build sessions with speakers
	sessions := make([]models.SessionWithSpeakers, 0, len(sessionDocs))
	for _, doc := range sessionDocs {
		var session models.Session
		if err := doc.DataTo(&session); err != nil {
			continue
		}
		session.ID = doc.Ref.ID

		// Get associated speakers
		speakers := make([]models.Speaker, 0)
		for _, speakerID := range session.SpeakerIDs {
			if speaker, ok := speakerMap[speakerID]; ok {
				speakers = append(speakers, speaker)
			}
		}

		sessions = append(sessions, models.SessionWithSpeakers{
			Session:  session,
			Speakers: speakers,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func GetSpeakers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	collectionPath := fs.GetCollectionPath("speakers")

	docs, err := fs.Client.Collection(collectionPath).Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch speakers: %v", err), http.StatusInternalServerError)
		return
	}

	speakers := make([]models.Speaker, 0, len(docs))
	for _, doc := range docs {
		var speaker models.Speaker
		if err := doc.DataTo(&speaker); err != nil {
			continue
		}
		speaker.ID = doc.Ref.ID
		speakers = append(speakers, speaker)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(speakers)
}
