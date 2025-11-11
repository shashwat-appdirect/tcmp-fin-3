package firestore

import (
	"sync"

	"cloud.google.com/go/firestore"
)

var (
	testClient    *firestore.Client
	testMockClient *MockFirestoreClient
	testClientID  string
	testMode      bool
	testMutex     sync.RWMutex
)

// SetTestClient sets a test Firestore client
func SetTestClient(client *firestore.Client, clientID string) {
	testMutex.Lock()
	defer testMutex.Unlock()
	testClient = client
	testMockClient = nil
	testClientID = clientID
	testMode = true
}

// SetMockClient sets a mock Firestore client for testing
func SetMockClient(mockClient *MockFirestoreClient, clientID string) {
	testMutex.Lock()
	defer testMutex.Unlock()
	testClient = nil
	testMockClient = mockClient
	testClientID = clientID
	testMode = true
}

// ClearTestClient clears the test client
func ClearTestClient() {
	testMutex.Lock()
	defer testMutex.Unlock()
	testClient = nil
	testMockClient = nil
	testClientID = ""
	testMode = false
}

// getClient returns the appropriate client (test or production)
func getClient() *firestore.Client {
	testMutex.RLock()
	defer testMutex.RUnlock()
	if testMode && testClient != nil {
		return testClient
	}
	return Client
}

// getMockClient returns the mock client if in test mode
func getMockClient() *MockFirestoreClient {
	testMutex.RLock()
	defer testMutex.RUnlock()
	if testMode && testMockClient != nil {
		return testMockClient
	}
	return nil
}

// getClientID returns the appropriate client ID (test or production)
func getClientID() string {
	testMutex.RLock()
	defer testMutex.RUnlock()
	if testMode && testClientID != "" {
		return testClientID
	}
	return ClientID
}


