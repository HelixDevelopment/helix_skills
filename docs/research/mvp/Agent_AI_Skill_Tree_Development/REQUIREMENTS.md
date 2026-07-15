# HelixKnowledge Skill Graph System — Consolidated Requirements (living doc)

Status: **requirements-gathering + research phase.** Foundation (extracted Go
backend) **does not compile** — see "Baseline". This doc is the single source of
truth as scope evolves; correct it here.

## Vision
A **universal**, self-growing Knowledge Skill Graph for AI CLI agents. Users
request knowledge for arbitrary technologies; the system **creates skills
dynamically on demand**, maps their dependency DAG, validates them (no bluff),
stores them (Postgres + pgvector), and serves them to AI coding agents.

## Canonical source-of-truth (original founding request) — reconciled
The operator supplied the ORIGINAL request used to prepare the MVP. It is the
source of truth. No separate `.txt` attachment is present in the repo — the
extracted MVP (SPEC/plan/research/project) IS the prepared start point.

CONFIRMS (already captured above): recursive skill DAG for AOSP/Android/Java/
Kotlin/C++ + every dependency tech (R2/R13); central registrar + fully-automatic
auto-growth/expansion (R2); exhaustive reviews, zero false/faulty/bluff (R8/R11);
Go core with CLI/TUI/REST (R3 core); full CLI-agent integration via plugin + MCP +
ACP services (R4); learn-from-real-codebase experience pipeline (record→triage→
process→universal knowledge) (R2); deep research incl. CodeGraph under-the-hood +
memory systems (done); vector DB + Postgres; Docker/Podman Compose via **systemctl
user-space** + `scripts/{start,stop,restart,status,install,...}`; docs to nano-detail
w/ diagrams+schemes+graphs, all SQL/template definitions; final **zip + tar.gz**
packaging retrievable to filesystem.

PINS / CORRECTS (update the plan):
- ⚠️ **Wire format = TOON, JSON fallback (NOT TOML).** Original: "Toon instead of
  JSON with JSON capability as fallback." MVP SPEC §2/§6 + the committed
  `api/openapi.yaml` used JSON+TOML → API serialization must become **TOON primary +
  JSON fallback**; openapi.yaml content-negotiation REVISION QUEUED. `config.toml`
  stays TOML. DEFAULT (override-able): on-disk skill files stay TOML+Markdown
  (human-editable, git-versionable R14).
  → **TOON is the REAL format `github.com/toon-format/toon`** (operator-confirmed via
  a direct repo link; `git ls-remote` verified reachable — default branch `main`, HEAD
  a19a117). The blueprint TXT's "interpret Toon as TOML" is **SUPERSEDED** — we
  implement/vendor a **Go TOON codec** and serve `application/toon` (primary) +
  `application/json` (fallback). Queued: TOON serialization layer + openapi.yaml revision.
- Transport: **Quic/HTTP3 + Brotli** (aligns; supersedes MVP HTTP2-default).
- **ACP** (Agent Client Protocol) confirmed alongside MCP + plugins.
- **Endless, deeply-recursive** skill branching: every technology appearing in any
  dependency chain gets its own Skill, recursively, without bound.
- Explicit deliverable: a **feasibility + step-by-step methodology + danger-zones**
  document ("Can we do this? Here's exactly how"), backed by deep web research.
- Skills must be directly usable by Claude Code / OpenCode / other CLI agents.

## Requirement clusters (from operator, in order received)
- **R1 — Standard.** Build+run clean; 4-layer test coverage + paired mutations
  (Constitution §1/§1.1/§11.4); fix flagged security defects (CORS reflect+credentials;
  api_key in query/log); dedupe `project/` vs `skill-system/` to one canonical tree.
- **R2 — Universal + dynamic.** Not Android-specific; create skills on demand for
  any technology set. Mentioned tech (Gin, quic-go, pgvector, tree-sitter, mcp-go,
  Bubble Tea, Cobra, brotli, TOML) must be **working POCs**, not stubs.
- **R3 — Clients.** CLI, TUI, REST API, **Web**, **Desktop (Win/macOS/Linux)**,
  **Mobile: Android, iOS, HarmonyOS, Aurora OS (auroraos.ru)**. **Maximize shared
  codebase** across all surfaces. (Web/Desktop/Mobile are greenfield.)
- **R4 — Agent interop.** Usable from Claude Code (toolkit + aliases), OpenCode,
  Kimi Code, and **HelixTrack / HelixAgent / HelixLLM** (`../helix_code/`).
- **R5 — Incorporation analysis.** In-depth research on how to integrate all of R4
  properly (→ 4 research agents in flight).
- **R6 — Canonical use case (wizard).** User opens a client → **skill-creation
  wizard** → enters a tech set (example: `android, android-aosp, java, kotlin,
  c++, cmake`) → submit → backend runs **create → map (DAG) → full processing**
  (generate content, resolve deps, validate via jury, embed, store) → progress
  reported back to the client.
- **R7 — Model access.** Obtain quality models via **LLMsVerifier**, **HelixLLM**,
  **Claude Toolkit aliases** — pluggable ModelProvider, not hardcoded OpenAI.
  Any required dependency not present in the parent dir must be **vendored as a
  submodule at `submodules/<snake_case_name>/`** per Constitution naming.
- **R8 — Exhaustive testing.** Every unit covered by ALL supported test types PLUS
  **Challenges** and **HelixQA test banks/suites**; coverage as close to **100%** as
  possible. (Helix ships `helix_qa` + `challenges` submodules and a `helixqa`
  binary — these are real components, not abstractions. Exact rules ← research agent 5.)
- **R9 — Submodule resolution + sync.** All deps and dependency submodules live under
  `submodules/<snake_case>/`, OR reuse the versions from the parent dir/submodules —
  **parent-dir versions have PRIORITY**. BOTH copies must always be in sync with
  main/master. (Exact Constitution rule ← research agent 5.)
- **R10 — Docs Chain.** Fully incorporate the **Docs Chain** submodule (`.docs_chain`)
  per Constitution. (Rules ← research agent 5.)
- **R11 — Zero bluff, anywhere.** No false results, faulty results, or faulty codebase;
  no bluff of any kind or form. Positive-evidence-only everywhere. This governs every
  gate, test, and status report (Constitution §7.1/§11.4/§1.1).
- **R12 — OpenDesign for all design.** Every client's design/styling/diagrams/
  illustrations MUST use **OpenDesign** (`git@github.com:nexu-io/open-design.git`).
  Deliver all Constitution-mandated design artifacts: wireframes, sketches, Figma,
  and exports in PDF, PSD, SVG (+ other Constitution-defined types). [Incorporation
  research: NEXT WAVE.]
- **R13 — Validation skill corpus.** The system must successfully create + validate
  this corpus (this IS the end-to-end proof): android, android-aosp, rockchip,
  orange-pi, java, c, c++, kotlin, python, postgres, bash, go, gin-gonic, flutter,
  angular, typescript, javascript, major web/mobile/desktop frameworks (all OSes),
  linux, macos, debugging, cmake, make, gcc, bazel, brotli, quic, http3, http,
  protocols, design-patterns, algorithms, security, snyk, sonarqube, maven, gradle.
- **R14 — Git-versioned + real-time growth.** Every created/obtained skill is
  **Git-versionable** — persisted as versioned files so the main repo grows as
  knowledge expands (TOML per SPEC §4.2, committed; kept in sync with the
  Postgres+pgvector index via import/export). **MCP / ACP / plugins** must trigger
  real-time skill creation + deep research, expanding/updating the knowledge base
  and skill graph live, improvements available immediately.

### Architecture notes forced by R13/R14
- Dual persistence: skills live as **versioned TOML files in-repo** (source of truth
  for git history) AND in **Postgres+pgvector** (query/search index); a sync path
  keeps them consistent (SPEC §6 import/export is the seam).
- The MCP server is not just read: it exposes **create/expand/learn** tools that
  kick off (async) generation + research jobs and stream progress — the R6 wizard
  and R14 real-time growth run through the same job pipeline.

## Baseline (captured evidence, 2026-07-15)
- Go 1.26.4 present. `go.mod` = `github.com/helixdevelopment/skill-system`, 18 deps,
  no `go.sum` (never resolved). 53 `.go` files, **0 tests**.
- `go build ./...` → **FAILS**: duplicate `color*` consts (common.go vs skill.go);
  unused imports (`strconv`, `context`); undefined `pgvector.RegisterTypes`,
  `stdlib.ConnPool`, `stdlib.GetPool`; `tls.Config.Allow0RTT` unknown field.
- `go mod tidy` pulled **five competing ORMs** (ent, gorm, uptrace/bun, go-pg,
  sqlx) alongside pgx — incoherent generated skeleton.
- Security: CORS reflect-origin + `Allow-Credentials: true` (HIGH); api_key via
  query param then logged (MEDIUM). Present in both duplicated copies.

## Global constraints (Constitution)
No bluffing / positive-evidence only; every gate paired with a mutation; no
guessing language; credentials never committed (templates only); multi-upstream
push; deep understanding + captured evidence before implementation; stop-and-fix
at root cause on any discovered defect.

## Open architecture decisions (pending research → operator sign-off)
1. **Shared-core tech** for CLI/TUI/Web/Desktop/Mobile(incl. Aurora/Harmony) —
   maximize reuse (research agent 3).
2. **Model-access abstraction** — LLMsVerifier/HelixLLM/alias backends (agent 4).
3. **Agent-interop layer** — one MCP server + CLI + aliases vs bespoke plugins
   (agent 2).
4. **Helix ecosystem integration** — HelixTrack/HelixAgent/HelixLLM surfaces
   (agent 1).
5. **Canonical location** for the deduped service tree (top-level `skill-system/`
   vs under this research folder).

## Research findings — CONFIRMED (agents 1-3 of 5, cited)
- **Interop is solved by MCP.** Claude Code, OpenCode, and Kimi Code all speak MCP
  natively over stdio. Ship **ONE stdio MCP server + a thin CLI + shell aliases** —
  no bespoke plugins. Register once in claude_toolkit's shared plugins
  (`~/.claude-shared/plugins/`), add a sibling entry in `helix_code/opencode.jsonc`
  (already has a `codegraph` MCP entry), and `mcpServers` for Kimi. Optional `SKILL.md`.
- **Models (R7, partial).** HelixAgent (`:7061`) and HelixLLM (`:8443`) are both
  **OpenAI-compatible** → one provider-agnostic client, swap base URL; use
  `/v1/embeddings` (768d) + `/v1/chat/completions`; `llms-verifier`/`helixqa` are
  PATH binaries. HelixAgent/HelixLLM are submodules under `helix_code/submodules/`.
  (LLMsVerifier specifics + ModelProvider abstraction ← agent 4, running.)
- **Helix ecosystem.** HelixTrack = Go/Gin REST `:8080` JWT **sibling** repo → integrate
  via REST + register our MCP server in `helix_track/.mcp.json`. Plug-in idiom for
  owned components: submodule + Makefile `build` + PATH binary + `.mcp.json` entry +
  `internal/<name>/` REST adapter.
- **Shared-core = contract-first thin clients.** Every surface is a thin client over
  the REST/MCP API; genuine shared logic is only ~15-30%. Make **OpenAPI + the MCP
  schema the single source of truth**, codegen clients, enforce with contract tests.
  Per surface: **CLI+TUI = Go** (reuse existing); **GUI apps = Flutter** — the ONLY
  framework covering Android + iOS + **HarmonyOS** (Flutter-OHOS) + **Aurora OS**
  (omprussia embedder) + desktop from one codebase; **Web = React**. Build order:
  CLI → Web → Desktop → Android/iOS → TUI → HarmonyOS → **Aurora last** (highest risk).
- Still running: agent 4 (model access), agent 5 (Constitution: test taxonomy /
  submodule sync / Docs Chain).

## Draft staged plan (to be finalized after research; foundation-first)
- **S0 Research** (in flight): 4 agents — Helix ecosystem, agent interop,
  shared-core, model access.
- **S1 Foundation**: make backend compile & run; pick ONE db layer (pgx); remove
  the 4 stray ORMs; fix security defects; dedupe to one tree; `go build ./...` +
  `go vet` green; docker-compose up (Postgres+pgvector) smoke test.
- **S2 Core correctness + tests**: 4-layer tests (unit/integration/e2e/contract)
  + paired mutations per Constitution; DAG logic, recursive CTE, search proven.
- **S3 Model + processing**: ModelProvider abstraction; wizard pipeline
  (create→map→validate→embed→store) end-to-end (R6).
- **S4 Interop**: MCP server + CLI + aliases; Helix* integration (R4).
- **S5 Shared core + clients**: shared-core lib → CLI/TUI/REST, then Web, Desktop,
  Mobile (Android/iOS first; Harmony/Aurora per research risk).
- **S6 Hardening/docs/packaging**: security, ops, docs to nano-detail, packaging.

Each stage lands only with captured build/test evidence; SDD dispatches per task.

## Addenda — 2026-07-15 (new operator mandates, folded in)
- **TOON correction (serialization):** supersedes blueprint "Toon = TOML". Real
  format = `github.com/toon-format/toon` (ls-remote OK, branch main, HEAD a19a117).
  API wire = `application/toon` primary + `application/json` fallback; needs a Go
  TOON codec. openapi.yaml content-negotiation revision QUEUED.
- **R15 — systemctl (user scope) + ops scripts:** system MUST integrate via
  `systemctl --user` with a user-scope unit, plus `scripts/` containing at least
  start, stop, restart, status, install (and uninstall/logs). Docker/Podman Compose
  (pgvector/pgvector:pg16) orchestrated underneath. Matches blueprint deploy
  section; now an explicit hard requirement. **Stream B dispatched.**
- **R16 — skill granularity & composition:** big technologies MUST be decomposable
  into smaller atomic skills wrapped by an umbrella/composite skill, with a
  first-class relationship mechanism (typed edges: composes/part_of, depends_on,
  requires, extends, prerequisite, related_to, alternative_to) usable as building
  blocks for whole stacks. Research + data-model design **Stream A dispatched**
  (research/skill_granularity_and_composition.md).
- **Security:** sync_submodules.sh hardened (fail-closed validation; paired attack
  proof) — committed c473d01.
- **Seed corpus:** R13 validation corpus + 8 real seed skills committed 0e0bc3b
  (43 nodes / 56 edges, DAG acyclic, independently re-verified).
