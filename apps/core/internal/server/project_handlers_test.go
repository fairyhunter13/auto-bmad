package server

import (
	"encoding/json"
	"os"
	"testing"
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
