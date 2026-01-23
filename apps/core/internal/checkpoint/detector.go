package checkpoint

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// GitDetectionResult represents the result of Git detection.
type GitDetectionResult struct {
	Found      bool   `json:"found"`
	Version    string `json:"version,omitempty"`
	Path       string `json:"path,omitempty"`
	Compatible bool   `json:"compatible"`
	MinVersion string `json:"minVersion"`
	Error      string `json:"error,omitempty"`
}

// GitMinimumVersion is the minimum required Git version.
const GitMinimumVersion = "2.0.0"

// DetectGit finds and validates the Git installation.
func DetectGit() (*GitDetectionResult, error) {
	result := &GitDetectionResult{
		MinVersion: GitMinimumVersion,
	}

	// Find git in PATH
	path, err := exec.LookPath("git")
	if err != nil {
		result.Found = false
		result.Error = "Git not found in PATH"
		return result, nil
	}
	result.Path = path
	result.Found = true

	// Get version
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		result.Error = "Failed to get Git version"
		return result, nil
	}

	// Parse version (e.g., "git version 2.39.0")
	version := parseGitVersion(string(output))
	result.Version = version

	// Check compatibility
	result.Compatible = isGitCompatible(version, GitMinimumVersion)

	return result, nil
}

// parseGitVersion extracts the semantic version from Git output.
func parseGitVersion(output string) string {
	// Match pattern "git version X.Y.Z" or "git version X.Y"
	re := regexp.MustCompile(`git version (\d+\.\d+\.?\d*)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// isGitCompatible checks if the version meets the minimum requirement.
func isGitCompatible(version, minVersion string) bool {
	return compareVersions(version, minVersion) >= 0
}

// GitRepoStatus represents the status of a Git repository.
type GitRepoStatus struct {
	IsGitRepo  bool   `json:"isGitRepo"`
	Branch     string `json:"branch,omitempty"`
	HasChanges bool   `json:"hasChanges"`
	Error      string `json:"error,omitempty"`
}

// GetRepoStatus returns the Git status for a given directory path.
func GetRepoStatus(path string) *GitRepoStatus {
	result := &GitRepoStatus{}

	// Check if .git directory exists
	cmd := exec.Command("git", "-C", path, "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		// Not a git repository
		result.IsGitRepo = false
		return result
	}
	result.IsGitRepo = true

	// Get current branch
	branchCmd := exec.Command("git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, err := branchCmd.Output()
	if err != nil {
		result.Error = "Failed to get branch name"
	} else {
		result.Branch = strings.TrimSpace(string(branchOutput))
	}

	// Check for uncommitted changes
	statusCmd := exec.Command("git", "-C", path, "status", "--porcelain")
	statusOutput, err := statusCmd.Output()
	if err != nil {
		result.Error = "Failed to get status"
	} else {
		result.HasChanges = len(strings.TrimSpace(string(statusOutput))) > 0
	}

	return result
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
