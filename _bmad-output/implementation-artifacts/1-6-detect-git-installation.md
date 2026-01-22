# Story 1.6: Detect Git Installation

Status: review

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

- [x] **Task 1: Implement Git detection in Golang** (AC: #1, #2)
  - [x] Create `internal/checkpoint/detector.go`
  - [x] Execute `git --version` and parse output
  - [x] Handle "command not found" error gracefully
  - [x] Return structured detection result

- [x] **Task 2: Implement version parsing and comparison** (AC: #1)
  - [x] Parse version from output (e.g., "git version 2.39.0")
  - [x] Compare against minimum version (2.0)
  - [x] Return compatibility status

- [x] **Task 3: Integrate with detectDependencies handler** (AC: #1)
  - [x] Add Git result to existing handler from Story 1.4
  - [x] Return both OpenCode and Git status

- [x] **Task 4: Update frontend types** (AC: all)
  - [x] Add GitDetection type to DependencyDetectionResult
  - [x] Update UI to display Git status

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

claude-3-7-sonnet-20250219

### Completion Notes List

- ✅ Implemented Git detection in `internal/checkpoint/detector.go` following the same pattern as OpenCode detection
- ✅ Created comprehensive unit tests for version parsing, comparison logic, and detection flow
- ✅ Integrated Git detection into existing `project.detectDependencies` handler alongside OpenCode detection
- ✅ Updated frontend TypeScript types to include Git detection results (made `git` field required in `DependencyDetectionResult`)
- ✅ All tests passing: Go tests (checkpoint, server packages) and TypeScript type checking
- ✅ Security compliance: No credential access, Git operations use system credentials transparently (NFR-S1, NFR-I7)
- ✅ Version comparison supports Git 2.0+ requirement with flexible parsing (handles various Git output formats)
- ✅ Error handling implemented for Git not found and version check failures

### File List

**Created:**
- `apps/core/internal/checkpoint/detector.go` - Git detection implementation
- `apps/core/internal/checkpoint/detector_test.go` - Comprehensive unit tests

**Modified:**
- `apps/core/internal/server/project_handlers.go` - Added Git detection to handler
- `apps/core/internal/server/project_handlers_test.go` - Updated tests to verify Git result
- `apps/desktop/src/renderer/src/types/dependencies.ts` - Made `git` field required, added documentation
- `apps/desktop/src/preload/index.ts` - Made `git` field required in API type
- `apps/desktop/src/preload/index.d.ts` - Made `git` field required in type declaration
- `_bmad-output/implementation-artifacts/sprint-status.yaml` - Updated story status to in-progress

### Change Log

| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Initial implementation of Git detection | Story 1-6 implementation complete |
