---
project_name: 'auto-bmad'
user_name: 'Hafiz'
date: '2026-01-21'
status: 'complete'
---

# Project Context for AI Agents

_Critical rules and patterns for implementing Auto-BMAD. Focus on unobvious details._

---

## Technology Stack & Versions

| Layer | Technology | Version | Notes |
|-------|------------|---------|-------|
| Desktop Shell | Electron | Latest | Use electron-vite template |
| UI Framework | React | 18+ | Functional components only |
| Language (FE) | TypeScript | 5.x | Strict mode enabled |
| Build | Vite | 5.x | Via electron-vite |
| Styling | Tailwind CSS | 3.x | Via shadcn/ui |
| Components | shadcn/ui | Latest | Radix primitives |
| State | Zustand | Latest | Feature-scoped stores |
| Backend | Golang | 1.21+ | Standard layout |
| Logging | slog | Built-in | Structured JSON |

---

## Critical Implementation Rules

### Golang Rules

**Package & Naming:**
- Packages: lowercase, single word (`journey`, `opencode`, not `journeyManager`)
- Files: snake_case (`journey_state.go`, not `journeyState.go`)
- Exported: PascalCase (`StartJourney`, `JourneyState`)
- Unexported: camelCase (`parseOutput`, `currentStep`)
- Interfaces: NO `I` prefix (`Executor`, not `IExecutor`)

**Error Handling:**
- Always wrap errors with context: `fmt.Errorf("starting journey: %w", err)`
- Use custom error types in `internal/common/errors.go`
- Map errors to JSON-RPC error codes before returning

**Logging:**
```go
// ✅ Correct - structured with context
slog.Info("journey started", "journeyId", j.ID, "workflow", j.Workflow)

// ❌ Wrong - unstructured
log.Printf("Journey %s started", j.ID)
```

**Concurrency:**
- Use `context.Context` for cancellation
- Prefer channels over shared memory
- Always handle goroutine cleanup

### TypeScript/React Rules

**File Naming:**
- Components: PascalCase (`JourneyCard.tsx`)
- Hooks: camelCase with `use` prefix (`useJourney.ts`)
- Utils: camelCase (`formatTime.ts`)
- Types: camelCase or `types.ts` (`journey.types.ts`)

**Component Structure:**
```typescript
// ✅ Correct - function declaration, named export
export function JourneyCard({ journey }: JourneyCardProps) {
  return <div>...</div>;
}

// ❌ Wrong - arrow function default export
export default ({ journey }) => <div>...</div>;
```

**Zustand Stores:**
```typescript
// ✅ Correct - feature-scoped, typed
export const useJourneyStore = create<JourneyStore>((set, get) => ({
  currentJourney: null,
  startJourney: (workflow) => set({ ... }),
}));

// ❌ Wrong - global store, untyped
export const useStore = create((set) => ({ ... }));
```

**State Updates:**
```typescript
// ✅ Correct - immutable update
set((state) => ({
  steps: state.steps.map((s, i) => i === index ? { ...s, status } : s)
}));

// ❌ Wrong - direct mutation
set((state) => { state.steps[index].status = status; });
```

### JSON-RPC Protocol Rules

**Method Naming:** `resource.action`
```typescript
// ✅ Correct
'journey.start', 'journey.pause', 'step.retry', 'checkpoint.create'

// ❌ Wrong
'startJourney', 'JourneyStart', 'journey/start'
```

**Event Naming:** `resource.event`
```typescript
// ✅ Correct
'journey.started', 'step.completed', 'yellowFlag.raised'

// ❌ Wrong  
'onJourneyStart', 'STEP_COMPLETED'
```

**Field Naming:** camelCase in JSON
```json
// ✅ Correct
{ "journeyId": "j-123", "currentStep": 3, "isRetryable": true }

// ❌ Wrong
{ "journey_id": "j-123", "CurrentStep": 3 }
```

**Error Response:** JSON-RPC 2.0 format
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32000,
    "message": "Journey not found",
    "data": { "journeyId": "j-invalid" }
  }
}
```

### IPC Security Rules

**Preload Script:**
```typescript
// ✅ Correct - explicit API exposure
contextBridge.exposeInMainWorld('api', {
  journey: {
    start: (workflow: string) => ipcRenderer.invoke('journey.start', workflow),
  }
});

// ❌ Wrong - exposing ipcRenderer directly
contextBridge.exposeInMainWorld('ipc', ipcRenderer);
```

**Renderer Access:**
```typescript
// ✅ Correct - use exposed API
window.api.journey.start('create-prd');

// ❌ Wrong - direct require (blocked by contextIsolation)
require('electron').ipcRenderer.send(...);
```

---

## Testing Rules

**File Location:** Co-located with source
```
journey/
├── orchestrator.go
├── orchestrator_test.go    ← Go tests
├── JourneyCard.tsx
└── JourneyCard.test.tsx    ← React tests
```

**Test Naming:**
- Go: `func TestOrchestratorStart(t *testing.T)`
- TS: `describe('JourneyCard', () => { it('renders status', ...) })`

**What to Test:**
- Go: Unit test each public function, integration test JSON-RPC handlers
- React: Component rendering, user interactions, store updates
- Skip: Electron main process (integration test manually)

---

## Code Quality Rules

**Before Every Commit:**
```bash
# Go
go fmt ./... && golangci-lint run

# TypeScript
npm run lint && npm run format
```

**Import Order (TypeScript):**
```typescript
// 1. React/external
import { useState } from 'react';
import { create } from 'zustand';

// 2. Internal absolute
import { Button } from '@/components/ui/button';

// 3. Relative
import { useJourney } from './hooks/useJourney';
import type { Journey } from './types';
```

**Component Organization:** By feature, not type
```
// ✅ Correct
features/journey/components/JourneyCard.tsx
features/journey/hooks/useJourney.ts

// ❌ Wrong
components/JourneyCard.tsx
hooks/useJourney.ts
```

---

## Data Format Rules

| Data Type | Format | Example |
|-----------|--------|---------|
| Dates | ISO 8601 | `2026-01-21T10:30:00Z` |
| IDs | String with prefix | `j-20260121-001`, `s-3` |
| Booleans | `true`/`false` | NOT `1`/`0` |
| Null | Explicit `null` | NOT omitted |
| Empty array | `[]` | NOT `null` |
| Durations | Milliseconds (int) | `5000` for 5 seconds |

---

## Critical Don't-Miss Rules

### NEVER Do These:

1. **Don't store state in React** — Filesystem is source of truth, React is view only
2. **Don't use `any` in TypeScript** — Use proper types or `unknown`
3. **Don't log to stdout in Go** — Stdout is reserved for JSON-RPC, use stderr or file
4. **Don't skip error wrapping** — Always add context to errors
5. **Don't mutate Zustand state directly** — Use immutable updates
6. **Don't expose ipcRenderer** — Use contextBridge with explicit API
7. **Don't use snake_case in JSON** — Always camelCase across IPC boundary

### Edge Cases to Handle:

1. **OpenCode crash mid-step** — Capture partial output, allow retry
2. **Git conflict on checkpoint** — Stash, commit, restore
3. **Filesystem permission denied** — Show user-friendly error with fix suggestion
4. **Journey interrupted (app close)** — Resume from last checkpoint on restart
5. **Yellow flag timeout** — Auto-pause after 5 minutes if no response

### Security Considerations:

1. **No API keys** — Auto-BMAD doesn't store credentials
2. **IPC isolation** — contextIsolation: true, nodeIntegration: false
3. **Validate all IPC input** — Don't trust renderer data
4. **Sanitize file paths** — Prevent path traversal in project detection
