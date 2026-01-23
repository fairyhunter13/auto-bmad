package project

import (
	"strings"
	"testing"
)

func TestSanitizeContext(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "plain text",
			input:    "This is a simple project description",
			expected: "This is a simple project description",
		},
		{
			name:     "text with newlines",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "HTML tags removed",
			input:    "This is <b>bold</b> and <i>italic</i> text",
			expected: "This is bold and italic text",
		},
		{
			name:     "script tags removed",
			input:    "Safe text <script>alert('XSS')</script> more text",
			expected: "Safe text more text",
		},
		{
			name:     "script tags case insensitive",
			input:    "Safe text <SCRIPT>alert('XSS')</SCRIPT> more text",
			expected: "Safe text more text",
		},
		{
			name:     "multiple script tags",
			input:    "<script>bad()</script>Good<script>bad2()</script>",
			expected: "Good",
		},
		{
			name:     "complex HTML structure",
			input:    "<div><p>Paragraph</p><a href='evil'>Link</a></div>",
			expected: "ParagraphLink",
		},
		{
			name:     "HTML entities preserved (not decoded)",
			input:    "5 &lt; 10 and 10 &gt; 5",
			expected: "5 &lt; 10 and 10 &gt; 5",
		},
		{
			name:     "control characters removed",
			input:    "Text\x00with\x01control\x02chars",
			expected: "Textwithcontrolchars",
		},
		{
			name:     "multiple spaces normalized",
			input:    "Too    many     spaces",
			expected: "Too many spaces",
		},
		{
			name:     "multiple newlines normalized",
			input:    "Line 1\n\n\n\n\nLine 2",
			expected: "Line 1\n\nLine 2",
		},
		{
			name:     "leading and trailing whitespace removed",
			input:    "   trimmed   ",
			expected: "trimmed",
		},
		{
			name:     "tabs normalized",
			input:    "Text\t\t\twith\t\ttabs",
			expected: "Text with tabs",
		},
		{
			name:     "unicode characters preserved",
			input:    "Unicode: 你好世界 🚀 émojis",
			expected: "Unicode: 你好世界 🚀 émojis",
		},
		{
			name:     "mixed attacks",
			input:    "<script>alert('xss')</script><b>Bold</b> with  spaces\x00\x01",
			expected: "Bold with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizeContext(tt.input)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSanitizeContext_MaxLength(t *testing.T) {
	// Test exact max length (should pass)
	t.Run("exact max length", func(t *testing.T) {
		input := strings.Repeat("a", MaxContextLength)
		result, err := SanitizeContext(input)
		if err != nil {
			t.Errorf("Unexpected error for max length: %v", err)
		}
		if len(result) != MaxContextLength {
			t.Errorf("Expected length %d, got %d", MaxContextLength, len(result))
		}
	})

	// Test over max length (should error)
	t.Run("over max length", func(t *testing.T) {
		input := strings.Repeat("a", MaxContextLength+1)
		_, err := SanitizeContext(input)
		if err != ErrContextTooLong {
			t.Errorf("Expected ErrContextTooLong, got: %v", err)
		}
	})

	// Test that oversized input is rejected even if HTML would make it smaller
	t.Run("oversized HTML input rejected", func(t *testing.T) {
		// Create input with HTML that, when removed, would be under max length
		// But raw input exceeds max length - should be rejected for DoS protection
		input := strings.Repeat("<b>a</b>", MaxContextLength/2)
		_, err := SanitizeContext(input)
		if err != ErrContextTooLong {
			t.Errorf("Expected ErrContextTooLong for oversized HTML input, got: %v", err)
		}
	})
}

func TestSanitizeContext_XSSVectors(t *testing.T) {
	xssVectors := []struct {
		name  string
		input string
	}{
		{
			name:  "basic script tag",
			input: "<script>alert('XSS')</script>",
		},
		{
			name:  "img onerror",
			input: "<img src=x onerror=alert('XSS')>",
		},
		{
			name:  "javascript protocol",
			input: "<a href='javascript:alert(1)'>Click</a>",
		},
		{
			name:  "svg onload",
			input: "<svg onload=alert('XSS')>",
		},
		{
			name:  "iframe",
			input: "<iframe src='javascript:alert(1)'>",
		},
		{
			name:  "nested tags",
			input: "<div><script>alert(1)</script></div>",
		},
		{
			name:  "event handlers",
			input: "<button onclick='alert(1)'>Click</button>",
		},
	}

	for _, tt := range xssVectors {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizeContext(tt.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Result should not contain script tags or event handlers
			if strings.Contains(strings.ToLower(result), "<script") {
				t.Errorf("Result contains script tag: %s", result)
			}
			if strings.Contains(strings.ToLower(result), "javascript:") {
				t.Errorf("Result contains javascript protocol: %s", result)
			}
			if strings.Contains(result, "<") || strings.Contains(result, ">") {
				t.Errorf("Result contains HTML tags: %s", result)
			}
		})
	}
}

func TestSanitizeContext_UTF8(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "emoji",
			input:    "Project with emoji 🚀 🎉",
			expected: "Project with emoji 🚀 🎉",
		},
		{
			name:     "chinese characters",
			input:    "项目描述",
			expected: "项目描述",
		},
		{
			name:     "arabic characters",
			input:    "وصف المشروع",
			expected: "وصف المشروع",
		},
		{
			name:     "mixed unicode",
			input:    "Mixed: English 中文 العربية 🌍",
			expected: "Mixed: English 中文 العربية 🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizeContext(tt.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTruncateToLength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "no truncation needed",
			input:    "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "exact length",
			input:    "exactly10!",
			maxLen:   10,
			expected: "exactly10!",
		},
		{
			name:     "truncate ASCII",
			input:    "this is too long",
			maxLen:   7,
			expected: "this is",
		},
		{
			name:     "truncate unicode",
			input:    "你好世界朋友们",
			maxLen:   4,
			expected: "你好世界",
		},
		{
			name:     "truncate emoji",
			input:    "🚀🎉🌍🔥💡",
			maxLen:   3,
			expected: "🚀🎉🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateToLength(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
			// Verify rune count doesn't exceed max
			if len([]rune(result)) > tt.maxLen {
				t.Errorf("Result rune count %d exceeds max %d", len([]rune(result)), tt.maxLen)
			}
		})
	}
}

func TestSanitizeProjectName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name",
			input:    "my-project",
			expected: "my-project",
		},
		{
			name:     "name with spaces",
			input:    "My Project Name",
			expected: "My Project Name",
		},
		{
			name:     "control characters removed",
			input:    "project\x00name\x01",
			expected: "projectname",
		},
		{
			name:     "whitespace trimmed",
			input:    "  project  ",
			expected: "project",
		},
		{
			name:     "unicode preserved",
			input:    "项目-🚀",
			expected: "项目-🚀",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSanitizeProjectName_MaxLength(t *testing.T) {
	// Test truncation of very long project names
	longName := strings.Repeat("a", MaxProjectNameLength+50)
	result := SanitizeProjectName(longName)

	if len(result) > MaxProjectNameLength {
		t.Errorf("Expected length <= %d, got %d", MaxProjectNameLength, len(result))
	}

	if len(result) != MaxProjectNameLength {
		t.Errorf("Expected exact truncation to %d, got %d", MaxProjectNameLength, len(result))
	}
}
