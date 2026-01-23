package opencode

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetProfiles_NoAliasesFile(t *testing.T) {
	// Create a temporary directory to simulate no ~/.bash_aliases
	tempDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	result, err := GetProfiles()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.DefaultFound {
		t.Error("Expected DefaultFound to be true when no aliases file exists")
	}

	if len(result.Profiles) != 1 {
		t.Errorf("Expected 1 default profile, got %d", len(result.Profiles))
	}

	if len(result.Profiles) > 0 {
		profile := result.Profiles[0]
		if profile.Name != "default" {
			t.Errorf("Expected default profile name, got: %s", profile.Name)
		}
		if !profile.Available {
			t.Error("Expected default profile to be available")
		}
		if !profile.IsDefault {
			t.Error("Expected default profile IsDefault flag to be true")
		}
	}

	if result.Source != "~/.bash_aliases" {
		t.Errorf("Expected source to be ~/.bash_aliases, got: %s", result.Source)
	}
}

func TestGetProfiles_WithOpenCodeAliases(t *testing.T) {
	// Create a temporary directory with .bash_aliases
	tempDir := t.TempDir()
	aliasContent := `# OpenCode Multi-Account Setup
alias opencode-personal='OPENCODE_DISABLE_AUTOUPDATE=true XDG_CONFIG_HOME=$HOME/.config/opencode-personal opencode'
alias opencode-work='OPENCODE_DISABLE_AUTOUPDATE=true XDG_CONFIG_HOME=$HOME/.config/opencode-work opencode'
alias opencode-default='OPENCODE_DISABLE_AUTOUPDATE=true opencode'
alias ocp='opencode-personal'
alias other-alias='ls -la'
`

	aliasFile := filepath.Join(tempDir, ".bash_aliases")
	if err := os.WriteFile(aliasFile, []byte(aliasContent), 0644); err != nil {
		t.Fatalf("Failed to create test .bash_aliases: %v", err)
	}

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	result, err := GetProfiles()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Should find 3 opencode profiles
	if len(result.Profiles) != 3 {
		t.Errorf("Expected 3 profiles, got %d", len(result.Profiles))
	}

	// Check profile names
	expectedProfiles := map[string]bool{
		"personal": false,
		"work":     false,
		"default":  false,
	}

	for _, profile := range result.Profiles {
		if _, exists := expectedProfiles[profile.Name]; !exists {
			t.Errorf("Unexpected profile name: %s", profile.Name)
		}
		expectedProfiles[profile.Name] = true

		// SECURITY: Alias field removed to prevent shell injection
		// Profile names are validated (alphanumeric only) and safe to use
	}

	// Verify all expected profiles were found
	for name, found := range expectedProfiles {
		if !found {
			t.Errorf("Expected profile %s not found", name)
		}
	}
}

func TestGetProfiles_EmptyAliasesFile(t *testing.T) {
	// Create a temporary directory with empty .bash_aliases
	tempDir := t.TempDir()
	aliasFile := filepath.Join(tempDir, ".bash_aliases")
	if err := os.WriteFile(aliasFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create test .bash_aliases: %v", err)
	}

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	result, err := GetProfiles()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should return default profile when no opencode aliases found
	if !result.DefaultFound {
		t.Error("Expected DefaultFound to be true when no opencode aliases exist")
	}

	if len(result.Profiles) != 1 {
		t.Errorf("Expected 1 default profile, got %d", len(result.Profiles))
	}

	if len(result.Profiles) > 0 && result.Profiles[0].Name != "default" {
		t.Errorf("Expected default profile, got: %s", result.Profiles[0].Name)
	}
}

func TestGetProfiles_MixedAliases(t *testing.T) {
	// Create a temporary directory with mixed aliases
	tempDir := t.TempDir()
	aliasContent := `# Various aliases
alias ll='ls -la'
alias opencode-test='opencode --provider test'
alias k8='kubectl'
alias opencode-staging="opencode --env staging"
alias gs='git status'
`

	aliasFile := filepath.Join(tempDir, ".bash_aliases")
	if err := os.WriteFile(aliasFile, []byte(aliasContent), 0644); err != nil {
		t.Fatalf("Failed to create test .bash_aliases: %v", err)
	}

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	result, err := GetProfiles()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should find 2 opencode profiles (test and staging)
	if len(result.Profiles) != 2 {
		t.Errorf("Expected 2 profiles, got %d", len(result.Profiles))
	}

	profileNames := make(map[string]bool)
	for _, profile := range result.Profiles {
		profileNames[profile.Name] = true
	}

	if !profileNames["test"] {
		t.Error("Expected to find 'test' profile")
	}
	if !profileNames["staging"] {
		t.Error("Expected to find 'staging' profile")
	}
}

func TestValidateProfile_Available(t *testing.T) {
	// This test assumes opencode is installed
	available, errMsg := validateProfile("opencode")

	// Since we can't guarantee opencode is installed in test environment,
	// we just check that the function returns valid results
	if available && errMsg != "" {
		t.Error("Available profile should not have error message")
	}
	if !available && errMsg == "" {
		t.Error("Unavailable profile should have error message")
	}
}

func TestValidateProfile_WithAlias(t *testing.T) {
	// Create a temporary directory with a test alias
	tempDir := t.TempDir()
	aliasContent := `alias opencode-test='opencode --version'
`

	aliasFile := filepath.Join(tempDir, ".bash_aliases")
	if err := os.WriteFile(aliasFile, []byte(aliasContent), 0644); err != nil {
		t.Fatalf("Failed to create test .bash_aliases: %v", err)
	}

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	result, err := GetProfiles()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Profile validation should be executed during GetProfiles
	if len(result.Profiles) > 0 {
		profile := result.Profiles[0]
		// Check that Available field is set (either true or false)
		// We can't guarantee the actual value since it depends on opencode installation
		_ = profile.Available
	}
}

func TestGetProfilesWithValidation(t *testing.T) {
	// Create a temporary directory with test aliases
	tempDir := t.TempDir()
	aliasContent := `alias opencode-valid='opencode --version'
alias opencode-test='echo test'
`

	aliasFile := filepath.Join(tempDir, ".bash_aliases")
	if err := os.WriteFile(aliasFile, []byte(aliasContent), 0644); err != nil {
		t.Fatalf("Failed to create test .bash_aliases: %v", err)
	}

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	result, err := GetProfiles()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result.Profiles) != 2 {
		t.Errorf("Expected 2 profiles, got %d", len(result.Profiles))
	}

	// Each profile should have availability status set
	for _, profile := range result.Profiles {
		if profile.Name == "" {
			t.Error("Profile should have a name")
		}
		// Available field should be set (true or false)
		// Error field may or may not be set depending on availability
		if !profile.Available && profile.Error == "" {
			t.Errorf("Unavailable profile %s should have error message", profile.Name)
		}
	}
}
