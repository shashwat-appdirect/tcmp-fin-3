package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	fs "github.com/tcmp-fin-3/backend/firestore"
	"github.com/tcmp-fin-3/backend/models"
)

// AdminMiddleware validates admin password
func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := r.Header.Get("X-Admin-Password")
		expectedPassword := os.Getenv("ADMIN_PASSWORD")

		if expectedPassword == "" {
			http.Error(w, "Admin password not configured", http.StatusInternalServerError)
			return
		}

		if password != expectedPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	expectedPassword := os.Getenv("ADMIN_PASSWORD")
	if expectedPassword == "" {
		http.Error(w, "Admin password not configured", http.StatusInternalServerError)
		return
	}

	if req.Password != expectedPassword {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func GetAdminAttendees(w http.ResponseWriter, r *http.Request) {
	GetAttendees(w, r)
}

func DeleteAttendee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	attendeeID := vars["id"]

	if attendeeID == "" {
		http.Error(w, "Attendee ID required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	collectionPath := fs.GetCollectionPath("attendees")

	_, err := fs.Client.Collection(collectionPath).Doc(attendeeID).Delete(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete attendee: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func AddOrUpdateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	collectionPath := fs.GetCollectionPath("sessions")

	session := models.Session{
		Title:       req.Title,
		Description: req.Description,
		Time:        req.Time,
		SpeakerIDs:  req.SpeakerIDs,
	}

	var docRef *firestore.DocumentRef
	if req.ID != "" {
		// Update existing
		docRef = fs.Client.Collection(collectionPath).Doc(req.ID)
		_, err := docRef.Set(ctx, session)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update session: %v", err), http.StatusInternalServerError)
			return
		}
		session.ID = req.ID
	} else {
		// Create new
		ref, _, err := fs.Client.Collection(collectionPath).Add(ctx, session)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create session: %v", err), http.StatusInternalServerError)
			return
		}
		docRef = ref
		session.ID = docRef.ID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func DeleteSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	collectionPath := fs.GetCollectionPath("sessions")

	_, err := fs.Client.Collection(collectionPath).Doc(sessionID).Delete(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete session: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func AddOrUpdateSpeaker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SpeakerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	collectionPath := fs.GetCollectionPath("speakers")

	speaker := models.Speaker{
		Name:     req.Name,
		Bio:      req.Bio,
		PhotoURL: req.PhotoURL,
	}

	var docRef *firestore.DocumentRef
	if req.ID != "" {
		// Update existing
		docRef = fs.Client.Collection(collectionPath).Doc(req.ID)
		_, err := docRef.Set(ctx, speaker)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update speaker: %v", err), http.StatusInternalServerError)
			return
		}
		speaker.ID = req.ID
	} else {
		// Create new
		ref, _, err := fs.Client.Collection(collectionPath).Add(ctx, speaker)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create speaker: %v", err), http.StatusInternalServerError)
			return
		}
		docRef = ref
		speaker.ID = docRef.ID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(speaker)
}

func DeleteSpeaker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	speakerID := vars["id"]

	if speakerID == "" {
		http.Error(w, "Speaker ID required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	collectionPath := fs.GetCollectionPath("speakers")

	_, err := fs.Client.Collection(collectionPath).Doc(speakerID).Delete(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete speaker: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	collectionPath := fs.GetCollectionPath("attendees")

	docs, err := fs.Client.Collection(collectionPath).Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch attendees: %v", err), http.StatusInternalServerError)
		return
	}

	// Count by designation
	designationCount := make(map[string]int)
	for _, doc := range docs {
		var attendee models.Attendee
		if err := doc.DataTo(&attendee); err != nil {
			continue
		}
		designationCount[attendee.Designation]++
	}

	// Convert to array format for pie chart
	stats := make([]map[string]interface{}, 0, len(designationCount))
	for designation, count := range designationCount {
		stats = append(stats, map[string]interface{}{
			"designation": designation,
			"count":       count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
