---
project_name: auto-bmad
author: Hafiz
date: 2026-01-20
version: 1.1
status: complete
classification:
  projectType: desktop_app
  projectCategory: developer_tool
  platformPriority: [linux, macos]
  domain: developer_tooling
  projectContext: greenfield
workflowType: prd
workflow: edit
inputDocuments:
  - _bmad-output/planning-artifacts/product-brief-auto-bmad-2026-01-20.md
  - _bmad-output/planning-artifacts/prd-validation-report.md
stepsCompleted:
  - step-e-01-discovery
  - step-e-02-review
  - step-e-03-edit
lastEdited: 2026-01-20
editHistory:
  - date: 2026-01-20
    changes: "Systematic improvements based on validation report: (1) Restructured User Journeys from narrative to structured format (47 density violations eliminated, ~60% compression), (2) Enhanced all 24 NFRs with measurement methods and context (closes measurability gaps), (3) Resolved Phase Jump-Back traceability gap by moving to Post-MVP"
---

# Product Requirements Document - Auto-BMAD

**Author:** Hafiz
**Date:** 2026-01-20
**Version:** 1.0

---

## Executive Summary

### Product Vision

**Auto-BMAD** is an autonomous BMAD workflow orchestration platform that transforms how solo developers interact with AI-assisted development workflows.

**Core Philosophy:**
> *"Set the destination, the system drives the journey. Flow like water."*

### The Problem

Solo developers running BMAD workflows face "death by a thousand interruptions" — constant babysitting of AI workflows, manual transitions between phases, and cognitive overload from tracking progress across multiple projects. Running 100 workflows/day becomes unsustainable.

### The Solution

Auto-BMAD is a **thin orchestration layer** (Electron + Golang) that leverages OpenCode CLI as the AI execution backbone. Users set a destination (e.g., "Architecture Document"), and Auto-BMAD autonomously orchestrates the entire journey through BMAD workflows.

### Key Differentiators

| Differentiator | Description |
|----------------|-------------|
| **Journey-Driven Architecture** | Destination-based, not step-based orchestration |
| **Trust Through Honesty** | 0% false success rate; honest failure > fake success |
| **Auto-Retry with Learning** | Failed workflows retry with accumulated feedback |
| **Artifacts as Memory** | BMAD outputs ARE the state; no separate memory layer |
| **Thin Layer Philosophy** | Orchestration only; OpenCode does the AI work |

### Target Users

**Primary:** Solo developers (like Hafiz) managing multiple projects, running 100+ BMAD workflows/day
**Secondary:** BMAD community adopters seeking autonomous workflow execution

### MVP Scope (5 Features)

1. **OpenCode CLI Integration** — Spawn, monitor, capture processes
2. **Journey Engine** — Route calculation and sequential execution
3. **Auto-Retry with Feedback** — Learn from failures
4. **Honest Failure Reporting** — 0% false success
5. **Basic Dashboard** — Progress visualization

### Platforms

- **MVP:** Linux, macOS (x86_64 + ARM64)
- **Post-MVP:** Windows

---

## Success Criteria

### User Success

**Primary User Success Indicators:**

| Indicator | Description | Target |
|-----------|-------------|--------|
| **Flow State Achievement** | Sustained periods of uninterrupted deep work | Minimum interruption time; 30+ minute sessions |
| **Cognitive Offload** | Mental burden transferred to system | 90% reduction in context switches |
| **Break-ability** | Ability to walk away mid-journey | Workflow continues or waits gracefully |
| **Adaptive Direction** | System readjusts based on feedback | Auto-redirect without manual re-triggering |

**User Satisfaction Measurement:**
- Post-journey satisfaction surveys
- "Would you trust Auto-BMAD with this journey again?" (Yes/No)
- Net Promoter Score for autonomous execution experience
- Qualitative feedback on "flow like water" experience

**User Success Moment:**
> *"I gave feedback once, walked away, came back to either completed work or an honest explanation. I didn't have to babysit."*

### Business Success

**North Star:**
> *Intention Fulfillment* - Success is NOT vanity metrics (stars, downloads). Success is fulfilling the original intention: "Stop babysitting. Flow like water."

**Business Objectives:**

| Objective | Success Criteria | Timeframe |
|-----------|------------------|-----------|
| **Personal productivity transformation** | Reclaim hours from workflow babysitting | MVP |
| **Sustainable workflow volume** | Handle 100 workflows/day without burnout | Post-MVP |
| **Trust establishment** | Confidence to delegate entire phases | MVP + 1 month |
| **Open-source contribution** | Enable BMAD community adoption | V2 |

**Anti-Metrics (Explicitly NOT Optimized):**
- GitHub stars
- User count
- Feature count
- "Looks successful" rate

### Technical Success

**Philosophy:** Keep technical metrics **minimal for MVP**. Focus on autonomous flow, not optimization.

| Metric | MVP Target | Rationale |
|--------|------------|-----------|
| **Startup Time** | < 5 seconds | Reasonable, not optimized |
| **Memory Usage** | < 500MB idle | Acceptable for always-running daemon |
| **OpenCode Process Reliability** | 99% successful spawns | Core functionality |
| **Git Checkpoint Reliability** | 100% commits succeed | Safety-critical |
| **Crash Recovery** | Resume from last checkpoint | State preservation |

**Technical Success Philosophy:**
> *"We optimize for autonomous flow and trust, not performance benchmarks. Technical metrics exist only to ensure the system doesn't get in the way."*

### Measurable Outcomes

**Intention-Based Success Measurement:**

| Scenario | User Intent | Outcome | Classification |
|----------|-------------|---------|----------------|
| Manual trigger | Intended manual | Manual execution | **Success** (user choice) |
| Auto trigger | Intended auto | Auto execution | **Success** |
| Auto trigger | Intended auto | Required manual | **Failure** |

**Core Metric:** *"Did the user have to manually intervene when they intended autonomous execution?"*

**Quantitative KPIs:**

| KPI | Baseline | MVP Target | Post-MVP Target |
|-----|----------|------------|-----------------|
| **Unintended manual interventions** | 100% (current state) | < 20% | < 5% |
| **Complete journey without intervention** | 0 | 1+ journeys | 10+ journeys/day |
| **False success rate** | N/A | **0%** | **0%** |
| **Failure honesty rate** | N/A | **100%** | **100%** |
| **Context switch reduction** | 0% | 90% | 95% |

**Trust Development Timeline:**

| Milestone | Trust Level | Indicator |
|-----------|-------------|-----------|
| **First successful journey** | Initial trust | "It actually works" |
| **First honest failure report** | Reinforced trust | "It tells me the truth" |
| **X hours of hands-off operation** | Established trust | "I can rely on it" |
| **Recovery from unexpected failure** | Deep trust | "It handles problems gracefully" |

**Trust Principle:**
> *"Trust = Honesty, Not Correctness. A system that fails honestly is more trustworthy than one that succeeds mysteriously."*

---

## Product Scope

### MVP - Minimum Viable Product

**Core Capabilities (P0):**

| Feature | Success Criteria |
|---------|------------------|
| **OpenCode CLI Integration** | Successfully spawn, monitor, and capture OpenCode output |
| **Journey Engine** | Execute 1 complete journey (e.g., Analysis → Architecture) |
| **BMAD Detection** | Auto-detect `_bmad/` and `_bmad-output/` folders |
| **Greenfield Support** | Full journey from scratch |
| **Brownfield Support** | Partial journey on existing projects |
| **Emergency Stop** | One-click halt, no confirmation dialogs |
| **Honest Failure Reporting** | 0% false success, 100% failure transparency |
| **Basic Dashboard** | Journey setup, progress visualization, intervention panel |
| **Git Safety (Flexible)** | Checkpoint commits, rollback capability |
| **Adaptive Redirection** | Accept feedback, auto-adjust direction |

**Platform Support (MVP):**
- Linux
- macOS
- Windows (Post-MVP)

**MVP Success Definition:**
> *"One user (Hafiz) can run 1 complete journey on both greenfield and brownfield projects without unintended manual intervention, with honest failure reporting when things go wrong."*

### Growth Features (Post-MVP / V2)

| Feature | Description | Priority |
|---------|-------------|----------|
| **Agent Negotiation Protocol** | Agents communicate and clarify autonomously | P2 |
| **Preference Learning** | System learns user's style over time | P2 |
| **Advanced Observability** | Decision audit trails, confidence heatmaps | P2 |
| **Cost Optimization** | AI-suggested token/cost optimizations | P2 |
| **Windows Support** | Full cross-platform coverage | P1 |

### Vision (Future / V3+)

| Feature | Description |
|---------|-------------|
| **Multi-Project Dashboard** | Single view across all active journeys |
| **Team Collaboration** | Shared journeys with role-based access |
| **Custom Workflow Builder** | Create BMAD workflows through UI |
| **Plugin Architecture** | Community extensions |

**Ultimate Vision:**
> *"Auto-BMAD becomes invisible infrastructure. Developers set intentions, production-ready artifacts appear—with complete transparency about how they got there."*

---

## User Journeys

### Journey 1: The Solo Polymath — Greenfield Success Path

**Persona:** Hafiz, Solo Developer  
**Context:** Starting a brand new project from scratch

**Scenario:**  
Solo developer starts new project, runs `bmad-init`, opens Auto-BMAD to create Architecture Document from scratch with zero existing artifacts.

**Actions:**
1. Create new project folder, run `bmad-init` (creates `_bmad/` folder)
2. Open Auto-BMAD → Dashboard detects Greenfield (no `_bmad-output/`)
3. Select "New Journey" → Set destination: Architecture Document
4. Auto-BMAD calculates route: Brainstorming → Product Brief → PRD → Architecture
5. Confirm and start journey
6. Auto-BMAD orchestrates: Spawn OpenCode with brainstorming-coach → save artifact → checkpoint commit
7. Product Brief workflow completes → checkpoint
8. PRD workflow requests input on success criteria → Dashboard turns yellow with notification
9. User provides feedback in 30 seconds → Auto-BMAD continues
10. All workflows complete → Dashboard green → Four artifacts in `_bmad-output/`
11. User reviews Architecture Document → Provides feedback: "Use SQLite instead of PostgreSQL for MVP"
12. Auto-BMAD adjusts Architecture Document in place (no restart)

**Outcome:**  
Complete journey from idea to Architecture Document in 4 hours without manual workflow babysitting. User attended meetings while Auto-BMAD executed autonomously.

**Requirements:**  
FR44-FR47 (Project detection), FR1-FR4 (Journey builder, route calculation), FR15-FR17, FR20-FR21 (OpenCode execution), FR9-FR10 (Feedback absorption), FR36-FR41 (Dashboard visualization)

---

### Journey 2: The Solo Polymath — Brownfield Discovery

**Persona:** Hafiz, Solo Developer  
**Context:** Existing project needs new feature development

**Scenario:**  
Developer returns to 6-month-old REST API project (50+ endpoints) to add webhook support. Needs project context refresh and pattern-aware development.

**Actions:**
1. Open Auto-BMAD on existing project folder
2. Dashboard detects Brownfield: `_bmad/` folder found, `_bmad-output/` folder with artifacts, codebase present
3. Auto-BMAD executes Brownfield Discovery Protocol:
   - **Phase 1 Detection:** Identify project structure (Language: Go, Framework: Chi, Database: PostgreSQL)
   - **Phase 2 Scan:** Codebase inventory (147 files, 50+ endpoints, 23 models)
   - **Phase 3 Examine:** Pattern analysis (Repository pattern, Service layer, Handler layer)
   - **Phase 4 Analyze:** Existing artifacts status (Last PRD: 6 months old)
4. Dashboard presents Project Snapshot with full context
5. User selects "Add Feature" → Describes: "Add webhook support"
6. Auto-BMAD suggests context-aware journey using existing patterns
7. Journey executes with full project understanding → New code follows existing patterns

**Outcome:**  
Context-aware feature development with seamless codebase integration. Auto-BMAD understands project architecture better than developer remembers.

**Requirements:**  
FR44-FR48 (Project detection, manual context description for MVP), FR1-FR4 (Journey planning), FR20-FR21 (Execution). Post-MVP: Automated brownfield scanner.

---

### Journey 3: The New Adopter — Onboarding

**Persona:** Alex, BMAD Community Member  
**Context:** Discovered Auto-BMAD on GitHub, wants to try it

**Scenario:**  
New user downloads Auto-BMAD, installs on macOS, runs first autonomous journey to build trust and validate setup.

**Actions:**
1. Download Electron app for macOS
2. Install via drag to Applications (one-click installation)
3. First launch triggers setup flow:
   - OpenCode Detection → Found and compatible
   - OpenCode Profiles → Lists available profiles for selection
   - Project Selection → Browse to project with `_bmad/` folder
   - BMAD Detection → Version compatible, ready to start
4. Auto-BMAD offers Guided First Journey: "Create a Product Brief"
5. User describes idea in 2 sentences
6. Dashboard shows real-time progress
7. User steps away (makes coffee)
8. Returns to Journey Complete notification
9. Review Product Brief artifact

**Outcome:**  
Trust established through successful first autonomous journey. Zero manual interventions required. User confident to delegate workflows.

**Requirements:**  
FR11-FR14 (OpenCode detection, profiles), FR45-FR46 (BMAD detection), FR1-FR4 (Journey execution), FR36-FR41 (Dashboard visualization). Post-MVP: Guided onboarding flow.

---

### Journey 4: The Troubleshooter — Failure Recovery

**Persona:** Hafiz, Solo Developer  
**Context:** A journey failed mid-execution

**Scenario:**  
User runs ambitious journey overnight. Journey fails at Architecture phase due to system resource constraints. User needs honest failure reporting and graceful recovery options.

**Actions:**
1. Dashboard shows red status:
   - Journey Failed at Phase: Solutioning
   - Workflow: create-architecture
   - Error: OpenCode process terminated unexpectedly
   - Last checkpoint: commit abc123
2. User clicks "Investigate"
3. Auto-BMAD presents Failure Analysis Panel:
   - **Timeline:** Step-by-step events leading to failure
   - **Root Cause:** Out of memory (OpenCode processing large PRD)
   - **What Was Saved:** Complete artifacts before failure, partial architecture (40%)
   - **Honest Assessment:** "Journey failed due to system resource constraints, not workflow logic. Partial architecture is salvageable."
4. Recovery Options presented:
   - Retry from failure point (resume with existing context)
   - Retry from phase start (restart architecture workflow)
   - Retry from previous phase (go back to PRD)
   - Abort journey (keep artifacts, start fresh later)
5. User selects "Retry from failure point"
6. Auto-BMAD restores checkpoint state and resumes
7. Journey completes → Dashboard green → All artifacts present

**Outcome:**  
Graceful recovery from honest failure. Trust reinforced through transparent reporting and successful retry. Zero data loss.

**Requirements:**  
FR28-FR35 (Failure detection, reporting, analysis, root cause, saved state), FR22-FR27 (Auto-retry, feedback accumulation, escalation), FR8 (Checkpoint restoration). NFR-R1 (Zero data loss).

---

### Journey Requirements Summary

| Capability | Greenfield | Brownfield | Onboarding | Recovery |
|------------|------------|------------|------------|----------|
| **Project Detection** | Greenfield | Brownfield scan | BMAD detection | - |
| **Journey Builder** | ✅ | ✅ | Guided | - |
| **Route Calculation** | ✅ | Context-aware | Simple | - |
| **Background Execution** | ✅ | ✅ | ✅ | Resume |
| **Checkpoint Commits** | ✅ | ✅ | ✅ | Restore |
| **Yellow Flag (Input)** | ✅ | ✅ | - | - |
| **Feedback Absorption** | ✅ | ✅ | - | - |
| **Codebase Scanning** | - | ✅ | - | - |
| **Pattern Analysis** | - | ✅ | - | - |
| **Failure Analysis** | - | - | - | ✅ |
| **Retry from Checkpoint** | - | - | - | ✅ |

### Core Capabilities Required

1. **Detection Engine:** Greenfield vs Brownfield, BMAD version, OpenCode compatibility
2. **Brownfield Scanner:** Language, framework, patterns, existing artifacts
3. **Journey Builder:** Destination selection, route calculation, context-aware planning
4. **Execution Engine:** Background operation, OpenCode spawning, checkpoint commits
5. **Feedback System:** Yellow flags, input absorption, direction adjustment
6. **Failure Handling:** Honest reporting, root cause analysis, recovery options
7. **Checkpoint System:** Git commits, state restoration, partial artifact preservation
8. **Progress Dashboard:** Real-time status, journey stats, completion summary
9. **Onboarding Flow:** Guided setup, first journey recommendation, trust building

---

## Innovation & Novel Patterns

### Detected Innovation Areas

**1. Journey-Driven Architecture (New Paradigm)**

Auto-BMAD introduces a fundamental shift from **step-based** to **destination-based** workflow orchestration:

| Traditional Approach | Auto-BMAD Approach |
|---------------------|-------------------|
| "Execute step 1, then step 2, then step 3..." | "Set destination, system calculates route" |
| User manages each transition | System manages entire journey |
| Mental model: Task list | Mental model: Navigation/GPS |

**Novelty Assessment:** Potentially first-of-kind in BMAD/AI workflow space. No known implementations of destination-driven workflow orchestration for AI agent systems.

**Commitment Level:** Full commitment — no fallback to step-based execution. Journey-driven is the core value proposition.

**Validation Approach:**
- Track journey execution like JIRA epics/stories
- Each workflow phase = epic milestone
- Each workflow step = story/task completion
- Success measured by journey completion, not step completion

---

**2. Trust Through Honesty + Auto-Retry Loop (Novel Pattern)**

Auto-BMAD transforms honest failure from a dead-end into a **learning feedback loop**:

```
┌─────────────────────────────────────────────────────────┐
│                 AUTO-RETRY WITH FEEDBACK                 │
├─────────────────────────────────────────────────────────┤
│                                                          │
│   Attempt 1 ──→ Honest Failure ──→ Capture Feedback     │
│                                           │              │
│                                           ▼              │
│   Attempt 2 ──→ Honest Failure ──→ Accumulate Feedback  │
│       (with Attempt 1 feedback)           │              │
│                                           ▼              │
│   Attempt 3 ──→ SUCCESS                                  │
│       (with Attempt 1+2 feedback)                        │
│                                                          │
│   Each retry learns from ALL previous attempts           │
└─────────────────────────────────────────────────────────┘
```

**Configuration Options:**
- Max retry attempts (configurable)
- Feedback accumulation strategy (append all / summarize / latest N)
- Retry triggers (automatic vs. user-approved)
- Escalation threshold (when to stop retrying and ask human)

**Novelty Assessment:** Unique combination of radical honesty (0% false success), persistent retry (don't give up easily), and feedback accumulation (each attempt smarter than last).

---

**3. Artifacts as Memory (Counter-Pattern to AI Bloat)**

While other AI orchestration tools build complex memory systems (vector DBs, conversation history, context windows), Auto-BMAD takes the opposite approach:

| Other Tools | Auto-BMAD |
|-------------|-----------|
| Separate memory layer | BMAD artifacts ARE the memory |
| Complex state management | Files on disk = state |
| Memory synchronization issues | Single source of truth |
| Memory cleanup/pruning needed | Git handles history |

**How It Works:**
- `_bmad-output/` folder contains all artifacts
- Each artifact is a checkpoint of understanding
- Git commits preserve state history
- Resume = read artifacts + continue

**Resume/Recovery Pattern:**
1. Read existing artifacts in `_bmad-output/`
2. Parse frontmatter for workflow state (stepsCompleted, etc.)
3. Load last checkpoint commit
4. Resume from exact failure point
5. No separate "memory reconstruction" needed

**Complementary Solutions (Non-Conflicting):**
- **Journey State File:** `_bmad-output/journey-state.yaml` for active journey tracking
- **Retry History:** Append retry feedback to artifact frontmatter
- **Checkpoint Metadata:** Git commit messages contain journey context

### Market Context & Competitive Landscape

| Tool | Approach | Auto-BMAD Differentiation |
|------|----------|---------------------------|
| **Auto-Claude** | Complex memory, parallel agents, heavy infra | Thin layer, sequential journeys, artifacts as memory |
| **Claude Code** | Interactive, human-in-loop every step | Autonomous with async human input |
| **Cursor/Copilot** | Code completion, not workflow orchestration | Full BMAD workflow orchestration |
| **n8n/Zapier** | Visual step-based automation | Journey-driven, AI-native |

**Unique Position:**
> *"Auto-BMAD is the first BMAD-native, journey-driven autonomous workflow orchestrator that builds trust through honesty rather than success theater."*

### Validation Approach

| Innovation | Validation Method | Success Indicator |
|------------|-------------------|-------------------|
| **Journey-Driven** | Epic/story tracking simulation | Users intuitively set destinations, not steps |
| **Honest + Retry** | Retry-to-success metrics | Failures decrease across retries; eventual success |
| **Artifacts as Memory** | Resume-from-artifact tests | 100% state recovery from files alone |

### Risk Mitigation

| Risk | Mitigation |
|------|------------|
| **Journey-driven too abstract** | Clear journey visualization; show calculated route before execution |
| **Auto-retry becomes infinite loop** | Configurable max retries; escalation to human after threshold |
| **Artifacts insufficient for complex state** | Journey-state.yaml for active tracking; git commits for history |
| **Honest failures discourage users** | Frame as "learning loop"; show retry progress; celebrate eventual success |

---

## Desktop App Specific Requirements

### Project-Type Overview

Auto-BMAD is an **Electron-based desktop application** in the developer tool category (like VS Code, Docker Desktop). It serves as a visual orchestration layer on top of OpenCode CLI and BMAD workflows, providing autonomous journey execution with a minimal, trust-building user interface.

**Core Dependencies:**
- **Electron:** Latest stable version
- **Git:** Required for checkpoint commits and state management
- **OpenCode CLI:** Required for AI agent execution (user-configured via `~/.bash_aliases`)

### Platform Support

| Platform | Architecture | MVP Priority | Notes |
|----------|--------------|--------------|-------|
| **Linux** | x86_64 | P0 | Primary development platform |
| **Linux** | ARM64 | P0 | Raspberry Pi, ARM servers |
| **macOS** | x86_64 (Intel) | P0 | Legacy Mac support |
| **macOS** | ARM64 (Apple Silicon) | P0 | M1/M2/M3 Macs |
| **Windows** | x86_64 | P1 (Post-MVP) | Deferred |
| **Windows** | ARM64 | P2 | Future consideration |

**Build Targets:**
- `.deb` and `.rpm` for Linux
- `.dmg` and `.pkg` for macOS
- `.exe` installer for Windows (Post-MVP)
- AppImage for universal Linux support

### System Integration

| Feature | Behavior | Default | Configurable |
|---------|----------|---------|--------------|
| **System Tray** | Run minimized in system tray | Yes | Yes |
| **Startup Launch** | Launch at system startup | No | Yes (opt-in) |
| **Global Hotkeys** | Quick access shortcuts | None | No |
| **File Associations** | `.bmad`, `.journey` files open in Auto-BMAD | Yes | Yes |
| **Terminal Integration** | Display OpenCode output | System default | Yes (configurable) |

**File Association Details:**
- `.bmad` files → Open project in Auto-BMAD
- `.journey` files → Resume/view journey state
- Double-click in file manager opens Auto-BMAD with context

**Terminal Configuration:**
- Default: System default terminal
- Configurable: User can specify preferred terminal (iTerm2, Alacritty, Kitty, etc.)
- Purpose: View raw OpenCode output when needed

### Update Strategy

**Auto-Update Mechanism:**

| Aspect | Specification |
|--------|---------------|
| **Self-Update** | Yes, Auto-BMAD updates itself |
| **Update Check** | On startup + periodic (configurable interval) |
| **User Control** | Can disable auto-updates entirely |
| **Update Prompt** | Notification with "Update Now" / "Later" / "Skip Version" |

**Release Channels:**

| Channel | Description | Audience |
|---------|-------------|----------|
| **Stable** | Production-ready releases | Default for all users |
| **Beta** | Feature-complete, testing phase | Early adopters |
| **Nightly** | Latest development builds | Contributors, testers |

**OpenCode CLI Updates:**
- Auto-BMAD does NOT manage OpenCode CLI updates
- User is responsible for OpenCode CLI version management
- Auto-BMAD detects OpenCode version and warns if incompatible

### Offline Capabilities

**Network Dependency:**

| Scenario | Behavior |
|----------|----------|
| **Online (normal)** | Full functionality |
| **Offline (cloud providers)** | Cannot execute journeys |
| **Offline (local providers)** | Can execute journeys (e.g., local Ollama) |
| **Connection lost mid-journey** | Warning notification + graceful handling |

**Graceful Degradation:**
- Network issues detected → Warning notification displayed
- Options: Pause Journey / Continue Anyway / Queue for Later
- User informed of potential failure or incomplete results

**Queue for Reconnection:**
- Journey saved to queue when connection lost
- Notification when connection restored
- User can review queue and start/cancel queued journeys

**Local Storage:**
- All artifacts stored locally in `_bmad-output/`
- No cloud storage feature
- Git provides backup/sync if user configures remote

### Dependency Detection

**Startup Checks:**
1. Check Electron environment
2. Detect Git installation (error with guide if not found)
3. Detect OpenCode CLI (check PATH and `~/.bash_aliases`)
4. Verify OpenCode version compatibility (warning if incompatible)
5. Ready for use

**OpenCode Profile Detection:**
- Parse `~/.bash_aliases` for OpenCode profile definitions
- Detect environment variables: `XDG_CONFIG_HOME`, `XDG_DATA_HOME`
- Present available profiles in UI for user selection
- Remember last-used profile per project

### Technical Architecture Considerations

**Electron Architecture:**
- **Renderer Process:** React/Vue/Svelte UI, Dashboard components, Journey visualization
- **Main Process:** Window management, System tray, File associations, Auto-updater
- **Golang Backend:** Journey Engine, OpenCode CLI spawning, Git operations, BMAD detection

**Process Communication:**
- Electron Main ↔ Renderer: IPC (Inter-Process Communication)
- Electron Main ↔ Golang Backend: Embedded binary with stdin/stdout or local HTTP/WebSocket

**Security Considerations:**
- Electron security best practices (contextIsolation, nodeIntegration: false)
- Secure IPC communication
- No sensitive data in renderer process
- Git credentials handled by system Git, not Auto-BMAD

### Implementation Considerations

**Build & Distribution:**
- Use `electron-builder` for cross-platform builds
- Code signing for macOS (required for Gatekeeper)
- Notarization for macOS distribution
- Linux: Provide .deb, .rpm, and AppImage

**Performance:**
- Minimize Electron footprint when in system tray
- Lazy-load UI components
- Golang backend handles heavy lifting
- Efficient IPC to avoid UI blocking

---

## Project Scoping & Phased Development

### MVP Strategy & Philosophy

**MVP Approach:** Learning MVP — Validate that journey-driven autonomous execution produces trustworthy results

**Core Hypothesis to Validate:**
> *"Can Auto-BMAD run BMAD workflows autonomously with results that are correct, accurate, valid, and honest?"*

**MVP Success Criteria:**
- 80%+ journeys complete without manual intervention
- 0% false success (honest about all failures)
- Output artifacts are usable (not garbage)
- User trusts the system enough to walk away

**Resource Requirements:**
- Solo developer (Hafiz)
- Timeline: Focus on validation, not feature completeness

### MVP Feature Set (Phase 1)

**Core User Journeys Supported:**
- Greenfield Success Path (basic)
- Brownfield (manual project description, no auto-scan)
- Failure Recovery (honest reporting, manual retry positioning)
- Onboarding (documentation only)

**Must-Have Capabilities (5 Features):**

| Feature | Description | Success Criteria |
|---------|-------------|------------------|
| **1. OpenCode CLI Integration** | Spawn, monitor, capture OpenCode processes | 99% successful process spawns |
| **2. Journey Engine** | Calculate route from start to destination, execute sequentially | Complete 1 full journey |
| **3. Auto-Retry with Feedback** | On failure, retry with accumulated feedback | Retry success rate improves |
| **4. Honest Failure Reporting** | Never report success when failed | 0% false success rate |
| **5. Basic Dashboard** | Journey setup, progress visualization, failure display | User can see state and history |

**Explicitly NOT in MVP:**

| Feature | Status | Rationale |
|---------|--------|-----------|
| Git Safety | Manual | User uses IDE to manage branches |
| BMAD Detection | Manual | User selects project folder |
| Brownfield Scanner | Manual | User describes project context |
| Emergency Stop | Manual | Kill process via OS |
| Adaptive Redirection | Manual | Restart journey with new params |
| Phase Jump-Back | Manual | Trigger specific workflow manually |
| Onboarding Flow | Docs | README/documentation |
| Queue for Reconnection | Manual | User restarts when online |
| System Tray | Deferred | Window only for MVP |
| File Associations | Deferred | Manual project opening |
| Auto-Updates | Deferred | Manual updates for MVP |
| Multiple Profiles | Deferred | Single profile for MVP |

### Post-MVP Features

**Phase 2 (Growth) — After Core Validation:**

| Feature | Priority | Description |
|---------|----------|-------------|
| **BMAD Detection** | P1 | Auto-detect `_bmad/` folder |
| **Git Safety** | P1 | Checkpoint commits, rollback |
| **Emergency Stop** | P1 | One-click halt |
| **Adaptive Redirection** | P1 | Feedback absorption, direction adjustment |
| **Phase Jump-Back** | P1 | Return to previous phases |
| **Brownfield Scanner** | P2 | Detection → Scan → Examine → Analyze |
| **System Tray** | P2 | Desktop integration |
| **Auto-Updates** | P2 | Self-updating capability |
| **Multiple OpenCode Profiles** | P2 | Profile selection per project |

**Phase 3 (Expansion) — Platform Maturity:**

| Feature | Description |
|---------|-------------|
| **Onboarding Flow** | Guided first journey experience |
| **File Associations** | `.bmad` and `.journey` files |
| **Queue for Reconnection** | Network resilience |
| **Windows Support** | Platform expansion |
| **Agent Negotiation** | Agents communicate autonomously |
| **Preference Learning** | System learns user style |
| **Advanced Observability** | Decision audit trails |

### Risk Mitigation Strategy

**Technical Risk: Autonomous Execution Quality**

| Risk | Mitigation |
|------|------------|
| OpenCode produces garbage output | Auto-retry with feedback; validate artifacts |
| Journey engine sequences incorrectly | Start with simple journeys; expand after validation |
| Feedback accumulation bloats context | Summarize feedback; cap context size |
| Process spawning unreliable | Robust error handling; clear failure reporting |

**Market Risk: Users Don't Trust Autonomy**

| Risk | Mitigation |
|------|------------|
| Users won't walk away | Start with short journeys; build trust incrementally |
| Honest failures feel broken | Frame as "learning loop"; show retry progress |
| "Why not just use Claude Code?" | Emphasize BMAD-native, journey-driven |

**Resource Risk: Solo Developer**

| Risk | Mitigation |
|------|------------|
| Scope creep | Ruthless MVP focus (5 features only) |
| Burnout | Ship fast, validate, iterate |
| Technical blockers | Proven tech stack; minimize unknowns |

### Development Sequence

**Phase 1: MVP (Validate Core) — Weeks 1-10**
- Week 1-2: OpenCode CLI Integration
- Week 3-4: Journey Engine (Basic)
- Week 5-6: Auto-Retry + Honest Failure
- Week 7-8: Basic Dashboard
- Week 9-10: Integration + Validation (10 complete journeys)

**Phase 2: Growth — After Validation**
- BMAD Detection, Git Safety, Emergency Stop
- Adaptive Redirection, Phase Jump-Back
- Brownfield Scanner, System Tray, Auto-Updates

**Phase 3: Expansion — Platform Maturity**
- Windows Support, Advanced Features
- Onboarding Flow, File Associations
- Agent Negotiation, Preference Learning

---

## Functional Requirements

### Journey Management

- **FR1:** User can create a new journey by specifying a destination (target artifact/phase)
- **FR2:** System can calculate a route from current state to destination
- **FR3:** User can view the calculated route before starting execution
- **FR4:** User can start journey execution with one action
- **FR5:** User can pause an active journey
- **FR6:** User can cancel an active journey
- **FR7:** User can view journey history (past journeys and outcomes)
- **FR8:** System can resume a paused journey from last checkpoint
- **FR9:** User can provide feedback during journey execution
- **FR10:** System can adjust journey direction based on user feedback

### OpenCode Integration

- **FR11:** System can detect installed OpenCode CLI
- **FR12:** System can detect OpenCode CLI version and verify compatibility
- **FR13:** System can detect available OpenCode profiles from user configuration
- **FR14:** User can select which OpenCode profile to use for a project
- **FR15:** System can spawn OpenCode CLI processes with specified configuration
- **FR16:** System can monitor active OpenCode processes
- **FR17:** System can capture OpenCode output (stdout, stderr)
- **FR18:** System can terminate OpenCode processes
- **FR19:** System can pass BMAD workflow/agent context to OpenCode

### Execution & Retry

- **FR20:** System can execute BMAD workflows sequentially per journey route
- **FR21:** System can detect workflow completion (success or failure)
- **FR22:** System can automatically retry failed workflows
- **FR23:** System can accumulate feedback from previous retry attempts
- **FR24:** System can pass accumulated feedback to retry attempts
- **FR25:** User can configure maximum retry attempts
- **FR26:** System can escalate to user after exceeding retry threshold
- **FR27:** System can track retry count and history per workflow

### Failure & Reporting

- **FR28:** System can detect workflow failures accurately (0% false success)
- **FR29:** System can report failures with clear explanations
- **FR30:** System can provide failure timeline (events leading to failure)
- **FR31:** System can identify root cause category (system, network, logic, etc.)
- **FR32:** System can report what was saved before failure
- **FR33:** System can report what was lost due to failure
- **FR34:** User can view failure details in dashboard
- **FR35:** System can suggest recovery options based on failure type

### Dashboard & Visualization

- **FR36:** User can view current journey status (idle, running, paused, failed, complete)
- **FR37:** User can view journey progress (current phase, step, percentage)
- **FR38:** User can view active workflow details
- **FR39:** User can view generated artifacts
- **FR40:** User can view retry progress (attempt number, feedback history)
- **FR41:** User can view journey completion summary (duration, artifacts, retries)
- **FR42:** User can access journey history
- **FR43:** System can display notifications for important events

### Project & Configuration

- **FR44:** User can select a project folder to work with
- **FR45:** System can detect if project has `_bmad/` folder
- **FR46:** System can detect if project has `_bmad-output/` folder
- **FR47:** System can detect project type (greenfield vs brownfield)
- **FR48:** User can describe project context manually (for brownfield without auto-scan)
- **FR49:** User can configure Auto-BMAD settings (retry limits, notifications, etc.)
- **FR50:** System can persist user preferences across sessions
- **FR51:** System can detect network connectivity status
- **FR52:** System can warn user of network issues during journey

### FR Priority Mapping

**MVP (Phase 1):** FR1-FR7, FR9, FR11-FR48, FR51-FR52 (45 FRs)
**Post-MVP (Phase 2+):** FR8, FR10, FR49-FR50 (4 FRs)

---

## Non-Functional Requirements

### Performance

| NFR ID | Requirement | Target | Measurement Method | Context |
|--------|-------------|--------|-------------------|---------|
| **NFR-P1** | Application startup time | < 5 seconds | Measured from process launch to dashboard ready state using system timer | Critical for developer flow; measured on reference hardware (4-core CPU, 8GB RAM) |
| **NFR-P2** | Memory usage (idle) | < 500 MB | Measured via system memory profiler (htop/Activity Monitor) after 5 minutes idle | Acceptable for always-running daemon; excludes OpenCode process memory |
| **NFR-P3** | Dashboard UI updates | Near real-time (< 100ms) | Measured from state change event to UI render completion using browser DevTools Performance | Ensures responsive feedback; applies to journey status, progress bar, notifications |
| **NFR-P4** | OpenCode process spawn time | < 2 seconds | Measured from spawn command to process ready state via process manager | Critical for journey flow; includes profile loading and initialization |
| **NFR-P5** | Feedback submission response | < 500ms acknowledgment | Measured from submit action to confirmation display using IPC timing | User perception of responsiveness; does not include processing time |
| **NFR-P6** | Journey state save | < 1 second | Measured from save trigger to file write completion using filesystem monitor | Ensures rapid checkpointing; applies to journey-state.yaml writes |

### Reliability

| NFR ID | Requirement | Target | Measurement Method | Context |
|--------|-------------|--------|-------------------|---------|
| **NFR-R1** | Data loss on crash | **Zero tolerance** — no loss of any progress | Verified via crash simulation tests; all in-flight data must be recoverable from checkpoint | Safety-critical; applies to journey state, artifacts, user feedback |
| **NFR-R2** | Checkpoint frequency | Continuous (after every significant state change) | Measured via checkpoint log timestamps; max 30 seconds between checkpoints | Ensures minimal recovery window; triggered by workflow completion, feedback submission, artifact creation |
| **NFR-R3** | Crash recovery | Resume from last saved state | Verified via recovery tests; 100% state restoration from checkpoint files and git commits | Must restore journey position, accumulated feedback, partial artifacts |
| **NFR-R4** | OpenCode crash detection | < 5 seconds | Measured from process exit to detection notification using health check polling (1-second interval) | Enables rapid failure response; polling interval: 1 second |
| **NFR-R5** | OpenCode auto-retry on crash | Automatic restart and continue | Verified via process manager logs; restart within 10 seconds of detection | Applies to transient failures; escalates to user after 3 consecutive crashes |
| **NFR-R6** | Journey state persistence | Before any risky operation | Verified via checkpoint audit trail; state saved before OpenCode spawn, artifact modification, network operations | Risk mitigation; defines risky operations as: process spawn, file operations, network calls |
| **NFR-R7** | Graceful shutdown | Save all state before exit | Verified via shutdown tests; all state persisted within 2 seconds of exit command | Applies to user-initiated exit, system shutdown; force-kill after 5-second timeout |

**Checkpoint Strategy:**
- Checkpoint after each workflow phase completion
- Checkpoint after each user feedback submission
- Checkpoint before spawning new OpenCode process
- Checkpoint after each artifact creation/modification
- Checkpoint on periodic interval (every 30 seconds minimum)

### Integration

| NFR ID | Requirement | Target | Measurement Method | Context |
|--------|-------------|--------|-------------------|---------|
| **NFR-I1** | OpenCode CLI version | Any installed version (use global default) | Detected via `opencode --version` command on startup | Compatible with OpenCode CLI v0.1.0+; warns if version incompatibility detected |
| **NFR-I2** | OpenCode profiles | Support multiple profiles with load-balancing | Verified via profile rotation tests; round-robin distribution across configured profiles | Detected from `~/.bash_aliases`; Post-MVP feature for load balancing |
| **NFR-I3** | Git version | Standard user Git installation | Detected via `git --version` command on startup; requires Git 2.0+ | Uses system Git; no embedded Git binary |
| **NFR-I4** | BMAD version | Current `_bmad/` in project directory | Detected via `_bmad/_config/manifest.yaml` version field | Compatible with BMAD 6.0.0+; warns if version mismatch |
| **NFR-I5** | Profile misconfiguration | Warn user with clear message | Verified via error handling tests; displays profile name, error type, resolution steps within 1 second of detection | Triggered when profile not found, credentials missing, or spawn fails |
| **NFR-I6** | OpenCode output capture | Complete stdout/stderr capture | Verified via byte-level comparison tests; 100% capture with < 1% data loss measured via stream monitoring | Ensures full output available for retry feedback and debugging |
| **NFR-I7** | Git operations | Use system Git credentials | Verified via integration tests; uses user's configured git config and SSH keys | Auto-BMAD never stores or manages credentials; delegates to system Git |

**OpenCode Profile Load-Balancing:**
- Detect all available OpenCode profiles from `~/.bash_aliases`
- Allow user to select multiple profiles for a journey
- Distribute workflow executions across selected profiles
- Handle profile failures by falling back to other profiles

### Security

| NFR ID | Requirement | Target | Measurement Method | Context |
|--------|-------------|--------|-------------------|---------|
| **NFR-S1** | API key storage | None — Auto-BMAD never stores API keys | Verified via codebase audit and filesystem scans; zero API keys in config, logs, or storage | All API keys managed by OpenCode CLI; Auto-BMAD is credential-agnostic |
| **NFR-S2** | Credential handling | Delegated to OpenCode | Verified via architecture review; Auto-BMAD passes no credentials to OpenCode processes | OpenCode manages all AI provider credentials via its own configuration |
| **NFR-S3** | Update verification | Signature verification required | Verified via update tests; auto-updater rejects unsigned packages using electron-updater signature validation | Uses code signing certificates; prevents man-in-the-middle attacks |
| **NFR-S4** | Code signing | macOS: Signed and notarized | Verified via codesign and spctl tools; `spctl --assess --verbose` returns accepted | Required for macOS Gatekeeper; also applies to Windows Authenticode (Post-MVP) |
| **NFR-S5** | IPC security | Electron best practices (contextIsolation) | Verified via Electron security checklist; contextIsolation: true, nodeIntegration: false in renderer | Prevents renderer process from accessing Node.js APIs; uses secure IPC channels |

### Usability

| NFR ID | Requirement | Target | Measurement Method | Context |
|--------|-------------|--------|-------------------|---------|
| **NFR-U1** | Keyboard navigation | Full keyboard support | Verified via manual keyboard-only testing; 100% of UI actions accessible via keyboard (Tab, Arrow keys, Enter, Esc) | Includes journey setup, navigation, settings, all buttons and inputs |
| **NFR-U2** | Color contrast | WCAG AA minimum | Verified via automated contrast checker tools (e.g., axe DevTools); all text/background pairs ≥ 4.5:1 ratio | Ensures readability for visual impairments; applies to all UI states (normal, hover, focus) |
| **NFR-U3** | Font sizing | Configurable (Post-MVP) | Verified via UI settings tests; font scale range 80%-150% in 10% increments | Deferred to Post-MVP; MVP uses system default font size |
| **NFR-U4** | Screen reader | Basic support (Post-MVP) | Verified via screen reader testing (NVDA/VoiceOver); all interactive elements have ARIA labels | Deferred to Post-MVP; MVP focuses on visual interface |

### NFR Priority Summary

**MVP:** NFR-P1 to NFR-P6, NFR-R1 to NFR-R7, NFR-I1 to NFR-I7, NFR-S1 to NFR-S5, NFR-U1 to NFR-U2 (22 NFRs)
**Post-MVP:** NFR-U3, NFR-U4 (2 NFRs)
