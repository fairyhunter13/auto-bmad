# Story 1.4: Detect and Validate OpenCode CLI

Status: ready-for-dev

## Story

As a **user launching Auto-BMAD**,
I want **the system to detect my OpenCode CLI installation and verify compatibility**,
So that **I know my environment is correctly configured before starting journeys**.

## Acceptance Criteria

1. **Given** OpenCode CLI is installed and in PATH  
   **When** Auto-BMAD performs dependency detection  
   **Then** OpenCode is detected with version displayed  
   **And** compatibility is verified against minimum version (v0.1.0+)  
   **And** the detection result is returned via JSON-RPC `project.detectDependencies`

2. **Given** OpenCode CLI is NOT in PATH  
   **When** Auto-BMAD performs dependency detection  
   **Then** a clear error message is displayed: "OpenCode CLI not found"  
   **And** installation instructions are suggested

3. **Given** OpenCode version is incompatible  
   **When** Auto-BMAD performs dependency detection  
   **Then** a warning is displayed with detected version and minimum required version

## Tasks / Subtasks

- [ ] **Task 1: Implement OpenCode detection in Golang** (AC: #1, #2)
  - [ ] Create `internal/opencode/detector.go`
  - [ ] Execute `opencode --version` and parse output
  - [ ] Handle "command not found" error gracefully
  - [ ] Return structured detection result

- [ ] **Task 2: Implement version parsing and comparison** (AC: #1, #3)
  - [ ] Parse semantic version from output (e.g., "opencode v0.1.5")
  - [ ] Compare against minimum version (v0.1.0)
  - [ ] Return compatibility status

- [ ] **Task 3: Create JSON-RPC handler** (AC: #1)
  - [ ] Register `project.detectDependencies` method
  - [ ] Return OpenCode status in response
  - [ ] Include version, path, and compatibility

- [ ] **Task 4: Add frontend API and types** (AC: all)
  - [ ] Add `window.api.project.detectDependencies()` to preload
  - [ ] Create TypeScript types for detection result
  - [ ] Handle error states in UI

## Dev Notes

### Architecture Requirements

**Source:** [architecture.md#Technical Constraints & Dependencies]

| Dependency | Minimum Version | Detection Command |
|------------|-----------------|-------------------|
| OpenCode CLI | v0.1.0+ | `opencode --version` |

### Detection Implementation

```go
// internal/opencode/detector.go

package opencode

import (
    "os/exec"
    "regexp"
    "strings"
)

type DetectionResult struct {
    Found       bool   `json:"found"`
    Version     string `json:"version,omitempty"`
    Path        string `json:"path,omitempty"`
    Compatible  bool   `json:"compatible"`
    MinVersion  string `json:"minVersion"`
    Error       string `json:"error,omitempty"`
}

const MinimumVersion = "0.1.0"

func Detect() (*DetectionResult, error) {
    result := &DetectionResult{
        MinVersion: MinimumVersion,
    }
    
    // Find opencode in PATH
    path, err := exec.LookPath("opencode")
    if err != nil {
        result.Found = false
        result.Error = "OpenCode CLI not found in PATH"
        return result, nil
    }
    result.Path = path
    result.Found = true
    
    // Get version
    cmd := exec.Command("opencode", "--version")
    output, err := cmd.Output()
    if err != nil {
        result.Error = "Failed to get OpenCode version"
        return result, nil
    }
    
    // Parse version (e.g., "opencode v0.1.5" or "opencode version 0.1.5")
    version := parseVersion(string(output))
    result.Version = version
    
    // Check compatibility
    result.Compatible = isCompatible(version, MinimumVersion)
    
    return result, nil
}

func parseVersion(output string) string {
    // Match patterns like "v0.1.5", "0.1.5", "version 0.1.5"
    re := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)
    matches := re.FindStringSubmatch(output)
    if len(matches) > 1 {
        return matches[1]
    }
    return strings.TrimSpace(output)
}

func isCompatible(version, minVersion string) bool {
    // Simple semver comparison
    return compareVersions(version, minVersion) >= 0
}

func compareVersions(v1, v2 string) int {
    // Parse and compare major.minor.patch
    // Return: -1 if v1 < v2, 0 if equal, 1 if v1 > v2
    // Implementation left to developer
}
```

### JSON-RPC Handler

```go
// internal/server/handlers.go (additions)

func (s *Server) handleDetectDependencies(params json.RawMessage) (interface{}, error) {
    opencodeResult, _ := opencode.Detect()
    gitResult, _ := git.Detect()  // Story 1.6
    
    return map[string]interface{}{
        "opencode": opencodeResult,
        "git":      gitResult,
    }, nil
}

// Register in server initialization
server.RegisterHandler("project.detectDependencies", s.handleDetectDependencies)
```

### Frontend Types

```typescript
// src/renderer/types/dependencies.ts

export interface OpenCodeDetection {
  found: boolean;
  version?: string;
  path?: string;
  compatible: boolean;
  minVersion: string;
  error?: string;
}

export interface DependencyDetectionResult {
  opencode: OpenCodeDetection;
  git: GitDetection; // Story 1.6
}
```

### Preload API Addition

```typescript
// src/preload/index.ts (additions)

const api = {
  // ... existing
  project: {
    detectDependencies: (): Promise<DependencyDetectionResult> => 
      ipcRenderer.invoke('rpc:call', 'project.detectDependencies'),
  },
};
```

### Error Messages (User-Facing)

| Condition | Message |
|-----------|---------|
| Not found | "OpenCode CLI not found. Please install OpenCode and ensure it's in your PATH." |
| Version too old | "OpenCode version {detected} is below minimum required version {min}. Please update OpenCode." |
| Detection failed | "Could not verify OpenCode installation. Please check that 'opencode --version' works." |

### File Structure

```
apps/core/internal/
└── opencode/
    ├── detector.go       # Detection logic
    ├── detector_test.go  # Unit tests
    └── version.go        # Version comparison utilities
```

### Testing Requirements

1. Mock `exec.LookPath` and `exec.Command` for unit tests
2. Test version parsing with various formats
3. Test version comparison logic
4. Test error handling when OpenCode not found

### Dependencies

- **Story 1.2**: JSON-RPC server must be implemented
- **Story 1.3**: IPC bridge must be working

### References

- [architecture.md#External Dependencies] - Minimum versions
- [prd.md#FR11] - System can detect installed OpenCode CLI
- [prd.md#FR12] - System can detect OpenCode CLI version and verify compatibility

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Completion Notes List

- 

### Change Log

| Date | Change | Reason |
|------|--------|--------|
