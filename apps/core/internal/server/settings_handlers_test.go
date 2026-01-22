package server

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/state"
)

// TestRegisterSettingsHandlers verifies handler registration
func TestRegisterSettingsHandlers(t *testing.T) {
	// Override home directory for test
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	// Create server
	srv := New(nil, nil, log.New(io.Discard, "", 0))

	// Register handlers
	if err := RegisterSettingsHandlers(srv); err != nil {
		t.Fatalf("RegisterSettingsHandlers failed: %v", err)
	}

	// Verify handlers are registered
	srv.mu.RLock()
	defer srv.mu.RUnlock()

	if _, ok := srv.handlers["settings.get"]; !ok {
		t.Error("settings.get handler not registered")
	}
	if _, ok := srv.handlers["settings.set"]; !ok {
		t.Error("settings.set handler not registered")
	}
	if _, ok := srv.handlers["settings.reset"]; !ok {
		t.Error("settings.reset handler not registered")
	}
}

// TestSettingsGetHandler verifies settings.get returns default settings
func TestSettingsGetHandler(t *testing.T) {
	// Override home directory for test
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	srv := New(nil, nil, log.New(io.Discard, "", 0))
	if err := RegisterSettingsHandlers(srv); err != nil {
		t.Fatalf("RegisterSettingsHandlers failed: %v", err)
	}

	// Get the handler
	srv.mu.RLock()
	handler := srv.handlers["settings.get"]
	srv.mu.RUnlock()

	// Call handler
	result, err := handler(nil)
	if err != nil {
		t.Fatalf("settings.get failed: %v", err)
	}

	// Verify result is Settings
	settings, ok := result.(*state.Settings)
	if !ok {
		t.Fatalf("Result is not *Settings, got %T", result)
	}

	// Verify default values
	if settings.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", settings.MaxRetries)
	}
	if settings.Theme != "system" {
		t.Errorf("Theme = %q, want \"system\"", settings.Theme)
	}
}

// TestSettingsSetHandler verifies settings.set updates and persists settings
func TestSettingsSetHandler(t *testing.T) {
	// Override home directory for test
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	srv := New(nil, nil, log.New(io.Discard, "", 0))
	if err := RegisterSettingsHandlers(srv); err != nil {
		t.Fatalf("RegisterSettingsHandlers failed: %v", err)
	}

	// Get the handler
	srv.mu.RLock()
	handler := srv.handlers["settings.set"]
	srv.mu.RUnlock()

	// Prepare update params
	updates := map[string]interface{}{
		"maxRetries": 5,
		"theme":      "dark",
	}
	params, _ := json.Marshal(updates)

	// Call handler
	result, err := handler(params)
	if err != nil {
		t.Fatalf("settings.set failed: %v", err)
	}

	// Verify result is updated settings
	settings, ok := result.(*state.Settings)
	if !ok {
		t.Fatalf("Result is not *Settings, got %T", result)
	}

	if settings.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want 5", settings.MaxRetries)
	}
	if settings.Theme != "dark" {
		t.Errorf("Theme = %q, want \"dark\"", settings.Theme)
	}

	// Verify persistence
	configPath := filepath.Join(tmpDir, ".autobmad", "_bmad-output", ".autobmad", "config.json")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Config file not saved: %v", err)
	}
}

// TestSettingsSetInvalidJSON verifies error handling for invalid params
func TestSettingsSetInvalidJSON(t *testing.T) {
	// Override home directory for test
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	srv := New(nil, nil, log.New(io.Discard, "", 0))
	if err := RegisterSettingsHandlers(srv); err != nil {
		t.Fatalf("RegisterSettingsHandlers failed: %v", err)
	}

	srv.mu.RLock()
	handler := srv.handlers["settings.set"]
	srv.mu.RUnlock()

	// Pass invalid JSON
	_, err := handler([]byte("{invalid json"))
	if err == nil {
		t.Error("settings.set should fail with invalid JSON")
	}

	// Verify it's a JSON-RPC error
	if _, ok := err.(*Error); !ok {
		t.Errorf("Error should be *Error, got %T", err)
	}
}

// TestSettingsResetHandler verifies settings.reset restores defaults
func TestSettingsResetHandler(t *testing.T) {
	// Override home directory for test
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	srv := New(nil, nil, log.New(io.Discard, "", 0))
	if err := RegisterSettingsHandlers(srv); err != nil {
		t.Fatalf("RegisterSettingsHandlers failed: %v", err)
	}

	// First, modify settings
	srv.mu.RLock()
	setHandler := srv.handlers["settings.set"]
	resetHandler := srv.handlers["settings.reset"]
	srv.mu.RUnlock()

	updates := map[string]interface{}{
		"maxRetries": 10,
		"theme":      "dark",
	}
	params, _ := json.Marshal(updates)
	if _, err := setHandler(params); err != nil {
		t.Fatalf("Failed to set settings: %v", err)
	}

	// Call reset
	result, err := resetHandler(nil)
	if err != nil {
		t.Fatalf("settings.reset failed: %v", err)
	}

	// Verify result is default settings
	settings, ok := result.(*state.Settings)
	if !ok {
		t.Fatalf("Result is not *Settings, got %T", result)
	}

	if settings.MaxRetries != 3 {
		t.Errorf("After reset, MaxRetries = %d, want 3", settings.MaxRetries)
	}
	if settings.Theme != "system" {
		t.Errorf("After reset, Theme = %q, want \"system\"", settings.Theme)
	}
}
