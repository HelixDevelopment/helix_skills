# G136 — Retroactive Severity Assessment for G52–G137

**Revision:** 1
**Last modified:** 2026-07-17
**Status:** preliminary assessment (not independently re-verified per item;
see GAPS register Summary-counts note)
**Method:** category-based rubric applied to each item's stated type +
evidence as recorded in the register at date of assessment.

---

## §1. Rubric

| Severity | Criteria |
|----------|----------|
| **CRITICAL** | Security hole (auth bypass, RCE, data-loss); constitutional §X violation with active impact; running binary !== hardened code |
| **HIGH** | Core feature broken/missing; latent security (reachable by normal ops); spec-drift making contracts unusable; §11.4.108 SOURCE≠RUNTIME |
| **MEDIUM** | Functional weakness, degraded UX, test gap, process gap, doc gap (non-spec-breaking) |
| **LOW** | Doc nit, naming, cosmetic, latent foot-gun not reachable today |
| **N/A** | Already fixed/closed; planning item |

---

## §2. G52–G58 (Additional findings, post-R23 audit)

| ID | Type | Proposed severity | Rationale | Current status |
|----|------|-------------------|-----------|----------------|
| G52 | Bug | MEDIUM | Health-probe misalignment (TUI/CLI vs server) — user-visible break in client connectivity | Fixed |
| G53 | Bug | MEDIUM | Catalog-query correctness in VectorIndex readiness — degraded robustness | Fixed |
| G54 | Task | LOW | Pre-existing `gofmt` drift in pipeline.go — cosmetic only | Queued |
| G55 | Bug | MEDIUM | Phantom OpenAPI routes — spec-drift per G09 | Queued |
| G56 | Bug | HIGH | Compose deployment contract broken (double ENTRYPOINT, worker runs server binary, nonexistent flags) — blocks clean deployment | Queued |
| G57 | Bug | MEDIUM | ACP stdio transport unwired — dead flagship (§11.4.108) | Fixed |
| G58 | Bug | MEDIUM | Migration 003 pg_trgm + Store.Search restore on clean deploy — test-only gap exposed | Fixed |

---

## §3. G59–G68 (Session-discovered findings this round)

| ID | Type | Proposed severity | Rationale | Current status |
|----|------|-------------------|-----------|----------------|
| G59 | Bug | **HIGH** | Embedding not written on create/update — G29 hybrid search degrades to keyword-only post-merge | Queued |
| G60 | Bug | MEDIUM | Conflict oracle uses ranked Search instead of exact GetByName — latent, only after G59+G29 | Queued |
| G61 | Task | MEDIUM | Two divergent /health implementations — composes G01-O3 | Queued |
| G62 | Task | LOW | gofmt drift across 18 files — cosmetic, zero behaviour change | Queued |
| G63 | Bug | HIGH | 4th divergent route surface + unreachable registry CLI/TUI — §11.4.108 SOURCE≠RUNTIME, operator blocked | Operator-blocked |
| G64 | Bug | MEDIUM | migrate.sh unconditionally deletes schema_migrations row on down — migration state desync | Queued |
| G65 | Bug | MEDIUM | stop.sh rejects --compose arg; restore.sh swallows failure — silent skip | Queued |
| G66 | Bug | MEDIUM | Seed corpus missing 3 prerequisite TOMLs — blocks clean import | Queued |
| G67 | Task | LOW | qa-results/ gitignore policy conflict — procedure only, no code impact | Operator-blocked |
| G68 | Bug | MEDIUM | RemoveDependency coarse delete + dead/unwired — §11.4.124, needs verification | Queued |

---

## §4. G69–G92 (GitHub skills ingestion epic)

All are **Feature** type — severity reflects "importance to the roadmap", not
"severity of a defect":

| ID | Type | Proposed severity | Rationale | Status |
|----|------|-------------------|-----------|--------|
| G69 | Feature | HIGH | Operator-mandated net-new capability | Queued |
| G70 | Task | MEDIUM | Migration — schema prerequisite | Queued |
| G71 | Task | MEDIUM | Migration — proposal store | Queued |
| G72 | Task | MEDIUM | Config surface — shared with G94 | Queued |
| G73 | Task | LOW | Audit event constants — mechanical | Queued |
| G74 | Feature | MEDIUM | Source registry CRUD — overlaps G97 | Queued |
| G75 | Feature | MEDIUM | GitHub REST fetch client | Queued |
| G76 | Feature | LOW | Shallow-clone fallback | Queued |
| G77 | Feature | MEDIUM | SKILL.md parser | Queued |
| G78 | Feature | MEDIUM | ParsedSkill→models.Skill mapper | Queued |
| G79 | Feature | MEDIUM | Dedup classifier | Queued |
| G80 | Feature | MEDIUM | ImportSkillModel store method | Queued |
| G81 | Feature | MEDIUM | Enhancement delta extraction | Queued |
| G82 | Feature | MEDIUM | Per-source scan orchestrator | Queued |
| G83 | Task | MEDIUM | Worker wiring for source rescan | Queued |
| G84 | Feature | MEDIUM | REST routes | Queued |
| G85 | Feature | MEDIUM | CLI commands | Queued |
| G86 | Feature | MEDIUM | MCP tools | Queued |
| G87 | Feature | LOW | TUI read-only pane | Queued |
| G88 | Task | MEDIUM | e2e test with real anthropics/skills | Queued |
| G89 | Task | LOW | Stress + chaos test suite | Queued |
| G90 | Task | LOW | Vendor Challenges + HelixQA submodules | Queued |
| G91 | Task | LOW | HelixQA Challenge bank entry | Queued |
| G92 | Task | LOW | Docs sync for new surfaces | Queued |

---

## §5. G93–G122 (Unified multi-source ingestion epic)

| ID | Type | Proposed severity | Rationale | Status |
|----|------|-------------------|-----------|--------|
| G93 | Feature | HIGH | Net-new capability, operator-mandated | Queued |
| G94 | Task | MEDIUM | Config section — overlaps G72 | Queued |
| G95 | Task | MEDIUM | Ingestion schema migration | Queued |
| G96 | Task | MEDIUM | Router duplication — UNCONFIRMED: duplicates G01-O3 | Queued |
| G97 | Feature | MEDIUM | Source interface + types | Queued |
| G98 | Feature | MEDIUM | Filesystem source (bulk) | Queued |
| G99 | Feature | MEDIUM | HTTP/website source | Queued |
| G100 | Feature | MEDIUM | PDF source | Queued |
| G101 | Feature | LOW | OpenAPI/API-schema source | Queued |
| G102 | Feature | LOW | FTP source | Queued |
| G103 | Feature | LOW | SMB source | Queued |
| G104 | Feature | LOW | WebDAV source | Queued |
| G105 | Task | LOW | NFS honest-gap documentation | Queued |
| G106 | Feature | MEDIUM | HTML EXTRACT+NORMALIZE | Queued |
| G107 | Feature | MEDIUM | PDF EXTRACT+NORMALIZE | Queued |
| G108 | Feature | LOW | OpenAPI EXTRACT+NORMALIZE | Queued |
| G109 | Feature | MEDIUM | LLM-REFINE stage (interface) | Queued |
| G110 | Feature | MEDIUM | DEDUP stage | Queued |
| G111 | Feature | MEDIUM | CREATE/EXTEND stage | Queued |
| G112 | Feature | MEDIUM | Ingestion job orchestration | Queued |
| G113 | Feature | MEDIUM | Recursive directory watcher | Queued |
| G114 | Task | MEDIUM | Worker JobTypeIngestSource | Queued |
| G115 | Feature | MEDIUM | REST /api/v1/ingest/* | Queued |
| G116 | Feature | MEDIUM | CLI ingest command group | Queued |
| G117 | Feature | MEDIUM | MCP skill_ingest_source tool | Queued |
| G118 | Task | LOW | Anti-bluff test suite | Queued |
| G119 | Feature | LOW | Periodic polling (deferred) | Queued |
| G120 | Feature | LOW | Deep-research-extend (deferred) | Queued |
| G121 | Feature | LOW | TUI ingestion pane (deferred) | Queued |
| G122 | Feature | LOW | Source-removal → staleness policy (deferred) | Queued |

---

## §6. G123–G137 (Architecture reconciliation + catalog + post-session)

| ID | Type | Proposed severity | Rationale | Status |
|----|------|-------------------|-----------|--------|
| G123 | Task | **HIGH** | Architectural overlap between G69/G93 — if unresolved, blocks both epics | Queued |
| G124 | Feature | HIGH | Docs/skills always-in-sync catalog — operator-mandated | Queued |
| G125 | Task | MEDIUM | Generator impl — already DONE in skillscatalog/ | DONE |
| G126 | Task | MEDIUM | CLI subcommand | Queued |
| G127 | Feature | MEDIUM | REST catalog endpoints | Queued |
| G128 | Feature | MEDIUM | MCP catalog tools | Queued |
| G129 | Task | MEDIUM | Worker reconciliation loop | Queued |
| G130 | Task | LOW | Write-path signal wiring | Queued |
| G131 | Task | MEDIUM | Pre-commit guard hook | Queued |
| G132 | Task | LOW | Hook auto-propagation | Queued |
| G133 | Task | MEDIUM | Docs Chain context (blocked) | Queued |
| G134 | Task | MEDIUM | Anti-bluff proof plan — largely done in G125 test suite | Partial |
| G135 | Task | LOW | HelixQA Challenge bank entry | Queued |
| G136 | Task | LOW | This severity assessment itself | Queued |
| G137 | Bug | **HIGH** | Autoexpand gap-detection inert against store-constructed graphs — core feature silently no-op | Queued |

---

## §7. Summary counts (aligned with GAPS register)

Counts match the Summary counts table in `GAPS_AND_RISKS_REGISTER.md` (2026-07-17 update). All sub-items of feature epics are enumerated individually.

| Severity | Count | Notes |
|----------|-------|-------|
| **CRITICAL** | 4 | G01–G04 |
| **HIGH** | 78 | G05–G15 (11), G29, G31, G32, G35, G39–G43 (5), G57, G59, G63, G69–G92 (24), G93–G122 (30), G137 |
| **MEDIUM** | 35 | G16–G23 (8), G28, G30, G34, G44, G45, G47, G51, G55, G56, G58, G60, G61, G64, G66, G123, G124–G135 (12) |
| **LOW** | 18 | G24–G27 (4), G33, G36, G37, G38, G46, G48, G49, G52, G53, G54, G62, G65, G67, G68 |
| **N/A** | 1 | G136 (meta-assessment task itself) |
| **TOTAL** | **136** | G01–G135 + G137; G136 deliberately unrated |

*Note: Feature-epic sub-items (G70–G92, G94–G122, G125–G135) inherit their umbrella's
HIGH/MEDIUM umbrella-level severity for the purpose of the total count. Individual
sub-items within each epic may later be downgraded to MEDIUM or LOW upon
evidence-based re-assessment (see the per-item tables above for proposed
individual severities).*

---

## §8. Action items

1. **Review this assessment** for items whose evidence is thin or status has changed since filing.
2. **Update the GAPS register's Summary table** each time an item's severity is changed or the item is closed.
3. **Finalize per-epic sub-item severities** when G69, G93, and G124 move from Queued to In Progress — individual sub-items may need separate assessments at that point.

---

*End of severity assessment. Update in GAPS_AND_RISKS_REGISTER.md's Summary
counts table after operator/peer review.*
