package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	fs "github.com/tcmp-fin-3/backend/firestore"
	"github.com/tcmp-fin-3/backend/models"
)

func GetAttendees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	collectionPath := fs.GetCollectionPath("attendees")

	docs, err := fs.Client.Collection(collectionPath).Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch attendees: %v", err), http.StatusInternalServerError)
		return
	}

	attendees := make([]models.Attendee, 0, len(docs))
	for _, doc := range docs {
		var attendee models.Attendee
		if err := doc.DataTo(&attendee); err != nil {
			continue
		}
		attendee.ID = doc.Ref.ID
		attendees = append(attendees, attendee)
	}

	response := models.AttendeesResponse{
		Count:     len(attendees),
		Attendees: attendees,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RegisterAttendee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AttendeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Designation == "" {
		http.Error(w, "Name, email, and designation are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	collectionPath := fs.GetCollectionPath("attendees")

	// Check if email already exists
	query := fs.Client.Collection(collectionPath).Where("email", "==", req.Email).Limit(1)
	existing, err := query.Documents(ctx).GetAll()
	if err == nil && len(existing) > 0 {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	attendee := models.Attendee{
		Name:        req.Name,
		Email:       req.Email,
		Designation: req.Designation,
		RegisteredAt: time.Now(),
	}

	docRef, _, err := fs.Client.Collection(collectionPath).Add(ctx, attendee)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to register: %v", err), http.StatusInternalServerError)
		return
	}

	attendee.ID = docRef.ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attendee)
}
