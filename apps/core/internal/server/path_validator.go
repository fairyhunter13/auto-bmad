package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// ErrInvalidPath indicates the path failed validation
	ErrInvalidPath = errors.New("invalid path")

	// ErrPathNotAbsolute indicates the path is not absolute
	ErrPathNotAbsolute = errors.New("path must be absolute")

	// ErrPathNotDirectory indicates the path is not a directory
	ErrPathNotDirectory = errors.New("path must be a directory")

	// ErrForbiddenPath indicates access to this path is forbidden
	ErrForbiddenPath = errors.New("access to this path is forbidden")
)

// PathValidator validates and sanitizes file system paths to prevent path traversal attacks
type PathValidator struct {
	// AllowNonExistent allows validating paths that don't exist yet
	AllowNonExistent bool

	// RequireDirectory requires the path to be a directory
	RequireDirectory bool
}

// NewPathValidator creates a new path validator with default settings
func NewPathValidator() *PathValidator {
	return &PathValidator{
		AllowNonExistent: false,
		RequireDirectory: true,
	}
}

// Validate validates and sanitizes a file system path
func (v *PathValidator) Validate(inputPath string) (string, error) {
	if inputPath == "" {
		return "", ErrInvalidPath
	}

	// Clean the path (removes .., ., etc.)
	cleanPath := filepath.Clean(inputPath)

	// Ensure path is absolute
	if !filepath.IsAbs(cleanPath) {
		return "", ErrPathNotAbsolute
	}

	// Resolve symlinks to get the real path
	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		// If path doesn't exist and we allow non-existent paths, use cleaned path
		if v.AllowNonExistent && os.IsNotExist(err) {
			realPath = cleanPath
		} else {
			return "", fmt.Errorf("failed to resolve path: %w", err)
		}
	}

	// Check if path exists (if required)
	if !v.AllowNonExistent {
		info, err := os.Stat(realPath)
		if err != nil {
			if os.IsNotExist(err) {
				return "", fmt.Errorf("path does not exist: %s", realPath)
			}
			return "", fmt.Errorf("failed to stat path: %w", err)
		}

		// Check if it's a directory (if required)
		if v.RequireDirectory && !info.IsDir() {
			return "", ErrPathNotDirectory
		}
	}

	// Check forbidden paths (system directories)
	if err := v.checkForbiddenPaths(realPath); err != nil {
		return "", err
	}

	return realPath, nil
}

// checkForbiddenPaths blocks access to sensitive system directories
func (v *PathValidator) checkForbiddenPaths(path string) error {
	// Normalize path separators for comparison
	normalizedPath := filepath.ToSlash(strings.ToLower(path))

	// Define forbidden path prefixes by OS
	var forbiddenPrefixes []string

	switch runtime.GOOS {
	case "windows":
		forbiddenPrefixes = []string{
			"c:/windows/system32",
			"c:/windows/syswow64",
			"c:/program files/windowsapps",
		}
	case "darwin": // macOS
		forbiddenPrefixes = []string{
			"/system",
			"/private/var/db",
			"/private/var/root",
			"/library/application support/apple",
		}
	case "linux":
		forbiddenPrefixes = []string{
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

	// Check if path starts with any forbidden prefix
	for _, prefix := range forbiddenPrefixes {
		if strings.HasPrefix(normalizedPath, prefix) {
			return fmt.Errorf("%w: access to system directory %s is not allowed", ErrForbiddenPath, prefix)
		}
	}

	return nil
}

// ValidateProjectPath is a convenience function for validating project paths
// It ensures the path is absolute, exists, is a directory, and resolves symlinks
func ValidateProjectPath(path string) (string, error) {
	validator := &PathValidator{
		AllowNonExistent: false,
		RequireDirectory: true,
	}
	return validator.Validate(path)
}
