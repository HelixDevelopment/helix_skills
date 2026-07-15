# R19 — Anthropic Messages API as a first-class `LLMClient` provider + embeddings honesty + interop-surface sub-decision

**Revision:** 1
**Last modified:** 2026-07-15T16:55:00Z
**Status:** design-research, no code landed
**Scope:** the Go backend `github.com/helixdevelopment/skill-system` under `project/`, package
`internal/autoexpand` (new `AnthropicLLM` + `NewLLMClientFromConfig`) and `internal/config`
(new `AutoExpandConfig` fields).
**Authority / mandates served:** R19 (`REQUIREMENTS.md:228-240`, operator mandate 2026-07-15) —
besides full OpenAPI-contract compatibility (G09), the system MUST support Anthropic's API(s);
composes R7 (`REQUIREMENTS.md:70-73`, pluggable ModelProvider, not hardcoded OpenAI) + R4
(`REQUIREMENTS.md:61-62`, agent interop). Constitution §11.4.6 (no-guessing), §11.4.8/§11.4.150
(deep-research-before-implementation), §11.4.10 (credentials — env-var NAME only, never a value),
§11.4.28/§11.4.161 (decoupling + rootless-container mandate, cited for the vendor decision),
§11.4.66 (interactive-clarification for the genuine secondary-facet sub-decision), §11.4.99
(latest-source verification for a risk-classified LLM-provider API — done below, footer cites
sources + date), §11.4.108 (four-layer runtime-signature), §11.4.115 (RED-first + polarity
switch), §11.4.186 (cross-doc lockstep with `testing_infrastructure_plan.md`), §11.4.197
(research fully completed/wired, never left un-wired).
**Read discipline:** every Go fact below was read from the committed baseline ref `255061b` via
`git show 255061b:…project/<path>` (the working tree has an uncommitted G02 change and a live Go
mutator, and was NOT read for `project/` source). No existing file was modified; this is the
single new deliverable.
**Composes with (read in full, no re-derivation):** `research/g20_autoexpand_realgrowth_design.md`
(G20 — already removes the `*OpenAILLM` concrete-type coupling from `Pipeline.DraftSkill`,
programs the pipeline to the bare `LLMClient.Generate` primitive, and explicitly flags "no
`NewLLMClientFromConfig` factory exists anywhere in the tree" as a tracked, un-designed wiring
gap — **this document is that factory's design**, the R19 plug-in point G20 named but did not
build); `research/g10_embedding_provider_design.md` (G10 — the sibling embedding-provider-factory
design (`NewEmbedderFromConfig`, `internal/db/embedding.go:293-309`) whose registry shape +
OpenAI-compatible-base-URL-swap pattern for `helixllm`/`local` this document reuses on the LLM
side rather than re-deriving a second, divergent factory shape).

---

## 0. One-paragraph problem statement

Today, `internal/autoexpand` has exactly ONE `LLMClient` implementation — `*OpenAILLM` — and, per
G20's own audit, ZERO factory that builds an `LLMClient` from `config.AutoExpandConfig` at all
(unlike the embedding side, which already has `NewEmbedderFromConfig`). R19 requires Anthropic's
Messages API to become a first-class, config-selectable `LLMClient` provider — usable
interchangeably with OpenAI for both the G05 jury (independent LLM voters) and G20's auto-growth
generation — without inventing a second interface shape and without silently assuming Anthropic
also offers embeddings (it does not — verified below, §4). A secondary, genuinely open question
(whether the skill system should ALSO expose an Anthropic-Messages-API-*shaped* HTTP surface for
R4 interop) is analyzed and surfaced as an explicit operator ask, not resolved unilaterally (§5).

---

## 1. The `LLMClient` interface TODAY (file:line) + what `AnthropicLLM` must satisfy

`internal/autoexpand/llm.go:26-29`:

```go
26:type LLMClient interface {
27:	// Generate creates text from a prompt with a token limit.
28:	Generate(ctx context.Context, prompt string, maxTokens int) (string, error)
29:}
```

One method. `*OpenAILLM` (`llm.go:36-55`, fields `apiKey`/`model`/`baseURL`/`httpClient`/`logger`
+ a request-rate-limit semaphore) is today's only implementor; its `Generate`
(`llm.go:88-155`) builds one `POST {baseURL}/chat/completions` request, parses
`openAIChatResponse.Choices[0].Message.Content`, and returns that string. **Per G20 §2.1**
(already decided, not re-derived here), `Pipeline.DraftSkill` now calls `p.llm.Generate(...)`
directly — it does **not** call any Anthropic- or OpenAI-specific method — so `AnthropicLLM`
needs to satisfy exactly this one method to be a fully pluggable third provider:

```go
Generate(ctx context.Context, prompt string, maxTokens int) (string, error)
```

**Contract `AnthropicLLM.Generate` must honor** (the same contract `*OpenAILLM.Generate` already
honors, so no caller-side branching is ever needed): (1) on success, return the generated text and
a `nil` error; (2) on ANY failure to produce genuine text — HTTP error, malformed response, or a
policy refusal — return `("", non-nil error)`, **never** `("", nil)`. The last clause is the one
place `AnthropicLLM` needs a check `*OpenAILLM` does not (Anthropic's Messages API can return
HTTP 200 with no usable text — see §2.2's refusal handling); getting it wrong would silently
reproduce the exact "no genuine content ⇒ still looks like a normal empty-but-successful call"
shape G20 spent its whole design forbidding at the `Pipeline` layer. `AnthropicLLM` closes that
door at the provider layer instead of trusting every future caller to re-check it.

---

## 2. DECISION

### 2.1 `AnthropicLLM` — a hand-rolled `net/http` client for `POST /v1/messages`, matching the house style

**Endpoint:** `POST {baseURL}/v1/messages`, default `baseURL = "https://api.anthropic.com"`
(overridable via `SetBaseURL`, mirroring `OpenAILLM.SetBaseURL`, `llm.go:67-70` — needed for tests
and for any future Claude-Platform-on-AWS / proxy routing, never hardcoded elsewhere).

**Required headers** (§11.4.99-verified today, 2026-07-15, against `platform.claude.com/docs/en/api/overview`):

| Header | Value |
|---|---|
| `x-api-key` | the resolved API key (§2.3 — never a literal in source or committed config) |
| `anthropic-version` | `2023-06-01` (current stable API version string — verified live) |
| `content-type` | `application/json` |

**Request body** (§11.4.99-verified against `platform.claude.com/docs/en/api/messages/create`):

```go
type anthropicMessageRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}
type anthropicMessage struct {
	Role    string `json:"role"`    // "user"
	Content string `json:"content"` // the GeneratePrompt(...) output, unchanged from today
}
```

**MUST NOT set `temperature`/`top_p`/`top_k`.** `*OpenAILLM`'s own request struct hardcodes
`Temperature: 0.7` (`llm.go`, `openAIChatRequest.Temperature`) — copying that literally into
`AnthropicLLM` would be a **live correctness bug**, not a harmless port: the currently-verified
model catalog (skill-cached 2026-06-24, re-confirmed live today) states plainly that Claude Fable
5, Claude Opus 4.7/4.8, and Claude Sonnet 5 **reject** any non-default `temperature`/`top_p`/`top_k`
with an HTTP 400. Since `AutoExpandConfig.LLMModel` is operator-configurable and will point at
whichever current-or-future model the operator names, `AnthropicLLM.Generate` omits all three
sampling parameters unconditionally — the one place this design deliberately diverges from the
`*OpenAILLM` template it otherwise mirrors, and the exact class of stale-training-prior mistake
§11.4.99 exists to catch.

**Response body** (§11.4.99-verified, same source):

```go
type anthropicContentBlock struct {
	Type string `json:"type"`           // "text" (only type this client reads)
	Text string `json:"text,omitempty"`
}
type anthropicStopDetails struct {
	Type        string `json:"type"`        // "refusal"
	Category    string `json:"category"`    // "cyber" | "bio" | "reasoning_extraction" | "frontier_llm" | ""
	Explanation string `json:"explanation"`
}
type anthropicMessageResponse struct {
	ID           string                  `json:"id"`
	Type         string                  `json:"type"`         // "message"
	Role         string                  `json:"role"`         // "assistant"
	Content      []anthropicContentBlock `json:"content"`
	Model        string                  `json:"model"`
	StopReason   string                  `json:"stop_reason"`  // end_turn|max_tokens|stop_sequence|tool_use|pause_turn|refusal
	StopSequence *string                 `json:"stop_sequence"`
	StopDetails  *anthropicStopDetails   `json:"stop_details"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}
```

**Error envelope on non-2xx** (§11.4.99-verified against the same source family — the standard
Anthropic error shape documented consistently across every current endpoint):

```go
type anthropicErrorResponse struct {
	Type  string `json:"type"` // "error"
	Error struct {
		Type    string `json:"type"`    // e.g. "invalid_request_error", "rate_limit_error"
		Message string `json:"message"`
	} `json:"error"`
	RequestID string `json:"request_id"`
}
```

### 2.2 Error handling — never a silent empty-success, mirroring G20's fail-closed discipline

```go
func (c *AnthropicLLM) Generate(ctx context.Context, prompt string, maxTokens int) (string, error) {
	reqBody := anthropicMessageRequest{
		Model:     c.model,
		MaxTokens: maxTokens,
		Messages:  []anthropicMessage{{Role: "user", Content: prompt}},
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal anthropic request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/messages", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create anthropic request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicAPIVersion) // "2023-06-01"

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic API request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return "", fmt.Errorf("read anthropic response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errEnv anthropicErrorResponse
		if json.Unmarshal(respBody, &errEnv) == nil && errEnv.Error.Type != "" {
			return "", fmt.Errorf("anthropic API returned %d [%s]: %s (request_id=%s)",
				resp.StatusCode, errEnv.Error.Type, errEnv.Error.Message, resp.Header.Get("request-id"))
		}
		return "", fmt.Errorf("anthropic API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result anthropicMessageResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal anthropic response: %w", err)
	}

	// A "refusal" is HTTP 200 with no usable content — the exact shape that,
	// left unchecked, silently reproduces the "empty-but-successful call"
	// class G20 forbids at the Pipeline layer. Fail here, not there.
	if result.StopReason == "refusal" {
		category, explanation := "", ""
		if result.StopDetails != nil {
			category, explanation = result.StopDetails.Category, result.StopDetails.Explanation
		}
		return "", fmt.Errorf("anthropic refused the request (category=%q): %s", category, explanation)
	}

	var textParts []string
	for _, block := range result.Content {
		if block.Type == "text" {
			textParts = append(textParts, block.Text)
		}
	}
	if len(textParts) == 0 {
		return "", fmt.Errorf("anthropic response contained no text content (stop_reason=%q)", result.StopReason)
	}
	return strings.Join(textParts, ""), nil
}
```

`NewAnthropicLLM(apiKey, model string, logger *zap.Logger) *AnthropicLLM` mirrors
`NewOpenAILLM` (`llm.go:51-65`) exactly in shape: defaults an empty `model` to
`"claude-opus-4-8"` (the current-catalog general-purpose default, §11.4.99-verified;
operator-overridable via the `llm_model` TOML key — never a hardcoded mandate), sets
`baseURL = "https://api.anthropic.com"`, and a `120*time.Second` `http.Client` timeout matching
`OpenAILLM`'s own. `SetBaseURL`/`SetHTTPClient` are provided for test/proxy overrides, same as the
OpenAI sibling.

### 2.3 `NewLLMClientFromConfig` — the R19 plug-in point G20 named but did not build

G20's own honest-gaps section (`g20_autoexpand_realgrowth_design.md §6`) states this factory does
not exist and explicitly asks it "reuse the same registry shape as G10 §2.1/§2.4." It does:

**New `AutoExpandConfig` fields** (`internal/config/config.go:124-131`, additive — nothing
existing removed or renamed, §11.4.124):

```go
// AutoExpandConfig controls the automatic skill-tree expansion.
type AutoExpandConfig struct {
	Enabled            bool   `toml:"enabled"`
	MaxDepth           int    `toml:"max_depth"`
	MaxNewSkillsPerRun int    `toml:"max_new_skills_per_run"`
	LLMProvider        string `toml:"llm_provider"` // "openai" | "anthropic" | "local" | "helixllm"
	LLMModel           string `toml:"llm_model"`
	LLMAPIKey          string `toml:"llm_api_key"`  // ${ANTHROPIC_API_KEY}-style; NEVER a literal secret (§11.4.10)
	LLMBaseURL         string `toml:"llm_base_url"` // required for "local"/"helixllm"; ignored otherwise
}
```

`config.go`'s existing `substituteEnv` (`config.go:311-333`) already resolves
`${VAR}`/`${VAR:-default}` for `Embedding.APIKey` (`cfg.Embedding.APIKey = sub(cfg.Embedding.APIKey)`,
`config.go:328`) — this design adds the symmetric line `cfg.AutoExpand.LLMAPIKey =
sub(cfg.AutoExpand.LLMAPIKey)` beside the two `LLMProvider`/`LLMModel` substitutions already there
(`config.go:330-331`). **Zero new config-parsing code** — the interpolation mechanism the codebase
already has is reused verbatim, satisfying §11.4.10 by construction (only the literal string
`${ANTHROPIC_API_KEY}` — a variable **name** — ever appears in a config file; the resolved value
lives only in the process's environment at runtime).

**The factory** (`internal/autoexpand/llm.go`, new — mirrors `NewEmbedderFromConfig`'s shape,
`internal/db/embedding.go:293-309`, exactly):

```go
func NewLLMClientFromConfig(cfg config.AutoExpandConfig, logger *zap.Logger) (LLMClient, error) {
	switch cfg.LLMProvider {
	case "openai":
		if cfg.LLMAPIKey == "" {
			logger.Warn("openai LLM client created without API key; requests will fail")
		}
		return NewOpenAILLM(cfg.LLMAPIKey, cfg.LLMModel, logger), nil
	case "anthropic":
		if cfg.LLMAPIKey == "" {
			logger.Warn("anthropic LLM client created without API key; requests will fail")
		}
		return NewAnthropicLLM(cfg.LLMAPIKey, cfg.LLMModel, logger), nil
	case "local", "helixllm":
		// Both are OpenAI-compatible chat-completions endpoints reached by
		// base-URL swap — REQUIREMENTS.md's own confirmed research already
		// states this: "HelixAgent (:7061) and HelixLLM (:8443) are both
		// OpenAI-compatible → one provider-agnostic client, swap base URL."
		if cfg.LLMBaseURL == "" {
			return nil, fmt.Errorf("llm_provider %q requires llm_base_url", cfg.LLMProvider)
		}
		client := NewOpenAILLM(cfg.LLMAPIKey, cfg.LLMModel, logger)
		client.SetBaseURL(cfg.LLMBaseURL)
		return client, nil
	default:
		return nil, fmt.Errorf("unsupported llm_provider: %q (expected "+
			"\"openai\", \"anthropic\", \"local\", or \"helixllm\")", cfg.LLMProvider)
	}
}
```

The `default` branch fails closed exactly like `NewEmbedderFromConfig`'s own `default`
(`embedding.go:305-307`) — an unrecognized provider string errors at construction, it never
silently falls back to a default provider. Empty-API-key behaviour **intentionally mirrors**
`NewEmbedderFromConfig`'s existing `"openai"` case (`embedding.go:296-298`, Warn-not-fail) rather
than inventing a second, stricter policy: G20 §2.2 (already decided) means any real
`Generate()` failure — including a 401 from a missing key — is now surfaced as a genuine,
captured `result.Errors` entry, never fabricated into a placeholder skill. Duplicating a
hard-fail-at-construction check here would be redundant defense with no additional safety, and
would diverge from the one factory-behaviour convention this codebase has already established.

### 2.4 Composition with G05 (jury) and G20 (auto-growth generation)

Neither G05's jury nor G20's `Pipeline.DraftSkill` (per G20 §2.1, already decided) references any
concrete provider type — both are written against the bare `LLMClient` interface. Once
`NewLLMClientFromConfig` returns an `*AnthropicLLM`, it plugs into **both** unchanged: a jury
entry configured with `provider = "anthropic"` casts a vote via the identical `Generate(ctx,
prompt, maxTokens)` call every other juror uses (composes G05, not re-derived here — this design
does not touch jury-quorum/threshold logic); `Pipeline.DraftSkill` calls
`p.llm.Generate(ctx, prompt, 4000)` exactly as it does for `*OpenAILLM` today. **No `*AnthropicLLM`
type assertion is ever introduced anywhere** — that would silently reintroduce the exact
concrete-type coupling G20 §2.1 already removed for `*OpenAILLM`. The R7 pluggability
`GAPS_AND_RISKS_REGISTER.md` demands (`REQUIREMENTS.md:70-73`) now has a second real, non-fake
implementor proving the interface is genuinely provider-agnostic, not merely "abstract in theory,
OpenAI-only in practice."

### 2.5 Credentials (§11.4.10)

The only credential-related literal in this entire design is the environment-variable **name**
`ANTHROPIC_API_KEY` (matching the ecosystem-standard name the `claude-api` skill's own
documentation uses throughout). No value, no example key, no placeholder-that-looks-real appears
anywhere in this document or in the proposed code. The operator supplies the value via
`.env`/shell export; `config.toml` carries only `llm_api_key = "${ANTHROPIC_API_KEY}"`.

### 2.6 Vendor decision — thin `net/http` client, NOT `anthropic-sdk-go` (§11.4.28/§11.4.161)

**DECISION: hand-roll the `net/http` client above; do NOT add `github.com/anthropics/anthropic-sdk-go`
as a new dependency for this requirement.**

Rationale:

1. **House-style consistency (§11.4.28).** Every existing provider client in this codebase —
   `*OpenAILLM` (`llm.go`) and `*OpenAIEmbedder` (`embedding.go`) — is a minimal hand-rolled
   `net/http` transport (~70-120 lines each), not a vendored official SDK. `AnthropicLLM` as
   designed above is the same size and shape. Introducing the FIRST officially-vendored SDK for
   exactly one of three-plus providers breaks that consistency for no functional gain: the
   Messages API surface this design actually needs — one endpoint, three headers, one small JSON
   envelope — is precisely what `anthropic-sdk-go` also does under the hood, plus streaming, tool
   use, Batches, Files, Skills, and Managed Agents surface this codebase's G05/G20 use case (a
   single blocking prompt→text call) does not touch.
2. **Dependency-weight discipline (adjacent to G14).** `GAPS_AND_RISKS_REGISTER.md`'s G14 finding
   already tracks 7 declared-but-unvendored dependencies as an open escalation
   (`testing_infrastructure_plan.md:46-48`). R7's own text (`REQUIREMENTS.md:72-73`) reads "*Any
   required dependency not present in the parent dir must be vendored as a submodule*" — taken
   literally, that clause would route `anthropic-sdk-go` through the SAME `submodules/<snake_case>/`
   vendoring discipline as a Helix-owned component, adding a NEW G14-class escalation for a
   dependency this design does not need at all. Choosing the thin client sidesteps that question
   entirely rather than answering it by accident.
3. **§11.4.161 (rootless containers) — checked, not implicated either way.** A Go module
   dependency (vendored `go.sum` entry, resolved at `go build` time) has no container-runtime
   footprint distinct from any of this project's other 18 `go.mod` dependencies; neither choice
   changes `deploy/docker-compose.yml` or the Podman rootless build. Cited for completeness, not
   because it discriminates between the two options.
4. **Flag, not foreclose.** If a FUTURE requirement needs richer Anthropic surface (streaming
   responses to a client, tool use, Batches for bulk skill drafting) beyond G05/G20's simple
   blocking-prompt need, escalating to `anthropic-sdk-go` is the right call THEN — as a normal
   `go.mod` dependency addition (it is third-party, not a Helix-owned component, so §11.4.28's
   submodule rule for OWNED components does not strictly apply to it; whether R7's broader
   "any required dependency" clause was meant to reach third-party SDKs too is the kind of
   ambiguity that should be resolved with the operator at that time, not guessed now). This
   design closes today's requirement with a strictly smaller footprint: zero new `go.mod` line,
   zero submodule question raised.

---

## 3. WHY (§11.4.8 external precedent, folded from §2's own citations)

**REQUIREMENTS.md's own confirmed research is the direct precedent for the `local`/`helixllm`
base-URL-swap branch** — "HelixAgent (`:7061`) and HelixLLM (`:8443`) are both OpenAI-compatible →
one provider-agnostic client, swap base URL" (`REQUIREMENTS.md:157-160`) — no new research
needed there, it is already-confirmed fact reused, not re-derived. **For the Anthropic wire
format itself**, the applicable precedent is Anthropic's own official documentation (§11.4.99,
cited below) — there is no ambiguity to resolve from a third-party source; the request/response
shapes above are transcribed, not inferred. **NO external solution found beyond official docs —
original work** for the factory-dispatch shape, which is a direct structural reuse of this same
codebase's own `NewEmbedderFromConfig` (`embedding.go:293-309`, cited in full in G10's design),
per the same "one already-established house pattern, applied a third time" reasoning G20 §3
already used for its own fail-closed-skip decision.

---

## 4. EMBEDDINGS — Anthropic offers none; embeddings stay on the G10 provider set

**FACT (§11.4.99-verified today, 2026-07-15, `platform.claude.com/docs/en/build-with-claude/embeddings`):**

> "Anthropic does not offer its own embedding model. One embeddings provider that has a wide
> variety of options and capabilities … is Voyage AI."

This is stated as flatly as an official vendor page states anything — there is no first-party
Anthropic embeddings endpoint, no `/v1/embeddings`-equivalent on `api.anthropic.com`, and no
"Anthropic embedding model" of any kind to configure. Consequences for this design:

- **The LLM-provider factory (§2.3) and the embedder factory (G10 `NewEmbedderFromConfig`,
  `embedding.go:293-309`) are, and remain, two SEPARATE dispatch points** keyed off two SEPARATE
  config sections (`config.AutoExpandConfig.LLMProvider` vs. `config.EmbeddingConfig.Provider`).
  `"anthropic"` is a valid value for the FORMER only; it is never a valid value for the latter,
  and `NewEmbedderFromConfig`'s existing `default` branch (`embedding.go:305-307`) already fails
  closed on any string it does not recognize, so `provider = "anthropic"` in the `[embedding]`
  TOML section errors today and MUST continue to error — never silently reinterpreted as "route
  Anthropic-Messages-API text through some embedding-shaped adapter," which would be a fabricated
  capability Anthropic does not offer.
- **Embeddings stay on G10's already-designed provider set** — `openai` / `local` / the
  `helixllm`-class OpenAI-compatible branch G10 §2.4 already plans, plus Voyage AI as the
  vendor-recommended third-party option `SPEC.md:395`'s sample config already names alongside
  `openai`/`local`. **Voyage integration is explicitly OUT OF SCOPE for this R19 design** — it is
  an embedding-provider addition, owned by G10 §2.4's already-designed extension point, not a
  re-derivation this Anthropic-focused document should duplicate. (Honest observation, not a
  design decision: Voyage's HTTP embeddings shape, per the same fetched page, is CLOSE to but not
  byte-identical to OpenAI's — `{"input":[...],"model":"voyage-4"}` → `{"data":[{"embedding":...,
  "index":...}],"model":...,"usage":{"total_tokens":...}}`, lacking OpenAI's `prompt_tokens` field
  and (per Voyage's own FAQ) expecting an `input_type` field for best retrieval quality — so it
  would need its own small adapter, not a bare base-URL swap onto `OpenAIEmbedder`. Flagged for
  G10's scope, not designed here.)
- **Never faked:** this design introduces no code path, adapter, or documentation that implies
  Anthropic embeddings exist. The `AutoExpandConfig`/`EmbeddingConfig` split above is the
  mechanical guarantee.

---

## 5. SECONDARY FACET — an Anthropic-Messages-API-*compatible SERVER surface? Genuine operator sub-decision, not resolved here (§11.4.66)

R19's text (`REQUIREMENTS.md:236-238`) explicitly flags this as a facet needing a decision, not a
foregone conclusion: **should the skill system ALSO expose its OWN `POST /v1/messages`-*shaped*
HTTP endpoint**, so that a tool speaking ONLY Anthropic's wire format (not MCP) could call the
skill system as if it were an Anthropic-compatible backend? This is a DIFFERENT direction from
everything in §2 above — §2 makes the skill system an Anthropic API **client** (consuming
Anthropic to generate/validate skills); this facet would make the skill system an Anthropic API
**server-shape mimic** (a second, parallel interop protocol alongside its own REST/OpenAPI surface
and its MCP stdio server).

**Analysis (not a decision):**

- **R4's interop question is ALREADY confirmed answered — by a different mechanism.**
  `REQUIREMENTS.md`'s own "Research findings — CONFIRMED" section states plainly: "Interop is
  solved by MCP. Claude Code, OpenCode, and Kimi Code all speak MCP natively over stdio. Ship ONE
  stdio MCP server + a thin CLI + shell aliases — no bespoke plugins." Every CLI agent R4 names by
  name is already covered by the MCP path a prior, dedicated research agent already confirmed and
  landed as the interop decision. Nothing in R4/R7/R19's own text names a consumer that can speak
  ONLY the Anthropic Messages API wire shape and cannot speak MCP.
- **A third protocol surface is a real, ongoing cost, not a one-time add.** The skill system
  already carries two protocol surfaces to keep honest and in sync: its native REST/OpenAPI API
  (`api/openapi.yaml`, G09's "route table == spec" discipline) and the MCP stdio server. Adding a
  THIRD — an Anthropic-Messages-*shaped* HTTP endpoint — would need its OWN drift discipline (a
  `/v1/messages`-compatible route whose request/response shape must track Anthropic's own schema
  evolution, not just this project's OpenAPI spec) for a consumer nobody has named yet.
- **RECOMMENDATION (not a mandate): do not build it as part of landing R19.** The already-decided
  MCP path satisfies every currently-named R4 consumer; building a speculative second interop
  surface ahead of a named need is exactly the kind of scope expansion §11.4.66 exists to gate
  behind an explicit ask rather than silently deciding either way.
- **THE EXPLICIT ASK (surfaced, not resolved):** *"Do you want the skill system to ALSO expose a
  `POST /v1/messages`-shaped HTTP endpoint, so a tool that can ONLY speak Anthropic's wire format
  (not MCP) can call it directly? If yes — which consumer/tool needs this (since MCP already
  covers Claude Code, OpenCode, and Kimi Code per the confirmed research), so the surface is
  scoped to a real requirement rather than built speculatively?"* This design deliberately stops
  here rather than guessing the answer either way (§11.4.6).

---

## 6. RUNTIME SIGNATURE (definition-of-done, §11.4.108)

| Scenario | Observable on a clean deployment |
|---|---|
| `llm_provider = "anthropic"`, real key reachable | `NewLLMClientFromConfig` returns a non-nil `*AnthropicLLM`; a real `Generate(ctx, prompt, maxTokens)` call against `api.anthropic.com` returns non-empty generated text and a `nil` error; feeding that client into `Pipeline.Run` (composes G20) produces the SAME anti-bluff guarantees G20 §4 already specifies (real LLM content, resources persisted in the same transaction, draft stays `draft` until a real G05 jury verdict) — with ZERO `*OpenAILLM`/`*AnthropicLLM` type assertions anywhere in `pipeline.go` (grep gate, extending G20's own). |
| `llm_provider = "anthropic"`, `llm_api_key` unset/empty | `NewLLMClientFromConfig` still returns a non-nil `*AnthropicLLM` (Warn-not-fail, §2.3) — but the FIRST real `Generate()` call returns a `401`-class wrapped error citing the Anthropic error envelope's `type`+`message`+`request_id`; composes with G20 §2.2 — that error lands in `result.Errors`, **zero** skill/resource rows are created for that gap, never a fabricated draft. |
| Non-2xx from Anthropic (429/500/etc.) | `Generate` returns a non-nil error naming the HTTP status + the envelope's `error.type`/`error.message` + `request_id`; no partial/garbage text is ever returned as if it were success. |
| `stop_reason == "refusal"` (HTTP 200, no usable content) | `Generate` returns `("", non-nil error)` citing the refusal category — never `("", nil)`; composes with G20 §2.2's real-fault-not-fabrication branch identically to any other `DraftSkill` error. |
| `llm_provider = "local"` / `"helixllm"` with `llm_base_url` set | `NewLLMClientFromConfig` returns an `*OpenAILLM` whose `baseURL` equals the configured value (verifiable via a getter/exported field in tests); Generate against a real reachable HelixAgent/HelixLLM OpenAI-compatible endpoint (`REQUIREMENTS.md:157-160`) succeeds identically to the `"openai"` path. |
| Unknown `llm_provider` string | `NewLLMClientFromConfig` returns `(nil, non-nil error)` — fail-closed, mirrors `NewEmbedderFromConfig`'s own `default` branch behaviour exactly. |
| `embedding.provider = "anthropic"` (misconfiguration) | `NewEmbedderFromConfig` (G10, unmodified by this design) still errors on its `default` branch — proving §4's separation-of-factories holds; no code path in this design ever makes that string succeed. |

---

## 7. TEST-CASE COUNT (RED-first per §11.4.115; each PASS cites captured evidence per §11.4.5/§11.4.69)

**Unit (no network) — 9**
1. `TestNewLLMClientFromConfig_Anthropic_ReturnsAnthropicLLM` — `LLMProvider="anthropic"` ⇒
   non-nil `*AnthropicLLM` with the configured model/key. **[RED on baseline: the function does
   not exist anywhere in the tree — G20 §6's own audit confirms this, package does not compile
   against it today]**
2. `TestNewLLMClientFromConfig_Openai_StillWorks` — `LLMProvider="openai"` ⇒ unchanged behaviour
   (regression proof the new factory doesn't disturb the existing provider).
3. `TestNewLLMClientFromConfig_LocalOrHelixLLM_UsesConfiguredBaseURL` — `LLMBaseURL` set ⇒
   returned `*OpenAILLM`'s base URL equals it (table-driven over `"local"`/`"helixllm"`).
4. `TestNewLLMClientFromConfig_LocalWithoutBaseURL_FailsClosed` — `LLMBaseURL == ""` ⇒ non-nil
   error, nil client.
5. `TestNewLLMClientFromConfig_UnknownProvider_FailsClosed` — mirrors
   `NewEmbedderFromConfig`'s own `default`-branch test shape.
6. `TestAnthropicLLM_Generate_MapsRequestFieldsAndHeaders` — fake `http.RoundTripper` captures
   the outbound request; asserts `model`/`max_tokens`/`messages` body fields AND
   `x-api-key`/`anthropic-version`/`content-type` headers are present; asserts NO
   `temperature`/`top_p`/`top_k` key is ever marshaled into the body. **[the §2.1 correctness
   point this design catches that a verbatim `*OpenAILLM` port would have missed]**
7. `TestAnthropicLLM_Generate_ParsesTextContent` — a recorded fixture response (the exact
   §11.4.99-verified schema above) ⇒ `Generate` returns the joined text, `nil` error.
8. `TestAnthropicLLM_Generate_NonOKStatus_ReturnsWrappedError` — stub 401 + the Anthropic error
   envelope ⇒ error message contains `error.type`+`error.message`+`request_id`; response is never
   parsed as success.
9. `TestAnthropicLLM_Generate_RefusalStopReason_ReturnsError` — stub 200 with
   `stop_reason:"refusal"`, `content:[]` ⇒ `Generate` returns `("", non-nil error)`, **never**
   `("", nil)`. **[RED on a naive port of `*OpenAILLM.Generate`'s shape, which has no refusal
   check at all — this is the concrete defect §2.2 exists to prevent]**

**Integration (build tag `integration`, real network — §11.4.27 no mocks beyond the sanctioned
unit-level fake transport above) — 2**
10. `TestIntegration_AnthropicLLM_Generate_LiveKey` — `//go:build integration`; `t.Skip` with
    reason if `ANTHROPIC_API_KEY` is unset (§11.4.3 SKIP-with-reason, never a silent green); else
    one real, low-`max_tokens` `Generate()` call against `api.anthropic.com`, captured response
    saved to `qa-results/<run-id>/anthropic_live_response.json`.
11. `TestIntegration_Pipeline_Run_ViaAnthropicProvider_NoPlaceholder` — (composes G20's own
    integration test 9) seeds a gap, constructs `Pipeline` via `NewLLMClientFromConfig` with
    `llm_provider="anthropic"` and a real key (SKIP-with-reason if absent), runs `Pipeline.Run`,
    asserts the created skill's content does NOT contain G20's named anti-bluff placeholder
    string AND that resources persisted in the same transaction — proving R7 pluggability
    end-to-end through a SECOND real provider, not just G20's own fake `LLMClient` double.

**Paired §1.1 mutations — 2** (each MUST flip its guard RED then restore GREEN)
- M1: route the `"anthropic"` branch of `NewLLMClientFromConfig` through an accidental
  `*OpenAILLM`-only code path (e.g. copy-paste the `"openai"` case's return) → test 1 FAILs
  (proves the factory really returns a distinct `*AnthropicLLM`, not a mislabeled OpenAI client).
- M2: delete the `result.StopReason == "refusal"` check in `AnthropicLLM.Generate` → test 9 FAILs.

**Total: 13 cases** (9 unit + 2 integration + 2 paired mutations).

### 7.1 Reconciliation with `research/testing_infrastructure_plan.md` (§11.4.186)

**No existing row covers this.** The plan's coverage matrix (`testing_infrastructure_plan.md:285-313`)
has a row per `GAPS_AND_RISKS_REGISTER.md` finding G01–G27; grepping both documents for
`anthropic`/`R19` finds exactly ONE hit — the G10 row's own citation of `SPEC.md:396`'s
`anthropic` embedding-config mention (`GAPS_AND_RISKS_REGISTER.md:149`, about the EMBEDDING
factory, not the LLM factory this design builds). **This is confirmed to be a genuinely NEW row**,
exactly as the dispatching task anticipated. Recommended addition, immediately after the existing
G20 row (`testing_infrastructure_plan.md:306`), continuing the register's `G01`–`G27` numbering
convention the plan's matrix structurally depends on:

| Gap | Sev | Test types (all applicable) | Challenge / HelixQA | Captured evidence each PASS cites |
|---|---|---|---|---|
| **G28** (NEW — R19) No `LLMClient` provider for Anthropic's Messages API; no `NewLLMClientFromConfig` factory exists at all (any provider) | MED-HIGH | unit(factory dispatch incl. anthropic/local/helixllm; request/response field-mapping; no-sampling-params; refusal + non-2xx fail-closed), integration(live `Generate` behind `integration` tag, SKIP-with-reason without `ANTHROPIC_API_KEY`; end-to-end `Pipeline.Run` via the anthropic provider, no-placeholder assertion), mutation(mislabel the anthropic branch as OpenAI → factory-distinctness test FAILs; drop the refusal check → refusal test FAILs) | Challenge: recommend **CH-G28** (submit `llm_provider=anthropic` + a stubbed refusal response, assert `Pipeline.Run` creates zero skill rows — same anti-bluff shape as CH-G20); HelixQA: extend HQA-WIZARD-ANDROID's provider matrix with an anthropic-backed run (operator-scheduled live window per §11.4.99/live-LLM honest boundary) | factory-dispatch table; captured request/response JSON fixture (headers + no-sampling-params assertion); refusal-path error message; live `Generate()` response (operator-scheduled) |

**Recommend this finding be formally registered as `G28` in `GAPS_AND_RISKS_REGISTER.md`** (the
next available number after G27) so the plan's row above has a stable anchor consistent with
every other row's `G01`–`G27` convention — an extension, per §11.4.186, not a contradiction of
anything already landed; no existing row's content changes.

---

## 8. HONEST GAPS (§11.4.6) — never guessed

- **Depends on G20's `Pipeline.DraftSkill` fix actually landing first.** This design's §2.4
  composition claim ("no type assertion is ever introduced") is true against G20's *designed*
  `DraftSkill` (§2.1 of that document); at `255061b`, `DraftSkill` still contains the
  `*OpenAILLM`-only assertion G20 itself flags as the current, unfixed state
  (`pipeline.go:215-218`). This design does not re-land G20's fix — it is designed to compose
  correctly the moment G20's own fix lands, not before.
- **Depends on G05's jury design for the "second juror" claim.** §2.4 states an Anthropic-backed
  juror votes identically to any other provider because G05's jury is (per G20's own citation)
  already written against the bare `LLMClient` interface — this document did not independently
  read a `research/g05_*.md` design doc (none was supplied) and does not re-derive G05's
  quorum/threshold logic; the composition claim is inherited from G20's own characterization of
  G05, not independently re-verified here.
- **Live-provider proof is operator-scheduled.** Tests 10/11 need a real, reachable
  `ANTHROPIC_API_KEY` — exactly the same honest boundary G10 (§5) and G20 (§6) already record for
  their own live-dependency tests. A mock-transport unit-test GREEN (tests 1–9) proves the
  request/response CONTRACT; it is never conflated with genuine live-quota/latency proof, which
  needs an operator-scheduled acceptance window.
- **The `AutoExpandConfig.LLMAPIKey`/`LLMBaseURL` field additions are new, not yet landed.**
  Like G10's own analogous embedder-construction wiring gap, these fields do not exist in the
  `255061b` baseline; this design specifies them precisely (§2.3) rather than assuming they
  already exist, and flags their absence so the work is not silently treated as already wired
  (§11.4.197).
- **§5's secondary facet is explicitly NOT decided here.** A recommendation is given (do not
  build an Anthropic-Messages-compatible server surface absent a named consumer), but the actual
  decision is the operator's per the explicit ask in §5 — this document does not silently assume
  either "yes, build it" or "no, never build it."
- **Voyage AI embeddings are out of scope** — flagged in §4 as owned by G10 §2.4's already-planned
  provider-factory extension, not re-designed here; the observed shape difference from OpenAI's
  embeddings response is noted as a fact for G10's future scope, not resolved.
- **Not in scope / not claimed:** improving G05 jury-vote QUALITY when an Anthropic model is one
  of the voters (a jury-design concern, not a provider-plumbing one); choosing which provider
  becomes the operator's DEFAULT for new deployments (a product decision, explicitly out of
  scope, same boundary G20 §6 already draws for its own provider-choice question).

*Positive-evidence-only. Every Go fact above is pinned to a file:line read from
`git show 255061b:…`; every Anthropic API fact is pinned to a live WebFetch performed today
(2026-07-15) against the URLs cited in the Sources-verified footer below, per §11.4.99. No
forbidden hedging vocabulary used.*

---

## Sources verified

- `https://platform.claude.com/docs/en/api/overview` — fetched 2026-07-15. Confirmed: endpoint
  base `https://api.anthropic.com`, required headers (`x-api-key` or `Authorization: Bearer`,
  `anthropic-version`, `content-type`), 32 MB request-size limit for Messages, response headers
  `request-id`/`anthropic-organization-id`.
- `https://platform.claude.com/docs/en/api/messages/create` — fetched 2026-07-15. Confirmed:
  `POST /v1/messages` request fields (`max_tokens`, `messages`), response fields (`id`, `content`
  array of typed blocks, `model`, `role`, `stop_reason`, `stop_sequence`, `stop_details`, `type`,
  `usage`), the full `StopReason` closed set including `"refusal"`, and a live example response
  matching the struct shapes transcribed in §2.1/§2.2 above.
- `https://platform.claude.com/docs/en/build-with-claude/embeddings` — fetched 2026-07-15.
  Confirmed verbatim: "Anthropic does not offer its own embedding model," Voyage AI is the
  vendor-recommended third-party provider.
- `claude-api` skill (bundled, session-loaded 2026-07-15; internal cache date 2026-06-24, within
  the §11.4.99 90-day risk-classified staleness window and independently re-confirmed live above)
  — cross-checked the current model catalog (`claude-opus-4-8` default), the
  no-sampling-parameters-on-current-models breaking change, and the `anthropic-version:
  2023-06-01` header convention used consistently across every current code example.
