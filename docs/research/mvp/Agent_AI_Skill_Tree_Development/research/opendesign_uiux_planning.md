# OpenDesign UI/UX Planning — Skill Tree Visualization & Management

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

**Date:** 2026-07-17
**Status:** DESIGN
**Scope:** UI/UX planning for the skill tree visualization and management interface

---

## 1. Overview

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

The skill system needs UI/UX surfaces for:
1. **Skill Tree Visualization** — Interactive DAG visualization
2. **Skill Management** — CRUD operations via GUI
3. **Gap Analysis Dashboard** — Visual gap detection and auto-expand status
4. **Evidence Browser** — Code evidence and pattern exploration
5. **Search Interface** — Hybrid vector+keyword search

---

## 2. Current Surfaces

### 2.1 CLI (`cmd/cli/`)
- Cobra-based command groups: `skills`, `search`, `registry`, `expand`, `learn`, `source`
- Functional but text-only
- Status: ✅ Implemented (basic), ⚠️ Some commands unreachable (G63)

### 2.2 TUI (`cmd/tui/`)
- Bubble Tea-based interactive terminal UI
- Panels: browse, search, tree, registry
- Status: ⚠️ Partially implemented, health indicator disconnected (G63)

### 2.3 REST API (`/api/v1/`)
- Full CRUD spec exists in OpenAPI
- Status: ⚠️ Most endpoints unwired (G01-O3 dead server consolidation)

### 2.4 MCP Server
- 7 tools registered
- Status: ✅ Implemented (basic)

---

## 3. Proposed UI/UX Architecture

### 3.1 Technology Choices

| Surface | Technology | Rationale |
|---------|------------|-----------|
| Web GUI | React + D3.js | Industry standard for DAG visualization |
| Desktop | Tauri (Rust+React) | Lightweight, cross-platform |
| Mobile | Flutter | Already chosen for R3 (Android/iOS/HarmonyOS/Aurora) |
| TUI | Bubble Tea (existing) | Already in stack |

### 3.2 Component Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Web GUI (React)                        │
│                                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │ SkillTree   │  │ GapDashboard│  │ Evidence    │     │
│  │ Visualizer  │  │             │  │ Browser     │     │
│  │ (D3.js)     │  │             │  │             │     │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘     │
│         │                │                │              │
│  ┌──────┴────────────────┴────────────────┴──────┐      │
│  │              API Client Layer                   │      │
│  │  (OpenAPI codegen → TypeScript client)          │      │
│  └────────────────────┬───────────────────────────┘      │
│                       │                                   │
└───────────────────────┼───────────────────────────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │  REST API        │
              │  /api/v1/*       │
              └──────────────────┘
```

---

## 4. Key UI Components

### 4.1 Skill Tree Visualizer

**Purpose:** Interactive DAG visualization of the skill dependency graph.

**Features:**
- Zoom/pan navigation
- Node grouping by domain/complexity
- Edge coloring by relation type (requires=red, extends=blue, recommends=gray, composes=green)
- Click-to-expand node details
- Search/filter by name, tag, domain
- Layout algorithms: force-directed, hierarchical, circular

**Technology:** D3.js force-directed graph with React wrapper

**Data Source:** `GET /api/v1/skills/:name/tree?recursive=true&format=json`

### 4.2 Gap Dashboard

**Purpose:** Visual representation of knowledge gaps and auto-expand progress.

**Features:**
- Gap count by domain
- Auto-expand job status (queued/running/done/failed)
- Coverage heatmap (domain × completeness)
- Manual trigger for gap analysis
- History of auto-expanded skills

**Data Source:** `GET /api/v1/registry/coverage`, `GET /api/v1/registry/missing`

### 4.3 Evidence Browser

**Purpose:** Explore code evidence linked to skills.

**Features:**
- Code snippet viewer with syntax highlighting
- Pattern search across evidence
- Source project/file navigation
- Validation status indicators
- Link to skill details

**Data Source:** `GET /api/v1/evidences/:skill_name`

### 4.4 Skill Editor

**Purpose:** Create and edit skills via GUI.

**Features:**
- Markdown editor with preview
- TOML/JSON toggle for metadata
- Dependency picker (autocomplete from existing skills)
- Resource URL validator
- Real-time validation feedback

**Data Source:** `POST /api/v1/skills`, `PUT /api/v1/skills/:name`

### 4.5 Search Interface

**Purpose:** Hybrid vector+keyword search across skills.

**Features:**
- Search bar with autocomplete
- Filter by domain, complexity, status
- Result ranking (relevance, freshness, coverage)
- Similar skills recommendations
- Search history

**Data Source:** `POST /api/v1/search`

---

## 5. Implementation Phases

### Phase 1: Web GUI Foundation (P0)
1. Set up React project with TypeScript
2. Generate OpenAPI TypeScript client
3. Implement skill tree visualizer (D3.js)
4. Basic search interface
5. Skill detail view

### Phase 2: Management Features (P1)
1. Skill editor (create/update)
2. Gap dashboard
3. Evidence browser
4. Registry overview

### Phase 3: Advanced Features (P2)
1. Real-time updates (WebSocket)
2. Collaborative editing
3. Export/import via GUI
4. Plugin system for custom visualizations

### Phase 4: Mobile/Desktop (P3)
1. Flutter mobile app (Android/iOS)
2. Tauri desktop app
3. HarmonyOS/Aurora ports (per G15 feasibility)

---

## 6. Design Principles

1. **Contract-first**: All UI surfaces consume the same OpenAPI/MCP contract
2. **Progressive disclosure**: Simple view by default, advanced features on demand
3. **Accessibility**: WCAG 2.1 AA compliance
4. **Responsive**: Mobile-first design
5. **Offline-capable**: Local caching for read operations

---

## 7. Dependencies

| Dependency | Status | Blocking |
|------------|--------|----------|
| G01-O3 (API consolidation) | OPEN | All REST-based features |
| G09 (OpenAPI drift) | OPEN | Client codegen |
| G29 (hybrid search) | QUEUED | Search interface quality |
| OpenAPI spec finalized | ⚠️ | Client codegen |

---

## 8. Honest Gaps

1. **Aurora/HarmonyOS** (G15): Flutter-OHOS and omprussia embedder feasibility unproven. Mobile phase blocked until spike completes.
2. **Real-time updates**: WebSocket support requires server-side changes not yet designed.
3. **Collaborative editing**: Requires conflict resolution strategy (CRDT/OT) not yet chosen.
