# CONTINUATION — HelixKnowledge Skill Graph System (MVP)

**Revision:** 4
**Last modified:** 2026-07-15T16:43:28Z
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
  at `5504400`; g12/g10/ops in the next doc batch. Register STATUS lines for G06/G07/G10/G11/G12/G13.
- **G02/G03/G05/G16/G21 Go impl DONE + INDEPENDENTLY VERIFIED (not yet committed):** build=0, vet=0,
  `go test ./...`=0 (all pkgs PASS), host-exec paths DELETED (grep=0 for executeProcess/executeGoSnippet/
  WASMSandbox/LD_PRELOAD/os/exec), zero mutation residue. Fail-closed StaticValidator (executes nothing),
  empty-jury BLOCKS, SSRF egress guard, MCP `skill_create` validates-before-persist + persists-as-`draft`.
  Honest boundaries disclosed: REST `api.Server` create path is dead code until G01 consolidation; MCP path
  parks-as-draft (not hard-reject); PodmanIsolatedExecutor unimplemented (SKIPs honestly). **GATED on the
  §11.4.209 Fable-xhigh review (running) → then commit.**
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
- **Live fleet (moment-valid, HEAD `5504400`; G02 Go change uncommitted in working tree):** the Go tree
  is QUIESCENT (G02 mutator done). RUNNING: the §11.4.209 Fable-xhigh review of the G02 change (gates its
  commit) + the next design wave (G09 OpenAPI drift, G15 Aurora/HarmonyOS client feasibility, G20 auto-expand
  + G25/G26). **Bottleneck honesty (§11.4.6):** the design lane is well ahead; the single-threaded Go-impl
  lane (§11.4.84 one mutator, gated by the Fable review) is the true throughput constraint — more design
  docs are queued-not-abandoned (§11.4.197), not extra Go throughput.
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

## R18 — full documentation-delivery mandate (2026-07-15, captured)

Whole project delivered with: API docs (static + interactive), user manuals, guides,
tutorials, FAQs; all diagrams/schemes/graphs with OpenDesign illustrations (R12/§11.4.162/
§11.4.190); all SQL definitions, templates, materials; **exported to every mandatory format
(§11.4.65)**; ALWAYS up to date + in sync via hooks + the Docs Chain submodule (R10/§11.4.106).
Architecture design in flight (`research/` R18 doc). Composes R10 + R12 + §11.4.12/.44/.45/
.53/.56/.57/.59/.60/.65/.106/.168/.170/.190.
