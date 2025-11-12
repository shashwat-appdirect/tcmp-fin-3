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
	Client   *firestore.Client
	ClientID string
)

// InitFirestore initializes the Firestore client using service account credentials
func InitFirestore(ctx context.Context) error {
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsPath == "" {
		return fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable not set")
	}

	// Read the service account JSON to extract client_id
	credentialsData, err := os.ReadFile(credentialsPath)
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds struct {
		ClientID string `json:"client_id"`
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(credentialsData, &creds); err != nil {
		return fmt.Errorf("failed to parse credentials: %w", err)
	}

	ClientID = creds.ClientID

	// Initialize Firestore client
	opt := option.WithCredentialsFile(credentialsPath)
	client, err := firestore.NewClient(ctx, creds.ProjectID, opt)
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}

	Client = client
	return nil
}

// GetCollectionPath returns the path to a subcollection under events/{client_id}
func GetCollectionPath(collectionName string) string {
	return fmt.Sprintf("events/%s/%s", ClientID, collectionName)
}

