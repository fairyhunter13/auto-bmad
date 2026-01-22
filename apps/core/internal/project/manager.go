package project

import (
	"os"
	"path/filepath"
)

var (
	// Global recent manager instance
	recentManager *RecentManager
)

// InitRecentManager initializes the global recent manager
func InitRecentManager() error {
	// Use user config directory for storing recent projects
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "auto-bmad")
	configPath := filepath.Join(configDir, "recent-projects.json")

	recentManager = NewRecentManager(configPath, 10) // Keep last 10 projects
	return nil
}

// GetRecentManager returns the global recent manager instance
func GetRecentManager() *RecentManager {
	if recentManager == nil {
		// Auto-initialize if not done yet
		InitRecentManager()
	}
	return recentManager
}
