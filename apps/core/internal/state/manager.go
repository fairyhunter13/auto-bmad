package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// StateManager handles settings persistence and retrieval.
// It provides thread-safe access to user settings stored in config.json.
type StateManager struct {
	settings   *Settings
	configPath string
	mu         sync.RWMutex
}

// NewStateManager creates a new StateManager instance.
// It creates the config directory if it doesn't exist and loads existing settings.
// If no settings file exists, default settings are used.
func NewStateManager(projectPath string) (*StateManager, error) {
	configDir := filepath.Join(projectPath, "_bmad-output", ".autobmad")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("creating config directory: %w", err)
	}

	sm := &StateManager{
		configPath: filepath.Join(configDir, "config.json"),
		settings:   DefaultSettings(),
	}

	// Load existing settings (ignore if file doesn't exist)
	if err := sm.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading settings: %w", err)
	}

	return sm, nil
}

// load reads settings from the config file.
// If the file doesn't exist, default settings are kept.
func (sm *StateManager) load() error {
	data, err := os.ReadFile(sm.configPath)
	if err != nil {
		return err
	}

	// Start with defaults to ensure all fields have values
	settings := DefaultSettings()
	if err := json.Unmarshal(data, settings); err != nil {
		return err
	}

	sm.settings = settings
	return nil
}

// save writes settings to the config file using atomic write.
// It writes to a temporary file first, then renames it to ensure atomicity.
func (sm *StateManager) save() error {
	data, err := json.MarshalIndent(sm.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling settings: %w", err)
	}

	// Atomic write: write to temp file, then rename
	tempPath := sm.configPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("writing temp file: %w", err)
	}

	if err := os.Rename(tempPath, sm.configPath); err != nil {
		return fmt.Errorf("renaming temp file: %w", err)
	}

	return nil
}

// Get returns a copy of the current settings.
// The returned settings can be safely modified without affecting internal state.
func (sm *StateManager) Get() *Settings {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Return a copy to prevent mutation
	settingsCopy := *sm.settings
	// Deep copy the map
	settingsCopy.ProjectProfiles = make(map[string]string, len(sm.settings.ProjectProfiles))
	for k, v := range sm.settings.ProjectProfiles {
		settingsCopy.ProjectProfiles[k] = v
	}
	return &settingsCopy
}

// Set updates settings with the provided values and saves to disk.
// Only recognized fields are updated; unknown fields are ignored.
// Input validation is performed to prevent invalid or malicious values.
// Validation is atomic: all fields are validated before any changes are applied.
func (sm *StateManager) Set(updates map[string]interface{}) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// PHASE 1: Validate ALL fields first (atomic validation)
	for key, value := range updates {
		switch key {
		case "maxRetries":
			val := toInt(value)
			if val < 0 || val > 10 {
				return fmt.Errorf("maxRetries must be 0-10, got %d", val)
			}

		case "retryDelay":
			val := toInt(value)
			if val < 0 || val > 60000 {
				return fmt.Errorf("retryDelay must be 0-60000ms, got %d", val)
			}

		case "stepTimeoutDefault":
			val := toInt(value)
			if val < 1000 || val > 3600000 {
				return fmt.Errorf("stepTimeoutDefault must be 1000-3600000ms (1s-1h), got %d", val)
			}

		case "heartbeatInterval":
			val := toInt(value)
			if val < 1000 || val > 300000 {
				return fmt.Errorf("heartbeatInterval must be 1000-300000ms (1s-5min), got %d", val)
			}

		case "theme":
			if v, ok := value.(string); ok {
				if v != "light" && v != "dark" && v != "system" {
					return fmt.Errorf("theme must be 'light', 'dark', or 'system', got %q", v)
				}
			}

		case "lastProjectPath":
			if v, ok := value.(string); ok {
				// Prevent path traversal attacks
				if strings.Contains(v, "..") {
					return fmt.Errorf("lastProjectPath contains path traversal sequence '..'")
				}
			}

		case "recentProjectsMax":
			val := toInt(value)
			if val < 1 || val > 50 {
				return fmt.Errorf("recentProjectsMax must be 1-50, got %d", val)
			}

		case "projectProfiles":
			if v, ok := value.(map[string]interface{}); ok {
				for k := range v {
					// Validate project path doesn't contain path traversal
					if strings.Contains(k, "..") {
						return fmt.Errorf("projectProfiles key contains path traversal sequence '..'")
					}
				}
			}
		}
	}

	// PHASE 2: Apply updates (only after all validations pass)
	for key, value := range updates {
		switch key {
		case "maxRetries":
			sm.settings.MaxRetries = toInt(value)
		case "retryDelay":
			sm.settings.RetryDelay = toInt(value)
		case "desktopNotifications":
			if v, ok := value.(bool); ok {
				sm.settings.DesktopNotifications = v
			}
		case "soundEnabled":
			if v, ok := value.(bool); ok {
				sm.settings.SoundEnabled = v
			}
		case "stepTimeoutDefault":
			sm.settings.StepTimeoutDefault = toInt(value)
		case "heartbeatInterval":
			sm.settings.HeartbeatInterval = toInt(value)
		case "theme":
			if v, ok := value.(string); ok {
				sm.settings.Theme = v
			}
		case "showDebugOutput":
			if v, ok := value.(bool); ok {
				sm.settings.ShowDebugOutput = v
			}
		case "lastProjectPath":
			if v, ok := value.(string); ok {
				sm.settings.LastProjectPath = v
			}
		case "recentProjectsMax":
			sm.settings.RecentProjectsMax = toInt(value)
		case "projectProfiles":
			if v, ok := value.(map[string]interface{}); ok {
				profiles := make(map[string]string)
				for k, pv := range v {
					if s, ok := pv.(string); ok {
						profiles[k] = s
					}
				}
				sm.settings.ProjectProfiles = profiles
			}
		}
	}

	return sm.save()
}

// toInt converts interface{} to int, handling both int and float64 types
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case int64:
		return int(val)
	default:
		return 0
	}
}

// Reset restores all settings to their default values and saves to disk.
func (sm *StateManager) Reset() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.settings = DefaultSettings()
	return sm.save()
}
