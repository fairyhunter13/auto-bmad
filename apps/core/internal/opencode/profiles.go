package opencode

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Profile represents an OpenCode profile from .bash_aliases.
type Profile struct {
	Name      string `json:"name"`
	Alias     string `json:"alias"` // Full alias command
	Available bool   `json:"available"`
	Error     string `json:"error,omitempty"`
	IsDefault bool   `json:"isDefault"`
}

// ProfilesResult represents the result of profile detection.
type ProfilesResult struct {
	Profiles     []Profile `json:"profiles"`
	DefaultFound bool      `json:"defaultFound"`
	Source       string    `json:"source"` // e.g., "~/.bash_aliases"
}

// GetProfiles detects and returns all available OpenCode profiles.
func GetProfiles() (*ProfilesResult, error) {
	result := &ProfilesResult{
		Profiles: []Profile{},
		Source:   "~/.bash_aliases",
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Cannot get home dir, return default
		result.DefaultFound = true
		result.Profiles = append(result.Profiles, Profile{
			Name:      "default",
			Available: true,
			IsDefault: true,
		})
		return result, nil
	}

	aliasFile := filepath.Join(homeDir, ".bash_aliases")

	file, err := os.Open(aliasFile)
	if err != nil {
		// No aliases file - use default
		result.DefaultFound = true
		result.Profiles = append(result.Profiles, Profile{
			Name:      "default",
			Available: true,
			IsDefault: true,
		})
		return result, nil
	}
	defer file.Close()

	// Pattern: alias opencode-{name}='...' or alias opencode-{name}="..."
	aliasPattern := regexp.MustCompile(`^alias\s+opencode-(\w+)=['"](.+?)['"]`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		matches := aliasPattern.FindStringSubmatch(line)
		if len(matches) == 3 {
			profileName := matches[1]
			aliasCmd := matches[2]

			profile := Profile{
				Name:      profileName,
				Alias:     aliasCmd,
				IsDefault: false,
			}

			// Validate profile availability
			profile.Available, profile.Error = validateProfile(aliasCmd)

			result.Profiles = append(result.Profiles, profile)
		}
	}

	// If no profiles found, add default
	if len(result.Profiles) == 0 {
		result.DefaultFound = true
		result.Profiles = append(result.Profiles, Profile{
			Name:      "default",
			Available: true,
			IsDefault: true,
		})
	}

	return result, nil
}

// validateProfile checks if a profile alias command is valid.
// It attempts to parse the command and check if it references opencode.
// Returns (available, errorMessage)
func validateProfile(aliasCmd string) (bool, string) {
	// Basic validation: check if the command contains "opencode"
	if !strings.Contains(aliasCmd, "opencode") {
		return false, "alias does not reference opencode command"
	}

	// For now, mark all opencode aliases as available
	// Actual runtime validation would happen when the profile is used
	// This prevents expensive shell executions during profile detection
	return true, ""
}
