---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8]
inputDocuments:
  - _bmad-output/planning-artifacts/prd.md
  - _bmad-output/planning-artifacts/ux-design-specification.md
  - _bmad-output/planning-artifacts/product-brief-auto-bmad-2026-01-20.md
  - _bmad-output/planning-artifacts/prd-validation-report-v1.1.md
workflowType: 'architecture'
project_name: 'auto-bmad'
user_name: 'Hafiz'
date: '2026-01-21'
status: 'complete'
completedAt: '2026-01-21'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**

Auto-BMAD has 52 functional requirements organized into 6 categories:

| Category | Count | Architectural Significance |
|----------|-------|---------------------------|
| Journey Management (FR1-FR10) | 10 | Core orchestration engine, state machine design |
| OpenCode Integration (FR11-FR19) | 9 | Process spawning, IPC, output capture |
| Execution & Retry (FR20-FR27) | 8 | Workflow sequencing, feedback accumulation |
| Failure & Reporting (FR28-FR35) | 8 | Error handling, honest failure protocol |
| Dashboard & Visualization (FR36-FR43) | 8 | Real-time UI updates, notification system |
| Project & Configuration (FR44-FR52) | 9 | Filesystem detection, settings persistence |

**Non-Functional Requirements:**

| Category | Key Constraints | Architectural Impact |
|----------|-----------------|---------------------|
| Performance | Startup < 5s, UI < 100ms, Memory < 500MB | Lazy loading, efficient IPC, minimal footprint |
| Reliability | Zero data loss, continuous checkpoints | Robust state persistence, Git integration |
| Integration | OpenCode CLI, Git 2.0+, BMAD 6.0+ | Dependency detection, version compatibility |
| Security | No API keys, code signing, IPC isolation | Electron security best practices |
| Usability | Keyboard nav, WCAG AA | Accessible component architecture |

**Scale & Complexity:**

- Primary domain: Desktop Application (Electron + Golang)
- Complexity level: Medium-High
- Estimated architectural components: 15-20

### Technical Constraints & Dependencies

**Hard Constraints:**

| Constraint | Rationale |
|------------|-----------|
| Electron frontend | Desktop app requirement, cross-platform UI |
| Golang backend | Performance, single binary distribution, process management |
| OpenCode CLI integration | AI execution backbone (not embedded) |
| Git for checkpoints | State safety, rollback capability |
| Linux + macOS (MVP) | Primary target platforms |

**External Dependencies:**

| Dependency | Minimum Version | Detection |
|------------|-----------------|-----------|
| OpenCode CLI | v0.1.0+ | `opencode --version` |
| Git | 2.0+ | `git --version` |
| BMAD | 6.0.0+ | `_bmad/_config/manifest.yaml` |

### Cross-Cutting Concerns Identified

| Concern | Components Affected | Decision Required |
|---------|---------------------|-------------------|
| **IPC Protocol** | Electron ↔ Golang | Communication mechanism (embedded binary, HTTP, WebSocket) |
| **State Management** | UI, Backend, Filesystem | Source of truth, sync strategy |
| **Error Handling** | All layers | Propagation pattern, user-facing messages |
| **Process Lifecycle** | Backend, OpenCode | Spawn, monitor, terminate, crash recovery |
| **Checkpointing** | Backend, Git | Frequency, content, rollback strategy |
| **Logging** | All layers | Structured logging, audit trails |
| **Notifications** | UI, OS | Desktop notifications, in-app feedback |

## Starter Template Evaluation

### Technical Preferences

| Preference | Choice | Rationale |
|------------|--------|-----------|
| TypeScript | ✅ Yes | Type safety for Electron/React frontend |
| Golang Skill Level | Expert | Standard Go layout appropriate, idiomatic patterns |
| Monorepo Structure | ✅ Yes | Single repository with `/apps` structure for cohesion |

### Starter Template Selection

**Selected:** `electron-vite` (create @electron-vite)

| Option Evaluated | Verdict | Rationale |
|------------------|---------|-----------|
| **electron-vite** | ✅ SELECTED | Fast DX (Vite HMR), clean structure, easy shadcn/ui integration |
| electron-react-boilerplate | ❌ Rejected | Webpack-based (slower), too opinionated |
| Electron Forge + Vite | ❌ Rejected | Less React-focused, Forge-specific patterns |

### Monorepo Structure

```
auto-bmad/
├── apps/
│   ├── desktop/              # Electron + React + TypeScript + Vite
│   │   ├── src/
│   │   │   ├── main/         # Electron main process
│   │   │   ├── preload/      # IPC bridge (contextIsolation)
│   │   │   └── renderer/     # React UI (shadcn/ui + Zustand)
│   │   └── package.json
│   │
│   └── core/                 # Golang backend
│       ├── cmd/autobmad/     # Main entry point
│       ├── internal/         # Private packages
│       │   ├── journey/      # Journey orchestration engine
│       │   ├── opencode/     # OpenCode CLI integration
│       │   ├── checkpoint/   # Git operations
│       │   ├── state/        # State management
│       │   └── server/       # IPC server
│       └── go.mod
│
├── packages/
│   └── shared-types/         # Shared TypeScript types (optional)
│
├── scripts/                  # Build, dev, package scripts
└── package.json              # Root workspace config (pnpm/npm)
```

### Decisions Made by Starter

| Decision | Technology | Version |
|----------|------------|---------|
| Language (Frontend) | TypeScript | 5.x (strict mode) |
| Styling | Tailwind CSS | 3.x (via shadcn/ui) |
| Build Tool | Vite | 5.x |
| Packaging | electron-builder | Latest |
| Golang Layout | Standard (cmd/internal/pkg) | N/A |

### Initialization Commands

```bash
# 1. Create monorepo structure
mkdir -p auto-bmad/{apps,packages,scripts}
cd auto-bmad

# 2. Initialize Electron app with electron-vite
npm create @electron-vite/app@latest apps/desktop -- --template react-ts

# 3. Initialize Golang backend (expert layout)
mkdir -p apps/core/cmd/autobmad apps/core/internal
cd apps/core
go mod init github.com/fairyhunter13/auto-bmad/apps/core

# 4. Add shadcn/ui to desktop app
cd ../desktop
npx shadcn-ui@latest init
```

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**
- IPC Protocol between Electron and Golang
- Backend embedding and distribution model
- State management source of truth
- OpenCode CLI integration pattern
- Checkpoint strategy for zero data loss

**Deferred Decisions (Post-MVP):**
- Windows platform support
- Plugin/extension architecture
- Multi-journey parallelism

### Communication Architecture

**IPC Protocol: stdio/JSON-RPC**

| Aspect | Decision |
|--------|----------|
| Protocol | JSON-RPC 2.0 over stdin/stdout |
| Direction | Bidirectional (request/response + event streaming) |
| Serialization | JSON |
| Error Handling | JSON-RPC error codes + custom application codes |

**Message Flow:**
```
Electron Main Process
    │
    ├── spawn(golang-binary)
    │
    ├── stdin  → JSON-RPC requests  → Golang
    └── stdout ← JSON-RPC responses ← Golang
              ← Event stream (newline-delimited JSON)
```

**Backend Embedding: Embedded Binary**

| Aspect | Decision |
|--------|----------|
| Packaging | Golang binary in Electron `resources/bin/` |
| Platform Binaries | `autobmad-linux`, `autobmad-darwin` |
| Launch | Electron main process spawns on app start |
| Lifecycle | Tied to Electron app lifecycle |

### State Architecture

**Source of Truth: Filesystem (Artifacts)**

| Aspect | Decision |
|--------|----------|
| Primary State | BMAD artifacts in `_bmad-output/` |
| Journey State | `_bmad-output/.autobmad/journey-state.json` |
| Configuration | `_bmad-output/.autobmad/config.json` |
| Crash Recovery | Read filesystem state on restart |

**State Flow:**
```
Filesystem (artifacts)
    ↑ write
    │
Golang Backend (orchestrator)
    ↑ JSON-RPC
    │
Electron Main (IPC bridge)
    ↑ contextBridge
    │
React UI (Zustand store) ← view only, not source of truth
```

### OpenCode Integration

**Process Model: One-shot per Step**

| Aspect | Decision |
|--------|----------|
| Invocation | New process per workflow step |
| Command | `opencode -p "{prompt}" --non-interactive` |
| Output Capture | Stream stdout/stderr in real-time |
| Completion | Exit code + output parsing |
| Timeout | Configurable per step type (default: 5 min) |

**Process Lifecycle:**
```
Step Start
    ↓
Spawn OpenCode (with prompt + context)
    ↓
Stream output → Parse → Update UI
    ↓
Wait for exit
    ↓
Exit Code 0? → Step Complete
Exit Code ≠ 0? → Step Failed → Capture error → Allow retry
```

### Checkpoint Strategy

**Approach: Event-based Git Commits**

| Event | Action |
|-------|--------|
| Step Completion (success) | Commit all changed artifacts |
| Step Completion (failure) | Commit with failure state |
| User-initiated Pause | Commit current state |
| Yellow Flag Raised | Commit before prompting user |
| Journey Completion | Final commit + optional tag |

**Git Strategy:**

| Aspect | Decision |
|--------|----------|
| Branch | `autobmad/journey-{timestamp}` |
| Commit Message | `[AutoBMAD] Step {n}: {name} - {status}` |
| Rollback | `git reset --hard {checkpoint-sha}` |
| Cleanup | Merge to main or delete branch on journey complete |

### Decision Impact Analysis

**Implementation Sequence:**
1. Golang binary scaffold with JSON-RPC server
2. Electron spawn + IPC bridge
3. State management (filesystem read/write)
4. OpenCode process spawning
5. Checkpoint integration (Git operations)
6. UI integration (Zustand ← JSON-RPC events)

**Cross-Component Dependencies:**
```
┌─────────────────────────────────────────────────────────┐
│                     Electron App                        │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │ React UI    │◄──►│ Main Process│◄──►│ Preload     │ │
│  │ (Zustand)   │    │ (spawn)     │    │ (bridge)    │ │
│  └─────────────┘    └──────┬──────┘    └─────────────┘ │
└────────────────────────────┼────────────────────────────┘
                             │ stdio/JSON-RPC
┌────────────────────────────┼────────────────────────────┐
│                     Golang Backend                      │
│  ┌─────────────┐    ┌──────┴──────┐    ┌─────────────┐ │
│  │ Journey     │◄──►│ JSON-RPC    │◄──►│ State       │ │
│  │ Orchestrator│    │ Server      │    │ Manager     │ │
│  └──────┬──────┘    └─────────────┘    └──────┬──────┘ │
│         │                                      │        │
│  ┌──────┴──────┐                       ┌──────┴──────┐ │
│  │ OpenCode    │                       │ Checkpoint  │ │
│  │ Executor    │                       │ (Git)       │ │
│  └──────┬──────┘                       └──────┬──────┘ │
└─────────┼──────────────────────────────────────┼───────┘
          │ spawn                                │ git
          ▼                                      ▼
    ┌───────────┐                         ┌───────────┐
    │ OpenCode  │                         │ Filesystem│
    │ CLI       │                         │ (artifacts)│
    └───────────┘                         └───────────┘
```

## Implementation Patterns & Consistency Rules

### Pattern Categories Defined

**Critical Conflict Points Addressed:** 8 areas where AI agents could diverge

| Area | Resolution |
|------|------------|
| JSON-RPC method naming | `resource.action` convention |
| JSON field naming | camelCase throughout |
| Error structure | JSON-RPC 2.0 compliant |
| File naming | Language-idiomatic conventions |
| Test organization | Co-located tests |
| Event naming | `resource.event` convention |
| State management | Feature-based Zustand stores |
| Component organization | By feature |

### JSON-RPC Protocol Patterns

**Method Naming Convention:** `resource.action`

| Resource | Methods |
|----------|---------|
| `journey` | `journey.start`, `journey.pause`, `journey.resume`, `journey.abort`, `journey.getState` |
| `step` | `step.retry`, `step.skip`, `step.getStatus` |
| `opencode` | `opencode.execute`, `opencode.cancel` |
| `checkpoint` | `checkpoint.create`, `checkpoint.rollback`, `checkpoint.list` |
| `config` | `config.get`, `config.set` |
| `project` | `project.detect`, `project.validate` |

**Event Naming Convention:** `resource.event`

| Event | Payload |
|-------|---------|
| `journey.started` | `{ journeyId, workflow, startedAt }` |
| `journey.completed` | `{ journeyId, status, completedAt }` |
| `step.started` | `{ journeyId, stepIndex, stepName }` |
| `step.completed` | `{ journeyId, stepIndex, status, output }` |
| `step.failed` | `{ journeyId, stepIndex, error, retryable }` |
| `opencode.output` | `{ journeyId, stepIndex, chunk, stream }` |
| `yellowFlag.raised` | `{ journeyId, stepIndex, reason, options }` |
| `checkpoint.created` | `{ journeyId, commitSha, message }` |

**JSON Field Naming:** camelCase

```json
{
  "journeyId": "j-20260121-001",
  "currentStep": 3,
  "stepName": "create-prd",
  "isRetryable": true,
  "createdAt": "2026-01-21T10:30:00Z"
}
```

### Error Handling Patterns

**Error Response Structure (JSON-RPC 2.0):**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32000,
    "message": "Journey not found",
    "data": {
      "journeyId": "j-invalid",
      "suggestion": "Check journey ID or start a new journey"
    }
  }
}
```

**Application Error Codes:**

| Code | Category | Example |
|------|----------|---------|
| -32000 to -32099 | Journey errors | Not found, invalid state |
| -32100 to -32199 | OpenCode errors | Spawn failed, timeout, crash |
| -32200 to -32299 | Checkpoint errors | Git failed, conflict |
| -32300 to -32399 | Validation errors | Invalid config, missing deps |
| -32600 to -32699 | JSON-RPC standard | Invalid request, method not found |

**Error Propagation Pattern:**

```
OpenCode exit code ≠ 0
    ↓
Golang: Capture stderr, wrap in AppError
    ↓
JSON-RPC: Send error response with code + context
    ↓
Electron: Parse error, update Zustand store
    ↓
React: Display user-friendly message + retry option
```

### Naming Conventions

**TypeScript/React:**

| Element | Convention | Example |
|---------|------------|---------|
| Component files | PascalCase | `JourneyCard.tsx` |
| Component exports | PascalCase | `export function JourneyCard()` |
| Hook files | camelCase | `useJourney.ts` |
| Hook exports | camelCase with `use` | `export function useJourney()` |
| Utility files | camelCase | `formatTime.ts` |
| Type files | camelCase or `types.ts` | `journey.types.ts` |
| Type exports | PascalCase | `export interface Journey` |
| Constants | SCREAMING_SNAKE | `export const MAX_RETRIES = 3` |
| Zustand stores | `use{Resource}Store` | `useJourneyStore` |

**Golang:**

| Element | Convention | Example |
|---------|------------|---------|
| Package names | lowercase, single word | `journey`, `opencode` |
| Files | snake_case | `journey_state.go` |
| Exported types | PascalCase | `type JourneyState struct` |
| Exported funcs | PascalCase | `func StartJourney()` |
| Unexported | camelCase | `func parseOutput()` |
| Constants (exported) | PascalCase | `const MaxRetries = 3` |
| Constants (unexported) | camelCase | `const defaultTimeout` |
| Interfaces | PascalCase, no `I` prefix | `type Executor interface` |

### Structure Patterns

**Test Organization:** Co-located

```
# TypeScript
src/features/journey/
├── JourneyCard.tsx
├── JourneyCard.test.tsx    ← co-located
├── useJourney.ts
└── useJourney.test.ts      ← co-located

# Golang
internal/journey/
├── state.go
├── state_test.go           ← co-located
├── orchestrator.go
└── orchestrator_test.go    ← co-located
```

**Component Organization:** By feature

```
src/
├── features/               # Feature modules
│   ├── journey/           # Journey feature
│   │   ├── components/    # Journey-specific components
│   │   ├── hooks/         # Journey-specific hooks
│   │   ├── store.ts       # Journey Zustand store
│   │   └── types.ts       # Journey types
│   ├── dashboard/         # Dashboard feature
│   └── settings/          # Settings feature
├── components/
│   ├── ui/                # shadcn/ui components
│   └── common/            # Shared app components
├── hooks/                  # Shared hooks
├── lib/                    # Utilities, IPC client
└── types/                  # Shared types
```

### Data Format Patterns

| Pattern | Standard |
|---------|----------|
| Dates | ISO 8601: `2026-01-21T10:30:00Z` |
| Booleans | `true` / `false` (not 0/1) |
| Null values | Explicit `null`, never omitted |
| Empty arrays | `[]` (not `null`) |
| IDs | String prefixed: `j-{timestamp}`, `s-{index}` |
| Durations | Milliseconds as integer |

### Logging Patterns

**Log Levels:**

| Level | Usage |
|-------|-------|
| `debug` | Detailed debugging (dev only) |
| `info` | Normal operations, state changes |
| `warn` | Recoverable issues, retries |
| `error` | Failures requiring attention |

**Golang Logging (slog):**

```go
slog.Info("journey started",
    "journeyId", j.ID,
    "workflow", j.Workflow,
)

slog.Error("opencode execution failed",
    "journeyId", j.ID,
    "stepIndex", step.Index,
    "error", err,
)
```

**TypeScript Logging:**

```typescript
console.info('[Journey]', 'Started', { journeyId, workflow });
console.error('[OpenCode]', 'Failed', { journeyId, stepIndex, error });
```

### State Management Patterns

**Zustand Store Structure:**

```typescript
// Feature-scoped stores
export const useJourneyStore = create<JourneyStore>((set, get) => ({
  // State
  currentJourney: null,
  steps: [],
  status: 'idle',
  
  // Actions (camelCase verbs)
  startJourney: (workflow) => { ... },
  pauseJourney: () => { ... },
  updateStep: (index, status) => { ... },
  
  // Computed (get prefix)
  getCurrentStep: () => get().steps[get().currentStepIndex],
}));
```

**State Update Pattern:** Immutable updates

```typescript
// ✅ Correct
set((state) => ({
  steps: state.steps.map((s, i) => 
    i === index ? { ...s, status: newStatus } : s
  )
}));

// ❌ Wrong - direct mutation
set((state) => {
  state.steps[index].status = newStatus; // NO!
});
```

### Enforcement Guidelines

**All AI Agents MUST:**

1. Follow JSON-RPC 2.0 specification for all IPC communication
2. Use camelCase for all JSON fields crossing the IPC boundary
3. Use `resource.action` naming for methods, `resource.event` for events
4. Co-locate test files with source files
5. Organize React components by feature, not by type
6. Use structured logging with consistent levels
7. Return ISO 8601 dates, explicit nulls, empty arrays (not null)

**Pattern Verification:**

- TypeScript: ESLint + Prettier enforces naming/formatting
- Golang: `go fmt` + `golangci-lint` enforces conventions
- JSON-RPC: Schema validation in IPC layer
- Pre-commit hooks: Run linters before commit

## Project Structure & Boundaries

### Requirements → Structure Mapping

| FR Category | Golang Package | React Feature | Key Files |
|-------------|----------------|---------------|-----------|
| Journey Management (FR1-10) | `internal/journey` | `features/journey` | `orchestrator.go`, `JourneyPanel.tsx` |
| OpenCode Integration (FR11-19) | `internal/opencode` | `features/journey` | `executor.go`, `OutputViewer.tsx` |
| Execution & Retry (FR20-27) | `internal/journey` | `features/journey` | `state_machine.go`, `StepControls.tsx` |
| Failure & Reporting (FR28-35) | `internal/journey` | `features/journey` | `error_handler.go`, `ErrorDisplay.tsx` |
| Dashboard & Viz (FR36-43) | `internal/server` | `features/dashboard` | `events.go`, `Dashboard.tsx` |
| Project & Config (FR44-52) | `internal/project` | `features/settings` | `detector.go`, `Settings.tsx` |

### Complete Project Directory Structure

```
auto-bmad/
├── .github/
│   └── workflows/
│       └── ci.yml                    # CI pipeline (lint, test, build)
│
├── apps/
│   ├── desktop/                      # Electron + React + TypeScript
│   │   ├── package.json
│   │   ├── electron.vite.config.ts
│   │   ├── tsconfig.json
│   │   ├── tsconfig.node.json
│   │   ├── tailwind.config.js
│   │   ├── postcss.config.js
│   │   ├── components.json           # shadcn/ui config
│   │   │
│   │   ├── src/
│   │   │   ├── main/                 # Electron main process
│   │   │   │   ├── index.ts          # App entry, window management
│   │   │   │   ├── backend.ts        # Spawn & manage Golang binary
│   │   │   │   ├── ipc.ts            # JSON-RPC client over stdio
│   │   │   │   └── menu.ts           # App menu configuration
│   │   │   │
│   │   │   ├── preload/              # Context bridge
│   │   │   │   ├── index.ts          # Expose IPC to renderer
│   │   │   │   └── types.ts          # Preload API types
│   │   │   │
│   │   │   └── renderer/             # React UI
│   │   │       ├── index.html
│   │   │       ├── main.tsx          # React entry
│   │   │       ├── App.tsx           # Root component
│   │   │       │
│   │   │       ├── features/
│   │   │       │   ├── journey/      # Journey management
│   │   │       │   │   ├── components/
│   │   │       │   │   │   ├── JourneyPanel.tsx
│   │   │       │   │   │   ├── JourneyPanel.test.tsx
│   │   │       │   │   │   ├── StepList.tsx
│   │   │       │   │   │   ├── StepCard.tsx
│   │   │       │   │   │   ├── OutputViewer.tsx
│   │   │       │   │   │   ├── YellowFlagModal.tsx
│   │   │       │   │   │   └── StepControls.tsx
│   │   │       │   │   ├── hooks/
│   │   │       │   │   │   ├── useJourney.ts
│   │   │       │   │   │   └── useJourney.test.ts
│   │   │       │   │   ├── store.ts           # useJourneyStore
│   │   │       │   │   └── types.ts
│   │   │       │   │
│   │   │       │   ├── dashboard/    # Main dashboard
│   │   │       │   │   ├── components/
│   │   │       │   │   │   ├── Dashboard.tsx
│   │   │       │   │   │   ├── StatusIndicator.tsx
│   │   │       │   │   │   ├── WorkflowSelector.tsx
│   │   │       │   │   │   └── RecentJourneys.tsx
│   │   │       │   │   ├── hooks/
│   │   │       │   │   │   └── useDashboard.ts
│   │   │       │   │   ├── store.ts           # useDashboardStore
│   │   │       │   │   └── types.ts
│   │   │       │   │
│   │   │       │   └── settings/     # Configuration
│   │   │       │       ├── components/
│   │   │       │       │   ├── Settings.tsx
│   │   │       │       │   ├── ProjectConfig.tsx
│   │   │       │       │   └── DependencyCheck.tsx
│   │   │       │       ├── hooks/
│   │   │       │       │   └── useSettings.ts
│   │   │       │       ├── store.ts           # useSettingsStore
│   │   │       │       └── types.ts
│   │   │       │
│   │   │       ├── components/
│   │   │       │   ├── ui/           # shadcn/ui components
│   │   │       │   │   ├── button.tsx
│   │   │       │   │   ├── card.tsx
│   │   │       │   │   ├── dialog.tsx
│   │   │       │   │   ├── progress.tsx
│   │   │       │   │   └── toast.tsx
│   │   │       │   │
│   │   │       │   └── common/       # Shared app components
│   │   │       │       ├── Layout.tsx
│   │   │       │       ├── Header.tsx
│   │   │       │       ├── Sidebar.tsx
│   │   │       │       └── ErrorBoundary.tsx
│   │   │       │
│   │   │       ├── hooks/            # Shared hooks
│   │   │       │   ├── useIpc.ts     # JSON-RPC hook
│   │   │       │   └── useNotification.ts
│   │   │       │
│   │   │       ├── lib/              # Utilities
│   │   │       │   ├── ipc-client.ts # JSON-RPC client
│   │   │       │   ├── utils.ts      # shadcn/ui cn()
│   │   │       │   └── constants.ts
│   │   │       │
│   │   │       ├── types/            # Shared types
│   │   │       │   ├── journey.ts
│   │   │       │   ├── ipc.ts
│   │   │       │   └── index.ts
│   │   │       │
│   │   │       └── styles/
│   │   │           └── globals.css   # Tailwind imports
│   │   │
│   │   ├── resources/                # Static resources
│   │   │   ├── bin/                  # Golang binaries (gitignored, built)
│   │   │   │   ├── .gitkeep
│   │   │   │   ├── autobmad-linux    # Linux binary
│   │   │   │   └── autobmad-darwin   # macOS binary
│   │   │   └── icons/
│   │   │       └── icon.png
│   │   │
│   │   └── electron-builder.yml      # Build config
│   │
│   └── core/                         # Golang backend
│       ├── go.mod
│       ├── go.sum
│       ├── Makefile                  # Build commands
│       │
│       ├── cmd/
│       │   └── autobmad/
│       │       └── main.go           # Entry point, starts JSON-RPC server
│       │
│       └── internal/
│           ├── server/               # JSON-RPC server
│           │   ├── server.go         # stdio server, request routing
│           │   ├── server_test.go
│           │   ├── handlers.go       # Method handlers
│           │   ├── events.go         # Event emission
│           │   └── types.go          # Request/response types
│           │
│           ├── journey/              # Journey orchestration
│           │   ├── orchestrator.go   # Main orchestration logic
│           │   ├── orchestrator_test.go
│           │   ├── state_machine.go  # Journey state transitions
│           │   ├── state_machine_test.go
│           │   ├── step.go           # Step execution
│           │   ├── step_test.go
│           │   └── types.go          # Journey, Step structs
│           │
│           ├── opencode/             # OpenCode CLI integration
│           │   ├── executor.go       # Process spawning
│           │   ├── executor_test.go
│           │   ├── parser.go         # Output parsing
│           │   ├── parser_test.go
│           │   └── types.go          # Execution config, result
│           │
│           ├── checkpoint/           # Git operations
│           │   ├── checkpoint.go     # Commit, rollback operations
│           │   ├── checkpoint_test.go
│           │   ├── git.go            # Git command wrapper
│           │   └── types.go          # Checkpoint metadata
│           │
│           ├── state/                # State management
│           │   ├── manager.go        # Read/write journey state
│           │   ├── manager_test.go
│           │   ├── files.go          # Filesystem operations
│           │   └── types.go          # State structs
│           │
│           ├── project/              # Project detection
│           │   ├── detector.go       # BMAD project detection
│           │   ├── detector_test.go
│           │   ├── validator.go      # Dependency validation
│           │   └── types.go          # Project config
│           │
│           └── common/               # Shared utilities
│               ├── errors.go         # Application error types
│               ├── logger.go         # slog configuration
│               └── constants.go
│
├── packages/                         # Shared packages (optional)
│   └── shared-types/                 # TypeScript types (if needed)
│       ├── package.json
│       └── src/
│           └── index.ts
│
├── scripts/                          # Build & dev scripts
│   ├── build-backend.sh              # Cross-compile Golang
│   ├── dev.sh                        # Start dev environment
│   └── package.sh                    # Build distributable
│
├── .gitignore
├── .prettierrc
├── .eslintrc.js
├── package.json                      # Root workspace config
├── pnpm-workspace.yaml               # pnpm workspace
└── README.md
```

### Architectural Boundaries

**IPC Boundary (Electron ↔ Golang):**

```
┌─────────────────────────────────────────────────────────────┐
│                    Electron Main Process                     │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ src/main/backend.ts                                    │ │
│  │   - spawn(resources/bin/autobmad-{platform})          │ │
│  │   - pipe stdin/stdout                                  │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ src/main/ipc.ts                                        │ │
│  │   - sendRequest(method, params) → Promise<result>     │ │
│  │   - onEvent(event, callback)                          │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                    stdio (JSON-RPC 2.0)
                              │
┌─────────────────────────────────────────────────────────────┐
│                      Golang Backend                          │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ internal/server/server.go                              │ │
│  │   - Read stdin → Parse JSON-RPC → Route to handler    │ │
│  │   - Write stdout ← JSON-RPC response                  │ │
│  │   - Write stdout ← Events (newline-delimited)         │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

**Context Bridge Boundary (Main ↔ Renderer):**

```typescript
// src/preload/index.ts
contextBridge.exposeInMainWorld('api', {
  journey: {
    start: (workflow: string) => ipcRenderer.invoke('journey.start', workflow),
    pause: () => ipcRenderer.invoke('journey.pause'),
    resume: () => ipcRenderer.invoke('journey.resume'),
    abort: () => ipcRenderer.invoke('journey.abort'),
  },
  on: (event: string, callback: Function) => { ... },
  off: (event: string, callback: Function) => { ... },
});
```

**Data Boundary (Filesystem):**

```
Project Root/
├── _bmad/                           # BMAD configuration (read-only)
│   ├── _config/
│   │   └── manifest.yaml
│   └── bmm/
│       └── workflows/
│
├── _bmad-output/                    # BMAD artifacts (source of truth)
│   ├── planning-artifacts/
│   │   ├── prd.md
│   │   ├── architecture.md
│   │   └── ...
│   │
│   └── .autobmad/                   # Auto-BMAD state (hidden)
│       ├── journey-state.json       # Current journey state
│       ├── config.json              # User preferences
│       └── logs/                    # Session logs
│           └── journey-{id}.log
```

### Integration Points

**Internal Communication:**

| From | To | Mechanism |
|------|-----|-----------|
| React UI | Zustand Store | Direct import |
| Zustand Store | Preload API | `window.api.*` |
| Preload | Main Process | `ipcRenderer.invoke()` |
| Main Process | Golang | stdin JSON-RPC |
| Golang | Main Process | stdout JSON-RPC + Events |
| Main Process | React UI | `webContents.send()` → Zustand |

**External Integrations:**

| Integration | Location | Mechanism |
|-------------|----------|-----------|
| OpenCode CLI | `internal/opencode/executor.go` | `exec.Command()` subprocess |
| Git | `internal/checkpoint/git.go` | `exec.Command()` subprocess |
| Filesystem | `internal/state/files.go` | Standard file I/O |
| OS Notifications | `src/main/index.ts` | Electron `Notification` API |

### Development Workflow

**Start Development:**
```bash
# Terminal 1: Golang backend (hot reload with air)
cd apps/core && air

# Terminal 2: Electron + React (Vite HMR)
cd apps/desktop && npm run dev
```

**Build for Distribution:**
```bash
# 1. Build Golang binaries
./scripts/build-backend.sh

# 2. Build Electron app
cd apps/desktop && npm run build

# 3. Package distributable
./scripts/package.sh
```

## Architecture Validation Results

### Coherence Validation ✅

**Decision Compatibility:** All technology choices work together without conflicts. Electron + Golang via stdio/JSON-RPC is a proven pattern. React + Vite + shadcn/ui is a standard modern stack.

**Pattern Consistency:** Naming conventions, error handling, and communication patterns are consistent across all layers. camelCase for JSON, standard conventions for Go and TypeScript.

**Structure Alignment:** Project structure directly supports all architectural decisions. Clear separation between Electron frontend and Golang backend with well-defined IPC boundary.

### Requirements Coverage Validation ✅

**Functional Requirements:** All 52 FRs are architecturally supported.

| Category | Coverage |
|----------|----------|
| Journey Management (FR1-10) | ✅ 100% |
| OpenCode Integration (FR11-19) | ✅ 100% |
| Execution & Retry (FR20-27) | ✅ 100% |
| Failure & Reporting (FR28-35) | ✅ 100% |
| Dashboard & Visualization (FR36-43) | ✅ 100% |
| Project & Configuration (FR44-52) | ✅ 100% |

**Non-Functional Requirements:** All 24 NFRs are architecturally addressed.

| Category | Status |
|----------|--------|
| Performance (NFR-P1 to P3) | ✅ Supported |
| Reliability (NFR-R1 to R3) | ✅ Supported |
| Integration (NFR-I1 to I3) | ✅ Supported |
| Security (NFR-S1 to S5) | ✅ Supported |
| Usability (NFR-U1 to U4) | ✅ Supported |

### Implementation Readiness Validation ✅

**Decision Completeness:** All critical architectural decisions are documented with technology versions and rationale.

**Structure Completeness:** Complete project tree with ~80 files defined, all integration points mapped.

**Pattern Completeness:** Comprehensive patterns for naming, structure, communication, and processes.

### Gap Analysis Results

**Critical Gaps:** None identified. Architecture is complete for MVP implementation.

**Important Gaps (Non-blocking):**
- Keyboard shortcuts: Define during UI implementation
- Log rotation: Implement in logger module
- Auto-update: Add electron-updater post-MVP

**Deferred to Post-MVP:**
- Windows platform support
- Plugin/extension architecture
- Multi-journey parallelism

### Architecture Completeness Checklist

**✅ Requirements Analysis**
- [x] Project context thoroughly analyzed
- [x] Scale and complexity assessed (Medium-High, 15-20 components)
- [x] Technical constraints identified (Electron, Golang, OpenCode CLI)
- [x] Cross-cutting concerns mapped (IPC, state, checkpoints, errors)

**✅ Architectural Decisions**
- [x] IPC Protocol: stdio/JSON-RPC 2.0
- [x] Backend Embedding: Embedded binary
- [x] State Management: Filesystem as source of truth
- [x] OpenCode Integration: One-shot per step
- [x] Checkpoint Strategy: Event-based Git commits

**✅ Implementation Patterns**
- [x] Naming conventions: All layers covered
- [x] Structure patterns: Co-located tests, feature-based
- [x] Communication patterns: JSON-RPC methods and events
- [x] Process patterns: Error handling, logging, state updates

**✅ Project Structure**
- [x] Complete directory structure: ~80 files defined
- [x] Component boundaries: Go packages, React features
- [x] Integration points: IPC, preload, filesystem
- [x] Requirements mapping: All FRs mapped to locations

### Architecture Readiness Assessment

**Overall Status:** ✅ READY FOR IMPLEMENTATION

**Confidence Level:** HIGH

**Key Strengths:**
- Clean 3-layer architecture with clear boundaries
- Filesystem as source of truth ensures zero data loss
- Event-based checkpoints provide rollback capability
- JSON-RPC 2.0 provides standard, debuggable IPC
- Feature-based organization scales with requirements

**Areas for Future Enhancement:**
- Windows platform support
- Electron auto-update mechanism
- Plugin architecture for extensibility
- Journey templates for common workflows

## Architecture Completion Summary

### Workflow Completion

**Architecture Decision Workflow:** COMPLETED ✅
**Total Steps Completed:** 8
**Date Completed:** 2026-01-21
**Document Location:** `_bmad-output/planning-artifacts/architecture.md`

### Final Architecture Deliverables

**📋 Complete Architecture Document**
- All architectural decisions documented with specific versions
- Implementation patterns ensuring AI agent consistency
- Complete project structure with all files and directories
- Requirements to architecture mapping
- Validation confirming coherence and completeness

**🏗️ Implementation Ready Foundation**
- 5 critical architectural decisions made
- 8 implementation pattern categories defined
- 15-20 architectural components specified
- 52 FRs + 24 NFRs fully supported

**📚 AI Agent Implementation Guide**
- Technology stack with verified versions
- Consistency rules that prevent implementation conflicts
- Project structure with clear boundaries
- Integration patterns and communication standards

### Implementation Handoff

**For AI Agents:**
This architecture document is your complete guide for implementing Auto-BMAD. Follow all decisions, patterns, and structures exactly as documented.

**First Implementation Priority:**
```bash
# Initialize monorepo
mkdir -p auto-bmad/{apps,packages,scripts}
cd auto-bmad

# Create Electron app
npm create @electron-vite/app@latest apps/desktop -- --template react-ts

# Create Golang backend
mkdir -p apps/core/cmd/autobmad apps/core/internal
cd apps/core && go mod init github.com/fairyhunter13/auto-bmad/apps/core

# Add shadcn/ui
cd ../desktop && npx shadcn-ui@latest init
```

**Development Sequence:**
1. Initialize project using documented starter template
2. Set up development environment per architecture
3. Implement Golang JSON-RPC server (`internal/server`)
4. Implement Electron spawn + IPC bridge (`src/main`)
5. Build core features following established patterns
6. Maintain consistency with documented rules

### Quality Assurance Checklist

**✅ Architecture Coherence**
- [x] All decisions work together without conflicts
- [x] Technology choices are compatible
- [x] Patterns support the architectural decisions
- [x] Structure aligns with all choices

**✅ Requirements Coverage**
- [x] All 52 functional requirements are supported
- [x] All 24 non-functional requirements are addressed
- [x] Cross-cutting concerns are handled
- [x] Integration points are defined

**✅ Implementation Readiness**
- [x] Decisions are specific and actionable
- [x] Patterns prevent agent conflicts
- [x] Structure is complete and unambiguous
- [x] Examples are provided for clarity

### Project Success Factors

**🎯 Clear Decision Framework**
Every technology choice was made collaboratively with clear rationale, ensuring all stakeholders understand the architectural direction.

**🔧 Consistency Guarantee**
Implementation patterns and rules ensure that multiple AI agents will produce compatible, consistent code that works together seamlessly.

**📋 Complete Coverage**
All project requirements are architecturally supported, with clear mapping from business needs to technical implementation.

**🏗️ Solid Foundation**
The chosen starter template and architectural patterns provide a production-ready foundation following current best practices.

---

**Architecture Status:** ✅ READY FOR IMPLEMENTATION

**Next Phase:** Begin implementation using the architectural decisions and patterns documented herein.

**Document Maintenance:** Update this architecture when major technical decisions are made during implementation.

