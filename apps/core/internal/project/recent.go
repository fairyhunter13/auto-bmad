package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// RecentProject represents a recently opened BMAD project
type RecentProject struct {
	Path       string    `json:"path"`
	Name       string    `json:"name"`
	LastOpened time.Time `json:"lastOpened"`
	Context    string    `json:"context,omitempty"`
}

// RecentManager manages the list of recently opened projects
type RecentManager struct {
	configPath string
	projects   []RecentProject
	maxRecent  int
	mu         sync.RWMutex
}

// NewRecentManager creates a new RecentManager instance
func NewRecentManager(configPath string, maxRecent int) *RecentManager {
	rm := &RecentManager{
		configPath: configPath,
		maxRecent:  maxRecent,
		projects:   []RecentProject{},
	}

	// Load existing projects if file exists
	rm.load()

	return rm
}

// Add adds a project to the recent list, or updates its timestamp if it exists
func (rm *RecentManager) Add(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Remove if already exists (to update timestamp and move to front)
	rm.removeWithoutLock(path)

	// Create new entry
	project := RecentProject{
		Path:       path,
		Name:       filepath.Base(path),
		LastOpened: time.Now(),
	}

	// Add to front of list
	rm.projects = append([]RecentProject{project}, rm.projects...)

	// Trim to max size
	if len(rm.projects) > rm.maxRecent {
		rm.projects = rm.projects[:rm.maxRecent]
	}

	return rm.save()
}

// Remove removes a project from the recent list
func (rm *RecentManager) Remove(path string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.removeWithoutLock(path)
	return rm.save()
}

// removeWithoutLock removes a project without acquiring lock (internal use)
func (rm *RecentManager) removeWithoutLock(path string) {
	for i, p := range rm.projects {
		if p.Path == path {
			rm.projects = append(rm.projects[:i], rm.projects[i+1:]...)
			return
		}
	}
}

// SetContext sets the context description for a project
func (rm *RecentManager) SetContext(path string, context string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for i := range rm.projects {
		if rm.projects[i].Path == path {
			rm.projects[i].Context = context
			return rm.save()
		}
	}

	return fmt.Errorf("project not found: %s", path)
}

// GetAll returns all recent projects, sorted by most recent first
func (rm *RecentManager) GetAll() ([]RecentProject, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]RecentProject, len(rm.projects))
	copy(result, rm.projects)

	return result, nil
}

// load loads projects from the config file
func (rm *RecentManager) load() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// If file doesn't exist, start with empty list
	if _, err := os.Stat(rm.configPath); os.IsNotExist(err) {
		rm.projects = []RecentProject{}
		return nil
	}

	data, err := os.ReadFile(rm.configPath)
	if err != nil {
		return fmt.Errorf("reading recent projects file: %w", err)
	}

	if len(data) == 0 {
		rm.projects = []RecentProject{}
		return nil
	}

	var projects []RecentProject
	if err := json.Unmarshal(data, &projects); err != nil {
		// If file is corrupted, start fresh
		rm.projects = []RecentProject{}
		return fmt.Errorf("parsing recent projects (starting fresh): %w", err)
	}

	rm.projects = projects
	return nil
}

// save persists projects to the config file
func (rm *RecentManager) save() error {
	// Ensure directory exists
	dir := filepath.Dir(rm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(rm.projects, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling recent projects: %w", err)
	}

	if err := os.WriteFile(rm.configPath, data, 0644); err != nil {
		return fmt.Errorf("writing recent projects file: %w", err)
	}

	return nil
}
