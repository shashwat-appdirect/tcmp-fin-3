package models

type Speaker struct {
	ID       string `firestore:"id" json:"id"`
	Name     string `firestore:"name" json:"name"`
	Bio      string `firestore:"bio" json:"bio"`
	PhotoURL string `firestore:"photoUrl,omitempty" json:"photoUrl,omitempty"`
}

type SpeakerRequest struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	PhotoURL string `json:"photoUrl,omitempty"`
}
