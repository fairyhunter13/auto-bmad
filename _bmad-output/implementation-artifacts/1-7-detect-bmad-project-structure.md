# Story 1.7: Detect BMAD Project Structure

Status: ready-for-dev

## Story

As a **user opening a project folder**,
I want **Auto-BMAD to detect if it's a BMAD project and identify greenfield vs brownfield**,
So that **the system knows what journey options are available**.

## Acceptance Criteria

1. **Given** a folder with `_bmad/` folder present  
   **When** Auto-BMAD scans the project  
   **Then** the project is identified as a BMAD project  
   **And** BMAD version is read from `_bmad/_config/manifest.yaml`  
   **And** compatibility is verified (6.0.0+)

2. **Given** a folder with `_bmad/` but NO `_bmad-output/` folder  
   **When** Auto-BMAD scans the project  
   **Then** project type is identified as "greenfield"  
   **And** full journey options are available

3. **Given** a folder with both `_bmad/` and `_bmad-output/` with artifacts  
   **When** Auto-BMAD scans the project  
   **Then** project type is identified as "brownfield"  
   **And** existing artifacts are listed for context  
   **And** partial journey options are available

4. **Given** a folder without `_bmad/` folder  
   **When** Auto-BMAD scans the project  
   **Then** a message indicates "Not a BMAD project"  
   **And** user is prompted to run `bmad-init` first

## Tasks / Subtasks

- [ ] **Task 1: Implement project structure detection** (AC: #1, #4)
  - [ ] Create `internal/project/detector.go`
  - [ ] Check for `_bmad/` folder existence
  - [ ] Check for `_bmad/_config/manifest.yaml`
  - [ ] Parse BMAD version from manifest

- [ ] **Task 2: Implement greenfield/brownfield detection** (AC: #2, #3)
  - [ ] Check for `_bmad-output/` folder
  - [ ] Scan for existing artifacts if brownfield
  - [ ] Categorize artifacts by type (PRD, architecture, etc.)

- [ ] **Task 3: Create JSON-RPC handler** (AC: all)
  - [ ] Register `project.scan` method
  - [ ] Accept project path as parameter
  - [ ] Return comprehensive project info

- [ ] **Task 4: Add frontend API** (AC: all)
  - [ ] Add `window.api.project.scan(path)` to preload
  - [ ] Create TypeScript types for project info

## Dev Notes

### Architecture Requirements

**Source:** [architecture.md#External Dependencies]

| Dependency | Minimum Version | Detection |
|------------|-----------------|-----------|
| BMAD | 6.0.0+ | `_bmad/_config/manifest.yaml` |

### Detection Implementation

```go
// internal/project/detector.go

package project

import (
    "os"
    "path/filepath"
    
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
    IsBMAD           bool        `json:"isBmad"`
    ProjectType      ProjectType `json:"projectType"`
    BmadVersion      string      `json:"bmadVersion,omitempty"`
    BmadCompatible   bool        `json:"bmadCompatible"`
    MinBmadVersion   string      `json:"minBmadVersion"`
    Path             string      `json:"path"`
    HasBmadFolder    bool        `json:"hasBmadFolder"`
    HasOutputFolder  bool        `json:"hasOutputFolder"`
    ExistingArtifacts []Artifact `json:"existingArtifacts,omitempty"`
    Error            string      `json:"error,omitempty"`
}

const MinBmadVersion = "6.0.0"

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

func detectArtifactType(filename string) string {
    switch {
    case strings.Contains(filename, "prd"):
        return "prd"
    case strings.Contains(filename, "architecture"):
        return "architecture"
    case strings.Contains(filename, "epic"):
        return "epics"
    case strings.Contains(filename, "ux"):
        return "ux-design"
    case strings.Contains(filename, "product-brief"):
        return "product-brief"
    default:
        return "other"
    }
}
```

### JSON-RPC Handler

```go
// internal/server/handlers.go (additions)

type ScanParams struct {
    Path string `json:"path"`
}

func (s *Server) handleProjectScan(params json.RawMessage) (interface{}, error) {
    var p ScanParams
    if err := json.Unmarshal(params, &p); err != nil {
        return nil, err
    }
    
    return project.Scan(p.Path)
}

// Register
server.RegisterHandler("project.scan", s.handleProjectScan)
```

### Frontend Types

```typescript
// src/renderer/types/project.ts

export type ProjectType = 'not-bmad' | 'greenfield' | 'brownfield';

export interface Artifact {
  name: string;
  path: string;
  type: string;
  modified?: string;
}

export interface ProjectScanResult {
  isBmad: boolean;
  projectType: ProjectType;
  bmadVersion?: string;
  bmadCompatible: boolean;
  minBmadVersion: string;
  path: string;
  hasBmadFolder: boolean;
  hasOutputFolder: boolean;
  existingArtifacts?: Artifact[];
  error?: string;
}
```

### Preload API Addition

```typescript
// src/preload/index.ts (additions)

const api = {
  project: {
    // ... existing
    scan: (path: string): Promise<ProjectScanResult> => 
      ipcRenderer.invoke('rpc:call', 'project.scan', { path }),
  },
};
```

### Journey Options by Project Type

| Project Type | Available Journeys |
|--------------|-------------------|
| **Greenfield** | Full BMAD journey (Brainstorming → Implementation) |
| **Brownfield** | Continue from existing artifacts, skip completed phases |
| **Not BMAD** | None - prompt to run `bmad-init` |

### Artifact Type Detection

| Filename Pattern | Type |
|------------------|------|
| `*prd*.md` | prd |
| `*architecture*.md` | architecture |
| `*epic*.md` | epics |
| `*ux*.md` | ux-design |
| `*product-brief*.md` | product-brief |
| Other `.md` files | other |

### File Structure

```
apps/core/internal/
└── project/
    ├── detector.go       # Project detection
    ├── detector_test.go  # Unit tests
    └── types.go          # Shared types
```

### Go Dependency

Add YAML parsing:
```bash
cd apps/core
go get gopkg.in/yaml.v3
```

### Testing Requirements

1. Test greenfield detection (no _bmad-output/)
2. Test brownfield detection (with artifacts)
3. Test not-BMAD detection (no _bmad/)
4. Test version parsing from manifest.yaml
5. Test artifact type detection

### Dependencies

- **Story 1.3**: IPC bridge must be working

### References

- [architecture.md#External Dependencies] - BMAD 6.0.0+ requirement
- [prd.md#FR45] - Detect _bmad/ folder
- [prd.md#FR46] - Detect _bmad-output/ folder
- [prd.md#FR47] - Detect project type (greenfield vs brownfield)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Completion Notes List

- 

### Change Log

| Date | Change | Reason |
|------|--------|--------|
