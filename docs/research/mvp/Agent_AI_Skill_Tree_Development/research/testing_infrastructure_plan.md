# Testing-Infrastructure & Anti-Bluff Harness Plan — HelixKnowledge Skill Graph System

**Revision:** 1
**Last modified:** 2026-07-15T20:40:00Z
**Description:** Concrete, codebase-grounded plan for the Go backend's test harness,
the 13 R8 test types + Challenges + HelixQA banks, the anti-bluff paired-mutation
proof, the per-gap (G01–G27) coverage matrix (R17), and local gate wiring under a
CI-disabled ecosystem (§11.4.156).
**Authority / mandates served:** R8 (13 test types + Challenges + HelixQA, ~100%
coverage), R17 (every gaps-register finding heavily covered with all applicable
test types, fully validated/verified, zero bluff), R11 (positive-evidence-only),
Constitution §1/§1.1, §11.4.27, §11.4.107(10), §11.4.108, §11.4.169, §11.4.156.
**Scope:** the Go backend `github.com/helixdevelopment/skill-system` under `project/`.
Read-only design doc; no Go source or other doc modified.

---

## 0. Honest baseline (captured evidence, 2026-07-15)

Real state read from the tree — not assumed:

- **Module:** `github.com/helixdevelopment/skill-system`, **Go 1.25.5** (`project/go.mod`).
- **Test files present:** **7 `_test.go` files across 5 packages** (Stream D):
  `internal/api/middleware_test.go` + `internal/api/auth_wiring_test.go`,
  `internal/skill/graph_test.go`, `internal/config/config_test.go`,
  `internal/codeanalysis/analyzer_test.go`,
  `internal/validation/pipeline_test.go` + `internal/validation/sandbox_test.go`.
  (The task's "6" counts `internal/api` once; the tree has two files there.)
- **`go test ./...` → exit 0** captured just now: `internal/api`, `internal/codeanalysis`,
  `internal/config`, `internal/skill`, `internal/validation` all `ok`; the other **11
  packages report `[no test files]`** (`cmd/{cli,cli/commands,server,tui,worker}`,
  `internal/{autoexpand,db,mcp,models,registry,worker}`).
- **Test style already established (do not fork it):** pure stdlib `testing`,
  table-driven, `net/http/httptest` for HTTP, **honest `t.Skip()` with a reason**
  for infra-bound paths (see `graph_test.go` `TestHasCycle_RequiresLiveDatabase`,
  which explicitly refuses to fake a `pgx.Tx` — §11.4.27/§11.4.3). **No testify, no
  sqlmock, no gomock anywhere.**
- **Real infra exists:** `project/deploy/docker-compose.yml` → `pgvector/pgvector:pg16`
  with `pg_isready` healthcheck (Stream B); `project/deploy/systemd/helix-skills.service`;
  `scripts/{start,stop,restart,status,install,uninstall,logs}.sh` + `sync_submodules.sh`
  (hardened) + `migrate.sh` + `package.sh`. **Rival copy still present:**
  `project/docker-compose.yml` (G13).
- **Makefile targets already there:** `test` (`go test -v -race -coverprofile=coverage.out ./...`),
  `test-unit` (`-short`), `test-integration` (`-run Integration`), `coverage`
  (`go tool cover -func`), `coverage-html`, `lint`, `vet`.
- **Deps declared but unvendored** (`project/helix-deps.yaml`, G14): `helix_qa`
  (`HelixDevelopment/HelixQA`), `challenges` (`vasic-digital/Challenges`),
  `docs_chain`, `helix_llm`, `helix_agent`, `llms_verifier`, `embeddings`.
- **Seed corpus for tests:** `seed/skills/*.toml` (8 real skills), `seed/CORPUS.yaml`,
  `seed/validate_dag.py` — the fixture source for integration/e2e closure tests.

The gaps register's own headline is the load-bearing fact for this plan: **the
running binary (`cmd/server/main.go`) is not the hardened codebase** — most tests
today prove *dead code*. Every gate below is written to assert against the **wired
runtime** (§11.4.108), not source presence.

---

## 1. Go test harness layout

### 1.1 Runner + assertion library — recommendation: **stdlib `testing`, keep it**

**Decision: stdlib `testing` + `net/http/httptest` + real `pgx` against the container.
Do NOT add testify/gomock/sqlmock.** Grounded reasons:

- All 7 landed test files are pure stdlib table-driven; adding testify now forks a
  working convention mid-stream (churn with zero correctness gain).
- **§11.4.27 (no fakes beyond unit):** `sqlmock`/`gomock` for the DB/HTTP-client
  seams would test the mock, not the real system — exactly the bluff `graph_test.go`
  already refuses. Integration/e2e MUST hit real Postgres+pgvector via the
  `deploy/` compose. Mocks are confined to *unit* tests of pure logic (the fake
  `ModelProvider` in P3.T1 is the sanctioned exception — a unit-level test double).
- stdlib `t.Run` table tests + `t.Setenv`/`t.Cleanup`/`t.Skip` already cover every
  ergonomic need seen in the codebase; `go test -race` is the concurrency oracle.

Optional additive (no convention change): `net/http/httptest` (already used),
`testing/fstest` for embed-FS migration tests (G23), and Go **native fuzz**
(`func FuzzXxx(f *testing.F)`) for TOON/snippet fuzzing — all stdlib.

### 1.2 Per-package `_test.go` conventions (extend, don't reinvent)

- **File naming:** `<unit>_test.go` in the *same package* for white-box unit tests
  (matches `middleware_test.go` testing unexported `redactQuery`, `CORS`); use an
  `_test` external package only where black-box API framing is clearer.
- **Integration/e2e naming for the Makefile filter:** name integration tests
  `TestIntegration_*` and e2e `TestE2E_*` so `make test-integration`
  (`-run Integration`) and a new `-run E2E` filter select them; guard them with a
  build tag `//go:build integration` **and** a runtime skip:
  ```go
  //go:build integration
  func TestIntegration_AddDependency_PersistsEdge(t *testing.T) {
      dsn := os.Getenv("HELIX_TEST_DATABASE_URL")
      if dsn == "" { t.Skip("no HELIX_TEST_DATABASE_URL; start deploy/ compose first (§11.4.3 SKIP-with-reason)") }
      // real pgx pool → real migrate up → insert → read back identical
  }
  ```
  The skip-with-reason (never a silent pass) is the §11.4.3 discipline the codebase
  already models. A missing DB is a SKIP, **never** a green PASS.
- **Table-driven default**, one behavior per `t.Run` case, error assertions via
  `errors.Is` against sentinel errors (as `graph_test.go` does with
  `ErrCycleDetected`/`ErrInvalidSkill`).
- **Fixtures:** integration/e2e load the real `seed/skills/*.toml` corpus (the known
  android closure) rather than inventing graphs — the fixture already ships and is
  DAG-verified (`seed/validate_dag.py`).
- **Directory additions** (no restructure of `internal/`):
  ```
  project/
    test/
      integration/        # //go:build integration; real Postgres+pgvector
      e2e/                # //go:build e2e; server boot → HTTP → DB round-trips
      challenges/         # Challenge bank (see §3)
      helixqa/            # HelixQA banks + runner glue (see §3)
      fixtures/           # golden TOON vectors, malicious-snippet corpus, OpenAPI route table
    scripts/gates/        # local pre/post/runtime gate scripts (see §5)
    scripts/mutation/     # paired-mutation harness (see §1.3)
  ```

### 1.3 Anti-bluff PAIRED-MUTATION pattern (§1.1) — the proof each gate is not a bluff-gate

**Live proof already in the tree:** `internal/api/middleware_test.go`
`TestCORS_NonAllowlistedOriginIsNeverReflected` carries the explicit contract in its
doc comment — *"If this test is broken by re-introducing 'reflect any Origin'
behavior, it must fail — see the paired-mutation demonstration"*; and
`auth_wiring_test.go` `TestResolveAPIKeyAuth_FailsClosedWhenUnconfigured` asserts the
503 fail-closed path. These are the reference anchors: the mutation for CORS is
"re-introduce reflect-origin in `CORS()` (`middleware.go`)"; for auth it is "make
`ResolveAPIKeyAuth` fail-open on an empty key set."

**Mechanical loop for every gate** (`scripts/mutation/run_mutation.sh <mutation-id>`):

1. **Baseline GREEN** — `go test ./<pkg>` passes on the clean tree.
2. **MUTATE the protection** — apply one surgical diff that re-introduces the exact
   defect the test guards (a registry of `mutation-id → sed/patch` pairs, each with a
   `// MUTATION:<id>` marker so §11.4.84 mutation-residue scanning can catch a leak).
3. **Assert RED** — `go test -run <TestName> ./<pkg>` MUST now **FAIL**; the harness
   records exit≠0 as the proof the test is load-bearing.
4. **RESTORE** — `git checkout -- <file>` (or reverse-patch); working tree
   verifiably clean (`git status --porcelain` empty) before anything else runs
   (§11.4.84 quiescence — no mutation may co-commit).
5. **Assert GREEN again** — re-run; PASS. Captured evidence = the paired
   `mutation_<id>_RED.log` + `mutation_<id>_GREEN.log`.

A gate whose paired mutation does NOT flip it RED is itself a bluff (§11.4.107(10)
applied to gates) and is a release blocker. This is the RED-polarity (§11.4.115)
discipline expressed for gates rather than fixes.

**Analyzer self-validation (§11.4.107(10)):** every *analyzer* we build (the TOON
codec's golden comparator, the sandbox-escape detector, the contract route-parity
checker, the freeze/liveness-style content oracles if any) ships a **golden-good
fixture (MUST PASS) + golden-bad fixture (MUST FAIL) + negative-control (a
legitimately-different-but-valid input that MUST PASS)**. An analyzer that passes its
golden-bad is the bluff.

### 1.4 Coverage gate — threshold + measurement

- **Measured** by `go test -race -coverprofile=coverage.out ./...` +
  `go tool cover -func=coverage.out` (both already in the Makefile).
- **Gate model — ratcheting floor, never a fixed 100% lie:** a per-package floor
  recorded in `scripts/gates/coverage_floor.tsv` (package → min %). The gate
  (`scripts/gates/coverage_gate.sh`) FAILs if any package drops below its recorded
  floor (§11.4.50 baseline ratchets upward; §11.4.110(3) coverage-completeness). New
  work raises the floor; nothing lowers it without an evidence-backed justification
  row. R8's "~100%" is the *destination*, the floor is the *mechanism*.
- **Honest exclusions (documented, not silent):** `cmd/*/main.go` process wiring and
  any generated code (future OpenAPI client) are excluded with a cited reason;
  exclusions live in the same TSV so they are auditable.
- **The floor is NOT the correctness bar.** Line coverage can be 100% and still
  bluff; the *real* bar is **behavioral coverage + a passing paired mutation per
  invariant** (§1.3). Coverage % is a necessary floor, the mutation is the proof.

---

## 2. The 13 R8 test types × THIS backend

The operator's R8/§11.4.27 list is used verbatim: **unit · integration · e2e ·
full-automation · security · ddos · scaling · chaos · stress · performance ·
benchmarking · ui · ux**. The plan's §11.4.169 list adds *contract, property, fuzz,
mutation, regression, smoke, acceptance* — those are **not dropped**: they are
cross-cutting families folded into the types below (contract → integration; mutation
→ §1.3 harness applied to every type; fuzz/property → unit+security; smoke/acceptance
→ e2e/full-automation; regression → every fixed gap gets a permanent guard §11.4.135).

**Real infrastructure rule (§11.4.27):** everything beyond *unit* runs against the
real Postgres+pgvector from `deploy/docker-compose.yml` — no DB mock. LLM-dependent
paths use the P3 `ModelProvider` (real HelixLLM/LLMsVerifier where reachable, else an
honest SKIP-with-reason — never a fabricated jury verdict).

| # | Type | What it means for this Go HTTP+Postgres+pgvector service | Real infra needed | Status |
|---|------|----------------------------------------------------------|-------------------|--------|
| 1 | **unit** | Pure logic: `CORS`/`redactQuery`/`ResolveAPIKeyAuth` (api), tree assembly + cycle guards (skill), `interpolate`/config-enum validation (config), `detectLanguage`/regex extraction (codeanalysis), `extractCodeBlocks`/`normalizeLanguage` (validation), TOON marshal round-trip, provider-factory dim check (db). Fake `ModelProvider` allowed here only. | none | **APPLICABLE-NOW** (5 pkgs done; extend to models/registry/mcp/autoexpand/db/worker) |
| 2 | **integration** | Real pgx against pgvector: migrate up → insert skill+typed edges → read back identical; recursive-CTE closure over the seed android graph; HNSW vector search returns ranked distances; TOON/JSON content-negotiation over a live Gin router. **Contract tests live here:** every response schema-validates against `api/openapi.yaml` and the live route table == the spec (G09). | `deploy/` compose (pgvector:pg16) | **APPLICABLE-NOW** |
| 3 | **e2e** | Server boots (wired `internal/api.Server` post-G01) → HTTP client drives create → (jury) → store → `get`/`tree`/`search`; MCP stdio handshake → `create_skill` → immediately `get_skill` in-session (R14). Full jury e2e needs live LLM. | compose + built binary; jury leg needs P3 LLM | **APPLICABLE-NOW (backend leg); jury/LLM leg PARTIAL** |
| 4 | **full-automation** | The whole thing re-runnable end-to-end with zero human action, `-count=3` stable, self-cleaning state (§11.4.98): compose up → migrate → seed → drive wizard tech-set → assert connected sub-DAG + closures → compose down. | compose + built binary | **APPLICABLE-NOW (backend); LLM-gated legs SKIP-with-reason** |
| 5 | **security** | CORS non-reflection + no wildcard+credentials (done); auth 401/503 fail-closed (done); SSRF guard on resource-verify (G21 → block `169.254.169.254`); sandbox escape/egress/fs-write blocked or SKIP not execute (G02); api_key never logged (done); `/metrics` gating (G24); table-name injection rejected (G27); no secret in tracked files (grep). | compose (for wired-server tests); network-egress sandbox for G02 | **APPLICABLE-NOW** |
| 6 | **ddos** | Flood the HTTP surface (connection/request storms) and assert graceful refusal (429) not collapse; body-size cap (413). **Mitigation is unbuilt (G22/P7.T5)** → the guard test is RED-first now (proves the surface is unprotected), flips GREEN when token-bucket rate-limiting lands. | compiled server | **PARTIAL — surface testable now (RED-first); GREEN blocked on G22** |
| 7 | **scaling** | *Dataset scaling* (applicable now): load 10k+ skills/edges, assert recursive-CTE + HNSW latency stays bounded and connection pool (`max_connections=25`) does not exhaust. *Horizontal cluster scaling*: N/A to a `systemctl --user` single-binary deploy. | compose + large seed | **PARTIAL — dataset-scaling APPLICABLE-NOW; horizontal-cluster SKIP (no cluster target, R15 is single-node)** |
| 8 | **chaos** | Worker panic-safety (G11: malformed coverage map ⇒ logged error not process death; per-goroutine `recover()`); mid-write `SIGKILL` on a pipeline commit ⇒ consistent recovery; DB connection drop mid-query ⇒ typed error not hang; migrations-missing boot ⇒ fail-fast (G23). Cleanup in `trap`/`t.Cleanup` (§11.4.14). | compose (killable container) | **APPLICABLE-NOW** |
| 9 | **stress** | Sustained N≥100 concurrent create/search/tree over the real DB; no deadlock, no leaked pgx conns, latency p50/p95/p99 recorded; `go test -race` clean under load. | compose | **APPLICABLE-NOW** |
| 10 | **performance** | Endpoint latency budgets (list/get/search/tree); HNSW query time vs corpus size; embedding call latency; recursive-CTE depth cost. Captured as `latency.json`. | compose | **APPLICABLE-NOW** |
| 11 | **benchmarking** | `go test -bench` micro-benchmarks: TOON-vs-JSON encode size + token-count (the R-mandate justification, G08), TOON marshal/unmarshal throughput, embedding/vector-insert, CTE closure. `-benchmem` for allocs. | none for codec bench; compose for DB bench | **APPLICABLE-NOW** |
| 12 | **ui** | Rendered-UI proof (§11.4.170 host-rendered pixels) for R3 clients. **No GUI client exists** (Web/Desktop/Mobile = P8 greenfield). The one sliver available now: the Bubble Tea **TUI** (`cmd/tui`) can get `teatest`-style golden-frame snapshot tests. | none (TUI); browsers/Flutter later | **SKIP-WITH-REASON (§11.4.3): R3 GUI clients not built (P8); host-rendered pixel proof binds when they ship. TUI snapshot = optional now.** |
| 13 | **ux** | Wizard-flow UX, accessibility (WCAG), OpenDesign token provenance (R12/§11.4.190) — all require the clients (P8) and OpenDesign vendoring (P10). | clients + OpenDesign | **SKIP-WITH-REASON (§11.4.3): depends on P8 clients + P10 OpenDesign; binds the moment a client ships.** |

**Count:** **9 APPLICABLE-NOW** (unit, integration, e2e, full-automation, security,
chaos, stress, performance, benchmarking) + **2 PARTIAL** (ddos — surface now,
mitigation RED-until-G22; scaling — dataset now, horizontal-cluster SKIP) + **2
SKIP-WITH-REASON** (ui, ux — no clients until P8/P10). No type is faked; every SKIP
carries its §11.4.3 reason and the exact unblock condition.

---

## 3. Challenges + HelixQA banks (R8)

Both are **real Helix components** (declared in `helix-deps.yaml`, currently
unvendored — blocked by the G14 submodule escalation; until vendored, the banks are
authored in-repo against the local runner interface and wired the moment the
submodules resolve — never faked as "passing" while absent).

### 3.1 Challenge bank — `project/test/challenges/`

A Challenge is an **adversarial end-to-end probe** that scores PASS only on positive
captured evidence (R11/§11.4.5) — it must try to *break* the claim, not confirm it.
Structure: one YAML/Go entry per invariant, referencing the gap it proves.

Example single Challenge entry (unauthorised-access probe, proves G01):
```yaml
# project/test/challenges/CH-G01-unauth-access.yaml
id: CH-G01-UNAUTH-ACCESS
proves: G01
setup: compose up + wired server with HELIX_API_KEYS set, HELIX_AUTH_DISABLED=false
steps:
  - name: no-key request is rejected
    request:  { method: GET, path: /api/v1/skills }        # no X-API-Key
    expect:   { status: 401 }
  - name: empty-keyset server refuses to serve open
    setup:    HELIX_API_KEYS="" HELIX_AUTH_DISABLED=false
    request:  { method: GET, path: /api/v1/skills }
    expect:   { status: 503 }                               # fail-closed, not open
  - name: disallowed Origin never reflected
    request:  { method: GET, path: /api/v1/skills, headers: { Origin: https://evil.example } }
    expect:   { header_absent: Access-Control-Allow-Origin }
evidence: challenges/evidence/CH-G01/<run-id>/{responses.json,headers.json}
pass_rule: all steps satisfied AND evidence files non-empty
```
Fabricated-skill Challenges (CH-G03/CH-G05/CH-G20) submit a **knowingly false skill**
and assert the jury/pipeline **rejects** it — the anti-bluff core of the product.

### 3.2 HelixQA bank — `project/test/helixqa/`

HelixQA runs autonomous QA sessions over the running system with captured evidence.
A HelixQA entry = a scenario bank (real request → real DB/LLM/response → recorded
transcript) that the `helixqa` binary drives. Structure mirrors §11.4.116's
append-only event stream: each verdict event carries its evidence path.

Example single HelixQA entry:
```yaml
# project/test/helixqa/HQA-wizard-android.yaml
id: HQA-WIZARD-ANDROID
scenario: R6 wizard end-to-end
input:  { techs: [android, android_aosp, java, kotlin, cpp, cmake] }
drive:  POST wizard → poll job → GET tree
assert:
  - connected sub-DAG returned, acyclic
  - every input tech resolves to a real skill node with provenance + real sources
  - no placeholder/"auto-generated to fill a gap" content (G20 anti-bluff)
evidence: helixqa/evidence/HQA-WIZARD-ANDROID/<run-id>/{transcript.jsonl,dag.json}
verdict: PASS|FAIL|SKIP-with-reason   # SKIP only if LLM provider unreachable
```
**Honest boundary:** HelixQA scenarios that need a live LLM jury (P3) or the
device-lab (Aurora/Harmony builds, G15) record an honest `operator_attended` /
`SKIP-with-reason` verdict until that dependency is live — never a fabricated PASS.

Both banks are wired into the local runner (§5), not a disabled CI pipeline.

---

## 4. Per-gap → test-type coverage matrix (R17 core)

Every G01–G27 finding gets a row: the test types that PROVE its fix, the
Challenge/HelixQA entry, and the captured-evidence each PASS must cite
(§11.4.5/§11.4.69). This is the R17 deliverable skeleton — each row's fix is "done"
only when its runtime signature verifies on a clean deployment (§11.4.108) and its
paired mutation flips RED (§1.3).

| Gap | Sev | Test types (all applicable) | Challenge / HelixQA | Captured evidence each PASS cites |
|-----|-----|-----------------------------|---------------------|-----------------------------------|
| **G01** dual servers / unwired hardened API, fail-open auth+wildcard CORS | CRIT | unit(done: CORS+auth), integration(auth 401/503, route==spec), security(no-ACAO for disallowed origin, no wildcard+creds), e2e, contract, regression, **mutation** (reflect-origin / drop-auth → FAIL) | Challenge: **CH-G01** unauth-access probe; HelixQA: yes | wired-server `responses.json`+`headers.json`; `grep 'Allow-Origin", "*"' cmd/server`=0; `internal/api` imported; mutation RED/GREEN logs |
| **G02** sandbox = host RCE, false isolation | CRIT | security(egress blocked, fs-write-outside blocked, fork/resource-limit), integration(no-runtime ⇒ SKIP not execute), **fuzz**(malicious-snippet corpus), mutation(remove isolation flag → egress test FAILs), regression | Challenge: **escape-attempt bank**; HelixQA: yes | egress-attempt `denied.log`; container/namespace id; SKIP-with-reason record on no-runtime; fuzz corpus + crash-free run |
| **G03** jury+auto-growth dead code; worker stubs | CRIT | unit(stage state machine), integration(draft→jury→merge on real DB), e2e(validated only after ≥2 approvals), mutation(strip a stage / fabricated skill passes → FAIL), regression | Challenge: **fabricated-skill-must-fail**; HelixQA: autonomous session over pipeline | pipeline `stage_transitions.json`; juror votes log; "no `validated` without recorded verdict" gated assertion |
| **G04** zero tests (this plan closes it) | CRIT | ALL 13 are the deliverable; prioritize unit+integration+contract+security+mutation | Challenge: yes; HelixQA: yes | `go test ./... -coverprofile` output; coverage-floor TSV; paired-mutation RED/GREEN pairs |
| **G05** empty jury auto-approves | HIGH | unit(empty jury ⇒ blocked not pass), integration(2-of-3 real approvals), mutation(flip auto-pass back → FAIL) | Challenge: yes; HelixQA: yes | empty-jury `rejected.json`; 2-of-3 vote transcript |
| **G06** `GetDependencyTree` depth-1 truncation | HIGH | unit(recursive tree assembly), integration(seed android returns known N-level tree), **property**(node count == closure size), regression, mutation(revert to depth-1 → FAIL) | Challenge: yes | tree JSON with depth>1 over seed corpus; node-count equality assertion log |
| **G07** TOML/JSON dep+resource round-trip drops edges | HIGH | unit(convert preserves names), integration(create-with-deps ⇒ edges in DB), **property/round-trip**(export→import→export byte-stable), contract, mutation(drop `DependsOnName` → FAIL) | Challenge: yes; HelixQA: yes | byte-stable round-trip diff (empty); `skill_dependencies` row dump; MCP `ImportFromTOML` edge-fidelity read (closes the UNCONFIRMED) |
| **G08** no TOON codec; wire-format mandate unmet | HIGH | unit(struct→TOON→struct), contract(**golden TOON vectors byte-for-byte**), integration(`Accept: application/toon` ⇒ TOON; unknown ⇒ JSON fallback + right `Content-Type`), fuzz(malformed TOON rejected), benchmark(TOON vs JSON size/tokens), mutation(swap codec to JSON → golden FAILs) | Challenge: yes | golden TOON fixtures; content-negotiation `headers.json`; TOON-vs-JSON token-count table |
| **G09** OpenAPI↔impl drift | HIGH | **contract**(schema-validate every response + route-table == `openapi.yaml`), integration, regression, smoke, mutation(rename a route → parity FAILs) | Challenge: yes; HelixQA: yes | route-parity report (live routes vs spec); per-endpoint schema-validation log |
| **G10** embedding dim: no model↔column assertion; `vector(768)` hardcoded | HIGH | unit(provider factory incl. helixllm; OpenAI length-mismatch rejected), integration(startup FAILs on dim mismatch; correct dim inserts), contract(config schema), mutation(column→1536 keep config 768 → startup assertion FAILs) | Challenge: yes | startup dim-assertion log (`pg_attribute` query); mismatch-rejected error; provider-factory table |
| **G11** worker no-op + panic on unchecked assertion | HIGH | unit(cycle calls pipeline; coverage type-safety comma-ok), integration(worker creates real skill from seeded gap), **chaos**(malformed coverage map ⇒ logged error not panic), mutation(reintroduce bare assertion → panic test FAILs) | Challenge: yes; HelixQA: yes | worker log showing recover()+continue; created-skill row; chaos malformed-input `no-panic.log` |
| **G12** tree-sitter stub; regex-only; Kotlin/C# unsupported | HIGH | unit(per-language extraction incl. kotlin), integration(parse real Android/Kotlin repo → real symbols), fuzz(malformed source no crash), mutation(remove a grammar → extraction FAILs) | Challenge: yes | extracted-symbols JSON from a real repo (no lorem); per-language coverage incl. kotlin |
| **G13** two rival `docker-compose.yml` | HIGH | **smoke**(compose up → `pg_isready`, `CREATE EXTENSION vector`, `SELECT 1`), integration(scripts reference the one file), regression(grep gate: no second compose), acceptance(`systemctl --user` up/down) | Challenge: yes | `pg_isready` output; `grep -rl docker-compose scripts/`=1 path; systemd up/down transcript |
| **G14** submodule policy conflict; 7 deps unvendored | HIGH | integration(`sync_submodules.sh --dry-run` resolves each dep to one canonical + pinned mirror), security(fail-closed on unexpected path — hardened), regression(`ls-remote` per dep), mutation(introduce nested own-org submodule → sync FAILs) | Challenge: yes | dry-run resolution report; `ls-remote` reachability per dep; attack-path fail-closed proof |
| **G15** Aurora/HarmonyOS client feasibility | HIGH | acceptance(one build artifact per OS OR documented blocker), smoke(thin client hits `/health`+wizard), contract(generated client compiles vs OpenAPI) | Challenge: per-OS build feasibility; HelixQA: **device-lab dependent — operator-attended where autonomous build infeasible** | build artifact OR §11.4.112 evidence-backed blocker; contract-compile log |
| **G16** "wasm" never uses WASM; Docker dual `/tmp` mount broken | MED | integration(Docker path runs the mounted file), unit(mount args well-formed), mutation(reintroduce dual /tmp → FAILs) | (folds into G02 bank) | Docker-run stdout of mounted file; mount-arg assertion |
| **G17** weak default DB password; config enums unvalidated | MED | unit(invalid provider/sandbox/level rejected; empty password rejected in prod mode), security(no secret in tracked files — pre-commit grep), mutation(add invalid enum → validate FAILs) | Challenge: yes | validate() rejection logs; tracked-file secret-grep=0 |
| **G18** CORS allowlist unreachable on live path; SPEC omits `allowed_origins` | MED | integration(config allowlist honoured end-to-end), security(non-allowlisted origin blocked), contract(config documents the key) | (folds into G01 bank) | live-server allowlist honoured `headers.json`; config-doc presence |
| **G19** SPEC §8 sample uses `--` (invalid TOML) | MED | unit/lint(parse every ```toml block in docs), regression | — | docs-lint output (all `toml` fences parse) |
| **G20** auto-expand fabricates placeholder skills; concrete-`*OpenAILLM` coupling; resources not persisted | MED | unit(nil LLM ⇒ no placeholder persisted; interface pluggability), integration(draft ⇒ resources persisted same tx), mutation(reintroduce placeholder-persist → anti-bluff test FAILs) | Challenge: yes | "no placeholder row" assertion; persisted-resources row dump |
| **G21** resource verify shallow (HEAD-only, fail-open, SSRF) | MED | unit(dead URL fails, mismatched hash fails, fetch error fails-closed), **security**(SSRF to 169.254.169.254 blocked), integration, mutation(flip fail-open back → FAILs) | Challenge: yes | SSRF-blocked `denied.log`; hash-mismatch reject; fail-closed error |
| **G22** no rate-limit/auth on live server; body limit only; Brotli flush errors ignored | MED | **load/ddos**(429 over-limit), integration(413 over-size), unit(Brotli error handled), security, regression | (folds into CH-G01/ddos bank) | 429 response under flood; 413 over-size; Brotli-error-abort log |
| **G23** migrations cwd-relative; failure only warns | MED | integration(missing migrations dir ⇒ startup FAILs), smoke(`migrate up` on fresh pgvector, `\d+` verified), regression | Challenge: yes | fail-fast exit≠0 log; `\d+` schema dump |
| **G24** health/metrics/version unauth; `/metrics` leaks internals | LOW | security(anonymous `/metrics` denied where required), contract, regression | — | `/metrics` 401 (gated) response; contract alignment |
| **G25** `RemoveDependency` ignores name-lookup errors → empty audit names | LOW | unit(audit detail records missing name), regression | — | audit row showing explicit not-found condition |
| **G26** `${VAR:-default}` can't resolve intentional-empty | LOW | unit(empty-override honoured via LookupEnv; unset uses default), regression | — | config-resolution table (empty vs unset) |
| **G27** `sanitizeTableName` strips instead of rejecting | LOW | unit(invalid table name rejected), security, regression | — | rejection error for `"skills; DROP"` |

**Every G01–G27 finding has a coverage-matrix row (27/27).** Rows fold sibling gaps
into shared Challenge banks where the register itself pairs them (G16→G02, G18→G01,
G22→G01/ddos) — the fold is explicit, never a silent drop.

---

## 5. CI/local gate wiring (§11.4.156 — CI/CD is DISABLED)

**Reconciliation:** §11.4.156 disables all GitHub Actions / GitLab pipelines in this
ecosystem, so gates run as **local shell scripts** invoked by the Makefile and the
Constitution's pre-build/post-build/runtime seams — **no `.github/workflows`, no
`.gitlab-ci.yml`**. The Makefile is the developer entry point; the gate scripts are
the mechanical enforcers (each with a paired mutation per §1.3).

**Pre-build gate — `scripts/gates/pre_build.sh`** (grep-speed, always-on):
- `go vet ./...` = 0; `gofmt -l` empty.
- `go test -short -race ./...` = 0 (unit layer).
- Security greps: `grep 'Allow-Origin", "*"' cmd/server` = 0 (G01); no secret in
  tracked files (G17); no `push --force` in scripts (§11.4.113).
- Contract route-parity check (live route table vs `api/openapi.yaml`) once G09 lands.
- Coverage floor: `scripts/gates/coverage_gate.sh` (§1.4).
- **Paired-mutation sweep:** `scripts/mutation/run_all.sh` runs every registered
  mutation, asserts each flips its guard RED then restores GREEN; residue scan
  (§11.4.84) confirms a clean tree afterward.

**Post-build gate — `scripts/gates/post_build.sh`** (artifact layer, §11.4.108(2)):
- Binaries built (`make build`); `--version` runs.
- Compose up (`deploy/`), `pg_isready`, `CREATE EXTENSION vector`, `SELECT 1`
  (G13 smoke); integration + e2e tags run against the live container
  (`make test-integration`, new `-run E2E`).
- `package.sh` produces zip + tar.gz; extract-and-build smoke (P13.T3).

**Runtime gate — `scripts/gates/runtime_signature.sh`** (§11.4.108(3), on a clean
deploy): asserts the *wired* runtime, not source — e.g. `GET /api/v1/skills` without
a key → 401 on the deployed binary; `internal/api.Server` is the live surface;
no-placeholder-content in any created skill (G20). This is the definition-of-done
oracle for the security/pipeline gaps.

**Challenges + HelixQA runner — `scripts/gates/run_qa.sh`:** drives
`test/challenges/*` and `test/helixqa/*` against the live compose stack, writing
captured evidence to `test/{challenges,helixqa}/evidence/<run-id>/`. Blocked banks
(unvendored HelixQA/Challenges submodules, G14) SKIP-with-reason honestly and the
gate reports the skip — never a fabricated pass.

**Makefile additions (thin wrappers, no new tool):** `make gate-pre`, `make
gate-post`, `make gate-runtime`, `make mutation`, `make qa` → the scripts above; the
existing `test`/`coverage`/`vet`/`lint` stay as-is.

---

## 6. Honest boundaries (§11.4.6) — what CANNOT be built yet, stated explicitly

These are flagged as tracked gaps, not silently omitted:

- **LLM-dependent proof (jury e2e, wizard full-run, auto-expand real content):**
  needs P3 `ModelProvider` wired to a reachable HelixLLM/LLMsVerifier/alias. Until
  then those e2e/HelixQA legs record an honest **SKIP-with-reason** (`llm_unreachable`)
  — a real jury verdict cannot be fabricated (G05/G20 forbid it). Unblock: P3 lands +
  a reachable provider (operator-scheduled live window for genuine per-provider proof).
- **ui / ux test types:** no GUI clients exist (Web/Desktop/Mobile = P8 greenfield;
  OpenDesign = P10). Host-rendered pixel proof (§11.4.170) and UX/accessibility bind
  the moment a client ships; today only the TUI admits a snapshot sliver. This is a
  **type-level SKIP-with-reason**, not zero coverage of shipped surfaces.
- **Aurora/HarmonyOS acceptance (G15):** device-lab / omprussia-embedder dependent;
  where autonomous build is infeasible the HelixQA verdict is **operator-attended**,
  and a failed feasibility spike is a §11.4.112 evidence-backed blocker, never a
  bluffed build.
- **Horizontal scaling:** the R15 deploy model is single-node `systemctl --user`;
  cluster-scaling has no target to test. Dataset-scaling (10k+ skills) IS testable
  now; cluster-scaling is out-of-scope-until-a-cluster-exists (not "passing").
- **Challenges/HelixQA banks are authored but not yet executed against the real
  submodules:** blocked by the G14 vendoring escalation. The banks + runner glue are
  built now (so nothing is un-wired research per §11.4.197); the *pass* is recorded
  only after the real `helixqa`/`challenges` binaries run — until then the runner
  SKIPs-with-reason (`submodule_absent`).
- **`ddos` mitigation (G22) and `security` for the sandbox (G02) depend on unbuilt
  code:** their guard tests are authored RED-first now (proving the surface is
  currently unprotected — real captured evidence of the defect, §11.4.115), flipping
  GREEN when the mitigation lands. RED-first is coverage, not a gap.

*Positive-evidence-only. Every APPLICABLE-NOW claim above is backed by a real file or
the captured `go test ./...` = exit 0 baseline; every deferral carries its §11.4.3/§11.4.6
reason and unblock condition. No test type is faked; no G-finding is left without a row.*
