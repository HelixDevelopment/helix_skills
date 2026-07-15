# CONTINUATION — HelixKnowledge Skill Graph System (MVP)

**Revision:** 7
**Last modified:** 2026-07-15T17:53:32Z
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
