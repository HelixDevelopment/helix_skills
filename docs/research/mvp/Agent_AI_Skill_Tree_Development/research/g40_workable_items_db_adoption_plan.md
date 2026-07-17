# G40 — Workable-Items DB Adoption Plan

**Revision:** 1
**Last modified:** 2026-07-17T20:30:00Z
**Scope:** ACTIONABLE PLAN. Companion to the design doc at
`g40_workable_items_db_adoption_design.md`.
**Closes:** G40 (HIGH, §11.4.93/§11.4.95), G45 (MEDIUM, status/type
closed-vocabulary), G47 (MEDIUM, id-scheme operator-decision).

---

## 1. Current state

Workable items are scattered across markdown files:

| Source | Items | Format |
|---|---|---|
| `GAPS_AND_RISKS_REGISTER.md` | G01–G137 (136 items) | `### G##` headers with structured fields |
| `CONTINUATION.md` | Session continuation state | Free-form markdown |
| `requests/history.md` | Request history log | Timestamped entries |

**Problems:**
- No programmatic lifecycle enforcement (status transitions are prose, not code)
- No query/filter/sort capability across items
- No automated staleness detection
- Closed-vocabulary (§11.4.15/§11.4.16/§11.4.34) not enforced
- Cross-referencing between items is manual prose

---

## 2. Target state

A SQLite `workable_items.db` as single-source-of-truth, with:
- Every G## item as a row with structured fields
- Programmatic status transitions via closed-vocabulary enum
- Queryable via CLI, REST, MCP, and hooks
- Auto-sync from markdown → DB on read (one-way import)
- DB → markdown export for human-readable docs (one-way export)

---

## 3. Schema design

```sql
-- Core items table
CREATE TABLE workable_items (
    id TEXT PRIMARY KEY,              -- 'G01', 'G40', 'ATM-001', etc.
    type TEXT NOT NULL CHECK(type IN ('Bug', 'Task', 'Feature', 'Meta')),
    severity TEXT CHECK(severity IN ('CRITICAL', 'HIGH', 'MEDIUM', 'LOW', 'N/A')),
    status TEXT NOT NULL CHECK(status IN (
        'FILED', 'IN_PROGRESS', 'QUEUED', 'BLOCKED',
        'FIXED', 'CLOSED', 'SUPERSEDED', 'OPERATOR_BLOCKED',
        'UNCONFIRMED', 'COMPLETED'
    )),
    title TEXT NOT NULL,
    description TEXT,
    evidence TEXT,
    decision TEXT,
    depends_on TEXT,                   -- comma-separated G## ids
    composes_with TEXT,                -- comma-separated G## ids
    category TEXT,                     -- 'security', 'gap', 'inconsistency', etc.
    created_by TEXT DEFAULT 'operator',
    assigned_to TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    resolved_at TEXT,
    resolution_commit TEXT,            -- git commit hash
    constitution_refs TEXT,             -- comma-separated §11.4.NNN refs
    notes TEXT
);

-- Status transition log
CREATE TABLE status_transitions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_id TEXT NOT NULL REFERENCES workable_items(id),
    from_status TEXT,
    to_status TEXT NOT NULL,
    reason TEXT,
    actor TEXT,
    timestamp TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Cross-references
CREATE TABLE item_cross_refs (
    item_id TEXT NOT NULL REFERENCES workable_items(id),
    ref_id TEXT NOT NULL REFERENCES workable_items(id),
    ref_type TEXT NOT NULL CHECK(ref_type IN ('composes', 'depends', 'blocks', 'duplicates', 'supersedes')),
    PRIMARY KEY (item_id, ref_id, ref_type)
);

-- Indexes
CREATE INDEX idx_items_status ON workable_items(status);
CREATE INDEX idx_items_severity ON workable_items(severity);
CREATE INDEX idx_items_type ON workable_items(type);
CREATE INDEX idx_items_category ON workable_items(category);
CREATE INDEX idx_transitions_item ON status_transitions(item_id);
```

---

## 4. Migration strategy

### Phase 1: Schema + Import (source-side, no DB required)

1. Create `internal/workable/db.go` — SQLite connection + schema init
2. Create `internal/workable/models.go` — Go structs matching the schema
3. Create `internal/workable/import.go` — parse `GAPS_AND_RISKS_REGISTER.md`
   `### G##` headers into `WorkableItem` structs
4. Create `internal/workable/export.go` — generate markdown from DB rows
5. Unit tests: golden-good fixture (G01–G05 round-trip), golden-bad fixture
   (malformed entry → clear error), paired §1.1 mutation

### Phase 2: CLI + REST (source-side)

6. Wire `workable-items migrate` CLI command — one-shot markdown → DB import
7. Wire `workable-items list|show|update` CLI commands
8. Wire `POST /api/v1/workable-items/import` + `GET /api/v1/workable-items`
   REST handlers
9. Integration tests: import → query → export round-trip

### Phase 3: Hooks + Automation (requires running system)

10. Wire `PreToolUse` hook: on commit touching `GAPS_AND_RISKS_REGISTER.md`,
    verify DB is in sync (fingerprint check)
11. Wire `PostToolUse` hook: on workable-item status change, update DB + log
    transition
12. Add periodic reconciliation worker: verify DB ↔ markdown consistency

### Phase 4: Decommission Markdown Sources

13. Once DB is proven authoritative, mark markdown files as GENERATED
14. Auto-generate `GAPS_AND_RISKS_REGISTER.md` from DB (like skills catalog)
15. Remove manual editing of markdown sources

---

## 5. Integration points

### CLI

```
workable-items migrate          # one-shot markdown → DB import
workable-items list [--status=X] [--severity=Y] [--type=Z]
workable-items show G40
workable-items update G40 --status=IN_PROGRESS
workable-items export           # DB → markdown
workable-items verify           # check DB ↔ markdown consistency
```

### REST

```
GET  /api/v1/workable-items          # list with filters
GET  /api/v1/workable-items/:id      # single item
POST /api/v1/workable-items/import   # markdown → DB
PUT  /api/v1/workable-items/:id      # update fields
GET  /api/v1/workable-items/export   # DB → markdown
```

### MCP

```
workable_item_get      # read single item
workable_item_list     # list with filters
workable_item_update   # update status/fields
```

### Hooks

- `PreToolUse`: on commit touching register, verify DB sync
- `PostToolUse`: on item status change, update DB + log transition

---

## 6. Lifecycle enforcement

Valid status transitions (closed-vocabulary per §11.4.15/§11.4.16/§11.4.34):

```
FILED → IN_PROGRESS, QUEUED, BLOCKED, OPERATOR_BLOCKED, UNCONFIRMED
IN_PROGRESS → QUEUED, BLOCKED, FIXED, CLOSED
QUEUED → IN_PROGRESS, BLOCKED, OPERATOR_BLOCKED
BLOCKED → IN_PROGRESS, QUEUED
OPERATOR_BLOCKED → IN_PROGRESS, QUEUED
UNCONFIRMED → FILED, CLOSED
FIXED → CLOSED
CLOSED → (terminal)
SUPERSEDED → (terminal)
COMPLETED → (terminal)
```

---

## 7. Phased rollout

| Phase | Scope | Gate | Duration est. |
|---|---|---|---|
| 1 | Schema + import + export | Source-side only | 1–2 sessions |
| 2 | CLI + REST | Source-side only | 1 session |
| 3 | Hooks + automation | Running system | 1 session |
| 4 | Decommission markdown | DB proven | 1 session |

---

## 8. Acceptance criteria

### Phase 1
- [ ] `workable_items.db` schema created + versioned
- [ ] Import parses all 136 G## items from register
- [ ] Export generates equivalent markdown from DB
- [ ] Round-trip test: import → export → diff = 0
- [ ] Golden-good + golden-bad + paired mutation tests

### Phase 2
- [ ] CLI `list|show|update` functional
- [ ] REST `GET|POST|PUT` functional
- [ ] Integration tests pass

### Phase 3
- [ ] PreToolUse hook blocks stale commits
- [ ] PostToolUse hook logs transitions
- [ ] Reconciliation worker detects drift

### Phase 4
- [ ] `GAPS_AND_RISKS_REGISTER.md` is GENERATED from DB
- [ ] Manual edits to markdown are rejected by hook
- [ ] All item management goes through DB/CLI/REST/MCP

---

## 9. Risks and mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| Markdown parser fragility | Import fails on edge cases | Golden fixtures + fuzz testing |
| DB schema evolution | Migration needed on schema change | Versioned migrations (same pattern as skill-system) |
| Markdown ↔ DB drift | Two sources of truth diverge | Fingerprint check + reconciliation worker |
| Operator preference for markdown | DB adoption stalled | Phase 4 is opt-in; markdown remains usable through Phase 3 |
| Constitution engine incompatibility | `constitution/scripts/workable-items/` schema differs | Adapter layer; design doc §1.1 confirms engine builds |

---

## 10. Dependencies

- G40 design doc: `research/g40_workable_items_db_adoption_design.md`
- Constitution engine: `constitution/scripts/workable-items/` (builds cleanly)
- G45 (closed-vocabulary): resolved by schema CHECK constraints
- G47 (id-scheme): resolved by `id TEXT PRIMARY KEY` accepting both `G##` and `ATM-NNN`
