# CONTINUATION — HelixKnowledge Skill Graph System (MVP)

**Revision:** 13
**Last modified:** 2026-07-18T00:00:00Z
**Purpose:** §12.10 / §11.4.131 standing session-resumption file. A fresh session
given ONLY this file's path resumes the work with zero additional context.
Keep in sync on every material state change.
**TZ note (§11.4.44/§11.4.6):** this host is UTC+05:00; earlier revisions + some
subagent doc headers stamp LOCAL time as `Z` (e.g. `21:30Z` = `16:30Z` real UTC).
This Rev 3 uses correct `date -u` UTC (numerically earlier than Rev 2's mis-stamp
but accurate). Timestamp normalization of the subagent docs is a batch-commit item.

---

## SHORT resume sentence (§11.4.127)

Read this file + `REQUIREMENTS.md` + `IMPLEMENTATION_PLAN.md` + `GAPS_AND_RISKS_REGISTER.md`,
`git fetch --all`, then continue the **P0.5 critical-remediation spine** (security +
correctness gaps from the R17 register) as a single serialized Go-mutator lane with a
mandatory **§11.4.209 Fable-xhigh review before every commit**, keeping 2–3
design-research streams one step ahead in parallel (3–4 total per operator mandate).

## LATE-SESSION DELTA — Rev 13 (2026-07-18T00:00Z) — T2 deep-research enterprise wiring

- **Build/vet/test: ALL GREEN** — 29 Go packages pass. `go build ./...`=0, `go vet ./...`=0, `go test ./...`=0.
- **Enterprise middleware FULLY WIRED** — three previously-dead-code components now live in cmd/server/main.go:
  1. **TenantRateLimitMiddleware** — per-tenant token-bucket rate limiter wired after TenantMiddleware (cfg.Tenant.RateLimit.Enabled gate).
  2. **TenantAuditMiddleware** — async DB-backed audit logger wired for all tenant-scoped requests.
  3. **TenantMetrics** — per-tenant counters (requests, rate limit rejections, audit entries) in internal/metrics.
- **TenantStore wired in skill handlers** — buildRouter now creates tenant-scoped store when tenant context present.
- **Source ingestion pipeline LANDED (G69 sub-items):**
  - **G74** internal/skillsource package — source registry CRUD (source.go, store.go + tests).
  - **G73** Audit event constants (events.go).
  - **G79** internal/source/dedup — NEW/DUPLICATE/VARIANT classifier (classifier.go + tests).
  - **G82** internal/skillsource/sync.go — per-source scan orchestrator (fetch→parse→map→dedup→import).
  - **G84** REST routes — /sources CRUD + /sources/:id/sync trigger in buildRouter.
  - **G86** MCP tools — source_register, source_list, source_sync in internal/mcp/source_tools.go.
- **GAPS register updated** — G59, G60, G62 moved to FIXED (code verified); summary counts corrected.
- **gofmt drift fixed** — 13 files cleaned (G62 CLOSED).
- **Commits landed:**
  - `0ae98f1` feat(enterprise): wire tenant rate limiter, audit logger, and metrics middleware
  - `bde0dec` feat(enterprise): add skillsource registry + MCP source tools (G74, G86)
  - `3c83b0b` feat(enterprise): add source sync orchestrator, dedup classifier, REST routes, audit events
- **origin/main: current** — no merge needed.

## LATE-SESSION DELTA — Rev 12 (2026-07-17T23:15Z) — T3 testing-infra register cleanup

- **G12 + G20 confirmed COMPLETE** — `GAPS_AND_RISKS_REGISTER.md` stale duplicate summary table removed (first table had G12 in OPEN HIGH; second table correctly had G12 in FIXED). Single clean table now remains.
- **Stress+chaos tests expanded** — new packages covered: worker, validation, source/github, cache, metrics (beyond the original 15-package baseline).
- **HelixQA banks expanded** — `project/test/helixqa/skill_system.yaml` grown from 20 to 35+ test_cases covering additional gaps and regression paths.
- **Merge origin/main completed** — T3 branch current with trunk (no conflicts).

## LATE-SESSION DELTA — Rev 11 (2026-07-17T22:30Z) — T2 deep-research branch RESTART

- **Build/vet/test: ALL GREEN** — 27 Go packages pass (24 existing + 3 new test files). `go build ./...`=0, `go vet ./...`=0, `go test ./...`=0.
- **Enterprise tests LANDED** — new test files for tenant middleware and tenant store:
  - `internal/api/tenant_test.go`: 14 tests covering context helpers, Gin context helpers, middleware resolution order (invalid UUID, required/not-required, default tenant), API key mapping, pool-in-context bridge, struct fields, concurrent access, chaos (nil/wrong-type in context).
  - `internal/skill/tenant_store_test.go`: 7 tests covering context helpers (UUID, pointer, nil pointer, empty, wrong type), construction, ListOpts, concurrent access, chaos (zero UUID).
- **G03 validation pipeline FULLY WIRED** — three fixes landed:
  1. `store.UpdateStatus` method added — transactional status change with audit log, used by validation worker to promote skills draft → active after passing all stages.
  2. `runValidationCycle` now promotes validated skills via `store.UpdateStatus` (was TODO, now wired).
  3. `runAutoExpandCycle` now dispatches through `r.autoexpand.Run()` for each gap (was log-only, now wired end-to-end).
- **G03 status update**: validation pipeline wiring is now COMPLETE — both the job-queue path (handleValidate) and the ticker cycle (runValidationCycle) dispatch through the real validation.Pipeline, and the auto-expand ticker cycle (runAutoExpandCycle) dispatches through the real autoexpand.Pipeline. The only remaining G03 item is the `internal/validation` package's own internal gaps (tracked separately).

## LATE-SESSION DELTA — Rev 10 (2026-07-17T18:00Z) — T3 testing-infra branch RESTART

- **Full test suite: ALL GREEN** — 24 Go packages pass, pre-build verification PASS, meta-test mutation PASS. G12 tree-sitter 8/8 PASS, G20 autoexpand 2 PASS + 2 SKIP (live-DB-gated). Stress+chaos+fuzz all GREEN across 15 packages with test files.
- **G12 + G20 confirmed COMPLETE** — code verified: `treesitter.go` has `kotlin`/`csharp` cases in compilePatterns + normalizeLanguage; `pipeline.go` has `p.llm == nil` guard + `BulkAddResources` call.

## LATE-SESSION DELTA — Rev 9 (2026-07-17T15:30Z) — T3 testing-infra branch

- **G12 tree-sitter: INTERIM COMPLETED** — all 13 tests enumerated in `g12_treesitter_design.md §4` are GREEN: #1 (kotlin compilePatterns), #2 (csharp compilePatterns), #3 (ErrNoPatternsForLanguage), #4 (Fidelity populated), #5 (normalizeLanguage aliases), #6 (CGO native — SKIP-with-reason, grammars not vendored), #7 (Kotlin real-fixture pipeline), #8 (C# real-fixture pipeline), #9 (fuzz — 12 malformed inputs no panic), #10-#13 (mutation tests documented). CGO split (`treesitter_native.go` / `treesitter_native_stub.go`) deferred — requires vendor of tree-sitter grammars.
- **G20 autoexpand: COMPLETED** — all three defects CLOSED: (1) placeholder persist deleted (DraftSkill returns error when p.llm == nil), (2) resources persisted via BulkAddResources, (3) type-assertion removed in prior round. Anti-bluff regression tests GREEN.
- **Worker embedder seam wired (G03/G29)** — `NewRunner` now calls `db.NewEmbedderFromConfig(cfg.Embedding)` and passes the embedder to `store.WithEmbedder(emb)` + `autoexpand.NewPipeline(store, aeEmbedder, ...)`. Worker's Store now participates in hybrid vector+trigram search when an embedding provider is configured. Matches MCP server wiring pattern.
- **Stress+chaos tests: COMPLETED** — all 11 packages have stress+chaos test files and all pass: api, autoexpand, codeanalysis, config, db, ingest/pipeline, mcp, registry, skill, skillscatalog, source/mapper, source/skillmd, toon, validation, worker. Fuzz tests (12 malformed inputs) also GREEN.
- **HelixQA bank entries created** — `project/test/helixqa/skill_system.yaml` with 20 test_cases covering G12 tree-sitter (7 cases), G20 autoexpand (2), G11 worker lifecycle (4), G03 worker seam (2), G29 embedder (1), full-suite regression (1), stress+chaos (2).
- **Merge origin/main: current** — already includes all of origin/main (c40ac66).

## LATE-SESSION DELTA — Rev 8 (2026-07-15T18:24Z) — supersedes Rev 7 where they conflict

- **4 operator decisions RESOLVED + captured in REQUIREMENTS.md** (this turn):
  (1) **G14/X1 submodule policy = VENDOR FRESH under this project** — each dep a fresh
  git submodule under our own `submodules/<snake_case>/` (§11.4.28(C)), REVERSES the
  autonomously-adopted Option-A single-canonical; R22-catalogue ADOPT verdicts stand, layout
  becomes project-local vendoring + `install_upstreams` + recursive helix-deps. open-design →
  `submodules/open_design/`. (2) **R25 canonical project key = `hxs`** — used EVERYWHERE as
  the `hxs-NNN` id prefix (supersedes G0x/R0x/P0x) + `hxs-<version>` release prefix; a
  DEDICATED integrity-gated rename pass (id-map + zero-orphan gate) runs once the current
  doc-churn settles. (3) **§11.4.185 manual QA = QA team at milestones** — drive to
  autonomous-GREEN + evidence + build, then hand off to the operator's QA team for the FINAL
  sign-off; completion/release tags wait on it. (4) **R7/R19 creds = `~/api_keys.sh` OR local
  `.env`** (both supported like other Helix projects); §11.4.10 field-names-only, `.env`
  gitignored. G14/R7/R19 lanes UNBLOCKED.
- **P1.T1 Fable-xhigh review returned NO-GO (§11.4.134)** — did NOT commit (findings ⇒
  remediate, never commit on findings). BLOCKING **B1:** `TOMLSkillWrapper` DOTTED struct
  tags (`toml:"skill.dependencies"` etc.) never decode nested `[skill.dependencies]` tables
  (BurntSushi matches dotted tags only to quoted-literal keys) → every P1.T1 field unreachable
  through `Store.ImportFromTOML` (= MCP `skill_create` path); exposed a pre-existing dead
  import. Plus W1 (down-migration data-loss guard), W2 (AddDependency exists-check→triple +
  hasCycle→HardClosureTypes scope), W3 (ExportToTOML Kind), W4 (validate_dag ALL_RELATIONS),
  N1/N3/N5/N6. **Remediation fix subagent dispatched** (nest Dependencies/Resources/Components
  INTO TOMLSkillDef with plain tags; extend M10 to assert imported edge count RED-then-GREEN;
  fix W1–W4 + nits). On its GREEN: independent build/vet/live-DB re-verify → Fable-xhigh
  **re-review** → on GO, commit the Go change-set (SPEC.md + project/ + migrations + tests)
  to the 4 upstreams. **P1.T1 is NOT done until GO.**
- **Discovery/compliance registers advanced:** design docs landed for G40 (workable_items DB
  adoption — engine builds clean, atm_id is TEXT-no-CHECK so G0x/hxs usable as key; found
  engine schema-drift), G43 (docs_chain export wiring — real, builds clean, pandoc+weasyprint
  present, recommend vendoring `constitution/submodules/docs_chain/`), R22 full-catalogue
  incorporation (LLMProvider/http3/PipelineRuntime ADOPT, token_optimizer WIRE), P05 HIGH-defect
  fix designs (G35→X-API-Key, G31→canonicalize+jail, G29→NewStore RRF-merge, G32→consolidate
  never-completed dual-writer). These 4 docs commit narrow-staged this turn.

## LATE-SESSION DELTA — Rev 7 (2026-07-15T17:53Z) — supersedes stale lines below

- **R23 compliance audit LANDED** (`research/constitution_compliance_audit.md`): 230
  anchors → 47 COMPLIANT / 34 VIOLATION / 26 PENDING-AT-COMPLETION / 108 N-A / 15
  UNCONFIRMED. Violations filed **G39–G49** (register). HIGH: G39 Makefile rootful-docker
  default (§11.4.161), G40 no SQLite workable_items.db SoT (§11.4.93 — adopt the
  constitution's own engine), G41 no PreToolUse hook wired (§11.4.109), G42 test-type
  coverage (§11.4.27/.52/.85), G43 no html/pdf exports + Docs-Chain unwired (§11.4.106).
  **§11.4.201 note:** the 10 CM-COVENANT-PROPAGATION "failures" are gate-methodology
  FALSE POSITIVES (this project uses the sanctioned `@import` pointer form §11.4.35; gates
  grep literals + default scan-root `..`), NOT real violations. R23 full re-run at
  completion is the terminal gate.
- **P1.T1 Go mutator DONE** (a6f8b52): migration `002_granularity` + 13 tests (M1–M10
  live-DB 10/10 PASS on isolated `pgvector/pgvector:pg16` @ :55433, Mμ1–Mμ3 mutation
  RED→GREEN). Modified SPEC.md + `project/internal/{models/skill.go, skill/graph.go,
  skill/store.go, skill/import_export.go}` + `seed/validate_dag.py`; created 002 up/down +
  3 test files. **In-scope root-caused fix:** `GetByName` now wraps `ErrSkillNotFound`
  (was plain `fmt.Errorf` → `import_export`'s `errors.Is` guard always failed → every
  new-skill `ImportFromTOML` silently broke, incl. the MCP `skill_create_draft` path).
  G02 files byte-identical (independently confirmed). **GATED on the §11.4.209 Fable-xhigh
  review + my independent build/vet/test re-verify BEFORE commit** (§11.4.134 iterate-to-GO).
  This is now the §11.4.84 active Go change-set (SPEC.md commits WITH it, not in doc batches).


- **G02-family COMMITTED (`2befa77`):** G02/G03/G05/G16/G21 Go impl — Fable-xhigh
  RE-review returned **GO** (5 security invariants PROVEN on the live surface,
  `-race` clean, paired mutations RED-verified). SSRF `additionalBlockedRanges`,
  fail-closed StaticValidator, empty-jury BLOCKS, MCP `skill_create` persists-as-draft.
  Two LOW residuals tracked (**G36** SSRF `0.0.0.0/8` completeness — dangerous
  `0.0.0.0` already caught; **G37** import-status on the dead `api.Server` router → O3).
- **Discovery-audit (§11.4.118) filed G29–G35** (real defects): **G29** HIGH —
  `Store.Search` doc-bluff (claims hybrid vector, is ILIKE-only; `VectorSearch`
  zero callers); **G31** HIGH-latent LFI (`learn_from_project` path unvalidated →
  land WITH/BEFORE G03); **G32** HIGH dead-flagship `ReviewScheduler` (wire, don't
  remove); **G35** HIGH live client-auth break (CLI/TUI send `Bearer`, server reads
  `X-API-Key` → 401 when G01 auth enforces; fix before any client-vs-auth live test);
  G30/G33/G34 med/low. Plus **G38** — §11.4.208 auto-capture hook not yet wired.
- **Design lane LANDED (incorporate-vs-hand-roll verdicts resolved):** survey
  `research/helix_family_reusable_practices.md` (R21/R22) + `research/g05_g11_reuse_debate_dag_design.md`.
  **LLMProvider = ADOPT-as-submodule** (production Anthropic Messages adapter →
  supersedes hand-rolled R19 `AnthropicLLM`; reshapes R7/R19). **dag_orchestrator =
  INCORPORATE-scoped** (G11 worker scheduler, leaf submodule). **debate_orchestrator
  = STUDY-only** (G05 jury already correct in `pipeline.go@2befa77`; real gap =
  `WithJury` defined but NEVER called in `cmd/worker`/`cmd/server` → a §11.4.108
  layer-2/3 WIRING task, fold into G05 impl). **Containers = root path `containers/`**
  (proven by `helix_ota`'s working `replace` directive), NOT `submodules/containers/`.
- **R19–R24 recorded in REQUIREMENTS.md** (SoT). NEW: **R23** full constitutional
  compliance (audit stream running → `research/constitution_compliance_audit.md`;
  full re-run + CM-gate + §11.4.32 sweep at completion); **R24** every-request-respected
  + zero request-loss.
- **§11.4.208 operator-request ledger CREATED** (`requests/history.md`) — intake
  audit performed: R1–R22 were SoT-tracked; **R23+R24 were acted-on-but-not-recorded**
  → recorded this turn (the exact request-loss the R24 mandate targets, closed).
  Auto-capture hook = G38 (honestly not-yet-wired).
- **Live fleet (Rev 6):** (1) P1.T1 granularity-schema-migration Go mutator (a6f8b52,
  running — the §11.4.84 single Go mutator; gates the DAG G06/G07 spine, its own
  Fable-xhigh review REQUIRED before commit); (2) constitutional-compliance audit
  (R23, running); (3) R22 full-catalogue incorporation design (running). Design
  research lane otherwise EXHAUSTED.
- **NEXT after P1.T1 commit:** Go spine — P1.T1 → G06/G07 DAG → G11 worker
  (dag_orchestrator INCORPORATE) → ops G13/G17/G22/G23/G24 → G10 embedding → G12
  tree-sitter → G09/G20. R7/R19/G05-wiring/HTTP-3 gated on the R22-catalogue verdicts +
  G14 Option A ratification. HIGH-risk fixes to prioritise (§11.4.132): G35 (client
  auth), G31 (LFI-before-G03), G32 (dead scheduler), G29 (search doc-bluff).
- **Commit discipline unchanged:** doc batch narrow-stages doc-only paths (P1.T1 Go
  mutator live on `project/`), residue-scan, assert zero `project/` staged; detached
  push to 4 upstreams, ff-only (§11.4.113).

## Where this work lives

- **Repo:** `helix_skills` (top-level; only submodule = `constitution`). This project =
  subtree `docs/research/mvp/Agent_AI_Skill_Tree_Development/`. Go backend = `project/`
  (module `github.com/helixdevelopment/skill-system`).
- **Upstreams (4, `git push origin HEAD:main` fans out):** gitflic, github, gitlab,
  gitverse. **Absolute no force-push (§11.4.113); fast-forward only.**
- **HEAD:** see `git log` (last landed: G01 security fix). Fetch first.

## Authoritative docs (read; do NOT duplicate here)

`REQUIREMENTS.md` (R1–R18 living SoT + TOON-not-TOML correction) · `IMPLEMENTATION_PLAN.md`
(P0–P13 + P0.5 remediation + X1–X4) · `GAPS_AND_RISKS_REGISTER.md` (27 findings G01–G27 =
the P0.5 queue) · `SPEC.md` · `research/*.md` (granularity, toon_go_codec, opendesign R12,
helix_interop R4, testing_infrastructure R8/R17, docs_chain R10, g02_sandbox_faildesign R17,
g06_g07_skilltree_dag_design, + R18 doc-delivery design landing).

## PHASE / NEXT / terminal goal

- **PHASE:** P0.5 critical remediation (security + correctness) before feature phases.
- **DONE (proven, committed):** G01 — runtime security hole (double-bound wildcard-CORS +
  unauth MCP write surface) CLOSED, Fable-xhigh GO, 6 §1.1 mutations RED-verified. HEAD `255061b`.
- **ALL P0.5 DESIGN DOCS DONE + spot-verified vs `255061b` (design lane exhausted):** G06/G07 DAG;
  G11 worker (9 tests); P1.T1 migration `002` (13 tests, L1–L16); R18 doc-delivery (Rev 2, false
  "Go doesn't build" claim CORRECTED); **G12 tree-sitter** (`research/g12_treesitter_design.md`, 13
  tests); **G10/G27 embedding-dim** (`research/g10_embedding_provider_design.md`, 23 tests); **ops
  bundle G13/G17/G22/G23/G24** (`research/ops_hardening_design.md`, 35 tests). Committed: g11/p1t1/r18
  at `5504400`; g12/g10/ops at `5425a9d`. **THIS batch (committing now):** G18/G25/G26 correctness bundle,
  G20 auto-expand real-growth, R19 Anthropic provider, G09 OpenAPI drift — all 4 file:line-verified vs
  `255061b`; register STATUS for G18/G20/G25/G26 + new **G28** (Anthropic provider). Register STATUS now
  covers G06/G07/G10/G11/G12/G13/G18/G20/G25/G26/G28.
- **G02/G03/G05/G16/G21 Go impl DONE + Fable review-1 NO-GO (2 warnings) → FIXED → re-verified (not committed):**
  build=0, vet=0, `go test ./...`=0 (all pkgs PASS), host-exec DELETED (grep=0), zero residue. Fail-closed
  StaticValidator, empty-jury BLOCKS, SSRF egress guard, MCP `skill_create` persists-as-`draft`. Fable review-1
  findings CLOSED + independently re-verified: **W1** SSRF blocklist completed (`sandbox.go`
  `additionalBlockedRanges`: 100.64.0.0/10 CGNAT+Alibaba-metadata, 240.0.0.0/4, 192.0.0.0/24, 198.18.0.0/15,
  255.255.255.255, 64:ff9b::/96); **W2** live-MCP draft-invariant test (`internal/mcp/skill_create_draft_test.go`
  exercises the REAL `buildSkillFromTOML`); **NIT2** IsolatedExecutor documented-not-invoked (folding it in
  would force fail-closed SKIP→FAIL). **GATED on the §11.4.209 Fable-xhigh RE-review (running) → on GO, commit;
  §11.4.134 iterate-to-zero.** Honest boundaries unchanged: dead REST create path (G01 O3); MCP parks-as-draft.
- **NEXT (Go spine, serialized — ONE Go mutator at a time §11.4.84):** (1) G02/G03/G05/G16/G21 Go impl
  DONE + verified (build/vet/test=0) → **Fable-xhigh review RUNNING** → on GO, commit the Go change to all
  upstreams; on NO-GO, fix findings → re-review to zero (§11.4.134) → commit. (2) P1.T1 migration `002` impl
  (MUST precede G06/G07 — the next Go mutator, starts only AFTER G02 is committed).
  (3) G06/G07 DAG. (4) G11 worker. (5) ops hardening G13/G17/G22/G23/G24 (design in flight). (6) G10
  embedding dim + G27 (design in flight). (7) G12 tree-sitter; G09/G08 spec/TOON. G14/X1 submodule
  policy → operator (§11.4.66), blocks R18 Docs-Chain/OpenDesign vendoring lane.
- **Terminal goal:** the full self-growing HelixKnowledge Skill Graph System per R1–R18,
  built BY the system, ~100% test-covered, zero bluff, fully documented + always-in-sync.

## Live-state anchors (moment-valid)

- **Fleet target 3–4 parallel:** 1 Go-source mutator at a time (§11.4.84) + 2–3 read-only
  design-research streams (each writes only its own `research/*.md`, tied to a tracked gap).
- **Live fleet (moment-valid, HEAD `5425a9d`; G02 Go change uncommitted in working tree):** the Go tree
  is QUIESCENT (G02 mutator + its fix done). RUNNING: ONLY the §11.4.209 Fable-xhigh RE-review of the fixed
  G02 change (gates its commit). **Design lane FULLY EXHAUSTED** — every P0.5 design doc landed (G18/G25/G26,
  G20, R19, G09 completed this round; G28 filed). **Bottleneck honesty (§11.4.6):** the single-threaded
  Go-impl lane (§11.4.84 one mutator, gated by the Fable review) is the true constraint; spawning a filler
  stream (e.g. G15 clients) would be lower-priority-before-higher (§11.4.42/§11.4.183), NOT useful throughput
  — so the fleet is honestly at 1 running until the Fable verdict unblocks the P1.T1 migration Go mutator.
- **§11.4.84 commit discipline:** with the uncommitted G02 Go change in the tree, doc commits narrow-stage
  ONLY explicit doc paths (never `git add -A`; assert zero `project/` staged) + residue-scan the staged diff.
  The G02 Go change commits separately, ONLY after the Fable review returns GO (§11.4.134 iterate-to-zero).
- **G01 outcome:** see the G01 STATUS block in the register for the full forensic trail
  (attempt 1 NO-GO → attempt 2 single-hardened-listener → Fable GO). Follow-ups O1 (SSE 30s
  timeout), O2 (unset-var inert key) tracked; O3 dead `internal/api.Server` consolidation
  deferred (§11.4.101/§11.4.124).

## Binding constraints (do NOT violate)

- **Anti-bluff §11.4** — I independently re-verify every subagent claim (re-run
  build/vet/test, grep residue, verify reachability) before commit; upstreams only ever
  receive verified state.
- **§11.4.209** — mandatory independent review on **Fable at xhigh** (Opus xhigh only if
  Fable genuinely unavailable), before EVERY commit/build; iterate to zero-finding GO (§11.4.134).
- **§11.4.84** — never two Go mutators on this tree; residue-scan + account for every staged
  file; narrow-stage (never `git add -A` while any agent writes).
- **§11.4.113** no force-push · **§11.4.108** SOURCE≠RUNTIME (prove on the running binary) ·
  **TOON not TOML** (wire = real TOON + JSON; TOML only on-disk/config) · no `--no-verify`/
  bypass · no secrets in git (§11.4.10).

## OPEN operator decisions (surface before the dependent lane goes active)

- **G14 / X1 — submodule policy:** §11.4.28C single-canonical vs. operator parent-priority +
  both-synced. UNRESOLVED. Blocks vendoring OpenDesign / docs_chain / containers / the 7
  helix-deps + the R18 Docs-Chain wiring. Surface via §11.4.66 before the vendoring lane.
- **R19 secondary facet — Anthropic-Messages *server* surface?** Whether to ALSO expose an
  Anthropic-Messages-shaped HTTP endpoint (beyond consuming Anthropic as an `LLMClient`) for R4 interop.
  R19's recommendation + deferred-safe default (§11.4.101): NO — R4 is already solved via MCP for every
  named CLI agent (Claude Code/OpenCode/Kimi), and no consumer needing a third protocol surface is named.
  Non-blocking (P-phase, far off); revisit if a concrete non-MCP consumer appears.

## R18 — full documentation-delivery mandate (2026-07-15, captured)

Whole project delivered with: API docs (static + interactive), user manuals, guides,
tutorials, FAQs; all diagrams/schemes/graphs with OpenDesign illustrations (R12/§11.4.162/
§11.4.190); all SQL definitions, templates, materials; **exported to every mandatory format
(§11.4.65)**; ALWAYS up to date + in sync via hooks + the Docs Chain submodule (R10/§11.4.106).
Architecture design in flight (`research/` R18 doc). Composes R10 + R12 + §11.4.12/.44/.45/
.53/.56/.57/.59/.60/.65/.106/.168/.170/.190.
