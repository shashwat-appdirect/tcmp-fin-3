package models

type Speaker struct {
	ID      string `json:"id" firestore:"-"`
	Name    string `json:"name" firestore:"name"`
	Bio     string `json:"bio" firestore:"bio"`
	PhotoURL string `json:"photoUrl" firestore:"photoUrl"`
}

type SpeakerRequest struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	PhotoURL string `json:"photoUrl"`
}

