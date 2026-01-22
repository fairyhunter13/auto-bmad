// Package state provides settings persistence and state management.
package state

// Settings represents user-configurable settings for Auto-BMAD.
// All settings are persisted to _bmad-output/.autobmad/config.json.
type Settings struct {
	// Retry settings
	MaxRetries int `json:"maxRetries"` // Default: 3
	RetryDelay int `json:"retryDelay"` // Default: 5000 (ms)

	// Notification settings
	DesktopNotifications bool `json:"desktopNotifications"` // Default: true
	SoundEnabled         bool `json:"soundEnabled"`         // Default: false

	// Timeout settings
	StepTimeoutDefault int `json:"stepTimeoutDefault"` // Default: 300000 (5 min)
	HeartbeatInterval  int `json:"heartbeatInterval"`  // Default: 60000 (60s)

	// UI preferences
	Theme           string `json:"theme"`           // Default: "system"
	ShowDebugOutput bool   `json:"showDebugOutput"` // Default: false

	// Project memory
	LastProjectPath   string            `json:"lastProjectPath,omitempty"`
	ProjectProfiles   map[string]string `json:"projectProfiles"`   // path -> profile name
	RecentProjectsMax int               `json:"recentProjectsMax"` // Default: 10
}

// DefaultSettings returns a new Settings instance with sensible defaults.
func DefaultSettings() *Settings {
	return &Settings{
		MaxRetries:           3,
		RetryDelay:           5000,
		DesktopNotifications: true,
		SoundEnabled:         false,
		StepTimeoutDefault:   300000,
		HeartbeatInterval:    60000,
		Theme:                "system",
		ShowDebugOutput:      false,
		ProjectProfiles:      make(map[string]string),
		RecentProjectsMax:    10,
	}
}
