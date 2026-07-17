# Design: CodeGraph Integration (§11.4.78/§11.4.80)

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

**Date:** 2026-07-17
**Status:** DESIGN
**Scope:** Wire CodeGraph MCP server integration for code index and sync automation
**References:** §11.4.78 (CodeGraph MCP), §11.4.79 (code index), §11.4.80 (sync automation)

---

## 1. Overview

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

CodeGraph is an MCP server providing code intelligence (AST parsing, symbol indexing, dependency analysis). The skill system needs to integrate with CodeGraph for:

1. **Code Index** (§11.4.79): Use CodeGraph to index codebases and extract skill-relevant patterns
2. **Sync Automation** (§11.4.80): Keep skill evidence synchronized with code changes
3. **MCP Co-registration** (§11.4.78): Register skill-system MCP alongside CodeGraph in shared plugin configs

---

## 2. Current State

### What Exists
- `.codegraph/` directory with `.gitignore` (empty, placeholder)
- `internal/codeanalysis/` package with tree-sitter stub (G12 — native parsing always errors)
- `internal/codeanalysis/analyzer.go` — regex-based code analysis (functional but limited)
- MCP server in `internal/mcp/` — 7 tools registered
- `learn_from_project` MCP tool — submits project for analysis

### What's Missing
- No CodeGraph MCP client integration
- No code index synchronization
- No automated skill-evidence refresh on code changes
- tree-sitter native parsing (G12 — DESIGN DONE, impl PENDING)

---

## 3. Architecture

### 3.1 Integration Points

```
┌─────────────────────────────────────────────────────────┐
│                    Skill System                          │
│                                                          │
│  ┌──────────┐    ┌──────────┐    ┌──────────────────┐  │
│  │ MCP      │    │ Worker   │    │ CodeAnalysis     │  │
│  │ Server   │    │ Runner   │    │ (tree-sitter)    │  │
│  └────┬─────┘    └────┬─────┘    └────────┬─────────┘  │
│       │               │                    │             │
│       │    ┌──────────┴──────────┐         │             │
│       │    │                     │         │             │
│       │    ▼                     ▼         ▼             │
│  ┌────┴─────────────────────────────────────────────┐   │
│  │              CodeGraph Integration Layer          │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌──────────┐ │   │
│  │  │ IndexClient │  │ SyncManager │  │ PatternDB│ │   │
│  │  └──────┬──────┘  └──────┬──────┘  └────┬─────┘ │   │
│  └─────────┼────────────────┼───────────────┼───────┘   │
│            │                │               │            │
└────────────┼────────────────┼───────────────┼────────────┘
             │                │               │
             ▼                ▼               ▼
    ┌────────────────────────────────────────────────┐
    │              CodeGraph MCP Server               │
    │  (external, co-registered in .claude/settings)  │
    └────────────────────────────────────────────────┘
```

### 3.2 Components

#### 3.2.1 CodeGraph MCP Client (`internal/codegraph/client.go`)
- Thin MCP client that connects to CodeGraph server via stdio/HTTP
- Exposes: `IndexProject(path)`, `QuerySymbols(pattern)`, `GetDependencies(file)`, `WatchChanges(path)`
- Config-driven: `config.toml` `[codegraph]` section

#### 3.2.2 Index Manager (`internal/codegraph/index.go`)
- Manages code index lifecycle: create, update, query, delete
- Maps CodeGraph symbols → skill evidence entries
- Deduplicates evidence by `(source_project, source_file, pattern)`

#### 3.2.3 Sync Manager (`internal/codegraph/sync.go`)
- Watches for code changes (fsnotify or CodeGraph webhook)
- Triggers re-indexing on change
- Updates skill evidence in the DB
- Marks affected skills as `stale` in registry

#### 3.2.4 Pattern Extractor (`internal/codegraph/patterns.go`)
- Extracts skill-relevant patterns from CodeGraph index:
  - API usage patterns
  - Architecture patterns
  - Dependency patterns
  - Configuration patterns
- Maps patterns to existing skills or suggests new skills

---

## 4. Configuration

```toml
[codegraph]
enabled = true
transport = "stdio"  # stdio | http
endpoint = ""        # HTTP endpoint (if transport=http)
sync_interval_seconds = 300
auto_index_on_learn = true
watch_enabled = false  # requires fsnotify
```

---

## 5. Implementation Plan

### Phase 1: MCP Client (P0)
1. Create `internal/codegraph/client.go` — MCP client wrapper
2. Add `[codegraph]` config section
3. Wire into `cmd/server` and `cmd/worker`
4. Test: unit (mock MCP server), integration (real CodeGraph)

### Phase 2: Index Integration (P1)
1. Create `internal/codegraph/index.go` — index management
2. Enhance `learn_from_project` to use CodeGraph for symbol extraction
3. Store CodeGraph-derived evidence in `evidences` table
4. Test: unit, integration (real codebase → evidence)

### Phase 3: Sync Automation (P2)
1. Create `internal/codegraph/sync.go` — change detection
2. Add worker job `JobTypeCodeSync`
3. Wire fsnotify watcher (optional) or periodic polling
4. Test: unit, integration (code change → evidence update)

### Phase 4: Co-registration (P3)
1. Update `.claude/settings.local.json` to register skill-system MCP
2. Create `scripts/codegraph_sync.sh` for manual sync
3. Document in README

---

## 6. Dependencies

| Dependency | Status | Blocking |
|------------|--------|----------|
| G12 (tree-sitter) | DESIGN DONE | Phase 2 (native parsing for better patterns) |
| G29 (hybrid search) | QUEUED | Phase 2 (search CodeGraph-derived evidence) |
| G59 (embedding ingestion) | QUEUED | Phase 2 (embed CodeGraph patterns) |
| CodeGraph MCP server | External | All phases |

---

## 7. Test Plan

| Test Type | Scope | Priority |
|-----------|-------|----------|
| Unit | Client mock, pattern extraction | P0 |
| Integration | Real CodeGraph MCP, real DB | P1 |
| E2E | Code change → evidence update → skill staleness | P2 |
| Contract | MCP tool schema matches CodeGraph API | P1 |
| Mutation | Remove sync trigger → staleness test fails | P2 |
| Security | CodeGraph endpoint validation, SSRF guard | P0 |

---

## 8. Honest Gaps

1. **CodeGraph availability**: Assumes CodeGraph MCP server is installed and running. If not available, integration degrades gracefully (regex-only analysis).
2. **Real-time sync**: fsnotify-based watching is best-effort; periodic polling is the fallback.
3. **Pattern quality**: CodeGraph-derived patterns are only as good as the underlying AST analysis. Complex patterns (e.g., architectural decisions) may require LLM refinement.
