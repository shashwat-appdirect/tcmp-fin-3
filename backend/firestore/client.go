package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var (
	Client    *firestore.Client
	ClientID  string
	projectID string
)

// InitializeFirestore initializes the Firestore client and extracts client_id
func InitializeFirestore(ctx context.Context, serviceAccountPath string) error {
	// Read service account JSON
	data, err := os.ReadFile(serviceAccountPath)
	if err != nil {
		return fmt.Errorf("failed to read service account file: %w", err)
	}

	// Parse service account JSON to extract client_id and project_id
	var sa map[string]interface{}
	if err := json.Unmarshal(data, &sa); err != nil {
		return fmt.Errorf("failed to parse service account JSON: %w", err)
	}

	// Extract client_id
	if clientID, ok := sa["client_id"].(string); ok {
		ClientID = clientID
	} else {
		return fmt.Errorf("client_id not found in service account JSON")
	}

	// Extract project_id
	if projID, ok := sa["project_id"].(string); ok {
		projectID = projID
	} else {
		return fmt.Errorf("project_id not found in service account JSON")
	}

	// Initialize Firestore client
	Client, err = firestore.NewClient(ctx, projectID, option.WithCredentialsFile(serviceAccountPath))
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}

	return nil
}

// GetAttendeesPath returns the path to the attendees subcollection
func GetAttendeesPath() string {
	return fmt.Sprintf("clients/%s/attendees", ClientID)
}

// GetSessionsPath returns the path to the sessions subcollection
func GetSessionsPath() string {
	return fmt.Sprintf("clients/%s/sessions", ClientID)
}

// GetSpeakersPath returns the path to the speakers subcollection
func GetSpeakersPath() string {
	return fmt.Sprintf("clients/%s/speakers", ClientID)
}

// GetAttendeesCollection returns the attendees collection reference
func GetAttendeesCollection() CollectionRefInterface {
	mockClient := getMockClient()
	if mockClient != nil {
		clientID := getClientID()
		return mockClient.Collection("clients").Doc(clientID).Collection("attendees")
	}
	client := getClient()
	clientID := getClientID()
	return &RealCollectionRef{ref: client.Collection("clients").Doc(clientID).Collection("attendees")}
}

// GetSessionsCollection returns the sessions collection reference
func GetSessionsCollection() CollectionRefInterface {
	mockClient := getMockClient()
	if mockClient != nil {
		clientID := getClientID()
		return mockClient.Collection("clients").Doc(clientID).Collection("sessions")
	}
	client := getClient()
	clientID := getClientID()
	return &RealCollectionRef{ref: client.Collection("clients").Doc(clientID).Collection("sessions")}
}

// GetSpeakersCollection returns the speakers collection reference
func GetSpeakersCollection() CollectionRefInterface {
	mockClient := getMockClient()
	if mockClient != nil {
		clientID := getClientID()
		return mockClient.Collection("clients").Doc(clientID).Collection("speakers")
	}
	client := getClient()
	clientID := getClientID()
	return &RealCollectionRef{ref: client.Collection("clients").Doc(clientID).Collection("speakers")}
}

