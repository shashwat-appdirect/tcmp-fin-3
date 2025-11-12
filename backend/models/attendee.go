package models

import (
	"time"
)

type Attendee struct {
	ID          string    `firestore:"id" json:"id"`
	Name        string    `firestore:"name" json:"name"`
	Email       string    `firestore:"email" json:"email"`
	Designation string    `firestore:"designation" json:"designation"`
	RegisteredAt time.Time `firestore:"registeredAt" json:"registeredAt"`
}

type AttendeeRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Designation string `json:"designation"`
}

type AttendeesResponse struct {
	Count     int        `json:"count"`
	Attendees []Attendee `json:"attendees"`
}
