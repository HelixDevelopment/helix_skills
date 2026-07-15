# R7/R19 Credential-Loading Layer Design — `~/api_keys.sh` + `.env` dual-source provider-key loader

**Revision:** 1
**Last modified:** 2026-07-15T18:32:32Z

> **Operator mandate (2026-07-15, verbatim):** "All API tokens can be taken /
> loaded from `api_keys.sh` located in home directory of the host or local
> `.env` file! Both MUST BE fully supported like other Helix projects do!"

**Classification:** design/research only. No Go code written, nothing under
`project/` modified, no git operations. This document specifies the layer; a
separate landed work-item implements it.

---

## 1. Goal

Give the HelixKnowledge Skill Graph System (`github.com/helixdevelopment/skill-system`)
a credential-loading layer that resolves provider API tokens from **both**
`$HOME/api_keys.sh` **and** a project-local `.env`, both fully supported, with a
**defined, deterministic precedence** (§11.4.6), mirroring the established
Helix-family pattern (`helix_code/internal/secrets/loader.go`). The resolved
tokens must reach:

- the existing TOML config layer's `${VAR}` interpolation
  (`project/internal/config/config.go:343-365`), so `config.toml` can carry
  `llm_api_key = "${ANTHROPIC_API_KEY}"` (R19 §2.5) and `api_key = "${OPENAI_API_KEY}"`
  (`config.go:110`, `EmbeddingConfig.APIKey`); and
- the LLMProvider layer — R19's interim thin `net/http` Anthropic client
  (`NewAnthropicLLM(apiKey, model string, ...)`) **and** R22's adopted
  `HelixDevelopment/LLMProvider` submodule (`pkg/providers/anthropic.NewProvider(apiKey, baseURL, model string)`,
  `pkg/providers/openai.NewProvider(...)`), both of which take the key as a
  **constructor string argument** and never read the environment themselves
  (verified: `r22_full_catalogue_incorporation_design.md:77-82`). The loader's
  job is therefore to land the value into the process environment **before**
  `config.Load()` runs, so config interpolation → `cfg.AutoExpand.LLMAPIKey` /
  `cfg.Embedding.APIKey` → the constructor argument.

Hard invariant: variable **NAMES** may appear in tracked files; secret **VALUES**
never do (§11.4.10 / §11.4.30). `.env` / `.env.*` / `*.env` and a real
`api_keys.sh` are git-ignored; only `*.example` placeholders are tracked.

---

## 2. Established Helix pattern (cited)

The canonical reference implementation is **`helix_code`** — the only sibling
that resolves provider keys from these two sources in Go. Cited files (all read
directly; line numbers verified):

### 2.1 The Go in-process loader — `helix_code/helix_code/internal/secrets/loader.go`

`LoadAPIKeys()` (`loader.go:30-48`) is the authoritative pattern:

- **Source order:** `$HOME/api_keys.sh` first (`loader.go:31-37`) — if it exists,
  parse it and **return**; otherwise walk up from cwd to find a `.env`
  (`loader.go:39-47`, `findEnvFile()` at `:133-143`). Returns an error only when
  **neither** source is found (`:45`).
- **Parse, never source/execute.** `loadFromShell` (`:54-81`) is a line parser
  for `export VAR=value` lines; `loadFromEnv` (`:84-106`) parses `VAR=value`. It
  strips single/double quotes (`stripQuotes`, `:122-131`), skips comments and
  blanks, and in the shell file **skips any line without the `export ` prefix**
  (`:67-69`) so the two formats stay distinct. No shell is ever executed.
- **Gap-fill precedence.** `setIfAbsent` (`:115-120`): a var already present and
  **non-empty** in the process environment is never overwritten — "the file only
  fills gaps" (`:108-111`). A present-but-empty var is treated as a gap.
- **No secret ever logged** (`loader.go:29`, their CONST-042 ≈ our §11.4.10).
- **Idempotent**, applies via `os.Setenv`.

### 2.2 Startup wiring — `helix_code/helix_code/cmd/server/main.go`

`loadAPIKeysAtStartup()` (`main.go:81-83`) = `secrets.LoadAPIKeys() == nil`, and
`main()` calls it (`main.go:99`) **before** `config.Get()` (`:105`) so "a key
supplied only via those files becomes visible to config … AND to the … funnel's
key-presence gate" (`main.go:70-76`). A missing source is **non-fatal**
(`:77`). The boolean return exists purely so a test can assert the wiring is
live — an explicit anti-bluff hook (`:78-80`). The CLI path
(`cmd/cli/main.go:3090-3136`) is identical.

### 2.3 Two coexisting env-var naming conventions

- **Canonical:** `<UPPER_SNAKE>_API_KEY` (e.g. `ANTHROPIC_API_KEY`,
  `OPENAI_API_KEY`) — what Go provider constructors read.
- **Legacy `~/api_keys.sh` convention:** `ApiKey_<TitleCase>` (e.g.
  `ApiKey_Anthropic`, `ApiKey_OpenAI`).
- `DeriveKeyEnvAliases(name)` (`internal/llm/verifier_dynamic_catalogue.go:56-73`)
  derives, from a provider **name** alone, the ordered alias list
  `[<UPPER>_API_KEY, ApiKey_<TitleCase>, …catalogue aliases]`.
- The shell wrapper `scripts/load_api_keys.sh:80-121`
  (`helixcode_normalise_api_keys`) translates `ApiKey_<Provider>` →
  `<PROVIDER>_API_KEY`, gap-fill only (canonical wins if already set,
  `:114-118`), across a 29-entry provider table (`:82-110`). Its comment
  (`:68-75`) names the exact failure this prevents: an operator whose
  `api_keys.sh` exports only `ApiKey_Anthropic` would leave `ANTHROPIC_API_KEY`
  unresolved and the feature silently keyless — "a readiness bluff".

### 2.4 The source-able shell wrapper — `helix_code/scripts/load_api_keys.sh`

A **separate** convenience for shell entrypoints and test harnesses.
`helixcode_load_api_keys` (`:16-66`) prefers `$HOME/api_keys.sh` (real `source`,
`:31-36`) else walks up for `.env` with `set -a` auto-export (`:41-61`). It runs
in a real shell, so it **does** execute the file; it defends against a caller's
`set -u` aborting on unbound `${…}` references by toggling `set +u` around the
sourcing (`:26-34`). Auto-runs when sourced unless `HELIXCODE_LOAD_API_KEYS=0`
(`:124-127`). Project-agnostic per their CONST-051(B) (`:3-7`).

### 2.5 Test shape — `helix_code/helix_code/internal/secrets/loader_test.go`

Hermetic, no real secrets. `withIsolatedEnv` (`:12-47`) points `HOME` at a
`t.TempDir()` and `chdir`s to another, saving/restoring env. `writeFile` uses
`0o600` (`:49-54`). Cases: from-shell (`:56`), from-env (`:71`),
prefers-shell-over-env (`:86`), strips-quotes (`:102`), ignores-comments
(`:121`), ignores-blank (`:137`), skips-non-`export` lines (`:153`),
neither-found-returns-error (`:174`). Fixture values are obviously fake
(`bar`, `from_sh`, `from_env`) — no secret VALUE ever enters a fixture.

**Cross-family corroboration (names only):** `helix_ota`, `helix_translate`,
`helix_terminator`, `helix_vpn` all vendor the `vasic-digital/LLMProvider` copy
(`r22_…:57-59`); `helix_code` and this project resolve the
`HelixDevelopment/LLMProvider` copy. The `~/api_keys.sh` + `.env` loader itself
was located only in `helix_code` (grep of all siblings, §research). Where a
sibling other than `helix_code` implements it differently I did **not** find a
divergent Go loader — `helix_code` is the sole in-Go reference and is treated as
canonical here.

---

## 3. Sources & precedence (DEFINED — §11.4.6)

Three sources can supply a variable. The precedence, highest wins:

```
1. process environment (a real shell `export` / launcher-set var)   ← HIGHEST
2. $HOME/api_keys.sh   (host-home shell file, `export VAR=value`)
3. <repo>/.env         (project-local dotenv, VAR=value)            ← LOWEST
```

**Rationale.** (1) A var an operator/launcher already exported is an explicit,
in-the-moment intent and must never be silently overwritten by a file — this is
the sibling's `setIfAbsent` gap-fill rule (`loader.go:108-120`) and the least
surprising. (2) `api_keys.sh` is the operator's host-wide credential vault; a
project `.env` is a per-checkout override for local dev, so the host vault
outranks the checkout when both name the same var. This matches the sibling's
"shell file checked first" ordering and its
`TestLoadAPIKeys_PrefersShellOverEnv`.

### 3.1 Exclusive vs. union — the one deliberate refinement (decision point)

The sibling loader is **exclusive**: if `api_keys.sh` exists it returns without
ever reading `.env` (`loader.go:35`). That means a var present **only** in
`.env` is invisible whenever `api_keys.sh` exists — arguably *not* "both fully
supported" for that var.

- **Option A — sibling-exact (exclusive):** `api_keys.sh` present ⇒ `.env`
  ignored entirely. Simplest, byte-for-byte matches `helix_code`.
- **Option B — union with the same precedence (RECOMMENDED):** read `api_keys.sh`
  first (gap-fills over the process env), then **also** read `.env` (gap-fills
  what remains). Identical result to Option A whenever both files define the
  same var; differs only when `.env` supplies a var `api_keys.sh` lacks — which
  Option A drops. Option B is a strict superset and is the most literal reading
  of "**both** MUST BE fully supported".

**Recommendation:** adopt **Option B**, because the mandate's emphasis is "both
… fully supported"; flag the divergence from the literal sibling behaviour in
the commit and this doc (see §11 Open questions). Precedence is unchanged and
fully deterministic either way.

---

## 4. `api_keys.sh` consumption — security analysis (parse, do not source)

`~/api_keys.sh` is a POSIX shell file of `export KEY=VALUE` lines living in the
host home directory. Two ways a Go process could consume it:

| Approach | Mechanism | Verdict |
|---|---|---|
| **Source / exec** | `bash -c 'source ~/api_keys.sh; env'`, diff the env | **REJECTED** |
| **Parse** | line-scan `export VAR=value`, `os.Setenv` (sibling `loadFromShell`) | **ADOPT** |

**Why parse, not source (§11.4.10 blast-radius).** Sourcing executes whatever is
in the file. A tampered or careless `api_keys.sh` could contain `rm -rf`, a
`curl … | sh` exfil, or a fork bomb; even an honest file that references an
undefined var aborts under `set -u` (the exact footgun the sibling shell wrapper
works around at `load_api_keys.sh:16-34`). Executing host-home shell as the app
is an arbitrary-code-execution surface with the app's full privileges. The Go
in-process loader has no in-process shell, and shelling out only to read
key=value pairs trades a security hole for zero benefit. **Parse.**

**What the parser accepts** (port of `loader.go:54-81`): lines beginning
`export ` (after trim); one `key=value` per line at the first `=`; single- and
double-quote stripping; `#`-comment and blank-line skipping; **non-`export`
lines skipped** so a stray `VAR=…` cannot leak. Values are never logged.

**Documented limitation — no shell evaluation in-process.** The parser does
**not** expand `${OTHER}` references, `$(command)` substitution, or line
continuations. So an `api_keys.sh` line like
`export VERTEX_API_KEY=${ApiKey_Google_Vertex_AI}` stores the **literal**
`${ApiKey_Google_Vertex_AI}` string, not its value (the sibling's real-source
shell wrapper *does* expand it — `load_api_keys.sh:17-25`). This is a security
FEATURE (no eval), stated honestly as a compatibility gap (§11.4.6). Operators
who genuinely need shell-eval semantics use the source-able wrapper (§8.3) at a
**shell** entrypoint, where a real shell — not the app — does the evaluation.
If bounded in-process eval is ever required it MUST be bounded (scrubbed env,
`context` timeout, locked-down `sh -c`, capture only the env delta) — but that
reintroduces exec risk and is **not** recommended; default remains parse.

**File-permission hygiene (advisory, not enforced on a host-owned file).** The
loader MAY `os.Stat` `api_keys.sh` and emit a NAME-only warning (never a value)
if it is group/world-readable, mirroring §11.4.10's `chmod 600` expectation. It
must not refuse to load on a permission finding (the file is host-operator
owned), only surface it.

---

## 5. `.env` consumption

Standard dotenv `VAR=value`, one per line, `#` comments, blanks skipped, quote
stripping — the sibling `loadFromEnv` (`loader.go:84-106`). Discovery: walk up
from cwd to the first `.env` (`findEnvFile`, `loader.go:133-143`); the sibling
shell wrapper additionally anchors on the `.gitmodules` meta-repo root
(`load_api_keys.sh:41-50`) — for this project the plain walk-up is sufficient
(the repo root holds both `.gitmodules`-free `go.mod` and, if present, `.env`).
`.env` is also parsed, never sourced (a dotenv file is data, not a script; same
blast-radius argument). Same gap-fill precedence via `setIfAbsent`.

The project already has a rich `.env.example` (DB/HTTP/worker settings,
`project/.env.example`) whose values feed the compose/deploy stack, **not** the
Go binary's config loader. Provider keys are a **new section** to add there
(§7), NAMES only.

---

## 6. Env-var name registry (NAMES ONLY — no values anywhere)

### 6.1 Provider keys the loader targets

Canonical names (what `config.toml` `${…}` interpolation and the provider
constructors read):

| Canonical name | Consumed by | Status |
|---|---|---|
| `ANTHROPIC_API_KEY` | R19 `NewAnthropicLLM` / R22 `providers/anthropic.NewProvider`; `config.toml` `llm_api_key = "${ANTHROPIC_API_KEY}"` | **First-class (R19)** |
| `OPENAI_API_KEY` | `EmbeddingConfig.APIKey` (`config.go:110`) & `openai` LLM provider; `config.toml` `api_key = "${OPENAI_API_KEY}"` | **First-class** |
| `GEMINI_API_KEY`, `DEEPSEEK_API_KEY`, `MISTRAL_API_KEY`, `GROQ_API_KEY`, `XAI_API_KEY`, `OPENROUTER_API_KEY`, `QWEN_API_KEY` | LLMProvider catalogue (R22 natively-wired set, `verifier_dynamic_catalogue.go:30-43`) | Optional / forward-compatible |

Legacy `~/api_keys.sh` aliases the normaliser bridges to the canonical names
above: `ApiKey_Anthropic`, `ApiKey_OpenAI`, `ApiKey_Gemini`, `ApiKey_DeepSeek`,
`ApiKey_Mistral_AiStudio`, `ApiKey_Groq`, `ApiKey_XAI`, `ApiKey_OpenRouter`,
`ApiKey_Qwen` (mapping table modelled on `load_api_keys.sh:82-110`, scoped to the
wired set).

### 6.2 App-own secrets (already defined, loadable from the same files)

Distinct from provider keys but resolvable through the same sources +
`config.go` interpolation/overrides: `HELIX_API_KEYS` (`config.go:398`,
X-API-Key allowlist), `HELIX_DB_PASSWORD` / `HELIX_DB_USER` / … (`config.go:372-403`).
These already work via `applyEnvOverrides` + `${VAR}`; the new loader simply
makes them resolvable from `api_keys.sh`/`.env` too, for free.

### 6.3 Name-normalisation step (bridge `ApiKey_*` → canonical)

After `LoadAPIKeys()`, run `NormalizeProviderKeys()` — a Go port of
`helixcode_normalise_api_keys` (`load_api_keys.sh:80-121`) scoped to the §6.1
wired set: for each `(ApiKey_<Provider>, <PROVIDER>_API_KEY)` pair, if the
canonical is unset/empty and the `ApiKey_` alias is non-empty, set the canonical
(gap-fill; caller/canonical wins). Without this an operator whose `api_keys.sh`
carries only `ApiKey_Anthropic` leaves `${ANTHROPIC_API_KEY}` unresolved — the
readiness bluff the sibling explicitly guards (`load_api_keys.sh:68-75`).
Alternatively (equivalent, no static table) resolve at read time via a
`DeriveKeyEnvAliases`-style helper (`verifier_dynamic_catalogue.go:56-73`); the
normalise-on-load table is simpler and is recommended for the MVP.

---

## 7. `.gitignore` + `.env.example` plan

### 7.1 `.gitignore` tightening (§11.4.30)

Current `project/.gitignore` ignores only `.env` and `.env.local`. Broaden to
the §11.4.30 forbidden set while keeping `*.example` tracked:

```gitignore
# Secrets — VALUES never tracked (§11.4.10 / §11.4.30)
.env
.env.*
*.env
api_keys.sh          # a REAL host-style key file must never be committed
!.env.example        # NAMES-only placeholder is the sole tracked variant
!deploy/.env.example
```

Anti-bluff (§11.4.30): the ignore line is necessary but not sufficient — a gate
must assert **no file matching the forbidden pattern is currently tracked** (§9).

### 7.2 `.env.example` — add a provider-keys section (NAMES only)

Append to the existing `project/.env.example`:

```dotenv
# =============================================================================
# LLM / Embedding provider API keys  (NAMES ONLY — real values go in .env or
# ~/api_keys.sh, both git-ignored). §11.4.10.
# =============================================================================
ANTHROPIC_API_KEY=        # R19 Anthropic Messages API (llm_provider = "anthropic")
OPENAI_API_KEY=           # embeddings + openai LLM provider
# GEMINI_API_KEY=
# DEEPSEEK_API_KEY=
# MISTRAL_API_KEY=
```

### 7.3 New tracked `api_keys.sh.example` (NAMES only)

A committed placeholder documenting the host-home convention, both naming
schemes, empty values:

```sh
# ~/api_keys.sh — copy to $HOME/api_keys.sh, fill real values (NEVER commit the
# real file; it is git-ignored). Loaded by the app in-process (parsed, not
# sourced). §11.4.10 — names only here.
export ANTHROPIC_API_KEY=      # canonical name the app reads
export OPENAI_API_KEY=
# Legacy ApiKey_<Provider> aliases are auto-normalised to the canonical names:
# export ApiKey_Anthropic=
# export ApiKey_OpenAI=
```

---

## 8. Wiring into `config.go` + LLMProvider

### 8.1 New package `internal/secrets` (port of the sibling)

Add `project/internal/secrets/loader.go` under module
`github.com/helixdevelopment/skill-system/internal/secrets`, porting
`loader.go` (§2.1) plus §3.1 Option B (read both, gap-fill) and the §6.3
`NormalizeProviderKeys()`. Zero new third-party deps — pure stdlib (`bufio`,
`os`, `path/filepath`, `strings`), exactly as the sibling. (Note: the project
deliberately has **no** `godotenv` dependency; the hand-rolled parser is the
established, dependency-free, source-safe choice.)

### 8.2 Startup call — before `config.Load()`

In every entrypoint under `project/cmd/*/main.go`, as the **first** action:

```go
_ = secrets.LoadAPIKeys()      // non-fatal; values never logged (§11.4.10)
secrets.NormalizeProviderKeys()
cfg, err := config.Load("")    // now ${ANTHROPIC_API_KEY} etc. resolve
```

Ordering is load-bearing: `config.go`'s `substituteEnv` (`config.go:285-339`)
reads `os.Getenv` at `interpolate` (`:356`), so the env must be populated first —
identical to the sibling's `main.go:99` before `:105`. The call returns a bool
(sibling `main.go:81-83`) so a wiring test can assert the loader actually runs
on this project's server/CLI/worker paths (anti-bluff: proves wiring, not just
presence).

### 8.3 Config field + interpolation for the Anthropic key (R19)

Land R19 §2.5: add `LLMAPIKey`/`LLMBaseURL` to `AutoExpandConfig`
(`config.go:124-130`) and the symmetric substitution line
`cfg.AutoExpand.LLMAPIKey = sub(cfg.AutoExpand.LLMAPIKey)` beside the existing
`LLMProvider`/`LLMModel` subs (`config.go:325-326`). `config.toml` then carries
`llm_api_key = "${ANTHROPIC_API_KEY}"` — a **name**, resolved to the value only
in the process env at runtime. The existing fail-closed check for
uninterpolated `${` placeholders (`config.go:438-447`) should be extended to
`LLMAPIKey`/`EmbeddingConfig.APIKey` so an unset referenced var fails closed
rather than shipping a literal `${…}` as a "key" (values never echoed —
field-name + index only, per the existing pattern).

### 8.4 Feeding LLMProvider (R19 interim client and R22 submodule)

Both consumers take the resolved key as a **string constructor argument**:

- R19 interim: `NewAnthropicLLM(cfg.AutoExpand.LLMAPIKey, cfg.AutoExpand.LLMModel, logger)`
  (`r19_…:237,285-290`).
- R22 submodule: `providers/anthropic.NewProvider(apiKey, baseURL, model)` /
  `providers/openai.NewProvider(...)` (`r22_…:77-82`).

The submodule/client never touches the environment — so the **only** thing the
credential layer must guarantee is that `cfg.AutoExpand.LLMAPIKey` /
`cfg.Embedding.APIKey` are populated, which §8.2+§8.3 achieve. The loader design
is identical whether R19's thin client or R22's submodule is the active
consumer; R22 adoption is gated on the G14/X1 multi-root submodule decision
(`r22_…:39-59,320`) and does **not** block this layer.

---

## 9. Leak-audit gate (§11.4.10.A) — before storing any provided secret

When the operator hands over a real token to store in `.env` / `~/api_keys.sh`,
run this audit **before** writing it (values never echoed — operate on a local
SHA-256 fingerprint, never print the value):

1. **Tree scan:** `git ls-files -z | xargs -0 grep -lF -- "$VALUE"` — any tracked
   file containing the value.
2. **History scan:** `git log -S"$VALUE" --all --source --remotes` — any commit
   that ever added/removed it.
3. **Surface findings to the operator first** — operator chooses rotate /
   accept-as-compromised / abort (§11.4.10.A(3)).
4. **On any hit:** open a §6/§7 sixth-law-incident record, redact the tracked
   file in-place to `<redacted-per-§11.4.10>`, and record OPERATOR ACTION
   REQUIRED = rotate the key (§11.4.10 clause 7).
5. **Extend the pre-push hook** credential-pattern grep to catch the escaped
   class in the same change.

### 9.1 Standing gate — `CM-CRED-NO-VALUE-IN-TREE` (anti-bluff)

A pre-build/pre-commit gate that:

- asserts no `.env` / `.env.*` / `*.env` / real `api_keys.sh` is **tracked**
  (§7.1, §9 of the constitution);
- scans tracked files for high-entropy provider-key **shapes** (`sk-ant-…`,
  `sk-…`, `AIza…`, generic ≥32-char base64 assigned to a `*_API_KEY`/`ApiKey_*`
  name), reporting NAME + path + fingerprint, never the value.

Self-validation (§11.4.107(10)): golden-good (clean tree → PASS), golden-bad (a
temp tracked file carrying a **fake** `sk-ant-FAKE-000…` value → FAIL, offender
pinpointed), negative-control (the NAMES-only `.env.example` / `api_keys.sh.example`
→ PASS). A gate that passes its golden-bad is itself the bluff. Paired §1.1
mutation: strip the scanner's shape pattern → its golden-bad case FAILs.

---

## 10. Test / validation plan (anti-bluff)

All fixtures use obviously-fake tokens (e.g. `test-not-a-real-key`,
`sk-ant-FAKE-000`); **no secret VALUE ever enters a fixture or a log**
(§11.4.10). Tests are hermetic — `t.TempDir()` for `HOME` and cwd, env
save/restore — modelled on the sibling `loader_test.go` (§2.5).

| # | Proof required | Test | Captured-evidence shape (no values) |
|---|---|---|---|
| a | key present ONLY in `api_keys.sh` loads | write `$HOME/api_keys.sh` = `export ANTHROPIC_API_KEY=test-not-a-real-key`; no `.env`; call `LoadAPIKeys()`; assert `os.Getenv("ANTHROPIC_API_KEY") == "test-not-a-real-key"` | test PASS line "loaded from api_keys.sh"; value asserted equal to the **known fixture constant**, never printed |
| b | key present ONLY in `.env` loads | write `<cwd>/.env` = `ANTHROPIC_API_KEY=test-not-a-real-key`; no `api_keys.sh`; assert resolves | PASS line "loaded from .env" |
| c | precedence resolves as specified | table-driven matrix: process-env vs `api_keys.sh` vs `.env` — assert **process-env wins over both**, **`api_keys.sh` wins over `.env`** (Option B: `.env`-only var still loads when `api_keys.sh` exists) | precedence truth-table, all rows PASS; mirrors sibling `TestLoadAPIKeys_PrefersShellOverEnv` |
| d | leak-audit gate catches a committed value | plant a **fake** key value in a temp tracked file → run `CM-CRED-NO-VALUE-IN-TREE` → assert **FAIL (RED)**; remove → assert **PASS (GREEN)** | gate exit codes (RED→GREEN, §11.4.115), matched **fingerprint** only |

Additional coverage porting the sibling suite: strips-quotes, ignores-comments,
ignores-blank, skips-non-`export`-lines, neither-found-returns-error,
`NormalizeProviderKeys` maps `ApiKey_Anthropic`→`ANTHROPIC_API_KEY` (gap-fill,
canonical-wins), parse-not-source (an `api_keys.sh` line
`export EVIL=$(touch /tmp/pwned)` stores the literal string and **creates no
file** — proves no shell executed). Paired §1.1 mutations: (i) invert
`setIfAbsent` so file overwrites process env → precedence test FAILs; (ii) drop
the `export ` skip → non-`export` line pollutes env, test FAILs; (iii) strip the
gate shape pattern → §9.1 golden-bad FAILs. Wiring assertion: a
`loadapikeys_wiring_test.go` (sibling name) asserts the loader is invoked on the
server/CLI/worker startup paths.

---

## 11. Open questions / UNCERTAIN items

- **[DECISION — recommended, needs operator sign-off] Option A (exclusive,
  sibling-exact) vs Option B (union gap-fill).** §3.1. Recommend **B** as the
  literal reading of "both fully supported"; it diverges from `helix_code`'s
  exclusive loader only in that a `.env`-only var loads even when `api_keys.sh`
  exists. If strict sibling parity is preferred over the mandate's "both"
  wording, choose A. Precedence is deterministic either way.
- **[SCOPE] Normalisation table breadth.** §6.3 scopes `ApiKey_*`→canonical to
  the wired set {anthropic, openai, +7}. The sibling table has 29 entries
  (`load_api_keys.sh:82-110`). Recommend the minimal wired set for the MVP,
  extensible as providers are wired — but confirm no other provider is expected
  at launch.
- **[DEPENDENCY, not a blocker] R22 LLMProvider submodule adoption is gated on
  G14/X1** (`r22_…:39-59,320`). Until it lands, R19's thin `net/http` Anthropic
  client is the consumer. The loader is identical for both; no rework when R22
  lands.
- **[UNCERTAIN] Shell-eval semantics in `api_keys.sh`.** §4: the parse-don't-source
  loader cannot expand `export X=${Y}` / `$(...)` in-process. Confirmed the
  sibling Go loader has the same limitation (it also parses). If any operator's
  real `api_keys.sh` relies on cross-variable expansion, they must use the
  source-able shell wrapper (§8.3 / `load_api_keys.sh`) at a shell entrypoint —
  flag whether the project needs that wrapper shipped too, or whether the
  in-process parser suffices for the MVP.
- **[CONFIRM] Whether to mirror the source-able `scripts/load_api_keys.sh`**
  wrapper into `project/scripts/` for shell entrypoints/CI. Recommended as a thin
  convenience, but the Go in-process loader is authoritative for the running app;
  low priority for the MVP.

**No BLOCKERS.** Every source needed to implement the layer exists and was read
directly; the design is a faithful port of the established, cited `helix_code`
pattern with one flagged, operator-gated refinement (Option B) and the R19
config-field additions already specified in `r19_anthropic_api_support_design.md`.

## 12. Conductor resolutions of the flagged decision points (§11.4.101)

All three §11 UNCERTAIN points are reversible + evidence-backed, so they are
resolved autonomously per §11.4.101 (no operator block needed) — recorded here
as FACT, not guess (§11.4.6):

1. **Precedence = UNION (Option B), order `process-env > ~/api_keys.sh > .env`.**
   Rationale: the operator's literal words are "Both MUST BE fully supported."
   Exclusive mode (sibling `helix_code` behaviour: api_keys.sh present ⇒ .env
   never read) would make `.env` *unsupported* whenever api_keys.sh exists — the
   opposite of "both fully supported." A var present only in `.env` MUST load.
   This is a deliberate, documented divergence from strict sibling parity.
   REVERSIBLE (one precedence constant) if the operator later prefers exclusive.
2. **`~/api_keys.sh` = PARSE, never source (no in-process shell eval).** Matches
   the sibling Go loader + §11.4.10 blast-radius. Cross-variable forms
   (`export X=${Y}`) are therefore NOT expanded in-app; an operator needing that
   uses a source-able shell wrapper OUTSIDE the app process. Security-first,
   non-negotiable for the in-app loader.
3. **Normalisation table = minimal wired set for the MVP** ({anthropic, openai}
   + the 7 already-referenced providers), expanded per-provider as providers are
   added — not the sibling's full 29-entry table up front. YAGNI; expansion is a
   pure data addition.

These resolutions bind the R7/R19 implementation task; none is an open blocker.
