package server

import (
	"encoding/json"
	"testing"
)

func TestHandleGetProfiles(t *testing.T) {
	result, err := handleGetProfiles(nil)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Try to marshal result to ensure it's serializable
	_, err = json.Marshal(result)
	if err != nil {
		t.Errorf("Result is not JSON serializable: %v", err)
	}
}

func TestRegisterOpenCodeHandlers(t *testing.T) {
	srv := &Server{
		handlers: make(map[string]Handler),
	}

	RegisterOpenCodeHandlers(srv)

	// Verify handler is registered
	if _, exists := srv.handlers["opencode.getProfiles"]; !exists {
		t.Error("Expected opencode.getProfiles handler to be registered")
	}
}
