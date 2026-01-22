package project

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRecentManager_Add(t *testing.T) {
	// Create temp directory for test config
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "recent.json")

	rm := NewRecentManager(configPath, 5)

	// Add first project
	err := rm.Add("/home/user/project1")
	if err != nil {
		t.Fatalf("Failed to add project: %v", err)
	}

	projects, err := rm.GetAll()
	if err != nil {
		t.Fatalf("Failed to get projects: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("Expected 1 project, got %d", len(projects))
	}

	if projects[0].Path != "/home/user/project1" {
		t.Errorf("Expected path /home/user/project1, got %s", projects[0].Path)
	}

	if projects[0].Name != "project1" {
		t.Errorf("Expected name project1, got %s", projects[0].Name)
	}
}

func TestRecentManager_Add_UpdatesTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "recent.json")

	rm := NewRecentManager(configPath, 5)

	// Add project first time
	err := rm.Add("/home/user/project1")
	if err != nil {
		t.Fatalf("Failed to add project: %v", err)
	}

	projects1, _ := rm.GetAll()
	firstTime := projects1[0].LastOpened

	// Wait a bit and add again
	time.Sleep(10 * time.Millisecond)

	err = rm.Add("/home/user/project1")
	if err != nil {
		t.Fatalf("Failed to re-add project: %v", err)
	}

	projects2, _ := rm.GetAll()
	secondTime := projects2[0].LastOpened

	if !secondTime.After(firstTime) {
		t.Error("Expected timestamp to be updated")
	}

	// Should still have only 1 project
	if len(projects2) != 1 {
		t.Errorf("Expected 1 project after re-add, got %d", len(projects2))
	}
}

func TestRecentManager_Add_RespectMaxLimit(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "recent.json")

	rm := NewRecentManager(configPath, 3)

	// Add 5 projects (max is 3)
	for i := 1; i <= 5; i++ {
		path := filepath.Join("/home/user", "project", string(rune('0'+i)))
		err := rm.Add(path)
		if err != nil {
			t.Fatalf("Failed to add project %d: %v", i, err)
		}
	}

	projects, _ := rm.GetAll()

	if len(projects) != 3 {
		t.Errorf("Expected max 3 projects, got %d", len(projects))
	}

	// Most recent should be first
	if projects[0].Path != "/home/user/project/5" {
		t.Errorf("Expected most recent project first, got %s", projects[0].Path)
	}
}

func TestRecentManager_Remove(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "recent.json")

	rm := NewRecentManager(configPath, 5)

	// Add projects
	rm.Add("/home/user/project1")
	rm.Add("/home/user/project2")
	rm.Add("/home/user/project3")

	// Remove middle one
	err := rm.Remove("/home/user/project2")
	if err != nil {
		t.Fatalf("Failed to remove project: %v", err)
	}

	projects, _ := rm.GetAll()

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects after removal, got %d", len(projects))
	}

	// Verify the right one was removed
	for _, p := range projects {
		if p.Path == "/home/user/project2" {
			t.Error("Project2 should have been removed")
		}
	}
}

func TestRecentManager_SetContext(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "recent.json")

	rm := NewRecentManager(configPath, 5)

	// Add project
	rm.Add("/home/user/project1")

	// Set context
	err := rm.SetContext("/home/user/project1", "This is my awesome project")
	if err != nil {
		t.Fatalf("Failed to set context: %v", err)
	}

	projects, _ := rm.GetAll()

	if projects[0].Context != "This is my awesome project" {
		t.Errorf("Expected context to be set, got: %s", projects[0].Context)
	}
}

func TestRecentManager_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "recent.json")

	// First manager instance
	rm1 := NewRecentManager(configPath, 5)
	rm1.Add("/home/user/project1")
	rm1.Add("/home/user/project2")

	// Create second manager instance (simulates app restart)
	rm2 := NewRecentManager(configPath, 5)
	projects, err := rm2.GetAll()
	if err != nil {
		t.Fatalf("Failed to load projects: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("Expected 2 persisted projects, got %d", len(projects))
	}
}

func TestRecentManager_InvalidPath(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "recent.json")

	rm := NewRecentManager(configPath, 5)

	// Try to add empty path
	err := rm.Add("")
	if err == nil {
		t.Error("Expected error for empty path")
	}
}

func TestRecentManager_CorruptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "recent.json")

	// Write corrupted JSON
	os.WriteFile(configPath, []byte("invalid json {{{"), 0644)

	// Should handle gracefully
	rm := NewRecentManager(configPath, 5)
	projects, err := rm.GetAll()

	// Should return empty list or error, not crash
	if err == nil && len(projects) != 0 {
		t.Error("Expected empty list or error for corrupted file")
	}
}
