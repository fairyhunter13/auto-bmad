package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScan_NotBMADProject(t *testing.T) {
	// Create temp directory without _bmad folder
	tmpDir := t.TempDir()

	result, err := Scan(tmpDir)
	require.NoError(t, err)

	assert.False(t, result.IsBMAD)
	assert.False(t, result.HasBmadFolder)
	assert.Equal(t, TypeNotBMAD, result.ProjectType)
	assert.Equal(t, tmpDir, result.Path)
	assert.Equal(t, MinBmadVersion, result.MinBmadVersion)
	assert.Contains(t, result.Error, "Not a BMAD project")
	assert.Contains(t, result.Error, "bmad-init")
}

func TestScan_GreenfieldProject(t *testing.T) {
	// Create temp directory with _bmad but no _bmad-output
	tmpDir := t.TempDir()
	bmadPath := filepath.Join(tmpDir, "_bmad")
	require.NoError(t, os.MkdirAll(bmadPath, 0755))

	// Create manifest with version
	manifestPath := filepath.Join(bmadPath, "_config")
	require.NoError(t, os.MkdirAll(manifestPath, 0755))
	manifestFile := filepath.Join(manifestPath, "manifest.yaml")
	manifestContent := "version: 6.1.0\n"
	require.NoError(t, os.WriteFile(manifestFile, []byte(manifestContent), 0644))

	result, err := Scan(tmpDir)
	require.NoError(t, err)

	assert.True(t, result.IsBMAD)
	assert.True(t, result.HasBmadFolder)
	assert.False(t, result.HasOutputFolder)
	assert.Equal(t, TypeGreenfield, result.ProjectType)
	assert.Equal(t, "6.1.0", result.BmadVersion)
	assert.True(t, result.BmadCompatible)
	assert.Empty(t, result.Error)
	assert.Empty(t, result.ExistingArtifacts)
}

func TestScan_BrownfieldProject(t *testing.T) {
	// Create temp directory with _bmad and _bmad-output with artifacts
	tmpDir := t.TempDir()
	bmadPath := filepath.Join(tmpDir, "_bmad")
	require.NoError(t, os.MkdirAll(bmadPath, 0755))

	// Create manifest
	manifestPath := filepath.Join(bmadPath, "_config")
	require.NoError(t, os.MkdirAll(manifestPath, 0755))
	manifestFile := filepath.Join(manifestPath, "manifest.yaml")
	manifestContent := "version: 6.0.0\n"
	require.NoError(t, os.WriteFile(manifestFile, []byte(manifestContent), 0644))

	// Create _bmad-output with artifacts
	outputPath := filepath.Join(tmpDir, "_bmad-output", "planning-artifacts")
	require.NoError(t, os.MkdirAll(outputPath, 0755))

	// Create some artifacts
	prdFile := filepath.Join(outputPath, "prd.md")
	require.NoError(t, os.WriteFile(prdFile, []byte("# PRD"), 0644))

	archFile := filepath.Join(outputPath, "architecture.md")
	require.NoError(t, os.WriteFile(archFile, []byte("# Architecture"), 0644))

	result, err := Scan(tmpDir)
	require.NoError(t, err)

	assert.True(t, result.IsBMAD)
	assert.True(t, result.HasBmadFolder)
	assert.True(t, result.HasOutputFolder)
	assert.Equal(t, TypeBrownfield, result.ProjectType)
	assert.Equal(t, "6.0.0", result.BmadVersion)
	assert.True(t, result.BmadCompatible)
	assert.Len(t, result.ExistingArtifacts, 2)

	// Check artifact types
	var foundPRD, foundArch bool
	for _, artifact := range result.ExistingArtifacts {
		if artifact.Type == "prd" {
			foundPRD = true
			assert.Equal(t, "prd.md", artifact.Name)
			assert.Equal(t, filepath.Join("planning-artifacts", "prd.md"), artifact.Path)
		}
		if artifact.Type == "architecture" {
			foundArch = true
			assert.Equal(t, "architecture.md", artifact.Name)
		}
	}
	assert.True(t, foundPRD, "Should detect PRD artifact")
	assert.True(t, foundArch, "Should detect architecture artifact")
}

func TestScan_IncompatibleVersion(t *testing.T) {
	tmpDir := t.TempDir()
	bmadPath := filepath.Join(tmpDir, "_bmad")
	require.NoError(t, os.MkdirAll(bmadPath, 0755))

	// Create manifest with old version
	manifestPath := filepath.Join(bmadPath, "_config")
	require.NoError(t, os.MkdirAll(manifestPath, 0755))
	manifestFile := filepath.Join(manifestPath, "manifest.yaml")
	manifestContent := "version: 5.9.9\n"
	require.NoError(t, os.WriteFile(manifestFile, []byte(manifestContent), 0644))

	result, err := Scan(tmpDir)
	require.NoError(t, err)

	assert.True(t, result.IsBMAD)
	assert.Equal(t, "5.9.9", result.BmadVersion)
	assert.False(t, result.BmadCompatible)
}

func TestScan_GreenfieldWithEmptyOutputFolder(t *testing.T) {
	// _bmad-output exists but has no artifacts
	tmpDir := t.TempDir()
	bmadPath := filepath.Join(tmpDir, "_bmad")
	require.NoError(t, os.MkdirAll(bmadPath, 0755))

	// Create manifest
	manifestPath := filepath.Join(bmadPath, "_config")
	require.NoError(t, os.MkdirAll(manifestPath, 0755))
	manifestFile := filepath.Join(manifestPath, "manifest.yaml")
	manifestContent := "version: 6.1.0\n"
	require.NoError(t, os.WriteFile(manifestFile, []byte(manifestContent), 0644))

	// Create empty _bmad-output
	outputPath := filepath.Join(tmpDir, "_bmad-output", "planning-artifacts")
	require.NoError(t, os.MkdirAll(outputPath, 0755))

	result, err := Scan(tmpDir)
	require.NoError(t, err)

	assert.True(t, result.HasOutputFolder)
	assert.Equal(t, TypeGreenfield, result.ProjectType, "Empty output folder should be greenfield")
	assert.Empty(t, result.ExistingArtifacts)
}

func TestReadBmadVersion(t *testing.T) {
	tmpDir := t.TempDir()
	manifestFile := filepath.Join(tmpDir, "manifest.yaml")

	t.Run("valid manifest", func(t *testing.T) {
		content := "version: 6.2.1\nname: test\n"
		require.NoError(t, os.WriteFile(manifestFile, []byte(content), 0644))

		version, err := readBmadVersion(manifestFile)
		require.NoError(t, err)
		assert.Equal(t, "6.2.1", version)
	})

	t.Run("missing file", func(t *testing.T) {
		version, err := readBmadVersion(filepath.Join(tmpDir, "nonexistent.yaml"))
		assert.Error(t, err)
		assert.Empty(t, version)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		content := "invalid: yaml: content::"
		require.NoError(t, os.WriteFile(manifestFile, []byte(content), 0644))

		version, err := readBmadVersion(manifestFile)
		assert.Error(t, err)
		assert.Empty(t, version)
	})
}

func TestDetectArtifactType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"prd.md", "prd"},
		{"product-prd-v2.md", "prd"},
		{"architecture.md", "architecture"},
		{"system-architecture-design.md", "architecture"},
		{"epic-1.md", "epics"},
		{"epic-list.md", "epics"},
		{"ux-design.md", "ux-design"},
		{"product-brief.md", "product-brief"},
		{"random-doc.md", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := detectArtifactType(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestScanArtifacts_MultipleTypes(t *testing.T) {
	tmpDir := t.TempDir()
	planningPath := filepath.Join(tmpDir, "planning-artifacts")
	require.NoError(t, os.MkdirAll(planningPath, 0755))

	// Create various artifacts
	files := map[string]string{
		"prd.md":           "prd",
		"architecture.md":  "architecture",
		"epic-1.md":        "epics",
		"ux-design.md":     "ux-design",
		"product-brief.md": "product-brief",
		"notes.md":         "other",
		"readme.txt":       "", // Should be ignored (not .md)
	}

	for filename := range files {
		filePath := filepath.Join(planningPath, filename)
		require.NoError(t, os.WriteFile(filePath, []byte("content"), 0644))
	}

	artifacts := scanArtifacts(tmpDir)

	// Should only find .md files (6 total)
	assert.Len(t, artifacts, 6)

	// Verify all types are detected
	typeMap := make(map[string]bool)
	for _, artifact := range artifacts {
		typeMap[artifact.Type] = true
	}

	assert.True(t, typeMap["prd"])
	assert.True(t, typeMap["architecture"])
	assert.True(t, typeMap["epics"])
	assert.True(t, typeMap["ux-design"])
	assert.True(t, typeMap["product-brief"])
	assert.True(t, typeMap["other"])
}
// TestCompareVersions tests the semantic version comparison logic
func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int // -1, 0, or 1
	}{
		{"equal versions", "6.0.0", "6.0.0", 0},
		{"v1 greater major", "7.0.0", "6.0.0", 1},
		{"v1 lesser major", "5.0.0", "6.0.0", -1},
		{"v1 greater minor", "6.1.0", "6.0.0", 1},
		{"v1 lesser minor", "6.0.0", "6.1.0", -1},
		{"v1 greater patch", "6.0.1", "6.0.0", 1},
		{"v1 lesser patch", "6.0.0", "6.0.1", -1},
		
		// CRITICAL BUG FIX: Version 10+ comparison
		{"version 10 vs 6 (major)", "10.0.0", "6.0.0", 1},
		{"version 10 vs 9 (major)", "10.0.0", "9.0.0", 1},
		{"version 20 vs 9 (major)", "20.0.0", "9.0.0", 1},
		{"version 6 vs 10 (major)", "6.0.0", "10.0.0", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			assert.Equal(t, tt.expected, result, 
				"compareVersions(%s, %s) = %d, want %d", 
				tt.v1, tt.v2, result, tt.expected)
		})
	}
}

// TestIsVersionCompatible tests the version compatibility check
func TestIsVersionCompatible(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		minVersion string
		compatible bool
	}{
		{"exact match", "6.0.0", "6.0.0", true},
		{"greater major", "7.0.0", "6.0.0", true},
		{"greater minor", "6.1.0", "6.0.0", true},
		{"greater patch", "6.0.1", "6.0.0", true},
		{"lesser major", "5.0.0", "6.0.0", false},
		{"lesser minor", "6.0.0", "6.1.0", false},
		{"lesser patch", "6.0.0", "6.0.1", false},
		
		// CRITICAL BUG FIX: BMAD version 10+ should be compatible with 6.0.0
		{"version 10 compatible", "10.0.0", "6.0.0", true},
		{"version 15 compatible", "15.2.1", "6.0.0", true},
		{"version 100 compatible", "100.0.0", "6.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isVersionCompatible(tt.version, tt.minVersion)
			assert.Equal(t, tt.compatible, result,
				"isVersionCompatible(%s, %s) = %v, want %v",
				tt.version, tt.minVersion, result, tt.compatible)
		})
	}
}

// TestParseVersion tests the version parsing logic
func TestParseVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected [3]int
	}{
		{"standard version", "6.0.0", [3]int{6, 0, 0}},
		{"version with minor", "6.1.0", [3]int{6, 1, 0}},
		{"version with patch", "6.1.2", [3]int{6, 1, 2}},
		{"double digit major", "10.5.3", [3]int{10, 5, 3}},
		{"triple digit major", "123.45.67", [3]int{123, 45, 67}},
		{"incomplete version", "6.1", [3]int{6, 1, 0}},
		{"major only", "6", [3]int{6, 0, 0}},
		{"invalid version", "invalid", [3]int{0, 0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseVersion(tt.version)
			assert.Equal(t, tt.expected, result,
				"parseVersion(%s) = %v, want %v",
				tt.version, result, tt.expected)
		})
	}
}
