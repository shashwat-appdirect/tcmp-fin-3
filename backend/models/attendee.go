package models

import "time"

type Attendee struct {
	ID          string    `json:"id" firestore:"-"`
	FullName    string    `json:"fullName" firestore:"fullName"`
	Email       string    `json:"email" firestore:"email"`
	Designation string    `json:"designation" firestore:"designation"`
	RegisteredAt time.Time `json:"registeredAt" firestore:"registeredAt"`
}

type RegisterRequest struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	Designation string `json:"designation"`
}

