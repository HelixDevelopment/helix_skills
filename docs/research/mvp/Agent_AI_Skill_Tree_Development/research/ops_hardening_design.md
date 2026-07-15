# Ops-Hardening Design — G13 / G17 / G22 / G23 / G24 (server/ops bundle)

**Revision:** 1
**Last modified:** 2026-07-15T16:20:00Z
**Status:** design-research, no code landed
**Description:** Design-decision document for the OPS-HARDENING gap cluster of the
HelixKnowledge Skill Graph System Go backend — G13 (rival docker-compose), G17
(weak/committed default DB password + missing config-enum validation), G22
(no rate-limit/auth-body-cap on the live server + ignored Brotli flush errors),
G23 (cwd-relative migrations that only warn on failure), G24 (unauthenticated
health/metrics/version + public Prometheus exposition). Every root cause is proven
from the committed baseline ref `255061b` (NOT the working tree, which another stream
is mid-mutating — §11.4.84/§11.4.119). Design research per §11.4.8/§11.4.150/§11.4.197.
**Authority / mandates served:** R1 (security posture), R15 (systemctl --user single-node
deploy), R17 (exhaustive gaps remediation + total test coverage), §11.4.10 (credentials),
§11.4.108 (four-layer runtime-signature-as-done), §11.4.201 (guard asserts the REAL
condition), §11.4.122 (no silent component removal), §11.4.186 (cross-doc consistency).
**Baseline ref:** `255061b` (`docs(mvp/skill-graph): G06/G07 DAG design …`), with the
G01 single-hardened-listener fix (`1a1a3f3`) already an ancestor — verified
(`git merge-base --is-ancestor 1a1a3f3 255061b` → YES).
**Scope:** read-only design; the single new file is this document. No Go/config/compose
source is modified by this deliverable.

---

## 0. Composition with already-landed work (read first)

Two prior fixes are LIVE in `255061b` and constrain every decision below:

- **G01 — one hardened listener (commit `1a1a3f3`, ancestor of `255061b`).**
  `cmd/server/main.go` `buildRouter` (from **`main.go:140`**) is the single live router
  for every HTTP-serving mode. Its middleware chain today is, in order:
  `gin.Recovery()` (`main.go:146`), `apiLoggingMiddleware(logger)` (`main.go:147`),
  `api.CORS(cfg.Server.AllowedOrigins)` (`main.go:151`), then the fail-closed auth
  resolved ONCE at **`main.go:160`** (`authMW := api.ResolveAPIKeyAuth(...)`) and applied
  to the `/api/v1` group (`main.go:~185`) and to the mounted MCP routes via
  `mcpServer.RegisterHTTPRoutes(router, authMW)` (**`main.go:269`**). `setupAPI`
  (**`main.go:298`**) starts the ONE `http.Server` and treats a bind failure as
  `logger.Fatal`. **Consequence for this cluster:** every new middleware for G22 (rate
  limit, body cap, Brotli-error handling) and every new route for G24 (`/metrics`,
  `/version`) MUST be added to `buildRouter` — the live surface — NOT to the still-dead
  `internal/api/server.go`. Adding hardening to `internal/api/server.go` would repeat the
  exact "fixed but not live" failure G01 closed (§11.4.108 SOURCE≠RUNTIME).
- **G18 / config-driven CORS (part of G01).** `internal/api.CORS` (`middleware.go:372`)
  is fail-closed on an empty allowlist and is already fed from `cfg.Server.AllowedOrigins`
  (`config.go:59-63`) on the live path (`main.go:151`). The `config.go` `substituteEnv`
  already interpolates `${VAR}` in `APIKeys`, `AllowedOrigins`, and `Database.Password`,
  and `validate()` (`config.go:428`) already carries a §11.4.10 defense-in-depth check that
  rejects an api-key/origin entry still containing `${` after interpolation (`config.go:443-452`).
  **Consequence for G17:** the config-enum + password fail-closed logic composes INTO the
  existing `validate()` function — it extends, never rewrites, the landed §11.4.10 checks.

---

## G13 — Two rival `docker-compose.yml` files

### Root cause (proven from `255061b`)
Two compose files are tracked, with divergent scope and divergent script references:

- **`project/docker-compose.yml` — 9198 bytes** (`git show 255061b:…/project/docker-compose.yml`).
  Declares `db` (`pgvector/pgvector:pg16`, container `skill-db`, user default
  `${DB_USER:-skilluser}`, password default `${DB_PASSWORD:-skillpassword}`),
  a full `api` service (`build: context: . dockerfile: Dockerfile`), a `worker` service,
  and `prometheus`+`grafana` under a `monitoring` profile. Subnet `172.20.0.0/16`.
- **`project/deploy/docker-compose.yml` — 3332 bytes**. Declares ONLY `postgres`
  (`pgvector/pgvector:pg16`, container `${COMPOSE_PROJECT_NAME:-helix-skills}-postgres`);
  the `app` service is intentionally **commented out** with an explicit anti-bluff note
  ("shipping a fake/placeholder 'app' image here would be a bluff"). Carries `name:`,
  `.env.example`, and `systemd/helix-skills.service` alongside it.

The scripts do NOT agree on which file is canonical:
- `scripts/_lib.sh:57` → `HX_COMPOSE_FILE="${HX_DEPLOY_DIR}/docker-compose.yml"` (deploy).
- `deploy/systemd/helix-skills.service:12` → references `deploy/docker-compose.yml`.
- `start.sh`/`stop.sh`/`restart.sh`/`status.sh`/`logs.sh` → all reference `deploy/docker-compose.yml`.
- **BUT** `scripts/package.sh:75` → `cp "$PROJECT_DIR/docker-compose.yml" "$package_dir/"`
  (the ROOT file), and `scripts/backup.sh:136-137` → `cp "$INSTALL_DIR/docker-compose.yml" …`
  (a root-level file in the install dir).

The two files also disagree on the DB user default (`skilluser` vs the config's `skill`),
container names, subnet, and the presence of the app/worker/monitoring services — the exact
divergence-drift the §11.4.186 anti-divergence mandate forbids. `IMPLEMENTATION_PLAN.md:258`
(P12.T4) requires exactly one canonical compose.

### DECISION
Make **`project/deploy/`** the single canonical ops home (it already carries
`.env.example` + `systemd/` + the hardened `_lib.sh` path resolution). Then:
1. Fold the app/worker/monitoring service DEFINITIONS from the root file into
   `deploy/docker-compose.yml` as **compose profiles** (`profiles: [app]`,
   `profiles: [monitoring]`) — preserving them, NOT deleting them (§11.4.122: those are
   already-declared end-user-visible services; a silent drop is a release blocker). The
   `postgres`-only default stays the standing datastore layer; `--profile app` /
   `--profile monitoring` opt in.
2. Point `package.sh:75` and `backup.sh:136-137` at the canonical `deploy/docker-compose.yml`.
3. Remove `project/docker-compose.yml` ONLY after (1)+(2) land, as its own descriptive
   commit citing this design (§11.4.124 dead-code investigate-before-remove: the root file
   is the SOURCE of the app/worker/monitoring definitions, so it is folded-forward, never
   dropped-on-sight).
- **Alternatives rejected:** keeping both "for dev vs deploy" — the rival-copy anti-pattern
  itself; use profiles/overrides for the two modes instead.

### WHY (external precedent, §11.4.8)
Compose supports first-class `profiles:` for optional service sets and a documented single
base file + override model (`compose.override.yml`) — the canonical way to express
"one file, several activation modes" without divergent copies (Docker Compose "Using profiles"
and "Multiple Compose files" docs). NO exotic mechanism is needed; this is the standard
one-canonical-file pattern.

### RUNTIME SIGNATURE (§11.4.108)
On a clean checkout: `git ls-files 'project/**/docker-compose*.yml'` returns **exactly one**
path (`project/deploy/docker-compose.yml`), AND `grep -rl 'docker-compose\.yml' project/scripts
project/deploy/systemd` resolves every reference to that one path (zero references to a
root-level `project/docker-compose.yml`), AND `compose -f deploy/docker-compose.yml up -d`
→ `pg_isready` OK + `CREATE EXTENSION IF NOT EXISTS vector` + `SELECT 1` all succeed.

---

## G17 — Weak/committed default DB password; config enums unvalidated

### Root cause (proven from `255061b`)
- **`internal/config/config.go:177`** — `defaultConfig()` sets
  `Password:       "secret",` inside `DatabaseConfig`. A binary launched with no override
  boots against a hardcoded, committed credential.
- **`config/config.toml:35`** — the shipped template contains `password = "secret"`
  (a real, tracked secret literal, not a `${VAR}` placeholder). The compose defaults
  compound this with `${DB_PASSWORD:-skillpassword}` (root compose) / `${DB_PASSWORD:-skillpassword}`
  (`deploy/docker-compose.yml`) and `DB_PASSWORD=skillpassword` in `deploy/.env.example`.
- **`internal/config/config.go:428-486`** — `validate()` checks ports, the `${` fail-closed
  placeholder residue (`config.go:443-452`), `Embedding.Dimensions>0`, `Validation.JurySize>0`,
  `Validation.ApprovalThreshold>0`, `AutoExpand.MaxDepth>0`, `AutoExpand.MaxNewSkillsPerRun>0`,
  `Registry.CoverageThreshold∈[0,1]` — but it **never** validates the closed-set string
  fields against their allowed sets: `Embedding.Provider` (`config.go` doc: `"openai" | "local"`),
  `Validation.SandboxType` (`"wasm" | "docker" | "none"`), `Logging.Level`
  (`"debug"|"info"|"warn"|"error"`), `MCP.Transport` (`"stdio" | "http"`; `main.go` also uses
  `"both"`). A typo (`provider = "opennai"`) passes `validate()` and fails late, deep in the
  provider-factory call stack.

### DECISION
Extend the existing `validate()` (compose into it, never rewrite it):

1. **Fail-closed on a default/committed DB password.** Remove `Password: "secret"` from
   `defaultConfig()` (`config.go:177`) so there is NO working default. In `validate()`, reject
   `Database.Password` when it is empty OR equals a KNOWN committed sentinel from the tracked
   tree (`"secret"`, `"skillpassword"`) UNLESS an explicit, deliberately-set dev opt-in signal
   is present (see honest gap G17-a). The rejection message names ONLY the field
   (`database.password`) and the failure class — never the value (§11.4.10 field-names-only;
   §11.4.201 the guard asserts the REAL value, not a proxy). Rotate `config.toml:35` to a
   commented `# password = "${HELIX_DB_PASSWORD}"` placeholder so no live secret is tracked.
2. **Closed-set enum validation.** In `validate()`, reject `Embedding.Provider ∉
   {openai, local}`, `Validation.SandboxType ∉ {wasm, docker, none}`, `Logging.Level ∉
   {debug, info, warn, error}`, `MCP.Transport ∉ {stdio, http, both}`. Each rejection cites the
   field name and the allowed set (actionable, §11.4.6). The allowed sets are DATA read from the
   field doc-comments already in `config.go`, not invented.
- **Alternatives rejected:** a documented default password — a standing §11.4.10 credential-hygiene
  risk; late-failing enum typos — confusing runtime errors.

### WHY (external precedent, §11.4.8)
Fail-closed-on-default-credential is standard secure-config hygiene (OWASP ASVS V2/V6:
"no default/shipped credentials"; 12-factor "config in the environment"). Enum-validate-at-load
is the standard "validate the whole config once at startup, fail fast" pattern (the same
discipline `validate()` already applies to numeric fields). NO external solution needed beyond
these established patterns — the design is the direct application of the field-name-only §11.4.10
rule to the password and the existing `validate()` shape to the enums.

### RUNTIME SIGNATURE (§11.4.108)
On a clean deployment: a config with `password=""` OR `password="secret"` (dev opt-in absent)
causes `config.Load` to return a non-nil error and the process exits non-zero BEFORE binding
the port (observable: no listener on `cfg.Server.HTTPPort`; the error names `database.password`
with no value echoed). A config with `provider="opennai"` (or any out-of-set enum) causes the
same load-time non-zero exit naming `embedding.provider` and its allowed set. A valid config
+ real env-supplied password boots and serves. `git grep -n '"secret"\|skillpassword'
project/config project/internal/config` returns zero tracked live-secret literals.

---

## G22 — No rate-limit / auth-body-cap on the live server; Brotli flush errors ignored

### Root cause (proven from `255061b`)
- **Live path has neither rate-limiting nor a body cap.** `buildRouter` (`main.go:140-296`)
  installs `gin.Recovery`, `apiLoggingMiddleware`, `api.CORS`, and the fail-closed auth
  (`main.go:146-160`) — and NOTHING else. There is no token-bucket limiter and no
  `MaxBodySize`. `go.mod` declares no rate-limit dependency (no `golang.org/x/time`, no
  gin rate-limit middleware).
- **The 100 MB body cap exists only in the DEAD path.** `internal/api/middleware.go:480`
  `func MaxBodySize(maxBytes int64)` wraps `http.MaxBytesReader`; it is applied ONLY in the
  unwired `internal/api/server.go:153` (`s.router.Use(MaxBodySize(100 * 1024 * 1024))`), which
  `cmd/server` never constructs. On the live router the request body is unbounded.
- **Brotli flush/close errors are discarded.** `internal/api/middleware.go:104-108` —
  the deferred closure calls `bw.writer.Flush()` and `bw.writer.Close()` and discards both
  return values, so a compression error yields a silently truncated response with a 200 status.

Post-G03 wiring will make `/expand`, `/learn`, and validation endpoints code-executing; an
unauthenticated-flood + unbounded-body + code-executing surface is a DoS/abuse vector, and the
silent Brotli truncation corrupts responses.

### DECISION
Add to the SINGLE live `buildRouter` chain (composing with the G01 order — rate-limit BEFORE
the auth work so throttling precedes credential checks; body-cap early):

1. **Token-bucket rate limiting** via `golang.org/x/time/rate`, keyed per client (API key when
   present, else remote IP), returning HTTP 429 when the bucket is empty. A bounded, TTL-reaped
   map of `*rate.Limiter` per key (memory-bounded; stale keys evicted) so one client's flood
   cannot starve another (per-key isolation) and the map cannot grow unbounded.
2. **Body cap on the live path** — apply the existing `internal/api.MaxBodySize(100MB)`
   (already implemented, `middleware.go:480`) to `buildRouter` so oversized bodies return 413.
   This wires an existing, tested helper to the live surface (§11.4.108 SOURCE→RUNTIME) rather
   than duplicating it.
3. **Brotli error handling** — capture the `Flush()`/`Close()` return values at
   `middleware.go:106-107`; on a non-nil error, abort the response (do not emit a 200 over a
   truncated body). A compression failure becomes a visible error, never a silent corruption.
- **Alternatives rejected:** relying on an upstream proxy for limits — R15's `systemctl --user`
  single-binary deploy has no mandated proxy; the app MUST be safe standalone.

### WHY (external precedent, §11.4.8)
`golang.org/x/time/rate` is the canonical Go token-bucket limiter (`NewLimiter(r, b)` +
`Allow()`), designed for exactly this HTTP-throttle use, goroutine-safe; the per-client-key
`map[string]*rate.Limiter` + 429 pattern is the documented middleware idiom (verified against
`pkg.go.dev/golang.org/x/time/rate`). `http.MaxBytesReader` (already used by the existing
`MaxBodySize`) is the stdlib body-cap. Checking `io`/compressor `Close()` errors is the standard
Go correctness rule (an ignored `Close()` on a flushing writer drops buffered/erroring bytes).

### RUNTIME SIGNATURE (§11.4.108)
On a clean deployment of the live binary: a burst above the configured rate against
`/api/v1/skills` yields HTTP **429** (captured status), a request body above the cap yields HTTP
**413**, two distinct client keys each keep their own budget (one flooding does not 429 the
other), and a forced Brotli error produces a non-2xx response (never a 200 over a truncated
body). Baseline captures 200/no-limit for all of these — the RED-first evidence (§11.4.115) that
the surface is currently unprotected.

---

## G23 — Migrations loaded from a cwd-relative path; failure only warns and the server continues

### Root cause (proven from `255061b`)
- **`cmd/server/main.go:86`** — `if err := db.Migrate(ctx, pool, "./migrations"); err != nil {`
  passes a cwd-RELATIVE path. `db.Migrate` (`internal/db/migrations.go`) resolves it via
  `os.ReadDir(dir)` inside `discoverMigrations`, so a server started from any directory other
  than the one containing `./migrations` finds an empty set and applies nothing.
- **`cmd/server/main.go:87-88`** — on migration failure the code logs `logger.Warn("Migration
  failed", …)` and **falls through to continue booting**. The server then binds the port and
  serves traffic against an un-migrated (schema-less) DB, failing every query at runtime — the
  §11.4.108 "runs green but broken" hazard. The migrations tree on disk is
  `project/migrations/001_initial.{up,down}.sql`; `git grep '"embed"\|embed.FS' project` returns
  zero — nothing is embedded today.

### DECISION
1. **Deterministic path resolution via `embed.FS`.** Embed the migration SQL into the binary
   (`//go:embed migrations/*.sql`) and pass the `embed.FS` (or an `fs.FS`) to `db.Migrate` so
   discovery is cwd-independent and the migrations always travel with the binary. This requires
   `Migrate`/`discoverMigrations` to accept an `fs.FS` (both `embed.FS` and `os.DirFS` satisfy
   `fs.ReadDirFS`), replacing the `os.ReadDir`/`os.ReadFile` calls with `fs.ReadDir`/`fs.ReadFile`.
2. **Fail-CLOSED on migration error.** Replace the `main.go:86-88` warn-and-continue with a
   fatal path: a migration error (or an unresolvable migration set) causes a non-zero exit
   BEFORE the port binds. A server that cannot bring its schema to the required version is NOT
   serving-ready and must not pretend to be (mirrors the G01 fatal-on-bind-failure discipline).
- **Alternatives rejected:** warn-and-continue — hides a fatal state; a cwd-relative external
  dir — non-deterministic across working directories and packaging layouts.

### WHY (external precedent, §11.4.8)
`embed.FS` is the stdlib mechanism for shipping SQL/assets inside a Go binary; it implements
`io/fs.FS`, so code that used `os.ReadDir`/`os.ReadFile` moves to `fs.ReadDir`/`fs.ReadFile` with
identical semantics (verified against `pkg.go.dev/embed`). Fail-fast-on-migration-error is the
standard migration-tool contract (golang-migrate, goose: a failed migration aborts startup, it
does not degrade to serving). The testing plan already anticipates this exact approach
(`testing_infrastructure_plan.md:77` names `testing/fstest` "for embed-FS migration tests (G23)").

### RUNTIME SIGNATURE (§11.4.108)
On a clean deployment: starting the binary from ANY working directory applies the embedded
migrations to the required version (observable: `\d+` on a fresh pgvector DB shows the expected
tables, and `SELECT MAX(version) FROM schema_migrations` equals the highest embedded version);
a forced migration failure (e.g. an intentionally-broken migration or an unreachable DB during
migrate) causes a non-zero process exit with NO listener bound on `cfg.Server.HTTPPort`. Baseline
captures the process continuing-to-serve after a warn — the RED-first evidence of the hazard.

---

## G24 — Health/metrics/version unauthenticated; `/metrics` exposes Prometheus internals publicly

### Root cause (proven from `255061b`) — with an honest correction to the register's line note
- **Live router today registers ONLY `/health` (open) among the System routes.** `buildRouter`
  registers `router.GET("/health", …)` at **`main.go:163`** (open, hand-rolled) and
  `router.GET("/", …)` at `main.go:271` (open server-info). It does **NOT** register `/metrics`
  or `/version` at all. (The register's `main.go:151` note for `/metrics`+`/version` predates the
  G01 rewrite; in `255061b` the live router exposes neither — stated as FACT, §11.4.6.)
- **The dead hardened server registers all three OUTSIDE auth.** `internal/api/server.go:158-161`
  (`SetupRoutes`): `s.router.GET("/health", …)`, `s.router.GET("/metrics", s.handleMetrics())`,
  `s.router.GET("/version", …)` are registered at the router root, BEFORE the `/api/v1` auth
  group (`server.go:163-172`) — so if this server were the wired one, all three would be
  unauthenticated.
- **`/metrics` serves the full Prometheus default registry.** `internal/api/system_handler.go:111-118`
  `handleMetrics` returns `promhttp.Handler()` (the complete default-registry exposition:
  `helix_api_http_requests_total`, latency histograms, `helix_api_goroutines_count`,
  `helix_api_memory_usage_bytes`, Go runtime metrics) to any caller.
- **The contract says all three are authenticated.** `api/openapi.yaml:41-42` sets a GLOBAL
  `security: [ApiKeyAuth: []]`; `/health` (`openapi.yaml:958`), `/metrics` (`openapi.yaml:983`),
  and `/version` (`openapi.yaml:1002`) EACH carry a `'401': Unauthorized` response — so the spec
  marks every System endpoint as key-guarded. This is a live cross-doc divergence (§11.4.186):
  the contract-401 on `/health` contradicts the operational need for an open liveness probe.

### DECISION
1. **Keep `/health` OPEN** (liveness) — an orchestrator/systemd probe must reach it without a
   key. Correct `api/openapi.yaml` to document `/health` as unauthenticated (remove its `401`;
   add `security: []` override on that operation) so the contract matches the operable posture
   (§11.4.186 divergence closed at the doc, not by breaking probes).
2. **Gate `/metrics` behind auth** on whichever router is live post-G09/route-parity. Because the
   live `buildRouter` exposes no `/metrics` today, `/metrics` MUST be registered UNDER `authMW`
   from the start (register it inside the auth-guarded group, never at the router root), OR bind
   the metrics endpoint to a private/loopback interface separate from the public listener
   (single-node R15). Either way, an anonymous `GET /metrics` returns 401 (or is unreachable on
   the public interface).
3. **Align `/version`** with the contract — register it under `authMW` (the spec marks it 401).
- **Alternatives rejected:** authing `/health` — breaks orchestrator/systemd liveness probes;
  leaving `/metrics` public — leaks internal counters, versions, goroutine/memory gauges to
  anonymous callers.

### WHY (external precedent, §11.4.8)
Public Prometheus `/metrics` is a recognized info-leak (Prometheus "Securing Prometheus" guidance:
the exposition endpoint should be access-controlled or bound to a private interface); keeping
liveness (`/health`) open while gating telemetry (`/metrics`) is the standard split (Kubernetes
liveness/readiness probes are unauthenticated by design, metrics scraping is network-restricted).
NO exotic mechanism needed.

### RUNTIME SIGNATURE (§11.4.108)
On a clean deployment of the live binary: anonymous `GET /health` → **200** (probe works);
anonymous `GET /metrics` → **401** (or connection refused on the public interface if bound
privately); anonymous `GET /version` → **401**; an authenticated `GET /metrics` (valid X-API-Key)
→ 200 with the Prometheus exposition; and the live System-route table matches
`api/openapi.yaml`'s (corrected) auth posture (contract parity).

---

## Cross-gap interaction summary (§11.4.92 blast-radius)

- **G22 + G24 share the `buildRouter` chain (`main.go:140-296`).** The rate-limiter,
  body-cap, and the newly-gated `/metrics`/`/version` all attach to the SAME live router the
  G01 fix produced. Ordering: `Recovery → logging → CORS(G01) → rate-limit(G22) →
  body-cap(G22) → auth(G01)`; `/metrics`+`/version`(G24) registered INSIDE the auth group;
  `/health` stays open. They do not conflict — they compose as additional links on one chain.
- **G17 extends the existing `validate()` (`config.go:428`)** already carrying the G01/§11.4.10
  `${` residue checks — the password + enum logic is appended, not a rewrite.
- **G23 is isolated to `main.go:84-90` + `internal/db/migrations.go`** (migration path + fatal
  policy) and touches no HTTP surface — no interaction with G22/G24.
- **G13 is ops-file/script only** (`deploy/` canonicalization); its only code-adjacent edge is
  that `deploy/docker-compose.yml` is where the G17-hardened DB password default surfaces
  (`${DB_PASSWORD:-skillpassword}`) — the canonical compose should carry the same
  no-committed-default posture in its `.env.example` guidance.

---

## Test-case enumeration (RED-first §11.4.115; each PASS cites captured evidence §11.4.5/§11.4.69)

Every test is table-driven stdlib `testing` (the codebase convention; `testing_infrastructure_plan.md:61-78`),
integration/e2e against real `pgvector` from `deploy/docker-compose.yml` (no DB mock, §11.4.27),
and each guarded invariant carries a paired §1.1 mutation registered in the plan's
`scripts/mutation/` harness (`plan §1.3`).

**G13 — 5 cases (1 mutation):**
1. regression/unit: `git ls-files` shows exactly one tracked compose under `project/`.
2. integration: every `scripts/` + `_lib.sh:57` + `package.sh:75` + `backup.sh:136` compose
   reference resolves to the one canonical `deploy/docker-compose.yml`.
3. smoke: canonical compose up → `pg_isready`, `CREATE EXTENSION vector`, `SELECT 1`.
4. acceptance: `systemctl --user` start/stop cycle drives the one file.
5. **mutation:** re-introduce a second root compose → the single-compose grep gate FAILs.

**G17 — 11 cases (2 mutations):**
1. unit: empty `database.password` rejected (dev opt-in absent).
2. unit: committed sentinel `"secret"` rejected (dev opt-in absent).
3. unit: committed sentinel `"skillpassword"` rejected (dev opt-in absent).
4. unit: explicit dev opt-in signal present → default password accepted (honest escape).
5. unit: invalid `embedding.provider` rejected; 6. unit: invalid `validation.sandbox_type`
   rejected; 7. unit: invalid `logging.level` rejected; 8. unit: invalid `mcp.transport` rejected.
9. unit: all-valid enums + real env password → `Load` succeeds.
10. security: `git grep` finds no live secret literal in tracked config (`config.toml`/`config.go`).
11. **mutation A:** re-add `Password:"secret"` default / drop the reject → the password test FAILs;
    **mutation B:** drop one enum check → that enum test FAILs. (2 mutations)

**G22 — 7 cases (2 mutations):**
1. unit: token-bucket `Allow()` permits under the configured rate.
2. ddos/load: burst above rate on the live path → 429 (RED-first now).
3. integration: body above the cap on the live path → 413 (RED-first now).
4. unit: forced Brotli `Flush`/`Close` error → response aborted, not a 200 over truncated bytes.
5. security: per-key isolation — client A's flood does not 429 client B.
6. **mutation A:** remove the rate-limit middleware → the 429 test FAILs;
   **mutation B:** revert Brotli to discard errors → the Brotli test FAILs. (2 mutations)

**G23 — 6 cases (1 mutation):**
1. integration: unresolvable/failed migration at startup → process exits non-zero, no listener.
2. integration: migrate applies on fresh pgvector; `\d+`/`schema_migrations` show expected version.
3. unit (`testing/fstest`): `embed.FS`-backed discovery is deterministic (no cwd dependence).
4. smoke: server started from a different cwd still migrates to the required version.
5. chaos: migration mid-apply failure → server exits, never serves an un-migrated schema.
6. **mutation:** revert to warn-and-continue → the startup-fails test FAILs (server serves on
   empty schema).

**G24 — 6 cases (1 mutation):**
1. security: anonymous `GET /metrics` → 401 (or unreachable on the public interface).
2. security: anonymous `GET /version` → 401 (contract-aligned).
3. integration: anonymous `GET /health` → 200 (liveness probe preserved).
4. contract: live System-route table == corrected `api/openapi.yaml` auth posture.
5. security: authenticated `GET /metrics` → 200 exposition (positive path).
6. **mutation:** move `/metrics` outside the auth guard → the anonymous-denied test FAILs.

**TOTAL: 35 test cases, including 7 paired §1.1 mutations** (G13:5/1, G17:11/2, G22:7/2,
G23:6/1, G24:6/1).

### Reconciliation with `research/testing_infrastructure_plan.md` (§11.4.186)
- **G13 → plan matrix row `testing_infrastructure_plan.md:299`** (smoke/integration/regression/
  acceptance). My 5 cases map 1:1; NO new test type — consistent, no extension.
- **G17 → plan row `:303`** (unit invalid provider/sandbox/level; empty-password prod-mode;
  security no-secret-grep; mutation invalid-enum). **EXTENSION flagged (§11.4.186):** row `:303`
  lists provider/sandbox/level — I ADD `mcp.transport` to the enum set, and I extend the password
  case from "empty password" to also reject the committed sentinels (`secret`/`skillpassword`)
  plus the explicit dev opt-in case. These extend row `:303`; the row's test types are unchanged.
- **G22 → plan row `:308`** (load/ddos 429; integration 413; unit Brotli; security; regression;
  folds into CH-G01/ddos bank) AND plan §2 table row #6 (`:195`, "ddos — RED-first now, GREEN
  blocked on G22"). My cases map directly; **EXTENSION flagged:** I add the per-key-isolation
  security case (row `:308` does not name it explicitly).
- **G23 → plan row `:309`** (integration missing-dir→startup-fails; smoke migrate-up+`\d+`;
  regression) AND plan `§1.1:77` (`testing/fstest` for embed-FS migration tests, G23). My cases
  map; the `embed.FS`/`fstest` and chaos-mid-apply cases reconcile to `:77` (already anticipated)
  — no new type beyond what the plan names.
- **G24 → plan row `:310`** (security anonymous `/metrics` denied; contract; regression).
  **EXTENSION flagged:** I add the `/health`-stays-open case, the `/version` alignment case, and
  the authenticated-`/metrics` positive case; and I flag the `api/openapi.yaml` `/health`-401
  correction as a cross-doc (§11.4.186) fix the row does not currently carry.

---

## Honest gaps (§11.4.6) — what needs live infra or an operator decision, stated precisely

- **G13-a (operator decision, §11.4.122):** the root `project/docker-compose.yml` is the ONLY
  declaration of the `api`, `worker`, `prometheus`, and `grafana` services. This design folds them
  forward into `deploy/` as profiles rather than deleting them; whether the `app`/`worker`
  services should be ACTIVATED now (the `deploy/` file deliberately comments the `app` service out
  as "no real image yet — a fake would be a bluff") is an operator/build-readiness decision, not
  an autonomous one. UNCONFIRMED until the Go binary builds a real image (the `deploy/` note's own
  gate). The design does NOT silently drop those service definitions.
- **G17-a (design addition to verify):** there is NO existing "environment"/dev-vs-prod field in
  `config.go` today (verified: `Config` has no such field). The fail-closed-on-default-password
  rule REQUIRES a new explicit opt-in signal (a config field or an env var such as
  `HELIX_ALLOW_DEFAULT_DB_PASSWORD`). The exact name/shape is a design choice to confirm with the
  operator; this document does not assume one already exists.
- **G17-b:** the committed-sentinel denylist (`secret`, `skillpassword`) is enumerated from the
  captured tracked tree (`config.go:177`, `config.toml:35`, both compose files, `.env.example`).
  Whether to hardcode that denylist versus a generic "must differ from any value tracked in the
  repo" check is a design choice; the denylist is the conservative, evidence-grounded default.
- **G22-a (needs live calibration, §11.4.6):** the concrete rate-limit numbers (`r` tokens/s,
  burst `b`) are UNCONFIRMED — they MUST be calibrated on the target R15 single-node deploy under
  a real load profile, NOT hardcoded from literature. The 429/413/isolation BEHAVIOR is provable
  now (RED-first); the specific thresholds require a live load test to fix.
- **G22-b (needs live container):** the 413 body-cap and 429 flood behaviors on the LIVE binary
  are provable only against the wired server + real pgvector container; the unit-level token-bucket
  and Brotli-error tests run without infra.
- **G23-a (needs live pgvector):** the fail-fast non-zero-exit-on-migration-failure and the
  `\d+`/`schema_migrations` version assertion require the `deploy/` container; the `embed.FS`
  discovery determinism is provable via `testing/fstest` without a DB.
- **G24-a (cross-doc contract tension, §11.4.186):** `api/openapi.yaml:41-42` applies a GLOBAL
  `ApiKeyAuth` and marks `/health` with a `401` (`:980`). Keeping `/health` open (the register's
  own decision and the operational requirement) CONTRADICTS that contract, so the spec MUST be
  corrected (document `/health` as `security: []`). This is a real divergence to close at the doc;
  the design states it rather than silently serving `/health` in violation of the tracked spec.
- **G24-b (topology assumption to verify):** the "bind `/metrics` to a private interface" option
  assumes the R15 single-node `systemctl --user` topology has a usable loopback/private bind
  distinct from the public listener. The auth-gated `/metrics` option needs no such assumption and
  is the portable default; the private-bind option is offered only where the deploy topology
  supports it.

*Positive-evidence-only. Every root cause above is quoted from `git show 255061b:…`; every
"live vs dead" claim is grounded in the exact file:line; every deferral carries its §11.4.6 reason
and unblock condition. No word from the §11.4.6 forbidden set is used to describe a cause.*
