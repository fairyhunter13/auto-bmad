package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNewStateManager verifies StateManager initialization
func TestNewStateManager(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")

	// Create StateManager
	sm, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() failed: %v", err)
	}
	if sm == nil {
		t.Fatal("NewStateManager() returned nil")
	}

	// Verify config directory was created
	configDir := filepath.Join(projectPath, "_bmad-output", ".autobmad")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Config directory not created: %s", configDir)
	}

	// Verify default settings are loaded
	settings := sm.Get()
	if settings.MaxRetries != 3 {
		t.Error("Default settings not loaded properly")
	}
}

// TestStateManagerSaveLoad verifies settings persistence
func TestStateManagerSaveLoad(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")

	// Create first manager instance
	sm1, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() failed: %v", err)
	}

	// Modify settings
	updates := map[string]interface{}{
		"maxRetries":           5,
		"desktopNotifications": false,
		"theme":                "dark",
	}

	if err := sm1.Set(updates); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	// Create second manager instance (simulates restart)
	sm2, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() on restart failed: %v", err)
	}

	// Verify settings were persisted
	settings := sm2.Get()
	if settings.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want 5", settings.MaxRetries)
	}
	if settings.DesktopNotifications != false {
		t.Error("DesktopNotifications should be false")
	}
	if settings.Theme != "dark" {
		t.Errorf("Theme = %q, want \"dark\"", settings.Theme)
	}
}

// TestStateManagerAtomicWrite verifies atomic write behavior
func TestStateManagerAtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")

	sm, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() failed: %v", err)
	}

	// Make a change
	updates := map[string]interface{}{"maxRetries": 7}
	if err := sm.Set(updates); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	// Verify temp file doesn't exist (should be cleaned up after rename)
	configPath := filepath.Join(projectPath, "_bmad-output", ".autobmad", "config.json")
	tempPath := configPath + ".tmp"

	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		t.Error("Temp file still exists after save - atomic write not working")
	}

	// Verify final file exists
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Config file doesn't exist: %v", err)
	}
}

// TestStateManagerSavePerformance verifies save completes within 1 second (NFR-P6)
func TestStateManagerSavePerformance(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")

	sm, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() failed: %v", err)
	}

	// Measure save time
	start := time.Now()
	updates := map[string]interface{}{
		"maxRetries": 8,
		"theme":      "light",
	}
	if err := sm.Set(updates); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}
	duration := time.Since(start)

	// Verify save time < 1 second (NFR-P6)
	if duration > time.Second {
		t.Errorf("Save took %v, must be < 1 second (NFR-P6)", duration)
	}
}

// TestStateManagerGet verifies Get returns a copy, not reference
func TestStateManagerGet(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")

	sm, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() failed: %v", err)
	}

	// Get settings
	settings1 := sm.Get()
	settings1.MaxRetries = 99 // Try to mutate

	// Get settings again
	settings2 := sm.Get()

	// Verify mutation didn't affect internal state
	if settings2.MaxRetries == 99 {
		t.Error("Get() returned a reference instead of a copy - internal state was mutated")
	}
	if settings2.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3 (default)", settings2.MaxRetries)
	}
}

// TestStateManagerReset verifies Reset restores defaults
func TestStateManagerReset(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")

	sm, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() failed: %v", err)
	}

	// Modify settings
	updates := map[string]interface{}{
		"maxRetries": 10,
		"theme":      "dark",
	}
	if err := sm.Set(updates); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	// Reset
	if err := sm.Reset(); err != nil {
		t.Fatalf("Reset() failed: %v", err)
	}

	// Verify defaults are restored
	settings := sm.Get()
	if settings.MaxRetries != 3 {
		t.Errorf("After Reset, MaxRetries = %d, want 3", settings.MaxRetries)
	}
	if settings.Theme != "system" {
		t.Errorf("After Reset, Theme = %q, want \"system\"", settings.Theme)
	}

	// Verify reset is persisted
	sm2, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() after reset failed: %v", err)
	}

	settings2 := sm2.Get()
	if settings2.MaxRetries != 3 {
		t.Error("Reset was not persisted")
	}
}

// TestStateManagerSetValidation verifies Set handles various field types
func TestStateManagerSetValidation(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")

	sm, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() failed: %v", err)
	}

	// Test all supported field types
	updates := map[string]interface{}{
		"maxRetries":           7,
		"retryDelay":           10000,
		"desktopNotifications": true,
		"soundEnabled":         true,
		"stepTimeoutDefault":   600000,
		"heartbeatInterval":    30000,
		"theme":                "dark",
		"showDebugOutput":      true,
		"lastProjectPath":      "/home/test/project",
		"recentProjectsMax":    5,
	}

	if err := sm.Set(updates); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	// Verify all fields were set
	settings := sm.Get()
	if settings.MaxRetries != 7 {
		t.Errorf("MaxRetries = %d, want 7", settings.MaxRetries)
	}
	if settings.RetryDelay != 10000 {
		t.Errorf("RetryDelay = %d, want 10000", settings.RetryDelay)
	}
	if settings.DesktopNotifications != true {
		t.Error("DesktopNotifications should be true")
	}
	if settings.SoundEnabled != true {
		t.Error("SoundEnabled should be true")
	}
	if settings.Theme != "dark" {
		t.Errorf("Theme = %q, want \"dark\"", settings.Theme)
	}
	if settings.ShowDebugOutput != true {
		t.Error("ShowDebugOutput should be true")
	}
	if settings.LastProjectPath != "/home/test/project" {
		t.Errorf("LastProjectPath = %q, want \"/home/test/project\"", settings.LastProjectPath)
	}
}

// TestStateManagerLoadWithMissingFile verifies defaults are used when file doesn't exist
func TestStateManagerLoadWithMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "nonexistent-project")

	sm, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() should succeed even with missing file: %v", err)
	}

	// Verify defaults are loaded
	settings := sm.Get()
	if settings.MaxRetries != 3 {
		t.Error("Should use defaults when config file doesn't exist")
	}
}

// TestStateManagerLoadWithCorruptedFile verifies error handling for corrupted JSON
func TestStateManagerLoadWithCorruptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	configDir := filepath.Join(projectPath, "_bmad-output", ".autobmad")

	// Create config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Write corrupted JSON
	configPath := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(configPath, []byte("{invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write corrupted file: %v", err)
	}

	// NewStateManager should return error for corrupted file
	_, err := NewStateManager(projectPath)
	if err == nil {
		t.Error("NewStateManager() should fail with corrupted JSON file")
	}
}

// TestStateManagerProjectProfiles verifies project profile persistence
func TestStateManagerProjectProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")

	sm, err := NewStateManager(projectPath)
	if err != nil {
		t.Fatalf("NewStateManager() failed: %v", err)
	}

	// Add project profiles via Set
	settings := sm.Get()
	settings.ProjectProfiles["/project1"] = "claude-sonnet"
	settings.ProjectProfiles["/project2"] = "gpt-4"

	// Save entire settings (simulating frontend update)
	data, _ := json.Marshal(settings)
	var updates map[string]interface{}
	json.Unmarshal(data, &updates)

	if err := sm.Set(updates); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	// Verify persistence
	sm2, _ := NewStateManager(projectPath)
	restored := sm2.Get()

	if restored.ProjectProfiles["/project1"] != "claude-sonnet" {
		t.Error("Project profile 1 not persisted")
	}
	if restored.ProjectProfiles["/project2"] != "gpt-4" {
		t.Error("Project profile 2 not persisted")
	}
}
