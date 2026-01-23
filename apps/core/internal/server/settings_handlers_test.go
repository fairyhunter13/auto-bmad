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
	// Use temp directory as project path
	tmpDir := t.TempDir()

	// Create server with same project path (don't use newTestServer which creates its own tmpDir)
	srv := New(nil, nil, log.New(io.Discard, "", 0), tmpDir)

	// Register handlers with project path
	if err := RegisterSettingsHandlers(srv, tmpDir); err != nil {
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
	// Use temp directory as project path
	tmpDir := t.TempDir()

	srv := New(nil, nil, log.New(io.Discard, "", 0), tmpDir)
	if err := RegisterSettingsHandlers(srv, tmpDir); err != nil {
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
	// Use temp directory as project path
	tmpDir := t.TempDir()
	t.Logf("tmpDir = %s", tmpDir)
	t.Logf("srv.ProjectPath() = will check after creation")

	srv := New(nil, nil, log.New(io.Discard, "", 0), tmpDir)
	t.Logf("srv.ProjectPath() = %s", srv.ProjectPath())

	if err := RegisterSettingsHandlers(srv, tmpDir); err != nil {
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

	// Verify persistence (new path: <project>/_bmad-output/.autobmad/config.json)
	configDir := filepath.Join(tmpDir, "_bmad-output", ".autobmad")
	configPath := filepath.Join(configDir, "config.json")

	// List what's actually in tmpDir
	entries, _ := os.ReadDir(tmpDir)
	t.Logf("Contents of %s:", tmpDir)
	for _, e := range entries {
		t.Logf("  - %s (isDir: %v)", e.Name(), e.IsDir())
		if e.IsDir() {
			subEntries, _ := os.ReadDir(filepath.Join(tmpDir, e.Name()))
			for _, se := range subEntries {
				t.Logf("    - %s/%s (isDir: %v)", e.Name(), se.Name(), se.IsDir())
				if se.IsDir() {
					subSubEntries, _ := os.ReadDir(filepath.Join(tmpDir, e.Name(), se.Name()))
					for _, sse := range subSubEntries {
						t.Logf("      - %s/%s/%s", e.Name(), se.Name(), sse.Name())
					}
				}
			}
		}
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Config file not saved: %v", err)
	}
}

// TestSettingsSetInvalidJSON verifies error handling for invalid params
func TestSettingsSetInvalidJSON(t *testing.T) {
	// Use temp directory as project path
	tmpDir := t.TempDir()

	srv := New(nil, nil, log.New(io.Discard, "", 0), tmpDir)
	if err := RegisterSettingsHandlers(srv, tmpDir); err != nil {
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
	// Use temp directory as project path
	tmpDir := t.TempDir()

	srv := New(nil, nil, log.New(io.Discard, "", 0), tmpDir)
	if err := RegisterSettingsHandlers(srv, tmpDir); err != nil {
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

// TestSettingsHandlersPersistenceAcrossRestart verifies settings persist across server restart.
// This integration test ensures the full RPC stack works end-to-end for persistence.
func TestSettingsHandlersPersistenceAcrossRestart(t *testing.T) {
	// Use temp directory as project path (shared across server instances)
	tmpDir := t.TempDir()
	t.Logf("Using project path: %s", tmpDir)

	// PHASE 1: Create first server instance and set settings
	t.Log("PHASE 1: Setting up first server and modifying settings")
	srv1 := New(nil, nil, log.New(io.Discard, "", 0), tmpDir)
	if err := RegisterSettingsHandlers(srv1, tmpDir); err != nil {
		t.Fatalf("RegisterSettingsHandlers (server 1) failed: %v", err)
	}

	// Set custom settings via RPC handler
	srv1.mu.RLock()
	setHandler1 := srv1.handlers["settings.set"]
	srv1.mu.RUnlock()

	customSettings := map[string]interface{}{
		"maxRetries":           7,
		"retryDelay":           8000,
		"theme":                "dark",
		"desktopNotifications": false,
		"lastProjectPath":      "/custom/project/path",
	}
	params, _ := json.Marshal(customSettings)

	result1, err := setHandler1(params)
	if err != nil {
		t.Fatalf("settings.set failed on server 1: %v", err)
	}

	// Verify settings were set on server 1
	settings1, ok := result1.(*state.Settings)
	if !ok {
		t.Fatalf("Result is not *Settings, got %T", result1)
	}

	if settings1.MaxRetries != 7 {
		t.Errorf("Server 1: MaxRetries = %d, want 7", settings1.MaxRetries)
	}
	if settings1.Theme != "dark" {
		t.Errorf("Server 1: Theme = %q, want \"dark\"", settings1.Theme)
	}

	t.Log("PHASE 1 complete: Settings saved on server 1")

	// Simulate server shutdown (srv1 goes out of scope, no cleanup needed for test)
	srv1 = nil
	t.Log("Server 1 stopped (simulating restart)")

	// PHASE 2: Create second server instance (simulates restart)
	t.Log("PHASE 2: Starting second server with same project path")
	srv2 := New(nil, nil, log.New(io.Discard, "", 0), tmpDir)
	if err := RegisterSettingsHandlers(srv2, tmpDir); err != nil {
		t.Fatalf("RegisterSettingsHandlers (server 2) failed: %v", err)
	}

	// Get settings via RPC handler (should load persisted settings)
	srv2.mu.RLock()
	getHandler2 := srv2.handlers["settings.get"]
	srv2.mu.RUnlock()

	result2, err := getHandler2(nil)
	if err != nil {
		t.Fatalf("settings.get failed on server 2: %v", err)
	}

	// Verify settings were persisted and loaded on server 2
	settings2, ok := result2.(*state.Settings)
	if !ok {
		t.Fatalf("Result is not *Settings, got %T", result2)
	}

	// Verify ALL custom settings persisted
	if settings2.MaxRetries != 7 {
		t.Errorf("Server 2: MaxRetries = %d, want 7 (not persisted!)", settings2.MaxRetries)
	}
	if settings2.RetryDelay != 8000 {
		t.Errorf("Server 2: RetryDelay = %d, want 8000 (not persisted!)", settings2.RetryDelay)
	}
	if settings2.Theme != "dark" {
		t.Errorf("Server 2: Theme = %q, want \"dark\" (not persisted!)", settings2.Theme)
	}
	if settings2.DesktopNotifications != false {
		t.Errorf("Server 2: DesktopNotifications = %v, want false (not persisted!)", settings2.DesktopNotifications)
	}
	if settings2.LastProjectPath != "/custom/project/path" {
		t.Errorf("Server 2: LastProjectPath = %q, want \"/custom/project/path\" (not persisted!)", settings2.LastProjectPath)
	}

	t.Log("PHASE 2 complete: All settings correctly loaded on server 2")

	// PHASE 3: Verify config file exists at correct location
	configPath := filepath.Join(tmpDir, "_bmad-output", ".autobmad", "config.json")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Config file not found at %s: %v", configPath, err)
	} else {
		t.Logf("Config file verified at: %s", configPath)
	}

	// PHASE 4: Test that modifications on server 2 also persist
	t.Log("PHASE 4: Testing modifications on server 2")
	srv2.mu.RLock()
	setHandler2 := srv2.handlers["settings.set"]
	srv2.mu.RUnlock()

	secondUpdate := map[string]interface{}{
		"maxRetries": 9,
		"theme":      "light",
	}
	params2, _ := json.Marshal(secondUpdate)

	if _, err := setHandler2(params2); err != nil {
		t.Fatalf("settings.set failed on server 2: %v", err)
	}

	// Create third server to verify second update persisted
	srv3 := New(nil, nil, log.New(io.Discard, "", 0), tmpDir)
	if err := RegisterSettingsHandlers(srv3, tmpDir); err != nil {
		t.Fatalf("RegisterSettingsHandlers (server 3) failed: %v", err)
	}

	srv3.mu.RLock()
	getHandler3 := srv3.handlers["settings.get"]
	srv3.mu.RUnlock()

	result3, err := getHandler3(nil)
	if err != nil {
		t.Fatalf("settings.get failed on server 3: %v", err)
	}

	settings3, ok := result3.(*state.Settings)
	if !ok {
		t.Fatalf("Result is not *Settings, got %T", result3)
	}

	if settings3.MaxRetries != 9 {
		t.Errorf("Server 3: MaxRetries = %d, want 9 (second update not persisted!)", settings3.MaxRetries)
	}
	if settings3.Theme != "light" {
		t.Errorf("Server 3: Theme = %q, want \"light\" (second update not persisted!)", settings3.Theme)
	}
	// Verify first update's non-overwritten values still persist
	if settings3.RetryDelay != 8000 {
		t.Errorf("Server 3: RetryDelay = %d, want 8000 (original value lost!)", settings3.RetryDelay)
	}

	t.Log("PHASE 4 complete: Multiple restarts with updates work correctly")
	t.Log("✅ Integration test PASSED: Settings persistence across restarts verified")
}
