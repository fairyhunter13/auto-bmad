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
