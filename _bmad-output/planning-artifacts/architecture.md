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
lastReviewed: '2026-01-22'
reviewType: 'adversarial'
resilienceSectionsAdded: true
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

## System Resilience Architecture

This section addresses failure scenarios, recovery mechanisms, and operational resilience that ensure Auto-BMAD meets its zero-data-loss and high-reliability requirements.

### OpenCode Process Lifecycle & Failure Handling

**CRITICAL: OpenCode is not a reliable black box. The architecture must handle all failure modes.**

#### Process States

```
┌─────────────────────────────────────────────────────────────────┐
│                    OpenCode Process Lifecycle                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────┐    spawn    ┌─────────┐   output   ┌──────────┐   │
│  │ PENDING │ ──────────► │ RUNNING │ ─────────► │ STREAMING│   │
│  └─────────┘             └────┬────┘            └─────┬────┘   │
│                               │                       │         │
│         ┌─────────────────────┼───────────────────────┤         │
│         │                     │                       │         │
│         ▼                     ▼                       ▼         │
│  ┌──────────┐         ┌──────────┐            ┌──────────┐     │
│  │ TIMEOUT  │         │ CRASHED  │            │ COMPLETED│     │
│  └────┬─────┘         └────┬─────┘            └────┬─────┘     │
│       │                    │                       │            │
│       ▼                    ▼                       ▼            │
│  ┌──────────────────────────────────────────────────────┐      │
│  │              CLEANUP (kill process, save state)       │      │
│  └──────────────────────────────────────────────────────┘      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Failure Detection Mechanisms

| Failure Type | Detection Method | Timeout | Recovery Action |
|--------------|------------------|---------|-----------------|
| **Spawn Failure** | `exec.Command` returns error | Immediate | Log error, mark step failed, offer retry |
| **Hang (no output)** | Heartbeat timer (no stdout for N seconds) | 60s default | Send SIGTERM, wait 5s, SIGKILL |
| **Hang (with output)** | Step timeout exceeded | 5min default (configurable) | Same as above |
| **Crash (unexpected exit)** | Exit code ≠ 0 with stderr | Immediate | Capture stderr, parse error, mark failed |
| **Partial Output** | Exit during streaming | Immediate | Save partial output, mark incomplete |
| **Garbage Output** | Output validation fails | Post-completion | Mark step as "needs review" (yellow flag) |
| **Buffer Blocked** | Write to stdin times out | 10s | Log warning, continue (non-fatal) |

#### Process Monitor Component

```go
// internal/opencode/monitor.go

type ProcessMonitor struct {
    process      *exec.Cmd
    heartbeat    *time.Timer     // Reset on each stdout line
    stepTimeout  *time.Timer     // Overall step timeout
    outputBuffer *RingBuffer     // Last N lines for crash diagnosis
    state        ProcessState
    mu           sync.Mutex
}

type ProcessState int
const (
    StatePending ProcessState = iota
    StateRunning
    StateStreaming
    StateCompleted
    StateFailed
    StateTimeout
    StateCrashed
)

// Monitor goroutine runs alongside process
func (m *ProcessMonitor) Run(ctx context.Context) error {
    // 1. Reset heartbeat on each stdout line
    // 2. Check stepTimeout hasn't exceeded
    // 3. Detect process exit
    // 4. Handle cleanup on any termination
}
```

#### Output Validation Rules

| Validation | Criteria | On Failure |
|------------|----------|------------|
| **Non-empty** | Output has at least 1 line | Yellow flag: "OpenCode produced no output" |
| **No error markers** | No `Error:`, `FATAL:`, `panic:` in output | Capture context, mark step failed |
| **Expected artifacts** | Check if expected files were created | Yellow flag: "Expected artifact not found" |
| **Frontmatter valid** | If markdown output, frontmatter parses | Yellow flag: "Artifact may be corrupted" |

#### Timeout Configuration

```json
// journey-state.json
{
  "timeouts": {
    "stepDefault": 300000,       // 5 minutes (ms)
    "heartbeatInterval": 60000,  // 60 seconds
    "gracefulShutdown": 5000,    // 5 seconds SIGTERM→SIGKILL
    "spawnTimeout": 10000        // 10 seconds to spawn
  },
  "stepOverrides": {
    "create-architecture": 600000,  // 10 min for complex steps
    "brainstorming": 900000         // 15 min for creative steps
  }
}
```

### IPC Protocol Resilience

**CRITICAL: The stdio/JSON-RPC channel must not deadlock or lose messages.**

#### Message Framing Protocol

**Problem:** Newline-delimited JSON breaks when messages contain newlines (e.g., code blocks).

**Solution:** Length-prefixed framing with JSON envelope.

```
┌────────────────────────────────────────────────────────────────┐
│                     Message Frame Format                        │
├────────────────────────────────────────────────────────────────┤
│  [4 bytes: length]  [N bytes: JSON payload]  [1 byte: newline] │
│                                                                 │
│  Example:                                                       │
│  \x00\x00\x00\x2F{"jsonrpc":"2.0","method":"ping","id":1}\n    │
│  └─── 47 ────────┘└──────────── 47 bytes ───────────────┘      │
└────────────────────────────────────────────────────────────────┘
```

**Fallback:** If length-prefix adds complexity, use Base64 encoding for payloads with newlines:

```json
{
  "jsonrpc": "2.0",
  "method": "opencode.output",
  "params": {
    "chunk": "SGVsbG8gV29ybGQK...",  // Base64 encoded
    "encoding": "base64"
  }
}
```

#### Backpressure & Flow Control

| Mechanism | Implementation | Trigger |
|-----------|----------------|---------|
| **Acknowledgments** | Every 100 events, backend waits for `ack` | Event count threshold |
| **Buffer High-Water Mark** | Pause event emission at 1000 queued | Queue size |
| **Reader Slow Detection** | If stdin read takes > 1s, log warning | Read latency |
| **Graceful Degradation** | Drop non-critical events (debug logs) | Queue > 5000 |

```go
// internal/server/flow_control.go

type FlowController struct {
    eventQueue     chan Event
    pendingAcks    int
    maxPendingAcks int  // 100
    highWaterMark  int  // 1000 events
    dropThreshold  int  // 5000 events
}

func (fc *FlowController) Emit(event Event) error {
    if len(fc.eventQueue) > fc.dropThreshold {
        if event.Priority == PriorityDebug {
            return nil // Drop low-priority events
        }
    }
    
    if fc.pendingAcks >= fc.maxPendingAcks {
        // Block until ack received
        <-fc.ackReceived
    }
    
    fc.eventQueue <- event
    fc.pendingAcks++
    return nil
}
```

#### Corrupted Message Recovery

| Scenario | Detection | Recovery |
|----------|-----------|----------|
| **Partial JSON** | JSON parse error | Discard, log, request resend |
| **Wrong length** | Payload length ≠ header | Resync: scan for next valid frame |
| **Invalid UTF-8** | Encoding validation | Replace invalid bytes, log warning |
| **Missing fields** | Schema validation | Return JSON-RPC error, continue |

```go
// Resync algorithm
func (r *MessageReader) Resync() error {
    // Scan forward looking for valid JSON-RPC frame
    for {
        // Look for `{"jsonrpc":"2.0"` pattern
        if found {
            // Attempt to parse from here
            // If valid, resume normal operation
            // If invalid, continue scanning
        }
    }
}
```

#### Buffer Size Specifications

| Buffer | Size | Rationale |
|--------|------|-----------|
| **Golang stdout pipe** | 64KB (OS default) | Sufficient for JSON-RPC |
| **Event queue** | 10,000 events | ~10MB at 1KB/event avg |
| **OpenCode output ring buffer** | 1000 lines | Diagnostic context on crash |
| **IPC read buffer** | 1MB | Handle large payloads (code artifacts) |

### Git Checkpoint Error Handling & Recovery

**CRITICAL: Checkpoint operations must never cause data loss.**

#### Failure Modes & Recovery

| Failure | Detection | Recovery | User Impact |
|---------|-----------|----------|-------------|
| **Disk full** | `git commit` returns error | Alert user, pause journey | Yellow flag: "Disk space low" |
| **Git lock contention** | `.git/index.lock` exists | Wait up to 30s, retry 3x | Transparent retry |
| **Branch already exists** | `git checkout -b` fails | Append timestamp: `journey-{ts}-{n}` | None |
| **Dirty index** | `git status --porcelain` non-empty | Stash user changes, commit, unstash | Warning: "Stashed your changes" |
| **Merge conflict** | `git merge` fails | Abort merge, stay on journey branch | Yellow flag with options |
| **Corrupted repo** | `git status` returns error | Alert user, disable checkpoints | Error: "Git repo issue" |

#### Atomic Checkpoint Operation

```go
// internal/checkpoint/checkpoint.go

func (c *Checkpointer) CreateCheckpoint(journeyID string, message string) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // 1. Validate git repo is healthy
    if err := c.validateRepo(); err != nil {
        return c.handleCorruptedRepo(err)
    }
    
    // 2. Handle dirty index from user
    stashed, err := c.stashIfDirty()
    if err != nil {
        return err
    }
    defer func() {
        if stashed {
            c.unstash()
        }
    }()
    
    // 3. Retry loop for transient failures
    for attempt := 0; attempt < 3; attempt++ {
        if err := c.tryCommit(message); err != nil {
            if isTransient(err) {
                time.Sleep(time.Second * time.Duration(attempt+1))
                continue
            }
            return err
        }
        return nil
    }
    
    return ErrCheckpointFailed
}
```

#### Branch Naming with Collision Handling

```
autobmad/journey-{YYYYMMDD-HHMMSS}        # Primary attempt
autobmad/journey-{YYYYMMDD-HHMMSS}-1      # First collision
autobmad/journey-{YYYYMMDD-HHMMSS}-2      # Second collision
```

#### Pre-Commit Validation Checklist

| Check | Command | On Failure |
|-------|---------|------------|
| Repo exists | `git rev-parse --git-dir` | Disable checkpoints |
| Not bare | `git rev-parse --is-bare-repository` | Disable checkpoints |
| No lock | Check `.git/index.lock` | Wait and retry |
| Branch valid | `git rev-parse HEAD` | Init if empty repo |
| Disk space | `df` check for > 100MB free | Yellow flag warning |

### Memory Management Strategy

**CRITICAL: Long-running journeys must not exhaust memory.**

#### Memory Budget by Component

| Component | Idle Budget | Active Budget | Max Growth |
|-----------|-------------|---------------|------------|
| **Electron Main** | 100MB | 150MB | +50MB |
| **React Renderer** | 150MB | 300MB | +150MB |
| **Golang Backend** | 50MB | 200MB | +150MB |
| **Total** | 300MB | 650MB | < 800MB |

#### Output Streaming Strategy

**Problem:** OpenCode output can be megabytes. Holding in memory causes OOM.

**Solution:** Stream to disk, keep only tail in memory.

```go
// internal/opencode/output_handler.go

type OutputHandler struct {
    memoryBuffer  *RingBuffer     // Last 500 lines in memory
    diskFile      *os.File        // Full output on disk
    diskPath      string          // e.g., .autobmad/logs/step-{n}.log
    totalLines    int
    totalBytes    int64
}

func (h *OutputHandler) WriteLine(line string) error {
    // Always write to disk
    if _, err := h.diskFile.WriteString(line + "\n"); err != nil {
        return err
    }
    h.totalBytes += int64(len(line) + 1)
    h.totalLines++
    
    // Keep tail in memory for UI
    h.memoryBuffer.Write(line)
    
    return nil
}

func (h *OutputHandler) GetRecentOutput(lines int) []string {
    return h.memoryBuffer.Last(lines)
}

func (h *OutputHandler) GetFullOutputPath() string {
    return h.diskPath
}
```

#### Log Rotation Policy

| Log Type | Max Size | Max Files | Rotation Trigger |
|----------|----------|-----------|------------------|
| **Step output** | 10MB | Per step (no limit) | Step completion |
| **Journey log** | 50MB | 10 journeys | Journey completion |
| **Application log** | 20MB | 5 files | Size exceeded |
| **Debug log** | 100MB | 3 files | Size exceeded |

```go
// internal/common/logger.go

type RotatingLogger struct {
    maxSize   int64  // bytes
    maxFiles  int
    currentFile *os.File
    currentSize int64
}

func (l *RotatingLogger) rotate() error {
    // 1. Close current file
    // 2. Rename: app.log → app.log.1, app.log.1 → app.log.2, etc.
    // 3. Delete oldest if > maxFiles
    // 4. Create new app.log
}
```

#### Journey State Pruning

```json
// .autobmad/config.json
{
  "retention": {
    "completedJourneys": 100,       // Keep last N journey states
    "journeyLogsMaxAge": "30d",     // Delete logs older than
    "archiveAfter": "7d",           // Move to archive after
    "archiveLocation": ".autobmad/archive/"
  }
}
```

#### Memory Monitoring

```go
// internal/common/memory_monitor.go

type MemoryMonitor struct {
    warningThreshold  uint64  // 600MB
    criticalThreshold uint64  // 750MB
    checkInterval     time.Duration
}

func (m *MemoryMonitor) Run(ctx context.Context) {
    ticker := time.NewTicker(m.checkInterval)
    for {
        select {
        case <-ticker.C:
            var mem runtime.MemStats
            runtime.ReadMemStats(&mem)
            
            if mem.Alloc > m.criticalThreshold {
                // Force GC, emit warning event
                runtime.GC()
                emitEvent("system.memoryWarning", "critical")
            } else if mem.Alloc > m.warningThreshold {
                emitEvent("system.memoryWarning", "warning")
            }
        case <-ctx.Done():
            return
        }
    }
}
```

### Yellow Flag Lifecycle & Timeout Policy

**CRITICAL: Yellow flags must not block indefinitely or leak resources.**

#### Yellow Flag State Machine

```
┌─────────────────────────────────────────────────────────────────┐
│                    Yellow Flag Lifecycle                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Step Needs Input                                                │
│        │                                                         │
│        ▼                                                         │
│  ┌──────────┐    30 min    ┌──────────────┐    4 hr    ┌──────┐ │
│  │ RAISED   │ ───────────► │ REMINDER     │ ────────► │URGENT│ │
│  │          │              │ (notification)│           │      │ │
│  └────┬─────┘              └──────┬───────┘           └───┬──┘ │
│       │                          │                        │     │
│       │ user responds            │ user responds          │     │
│       ▼                          ▼                        ▼     │
│  ┌──────────┐              ┌──────────┐            ┌──────────┐ │
│  │ RESOLVED │              │ RESOLVED │            │AUTO-DECIDE│ │
│  └──────────┘              └──────────┘            └─────┬────┘ │
│                                                          │      │
│                                                          ▼      │
│                            48 hr from RAISED      ┌──────────┐  │
│                            ───────────────────►   │ ARCHIVED │  │
│                            (no response)          └──────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Timeout Policy

| Timeout | Duration | Action | User Notification |
|---------|----------|--------|-------------------|
| **Initial** | 0 | Raise yellow flag | Desktop notification + in-app |
| **Reminder** | 30 min | Send reminder | Desktop notification |
| **Escalation** | 4 hr | Offer "Let AI Decide" | Urgent notification |
| **Auto-Archive** | 48 hr | Archive journey | Final notification |

#### Resource State During Yellow Flag

| Resource | State During Wait | Cleanup |
|----------|-------------------|---------|
| **OpenCode Process** | Terminated (exit or killed) | Already cleaned up |
| **Golang Backend** | Running, minimal CPU | N/A |
| **Journey State** | Persisted to disk | N/A |
| **UI** | Showing yellow flag modal | Updates on resolution |
| **Memory** | Output already streamed to disk | N/A |

```go
// internal/journey/yellow_flag.go

type YellowFlag struct {
    JourneyID    string
    StepIndex    int
    Reason       string
    Options      []YellowFlagOption
    RaisedAt     time.Time
    State        YellowFlagState
    ReminderSent bool
    EscalatedAt  *time.Time
}

type YellowFlagState int
const (
    YellowFlagRaised YellowFlagState = iota
    YellowFlagReminded
    YellowFlagEscalated
    YellowFlagResolved
    YellowFlagAutoDecided
    YellowFlagArchived
)

func (yf *YellowFlag) CheckTimeouts() {
    elapsed := time.Since(yf.RaisedAt)
    
    switch {
    case elapsed > 48*time.Hour && yf.State < YellowFlagArchived:
        yf.archive()
    case elapsed > 4*time.Hour && yf.State < YellowFlagEscalated:
        yf.escalate()
    case elapsed > 30*time.Minute && yf.State < YellowFlagReminded:
        yf.sendReminder()
    }
}

func (yf *YellowFlag) escalate() {
    yf.State = YellowFlagEscalated
    yf.EscalatedAt = timePtr(time.Now())
    emitEvent("yellowFlag.escalated", YellowFlagPayload{
        JourneyID: yf.JourneyID,
        StepIndex: yf.StepIndex,
        Options:   append(yf.Options, YellowFlagOption{
            ID:    "auto-decide",
            Label: "Let AI Decide",
            Hint:  "AI will choose based on context",
        }),
    })
}
```

#### Yellow Flag Options Schema

```json
{
  "yellowFlag": {
    "journeyId": "j-20260121-001",
    "stepIndex": 3,
    "reason": "Multiple architecture patterns possible",
    "raisedAt": "2026-01-21T10:30:00Z",
    "options": [
      {
        "id": "option-1",
        "label": "Use Microservices",
        "hint": "Better for scaling, more complex",
        "confidence": 0.7
      },
      {
        "id": "option-2",
        "label": "Use Monolith",
        "hint": "Simpler, faster to implement",
        "confidence": 0.6
      },
      {
        "id": "retry",
        "label": "Retry Step",
        "hint": "Try again with different approach"
      },
      {
        "id": "skip",
        "label": "Skip Step",
        "hint": "Continue without this step"
      }
    ],
    "timeoutPolicy": {
      "reminderAt": "2026-01-21T11:00:00Z",
      "escalateAt": "2026-01-21T14:30:00Z",
      "archiveAt": "2026-01-23T10:30:00Z"
    }
  }
}
```

### System Crash Recovery Sequence

**CRITICAL: Users must never lose work due to crashes.**

#### Crash Scenarios & Recovery

| Crash Type | Detection | Recovery Sequence |
|------------|-----------|-------------------|
| **Electron Main Crash** | Process monitor, OS signal | Backend continues, restart Electron, reconnect |
| **Golang Backend Crash** | Electron detects stdio EOF | Save UI state, restart backend, restore from disk |
| **OpenCode Crash** | Exit code, process monitor | Capture partial output, mark step failed, offer retry |
| **System Power Loss** | Next startup | Read journey-state.json, validate vs git, offer resume |
| **User Force Quit** | SIGTERM/SIGKILL | Same as power loss |

#### Startup Recovery Sequence

```
┌─────────────────────────────────────────────────────────────────┐
│                    Startup Recovery Sequence                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  App Launch                                                      │
│      │                                                           │
│      ▼                                                           │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │ 1. Check for .autobmad/journey-state.json                │   │
│  └─────────────────────────┬────────────────────────────────┘   │
│                            │                                     │
│         ┌──────────────────┴──────────────────┐                 │
│         │                                      │                 │
│         ▼ exists                               ▼ not exists      │
│  ┌─────────────────┐                    ┌─────────────────┐     │
│  │ 2. Parse state  │                    │ Normal startup  │     │
│  └────────┬────────┘                    └─────────────────┘     │
│           │                                                      │
│           ▼                                                      │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ 3. Validate state:                                      │    │
│  │    - JSON schema valid?                                  │    │
│  │    - Referenced files exist?                             │    │
│  │    - Git branch exists?                                  │    │
│  │    - Checksum matches (if stored)?                       │    │
│  └────────┬────────────────────────────────┬───────────────┘    │
│           │ valid                          │ corrupted           │
│           ▼                                ▼                     │
│  ┌─────────────────┐               ┌─────────────────┐          │
│  │ 4. Check status │               │ Attempt Git     │          │
│  │    of journey   │               │ checkpoint      │          │
│  └────────┬────────┘               │ recovery        │          │
│           │                        └────────┬────────┘          │
│           │                                 │                    │
│  ┌────────┴────────────────────────────────┐                    │
│  │                                          │                    │
│  ▼ running/paused                          ▼ completed           │
│  ┌─────────────────┐               ┌─────────────────┐          │
│  │ 5. Show resume  │               │ Normal startup  │          │
│  │    dialog       │               │ (mark seen)     │          │
│  └─────────────────┘               └─────────────────┘          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### State File Schema with Integrity

```json
// .autobmad/journey-state.json
{
  "version": 1,
  "checksum": "sha256:abc123...",  // Hash of state content (excluding checksum)
  "lastModified": "2026-01-21T10:35:00Z",
  "journey": {
    "id": "j-20260121-001",
    "workflow": "create-prd",
    "status": "running",
    "currentStepIndex": 3,
    "startedAt": "2026-01-21T10:00:00Z",
    "gitBranch": "autobmad/journey-20260121-100000",
    "lastCheckpoint": "abc123def456"
  },
  "steps": [
    {
      "index": 0,
      "name": "init",
      "status": "completed",
      "checkpointSha": "sha1:..."
    }
  ],
  "yellowFlag": null
}
```

#### Orphan Process Cleanup

```go
// internal/opencode/cleanup.go

func CleanupOrphanedProcesses() error {
    // 1. Read PID file if exists: .autobmad/opencode.pid
    // 2. Check if process is still running
    // 3. If running, send SIGTERM
    // 4. Wait 5s, send SIGKILL if still alive
    // 5. Delete PID file
    
    pidFile := ".autobmad/opencode.pid"
    if _, err := os.Stat(pidFile); os.IsNotExist(err) {
        return nil // No orphan
    }
    
    pidBytes, _ := os.ReadFile(pidFile)
    pid, _ := strconv.Atoi(string(pidBytes))
    
    process, err := os.FindProcess(pid)
    if err != nil {
        os.Remove(pidFile)
        return nil
    }
    
    // Try graceful shutdown first
    process.Signal(syscall.SIGTERM)
    
    // Wait with timeout
    done := make(chan error)
    go func() {
        _, err := process.Wait()
        done <- err
    }()
    
    select {
    case <-done:
        // Process exited
    case <-time.After(5 * time.Second):
        // Force kill
        process.Signal(syscall.SIGKILL)
    }
    
    os.Remove(pidFile)
    return nil
}
```

### Security Model & Input Validation

**CRITICAL: IPC messages must not be blindly trusted.**

#### Threat Model

| Threat | Vector | Mitigation |
|--------|--------|------------|
| **Path Traversal** | Malicious file paths in `project.detect` | Allowlist validation |
| **Command Injection** | Malformed prompts to OpenCode | Escape/quote all inputs |
| **Git Injection** | Malicious branch names, commit messages | Sanitize special characters |
| **IPC Spoofing** | Malicious renderer sends fake requests | Context isolation + validation |
| **Resource Exhaustion** | Huge payloads in IPC | Size limits on all inputs |
| **Sensitive Data Exposure** | Logging API keys | Never log prompt content |

#### Path Validation

```go
// internal/common/security.go

var ErrPathTraversal = errors.New("path traversal detected")

func ValidatePath(basePath, requestedPath string) error {
    // 1. Resolve to absolute path
    absBase, _ := filepath.Abs(basePath)
    absRequested, _ := filepath.Abs(requestedPath)
    
    // 2. Check requested is within base
    if !strings.HasPrefix(absRequested, absBase) {
        return ErrPathTraversal
    }
    
    // 3. Check for suspicious patterns
    suspicious := []string{"..", "~", "$", "`", "|", ";", "&"}
    for _, s := range suspicious {
        if strings.Contains(requestedPath, s) {
            return ErrPathTraversal
        }
    }
    
    return nil
}
```

#### Input Validation for JSON-RPC Handlers

| Method | Validation Rules |
|--------|------------------|
| `journey.start` | Workflow name must exist in manifest, no path separators |
| `project.detect` | Path must be within allowed directories |
| `checkpoint.rollback` | SHA must match `^[a-f0-9]{40}$` |
| `config.set` | Key must be in allowlist, value type must match schema |
| `opencode.execute` | Prompt length < 100KB, no null bytes |

```go
// internal/server/validators.go

type JourneyStartRequest struct {
    Workflow string `json:"workflow" validate:"required,alphanum,max=100"`
}

type CheckpointRollbackRequest struct {
    SHA string `json:"sha" validate:"required,hexadecimal,len=40"`
}

type ProjectDetectRequest struct {
    Path string `json:"path" validate:"required,filepath"`
}

func (h *Handler) validateRequest(method string, params json.RawMessage) error {
    switch method {
    case "journey.start":
        var req JourneyStartRequest
        if err := json.Unmarshal(params, &req); err != nil {
            return ErrInvalidParams
        }
        return validate.Struct(req)
        
    case "checkpoint.rollback":
        var req CheckpointRollbackRequest
        // ... validation
        
    case "project.detect":
        var req ProjectDetectRequest
        if err := json.Unmarshal(params, &req); err != nil {
            return ErrInvalidParams
        }
        // Additional path traversal check
        return ValidatePath(h.allowedBasePath, req.Path)
    }
    return nil
}
```

#### Electron Security Configuration

```typescript
// src/main/index.ts

const mainWindow = new BrowserWindow({
    webPreferences: {
        // MANDATORY security settings
        contextIsolation: true,      // Separate contexts
        nodeIntegration: false,      // No Node in renderer
        sandbox: true,               // OS-level sandboxing
        webSecurity: true,           // Same-origin policy
        allowRunningInsecureContent: false,
        
        preload: path.join(__dirname, '../preload/index.js'),
    },
});

// Restrict navigation
mainWindow.webContents.on('will-navigate', (event, url) => {
    // Only allow navigation within app
    if (!url.startsWith('file://')) {
        event.preventDefault();
    }
});

// Block new windows
mainWindow.webContents.setWindowOpenHandler(() => {
    return { action: 'deny' };
});
```

#### Preload Script Security

```typescript
// src/preload/index.ts

import { contextBridge, ipcRenderer } from 'electron';

// ONLY expose specific methods, NEVER expose ipcRenderer directly
contextBridge.exposeInMainWorld('api', {
    journey: {
        start: (workflow: string) => {
            // Validate before sending
            if (typeof workflow !== 'string' || workflow.length > 100) {
                throw new Error('Invalid workflow');
            }
            return ipcRenderer.invoke('journey.start', workflow);
        },
        pause: () => ipcRenderer.invoke('journey.pause'),
        resume: () => ipcRenderer.invoke('journey.resume'),
        abort: () => ipcRenderer.invoke('journey.abort'),
    },
    
    // Event subscriptions with cleanup
    on: (channel: string, callback: Function) => {
        const allowedChannels = [
            'journey.started', 'journey.completed', 'step.started',
            'step.completed', 'step.failed', 'yellowFlag.raised',
        ];
        if (!allowedChannels.includes(channel)) {
            throw new Error(`Channel not allowed: ${channel}`);
        }
        const subscription = (_event: any, ...args: any[]) => callback(...args);
        ipcRenderer.on(channel, subscription);
        return () => ipcRenderer.removeListener(channel, subscription);
    },
});
```

#### Sensitive Data Handling

| Data Type | Handling | Storage |
|-----------|----------|---------|
| **API Keys** | Never logged, never in state | Only in OpenCode config |
| **Prompt Content** | Truncate in logs (first 100 chars) | Full text only in artifacts |
| **File Paths** | Log relative paths only | Full paths in state |
| **Error Messages** | Sanitize before UI display | Full in debug logs |

```go
// internal/common/logger.go

func SanitizeForLog(content string) string {
    // Truncate long content
    if len(content) > 100 {
        content = content[:100] + "..."
    }
    
    // Redact potential secrets
    patterns := []string{
        `(?i)(api[_-]?key|secret|password|token)\s*[:=]\s*\S+`,
        `sk-[a-zA-Z0-9]{32,}`,  // OpenAI key pattern
    }
    for _, p := range patterns {
        re := regexp.MustCompile(p)
        content = re.ReplaceAllString(content, "[REDACTED]")
    }
    
    return content
}
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

**Critical Gaps:** ✅ ADDRESSED (2026-01-22 Adversarial Review)

The following critical and high-severity gaps were identified and resolved:

| Issue | Severity | Resolution |
|-------|----------|------------|
| OpenCode CLI failure handling | Critical | Added Process Lifecycle & Failure Handling section |
| IPC backpressure/flow control | Critical | Added IPC Protocol Resilience section |
| Git checkpoint error handling | High | Added Git Checkpoint Error Handling section |
| Memory budget for long journeys | High | Added Memory Management Strategy section |
| Yellow flag timeout policy | High | Added Yellow Flag Lifecycle & Timeout Policy section |
| Crash recovery sequence | High | Added System Crash Recovery Sequence section |
| Security model/input validation | High | Added Security Model & Input Validation section |

**Medium Gaps (Address During Implementation):**
- UI-Backend state synchronization: Define event subscription model with sequence numbers
- API contract versioning: Add version handshake on connection
- Integration test strategy: Add test harness specification
- Corruption detection: Implement atomic writes with checksums
- Brownfield detection: Explicitly mark as post-MVP or expand project structure
- Network connectivity handling: Add network monitor component

**Low-Priority Gaps (Nice to Have):**
- OpenCode profile detection algorithm
- Emergency stop implementation details
- Keyboard shortcuts: Define during UI implementation
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

**✅ System Resilience (Added 2026-01-22)**
- [x] OpenCode process lifecycle: States, timeouts, failure detection, recovery
- [x] IPC protocol resilience: Message framing, backpressure, flow control
- [x] Git checkpoint robustness: Error handling, retry logic, collision handling
- [x] Memory management: Budgets, streaming, log rotation, pruning
- [x] Yellow flag lifecycle: Timeout policy, escalation, auto-archive
- [x] Crash recovery: Startup sequence, state validation, orphan cleanup
- [x] Security model: Input validation, path traversal, Electron hardening

### Architecture Readiness Assessment

**Overall Status:** ✅ READY FOR IMPLEMENTATION

**Confidence Level:** HIGH

**Key Strengths:**
- Clean 3-layer architecture with clear boundaries
- Filesystem as source of truth ensures zero data loss
- Event-based checkpoints provide rollback capability
- JSON-RPC 2.0 provides standard, debuggable IPC
- Feature-based organization scales with requirements
- Comprehensive resilience architecture covering all failure modes
- Security-first design with input validation and path protection
- Memory-conscious design for long-running operations

**Areas for Future Enhancement:**
- Windows platform support
- Electron auto-update mechanism
- Plugin architecture for extensibility
- Journey templates for common workflows
- Advanced brownfield discovery protocol

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
- **System Resilience Architecture** (added 2026-01-22 after adversarial review)

**🏗️ Implementation Ready Foundation**
- 5 critical architectural decisions made
- 8 implementation pattern categories defined
- 7 resilience subsystems specified (OpenCode, IPC, Git, Memory, Yellow Flag, Recovery, Security)
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

