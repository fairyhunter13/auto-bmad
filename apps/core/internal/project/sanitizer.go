package project

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	// MaxContextLength is the maximum allowed length for project context descriptions
	MaxContextLength = 500

	// MaxProjectNameLength is the maximum allowed length for project names
	MaxProjectNameLength = 255
)

var (
	// htmlTagRegex matches HTML tags for removal
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

	// scriptTagRegex matches script tags and their content for removal
	scriptTagRegex = regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)

	// ErrContextTooLong indicates the context exceeds maximum length
	ErrContextTooLong = fmt.Errorf("context exceeds maximum length of %d characters", MaxContextLength)
)

// SanitizeContext sanitizes user input for project context descriptions.
// It prevents stored XSS attacks by:
// - Limiting length to prevent DoS
// - Stripping HTML/JavaScript tags
// - Removing control characters
// - Normalizing whitespace
func SanitizeContext(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	// Check length before processing (performance optimization)
	if len(input) > MaxContextLength {
		return "", ErrContextTooLong
	}

	// Remove script tags first (before general HTML tag removal)
	sanitized := scriptTagRegex.ReplaceAllString(input, "")

	// Remove HTML tags
	sanitized = htmlTagRegex.ReplaceAllString(sanitized, "")

	// Remove control characters (except newlines and tabs)
	sanitized = removeControlCharacters(sanitized)

	// Normalize whitespace
	sanitized = normalizeWhitespace(sanitized)

	// Trim leading/trailing whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Final length check after sanitization
	if len(sanitized) > MaxContextLength {
		// Truncate to max length
		sanitized = truncateToLength(sanitized, MaxContextLength)
	}

	return sanitized, nil
}

// removeControlCharacters removes control characters except newlines (\n), tabs (\t), and carriage returns (\r)
func removeControlCharacters(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))

	for _, r := range s {
		// Keep newlines, tabs, carriage returns, and non-control characters
		if r == '\n' || r == '\t' || r == '\r' || !unicode.IsControl(r) {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// normalizeWhitespace normalizes consecutive whitespace to single spaces
func normalizeWhitespace(s string) string {
	// Replace multiple spaces/tabs with single space
	s = regexp.MustCompile(`[ \t]+`).ReplaceAllString(s, " ")

	// Replace multiple newlines with double newline (preserve paragraph breaks)
	s = regexp.MustCompile(`\n{3,}`).ReplaceAllString(s, "\n\n")

	return s
}

// truncateToLength truncates a string to the specified length, respecting UTF-8 character boundaries
func truncateToLength(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	// Ensure we don't cut in the middle of a UTF-8 character
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	return string(runes[:maxLen])
}

// SanitizeProjectName sanitizes project names (from file paths)
// This is less strict since it comes from the file system, but still validates length
func SanitizeProjectName(name string) string {
	// Remove control characters
	sanitized := removeControlCharacters(name)

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Truncate if too long
	if len(sanitized) > MaxProjectNameLength {
		sanitized = truncateToLength(sanitized, MaxProjectNameLength)
	}

	return sanitized
}
