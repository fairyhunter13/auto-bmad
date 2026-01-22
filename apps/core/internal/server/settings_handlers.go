package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/state"
)

// settingsManager is the global settings manager instance
var settingsManager *state.StateManager

// RegisterSettingsHandlers registers all settings-related JSON-RPC handlers.
// These handlers manage user settings persistence and retrieval.
// Settings are stored globally (not per-project) to remember user preferences across sessions.
func RegisterSettingsHandlers(s *Server) error {
	// Use user's home directory for global settings
	// This allows settings to persist across different projects
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	// Store settings in ~/.autobmad directory
	settingsPath := filepath.Join(homeDir, ".autobmad")

	// Create state manager
	sm, err := state.NewStateManager(settingsPath)
	if err != nil {
		return fmt.Errorf("creating state manager: %w", err)
	}

	// Store globally for access by handlers
	settingsManager = sm

	// Register handlers
	s.RegisterHandler("settings.get", handleSettingsGet(sm))
	s.RegisterHandler("settings.set", handleSettingsSet(sm))
	s.RegisterHandler("settings.reset", handleSettingsReset(sm))

	return nil
}

// handleSettingsGet returns the current settings.
// Method: settings.get
// Params: none
// Result: Settings object
func handleSettingsGet(sm *state.StateManager) Handler {
	return func(params json.RawMessage) (interface{}, error) {
		return sm.Get(), nil
	}
}

// handleSettingsSet updates settings with provided values.
// Method: settings.set
// Params: map of setting keys to values
// Result: Updated Settings object
func handleSettingsSet(sm *state.StateManager) Handler {
	return func(params json.RawMessage) (interface{}, error) {
		// Parse update map
		var updates map[string]interface{}
		if err := json.Unmarshal(params, &updates); err != nil {
			return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid params", err.Error())
		}

		// Apply updates
		if err := sm.Set(updates); err != nil {
			return nil, NewErrorWithData(ErrCodeInternalError, "Failed to save settings", err.Error())
		}

		// Return updated settings
		return sm.Get(), nil
	}
}

// handleSettingsReset resets all settings to defaults.
// Method: settings.reset
// Params: none
// Result: Default Settings object
func handleSettingsReset(sm *state.StateManager) Handler {
	return func(params json.RawMessage) (interface{}, error) {
		if err := sm.Reset(); err != nil {
			return nil, NewErrorWithData(ErrCodeInternalError, "Failed to reset settings", err.Error())
		}

		// Return reset settings
		return sm.Get(), nil
	}
}
