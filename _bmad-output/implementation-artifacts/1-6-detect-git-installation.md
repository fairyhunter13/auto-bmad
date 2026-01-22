# Story 1.6: Detect Git Installation

Status: ready-for-dev

## Story

As a **user launching Auto-BMAD**,
I want **the system to verify Git is installed**,
So that **checkpoint commits and rollback features will work correctly**.

## Acceptance Criteria

1. **Given** Git is installed (version 2.0+)  
   **When** Auto-BMAD performs dependency detection  
   **Then** Git version is detected and displayed  
   **And** the detection result is included in `project.detectDependencies` response

2. **Given** Git is NOT installed  
   **When** Auto-BMAD performs dependency detection  
   **Then** a clear error message is displayed: "Git not found"  
   **And** the error blocks journey start with explanation: "Git is required for checkpoint safety"

3. **Given** Git credentials are configured  
   **When** Auto-BMAD detects Git  
   **Then** Auto-BMAD does NOT access or store credentials (NFR-S1, NFR-I7)  
   **And** Git operations use system credentials transparently

## Tasks / Subtasks

- [ ] **Task 1: Implement Git detection in Golang** (AC: #1, #2)
  - [ ] Create `internal/checkpoint/detector.go`
  - [ ] Execute `git --version` and parse output
  - [ ] Handle "command not found" error gracefully
  - [ ] Return structured detection result

- [ ] **Task 2: Implement version parsing and comparison** (AC: #1)
  - [ ] Parse version from output (e.g., "git version 2.39.0")
  - [ ] Compare against minimum version (2.0)
  - [ ] Return compatibility status

- [ ] **Task 3: Integrate with detectDependencies handler** (AC: #1)
  - [ ] Add Git result to existing handler from Story 1.4
  - [ ] Return both OpenCode and Git status

- [ ] **Task 4: Update frontend types** (AC: all)
  - [ ] Add GitDetection type to DependencyDetectionResult
  - [ ] Update UI to display Git status

## Dev Notes

### Architecture Requirements

**Source:** [architecture.md#External Dependencies]

| Dependency | Minimum Version | Detection Command |
|------------|-----------------|-------------------|
| Git | 2.0+ | `git --version` |

### Detection Implementation

```go
// internal/checkpoint/detector.go

package checkpoint

import (
    "os/exec"
    "regexp"
)

type GitDetectionResult struct {
    Found      bool   `json:"found"`
    Version    string `json:"version,omitempty"`
    Path       string `json:"path,omitempty"`
    Compatible bool   `json:"compatible"`
    MinVersion string `json:"minVersion"`
    Error      string `json:"error,omitempty"`
}

const GitMinimumVersion = "2.0.0"

func DetectGit() (*GitDetectionResult, error) {
    result := &GitDetectionResult{
        MinVersion: GitMinimumVersion,
    }
    
    // Find git in PATH
    path, err := exec.LookPath("git")
    if err != nil {
        result.Found = false
        result.Error = "Git not found in PATH"
        return result, nil
    }
    result.Path = path
    result.Found = true
    
    // Get version
    cmd := exec.Command("git", "--version")
    output, err := cmd.Output()
    if err != nil {
        result.Error = "Failed to get Git version"
        return result, nil
    }
    
    // Parse version (e.g., "git version 2.39.0")
    version := parseGitVersion(string(output))
    result.Version = version
    
    // Check compatibility
    result.Compatible = isGitCompatible(version, GitMinimumVersion)
    
    return result, nil
}

func parseGitVersion(output string) string {
    // Match pattern "git version X.Y.Z"
    re := regexp.MustCompile(`git version (\d+\.\d+\.?\d*)`)
    matches := re.FindStringSubmatch(output)
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}

func isGitCompatible(version, minVersion string) bool {
    // Compare major.minor
    return compareVersions(version, minVersion) >= 0
}
```

### Updated detectDependencies Handler

```go
// internal/server/handlers.go (update)

func (s *Server) handleDetectDependencies(params json.RawMessage) (interface{}, error) {
    opencodeResult, _ := opencode.Detect()
    gitResult, _ := checkpoint.DetectGit()
    
    return map[string]interface{}{
        "opencode": opencodeResult,
        "git":      gitResult,
    }, nil
}
```

### Frontend Types Update

```typescript
// src/renderer/types/dependencies.ts (update)

export interface GitDetection {
  found: boolean;
  version?: string;
  path?: string;
  compatible: boolean;
  minVersion: string;
  error?: string;
}

export interface DependencyDetectionResult {
  opencode: OpenCodeDetection;
  git: GitDetection;
}
```

### Security Note (CRITICAL)

**Source:** [prd.md#NFR-S1, NFR-I7]

Auto-BMAD MUST NOT:
- Read Git credentials
- Store Git credentials
- Intercept credential prompts

Git operations MUST use system credentials transparently (SSH keys, credential helpers, etc.).

### Error Messages (User-Facing)

| Condition | Message |
|-----------|---------|
| Not found | "Git not found. Please install Git (version 2.0 or later) to enable checkpoint safety." |
| Version too old | "Git version {detected} is below minimum required version 2.0. Please update Git." |
| Journey block | "Cannot start journey: Git is required for checkpoint safety and rollback capabilities." |

### File Structure

```
apps/core/internal/
└── checkpoint/
    ├── detector.go       # Git detection
    ├── detector_test.go  # Unit tests
    └── checkpoint.go     # Future: Git operations (Story 3.6)
```

### Testing Requirements

1. Test version parsing with various Git output formats
2. Test handling when Git not installed
3. Test version comparison logic
4. Verify no credential access

### Dependencies

- **Story 1.4**: OpenCode detection (shares handler)
- **Story 1.3**: IPC bridge must be working

### References

- [architecture.md#External Dependencies] - Minimum Git version
- [architecture.md#Git Checkpoint Error Handling] - Git operations context
- [prd.md#NFR-I3] - Git version 2.0+
- [prd.md#NFR-I7] - Git operations use system credentials
- [prd.md#NFR-S1] - No credential storage

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Completion Notes List

- 

### Change Log

| Date | Change | Reason |
|------|--------|--------|
