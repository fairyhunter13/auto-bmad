---
stepsCompleted: [1, 2, 3, 4, 5]
inputDocuments:
  - _bmad-output/analysis/brainstorming-session-2026-01-20.md
date: 2026-01-20
author: Hafiz
project_name: auto-bmad
---

# Product Brief: Auto-BMAD

## Executive Summary

**Auto-BMAD** is an autonomous workflow orchestration platform that transforms BMAD from a manually-triggered methodology into a destination-driven autonomous system. Built for solo developers who are tired of babysitting workflows, Auto-BMAD executes entire development phases autonomously while maintaining radical transparency about its actions, decisions, and failures.

The core philosophy: **"Set the destination, the system drives the journey."** Like an autonomous vehicle, Auto-BMAD handles the routine navigation while humans focus on strategic direction—intervening only when genuinely needed.

**Key Value Proposition:** Production-ready autonomous workflow execution with trust built on honesty, not correctness. When Auto-BMAD succeeds, it delivers. When it fails, it explains exactly why—no hidden failures, no fake success.

---

## Core Vision

### Problem Statement

Solo developers using BMAD face **death by a thousand interruptions**:

- **Repetitive approvals** drain time and cognitive energy
- **Manual refixes** create frustrating rework loops
- **Context-switching** destroys deep work and flow state
- **Mental load of tracking** workflow state across sessions
- **Manual triggering** adds friction at every single step

The cumulative effect: developers spend more time orchestrating the methodology than doing meaningful work.

### Problem Impact

Without autonomous orchestration:
- **Productivity loss**: Hours spent on workflow babysitting instead of creative problem-solving
- **Cognitive drain**: Mental energy wasted on tracking state and remembering "where was I?"
- **Workflow abandonment**: BMAD's full power remains untapped because the friction is too high
- **Scaling impossibility**: As work volume increases, manual orchestration becomes unsustainable

### Why Existing Solutions Fall Short

**Auto-Claude and similar tools** have attempted autonomous development, but they fail for BMAD users:

| Problem | Impact |
|---------|--------|
| **Bloatware architecture** | Complicated memory systems and MCP overhead that BMAD doesn't need |
| **High token consumption** | Notoriously expensive with diminishing returns |
| **Not production-ready** | Results require significant human cleanup |
| **Hidden failures** | Systems optimize for "looking successful" over being honest |
| **No BMAD awareness** | Don't understand workflow phases, artifacts, or agent handoffs |

**The key insight:** BMAD already produces artifacts that track state and plan next moves. We don't need a separate memory layer—we need an orchestrator that understands BMAD's native artifact system.

### Proposed Solution

**Auto-BMAD** is a lean, BMAD-native autonomous workflow orchestrator with three core capabilities:

1. **Journey-Driven Execution**
   - Greenfield projects: Run Phase 1 (Analysis) through final phase autonomously
   - Brownfield projects: Run specific phases (e.g., implementation) autonomously
   - User sets destination, system navigates the route

2. **Artifact-Based State Management**
   - No complex memory systems—BMAD artifacts ARE the memory
   - Workflow status tracked in `_bmad-output/` structure
   - Checkpoints via git commits for safe rollback

3. **Radical Honesty Protocol**
   - Trust built on transparency, not success rates
   - Every action logged, every decision explained
   - Failures reported with root cause—never hidden
   - "Honest failure > fake success"

**Tech Stack:**
- **Backend:** Golang (performance, simplicity, single binary)
- **Frontend:** Electron desktop app (visual dashboard)
- **Integration:** OpenCode CLI with multi-account profiles

### Key Differentiators

| Differentiator | Description |
|----------------|-------------|
| **BMAD-Native** | Built specifically for BMAD workflows—not a generic AI wrapper |
| **Simplicity Over Bloat** | No unnecessary memory layers, MCP overhead, or complexity |
| **Honesty Over Correctness** | Trust = transparency about success AND failure |
| **Artifact-Driven** | Uses BMAD's existing artifact system as state management |
| **Production-Ready Focus** | Delivers working results or explains exactly why not |
| **Journey-Driven** | Destination-based thinking, not step-based checklists |

**Unfair Advantage:** Built by an experienced BMAD power user who knows the workflows intimately and feels the pain firsthand.

**Why Now:** BMAD has reached maturity (even in alpha), and increasing work volume makes manual orchestration unsustainable. The time is right to transform BMAD from a methodology into an autonomous system.

---

## Target Users

### Primary Users

#### Persona: "The Solo Polymath" — High-Volume BMAD Power User

**Representative Profile:**
- **Name:** Hafiz (and developers like him)
- **Role:** Solo developer managing multiple projects simultaneously
- **Project Scope:** Ranges from large enterprise systems to small utility tools
- **BMAD Proficiency:** Experienced user who knows the workflows intimately

**Work Context:**
- Triggers up to **100 workflows per day**
- Average workflow duration: ~10 minutes
- Juggles **multiple projects simultaneously**
- Primary interruption source: Higher-priority work demanding attention
- Cannot afford to babysit workflows—needs to trust the system and move on

**Pain Experience:**
- **Current state:** "Hell on earth" — constant workflow babysitting prevents breaks, destroys focus, and creates cognitive overload
- **Manual triggering:** Every workflow requires attention, approval, and context-switching
- **Mental load:** Tracking workflow state across multiple projects simultaneously is exhausting
- **Lost productivity:** Hours spent orchestrating instead of creating

**Desired State:**
- **Vision:** "Flow like water" — natural, effortless workflow execution
- **Autonomy:** Set destination, trust the journey, intervene only when needed
- **Freedom:** Take breaks without workflows grinding to a halt
- **Confidence:** Know that failures will be reported honestly, not hidden

**Success Criteria:**
> *"I want to set a destination, walk away, and come back to either completed work or an honest explanation of what went wrong."*

---

### Secondary Users

#### BMAD Community Adopters

**Profile:**
- Developers who discover Auto-BMAD through the public repository
- Existing BMAD users frustrated with manual workflow orchestration
- May have lower workflow volume than primary persona but same pain points

**Needs:**
- Clear documentation and easy setup
- Works out-of-the-box with standard BMAD installations
- Gradual adoption path (can start with semi-autonomous before full autonomy)

**Value Proposition:**
- Same "autonomy with trust" benefits
- Learns from primary user's workflow patterns and optimizations
- Community-driven improvements and feedback

---

### User Journey

#### Discovery Phase
1. **Trigger:** Developer experiences "workflow babysitting fatigue" with manual BMAD
2. **Search:** Looks for BMAD automation solutions, finds Auto-BMAD repository
3. **Evaluation:** Reads README, sees "autonomy with trust" philosophy, resonates with pain points

#### Onboarding Phase
1. **Installation:** Downloads Electron app, points to existing project with `_bmad/` folder
2. **Detection:** Auto-BMAD auto-detects BMAD installation and available workflows
3. **First Run:** Starts with semi-autonomous mode to build confidence
4. **Configuration:** Sets cost boundaries, checkpoint preferences, notification thresholds

#### Core Usage Phase
1. **Destination Setting:** "Run Enterprise Method from brainstorm to architecture"
2. **Monitoring:** Glances at dashboard occasionally—green means good
3. **Intervention:** Only when flagged (yellow) or critical (red)
4. **Review:** Async review of completed work, provides feedback that improves future runs

#### "Aha!" Moment
> *"I set up three project journeys before lunch, came back to two completed PRDs and one honest failure report explaining exactly what went wrong. I fixed the issue in 5 minutes. This used to take me all day."*

#### Long-Term Integration
1. **Trust builds:** System learns preferences, failures decrease
2. **Volume scales:** Handles increasing project load without proportional time investment
3. **Flow achieved:** BMAD becomes invisible infrastructure, not active burden
4. **Advocacy:** Recommends to other BMAD users, contributes feedback

---

## Success Metrics

### North Star: Intention Fulfillment

Success is not measured by vanity metrics (stars, downloads, adoption). Success is measured by whether Auto-BMAD fulfills its original intention:

> *"Stop babysitting BMAD workflows. Set the destination, let the system drive the journey. Achieve autonomy with trust. Flow like water."*

---

### Primary Success Indicators

#### Autonomy Achievement

| Metric | Target | Description |
|--------|--------|-------------|
| **Autonomous execution rate** | 100% | All workflows execute without manual step-by-step triggering |
| **Human-in-the-loop availability** | Always | User can redirect, pause, or intervene at any moment |
| **Hands-off journey completion** | Yes | Set destination → walk away → return to results |

#### Trust & Honesty (Non-Negotiable)

| Metric | Target | Description |
|--------|--------|-------------|
| **False success rate** | 0% | Zero hidden failures — if it says "done," it's truly done |
| **Failure honesty rate** | 100% | Every failure includes accurate root cause diagnosis |
| **Transparency completeness** | 100% | No missing crumbs — full visibility into all actions and decisions |

---

### Intention Fulfillment Indicators

These indicators prove the original pain points are solved:

| Original Pain | Success Indicator | Fulfilled When... |
|---------------|-------------------|-------------------|
| "Can't take breaks" | **Break-ability** | User walks away mid-journey; workflow continues or waits gracefully |
| "Babysitting is hell" | **Intervention rarity** | Interventions are exceptional, not routine |
| "Context-switching kills flow" | **Flow preservation** | Complete deep work sessions without BMAD interruptions |
| "Mental load of tracking" | **Cognitive offload** | System maintains all state; user's mind is free |
| "100 workflows/day manually" | **Volume scalability** | Handle 100 workflows/day with effort of 10 manual workflows |

---

### Business Objectives

Since Auto-BMAD is an open-source tool built for personal productivity:

| Objective | Success Criteria |
|-----------|------------------|
| **Personal productivity transformation** | Reclaim hours previously spent on workflow babysitting |
| **Sustainable workflow volume** | Scale to increasing project load without burnout |
| **Trust establishment** | Confidence to delegate entire phases to autonomous execution |
| **Open-source contribution** | Enable other BMAD users to experience the same liberation |

---

### Key Performance Indicators (KPIs)

#### Quantitative KPIs

| KPI | Baseline | Target | Timeframe |
|-----|----------|--------|-----------|
| **Daily workflows without intervention** | 0 | 100 | MVP |
| **Average journey completion (hands-off)** | N/A | Full phase completion | MVP |
| **False success incidents** | N/A | 0 | Always |
| **Time spent on workflow orchestration** | Hours/day | Minutes/day | Post-MVP |

#### Qualitative KPIs

| KPI | Indicator |
|-----|-----------|
| **"Flow like water" achieved** | User reports sustained flow states during BMAD work |
| **"Autonomy with trust" established** | User comfortable delegating multi-hour journeys |
| **"Hell on earth" eliminated** | No more babysitting-induced frustration or burnout |

---

### Anti-Metrics (What We Explicitly Don't Optimize For)

| Anti-Metric | Why We Ignore It |
|-------------|------------------|
| **GitHub stars** | Vanity metric — doesn't indicate intention fulfillment |
| **User count** | Building for quality experience, not mass adoption |
| **Feature count** | Simplicity over bloat — fewer features done right |
| **"Looks successful" rate** | We optimize for honesty, not appearance |

---

## MVP Scope

### Core Architecture Principle

Auto-BMAD is a **thin orchestration layer** that leverages OpenCode CLI as the AI execution backbone. It does not run AI directly — it sequences BMAD workflows and manages journey state while OpenCode handles the actual AI agent execution.

```
┌─────────────────────────────────────────────────────────┐
│                     AUTO-BMAD MVP                        │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────┐    │
│  │            ELECTRON DASHBOARD                    │    │
│  │  • Journey setup (set destination)              │    │
│  │  • Progress visualization                       │    │
│  │  • Emergency stop button                        │    │
│  │  • Intervention flags                           │    │
│  └─────────────────────────────────────────────────┘    │
│                          ↕                               │
│  ┌─────────────────────────────────────────────────┐    │
│  │            GOLANG BACKEND                        │    │
│  │  • Journey Engine (workflow sequencing)         │    │
│  │  • BMAD Detection (_bmad/ awareness)            │    │
│  │  • State Management (artifact-driven)           │    │
│  │  • Git Integration (flexible checkpoints)       │    │
│  │  • Honesty Protocol (failure reporting)         │    │
│  └─────────────────────────────────────────────────┘    │
│                          ↕                               │
│  ┌─────────────────────────────────────────────────┐    │
│  │            OPENCODE CLI (Backbone)               │    │
│  │  • AI agent execution                           │    │
│  │  • Multi-account profiles                       │    │
│  │  • Actual workflow processing                   │    │
│  └─────────────────────────────────────────────────┘    │
│                          ↕                               │
│  ┌─────────────────────────────────────────────────┐    │
│  │            BMAD WORKFLOWS                        │    │
│  │  • _bmad/ folder structure                      │    │
│  │  • Workflow definitions                         │    │
│  │  • Agent configurations                         │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

---

### Core Features (P0 - MVP)

#### 1. OpenCode CLI Integration
- **Priority:** P0 (Foundation)
- **Description:** Connect to OpenCode as the AI execution backbone
- **Requirements:**
  - Detect OpenCode installation and configuration
  - Support multi-account profiles (from `~/.bash_aliases`)
  - Execute BMAD agent commands through OpenCode
  - Capture and parse OpenCode output for state management

#### 2. Journey Engine
- **Priority:** P0 (Core Value)
- **Description:** Destination-driven workflow sequencing
- **Requirements:**
  - User sets destination (e.g., "Run Enterprise Method to Architecture phase")
  - System calculates workflow sequence to reach destination
  - Execute workflows in sequence, passing artifacts between phases
  - Handle greenfield (full journey) and brownfield (partial journey) modes

#### 3. BMAD Detection
- **Priority:** P0 (Usability)
- **Description:** Auto-detect BMAD installation and available workflows
- **Requirements:**
  - Detect `_bmad/` folder in project root
  - Parse workflow manifest and available agents
  - Detect `_bmad-output/` for existing artifacts (brownfield)
  - Present available journeys based on detected configuration

#### 4. Emergency Stop
- **Priority:** P0 (Safety)
- **Description:** Immediate halt capability
- **Requirements:**
  - Always-visible stop button in UI
  - One-click stop with no confirmation dialogs
  - Graceful termination of current OpenCode process
  - State preservation for potential resume

#### 5. Honest Failure Reporting
- **Priority:** P0 (Trust Foundation)
- **Description:** Transparent failure communication with root cause
- **Requirements:**
  - Capture all failure information from OpenCode execution
  - Parse and present root cause analysis
  - Never hide failures or claim false success
  - Provide actionable information for resolution

#### 6. Basic Dashboard (Electron)
- **Priority:** P0 (Usability)
- **Description:** Visual interface for journey management
- **Requirements:**
  - Journey setup: select project, choose destination
  - Progress visualization: current phase, completed steps, ETA
  - Status indicators: green (good), yellow (flagged), red (critical)
  - Intervention panel: view flags, provide input when needed

#### 7. Git Safety (Flexible)
- **Priority:** P0 (Safety)
- **Description:** Checkpoint-based safety with flexible worktree usage
- **Requirements:**
  - Use existing worktree if available
  - Create checkpoint commits at major decision points
  - Enable rollback to previous checkpoints
  - Never force worktree creation (user's choice)

---

### Out of Scope for MVP

| Feature | Reason for Deferral | Post-MVP Priority |
|---------|---------------------|-------------------|
| **Agent Negotiation** | Requires foundation; sequential handoffs sufficient for MVP | P2 |
| **Preference Learning** | Needs usage data to learn from; manual config works | P2 |
| **Advanced Observability** | Basic dashboard sufficient; analytics are nice-to-have | P2 |
| **Multi-Project Parallel** | User can run multiple Auto-BMAD windows; simple solution | P3 |
| **Mandatory Worktree** | Flexibility preferred over enforcement | N/A |
| **Cost Optimization AI** | Basic tracking sufficient; optimization is enhancement | P2 |
| **Complex Memory Systems** | BMAD artifacts ARE the memory; no separate layer needed | N/A |

---

### MVP Success Criteria

#### Functional Success
| Criteria | Validation |
|----------|------------|
| **Journey completion** | Can execute full BMAD phase (e.g., Analysis → PRD) hands-off |
| **OpenCode integration** | Successfully triggers and monitors OpenCode execution |
| **BMAD detection** | Correctly identifies `_bmad/` structure and available workflows |
| **Failure honesty** | All failures reported with accurate root cause; zero hidden failures |
| **Emergency stop** | Halts execution immediately when triggered |

#### User Experience Success
| Criteria | Validation |
|----------|------------|
| **"Walk away" capability** | User can set destination and leave; returns to results |
| **Intervention rarity** | Interventions needed only for genuine decisions, not routine approvals |
| **Trust establishment** | User confident that "done" means done; failures are honest |

#### Technical Success
| Criteria | Validation |
|----------|------------|
| **Stability** | No crashes during normal operation |
| **State preservation** | Can resume after interruption; no lost progress |
| **Checkpoint integrity** | Git commits created; rollback functions correctly |

---

### Future Vision

#### Post-MVP Enhancements (V2)

| Feature | Description |
|---------|-------------|
| **Agent Negotiation Protocol** | Agents communicate and clarify with each other autonomously |
| **Preference Learning** | System learns user's style and preferences over time |
| **Advanced Observability** | Decision audit trails, confidence heatmaps, journey analytics |
| **Cost Optimization** | AI-suggested token/cost optimizations based on patterns |

#### Long-Term Vision (V3+)

| Feature | Description |
|---------|-------------|
| **Multi-Project Dashboard** | Single view across all active project journeys |
| **Team Collaboration** | Shared journeys with role-based access |
| **Custom Workflow Builder** | Create new BMAD workflows through UI |
| **Plugin Architecture** | Extend Auto-BMAD with community plugins |

#### Ultimate Vision

> *Auto-BMAD becomes the invisible infrastructure that makes BMAD feel like magic. Developers set intentions, and production-ready artifacts appear — with complete transparency about how they got there.*
