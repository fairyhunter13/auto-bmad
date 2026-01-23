package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/project"
)

func TestHandleGetRecent(t *testing.T) {
	// Initialize recent manager
	project.InitRecentManager()

	// Add some test projects
	rm := project.GetRecentManager()
	rm.Add("/home/user/project1")
	rm.Add("/home/user/project2")

	// Call handler
	result, err := handleGetRecent(nil)
	if err != nil {
		t.Fatalf("handleGetRecent failed: %v", err)
	}

	projects, ok := result.([]project.RecentProject)
	if !ok {
		t.Fatalf("Expected []RecentProject, got %T", result)
	}

	if len(projects) < 2 {
		t.Errorf("Expected at least 2 projects, got %d", len(projects))
	}
}

func TestHandleAddRecent(t *testing.T) {
	project.InitRecentManager()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "addrecent-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Get absolute path
	absPath, _ := filepath.Abs(tempDir)

	params := AddRecentParams{
		Path: absPath,
	}
	paramsJSON, _ := json.Marshal(params)

	result, err := handleAddRecent(paramsJSON)
	if err != nil {
		t.Fatalf("handleAddRecent failed: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verify project was added (using the validated path which may be resolved)
	rm := project.GetRecentManager()
	projects, _ := rm.GetAll()

	found := false
	for _, p := range projects {
		// Compare resolved paths
		if p.Path == absPath || filepath.Clean(p.Path) == filepath.Clean(absPath) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Project was not added to recent list. Added path: %s, Projects: %+v", absPath, projects)
	}
}

func TestHandleAddRecent_EmptyPath(t *testing.T) {
	params := AddRecentParams{
		Path: "",
	}
	paramsJSON, _ := json.Marshal(params)

	_, err := handleAddRecent(paramsJSON)
	if err == nil {
		t.Error("Expected error for empty path")
	}
}

func TestHandleRemoveRecent(t *testing.T) {
	project.InitRecentManager()
	rm := project.GetRecentManager()

	// Add a project first
	rm.Add("/home/user/to-remove")

	params := RemoveRecentParams{
		Path: "/home/user/to-remove",
	}
	paramsJSON, _ := json.Marshal(params)

	result, err := handleRemoveRecent(paramsJSON)
	if err != nil {
		t.Fatalf("handleRemoveRecent failed: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verify project was removed
	projects, _ := rm.GetAll()
	for _, p := range projects {
		if p.Path == "/home/user/to-remove" {
			t.Error("Project should have been removed")
		}
	}
}

func TestHandleSetContext(t *testing.T) {
	project.InitRecentManager()
	rm := project.GetRecentManager()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "setcontext-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Get absolute path and resolve it (like the validator does)
	absPath, _ := filepath.Abs(tempDir)
	resolvedPath, _ := filepath.EvalSymlinks(absPath)

	// Add a project first (directly to recent manager, bypassing validation)
	rm.Add(resolvedPath)

	params := SetContextParams{
		Path:    absPath,
		Context: "This is a test project for Auto-BMAD",
	}
	paramsJSON, _ := json.Marshal(params)

	result, err := handleSetContext(paramsJSON)
	if err != nil {
		t.Fatalf("handleSetContext failed: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verify context was set (check both original and resolved paths)
	projects, _ := rm.GetAll()
	for _, p := range projects {
		if p.Path == resolvedPath || p.Path == absPath {
			if p.Context != "This is a test project for Auto-BMAD" {
				t.Errorf("Expected context to be set, got: %s", p.Context)
			}
			return
		}
	}

	t.Errorf("Project not found after setting context. Looking for: %s or %s, Projects: %+v", absPath, resolvedPath, projects)
}

func TestHandleSetContext_NonexistentProject(t *testing.T) {
	project.InitRecentManager()

	params := SetContextParams{
		Path:    "/home/user/nonexistent",
		Context: "This should fail",
	}
	paramsJSON, _ := json.Marshal(params)

	_, err := handleSetContext(paramsJSON)
	if err == nil {
		t.Error("Expected error for nonexistent project")
	}
}
