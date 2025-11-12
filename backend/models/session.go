package models

type Session struct {
	ID          string   `firestore:"id" json:"id"`
	Title       string   `firestore:"title" json:"title"`
	Description string   `firestore:"description" json:"description"`
	Time        string   `firestore:"time" json:"time"`
	SpeakerIDs  []string `firestore:"speakerIds" json:"speakerIds"`
}

type SessionRequest struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Time        string   `json:"time"`
	SpeakerIDs  []string `json:"speakerIds"`
}

type SessionWithSpeakers struct {
	Session
	Speakers []Speaker `json:"speakers"`
}
