package state

import (
	"encoding/json"
	"testing"
)

// TestDefaultSettings verifies that DefaultSettings returns proper defaults
func TestDefaultSettings(t *testing.T) {
	settings := DefaultSettings()

	if settings == nil {
		t.Fatal("DefaultSettings() returned nil")
	}

	// Verify retry settings
	if settings.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", settings.MaxRetries)
	}
	if settings.RetryDelay != 5000 {
		t.Errorf("RetryDelay = %d, want 5000", settings.RetryDelay)
	}

	// Verify notification settings
	if settings.DesktopNotifications != true {
		t.Error("DesktopNotifications should be true by default")
	}
	if settings.SoundEnabled != false {
		t.Error("SoundEnabled should be false by default")
	}

	// Verify timeout settings
	if settings.StepTimeoutDefault != 300000 {
		t.Errorf("StepTimeoutDefault = %d, want 300000", settings.StepTimeoutDefault)
	}
	if settings.HeartbeatInterval != 60000 {
		t.Errorf("HeartbeatInterval = %d, want 60000", settings.HeartbeatInterval)
	}

	// Verify UI preferences
	if settings.Theme != "system" {
		t.Errorf("Theme = %q, want \"system\"", settings.Theme)
	}
	if settings.ShowDebugOutput != false {
		t.Error("ShowDebugOutput should be false by default")
	}

	// Verify project memory
	if settings.ProjectProfiles == nil {
		t.Error("ProjectProfiles should be initialized, not nil")
	}
	if len(settings.ProjectProfiles) != 0 {
		t.Errorf("ProjectProfiles should be empty, got %d entries", len(settings.ProjectProfiles))
	}
	if settings.RecentProjectsMax != 10 {
		t.Errorf("RecentProjectsMax = %d, want 10", settings.RecentProjectsMax)
	}
}

// TestSettingsJSONMarshaling verifies settings can be marshaled to JSON
func TestSettingsJSONMarshaling(t *testing.T) {
	settings := DefaultSettings()
	settings.LastProjectPath = "/home/user/project"
	settings.ProjectProfiles["/home/user/project"] = "claude-sonnet"
	settings.MaxRetries = 5

	// Marshal to JSON
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("Failed to marshal settings: %v", err)
	}

	// Unmarshal back
	var restored Settings
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Failed to unmarshal settings: %v", err)
	}

	// Verify field names are camelCase in JSON
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// Check camelCase field names
	if _, ok := jsonMap["maxRetries"]; !ok {
		t.Error("JSON should have 'maxRetries' field (camelCase)")
	}
	if _, ok := jsonMap["desktopNotifications"]; !ok {
		t.Error("JSON should have 'desktopNotifications' field (camelCase)")
	}

	// Verify values are preserved
	if restored.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d after unmarshal, want 5", restored.MaxRetries)
	}
	if restored.LastProjectPath != "/home/user/project" {
		t.Errorf("LastProjectPath = %q, want %q", restored.LastProjectPath, "/home/user/project")
	}
	if restored.ProjectProfiles["/home/user/project"] != "claude-sonnet" {
		t.Error("ProjectProfiles not preserved after unmarshal")
	}
}
