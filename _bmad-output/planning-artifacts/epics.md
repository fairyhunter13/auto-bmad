---
stepsCompleted: [1, 2, 3, 4]
inputDocuments:
  - _bmad-output/planning-artifacts/prd.md
  - _bmad-output/planning-artifacts/architecture.md
  - _bmad-output/planning-artifacts/ux-design-specification.md
workflowType: epics-and-stories
project_name: auto-bmad
user_name: Hafiz
date: 2026-01-21
status: complete
completedAt: 2026-01-21
---

# Auto-BMAD - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for Auto-BMAD, decomposing the requirements from the PRD, UX Design, and Architecture into implementable stories.

## Requirements Inventory

### Functional Requirements

**Journey Management (FR1-FR10):**
- FR1: User can create a new journey by specifying a destination (target artifact/phase)
- FR2: System can calculate a route from current state to destination
- FR3: User can view the calculated route before starting execution
- FR4: User can start journey execution with one action
- FR5: User can pause an active journey
- FR6: User can cancel an active journey
- FR7: User can view journey history (past journeys and outcomes)
- FR8: System can resume a paused journey from last checkpoint
- FR9: User can provide feedback during journey execution
- FR10: System can adjust journey direction based on user feedback

**OpenCode Integration (FR11-FR19):**
- FR11: System can detect installed OpenCode CLI
- FR12: System can detect OpenCode CLI version and verify compatibility
- FR13: System can detect available OpenCode profiles from user configuration
- FR14: User can select which OpenCode profile to use for a project
- FR15: System can spawn OpenCode CLI processes with specified configuration
- FR16: System can monitor active OpenCode processes
- FR17: System can capture OpenCode output (stdout, stderr)
- FR18: System can terminate OpenCode processes
- FR19: System can pass BMAD workflow/agent context to OpenCode

**Execution & Retry (FR20-FR27):**
- FR20: System can execute BMAD workflows sequentially per journey route
- FR21: System can detect workflow completion (success or failure)
- FR22: System can automatically retry failed workflows
- FR23: System can accumulate feedback from previous retry attempts
- FR24: System can pass accumulated feedback to retry attempts
- FR25: User can configure maximum retry attempts
- FR26: System can escalate to user after exceeding retry threshold
- FR27: System can track retry count and history per workflow

**Failure & Reporting (FR28-FR35):**
- FR28: System can detect workflow failures accurately (0% false success)
- FR29: System can report failures with clear explanations
- FR30: System can provide failure timeline (events leading to failure)
- FR31: System can identify root cause category (system, network, logic, etc.)
- FR32: System can report what was saved before failure
- FR33: System can report what was lost due to failure
- FR34: User can view failure details in dashboard
- FR35: System can suggest recovery options based on failure type

**Dashboard & Visualization (FR36-FR43):**
- FR36: User can view current journey status (idle, running, paused, failed, complete)
- FR37: User can view journey progress (current phase, step, percentage)
- FR38: User can view active workflow details
- FR39: User can view generated artifacts
- FR40: User can view retry progress (attempt number, feedback history)
- FR41: User can view journey completion summary (duration, artifacts, retries)
- FR42: User can access journey history
- FR43: System can display notifications for important events

**Project & Configuration (FR44-FR52):**
- FR44: User can select a project folder to work with
- FR45: System can detect if project has `_bmad/` folder
- FR46: System can detect if project has `_bmad-output/` folder
- FR47: System can detect project type (greenfield vs brownfield)
- FR48: User can describe project context manually (for brownfield without auto-scan)
- FR49: User can configure Auto-BMAD settings (retry limits, notifications, etc.)
- FR50: System can persist user preferences across sessions
- FR51: System can detect network connectivity status
- FR52: System can warn user of network issues during journey

### Non-Functional Requirements

**Performance (NFR-P1 to NFR-P6):**
- NFR-P1: Application startup time < 5 seconds
- NFR-P2: Memory usage (idle) < 500 MB
- NFR-P3: Dashboard UI updates near real-time (< 100ms)
- NFR-P4: OpenCode process spawn time < 2 seconds
- NFR-P5: Feedback submission response < 500ms acknowledgment
- NFR-P6: Journey state save < 1 second

**Reliability (NFR-R1 to NFR-R7):**
- NFR-R1: Zero data loss on crash
- NFR-R2: Checkpoint frequency: continuous (after every significant state change)
- NFR-R3: Crash recovery: resume from last saved state
- NFR-R4: OpenCode crash detection < 5 seconds
- NFR-R5: OpenCode auto-retry on crash: automatic restart and continue
- NFR-R6: Journey state persistence before any risky operation
- NFR-R7: Graceful shutdown: save all state before exit

**Integration (NFR-I1 to NFR-I7):**
- NFR-I1: OpenCode CLI version: any installed version (use global default)
- NFR-I2: OpenCode profiles: support multiple profiles with load-balancing
- NFR-I3: Git version: standard user Git installation (2.0+)
- NFR-I4: BMAD version: current `_bmad/` in project directory (6.0.0+)
- NFR-I5: Profile misconfiguration: warn user with clear message
- NFR-I6: OpenCode output capture: complete stdout/stderr capture
- NFR-I7: Git operations: use system Git credentials

**Security (NFR-S1 to NFR-S5):**
- NFR-S1: API key storage: none — Auto-BMAD never stores API keys
- NFR-S2: Credential handling: delegated to OpenCode
- NFR-S3: Update verification: signature verification required
- NFR-S4: Code signing: macOS signed and notarized
- NFR-S5: IPC security: Electron best practices (contextIsolation)

**Usability (NFR-U1 to NFR-U4):**
- NFR-U1: Keyboard navigation: full keyboard support
- NFR-U2: Color contrast: WCAG AA minimum
- NFR-U3: Font sizing: configurable (Post-MVP)
- NFR-U4: Screen reader: basic support (Post-MVP)

### Additional Requirements

**From Architecture:**
- Starter Template: electron-vite + React + TypeScript + Golang monorepo structure
- IPC Protocol: stdio/JSON-RPC 2.0 between Electron Main and Golang backend
- Backend Embedding: Golang binary in Electron `resources/bin/` directory
- State Source of Truth: Filesystem (`_bmad-output/.autobmad/journey-state.json`)
- OpenCode Process Model: One-shot per step with `opencode -p "{prompt}" --non-interactive`
- Checkpoint Strategy: Event-based Git commits (step completion, pause, yellow flag, journey end)
- Git Branch Strategy: `autobmad/journey-{timestamp}` branches with descriptive commit messages
- JSON Field Naming: camelCase throughout IPC boundary
- Error Codes: JSON-RPC 2.0 compliant with application-specific ranges (-32000 to -32399)
- Test Organization: Co-located tests (`.test.tsx` / `_test.go` next to source)

**From UX Design:**
- Color-coded Status System: 🟢 (success/walk away), 🟡 (needs input), 🔴 (failed), 🔵 (running)
- Kanban Dashboard Layout: Journey cards in status columns (Running | Yellow Flag | Paused | Complete | Failed)
- Command Palette: Cmd/Ctrl+K for quick actions and universal search
- Yellow Flag Response Loop: Notification → Context → Input → Continue in < 30 seconds
- Progressive Disclosure: Simple journey setup by default, advanced options expandable
- Split View Layout: Journey timeline on left, status/output/feedback on right
- Status Bar: Persistent bottom bar showing active journey phase, elapsed time, provider status
- Desktop Notifications: Non-intrusive alerts for yellow flags and failures
- "Walk Away Confidence" Design: Visual language communicating "I've got this, you can go"

### FR Coverage Map

**Journey Management:**
- FR1: Epic 2 - Create journey by destination
- FR2: Epic 2 - Calculate route
- FR3: Epic 2 - View calculated route
- FR4: Epic 2 - Start journey execution
- FR5: Epic 3 - Pause journey
- FR6: Epic 3 - Cancel journey
- FR7: Epic 6 - View journey history
- FR8: Epic 4 - Resume from checkpoint
- FR9: Epic 4 - Provide feedback
- FR10: Epic 4 - Adjust direction

**OpenCode Integration:**
- FR11: Epic 1 - Detect OpenCode CLI
- FR12: Epic 1 - Verify version compatibility
- FR13: Epic 1 - Detect available profiles
- FR14: Epic 1 - Select profile for project
- FR15: Epic 3 - Spawn OpenCode processes
- FR16: Epic 3 - Monitor processes
- FR17: Epic 3 - Capture output
- FR18: Epic 3 - Terminate processes
- FR19: Epic 3 - Pass BMAD context

**Execution & Retry:**
- FR20: Epic 3 - Execute workflows sequentially
- FR21: Epic 3 - Detect completion
- FR22: Epic 5 - Auto-retry failed workflows
- FR23: Epic 5 - Accumulate feedback
- FR24: Epic 5 - Pass accumulated feedback
- FR25: Epic 5 - Configure max retries
- FR26: Epic 5 - Escalate after threshold
- FR27: Epic 5 - Track retry history

**Failure & Reporting:**
- FR28: Epic 5 - Detect failures accurately
- FR29: Epic 5 - Report with explanations
- FR30: Epic 5 - Failure timeline
- FR31: Epic 5 - Root cause category
- FR32: Epic 5 - Report saved state
- FR33: Epic 5 - Report lost state
- FR34: Epic 5 - View failure details
- FR35: Epic 5 - Suggest recovery options

**Dashboard & Visualization:**
- FR36: Epic 2 - View journey status
- FR37: Epic 2 - View progress
- FR38: Epic 2 - View active workflow
- FR39: Epic 2 - View artifacts
- FR40: Epic 6 - View retry progress
- FR41: Epic 6 - Completion summary
- FR42: Epic 6 - Access history
- FR43: Epic 4 - Display notifications

**Project & Configuration:**
- FR44: Epic 1 - Select project folder
- FR45: Epic 1 - Detect _bmad/ folder
- FR46: Epic 1 - Detect _bmad-output/ folder
- FR47: Epic 1 - Detect project type
- FR48: Epic 1 - Manual context description
- FR49: Epic 1 - Configure settings
- FR50: Epic 1 - Persist preferences
- FR51: Epic 1 - Detect network status
- FR52: Epic 1 - Warn network issues

## Epic List

### Epic 1: Project Foundation & OpenCode Integration

**Goal:** Users can open any BMAD project and have Auto-BMAD detect the project structure, verify dependencies (Git, OpenCode), and be ready to start journeys.

**FRs Covered:** FR11-FR14, FR44-FR52 (18 FRs)
**NFRs Addressed:** NFR-I1 to NFR-I7, NFR-S1, NFR-S2, NFR-S5

**Implementation Notes:**
- Story 1.1 uses starter template from Architecture (electron-vite + Golang monorepo)
- Includes JSON-RPC server foundation for IPC
- Establishes filesystem state management pattern

---

### Epic 2: Journey Planning & Visualization

**Goal:** Users can create journeys by selecting destinations, view calculated routes before starting, and understand exactly what will happen when they click "Start."

**FRs Covered:** FR1-FR4, FR36-FR39 (8 FRs)
**NFRs Addressed:** NFR-P1, NFR-P3, NFR-U1, NFR-U2

**Implementation Notes:**
- Builds "walk away confidence" through route transparency
- Kanban dashboard with status columns
- Color-coded status indicators (🟢🟡🔴🔵)

---

### Epic 3: Autonomous Workflow Execution

**Goal:** Users can start journeys and have Auto-BMAD execute BMAD workflows autonomously via OpenCode, with real-time output capture and checkpoint safety.

**FRs Covered:** FR5-FR6, FR15-FR21 (9 FRs)
**NFRs Addressed:** NFR-P4, NFR-P6, NFR-R1 to NFR-R7, NFR-I6

**Implementation Notes:**
- One-shot OpenCode process per step
- Event-based Git checkpoint commits
- Real-time stdout/stderr streaming

---

### Epic 4: Feedback System & Adaptive Direction

**Goal:** Users can respond to yellow flags (input requests) quickly, provide feedback during execution, and have the system adapt direction based on their input.

**FRs Covered:** FR8-FR10, FR43 (4 FRs)
**NFRs Addressed:** NFR-P5

**Implementation Notes:**
- Yellow flag notification system
- Context-aware feedback prompts
- < 30 second response loop target

---

### Epic 5: Auto-Retry & Failure Recovery

**Goal:** Users experience intelligent auto-retry on failures, accumulating feedback across attempts, with honest failure reporting when escalation is needed, and clear recovery options.

**FRs Covered:** FR22-FR35 (14 FRs)
**NFRs Addressed:** NFR-R4, NFR-R5

**Implementation Notes:**
- Configurable retry limits
- Feedback accumulation across attempts
- Root cause categorization
- Recovery option suggestions

---

### Epic 6: Dashboard & Journey History

**Goal:** Users can view complete journey history, completion summaries, retry statistics, and access past artifacts — the full observability experience.

**FRs Covered:** FR7, FR40-FR42 (4 FRs)
**NFRs Addressed:** NFR-U3, NFR-U4 (Post-MVP items)

**Implementation Notes:**
- Journey history persistence
- Completion statistics and summaries
- Artifact access and navigation

---

## Epic 1: Project Foundation & OpenCode Integration

**Goal:** Users can open any BMAD project and have Auto-BMAD detect the project structure, verify dependencies (Git, OpenCode), and be ready to start journeys.

### Story 1.1: Initialize Monorepo with Electron-Vite and Golang

As a **developer setting up Auto-BMAD**,
I want **the project scaffolded with the correct monorepo structure**,
So that **I have a working foundation for Electron + React + Golang development**.

**Acceptance Criteria:**

**Given** a fresh development environment
**When** I run the initialization commands
**Then** the monorepo structure is created with:
- `apps/desktop/` containing electron-vite React TypeScript app
- `apps/core/` containing Golang module with `cmd/autobmad/main.go`
- `packages/` directory for shared types
- `scripts/` directory for build scripts
**And** `npm run dev` starts the Electron app without errors
**And** `go build ./cmd/autobmad` compiles the Golang binary
**And** shadcn/ui is initialized with Tailwind CSS configured

---

### Story 1.2: Implement JSON-RPC Server Foundation

As a **developer building the backend**,
I want **a JSON-RPC 2.0 server that reads from stdin and writes to stdout**,
So that **the Electron main process can communicate with the Golang backend**.

**Acceptance Criteria:**

**Given** the Golang binary is running
**When** a valid JSON-RPC request is sent via stdin
**Then** a valid JSON-RPC response is returned via stdout
**And** the response follows JSON-RPC 2.0 specification exactly
**And** invalid requests return proper error responses with code -32600

**Given** a `system.ping` method is called
**When** the request is processed
**Then** the response contains `{"result": "pong"}` with matching request ID

**Given** an unknown method is called
**When** the request is processed
**Then** the response contains error code -32601 (Method not found)

---

### Story 1.3: Create Electron IPC Bridge

As a **developer integrating Electron with Golang**,
I want **the Electron main process to spawn the Golang binary and communicate via JSON-RPC**,
So that **the UI can request operations from the backend**.

**Acceptance Criteria:**

**Given** the Electron app starts
**When** the main process initializes
**Then** the Golang binary is spawned from `resources/bin/autobmad-{platform}`
**And** stdin/stdout pipes are established for JSON-RPC communication
**And** the preload script exposes a secure `window.api` object

**Given** the renderer calls `window.api.system.ping()`
**When** the request is processed through the IPC bridge
**Then** the response `"pong"` is returned to the renderer
**And** contextIsolation is enabled (NFR-S5)
**And** nodeIntegration is disabled in the renderer

**Given** the Golang process crashes
**When** the main process detects the crash
**Then** an error event is emitted to the renderer
**And** the crash is logged for debugging

---

### Story 1.4: Detect and Validate OpenCode CLI

As a **user launching Auto-BMAD**,
I want **the system to detect my OpenCode CLI installation and verify compatibility**,
So that **I know my environment is correctly configured before starting journeys**.

**Acceptance Criteria:**

**Given** OpenCode CLI is installed and in PATH
**When** Auto-BMAD performs dependency detection
**Then** OpenCode is detected with version displayed
**And** compatibility is verified against minimum version (v0.1.0+)
**And** the detection result is returned via JSON-RPC `project.detectDependencies`

**Given** OpenCode CLI is NOT in PATH
**When** Auto-BMAD performs dependency detection
**Then** a clear error message is displayed: "OpenCode CLI not found"
**And** installation instructions are suggested

**Given** OpenCode version is incompatible
**When** Auto-BMAD performs dependency detection
**Then** a warning is displayed with detected version and minimum required version

---

### Story 1.5: Detect OpenCode Profiles

As a **user with multiple OpenCode profiles configured**,
I want **Auto-BMAD to detect my available profiles from ~/.bash_aliases**,
So that **I can select which profile to use for my project**.

**Acceptance Criteria:**

**Given** the user has OpenCode profiles defined in `~/.bash_aliases`
**When** Auto-BMAD parses the configuration
**Then** all profile names are extracted and listed
**And** the default profile is identified
**And** the list is returned via JSON-RPC `opencode.getProfiles`

**Given** a profile is misconfigured or has missing credentials
**When** Auto-BMAD validates the profile
**Then** a clear warning message is displayed (NFR-I5)
**And** the profile is marked as "unavailable" in the UI

**Given** no profiles are found
**When** Auto-BMAD parses the configuration
**Then** the system uses the global OpenCode default
**And** a message indicates "Using default OpenCode configuration"

---

### Story 1.6: Detect Git Installation

As a **user launching Auto-BMAD**,
I want **the system to verify Git is installed**,
So that **checkpoint commits and rollback features will work correctly**.

**Acceptance Criteria:**

**Given** Git is installed (version 2.0+)
**When** Auto-BMAD performs dependency detection
**Then** Git version is detected and displayed
**And** the detection result is included in `project.detectDependencies` response

**Given** Git is NOT installed
**When** Auto-BMAD performs dependency detection
**Then** a clear error message is displayed: "Git not found"
**And** the error blocks journey start with explanation: "Git is required for checkpoint safety"

**Given** Git credentials are configured
**When** Auto-BMAD detects Git
**Then** Auto-BMAD does NOT access or store credentials (NFR-S1, NFR-I7)
**And** Git operations use system credentials transparently

---

### Story 1.7: Detect BMAD Project Structure

As a **user opening a project folder**,
I want **Auto-BMAD to detect if it's a BMAD project and identify greenfield vs brownfield**,
So that **the system knows what journey options are available**.

**Acceptance Criteria:**

**Given** a folder with `_bmad/` folder present
**When** Auto-BMAD scans the project
**Then** the project is identified as a BMAD project
**And** BMAD version is read from `_bmad/_config/manifest.yaml`
**And** compatibility is verified (6.0.0+)

**Given** a folder with `_bmad/` but NO `_bmad-output/` folder
**When** Auto-BMAD scans the project
**Then** project type is identified as "greenfield"
**And** full journey options are available

**Given** a folder with both `_bmad/` and `_bmad-output/` with artifacts
**When** Auto-BMAD scans the project
**Then** project type is identified as "brownfield"
**And** existing artifacts are listed for context
**And** partial journey options are available

**Given** a folder without `_bmad/` folder
**When** Auto-BMAD scans the project
**Then** a message indicates "Not a BMAD project"
**And** user is prompted to run `bmad-init` first

---

### Story 1.8: Project Folder Selection UI

As a **user starting Auto-BMAD**,
I want **to select a project folder and optionally describe its context**,
So that **Auto-BMAD knows which project to work with**.

**Acceptance Criteria:**

**Given** the user opens Auto-BMAD
**When** no project is selected
**Then** a project selection screen is displayed
**And** a "Select Project Folder" button opens the native file picker
**And** recently opened projects are listed for quick access

**Given** the user selects a valid BMAD project folder
**When** the selection is confirmed
**Then** project detection runs automatically
**And** results are displayed (BMAD version, greenfield/brownfield, dependencies)
**And** the project is saved to recent projects list

**Given** a brownfield project is selected
**When** the user wants to provide additional context
**Then** a text area allows manual project description (FR48)
**And** the description is saved with the project configuration

---

### Story 1.9: Settings Persistence

As a **user who has configured Auto-BMAD**,
I want **my settings to be saved and restored across sessions**,
So that **I don't have to reconfigure every time I open the app**.

**Acceptance Criteria:**

**Given** the user changes settings (retry limits, notification preferences)
**When** the changes are made
**Then** settings are saved to `_bmad-output/.autobmad/config.json`
**And** the save completes within 1 second (NFR-P6)

**Given** the user restarts Auto-BMAD
**When** the app initializes
**Then** previous settings are restored automatically
**And** the last-used project folder is remembered
**And** the last-used OpenCode profile per project is restored

**Given** no settings file exists
**When** the app initializes
**Then** sensible defaults are applied
**And** a new settings file is created on first change

---

### Story 1.10: Network Status Detection

As a **user running journeys with cloud AI providers**,
I want **Auto-BMAD to monitor network connectivity and warn me of issues**,
So that **I understand why journeys might fail and can take action**.

**Acceptance Criteria:**

**Given** the user is online
**When** Auto-BMAD checks network status
**Then** status indicator shows "connected"
**And** cloud provider workflows are enabled

**Given** the user goes offline
**When** network status changes
**Then** a non-intrusive warning notification appears
**And** status indicator updates to "offline"
**And** if a journey is running, user is asked: "Pause journey or continue anyway?"

**Given** the user is using a local provider (e.g., Ollama)
**When** network status is offline
**Then** journeys can still execute
**And** UI indicates "Using local provider - offline operation available"

---

## Epic 2: Journey Planning & Visualization

**Goal:** Users can create journeys by selecting destinations, view calculated routes before starting, and understand exactly what will happen when they click "Start."

### Story 2.1: Journey State Management

As a **developer building the journey system**,
I want **a well-defined journey state structure persisted to the filesystem**,
So that **journeys can be tracked, resumed, and recovered from crashes**.

**Acceptance Criteria:**

**Given** a journey is created
**When** the state is initialized
**Then** a `journey-state.json` file is created in `_bmad-output/.autobmad/`
**And** the state includes: journeyId, destination, route, currentStep, status, createdAt

**Given** the journey state changes
**When** an update occurs
**Then** the file is updated within 1 second (NFR-P6)
**And** the update is atomic (no partial writes)

**Given** Auto-BMAD restarts after a crash
**When** the app initializes
**Then** existing journey state is read from filesystem
**And** the journey can be resumed from last saved state (NFR-R3)

**Given** no journey is active
**When** the state file is checked
**Then** status shows "idle" or file doesn't exist

---

### Story 2.2: BMAD Workflow Discovery

As a **user planning a journey**,
I want **Auto-BMAD to discover all available BMAD workflows in my project**,
So that **I can see what destinations are possible**.

**Acceptance Criteria:**

**Given** a valid BMAD project is loaded
**When** workflow discovery runs
**Then** all workflows from `_bmad/bmm/workflows/` are discovered
**And** workflow metadata (name, description, phase) is extracted
**And** workflows are organized by phase (Analysis, Planning, Solutioning, Implementation)

**Given** the workflow manifest exists at `_bmad/_config/workflow-manifest.csv`
**When** discovery runs
**Then** the manifest is used to identify available workflows
**And** workflow dependencies are extracted

**Given** workflows are discovered
**When** the list is returned via JSON-RPC `workflow.list`
**Then** each workflow includes: id, name, phase, description, inputs, outputs

---

### Story 2.3: Destination Selection UI

As a **user starting a new journey**,
I want **to select my destination from a clear list of options**,
So that **I can tell Auto-BMAD where I want to end up**.

**Acceptance Criteria:**

**Given** the user clicks "New Journey"
**When** the destination selection screen appears
**Then** available destinations are displayed grouped by phase
**And** each destination shows: name, description, typical duration
**And** keyboard navigation works (NFR-U1)

**Given** a greenfield project
**When** destinations are displayed
**Then** all destinations are available (Brainstorming → Architecture → Implementation)

**Given** a brownfield project with existing artifacts
**When** destinations are displayed
**Then** completed artifacts are shown as "already exists"
**And** partial journeys are suggested (e.g., "Continue from PRD to Architecture")

**Given** the user selects a destination
**When** the selection is confirmed
**Then** route calculation is triggered automatically

---

### Story 2.4: Route Calculation Engine

As a **user who selected a destination**,
I want **Auto-BMAD to calculate the optimal route of workflows**,
So that **I know exactly what steps will be executed**.

**Acceptance Criteria:**

**Given** destination "Architecture Document" is selected
**When** route calculation runs on a greenfield project
**Then** the route includes: Brainstorming → Product Brief → PRD → Architecture
**And** each step shows: workflow name, estimated duration, required inputs

**Given** destination "Architecture Document" is selected
**When** route calculation runs on a brownfield with existing PRD
**Then** the route skips completed steps: Architecture only
**And** existing artifacts are marked as "input from previous work"

**Given** a destination requires intermediate artifacts
**When** those artifacts don't exist
**Then** the route automatically includes prerequisite workflows
**And** the full dependency chain is resolved

**Given** route calculation completes
**When** the result is returned via JSON-RPC `journey.calculateRoute`
**Then** response includes: steps[], totalEstimatedTime, requiredInputs[]

---

### Story 2.5: Route Preview Display

As a **user about to start a journey**,
I want **to see a clear preview of the calculated route**,
So that **I have "walk away confidence" knowing exactly what will happen**.

**Acceptance Criteria:**

**Given** a route is calculated
**When** the route preview is displayed
**Then** all steps are shown in a visual timeline/list
**And** each step shows: number, workflow name, phase badge, estimated time
**And** the total journey duration estimate is displayed

**Given** the route preview is shown
**When** the user reviews it
**Then** a message confirms: "I'll notify you if I need input (yellow flag)"
**And** current notification settings are displayed
**And** "Start Journey" button is prominently displayed

**Given** the user wants to customize the route
**When** they expand "Advanced Options"
**Then** they can: skip optional steps, adjust retry limits, select OpenCode profile
**And** the route preview updates to reflect changes

---

### Story 2.6: Journey Dashboard Layout

As a **user monitoring journeys**,
I want **a kanban-style dashboard showing all journey states**,
So that **I can see everything at a glance**.

**Acceptance Criteria:**

**Given** the user opens the dashboard
**When** the main view loads
**Then** journeys are displayed as cards in status columns:
- Running (🔵)
- Yellow Flag (🟡)
- Paused (⏸️)
- Complete (🟢)
- Failed (🔴)
**And** the layout is responsive and uses shadcn/ui Card components

**Given** a journey card is displayed
**When** the user views it
**Then** the card shows: journey name, destination, current step, elapsed time, status icon
**And** clicking the card expands to show more details

**Given** no journeys exist
**When** the dashboard loads
**Then** an empty state is displayed with "Start Your First Journey" call-to-action

**Given** the UI updates
**When** journey state changes
**Then** cards move between columns in near real-time (< 100ms, NFR-P3)
**And** transitions are smooth and non-jarring

---

### Story 2.7: Journey Status Indicators

As a **user glancing at Auto-BMAD**,
I want **instant visual comprehension of journey status via colors**,
So that **I know in < 1 second if I can walk away or need to act**.

**Acceptance Criteria:**

**Given** a journey is running normally
**When** the status is displayed
**Then** the indicator is 🔵 Blue with label "Running"
**And** the color meets WCAG AA contrast requirements (NFR-U2)

**Given** a journey needs user input
**When** the status changes to yellow flag
**Then** the indicator is 🟡 Amber with label "Needs Input"
**And** the card pulses subtly to draw attention

**Given** a journey completed successfully
**When** the status is displayed
**Then** the indicator is 🟢 Green with label "Complete"
**And** "Walk away safely" is communicated visually

**Given** a journey failed
**When** the status is displayed
**Then** the indicator is 🔴 Red with label "Failed"
**And** "Click for details" is indicated

**Given** any status indicator
**When** screen reader accesses it
**Then** appropriate ARIA labels are present (e.g., "Journey status: Running")

---

### Story 2.8: Artifact Viewer

As a **user reviewing journey progress**,
I want **to view generated artifacts within Auto-BMAD**,
So that **I can check outputs without leaving the app**.

**Acceptance Criteria:**

**Given** a journey has generated artifacts
**When** the user clicks "View Artifacts"
**Then** a list of artifacts in `_bmad-output/planning-artifacts/` is displayed
**And** each artifact shows: filename, type, last modified, size

**Given** the user selects an artifact
**When** it opens in the viewer
**Then** markdown content is rendered with proper formatting
**And** frontmatter is displayed in a collapsible header
**And** a "Open in Editor" button launches the system default editor

**Given** an artifact is being generated
**When** the journey is running
**Then** a "Generating..." indicator shows on the artifact
**And** the viewer auto-refreshes when generation completes

---

### Story 2.9: Start Journey Action

As a **user ready to begin a journey**,
I want **to start execution with a single click**,
So that **I can begin autonomous operation immediately**.

**Acceptance Criteria:**

**Given** a route is calculated and previewed
**When** the user clicks "Start Journey"
**Then** the journey state transitions to "running"
**And** the first workflow step begins execution
**And** the UI navigates to the journey detail view

**Given** the journey starts
**When** the initial state is saved
**Then** a checkpoint is created before execution begins (NFR-R6)
**And** a Git branch `autobmad/journey-{timestamp}` is created

**Given** the journey is starting
**When** the UI updates
**Then** a confirmation message appears: "Journey started. You can close this window."
**And** the message includes: "I'll notify you if I need input"

**Given** prerequisites are missing (e.g., OpenCode not found)
**When** the user clicks "Start Journey"
**Then** the start is blocked with a clear error message
**And** resolution steps are suggested

---

## Epic 3: Autonomous Workflow Execution

**Goal:** Users can start journeys and have Auto-BMAD execute BMAD workflows autonomously via OpenCode, with real-time output capture and checkpoint safety.

### Story 3.1: OpenCode Process Spawning

As a **user starting a workflow step**,
I want **Auto-BMAD to spawn OpenCode with the correct configuration**,
So that **the AI can execute the BMAD workflow**.

**Acceptance Criteria:**

**Given** a workflow step is ready to execute
**When** OpenCode is spawned
**Then** the process starts within 2 seconds (NFR-P4)
**And** the command uses: `opencode -p "{prompt}" --non-interactive`
**And** the selected profile is used if configured

**Given** OpenCode spawning fails
**When** the error is detected
**Then** the failure is captured with error message
**And** the step status is set to "failed"
**And** retry logic is available (Epic 5)

**Given** the spawned process is running
**When** state is checked
**Then** process PID is tracked for monitoring
**And** spawn timestamp is recorded

---

### Story 3.2: Output Streaming & Capture

As a **user monitoring a running workflow**,
I want **to see OpenCode output in real-time**,
So that **I know what's happening without waiting for completion**.

**Acceptance Criteria:**

**Given** OpenCode is running
**When** output is produced on stdout/stderr
**Then** output is captured 100% completely (NFR-I6)
**And** chunks are streamed to the UI via JSON-RPC events
**And** output is also logged to `_bmad-output/.autobmad/logs/journey-{id}.log`

**Given** the UI receives output chunks
**When** they are displayed
**Then** the output viewer updates in near real-time (< 100ms)
**And** stdout and stderr are visually distinguished
**And** the viewer auto-scrolls to latest output

**Given** output is very large
**When** it exceeds display limits
**Then** older output is truncated in memory (keep last 10,000 lines)
**And** full output remains in log file

---

### Story 3.3: Process Monitoring

As a **system maintaining journey health**,
I want **to continuously monitor OpenCode process status**,
So that **crashes are detected quickly and handled gracefully**.

**Acceptance Criteria:**

**Given** an OpenCode process is running
**When** monitoring checks status
**Then** process health is verified every 1 second
**And** memory/CPU usage is tracked (optional display)

**Given** the OpenCode process crashes unexpectedly
**When** the crash is detected
**Then** detection occurs within 5 seconds (NFR-R4)
**And** a `step.failed` event is emitted with crash details
**And** partial output before crash is preserved

**Given** the OpenCode process becomes unresponsive
**When** no output is received for timeout period (configurable, default 5 min)
**Then** a timeout warning is raised
**And** user is asked: "Process may be stuck. Continue waiting or terminate?"

---

### Story 3.4: Workflow Completion Detection

As a **system managing workflow steps**,
I want **to accurately detect when a workflow completes successfully or fails**,
So that **the journey can proceed or handle failure appropriately**.

**Acceptance Criteria:**

**Given** OpenCode process exits with code 0
**When** completion is detected
**Then** step status is set to "completed"
**And** output is parsed for any warnings or notes
**And** a `step.completed` event is emitted

**Given** OpenCode process exits with non-zero code
**When** completion is detected
**Then** step status is set to "failed"
**And** stderr is captured as error message
**And** a `step.failed` event is emitted with `retryable: true`

**Given** completion detection runs
**When** determining success
**Then** there is 0% false success rate (FR28)
**And** if uncertain, the step is marked as "needs review" rather than "success"

---

### Story 3.5: Sequential Workflow Execution

As a **user running a multi-step journey**,
I want **workflows to execute in sequence according to the route**,
So that **each step builds on the previous outputs**.

**Acceptance Criteria:**

**Given** a journey with multiple steps in the route
**When** execution proceeds
**Then** steps execute in order (step 1, then step 2, etc.)
**And** each step waits for the previous to complete before starting

**Given** step N completes successfully
**When** the orchestrator proceeds
**Then** step N+1 begins automatically
**And** previous step outputs are available as context for next step
**And** journey progress updates (currentStep, percentage)

**Given** step N fails
**When** failure is detected
**Then** execution pauses (does not proceed to N+1)
**And** retry/recovery options are presented (Epic 5)

**Given** all steps complete successfully
**When** the journey finishes
**Then** journey status is set to "complete"
**And** a `journey.completed` event is emitted
**And** a final checkpoint commit is created

---

### Story 3.6: Git Checkpoint Commits

As a **user who might need to rollback**,
I want **Auto-BMAD to create Git checkpoints at key moments**,
So that **I can recover from any point with zero data loss**.

**Acceptance Criteria:**

**Given** a step completes (success or failure)
**When** the checkpoint is triggered
**Then** all changed files in `_bmad-output/` are committed
**And** commit message follows: `[AutoBMAD] Step {n}: {name} - {status}`
**And** commit is on branch `autobmad/journey-{timestamp}`

**Given** the user pauses a journey
**When** pause is triggered
**Then** a checkpoint commit is created before pausing
**And** commit message includes: `[AutoBMAD] Journey paused by user`

**Given** a yellow flag is raised (Epic 4)
**When** the flag is raised
**Then** a checkpoint commit is created before prompting user

**Given** the journey completes
**When** final checkpoint is created
**Then** an optional tag is created: `autobmad-complete-{timestamp}`
**And** user is prompted: "Merge to main branch?" (optional)

**Given** checkpoint commit is needed
**When** the commit is created
**Then** commit succeeds 100% of the time (NFR reliability)
**And** if Git fails, error is logged and user is warned

---

### Story 3.7: Pause Journey

As a **user who needs to step away**,
I want **to pause an active journey with state preserved**,
So that **I can resume later without losing progress**.

**Acceptance Criteria:**

**Given** a journey is running
**When** the user clicks "Pause"
**Then** the current step is allowed to complete (or immediately if between steps)
**And** journey status changes to "paused"
**And** a checkpoint commit is created

**Given** a journey is paused
**When** the state is saved
**Then** all context is preserved: currentStep, accumulated feedback, output logs
**And** the journey can be resumed from this exact point

**Given** a paused journey
**When** the UI displays it
**Then** a "Resume" button is prominently displayed
**And** time paused is tracked and displayed

---

### Story 3.8: Cancel Journey

As a **user who wants to stop a journey entirely**,
I want **to cancel the journey and optionally keep or discard outputs**,
So that **I have full control over my project**.

**Acceptance Criteria:**

**Given** a journey is running or paused
**When** the user clicks "Cancel"
**Then** a confirmation dialog appears: "Cancel this journey?"
**And** options are presented: "Keep generated artifacts" / "Rollback to start"

**Given** the user confirms cancellation with "Keep artifacts"
**When** cancellation proceeds
**Then** the current process is terminated
**And** a final checkpoint commit is created
**And** journey status is set to "cancelled"

**Given** the user confirms cancellation with "Rollback to start"
**When** cancellation proceeds
**Then** `git reset --hard` to the initial journey checkpoint
**And** all generated artifacts are removed
**And** journey status is set to "cancelled"

---

### Story 3.9: Process Termination

As a **system managing OpenCode processes**,
I want **to cleanly terminate processes when needed**,
So that **resources are freed and no zombie processes remain**.

**Acceptance Criteria:**

**Given** a process needs to be terminated (pause, cancel, timeout)
**When** termination is triggered
**Then** SIGTERM is sent first, allowing graceful shutdown
**And** if process doesn't exit within 5 seconds, SIGKILL is sent
**And** process resources are cleaned up

**Given** the Electron app is closing
**When** graceful shutdown begins
**Then** any running OpenCode process is terminated
**And** current journey state is saved before exit (NFR-R7)
**And** shutdown completes within 5 seconds

---

### Story 3.10: BMAD Context Passing

As a **workflow step being executed**,
I want **to receive proper BMAD workflow and agent context**,
So that **OpenCode knows which workflow to run and how**.

**Acceptance Criteria:**

**Given** a workflow step is about to execute
**When** the OpenCode command is constructed
**Then** the prompt includes: workflow path, agent to invoke, any user context
**And** the working directory is set to the project root

**Given** previous steps generated artifacts
**When** context is passed to current step
**Then** references to previous artifacts are included
**And** the prompt indicates: "Continue from {previous_artifact}"

**Given** the user provided feedback (Epic 4)
**When** context is passed
**Then** accumulated feedback is included in the prompt
**And** feedback is formatted for the AI to understand

---

## Epic 4: Feedback System & Adaptive Direction

**Goal:** Users can respond to yellow flags (input requests) quickly, provide feedback during execution, and have the system adapt direction based on their input.

### Story 4.1: Yellow Flag Detection

As a **system executing workflows**,
I want **to detect when a workflow needs user input**,
So that **the user can be notified and provide feedback**.

**Acceptance Criteria:**

**Given** OpenCode is executing a workflow
**When** the workflow reaches a decision point requiring user input
**Then** the system detects the yellow flag condition
**And** execution pauses at that point
**And** a `yellowFlag.raised` event is emitted

**Given** a yellow flag is detected
**When** the event is emitted
**Then** the event includes: journeyId, stepIndex, reason, options[], context
**And** the reason is human-readable (e.g., "Choose database type")
**And** options are extracted if the workflow provides choices

**Given** the workflow output contains a BMAD decision prompt
**When** parsing the output
**Then** the prompt text is extracted for display to user
**And** any suggested options are parsed

---

### Story 4.2: Desktop Notifications

As a **user doing other work while a journey runs**,
I want **to receive desktop notifications for yellow flags and failures**,
So that **I know when Auto-BMAD needs my attention without constantly checking**.

**Acceptance Criteria:**

**Given** a yellow flag is raised
**When** the notification is triggered
**Then** a desktop notification appears (non-intrusive, doesn't steal focus)
**And** notification shows: "Journey needs input: {brief context}"
**And** clicking the notification brings Auto-BMAD to foreground

**Given** a journey fails
**When** the notification is triggered
**Then** a desktop notification appears with: "Journey failed: {brief reason}"
**And** notification persists until dismissed or clicked

**Given** notification preferences are configured
**When** the user has disabled notifications
**Then** no desktop notifications are sent
**And** in-app indicators still update (badge, status)

**Given** the notification appears
**When** the user doesn't respond
**Then** the notification persists (no auto-dismiss)
**And** a badge appears on the app icon/taskbar

---

### Story 4.3: Yellow Flag Modal UI

As a **user responding to a yellow flag**,
I want **a clear, context-rich modal that helps me provide feedback quickly**,
So that **I can respond in < 30 seconds and get back to my other work**.

**Acceptance Criteria:**

**Given** a yellow flag is raised
**When** the user opens Auto-BMAD
**Then** the yellow flag modal appears automatically (or is prominently indicated)
**And** the modal shows: question/prompt, context from workflow, any suggested options

**Given** the yellow flag modal is displayed
**When** the user reviews it
**Then** the input mechanism matches the question type:
- Text field for open-ended questions
- Radio buttons for single-choice options
- Checkboxes for multi-select options
- File picker if file input is needed
**And** keyboard shortcuts work (Enter to submit, Esc to defer)

**Given** the modal shows context
**When** the user needs more information
**Then** an expandable "More Context" section shows recent output
**And** the current step and journey progress are visible

---

### Story 4.4: Feedback Submission

As a **user providing feedback to a yellow flag**,
I want **instant acknowledgment and journey continuation**,
So that **I feel confident my input was received and can return to other work**.

**Acceptance Criteria:**

**Given** the user enters feedback in the yellow flag modal
**When** they click "Submit" or press Enter
**Then** acknowledgment appears within 500ms (NFR-P5)
**And** the modal closes
**And** journey status changes from "yellow flag" to "running"

**Given** feedback is submitted
**When** the journey resumes
**Then** the feedback is passed to the workflow via OpenCode
**And** the feedback is logged in journey history
**And** a checkpoint commit is created before resuming

**Given** the user wants to defer responding
**When** they click "Remind me later" or close the modal
**Then** the yellow flag remains active
**And** the journey stays paused
**And** a reminder notification is scheduled (configurable interval)

---

### Story 4.5: Journey Resume from Checkpoint

As a **user returning to a paused journey**,
I want **to resume from exactly where I left off**,
So that **no work is lost and I can continue seamlessly**.

**Acceptance Criteria:**

**Given** a journey was paused (manually or due to yellow flag)
**When** the user clicks "Resume"
**Then** the journey resumes from the last saved checkpoint
**And** the current step continues or next step begins
**And** all accumulated context (feedback, outputs) is preserved

**Given** the app crashed during a journey
**When** Auto-BMAD restarts
**Then** the journey state is read from `journey-state.json`
**And** a prompt asks: "Resume interrupted journey from step {N}?"
**And** clicking "Resume" continues from last checkpoint

**Given** resume is triggered
**When** the journey restarts
**Then** previous output logs are still accessible
**And** the UI shows "Resuming from step {N}..."

---

### Story 4.6: Adaptive Direction

As a **user who provided feedback that changes the journey direction**,
I want **the system to recalculate the route if needed**,
So that **my input actually influences what happens next**.

**Acceptance Criteria:**

**Given** user feedback indicates a direction change (e.g., "Skip testing phase")
**When** the feedback is processed
**Then** the route is recalculated to reflect the new direction
**And** the UI shows "Route updated" with the new plan
**And** the change is logged in journey history

**Given** user feedback is incorporated
**When** subsequent steps execute
**Then** the feedback context is included in prompts
**And** AI agents can reference the user's stated preferences

**Given** route recalculation occurs
**When** the new route is displayed
**Then** removed steps are shown as "skipped by user"
**And** added steps are highlighted as "added based on feedback"

---

## Epic 5: Auto-Retry & Failure Recovery

**Goal:** Users experience intelligent auto-retry on failures, accumulating feedback across attempts, with honest failure reporting when escalation is needed, and clear recovery options.

### Story 5.1: Auto-Retry on Failure

As a **user running a journey**,
I want **failed steps to automatically retry without my intervention**,
So that **transient failures are handled and I don't have to babysit**.

**Acceptance Criteria:**

**Given** a workflow step fails with a retryable error
**When** auto-retry is triggered
**Then** the step restarts automatically after a brief delay (5 seconds)
**And** the retry is logged but user is NOT notified for first retry
**And** the UI shows "Retrying... (attempt 2 of {max})"

**Given** auto-retry succeeds
**When** the step completes
**Then** the journey continues to the next step
**And** the retry is recorded in history but not highlighted
**And** user experience feels seamless (NFR-R5)

**Given** an OpenCode process crashes
**When** crash is detected (within 5 seconds per NFR-R4)
**Then** auto-retry is triggered automatically
**And** the crash is logged with any captured output

---

### Story 5.2: Feedback Accumulation

As a **system retrying a failed step**,
I want **to accumulate feedback from previous attempts**,
So that **each retry is smarter and more likely to succeed**.

**Acceptance Criteria:**

**Given** a step fails and is retried
**When** the retry executes
**Then** the prompt includes: "Previous attempt failed with: {error summary}"
**And** any partial output from previous attempt is referenced

**Given** multiple retries have occurred
**When** attempt N executes
**Then** feedback from ALL previous attempts (1 to N-1) is accumulated
**And** the prompt indicates patterns: "This step has failed {N-1} times"

**Given** user provided feedback during retries
**When** that feedback is available
**Then** it's included in subsequent retry prompts
**And** the AI can learn from user guidance

---

### Story 5.3: Retry Configuration

As a **user who wants control over retry behavior**,
I want **to configure maximum retry attempts**,
So that **I can balance persistence with fast failure**.

**Acceptance Criteria:**

**Given** the settings panel
**When** the user configures retry settings
**Then** options include: max retries (default: 3), retry delay (default: 5s)
**And** settings can be configured globally or per-project

**Given** a per-journey override is needed
**When** the user starts a journey with "Advanced Options"
**Then** retry limit can be overridden for that journey
**And** the override is displayed in the route preview

**Given** retry limit is set to 0
**When** a step fails
**Then** no auto-retry occurs
**And** failure is immediately escalated to user

---

### Story 5.4: Retry Escalation

As a **user whose step has exceeded retry limits**,
I want **to be notified and given options**,
So that **I can decide how to proceed**.

**Acceptance Criteria:**

**Given** a step has failed {max_retries} times
**When** the threshold is exceeded
**Then** a notification is sent: "Step failed after {N} attempts"
**And** journey status changes to "needs intervention"
**And** the escalation modal appears when user opens Auto-BMAD

**Given** escalation occurs
**When** the user reviews options
**Then** choices include:
- "Retry again" (resets counter, tries once more)
- "Retry with feedback" (user provides guidance)
- "Skip this step" (if optional)
- "Abort journey"
**And** each option has clear consequences explained

**Given** the user chooses "Retry with feedback"
**When** feedback is provided
**Then** the feedback is added to accumulated context
**And** retry proceeds with enhanced prompt

---

### Story 5.5: Retry Tracking

As a **user reviewing journey progress**,
I want **to see retry counts and history**,
So that **I understand how the journey handled difficulties**.

**Acceptance Criteria:**

**Given** a step has been retried
**When** the step details are viewed
**Then** retry count is displayed: "Completed after 2 retries"
**And** each attempt is logged with: timestamp, duration, outcome, error (if failed)

**Given** the journey completes
**When** the summary is displayed
**Then** total retries handled is shown
**And** steps that required retries are highlighted
**And** user can expand to see retry details

**Given** retry tracking is active
**When** retries occur
**Then** tracking is stored in journey-state.json
**And** persists across app restarts

---

### Story 5.6: Accurate Failure Detection

As a **user trusting Auto-BMAD's reports**,
I want **0% false success rate**,
So that **when Auto-BMAD says something succeeded, it actually did**.

**Acceptance Criteria:**

**Given** a workflow step completes
**When** success is determined
**Then** multiple signals are checked: exit code, output patterns, artifact existence
**And** if ANY signal indicates failure, step is marked failed (not success)

**Given** output is ambiguous
**When** success cannot be confidently determined
**Then** step is marked "needs review" (not success)
**And** user is asked to verify the output

**Given** an artifact should have been generated
**When** completion is checked
**Then** artifact existence is verified
**And** if missing, step is marked failed despite exit code 0

**Given** the system errs
**When** in doubt
**Then** it errs toward "failed" not "success"
**And** honest uncertainty is communicated

---

### Story 5.7: Failure Explanation

As a **user whose journey failed**,
I want **clear, human-readable explanations**,
So that **I understand what went wrong without digging through logs**.

**Acceptance Criteria:**

**Given** a step fails
**When** the failure explanation is generated
**Then** it includes: what was being attempted, what went wrong, when it happened
**And** language is human-readable (not raw error codes)

**Given** a failure timeline is requested
**When** the timeline is displayed
**Then** events leading to failure are shown chronologically
**And** each event shows: timestamp, action, result
**And** the failure point is clearly marked

**Given** technical details are needed
**When** the user expands "Technical Details"
**Then** raw error output, exit codes, and stack traces are available
**And** "Copy to Clipboard" allows easy sharing for support

---

### Story 5.8: Root Cause Categorization

As a **user understanding a failure**,
I want **the failure categorized by root cause**,
So that **I know if it's something I can fix or a system issue**.

**Acceptance Criteria:**

**Given** a failure occurs
**When** root cause analysis runs
**Then** failure is categorized into one of:
- **Network**: Connection issues, API timeouts
- **System**: Out of memory, disk full, process crash
- **Logic**: Workflow logic error, invalid input
- **Authentication**: API key issues (in OpenCode)
- **Timeout**: Process exceeded time limit
- **Unknown**: Cannot determine cause

**Given** a category is assigned
**When** displayed to user
**Then** category is shown with icon and color
**And** category-specific guidance is provided
**And** common fixes for that category are suggested

---

### Story 5.9: Saved/Lost State Reporting

As a **user whose journey failed**,
I want **to know what was saved before failure and what was lost**,
So that **I can assess the damage and decide next steps**.

**Acceptance Criteria:**

**Given** a failure occurs mid-journey
**When** the failure report is generated
**Then** "Saved" section lists: completed steps, generated artifacts, checkpoints
**And** "Lost" section lists: current step progress, uncommitted changes

**Given** checkpoints were created
**When** the saved state is reported
**Then** last checkpoint commit SHA is shown
**And** "Rollback to checkpoint" option is available

**Given** partial artifacts were generated
**When** reporting lost state
**Then** partial files are identified
**And** user is warned: "This artifact may be incomplete"

---

### Story 5.10: Failure Details View

As a **user viewing the dashboard**,
I want **to access detailed failure information easily**,
So that **I can investigate and decide on recovery**.

**Acceptance Criteria:**

**Given** a journey has failed
**When** the journey card shows 🔴 status
**Then** clicking the card opens failure details view
**And** summary shows: step that failed, attempt count, root cause category

**Given** failure details view is open
**When** the user reviews it
**Then** sections include:
- Summary (what failed, when)
- Explanation (human-readable)
- Timeline (events leading to failure)
- Saved/Lost (what's preserved)
- Recovery Options (next steps)
- Technical Details (expandable)

---

### Story 5.11: Recovery Option Suggestions

As a **user deciding how to recover from failure**,
I want **clear, actionable recovery options**,
So that **I can get back on track quickly**.

**Acceptance Criteria:**

**Given** a failure has occurred
**When** recovery options are displayed
**Then** options are contextual to the failure type:
- Network failure: "Check connection and retry"
- Timeout: "Retry with extended timeout"
- Logic error: "Provide feedback and retry"
- System error: "Check system resources"
**And** each option has a button to execute it

**Given** the user selects a recovery option
**When** the action executes
**Then** the journey transitions appropriately:
- Retry: restarts the failed step
- Skip: moves to next step (if allowed)
- Rollback: reverts to checkpoint
- Abort: ends the journey
**And** the choice is logged in journey history

**Given** "Retry with feedback" is selected
**When** the feedback modal opens
**Then** the failure context is pre-populated
**And** user can add guidance for the retry

---

## Epic 6: Dashboard & Journey History

**Goal:** Users can view complete journey history, completion summaries, retry statistics, and access past artifacts — the full observability experience.

### Story 6.1: Journey History Persistence

As a **user who has completed journeys**,
I want **journey history persisted across sessions**,
So that **I can review past work anytime**.

**Acceptance Criteria:**

**Given** a journey completes (success, failure, or cancelled)
**When** the journey ends
**Then** journey metadata is saved to `_bmad-output/.autobmad/history.json`
**And** metadata includes: journeyId, destination, status, duration, startedAt, completedAt, retryCount

**Given** Auto-BMAD restarts
**When** the history is loaded
**Then** all past journeys are available
**And** journeys are sorted by most recent first

**Given** history grows large (100+ journeys)
**When** storage is managed
**Then** old journeys are archived (metadata kept, logs compressed)
**And** user can configure retention period

---

### Story 6.2: Journey History List View

As a **user reviewing past work**,
I want **a clear list of all past journeys**,
So that **I can find and review any previous journey**.

**Acceptance Criteria:**

**Given** the user navigates to History view
**When** the view loads
**Then** past journeys are displayed as a list/table
**And** each row shows: date, destination, status icon, duration, retry count

**Given** a journey in the list
**When** the user clicks it
**Then** the journey detail view opens
**And** all information from that journey is accessible

**Given** the history list is displayed
**When** sorted by default
**Then** most recent journeys appear first
**And** user can re-sort by: date, destination, status, duration

---

### Story 6.3: Retry Progress Display

As a **user reviewing a journey (past or current)**,
I want **to see retry progress details**,
So that **I understand how difficulties were handled**.

**Acceptance Criteria:**

**Given** a journey with retry attempts
**When** retry progress is displayed
**Then** current attempt number is shown: "Attempt 2 of 3"
**And** a progress indicator shows retry status

**Given** retry details are expanded
**When** the user views them
**Then** each attempt shows: timestamp, duration, outcome, error (if failed)
**And** accumulated feedback is shown for each attempt

**Given** a journey completed after retries
**When** the summary is viewed
**Then** "Completed after {N} retries" is highlighted
**And** the system's persistence is celebrated (builds trust)

---

### Story 6.4: Journey Completion Summary

As a **user whose journey completed**,
I want **a summary of what was accomplished**,
So that **I appreciate the value delivered**.

**Acceptance Criteria:**

**Given** a journey completes successfully
**When** the completion summary is displayed
**Then** it shows:
- Destination reached
- Total duration
- Steps completed
- Artifacts generated (with links)
- Retries handled (if any)
- Time saved estimate (based on manual workflow time)

**Given** the completion summary is shown
**When** the user reviews it
**Then** a celebratory but calm visual indicates success
**And** "View Artifacts" button provides quick access
**And** "Start New Journey" allows immediate next action

**Given** journey stats are displayed
**When** time saved is calculated
**Then** estimate is based on average manual BMAD workflow time
**And** cumulative time saved across all journeys is tracked

---

### Story 6.5: History Search & Filter

As a **user with many past journeys**,
I want **to search and filter history**,
So that **I can find specific journeys quickly**.

**Acceptance Criteria:**

**Given** the history view
**When** the user types in the search box
**Then** journeys are filtered by: destination name, project name, date
**And** search is fuzzy (typo-tolerant)

**Given** filter options are available
**When** the user applies filters
**Then** filters include:
- Status: Complete, Failed, Cancelled
- Date range: Today, This week, This month, Custom
- Destination type: PRD, Architecture, Stories, etc.
**And** filters can be combined

**Given** search results are displayed
**When** no results match
**Then** "No journeys found" message appears
**And** suggestions to broaden search are offered

---

### Story 6.6: Command Palette

As a **power user who prefers keyboard navigation**,
I want **a command palette for quick actions**,
So that **I can work faster without mouse navigation**.

**Acceptance Criteria:**

**Given** the user presses Cmd/Ctrl+K
**When** the command palette opens
**Then** a search input appears with recent commands
**And** typing filters available commands
**And** arrow keys navigate, Enter executes

**Given** the command palette is open
**When** the user searches
**Then** available commands include:
- "New Journey to [destination]"
- "Open Project..."
- "View History"
- "Settings"
- "Pause Journey"
- "Resume Journey"
**And** commands are fuzzy-matched

**Given** a command is selected
**When** Enter is pressed
**Then** the command executes immediately
**And** the palette closes
**And** focus moves to appropriate UI element

