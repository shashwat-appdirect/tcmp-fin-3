package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"event-registration-backend/firestore"
	"event-registration-backend/models"
)

// RegisterAttendee handles attendee registration
func RegisterAttendee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.FullName == "" || req.Email == "" || req.Designation == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	ctx := r.Context()
	attendeesRef := firestore.GetAttendeesCollection()
	query := attendeesRef.Where("email", "==", req.Email).Limit(1)
	docs, err := query.Documents(ctx).GetAll()
	if err == nil && len(docs) > 0 {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Create attendee
	attendee := models.Attendee{
		FullName:     req.FullName,
		Email:        req.Email,
		Designation:  req.Designation,
		RegisteredAt: time.Now(),
	}

	docRef, _, err := attendeesRef.Add(ctx, attendee)
	if err != nil {
		http.Error(w, "Failed to register attendee", http.StatusInternalServerError)
		return
	}

	attendee.ID = docRef.GetID()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendee)
}

// GetAttendeeCount returns the total count of registered attendees
func GetAttendeeCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	attendeesRef := firestore.GetAttendeesCollection()
	docs, err := attendeesRef.Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, "Failed to get attendee count", http.StatusInternalServerError)
		return
	}

	count := len(docs)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}

// GetAttendees returns all attendees (admin only)
func GetAttendees(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	attendeesRef := firestore.GetAttendeesCollection()
	docs, err := attendeesRef.Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, "Failed to get attendees", http.StatusInternalServerError)
		return
	}

	var attendees []models.Attendee
	for _, doc := range docs {
		var attendee models.Attendee
		if err := doc.DataTo(&attendee); err != nil {
			continue
		}
		attendee.ID = doc.GetID()
		attendees = append(attendees, attendee)
	}

	// Ensure we return an empty array, not null
	if attendees == nil {
		attendees = []models.Attendee{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendees)
}

// GetAttendeeStats returns attendee breakdown by designation
func GetAttendeeStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	attendeesRef := firestore.GetAttendeesCollection()
	docs, err := attendeesRef.Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, "Failed to get attendee stats", http.StatusInternalServerError)
		return
	}

	stats := make(map[string]int)
	for _, doc := range docs {
		var attendee models.Attendee
		if err := doc.DataTo(&attendee); err != nil {
			continue
		}
		stats[attendee.Designation]++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

