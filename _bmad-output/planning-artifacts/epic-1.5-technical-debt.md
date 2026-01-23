# Epic 1.5: Technical Debt & Architecture Fixes

**Epic Status:** Planned  
**Priority:** Medium  
**Created:** 2026-01-23  
**Target:** Post-Epic 2  

---

## Overview

This epic addresses **technical debt and non-critical issues** identified during the Epic 1 code review. All **critical security vulnerabilities** have been fixed in Epic 1 (8 commits, 100% tests passing). This epic focuses on **architectural compliance**, **test coverage gaps**, and **code quality improvements**.

### Why This Epic Exists

During Epic 1 implementation and code review, we identified:
- **2 Critical Issues** (architectural violations, not security)
- **7 High Priority Issues** (should fix but not blocking)
- **30+ Medium Priority Issues** (recommended improvements)
- **38+ Low Priority Issues** (optional enhancements)

These issues are **deferred to Epic 1.5** to allow Epic 2 (workflow execution) to proceed without delay.

---

## Success Criteria

✅ All critical architectural violations resolved  
✅ All high priority issues addressed  
✅ Test coverage > 80% across all packages  
✅ 100% compliance with architecture.md specifications  
✅ Zero failing tests  

---

## Stories

### 🔴 Story 1.5.1: Fix Settings Path Architecture Violation [CRITICAL]

**Priority:** Critical  
**Effort:** 4-8 hours  
**Complexity:** High  
**Blocking:** Epic 2? No (can be deferred)

#### Problem Statement

**Current Implementation:**
```
~/.autobmad/_bmad-output/.autobmad/config.json  # WRONG - Global settings
```

**Required by Architecture:**
```
<project-root>/_bmad-output/.autobmad/config.json  # Correct - Project-local
```

#### Impact

- **Violates architecture.md Line 222** (`Configuration | _bmad-output/.autobmad/config.json`)
- **Breaks Story 1-9 AC #1:** "settings are saved to `_bmad-output/.autobmad/config.json`"
- **Breaks Story 1-9 AC #2:** "last-used OpenCode profile **per project**" cannot work with global settings
- **Breaks multi-project workflows:** Users working on multiple BMAD projects will have conflicting settings

#### Root Cause

**File:** `apps/core/internal/server/settings_handlers.go:27`
```go
// WRONG:
settingsPath := filepath.Join(homeDir, ".autobmad")

// SHOULD BE:
settingsPath := filepath.Join(projectPath, "_bmad-output", ".autobmad")
```

Developer added comment: "Settings are stored globally (not per-project) to remember user preferences across sessions" - this contradicts the architecture specification and the Settings struct's `ProjectProfiles map[string]string` field.

#### Solution Design

**Phase 1: Pass Project Path Through Stack**

1. **Electron Main Process** (`apps/desktop/src/main/index.ts`)
   ```typescript
   // Add project path to backend startup
   const projectPath = await selectProjectFolder(); // From user selection
   backend.start({ projectPath });
   ```

2. **Backend Process Manager** (`apps/desktop/src/main/backend.ts`)
   ```typescript
   export function start(options: { projectPath: string }) {
     const child = spawn('autobmad', ['--project-path', options.projectPath]);
   }
   ```

3. **Go Backend Main** (`apps/core/cmd/autobmad/main.go`)
   ```go
   var projectPath string
   flag.StringVar(&projectPath, "project-path", "", "Path to BMAD project root")
   flag.Parse()
   
   if projectPath == "" {
       log.Fatal("--project-path is required")
   }
   
   // Pass to server
   srv := server.NewServer(projectPath)
   ```

4. **Server Constructor** (`apps/core/internal/server/server.go`)
   ```go
   type Server struct {
       projectPath string  // NEW field
       // ... existing fields
   }
   
   func NewServer(projectPath string) *Server {
       return &Server{
           projectPath: projectPath,
           // ...
       }
   }
   ```

5. **Settings Handlers** (`apps/core/internal/server/settings_handlers.go`)
   ```go
   func (s *Server) RegisterSettingsHandlers() error {
       // Use s.projectPath instead of homeDir
       settingsPath := filepath.Join(s.projectPath, "_bmad-output", ".autobmad")
       // ...
   }
   ```

**Phase 2: Migration Strategy**

For users upgrading from Epic 1 → Epic 1.5:

1. Check if old global config exists: `~/.autobmad/_bmad-output/.autobmad/config.json`
2. If yes, migrate to new location on first project open
3. Add warning log: "Migrated settings from global to project-local storage"
4. Leave old file in place (don't delete - user might have multiple versions)

**Phase 3: Testing**

1. Update all settings handler tests to use project paths
2. Add integration test for multi-project scenarios
3. Test migration from global → project-local
4. Verify project profiles work correctly per-project

#### Files to Modify

**Backend (Go):**
- ✏️ `apps/core/cmd/autobmad/main.go` - Add `--project-path` flag
- ✏️ `apps/core/internal/server/server.go` - Add projectPath field
- ✏️ `apps/core/internal/server/settings_handlers.go` - Use project path
- ✏️ `apps/core/internal/server/settings_handlers_test.go` - Update tests
- ✏️ `apps/core/internal/state/manager.go` - Add migration logic (optional)

**Frontend (TypeScript):**
- ✏️ `apps/desktop/src/main/index.ts` - Pass project path to backend
- ✏️ `apps/desktop/src/main/backend.ts` - Accept project path parameter

**Tests:**
- ✏️ `apps/desktop/src/main/backend.test.ts` - Mock project path
- ➕ `apps/core/internal/server/settings_migration_test.go` - New migration tests

#### Acceptance Criteria

- [ ] Settings stored at `<project>/_bmad-output/.autobmad/config.json`
- [ ] Different projects can have different settings
- [ ] Project profiles map works correctly (per-project)
- [ ] Migration from global settings (if exists) works seamlessly
- [ ] All existing tests pass with updated paths
- [ ] New integration test proves multi-project isolation

#### Risks

- **High Complexity:** Requires changes across 3 layers (Electron → Go → File System)
- **Breaking Change:** Users upgrading will need migration (mitigated by auto-migration)
- **Testing Burden:** Need to test project selection flow (not yet fully implemented)

#### Dependencies

- **Epic 1 Complete:** ✅ (All 8 commits merged)
- **Project Selection:** ⚠️ Might need to ensure project path is available early in app lifecycle

---

### 🟠 Story 1.5.2: Fix Failing SettingsScreen Test [HIGH]

**Priority:** High  
**Effort:** 1 hour  
**Complexity:** Low

#### Problem Statement

**Failing Test:**
```
FAIL src/renderer/src/screens/SettingsScreen.test.tsx > updates settings when input changes

Expected: { maxRetries: 5 }
Received: { maxRetries: 0 }  // Test receives intermediate value
```

#### Root Cause

**File:** `apps/desktop/src/renderer/src/screens/SettingsScreen.tsx:129`
```typescript
onChange={(e) => handleChange('maxRetries', parseInt(e.target.value) || 0)}
```

**Test:** `SettingsScreen.test.tsx`
```typescript
await user.clear(input);  // Triggers onChange with "" → parseInt("") = NaN → 0
await user.type(input, '5');  // Should update to 5
```

The test uses `userEvent.clear()` which triggers an intermediate `onChange` event with empty string, causing `parseInt("") = NaN`, falling back to `0`.

#### Solution Options

**Option A: Fix Test (Quick)**
```typescript
// Replace clear + type with direct fireEvent
fireEvent.change(input, { target: { value: '5' } });
```

**Option B: Fix Implementation (Better UX)**
```typescript
// Debounce changes to avoid saving on every keystroke
const debouncedChange = useMemo(
  () => debounce(handleChange, 300),
  []
);

onChange={(e) => debouncedChange('maxRetries', parseInt(e.target.value) || 3)}
```

**Recommendation:** Use **Option A** for quick fix, then implement **Option B** for better UX (fewer disk writes).

#### Acceptance Criteria

- [ ] All 9 SettingsScreen tests pass consistently
- [ ] No race conditions or flaky tests
- [ ] (Optional) Debounce implemented to reduce disk I/O

---

### 🟠 Story 1.5.3: Add Input Validation to Settings [HIGH]

**Priority:** High (Security/Data Integrity)  
**Effort:** 2-3 hours  
**Complexity:** Medium

#### Problem Statement

**Current State:** Settings can be set to **dangerous or nonsensical values** with NO validation.

**Attack Vectors:**
1. `maxRetries: -1` or `999999999` → Resource exhaustion
2. `theme: "INVALID"` → Frontend expects `"light"|"dark"|"system"`
3. `retryDelay: -5000` → Negative timeout breaks logic
4. `lastProjectPath: "../../etc/passwd"` → Path injection

#### Solution

Add validation in `StateManager.Set()`:

```go
// apps/core/internal/state/manager.go

func (sm *StateManager) Set(updates map[string]interface{}) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    for key, value := range updates {
        switch key {
        case "maxRetries":
            val := toInt(value)
            if val < 0 || val > 10 {
                return fmt.Errorf("maxRetries must be 0-10, got %d", val)
            }
            sm.settings.MaxRetries = val
            
        case "retryDelay":
            val := toInt(value)
            if val < 0 || val > 60000 {
                return fmt.Errorf("retryDelay must be 0-60000ms, got %d", val)
            }
            sm.settings.RetryDelay = val
            
        case "theme":
            val := value.(string)
            if val != "light" && val != "dark" && val != "system" {
                return fmt.Errorf("theme must be light|dark|system, got %s", val)
            }
            sm.settings.Theme = val
            
        case "lastProjectPath":
            val := value.(string)
            if val != "" {
                // Validate path (no traversal, must be absolute)
                if !filepath.IsAbs(val) {
                    return fmt.Errorf("lastProjectPath must be absolute, got %s", val)
                }
                // Reuse path validator from Story 1-8
                if err := server.ValidateProjectPath(val); err != nil {
                    return fmt.Errorf("invalid project path: %w", err)
                }
            }
            sm.settings.LastProjectPath = val
            
        // ... other fields with validation
        }
    }
    
    return sm.save()
}
```

#### Testing

Add tests in `manager_test.go`:
```go
func TestStateManagerInvalidInputs(t *testing.T) {
    tests := []struct {
        name    string
        updates map[string]interface{}
        wantErr string
    }{
        {"negative retries", map[string]interface{}{"maxRetries": -1}, "must be 0-10"},
        {"huge retries", map[string]interface{}{"maxRetries": 999999}, "must be 0-10"},
        {"invalid theme", map[string]interface{}{"theme": "RAINBOW"}, "must be light|dark|system"},
        {"negative delay", map[string]interface{}{"retryDelay": -5000}, "must be 0-60000"},
        {"relative path", map[string]interface{}{"lastProjectPath": "../etc"}, "must be absolute"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sm := createTestStateManager(t)
            err := sm.Set(tt.updates)
            if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
                t.Errorf("Expected error containing %q, got %v", tt.wantErr, err)
            }
        })
    }
}
```

#### Acceptance Criteria

- [ ] All numeric fields have min/max bounds validation
- [ ] Theme field has enum validation
- [ ] Path fields use existing path validator
- [ ] Invalid inputs return descriptive errors
- [ ] At least 10 validation tests added
- [ ] Frontend displays validation errors to user

---

### 🟠 Story 1.5.4: Add Integration Test for Settings Persistence [HIGH]

**Priority:** High  
**Effort:** 2 hours  
**Complexity:** Medium

#### Problem Statement

**AC #2 from Story 1-9:** "Settings persist across restarts"

**Current Test Coverage:**
- ✅ StateManager layer: `TestStateManagerSaveLoad` simulates restart
- ❌ JSON-RPC handler layer: No test for full server restart

**Gap:** A bug in handler initialization (e.g., forgetting to call `RegisterSettingsHandlers()`) would not be caught.

#### Solution

Add integration test in `settings_handlers_test.go`:

```go
func TestSettingsHandlersPersistenceAcrossRestart(t *testing.T) {
    // Phase 1: First server instance
    srv1 := createTestServer(t)
    defer srv1.Close()
    
    // Set some settings
    updates := map[string]interface{}{
        "maxRetries": 7,
        "theme": "dark",
        "desktopNotifications": false,
    }
    
    req1 := createJSONRPCRequest("settings.set", updates)
    resp1, err := srv1.HandleRequest(req1)
    require.NoError(t, err)
    
    // Phase 2: Simulate server restart (create new instance with same project path)
    srv1.Close()
    
    srv2 := createTestServer(t) // Uses same temp directory
    defer srv2.Close()
    
    // Phase 3: Retrieve settings from new instance
    req2 := createJSONRPCRequest("settings.get", nil)
    resp2, err := srv2.HandleRequest(req2)
    require.NoError(t, err)
    
    // Verify persistence
    var settings Settings
    json.Unmarshal(resp2.Result, &settings)
    
    assert.Equal(t, 7, settings.MaxRetries, "maxRetries should persist")
    assert.Equal(t, "dark", settings.Theme, "theme should persist")
    assert.False(t, settings.DesktopNotifications, "notifications should persist")
}
```

#### Acceptance Criteria

- [ ] Integration test proves settings persist across server restarts
- [ ] Test uses JSON-RPC handlers (not just StateManager)
- [ ] Test covers at least 5 different setting types
- [ ] Test runs in < 1 second

---

### 🟡 Story 1.5.5: Add Git Status UI Component [MEDIUM]

**Priority:** Medium  
**Effort:** 2-3 hours  
**Complexity:** Medium  
**Source:** Story 1-6 Code Review

#### Problem Statement

Story 1-6 (Git Detection) detects Git installation but **does not provide UI to show Git status** to the user.

**From Code Review:**
> "No visual feedback for Git repository status in the UI. Users can't see if they're in a Git repo, on which branch, or if there are uncommitted changes."

#### Solution

Add a **GitStatusBar** component to ProjectSelectScreen:

```tsx
// apps/desktop/src/renderer/src/components/GitStatusBar.tsx

export function GitStatusBar({ projectPath }: { projectPath: string }) {
  const [status, setStatus] = useState<GitStatus | null>(null);
  
  useEffect(() => {
    window.api.git.getStatus(projectPath).then(setStatus);
  }, [projectPath]);
  
  if (!status?.isGitRepo) return null;
  
  return (
    <div className="flex items-center gap-2 text-sm text-muted-foreground">
      <GitBranch className="h-4 w-4" />
      <span>{status.branch}</span>
      {status.hasChanges && (
        <Badge variant="warning">Uncommitted changes</Badge>
      )}
    </div>
  );
}
```

#### Acceptance Criteria

- [ ] UI shows Git branch name when in a Git repo
- [ ] UI shows warning if uncommitted changes exist
- [ ] UI hides gracefully when not in a Git repo
- [ ] Component tested with Git/non-Git projects

---

### 🟡 Story 1.5.6: Improve Test Coverage Metrics [MEDIUM]

**Priority:** Medium  
**Effort:** 1 hour  
**Complexity:** Low

#### Problem Statement

**Current State:** No coverage metrics reported in CI/CD or story completion notes.

**Story 1-9 Completion Notes:**
> "Comprehensive test coverage (12 unit tests, 100% pass rate)"

But no line/branch coverage % reported.

#### Solution

1. **Add Coverage to CI:**
   ```yaml
   # .github/workflows/test.yml
   - name: Run Go tests with coverage
     run: cd apps/core && go test -coverprofile=coverage.out ./...
   
   - name: Upload coverage to Codecov
     uses: codecov/codecov-action@v3
     with:
       file: ./apps/core/coverage.out
   ```

2. **Add Frontend Coverage:**
   ```json
   // apps/desktop/package.json
   "scripts": {
     "test:coverage": "vitest run --coverage"
   }
   ```

3. **Add Coverage Badge to README:**
   ```markdown
   [![Coverage](https://codecov.io/gh/yourorg/auto-bmad/branch/main/graph/badge.svg)](https://codecov.io/gh/yourorg/auto-bmad)
   ```

#### Acceptance Criteria

- [ ] Go tests report coverage % in CI
- [ ] Frontend tests report coverage % in CI
- [ ] Coverage > 80% for all packages
- [ ] Coverage badge added to README

---

### 🟡 Story 1.5.7: Add Concurrency Test for StateManager [MEDIUM]

**Priority:** Medium  
**Effort:** 1 hour  
**Complexity:** Low

#### Problem Statement

StateManager uses `sync.RWMutex` for thread safety, but **no test validates concurrent access**.

**Risk:** Race conditions, deadlocks, data corruption.

#### Solution

```go
func TestStateManagerConcurrency(t *testing.T) {
    sm := createTestStateManager(t)
    
    // Launch 10 goroutines writing simultaneously
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(val int) {
            defer wg.Done()
            sm.Set(map[string]interface{}{"maxRetries": val})
        }(i)
    }
    
    // Launch 10 goroutines reading simultaneously
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            _ = sm.Get()
        }()
    }
    
    wg.Wait()
    
    // Verify no panic, no corruption
    settings := sm.Get()
    if settings.MaxRetries < 0 || settings.MaxRetries > 10 {
        t.Errorf("Concurrent access corrupted settings")
    }
}
```

**Run with race detector:**
```bash
go test -race ./internal/state/
```

#### Acceptance Criteria

- [ ] Test launches 20+ concurrent goroutines
- [ ] Test passes with `-race` flag (no data races)
- [ ] No deadlocks or panics

---

### 🟡 Story 1.5.8: Wire Up Project Profiles to Project Opening [MEDIUM]

**Priority:** Medium  
**Effort:** 2 hours  
**Complexity:** Medium

#### Problem Statement

**Story 1-9 AC #2:** "last-used OpenCode profile **per project** is restored"

**Current State:**
- ✅ Settings has `ProjectProfiles map[string]string`
- ✅ Can persist profiles
- ❌ **Not wired up to project detection flow**

When a user opens a project, the last-used profile is **not** restored.

#### Solution

1. **Add RPC Handler:**
   ```go
   // apps/core/internal/server/project_handlers.go
   
   func (s *Server) handleProjectGetLastProfile(params json.RawMessage) (interface{}, error) {
       var req struct {
           Path string `json:"path"`
       }
       if err := json.Unmarshal(params, &req); err != nil {
           return nil, err
       }
       
       settings := s.stateManager.Get()
       profile := settings.ProjectProfiles[req.Path]
       
       return map[string]string{"profile": profile}, nil
   }
   ```

2. **Frontend Integration:**
   ```typescript
   // In ProjectSelectScreen.tsx
   
   const openProject = async (path: string) => {
       const { profile } = await window.api.project.getLastProfile(path);
       if (profile) {
           setSelectedProfile(profile);
       }
       // ... continue with project opening
   };
   ```

3. **Save Profile on Change:**
   ```typescript
   const handleProfileChange = async (profile: string) => {
       await window.api.settings.set({
           projectProfiles: { [currentProject]: profile }
       });
   };
   ```

#### Acceptance Criteria

- [ ] Last-used profile restored when opening a project
- [ ] Profile changes saved to settings
- [ ] Different projects remember different profiles
- [ ] Works with multiple projects in same session

---

### 🟢 Story 1.5.9: Add Field-Level Error Handling in SettingsScreen [LOW]

**Priority:** Low  
**Effort:** 1 hour  
**Complexity:** Low

#### Problem Statement

**Current:** Errors displayed globally - user doesn't know which field is invalid.

**Example:**
```
❌ Error: Invalid settings  (Which field? What's wrong?)
```

#### Solution

```tsx
const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

const handleChange = async (key: string, value: unknown) => {
  try {
    await window.api.settings.set({ [key]: value });
    setFieldErrors(prev => ({ ...prev, [key]: '' })); // Clear error
  } catch (err) {
    setFieldErrors(prev => ({ 
      ...prev, 
      [key]: err.message // "maxRetries must be 0-10"
    }));
  }
};

// In render:
<Input
  error={fieldErrors.maxRetries}
  helperText={fieldErrors.maxRetries}
/>
```

#### Acceptance Criteria

- [ ] Errors displayed next to the offending field
- [ ] Error messages are descriptive (not generic)
- [ ] Errors clear when field is fixed
- [ ] Tested with validation from Story 1.5.3

---

### 🟢 Story 1.5.10: Add Debouncing to Settings Inputs [LOW]

**Priority:** Low (UX improvement)  
**Effort:** 1 hour  
**Complexity:** Low

#### Problem Statement

**Current:** Settings save on **every keystroke** → Excessive disk I/O.

**User types "500" in retry delay field:**
```
onChange: "5"   → Save to disk
onChange: "50"  → Save to disk
onChange: "500" → Save to disk  (3 writes total)
```

#### Solution

```tsx
import { useDebouncedCallback } from 'use-debounce';

const debouncedChange = useDebouncedCallback(
  (key: string, value: unknown) => {
    window.api.settings.set({ [key]: value });
  },
  300 // Wait 300ms after last keystroke
);

// In input:
onChange={(e) => {
  const value = parseInt(e.target.value);
  setLocalValue(value); // Update UI immediately
  debouncedChange('maxRetries', value); // Save after 300ms
}}
```

#### Acceptance Criteria

- [ ] Changes debounced by 300ms
- [ ] UI updates immediately (optimistic)
- [ ] Only 1 disk write per "burst" of changes
- [ ] Tested with rapid typing

---

## Summary of Technical Debt

### By Priority

| Priority | Count | Effort (hours) |
|----------|-------|----------------|
| 🔴 Critical | 1 | 4-8 |
| 🟠 High | 4 | 7-9 |
| 🟡 Medium | 4 | 7-9 |
| 🟢 Low | 2 | 2 |
| **TOTAL** | **11** | **20-28** |

### By Category

| Category | Stories | Examples |
|----------|---------|----------|
| Architecture Compliance | 1 | Settings path, project-local config |
| Test Coverage | 4 | Integration tests, concurrency tests, coverage metrics |
| Input Validation | 1 | Bounds checking, enum validation |
| UX Improvements | 3 | Field errors, debouncing, Git status |
| Code Quality | 2 | Coverage reporting, test fixes |

---

## Execution Strategy

### Phase 1: Critical Fixes (Before Epic 2.1)
- ✅ Story 1.5.1: Settings path architecture
- ✅ Story 1.5.2: Failing test fix

**Reason:** These are architectural violations that should be fixed before building more features.

### Phase 2: High Priority (During Epic 2)
- Story 1.5.3: Input validation
- Story 1.5.4: Integration tests
- Story 1.5.5: Git status UI
- Story 1.5.6: Coverage metrics

**Reason:** These improve quality and can be done in parallel with Epic 2 feature work.

### Phase 3: Medium/Low Priority (Post-Epic 2)
- Stories 1.5.7 - 1.5.10: UX improvements, polishing

**Reason:** Nice-to-have improvements that don't block core functionality.

---

## Dependencies

### Blocking Epic 2?
**NO** - All critical security issues are fixed. Epic 2 can proceed.

### What Blocks This Epic?
- **Epic 1 Complete:** ✅ (8 commits merged)
- **Epic 2 Complete:** ⏳ (for full integration testing)

---

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Settings path refactor breaks existing users | Medium | High | Add auto-migration from global → project-local |
| Complex cross-layer changes introduce bugs | Medium | Medium | Comprehensive integration tests |
| Delayed indefinitely (technical debt grows) | High | Medium | **Assign to Epic 1.5 sprint explicitly** |

---

## References

- **Epic 1 Code Review:** All story files in `_bmad-output/implementation-artifacts/`
- **Architecture Spec:** `_bmad/architecture.md`
- **Story 1-9 Review:** Most critical issues documented in `1-9-settings-persistence.md`

---

**Created:** 2026-01-23  
**Next Review:** After Epic 2 completion  
**Owner:** TBD  
