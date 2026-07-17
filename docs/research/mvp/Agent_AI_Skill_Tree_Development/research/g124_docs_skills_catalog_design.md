# G124 — Docs/Skills Catalog: Design Research

**Revision:** 2
**Last modified:** 2026-07-17
**Status:** design — G125 generator CORE is built (at `project/internal/skillscatalog/`); integration layers G126–G132 are queued
**Composes with:** G125–G135, G43 (Docs Chain export), R18 (documentation delivery mandate)
**Upstream generator:** `project/internal/skillscatalog/` (G125, implemented)

---

## §1. Objective

Maintain `docs/skills/` as a complete, structurally-organized (tree: README →
INDEX → by-domain → by-kind → per-skill detail) catalog of EVERY skill in
the System with details + descriptions (name, kind, description,
dependencies/6 relation types, resources), GENERATED from the DB/skill store,
exported md+html+pdf (§11.4.65), and AUTOMATICALLY kept in sync via
§11.4.106 Docs Chain context + §11.4.86 sha256 roster fingerprint +
§11.4.109/§11.4.164 hooks, configurable/triggerable via CLI + REST + all
clients.

Per G124 in GAPS_AND_RISKS_REGISTER.md: the full umbrella item tracking
sub-items G125–G135.

---

## §2. Architecture

```
┌─────────────────────────────────────────────────────┐
│        skillscatalog Go package (G125)              │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │
│  │ load.go  │→│ render.go │→│ generate.go      │   │
│  │ (DB→rec) │  │ (rec→md) │  │ (orchestrate+fp) │   │
│  └──────────┘  └──────────┘  └──────────────────┘   │
│       ↑              ↑              ↑                │
│       │              │              │                │
│  ┌────┴─────┐   ┌────┴──────┐  ┌───┴────────┐       │
│  │Store     │   │model.go   │  │fingerprint │       │
│  │read ops  │   │+ record   │  │.go (sha256 │       │
│  │(ext.)    │   │types      │  │+ sidecar)  │       │
│  └──────────┘   └───────────┘  └────────────┘       │
└──────────────────────┬──────────────────────────────┘
                       │ Generate() / Verify()
                       ▼
┌──────────────────────────────────────────────────────┐
│  docs/skills/  (output tree)                         │
│  ├── README.md          — counts, summary, fingerprint│
│  ├── INDEX.md           — full flat table             │
│  ├── .catalog_fingerprint — sidecar (§11.4.86)        │
│  ├── by-domain/         — page per domain              │
│  ├── by-kind/           — atomic/composite/umbrella   │
│  └── skill/             — one detail page per skill   │
└──────────────────────┬───────────────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────────────┐
│  Integration layers (G126–G135)                      │
│  ├── CLI: skill-system docs skills-catalog generate  │
│  ├── REST: POST /api/v1/skills/catalog/regenerate   │
│  ├── MCP:  tool skill_catalog_status                 │
│  ├── Worker: periodic reconciliation (default 60s)   │
│  ├── Hook: guard-skills-catalog-fresh.sh             │
│  └── Docs Chain: .docs_chain/contexts/skills_catalog │
└──────────────────────────────────────────────────────┘
```

---

## §3. Generator design (G125 — implemented in `project/internal/skillscatalog/`)

### §3.1 Source layout

```
project/internal/skillscatalog/
  doc.go              — Package-level doc comment
  model.go            — GeneratorVersion = "skills-catalog/v3", slugify, escape funcs
  generate.go         — Generate(), Verify(), writeTree(), writeByDomain(), writeByKind()
  load.go             — loadRoster(), groupAndSortDeps(), verifyNoDanglingEdges()
  fingerprint.go      — computeRosterFingerprint(), computeSidecarIdentity(), read/write sidecar
  render.go           — renderREADME(), renderIndex(), renderDomainPage(), renderSkillDetail()
  generate_test.go    — Full real-DB test suite (18+ test functions, ~100KB)
  testdb_helper_test.go — Throwaway-DB helpers
```

### §3.2 Data loading
- `load.go` calls `Store.ListSkills()` (existing Store read method) for the
  full roster, then for each skill calls `GetByName()` which returns
  dependencies (6 types), dependents, and resources — all through existing,
  tested Store read paths.
- Deterministic sort keys: name → domain → kind. Every record is enriched
  with pre-sorted `DepsByType`, `Dependents`, `Resources` so renderers never
  re-sort.

### §3.2 Rendering
- `render.go` produces 5 file classes:
  1. **README.md** — header block, summary counts (by kind × by status),
     generator version, roster fingerprint (truncated 12-char prefix), links
     to INDEX.md and by-domain/by-kind landing pages.
  2. **INDEX.md** — flat table of ALL skills with columns: Name, Domain,
     Kind, Complexity, Status, Version, Description. Deterministic sort.
  3. **by-domain/<slug>.md** — one page per unique `metadata.domain` value.
     Table per domain with Name/Title/Description. An `_unclassified.md`
     page for skills with empty domain.
  4. **by-kind/<kind>.md** — three pages (atomic/composite/umbrella). Each
     lists skills of that kind with the same table format.
  5. **skill/<slug>.md** — one detail page per skill: full metadata,
     dependencies grouped by relation type (6 types in canonical order:
     requires → extends → composes → recommends → related_to →
     alternative_to), dependents list, resources table, content (full or
     excerpt per `EmbedFullContent` toggle).

### §3.3 Fingerprint system
- `fingerprint.go` implements §11.4.86 sha256 roster fingerprinting:
  - **Roster fingerprint:** concatenation of deterministic tuples
    `(name, kind, domain, content_hash)` for every skill, sorted by name,
    hashed with SHA-256.
  - **Composite identity:** `GeneratorVersion + roster_fingerprint + cfg
    output-affecting fields` — used as the sidecar value so a version bump
    or config change also triggers regeneration.
  - Stored in `docs/skills/.catalog_fingerprint`.
  - `Generate()` short-circuits when the composite identity matches the
    on-disk sidecar (detected as a no-op — returns `false, fingerprint`).
  - `Verify()` is read-only: recomputes current identity and compares
    against sidecar WITHOUT writing anything.

### §3.4 Generated-file discipline
- Every generated file carries a GENERATED FILE banner (§2.5 per DESIGN):
  "> **GENERATED FILE — DO NOT HAND-EDIT.** Regenerated from the live skill
  graph by the `skills-catalog` generator."
- `clearGeneratedMarkdown()` removes stale pages on every generation pass —
  if a skill or domain disappears, its page is removed, never left dangling.

---

## §4. Integration layers (G126–G132)

### §4.1 CLI subcommand (G126)
```bash
# Generate (or refresh) the catalog
skill-system docs skills-catalog generate [--output <dir>] [--force]

# Verify freshness — exit 0 if up-to-date, exit 1 if stale
skill-system docs skills-catalog verify [--output <dir>]
```

- Implementation: a `docs` command group in `cmd/cli/commands/`.
- `generate` calls `skillscatalog.Generate()` with a DB-connected Store.
- `verify` calls `skillscatalog.Verify()` and maps the boolean to exit code
  0/1 (matching Docs Chain's own `verify` exit-code contract).

### §4.2 REST endpoints (G127)
```
POST   /api/v1/skills/catalog/regenerate  — triggers generation
GET    /api/v1/skills/catalog/status       — returns inSync bool + fingerprint
```

- Wired alongside existing `skills_handler.go` in `internal/api/`.
- Both endpoints carry the same auth middleware as the rest of `/api/v1`.
- The write endpoint (`regenerate`) must respect G01's write-tool auth
  posture — only actions the current auth level permits.

### §4.3 MCP tools (G128)
```
skill_catalog_status  (read-only, no auth gate concerns)
  Returns: { inSync: bool, fingerprint: string, skills_count: int }
```

- Added alongside `skill_search`, `skill_get`, `skill_tree`, `skill_create`
  in `internal/mcp/tools.go` (or a new `internal/mcp/catalog_tools.go`).
- The write-capable `skill_catalog_regenerate` tool is deferred until G01's
  write-tool auth consolidation is fully settled.

### §4.4 Worker reconciliation loop (G129)
- A periodic worker job (default interval 60s, configurable) calls
  `skillscatalog.Verify()` against the on-disk output directory.
- On drift detected (inSync=false), the worker calls `Generate()` to
  converge.
- Implemented as a new `JobTypeCatalogReconcile` in `internal/worker/`.
- Same pattern as the existing auto-expand and validation worker cycles.

### §4.5 Write-trigger signal (G130)
- Every skill-graph write path signals the G129 worker queue:
  - `Store.Create()` / `Store.CreateFromTOML()` / `Store.ImportFromTOML()`
  - `AddDependency()` / `RemoveDependency()`
  - REST create/update/delete handlers
  - MCP create/update/import/delete tools
- Signal mechanism: atomically touch a "catalog-dirty" flag file or enqueue
  a `JobTypeCatalogReconcile` with immediate priority.
- The worker picks up the signal on its next tick and converges.

### §4.6 Pre-commit guard hook (G131–G132)
- `guard-skills-catalog-fresh.sh` — a §11.4.109/§11.4.201-compliant
  PreToolUse / pre-commit guard that:
  1. Reads the on-disk `.catalog_fingerprint`.
  2. Calls `skillscatalog.Verify()` against the current DB state.
  3. Compares: if the fingerprint would change (catalog is stale), BLOCKS
     the commit with a clear message to regenerate.
- Registered via the existing `post_update_hook.sh` auto-propagation seam
  (G132) so a fresh clone / new session gets the hook automatically.

---

## §5. Docs Chain export (G133)

Once the Docs Chain submodule is incorporated (X1/P13.T1 — G43 context),
wire `docs/skills/` through Docs Chain for HTML + PDF export:

- `.docs_chain/contexts/skills_catalog.yaml`:
  ```yaml
  context: skills-catalog
  root: ../../docs/skills
  derive:
    - format: html
      via: pandoc
      input: README.md
      output: ../../docs/skills/skills-catalog.html
    - format: pdf
      via: weasyprint
      input: README.md
      output: ../../docs/skills/skills-catalog.pdf
  ```

- Per-skill exports generate their own derive entries or one combined
  document with the INDEX.md as a single-file table of contents.

**Blocked** on X1/P13.T1 — currently draftable source-side but not
functionally verifiable until Docs Chain is incorporated.

---

## §6. Anti-bluff proof plan (§6 per DESIGN)

| # | Test type | What it proves | Go test function | Status |
|---|-----------|----------------|------------------|--------|
| P1 | Golden-good fixture | Deterministic generation: same input -> same bytes | TestSkillsCatalog_GoldenGood (8 sub-tests) | G125 done |
| P2 | Golden-bad: dangling edge | Detecting a stale sidecar | TestSkillsCatalog_GoldenBad_DanglingDependencyEdge | G125 done |
| P3 | Golden-bad: empty name | Defensive against G33-class export gaps | TestSkillsCatalog_GoldenBad_EmptySkillName | G125 done |
| P4 | Golden-bad: slug collision | Two names slugify to same filename | TestSkillsCatalog_GoldenBad_NameSlugCollision + DomainSlugCollision | G125 done |
| P5 | Fingerprint drift (add/modify/remove) | Fingerprint changes on mutation, stable otherwise | TestSkillsCatalog_FingerprintDrift_AddModifyRemove_StableOtherwise | G125 done |
| P6 | Fingerprint drift (title-only) | Title-only DB edit IS detected as drift (F1 fix) | TestSkillsCatalog_FingerprintDrift_TitleOnlyChange_Detected | G125 done |
| P7 | Fingerprint drift (id-only) | Drop+recreate same content but new UUID detected | TestSkillsCatalog_FingerprintDrift_IDChange_SameContent_Detected | G125 done |
| P8 | Timestamp churn guard | Touch-only updated_at bump does NOT trigger drift (F1 fix) | TestSkillsCatalog_TimestampChurn_TouchOnlyUpdatedAt_NoDrift | G125 done |
| P9 | Real-DB end-to-end | Generator works against seeded pgvector DB w/ ImportFromTOML | TestSkillsCatalog_RealSeedImport_NoDanglingInternalLinks | G125 done |
| P10 | Forged section heading sentinel | Sentinel distinguishes authentic footer from forged prose | TestSkillsCatalog_ForgedSectionHeading_SentinelDistinguishesAuthenticFooter | G125 done |
| P11 | Markdown injection (detail pages) | Title/dep-name injection neutralized | TestSkillsCatalog_MarkdownInjection_SkillDetailSurfaces_Escaped | G125 done |
| P12 | Markdown injection (domain pages) | Domain injection neutralized on by-domain + README | TestSkillsCatalog_MarkdownInjection_DomainSurfaces_Escaped | G125 done |
| P13 | Raw HTML injection | <img> tags escaped at every render site | TestSkillsCatalog_RawHTMLInjection_NameSurfaces_Escaped | G125 done |
| P14 | Table-cell escaping | Pipe/backslash/newline escaped | TestSkillsCatalog_TableCellEscaping_PipeBackslashNewline | G125 done |
| P15 | Length-prefixed tuple integrity | Netstring prevents boundary-forgery collision (F-A fix) | TestFingerprint_LengthPrefix_PreventsBoundaryForgedCollision | G125 done |
| P16 | Tags fingerprint coll. render equiv | Tags=[] vs Tags=[""] same output (R3-R1(b) fix) | TestSkillsCatalog_TagsFingerprintCollision_RenderEquivalence | G125 done |
| P17 | CLI integration | CLI exits 0/1 matching Verify | G126 PENDING | G126 |
| P18 | Hook blocking test | Guard blocks commit when catalog stale | G131 PENDING | G131 |

G125 implements P1-P16 in `project/internal/skillscatalog/generate_test.go`
(~99KB of test code). P17-P18 will land with G126/G131.


