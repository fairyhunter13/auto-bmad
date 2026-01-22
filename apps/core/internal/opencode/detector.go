package opencode

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// DetectionResult represents the result of OpenCode CLI detection.
type DetectionResult struct {
	Found      bool   `json:"found"`
	Version    string `json:"version,omitempty"`
	Path       string `json:"path,omitempty"`
	Compatible bool   `json:"compatible"`
	MinVersion string `json:"minVersion"`
	Error      string `json:"error,omitempty"`
}

// MinimumVersion is the minimum required OpenCode version.
const MinimumVersion = "0.1.0"

// Detect finds and validates the OpenCode CLI installation.
func Detect() (*DetectionResult, error) {
	result := &DetectionResult{
		MinVersion: MinimumVersion,
	}

	// Find opencode in PATH
	path, err := exec.LookPath("opencode")
	if err != nil {
		result.Found = false
		result.Error = "OpenCode CLI not found in PATH"
		return result, nil
	}
	result.Path = path
	result.Found = true

	// Get version
	cmd := exec.Command("opencode", "--version")
	output, err := cmd.Output()
	if err != nil {
		result.Error = "Failed to get OpenCode version"
		return result, nil
	}

	// Parse version (e.g., "opencode v0.1.5" or "opencode version 0.1.5")
	version := parseVersion(string(output))
	result.Version = version

	// Check compatibility
	result.Compatible = isCompatible(version, MinimumVersion)

	return result, nil
}

// parseVersion extracts the semantic version from OpenCode output.
func parseVersion(output string) string {
	// Match patterns like "v0.1.5", "0.1.5", "version 0.1.5"
	re := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1]
	}
	return strings.TrimSpace(output)
}

// isCompatible checks if the version meets the minimum requirement.
func isCompatible(version, minVersion string) bool {
	return compareVersions(version, minVersion) >= 0
}

// compareVersions compares two semantic versions.
// Returns: -1 if v1 < v2, 0 if equal, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	// Parse version components
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Ensure both have 3 parts (major.minor.patch)
	for len(parts1) < 3 {
		parts1 = append(parts1, "0")
	}
	for len(parts2) < 3 {
		parts2 = append(parts2, "0")
	}

	// Compare each component
	for i := 0; i < 3; i++ {
		num1, err1 := strconv.Atoi(parts1[i])
		num2, err2 := strconv.Atoi(parts2[i])

		// If parsing fails, treat as 0
		if err1 != nil {
			num1 = 0
		}
		if err2 != nil {
			num2 = 0
		}

		if num1 > num2 {
			return 1
		}
		if num1 < num2 {
			return -1
		}
	}

	return 0
}
