package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBrownfield_ImplementationArtifacts tests detection of implementation artifacts
func TestBrownfield_ImplementationArtifacts(t *testing.T) {
	tmpDir := t.TempDir()
	bmadPath := filepath.Join(tmpDir, "_bmad")
	require.NoError(t, os.MkdirAll(bmadPath, 0755))

	// Create manifest
	manifestPath := filepath.Join(bmadPath, "_config")
	require.NoError(t, os.MkdirAll(manifestPath, 0755))
	manifestFile := filepath.Join(manifestPath, "manifest.yaml")
	manifestContent := "version: 6.1.0\n"
	require.NoError(t, os.WriteFile(manifestFile, []byte(manifestContent), 0644))

	// Create implementation artifacts
	implPath := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	require.NoError(t, os.MkdirAll(implPath, 0755))

	storyFile := filepath.Join(implPath, "1-1-user-auth.md")
	require.NoError(t, os.WriteFile(storyFile, []byte("# Story"), 0644))

	result, err := Scan(tmpDir)
	require.NoError(t, err)

	// Should detect as greenfield (we only scan planning-artifacts)
	assert.Equal(t, TypeGreenfield, result.ProjectType)
}

// TestBrownfield_MixedArtifacts tests detection with multiple artifact types
func TestBrownfield_MixedArtifacts(t *testing.T) {
	tmpDir := t.TempDir()
	bmadPath := filepath.Join(tmpDir, "_bmad")
	require.NoError(t, os.MkdirAll(bmadPath, 0755))

	// Create manifest
	manifestPath := filepath.Join(bmadPath, "_config")
	require.NoError(t, os.MkdirAll(manifestPath, 0755))
	manifestFile := filepath.Join(manifestPath, "manifest.yaml")
	manifestContent := "version: 6.1.0\n"
	require.NoError(t, os.WriteFile(manifestFile, []byte(manifestContent), 0644))

	// Create planning artifacts
	planningPath := filepath.Join(tmpDir, "_bmad-output", "planning-artifacts")
	require.NoError(t, os.MkdirAll(planningPath, 0755))

	files := []string{
		"prd.md",
		"architecture.md",
		"epic-1-foundation.md",
		"epic-2-features.md",
		"ux-design-system.md",
	}

	for _, file := range files {
		filePath := filepath.Join(planningPath, file)
		require.NoError(t, os.WriteFile(filePath, []byte("content"), 0644))
	}

	result, err := Scan(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, TypeBrownfield, result.ProjectType)
	assert.Len(t, result.ExistingArtifacts, 5)

	// Verify all artifact types are detected
	typeCount := make(map[string]int)
	for _, artifact := range result.ExistingArtifacts {
		typeCount[artifact.Type]++
	}

	assert.Equal(t, 1, typeCount["prd"])
	assert.Equal(t, 1, typeCount["architecture"])
	assert.Equal(t, 2, typeCount["epics"])
	assert.Equal(t, 1, typeCount["ux-design"])
}

// TestBrownfield_NonMarkdownFiles tests that non-.md files are ignored
func TestBrownfield_NonMarkdownFiles(t *testing.T) {
	tmpDir := t.TempDir()
	bmadPath := filepath.Join(tmpDir, "_bmad")
	require.NoError(t, os.MkdirAll(bmadPath, 0755))

	// Create manifest
	manifestPath := filepath.Join(bmadPath, "_config")
	require.NoError(t, os.MkdirAll(manifestPath, 0755))
	manifestFile := filepath.Join(manifestPath, "manifest.yaml")
	manifestContent := "version: 6.1.0\n"
	require.NoError(t, os.WriteFile(manifestFile, []byte(manifestContent), 0644))

	// Create artifacts with different extensions
	planningPath := filepath.Join(tmpDir, "_bmad-output", "planning-artifacts")
	require.NoError(t, os.MkdirAll(planningPath, 0755))

	// Create .md file
	mdFile := filepath.Join(planningPath, "prd.md")
	require.NoError(t, os.WriteFile(mdFile, []byte("content"), 0644))

	// Create non-.md files (should be ignored)
	txtFile := filepath.Join(planningPath, "notes.txt")
	require.NoError(t, os.WriteFile(txtFile, []byte("content"), 0644))

	pdfFile := filepath.Join(planningPath, "design.pdf")
	require.NoError(t, os.WriteFile(pdfFile, []byte("content"), 0644))

	result, err := Scan(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, TypeBrownfield, result.ProjectType)
	assert.Len(t, result.ExistingArtifacts, 1, "Should only count .md files")
	assert.Equal(t, "prd.md", result.ExistingArtifacts[0].Name)
}

// TestBrownfield_SubdirectoriesIgnored tests that subdirectories are not scanned
func TestBrownfield_SubdirectoriesIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	bmadPath := filepath.Join(tmpDir, "_bmad")
	require.NoError(t, os.MkdirAll(bmadPath, 0755))

	// Create manifest
	manifestPath := filepath.Join(bmadPath, "_config")
	require.NoError(t, os.MkdirAll(manifestPath, 0755))
	manifestFile := filepath.Join(manifestPath, "manifest.yaml")
	manifestContent := "version: 6.1.0\n"
	require.NoError(t, os.WriteFile(manifestFile, []byte(manifestContent), 0644))

	// Create planning artifacts
	planningPath := filepath.Join(tmpDir, "_bmad-output", "planning-artifacts")
	require.NoError(t, os.MkdirAll(planningPath, 0755))

	// Create file at root
	rootFile := filepath.Join(planningPath, "prd.md")
	require.NoError(t, os.WriteFile(rootFile, []byte("content"), 0644))

	// Create subdirectory with file
	subDir := filepath.Join(planningPath, "archive")
	require.NoError(t, os.MkdirAll(subDir, 0755))
	subFile := filepath.Join(subDir, "old-prd.md")
	require.NoError(t, os.WriteFile(subFile, []byte("content"), 0644))

	result, err := Scan(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, TypeBrownfield, result.ProjectType)
	assert.Len(t, result.ExistingArtifacts, 1, "Should not scan subdirectories")
	assert.Equal(t, "prd.md", result.ExistingArtifacts[0].Name)
}

// TestGreenfield_OnlyNonPlanningFolders tests greenfield detection with other folders
func TestGreenfield_OnlyNonPlanningFolders(t *testing.T) {
	tmpDir := t.TempDir()
	bmadPath := filepath.Join(tmpDir, "_bmad")
	require.NoError(t, os.MkdirAll(bmadPath, 0755))

	// Create manifest
	manifestPath := filepath.Join(bmadPath, "_config")
	require.NoError(t, os.MkdirAll(manifestPath, 0755))
	manifestFile := filepath.Join(manifestPath, "manifest.yaml")
	manifestContent := "version: 6.1.0\n"
	require.NoError(t, os.WriteFile(manifestFile, []byte(manifestContent), 0644))

	// Create _bmad-output but with different folders
	implPath := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	require.NoError(t, os.MkdirAll(implPath, 0755))

	// Create file in implementation folder
	storyFile := filepath.Join(implPath, "1-1-user-auth.md")
	require.NoError(t, os.WriteFile(storyFile, []byte("# Story"), 0644))

	result, err := Scan(tmpDir)
	require.NoError(t, err)

	// Should be greenfield because planning-artifacts is empty
	assert.Equal(t, TypeGreenfield, result.ProjectType)
	assert.Empty(t, result.ExistingArtifacts)
}
