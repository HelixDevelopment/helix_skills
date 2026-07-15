# Helix-Family Reusable Practices Survey (R21 + R22)

**Revision:** 2
**Last modified:** 2026-07-15T18:10:00Z
**Status:** survey / catalogue-first research (§11.4.74), no code written, no
existing file modified.
**Scope:** operator mandate R21 (`REQUIREMENTS.md:249-254`) — survey the
sibling `helix_*` projects under
`/home/milos/Factory/projects/tools_and_research/` for reusable universal
architecture / submodules / practices the HelixKnowledge Skill Graph System
MUST adopt, per §11.4.74 (catalogue-first, extend-don't-reimplement) and
§11.4.28 (owned submodules are equal codebase) — **broadened by operator
mandate R22** (same date) to (a) ground every verdict in the real
`vasic-digital`/`HelixDevelopment` repo contents (README + go.mod/
helix-deps.yaml, not the sibling checkouts alone), (b) add a DECISIVE
verdict for `MCP_Module`, `DebateOrchestrator`, and `token_optimizer`
alongside the R21 items, and (c) FLAG (name + one-liner, no full design) a
second tier: `DagOrchestrator`, `PipelineRuntime`, `VisionEngine`/`Panoptic`,
the Kotlin-Multiplatform (KMP) suite, and the generic `config`/`security`/
`Storage`/`Filesystem` modules.
**Siblings surveyed (read-only):** `helix_code`, `helix_ota`,
`helix_terminator`, `helix_track`, `helix_translate`, `helix_vpn`, plus the
target's own `constitution/submodules/` tree, plus 5 live `gh repo
view`/`gh api` lookups against GitHub for repos not checked out on this host
(§10A).
**Method:** every claim below is cited to a real file path (and, where the
sibling is a checked-out git submodule, its `go.mod`/README/source), or to a
live `gh` API read where noted. Nothing is invented; unread material is
marked `UNCONFIRMED:` per §11.4.6.

---

## 0. TL;DR

- **Containers (R20):** ADOPT `vasic-digital/containers` at the **root path**
  `containers/` (not `submodules/containers/`) — that is the layout the only
  sibling with REAL production Go code wired against it (`helix_ota`)
  actually uses (`server/go.mod:66` `replace digital.vasic.containers =>
  ../containers`).
- **HelixQA (R8):** target's `helix-deps.yaml` already vendors `helix_qa` +
  `challenges` correctly, but is **missing `containers`** — HelixQA's own
  README states it is "Built on `digital.vasic.challenges` and
  `digital.vasic.containers` — both incorporated at the **parent project's
  root**" (`helix_ota/submodules/helixqa/README.md:19-22`). Without
  `containers` at the target's root, HelixQA cannot boot its own test infra.
- **Docs Chain (R10):** already fully researched
  (`research/docs_chain_incorporation.md`) — this survey adds one new
  cross-check: 3 of 4 siblings that vendor it (`helix_ota`, `helix_track`,
  `helix_translate`) mount it at the **root** path `docs_chain/`; only
  `helix_vpn` uses `submodules/docs_chain/`. The target's planned
  `submodules/docs_chain/` (per its own `project/helix-deps.yaml` `layout:
  grouped`) is the **minority** convention in the family (not wrong per
  §11.4.28C, just worth a deliberate choice, not a default).
- **LLMProvider (R7/R19 — CRITICAL REDIRECT):** `vasic-digital/LLMProvider`
  (module `digital.vasic.llmprovider`) is a real, MIT-licensed, 50+-provider
  Go LLM client library — including a **production `pkg/providers/anthropic`**
  (Messages API, streaming, tool-calls) and **`pkg/providers/openai`** — plus
  circuit-breaker + health-monitor + retry, all satisfying one common
  `LLMProvider` interface. This **supersedes** the plan to hand-roll a
  `AnthropicLLM` for R19: vendor the module instead of writing a second
  Anthropic client. See §4 for the full verdict.
- **open-design (R12):** already fully researched
  (`research/opendesign_incorporation.md`) — corroborated here: 2 siblings
  (`helix_terminator`, `helix_vpn`) vendor it at the **identical** path
  `submodules/open-design/`, confirming that path as the family convention.
- **http3:** `helix_ota` vendors `vasic-digital/http3` at
  `submodules/http3/` and uses it in **live production code**
  (`server/internal/transport/transport.go`,
  `server/cmd/ota-server/main.go`) — not a stub. Target currently uses
  `quic-go` directly (per task brief); recommendation is STUDY-pattern-only
  pending a size comparison (honest gap, §7).
- **HelixLLM + Helix-Track/Core (R4):** already fully researched
  (`research/helix_interop_incorporation.md`) — this survey adds one new
  data point: `helix_terminator` vendors `Helix-Track/Core` as a **build
  dependency** (`submodules/helixtrack-core/`), a different pattern than the
  target's own design (Core as a remote REST+JWT peer service, never
  vendored). Both are legitimate; pick REST-peer unless target ever needs to
  import Core's Go types directly.
- **helix_track client matrix (R3/G15):** helix_track ships **one submodule
  per platform** (`Web-Client`, `Desktop-Client`, `Android-Client`,
  `iOS-Client`, `Harmony-OS-Client`, `Aurora-Client`, plus `Core` +
  `Screensaver`) — Android and iOS are **separate native codebases**
  (Kotlin/Gradle, Swift), Desktop is **Tauri wrapping the Angular Web app**
  (a genuine cross-surface shared-code win), and **Harmony-OS-Client +
  Aurora-Client are EMPTY governance stubs with zero source code** — they are
  NOT a validated reference for HarmonyOS/Aurora; helix_track is exactly as
  unproven there as the target project. helix_track uses **no Flutter
  anywhere** — worth flagging against the target's own prior "Flutter for
  Android/iOS/HarmonyOS/Aurora/Desktop" lean.
- **Other:** all 7 siblings + target vendor `HelixDevelopment/HelixConstitution`
  at the identical path `constitution/` (100% consistent); `vasic-digital/Auth`
  (JWT/API-key/OAuth2/middleware) is a family-standard building block, not yet
  required by target's stated requirements — STUDY-pattern-only.
- **R22 additions — MCP_Module: STUDY-only** (target's `mark3labs/mcp-go`-based
  MCP layer is already working + hardened; switching would be a rewrite with
  no functional gain). **DebateOrchestrator: ADOPT-core + EXTEND-stubs** for
  the G05 LLM jury. **token_optimizer: ALREADY VENDORED** inside the target's
  own `constitution/submodules/token_optimizer/` — wire it, don't add it.
  **`vasic-digital/design_system` (found via live `gh` lookup, R12): ADOPT
  alongside `open-design`** — a real, already-extracted, light+dark,
  Tailwind-v4 CSS design-token package used by 4 sibling Helix web surfaces.
  See §5, §6A–§6C, §8A, §10A.

---

## 1. Containers (R20) — how the family wires `vasic-digital/containers`

**Real usage found in `helix_ota` (the only sibling with live Go code
against it, not just a declared submodule):**

- `.gitmodules` declares it **twice** — once at root path `containers`
  (`helix_ota/.gitmodules:3-4`, path=`containers`) and once at
  `submodules/containers` (`.gitmodules:14-15`) — both pointing to the
  identical URL `git@github.com:vasic-digital/containers.git`. Both
  directories are populated on disk. This is an unresolved duplication in
  the sibling itself (see honest gap §7) — the load-bearing copy is the
  **root** one.
- `helix_ota/server/go.mod:32,66`: `digital.vasic.containers
  v0.0.0-00010101000000-000000000000` + `replace digital.vasic.containers
  => ../containers` — the Go build resolves the dependency to the
  **root-level** `containers/` directory (sibling of `server/`), not
  `submodules/containers/`.
- Real production/test usage (`helix_ota/server/internal/store/postgres_integration_test.go:1-20`):
  imports `digital.vasic.containers/pkg/boot` + `pkg/compose` to boot a
  **real PostgreSQL** on demand for integration tests — the file's own
  header comment states the intent verbatim: "never a manual
  `podman`/`compose` step, never a fake … reuse the containers catalogue
  brick, never reimplement orchestration" (citing §11.4.74 by name).
  Same pattern repeats across `internal/rollout/*_integration_test.go`.
- Additional consumers of the same root `containers/` tree inside
  `helix_ota`: `pkg/boot`, `pkg/compose`, `pkg/lifecycle`, `pkg/lazyservice`,
  `pkg/orchestrator`, `pkg/remote`, plus two standalone binaries
  `containers/cmd/deploy-stack/main.go` and
  `containers/cmd/ota-device-emu-boot/main.go`.
- `helix_track` and `helix_translate` also mount it at the **root** path
  `containers/` (`.gitmodules`); `helix_terminator` and `helix_vpn` mount it
  under `submodules/containers/`. Root-path is the more common convention
  among siblings that actually run infra-heavy Go tests (`helix_ota`,
  `helix_track`).

**Recommendation: ADOPT-as-submodule**, root path `containers/`
(`<repo_root>/containers/`, ungrouped per §11.4.28C), wired via a Go
`replace digital.vasic.containers => ../containers` directive exactly like
`helix_ota/server/go.mod:66`. Consume `pkg/boot` + `pkg/compose` for every
integration test that needs a real Postgres/pgvector instance instead of a
docker-compose-by-hand step (directly serves R1's "docker-compose up smoke
test" and R15's "Docker/Podman Compose orchestrated underneath").

**Serves:** R20 (explicit mandate), R9/G13 ops-hardening compose profiles,
and — critically — **R8/HelixQA** (§2 below), because HelixQA cannot boot
its own test infra without `containers` present at the consumer's root.

**Gated on G14/X1:** YES — R20 explicitly says "vendoring is gated on the
G14/X1 submodule-policy decision" (`REQUIREMENTS.md:247`). This survey's
finding narrows that decision: whichever layout X1 picks for grouped vs.
ungrouped deps in general, `containers` specifically should follow
`helix_ota`'s **root-path** precedent, because that is the only
precedent backed by working, tested Go code.

---

## 2. HelixQA (R8/R11/§11.4.27) — how `helix_ota` wires `HelixDevelopment/HelixQA`

**What it is (verified `helix_ota/submodules/helixqa/README.md:1-25`):**
"an anti-bluff QA orchestration framework for cross-platform testing with
real-time crash detection, step validation, evidence collection, and
automated ticket generation." Built on `digital.vasic.challenges` and
`digital.vasic.containers` — the README states explicitly both **"MUST NOT
introduce its own `.gitmodules` entries for those repos"** and must be
"incorporated at the parent project's root per CONST-051(C)" — this is the
**exact §11.4.28C no-nested-own-org-chain rule in action**, verified from a
real consumer, not just from constitution text.

**Real CLI surface (`helix_ota/submodules/helixqa/cmd/`, 30+ binaries):**
the main `helixqa` orchestrator binary, `helixqa-bank-session` (dispatches a
named test bank), and a family of `helixqa-verify-*` binaries per capability
class (`helixqa-verify-coder-bench`, `-chaos`, `-concurrency`, `-ddos`,
`-memory`, `-race`, `-rag`, `-vision`, `-whisper`, `-tesseract`,
`-mcp-gateway`, `-embeddings`, `-netprov`, `-a2a`), plus recording/analysis
tools (`qa-audio-probe`, `recording-analyzer`, `helixqa-omniparser`,
`helixqa-uitars`). Test banks live under `helixqa/banks/*.{yaml,json}` (e.g.
`atmosphere.yaml`, `admin-operations.yaml`, `all-formats.yaml`) — this **is**
the "HelixQA test banks/suites" R8 refers to, confirmed as a real,
populated directory, not an abstraction.

**Honest gap found in the sibling itself:** `helix_ota`'s own top-level
`scripts/` and `Makefile` contain **zero** invocations of `helixqa` or its
banks (`grep -rli helixqa scripts/` → 2 unrelated hits in
`codegraph_validate.sh`/`sync_md_siblings.sh`; `grep helixqa Makefile` → 0
hits). `helix_ota` vendors HelixQA but does not visibly wire a
project-level "run the QA session" entry point at its own root — the
consumable pattern is (a) vendor `helix_qa` + its declared deps
(`challenges`, `containers`) at the project root, (b) author
project-specific `banks/*.yaml`, (c) invoke the `helixqa` /
`helixqa-bank-session` binary against those banks directly — there is no
extra sibling-side wrapper script to imitate.

**Target's current state:** `project/helix-deps.yaml` already vendors
`helix_qa` (`git@github.com:HelixDevelopment/HelixQA.git`) and `challenges`
(`git@github.com:vasic-digital/Challenges.git`) — correct per this
convention — but does **NOT** vendor `containers`, which HelixQA's own
README says is a hard structural dependency for booting its test infra.

**Recommendation: ADOPT-as-submodule** (already declared, correctly) **+ add
`containers`** to `project/helix-deps.yaml` (per §1) so HelixQA's structural
dependency is satisfied. Author skill-graph-specific `banks/*.yaml`
(candidate bank names: `skillgraph-wizard.yaml` for the R6 create→map→process
flow, `skillgraph-mcp.yaml` for the MCP tool surface) and invoke via the
vendored `helixqa` binary directly — no extra orchestration wrapper needed.

**Serves:** R8, R11.
**Gated on G14/X1:** the `containers` half is gated (§1); `helix_qa` +
`challenges` vendoring itself is not gated — already correctly present.

---

## 3. Docs Chain (R10/§11.4.106) — cross-check against the already-completed research

`research/docs_chain_incorporation.md` (Revision 1, this project) already
did the deep incorporation research for this project directly against the
real engine at `helix_code/submodules/docs_chain`. This survey adds **one**
new cross-family data point, not a re-derivation:

| Sibling | Mount path | Layout |
|---|---|---|
| `helix_ota` | `docs_chain/` (root) | ungrouped |
| `helix_track` | `docs_chain/` (root) | ungrouped |
| `helix_translate` | `docs_chain/` (root) | ungrouped |
| `helix_vpn` | `submodules/docs_chain/` | grouped |
| `helix_code` | `submodules/docs_chain/` | grouped |
| target (`project/helix-deps.yaml`) | `submodules/docs_chain/` (declared) | grouped |

**Finding:** root-path (ungrouped) is the **majority** convention (3 of 5
siblings that vendor it), not the minority the target's manifest currently
declares. This does not make the target's declared `layout: grouped`
wrong — §11.4.28C permits either — but it means the target is currently
diverging from the majority-family convention for this specific dependency,
and that choice should be made deliberately (e.g. as part of the same G14/X1
decision, or explicitly re-affirmed), not left as an unexamined default.

**Recommendation:** unchanged from `docs_chain_incorporation.md` (incorporate
by reference, engine must be `go build`'t on this host — the committed
binary is wrong-arch); **additionally** reconsider `layout: grouped` vs. the
family-majority root-path convention as part of the G14/X1 decision.

**Serves:** R10, §11.4.106.
**Gated on G14/X1:** the layout choice, yes; the incorporation mechanics,
already fully designed independent of X1.

---

## 4. LLMProvider (R7/R19) — THE CRITICAL REDIRECT VERDICT

### 4.1 What `vasic-digital/LLMProvider` actually is (verified from a checked-out copy)

Checked out identically at `helix_terminator/submodules/llmprovider/`,
`helix_ota/submodules/llm_provider/`, `helix_translate/llm_provider/`,
`helix_vpn/submodules/llm_provider/` (module `digital.vasic.llmprovider`,
`go.mod:1`, Go 1.25.3, MIT license, `doc.go:1-60`).

**It is not a duplicate of HelixLLM/HelixAgent.** It is a **generic Go LLM
provider abstraction + resilience library**:

- **Core interface** (`doc.go:16-23`):
  ```go
  type LLMProvider interface {
      Complete(ctx context.Context, req *models.LLMRequest) (*models.LLMResponse, error)
      CompleteStream(ctx context.Context, req *models.LLMRequest) (<-chan *models.LLMResponse, error)
      HealthCheck() error
      GetCapabilities() *models.ProviderCapabilities
      ValidateConfig(config map[string]interface{}) (bool, []string)
  }
  ```
- **Resilience primitives shipped alongside:** `CircuitBreaker`
  (closed/open/half-open, `circuit_breaker.go`), `HealthMonitor`
  (`health_monitor.go`), `RetryConfig` (exponential backoff + jitter +
  HTTP-status-aware retry, `retry.go`), `LazyProvider` (deferred init).
- **50+ CONCRETE provider implementations already written**, each in its own
  `pkg/providers/<name>/` package: **`anthropic`** (`pkg/providers/anthropic/anthropic.go`
  — real `https://api.anthropic.com/v1/messages`, `APIVersion = "2023-06-01"`,
  full `Tool`/`ToolChoice`/streaming/system-prompt request shape),
  **`openai`** (`pkg/providers/openai/openai.go` — `/v1/chat/completions`,
  `/v1/models`), a separate **`claude`** package
  (`pkg/providers/claude/claude.go` — Claude-Code-OAuth-aware variant that
  explicitly detects and rejects the "This credential is only authorized for
  use with Claude Code" restriction, i.e. it already handles the exact OAuth
  vs. API-key distinction the toolkit-alias P3 chain needs), plus `deepseek`,
  `gemini` (API + CLI-stub variants), `mistral`, `cohere`, `groq`, `xai`,
  `ollama`, `openrouter`, `huggingface`, `together`, `fireworks`, `cerebras`,
  `sambanova`, `nvidia`, `kimi`, `qwen` (API + CLI-stub), `zai`, `zhipu`,
  `perplexity`, `replicate`, `novita`, `cloudflare`, `githubmodels`, `modal`,
  `ai21`, `codestral`, `sarvam`, `upstage`, `siliconflow`, `hyperbolic`,
  `venice`, `nlpcloud`, `publicai`, `chutes`, `vulavula`, `nia`, `junie`
  (API + CLI-stub), `generic` (catch-all OpenAI-shaped fallback).

**Org discrepancy (flagged, not resolved here):** `helix_terminator`,
`helix_ota`, `helix_translate`, `helix_vpn` all point at
`git@github.com:vasic-digital/LLMProvider.git`; `helix_code` and
`helix_track` point at `git@github.com:HelixDevelopment/LLMProvider.git` —
same module name (`digital.vasic.llmprovider`) declared in both checkouts,
but the two org-URL copies were **not diffed against each other** in this
survey (see §7 honest gaps). Treat as the same logical project mirrored
across two owned orgs unless a diff proves otherwise.

### 4.2 Target's current hand-rolled state (`internal/autoexpand/llm.go`)

```go
// LLMClient abstracts LLM API calls for skill generation.
type LLMClient interface {
    Generate(ctx context.Context, prompt string, maxTokens int) (string, error)
}
```

Single-method, no streaming, no health check, no capability negotiation, no
circuit breaker, no retry — `OpenAILLM` hand-rolls its own `net/http` calls.
R19's own design doc (`research/r19_anthropic_api_support_design.md`) plans a
second hand-rolled `AnthropicLLM` implementing the same narrow interface.

### 4.3 Verdict (FACT, not a guess — §11.4.6)

**Vendoring `vasic-digital/LLMProvider` supersedes the plan to hand-roll a
second (`AnthropicLLM`) provider for R19.** The module already ships a
production-shape Anthropic Messages-API client (streaming, tools, system
prompt) and an OpenAI client, both satisfying one common resilient interface
with circuit-breaker/health-monitor/retry included for free — capabilities
target's current `LLMClient` has none of.

**Adapter cost (the honest tradeoff):** target's `LLMClient.Generate(ctx,
prompt, maxTokens) (string, error)` is narrower than llmprovider's
`Complete(ctx, *models.LLMRequest) (*models.LLMResponse, error)`. Two paths,
either legitimate:

1. **Thin adapter (recommended short-term):** wrap an llmprovider
   `anthropic.Provider` / `openai.Provider` behind a small
   `func (a *providerAdapter) Generate(ctx, prompt, maxTokens) (string,
   error)` that builds an `LLMRequest{Prompt: prompt, MaxTokens: maxTokens}`
   and returns `resp.Content` — existing call sites in
   `internal/autoexpand/pipeline.go` need zero changes.
2. **Widen the interface (recommended follow-on):** replace
   `internal/autoexpand.LLMClient` with llmprovider's own `LLMProvider`
   interface directly, gaining streaming + capability negotiation +
   circuit-breaker for the whole auto-growth pipeline (R14's "real-time
   growth" job pipeline is exactly the kind of long-running work that
   benefits from circuit-breaker + health-monitor).

**This is a redirect, not an addition:** the R19 design doc's plug-in point
(`NewLLMClientFromConfig` factory) should construct an llmprovider provider
+ adapter, **not** a second hand-rolled HTTP client. R7's "pluggable
ModelProvider, not hardcoded OpenAI" is satisfied more completely by
llmprovider's provider registry (50+ providers, one interface) than by
hand-rolling each one.

**Recommendation: ADOPT-as-submodule**, `submodules/llmprovider/`
(matching `helix_terminator`'s path exactly), consumed via the adapter
pattern above. Does **not** replace the P3 `ModelProvider` chain design in
`helix_interop_incorporation.md` (HelixLLM/LLMsVerifier/toolkit-alias/OpenAI)
— it is the concrete Go client layer that chain's HTTP calls should be built
on, in place of the hand-rolled ones.

**Serves:** R7 (directly, supersedes the hand-rolled approach), R19
(directly, supersedes `AnthropicLLM`), R17 resilience gaps (circuit
breaker/health monitor were flagged gaps in the gaps register — this closes
several with an off-the-shelf implementation rather than new code), and R8
indirectly (llmprovider ships its own `challenges/` directory — one more
family precedent for anti-bluff test coverage of a provider abstraction).

**Gated on G14/X1:** not directly — this is a new dependency addition, not
an existing one whose layout is contested. It should, however, follow
whichever root-vs-submodules convention the G14/X1 decision settles on for
consistency (§1).

---

## 5. open-design + `vasic-digital/design_system` (R12/§11.4.162)

`research/opendesign_incorporation.md` (Revision 1) already fully surveyed
`nexu-io/open-design` directly (cloned + inspected). This survey corroborates
the recommended mount path from two independent sibling precedents:
`helix_terminator/.gitmodules:7-9` and `helix_vpn/.gitmodules:872-875` both
mount it at the **identical** path `submodules/open-design/` (URL
`git@github.com:nexu-io/open-design.git`) — the exact path that document's
§3.1 already recommends.

**New finding (R22): `vasic-digital/design_system` is a real, second,
directly-relevant repo** not present in any locally-checked-out sibling but
confirmed via a live `gh repo view vasic-digital/design_system` +
`gh api repos/vasic-digital/design_system/{contents,readme}` read (§10A
method). Its README (`@vasic-digital/design-system`, Revision 1,
2026-07-14) states it is: **"A reusable, decoupled, OpenDesign-driven design
system shared across the Helix web surfaces (HelixOTA · HelixCode ·
HelixTrack · HelixQA) and the vasic-digital org … extracted from the
production HelixOTA website design library."** Verified contents
(`gh api .../contents`): `tokens/core.css` (theme-invariant type/spacing/
radius/elevation/motion, no brand color), `tokens/themes/*.css` (3 shipped
brand themes: `helix-green` default, `vasic-red`, `helix-ota-blue`, **each
light AND dark first-class** via all three sanctioned mechanisms —
`prefers-color-scheme`, `data-theme`, `.dark` class), `tailwind/tailwind-v4.css`,
`fonts/fonts.css`, `components/css/components.css` (framework-agnostic
`.ds-*` primitives), `components/angular/*` (Angular-specific adapters —
`ThemeService`, `I18nService`, theme-toggle + language-picker components),
`i18n/en.json`. Package-shaped (`package.json` present) — an npm-installable
asset, not a raw agent-driven generation tool like open-design. No
`helix-deps.yaml` was found at its root (§10A — either genuinely absent or
the fetch silently 404'd; treat its transitive-dep status as UNCONFIRMED,
not asserted as zero).

**Why this matters more than open-design alone for R12:** open-design is
the *raw* per-brand-DESIGN.md agent-driven generation substrate (150+ brand
systems, only 57/151 with a dark tier per `opendesign_incorporation.md §2.5`
— gap G5 there). `design_system` is the **already-curated, already-dual-themed,
already-Tailwind-v4-bound CONSUMABLE OUTPUT**, extracted from a real shipped
Helix website (HelixOTA) and already reused across 3 other Helix surfaces.
It directly closes `opendesign_incorporation.md`'s gap G5 (dark-mode
brand-dependency) for whichever brand theme the target adopts, and its
`.ds-*` CSS + Tailwind-v4 layer is framework-agnostic — consumable by the
target's own React web client exactly as easily as by design_system's native
Angular adapters.

**Recommendation: ADOPT-as-submodule** `vasic-digital/design_system`
**alongside** `open-design` (open-design = the design-generation tool +
skill corpus; design_system = the ready CSS/Tailwind/token consumable) —
both serve R12, neither substitutes for the other. Mount path: follow the
snake_case convention, `submodules/design_system/` (§11.4.29), pending the
same G14/X1 layout decision as the rest of the family's third-party/owned
deps.

**Serves:** R12 (directly, and closes the opendesign_incorporation.md G5
dark-mode gap). **Gated on G14/X1:** the layout choice, yes; the
incorporation itself, no (third-party-style vendor dep either way).

---

## 6. http3 — `helix_ota`'s live usage vs. target's direct `quic-go`

`helix_ota` vendors `vasic-digital/http3` at `submodules/http3/` (module
`digital.vasic.http3`), wired via `server/go.mod:33,68`: `digital.vasic.http3
v0.0.0-...` + `replace digital.vasic.http3 => ../submodules/http3`. This is
**live production code**, not a stub: `server/internal/transport/transport.go`
and `server/cmd/ota-server/main.go` both import it directly.

Target's own transport layer uses `quic-go` **directly** (stated in the
task brief for this survey; not independently re-derived from target's
source in this pass — see honest gap §7). Two legitimate outcomes depending
on facts not yet gathered:

- If target's own `quic-go` transport code is thin/duplicative of what
  `http3` already wraps (connection setup, TLS config, HTTP/3 server
  bootstrap) → **ADOPT-as-submodule**, replacing the ad-hoc wrapper with the
  family-standard one, reducing duplicate infra code (serves R2's "working
  POCs, not stubs" and reduces the Helix family's maintenance surface).
- If target's `quic-go` usage is already a mature, well-tested layer
  specific to the TOON/Brotli wire-format work (R2/addenda TOON note) →
  **STUDY-pattern-only**: read `helix_ota/server/internal/transport/transport.go`
  for the proven server-bootstrap pattern, keep target's own transport code.

**Recommendation:** STUDY-pattern-only for now; escalate to ADOPT after a
direct size/quality diff of target's transport code against
`http3/pkg/` (honest gap, not resolved in this pass — §7).

**Serves:** R2 (working POC quality bar). **Gated on G14/X1:** no.

**Transitive deps (§11.4.31, verified `helix_ota/submodules/http3/helix-deps.yaml`):**
leaf — `deps: []`, zero own-org dependencies.

---

## 6A. MCP_Module (R14) — STUDY-only, not a redirect

**What it is (verified `helix_code/submodules/mcp_module/README.md:1-40`,
module `digital.vasic.mcp`, Go 1.24):** a complete, generic MCP
implementation — JSON-RPC 2.0 server (stdio + HTTP/SSE) AND client (stdio +
HTTP/SSE), a thread-safe adapter `Registry` (stdio/Docker/HTTP adapter
types), config loading (JSON/YAML), i18n error-message seam. Status:
"FUNCTIONAL — all seven packages ship tested implementations,
`go test -race -count=1 ./pkg/...` green." Honors MCP protocol version
`2024-11-05`: `initialize`/`tools/list`/`tools/call`/`resources/list`/
`resources/read`/`prompts/list`/`prompts/get`.

**Target's current state (verified `project/internal/mcp/*.go`):** the
target does **NOT** hand-roll MCP from scratch — it already imports the
established third-party `github.com/mark3labs/mcp-go` SDK
(`internal/mcp/tools.go:9`, `internal/mcp/acp_adapter.go:12`) for the core
protocol types/dispatch, and builds its OWN value-add on top: a hardened
`HTTPTransport` that mounts MCP routes onto the **shared, already-secured**
API router (`internal/mcp/http_transport.go:15-24` — the file's own header
comment documents a real prior security bug it fixed: a second listener on
the same host:port raced the API server's listener, and whichever won the
bind decided the live CORS/auth posture; the fix collapses both onto one
listener), a custom `ACPAdapter` (JSON-RPC-over-stdio translation between
ACP and MCP method names — `acp_adapter.go:16-24`), and project-specific
tool/prompt definitions (`tools.go`, `prompts.go` — the 7 skill-graph tools
per `SPEC.md §7`).

**Verdict (FACT, not a guess):** `MCP_Module` (`digital.vasic.mcp`) and
target's `mark3labs/mcp-go`-based layer are **two independent, complete
implementations of the same protocol spec.** Target's is already working,
already has a captured security fix (the dual-listener CORS/auth bug), and
already has a custom ACP adapter `MCP_Module` does not document having.
Replacing a working, hardened implementation with a second complete
implementation for no new capability is **not** what §11.4.74
extend-don't-reimplement asks for — that clause exists to stop *reimplementing
from scratch*, and target did not: it built on an established upstream SDK
already. **STUDY-only**, not ADOPT, not EXTEND: `MCP_Module`'s
`pkg/adapter` (stdio/Docker/HTTP adapter *registry*, for proxying to OTHER
external MCP servers rather than serving your own tools) is worth reading
if/when the skill-graph ever needs to act as an MCP **client** fanning out
to multiple upstream MCP servers — a capability neither `mcp-go` nor
target's current code appears to need today (UNCONFIRMED whether target
ever needs this — no such requirement found in R1-R21).

**Serves:** R14 (MCP tool-trigger surface) — already served by the current
stack; no redirect. **Gated on G14/X1:** no (no adoption recommended).

**Transitive deps (§11.4.31, verified `mcp_module/helix-deps.yaml`):**
leaf — `deps: []`, own-org deps "none (genuine leaf)".

---

## 6B. DebateOrchestrator (G05 LLM jury) — ADOPT-core + EXTEND-stubs

**What it is (verified `helix_code/submodules/debate_orchestrator/README.md:1-45`,
module `digital.vasic.debate`, Go 1.26):** "multi-agent debate orchestration
primitives consumed by HelixAgent. The orchestrator coordinates multi-LLM
consensus + dissent across configurable agent pools, captures real
wall-clock latency per agent response, propagates `ProviderInvoker` errors
explicitly (no silent absorption)." Status: FUNCTIONAL — 14 packages
compile, `go test -race ./...` green.

**Honest package-tier table is itself load-bearing evidence (README, not
this survey's own claim):** REAL — `debate` root (`LessonBank` CRUD),
`agents`, `topology`, `gates` (permissive baseline), `orchestrator` (core
`Orchestrator.ConductDebate`, `AgentPool`, `APIAdapter`, session lifecycle,
`ProviderInvoker` wiring), `comprehensive` (real core, stubbed streaming).
STUB (constructors real, execution methods return explicit
`NotYetImplemented` errors, never silent success) — `validation`, `audit`,
`evaluation`, `reflexion`, `tools`.

**Relevance to the skill-graph's own G05 LLM jury:** target's own design
(cited in `research/helix_interop_incorporation.md §4`) already plans a
"multi-model **validation jury** (SPEC §8 `[validation] jury_size=3`)" for
skill-content correctness. `DebateOrchestrator`'s REAL core
(`Orchestrator.ConductDebate` across an `AgentPool` with per-agent latency
capture + explicit error propagation) is a **direct** structural match for
"N models vote/debate on whether a generated skill is correct" — exactly
the jury's job.

**The honest catch:** the packages the jury most needs for a **pass/fail
verdict** — `validation` (`Execute`) and `evaluation` (`RunBenchmark`) — are
the ones still marked STUB (`NotYetImplemented`). Adopting
`DebateOrchestrator` as-is gives a working multi-agent debate mechanism +
latency capture + lesson-bank persistence, but the skill-graph project would
need to **implement the STUB `validation`/`evaluation` packages itself**
(fork-and-extend upstream per §11.4.74) to get a structured jury verdict out
the other end — it is not a drop-in complete jury today.

**Verdict: ADOPT-as-submodule** (`orchestrator` + `debate` core, real,
tested) **+ EXTEND-upstream** the `validation`/`evaluation` STUB packages
with the skill-graph's own pass/fail scoring logic, contributing the
implementation back per §11.4.74 rather than forking silently. This
directly reduces the work of G05's jury design — the multi-agent
coordination + `ProviderInvoker` wiring (which composes cleanly with §4's
`LLMProvider` adoption — `ProviderInvoker` is the exact seam an
`llmprovider.Provider` adapter would plug into) does not need to be written
from scratch.

**Serves:** R8/R11 (jury correctness verdicts), G05. **Gated on G14/X1:**
no (new dependency, follow the chosen layout convention).

**Transitive deps (§11.4.31, verified `debate_orchestrator/helix-deps.yaml`):**
leaf — `deps: []`, own-org deps "none".

---

## 6C. token_optimizer (§11.4.141) — ALREADY VENDORED, needs wiring not adding

**Critical finding: this dependency is not missing — it is already present.**
`token_optimizer` (module `github.com/vasic-digital/token_optimizer`) is
already vendored inside the **target's own** constitution submodule tree at
`constitution/submodules/token_optimizer/` (verified
`ls constitution/submodules/` → `clickup_sync`, `continuum`,
`session_orchestrator`, `token_optimizer`) — landed there via the §11.4.28(C)
depth-1 reusable-engine carve-out, the same mechanism that landed
`continuum` (§11.4.207). The target does not need to add this dependency at
all; it needs to **use** it.

**What it does (verified `constitution/submodules/token_optimizer/README.md:1-25`,
Revision 2, 2026-07-08):** "a project-agnostic Go engine that minimizes the
token / cost / byte footprint of LLM request pipelines — tier routing with a
never-downgrade floor, a multi-layer cache, **shape-routed wire encoding**,
and telemetry." Package plan: `pkg/config` (tier/pricing/threshold registry
— **implemented**), `pkg/pipeline` (`Optimize(ctx, Request)` orchestrator),
`pkg/router` (tier decision + never-downgrade floor + failover), `pkg/cache`
(exact/semantic/artifact cache layers), `pkg/wire` (shape-routed encoder —
**`min(TOON, compactJSON)`, never-worse guard**).

**Direct relevance to the target project, not just to §11.4.141 generally:**
`pkg/wire`'s `min(TOON, compactJSON)` shape-routed encoder is the **exact
mechanism** the target's own R2 addenda already require — TOON as the
primary wire format with JSON fallback
(`REQUIREMENTS.md:32-42`, "TOON primary + JSON fallback … needs a Go TOON
codec"). This is not a generic §11.4.141 nicety for the skill-graph; it is
a load-bearing piece of the project's OWN stated wire-format requirement,
already built, already vendored, currently unused.

**Verdict: WIRE IT, do not re-add it and do not hand-roll a
TOON-vs-JSON size comparator.** Route the P3 `ModelProvider` chain's LLM
request/response bodies AND the REST/MCP API's TOON/JSON content negotiation
through `token_optimizer`'s `pkg/wire` + `pkg/pipeline` rather than a
bespoke implementation, satisfying §11.4.141 (heavy token-optimization,
always-on default per §11.4.198) directly from a component the project
already owns.

**Serves:** §11.4.141, §11.4.198, and directly the R2 TOON/JSON wire-format
requirement. **Gated on G14/X1:** no — already present, no vendoring
decision needed, only a wiring task.

**Transitive deps:** not independently re-derived in this pass (already
inside the constitution submodule's own §11.4.28(C) carve-out, which
mandates zero own-org deps for anything landed there — consistent with
`continuum`'s equivalent carve-out entry).

---

## 7. HelixLLM (R4) + Helix-Track/Core (R4 interop) — one new data point

`research/helix_interop_incorporation.md` (Revision 1) already fully
resolved HelixLLM's binding contract (`HELIX_LLM_LOCAL_OPENAI_ENDPOINT`,
default `http://localhost:18434`, OpenAI-compatible, no trailing `/v1`) and
concluded HelixTrack/Core should be integrated as a **remote REST+JWT peer
service**, never vendored, because no Go client import is currently needed.

**New cross-check from this survey:** `helix_terminator/.gitmodules:22-24`
vendors `Helix-Track/Core` as a **build dependency** —
`submodules/helixtrack-core/` (`git@github.com:Helix-Track/Core.git`) — a
different integration pattern than the target's own design. This is a real,
alternative precedent: if the skill-graph service ever needs to import
Core's Go domain types (`Ticket`, `TicketStatus`, etc.) directly rather than
calling its REST API, `helix_terminator`'s vendored-submodule pattern is the
one to follow. Until that need is proven, the REST+JWT peer-service design
already in `helix_interop_incorporation.md` §5 remains correct and is NOT
superseded by this finding.

**Recommendation:** unchanged (REST+JWT peer service); `helix_terminator`'s
vendoring pattern noted as the fallback if a Go-import need is later
proven.

**Serves:** R4. **Gated on G14/X1:** no (HelixTrack/Core stays a service,
not a vendored dependency, under the current design).

---

## 8. helix_track client matrix (R3/G15) — real precedent, with a major caveat

**Real structure** (`helix_track/.gitmodules:87-131`): one submodule **per
platform** — `Core` (Go backend), `Web-Client` → `web_client/`,
`Desktop-Client` → `desktop_client/`, `Android-Client` → `android_client/`,
`iOS-Client` → `ios_client/`, `Harmony-OS-Client` → `harmony_os_client/`,
`Aurora-OS-Client` → `aurora_os_client/` (note the repo itself is named
`Aurora-Client`), plus `Screensaver`.

**Verified per-client reality (not assumed from names):**

- **Web-Client:** real, substantial Angular 19 app (standalone components,
  Angular Material, RxJS, "custom HTTP/3 QUIC service", WebSocket;
  `helix_track/CLAUDE.md:173-189`). Has an explicit intra-client shared
  layer: `src/app/shared/` — "Shared/reusable components"
  (`CLAUDE.md:183`), but this is shared **within** the Angular app, not
  across platforms.
- **Desktop-Client:** **Tauri + Angular** (`CLAUDE.md:210` "Desktop Client
  (Tauri + Angular)") — real, populated directory
  (`angular.json`, `cypress/`, `cypress.config.ts`, implementation docs).
  This is a genuine cross-surface shared-code win: Desktop **wraps the same
  Angular Web-Client codebase** in a Tauri webview rather than being a
  separate rewrite — directly matches the "Tauri (web-view) reuses the Web
  `tokens.css`" pattern already identified independently in
  `opendesign_incorporation.md §3.3`, cross-confirming Tauri-wrapping-Web as
  the correct Desktop strategy for R3.
- **Android-Client:** real, populated **native Kotlin/Gradle** project
  (`build.gradle`, `app/`, `gradlew`, `gradle.properties` all present on
  disk) — not a stub, not Flutter.
- **iOS-Client:** real, populated **native Swift** project (`HelixTrack/`
  Xcode-project directory, `e2e-tests.sh`, `ai-qa-runner.js` present) — not
  a stub, not Flutter.
- **Harmony-OS-Client:** **EMPTY governance stub.** Directory contains only
  `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `QWEN.md`, `README.md`,
  `helix-deps.yaml` — the README is a single line: "HelixTrack Harmony OS
  Client." **Zero source code checked out.**
- **Aurora-OS-Client:** **EMPTY governance stub**, identical shape to
  Harmony — README: "HelixTrack Aurora Client." **Zero source code.**

**Verdict, stated as FACT not hedged:** helix_track's real, working
precedent for maximizing shared code is **"native per-platform clients +
Tauri wraps the Web app for Desktop,"** NOT a single cross-platform
framework. helix_track uses **no Flutter, no Kotlin Multiplatform, no
Compose Multiplatform anywhere** in its client matrix. This is a direct data
point against the target's own prior shared-core research finding
("GUI apps = Flutter — the ONLY framework covering Android + iOS +
HarmonyOS + Aurora OS + desktop from one codebase," `REQUIREMENTS.md:169-171`)
— the family's only real precedent for Desktop-code-sharing uses Tauri+Web,
not Flutter, and offers **zero** validated precedent for Flutter (or any
other single framework) on Android/iOS/HarmonyOS/Aurora, because
Harmony/Aurora are unbuilt everywhere in the family, including here.

**Cross-reference to `research/g15_aurora_harmonyos_client_feasibility.md`:**
that document independently found (via external web research, not this
sibling survey) that Flutter has **no official** HarmonyOS-NEXT target and
only a community/SIG-maintained, explicitly-beta fork
(`openharmony-sig/flutter_flutter`) covers it. This survey's finding
**reinforces** that risk from a different angle: not only is upstream
Flutter-on-HarmonyOS unofficial, but **no sibling project in the entire
Helix family has ever shipped it** — there is no internal precedent to fall
back on if the external fork proves unreliable.

**Recommendation:**
- **ADOPT the pattern** "Desktop = Tauri wrapping the Web client" for R3 —
  strong, working, cross-confirmed precedent (helix_track + the
  OpenDesign token-consumption research already assumed it independently).
- **STUDY-pattern-only** the overall "one submodule per platform, `Core` as
  a separate Go backend" repository organization — directly transferable to
  how target could lay out its own future client submodules, IF it chooses
  the per-platform-native strategy over Flutter.
- **NOT-applicable** as a validated reference for HarmonyOS/Aurora
  specifically — helix_track has shipped nothing there. The target's own
  `g15_aurora_harmonyos_client_feasibility.md` externally-researched
  findings remain the only real lead for those two platforms.

**Serves:** R3, G15. **Gated on G14/X1:** no (client architecture, not a
build-dependency layout question).

---

## 8A. Flagged items (R22 — name + one-liner only, full design deferred)

Per the operator's explicit instruction, these are named + given a single
applicability line and their real repo existence is confirmed (checked-out
copy or live `gh` lookup) — none are fully designed here; each needs its own
follow-on pass.

- **`HelixDevelopment/DagOrchestrator`** (checked out
  `helix_code/submodules/dag_orchestrator/`, module `dev.helix.dag`, real —
  "generic, reusable, agent-free pure-data DAG scheduler … dispatches the
  ready-set in topological order onto a bounded worker pool … dynamic node
  expansion," 8 passing race tests). **Applicability:** direct candidate for
  G11 worker scheduling AND the recursive skill-DAG's ready-set dispatch
  (R2/R16) — a follow-on pass should design its wiring against
  `research/g11_worker_design.md` and `research/g06_g07_skilltree_dag_design.md`.
  Leaf (`helix-deps.yaml: deps: []`); its own README notes it composes with
  `vasic-digital/Concurrency` (worker-pool substrate) when the consumer wants
  one, added by the consumer, not nested inside it.
- **`HelixDevelopment/PipelineRuntime`** (checked out
  `helix_code/submodules/pipeline_runtime/`, module `dev.helix.pipeline`,
  real — "staged streaming dataflow runtime: composable push-based operators
  … plus a flow-based-programming Component/Port/Network layer with bounded
  backpressure," 5 passing race tests). **Applicability:** candidate for the
  R6 wizard's create→map→process pipeline and/or R14's real-time
  auto-growth job pipeline (`Map`/`Filter`/`FlatMap`/`Window`/`MapReduce`
  operators with backpressure map directly onto "kick off async generation +
  research jobs and stream progress," `REQUIREMENTS.md:109-111`). Leaf
  (zero own-org deps).
- **`vasic-digital/VisionEngine` + `vasic-digital/Panoptic`** (both checked
  out in `helix_code/submodules/`). VisionEngine
  (`helix_code/submodules/vision_engine/README.md:1-10`): "Computer-vision +
  LLM-vision toolkit for UI analysis and navigation-graph construction" —
  ships an `Analyzer`/`VideoProcessor` interface + UI-element/screen-diff
  value types; already a declared transitive dep of `HelixQA` itself
  (`helixqa/helix-deps.yaml`: "Vision/OCR analysis engine for video-evidence
  validation"). Panoptic (`helix_code/submodules/panoptic/README.md`, MIT,
  Go 1.21+) is a broader AI-enhanced test-generation framework (real
  `ai_test_generation.yaml`/`ai_implementation_test.yaml` fixtures observed
  on disk). **Applicability:** both are candidates for §11.4.170's
  device-independent host-rendered visual-proof mandate once the target's
  R3 Web/Desktop/Mobile clients exist — VisionEngine for the OCR/vision
  layout oracle itself, Panoptic for AI-driven test-generation around it. No
  R3 client exists yet in this project, so neither is actionable before P8.
- **KMP suite** — `vasic-digital/UI-Components-KMP`, `Security-KMP`,
  `Concurrency-KMP`, `Storage-KMP` (none checked out in any surveyed
  sibling; **confirmed to genuinely exist** via live `gh repo view`, §10A —
  "Kotlin Multiplatform UI components: theme system, animations,
  accessibility utilities for Compose"; "KMP secure storage: AES encryption,
  platform Keychain/KeyStore integration"; "KMP concurrency utilities: lazy
  loading, platform synchronization, flow-based loaders"; "KMP network
  storage service interfaces and abstractions" respectively).
  **Applicability:** directly relevant to R3's Android/iOS/Desktop clients
  IF the target ever adopts Kotlin Multiplatform instead of (or alongside)
  Flutter for shared mobile/desktop code — a real, previously-unconsidered
  alternative to the Flutter-only lean in `REQUIREMENTS.md:169-171`, and
  notably a framework family with **actual reusable Helix-org building
  blocks already published**, unlike Flutter-on-HarmonyOS (§8, no family
  precedent at all). Not designed further here — a follow-on pass should
  weigh KMP-vs-Flutter for R3/G15 explicitly, now that KMP has confirmed
  real building blocks and Flutter-on-Harmony has none.
- **`vasic-digital/config` / `security` / `Storage` / `Filesystem`** (all
  checked out in `helix_code/submodules/`, each a generic reusable Go
  module — config-loading, security/crypto helpers, storage abstractions,
  filesystem abstractions respectively; not individually read beyond
  directory-listing confirmation in this pass). **Applicability:** generic
  infrastructure candidates for whichever of the target's own `internal/`
  packages currently hand-roll config loading / file I/O / storage
  abstraction — not yet diffed against target's actual code, so no
  ADOPT/EXTEND verdict is given here (honest gap, §11).

---

## 9. Other reusable practices

- **Constitution wiring (100% consistent):** all 7 surveyed siblings AND the
  target vendor `HelixDevelopment/HelixConstitution.git` at the identical
  path `constitution/`. HEAD commits differ across every project (expected —
  independent `git fetch`/pull timing per §11.4.26, not a defect):
  `helix_skills` `ba06f72e`, `helix_terminator` `e6504c27`, `helix_ota` /
  `helix_code` `7295a189` (these two happen to match — likely pulled at the
  same time), `helix_track` `e60cbde2`, `helix_translate` `2c338a4e`,
  `helix_vpn` `e7ff3c27`. No action required beyond the standing
  §11.4.26/§11.4.37 fetch-before-edit discipline; nothing here indicates the
  target is stale relative to any specific mandate.
- **`vasic-digital/Auth`** (module `digital.vasic.auth`): vendored by
  `helix_terminator` (`submodules/auth/`), `helix_track` (`auth`),
  `helix_code` (`submodules/auth`). Ships `pkg/jwt`, `pkg/apikey`,
  `pkg/oauth`, `pkg/middleware`, `pkg/token`, `pkg/i18n`
  (`helix_terminator/submodules/auth/README.md:1-10`). Target's
  requirements currently mandate only **security fixes** (CORS
  reflect-origin, `api_key`-in-query-then-logged — `REQUIREMENTS.md:121-122`),
  not a new auth subsystem. **STUDY-pattern-only now**; promote to
  ADOPT-as-submodule the moment target's REST/MCP surface needs bearer-token
  or API-key auth (not yet a stated requirement).
- **`helix-deps.yaml` as the dependency manifest format (§11.4.31):**
  universally present — not just at project roots but inside individual
  vendored submodules too (`llmprovider/helix-deps.yaml`,
  `containers/helix-deps.yaml`, `docs_chain/helix-deps.yaml` all exist and
  are leaves). Target's own `project/helix-deps.yaml` already follows this
  schema correctly — no gap.
- **Naming/layout convention (§11.4.29):** every sibling uses
  `submodules/<snake_case>/` for grouped deps or a bare root `<name>/` for
  ungrouped ones — consistent with target's manifest; no gap found.

---

## 10. Prioritized adoption table

| # | Subsystem | Recommendation | Requirement(s) served | Gated on G14/X1? |
|---|---|---|---|---|
| 1 | LLMProvider (`vasic-digital/LLMProvider`) | **ADOPT-as-submodule** `submodules/llmprovider/` — redirects R7/R19 away from hand-rolled Anthropic/OpenAI clients | R7, R19, R17 (resilience) | No (new dep; follow chosen layout convention) |
| 2 | Containers (`vasic-digital/containers`) | **ADOPT-as-submodule**, ROOT path `containers/` (per `helix_ota`'s live-code precedent) | R20, R9, R8 (HelixQA dependency) | **YES** |
| 3 | HelixQA (`HelixDevelopment/HelixQA`) | **ADOPT-as-submodule** (already declared correctly) + add `containers` dependency | R8, R11 | Partially (containers half) |
| 4 | Docs Chain (`vasic-digital/docs_chain`) | Incorporation mechanics already designed; **reconsider** `layout: grouped` vs. family-majority root-path convention | R10 | Layout choice, yes |
| 5 | open-design (`nexu-io/open-design`) | **ADOPT-as-submodule** `submodules/open-design/` — already fully designed, path now corroborated by 2 siblings | R12 | No |
| 6 | HelixLLM (`HelixDevelopment/HelixLLM`) | **ADOPT-as-submodule** (already declared, fully designed) | R4, R7 | No |
| 7 | Helix-Track/Core | **KEEP** as REST+JWT peer service (design unchanged); `helix_terminator`'s vendored-submodule pattern noted as fallback if Go-import need arises | R4 | No |
| 8 | http3 (`vasic-digital/http3`) | **STUDY-pattern-only**, pending size/quality diff vs. target's own `quic-go` transport | R2 | No |
| 9 | helix_track client matrix (Desktop=Tauri-wraps-Web pattern) | **ADOPT the Tauri-wraps-Web pattern** for Desktop; **STUDY-pattern-only** the one-submodule-per-platform layout; **NOT-applicable** for Harmony/Aurora validation | R3, G15 | No |
| 10 | `vasic-digital/Auth` | **STUDY-pattern-only** — promote to adopt only when an auth requirement exists | (none yet — future R?) | No |
| 11 | Constitution wiring | Already fully consistent — no action | (governance) | No |
| 12 | `vasic-digital/design_system` | **ADOPT-as-submodule** alongside open-design — ready CSS/Tailwind-v4/light+dark tokens, extracted from a real shipped Helix site | R12 | Layout choice, yes |
| 13 | `vasic-digital/MCP_Module` | **STUDY-only** — target's `mark3labs/mcp-go`-based, already-hardened MCP layer is not redirected | R14 | No |
| 14 | `HelixDevelopment/DebateOrchestrator` | **ADOPT-core + EXTEND-stubs** (`validation`/`evaluation` packages) for the G05 jury | R8, R11, G05 | No |
| 15 | `vasic-digital/token_optimizer` | **ALREADY VENDORED** in `constitution/submodules/` — WIRE, do not add; its `pkg/wire` directly implements R2's TOON/JSON shape-routing | §11.4.141, §11.4.198, R2 | No |
| 16 | `HelixDevelopment/DagOrchestrator` (FLAGGED) | Named candidate for G11 worker scheduling — full design deferred | G11 | No |
| 17 | `HelixDevelopment/PipelineRuntime` (FLAGGED) | Named candidate for R6 wizard / R14 job pipeline — full design deferred | R6, R14 | No |
| 18 | `VisionEngine` + `Panoptic` (FLAGGED) | Named candidates for §11.4.170 client visual proof — not actionable before R3 clients exist | §11.4.170 | No |
| 19 | KMP suite (FLAGGED) | Confirmed real; named as a Flutter alternative for R3/G15 — full weigh-off deferred | R3, G15 | No |
| 20 | `config`/`security`/`Storage`/`Filesystem` (FLAGGED) | Named generic-module candidates — not diffed against target's own code | (various `internal/`) | No |

---

## 10A. Method note — live `gh` lookups (R22)

Five repos named in the R22 mandate were not present as checked-out
submodules in any of the 6 surveyed siblings. Rather than mark all five
`UNCONFIRMED: not found` outright, this survey ran targeted, read-only
`gh repo view <org>/<repo> --json name,description` and
`gh api repos/<org>/<repo>/contents` / `.../readme` calls (no clone, no
write) to confirm real existence + read top-level structure/README where
the repo turned out to be real:

- `vasic-digital/UI-Components-KMP` — **EXISTS** (confirmed via `gh repo
  view`, description read).
- `vasic-digital/Security-KMP` — **EXISTS** (confirmed via `gh repo view`).
- `vasic-digital/Concurrency-KMP` — **EXISTS** (confirmed via `gh repo
  view`).
- `vasic-digital/Storage-KMP` — **EXISTS** (confirmed via `gh repo view`).
- `vasic-digital/design_system` — **EXISTS**, README + root directory
  listing read via `gh api` (content decoded from the API's base64
  payload) — see §5.
- `vasic-digital/DesignSystem` (capitalized variant, checked defensively in
  case of a naming-convention mismatch) — **DOES NOT EXIST** (`gh` returned
  "Could not resolve to a Repository"). The correct, real name is the
  lowercase `design_system`.

This is the only section of the survey relying on live network reads rather
than local file inspection; every fact drawn from it is labelled as such
inline (§5, §8A).

---

## 11. Honest gaps (§11.4.6 — nothing here is silently assumed)

1. **LLMProvider org duplication UNCONFIRMED:** `vasic-digital/LLMProvider`
   and `HelixDevelopment/LLMProvider` were not diffed against each other —
   both declare module `digital.vasic.llmprovider`, but content parity
   between the two org-owned copies is not independently proven in this
   survey.
2. **http3 vs. target's own `quic-go` code was not diffed.** The task brief
   states target "uses quic-go directly," but this survey did not open
   target's own transport source to compare size/maturity against
   `helix_ota/submodules/http3/pkg/` — the STUDY-vs-ADOPT call in §6 is
   provisional on that comparison.
3. **`helix_ota`'s duplicate `containers` + `submodules/containers`
   mounts were not reconciled** — both are populated on disk; this survey
   identified the root one as load-bearing (via `server/go.mod`'s `replace`
   directive) but did not determine whether the `submodules/` copy is dead
   weight or used by some other part of `helix_ota` not read here.
4. **HelixQA's own top-level orchestration entry point is UNCONFIRMED to
   exist anywhere in the family** — no sibling's `scripts/`/`Makefile` was
   found invoking `helixqa`/`helixqa-bank-session` against a project-owned
   bank; the "how you actually trigger a QA session end-to-end" step is
   inferred from the vendored binary's existence, not observed in a working
   invocation script.
5. **`Helix-Track/Core`'s vendored-submodule usage inside `helix_terminator`
   was not inspected** — only its `.gitmodules` entry was read; whether
   `helix_terminator` actually imports Core's Go types anywhere, or merely
   declared the submodule without wiring it, is UNCONFIRMED.
6. **`helix_terminator`'s own wiring code for `llmprovider`/`auth`/`helixllm`
   was not found** — its README (`helix_terminator/README.md`) is
   constitution-inheritance boilerplate only; no `go.mod`/source file was
   located at its root proving these submodules are actually imported by
   `helix_terminator`'s own build (only their *presence* as populated
   checkouts is confirmed, not their *consumption*). The load-bearing
   evidence for real Go-level wiring in this survey comes from `helix_ota`,
   not `helix_terminator`.
7. **`helix_code`'s own use of `llm_provider`, `helix_llm`, `helix_agent`
   was not re-derived here** — fully covered already in
   `helix_interop_incorporation.md`; not repeated.
8. **Constitution HEAD divergence across siblings is not itself
   investigated for cascade-gap risk** — noted as expected/benign per
   §11.4.26, but no check was run for whether any sibling is missing a
   specific recent anchor (e.g. §11.4.202/.207/.208/.209) that the target
   already has at `ba06f72e`.
9. **(R22) `vasic-digital/design_system`'s `helix-deps.yaml` was not
   conclusively read** — `gh api` returned no content for that path; this
   could mean the file is genuinely absent (a leaf) or the API call
   404'd/failed silently. Its transitive-dependency status is UNCONFIRMED,
   not asserted as zero, unlike the other DECISIVE items in this survey
   whose leaf status was positively confirmed from a real `helix-deps.yaml`
   body.
10. **(R22) `DebateOrchestrator`'s `ProviderInvoker` interface shape was not
    read line-by-line** — the recommendation that it composes cleanly with
    an `llmprovider.Provider` adapter (§6B) is an architectural inference
    from the README's description ("propagates `ProviderInvoker` errors
    explicitly"), not a verified interface-signature match. A follow-on
    pass implementing §6B must read `orchestrator/*.go` directly before
    coding the adapter.
11. **(R22) the KMP suite repos (`UI-Components-KMP`, `Security-KMP`,
    `Concurrency-KMP`, `Storage-KMP`) were confirmed to exist and their
    one-line descriptions were read, but their actual source/API surface
    was NOT inspected** (no `gh api contents` call was made for them,
    unlike `design_system`) — §8A's "Flutter alternative" framing is based
    on the repo names + one-line descriptions only; a follow-on KMP-vs-
    Flutter weigh-off must read their real contents first.
12. **(R22) `config`/`security`/`Storage`/`Filesystem` (the generic-module
    FLAG-list entries) were confirmed present on disk in `helix_code` by
    directory listing only** — no README/go.mod was read for any of the
    four, and none was diffed against any of the target's own `internal/`
    packages. No ADOPT/EXTEND/STUDY verdict is given for these four; §8A
    explicitly defers that call.

---

## Sources verified (local, 2026-07-15)

- `.gitmodules` of `helix_terminator`, `helix_ota`, `helix_track`,
  `helix_code`, `helix_translate`, `helix_vpn`.
- `helix_ota/server/go.mod:32-33,66,68` (replace directives for
  `digital.vasic.containers`, `digital.vasic.http3`).
- `helix_ota/server/internal/store/postgres_integration_test.go:1-20`
  (real containers pkg/boot + pkg/compose usage).
- `helix_ota/submodules/helixqa/README.md:1-25`; `helixqa/cmd/` listing;
  `helixqa/banks/` listing.
- `helix_terminator/submodules/llmprovider/{go.mod,doc.go,README.md}`;
  `pkg/providers/{anthropic,claude,openai}/*.go` (heads read directly).
- `helix_track/.gitmodules`; `harmony_os_client/README.md`;
  `aurora_os_client/README.md`; `CLAUDE.md:160-230`; `android_client/`,
  `ios_client/`, `desktop_client/` directory listings.
- `helix_terminator/README.md` (full, 55 lines); `submodules/auth/README.md:1-10`.
- Target: `REQUIREMENTS.md`, `project/helix-deps.yaml`,
  `project/internal/autoexpand/llm.go:1-40`, `project/internal/mcp/*.go`
  (`http_transport.go`, `acp_adapter.go`, `prompts.go`, `tools.go`),
  `research/{helix_interop_incorporation,opendesign_incorporation,docs_chain_incorporation,g15_aurora_harmonyos_client_feasibility}.md`.
- `git rev-parse HEAD` for `constitution/` in all 7 projects.
- **(R22 additions):** `helix_code/submodules/mcp_module/{README.md,go.mod,helix-deps.yaml}`;
  `helix_code/submodules/debate_orchestrator/{README.md,go.mod,helix-deps.yaml}`;
  `helix_code/submodules/dag_orchestrator/{README.md,helix-deps.yaml}`;
  `helix_code/submodules/pipeline_runtime/{README.md,helix-deps.yaml}`;
  `helix_code/submodules/vision_engine/README.md:1-10`;
  `helix_code/submodules/panoptic/README.md` (head);
  `helix_ota/submodules/helixqa/helix-deps.yaml` (full, 8 transitive deps);
  `helix_ota/containers/helix-deps.yaml`;
  `helix_ota/submodules/http3/helix-deps.yaml`;
  `helix_terminator/submodules/llmprovider/helix-deps.yaml`;
  `<target-repo>/constitution/submodules/token_optimizer/README.md:1-25`;
  live `gh repo view vasic-digital/{UI-Components-KMP,Security-KMP,Concurrency-KMP,Storage-KMP,design_system}`
  + `gh repo view vasic-digital/DesignSystem` (confirmed non-existent) +
  `gh api repos/vasic-digital/design_system/{contents,readme}` (2026-07-15,
  read-only, no clone, no write).
