package handlers

import (
	"testing"

	"event-registration-backend/firestore"
)

// setupTestClient creates a mock Firestore client for testing
func setupTestClient(t *testing.T) {
	testClientID := "test-client-id"
	mockClient := firestore.NewMockFirestoreClient(testClientID)
	firestore.SetMockClient(mockClient, testClientID)
}

func teardownTestClient() {
	firestore.ClearTestClient()
}

