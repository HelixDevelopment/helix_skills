# R22 Full-Catalogue Incorporation Design — LLMProvider, http3, PipelineRuntime, token_optimizer, VisionEngine/Panoptic, KMP, design_system + open-design

**Revision:** 1
**Last modified:** 2026-07-15T18:10:00Z
**Status:** design-only follow-on to the R21/R22 survey (`research/helix_family_reusable_practices.md`
Revision 2). No code written, no existing file modified. This is the deep-dive pass that survey's
§8A "flagged, full design deferred" items, plus a redirect-grade design for §4's LLMProvider
critical finding.
**Scope:** resolves per-module incorporation verdicts (ADOPT-as-submodule / EXTEND-upstream /
STUDY-only / N-A) for: **LLMProvider**, **http3**, **PipelineRuntime**, **token_optimizer**,
**VisionEngine + Panoptic**, the **KMP component suite**, **design_system + open-design**. Per
§11.4.74 (catalogue-first, extend-don't-reimplement) and §11.4.28 (owned submodules are equal
codebase) and §11.4.31 (recursive `helix-deps.yaml` transitive-dependency manifests).
**Explicitly out of scope (already resolved elsewhere, NOT re-derived or re-litigated):**
`debate_orchestrator` (STUDY-only for G05) and `dag_orchestrator` (INCORPORATE-as-submodule,
scoped) — both resolved in `research/g05_g11_reuse_debate_dag_design.md`; `containers`, `HelixQA`,
`Docs Chain`, `MCP_Module`, `Auth`, `HelixLLM`/`Helix-Track/Core`, the helix_track client matrix —
all resolved in `research/helix_family_reusable_practices.md` §§1-3,6A,7-9; the §11.4.28(C)
single-canonical dependency-resolution policy itself (**Option A adopted** — see
`research/g14_x1_submodule_policy_decision.md` §5) — assumed and applied here, not re-decided.
**Method:** every claim is cited to a real file path read directly in this pass (module source,
`go.mod`, `helix-deps.yaml`, `.gitmodules`), or to a live, read-only `gh repo view` / `gh api`
lookup run in this pass (distinguished inline from the survey's own earlier `gh` reads).
Unread/unverified material is marked `UNCONFIRMED:` per §11.4.6. No `likely`/`probably`/`seems`
anywhere below — every verdict is either a cited FACT or an explicit honest gap.

---

## 0. Reading map

Sections 1–7 below give one verdict per module (the task's numbered list). Section 8 is the
consolidated incorporation-order table. Section 9 — **the single most important output** — is the
concrete R7/R19 provider-layer redirect design (the LLMProvider adapter seam, the
`NewLLMClientFromConfig` rewrite, and a brand-new G05 jury `LLMValidator` wiring the project does
not have today). Section 10 is honest gaps.

---

## 1. LLMProvider (R7/R19) — verdict: **ADOPT-as-submodule**

### 1.1 Org-canonical resolution under the already-adopted Option A (FACT, newly resolved here)

The survey (§4.1) flagged, but left `UNCONFIRMED`, which of the two org-owned copies —
`vasic-digital/LLMProvider` or `HelixDevelopment/LLMProvider` — is canonical. This pass resolves
**which copy Option A's `find_ecosystem_copy()` actually finds** (not full byte-content parity,
which remains unread — see honest gap §10.1): `helix_code/.gitmodules:61-64` —

```
[submodule "dependencies/HelixDevelopment/llm_provider"]
	path = submodules/llm_provider
	url = git@github.com:HelixDevelopment/LLMProvider.git
	branch = master
```

`helix_code` is the currently-configured single `--ecosystem-root` per the adopted Option A
(`g14_x1_submodule_policy_decision.md` §5). Therefore, **under Option A as it stands today, the
operationally canonical copy the target project would resolve to is `HelixDevelopment/LLMProvider`
at `helix_code/submodules/llm_provider/`** — not the `vasic-digital` copy `helix_terminator`/
`helix_ota`/`helix_translate`/`helix_vpn` vendor. Both copies declare the identical Go module path
`digital.vasic.llmprovider` (`helix_code/submodules/llm_provider/go.mod:1`, verified directly), so
this is not a build break either way — but it settles which physical checkout the sync tooling
touches for this project.

### 1.2 What it is (verified directly, `helix_code/submodules/llm_provider/`)

- **Core interface** (`doc.go:16-23`): `Complete(ctx, *models.LLMRequest) (*models.LLMResponse,
  error)`, `CompleteStream(...) (<-chan *models.LLMResponse, error)`, `HealthCheck() error`,
  `GetCapabilities() *models.ProviderCapabilities`, `ValidateConfig(map[string]interface{}) (bool,
  []string)`.
- **`models.LLMRequest`** (`pkg/models/types.go:10-27`): `Prompt string`, `Messages []Message`,
  `ModelParams ModelParameters` (`Model`, `Temperature`, `MaxTokens`, `TopP`, `StopSequences`,
  `ProviderSpecific`), `Tools []Tool`, `ToolChoice interface{}`, plus session/audit fields not
  needed here (`SessionID`, `UserID`, `Status`, timestamps).
- **`models.LLMResponse`** (`pkg/models/types.go:44-59`): `Content string` (confirms the survey's
  §4.3 adapter sketch `resp.Content` — verified directly, not inferred), `TokensUsed int`,
  `FinishReason string`, `Confidence float64`, `ToolCalls []ToolCall`.
- **`pkg/providers/anthropic/anthropic.go`**: `NewProvider(apiKey, baseURL, model string) *Provider`
  (`:145`) and `NewProviderWithRetry(apiKey, baseURL, model string, retryConfig RetryConfig)
  *Provider` (`:150`); `(*Provider).Complete(ctx, *models.LLMRequest) (*models.LLMResponse, error)`
  (`:203`); `CompleteStream` (`:231`); `HealthCheck` (`:323`); `GetCapabilities` (`:344`).
- **`pkg/providers/openai/openai.go`**: identical constructor/method shape —
  `NewProvider(apiKey, baseURL, model string) *Provider` (`:130`), same five interface methods.
- **`pkg/providers/claude/claude.go`**: a *third*, OAuth-aware variant —
  `NewClaudeProviderWithOAuth(baseURL, model string, credReader OAuthCredentialReader)
  (*ClaudeProvider, error)` (`:209`) that explicitly detects and rejects
  `"This credential is only authorized for use with Claude Code"` (`:38`). This is **not** the
  provider R19's server-side Anthropic-API-key use case needs (the server never runs as the Claude
  Code CLI) — it is flagged here as a **distinct**, future-relevant building block for the
  project's own P3 toolkit-alias `ModelProvider` chain (`research/helix_interop_incorporation.md`),
  not consumed by the R19 redirect in §9 below.
- **Resilience primitives** (verified present, not re-derived): `circuit_breaker.go`,
  `health_monitor.go`, `retry.go` (exponential backoff + jitter + HTTP-status-aware retry) at the
  module root, composable over any `Provider` via `NewDefaultCircuitBreaker(name, provider)`
  (`doc.go` example, `:44`).

### 1.3 Transitive deps (§11.4.31, verified `helix_code/submodules/llm_provider/helix-deps.yaml`)

```yaml
schema_version: 1
deps: []
```
Leaf — zero own-org dependencies. No nested-chain risk under §11.4.28(C).

### 1.4 Incorporation seam

**Mount:** consumed BY REFERENCE from `helix_code/submodules/llm_provider/` under Option A — the
target project adds `llm_provider` to `project/helix-deps.yaml` and vendors **nothing new locally**
(identical mechanism already exercised for `containers`/`HelixQA`, per `g14_x1_submodule_policy_decision.md`
§4). **Does not supersede `*OpenAILLM` by file-deletion** — it supersedes the *plan* to hand-roll a
second Anthropic client (R19's `AnthropicLLM`, §2 design not-yet-implemented); the concrete
adapter/redirect design is in §9.

**Serves:** R7 (directly — a real 50+-provider registry behind one interface, satisfying "pluggable
ModelProvider, not hardcoded OpenAI" far more completely than a second hand-rolled client), R19
(directly — supersedes `AnthropicLLM`), R17 (closes circuit-breaker/health-monitor/retry gaps with
an off-the-shelf, already-tested implementation).

**Gated on G14/X1:** not directly gated (a new dependency addition, not a contested existing one),
but resolves under the already-adopted Option A with **zero new vendoring** exactly like
`containers`.

---

## 2. http3 — verdict: **ADOPT-as-submodule** (escalated from the survey's provisional STUDY-only)

### 2.1 Resolving the survey's honest gap #2 (the diff it flagged as not yet done)

The survey (§6, honest gap §11.2) explicitly deferred STUDY-vs-ADOPT pending "a direct size/quality
diff of target's transport code against `http3/pkg/`." This pass performs that diff.

**Target's own transport (`project/internal/api/http3.go`, read in full — 85 lines):**
`setupHTTP3(handler http.Handler) error` hand-builds a `tls.Config` (`MinVersion: TLS13`), a
`quic.Config` (`MaxIdleTimeout`, `HandshakeIdleTimeout`, `MaxIncomingStreams`,
`MaxIncomingUniStreams`, `Allow0RTT`), wraps them in a raw `github.com/quic-go/quic-go/http3.Server`,
starts it in a goroutine with ad-hoc error logging, and `shutdownHTTP3(ctx) error` closes it with a
5-second timeout via a `done chan error` select. **Zero dedicated test file exists** for this logic
— confirmed by directory listing (`internal/api/*http3*` → only `http3.go`, no `http3_test.go`).

**`vasic-digital/http3` (`helix_ota/submodules/http3/`, module `digital.vasic.http3`, read
directly):** `pkg/server/server.go` (204 lines) exposes the identical shape — a `Config` struct
(`Addr`, `Handler`, `TLSConf`, `QUICConfig *quic.Config`) with its own `Validate()`, `New(cfg)
(*Server, error)`, `(*Server).Start() error`, `(*Server).Shutdown(ctx) error`, `(*Server).Done()
<-chan error`, `(*Server).Addr() string` — narrower and more defensive (a dedicated `Validate`
step; a `Done()` channel the caller can select on instead of only a context-timeout race). It ships
**seven** dedicated test files: `server_test.go`, `server_branches_test.go`, `integration_test.go`,
`stress_test.go`, `fuzz_test.go`, `chaos_test.go`, `challenge_test.go` — the exact
stress/chaos/fuzz coverage §11.4.85 mandates for every fix/improvement, already paid for.

**Verdict (FACT, not a guess):** target's transport code is thin and duplicative of what `http3`
already wraps — same TLS/QUIC-config shape, same start/shutdown lifecycle, with **zero** of
`http3`'s test coverage. This resolves the survey's conditional §6 logic to its **ADOPT** branch.

### 2.2 Transitive deps (verified `helix_ota/submodules/http3/helix-deps.yaml`)

Leaf — `deps: []`, zero own-org dependencies (matches the survey's own citation, re-confirmed here).

### 2.3 Incorporation seam

Replace `internal/api/http3.go`'s `setupHTTP3`/`shutdownHTTP3` with a thin call into
`digital.vasic.http3/pkg/server`: `server.New(server.Config{Addr: addr, Handler: s.router, TLSConf:
tlsConfig, QUICConfig: quicConfig})` then `.Start()` in the existing goroutine pattern, `.Shutdown(ctx)`
in place of the hand-rolled `done`-channel select (the module's own `Done()` channel replaces it).
`s.router` (the shared, already-hardened `gin.Engine` per `internal/mcp/http_transport.go`'s own
collapsed-single-listener fix) is passed through unchanged — no coupling to the target's routing or
auth stack is introduced. Zero change to `go.mod`'s existing `github.com/quic-go/quic-go` require
(the module wraps the same upstream, does not replace it).

### 2.4 A genuinely new finding: this dependency does **not** resolve under the currently-configured single ecosystem root

Unlike `LLMProvider`/`containers`/`PipelineRuntime`/`VisionEngine`/`Panoptic` (all present under
`helix_code/submodules/`), **`http3` is not checked out anywhere under `helix_code`** — confirmed by
directory listing (`helix_code/submodules/` has no `http3` entry). It exists only under
`helix_ota/submodules/http3/` and `helix_terminator/submodules/http3/` (`.gitmodules`, verified via
directory search). Adopting it under Option A with **zero new vendoring** therefore requires
registering `helix_ota` (or `helix_terminator`) as a **second** `--ecosystem-root` /
`HELIX_ECOSYSTEM_ROOTS` entry — the *exact* multi-root extension
`g14_x1_submodule_policy_decision.md` §5 already identified as needed for `open-design` (§7 below).
This generalizes that finding: it is not an `open-design`-only need, it recurs for `http3`. Without
that second root, Option A's fallback branch fires and `http3` vendors fresh, locally, at
`skill-system/submodules/http3/` — legitimate under Option A, just not the zero-new-copy outcome.

**Serves:** R2 (working-POC quality bar; §11.4.85 stress/chaos coverage inherited for free).
**Gated on G14/X1:** the ecosystem-root registration is a direct, concrete instance of the G14/X1
"documentation/wiring obligation" already flagged in `g14_x1_submodule_policy_decision.md` §5
failure-mode 1 — not a new gate, but a second dependency that needs the same fix.

---

## 3. PipelineRuntime — verdict: **ADOPT-as-submodule, scoped to a genuinely new capability, deferred until that capability is built**

### 3.1 What it is (verified directly, `helix_code/submodules/pipeline_runtime/`, module `dev.helix.pipeline`)

A staged streaming-dataflow runtime: `Stream[T]` (`Subscribe(ctx, sink func(T) error) error`),
`Operator[T,U] func(Stream[T]) Stream[U]`, concrete operators `Map`/`Filter`/`FlatMap`/`Window`/
`MapReduce` (`pipeline.go:33-186`, all read directly), plus a flow-based-programming layer:
`Port` (bounded channel, `fbp.go:14-51`), `Component interface{ Name() string; Run(ctx) error }`,
`Network` (`Add`/`Connect`/`Run`, `fbp.go:62-100`) with bounded backpressure.

### 3.2 Deliberate scope non-overlap with the already-resolved `dag_orchestrator` verdict

`research/g05_g11_reuse_debate_dag_design.md` §2.3 already resolves `dag_orchestrator` as
"INCORPORATE-as-submodule (scoped)... replace `g11_worker_design.md`'s hand-designed
`supervise()`-wraps-nodes-directly + **flat sequential auto-expand loop**... driving each tick's
gap-set and draft-set" — a direct, cited description of `project/internal/autoexpand/pipeline.go`'s
`Pipeline.Run` (`:319-420`, read directly in this pass: a depth-layered BFS over `DetectGaps` →
`DraftSkill` → `store.Create` per gap, today a flat sequential `for`-loop with no concurrency and no
backpressure). **`dag_orchestrator` already claims this exact loop's scheduling machinery.**
Recommending `PipelineRuntime`'s operators for the *same* loop would directly conflict with an
already-decided verdict — this design deliberately does not do that.

### 3.3 The genuinely open niche: R6 wizard progress-streaming (confirmed not to exist today)

`REQUIREMENTS.md` R6 states the wizard flow ends with "progress reported back to the client."
Verified directly: `internal/api/expand_handler.go` (read in full) exposes only `POST
/api/v1/expand` (fire-and-forget job trigger) and `GET /api/v1/expand/status/:id` (poll-based
status read) — no SSE, no WebSocket, no streaming of any kind. A repo-wide search
(`grep -rln "progress\|SSE\|stream" internal --include=*.go`) turns up only unrelated hits (a
`models/skill.go` field comment, `sandbox.go`'s process-output-streaming, `middleware.go`'s
request-logging) — **zero** existing progress-streaming mechanism anywhere in the codebase. This is
not a redirect of existing code; it is a design proposal against a capability that genuinely does
not exist yet (see honest gap §10.4).

`PipelineRuntime`'s `Stream[T]`/FBP `Network` is a direct structural fit for this **specific** gap:
`worker.Runner`'s per-job execution (`internal/worker/runner.go`, `executeJob`/`recordSuccess`/
`recordFailure`, `:280-320,572-620`) could publish typed progress events (`GapDetected`,
`SkillDrafted`, `SkillStored`, `ExpansionComplete`) onto a bounded `Port`; the API layer's
`handleGetExpandStatus` (or a new SSE handler) subscribes via `Stream.Subscribe` and forwards each
event to the client as it arrives, replacing today's poll-only status read with genuine real-time
progress — directly serving R14's "expanding/updating... live, improvements available immediately."

### 3.4 Explicit mismatch noted (avoiding over-claim): the validation pipeline is NOT a good FBP fit

`internal/validation/pipeline.go`'s 4-stage sequential runner (`source_verification` →
`static_code_check` → `llm_jury` → `cross_reference`, `stageOrder` at `:29`) has fail-closed
short-circuit aggregation semantics (`computeOverallVerdict`, cited by `g05_g11_reuse_debate_dag_design.md`
§1.2 as already-correct) that do not map cleanly onto `PipelineRuntime`'s generic backpressure/flow
semantics — that pipeline's job is "any SKIP/FAIL is fail-closed," not "consume as fast as
downstream allows." Recommending `PipelineRuntime` there would be a semantic mismatch of the same
class that made `MCP_Module`/`debate_orchestrator` STUDY-only elsewhere in this research corpus;
this document does not recommend it for the validation pipeline.

### 3.5 Transitive deps (verified `helix_code/submodules/pipeline_runtime/helix-deps.yaml`)

Leaf — `deps: []`. README's own transitive note: if the consumer later wires fan-out onto
`vasic-digital/Concurrency` (worker-pool substrate), that dependency is added by the **consumer**,
never nested inside `pipeline_runtime` itself (§11.4.28(C)).

**Serves:** R14 (directly — the only concrete design in this corpus for its "stream progress"
clause). **Gated on G14/X1:** no (new dependency; present under `helix_code`, resolves with zero new
vendoring under the existing single ecosystem root — no multi-root extension needed, unlike `http3`
§2.4). **Timing:** deferred — implement when the progress-streaming capability is actually built
(no sooner), not immediately actionable today.

---

## 4. token_optimizer — verdict: **WIRE** (already vendored; confirms and extends the survey's §6C finding)

### 4.1 `pkg/wire`'s actual contract (verified directly, `constitution/submodules/token_optimizer/pkg/wire/`)

`Encoder` interface (`encode.go:48`): `Name() string`, `Encode(v any) ([]byte, error)`, `Decode(data
[]byte, v any) error`. One shipped implementation, `CompactJSON` (`encode.go:69-78`, a thin
`json.Marshal`/`json.Unmarshal` wrapper). `Selector` (`select.go:81`): `NewSelector(encoders
...Encoder) (*Selector, error)` validates non-empty/unique names and sorts deterministically
(§11.4.50); `Default()` (`select.go:120`) returns a `Selector` over `CompactJSON` alone; `Select(v
any) (Result, error)` (`select.go:135`) encodes with **every** registered encoder, keeps only the
ones that round-trip losslessly (`Decode(Encode(v))` deep-equals `v`), and returns the smallest —
never a smaller-but-lossy encoding (`ErrNoLosslessEncoder` is the hard floor, §11.4.6).

**Critical, load-bearing fact for this project (confirmed directly, `select.go:114-118`'s own doc
comment):** `token_optimizer` ships **no TOON encoder**. `Default()` is explicitly "the standalone,
dependency-free entry point; a consumer that also has a TOON (or other) encoder builds its Selector
with `NewSelector(CompactJSON{}, toonEncoder, ...)` instead." The `min(TOON, compactJSON)`
shape-routed encoding the survey's §6C described is a **capability of the `Selector` given both
encoders**, not something `token_optimizer` provides out of the box. The project's own
`research/toon_go_codec.md` (Revision 1, already decided: vendor `github.com/toon-format/toon-go`,
not hand-roll) is the **other half** of this wiring — a small adapter type implementing `wire.Encoder`
over `toon-go`'s `Marshal(v any, opts ...EncoderOption) ([]byte, error)` /
`Unmarshal(data []byte, v any, opts ...DecoderOption) error` is the missing piece, not yet designed
anywhere before this document (see §4.4).

### 4.2 Maturity finding beyond the survey (file-presence evidence, not run-verified — see honest gap §10.5)

The `token_optimizer` README's own table (Revision 2, 2026-07-08) describes `pkg/pipeline` /
`pkg/router` / `pkg/cache` / `pkg/telemetry` / `pkg/tier` as "planned." Directly listing those
directories in this pass shows each already contains substantial source **and** test files (e.g.
`pkg/pipeline/`: `pipeline.go`, `decision.go`, `cache.go`, plus five test files including
`evidence_wiring_test.go`, `savings_wiring_test.go`; `pkg/router/`: `router.go`, `failover.go`,
`loadbearing.go`, `evidence.go`; `pkg/transport/`: `negotiate.go`, `compress.go`, a `brotli/`
subpackage). The README's "planned" framing is stale relative to what is actually checked out —
this module is materially further along than its own last-written status line states. This is
**file-presence** evidence only; no `go build`/`go test` was run against it in this pass.

### 4.3 `pkg/transport/negotiate.go` vs. `http3` (§2) — distinct, non-overlapping concerns

`pkg/transport` is described in the README as "HTTP/3 + brotli seam (binary never on the
model-token channel)" — this is a **content-negotiation** concern (picking an encoding based on
`Accept`/`Accept-Encoding` for the LLM-token-cost-sensitive channel specifically), not a QUIC/HTTP3
**transport-server** concern. It does not overlap with, replace, or compete against §2's `http3`
verdict (the server bootstrap for `s.router`); the two operate at different layers of the same
request path and can be adopted independently.

### 4.4 Incorporation seam

1. A new small adapter (illustrative sketch, not final code) satisfying `wire.Encoder`:
   ```go
   type toonEncoder struct{}
   func (toonEncoder) Name() string { return "toon" }
   func (toonEncoder) Encode(v any) ([]byte, error) { return toon.Marshal(v) }
   func (toonEncoder) Decode(data []byte, v any) error { return toon.Unmarshal(data, v) }
   ```
2. Build the project's REST/MCP content-negotiation Selector once via
   `wire.NewSelector(wire.CompactJSON{}, toonEncoder{})` instead of `wire.Default()`.
3. Route `NegotiateResponse` (`internal/api/response.go`) and the P3 `ModelProvider` chain's
   request/response bodies through `Selector.Select(v)` rather than a bespoke TOON-vs-JSON size
   comparator — directly satisfies R2's addenda ("TOON primary + JSON fallback... needs a Go TOON
   codec," `REQUIREMENTS.md:32-42`) using components the project already owns (`toon-format/toon-go`
   per the already-decided `toon_go_codec.md`) plus (`token_optimizer` already vendored).

**Serves:** §11.4.141, §11.4.198, and the R2 TOON/JSON wire-format requirement directly.
**Gated on G14/X1:** no — already vendored inside `constitution/submodules/token_optimizer/` via
the §11.4.28(C) depth-1 carve-out (a different mechanism from Option A's ecosystem-root resolution
entirely — this dependency was never subject to the G14/X1 question in the first place).

---

## 5. VisionEngine + Panoptic — verdict: **confirmed real; N/A-until-R3-clients-exist**

### 5.1 VisionEngine (verified directly, `helix_code/submodules/vision_engine/`, module `digital.vasic.visionengine`, org `HelixDevelopment` per `helix_code/.gitmodules:65-67` `git@github.com:HelixDevelopment/VisionEngine.git`)

Four cooperating layers per its own README (verified, not re-derived from the survey): `pkg/analyzer`
(`Analyzer`/`VideoProcessor` interfaces, `UIElement`/`ScreenAnalysis`/`ScreenDiff`/`Rect`/`Size`/
`TextRegion`/`VisualIssue`/`ScreenIdentity`/`Action`/`KeyFrame` value types, `StubAnalyzer` reference
impl), `pkg/graph` (navigation-graph BFS + DOT/JSON/Mermaid export), `pkg/llmvision` (`VisionProvider`
+ OpenAI/Anthropic/Gemini/Qwen-VL/Kimi/StepGUI/Astica/Ollama adapters + `FallbackChain`), `pkg/config`.
Build-tag gated: `pkg/opencv` real GoCV only under `-tags vision`, buildable/testable without OpenCV
otherwise. Already a declared transitive dependency of `HelixQA` itself per the survey's own citation
(`helixqa/helix-deps.yaml`) — confirmed, not re-derived here.

### 5.2 Panoptic (verified directly, `helix_code/submodules/panoptic/`, module `panoptic` — literally unnamespaced, org `vasic-digital` per `helix_code/.gitmodules:293-295` `git@github.com:vasic-digital/Panoptic.git`)

A broader AI-enhanced multi-platform (web/desktop/mobile) test-generation + recording framework.
**A maturity signal worth flagging honestly:** its own README's badges and quick-start snippets still
reference the literal placeholder `your-org/panoptic` (`github.com/your-org/panoptic`) rather than
the real `vasic-digital/Panoptic` — unfinished README templating, not a functional defect, but a
signal this module has had less governance polish than `VisionEngine`/`LLMProvider`/`http3`.

### 5.3 Transitive deps (§11.4.31)

Both leaves — `vision_engine/helix-deps.yaml`: `deps: []`, own-org deps "none" (go.mod's only
requires are `stretchr/testify`, `gocv.io/x/gocv`, `golang.org/x/crypto` — no `replace` directives at
vasic-digital/HelixDevelopment). `panoptic/helix-deps.yaml`: `deps: []`, "own-org deps: none"
(Catalogue-Check comment, verified directly).

### 5.4 Verdict

**N/A-until-R3-clients-exist**, unchanged from the survey's §8A framing — both are genuine
candidates for §11.4.170's device-independent host-rendered visual-proof mandate (VisionEngine for
the OCR/vision layout oracle itself, Panoptic for AI-driven test-generation around it) the moment
the target's R3 Web/Desktop/Mobile clients exist, but there is **no R3 client in this project today**
(confirmed: `project/` contains only the Go backend — `cmd`, `internal`, `migrations`, `scripts`,
`config`, `deploy`, `seed`, `docs`; no client directory of any kind). Neither module is actionable
before that phase (P8+ per `REQUIREMENTS.md`'s draft staged plan).

**Serves:** §11.4.170 (future). **Gated on G14/X1:** no (both resolve via Option A with zero new
vendoring, present under `helix_code` already — moot until actionable).

---

## 6. KMP component suite — verdict: **ADOPT-consideration for Android/iOS/Desktop/Web; does NOT resolve HarmonyOS/Aurora; defer the weigh-off to the client-build phase**

### 6.1 New evidence beyond the survey's name+description-only pass (live `gh api`, this pass)

The survey's honest gap §11.11 explicitly flagged that the four KMP repos' "actual source/API
surface was NOT inspected." This pass fetched `vasic-digital/UI-Components-KMP`'s `build.gradle.kts`
and `settings.gradle.kts` directly (`gh api repos/vasic-digital/UI-Components-KMP/contents/{file}`,
2026-07-15, read-only):

```kotlin
group = "digital.vasic.uicomponents"
kotlin {
    androidTarget { ... }
    jvm("desktop") { ... }
    iosX64(); iosArm64(); iosSimulatorArm64()
    wasmJs { browser() }
    sourceSets { commonMain.dependencies {
        implementation(compose.runtime); implementation(compose.foundation)
        implementation(compose.material3); implementation(compose.ui)
    } }
}
```

This is a **real, non-trivial** Compose Multiplatform 1.11.0 project (`settings.gradle.kts`'s
`maven("https://maven.pkg.jetbrains.space/public/p/compose/dev")` repository, `gradlew`/`gradle/`
present per the earlier `contents` listing) targeting **Android + iOS (three ABIs) + Desktop (JVM) +
Web (WASM)** from one Kotlin codebase — a materially stronger evidentiary basis than the survey's
one-line description alone provided.

### 6.2 The honest, load-bearing limitation (a genuinely new finding — not previously assessed anywhere in this project's research corpus)

**`UI-Components-KMP`'s own build declares no HarmonyOS target and no Aurora OS target.** Compose
Multiplatform's official target list (Android/iOS/Desktop-JVM/Web-Wasm/Web-JS) does not include
either platform. This means KMP is in the **same structural class of gap** as
`research/g15_aurora_harmonyos_client_feasibility.md`'s already-researched finding for Flutter (no
official HarmonyOS-NEXT target, no Aurora target) — a real, working alternative for four of R3's six
platforms, but **zero** evidence, official or community, of covering the remaining two. Neither
`g15_aurora_harmonyos_client_feasibility.md` nor the R21/R22 survey previously compared KMP against
this specific gap; this document is the first to state it directly.

### 6.3 Verdict — does this change the survey's Flutter recommendation?

**Not a reversal, a genuine competing option now backed by real evidence.** The survey's existing
lean (`REQUIREMENTS.md:169-171`, "Flutter — the ONLY framework covering Android + iOS + HarmonyOS +
Aurora OS + desktop from one codebase") is **already contradicted by its own §8's finding** that no
Helix-family sibling has ever shipped Flutter-on-HarmonyOS/Aurora (the `helix_track` client-matrix
cross-check, §8 of the survey, not re-derived here). KMP does not fix the HarmonyOS/Aurora gap
either — but it is now a **confirmed real, previously-unconsidered, Helix-org-owned** (not
third-party) alternative for the other four platforms, with actual reusable building blocks already
published (`UI-Components-KMP` theme/animation/a11y, `Security-KMP` AES/Keychain/KeyStore,
`Concurrency-KMP` lazy-loading/flow-loaders, `Storage-KMP` network-storage abstractions — one-liners
confirmed via `gh repo view` in the survey, not re-verified deeper for these three in this pass, see
honest gap §10.3). **Recommendation: ADOPT-consideration, not ADOPT** — surface KMP as a
real, evidenced, Helix-owned competing option for the operator's R3/G15 sign-off at the client-build
phase (when Android/iOS/Desktop/Web client work actually begins), alongside Flutter and alongside
`helix_track`'s own precedent (native-per-platform + Tauri-wraps-Web for Desktop, survey §8) —
**do not** pick a winner now, and regardless of which wins, HarmonyOS/Aurora remain unsolved by
either and stay gated on `g15_aurora_harmonyos_client_feasibility.md`'s own findings.

**Serves:** R3, G15 (future). **Gated on G14/X1:** no (client-architecture question, not a
build-dependency layout question — same framing as the survey's §8 verdict on the client matrix).

---

## 7. design_system + open-design (R12) — verdict: **ADOPT-as-submodule, both**

### 7.1 design_system — one new confirmed FACT beyond the survey (live `gh api`, this pass)

The survey (honest gap §11.9) left `design_system`'s `helix-deps.yaml` status as `UNCONFIRMED` — "a
genuinely absent file, or the fetch silently 404'd." This pass re-ran the lookup directly:
`gh api repos/vasic-digital/design_system/contents/helix-deps.yaml` → HTTP 404 `"message":"Not
Found"`, **and** `gh api repos/vasic-digital/design_system/contents` (root listing) shows no
`helix-deps.yaml` among `[.gitignore, LICENSE, README.md, assets, components, docs, fonts, i18n,
manifest.json, package.json, scripts, tailwind, tokens, upstreams]`. **Upgraded from `UNCONFIRMED`
to CONFIRMED-ABSENT:** `design_system` genuinely does not ship a §11.4.31 manifest today. It does
carry an `upstreams/` directory (relevant to §11.4.36 `install_upstreams` on add) and a
`manifest.json` — a distinct, npm-package-shaped file, not the `helix-deps.yaml` schema; do not
conflate the two.

This is a real gap against §11.4.31, not a blocker: because `design_system` is `package.json`-shaped
(an npm-installable CSS/token asset, not a Go module with import-graph dependencies the way
`llmprovider`/`pipeline_runtime`/`http3` are), its "transitive dependency" surface is realistically
its `package.json` `dependencies` (unread in this pass — see honest gap §10.6), not a Go
`helix-deps.yaml`. Recommendation: when incorporating, either (a) treat absence-of-manifest as
"leaf, npm-shaped, no own-org Go deps" and proceed (consistent with `token_optimizer`'s own doc
being explicit that "if you find a project-specific string... it is a bug" — the absence here is
different in kind, not a hidden coupling), or (b) extend upstream to add one per §11.4.74, restoring
family-wide manifest coverage. Neither path blocks adoption.

### 7.2 open-design — already fully designed (`research/opendesign_incorporation.md` Revision 1); one discrepancy surfaced here, not previously reconciled

That document's §3.1 recommends mount path `submodules/open_design/` (renamed to snake_case per
§11.4.29's vendor-exception reading). The survey (§5) independently cross-checked two real sibling
precedents and found **both** `helix_terminator/.gitmodules:7-9` and `helix_vpn/.gitmodules:872-875`
mount it, unrenamed, at `submodules/open-design/` (hyphenated, matching the upstream repo name
`nexu-io/open-design` verbatim) — confirmed directly again in this pass by directory listing
(`helix_terminator/submodules/open-design/`, `helix_vpn/submodules/open-design/` both exist with the
hyphen, not an underscore). **This is a real, unreconciled discrepancy between the project's own two
research documents** (`opendesign_incorporation.md`'s recommended rename vs. the actual
100%-of-observed-siblings unrenamed convention) that neither document resolved against the other —
flagged here as a concrete open item for whoever lands the incorporation commit, not resolved
unilaterally in this design-only pass (§11.4.66-class sub-decision: rename-to-snake_case vs.
match-family-convention-verbatim).

### 7.3 Incorporation seam for both (R12 light+dark client tokens)

Both mount under `submodules/` (exact directory name per §7.2's flagged discrepancy), both third-party
(neither org is in the owned-org list — `nexu-io` for open-design, and `design_system`, while
`vasic-digital`-owned, is consumed as a **third-party-style** npm/CSS asset here, not co-engineered
Go code). **They are complementary, not substitutable, confirmed again in this pass:**
`open-design` is the raw per-brand `DESIGN.md` agent-driven **generation** substrate (150+ brand
systems, only 57/151 with a dark tier per `opendesign_incorporation.md` §2.5's own gap G5);
`design_system` is the already-curated, already-dual-themed (`tokens/themes/*.css`: `helix-green`
default + `vasic-red` + `helix-ota-blue`, each shipping all three sanctioned light/dark mechanisms —
`prefers-color-scheme`, `data-theme`, `.dark` class, per the survey's own `gh api contents` read, not
re-fetched here), Tailwind-v4-bound **consumable output**, already reused by three other Helix web
surfaces (HelixOTA/HelixCode/HelixTrack per its own README, per the survey). Adopting
`design_system` directly closes `opendesign_incorporation.md`'s own gap G5 for whichever brand theme
the target picks, while `open-design`'s `craft/accessibility-baseline.md` / `state-coverage.md` /
`anti-ai-slop.md` skills remain the review-rule source `design_system`'s CSS package does not
provide. Both serve R12; the Web (React) client per `opendesign_incorporation.md` §3.3 consumes
either's CSS-variable layer identically.

**Serves:** R12 (both). **Gated on G14/X1:** the mount-path/layout choice for both, yes (per the
survey's original framing); the incorporation mechanics for `open-design`, already fully designed
(Revision 1); `design_system`'s mechanics are newly confirmed here, not previously blocked on
anything beyond the same layout decision.

---

## 8. Consolidated incorporation-order table

| # | Module | Verdict | R-mapping | Transitive deps (§11.4.31) | Go-phase gate |
|---|---|---|---|---|---|
| 1 | LLMProvider (`HelixDevelopment/LLMProvider` — canonical under Option A today) | **ADOPT-as-submodule** — redirects R7/R19 away from a second hand-rolled Anthropic client (§9) | R7, R19, R17 | Leaf (`deps: []`) | Zero new vendoring — resolves at `helix_code/submodules/llm_provider/` under the existing single ecosystem root |
| 2 | http3 (`vasic-digital/http3`) | **ADOPT-as-submodule** (escalated from STUDY-only — diff done, §2.1) | R2 | Leaf | Requires registering `helix_ota` (or `helix_terminator`) as a **second** `--ecosystem-root` — not present under `helix_code` |
| 3 | PipelineRuntime (`HelixDevelopment/PipelineRuntime`) | **ADOPT-as-submodule**, scoped to R6 wizard progress-streaming (new capability, not yet built); explicitly NOT for the gap-expansion loop (already `dag_orchestrator`'s scope) nor the validation pipeline (semantic mismatch) | R6, R14 | Leaf | Zero new vendoring — present under `helix_code` |
| 4 | token_optimizer (`vasic-digital/token_optimizer`) | **WIRE**, already vendored in `constitution/submodules/`; needs a project-supplied `wire.Encoder` TOON adapter (§4.4), not a fresh vendor decision | §11.4.141, §11.4.198, R2 | N/A — §11.4.28(C) depth-1 carve-out, not Option A | Not gated — already present |
| 5 | VisionEngine + Panoptic | **N/A-until-R3-clients-exist** | §11.4.170 (future) | Both leaves | Zero new vendoring when actionable — present under `helix_code` |
| 6 | KMP suite (`UI-Components-KMP` + `Security-KMP` + `Concurrency-KMP` + `Storage-KMP`) | **ADOPT-consideration**, deferred weigh-off vs. Flutter at client-build phase; does NOT solve HarmonyOS/Aurora | R3, G15 (future) | Not independently confirmed (§10.3) | N/A — client-architecture decision, not a layout gate |
| 7 | open-design (`nexu-io/open-design`) | **ADOPT-as-submodule** — already fully designed; mount-path discrepancy flagged (§7.2), unresolved | R12 | Not a Go module (third-party design corpus) | Layout choice, yes (per survey) |
| 7b | design_system (`vasic-digital/design_system`) | **ADOPT-as-submodule** alongside open-design | R12 | Confirmed no `helix-deps.yaml` (npm-shaped; §7.1) | Layout choice, yes (per survey) |

---

## 9. Updated R7/R19 provider-layer design — the LLMProvider adapter seam (the single most important output of this document)

### 9.1 Why this redirects, rather than adds to, R19's existing design

`research/r19_anthropic_api_support_design.md` (Revision 1, status "design-research, no code
landed") already fully designed a hand-rolled `AnthropicLLM` (`internal/autoexpand/llm.go`, new)
plus a `NewLLMClientFromConfig` factory (§2.3 of that document, quoted in full below for the parts
this design changes). That document's §2.6 explicitly **declined** to vendor
`github.com/anthropics/anthropic-sdk-go`, reasoning: every existing provider client in this codebase
(`*OpenAILLM`, `*OpenAIEmbedder`) is a minimal ~70-120-line hand-rolled `net/http` transport, and
introducing "the FIRST officially-vendored SDK for exactly one of three-plus providers breaks that
consistency for no functional gain."

**That rationale does not extend to `LLMProvider`, and this section explains precisely why, rather
than asserting it:** `anthropic-sdk-go` is a raw, ungoverned, third-party SDK with zero relationship
to this project's own conventions. `LLMProvider` (`HelixDevelopment/LLMProvider` under Option A, §1)
is an **owned-org, equal-codebase submodule** (§11.4.28(A)) already vendored by four sibling Helix
projects, and its own provider packages (`pkg/providers/anthropic/anthropic.go`,
`pkg/providers/openai/openai.go`) are **the same size and shape** as this project's own hand-rolled
clients — each is a self-contained `net/http`-based transport of a few hundred lines, not a
monolithic official SDK wrapping streaming/tools/Batches/Files/Skills/Managed-Agents the way
`anthropic-sdk-go` does. Adopting `LLMProvider` is not "the first external SDK dependency" R19 §2.6
declined — it is the §11.4.74 catalogue-first mandate applied to a peer that has **already written
the same house-style client this project would otherwise write twice** (once here, already once at
`helix_terminator`/`helix_ota`/`helix_translate`/`helix_vpn`). R19 §2.6's conclusion stands
unchanged for `anthropic-sdk-go`; it simply never applied to `LLMProvider`, which was not evaluated
in that document at all (R19 predates the R22 catalogue mandate).

### 9.2 The redirected `NewLLMClientFromConfig` — same config surface, different construction

R19 §2.3 already designed the config plumbing this redirect reuses **unchanged**:

```go
// internal/config/config.go:124-131 (already exists in the project today, verified directly —
// LLMProvider + LLMModel fields present now; LLMAPIKey/LLMBaseURL are R19's additive, not-yet-
// landed fields, unchanged by this redirect)
type AutoExpandConfig struct {
	Enabled            bool
	MaxDepth           int
	MaxNewSkillsPerRun int
	LLMProvider        string // "openai" | "anthropic" | "local" | "helixllm"
	LLMModel           string
	LLMAPIKey          string // ${ANTHROPIC_API_KEY}-style; NEVER a literal (§11.4.10)
	LLMBaseURL         string // required for "local"/"helixllm"
}
```

**What changes is only the body of the factory's `"anthropic"` (and, optionally, `"openai"`) case** —
construct an `llmprovider` provider + a thin adapter, instead of a hand-rolled `*AnthropicLLM`:

```go
// internal/autoexpand/llmprovider_adapter.go (new file, illustrative sketch — not final code;
// package paths per §1.1's resolved canonical import "digital.vasic.llmprovider/pkg/providers/...")
package autoexpand

import (
	"context"
	"fmt"

	"digital.vasic.llmprovider/pkg/models"
	anthropicprovider "digital.vasic.llmprovider/pkg/providers/anthropic"
	openaiprovider "digital.vasic.llmprovider/pkg/providers/openai"
)

// llmProvider is the narrow subset of llmprovider.LLMProvider this adapter needs —
// Complete only; streaming/health-check/capabilities are NOT exposed through the
// existing LLMClient.Generate seam (a follow-on widening, §9.4, would change that).
type llmProvider interface {
	Complete(ctx context.Context, req *models.LLMRequest) (*models.LLMResponse, error)
}

// providerAdapter satisfies the existing autoexpand.LLMClient interface
// (Generate(ctx, prompt string, maxTokens int) (string, error)) by wrapping any
// llmprovider.LLMProvider implementation. Zero change to LLMClient itself, zero
// change to any caller (Pipeline.DraftSkill, a future LLMValidator juror, §9.3).
type providerAdapter struct {
	provider llmProvider
	model    string
}

// Generate honors the EXACT contract R19 §1 already specified for every LLMClient
// implementor: on success, non-empty text + nil error; on ANY failure to produce
// genuine text — HTTP error, malformed response, OR a policy refusal that returns
// HTTP 200 with no usable content — return ("", non-nil error), NEVER ("", nil).
// This is the one place AnthropicLLM's hand-rolled design (R19 §1) would have
// needed a check *OpenAILLM doesn't; the adapter closes that door here instead,
// so it applies uniformly to every llmprovider-backed provider, not just Anthropic.
func (a *providerAdapter) Generate(ctx context.Context, prompt string, maxTokens int) (string, error) {
	resp, err := a.provider.Complete(ctx, &models.LLMRequest{
		Prompt: prompt,
		ModelParams: models.ModelParameters{
			Model:     a.model,
			MaxTokens: maxTokens,
		},
	})
	if err != nil {
		return "", fmt.Errorf("llmprovider complete: %w", err)
	}
	if resp.Content == "" {
		// A 200-with-no-usable-text case (e.g. a policy refusal) MUST NOT look like
		// a normal empty-but-successful call — same anti-bluff floor R19 §1 states.
		return "", fmt.Errorf("llmprovider complete: empty content, finish_reason=%q", resp.FinishReason)
	}
	return resp.Content, nil
}

// NewLLMClientFromConfig — R19's plug-in point (§2.3), REDIRECTED: the "anthropic"
// and "openai" cases now construct llmprovider adapters instead of hand-rolled
// net/http clients. The "local"/"helixllm" case is UNCHANGED from R19 §2.3 (an
// OpenAI-compatible base-URL swap); this redirect does not touch it.
func NewLLMClientFromConfig(cfg AutoExpandConfigShape, logger Logger) (LLMClient, error) {
	switch cfg.LLMProvider {
	case "anthropic":
		if cfg.LLMAPIKey == "" {
			logger.Warn("anthropic LLM client created without API key; requests will fail")
		}
		baseURL := cfg.LLMBaseURL
		if baseURL == "" {
			baseURL = "https://api.anthropic.com" // R19 §2.1's already-verified default
		}
		p := anthropicprovider.NewProvider(cfg.LLMAPIKey, baseURL, cfg.LLMModel)
		return &providerAdapter{provider: p, model: cfg.LLMModel}, nil
	case "openai":
		if cfg.LLMAPIKey == "" {
			logger.Warn("openai LLM client created without API key; requests will fail")
		}
		baseURL := cfg.LLMBaseURL
		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}
		p := openaiprovider.NewProvider(cfg.LLMAPIKey, baseURL, cfg.LLMModel)
		return &providerAdapter{provider: p, model: cfg.LLMModel}, nil
	case "local", "helixllm":
		// UNCHANGED from R19 §2.3 — kept on *OpenAILLM (or an openaiprovider adapter
		// with baseURL swapped, equivalently) since HelixAgent/HelixLLM are already
		// confirmed OpenAI-compatible (REQUIREMENTS.md:157-160).
		if cfg.LLMBaseURL == "" {
			return nil, fmt.Errorf("llm_provider %q requires llm_base_url", cfg.LLMProvider)
		}
		p := openaiprovider.NewProvider(cfg.LLMAPIKey, cfg.LLMBaseURL, cfg.LLMModel)
		return &providerAdapter{provider: p, model: cfg.LLMModel}, nil
	default:
		return nil, fmt.Errorf("unsupported llm_provider: %q (expected "+
			"\"openai\", \"anthropic\", \"local\", or \"helixllm\")", cfg.LLMProvider)
	}
}
```

**Two honestly-stated implementation choices, not resolved unilaterally here (§11.4.66-class):**
(a) **minimal-diff** — redirect only the `"anthropic"` case through `llmprovider`, leave `"openai"`
on the existing `*OpenAILLM`; smaller PR, `*OpenAILLM` stays in the tree. (b) **full-redirect** —
both cases through `llmprovider` adapters, retiring `*OpenAILLM` entirely (subject to §11.4.124's
investigate-before-remove discipline — `*OpenAILLM` has real callers today, `git log`-cited removal
evidence would be required). (b) more completely satisfies R7's "pluggable ModelProvider, not
hardcoded OpenAI" per the survey's §4.3 framing; (a) is the safer first landing. **Does it supersede
`*OpenAILLM`?** Under (b), yes; under (a), not yet — it stays until a follow-on migration.

### 9.3 A capability the project does not have today: a real G05 jury `LLMValidator`

Confirmed by direct grep (`grep -rn "ValidateSkill\b" project --include=*.go`): the
`validation.LLMValidator` interface (`ValidateSkill(ctx, *models.Skill) (approved bool, feedback
string, err error)`, `pipeline.go:60-63`) has **zero concrete implementations anywhere in the
project** — only the interface declaration and its use inside `juryStage`
(`m.LLM.ValidateSkill(juryCtx, s)`, `pipeline.go:477`) exist. G05's already-correct two-factor
consensus (`g05_g11_reuse_debate_dag_design.md` §1.2, not re-derived here) has no real voter to
count votes from. This is a genuinely new design, not a redirect of existing code:

```go
// internal/validation/llmprovider_juror.go (new file, illustrative sketch)
type providerJuror struct {
	provider llmProvider // same narrow interface as §9.2's adapter
	model    string
}

func (j *providerJuror) ValidateSkill(ctx context.Context, s *models.Skill) (bool, string, error) {
	resp, err := j.provider.Complete(ctx, &models.LLMRequest{
		Prompt:      buildJuryPrompt(s), // a validation-specific prompt template, new
		ModelParams: models.ModelParameters{Model: j.model, MaxTokens: 1000},
	})
	if err != nil {
		return false, "", err // a real error is a non-approval, never a silent pass
	}
	return parseJuryVerdict(resp.Content) // approved bool + feedback string, new parser
}
```

Wired as `validation.JuryMember{Name: "anthropic-" + model, Weight: 1.0, LLM: &providerJuror{...}}`.
This gives G05's jury its **first** real, non-test-only voter, using the same `llmprovider` adoption
as §9.2 — and, because `llmprovider` ships `NewDefaultCircuitBreaker(name, provider)` for free
(§1.2), each juror can be circuit-breaker-wrapped with zero new code, directly closing an R17
resilience gap the survey (§4.3) already flagged. **This design does not depend on, and deliberately
avoids, any coupling to `DebateOrchestrator`** — `research/g05_g11_reuse_debate_dag_design.md` §1.3
already resolved `DebateOrchestrator` as STUDY-only for G05 specifically because its
`ConductDebate`/`ProviderInvoker` surface cannot address individual jurors by provider/model without
a fragile prompt-parsing workaround; the `providerJuror` design above sidesteps that entirely by
using `llmprovider`'s per-provider `Complete` call directly, one juror per provider, no debate
orchestration needed.

### 9.4 A deliberately NOT-taken widening (flagged, not designed)

`llmprovider.LLMProvider`'s full interface (`Complete`/`CompleteStream`/`HealthCheck`/
`GetCapabilities`/`ValidateConfig`) is strictly richer than the narrow `Generate` seam this section
adapts to. A follow-on could replace `autoexpand.LLMClient` with `llmprovider.LLMProvider` directly,
gaining streaming + capability negotiation + health-checking for the whole auto-growth pipeline
(R14's long-running job pipeline is exactly the kind of work that benefits) — the survey's §4.3
"path 2" already named this. This document deliberately does **not** design that widening: it is a
larger, callers-change-everywhere migration, out of scope for an adapter-seam design whose job is to
make the *minimal* correct redirect concrete.

### 9.5 Comparison — R19's hand-rolled `AnthropicLLM` design vs. this `LLMProvider` redirect

| Dimension | R19 hand-rolled `AnthropicLLM` (design-only, not landed) | This `LLMProvider` redirect |
|---|---|---|
| New dependency | None (net/http only) | `llmprovider` (already vendored elsewhere in the family, zero new vendoring under Option A, §1.1) |
| Lines of new project code | ~120-150 (a full `net/http` client, mirroring `OpenAILLM`) | ~40 (an adapter + factory-case rewrite) |
| Streaming | Not designed (R19 scope: `Generate` only) | Available on `llmprovider.Provider.CompleteStream` if §9.4 is later taken |
| Circuit breaker / health monitor / retry | Not designed | Free, via `llmprovider`'s root-level primitives |
| Providers gained per unit of new code | 1 (Anthropic only) | Effectively the same unit of code gains OpenAI, Anthropic, and (via §9.4) 48 more, if ever needed |
| G05 jury voter reuse | None — jury wiring untouched by R19 | Direct — §9.3's `providerJuror` reuses the identical adapter shape |
| Tension with R19 §2.6's stated rationale | N/A (its own design) | Resolved, not violated — §9.1 explains why §2.6's "no first vendored SDK" reasoning does not apply to an owned-org equal-codebase submodule |

**Recommendation:** adopt this redirect in place of implementing R19's `AnthropicLLM` design as
written. R19's config-plumbing design (§2.3's `AutoExpandConfig` additive fields, `substituteEnv`
wiring) is reused **verbatim** — nothing about it is wrong or wasted; only the construction inside
the factory's provider-specific branches changes.

---

## 10. Honest gaps (§11.4.6 — nothing here is silently assumed)

1. **LLMProvider org content-parity is still not diffed byte-for-byte.** §1.1 resolves *which*
   org-owned copy Option A's resolver would find first (`HelixDevelopment/LLMProvider` at
   `helix_code/submodules/llm_provider/`) — a narrower, now-answered question than the survey's — but
   whether that copy is byte-identical to `vasic-digital/LLMProvider` (vendored by
   `helix_terminator`/`helix_ota`/`helix_translate`/`helix_vpn`) remains unverified.
2. **`http3` and (implicitly) any future dependency present only under `helix_ota`/`helix_terminator`
   and not `helix_code`** will vendor fresh under Option A unless a second `--ecosystem-root` is
   registered — this document identifies the need (§2.4) but does not itself register that root or
   modify `sync_submodules.sh`'s invocation; that remains an implementation task for whoever adopts
   `http3`.
3. **`Security-KMP`, `Concurrency-KMP`, `Storage-KMP`'s actual source/API surface was NOT inspected
   in this pass** (only `UI-Components-KMP`'s `build.gradle.kts`/`settings.gradle.kts` were read via
   `gh api`) — their one-line descriptions (confirmed real repos via `gh repo view` in the survey,
   not re-verified deeper here) are the only evidence for them; a follow-on KMP-vs-Flutter weigh-off
   at the client-build phase must read all four before deciding anything.
4. **§3.3's progress-streaming design is a proposal against a capability that does not exist yet** —
   no `PipelineRuntime` code was written or proven to compile against `worker.Runner`/
   `expand_handler.go` in this pass; the fit is architectural, not implemented or tested.
5. **`token_optimizer`'s `pkg/pipeline`/`pkg/router`/`pkg/cache`/`pkg/telemetry`/`pkg/tier` maturity
   (§4.2) was assessed by directory-listing + signature-grep only** — no `go build`/`go test` was run
   against the module in this pass; the README's stale "planned" framing is contradicted by file
   presence, not by an executed test run.
6. **`design_system`'s actual `package.json` `dependencies` block was not read** in this pass — its
   confirmed-absent `helix-deps.yaml` (§7.1) is a manifest-format gap, not proof its real dependency
   surface is empty; treat as UNCONFIRMED, not asserted zero.
7. **§9.2's `providerAdapter`/§9.3's `providerJuror` are illustrative Go sketches, not compiled or
   tested code** — package import paths, exact `Logger`/`AutoExpandConfigShape` placeholder types,
   and the `buildJuryPrompt`/`parseJuryVerdict` helper functions are named but not implemented; a
   follow-on implementation pass must write and test them against the real
   `internal/config.AutoExpandConfig` and `internal/validation.LLMValidator` types, and must re-read
   `pipeline.go`'s `juryStage`/`computeOverallVerdict` (already resolved and cited, not re-derived
   here) before wiring `providerJuror` in.
8. **The default Anthropic model string / base URL used in §9.2's sketch (`"https://api.anthropic.com"`)
   reuses R19 §2.1's own already-§11.4.99-verified default** (`platform.claude.com/docs/en/api/overview`,
   verified 2026-07-15 per that document) — not independently re-verified against a live source in
   this pass; carried forward as a citation, not re-checked.
9. **`DebateOrchestrator`'s `ProviderInvoker` interface was not read in this pass either** (still the
   same gap `g05_g11_reuse_debate_dag_design.md` and the survey both already flagged) — irrelevant to
   §9.3's `providerJuror` design specifically, since that design does not depend on
   `DebateOrchestrator` at all, but the general gap remains open for whoever later revisits
   `debate_orchestrator`'s STUDY-only verdict.
10. **No code in this project was changed, built, or tested to produce this document** — every Go
    snippet above is a design sketch for a future implementation pass, consistent with this
    document's own stated deliverable contract (design-only, §11.4.197 "research… driven to a
    decision," not itself the implementation).

---

## Sources verified (this pass, 2026-07-15)

- `helix_code/.gitmodules:61-67,293-295` (LLMProvider org resolution; VisionEngine/Panoptic org
  URLs).
- `helix_code/submodules/llm_provider/{go.mod,doc.go,helix-deps.yaml}`;
  `pkg/models/types.go:10-59`; `pkg/providers/{anthropic,openai,claude}/*.go` (constructor +
  interface-method signatures read directly).
- `helix_code/submodules/pipeline_runtime/{README.md,go.mod,helix-deps.yaml,pipeline.go,fbp.go}`
  (full exported-symbol listing).
- `helix_ota/submodules/http3/{README.md,go.mod,helix-deps.yaml,pkg/server/server.go}` (exported
  symbols + line count); test-file directory listing (7 files).
- `helix_code/submodules/{vision_engine,panoptic}/{README.md,go.mod,helix-deps.yaml}`.
- `constitution/submodules/token_optimizer/{README.md,pkg/wire/*.go,pkg/{pipeline,router,transport,
  cache,telemetry,tier,config}/}` (directory listings + `wire/select.go`/`wire/encode.go` full read).
- Target project: `project/internal/autoexpand/llm.go` (full), `project/internal/autoexpand/pipeline.go`
  (full `Run` function + exported-symbol listing), `project/internal/validation/pipeline.go:1-110`
  (`LLMValidator`/`JuryMember`/`Pipeline` struct), `project/internal/worker/runner.go` (exported-symbol
  listing), `project/internal/api/{http3.go,expand_handler.go,server.go}` (full/grep),
  `project/internal/config/config.go` (grep for `AutoExpandConfig`), `project/go.mod`.
- `research/{helix_family_reusable_practices.md,g14_x1_submodule_policy_decision.md,
  g05_g11_reuse_debate_dag_design.md,r19_anthropic_api_support_design.md,toon_go_codec.md,
  opendesign_incorporation.md,g15_aurora_harmonyos_client_feasibility.md}` (cited sections, not
  re-derived).
- Live `gh api`/`gh repo view` (read-only, this pass, 2026-07-15):
  `repos/vasic-digital/UI-Components-KMP/contents/{settings.gradle.kts,build.gradle.kts}`;
  `repos/vasic-digital/design_system/contents` (root listing) and
  `repos/vasic-digital/design_system/contents/helix-deps.yaml` (confirmed 404).
- `find`/directory listings confirming `helix_terminator/submodules/open-design/`,
  `helix_vpn/submodules/open-design/` (hyphenated, unrenamed) and the absence of `http3`/
  `design_system` under any `helix_code/submodules/` checkout.
