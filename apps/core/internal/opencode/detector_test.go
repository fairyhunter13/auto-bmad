package opencode

import (
	"testing"
)

func TestDetect_OpencodeFound(t *testing.T) {
	// This test assumes opencode is installed on the system
	// In a real scenario, we would mock exec.LookPath and exec.Command
	result, err := Detect()
	if err != nil {
		t.Fatalf("Detect() returned unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Detect() returned nil result")
	}

	// The result should indicate whether opencode was found
	// We can't assert it's always true because it depends on the system
	if result.Found {
		if result.Path == "" {
			t.Error("Expected Path to be set when Found is true")
		}
		if result.Version == "" {
			t.Error("Expected Version to be set when Found is true")
		}
	} else {
		if result.Error == "" {
			t.Error("Expected Error to be set when Found is false")
		}
	}

	// MinVersion should always be set
	if result.MinVersion == "" {
		t.Error("Expected MinVersion to be set")
	}
}

func TestDetect_NotFound(t *testing.T) {
	// We can't easily test this without mocking
	// This is a placeholder for when we add dependency injection
	t.Skip("TODO: Implement mock-based test for opencode not found scenario")
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "version with v prefix",
			output:   "opencode v0.1.5",
			expected: "0.1.5",
		},
		{
			name:     "version without v prefix",
			output:   "opencode 0.1.5",
			expected: "0.1.5",
		},
		{
			name:     "version with word version",
			output:   "opencode version 0.1.5",
			expected: "0.1.5",
		},
		{
			name:     "multiline output",
			output:   "opencode v0.2.0\nBuild: 12345",
			expected: "0.2.0",
		},
		{
			name:     "no version found",
			output:   "opencode help",
			expected: "opencode help",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseVersion(tt.output)
			if result != tt.expected {
				t.Errorf("parseVersion(%q) = %q, expected %q", tt.output, result, tt.expected)
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
			v1:       "0.1.0",
			v2:       "0.1.0",
			expected: 0,
		},
		{
			name:     "v1 greater than v2",
			v1:       "0.2.0",
			v2:       "0.1.0",
			expected: 1,
		},
		{
			name:     "v1 less than v2",
			v1:       "0.1.0",
			v2:       "0.2.0",
			expected: -1,
		},
		{
			name:     "patch version difference",
			v1:       "0.1.5",
			v2:       "0.1.0",
			expected: 1,
		},
		{
			name:     "major version difference",
			v1:       "1.0.0",
			v2:       "0.9.9",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("compareVersions(%q, %q) = %d, expected %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestIsCompatible(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		minVersion string
		expected   bool
	}{
		{
			name:       "exact minimum version",
			version:    "0.1.0",
			minVersion: "0.1.0",
			expected:   true,
		},
		{
			name:       "above minimum version",
			version:    "0.2.0",
			minVersion: "0.1.0",
			expected:   true,
		},
		{
			name:       "below minimum version",
			version:    "0.0.9",
			minVersion: "0.1.0",
			expected:   false,
		},
		{
			name:       "patch above minimum",
			version:    "0.1.5",
			minVersion: "0.1.0",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCompatible(tt.version, tt.minVersion)
			if result != tt.expected {
				t.Errorf("isCompatible(%q, %q) = %v, expected %v", tt.version, tt.minVersion, result, tt.expected)
			}
		})
	}
}
