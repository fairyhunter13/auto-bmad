package server

import (
	"encoding/json"
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

	params := AddRecentParams{
		Path: "/home/user/test-project",
	}
	paramsJSON, _ := json.Marshal(params)

	result, err := handleAddRecent(paramsJSON)
	if err != nil {
		t.Fatalf("handleAddRecent failed: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verify project was added
	rm := project.GetRecentManager()
	projects, _ := rm.GetAll()

	found := false
	for _, p := range projects {
		if p.Path == "/home/user/test-project" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Project was not added to recent list")
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

	// Add a project first
	rm.Add("/home/user/my-project")

	params := SetContextParams{
		Path:    "/home/user/my-project",
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

	// Verify context was set
	projects, _ := rm.GetAll()
	for _, p := range projects {
		if p.Path == "/home/user/my-project" {
			if p.Context != "This is a test project for Auto-BMAD" {
				t.Errorf("Expected context to be set, got: %s", p.Context)
			}
			return
		}
	}

	t.Error("Project not found after setting context")
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
