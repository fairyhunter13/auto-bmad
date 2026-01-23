package server

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/state"
)

func TestHandleDetectDependencies(t *testing.T) {
	// Test that the handler returns a valid response
	result, err := handleDetectDependencies(nil)
	if err != nil {
		t.Fatalf("handleDetectDependencies() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("handleDetectDependencies() returned nil result")
	}

	// Convert result to map to check structure
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result to be a map, got %T", result)
	}

	// Check that opencode key exists
	if _, ok := resultMap["opencode"]; !ok {
		t.Error("Expected result to contain 'opencode' key")
	}

	// Check that git key exists
	if _, ok := resultMap["git"]; !ok {
		t.Error("Expected result to contain 'git' key")
	}

	// Verify opencode result has the expected structure
	opencodeData := resultMap["opencode"]
	opencodeBytes, err := json.Marshal(opencodeData)
	if err != nil {
		t.Fatalf("Failed to marshal opencode data: %v", err)
	}

	var opencodeResult struct {
		Found      bool   `json:"found"`
		Version    string `json:"version,omitempty"`
		Path       string `json:"path,omitempty"`
		Compatible bool   `json:"compatible"`
		MinVersion string `json:"minVersion"`
		Error      string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(opencodeBytes, &opencodeResult); err != nil {
		t.Fatalf("Failed to unmarshal opencode result: %v", err)
	}

	// MinVersion should always be set
	if opencodeResult.MinVersion == "" {
		t.Error("Expected opencode.minVersion to be set")
	}

	// If found, should have path and version
	if opencodeResult.Found {
		if opencodeResult.Path == "" {
			t.Error("Expected opencode.path to be set when found")
		}
		if opencodeResult.Version == "" {
			t.Error("Expected opencode.version to be set when found")
		}
	}

	// Verify git result has the expected structure
	gitData := resultMap["git"]
	gitBytes, err := json.Marshal(gitData)
	if err != nil {
		t.Fatalf("Failed to marshal git data: %v", err)
	}

	var gitResult struct {
		Found      bool   `json:"found"`
		Version    string `json:"version,omitempty"`
		Path       string `json:"path,omitempty"`
		Compatible bool   `json:"compatible"`
		MinVersion string `json:"minVersion"`
		Error      string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(gitBytes, &gitResult); err != nil {
		t.Fatalf("Failed to unmarshal git result: %v", err)
	}

	// MinVersion should always be set
	if gitResult.MinVersion == "" {
		t.Error("Expected git.minVersion to be set")
	}

	// If found, should have path and version
	if gitResult.Found {
		if gitResult.Path == "" {
			t.Error("Expected git.path to be set when found")
		}
		if gitResult.Version == "" {
			t.Error("Expected git.version to be set when found")
		}
	}
}

func TestRegisterProjectHandlers(t *testing.T) {
	s := &Server{
		handlers: make(map[string]Handler),
	}

	RegisterProjectHandlers(s)

	// Check that project.detectDependencies is registered
	s.mu.RLock()
	handler1, ok1 := s.handlers["project.detectDependencies"]
	handler2, ok2 := s.handlers["project.scan"]
	s.mu.RUnlock()

	if !ok1 {
		t.Error("Expected project.detectDependencies handler to be registered")
	}

	if handler1 == nil {
		t.Error("Expected detectDependencies handler to not be nil")
	}

	if !ok2 {
		t.Error("Expected project.scan handler to be registered")
	}

	if handler2 == nil {
		t.Error("Expected scan handler to not be nil")
	}
}

func TestHandleProjectScan_ValidPath(t *testing.T) {
	// Create temp directory with BMAD structure
	tmpDir := t.TempDir()

	// Create _bmad folder with manifest
	bmadPath := tmpDir + "/_bmad/_config"
	if err := MkdirAll(bmadPath, 0755); err != nil {
		t.Fatalf("Failed to create bmad path: %v", err)
	}

	manifestFile := bmadPath + "/manifest.yaml"
	manifestContent := []byte("version: 6.1.0\n")
	if err := WriteFile(manifestFile, manifestContent, 0644); err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	// Create params
	params, err := json.Marshal(map[string]string{
		"path": tmpDir,
	})
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	// Call handler
	result, err := handleProjectScan(params)
	if err != nil {
		t.Fatalf("handleProjectScan() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("handleProjectScan() returned nil result")
	}

	// Unmarshal result
	resultBytes, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	var scanResult struct {
		IsBMAD          bool   `json:"isBmad"`
		ProjectType     string `json:"projectType"`
		BmadVersion     string `json:"bmadVersion"`
		BmadCompatible  bool   `json:"bmadCompatible"`
		Path            string `json:"path"`
		HasBmadFolder   bool   `json:"hasBmadFolder"`
		HasOutputFolder bool   `json:"hasOutputFolder"`
	}

	if err := json.Unmarshal(resultBytes, &scanResult); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// Verify result
	if !scanResult.IsBMAD {
		t.Error("Expected isBmad to be true")
	}

	if scanResult.ProjectType != "greenfield" {
		t.Errorf("Expected projectType to be greenfield, got %s", scanResult.ProjectType)
	}

	if scanResult.BmadVersion != "6.1.0" {
		t.Errorf("Expected bmadVersion to be 6.1.0, got %s", scanResult.BmadVersion)
	}

	if !scanResult.BmadCompatible {
		t.Error("Expected bmadCompatible to be true")
	}

	if scanResult.Path != tmpDir {
		t.Errorf("Expected path to be %s, got %s", tmpDir, scanResult.Path)
	}

	if !scanResult.HasBmadFolder {
		t.Error("Expected hasBmadFolder to be true")
	}

	if scanResult.HasOutputFolder {
		t.Error("Expected hasOutputFolder to be false")
	}
}

func TestHandleProjectScan_MissingPath(t *testing.T) {
	// Test with missing path parameter
	params := json.RawMessage(`{}`)

	_, err := handleProjectScan(params)
	if err == nil {
		t.Fatal("Expected error when path is missing")
	}
}

func TestHandleProjectScan_InvalidJSON(t *testing.T) {
	// Test with invalid JSON
	params := json.RawMessage(`{invalid}`)

	_, err := handleProjectScan(params)
	if err == nil {
		t.Fatal("Expected error with invalid JSON")
	}
}

// Helper functions for testing (using os package)
func MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func TestRegisterProjectHandlers_ProfileHandlers(t *testing.T) {
	s := &Server{
		handlers: make(map[string]Handler),
	}

	RegisterProjectHandlers(s)

	s.mu.RLock()
	getHandler, getOk := s.handlers["project.getLastProfile"]
	setHandler, setOk := s.handlers["project.setLastProfile"]
	s.mu.RUnlock()

	if !getOk {
		t.Error("Expected project.getLastProfile handler to be registered")
	}
	if getHandler == nil {
		t.Error("Expected getLastProfile handler to not be nil")
	}

	if !setOk {
		t.Error("Expected project.setLastProfile handler to be registered")
	}
	if setHandler == nil {
		t.Error("Expected setLastProfile handler to not be nil")
	}
}

func TestHandleGetLastProfile_NoSettingsManager(t *testing.T) {
	// Ensure settingsManager is nil for this test
	oldSm := settingsManager
	settingsManager = nil
	defer func() { settingsManager = oldSm }()

	params := json.RawMessage(`{"path": "/test/project"}`)
	_, err := handleGetLastProfile(params)

	if err == nil {
		t.Fatal("Expected error when settingsManager is nil")
	}
}

func TestHandleGetLastProfile_MissingPath(t *testing.T) {
	params := json.RawMessage(`{}`)
	_, err := handleGetLastProfile(params)

	if err == nil {
		t.Fatal("Expected error when path is missing")
	}
}

func TestHandleGetLastProfile_InvalidJSON(t *testing.T) {
	params := json.RawMessage(`{invalid}`)
	_, err := handleGetLastProfile(params)

	if err == nil {
		t.Fatal("Expected error with invalid JSON")
	}
}

func TestHandleSetLastProfile_NoSettingsManager(t *testing.T) {
	// Ensure settingsManager is nil for this test
	oldSm := settingsManager
	settingsManager = nil
	defer func() { settingsManager = oldSm }()

	params := json.RawMessage(`{"path": "/test/project", "profile": "dev"}`)
	_, err := handleSetLastProfile(params)

	if err == nil {
		t.Fatal("Expected error when settingsManager is nil")
	}
}

func TestHandleSetLastProfile_MissingPath(t *testing.T) {
	params := json.RawMessage(`{"profile": "dev"}`)
	_, err := handleSetLastProfile(params)

	if err == nil {
		t.Fatal("Expected error when path is missing")
	}
}

func TestHandleSetLastProfile_InvalidJSON(t *testing.T) {
	params := json.RawMessage(`{invalid}`)
	_, err := handleSetLastProfile(params)

	if err == nil {
		t.Fatal("Expected error with invalid JSON")
	}
}

func TestHandleGetSetLastProfile_WithSettingsManager(t *testing.T) {
	// Create temp directory for settings
	tmpDir := t.TempDir()

	// Initialize settings manager
	sm, err := initTestSettingsManager(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create settings manager: %v", err)
	}

	// Set global settings manager
	oldSm := settingsManager
	settingsManager = sm
	defer func() { settingsManager = oldSm }()

	projectPath1 := "/test/project1"
	projectPath2 := "/test/project2"
	profile1 := "dev"
	profile2 := "prod"

	// Test setting profile for first project
	setParams1, _ := json.Marshal(map[string]string{
		"path":    projectPath1,
		"profile": profile1,
	})
	result1, err := handleSetLastProfile(setParams1)
	if err != nil {
		t.Fatalf("handleSetLastProfile() for project1 returned error: %v", err)
	}
	if result1.(map[string]string)["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", result1)
	}

	// Test setting profile for second project
	setParams2, _ := json.Marshal(map[string]string{
		"path":    projectPath2,
		"profile": profile2,
	})
	result2, err := handleSetLastProfile(setParams2)
	if err != nil {
		t.Fatalf("handleSetLastProfile() for project2 returned error: %v", err)
	}
	if result2.(map[string]string)["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", result2)
	}

	// Test getting profile for first project
	getParams1, _ := json.Marshal(map[string]string{"path": projectPath1})
	getResult1, err := handleGetLastProfile(getParams1)
	if err != nil {
		t.Fatalf("handleGetLastProfile() for project1 returned error: %v", err)
	}
	if getResult1.(map[string]string)["profile"] != profile1 {
		t.Errorf("Expected profile '%s', got '%s'", profile1, getResult1.(map[string]string)["profile"])
	}

	// Test getting profile for second project
	getParams2, _ := json.Marshal(map[string]string{"path": projectPath2})
	getResult2, err := handleGetLastProfile(getParams2)
	if err != nil {
		t.Fatalf("handleGetLastProfile() for project2 returned error: %v", err)
	}
	if getResult2.(map[string]string)["profile"] != profile2 {
		t.Errorf("Expected profile '%s', got '%s'", profile2, getResult2.(map[string]string)["profile"])
	}

	// Test getting profile for non-existent project (should return empty)
	getParams3, _ := json.Marshal(map[string]string{"path": "/non/existent"})
	getResult3, err := handleGetLastProfile(getParams3)
	if err != nil {
		t.Fatalf("handleGetLastProfile() for non-existent path returned error: %v", err)
	}
	if getResult3.(map[string]string)["profile"] != "" {
		t.Errorf("Expected empty profile for non-existent path, got '%s'", getResult3.(map[string]string)["profile"])
	}
}

func TestHandleSetLastProfile_UpdatesExistingProfile(t *testing.T) {
	// Create temp directory for settings
	tmpDir := t.TempDir()

	// Initialize settings manager
	sm, err := initTestSettingsManager(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create settings manager: %v", err)
	}

	// Set global settings manager
	oldSm := settingsManager
	settingsManager = sm
	defer func() { settingsManager = oldSm }()

	projectPath := "/test/project"
	profile1 := "dev"
	profile2 := "staging"

	// Set initial profile
	setParams1, _ := json.Marshal(map[string]string{
		"path":    projectPath,
		"profile": profile1,
	})
	_, err = handleSetLastProfile(setParams1)
	if err != nil {
		t.Fatalf("handleSetLastProfile() initial set returned error: %v", err)
	}

	// Update profile
	setParams2, _ := json.Marshal(map[string]string{
		"path":    projectPath,
		"profile": profile2,
	})
	_, err = handleSetLastProfile(setParams2)
	if err != nil {
		t.Fatalf("handleSetLastProfile() update returned error: %v", err)
	}

	// Verify update
	getParams, _ := json.Marshal(map[string]string{"path": projectPath})
	getResult, err := handleGetLastProfile(getParams)
	if err != nil {
		t.Fatalf("handleGetLastProfile() returned error: %v", err)
	}
	if getResult.(map[string]string)["profile"] != profile2 {
		t.Errorf("Expected profile '%s', got '%s'", profile2, getResult.(map[string]string)["profile"])
	}
}

// initTestSettingsManager creates a settings manager for testing
func initTestSettingsManager(projectPath string) (*state.StateManager, error) {
	return state.NewStateManager(projectPath)
}
