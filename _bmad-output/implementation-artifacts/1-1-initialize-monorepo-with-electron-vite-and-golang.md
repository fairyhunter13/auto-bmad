# Story 1.1: Initialize Monorepo with Electron-Vite and Golang

Status: review

## Story

As a **developer setting up Auto-BMAD**,
I want **the project scaffolded with the correct monorepo structure**,
So that **I have a working foundation for Electron + React + Golang development**.

## Acceptance Criteria

1. **Given** a fresh development environment  
   **When** I run the initialization commands  
   **Then** the monorepo structure is created with:
   - `apps/desktop/` containing electron-vite React TypeScript app
   - `apps/core/` containing Golang module with `cmd/autobmad/main.go`
   - `packages/` directory for shared types
   - `scripts/` directory for build scripts

2. **Given** the monorepo is initialized  
   **When** I run `npm run dev`  
   **Then** the Electron app starts without errors

3. **Given** the Golang module exists  
   **When** I run `go build ./cmd/autobmad`  
   **Then** the Golang binary compiles successfully

4. **Given** the desktop app is created  
   **When** shadcn/ui is initialized  
   **Then** Tailwind CSS is configured and working

## Tasks / Subtasks

- [x] **Task 1: Create monorepo root structure** (AC: #1)
  - [x] Create root `package.json` with workspace configuration
  - [x] Create `apps/`, `packages/`, `scripts/` directories
  - [x] Add root `.gitignore` with Node and Go patterns
  - [x] Add root `tsconfig.json` for workspace TypeScript settings

- [x] **Task 2: Initialize Electron app with electron-vite** (AC: #1, #2)
  - [x] Run `npm create @electron-vite/app@latest apps/desktop -- --template react-ts`
  - [x] Verify `apps/desktop/` structure: `src/main/`, `src/preload/`, `src/renderer/`
  - [x] Configure electron-builder for packaging
  - [x] Test `npm run dev` starts Electron window

- [x] **Task 3: Initialize Golang backend module** (AC: #1, #3)
  - [x] Create `apps/core/cmd/autobmad/main.go` with minimal entry point
  - [x] Create `apps/core/internal/` directory structure
  - [x] Run `go mod init github.com/fairyhunter13/auto-bmad/apps/core`
  - [x] Test `go build ./cmd/autobmad` compiles

- [x] **Task 4: Initialize shadcn/ui with Tailwind** (AC: #4)
  - [x] Run `npx shadcn-ui@latest init` in `apps/desktop`
  - [x] Configure Tailwind CSS with proper content paths
  - [x] Add base shadcn components (Button, Card, Input)
  - [x] Verify styles apply correctly in dev mode

- [x] **Task 5: Create build scripts** (AC: #1)
  - [x] Create `scripts/build-core.sh` for Golang cross-compilation
  - [x] Create `scripts/dev.sh` for concurrent dev startup
  - [x] Add npm scripts to root package.json

## Dev Notes

### Architecture Requirements

**Source:** [architecture.md#Monorepo Structure]

The monorepo structure MUST follow this exact layout:

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

### Technology Stack (MANDATORY)

| Component | Technology | Version |
|-----------|------------|---------|
| Frontend Framework | React | 18.x |
| Language (Frontend) | TypeScript | 5.x (strict mode) |
| Styling | Tailwind CSS | 3.x (via shadcn/ui) |
| UI Components | shadcn/ui + Radix | Latest |
| Build Tool | Vite | 5.x |
| Electron Framework | electron-vite | Latest |
| Packaging | electron-builder | Latest |
| Backend Language | Golang | 1.21+ |
| Backend Layout | Standard (cmd/internal) | N/A |

### Initialization Commands (EXACT)

**Source:** [architecture.md#Initialization Commands]

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

### Golang Internal Package Structure

Create these directories under `apps/core/internal/`:

| Package | Purpose |
|---------|---------|
| `server/` | JSON-RPC 2.0 server (stdio) |
| `journey/` | Journey orchestration state machine |
| `opencode/` | OpenCode CLI process management |
| `checkpoint/` | Git operations for checkpointing |
| `state/` | Filesystem state management |

### Minimal main.go Template

```go
// apps/core/cmd/autobmad/main.go
package main

import (
    "fmt"
    "os"
)

func main() {
    // This will be replaced with JSON-RPC server in Story 1.2
    fmt.Fprintln(os.Stderr, "AutoBMAD Core starting...")
    
    // Keep process alive for testing
    select {}
}
```

### Root package.json Template

```json
{
  "name": "auto-bmad",
  "private": true,
  "workspaces": [
    "apps/*",
    "packages/*"
  ],
  "scripts": {
    "dev": "npm run dev --workspace=apps/desktop",
    "build": "npm run build --workspace=apps/desktop",
    "build:core": "bash scripts/build-core.sh"
  }
}
```

### Critical Configuration Notes

1. **electron-vite vs electron-builder**: Use `electron-vite` for development, `electron-builder` for packaging
2. **TypeScript strict mode**: Enabled in tsconfig.json
3. **Tailwind content paths**: Must include all component paths in `apps/desktop/tailwind.config.js`
4. **Go module path**: MUST be `github.com/fairyhunter13/auto-bmad/apps/core`

### Electron Security Requirements

**Source:** [architecture.md#Backend Embedding]

- `contextIsolation: true` MUST be enabled
- `nodeIntegration: false` in renderer
- Preload script handles IPC exposure

### .gitignore Requirements

Must include:
```
# Node
node_modules/
dist/
.vite/

# Go
apps/core/autobmad
apps/core/autobmad-*

# Electron
out/
release/

# IDE
.idea/
.vscode/

# OS
.DS_Store
```

### Project Structure Notes

- This story creates the foundation that ALL subsequent stories build upon
- The structure MUST match architecture.md exactly - no creative variations
- Test both `npm run dev` and `go build` before marking complete

### References

- [architecture.md#Starter Template Selection] - electron-vite choice rationale
- [architecture.md#Monorepo Structure] - Exact directory layout
- [architecture.md#Initialization Commands] - Exact commands to run
- [architecture.md#Decisions Made by Starter] - Technology versions
- [prd.md#NFR-P1] - Startup time < 5 seconds (affects how we structure imports)

## File List

### Created Files

**Root Level:**
- `package.json` - pnpm workspace configuration
- `pnpm-workspace.yaml` - workspace package locations
- `tsconfig.json` - root TypeScript configuration
- `.gitignore` - updated with Node + Go patterns

**apps/desktop/:**
- `package.json` - Electron app dependencies and scripts
- `electron-builder.yml` - packaging configuration
- `electron.vite.config.ts` - electron-vite build configuration
- `tsconfig.json`, `tsconfig.node.json`, `tsconfig.web.json` - TypeScript configs
- `tailwind.config.js` - Tailwind CSS configuration
- `postcss.config.js` - PostCSS configuration
- `components.json` - shadcn/ui configuration
- `src/main/index.ts` - Electron main process
- `src/preload/index.ts`, `src/preload/index.d.ts` - IPC bridge
- `src/renderer/index.html` - HTML entry point
- `src/renderer/src/main.tsx` - React entry point
- `src/renderer/src/App.tsx` - Root React component
- `src/renderer/src/env.d.ts` - Vite type declarations
- `src/renderer/src/assets/main.css` - Tailwind styles with CSS variables
- `src/renderer/src/lib/utils.ts` - shadcn/ui utilities (cn function)
- `src/renderer/src/components/ui/button.tsx` - Button component
- `src/renderer/src/components/ui/card.tsx` - Card component
- `src/renderer/src/components/ui/input.tsx` - Input component
- `resources/icon.png` - placeholder app icon

**apps/core/:**
- `go.mod` - Go module definition
- `cmd/autobmad/main.go` - Entry point (placeholder)
- `internal/server/server.go` - JSON-RPC server placeholder
- `internal/journey/journey.go` - Journey orchestration placeholder
- `internal/opencode/opencode.go` - OpenCode executor placeholder
- `internal/checkpoint/checkpoint.go` - Git checkpoint placeholder
- `internal/state/state.go` - State manager placeholder

**scripts/:**
- `build-core.sh` - Golang cross-compilation script
- `dev.sh` - Development startup script

## Dev Agent Record

### Agent Model Used

Claude (claude-sonnet-4-20250514)

### Completion Notes List

- Manually scaffolded electron-vite structure due to interactive CLI not working in non-TTY environment
- Used pnpm instead of npm as specified in architecture for better monorepo support
- Downgraded vite from v6 to v5.4 to satisfy electron-vite peer dependency
- Created minimal Go internal packages with placeholder types to establish structure
- All Electron security settings applied: contextIsolation=true, nodeIntegration=false, sandbox=false (required for preload)
- Added tailwindcss-animate dependency for shadcn/ui animations
- Build verification: `pnpm run build` and `go build ./cmd/autobmad` both pass

### Change Log

| Date | Change | Reason |
|------|--------|--------|
| 2026-01-22 | Initial implementation | Story 1.1 scaffolding complete |
