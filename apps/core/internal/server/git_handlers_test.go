package server

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestRegisterGitHandlers(t *testing.T) {
	srv := newTestServer(t, nil, nil, log.New(io.Discard, "", 0))
	RegisterGitHandlers(srv)

	srv.mu.RLock()
	defer srv.mu.RUnlock()

	if _, ok := srv.handlers["git.getRepoStatus"]; !ok {
		t.Error("git.getRepoStatus handler not registered")
	}
}

func TestHandleGetRepoStatus_InvalidJSON(t *testing.T) {
	result, err := handleGetRepoStatus([]byte("{invalid json"))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestHandleGetRepoStatus_MissingPath(t *testing.T) {
	params := []byte(`{}`)
	result, err := handleGetRepoStatus(params)
	if err == nil {
		t.Error("Expected error for missing path")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestHandleGetRepoStatus_NonGitDirectory(t *testing.T) {
	// Create a non-git temp directory
	tmpDir := t.TempDir()

	params, _ := json.Marshal(map[string]string{"path": tmpDir})
	result, err := handleGetRepoStatus(params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Marshal to JSON and unmarshal to inspect
	resultBytes, _ := json.Marshal(result)
	var status struct {
		IsGitRepo  bool   `json:"isGitRepo"`
		Branch     string `json:"branch,omitempty"`
		HasChanges bool   `json:"hasChanges"`
		Error      string `json:"error,omitempty"`
	}
	if err := json.Unmarshal(resultBytes, &status); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if status.IsGitRepo {
		t.Error("Expected isGitRepo to be false for non-git directory")
	}
}

func TestHandleGetRepoStatus_GitRepository(t *testing.T) {
	// Create a git repo in temp directory
	tmpDir := t.TempDir()

	// Initialize git repo
	if err := os.MkdirAll(filepath.Join(tmpDir, ".git"), 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	// Create minimal git structure
	if err := os.WriteFile(filepath.Join(tmpDir, ".git", "HEAD"), []byte("ref: refs/heads/main\n"), 0644); err != nil {
		t.Fatalf("Failed to create HEAD file: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, ".git", "refs", "heads"), 0755); err != nil {
		t.Fatalf("Failed to create refs/heads directory: %v", err)
	}

	params, _ := json.Marshal(map[string]string{"path": tmpDir})
	result, err := handleGetRepoStatus(params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// The result should indicate it's a git repo (even if not fully initialized)
	if result == nil {
		t.Error("Expected non-nil result")
	}
}

func TestHandleGetRepoStatus_CurrentRepo(t *testing.T) {
	// Test with the actual auto-bmad repository (we know this is a git repo)
	// Get the repo root by going up from the test file location
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Navigate up to find the .git directory
	repoRoot := wd
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(repoRoot, ".git")); err == nil {
			break
		}
		repoRoot = filepath.Dir(repoRoot)
	}

	// Skip if we couldn't find a git repo
	if _, err := os.Stat(filepath.Join(repoRoot, ".git")); err != nil {
		t.Skip("Could not find git repository root")
	}

	params, _ := json.Marshal(map[string]string{"path": repoRoot})
	result, err := handleGetRepoStatus(params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Marshal result to JSON to inspect it
	resultBytes, _ := json.Marshal(result)
	t.Logf("Result: %s", string(resultBytes))

	// The result should be a GitRepoStatus
	if result == nil {
		t.Error("Expected non-nil result")
	}
}
