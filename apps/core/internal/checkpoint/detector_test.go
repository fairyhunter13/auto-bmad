package checkpoint

import (
	"testing"
)

func TestDetectGit(t *testing.T) {
	t.Run("should detect Git when installed", func(t *testing.T) {
		result, err := DetectGit()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
		if !result.Found {
			t.Error("expected Git to be found")
		}
		if result.Path == "" {
			t.Error("expected path to be set when Git is found")
		}
		if result.Version == "" {
			t.Error("expected version to be set when Git is found")
		}
		if result.MinVersion != GitMinimumVersion {
			t.Errorf("expected minVersion=%s, got %s", GitMinimumVersion, result.MinVersion)
		}
	})

	t.Run("should return not found when Git is missing", func(t *testing.T) {
		// This test would require mocking exec.LookPath
		// For now we'll skip it since we can't easily simulate Git not being installed
		t.Skip("requires mock to simulate missing Git")
	})

	t.Run("should set error when Git command fails", func(t *testing.T) {
		// This test would require mocking exec.Command
		t.Skip("requires mock to simulate Git command failure")
	})
}

func TestParseGitVersion(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "standard format",
			output:   "git version 2.39.0",
			expected: "2.39.0",
		},
		{
			name:     "with additional text",
			output:   "git version 2.39.0 (Apple Git-143)",
			expected: "2.39.0",
		},
		{
			name:     "older version",
			output:   "git version 2.0.0",
			expected: "2.0.0",
		},
		{
			name:     "newer version",
			output:   "git version 3.5.1",
			expected: "3.5.1",
		},
		{
			name:     "two digit patch",
			output:   "git version 2.34.12",
			expected: "2.34.12",
		},
		{
			name:     "minimal version",
			output:   "git version 1.9",
			expected: "1.9",
		},
		{
			name:     "empty output",
			output:   "",
			expected: "",
		},
		{
			name:     "invalid format",
			output:   "something else",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseGitVersion(tt.output)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIsGitCompatible(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		minVersion string
		expected   bool
	}{
		{
			name:       "exact minimum version",
			version:    "2.0.0",
			minVersion: "2.0.0",
			expected:   true,
		},
		{
			name:       "newer major version",
			version:    "3.0.0",
			minVersion: "2.0.0",
			expected:   true,
		},
		{
			name:       "newer minor version",
			version:    "2.39.0",
			minVersion: "2.0.0",
			expected:   true,
		},
		{
			name:       "older major version",
			version:    "1.9.5",
			minVersion: "2.0.0",
			expected:   false,
		},
		{
			name:       "older minor version",
			version:    "2.0.0",
			minVersion: "2.5.0",
			expected:   false,
		},
		{
			name:       "older patch version",
			version:    "2.0.0",
			minVersion: "2.0.1",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isGitCompatible(tt.version, tt.minVersion)
			if result != tt.expected {
				t.Errorf("version=%s, minVersion=%s: expected %v, got %v",
					tt.version, tt.minVersion, tt.expected, result)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{
			name:     "equal versions",
			v1:       "2.0.0",
			v2:       "2.0.0",
			expected: 0,
		},
		{
			name:     "v1 greater major",
			v1:       "3.0.0",
			v2:       "2.0.0",
			expected: 1,
		},
		{
			name:     "v1 lesser major",
			v1:       "1.0.0",
			v2:       "2.0.0",
			expected: -1,
		},
		{
			name:     "v1 greater minor",
			v1:       "2.5.0",
			v2:       "2.3.0",
			expected: 1,
		},
		{
			name:     "v1 lesser minor",
			v1:       "2.3.0",
			v2:       "2.5.0",
			expected: -1,
		},
		{
			name:     "v1 greater patch",
			v1:       "2.0.5",
			v2:       "2.0.3",
			expected: 1,
		},
		{
			name:     "v1 lesser patch",
			v1:       "2.0.3",
			v2:       "2.0.5",
			expected: -1,
		},
		{
			name:     "missing patch defaults to 0",
			v1:       "2.0",
			v2:       "2.0.0",
			expected: 0,
		},
		{
			name:     "complex comparison",
			v1:       "2.39.1",
			v2:       "2.0.0",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("compareVersions(%q, %q) = %d, want %d",
					tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}
