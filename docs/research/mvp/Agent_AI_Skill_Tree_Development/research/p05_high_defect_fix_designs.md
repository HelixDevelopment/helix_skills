# P05 — HIGH-Defect Fix Designs (G35, G31, G29, G32)

**Revision:** 1
**Last modified:** 2026-07-15T18:25:00Z
**Status:** design-only — no code changed by this document
**Authority:** Constitution §11.4.102 (systematic-debugging), §11.4.132 (risk-ordered validation
priority), §11.4.194 (exhaustive all-scenario review, prove-every-assumption), §11.4.197
(research/gaps driven to full completion), §11.4.124 (investigate-before-remove), §11.4.6
(no-guessing), §11.4.108 (four-layer fix verification), §11.4.110 (change-impact clash detection),
§11.4.115 (RED-then-GREEN polarity), §11.4.135 (permanent regression guard), §11.4.199
(exact-reproduction-sequence), §11.4.84 (one Go mutator at a time — Go IMPL is gated behind the
in-flight P1.T1 commit; this document is DESIGN ONLY).
**Scope:** the four HIGH-severity defects from the discovery-pressure sweep
(`research/p05_completion_audit_and_discovery.md` §A.2/A.5/A.6/A.7) as filed in
`GAPS_AND_RISKS_REGISTER.md` G29/G31/G32/G35. Every root-cause claim below was re-verified by
reading the actual source at `docs/research/mvp/Agent_AI_Skill_Tree_Development/project/` (current
working tree, not the register's prose) — file:line citations point to what was actually read in
this pass, not copied from the register/audit text.
**Read discipline (§11.4.199 / §11.4.6):** every fix design cross-checks its own root cause against
live source (not the register's summary), enumerates real call sites via direct search of the
current tree, and states explicitly wherever a claim is static-reasoning (no code run) rather than
captured-evidence (see Honest Gaps, §5).

---

## Order (risk/independence, per §11.4.132)

1. **G35** — client↔server auth-header mismatch (HIGH, live client break, zero external
   dependencies, cheapest fix — done first).
2. **G31** — `learn_from_project` path-traversal/LFI (HIGH-latent security, must land with/before
   G03, no dependency on G10/G11/G29).
3. **G32** — `ReviewScheduler` dead flagship (HIGH, depends on G11's advisory-lock scaffolding
   landing in the same window — see §4's sequencing note).
4. **G29** — `Store.Search` doc-bluff / dead `VectorSearch` (HIGH, depends on G10's embedder
   factory + boot-time dimension assertion having landed — largest blast surface of the four).

---

## §1. G35 — CLI/TUI send `Authorization: Bearer`, server reads `X-API-Key` only

### 1.1 Confirmed root cause (re-read, not assumed)

- **Server** — `internal/api/middleware.go:280-292` (`APIKeyAuth`): `key := c.GetHeader("X-API-Key")`
  is the **only** header read. The surrounding comment (`middleware.go:288-291`) is explicit and
  deliberate: *"Credentials must travel in the X-API-Key header"* — the `api_key` query-parameter
  fallback was intentionally removed for the same reason (secrets leaking via access logs/proxies).
- **CLI** — `cmd/cli/commands/common.go:54-56`, inside `(*APIClient).Request`:
  ```go
  if c.APIKey != "" {
      req.Header.Set("Authorization", "Bearer "+c.APIKey)
  }
  ```
  This is the CLI's **only** credential header, confirmed by reading the whole function
  (`common.go:39-70`) — no other header-set call touches auth.
- **TUI** — `cmd/tui/api_client.go:75-77`, inside the unexported `(*APIClient).request`, the
  identical pattern: `req.Header.Set("Authorization", "Bearer "+c.apiKey)`.
- **Governing contract, independently verified (not assumed from the register):**
  - `docs/research/mvp/Agent_AI_Skill_Tree_Development/SPEC.md:358`: *"All endpoints under
    `/api/v1`. Auth via `X-API-Key` header."*
  - `docs/research/mvp/Agent_AI_Skill_Tree_Development/api/openapi.yaml:18,42,1026-1030`:
    `securitySchemes.ApiKeyAuth = {type: apiKey, in: header, name: X-API-Key}`, referenced globally
    at line 42 as the API's one security requirement. **No `bearerAuth`/`Authorization` scheme
    exists anywhere in the committed OpenAPI document.**

Both governing documents agree unambiguously: `X-API-Key` is the one documented mechanism. The two
first-party clients are simply non-compliant with their own already-published contract.

### 1.2 Fix approach + exact seam

**DECISION: fix the two clients; do NOT widen the server.** Change:
- `cmd/cli/commands/common.go:55` → `req.Header.Set("X-API-Key", c.APIKey)`
- `cmd/tui/api_client.go:76` → `req.Header.Set("X-API-Key", c.apiKey)`

No production server change.

**Justification (§11.4.6 pick-and-justify, no guessing):** SPEC.md and the committed OpenAPI
document were independently re-read and both declare `X-API-Key` as the *sole* documented scheme.
**Alternative rejected** — widen `APIKeyAuth` to also accept `Authorization: Bearer`: this would (a)
contradict the already-published, already-committed API contract, creating a *second* G09-class
(OpenAPI↔implementation) drift the moment it landed, and (b) leave a permanent, undocumented,
untested second credential-acceptance path inside the fail-closed auth middleware. Fixing the
clients is spec-conformance; widening the server is scope creep that manufactures a new defect
class while "fixing" this one.

### 1.3 RED test + GREEN guard (§11.4.115)

**RED (reproduces on current code):** a contract test builds a real `gin.Engine` wired with the
actual, unmodified `api.APIKeyAuth([]string{"testkey"})` behind one protected route (no stub — the
real middleware). Separately, it drives the **real, production** client code —
`(*commands.APIClient).Request` and `(*main.APIClient).request` — against that test server with
`APIKey`/`apiKey = "testkey"` (§11.4.199: the test exercises the exact code path that broke, never a
hand-rolled HTTP request standing in for it). On current code: the client sends `Authorization:
Bearer testkey`; `APIKeyAuth`'s `c.GetHeader("X-API-Key")` returns `""`; the middleware responds
`401 missing_api_key` (`middleware.go:294-299`). The test asserts status `401` + error code
`missing_api_key` — true today, proving the defect is live, not hypothetical.

**GREEN guard (§11.4.115 polarity):** the identical test, `RED_MODE=0` post-fix, asserts a
non-`401` (the real protected-route body) is returned. One source file, two roles — the bug-catcher
is the permanent regression guard (§11.4.135).

### 1.4 Blast-radius PROOF (§11.4.194)

1. Every CLI subcommand (`cmd/cli/commands/{expand,learn,registry,search,skill}.go`) routes through
   the **one** shared `(*APIClient).Request` — confirmed by grepping for a second HTTP entry point
   in `cmd/cli`; none exists. One seam fixes every subcommand uniformly.
2. Every TUI screen (`cmd/tui/{browse,main,model,registry,search,tree}.go`) routes through the
   **one** shared `(*APIClient).request` (its `ListSkills`/`GetSkill`/etc. all call `c.request(...)`
   — confirmed by reading `api_client.go` in full). One seam, one fix.
3. **Riskiest blast-radius item:** any future client (Aurora/HarmonyOS per G15, or an external
   integrator) copying the CLI/TUI's *current* (wrong) pattern as a reference before this fix lands.
   This is a documentation/timing risk, not a code risk — mitigated by landing the fix + the
   contract test together so the corrected pattern is the only one present once merged.
   `grep -rn "X-API-Key"` across every non-test `.go` file returns exactly 3 hits, all inside
   `internal/api` (re-confirmed live) — there is no hidden third caller of the wrong header today.
4. `CORS`'s `Access-Control-Allow-Headers` allowlist (`middleware.go:410-411`) already permits both
   `Authorization` and `X-API-Key` in a browser preflight — that list is a browser-permission
   allowlist, not an auth-acceptance path; unaffected either way.
5. Must-not-break: neither client's `Accept`/`Content-Type` header logic is touched — confirmed by
   reading both functions end-to-end; the fix is a single-line, single-field change in each file.

### 1.5 Test plan

- Contract tests (integration-shaped, real `gin.Engine` + real client code): CLI variant, TUI
  variant — 2 tests.
- Unit tests asserting the literal header key/value set on the constructed `*http.Request` (no
  network) for each client — 2 tests.
- Paired §1.1 mutation: revert **one** client (e.g. CLI) to `Authorization: Bearer` → that client's
  contract test FAILs (401) while the TUI's stays green — proves the tests are correctly scoped
  per-client, not coincidentally passing together.
- Regression: re-run `cmd/server/security_test.go` / `internal/api/middleware_test.go` unchanged,
  confirm no interaction.

### 1.6 Impl-sequencing note

MUST land before any client-vs-auth-server LIVE test (§11.4.130, already stated in the register).
Independent of G31/G29/G32 — no ordering constraint against them.

---

## §2. G31 — `learn_from_project` `project_path` has zero validation (latent LFI)

### 2.1 Confirmed root cause (re-read, not assumed)

- `internal/mcp/tools.go:323` (inside `registerLearnFromProject`'s handler,
  `tools.go:305-359`): `projectPath, _ := request.GetArguments()["project_path"].(string)` — the
  comma-ok form is used defensively for *type* safety only; the resulting **value** receives zero
  path validation.
- `tools.go:344`: `job, err := s.skillStore.SubmitLearningJob(ctx, projectPath, languages)` — the
  unvalidated value passed straight through.
- **Today this is inert** — re-verified: `internal/skill/store.go:570-594` (`SubmitLearningJob`)
  only `json.Marshal`s `projectPath` into an `audit_log` row (`store.go:585-588`); it never touches
  the filesystem. No call from `SubmitLearningJob` to `codeanalysis.Analyzer` anywhere in the store
  package.
- **The latent, reachable sink once G03 wires the worker:**
  `internal/codeanalysis/analyzer.go:196-253` (`discoverFiles`) — `filepath.Walk(projectPath, ...)`
  (`analyzer.go:199`) with **no root-jail of its own**, followed by unconditional
  `os.ReadFile(path)` (`analyzer.go:239`) on every matched file. `discoverFiles` is called from
  `AnalyzeProject` (`analyzer.go:101,110`), which is itself **directly, publicly callable today**
  (it does not require G03's wiring to exist as a function — confirmed by reading `AnalyzeProject`'s
  full body, `analyzer.go:98-188`; it has no dependency on `SubmitLearningJob`/the worker at all).
- **Contrast (weaker, already-present, and re-read live for this pass):**
  `internal/api/learn_handler.go:47-50` — `strings.Contains(req.ProjectPath, "..") ||
  strings.Contains(req.ProjectPath, "~")`. This rejects `../etc/passwd` but does **not** reject an
  absolute path such as `/etc` or `/root/.ssh` (neither string contains `..` nor `~`), and this
  guard exists only on the REST path (currently dead per G01's O3 residual), never on the MCP path.

### 2.2 Fix approach + exact seam

**Enforce at the `Analyzer` boundary (defense-in-depth — never rely on a single caller-side guard,
per the register's own instruction), not only at the MCP handler:**

1. New method `func (a *Analyzer) resolveProjectRoot(projectPath string) (string, error)`, called as
   the **first** line of `AnalyzeProject` (before `discoverFiles` at `analyzer.go:110`). It:
   - `filepath.Abs(projectPath)` → `filepath.EvalSymlinks(abs)` → `filepath.Clean(...)`. Resolving
     symlinks is load-bearing: a `../` component *or* a symlink pointing outside an allowlisted root
     are both caught by comparing the **fully-resolved** path, never the literal input string.
   - Any resolution error (including a dangling symlink) is a **hard reject**, never a silent
     continue.
   - The resolved path must equal, or be a path-separator-bounded child of, at least one entry in a
     new config-driven allowlist `a.cfg.AllowedRoots []string` (canonicalized once, at
     `NewAnalyzer` construction time, the same way). **Bare string-prefix comparison is explicitly
     rejected** as insufficient — `/srv/allowed-evil` would wrongly satisfy a naive
     `strings.HasPrefix(resolved, "/srv/allowed")`. The correct check is
     `resolved == root || strings.HasPrefix(resolved, root+string(os.PathSeparator))` against
     *every* allowlisted root.
   - An **empty** `AllowedRoots` is fail-closed — it rejects every path, never "allow everything"
     (§11.4.6/§11.4.201: no implicit-allow default).
   - `AnalyzeProject` uses **only** the resolved, validated path from then on — both as the
     `filepath.Walk` root and as the `ProjectPath` field written into `AnalysisResult` — never the
     raw caller-supplied string again.
2. **Config seam (project-owned, §11.4.28):** extend `internal/config.CodeAnalysisConfig`
   (`config.go:133-138`, currently `Enabled/Languages/MaxFileSizeKB/ExcludePatterns`) with
   `AllowedRoots []string \`toml:"allowed_roots"\``. Canonicalized once inside `NewAnalyzer`
   (`analyzer.go:81-92`); a configured root that does not exist or cannot be resolved is a
   **construction-time (boot) error** — fail-closed, mirroring G10's `AssertEmbeddingDimension`
   boot-time-assertion pattern — never a silently-dropped entry.
3. **Second, independent, cheap check at the MCP call site** (`tools.go:323`): reject an empty
   string immediately via the existing `s.newToolError(...)` helper (already used at `tools.go:325`
   for the same purpose) *before* even reaching the store — fast-failing, but explicitly a *subset*
   of the analyzer's rejection set, never a superset, so it can never falsely reject something the
   analyzer would allow. The analyzer-layer check (step 1) remains the *authoritative* control, so a
   future second caller of `AnalyzeProject` (e.g. a CLI `learn` command bypassing MCP) is equally
   protected without a second, independently-maintained guard.
4. `internal/api/learn_handler.go:47-50`'s existing weaker REST guard is left in place (harmless,
   first-line, no regression) but documented as non-authoritative.

### 2.3 RED test + GREEN guard (§11.4.115)

Security tests call `(*Analyzer).AnalyzeProject(ctx, ...)` **directly** (no G03 wiring needed,
since `AnalyzeProject` is independently callable today):
- `AnalyzeProject(ctx, "../../../etc")`
- `AnalyzeProject(ctx, "/etc")`
- `AnalyzeProject(ctx, <temp-dir symlink pointing outside the allowed root>)`

**RED:** on current code (no `AllowedRoots` field exists at all today), every one of these
**succeeds** and returns file content read from outside the intended root — proving the
vulnerability is live *today* at the analyzer layer, independent of G03's production-reachability
timeline.

**GREEN guard:** identical calls, post-fix, each return a named error (e.g.
`ErrProjectPathNotAllowed`) and `discoverFiles` reads zero bytes outside the allowlist.

### 2.4 Blast-radius PROOF (§11.4.194)

1. **Riskiest blast-radius item:** `internal/codeanalysis/analyzer_test.go` (one of the 8 existing
   test files) very plausibly constructs `Analyzer`/calls `AnalyzeProject` against a temp-dir
   fixture with `AllowedRoots` unset today. A fail-closed default landed **without** updating that
   test file's fixtures would turn every existing `codeanalysis` test RED — a self-inflicted
   regression. The implementation MUST audit and update every `NewAnalyzer(...)` construction inside
   `analyzer_test.go` (add the fixture's temp dir to `AllowedRoots`) in the **same** change. This was
   not independently re-read in this design pass (see Honest Gaps §5.4) — flagged as a mandatory
   pre-implementation check, not assumed safe.
2. `internal/mcp/tools.go`'s `registerLearnFromProject` — the added cheap-reject only changes the
   error *message* for empty input; no other behavior change.
3. `internal/skill/store.go`'s `SubmitLearningJob` is unchanged by this fix — it is not the sink.
   G30's eventual real job-status/execution wiring is what will call `AnalyzeProject` in production;
   this fix's entire purpose is to already be in place before that lands.
4. `internal/worker/runner.go`'s stub `handleAutoExpand`/`runAutoExpandCycle` (currently never call
   `AnalyzeProject`, per G03/G11) are unaffected today; the moment G03 wires a real call, they
   inherit this protection for free — that is the reason the check lives at the analyzer layer and
   not only at the MCP handler.
5. `internal/codeanalysis/analyzer.go:692` (`MapToSkills`, currently uncalled, see §3.1 below for
   context) does not call `AnalyzeProject`/`discoverFiles` — it calls `store.Search`, unaffected by
   this fix.
6. No other caller of `AnalyzeProject`/`discoverFiles` exists in the tree today (re-confirmed: one
   definition, not yet invoked from `cmd/worker` or `cmd/server`).

### 2.5 Test plan

- Security: traversal `../../etc`, absolute `/etc`, absolute `/root/.ssh`, symlink-escape via a
  temp-dir symlink pointing outside the allowlisted root — all rejected.
- Unit: in-allowlist-root path accepted; nested child of an allowlisted root accepted;
  prefix-confusable sibling (`/srv/allowed-evil` vs allowlisted `/srv/allowed`) rejected — the exact
  boundary-vs-prefix bug class.
- Unit: empty `AllowedRoots` config rejects every path (fail-closed-by-default proof).
- Regression: full `analyzer_test.go` suite re-run green with fixtures updated per blast-radius item
  1.
- Paired §1.1 mutation: revert the boundary check from `HasPrefix(resolved, root+sep)` to a bare
  `HasPrefix(resolved, root)` → the prefix-confusable-sibling security test FAILs.

### 2.6 Impl-sequencing note

MUST land WITH or BEFORE G03 (per the register; re-affirmed here because `AnalyzeProject` becomes
*production-reachable* from an untrusted MCP caller only once G03 wires `SubmitLearningJob`/the
worker's `handleAutoExpand`-equivalent to actually invoke it). Note, however: the RED test above can
be authored and can prove the vulnerability **now**, independent of G03's landing — only
*production* reachability by a real attacker is gated on G03. No dependency on G29/G32/G35.

---

## §3. G29 — `Store.Search` doc-bluff; `Store.VectorSearch` dead with zero callers

### 3.1 Confirmed root cause (re-read, not assumed) + a 4th call site the register did not cite

- `internal/skill/store.go:50`: `// Search performs a hybrid search combining vector similarity and
  text matching.`
- `store.go:51-118` (the actual body): two successive raw-SQL statements — primary
  `similarity(...)`/trigram (`store.go:54-62`, `s.name % $1 OR s.title % $1 OR s.description ILIKE
  ...`) then, on zero rows, an `ILIKE`-only fallback (`store.go:88-95`). **Re-confirmed: zero
  references to `embedding`/`Embedder`/`pgvector` anywhere inside `Search`'s body.**
- `store.go:596-630` (`VectorSearch`) — a separate, correctly-implemented method:
  `pgvector.NewVector(embedding)` then `ORDER BY s.embedding <=> $1` cosine-distance KNN
  (`store.go:604`).
- **Re-verified live:** `grep -rn "\.VectorSearch("` across the whole `project/` tree returns
  nothing outside its own definition — zero callers, confirmed independently of the register.
- **Re-verified live call sites of `Store.Search`** (broader than the register's own citation of
  three): `internal/mcp/tools.go:56` (MCP `skill_search`), `cmd/server/main.go:207` (live REST
  `/api/v1/skills/search`), `internal/validation/pipeline.go:540,557` (jury duplicate-skill /
  existing-dependency checks), **and a fourth site the register's G29 entry did not enumerate:
  `internal/codeanalysis/analyzer.go:692`**, inside `(*Analyzer).MapToSkills`
  (`analyzer.go:674-699+`) — maps detected code patterns to existing skills via
  `store.Search(ctx, pType, 5)` per pattern type. `MapToSkills` itself is currently uncalled in
  production (consistent with G03's dead-autoexpand-pipeline finding), but it inherits the *same*
  trigram-only weakness the moment G03 wires codeanalysis into the worker — a blast-radius fact this
  pass adds.

### 3.2 Fix approach + exact seam

1. **Embedder-DI decision:** give `*skill.Store` an *optional* embedder via an additive constructor
   change: `func NewStore(pool *db.Pool, embedder db.Embedder) *Store`. `embedder` may be `nil`
   ("no embedder configured" — triggers the existing trigram/ILIKE-only path, unchanged behavior,
   now *honestly labelled*, see point 3). **Rejected alternative:** a package-level singleton/global
   embedder — rejected because it hides the dependency, breaks per-test isolation, and contradicts
   the codebase's own existing constructor-injection convention (`NewAnalyzer(cfg,logger)`,
   `NewRunner(pool,store,cfg,logger)`).
   - Every existing call site updated (all four re-confirmed live):
     `cmd/server/main.go:93`, `cmd/worker/main.go:89` — both resolve a real embedder via the
     **already-existing** `db.NewEmbedderFromConfig(cfg.Embedding)` (`internal/db/embedding.go:293`
     — currently unused in production anywhere; this fix is its first real caller), logging and
     degrading to `nil` on a resolution error (never fail-closed the whole store — a degraded search
     is a quality regression, §11.4.3 SKIP-with-reason posture, not a security concern);
     `internal/registry/registry.go:142` (`skill.NewStore(r.pool).GetCoverage(...)`, a throwaway
     `Store` that never calls `Search`) passes explicit `nil`; `cmd/server/security_test.go:48`
     passes explicit `nil` (its own assertions never call `Search`, confirmed by reading its usage).
2. **`Search`'s body:** when `s.embedder != nil`, embed `query` via
   `s.embedder.Embed(ctx, []string{query})`, run the existing `VectorSearch` KNN query **and** the
   existing trigram query, then merge via Reciprocal Rank Fusion (`score = Σ 1/(k + rank_i)` per
   result across the two rank lists, `k = 60`) — a parameter-light merge that needs no hand-tuned
   weight and never requires normalizing cosine-distance and trigram-similarity onto one shared
   scale. When `s.embedder == nil`, **or** the embed call itself errors (network/API failure), fall
   back to the existing trigram/ILIKE-only path **unchanged** — never fail the whole search on an
   embedder hiccup (§11.4.3).
3. **Additive field** `models.SearchResult.SearchMode string` (closed set:
   `"hybrid"` | `"trigram_no_embedder_configured"` | `"trigram_embed_error"`) — directly satisfies
   the register's own DECISION language ("say so honestly in the result set, never silently").
   Additive/non-breaking to existing JSON consumers.
4. Correct the doc-comment (`store.go:50`) to state actual post-fix behavior.
5. `Store.VectorSearch` (`store.go:596-630`) is **unchanged** — it becomes the KNN half `Search` now
   calls internally, resolving the §11.4.124 dead-code finding by wiring, not deleting, per the
   register's own DECISION.

### 3.3 RED test + GREEN guard (§11.4.115)

Two skills seeded: one whose name/title/description shares **no** substring with a query but is
semantically close (e.g. titled "Container orchestration health checks", queried with "kubernetes
liveness probes" — zero shared tokens) and one whose description `ILIKE`-matches the query text but
is topically unrelated. `Search` run against a real Postgres+pgvector instance (§11.4.27, this
specific ranking claim is genuinely integration-shaped) with a real/fixed-deterministic embedder
configured.

**RED (current code):** the shallow-keyword-match skill ranks first (or the semantically-close
skill is entirely absent if it shares zero trigrams); `VectorSearch` is never invoked — reproduces
"advertised capability not delivered."

**GREEN guard:** the semantically-close skill ranks *above* the shallow-keyword-match skill, and
`result.SearchMode == "hybrid"` for both. A separate, cheaper unit test with an injected spy
`db.Embedder` proves `VectorSearch`'s SQL is actually reached (call-count / SQL-text assertion)
without needing live pgvector for that narrower claim.

### 3.4 Blast-radius PROOF (§11.4.194)

1. `internal/mcp/tools.go:56` (`skill_search`, the documented "use this FIRST" MCP tool) — its
   result JSON (`tools.go:62-79`, `resultItem{Name,Title,Description,Status,Score}`) does **not**
   currently surface `SearchMode`. **Recommended: add a `search_mode` field to the tool's own
   result JSON**, so the honesty guarantee reaches the actual MCP consumer instead of stopping at
   the Go API boundary.
2. `cmd/server/main.go:207` (live REST search route) — must be checked at implementation time for
   the same parity need.
3. **Riskiest blast-radius item:** `internal/validation/pipeline.go:540,557` (jury duplicate-skill /
   existing-dependency checks). Today these run trigram/ILIKE only; once hybrid lands, a
   *semantically*-duplicate skill phrased differently will now be caught where it was previously
   invisible — a genuine behavior change. This is the register's own stated *intent* (closing the
   "silently weakening" gap it names), but it MUST be checked against `pipeline_test.go`'s existing
   fixtures/assertions at implementation time (not read in this pass — an explicit scope boundary,
   see Honest Gaps §5.3) to confirm no currently-passing "these are NOT duplicates" assertion would
   flip.
4. `internal/codeanalysis/analyzer.go:692` (`MapToSkills`, the newly-found 4th call site, currently
   unreached in production) — inherits hybrid matching for free the moment G03 wires it live; no
   separate code change required there.
5. `internal/api/search_handler.go` is **not** a caller of `Store.Search` at all — it calls
   `s.pool.SearchSkills`, a different, already-correct vector-capable path (confirmed unaffected,
   out of scope).
6. Must-not-break: `VectorSearch`'s own SQL/behavior is untouched; any future direct caller (none
   exist today) is unaffected.

### 3.5 Test plan

- Unit: RRF-merge math against fixed rank lists (`k=60`); nil-embedder fallback proves the *old*
  trigram-only path is byte-for-byte unchanged when no embedder is configured; spy-embedder test
  proves `VectorSearch`'s query is issued when an embedder is configured.
- Integration: the semantic-ranking RED/GREEN pair above; a **separate** false-positive guard for
  `pipeline.go:540,557` proving the new hybrid check does *not* flag two genuinely-different skills
  that merely embed near each other.
- Paired §1.1 mutation: revert `Search` to call the trigram path even when an embedder is configured
  → the semantic-ranking integration test FAILs.
- Regression: `pipeline_test.go` / `graph_test.go` re-run green, citing any fixture updates.

### 3.6 Impl-sequencing note

Composes **additively** with P1.T1's `kind` column (already present and already selected in
`Search`'s SQL today, `store.go:56,76` — P1.T1 is unaffected, this fix does not touch `kind` at all)
and with G10's embedding-dimension design (`db.NewEmbedderFromConfig`/`Embedder.Dimensions()` —
this fix is the *first* real production caller of that factory). Should land **after or alongside**
G10's boot-time `AssertEmbeddingDimension` fail-closed check, never before, so a misconfigured
dimension is caught at startup rather than surfacing as a confusing pgvector error inside
`VectorSearch`. No ordering dependency on G31/G32/G35.

---

## §4. G32 — `registry.ReviewScheduler` dead flagship

### 4.1 §11.4.124 git-history investigation (performed, not cited from the register)

```
git log --oneline --all -S"StartReviewScheduler"      -- project   → a2f2631 only
git log --oneline --all -S"registryReviewWorker"       -- project   → a2f2631 only
git log --oneline --all -S"runRegistryReview"          -- project   → a2f2631 only
git log --oneline --all -S"type ReviewScheduler struct" -- project  → a2f2631 only
```

**FACT:** `registry.ReviewScheduler` (the full, cron-based mechanism) **and**
`internal/worker/runner.go`'s own, separate, much weaker `registryReviewWorker`/`runRegistryReview`
ticker (already wired, already running today — unconditionally started from `Runner.Start()`,
`runner.go:147-149`) were introduced in the **exact same** single squashed commit (`a2f2631`,
"Extract Agent AI Skill Tree Development materials"). There is no earlier commit where
`ReviewScheduler` was called and a later commit removed the call — this is **never-completed, not
wired-then-regressed**. The codebase's original author built a more-capable mechanism
(`registry.ReviewScheduler`) and, in the *same* commit, wired a different, weaker, ad-hoc mechanism
directly inside `Runner` instead.

Re-reading `runRegistryReview` (`runner.go:509-527`) confirms exactly how much weaker: it calls
**only** `r.store.GetCoverage(ctx, "")` (a read) and writes one audit-log row — it never calls
`UpdateCoverage`, `CalculateMissingDeps`, or marks any skill stale. This is a **fourth** instance of
the audit's own "two/three rival implementations, the wrong one is wired" structural pattern (after
G01/G03/G29) — a genuinely new observation this pass adds, since the register's own G32 evidence
only checked "does anything reference `ReviewScheduler`'s constructors," not "does something else
already perform an equivalent, weaker function on the same tables."

### 4.2 Why the literal register fix-direction is rejected

Naively instantiating `registry.NewHourlyReviewScheduler()`/`StartReviewScheduler` from
`cmd/worker/main.go` (the register's literal fix-direction text) would create a **second**,
independent, differently-scheduled periodic writer against the *same* `skill_registry` rows the
already-running `registryReviewWorker` ticker writes: the runner's ticker keys off
`cfg.Registry.ReviewIntervalHours` (`config.go:148`, already declared); `ReviewScheduler`'s *own*
cron/ticker-fallback would key off whatever interval `cmd/worker/main.go` separately passed it — two
uncoordinated cadences hitting the same idempotent-but-redundant `UPDATE ... WHERE ... AND
stale=false` statements. Wasteful, confusing (two "which mechanism is authoritative" answers), and a
genuine §11.4.110 change-impact clash the register's own G32 evidence did not surface. **This is the
single most important finding of this design pass for G32.**

### 4.3 Fix approach + exact seam (RECONCILED)

**DECISION — REJECTED option:** start `ReviewScheduler` as a second, independent goroutine from
`cmd/worker/main.go` alongside the runner's existing ticker (creates the dual-writer clash above).

**DECISION — CHOSEN option:** retarget the already-wired, already-config-driven
`Runner.runRegistryReview` (`runner.go:509-527`) to call the full review logic instead of merely
`GetCoverage`, making the runner's own pre-existing ticker (`registryReviewWorker`,
`runner.go:414-434`) the single owner. Zero new goroutines, zero new config keys:

1. Add `registry *registry.Registry` as a new `Runner` field, constructed in `NewRunner` via
   `registry.NewRegistry(pool)` — trivial, dependency-free (`Registry{pool}`, `registry.go:17-24`,
   no coupling risk).
2. **Consolidate a pre-existing internal asymmetry, independently found by this re-read:**
   `ReviewScheduler.performReview`'s 5-step body (`review.go:122-153` — three `markXStale` calls +
   `UpdateCoverage` + `CalculateMissingDeps`) and `Registry.RunReviewOnce`'s 4-step body
   (`review.go:274-315`) are **not the same** — `RunReviewOnce` is missing the
   `markLowCoverageSkillsStale` step that `performReview` has. Extend `RunReviewOnce` to also call
   the low-coverage stale-mark (matching `performReview`'s full 5-step set exactly), then have
   `performReview` call `RunReviewOnce` internally instead of duplicating the 5 steps inline. This is
   a same-package, zero-external-caller-impact consolidation — neither method has any caller outside
   `registry.go`/`review.go` today (re-confirmed).
3. Retarget `Runner.runRegistryReview` to call `r.registry.RunReviewOnce(ctx)` instead of
   `r.store.GetCoverage(ctx, "")`. The existing audit-log write (`runner.go:522-526`) is kept, its
   payload switching from a raw coverage snapshot to `RunReviewOnce`'s own returned
   summary/error. (Noted, out of this task's scope: `runner.go:524` reuses
   `db.AuditEventExpansionStarted` — the wrong event constant for a registry-review completion log,
   an orthogonal pre-existing weakness flagged for a separate fix, not touched here.)
4. Wrap the retargeted call in the Postgres single-owner advisory lock **already designed** (not
   invented here) by `research/g11_worker_design.md` §3: `SELECT pg_try_advisory_lock($1)` keyed on
   a per-cycle constant, non-blocking, released at cycle end, honest skip-with-reason log for a
   non-holder. `g11_worker_design.md` §3 *already names `runRegistryReview` by name* as one of the
   write-cycles needing exactly this lock (cross-checked per the task's own instruction) — this
   design does not invent a second locking scheme; it specifies that G32's consolidated
   `RunReviewOnce` call and G11's own fix share the **same** lock-acquisition helper and the same
   lock key for the `"registry_review"` cycle-class, so whichever of G11/G32 lands first, the other
   composes without re-litigating the lock design.
5. `registry.ReviewScheduler`'s own cron/ticker-fallback/`Stop`/`IsRunning` machinery
   (`review.go:14-115,216-258`) is **not** instantiated by `cmd/worker/main.go` under this design
   (that is the rejected option). It remains formally unwired as a *type*, but its *capability* is
   no longer dead — it runs, consolidated, through `RunReviewOnce` via the runner's own ticker. Per
   §11.4.124/§11.4.90, the `ReviewScheduler` type's *scheduling wrapper* (not its review-logic, which
   is genuinely reused via `RunReviewOnce`) becomes a legitimate `Obsolete` candidate with reason
   `superseded-by-design-change` — actual removal requires a separate, explicit operator
   keep-or-remove decision per §11.4.122, not decided here, flagged honestly as the natural next
   step.

### 4.4 RED test + GREEN guard (§11.4.115 / §11.4.108)

Integration test against a real Postgres: seed a skill whose `skill_registry.last_review` is 31
days old and whose status is `active`. Run the runner's ticker cycle (or call `runRegistryReview`
directly) and assert `skill_registry.stale`.

**RED (current code):** stays `false` — `runRegistryReview` only reads coverage, never marks
anything stale, reproducing the exact "the periodic health sweep never executes" defect **on the
live, already-wired code path** — a stronger, more precise reproduction than merely proving
`ReviewScheduler` has zero callers, since it exercises the actual, already-running production
goroutine and shows it does not do what the worker's own package doc (`cmd/worker/main.go:1-6`)
promises.

**GREEN guard:** post-fix, `stale = true` after the cycle runs, and `skill_registry.coverage`/
`.missing_deps` show recalculated (non-placeholder) values — a §11.4.108 runtime-signature: the
scheduler tick is observable via a real `skill_registry` row transition on a clean deploy, not a log
line.

### 4.5 Blast-radius PROOF (§11.4.194)

1. **Riskiest blast-radius item:** the dual-writer clash described in §4.2 — avoided *by
   construction* in this design (single retargeted call site, no second goroutine) rather than
   merely noted. The naive literal reading of the register's own fix-direction text would have
   reintroduced exactly the "two rival implementations" class of defect the register elsewhere warns
   against.
2. `cmd/server/main.go:94` (`skillRegistry := registry.NewRegistry(pool)`, used only for on-demand
   REST query handlers — coverage/missing-deps/single-skill review-trigger) — unaffected: this fix
   touches `internal/worker/runner.go` and `internal/registry/review.go` only; the server's
   single-skill-scoped on-demand path is explicitly out of scope and does not need the advisory
   lock (single-row scope, not a whole-table sweep).
3. `Runner.GetMetrics()` (read at `cmd/worker/main.go:132-138`) — unaffected; `runRegistryReview`'s
   return-value change (now surfacing `RunReviewOnce`'s error, if any) threads into the existing
   `r.logger.Error(...)` pattern already used by every other cycle in `runner.go`, no new metrics
   field.
4. `registry.go:142`'s throwaway `skill.NewStore(r.pool).GetCoverage(...)` (used by
   `Registry.GetCoverageReport`) — unaffected, confirmed no overlap with this fix's call graph.
5. Must-not-break: `RunReviewOnce` has zero production callers today (re-confirmed), so extending it
   with the low-coverage stale-mark step is purely additive; the only *new* caller this fix
   introduces is `Runner.runRegistryReview`.
6. This fix's safety is directly gated on G11's advisory-lock scaffolding landing in the same
   window — see §4.7.

### 4.6 Test plan

- Integration: the RED/GREEN stale-marking cycle above, on a real DB.
- Integration: advisory-lock contention — two concurrent `RunReviewOnce` invocations against the
  same DB, asserting only one performs the writes and the other logs a skip-with-reason — proves
  single-owner.
- Unit: table-driven test asserting `RunReviewOnce` and `performReview` now perform the identical 5
  steps (closes the internal asymmetry this pass found).
- Paired §1.1 mutation: revert `runRegistryReview` to call `GetCoverage` only → the stale-marking
  integration test FAILs.
- Regression: existing `runner`/`registry` package tests re-run green.
- §11.4.108 runtime-signature: a fresh, clean-deployed worker process, left running one full
  `ReviewIntervalHours` tick, shows a real `skill_registry.stale` transition for a seeded 31-day-old
  skill.

### 4.7 Impl-sequencing note

This fix and G11's fix (typed `CoverageStats` accessor + panic-safety wrap + the same Postgres
advisory lock) touch the **exact same function** (`runRegistryReview`) and should land as one
coordinated change, in this order: (a) G11's typed-accessor + advisory-lock-helper scaffolding
first — G32's consolidated `RunReviewOnce` call needs the same lock helper to exist; (b) this G32
consolidation (extend `RunReviewOnce`, retarget `runRegistryReview` to call it inside the lock)
second, landed together or in the same PR window. Landing G32 alone without G11's advisory lock
would leave the newly more-powerful, newly write-heavy review cycle exposed to the exact
multi-replica race `g11_worker_design.md` §3 already warned about for this precise function. No
ordering dependency on G29/G31/G35.

---

## §5. Honest Gaps (§11.4.6)

1. **Design-only.** No Go code was written or compiled in this pass; every fix is specified at the
   file:line/seam level but not implemented or `go build`/`go vet`-proven. Implementation is
   explicitly gated behind the in-flight P1.T1 commit (§11.4.84, one Go mutator at a time), per this
   task's own constraint.
2. **G32's consolidation is wider than the register's own minimal fix-direction text.** The register
   says "call `registry.NewHourlyReviewScheduler()` ... from `cmd/worker/main.go`"; this design
   instead consolidates `RunReviewOnce`/`performReview` and retargets `runRegistryReview`, for the
   reasons in §4.2. This is a better-justified but *wider* change than the register anticipated —
   the register itself should be updated to cite this design doc as part of the eventual
   implementation commit.
3. **G29's `pipeline.go:540,557` behavior-change risk is identified, not resolved.** The exact impact
   on `pipeline_test.go`'s existing assertions requires reading that test file's actual bodies at
   implementation time — not done in this pass (out of this pass's read scope, consistent with the
   audit document's own Honest Gap #5, which states test-file assertions were enumerated but not
   read in depth).
4. **G31's `analyzer_test.go` fixture-update requirement (§2.4 item 1) is identified as the riskiest
   item, but the existing test bodies were not read in this pass** (test files were out of this
   design pass's read scope, which focused on the four production-code defects) — flagged as a
   mandatory pre-implementation check, not assumed safe.
5. **No live Postgres/pgvector instance was exercised in this pass.** Every RED-test reproduction
   claim above (e.g., "G32's stale flag stays false today," "G29's semantic match is not surfaced
   today") is derived from reading the code's control flow (e.g., `runRegistryReview` calls only
   `GetCoverage`, which cannot mutate `stale`), not from an actual captured test run. This is a
   static-reasoning claim, stated as such, per §11.4.199's exact-reproduction-sequence discipline —
   the actual RED test must still be authored and run against a real DB at implementation time to
   convert this static claim into captured evidence.
6. **The Postgres advisory-lock key-derivation scheme is specified at the design level only.** This
   doc recommends a stable-string-hash approach over a hand-maintained int64 enum, but the exact
   function (`pg_try_advisory_lock(hashtext(...))` vs a fixed constant table) is left as an
   implementation-time choice, to be cross-checked against G11's own eventual implementation (itself
   impl-pending) — not fully pinned here.
7. **G34** (the bare `rid.(string)` assertions in `middleware.go:184,258-268`) was re-read
   incidentally while confirming G35's evidence but is **not** one of this task's four assigned
   defects and is not designed here.
