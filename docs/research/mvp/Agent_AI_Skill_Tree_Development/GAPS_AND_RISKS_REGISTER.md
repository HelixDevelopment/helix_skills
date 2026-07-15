# GAPS_AND_RISKS_REGISTER — HelixKnowledge Skill Graph System

> Adversarial audit satisfying operator mandate **R17**. Every row carries
> concrete `file:line` evidence (positive-evidence-only, R11). Anything not
> directly verified is labelled **UNCONFIRMED**. Audit date: 2026-07-15.
> Scope audited: `IMPLEMENTATION_PLAN.md`, `REQUIREMENTS.md`, `SPEC.md`,
> `api/openapi.yaml`, and the full `project/` Go backend
> (`github.com/helixdevelopment/skill-system`, 53 `.go` files, 0 tests).
>
> **Method note:** `go build`/`go vet` are reported green by the task; this
> audit did not re-run them. Findings are about *design, behaviour, wiring,
> security, and contract fidelity*, not compilation.

## Summary counts

| Severity | Count | IDs |
|---|---|---|
| **CRITICAL** | 4 | G01, G02, G03, G04 |
| **HIGH** | 11 | G05, G06, G07, G08, G09, G10, G11, G12, G13, G14, G15 |
| **MEDIUM** | 8 | G16, G17, G18, G19, G20, G21, G22, G23 |
| **LOW** | 4 | G24, G25, G26, G27 |
| **TOTAL** | **27** | |

### Headline: the running binary is not the audited/hardened codebase

The single most important structural finding is that **`cmd/server/main.go`
runs its own ad-hoc API and never imports `internal/api`, `internal/validation`,
or `internal/autoexpand`.** The hardened, spec-shaped handlers (with API-key
auth, strict CORS, the "zero-bluff" jury pipeline, and the auto-growth pipeline)
all exist as source but are **dead code / unwired**. The R1 security fixes and
the R2/R8/R11 flagship features are therefore present in the tree but absent
from the artifact — a textbook §11.4.108 SOURCE≠RUNTIME gap. G01/G02/G03 all
flow from this.

---

## CRITICAL

### G01 — Two rival API servers; the hardened `internal/api` is unwired dead code, the live server has no auth + wildcard CORS
- **Category:** inconsistency / security
- **Severity:** critical — the R1-mandated CORS + api_key fixes are *not in the running binary*; the live REST surface is unauthenticated.
- **Evidence:**
  - Hardened server exists: `internal/api/server.go:82-210` (`Server`, `APIKeyAuth`, strict `CORS`), `internal/api/middleware.go:280-310` (`APIKeyAuth`), `internal/api/middleware.go:328-387` (allowlist `CORS`), `internal/api/skills_handler.go` (full CRUD).
  - `internal/api` has **zero importers** (`grep -rln skill-system/internal/api` → no match).
  - The actually-run server is `cmd/server/main.go:140-283` `setupAPI()`, which builds a *second* Gin router with its own `corsMiddleware()` that sets `Access-Control-Allow-Origin: *` **unconditionally** (`cmd/server/main.go:362-373`) and applies **no API-key auth** anywhere.
  - Even the hardened path is fail-open: auth is applied only `if len(s.cfg.APIKeys) > 0` (`internal/api/server.go:163`) — empty key set ⇒ all `/api/v1` open.
- **Why it matters:** OpenAPI declares `ApiKeyAuth` on every endpoint; the deployed server enforces none. Anyone who can reach the port can read/write. It also means the "security clean" P0.T2 gate is satisfiable only against dead code — an anti-bluff (R11) hazard.
- **DECISION:** Delete the `setupAPI()` router in `cmd/server/main.go`; wire `internal/api.Server` as the single REST surface, constructed with a real `Pool` adapter. Make auth **fail-closed**: if no keys configured, refuse to start (or bind loopback-only) rather than serve open. **Alternatives rejected:** (a) keep both and "pick at runtime" — guarantees drift and was the exact R1 dedupe smell; (b) add auth to the ad-hoc router — duplicates the hardened logic a third time.
- **STATUS (2026-07-15) — IN PROGRESS (attempt 1 = SOURCE-only, NO-GO at runtime layer; attempt 2 remediation in flight):** Attempt 1 hardened the `cmd/server` ad-hoc router (deleted `corsMiddleware()`/`*`; wired `internal/api.CORS(cfg.Server.AllowedOrigins)`, empty allowlist = fail-closed; fail-closed auth via shared `ResolveAPIKeyAuth` — no keys + `AuthDisabled=false` ⇒ every `/api/v1` request 503; `server.go:163` fail-open removed; keys from `HELIX_API_KEYS`/`HELIX_AUTH_DISABLED` env per §11.4.10). `go build/vet/test`=0, and the resolver+CORS+3 wiring tests are correct in isolation (full config-space enumerated: **no fail-open cell in the `/api/v1` wiring**). **BUT the §11.4.209 Fable-xhigh review returned NO-GO** on a CRITICAL §11.4.108 SOURCE≠RUNTIME defect I independently confirmed: in `http`/`both`/default modes the process co-binds a SECOND listener — the MCP HTTP router (`internal/mcp/http_transport.go`) — on the **byte-identical** `host:HTTPPort` (main.go:114/123/131 → RunHTTP/RunBoth, vs main.go:278 setupAPI), both in goroutines that **swallow the bind error** (main.go:289-291, http_transport.go:63-66). The MCP router carries unconditional wildcard CORS (http_transport.go:471, SSE :333) and **ZERO auth** on `/mcp/v1/tools/:name/call` — which dispatches the WRITE tools `skill_create`/`learn_from_project` (http_transport.go:89). So the process races two servers for one port; if the MCP one wins the bind (it starts first), the live surface is wildcard-CORS + no-auth and the hardened `/api/v1` routes 404 — the fix is dead at runtime. Attempt 1's "security fixed" claim held at SOURCE only. Also found: (MEDIUM) the new `api_keys` TOML field is excluded from `substituteEnv` `${VAR}` interpolation (a `${PROD_KEY}` placeholder becomes a literal valid credential); (MEDIUM) zero config-layer tests for the new env parsing; (LOW) new knobs undocumented in `config/config.toml`+`.env.example`. **Attempt-2 remediation (in flight):** collapse to ONE HTTP listener (mount MCP routes onto the hardened `setupAPI` router behind the SAME `api.CORS`+fail-closed auth via the existing `HTTPTransport.RegisterRoutes`), fail-HARD on bind error, delete the wildcard MCP `corsMiddleware`, interpolate/reject `${` in the new list fields, add config tests + ops docs; then re-run the Fable review to a zero-finding GO (§11.4.134). **STILL SEPARATELY OPEN — full single-surface consolidation of the DEAD `internal/api.Server`:** its `New` needs a 25-method `api.Pool` impl that nothing provides (probe: 0 impls of `SubmitProject`/`TriggerExpand`/`GetGapReport`/`ImportSkills`/`SearchSkills`); building the `apiadapter` (+ missing `Store` methods + job persistence) is deferred per §11.4.101 rather than stub-faked (§11.4.27/§11.4.108). G01 migrates to `Fixed` only after the runtime-layer security hole is closed AND the Fable review returns GO.
- **Test coverage:** integration (auth 401 on missing/invalid key, fail-closed on empty key set), security (disallowed Origin gets no ACAO header; wildcard never co-occurs with credentials), contract (routes == OpenAPI), regression, smoke, mutation (reintroduce reflect-origin / drop auth → test fails). **Challenges:** yes (end-to-end unauthorised-access probe). **HelixQA:** yes.

### G02 — The "sandbox" provides no isolation: default path executes arbitrary skill code on the host (RCE), with false security claims
- **Category:** security / danger-zone
- **Severity:** critical — latent (see reachability) but a full remote-code-execution primitive by design.
- **Evidence:**
  - Default `sandbox_type = "wasm"` (`internal/config/config.go:179`, `config/config.toml`), and `createDefaultSandbox` maps `wasm`→`NewWASMSandbox` (`internal/validation/pipeline.go:121-135`).
  - When no WASM runtime is found (the default on a bare host) `Execute` falls through to `executeProcess`, which runs code via the host interpreter: `python -c`, `bash -c`, `node -e` (`internal/validation/sandbox.go:86-105`, `206-290`, dispatch `237-245`). The Go path is literally `go run` on the host (`sandbox.go:143-203`).
  - Isolation claims are false: the "restrict network" comment sets `LD_PRELOAD=` (`sandbox.go:257-260`), which does nothing to network; there is no seccomp/namespace/cgroup. Comment even admits "In production, this would use a proper WASM compiler service" (`sandbox.go:118-120`).
  - The pipeline runs **every fenced code block in skill markdown content** regardless of language (`internal/validation/pipeline.go:336-374`, `extractCodeBlocks:383-420`) — documentation snippets get executed.
  - **Reachability (honest boundary):** currently *latent* because `internal/validation` has zero importers (see G03); it becomes live the moment P4 wires `Validate()` into the create/expand/MCP paths (which the plan intends). A skill body is fully attacker-controllable via `POST /skills` and MCP `skill_create`.
- **Why it matters:** "sandbox validation" that runs untrusted input on the host, unauthenticated (G01), is the worst-case security defect; the naming is itself an anti-bluff violation (R11).
- **DECISION:** Before wiring P4, replace the process-fallback with a real isolation boundary (rootless Podman/gVisor per Constitution §11.4.161, or a true WASM runtime with `wasmtime`), and **fail-closed**: if no isolated runtime is available, `SKIP` with reason — never silently execute on the host. Do **not** execute arbitrary documentation code blocks; only execute snippets explicitly tagged as runnable POCs. Delete the `LD_PRELOAD` "network restriction" line. **Alternatives rejected:** (a) keep process fallback "for dev" — one config typo = host RCE; (b) drop code execution entirely — loses the P4.T3 sandbox gate the zero-bluff pipeline needs.
- **Test coverage:** security (network-egress attempt from snippet is blocked; filesystem write outside sandbox blocked; `fork`/resource limits enforced), integration (no-runtime ⇒ SKIP not execute), fuzz (malicious snippet corpus), mutation (remove isolation flag → egress test fails), regression. **Challenges:** yes (escape-attempt bank). **HelixQA:** yes.

### G03 — Flagship pipelines are dead code: `internal/validation` (jury) and `internal/autoexpand` (auto-growth) are never instantiated; worker handlers are stubs
- **Category:** gap
- **Severity:** critical — R2 (dynamic creation), R8/R11 (zero-bluff validation), and the founding "central registrar + fully-automatic auto-growth jury" are unbuilt in the running system.
- **Evidence:**
  - `grep -rln skill-system/internal/validation` → **no match**; `grep -rln skill-system/internal/autoexpand` → **no match**. Neither `NewPipeline` is ever called.
  - Worker "handlers" are explicit stubs: `internal/worker/runner.go:317` (`// Job handlers (stub implementations…)`), `handleAutoExpand` returns `Success:true` with no work (`runner.go:320-335`), `handleValidate` likewise (`runner.go:337-349`), `handleCodeAnalysis` likewise (`runner.go:351-363`).
  - The worker *cycles* only log: `runAutoExpandCycle` iterates and increments a counter but creates nothing (`runner.go:440-481`); `runValidationCycle` only logs (`runner.go:483-507`).
  - The stub comments assert work that never happens ("Actual expansion is done by the autoExpandWorker polling loop", `runner.go:332`; "Actual validation is done by the validationWorker", `runner.go:347`) — bluff-comments (R11).
- **Why it matters:** every skill created (REST or MCP) is written straight to the DB as `draft` with **no** resource-verify / sandbox / jury / cross-ref. The "zero-bluff guarantee" (`internal/validation/pipeline.go:1-4`) is not in force anywhere.
- **DECISION:** Wire `validation.Pipeline.Validate` and `autoexpand.Pipeline.Run` into the worker's real `handleValidate`/`handleAutoExpand` and into the create path; delete the stub comments; the worker cycles must call the pipelines, not log. Gate a "no skill reaches `validated`/`active` without a recorded jury verdict" invariant. **Alternatives rejected:** leaving pipelines as libraries "to be wired later" is the §11.4.197 un-wired-research failure the Constitution forbids.
- **Test coverage:** unit (pipeline stage state machine), integration (draft → jury → merge on real DB), e2e (create → validated only after ≥2 approvals), mutation (a fabricated skill must be rejected; strip a stage → test fails), regression. **Challenges:** yes (fabricated-skill-must-fail). **HelixQA:** yes (autonomous QA session over the pipeline).

### G04 — Zero automated tests exist; the R1 security fixes and every behaviour have no proof
- **Category:** test-coverage
- **Severity:** critical — direct violation of R8 (~100% across 13 test types + Challenges + HelixQA) and R11 (positive evidence only).
- **Evidence:** `find project -name '*_test.go'` → **0 files**. The CORS allowlist (`middleware.go:328-387`) and api_key redaction/removal (`middleware.go:230-250`, `287-292`) compile but have **no** behavioural test. `IMPLEMENTATION_PLAN.md:117` (P0.T2) *promises* such a test; none exists. `REQUIREMENTS.md:113` records "0 tests".
- **Why it matters:** every "green" claim is unbacked; regressions are undetectable; the anti-bluff covenant is unmet at the test layer.
- **DECISION:** Land P0.T3 test bootstrap first, then per-package tables with paired mutations (§1.1) as each package is wired, starting with the security middleware (behavioural proof: disallowed origin rejected, credentials never with `*`, api_key never logged) and the DAG/graph correctness. Coverage gate in CI. **Alternatives rejected:** deferring tests to a final P11 phase — the plan already forbids this (`IMPLEMENTATION_PLAN.md:328`); untested code has been shipping bugs G06/G07/G11.
- **Test coverage:** all 13 types are the deliverable here; prioritise unit + integration + contract + security + mutation first. **Challenges:** yes. **HelixQA:** yes.

---

## HIGH

### G05 — `LLMJury` auto-approves when no jury is configured (default state)
- **Category:** danger-zone / security-of-correctness
- **Severity:** high — the zero-bluff gate defaults to "approve everything".
- **Evidence:** `internal/validation/pipeline.go:428-439` — `if len(p.jury) == 0 { … Consensus: true … "no jury configured, auto-approved" }`. No `ModelProvider` chain exists yet (P3 unbuilt), so the jury slice is empty by default.
- **Why it matters:** contradicts the founding "LLM jury (≥2 approvals)" and R11. A fabricated skill would pass the jury stage unconditionally.
- **DECISION:** Fail-closed — an empty jury when `validation.enabled` is a hard error (or forces `require_human_review`), never an auto-pass. Require `approval_threshold ≥ 2` real votes. **Alternatives rejected:** "auto-pass in dev" — same bluff class as G02's process fallback.
- **Test coverage:** unit (empty jury ⇒ fail/blocked, not pass), integration (2-of-3 real approvals required), mutation (flip auto-pass back → test fails). **Challenges:** yes. **HelixQA:** yes.

### G06 — `GetDependencyTree` returns only depth-1 children (recursive tree truncated)
- **Category:** existing-bug
- **Severity:** high — the core "recursive dependency DAG" feature is broken; REST `/skills/:id/tree` and MCP `skill_tree` both under-report.
- **Evidence:** `internal/skill/graph.go:280-307` builds `childrenMap` for all depths but attaches **only** `root.Children = childrenMap[rootSkill.ID]` (`graph.go:306`); grandchildren `Children` are never populated. MCP's recursive serializer (`internal/mcp/tools.go:226-246`) therefore also emits a 1-level tree despite recursing. Contrast `GetAllDependencies` (`graph.go:347-371`) which is correct (flat closure).
- **Why it matters:** the founding requirement is "endless, deeply-recursive skill branching"; the tree API silently returns a shallow slice — wrong results presented as complete.
- **DECISION:** Assemble the full tree from `childrenMap` recursively (attach children to every node by ID, cycle-guarded), or select `parent skill_id` in the CTE and build the tree in one pass. Add depth + cycle tests on the seed corpus. **Alternatives rejected:** documenting it as "direct deps only" — contradicts the API contract (`api/openapi.yaml:1355-1370` `SkillTreeNode` is recursive).
- **Test coverage:** unit (tree assembly), integration (seed android closure returns known N-level tree), property (tree node count == closure size), regression, mutation (revert to depth-1 → test fails). **Challenges:** yes.

### G07 — TOML/JSON dependency+resource round-trip is broken (edges silently dropped on import)
- **Category:** existing-bug
- **Severity:** high — breaks R14 git-versionable round-trip and the R6 wizard's DAG mapping.
- **Evidence:**
  - REST `POST /skills` discards `req.Deps` and `req.Resources` entirely — the built `models.Skill` sets only scalar fields (`internal/api/skills_handler.go:165-176`); `CreateDepsRequest`/`CreateResourceReq` (`skills_handler.go:33-45`) are parsed then ignored.
  - REST import: `convertTOMLWrapper` creates `SkillDependency{RelationType:…}` but **never sets `DependsOnName`/`DependsOn`** — the target name is thrown away (`skills_handler.go:548-565`, literally `_ = depName // placeholder for resolution`).
  - Export *does* write `DependsOnName` (`skills_handler.go:611-616`), so export→import is non-idempotent: names go out, come back empty.
  - **UNCONFIRMED:** the MCP create path uses `skillStore.ImportFromTOML` (`internal/mcp/tools.go:275`), a different function in `internal/skill/import_export.go` not read in this pass — needs verification that it preserves edge names.
- **Why it matters:** a skill graph whose *edges* vanish on ingest cannot resolve dependencies, run the recursive CTE, or drive the wizard; the DB drifts from the git-versioned TOML source of truth.
- **DECISION:** Resolve dependency names → IDs at create/import (insert edges via `Store.AddDependency`, which already validates + cycle-checks, `graph.go:21-96`); persist resources; make export→import a proven identity round-trip. **Alternatives rejected:** deferring edge resolution to a later "linker" pass — leaves the graph disconnected in the interim and hides the drift.
- **Test coverage:** unit (convertTOMLWrapper preserves names), integration (create with deps → edges in DB), property/round-trip (export→import→export byte-stable), contract, mutation (drop `DependsOnName` → round-trip test fails). **Challenges:** yes. **HelixQA:** yes.

### G08 — No TOON codec exists; wire format mandate unmet; OpenAPI still on JSON+TOML
- **Category:** gap / spec-drift
- **Severity:** high — R/founding mandate "TOON primary + JSON fallback"; danger-zone #2 (silent fallback).
- **Evidence:** `grep -rin toon --include='*.go'` → **0 matches** in the whole backend. Content negotiation only knows JSON/TOML (`internal/api/middleware.go:114-148`, `response.go`). `api/openapi.yaml` still advertises `application/json` + `application/toml` throughout (e.g. `openapi.yaml:21-29`, `111-113`, `1043-1051`), never `application/toon`. `IMPLEMENTATION_PLAN.md:335` flags the exact risk; P1.T2/P7.T4 are queued not done.
- **Why it matters:** clients told to expect TOON will get TOML/JSON with no signal — the silent-fallback bluff the plan warns about.
- **DECISION:** Implement/vendor a spec-conformant Go TOON codec (`github.com/toon-format/toon`) with its own golden test vectors before advertising `application/toon`; until it exists, the API MUST NOT claim TOON. Revise `openapi.yaml` content-negotiation in the same change (P7.T4). **Alternatives rejected:** "interpret TOON as TOML" — explicitly superseded (`REQUIREMENTS.md:38-40`).
- **Test coverage:** unit (round-trip struct→TOON→struct), contract (golden TOON fixtures byte-for-byte), integration (`Accept: application/toon` returns TOON; unknown → JSON fallback with correct `Content-Type`), fuzz (malformed TOON rejected), mutation (swap codec to JSON → golden test fails). **Challenges:** yes.

### G09 — Pervasive OpenAPI ↔ implementation drift; most documented endpoints are unimplemented or differently shaped
- **Category:** spec-drift
- **Severity:** high — the contract is the "single source of truth" for R3 thin clients (`REQUIREMENTS.md:156`); clients generated from it will not work.
- **Evidence (spec → live route in `cmd/server/main.go:169-243`):**
  - `POST /search` (`openapi.yaml:506-516`) → live is `GET /api/v1/skills/search?q=` (`main.go:182`). Also the hardened `internal/api/server.go:182` registers `GET /search` — both disagree with the POST spec.
  - `GET /registry/missing` (`openapi.yaml:610`) → hardened route is `GET /registry/missing-deps/:id` (`internal/api/server.go:189`); live server has `GET /api/v1/missing` (`main.go:234`).
  - `POST /expand/{name}` (`openapi.yaml:703`) → hardened `POST /expand` no name param (`server.go:198`); live server has **no** expand route at all.
  - `POST /learn` + `GET /evidences/{skill_name}` (`openapi.yaml:807, 874`) → hardened `POST /learn/projects`, `GET /learn/evidences/:skill_id` (`server.go:206-208`); live server has neither.
  - `/health`, `/metrics`, `/version` are documented under `/api/v1` with `ApiKeyAuth` + 401 (`openapi.yaml:36, 913-972`) but live are top-level and unauthenticated (`main.go:151`, `server.go:155-157`).
  - Live server implements only list/search/get/tree/coverage/missing (`main.go:170-243`) — no create/update/delete/import/export via REST.
- **Why it matters:** contract-first codegen (P8) will produce clients calling endpoints that 404 or expect the wrong verb/shape.
- **DECISION:** Treat `openapi.yaml` as authoritative, regenerate handlers/routes to match, and add a contract-test gate (schema-validate every response, assert route table == spec) that fails CI on drift. Do this **after** G01 collapses to one server. **Alternatives rejected:** hand-syncing docs to code — drift returns immediately without a gate.
- **Test coverage:** contract (spec-validation per endpoint + route-parity), integration, regression, smoke, mutation (rename a route → parity test fails). **Challenges:** yes. **HelixQA:** yes.

### G10 — Embedding dimension: no model↔column assertion; `vector(768)` hard-coded; OpenAI vector length unchecked; non-openai/local providers unsupported
- **Category:** danger-zone / inconsistency
- **Severity:** high — the 768/1536/384 conflict is "resolved" only by two unenforced constants that can silently disagree at runtime.
- **Evidence:**
  - SQL columns hard-code `vector(768)` (`migrations/001_initial.up.sql:14, 60`); config default `Dimensions: 768` (`internal/config/config.go:174`). They agree **by coincidence**, not by construction — the migration is not templated from config.
  - No startup assertion that `embedder.Dimensions() == column dim` exists (danger-zone #3, `IMPLEMENTATION_PLAN.md:336`, is unmet). `config.validate` only checks `dimensions > 0` (`config.go:407-409`).
  - `OpenAIEmbedder.Embed` never verifies the returned vector length equals `e.dimensions` (`internal/db/embedding.go:124-143`) — contrast `LocalEmbedder` which does (`embedding.go:257-262`). A model returning 1536 will be handed to a `vector(768)` insert → runtime failure.
  - `NewEmbedderFromConfig` supports only `"openai"`/`"local"` (`embedding.go:294-308`), but `SPEC.md:396` config lists `anthropic` and R7 mandates HelixLLM/LLMsVerifier providers — those configs error at startup.
- **Why it matters:** a single config edit (e.g. switching to a 1536-native model) produces opaque insert failures; the blueprint's 1536 assumption is silently incompatible.
- **DECISION:** Pin **768** as the shipped default (matches SPEC + current column + research "sweet spot", `SPEC.md:12`); make `vector(N)` a migration parameter driven by `config.embedding.dimensions`; add a **startup assertion** that queries `information_schema`/`pg_attribute` for the column's declared dim and fails fast on mismatch; add an OpenAI response-length check; extend the provider factory for HelixLLM/OpenAI-compatible providers (R7). **Alternatives rejected:** (a) support all dims live via runtime re-index — expensive, out of MVP scope; (b) leave 768 hard-coded and undocumented — reproduces the drift.
- **Test coverage:** unit (provider factory incl. helixllm; OpenAI length mismatch rejected), integration (startup fails on dim mismatch; correct dim inserts), contract (config schema), mutation (change column to 1536, keep config 768 → startup assertion fails). **Challenges:** yes.

### G11 — Worker does no real work and can panic the process (unchecked type assertions in a recover-less goroutine)
- **Category:** existing-bug / gap
- **Severity:** high — background auto-growth/validation/review are non-functional; a crash vector exists.
- **Evidence:** stubs at `internal/worker/runner.go:317-368`; cycles that only log (`runner.go:440-507`). `runRegistryReview` does `coverage["total_skills"].(int)` and `coverage["coverage_percentage"].(string)` (`runner.go:518-519`) — unchecked assertions; if `GetCoverage` returns a differently-typed/absent key the goroutine panics, and worker goroutines have **no `recover()`** (`runner.go:375-434`), so the process dies. (The API `Recovery()` middleware, `middleware.go:254-276`, does not cover worker goroutines.)
- **Why it matters:** the "central registrar + fully-automatic auto-growth" is the product's core promise (R2) and it is inert; the panic risk turns a data-shape change into an outage.
- **DECISION:** Implement the handlers/cycles against `autoexpand`/`validation` (see G03); replace unchecked assertions with comma-ok + typed struct returns from `GetCoverage`; add a `recover()` + restart wrapper around every worker goroutine. **Alternatives rejected:** keeping stubs "until P4/P5" — un-wired-research (§11.4.197) and the panic ships regardless.
- **Test coverage:** unit (cycle calls pipeline; coverage type-safety), integration (worker creates a real skill from a seeded gap), chaos (malformed coverage map ⇒ logged error not panic), mutation (reintroduce bare assertion → panic test fails). **Challenges:** yes. **HelixQA:** yes.

### G12 — tree-sitter is a stub: native parsing always fails; regex-only; Kotlin/C# unsupported despite being configured
- **Category:** gap
- **Severity:** high — R2 requires tree-sitter as a *working POC, not a stub*; learn-from-codebase (R2/R6/P5) rests on it.
- **Evidence:** `initNativeParser` **always** returns an error (`internal/codeanalysis/treesitter.go:106-131`); `parseNative`/`extractImportsNative`/`extractFunctionsNative`/`extractClassesNative` all return `"not implemented"` (`treesitter.go:160, 230, 235, 240`). Only regex fallback runs. `compilePatterns` has **no `kotlin` or `csharp` case** (`treesitter.go:264-296`), yet `kotlin` is in the default analysis languages (`config.go:194`) and normalizeLanguage maps `kt`→`kotlin` (`treesitter.go:558-559`) — Kotlin files yield an empty pattern set ⇒ zero extraction.
- **Why it matters:** "learn from real codebases" over Java/Kotlin/C++ (R13 corpus) will silently extract little/nothing for Kotlin and rely on brittle regex for the rest — evidence quality the jury cannot trust.
- **DECISION:** Land a CGO tree-sitter build (grammars for the R13 languages) behind a build tag, with the regex parser as an explicit, labelled fallback that reports reduced fidelity (never silently); add Kotlin/C# patterns to the fallback in the interim. Prove extraction on a real repo (P5.T2 gate). **Alternatives rejected:** shipping regex-only as "tree-sitter" — an anti-bluff/naming violation (R11).
- **Test coverage:** unit (per-language extraction incl. kotlin), integration (parse a real Android/Kotlin repo → real symbols), fuzz (malformed source doesn't crash), mutation (remove a grammar → extraction test fails). **Challenges:** yes.

### G13 — Two rival `docker-compose.yml` files (rival-copy risk)
- **Category:** ops
- **Severity:** high — P12.T4 explicitly requires one canonical compose; two exist now.
- **Evidence:** `project/docker-compose.yml` (9198 bytes) **and** `project/deploy/docker-compose.yml` (3332 bytes) both present. `IMPLEMENTATION_PLAN.md:258` (P12.T4) calls for exactly one; the scripts/systemd unit must reference a single file.
- **Why it matters:** divergent compose files drift (ports, image tags, env, volumes); operators/scripts pin the wrong one; reproducibility breaks.
- **DECISION:** Choose `project/deploy/` as the canonical ops home (it already carries `.env.example`, `systemd/`), delete/merge the root `docker-compose.yml`, and have `scripts/*` + the `systemd --user` unit reference only the canonical path. **Alternatives rejected:** keeping both "for dev vs deploy" — the exact rival-copy anti-pattern; use compose overrides/profiles instead if two modes are truly needed.
- **Test coverage:** smoke (compose up → `pg_isready`, `CREATE EXTENSION vector`, `SELECT 1`), integration (scripts reference the one file), regression (grep gate: no second compose), acceptance (systemctl --user up/down cycle). **Challenges:** yes.

### G14 — Submodule policy conflict (§11.4.28C single-canonical vs operator parent-priority+both-synced) unresolved; all 7 declared deps unvendored
- **Category:** supply-chain / ops
- **Severity:** high — an open governance escalation blocks P3/P7/P8/P10/P11/P13 dependency work.
- **Evidence:** `project/helix-deps.yaml` header asserts §11.4.28C "ONE canonical location … Nested own-org submodule chains are FORBIDDEN"; `REQUIREMENTS.md:76-79` (R9) wants "parent-dir versions have PRIORITY. BOTH copies must always be in sync"; `IMPLEMENTATION_PLAN.md:286` (X1) records the escalation as **open** ("no second-copy `--apply` until resolved"). The manifest declares 7 deps (llms_verifier, helix_llm, helix_agent, embeddings, helix_qa, challenges, docs_chain) all `layout: grouped` — none are vendored yet (`submodules/` absent).
- **Why it matters:** the two policies are literally contradictory (one canonical copy vs two synced copies); acting on either without a decision risks a §9.2 data-safety/duplication mistake and blocks R7/R8/R10.
- **DECISION (recommended reconciliation):** treat the **parent-ecosystem copy as the single logical canonical** and any `submodules/<name>/` as a **read-only mechanical mirror** pinned to the same commit by `sync_submodules.sh` (verify-only by default, `--apply` gated). This satisfies §11.4.28C's "one canonical" *semantically* (one source of truth) while honouring the operator's "parent priority + both in sync" *operationally* (the mirror is derived, never independently edited). Escalate this exact framing to the operator for sign-off before any `--apply`. **Alternatives rejected:** (a) two independently-editable copies — violates §11.4.28C and invites divergence; (b) delete the parent copy and vendor only under `submodules/` — violates the operator's parent-priority mandate.
- **Test coverage:** integration (`sync_submodules.sh` dry-run resolves each dep to exactly one canonical + a pinned mirror), security (fail-closed on unexpected path — already hardened, `REQUIREMENTS.md:197`), regression (`ls-remote` reachability per dep), mutation (introduce a nested own-org submodule → sync fails). **Challenges:** yes.

### G15 — Aurora OS / HarmonyOS client feasibility unproven (highest client risk)
- **Category:** danger-zone / gap
- **Severity:** high — R3 hard-requires Aurora + HarmonyOS clients; the plan itself ranks this the top danger.
- **Evidence:** `IMPLEMENTATION_PLAN.md:334` (danger-zone #1) and P8.T5 (`:215`) both flag Aurora/Harmony as "smallest ecosystem, highest risk"; `REQUIREMENTS.md:160` build order ends "Aurora last (highest risk)". The only proposed path is Flutter-OHOS + the omprussia embedder (`REQUIREMENTS.md:158-160`) — a spike, not a proven build. No client code exists (P8 not started).
- **Why it matters:** if the embedder path is infeasible, R3 cannot be met on those OSes and must be re-scoped honestly, not bluffed.
- **DECISION:** Run the Flutter-OHOS + omprussia spike **early** (before committing the shared-core client architecture), and risk-flag with the **exact blocker** (no bluffed build) if it fails, per §11.4.112 (structural-impossibility gets an evidence-backed won't-fix) rather than a silent gap. Because every surface is a thin OpenAPI/MCP client, the backend contract is the de-risking asset — freeze it first. **Alternatives rejected:** building Aurora last "and hoping" — the plan already warns against; deferring the spike hides the risk until the end.
- **Test coverage:** acceptance (one build artifact per OS or a documented blocker), smoke (thin client hits `/health` + wizard), contract (generated client compiles against OpenAPI). **Challenges:** yes (per-OS build feasibility). **HelixQA:** device-lab dependent — flag as operator-attended where autonomous build is infeasible.

---

## MEDIUM

### G16 — `sandbox_type = "wasm"` never actually uses WASM; Docker sandbox has conflicting mounts
- **Category:** weakness
- **Severity:** medium
- **Evidence:** even with a runtime present, the Go path is `go run` (`internal/validation/sandbox.go:132-136`) and non-Go langs fall to `executeProcess` (`sandbox.go:138-139`); "WASM" is a misnomer. `DockerSandbox.Execute` mounts `-v tmpDir:/tmp:ro` **and** `--tmpfs /tmp:…` (`sandbox.go:407-408`) — the tmpfs shadows the read-only bind, so `go run /tmp/main.go` reads an empty `/tmp` (`sandbox.go:387`).
- **Why it matters:** the isolation story is inconsistent and the Docker path is subtly broken for the Go case.
- **DECISION:** Rename/replace per G02; if Docker is the isolation boundary, mount code at a distinct path (e.g. `/work:ro`) and drop the conflicting tmpfs, or pass code via stdin. **Alternatives rejected:** leaving the dual `/tmp` mount — non-deterministic behaviour.
- **Test coverage:** integration (Docker path runs the mounted file), unit (mount args well-formed), mutation (reintroduce dual /tmp → test fails).

### G17 — Weak/committed default DB password; config validation misses provider/sandbox enums
- **Category:** weakness / security
- **Severity:** medium
- **Evidence:** `internal/config/config.go:167` defaults `Password: "secret"`; `config/config.toml` ships `password = "secret"`. `config.validate` (`config.go:390-433`) never validates `embedding.provider`, `validation.sandbox_type`, `logging.level`, or `mcp.transport` against their allowed sets — a typo (`provder = "opennai"`) fails late, deep in the call stack.
- **Why it matters:** weak default invites deployment as-is; unvalidated enums produce confusing runtime errors.
- **DECISION:** Require the DB password via env with **no** working default (fail-closed if unset in non-dev); validate all closed-set config fields in `validate()`. `deploy/.env` is correctly git-ignored (`.gitignore:21`) and untracked — keep it so. **Alternatives rejected:** documented default password — a standing credential-hygiene risk (§11.4.10).
- **Test coverage:** unit (invalid provider/sandbox/level rejected; empty password rejected in prod mode), security (no secret in tracked files — pre-commit grep), mutation (add an invalid enum → validate fails).

### G18 — CORS allowlist unreachable on the live path; SPEC config sample omits `allowed_origins`
- **Category:** weakness / spec-drift
- **Severity:** medium
- **Evidence:** the hardened, config-driven allowlist lives in `internal/api` (`middleware.go:328-387`, fed by `ServerConfig.AllowedOrigins`, `config.go:59-63`) but the live server uses `corsMiddleware()` wildcard (`cmd/server/main.go:362-373`) and never reads `AllowedOrigins`. `SPEC.md:376-384` config sample has no `allowed_origins`, and `config/config.toml` — the shipped template — likewise omits it, so operators won't know to set it.
- **Why it matters:** once G01 is fixed, a browser client still breaks unless `allowed_origins` is documented and set; today it's wildcard-open.
- **DECISION:** Wire `AllowedOrigins` end-to-end (config→ServerConfig→CORS), document it in `config.toml` + SPEC §8 with a safe example, default empty (fail-closed). **Alternatives rejected:** wildcard default — the security posture R1 removed.
- **Test coverage:** integration (config allowlist honoured), security (non-allowlisted origin blocked), contract (config documents the key).

### G19 — `SPEC.md §8` config sample uses `--` comments (invalid TOML)
- **Category:** spec-drift / doc
- **Severity:** medium (doc only; the real file is correct)
- **Evidence:** `SPEC.md:381, 396, 404, 407, 425` use `-- comment` syntax inside a `config.toml` block; TOML comments are `#`. The actual `config/config.toml` correctly uses `#`. A reader copy-pasting the SPEC sample gets a file that fails `toml.DecodeFile` (`config.go:251`).
- **Why it matters:** nano-detail-docs mandate (founding) is undermined by a sample that won't parse; erodes trust in the spec.
- **DECISION:** Fix the SPEC §8 sample to `#` comments and add a docs lint that TOML-parses fenced `toml` blocks in the docs (composes with Docs Chain, R10). **Alternatives rejected:** leaving it — a latent copy-paste footgun.
- **Test coverage:** unit/lint (parse every ```toml block in docs), regression.

### G20 — Auto-expand fabricates placeholder skills without an LLM; couples to concrete `*OpenAILLM`; drafted resources never persisted
- **Category:** weakness / gap (latent — package unwired per G03)
- **Severity:** medium
- **Evidence:** with `p.llm == nil` (the default until P3), `DraftSkill` returns `createMinimalDraft` (`internal/autoexpand/pipeline.go:209-213`), which stores boilerplate content ("This skill was auto-generated to fill a gap", `pipeline.go:282-311`) as a real skill — fake knowledge (R11). `DraftSkill` type-asserts `p.llm.(*OpenAILLM)` and errors on any other `LLMClient` (`pipeline.go:215-218`), defeating the interface and R7 pluggability. In `Run`, drafted `resources` get a `SkillID` assigned but are **never persisted** (`pipeline.go:401-403`).
- **Why it matters:** once wired, auto-growth would flood the graph with placeholder skills and drop their resources — the opposite of the zero-bluff promise.
- **DECISION:** Never persist a placeholder as a real skill — either produce genuine LLM content or mark the gap as unfilled; program to the `LLMClient` interface (remove the concrete assertion); persist resources in the same transaction as the skill. **Alternatives rejected:** keeping the minimal-draft fallback for "graceful degradation" — degrades into bluff data.
- **Test coverage:** unit (nil LLM ⇒ no placeholder persisted; interface pluggability), integration (draft → resources persisted), mutation (reintroduce placeholder-persist → anti-bluff test fails). **Challenges:** yes.

### G21 — Resource verification is shallow (HEAD-only, best-effort hash, fail-open on fetch errors)
- **Category:** weakness
- **Severity:** medium
- **Evidence:** `verifySingleResource` (`internal/validation/pipeline.go:259-303`) passes on any HEAD `< 400`; the content-hash check only runs when a prior hash exists and returns `nil` (pass) on any GET/read error (`pipeline.go:280-292`). SSRF is possible — arbitrary skill-supplied URLs are fetched server-side with no allowlist.
- **Why it matters:** "source verification" (stage 1 of the zero-bluff pipeline) can be satisfied by any reachable URL; a moved/altered doc without a stored hash passes; skill-controlled URLs enable SSRF.
- **DECISION:** Require a stored hash for `official-doc`/`code` resources, treat fetch/read errors as verification failures (not pass), and add an egress allowlist / block link-local + metadata IPs (SSRF guard). **Alternatives rejected:** HEAD-only reachability as sufficient — it proves nothing about content (R11).
- **Test coverage:** unit (dead URL fails, mismatched hash fails, fetch error fails-closed), security (SSRF to 169.254.169.254 blocked), integration, mutation (flip fail-open back → test fails). **Challenges:** yes.

### G22 — No rate limiting / auth on the live server; body limit only; Brotli flush errors ignored
- **Category:** weakness / performance
- **Severity:** medium
- **Evidence:** P7.T5 (auth + rate limiting, `IMPLEMENTATION_PLAN.md:203`) is unbuilt; the live server has neither (`cmd/server/main.go:140-283`). `MaxBodySize(100MB)` is applied only in the hardened (dead) path (`internal/api/server.go:149`), not live. Brotli `Flush()`/`Close()` return values are discarded (`internal/api/middleware.go:106-107`), so a compression error yields a silently truncated response.
- **Why it matters:** unauthenticated + unthrottled + code-executing endpoints (post-G03 wiring) are a DoS/abuse surface; silent truncation corrupts responses.
- **DECISION:** Add token-bucket rate limiting + the 100MB body cap to the unified server (G01); handle Brotli errors (abort the response on failure). **Alternatives rejected:** relying on an upstream proxy for limits — the app must be safe standalone per the deploy model (systemctl --user, R15).
- **Test coverage:** load (429 over-limit), integration (413 over-size), unit (Brotli error handled), security, regression.

### G23 — Migrations loaded from a cwd-relative path; failure only warns and the server continues
- **Category:** ops
- **Severity:** medium
- **Evidence:** `db.Migrate(ctx, pool, "./migrations")` (`cmd/server/main.go:84`) is cwd-relative; on failure the server logs `Warn` and keeps running (`main.go:85-88`), so it can serve traffic against a schema-less DB and fail every query.
- **Why it matters:** silent boot on a broken schema is a §11.4.108 runtime hazard; running from a different directory skips migrations entirely.
- **DECISION:** Resolve the migrations dir from config/embed (`embed.FS`), and **fail-fast** (exit non-zero) if migrations don't apply. **Alternatives rejected:** warn-and-continue — hides a fatal state.
- **Test coverage:** integration (missing migrations dir ⇒ startup fails), smoke (`migrate up` on fresh pgvector DB, `\d+` verified), regression.

---

## LOW

### G24 — Health/metrics/version unauthenticated; `/metrics` exposes Prometheus internals publicly
- **Category:** security
- **Severity:** low
- **Evidence:** `cmd/server/main.go:151` (`/health`), `internal/api/server.go:155-157` register these outside auth; OpenAPI marks them `ApiKeyAuth`/401 (`openapi.yaml:913-972`). `/metrics` returns full Prometheus exposition (`middleware.go:22-52` counters).
- **Impact:** minor info-leak (internal metrics, versions) to anonymous callers.
- **DECISION:** Keep `/health` open (liveness), but gate `/metrics` behind auth or bind it to a private interface; align `/version` with the contract. **Alternatives rejected:** authing `/health` — breaks orchestrator probes.
- **Test coverage:** security (anonymous `/metrics` denied where required), contract, regression.

### G25 — `RemoveDependency` ignores name-lookup errors → audit log with empty names
- **Category:** weakness
- **Severity:** low
- **Evidence:** `internal/skill/graph.go:103-104` discard the `Scan` error via `_ =`; if a skill is already gone, the audit entry records empty `from`/`to`.
- **Impact:** degraded audit fidelity (R11 evidence trail).
- **DECISION:** Capture names best-effort but record the not-found condition explicitly in the audit detail. **Alternatives rejected:** ignoring silently — weakens the audit trail.
- **Test coverage:** unit (audit detail records missing name), regression.

### G26 — `${VAR:-default}` cannot resolve to an intentionally-empty value; provider/model env-substitution edge cases
- **Category:** weakness
- **Severity:** low
- **Evidence:** `interpolate` treats any unset-or-empty env as "use default" (`internal/config/config.go:342-348`) — an env var explicitly set to `""` falls through to the default, so an operator cannot blank a value via env.
- **Impact:** surprising config behaviour for empty overrides.
- **DECISION:** Distinguish "unset" (`os.LookupEnv`) from "empty" so an explicit empty override is honoured. **Alternatives rejected:** documenting the quirk — still astonishing.
- **Test coverage:** unit (empty-override honoured; unset uses default), regression.

### G27 — `sanitizeTableName` silently strips instead of rejecting; `EmbedAsync` result-channel semantics
- **Category:** weakness
- **Severity:** low
- **Evidence:** `internal/db/vector.go:288-296` strips non-alnum chars rather than rejecting (`"skills; DROP"`→`"skillsDROP"`); safe today because callers pass internal constants only (`vector.go:216, 226`), but a future dynamic caller could hit a wrong-but-valid table name. `EmbedAsync` (`internal/db/embedding.go:359-405`) is correct (buffered to `len(texts)`), noted only as a caller-contract reminder.
- **Impact:** latent foot-gun if table names ever become user-influenced.
- **DECISION:** Reject invalid table names outright (return error) and keep the caller set to a fixed allowlist enum. **Alternatives rejected:** silent stripping — masks programmer error.
- **Test coverage:** unit (invalid table name rejected), security, regression.

---

## Adjudication of the 8 mandated open items

| # | Item | Finding | One-line decision |
|---|---|---|---|
| 1 | Embedding dimension (768/1536/384) | **G10** | Pin **768** default; template `vector(N)` from config; add startup model-dim==column-dim assertion + OpenAI length check; extend providers. |
| 2 | §11.4.28C single-canonical vs operator parent-priority+both-synced | **G14** | Parent copy = the one *logical* canonical; `submodules/<name>` = read-only mechanical mirror pinned by `sync_submodules.sh`; escalate this framing for sign-off before `--apply`. |
| 3 | TOON: no Go codec | **G08** | Implement/vendor a spec-conformant TOON codec with golden vectors before advertising `application/toon`; do not claim TOON until it exists; revise openapi.yaml in lockstep. |
| 4 | Zero automated tests vs R8 ~100% | **G04** | Bootstrap `go test`, then per-package tables + paired mutations, security middleware + DAG first; coverage gate in CI. |
| 5 | Security behavioural proof missing | **G01 + G04** | The CORS/api_key fixes aren't even wired (G01); wire the hardened server, then prove behaviour with security tests + paired mutations. |
| 6 | Aurora/HarmonyOS feasibility | **G15** | Run the Flutter-OHOS + omprussia spike early; risk-flag with the exact blocker (§11.4.112), never bluff a build; freeze the backend contract first. |
| 7 | Silent-failure / TODO / stub patterns | **G02, G03, G11, G12, G20** | Cited: `treesitter.go:130,160,230,235,240`; `runner.go:317,321,332,347`; `skills_handler.go:552`; `pipeline.go:431`; sandbox process-fallback. Fix per each finding. |
| 8 | Rival docker-compose copies | **G13** | Make `project/deploy/` the canonical ops home; delete/merge the root compose; scripts + systemd unit reference the one file. |

## Resolved / non-findings (verified, for the record)

- **DB-layer dedupe (P0.T1) is done:** `go.mod` shows only `pgx v5` + `pgvector-go`; the 5 stray ORMs (ent/gorm/bun/go-pg/sqlx) are gone.
- **`deploy/.env` is NOT a leak:** it is untracked and `.env` is git-ignored (`.gitignore:21`).
- **`config/config.toml` uses valid `#` comments** (only the SPEC §8 *sample* is wrong — G19).
- **Vector columns and config default agree at 768** today (the risk is the absence of an *assertion*, not a current mismatch — G10).
- **Parameterised SQL throughout the `db` layer** (pgx `$1..$n`); the only string-built identifiers are table names, guarded by `sanitizeTableName` against a fixed caller set (residual foot-gun noted in G27, not an active injection).

---

*Positive-evidence-only. Every "is/does" claim above is pinned to a file:line
that was read during this audit; the two UNCONFIRMED sub-points (MCP
`ImportFromTOML` edge fidelity in G07) are labelled as such and require a read of
`internal/skill/import_export.go` to close.*
