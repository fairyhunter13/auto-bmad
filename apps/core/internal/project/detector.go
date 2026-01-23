package project

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type ProjectType string

const (
	TypeNotBMAD    ProjectType = "not-bmad"
	TypeGreenfield ProjectType = "greenfield"
	TypeBrownfield ProjectType = "brownfield"
)

type Artifact struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Type     string `json:"type"` // prd, architecture, epics, etc.
	Modified string `json:"modified"`
}

type ProjectScanResult struct {
	IsBMAD            bool        `json:"isBmad"`
	ProjectType       ProjectType `json:"projectType"`
	BmadVersion       string      `json:"bmadVersion,omitempty"`
	BmadCompatible    bool        `json:"bmadCompatible"`
	MinBmadVersion    string      `json:"minBmadVersion"`
	Path              string      `json:"path"`
	HasBmadFolder     bool        `json:"hasBmadFolder"`
	HasOutputFolder   bool        `json:"hasOutputFolder"`
	ExistingArtifacts []Artifact  `json:"existingArtifacts,omitempty"`
	Error             string      `json:"error,omitempty"`
}

const MinBmadVersion = "6.0.0"

// Scan detects BMAD project structure and returns project information
func Scan(projectPath string) (*ProjectScanResult, error) {
	result := &ProjectScanResult{
		Path:           projectPath,
		MinBmadVersion: MinBmadVersion,
		ProjectType:    TypeNotBMAD,
	}

	// Check _bmad/ folder
	bmadPath := filepath.Join(projectPath, "_bmad")
	if _, err := os.Stat(bmadPath); os.IsNotExist(err) {
		result.IsBMAD = false
		result.HasBmadFolder = false
		result.Error = "Not a BMAD project. Run 'bmad-init' to initialize."
		return result, nil
	}
	result.HasBmadFolder = true
	result.IsBMAD = true

	// Read BMAD version from manifest
	manifestPath := filepath.Join(bmadPath, "_config", "manifest.yaml")
	if version, err := readBmadVersion(manifestPath); err == nil {
		result.BmadVersion = version
		result.BmadCompatible = isVersionCompatible(version, MinBmadVersion)
	}

	// Check _bmad-output/ folder
	outputPath := filepath.Join(projectPath, "_bmad-output")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		result.HasOutputFolder = false
		result.ProjectType = TypeGreenfield
		return result, nil
	}
	result.HasOutputFolder = true

	// Scan for existing artifacts
	artifacts := scanArtifacts(outputPath)
	if len(artifacts) == 0 {
		result.ProjectType = TypeGreenfield
	} else {
		result.ProjectType = TypeBrownfield
		result.ExistingArtifacts = artifacts
	}

	return result, nil
}

// readBmadVersion reads the BMAD version from manifest.yaml
func readBmadVersion(manifestPath string) (string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", err
	}

	var manifest struct {
		Version string `yaml:"version"`
	}
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return "", err
	}

	return manifest.Version, nil
}

// isVersionCompatible checks if version meets minimum requirement
// Uses proper semantic version comparison to handle versions like 10.0.0 correctly
func isVersionCompatible(version, minVersion string) bool {
	return compareVersions(version, minVersion) >= 0
}

// compareVersions compares two semantic versions (e.g., "10.0.0" vs "6.0.0")
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	// Compare each part (major, minor, patch)
	for i := 0; i < 3; i++ {
		if parts1[i] < parts2[i] {
			return -1
		}
		if parts1[i] > parts2[i] {
			return 1
		}
	}
	return 0
}

// parseVersion parses a semantic version string into [major, minor, patch]
// Returns [0, 0, 0] if parsing fails
func parseVersion(version string) [3]int {
	var parts [3]int
	segments := strings.Split(version, ".")

	for i := 0; i < 3 && i < len(segments); i++ {
		if num, err := strconv.Atoi(segments[i]); err == nil {
			parts[i] = num
		}
	}

	return parts
}

// scanArtifacts scans _bmad-output for existing artifacts
func scanArtifacts(outputPath string) []Artifact {
	artifacts := []Artifact{}

	// Scan planning-artifacts/
	planningPath := filepath.Join(outputPath, "planning-artifacts")
	if entries, err := os.ReadDir(planningPath); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
				artifactType := detectArtifactType(entry.Name())
				artifacts = append(artifacts, Artifact{
					Name: entry.Name(),
					Path: filepath.Join("planning-artifacts", entry.Name()),
					Type: artifactType,
				})
			}
		}
	}

	return artifacts
}

// detectArtifactType determines the type of artifact based on filename
func detectArtifactType(filename string) string {
	lower := strings.ToLower(filename)
	switch {
	case strings.Contains(lower, "prd"):
		return "prd"
	case strings.Contains(lower, "architecture"):
		return "architecture"
	case strings.Contains(lower, "epic"):
		return "epics"
	case strings.Contains(lower, "ux"):
		return "ux-design"
	case strings.Contains(lower, "product-brief"):
		return "product-brief"
	default:
		return "other"
	}
}
