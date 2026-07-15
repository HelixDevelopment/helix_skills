# G18 / G25 / G26 — CORS-Doc-Drift Reconciliation + Audit/Env-Substitution Correctness Bundle

**Revision:** 1
**Last modified:** 2026-07-15T16:45:00Z
**Status:** design-research (drives a later Go fix — no code landed by this doc)
**Authority:** remediates `GAPS_AND_RISKS_REGISTER.md` §MED G18 (lines 211-217) and
§LOW G25 (lines 271-277) + G26 (lines 279-285); reconciles with
`research/testing_infrastructure_plan.md` rows for G18 (line 304), G25 (line 311),
G26 (line 312), and the fold-note at line 316 (`G18→G01`).
**Anti-bluff note (§11.4.6/§11.4.123):** every root-cause and current-state claim
below cites a file:line read this session from the committed baseline ref
`255061b7c50b2d26c2fb1515fdff131a4c61520c` via `git show 255061b:<path>` — never
the uncommitted G02 working-tree change. Nothing is inferred from memory or
training data. Forbidden vocabulary (§11.4.6) is not used anywhere below; every
"is/does/confirms" statement is pinned to a quoted line.

---

## 0. Method

Read at `255061b` (which the repository's `git log --oneline` confirms is a
descendant of `1a1a3f3` — `fix(skill-graph/security): close runtime CORS+auth
hole — single hardened listener [G01]` — on the `main.go` file's own history:
`git log --oneline 255061b -- .../cmd/server/main.go` returns `1a1a3f3` as the
newest touch at-or-before `255061b`, i.e. `255061b`'s `main.go` **is** the
post-G01 file):

- `project/cmd/server/main.go` (415 lines, full read)
- `project/internal/config/config.go` (501 lines, full read)
- `project/internal/skill/graph.go` (451 lines, `RemoveDependency`/`AddDependency` read in full)
- `project/internal/api/middleware.go` (`ResolveAPIKeyAuth` + `CORS`, lines 330-419)
- `project/internal/api/middleware_test.go` (existing CORS unit-test roster)
- `project/cmd/server/security_test.go` (existing live-router security tests, full read)
- `project/internal/config/config_test.go` (existing `interpolate`/`substituteEnv` test roster)
- `project/internal/skill/graph_test.go` (existing dependency-graph test roster + its documented `_RequiresLiveDatabase` idiom)
- `project/config/config.toml` (the real shipped config template)
- `SPEC.md` §8 (the documented config sample)
- `research/testing_infrastructure_plan.md` (per-gap coverage-matrix rows)

---

## 1. G18 — CORS allowlist unreachable on the live path; SPEC config sample omits `allowed_origins`

### 1.1 Root cause (as originally filed)

Register evidence (`GAPS_AND_RISKS_REGISTER.md:214`):

> "the hardened, config-driven allowlist lives in `internal/api`
> (`middleware.go:328-387`, fed by `ServerConfig.AllowedOrigins`,
> `config.go:59-63`) but the live server uses `corsMiddleware()` wildcard
> (`cmd/server/main.go:362-373`) and never reads `AllowedOrigins`."

### 1.2 Current-state determination — FACT, verified by direct read

`255061b`'s `cmd/server/main.go` (post-`1a1a3f3`) `buildRouter` reads:

```go
// main.go:148-151
// Hardened, config-driven CORS allowlist (internal/api.CORS): an empty
// allowlist is fail-closed and no wildcard "*" origin is ever emitted with
// credentials. This replaces the previous wildcard corsMiddleware().
router.Use(api.CORS(cfg.Server.AllowedOrigins))
```

`grep -n corsMiddleware main.go` (this session) returns exactly one hit — the
comment on line 150 documenting the replacement; the function itself is gone
from the file. `cfg.Server.AllowedOrigins` is read directly at the call site.

`internal/config/config.go:59-63` still declares the field:

```go
// AllowedOrigins is the CORS allowlist of exact origins permitted to make
...
AllowedOrigins []string `toml:"allowed_origins"`
```

and it is now **wired through the whole config pipeline**, not merely declared:

- `substituteEnv` interpolates every entry (`config.go:309-311`):
  `for i := range cfg.Server.AllowedOrigins { cfg.Server.AllowedOrigins[i] = sub(cfg.Server.AllowedOrigins[i]) }`
- `validate` fails closed on a residual placeholder (`config.go:448-452`):
  `if strings.Contains(o, "${") { issues = append(issues, fmt.Sprintf("server.allowed_origins[%d] contains an uninterpolated ${...} placeholder", i)) }`

`internal/api/middleware.go`'s `CORS` function (lines 372-419) is the same
hardened implementation the register already cited as existing-but-unreachable
— it is now the **only** CORS middleware registered on the **only** router
(`buildRouter`), confirmed by the `setupAPI`/`RunStdio`/`waitForShutdown` call
graph in `main.go`: there is exactly one `ListenAndServe` call in the whole
file (`main.go:319`), a structural fact the `1a1a3f3` commit message itself
states as its purpose ("collapses to ONE hardened listener").

Existing test evidence (both files unchanged by this bundle, read to confirm
what is ALREADY proven, not to invent new claims):

- `internal/api/middleware_test.go` carries 8 tests directly against the
  `CORS()` function in isolation: `TestCORS_AllowlistedOrigin`,
  `TestCORS_NonAllowlistedOriginIsNeverReflected`,
  `TestCORS_WildcardAllowsAnyOriginWithoutCredentials`,
  `TestCORS_EmptyAllowlistFailsClosed`,
  `TestCORS_PreflightRejectedForDisallowedOrigin`,
  `TestCORS_PreflightAcceptedForAllowlistedOrigin`,
  `TestCORS_MultipleOriginsInAllowlist`,
  `TestCORS_NoOriginHeaderIsTreatedAsSameOriginOrNonBrowser` (lines 35-252).
- `cmd/server/security_test.go`'s `TestNoWildcardCORSOnLivePaths` (lines
  102-127) exercises the **real, production `buildRouter`** (not a mock) via
  `httptest`, across `/`, `/mcp/v1/tools`, `/api/v1/skills`, and
  `/mcp/v1/tools/skill_create/call`, with a disallowed `Origin` header, and
  asserts no live response ever carries `Access-Control-Allow-Origin: *`.

**Verdict (FACT):** the security-critical half of G18 — "the CORS allowlist is
unreachable on the live path" — **is resolved by `1a1a3f3` (G01)**. The wildcard
`corsMiddleware()` no longer exists; `cfg.Server.AllowedOrigins` is read by the
one live router; the empty-allowlist fail-closed case is proven end-to-end
against the real `buildRouter`, not merely against the isolated `CORS()`
function.

Two sub-clauses of G18's own title and its own DECISION text are, however,
**still open**, confirmed by direct reads:

**(a) `SPEC.md` §8 config sample still omits `allowed_origins`.** `SPEC.md`
lines 374-434 (the `## 8. Configuration (config.toml)` fenced block) enumerate
`[server]`, `[database]`, `[embedding]`, `[validation]`, `[autoexpand]`,
`[codeanalysis]`, `[mcp]`, `[registry]`, `[logging]` — the `[server]` stanza
(lines 377-384) has no `api_keys`, `auth_disabled`, or `allowed_origins` line
at all. `grep -n allowed_origins SPEC.md` returns zero hits. This is the exact
condition G18 named ("SPEC config sample omits `allowed_origins`") and it is
unchanged by `1a1a3f3` (that commit touched no `.md` file — its diff is
Go-source only, confirmed by `git show 1a1a3f3 --stat`).

**(b) The real shipped `config/config.toml` template, by contrast, now
documents it.** Lines 11-28 (post-`1a1a3f3`) carry a full fail-closed
explanation block for `api_keys`, `auth_disabled`, **and** `allowed_origins`,
each commented out with the safe empty default and a one-line semantic note
("Empty (the default) is fail-closed — no cross-origin browser access, and no
wildcard `*` is ever emitted."). This satisfies the "document it in
`config.toml`... with a safe example, default empty (fail-closed)" clause of
G18's own DECISION text.

**(c) No test exercises a *populated* `AllowedOrigins` against the *live*
`buildRouter`.** Every existing live-router test
(`newTestRouter` in `security_test.go:39-54`) hard-codes an **empty**
`cfg.Server.AllowedOrigins` (comment at line 45: "`cfg.Server.AllowedOrigins`
is empty -> fail-closed CORS (no `"*"`)"). The 8 `CORS()`-level unit tests in
`middleware_test.go` DO exercise non-empty allowlists, but only against the
`CORS()` middleware function called in isolation — never through the full
`buildRouter` assembly the way `TestNoWildcardCORSOnLivePaths` does for the
empty case. `research/testing_infrastructure_plan.md:304` lists
"integration(config allowlist honoured end-to-end)" as G18's required test
type and does **not** mark it done — unlike G01's own row
(`testing_infrastructure_plan.md:287`, which explicitly says
"unit(**done**: CORS+auth)"). The plan's own bookkeeping agrees with this
read: the integration-level, populated-allowlist, live-router proof is the one
piece of G18's required coverage that has not yet landed.

### 1.3 Decision

G18 is **narrowed, not reopened wholesale**. The runtime-security defect (the
reason G18 was `MED` severity, not merely a docs nit) is closed and should be
recorded as such. The residual is a **doc-drift + missing-integration-test**
item, correctly scoped as `Task`/`Bug`-low, not the original `MED`
weakness/spec-drift finding:

1. **Recommend closing the runtime-wiring clause of G18** per §11.4.90, citing
   commit `1a1a3f3` + the file:line evidence in §1.2 above as the closure
   proof — this is exactly the case the task brief anticipated ("If G18 is
   fully resolved by G01, say so as FACT... and recommend closing it") for the
   security-critical sub-clause. The correct §11.4.90 reason is
   `superseded-by-later-mandate` (G01's remediation superseded G18's own
   "wire `AllowedOrigins` end-to-end" ask) with `Superseding-item: G01 /
   1a1a3f3`.
2. **Open a narrowed successor** (`G18-residual`, or reuse the `G18` id with an
   updated title — the register/DB migration in P1.T1 governs which; this
   doc does not prescribe the ID mechanics) tracking exactly two remaining
   actions:
   - Add an `[server]` block to `SPEC.md` §8 documenting `api_keys`,
     `auth_disabled`, and `allowed_origins` with the SAME safe-default
     language already present in `config/config.toml` lines 11-28 (copy the
     semantics, not the comment marks — SPEC.md's sample additionally needs
     its pre-existing `--`-vs-`#` defect, G19, fixed in the same edit since
     both touch the identical fenced block; this bundle does not re-litigate
     G19's own decision, only notes the two edits land together to avoid a
     second churn of the same lines).
   - Add the one integration test in §1.4 below.
3. **Alternatives rejected:** leaving G18 fully open (ignores the proven,
   file:line-cited runtime fix that already landed — a §11.4.6 no-guessing-mandate
   violation in the other direction, since it would imply unproven uncertainty
   where a cited FACT exists); closing G18 with zero residual (ignores the
   SPEC.md drift that is still, factually, present).

### 1.4 Runtime signature (§11.4.108)

- **SOURCE:** `SPEC.md`'s `## 8. Configuration` fenced `toml` block contains an
  `allowed_origins` (and `api_keys`, `auth_disabled`) line under `[server]`.
- **RUNTIME-ON-CLEAN-TARGET:** a new test —
  `TestCORS_AllowlistHonoredEndToEndOnLiveRouter` (or equivalent name) in
  `cmd/server/security_test.go`, extending the existing `newTestRouter` helper
  with a populated `cfg.Server.AllowedOrigins = []string{"https://app.example.com"}`
  — asserts `GET /api/v1/skills` with `Origin: https://app.example.com` and a
  valid `X-API-Key` returns `Access-Control-Allow-Origin:
  https://app.example.com` + `Access-Control-Allow-Credentials: true`, AND the
  same request with `Origin: https://evil.example.com` returns neither header,
  run against the REAL `buildRouter` (never a mock).
- **USER-VISIBLE:** an operator who sets `allowed_origins =
  ["https://app.example.com"]` in `config.toml` per the (now-fixed) SPEC.md
  instructions gets a working credentialed cross-origin browser client against
  `/api/v1`, and a reader of SPEC.md alone (without reading the Go source)
  learns the key exists.

---

## 2. G25 — `RemoveDependency` ignores name-lookup errors → audit log with empty names

### 2.1 Root cause

`internal/skill/graph.go:99-131` (`RemoveDependency`, quoted in full — the
defect lines are 103-104):

```go
// graph.go:98-131
func (s *Store) RemoveDependency(ctx context.Context, skillID, dependsOn uuid.UUID) error {
	return s.pool.WithTx(ctx, func(tx pgx.Tx) error {
		// Get names for audit log
		var fromName, toName string
		_ = tx.QueryRow(ctx, `SELECT name FROM skills WHERE id = $1`, skillID).Scan(&fromName)
		_ = tx.QueryRow(ctx, `SELECT name FROM skills WHERE id = $1`, dependsOn).Scan(&toName)

		cmdTag, err := tx.Exec(ctx, `
			DELETE FROM skill_dependencies WHERE skill_id = $1 AND depends_on = $2
		`, skillID, dependsOn)
		if err != nil {
			return fmt.Errorf("delete dependency: %w", err)
		}
		if cmdTag.RowsAffected() == 0 {
			return fmt.Errorf("dependency not found: %s -> %s", skillID, dependsOn)
		}

		// Recalculate missing deps for the source skill
		if err := s.recalcMissingDeps(ctx, tx, skillID); err != nil {
			return fmt.Errorf("recalc missing deps: %w", err)
		}

		// Audit log
		if err := s.logAudit(ctx, tx, "dependency.removed", &skillID, map[string]interface{}{
			"from": fromName,
			"to":   toName,
		}); err != nil {
			return fmt.Errorf("log audit: %w", err)
		}

		return nil
	})
}
```

Lines 103-104 discard the `Scan` error via `_ =`. If either skill row is
absent (already deleted, or the caller supplied a stale ID) `fromName`/
`toName` remain their Go zero value — the empty string `""` — and the audit
detail persisted at lines 122-125 is `{"from":"","to":"<real-name-if-lookup-succeeded>"}`,
indistinguishable from a (schema-impossible, but unverifiable from the audit
record alone) skill legitimately named `""`.

### 2.2 Current-state determination

**Unresolved, untouched by G01.** `1a1a3f3`'s diff (confirmed via `git show
1a1a3f3 --stat`) touches only `cmd/server/main.go`, `internal/api/*`, and
`internal/config/config.go` — `internal/skill/graph.go` is not in that
commit's tree. Re-reading `graph.go` at `255061b` (the current baseline)
reproduces the exact `_ =`-discard lines quoted above, byte-for-byte matching
the register's original citation.

`internal/skill/graph_test.go` (full read, all 5 `Test*` functions) has **zero
coverage of `RemoveDependency`** — only `AddDependency`'s two DB-independent
guard clauses (`TestAddDependency_SelfReferenceIsRejectedAsCycle`,
`TestAddDependency_InvalidRelationTypeIsRejected`), `hasCycle` and
`AddDependency`'s DB-bound behavior (both honestly `t.Skip`-marked
`_RequiresLiveDatabase`, lines 68-72 and 78-81 — the file's own documented
rationale is that `pool.WithTx` against a nil pool panics, and faking a
`pgx.Tx` "would test the mock, not the production SQL... exactly the kind of
bluff the task forbids"), and `collectDepNames` (an unrelated pure
import/export helper). This confirms G25 is both functionally open and
zero-test-covered today.

### 2.3 Decision

Extract the "build the audit detail from two independent name-lookup
results" logic into a **pure helper function**, mirroring this file's own
established idiom for testability (`collectDepNames` in
`import_export.go` is exactly this pattern: pure logic extracted so it can be
unit-tested without a live `pgx.Tx`, per `graph_test.go:89-...`
`TestCollectDepNames`):

```go
// buildRemovalAuditDetail(fromName string, fromErr error, toName string, toErr error) map[string]interface{}
// Records an explicit, distinguishable marker when a name lookup failed
// instead of silently substituting the empty string.
func buildRemovalAuditDetail(fromName string, fromErr error, toName string, toErr error) map[string]interface{} {
	detail := map[string]interface{}{}
	if fromErr != nil {
		detail["from_lookup_error"] = fromErr.Error()
	} else {
		detail["from"] = fromName
	}
	if toErr != nil {
		detail["to_lookup_error"] = toErr.Error()
	} else {
		detail["to"] = toName
	}
	return detail
}
```

`RemoveDependency` then captures each `Scan` error explicitly (`fromErr :=
tx.QueryRow(...).Scan(&fromName)`, `toErr := tx.QueryRow(...).Scan(&toName)`)
and calls `s.logAudit(ctx, tx, "dependency.removed", &skillID,
buildRemovalAuditDetail(fromName, fromErr, toName, toErr))`. This satisfies
the register's own DECISION text verbatim — "Capture names best-effort but
record the not-found condition explicitly in the audit detail" — and does so
in a way that is unit-testable in isolation (no live Postgres needed for the
detail-building logic itself), while the end-to-end wiring is still proven
against a live database per §11.4.108 (source ≠ runtime).

**Alternatives rejected:** ignoring the error silently (the current, defective
behavior — weakens the audit trail per R11, explicitly rejected in the
register's own text); failing the whole `RemoveDependency` transaction when a
name lookup fails (over-corrects — the deletion itself may still be valid and
desired even when the audit's cosmetic name is unavailable; the register frames
this as an audit-fidelity gap, not a correctness-blocking one).

### 2.4 Runtime signature (§11.4.108)

- **SOURCE:** `graph.go`'s `RemoveDependency` calls `buildRemovalAuditDetail`
  with captured (non-discarded) `Scan` errors.
- **ARTIFACT:** `go build ./...` + `go vet ./...` clean; `grep -n '_ = tx.QueryRow' internal/skill/graph.go` returns zero hits inside `RemoveDependency`.
- **RUNTIME-ON-CLEAN-TARGET:** against a fresh/migrated Postgres, calling
  `RemoveDependency` where one endpoint skill row was deleted out-of-band
  before the `Scan` runs produces an `audit_log` row whose `detail` JSONB
  contains `"from_lookup_error"` (or `"to_lookup_error"`) — never a bare
  `"from":""` indistinguishable from success.
- **USER-VISIBLE:** an operator inspecting the audit trail for a
  `dependency.removed` event can tell "the skill legitimately had this name"
  apart from "the name lookup failed", restoring the evidentiary value R11
  requires of the audit log.

### 2.5 Test-case count (RED-first, §11.4.43/§11.4.115)

Reconciling with `research/testing_infrastructure_plan.md:311` — required
types: "unit(audit detail records missing name), regression":

| # | Test | Type | RED (today) behavior | GREEN (post-fix) behavior |
|---|------|------|----------------------|----------------------------|
| 1 | `TestBuildRemovalAuditDetail_SourceLookupFailureIsExplicit` | unit (pure function, no DB) | function does not exist yet → compile FAIL | `detail["from_lookup_error"]` set, `detail["from"]` absent |
| 2 | `TestBuildRemovalAuditDetail_TargetLookupFailureIsExplicit` | unit (pure function, no DB) | same | `detail["to_lookup_error"]` set, `detail["to"]` absent (independent of #1 — both `Scan` calls must be verified separately per §11.4.194's multi-factor mandate: a single combined test would not independently prove each of the two error paths) |
| 3 | `TestBuildRemovalAuditDetail_BothLookupsSucceed_UnchangedShape` | unit (pure function, regression floor) | n/a (new) | `detail == {"from":<name>,"to":<name>}`, byte-identical to today's success-path shape — proves the refactor does not regress the common case |
| 4 | `TestRemoveDependency_AuditRecordsExplicitLookupFailure_RequiresLiveDatabase` | integration, honestly `t.Skip`-marked absent a live/containerized Postgres (matching the file's own established idiom at lines 68-72/78-81) | today: persisted `audit_log.detail` is `{"from":"","to":...}` | post-fix: persisted `audit_log.detail` carries the explicit marker — proves the pure helper is actually WIRED into the live `RemoveDependency` path (§11.4.108 source≠runtime) |
| 5 | Paired §1.1 mutation | mutation | revert the call site to the bare `_ = tx.QueryRow(...).Scan(&fromName)` (bypass the helper) | tests #1/#2/#4 FAIL by reproducing the empty-string/no-marker output |

**Total: 4 new tests + 1 paired mutation = 5 test-cases.**

---

## 3. G26 — `${VAR:-default}` cannot resolve to an intentionally-empty value

### 3.1 Root cause

`internal/config/config.go`'s `interpolate` function (lines 346-370, defect at
361-367):

```go
// config.go:348-370
func interpolate(s string) (string, error) {
	result := envVarRegex.ReplaceAllStringFunc(s, func(match string) string {
		inner := match[2 : len(match)-1] // strip ${ and }

		// Check for default syntax: VAR:-default
		var envKey, defaultVal string
		if idx := strings.Index(inner, ":-"); idx >= 0 {
			envKey = inner[:idx]
			defaultVal = inner[idx+2:]
		} else {
			envKey = inner
		}

		if v := os.Getenv(envKey); v != "" {
			return v
		}
		if defaultVal != "" {
			return defaultVal
		}
		return ""
	})
	return result, nil
}
```

Line 361 (`if v := os.Getenv(envKey); v != ""`) treats "environment variable
unset" and "environment variable explicitly set to the empty string"
identically — both make `os.Getenv` return `""`, both fall through to the
default branch. An operator cannot use `SOME_VAR=""` to intentionally blank an
override; the documented default always wins in that case.

### 3.2 Current-state determination

**Unresolved, untouched in substance by G01.** `1a1a3f3` added two NEW
call-sites into `substituteEnv` (`config.go:309-311`, interpolating
`cfg.Server.AllowedOrigins`) but those call sites route through the SAME `sub`
closure (`config.go:290-297`) and the SAME `interpolate` function quoted
above — `1a1a3f3`'s diff does not touch `interpolate`'s body (confirmed: the
function's line range shifted down by exactly the 3 lines G01 inserted above
it for `AllowedOrigins` interpolation, but its logic — `os.Getenv(envKey); v
!= ""` — is byte-identical to the register's own citation of
`config.go:342-348` pre-`1a1a3f3`).

`internal/config/config_test.go`'s `TestInterpolate` (lines 17-45) enumerates
6 cases: plain string, set-with-value (no default), set-with-value
(has-default, default ignored), unset-with-default, unset-without-default,
embedded/multiple placeholders. **None of the 6 sets the environment variable
to an explicit empty string** (`t.Setenv(name, "")`) — the exact input that
discriminates this defect. `TestSubstituteEnv_AppliesAcrossAllDocumentedFields`
(lines 47-73) likewise only exercises unset-with-default
(`HELIX_TEST_DB_HOST_UNSET`) and unset-without-default
(`HELIX_TEST_APIKEY_UNSET`) — never an explicitly-empty override. This
confirms the gap is real, current, and presently uncovered by any test in the
suite.

### 3.3 Decision

Replace `os.Getenv(envKey); v != ""` with `os.LookupEnv(envKey)`, which
returns `(value, ok)` — `ok` is `true` whenever the variable is present in the
environment (including when its value is `""`) and `false` only when it is
genuinely absent:

```go
if v, ok := os.LookupEnv(envKey); ok {
	return v
}
if defaultVal != "" {
	return defaultVal
}
return ""
```

This is exactly the register's own DECISION text: "Distinguish 'unset'
(`os.LookupEnv`) from 'empty' so an explicit empty override is honoured."

Enumerating the full input space per §11.4.194 (a fix touching one branch of
a multi-factor condition must be verified against every combination, not only
the reported one) — six combinations of {var state} × {default present?}:

| Var state | Has default | Pre-fix result | Post-fix result | Discriminates the bug? |
|---|---|---|---|---|
| set, non-empty | no | value | value | no (unaffected) |
| set, non-empty | yes | value | value | no (unaffected) |
| unset | yes | default | default | no (must NOT regress — existing test line 28) |
| unset | no | `""` | `""` | no (must NOT regress — existing test line 29) |
| **set, empty** | **yes** | **default (WRONG)** | **`""` (correct)** | **yes — the bug** |
| set, empty | no | `""` | `""` | no (both paths already agree — non-discriminating, but still asserted for completeness so the fix's negative space is proven, not assumed) |

Four of these six rows are already covered by the existing `TestInterpolate`
table (lines 25-31); the fix must preserve all four unchanged (this is the
mutation's negative control) while correcting the fifth row and confirming the
sixth stays correct.

**Alternatives rejected:** documenting the quirk instead of fixing it (the
register's own rejected alternative — "still astonishing"); wrapping every
callsite in a special-case check instead of fixing `interpolate` itself (would
require touching all ~20 fields `substituteEnv` interpolates individually,
where fixing the one shared primitive fixes all of them at once).

### 3.4 Runtime signature (§11.4.108)

- **SOURCE:** `interpolate`'s body calls `os.LookupEnv`, not `os.Getenv() !=
  ""`.
- **ARTIFACT:** `go build ./...` clean; `grep -n 'os.LookupEnv' internal/config/config.go` returns a hit inside `interpolate`.
- **RUNTIME-ON-CLEAN-TARGET:** with `HELIX_TEST_VAR_EMPTY=""` exported in the
  process environment, `interpolate("${HELIX_TEST_VAR_EMPTY:-fallback}")`
  returns `("", nil)` — not `("fallback", nil)`.
- **USER-VISIBLE:** an operator who exports (e.g. via a systemd `EnvironmentFile`
  or `.env`) an explicit empty value for a config override — for example to
  intentionally disable a default TLS cert/key path or blank an optional
  provider setting — gets that blank value honored by `config.Load` instead of
  silently falling back to the documented default.

### 3.5 Test-case count (RED-first, §11.4.43/§11.4.115)

Reconciling with `research/testing_infrastructure_plan.md:312` — required
types: "unit(empty-override honoured via LookupEnv; unset uses default),
regression":

| # | Test | Type | RED (today) | GREEN (post-fix) |
|---|------|------|--------------|-------------------|
| 1 | New `TestInterpolate` table row: `"set-to-empty variable with default honors the empty override, not the default"` — `t.Setenv("HELIX_TEST_VAR_EMPTY", "")`, input `${HELIX_TEST_VAR_EMPTY:-fallback}`, want `""` | unit | returns `"fallback"` (FAIL) | returns `""` |
| 2 | New `TestInterpolate` table row: `"set-to-empty variable without default resolves to empty (non-discriminating, asserted for completeness)"` — `t.Setenv("HELIX_TEST_VAR_EMPTY2", "")`, input `${HELIX_TEST_VAR_EMPTY2}`, want `""` | unit (regression floor, both pre- and post-fix already agree here) | returns `""` (already passes) | returns `""` (unchanged) |
| 3 | New `TestSubstituteEnv_...` case: a `Config` struct field (e.g. `cfg.Database.Host`) set to `${HELIX_TEST_DB_HOST_EMPTY:-localhost-fallback}` with `HELIX_TEST_DB_HOST_EMPTY=""` exported — asserts the field resolves to `""` after `substituteEnv`, proving the fix propagates through the whole struct-walking layer, not merely the raw `interpolate` primitive | unit/integration (struct-level) | field becomes `"localhost-fallback"` (FAIL) | field becomes `""` |
| 4 | Paired §1.1 mutation | mutation | revert `os.LookupEnv` back to `os.Getenv(envKey); v != ""` | test #1 (and #3) FAIL by reproducing `"fallback"`/`"localhost-fallback"`; the four pre-existing `TestInterpolate` rows (26-29) and test #2 above stay green under BOTH implementations, proving the mutation discriminates exactly the intended invariant and nothing else |

**Total: 3 new tests + 1 paired mutation = 4 test-cases.**

---

## 4. Reconciliation summary vs `research/testing_infrastructure_plan.md`

| Gap | Plan row | Plan's required test types | This bundle's coverage |
|---|---|---|---|
| G18 | line 304 | integration(config allowlist honoured end-to-end), security(non-allowlisted origin blocked), contract(config documents the key) | integration = §1.4 new live-router test; security = same test's negative sub-case + the ALREADY-LANDED `TestNoWildcardCORSOnLivePaths`; contract = SPEC.md §8 update (§1.3 item 2) |
| G25 | line 311 | unit(audit detail records missing name), regression | §2.5 items 1-3 (unit, pure helper) + item 4 (integration, honestly `_RequiresLiveDatabase`) + item 5 (mutation = the permanent regression guard per §11.4.135) |
| G26 | line 312 | unit(empty-override honoured via LookupEnv; unset uses default), regression | §3.5 items 1-3 (unit) + item 4 (mutation = the permanent regression guard); "unset uses default" is the PRE-EXISTING `TestInterpolate` rows 28-29, which this bundle's mutation table proves are preserved, not merely assumed |

The plan's own fold-note (`testing_infrastructure_plan.md:316`, "G18→G01")
anticipated exactly the disposition reached in §1.3 above — this bundle
confirms, rather than contradicts, the plan's prior framing, and narrows G18's
open residual to precisely the two items the plan's own row 304 had not yet
marked done.

---

## 5. Honest gaps (§11.4.6)

- **G18 is not "fully" resolved — it is resolved for the security-critical
  runtime-reachability defect and open for the doc/test residual.** Both
  halves are stated as FACT above with file:line citations; neither is an
  unproven claim. Recommending a full §11.4.90 closure without the narrowed
  successor would itself be a loss-of-requirements violation (§11.4.197) —
  SPEC.md's `allowed_origins` omission is real, present, and unowned by any
  other open gap in the register (G19 covers the `--`-vs-`#` syntax defect in
  the SAME fenced block, not the missing key).
- **G25's fix requires a design decision (extracting `buildRemovalAuditDetail`)
  this doc did not find pre-existing in the codebase** — it is proposed here
  as the mechanism, following the file's own established
  `collectDepNames`-style pure-extraction idiom, but has not been implemented
  or tested against a live database as part of this design-research pass.
  Whether the eventual implementer names the map keys `from_lookup_error`/
  `to_lookup_error` exactly as proposed, or chooses an alternative
  distinguishable-sentinel scheme, is an implementation decision this bundle
  does not treat as fixed in stone — the REQUIREMENT (an explicit,
  distinguishable not-found marker, never a bare empty string) is the binding
  part.
- **G26's fix is a single-line, well-understood standard-library substitution**
  (`os.Getenv` → `os.LookupEnv`) with no ambiguity in Go's documented
  semantics; the only genuine open question is test-authoring effort, which
  §3.5 sizes at 3 new tests + 1 mutation.
- **No code was written, edited, or committed by this design-research pass.**
  All file:line citations above were read from the immutable baseline
  `255061b`; the current (uncommitted) working tree contains an unrelated G02
  change that this bundle never read.
- **This bundle does not verify `go build`/`go vet`/`go test` against these
  three files** — that is deliberately out of scope for a design-research
  document (§11.4.108's ARTIFACT/RUNTIME layers are proposed as future
  runtime signatures, not claimed as already-passing here).
