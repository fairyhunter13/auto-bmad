# Story 1.5: Detect OpenCode Profiles

Status: ready-for-dev

## Story

As a **user with multiple OpenCode profiles configured**,
I want **Auto-BMAD to detect my available profiles from ~/.bash_aliases**,
So that **I can select which profile to use for my project**.

## Acceptance Criteria

1. **Given** the user has OpenCode profiles defined in `~/.bash_aliases`  
   **When** Auto-BMAD parses the configuration  
   **Then** all profile names are extracted and listed  
   **And** the default profile is identified  
   **And** the list is returned via JSON-RPC `opencode.getProfiles`

2. **Given** a profile is misconfigured or has missing credentials  
   **When** Auto-BMAD validates the profile  
   **Then** a clear warning message is displayed (NFR-I5)  
   **And** the profile is marked as "unavailable" in the UI

3. **Given** no profiles are found  
   **When** Auto-BMAD parses the configuration  
   **Then** the system uses the global OpenCode default  
   **And** a message indicates "Using default OpenCode configuration"

## Tasks / Subtasks

- [ ] **Task 1: Parse ~/.bash_aliases for OpenCode profiles** (AC: #1, #3)
  - [ ] Read `~/.bash_aliases` file
  - [ ] Extract alias definitions matching opencode patterns
  - [ ] Parse profile names from alias definitions
  - [ ] Handle file not found gracefully

- [ ] **Task 2: Validate profile availability** (AC: #2)
  - [ ] Test each profile with `opencode --profile {name} --version`
  - [ ] Mark unavailable profiles with error reason
  - [ ] Return validation status per profile

- [ ] **Task 3: Create JSON-RPC handler** (AC: #1)
  - [ ] Register `opencode.getProfiles` method
  - [ ] Return profile list with availability status
  - [ ] Include default profile indicator

- [ ] **Task 4: Add frontend API** (AC: all)
  - [ ] Add `window.api.opencode.getProfiles()` to preload
  - [ ] Create TypeScript types for profile data

## Dev Notes

### Profile Detection Strategy

OpenCode profiles are typically defined as shell aliases in `~/.bash_aliases`:

```bash
# Example ~/.bash_aliases content
alias opencode-anthropic='ANTHROPIC_API_KEY=sk-xxx opencode'
alias opencode-openai='OPENAI_API_KEY=sk-xxx opencode'
alias opencode-local='opencode --provider ollama'
```

### Implementation

```go
// internal/opencode/profiles.go

package opencode

import (
    "bufio"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

type Profile struct {
    Name        string `json:"name"`
    Alias       string `json:"alias"`       // Full alias command
    Available   bool   `json:"available"`
    Error       string `json:"error,omitempty"`
    IsDefault   bool   `json:"isDefault"`
}

type ProfilesResult struct {
    Profiles     []Profile `json:"profiles"`
    DefaultFound bool      `json:"defaultFound"`
    Source       string    `json:"source"` // e.g., "~/.bash_aliases"
}

func GetProfiles() (*ProfilesResult, error) {
    result := &ProfilesResult{
        Profiles: []Profile{},
        Source:   "~/.bash_aliases",
    }
    
    homeDir, _ := os.UserHomeDir()
    aliasFile := filepath.Join(homeDir, ".bash_aliases")
    
    file, err := os.Open(aliasFile)
    if err != nil {
        // No aliases file - use default
        result.DefaultFound = true
        result.Profiles = append(result.Profiles, Profile{
            Name:      "default",
            Available: true,
            IsDefault: true,
        })
        return result, nil
    }
    defer file.Close()
    
    // Pattern: alias opencode-{name}='...'
    aliasPattern := regexp.MustCompile(`^alias\s+(opencode-\w+)=['"](.+)['"]`)
    
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        matches := aliasPattern.FindStringSubmatch(line)
        if len(matches) == 3 {
            aliasName := matches[1]
            profileName := strings.TrimPrefix(aliasName, "opencode-")
            
            profile := Profile{
                Name:  profileName,
                Alias: matches[2],
            }
            
            // Validate profile
            profile.Available, profile.Error = validateProfile(aliasName)
            
            result.Profiles = append(result.Profiles, profile)
        }
    }
    
    // If no profiles found, add default
    if len(result.Profiles) == 0 {
        result.DefaultFound = true
        result.Profiles = append(result.Profiles, Profile{
            Name:      "default",
            Available: true,
            IsDefault: true,
        })
    }
    
    return result, nil
}

func validateProfile(aliasName string) (bool, string) {
    // Quick validation - just check if the alias can be invoked
    // Note: In practice, we might just mark all as available
    // and let actual usage fail
    return true, ""
}
```

### JSON-RPC Handler

```go
// internal/server/handlers.go (additions)

func (s *Server) handleGetProfiles(params json.RawMessage) (interface{}, error) {
    return opencode.GetProfiles()
}

// Register
server.RegisterHandler("opencode.getProfiles", s.handleGetProfiles)
```

### Frontend Types

```typescript
// src/renderer/types/opencode.ts

export interface OpenCodeProfile {
  name: string;
  alias: string;
  available: boolean;
  error?: string;
  isDefault: boolean;
}

export interface ProfilesResult {
  profiles: OpenCodeProfile[];
  defaultFound: boolean;
  source: string;
}
```

### Preload API Addition

```typescript
// src/preload/index.ts (additions)

const api = {
  // ... existing
  opencode: {
    getProfiles: (): Promise<ProfilesResult> => 
      ipcRenderer.invoke('rpc:call', 'opencode.getProfiles'),
  },
};
```

### Edge Cases

| Scenario | Behavior |
|----------|----------|
| No ~/.bash_aliases file | Return single "default" profile |
| File exists but no opencode aliases | Return single "default" profile |
| Alias with special characters | Skip with warning in logs |
| Circular alias references | Treat as available (validation at runtime) |

### File Structure

```
apps/core/internal/
└── opencode/
    ├── detector.go       # From Story 1.4
    ├── profiles.go       # Profile detection
    └── profiles_test.go  # Unit tests
```

### Testing Requirements

1. Test parsing of various alias formats
2. Test handling of missing ~/.bash_aliases
3. Test default profile fallback
4. Mock file system for unit tests

### Dependencies

- **Story 1.4**: OpenCode detection must exist
- **Story 1.3**: IPC bridge must be working

### References

- [prd.md#FR13] - System can detect available OpenCode profiles
- [prd.md#FR14] - User can select which OpenCode profile to use
- [prd.md#NFR-I2] - Support multiple profiles with load-balancing
- [prd.md#NFR-I5] - Profile misconfiguration warning

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Completion Notes List

- 

### Change Log

| Date | Change | Reason |
|------|--------|--------|
