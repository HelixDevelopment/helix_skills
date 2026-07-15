# G10 / G27 — Embedding-dimension correctness + async-embed semantics: design-decision document

**Revision:** 1
**Last modified:** 2026-07-15T16:20:00Z
**Status:** design-research, no code landed
**Scope:** the Go backend `github.com/helixdevelopment/skill-system` under `project/`.
**Authority / mandates served:** G10 (`GAPS_AND_RISKS_REGISTER.md:139-149`) — embedding-dimension
correctness; G27 (`GAPS_AND_RISKS_REGISTER.md:287-293`) — `sanitizeTableName` reject-not-strip +
`EmbedAsync` result-channel semantics. Constitution §11.4.6 (no-guessing), §11.4.8/§11.4.150
(deep-research-before-implementation), §11.4.108 (four-layer runtime-signature), §11.4.115
(RED-first), §11.4.186 (cross-doc lockstep), §11.4.197 (research fully completed/wired), §11.4.201
(guard asserts the REAL condition, fail-closed).
**Read discipline:** every Go/SQL/config fact below was read from the committed baseline ref
`255061b` via `git show 255061b:…project/<path>` (the working tree is mid-mutation and was NOT
read). No existing file was modified; this is the single new deliverable.

---

## 0. One-paragraph problem statement

The system stores `vector(768)` columns and declares `dimensions = 768` in config, but **the two
agree only by coincidence, never by construction**: no code templates the DDL from config, no code
asserts at boot that the configured embedder's output width equals the column's declared width, and
the OpenAI embedder — unlike the local one — never checks that the returned vector is the length it
promised. A single config edit (a model whose native width ≠ 768, or one that ignores the
`dimensions` request parameter) produces opaque runtime insert failures deep in a background worker,
not a clear fail-fast at startup. The provider factory supports only `"openai"`/`"local"`, so R7's
mandated HelixLLM / LLMsVerifier / Claude-Toolkit providers error at construction. G27 adds two
embedding-adjacent latent foot-guns: `sanitizeTableName` silently *strips* invalid characters
(masking a programmer error into a wrong-but-valid table name) instead of *rejecting*, and
`EmbedAsync` — while correctly buffered — can silently drop an index (emit neither result nor error)
on a degenerate embedder response, and performs no per-vector width validation.

---

## 1. ROOT CAUSE — proven from `git show 255061b:…`, exact file:line

### 1.1 The `768` is hard-coded in the DDL, un-templated from config

`migrations/001_initial.up.sql:14` (skills table) and `:60` (evidences table):

```
14:    embedding     vector(768),
...
60:    embedding      vector(768),
```

The HNSW indexes are built over those literal columns (`001_initial.up.sql:93-98`). Migrations are
applied verbatim from disk — `internal/db/migrations.go` reads each `NNN_*.up.sql` with
`os.ReadFile` and `tx.Exec`s the raw SQL (`migrations.go:58-66`, `runMigrationSQL` at
`migrations.go:...` executes `tx.Exec(ctx, sql)`); **no substitution of a configured dimension into
the SQL exists anywhere in the migration path.** The `768` is a compile-time literal in a `.sql`
file, not a parameter.

Config default that "agrees by coincidence" — `internal/config/config.go:182-186`:

```
182:		Embedding: EmbeddingConfig{
183:			Provider:   "openai",
184:			Dimensions: 768,
185:			Model:      "text-embedding-3-small",
186:		},
```

The `EmbeddingConfig` type — `internal/config/config.go:105-112`:

```
105:// EmbeddingConfig selects the embedding provider and model.
106:type EmbeddingConfig struct {
107:	Provider      string `toml:"provider"`       // "openai" | "local"
108:	Dimensions    int    `toml:"dimensions"`     // e.g. 768
109:	Model         string `toml:"model"`          // e.g. "text-embedding-3-small"
110:	APIKey        string `toml:"api_key"`        // OpenAI API key (env override recommended)
111:	LocalEndpoint string `toml:"local_endpoint"` // URL for local model server
112:}
```

### 1.2 No model↔column dimension assertion exists (danger-zone #3 unmet)

`config.validate` checks only that the number is positive — `internal/config/config.go:461-463`:

```
461:	if cfg.Embedding.Dimensions <= 0 {
462:		issues = append(issues, fmt.Sprintf("invalid embedding.dimensions: %d", cfg.Embedding.Dimensions))
463:	}
```

There is **no** query of `information_schema` / `pg_attribute` for the `embedding` column's declared
dimension, and **no** comparison of `embedder.Dimensions()` against it, at any startup path. The
server boot sequence runs migrations and moves on without touching the embedder at all —
`cmd/server/main.go:84-90`:

```
84:	// 4. Run migrations
85:	ctx := context.Background()
86:	if err := db.Migrate(ctx, pool, "./migrations"); err != nil {
87:		logger.Warn("Migration failed", zap.Error(err))
88:	} else {
89:		logger.Info("Migrations completed")
90:	}
```

(Note also: migration failure is a **`logger.Warn` then continue** — a separate fail-open flagged as
G23; the dimension assertion designed below MUST be a hard fail, not a warn.)

**Load-bearing fact for the assertion's placement:** `NewEmbedderFromConfig` is **never called** in
`cmd/server/main.go` **nor** `cmd/worker/main.go` in the baseline (a whole-tree grep of all 61 Go
files at `255061b` for `NewEmbedderFromConfig` / `embedder.Embed` finds call sites ONLY inside
`internal/db/embedding.go` itself, and the only `db.Embedder` consumer is
`autoexpand.NewPipeline(store, embedder, …)` at `internal/autoexpand/pipeline.go:64`). The worker
`main` builds only the store, never the embedder — this exact wiring gap is independently recorded
in `research/g11_worker_design.md:150` and `:199`. **Consequence:** the boot-time assertion cannot
be bolted onto an existing embedder-construction site because there is none yet; it must be added
*together with* the embedder-construction wiring (a small but real new wiring item, tracked here so
it is not lost per §11.4.197).

### 1.3 OpenAI returned-vector length is never checked

`OpenAIEmbedder.Embed` validates the *count* of returned rows and each row's *index* but never the
*length* of any returned embedding — `internal/db/embedding.go:124-143`:

```
124:	if len(result.Data) != len(texts) {
125:		return nil, fmt.Errorf("OpenAI returned %d embeddings for %d inputs", len(result.Data), len(texts))
126:	}
127:
128:	// Sort by index to maintain input order.
129:	vectors := make([][]float32, len(texts))
130:	for _, d := range result.Data {
131:		if d.Index < 0 || d.Index >= len(texts) {
132:			return nil, fmt.Errorf("OpenAI returned invalid embedding index %d", d.Index)
133:		}
134:		vec := make([]float32, len(d.Embedding))
135:		for i, v := range d.Embedding {
136:			vec[i] = float32(v)
137:		}
138:		vectors[d.Index] = vec        // <-- len(d.Embedding) never compared to e.dimensions
139:	}
```

The code *requests* a width when configured — `internal/db/embedding.go:86-88`:

```
86:	if e.dimensions > 0 {
87:		reqBody.Dimensions = e.dimensions
88:	}
```

but never verifies the width it got back. Contrast the **local** embedder, which does exactly the
missing check — `internal/db/embedding.go:256-262`:

```
256:	// Validate dimensions.
257:	for i, vec := range result.Embeddings {
258:		if len(vec) != e.dimensions {
259:			return nil, fmt.Errorf("local embedder returned vector length %d at index %d, expected %d",
260:				len(vec), i, e.dimensions)
261:		}
262:	}
```

A model that returns 1536 (e.g. `text-embedding-ada-002`, which ignores the `dimensions` request
parameter) hands a 1536-float slice straight to a `vector(768)` insert → an opaque runtime error at
insert time, far from the config edit that caused it.

### 1.4 Only `"openai"`/`"local"` providers exist; R7 providers error at startup

`NewEmbedderFromConfig` — `internal/db/embedding.go:293-309`:

```
293:func NewEmbedderFromConfig(cfg config.EmbeddingConfig) (Embedder, error) {
294:	switch cfg.Provider {
295:	case "openai":
...
299:		return NewOpenAIEmbedder(cfg), nil
300:	case "local":
...
304:		return NewLocalEmbedder(cfg), nil
305:	default:
306:		return nil, fmt.Errorf("unsupported embedding provider: %q (expected "+
307:			"\"openai\" or \"local\")", cfg.Provider)
308:	}
309:}
```

But `SPEC.md:395` documents a third value in the sample config — `provider = "openai"  -- openai |
local | anthropic` — and R7 (`REQUIREMENTS.md:68-69`) mandates: *"Obtain quality models via
LLMsVerifier, HelixLLM, Claude Toolkit aliases — pluggable ModelProvider, not hardcoded OpenAI."*
Any of those provider strings hits the `default` and errors at construction.

### 1.5 `EmbedAsync` channel semantics (G27, embedding-adjacent)

`EmbedAsync` — `internal/db/embedding.go:359-405`. The **buffering IS correct** (the register's
verdict at `GAPS_AND_RISKS_REGISTER.md:290` — "correct (buffered to `len(texts)`)" — is confirmed):
`results := make(chan AsyncEmbedResult, len(texts))` (`embedding.go:369`) + `defer close(results)`
(`embedding.go:372`) guarantees no goroutine blocks on send and `close` is race-safe under
`wg.Wait()`.

Two **real, factual** latent semantics gaps remain (an honest refinement of the register's
"correct" note, not a contradiction of it):

**(a) Silent index-drop on an empty result** — `internal/db/embedding.go:395-397`:

```
395:				if len(vecs) > 0 {
396:					results <- AsyncEmbedResult{Index: idx, Vector: vecs[0]}
397:				}
```

If `embedder.Embed` returns `nil, nil` (a degenerate but reachable response — the interface permits
`([][]float32, error)` and every implementation early-returns `nil, nil` on empty input,
`embedding.go:70-72`, `:210-212`), the `else` branch emits **neither a result nor an error**. A
caller that ranges to close expecting one entry per input index gets a silently-missing index. This
is a fail-open on a degenerate embed response.

**(b) No per-vector width validation** — `EmbedAsync` (and `EmbedBatch`,
`internal/db/embedding.go:317-345`) pass whatever `Embed` returns straight through; with the §1.3
OpenAI gap unfixed, a wrong-width vector flows through the async path into the DB insert with no
guard. Fixing §1.3 at the embedder layer closes this by construction (the width check lives inside
`Embed`), which is why §3.1's assertion is placed at the provider layer, not duplicated per call
site.

Cancellation is otherwise sound: on `ctx.Done()` each not-yet-dispatched index emits exactly one
`{Index:i, Error: ctx.Err()}` and `continue`s (`embedding.go:378-383`), so no new goroutine is
spawned after cancel and the buffer (size `len(texts)`) never overflows.

### 1.6 `sanitizeTableName` strips instead of rejects (G27)

`internal/db/vector.go:286-296`:

```
286:// sanitizeTableName ensures the table name contains only alphanumeric
287:// characters and underscores to prevent SQL injection.
288:func sanitizeTableName(name string) string {
289:	var b strings.Builder
290:	for _, r := range name {
291:		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
292:			b.WriteRune(r)
293:		}
294:	}
295:	return b.String()
296:}
```

`"skills; DROP"` → `"skillsDROP"` (a valid-looking but wrong identifier). Callers guard only against
the empty string (e.g. `VectorSearch` at `vector.go:44-47`: `if sanitized == "" { return … invalid
table name }`), so a stripped-to-non-empty name passes the guard and queries the wrong table.
**Safe today** because every caller passes an internal constant (`FindSimilarSkills` →`"skills"`,
`vector.go:216`; `FindSimilarEvidences` →`"evidences"`, `vector.go:226`), so this is a **latent**
foot-gun, not an active injection (matching the register's `GAPS_AND_RISKS_REGISTER.md:290-291`
severity: low). A future dynamic/user-influenced caller would hit a wrong-but-valid name with no
error. This is a §11.4.201 violation-in-waiting: the guard does not assert the REAL condition
("is this a known table?") — it *transforms* the input into something that passes.

---

## 2. DECISION — the design (all decisions fail-CLOSED per §11.4.201)

### 2.1 A model→dimension registry (source of truth for expected width)

Add a small, declarative, in-code registry mapping `(provider, model)` → native/expected output
dimension, plus the rule for width-overridable models:

```
Registry entry: { provider, model, native_dim, dim_overridable bool, notes }
```

- `dim_overridable = true` (e.g. OpenAI `text-embedding-3-small` / `-3-large`, which honour the
  `dimensions` request parameter via Matryoshka truncation): the *effective* expected width is
  `min(config.dimensions, native_dim)` when `config.dimensions > 0`, else `native_dim`.
- `dim_overridable = false` (e.g. `text-embedding-ada-002`): the effective expected width is
  `native_dim`, and if `config.dimensions` is set to anything ≠ `native_dim` the config is REJECTED
  at validate time with a precise message (this catches the exact 768-vs-1536 trap in §1.3 before a
  single request is made).
- Unknown `(provider, model)` pairs: the registry returns "unknown, expected = config.dimensions"
  and the effective expected width falls back to `config.dimensions` — the config value is trusted
  ONLY as the last resort, and the fact that it was unverifiable from the registry is logged
  honestly (never silently trusted, §11.4.6). The boot-time column assertion (§2.2) still applies,
  so an unknown-model config that disagrees with the column is still caught.

This registry is consumer-owned DATA in the `db` package (not project-specific literals leaking into
a shared submodule — §11.4.28 N/A here, it is app-local). Concrete per-provider `native_dim` values
that require external verification are enumerated in §5 (Honest gaps) and MUST be filled from
latest official sources per §11.4.99, not from memory.

### 2.2 Boot-time (and migration-time) fail-CLOSED column↔model assertion

Add `AssertEmbeddingDimension(ctx, pool, table, expectedDim) error` in the `db` package that:

1. Resolves the column's **declared** dimension from the authoritative catalog. Primary method
   (portable, unambiguous): parse `format_type(atttypid, atttypmod)` for the column, which pgvector
   renders as `vector(768)` →parse the integer. Query:
   ```sql
   SELECT format_type(a.atttypid, a.atttypmod)
   FROM pg_attribute a
   WHERE a.attrelid = $1::regclass AND a.attname = 'embedding' AND NOT a.attisdropped;
   ```
   Secondary/faster method (confirm on live DB before relying on it, §5): read `atttypmod` directly.
2. Compares the parsed declared dimension to `expectedDim` (from §2.1: the effective expected width
   for the configured `(provider, model, dimensions)`).
3. On mismatch: returns an error that **fails the process** (hard stop, NOT a `logger.Warn`) — the
   §11.4.201 conservative-safe default. The error names both values and the resolved source
   (`"embedding column skills.embedding is vector(1536) but configured embedder text-embedding-3-small
   yields 768 — refusing to start"`). On a column-not-found / unreadable catalog: fail-closed with an
   honest "could not resolve column dimension" message (an unresolvable signal refuses, it does not
   assume-OK — §11.4.201(4)).

**Wiring (the §11.4.197 completeness item):** the assertion is called (a) at server boot in
`cmd/server/main.go` **after** `db.Migrate` and after constructing the embedder via
`NewEmbedderFromConfig`, and (b) at worker boot in `cmd/worker/main.go`, which ALSO gains the
embedder construction it currently lacks (§1.2). Both `skills` and `evidences` columns are asserted.
A separate migration-time variant runs the same assertion at the end of `Migrate` for the
just-created columns so a mismatch surfaces during `migrate up`, not only at serve time (this also
neutralises the §1.2 "migration failure only warns" interaction for the dimension case).

### 2.3 Support non-768 providers: parameterize the column dimension from config

- **Ship 768 as the default** (matches SPEC §12, the current columns, and the research "sweet spot")
  — decision inherited from the register `GAPS_AND_RISKS_REGISTER.md:148`.
- **Make `vector(N)` a migration parameter driven by `config.embedding.dimensions`.** Two viable
  mechanisms, decision = **(A)** for the MVP:
  - **(A) Templated migration (recommended):** the migration runner substitutes `{{.EmbeddingDim}}`
    (or an equivalent placeholder) in the `vector(…)` column DDL from `config.embedding.dimensions`
    at apply time, so the column is created at the configured width. This requires teaching
    `internal/db/migrations.go` a minimal, explicit template step for the initial schema ONLY (the
    `vector(N)` sites), keeping every other statement byte-verbatim. Cheapest correct path; no live
    re-index.
  - **(B) Documented supported-set + guard (fallback if templating is deemed too invasive for the
    first migration):** keep 768 literal, and REJECT at config-validate any `dimensions` value not
    in a documented supported set `{768}` for the MVP, with a clear "changing embedding width
    requires a schema migration" message. This is strictly weaker (no live non-768 support) but is
    honest and fail-closed; it is the explicit fallback, not the default.
- **Rejected alternatives** (from the register `GAPS_AND_RISKS_REGISTER.md:148`): (a) support all
  dims live via runtime re-index — expensive, out of MVP scope; (b) leave 768 hard-coded and
  undocumented — reproduces the drift.

### 2.4 Extend the provider factory for R7 providers (HelixLLM / OpenAI-compatible)

Extend `NewEmbedderFromConfig` (`embedding.go:293`) with an OpenAI-compatible provider branch that
covers HelixLLM / LLMsVerifier / Claude-Toolkit-fronted OpenAI-compatible embedding endpoints
(base-URL + auth injection reusing the `OpenAIEmbedder` transport, which already parameterizes
`baseURL`). Each new provider MUST register its `(provider, model)`→dim rows in the §2.1 registry so
the boot assertion covers it. Unknown provider strings still fail-closed at the `default` (that
error is correct — an unconfigured provider must not silently pick one). This closes the R7 gap
(`REQUIREMENTS.md:68-69`) and the `SPEC.md:395` `anthropic` mention. Honest boundary (§5): the exact
HelixLLM embedding wire contract is unvendored (G14) and UNCONFIRMED offline; the branch is designed
against the OpenAI-compatible shape and the real contract must be verified when the submodule
resolves.

### 2.5 OpenAI returned-length check (mirror the local embedder)

Add, immediately after building each `vec` in `OpenAIEmbedder.Embed` (at `embedding.go:134-138`), the
same guard the local embedder already has (`embedding.go:257-261`): if `len(d.Embedding) !=
e.dimensions` return a descriptive error naming got/expected/index. Fail-closed: a wrong-width vector
is an error, never stored. This makes §1.3's silent path impossible and, combined with §2.2, means a
mismatch is caught at *both* boot (column vs config) and per-request (response vs config).

### 2.6 `sanitizeTableName` → reject, not strip (§11.4.201)

Replace the stripping loop (`vector.go:288-296`) with a **validate-and-reject** function:
`validateTableName(name) (string, error)` returns the name unchanged iff it fully matches
`^[A-Za-z_][A-Za-z0-9_]*$` (leading-digit also rejected), else returns an error. Callers stop
treating `""` as the only failure and propagate the typed error. Additionally, since the only real
callers pass one of two constants, add a **closed allowlist enum** (`TableSkills`, `TableEvidences`)
and have `VectorSearch`/`HybridSearch`/`VectorSearchFiltered`/`VectorIndexStats`/
`WaitForVectorIndexReady` accept that enum type rather than a free `string`. The reject function is
the defence-in-depth guard; the enum makes a wrong table a compile-time impossibility. This is the
§11.4.201 fix: the guard now asserts the REAL condition (a valid, known identifier) and refuses on
violation instead of transforming bad input into passable input.

### 2.7 `EmbedAsync` result-completeness + no fail-open (§11.4.201)

- **Every input index emits exactly one `AsyncEmbedResult`.** Replace the `if len(vecs) > 0` branch
  (`embedding.go:395-397`) so the empty-result case emits an explicit error result
  (`AsyncEmbedResult{Index: idx, Error: errEmptyEmbedding}`), never nothing. Callers ranging to close
  get one entry per index — a guaranteed-complete result set (the fail-closed contract).
- **Keep the correct buffering** (`make(chan …, len(texts))` + `defer close`) — do not regress it.
- **Per-vector width** is enforced upstream in `Embed` (§2.5), so `EmbedAsync` needs no duplicate
  width check; document in the function comment that results carry validated-width vectors.
- Document the completeness contract in the `EmbedAsync` doc-comment (the register asked for this as
  a "caller-contract reminder", `GAPS_AND_RISKS_REGISTER.md:290`).

---

## 3. WHY (§11.4.8 external precedent) + RUNTIME SIGNATURES (§11.4.108)

### 3.1 Why — external precedent for pgvector dimension handling

- **pgvector's own contract:** a `vector(N)` column enforces N at INSERT — pgvector raises
  `ERROR: expected N dimensions, not M` on a width mismatch. The mature-stack precedent is therefore
  to **fail fast at boot** rather than discover this per-insert deep in a worker. (pgvector README /
  docs — the exact error string + whether `atttypmod` stores the raw dimension MUST be confirmed on
  the live `pgvector/pgvector:pg16` container per §11.4.99; marked UNCONFIRMED in §5.)
- **OpenAI Matryoshka (`dimensions` param):** `text-embedding-3-*` support requesting a shortened
  width; `ada-002` does not. This is the precedent for the `dim_overridable` flag in §2.1. Exact
  native widths MUST be verified from current OpenAI docs (§5), not asserted from memory.
- **Templated-DDL-from-config** is the standard way projects avoid the "column width is a magic
  literal" drift (the same class §11.4.186 forbids for docs).
- The **local embedder in this very codebase** (`embedding.go:256-262`) is the internal precedent for
  the missing OpenAI check — the fix is "make OpenAI do what local already does". No external source
  needed for that half: **"NO external solution found — original work"** (it is a one-line internal
  symmetry fix).

### 3.2 Runtime signatures — one machine-checkable observable per fix, on a clean deploy

| Fix | Runtime signature (asserted on a fresh `pgvector:pg16` deploy) |
|-----|----------------------------------------------------------------|
| §2.2 boot assertion (match) | Server/worker boots; log line `embedding dimension verified: skills.embedding=vector(768) == embedder(768)`; process healthy. |
| §2.2 boot assertion (mismatch → fail-closed) | With column `vector(1536)` + config `768`, the process **exits non-zero** at boot with `refusing to start: embedding column … vector(1536) but embedder yields 768`; `/health` never comes up. (RED-polarity: this is the captured proof the guard fires.) |
| §2.3 templated column | `SELECT format_type(atttypid,atttypmod) … = 'vector(768)'` on a fresh migrate; with `dimensions=1024` config + templating, the created column is `vector(1024)` and boot asserts equal. |
| §2.4 R7 provider | `provider = "helixllm"` (or OpenAI-compatible) constructs a non-nil `Embedder` with `Dimensions()==config.dimensions` (no `default`-branch error); unknown provider still errors. |
| §2.5 OpenAI length check | A stubbed OpenAI response of width 1536 with config 768 returns a `vector length 1536 … expected 768` error from `Embed` — the vector is never stored. |
| §2.6 reject-not-strip | `validateTableName("skills; DROP")` returns a non-nil error (not `"skillsDROP"`); `VectorSearch` with a bad table returns the typed error; the enum callers compile only with `TableSkills`/`TableEvidences`. |
| §2.7 EmbedAsync completeness | For N inputs, exactly N `AsyncEmbedResult`s are received before channel close, even when one embed returns empty (that index carries a non-nil `Error`). `go test -race` clean. |

---

## 4. TEST-CASE COUNT + reconciliation with `research/testing_infrastructure_plan.md`

### 4.1 Enumerated cases (RED-first per §11.4.115; each PASS cites captured evidence per §11.4.5/§11.4.69)

**Unit (no DB) — 12**
1. Registry: known overridable model → effective width = `min(config,native)`.
2. Registry: known non-overridable model + `config≠native` → config REJECTED at validate.
3. Registry: unknown `(provider,model)` → falls back to `config.dimensions` + honest "unverified" log flag.
4. `config.validate` rejects `dimensions` not in supported set under fallback design (§2.3B) OR accepts any positive under templated design (§2.3A) — one case per chosen path.
5. OpenAI `Embed` rejects a returned vector whose `len != e.dimensions` (stub transport, width 1536 vs 768). **[RED on baseline: currently PASSES-through — proves the gap]**
6. OpenAI `Embed` accepts a correct-width response.
7. Provider factory returns a real embedder for `openai`.
8. Provider factory returns a real embedder for `local`.
9. Provider factory returns a real embedder for the new R7/OpenAI-compatible provider (`helixllm`). **[RED on baseline: hits `default`]**
10. Provider factory still errors on a genuinely-unknown provider string (fail-closed preserved).
11. `validateTableName` rejects `"skills; DROP"`, `"1skills"`, `""`, `"a-b"`; accepts `"skills"`, `"evidences"`, `"a_b1"`. **[RED on baseline: `sanitizeTableName` returns `"skillsDROP"`]**
12. `EmbedAsync` emits exactly N results for N inputs incl. one empty-embed index (carrying an error), channel closes, `-race` clean. **[RED on baseline: empty index emits nothing]**

**Integration (real `pgvector:pg16` from `deploy/`) — 6**
13. Boot assertion PASSES when column dim == embedder dim (log signature present).
14. Boot assertion FAILS-CLOSED (process exit ≠ 0) when column `vector(1536)` vs config `768`. **[the §11.4.115 RED→GREEN artifact]**
15. Templated migration (§2.3A) creates `vector(N)` at the configured N; `format_type` == `vector(N)`.
16. A correct-width vector inserts and round-trips via `StoreSkillEmbedding` + `VectorSearch`.
17. A wrong-width vector insert is rejected by pgvector (captures the native error the boot-assertion pre-empts).
18. Migration-time assertion variant fails `migrate up` on a seeded mismatch.

**Contract — 1**
19. Config schema documents `provider ∈ {openai, local, helixllm,…}` + `dimensions`, and the `SPEC.md:395` sample parses/aligns (folds into the G09/G19 doc-lint gate).

**Paired §1.1 mutations — 4** (each MUST flip its guard RED then restore GREEN)
- M1: change the column to `vector(1536)` keeping config `768` → boot-assertion integration test (14) FAILs.
- M2: delete the OpenAI length check → unit test (5) FAILs.
- M3: revert `validateTableName` to the stripping loop → unit test (11) FAILs.
- M4: restore `if len(vecs) > 0` (drop the empty-emit) → `EmbedAsync` completeness test (12) FAILs.

**Total: 23 cases** (12 unit + 6 integration + 1 contract + 4 paired mutations).

### 4.2 Reconciliation with `research/testing_infrastructure_plan.md` (§11.4.186)

- **Coverage-matrix row G10** — `testing_infrastructure_plan.md:296`: *"unit (provider factory incl.
  helixllm; OpenAI length-mismatch rejected), integration (startup FAILs on dim mismatch; correct dim
  inserts), contract (config schema), mutation (column→1536 keep config 768 → startup assertion
  FAILs)"* with evidence *"startup dim-assertion log (`pg_attribute` query); mismatch-rejected error;
  provider-factory table"*. **This design fully realises that row**: cases 5/7/9 (unit factory +
  length), 13/14/16 (integration startup + inserts), 19 (contract), M1 (mutation) map 1:1;
  the plan's `pg_attribute`/startup-log evidence is exactly the §3.2 signature.
- **Coverage-matrix row G27** — `testing_infrastructure_plan.md:313`: *"unit (invalid table name
  rejected), security, regression"* with evidence *"rejection error for `\"skills; DROP\"`"*. Cases
  11 + M3 realise it. **Extension flagged per §11.4.186:** the plan's G27 row covers only the
  `sanitizeTableName` half. The `EmbedAsync` completeness half (cases 12 + M4) is an **extension**
  of the G27 row this design adds — the register itself scopes `EmbedAsync` as a caller-contract note
  (`GAPS_AND_RISKS_REGISTER.md:290`), so the plan's row does not yet enumerate it; recommend the G27
  row be extended to `"…; EmbedAsync emits one result per index incl. empty-embed error; mutation:
  drop empty-emit → FAILs"`.
- **Harness fit (no fork, §11.4.27):** all unit cases are pure stdlib `testing` table-driven with a
  fake transport (the sanctioned unit-level double, `testing_infrastructure_plan.md:190`); all
  integration cases run against the real `deploy/` `pgvector:pg16` with `//go:build integration` +
  `HELIX_TEST_DATABASE_URL` SKIP-with-reason (`testing_infrastructure_plan.md:88-98`); the four
  mutations register into `scripts/mutation/` per the §1.3 harness
  (`testing_infrastructure_plan.md:129-145`). No new test dependency is introduced.
- **Sequencing:** §2.3A templated-migration touches the initial schema and therefore composes with
  P1.T1 (kind-column / relation-type schema change, `research/p1t1_granularity_schema_migration.md`);
  the dimension parameterization MUST land **in lockstep** with that migration so the `vector(N)`
  template and the P1.T1 DDL edits do not conflict (§11.4.186 SPEC+migration lockstep — same
  discipline the G06/G07 doc applied). P1.T1 does not currently touch embedding DDL (grep of that doc
  for `embed`/`dimension` = 0 hits), so this is an additive, non-conflicting extension of the initial
  migration, flagged here so it is not lost.

---

## 5. HONEST GAPS (§11.4.6) — UNCONFIRMED / live-required, never guessed

- **UNCONFIRMED (live-DB required):** whether pgvector's `pg_attribute.atttypmod` stores the raw
  dimension (e.g. `768`) or an offset-encoded value. The §2.2 primary method uses
  `format_type(atttypid, atttypmod)` string-parse (`vector(768)`), which is robust regardless; the
  faster raw-`atttypmod` secondary path MUST be confirmed on the real `pgvector/pgvector:pg16`
  container before it is relied upon. Not guessed.
- **UNCONFIRMED (external, §11.4.99 latest-source required):** exact native output widths per
  provider/model (OpenAI `text-embedding-3-small` native 1536 / `-3-large` native 3072 /
  `ada-002` 1536-fixed-no-override) and which models honour the `dimensions` request parameter. These
  populate the §2.1 registry `native_dim` / `dim_overridable` columns and MUST be filled from current
  official docs, not from memory. Marked `PENDING_FORENSICS` until verified.
- **UNCONFIRMED (unvendored, G14):** the exact HelixLLM / LLMsVerifier / Claude-Toolkit embedding
  wire contract (endpoint shape, auth header, response envelope, whether it exposes a `dimensions`
  parameter). §2.4 is designed against the OpenAI-compatible shape; the real contract binds when the
  `embeddings`/`helix_llm` submodules resolve (`testing_infrastructure_plan.md:46-48`,
  `helix-deps.yaml`). Not fabricated.
- **Live-DB required for GREEN:** integration cases 13–18 (the fail-closed boot assertion, the
  templated-column width, the real insert round-trip) need the `deploy/` compose up; offline they are
  authored RED-first + SKIP-with-reason (`llm`/`db` unreachable), never a green PASS
  (`testing_infrastructure_plan.md:88-98`).
- **Wiring dependency (tracked, §11.4.197):** the boot assertion has no existing embedder-construction
  site to attach to — `NewEmbedderFromConfig` is called nowhere in `cmd/*` at `255061b`
  (whole-tree grep). Landing §2.2 REQUIRES adding embedder construction to `cmd/server/main.go` and
  `cmd/worker/main.go` (the latter is the same gap `research/g11_worker_design.md:150,199` records).
  Flagged so it is not lost as un-wired research.
- **Not in scope / not claimed:** live runtime re-index to change an already-populated column's width
  (register rejected-alternative (a), `GAPS_AND_RISKS_REGISTER.md:148`); a real per-provider live
  quota/latency proof (needs reachable providers). Neither is designed here; both are honest
  deferrals, not silent omissions.

*Positive-evidence-only. Every "is/does" claim above is pinned to a file:line read from
`git show 255061b:…`; every dimension/provider fact I could not verify offline is marked UNCONFIRMED
with its unblock condition. No forbidden hedging vocabulary used.*
