package models

type Session struct {
	ID          string `json:"id" firestore:"-"`
	Title       string `json:"title" firestore:"title"`
	Description string `json:"description" firestore:"description"`
	Time        string `json:"time" firestore:"time"`
	SpeakerID   string `json:"speakerId" firestore:"speakerId"`
	Speaker     *Speaker `json:"speaker,omitempty" firestore:"-"`
}

type SessionRequest struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
	SpeakerID   string `json:"speakerId"`
}

