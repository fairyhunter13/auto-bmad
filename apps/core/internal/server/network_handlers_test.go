package server

import (
	"io"
	"log"
	"testing"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/network"
)

// TestHandleNetworkGetStatus verifies network.getStatus handler
func TestHandleNetworkGetStatus(t *testing.T) {
	// Call handler with nil params
	result, err := handleNetworkGetStatus(nil)
	if err != nil {
		t.Fatalf("handleNetworkGetStatus failed: %v", err)
	}

	// Result should be a NetworkStatus struct
	status, ok := result.(network.NetworkStatus)
	if !ok {
		t.Fatalf("Expected NetworkStatus, got %T", result)
	}

	// Check status is one of valid values
	validStatuses := map[network.Status]bool{
		network.StatusOnline:   true,
		network.StatusOffline:  true,
		network.StatusChecking: true,
	}
	if !validStatuses[status.Status] {
		t.Errorf("Invalid status value: %v", status.Status)
	}
}

// TestHandleNetworkGetStatus_WithServer verifies handler registration
func TestHandleNetworkGetStatus_WithServer(t *testing.T) {
	stdinR, _ := io.Pipe()
	_, stdoutW := io.Pipe()

	srv := New(stdinR, stdoutW, log.New(io.Discard, "", 0))
	RegisterNetworkHandlers(srv)

	// Verify handler is registered
	srv.mu.RLock()
	_, ok := srv.handlers["network.getStatus"]
	srv.mu.RUnlock()

	if !ok {
		t.Error("network.getStatus handler not registered")
	}
}
