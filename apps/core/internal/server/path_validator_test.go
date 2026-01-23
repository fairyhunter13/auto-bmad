package server

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestPathValidator_Validate(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pathvalidator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test subdirectory
	testDir := filepath.Join(tempDir, "testproject")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tempDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name          string
		inputPath     string
		allowNonExist bool
		requireDir    bool
		wantErr       bool
		errContains   string
	}{
		{
			name:        "empty path",
			inputPath:   "",
			requireDir:  true,
			wantErr:     true,
			errContains: "invalid path",
		},
		{
			name:        "relative path",
			inputPath:   "./relative/path",
			requireDir:  true,
			wantErr:     true,
			errContains: "must be absolute",
		},
		{
			name:       "valid absolute directory",
			inputPath:  testDir,
			requireDir: true,
			wantErr:    false,
		},
		{
			name:        "non-existent path without allow flag",
			inputPath:   filepath.Join(tempDir, "nonexistent"),
			requireDir:  true,
			wantErr:     true,
			errContains: "failed to resolve path",
		},
		{
			name:          "non-existent path with allow flag",
			inputPath:     filepath.Join(tempDir, "nonexistent"),
			allowNonExist: true,
			requireDir:    false,
			wantErr:       false,
		},
		{
			name:        "file when directory required",
			inputPath:   testFile,
			requireDir:  true,
			wantErr:     true,
			errContains: "must be a directory",
		},
		{
			name:       "file when directory not required",
			inputPath:  testFile,
			requireDir: false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &PathValidator{
				AllowNonExistent: tt.allowNonExist,
				RequireDirectory: tt.requireDir,
			}

			result, err := v.Validate(tt.inputPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error containing %q but got: %v", tt.errContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if result == "" {
					t.Errorf("Expected non-empty result path")
				}
				if !filepath.IsAbs(result) {
					t.Errorf("Result path is not absolute: %s", result)
				}
			}
		})
	}
}

func TestPathValidator_ForbiddenPaths(t *testing.T) {
	v := NewPathValidator()
	v.AllowNonExistent = true // We're testing path patterns, not existence

	var forbiddenPaths []string

	switch runtime.GOOS {
	case "windows":
		forbiddenPaths = []string{
			"C:\\Windows\\System32",
			"C:\\Windows\\SysWOW64",
			"C:\\Program Files\\WindowsApps",
		}
	case "darwin":
		forbiddenPaths = []string{
			"/System",
			"/private/var/db",
			"/private/var/root",
		}
	case "linux":
		forbiddenPaths = []string{
			"/sys",
			"/proc",
			"/dev",
			"/boot",
			"/root",
			"/etc/shadow",
			"/etc/passwd",
			"/var/log",
		}
	}

	for _, path := range forbiddenPaths {
		t.Run("forbidden_"+path, func(t *testing.T) {
			_, err := v.Validate(path)
			if err == nil {
				t.Errorf("Expected error for forbidden path %s but got none", path)
				return
			}
			if !strings.Contains(err.Error(), "forbidden") && !strings.Contains(err.Error(), "not allowed") {
				t.Errorf("Expected 'forbidden' error for path %s, got: %v", path, err)
			}
		})
	}
}

func TestPathValidator_PathTraversalAttacks(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "pathtraversal-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a safe subdirectory
	safeDir := filepath.Join(tempDir, "safe")
	if err := os.MkdirAll(safeDir, 0755); err != nil {
		t.Fatalf("Failed to create safe dir: %v", err)
	}

	v := NewPathValidator()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "path with .. components is cleaned",
			path:    filepath.Join(safeDir, "..", "safe"),
			wantErr: false, // Should be cleaned to safeDir
		},
		{
			name:    "path with . components is cleaned",
			path:    filepath.Join(safeDir, ".", ".", "."),
			wantErr: false, // Should be cleaned to safeDir
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := v.Validate(tt.path)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				// Result should be cleaned and absolute
				if !filepath.IsAbs(result) {
					t.Errorf("Result is not absolute: %s", result)
				}
				// Should not contain .. or .
				if strings.Contains(result, "..") || strings.Contains(filepath.ToSlash(result), "/.") {
					t.Errorf("Result contains navigation components: %s", result)
				}
			}
		})
	}
}

func TestPathValidator_SymlinkResolution(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "symlink-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a real directory
	realDir := filepath.Join(tempDir, "real")
	if err := os.MkdirAll(realDir, 0755); err != nil {
		t.Fatalf("Failed to create real dir: %v", err)
	}

	// Create a symlink to the real directory
	symlinkPath := filepath.Join(tempDir, "symlink")
	if err := os.Symlink(realDir, symlinkPath); err != nil {
		// Skip test on systems that don't support symlinks
		t.Skipf("Skipping symlink test: %v", err)
	}

	v := NewPathValidator()

	// Validate the symlink
	result, err := v.Validate(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to validate symlink: %v", err)
	}

	// Result should be the real path, not the symlink
	if result == symlinkPath {
		t.Errorf("Expected symlink to be resolved, but got symlink path")
	}

	// Result should match the real directory (after cleaning)
	cleanRealDir, _ := filepath.EvalSymlinks(realDir)
	cleanResult, _ := filepath.EvalSymlinks(result)
	if cleanResult != cleanRealDir {
		t.Errorf("Expected resolved path %s, got %s", cleanRealDir, cleanResult)
	}
}

func TestValidateProjectPath(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "project-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test valid project path
	result, err := ValidateProjectPath(tempDir)
	if err != nil {
		t.Errorf("Unexpected error for valid project path: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result")
	}

	// Test non-existent path
	_, err = ValidateProjectPath(filepath.Join(tempDir, "nonexistent"))
	if err == nil {
		t.Errorf("Expected error for non-existent path")
	}

	// Test empty path
	_, err = ValidateProjectPath("")
	if err == nil {
		t.Errorf("Expected error for empty path")
	}
}
