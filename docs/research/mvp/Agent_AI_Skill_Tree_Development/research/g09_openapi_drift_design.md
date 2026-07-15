# G09 — OpenAPI ↔ implementation drift: reconciliation + contract-test gate design

**Revision:** 1
**Last modified:** 2026-07-15T17:00:00Z
**Description:** Design-research for gap **G09** (pervasive OpenAPI↔handler drift):
a git-`255061b`-verified drift table of every documented REST endpoint vs the live
route, the source-of-truth decision (spec authoritative → regenerate handlers), a
CI-failing contract-test gate (per-response schema-validation + route-table parity),
the §11.4.108 runtime signature, the enumerated test-case count reconciled against
`research/testing_infrastructure_plan.md`, and honest gaps.
**Authority / mandates served:** R3 (contract-first thin clients, P8), R11
(positive-evidence-only), R17 (every gap heavily covered), Constitution §11.4.8
(deep-web-research-before-implementation), §11.4.108 (runtime-signature-as-DoD),
§11.4.107(10) (self-validated analyzer), §11.4.186 (cross-document divergence gate),
§11.4.150, §1.1.
**Status:** design-research — NO code landed. Read-only. No Go source or other doc
modified; this is the ONE new file.
**Baseline ref:** every `file:line` below is verified against committed ref
`255061b` via `git show 255061b:<path>` (a live Go mutator holds an uncommitted G02
change in the working tree — the working tree was NOT read; the route/shape drift at
`255061b` is authoritative for this analysis).
**Non-contention:** a separate R19 stream owns Anthropic-provider support
(`research/r19_anthropic_api_support_design.md`). This doc's scope is the system's OWN
OpenAPI contract staying in sync with the REST handlers — distinct file, distinct
concern.

---

## 0. Honest baseline (captured from `255061b`)

- **Spec:** `api/openapi.yaml` (1795 lines, OpenAPI **3.1.0**, `openapi.yaml:1`),
  server base `http://localhost:8080/api/v1` (`openapi.yaml:36`), **global**
  `security: [ApiKeyAuth: []]` (`openapi.yaml:41-42`) → every documented path is
  under `/api/v1` and `X-API-Key`-authenticated. **24 documented operations** across
  21 declared paths (`openapi.yaml:59,204,339,402,461,485,531,565,604,628,646,682,
  713,764,796,836,885,916,958,983,1002`).
- **Live surface (post-G01 single hardened listener):** `cmd/server/main.go`
  `buildRouter` (`main.go:140-284`) — the ONE `ListenAndServe` (`main.go:319`).
  It registers **7 REST routes + `GET /`** + mounts the MCP `/mcp/v1/*` group
  (`main.go:269`). This is the only surface actually served.
- **Hardened rival server:** `internal/api/server.go` `SetupRoutes`
  (`server.go:154-216`) declares a near-complete CRUD surface but is **dead code**
  — zero importers, its own `ListenAndServe`/HTTP/3 unwired in every binary
  (G01 STATUS "O3" is still OPEN, `GAPS_AND_RISKS_REGISTER.md:49`). Its routes ALSO
  disagree with the spec (verb/path/param), so it cannot be adopted as-is either.
- **G02 composition (uncommitted, do-not-read):** a separate working-tree change adds
  fail-closed validate-before-persist to the create handlers but adds/removes NO
  routes; the route/shape drift analyzed here is unchanged by it.

The load-bearing fact (register headline, `GAPS_AND_RISKS_REGISTER.md:134`): the live
server "implements only list/search/get/tree/coverage/missing … no
create/update/delete/import/export via REST." A contract-first client generated from
`openapi.yaml` (P8/R3) calls endpoints that 404 or expect the wrong verb/shape.

---

## 1. Drift table — every documented endpoint vs the live route

Classification set: `matches` (verb+path+shape) · `wrong-verb` · `wrong-path` ·
`wrong-shape` (verb+path agree, response body / query-params disagree) ·
`unimplemented` (no live route) · `live-but-undocumented`. "Live" = `cmd/server/main.go`
`buildRouter` (the served surface). The hardened `internal/api/server.go` route is
cited where it exists (dead, but the intended O3 target).

| # | Documented (spec, `api/openapi.yaml`) | Live route (`cmd/server/main.go` `buildRouter`) | Hardened dead route (`internal/api/server.go`) | Class |
|---|---|---|---|---|
| 1 | `GET /api/v1/skills` — `listSkills`, 200 `SkillListResponse{total,limit,offset,skills}` (`openapi.yaml:60`, schema `openapi.yaml:1448-1462`) | `GET /api/v1/skills` (`main.go:190`) returns `{"skills":…,"count":…}` (`main.go:197`) | `GET ""` (`server.go:177`) | **wrong-shape** — verb+path match; body omits required `total/limit/offset`, adds undocumented `count` |
| 2 | `POST /api/v1/skills` — `createSkill`, 201/409/422 (`openapi.yaml:120`) | — none | `POST ""` (`server.go:178`) | **unimplemented** (live) |
| 3 | `GET /api/v1/skills/{name}` — `getSkill`, `?recursive`, 200 `SkillDetail` (`openapi.yaml:207`) | `GET /api/v1/skills/:name` returns `store.GetByName` struct, `recursive` ignored (`main.go:215-224`) | `GET /:id` (`server.go:179`, param `:id`≠`:name`) | **wrong-shape** — verb+path match; `recursive` query unhonored, body is raw Skill not `SkillDetail` |
| 4 | `PUT /api/v1/skills/{name}` — `updateSkill` (`openapi.yaml:264`) | — none | `PUT /:id` (`server.go:180`) | **unimplemented** (live) |
| 5 | `DELETE /api/v1/skills/{name}` — `deleteSkill`, 204 (`openapi.yaml:326`) | — none | `DELETE /:id` (`server.go:182`) | **unimplemented** (live) |
| 6 | `GET /api/v1/skills/{name}/tree` — `getSkillTree`, `?depth`(1–20, def 5), 200 `SkillTreeNode` (`openapi.yaml:340`) | `GET /api/v1/skills/:name/tree`, `depth` hard-coded 5, query ignored (`main.go:226-238`) | `GET /:id/tree` (`server.go:183`) | **wrong-shape** — verb+path match; `depth` query unhonored |
| 7 | `POST /api/v1/skills/import` — `importSkills` (`openapi.yaml:403`) | — none | `POST /import` (`server.go:184`) | **unimplemented** (live) |
| 8 | `GET /api/v1/skills/{name}/export` — `exportSkill` (`openapi.yaml:462`) | — none | `GET /:id/export` (`server.go:185`) | **unimplemented** (live) |
| 9 | `POST /api/v1/search` — `searchSkills`, body `SearchRequest` (`openapi.yaml:486`) | `GET /api/v1/skills/search?q=` (`main.go:200`) | `GET /search` (`server.go:189`) | **wrong-verb + wrong-path** — spec POST `/search`; live GET `/skills/search`; hardened GET `/search` (also wrong-verb) |
| 10 | `POST /api/v1/search/similar` — `searchSimilarSkills` (`openapi.yaml:532`) | — none | `POST /search/similar` (`server.go:190`) | **unimplemented** (live) |
| 11 | `GET /api/v1/registry` — `getRegistry` (`openapi.yaml:566`) | — none | `GET ""` on `/registry` (`server.go:195`) | **unimplemented** (live) |
| 12 | `GET /api/v1/registry/missing` — `getMissingDependencies` (`openapi.yaml:605`) | `GET /api/v1/missing` (`main.go:252`) | `GET /registry/missing-deps/:id` (`server.go:196`) | **wrong-path** — spec `/registry/missing`; live `/missing`; hardened `/registry/missing-deps/:id` (also wrong-path + spurious `:id`) |
| 13 | `GET /api/v1/registry/stale` — `getStaleSkills` (`openapi.yaml:629`) | — none | `GET /registry/stale` (`server.go:197`) | **unimplemented** (live) |
| 14 | `POST /api/v1/registry/review/{name}` — `triggerRegistryReview` (`openapi.yaml:647`) | — none | `POST /registry/review/:id` (`server.go:198`, `:id`≠`{name}`) | **unimplemented** (live) |
| 15 | `GET /api/v1/registry/coverage` — `getCoverageReport`, 200 `CoverageReport` (`openapi.yaml:683`) | `GET /api/v1/coverage` (`main.go:240`) | `GET /registry/coverage` (`server.go:199`) | **wrong-path** — spec `/registry/coverage`; live `/coverage` |
| 16 | `POST /api/v1/expand/{name}` — `triggerExpansion` (`openapi.yaml:714`) | — none | `POST /expand` (`server.go:205`, no `{name}`) | **unimplemented** (live) — hardened also **wrong-path** (drops `{name}`) |
| 17 | `GET /api/v1/expand/status/{id}` — `getExpansionStatus` (`openapi.yaml:765`) | — none | `GET /expand/status/:id` (`server.go:206`) | **unimplemented** (live) |
| 18 | `POST /api/v1/expand/gap-report` — `generateGapReport` (`openapi.yaml:797`) | — none | `GET /expand/gaps` (`server.go:207`) | **unimplemented** (live) — hardened also **wrong-verb+wrong-path** (`GET /gaps`≠`POST /gap-report`) |
| 19 | `POST /api/v1/learn` — `submitProjectForLearning` (`openapi.yaml:837`) | — none | `POST /learn/projects` (`server.go:213`) | **unimplemented** (live) — hardened also **wrong-path** (`/learn/projects`≠`/learn`) |
| 20 | `GET /api/v1/learn/status/{id}` — `getLearnStatus` (`openapi.yaml:886`) | — none | `GET /learn/status/:id` (`server.go:214`) | **unimplemented** (live) |
| 21 | `GET /api/v1/evidences/{skill_name}` — `getEvidenceForSkill` (`openapi.yaml:917`) | — none | `GET /learn/evidences/:skill_id` (`server.go:215`) | **unimplemented** (live) — hardened also **wrong-path** (`/learn/evidences/{skill_id}`≠`/evidences/{skill_name}`) |
| 22 | `GET /api/v1/health` — `getHealth`, ApiKeyAuth + `401` (`openapi.yaml:959`) | `GET /health` top-level, **open/no-auth** (`main.go:163`) | `GET /health` top-level, open (`server.go:159`) | **wrong-path** — spec `/api/v1/health` authed; live top-level `/health` unauthenticated (auth-posture drift too) |
| 23 | `GET /api/v1/metrics` — `getMetrics`, ApiKeyAuth + `401` (`openapi.yaml:984`) | — none | `GET /metrics` top-level, open (`server.go:160`) | **unimplemented** (live) — hardened exists but top-level+open (**wrong-path** + G24 auth drift) |
| 24 | `GET /api/v1/version` — `getVersion`, ApiKeyAuth + `401` (`openapi.yaml:1003`) | — none | `GET /version` top-level, open (`server.go:161`) | **unimplemented** (live) — hardened wrong-path+open |

### 1.1 Live-but-undocumented routes
| Live route | Location | Note |
|---|---|---|
| `GET /` server info | `main.go:272` | **live-but-undocumented** — not in `openapi.yaml`; allowlist candidate (see §3.2) |
| `GET /health` (top-level) | `main.go:163` | **live-but-undocumented** at that path — spec documents `/api/v1/health` (row 22) |
| `GET /api/v1/skills/search` | `main.go:200` | live target of row 9 (spec `POST /search`) |
| `GET /api/v1/missing` | `main.go:252` | live target of row 12 (spec `/registry/missing`) |
| `GET /api/v1/coverage` | `main.go:240` | live target of row 15 (spec `/registry/coverage`) |
| `/mcp/v1/*` group (messages, sse, tools, prompts) | mounted `main.go:269` | MCP transport — governed by the SEPARATE MCP tool schema (SPEC.md §7), NOT `openapi.yaml`; allowlisted out of REST parity (honest boundary) |

### 1.2 Drift summary (24 documented REST operations)
- **matches (verb+path+shape): 0.**
- **wrong-shape: 3** (rows 1, 3, 6 — the only verb+path agreements, each with body/query drift).
- **wrong-verb + wrong-path: 1** (row 9, search).
- **wrong-path: 3** (rows 12, 15, 22).
- **unimplemented in the live surface: 17** (rows 2,4,5,7,8,10,11,13,14,16,17,18,19,20,21,23,24).
- Every hardened `internal/api/server.go` route that *does* exist ALSO drifts on
  path param (`:id`≠`{name}`) or verb (`GET /search`, `GET /expand/gaps`) — the dead
  server is not a drop-in fix either.

---

## 2. DECISION — `openapi.yaml` is authoritative; regenerate handlers to match

**Verdict: the SPEC is the single source of truth; the handlers are regenerated /
rewritten to conform.** (Confirms the register's stated DECISION,
`GAPS_AND_RISKS_REGISTER.md:136`.)

**Justification (R3 / P8 contract-first).** `REQUIREMENTS.md:167-170` records the
CONFIRMED architecture finding: *"Shared-core = contract-first thin clients … Make
OpenAPI + the MCP schema the single source of truth, codegen clients, enforce with
contract tests."* P8.T1 (`IMPLEMENTATION_PLAN.md:253`) gates on a *generated OpenAPI
client that compiles and a smoke call that hits the API*. Every R3 surface (CLI, TUI,
Web, Desktop, Mobile incl. Aurora/HarmonyOS) is a thin client generated FROM the
spec. If the code were authoritative and the spec regenerated from it, the drift would
be *documented* rather than *fixed*, the generated clients would still be shaped by
whatever the handlers happen to return, and the "single source of truth" contract that
the multi-surface client strategy rests on would be inverted. The spec is also the
more complete artifact (24 operations incl. full CRUD/import/export/expand/learn) and
already encodes the intended TOON-primary/JSON-fallback negotiation and the auth
posture (all `/api/v1`, `X-API-Key`); the live surface is a 7-route subset. Authoring
direction therefore flows spec → handlers.

**Alternative rejected** (register, `GAPS_AND_RISKS_REGISTER.md:136`): hand-syncing
docs to code returns to drift the moment either side changes without a gate.

### 2.1 Reconciliation shape (design only)
1. **Generate the server contract from `openapi.yaml`** with `oapi-codegen` in
   **strict-server + gin-server** mode (Go 1.25.5, Gin already in use). Strict mode
   emits typed `…RequestObject`/`…ResponseObject` per operation, so a handler that
   returns the wrong shape does not compile — the drift becomes a **build-time**
   failure for the shape dimension, and a **contract-test** failure for the
   route/status dimension. Generated types + server interface land in a
   `internal/api/gen/` package excluded from the coverage floor
   (`testing_infrastructure_plan.md:165` already reserves the "future OpenAPI client"
   exclusion).
2. **Bind the generated `StrictServerInterface` onto the ONE wired surface.** This is
   the vehicle for G01's still-open **O3** consolidation
   (`GAPS_AND_RISKS_REGISTER.md:49`): the reconciled, spec-shaped handlers become the
   single served REST surface, replacing the ad-hoc `buildRouter` route block; the
   dead `internal/api/server.go` route set (which itself drifts) is retired under
   §11.4.124 investigate-before-remove rather than adopted.
3. **Path-param rename** `{name}`↔`:id` is resolved in the spec's favour: the spec
   keys skills by dotted `name` (`openapi.yaml:1034` `SkillNamePath`), and the live
   store already looks up by name (`store.GetByName`, `internal/skill/store.go:121`).
   The hardened server's `:id` handlers are the drift, not the target.
4. **Compose with G02 create-path validation.** When `POST /skills` (row 2) is
   implemented to the spec, its handler routes the request THROUGH the G02 fail-closed
   validate-before-persist path before write and persists as `draft`; the contract
   test for `createSkill` asserts BOTH the spec status matrix (201/409/422) AND the
   G02 runtime signature (a fabricated/invalid body → 422, never a persisted `active`
   skill). G09 (route/shape parity) and G02 (create-path validation) are orthogonal
   and compose on the same reconciled handler.

### 2.2 Sequencing vs G01 / G02
- **After G01** (done, `255061b`): the single hardened listener already exists, so the
  parity gate has exactly ONE surface to bind to — mandatory precondition, matching the
  register's *"Do this after G01 collapses to one server"*
  (`GAPS_AND_RISKS_REGISTER.md:136`).
- **With/after G02:** G02 lands first on the create path (uncommitted now); G09's
  `POST /skills` reconciliation consumes G02's validator. G09 does not block on G02 for
  the read endpoints (rows 1,3,6,9,12,15,22) — those reconcile independently.
- **Before P8:** the parity gate + generated types must be GREEN before any P8 client
  is generated, or the clients inherit the drift.

---

## 3. Contract-test GATE (fails CI on drift)

Two decidable checks compose into one deterministic PASS/FAIL/SKIP verdict, run as a
`//go:build integration` test against the booted wired server + real pgvector
(`testing_infrastructure_plan.md:192` places contract tests in the integration tier).

### 3.1 Check A — route-table parity (structural)
- **Live set:** enumerate `router.Routes()` (Gin exposes `[]gin.RouteInfo{Method,Path}`)
  from the booted `buildRouter`; normalize Gin `:name`→`{name}`.
- **Spec set:** walk `openapi.yaml` paths×methods via `kin-openapi` `openapi3.T`
  (`loader.LoadFromFile` → `doc.Paths`), prefix each with the `/api/v1` server base.
- **Assertion:** `live_set == spec_set` modulo an explicit, cited **allowlist**
  (`GET /`, `/mcp/v1/*`, and — until row 22 is reconciled — top-level `/health`).
  Any spec-only path (unimplemented) OR any live-only path not on the allowlist FAILs
  and is NAMED in the failure (§11.4.6 actionable). This is the §11.4.108 runtime
  signature (§4).

### 3.2 Check B — per-response schema-validation (behavioral)
- For every documented operation, drive a real request against the booted server with
  `Accept: application/json` and validate the response with `kin-openapi`
  `openapi3filter.ValidateResponse` against the loaded spec (status code + body against
  the operation's response schema). A body that omits a `required` field (row 1's
  missing `total/limit/offset`) or adds an undocumented field under a
  `additionalProperties:false` schema FAILs.
- Request-body validation (`openapi3filter.ValidateRequest`) for the 8 body-bearing
  operations asserts the handler accepts exactly the documented request schema.

### 3.3 Self-validated analyzer (§11.4.107(10))
The parity+schema analyzer itself ships a golden-good / golden-bad / negative-control
triple so the GATE cannot bluff:
- **golden-good:** a fixture router whose routes exactly equal the spec set → analyzer
  PASSes.
- **golden-bad:** a fixture router with one renamed route (`/skills`→`/skill`) and one
  shape-broken response (drop `total`) → analyzer MUST FAIL, naming both offenders.
- **negative-control:** a fixture router equal to the spec PLUS an allowlisted
  `GET /` → analyzer MUST PASS (proves the allowlist is honoured and the check does
  not false-positive on a legitimately-undocumented health/root route).
An analyzer that passes its golden-bad is itself the bluff.

### 3.4 Honest TOON boundary (composes G08)
`kin-openapi` validates JSON/JSON-Schema; it does NOT understand the spec's PRIMARY
`application/toon` media type (`openapi.yaml:22-28`). The G09 contract test therefore
schema-validates the **JSON representation** (`Accept: application/json`, the
documented fallback) against the shared schema, and DELEGATES the TOON-primary fidelity
to G08's golden-TOON-vector round-trip proving TOON⇄struct maps to the same shape
(`testing_infrastructure_plan.md:294`, the G08 row). G09 asserts route+JSON-shape
parity; G08 asserts the TOON wire format over that same shape. Neither is faked as the
other.

---

## 4. RUNTIME SIGNATURE (§11.4.108) — parity on a clean deployment

A fix in this design is DONE only when, on a **freshly-booted wired server** (clean
deployment, no stale binary), the parity report emits:

```
route_parity: PASS   spec_ops=24  live_ops=24  spec_only=0  live_only_unallowlisted=0
schema_validate: PASS  operations=24  invalid=0
allowlist: [ "GET /", "GET|POST /mcp/v1/*" ]   # cited, not silent
```

The single machine-checkable observable = **`spec_only==0 && live_only_unallowlisted==0
&& invalid==0`** read from the booted router's own `Routes()` + live responses — NOT a
grep of the source (a source grep is exactly the §11.4.108 SOURCE-layer bluff this
avoids). Pre-reconciliation this signature is RED (spec_only≈17, several
shape/verb/path failures); post-reconciliation it flips GREEN and stays a permanent
regression guard.

---

## 5. Enumerated test-case count (RED-first) + reconciliation

Every count is a design estimate from the matrix below (no test executed — §5.1). Each
case is authored RED-first against the pre-reconciliation `255061b` surface (the drift
IS the reproduction), flipping GREEN on the reconciled surface — the §11.4.115 polarity
discipline for contract gates.

| Family | Cases | Enumeration |
|---|---:|---|
| **contract — route parity** | 1 | Check A aggregate (§3.1): `live_set == spec_set ∪ allowlist` |
| **contract — response schema** | 24 | Check B per documented operation (§3.2), one per `operationId` |
| **contract — request schema** | 8 | body-bearing ops: createSkill, updateSkill, importSkills, searchSkills, searchSimilarSkills, triggerExpansion, generateGapReport, submitProjectForLearning |
| **contract — error-shape** | 6 | canonical `Error`/status matrix: 400,401,404,409,422,204 each schema-validated |
| **analyzer self-validation (§11.4.107(10))** | 3 | golden-good PASS, golden-bad FAIL, negative-control PASS (§3.3) |
| **integration** | 24 | booted server → each op driven against real pgvector; the 17 currently-unimplemented + create/expand/learn ops SKIP-with-reason (§11.4.3) until their handler/pipeline (G02/G03/G12) lands — never a green pass on an absent route |
| **regression** | 8 | 1 route-parity permanent guard + the 7 register-named mismatches each a named guard: search verb+path, registry/missing path, registry/coverage path, expand-absent, learn-absent, evidences-absent, health path/auth |
| **smoke** | 1 | booted server: non-empty route table; every spec path returns its documented status class (200/401), zero 5xx |
| **paired §1.1 mutation** | 4 | M-G09-1 rename a live route (`/skills`→`/skill`) → parity FAILs; M-G09-2 drop `total` from list handler → schema-validate FAILs; M-G09-3 revert search to `GET /search` → verb-parity FAILs; M-G09-4 add an un-allowlisted live route → parity FAILs |
| **TOTAL** | **79** | |

### 5.1 Reconciliation with `research/testing_infrastructure_plan.md` (§11.4.186)
- **G09 row (`testing_infrastructure_plan.md:295`)** lists the *families* — "contract
  (schema-validate every response + route-table == `openapi.yaml`), integration,
  regression, smoke, mutation (rename a route → parity FAILs)" — but supplies **no
  enumerated count**. This design **EXTENDS** that row with the 79-case breakdown
  above; every family the row names is present (contract ✓, integration ✓, regression
  ✓, smoke ✓, mutation ✓) and the mutation exemplar it names ("rename a route → parity
  FAILs") is M-G09-1.
- **Integration tier (`testing_infrastructure_plan.md:191-192`)** already declares
  "Contract tests live here: every response schema-validates against `api/openapi.yaml`
  and the live route table == the spec (G09)." — this design's Check A + Check B are
  the concrete realization of that line.
- **Post-G09 hook (`testing_infrastructure_plan.md:334`)** reserves "Contract
  route-parity check (live route table vs `api/openapi.yaml`) once G09 lands." — the
  §4 runtime signature is that check.
- **§11.4.186 flags (this design ADDS, the plan does not enumerate):** (a) the
  **analyzer self-validation triple** (golden-good/bad/negative-control) applied to the
  parity checker itself — the plan's §1.3 mandates it generically
  (`testing_infrastructure_plan.md:151-156`) but the G09 row does not list it; (b) the
  **TOON-vs-JSON split** — the G09 contract test validates the JSON representation and
  cross-references G08 for TOON fidelity (§3.4), a boundary the plan's G09 row leaves
  implicit; (c) the explicit **allowlist** (`GET /`, `/mcp/v1/*`) as cited data, not a
  silent skip. These are additive refinements, not contradictions — no plan row is
  weakened.

---

## 6. Honest gaps (§11.4.6)

1. **No test executed.** All 79 counts are a design estimate from the enumerated
   matrix, not RED-first captured evidence. `UNCONFIRMED` until authored and run.
2. **OpenAPI 3.1 tooling not exercised against THIS spec.** `openapi.yaml:1` is
   3.1.0. `kin-openapi` supports 3.1 validation; `oapi-codegen` 3.1 support is partial
   for some JSON-Schema keywords. I did NOT run codegen or the loader against this
   file, so whether the recursive `SkillTreeNode` `$ref` (`openapi.yaml:1415`) and the
   dual `application/toon`+`application/json` content blocks generate/validate cleanly
   is `UNCONFIRMED` — flagged, not assumed. Tracked follow-up before the gate is
   authored.
3. **TOON is not validatable by `kin-openapi`.** The gate can only schema-validate the
   JSON path (§3.4); TOON-primary correctness rests on G08's codec round-trip, which is
   a SEPARATE gap's proof. G09 does not prove TOON wire fidelity and does not claim to.
4. **Reconciliation TARGET surface not finally chosen.** Whether the spec-shaped
   handlers bind onto the reconstructed `buildRouter` or onto a revived
   `internal/api.Server` depends on the still-open G01 **O3** consolidation
   (`GAPS_AND_RISKS_REGISTER.md:49`, deferred per §11.4.101). The parity gate binds to
   whichever becomes the single wired surface; the drift analysis holds either way.
5. **Shape-drift field-diff is partial for rows 3 & 6.** The `wrong-shape` verdict for
   `getSkill`/`getSkillTree` rests on VERIFIED facts — the documented `recursive`/`depth`
   query params are ignored (`main.go:215-224`, `main.go:226-238`) and the handler
   returns the store struct directly — NOT on a field-by-field diff of the store struct
   JSON against `SkillDetail`/`SkillTreeNode`. Full field-diff is `PENDING_FORENSICS`
   (needs schema + struct side-by-side). Row 1's shape drift IS fully field-verified
   (`{skills,count}` vs required `{total,limit,offset,skills}`,
   `main.go:197` vs `openapi.yaml:1450-1462`).
6. **Auth-posture drift (row 22/23/24) overlaps G24.** The documented `401` on
   `/health`,`/metrics`,`/version` under `/api/v1` vs the live top-level open `/health`
   is BOTH a G09 path drift and a G24 auth-gating concern
   (`testing_infrastructure_plan.md:310`); this design records the path/shape dimension
   and defers the "should health be authed?" policy question to G24 rather than
   deciding it here.

---

## Sources verified (§11.4.8 / §11.4.99), 2026-07-15

- `getkin/kin-openapi` — `openapi3filter.ValidateResponse` / `ValidateRequest` validate
  `net/http` traffic against a loaded OpenAPI 3.0/3.1 spec:
  https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3filter ,
  https://github.com/getkin/kin-openapi
- `oapi-codegen/oapi-codegen` — generate Go **gin-server + strict-server** boilerplate
  from OpenAPI 3 (strict interface makes wrong-shape handlers a compile error):
  https://github.com/oapi-codegen/oapi-codegen
- `oapi-codegen/nethttp-middleware` — request-conformance middleware built on kin-openapi:
  https://pkg.go.dev/github.com/oapi-codegen/nethttp-middleware
- Jamie Tanna, "OpenAPI contract testing with Go's `net/http`" (`httptest-openapi`
  wrapper over kin-openapi's filter): https://www.jvt.me/posts/2022/05/22/go-openapi-contract-test/

**Negative finding (§11.4.99(B)):** no Go tool validates a custom `application/toon`
media type against an OpenAPI schema — the TOON leg of the contract is not covered by
any off-the-shelf validator and is delegated to G08's codec proof (§3.4). This is
original-work territory for the TOON dimension.
