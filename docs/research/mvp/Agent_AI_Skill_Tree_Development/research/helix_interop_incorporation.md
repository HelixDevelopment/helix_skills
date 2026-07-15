# Helix Ecosystem + Agent-CLI Incorporation Plan (R4)

**Revision:** 1
**Last modified:** 2026-07-15T00:00:00Z
**Scope:** R4 — how the HelixKnowledge Skill Graph System integrates with
HelixTrack, HelixAgent, HelixLLM (sibling ecosystem) and the agent CLIs
(Claude Code + toolkit/alias system, OpenCode, Kimi Code).
**Method:** Read-only inspection of real files under
`/home/milos/Factory/projects/tools_and_research/`. Every interface claim
below carries a `file:line` evidence pointer. Anything not proven from a real
file is labelled **UNCONFIRMED** with the reason. Prior research
`research/skillgraph_dim03_mcp_acp_protocols.md` (MCP is the universal interop
surface) and `SPEC.md §7` are built on, not redone.

---

## 0. TL;DR (verified)

- **Components that ACTUALLY exist in the sibling tree:**
  - `HelixLLM` → `helix_code/submodules/helix_llm/` (Go, OpenAI+Anthropic
    compatible LLM gateway). **EXISTS.**
  - `HelixAgent` → `helix_code/submodules/helix_agent/` (Go, ensemble
    multi-provider LLM service, REST+gRPC+MCP-bridge). **EXISTS.**
  - `LLMsVerifier` → `helix_code/submodules/llms_verifier/` (Go, model
    quality verifier/jury). **EXISTS** (named `llms_verifier`).
  - `HelixTrack` → **NOT** under `helix_code/`. It is a **separate sibling**
    at `helix_track/` (Go "JIRA alternative", issue/work-item tracker).
    **EXISTS**, at a different path than R4 implies.
  - Also present and relevant, already vendored as deps: `embeddings`,
    `helix_qa`, `challenges`, `docs_chain` under `helix_code/submodules/`.
- **HelixLLM binding (P3 ModelProvider):** OpenAI-compatible. Two real ports:
  the **binary default is `:8443` (HTTPS/TLS)**; the **operator-run local
  coder instance is `:18434` (plain HTTP)**. Bind via env
  `HELIX_LLM_LOCAL_OPENAI_ENDPOINT` (base URL, **NO** trailing `/v1`), POST
  `/v1/chat/completions`, `Authorization: Bearer $HELIX_LLM_API_KEY`. The
  plan's `:18434` is correct for the coder instance; the binary's own default
  is `:8443` — both verified.
- **HelixTrack interface:** **CONFIRMED as a service** (Go Core, Gin REST +
  JWT, `:8080`, `/health`, Consul discovery). The **exact skill→ticket CRUD
  route contract is PARTIALLY CONFIRMED** (request-envelope dispatch pattern +
  `Ticket`/`TicketStatus` domain models found; precise action-route names not
  fully enumerated — see §5).
- **HelixAgent interface:** **CONFIRMED.** REST default `:7061`
  (`/v1/chat/completions`, `/v1/ensemble/completions`, `/acp`), plus a gRPC
  `LLMFacade` and an MCP SSE bridge.
- **LLMsVerifier interface:** REST default **`:8080`** (README) — the plan's
  `:8081` is **UNCONFIRMED / likely wrong** (see §4). Primary Helix-native
  consumption is as a **Go config/scoring library**, not the REST port.

---

## 1. Component inventory (verified against real files)

| R4 name | Real location | What it is | Runtime interface | Evidence |
|---|---|---|---|---|
| HelixLLM | `helix_code/submodules/helix_llm/` | Single-binary Go LLM gateway: OpenAI **and** Anthropic compatible APIs, llama.cpp local inference, RAG, ReAct agents, multi-provider fallback | HTTP/3+HTTP/2 server; **binary default `HELIX_PORT=8443` (TLS)**; OpenAI routes `/v1/chat/completions`, `/v1/completions`, `/v1/models`, `/v1/embeddings`; Anthropic `/v1/messages`; health `/internal/health` | `helix_llm/README.md:5,37,44-49,86`; `helix_llm/.env.example:24` (`HELIX_PORT=8443`) |
| HelixLLM (coder instance) | runtime container `helixllm-coder` | llama.cpp OpenAI-compatible coder (Qwen3-Coder-30B) — the operator's live local model | **plain HTTP `http://localhost:18434`**, OpenAI `/v1/chat/completions`, `/v1/models`; live-proven HTTP 200 | `helix_code/RESUME.md:69-70`; `helix_code/internal/llm/tool_calling_concurrent_test.go:84` (`http://localhost:18434`); `helix_code/.superpowers/sdd/progress.md:214` (live POST `:18434/v1/chat/completions` → 200) |
| HelixAgent | `helix_code/submodules/helix_agent/` | Ensemble multi-provider LLM service (47+ providers), verification-scored fallback, sessions, webhooks | **REST default `:7061`** (`/v1/chat/completions`, `/v1/completions`, `/v1/ensemble/completions`, `/v1/models`, `/v1/providers`, `/acp`); **gRPC** `LLMFacade`; **MCP SSE bridge** binary | `helix_agent/README.md:324-333,339,352` (`:7061`); `helix_agent/internal/config/config.go:315` (`getEnv("PORT","7061")`); `helix_agent/cmd/grpc-server/main.go:26` (gRPC LLMFacade); `helix_agent/cmd/mcp-bridge/main.go:1-2` (MCP SSE bridge) |
| LLMsVerifier | `helix_code/submodules/llms_verifier/` | Enterprise LLM quality verifier / jury / router; "Do you see my code?" verification; config export | **REST default `:8080`** `/api/v1/verify`, `/api/v1/chat`, `/api/v1/models/:m/verify`, `/api/v1/config-exports/opencode`; **also a Go library** (`providers/config.go`) used in-process by the Helix fleet | `llms_verifier/README.md:153` (`port: 8080`), `:267,280,296,312`; `llms_verifier/llm-verifier/providers/config.go:15,22,40` |
| HelixTrack | `helix_track/` (**separate sibling**, not under `helix_code/`) | Open-source JIRA alternative: multi-platform issue/project tracker; **Core** Go backend + web/desktop/android/ios clients | **REST + JWT**, Core `htCore` on **`:8080`**, health `/health`, multi-space via `--space-root`; Consul discovery (`:8500`) with `:8080` fallback | `helix_track/CLAUDE.md:17-24`; `helix_track/compose.helixtrack.yml:19-51`; `helix_track/AGENTS.md:9` (`go build -o htCore main.go`); `helix_track/core/Application/CLIENT_INTEGRATION_GUIDE.md:53,58` (Consul `:8500`, fallback `:8080`) |

**Also present + already listed in `helix-deps.yaml` (relevant to R7/R8/R10):**
`embeddings` (`submodules/embeddings/`), `helix_qa` (`submodules/helix_qa/`),
`challenges` (`submodules/challenges/`), `docs_chain` (`submodules/docs_chain/`).

**Honest correction to R4's premise:** R4 says HelixTrack/HelixAgent/HelixLLM
are all in `../helix_code/`. Verified: **HelixAgent + HelixLLM are** (under
`submodules/`); **HelixTrack is NOT** — it is the top-level sibling
`../helix_track/`. Do not point a build/submodule resolver at
`helix_code/**/helix_track` — it does not exist there (only `helix_code/docs/helixtrack/`
documentation stubs exist).

---

## 2. Integration architecture (how it all fits)

The Skill Graph service integrates along two orthogonal axes:

1. **As a CLIENT of the Helix model layer (R7 — model access):** the P3
   `ModelProvider` interface fans out over a config-ordered chain:
   `HelixLLM (:18434 OpenAI-compat) → LLMsVerifier (scoring/jury) →
   claude-toolkit alias shellout (claudeN -p) → OpenAI (last resort)`.
   HelixAgent is an *optional* provider in the same chain (it is itself an
   OpenAI-compatible ensemble front — see §6).
2. **As a SERVER exposed to agents (R4 — agent interop):** the Skill Graph
   ships an **MCP server** (stdio primary, Streamable HTTP secondary — per
   `dim03` research and `SPEC.md §7`) that Claude Code, OpenCode, Kimi Code,
   **and HelixAgent's own MCP bridge** consume. Skill gaps/coverage are
   pushed to **HelixTrack** as tracked work items (§5).

MCP is the single interop surface for all agents (dim03 finding); the Helix
services are wired by HTTP client + config, never by importing project context
into them (§11.4.28 decoupling).

---

## 3. HelixLLM binding (P3 `ModelProvider`) — full contract

**This is the load-bearing integration.** The exact contract is already
implemented in the sibling `helix_agent` HelixLLM provider, so we mirror it
1:1 (do not reinvent):

`helix_code/submodules/helix_agent/internal/llm/providers/helixllm/provider.go`:

- **Endpoint resolution (precedence, verified `provider.go:26-63`):**
  1. `HELIX_LLM_ENDPOINT` — any OpenAI-compatible base URL.
  2. `HELIX_LLM_LOCAL_OPENAI_ENDPOINT` — **higher-precedence** seam for the
     local plain-HTTP coder (this is the one the fleet standardizes on).
  3. `HELIX_LLM_HOST` + `HELIX_LLM_PORT` → `http://$HOST:$PORT`.
  4. Fallbacks: `defaultHost=localhost`, `defaultPort=18434`
     (`provider.go:60-61`); binary TLS default `DefaultEndpoint =
     "https://localhost:8443"` (`provider.go:25`).
- **Base-URL gotcha (verified `provider.go:64-80`, `normalizeBase`):** the base
  URL MUST NOT carry a trailing `/v1` — the client appends the `/v1/...` path
  itself. `http://h:18434` and `http://h:18434/v1` both normalize to
  `http://h:18434`. A raw `.../v1` base double-appends → `/v1/v1` → 404.
- **Request shape (verified `provider.go:271-281,359-362`):**
  - `POST {base}/v1/chat/completions`, body = standard OpenAI chat request.
  - Headers: `Content-Type: application/json`; if a key is set,
    `Authorization: Bearer $HELIX_LLM_API_KEY` (`EnvHelixLLMApiKey`,
    `provider.go:222`); optional `X-Helix-LLM-Use-LlamaCpp: false`
    (`provider.go:281`).
  - Streaming: same path with `stream:true` + `Accept: text/event-stream`
    (`provider.go:360-361`).
- **Embeddings:** `POST {base}/v1/embeddings` (`provider.go:embeddingsEndpoint`).
- **Model discovery:** `GET {base}/v1/models` (`provider.go:modelsEndpoint`);
  the coder returns real model `Qwen3-Coder-30B-A3B-Instruct-Q4_K_M.gguf`
  (`.superpowers/sdd/progress.md:214`).
- **Health:** `GET {base}/internal/health` (`provider.go:healthEndpoint`).

### 3.1 Concrete P3 `ModelProvider` wiring

```go
// internal/provider/helixllm.go  (Skill Graph service — greenfield)
// Mirrors helix_agent/internal/llm/providers/helixllm/provider.go
type HelixLLMProvider struct {
    base   string // normalizeBase(os.Getenv("HELIX_LLM_LOCAL_OPENAI_ENDPOINT"))
    apiKey string // os.Getenv("HELIX_LLM_API_KEY") — optional for local coder
    hc     *http.Client
}
// base default: "http://localhost:18434"  (NO trailing /v1)
// chat:   POST base+"/v1/chat/completions"   OpenAI body, Bearer apiKey if set
// embed:  POST base+"/v1/embeddings"
// models: GET  base+"/v1/models"
// health: GET  base+"/internal/health"
```

**Config (`config.toml`, aligns with SPEC §8 `[embedding]`/new `[providers]`):**
```toml
[providers.helixllm]
enabled   = true
base_url  = "${HELIX_LLM_LOCAL_OPENAI_ENDPOINT}"  # default http://localhost:18434
api_key   = "${HELIX_LLM_API_KEY}"                 # optional for local coder
priority  = 1   # first in the fallback chain
```

**Reconciliation of the `:18434` vs `:8443` discrepancy (honest):**
- The Skill Graph `IMPLEMENTATION_PLAN.md:199` and `helix_code/RESUME.md:70`
  say `:18434` → **correct** for the *local coder instance* (plain HTTP,
  OpenAI-compat). Use this as the default.
- HelixLLM's *binary* README/`.env.example` default is `:8443` (TLS) → also
  correct, for a full `HELIX_MODE=full` gateway deployment. If the Skill Graph
  ever targets the TLS gateway, set `HELIX_LLM_ENDPOINT=https://host:8443` and
  handle TLS (the coder path is simpler and is the live default).
- **Both are real.** The env-driven `base_url` makes the Skill Graph agnostic
  to which one is running (§11.4.111 resolve-by-config, not hardcode).

### 3.2 Cross-check against LLMsVerifier (jury/scoring)

LLMsVerifier reads the **same** env var and **same** default endpoint —
verified `llms_verifier/llm-verifier/providers/config.go:15` (`const
helixLLMLocalOpenAIEndpointEnv = "HELIX_LLM_LOCAL_OPENAI_ENDPOINT"`) and `:22`
(`helixLLMDefaultBase = "http://localhost:18434"`), with a note that a live
`GET http://localhost:18434/v1/models` returned a real model list on
2026-07-08 (`config.go:18-20`). So HelixLLM and LLMsVerifier are already
wired to the identical coder endpoint — the Skill Graph joins that convention
verbatim (no new env var, no new default).

---

## 4. LLMsVerifier binding (validation jury + provider scoring)

**Role in the Skill Graph:** (a) the multi-model **validation jury** (SPEC §8
`[validation] jury_size=3`) and (b) the **provider-quality scorer** that ranks
the P3 fallback chain (R7). Two integration modes — pick by deployment:

### 4.1 Mode A — Go library (Helix-native, recommended)

LLMsVerifier is consumed **in-process** by HelixLLM/HelixAgent today via its
`providers` package (`llms_verifier/llm-verifier/providers/config.go`). If the
Skill Graph vendors `llms_verifier` (it is already in `helix-deps.yaml`), it
can call the verifier's Go API directly — no network hop, no port. This is how
HelixLLM "scores providers via LLMsVerifier, refreshed every 5 min"
(`helix_llm/README.md:11,108`).

### 4.2 Mode B — REST service

Verified REST surface (`llms_verifier/README.md`):
- `POST http://localhost:8080/api/v1/verify` (`README.md:267`)
- `POST http://localhost:8080/api/v1/chat` (`README.md:280`)
- `POST http://localhost:8080/api/v1/models/:model/verify` (`README.md:296`)
- `GET  http://localhost:8080/api/v1/models?verification_status=verified` (`README.md:304`)
- `POST http://localhost:8080/api/v1/config-exports/opencode` (`README.md:312`)
- Config default `port: 8080` (`README.md:153`).

**Discrepancy — the plan's `:8081` is UNCONFIRMED / likely wrong.**
`IMPLEMENTATION_PLAN.md:199` orders the chain `HelixLLM (:18434) →
LLMsVerifier (:8081)`. The LLMsVerifier README documents the API on **`:8080`**;
the only `:8081` hits in that repo are Kafka **schema-registry** and Flink
WebUI ports (`llms_verifier/docker-compose.messaging.yml:132-139`;
`llm-verifier/bigdata/streaming/flink_test.go:88`), not the verifier API.
**Action:** change the plan/config to `:8080` **or** confirm a project-specific
override before coding. Do not hardcode `:8081`.

**Contract test (evidence gate, per plan):** stub `POST /api/v1/verify` → assert
the jury verdict shape; a fallback test that fails provider-1 (HelixLLM 5xx)
and asserts provider-2 is selected from the verifier ranking.

---

## 5. HelixTrack binding (skill / issue sync) — PARTIALLY CONFIRMED

**What HelixTrack is (CONFIRMED):** an open-source JIRA alternative — Core Go
backend (`htCore`), RESTful API with JWT auth, multi-platform clients
(`helix_track/CLAUDE.md:17-24`). Runs on `:8080`, health `/health`, per-space
DB isolation via `--space-root` (`helix_track/compose.helixtrack.yml:44-51`).
Clients reach it via Consul service discovery (`:8500`) with a `:8080`
fallback (`core/Application/CLIENT_INTEGRATION_GUIDE.md:53,58`).

**Domain model (CONFIRMED, `helix_track/core/Application/internal/models/`):**
`Ticket{ id, title, description, ticketTypeID, ticketStatusID, projectID,
userID, creator, ticketNumber }` (`internal/models/ticket.go:12,25`),
`TicketStatus`, `CustomField` (`internal/models/customfield.go:3`), `Label` +
`LabelTicketMapping` (`internal/models/label.go:34`). Handlers exist:
`handleTicketStatusCreate/Read/...` (`internal/handlers/ticket_status_handler.go:16,106`).

**API style (CONFIRMED — request-envelope dispatch, NOT plain REST resources):**
handlers take `(c *gin.Context, req *models.Request)` and dispatch on an
action carried in the `Request` body (`internal/handlers/ticket_status_handler.go:17`).
The router exposes envelope endpoints like `POST /api/action`, `POST /api/data`,
`POST /api/submit`, plus `GET /tickets`, `GET /users` (route grep across
`core/Application`). **The precise action-name → ticket-CRUD mapping is
UNCONFIRMED** — the router file was not fully enumerated read-only, and much of
the raw route grep was polluted by test fixtures.

**Proposed integration (concrete, honest about the gap):**
- **Direction:** Skill Graph → HelixTrack. The Skill Graph's *registry* surfaces
  (`GET /registry/missing`, `GET /registry/stale`, gap reports — SPEC §6.3/§6.4)
  become HelixTrack **tickets**, one per missing/stale skill, so the human
  backlog and the knowledge-graph backlog stay one system.
- **Contract (to be finalized against the real router):** authenticate (JWT),
  then submit a `Ticket` via the envelope endpoint (candidate `POST /api/submit`
  or `POST /api/action` with `action=ticket.create`), fields mapped:
  `title = "Missing skill: <name>"`, `description = gap report`,
  `ticketTypeID = <task type>`, `projectID = <skill-graph project>`,
  `ticketStatusID = <queued>`. Idempotency key = skill name to avoid dupes.
- **Reverse direction (optional, UNCONFIRMED):** poll `GET /tickets` filtered by
  a `skill-graph` label to let humans request new skills from HelixTrack →
  enqueue an autoexpand job. Needs the filter/label API confirmed.

**Required before coding (§11.4.6 no-guessing):** confirm the exact ticket
create/read action route + JWT auth flow from
`helix_track/core/Application/internal/` router registration (or an OpenAPI —
**none found**, `find -iname '*openapi*'` empty). Treat the mapping above as a
**design**, not a verified contract, until that read is done. Contract tests
must run against a real `htCore` (§11.4.27 no fakes beyond unit).

---

## 6. HelixAgent binding (provider + task handoff + MCP)

**What HelixAgent is (CONFIRMED):** an ensemble LLM service that fronts 47+
providers with verification-scored fallback (`helix_agent/README.md:51,60`).
Three interfaces:

1. **REST `:7061` (CONFIRMED `internal/config/config.go:315`):** OpenAI-shaped —
   `POST /v1/chat/completions`, `POST /v1/completions`,
   `POST /v1/ensemble/completions`, `GET /v1/models`, `GET /v1/providers`,
   `GET /v1/health` (`README.md:324-333`). *Note:* `README.md:136` shows a
   generic env template `PORT=8080`; the **code default is `7061`**
   (`config.go:315`, corroborated `internal/config/config_test.go:43`). Confirm
   the running instance's `PORT` before wiring.
2. **gRPC `LLMFacade`** (`cmd/grpc-server/main.go:26,150,235`) — `Complete`,
   `CompleteStream`, `Chat`. For high-throughput internal calls.
3. **MCP SSE bridge** (`cmd/mcp-bridge/main.go:1-2`) — wraps stdio MCP servers
   over HTTP+SSE.

**Two integration options (both real, pick per need):**

- **(a) HelixAgent AS a `ModelProvider`** (drop-in): it is OpenAI-compatible on
  `:7061`, so the same P3 `HelixLLMProvider` code (§3.1) points at
  `http://localhost:7061` with `Authorization: Bearer $HELIXAGENT_API_KEY`
  (`HELIXAGENT_API_KEY`, `README.md:137`). Use when you want ensemble/jury
  routing done *for* you instead of the Skill Graph's own chain. Put it in the
  chain as an alternative to raw HelixLLM.
- **(b) HelixAgent AS a CLIENT of the Skill Graph MCP server** (task handoff):
  register the Skill Graph MCP server into HelixAgent's MCP bridge so a
  HelixAgent ReAct/ensemble run can call `skill_search` / `skill_get` /
  `learn_from_project` (SPEC §7 tools) during its own reasoning. This is the
  "agent task handoff" contract — HelixAgent delegates a knowledge lookup to
  the Skill Graph via MCP. Wire by adding our server to the bridge's stdio
  server list (same MCP `command`/`args` shape as §7.2 below).

**ACP note:** HelixAgent exposes `GET /acp` (route grep) and its HelixLLM
integration wires ACP via env `HELIX_LLM_USE_HELIXAGENT_ACP`
(`helix_agent/HELIXLLM_INTEGRATION_SUMMARY.md`). Per `dim03`, ACP is JSON-RPC
over stdio (NOT gRPC) — if the Skill Graph ever adds ACP, follow that, and MCP
tools ride *inside* the ACP session (dim03 §3.2). ACP is **not required** for
R4; MCP covers all named CLIs.

---

## 7. Agent-CLI MCP registration (paste-able)

The Skill Graph MCP server binary (SPEC §3 `cmd/server` with `[mcp] transport =
"stdio"`) is registered identically in shape across CLIs. Assume the built
binary is `skillgraph-mcp` (or `skill-system server --mcp-stdio`).

### 7.1 Claude Code (CONFIRMED shape — dim03 §4.1)

Project-scoped `.mcp.json` at repo root:
```json
{
  "mcpServers": {
    "skillgraph": {
      "command": "skillgraph-mcp",
      "args": ["--transport", "stdio"],
      "env": {
        "SKILLGRAPH_DB_URL": "postgres://skill:secret@localhost:5432/skilldb",
        "HELIX_LLM_LOCAL_OPENAI_ENDPOINT": "http://localhost:18434",
        "HELIX_LLM_API_KEY": "${HELIX_LLM_API_KEY}"
      }
    }
  }
}
```
Tools appear as `mcp__skillgraph__skill_search`, `mcp__skillgraph__skill_get`,
… (dim03 §4.2). Auto-approve read-only tools via settings `allowedTools:
["mcp__skillgraph__skill_search","mcp__skillgraph__skill_get","mcp__skillgraph__skill_tree","mcp__skillgraph__missing_skills","mcp__skillgraph__get_coverage"]`.

**Toolkit / alias consumption (CONFIRMED `claude_toolkit`):** the sibling
`claude_toolkit` installs multi-account aliases `claudeN` via a managed alias
file `~/.local/share/claude-multi-account/aliases.sh` sourced from
`.bashrc`/`.zshrc` (`claude_toolkit/AGENTS.md:42`), wrapping each `claudeN`
invocation with `cma_run` (`AGENTS.md:61`). Because every `claudeN` is just
Claude Code with a per-account `CLAUDE_CONFIG_DIR`/`HOME`, the **same project
`.mcp.json` is picked up by every alias automatically** — no per-alias MCP
config. The alias system is *also* usable as a **ModelProvider** in the P3
chain via headless shellout `claudeN -p "<prompt>"` (plan P3.T2;
constitution §11.4.187 `claude -p --output-format stream-json`).

### 7.2 OpenCode (CONFIRMED shape — dim03 §5.1)

`opencode.jsonc` (or `opencode.json`) under `mcp`:
```jsonc
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "skillgraph": {
      "type": "local",
      "command": ["skillgraph-mcp", "--transport", "stdio"],
      "enabled": true,
      "environment": {
        "SKILLGRAPH_DB_URL": "postgres://skill:secret@localhost:5432/skilldb",
        "HELIX_LLM_LOCAL_OPENAI_ENDPOINT": "http://localhost:18434"
      }
    }
  },
  "tools": { "skillgraph_*": true }
}
```
Per-agent tool gating via glob `skillgraph_*: true` (dim03 §5.2). For a remote
deployment use `"type": "remote"`, `"url": "http://host:8080/mcp"` (Streamable
HTTP) with `headers` for auth.

### 7.3 Kimi Code (**UNCONFIRMED** — no ecosystem precedent found)

**Honest gap:** neither the prior `dim03` research nor the sibling tree
(`helix_code/cli_agents*`, `claude_toolkit`) contains a Kimi Code MCP config —
`grep` for `kimi` found no config file. Kimi Code is not documented in any read
file here, so I will **not invent** its schema. Two verify-then-use paths:

- **Most-likely (many CLIs mirror Claude Code):** a `mcpServers` block using
  `command`/`args`/`env` (identical to §7.1). If Kimi Code follows the Claude
  Code convention, reuse the §7.1 JSON.
- **If Kimi Code follows an OpenCode-style config:** reuse §7.2.

**Required before claiming support:** fetch Kimi Code's live MCP docs
(§11.4.99 latest-source) and confirm the config path + key name. Until then,
Kimi Code support is **DESIGNED, not verified** — mark it as such in any status
doc (do not report it PASS).

---

## 8. Dependency resolution (§11.4.28C + R9) — helix-deps reconciliation

**Rule (§11.4.28C / R9):** every dependency resolves to ONE canonical location.
Order: **(1) reuse the PARENT copy if present** —
`helix_code/submodules/<name>/` (R9 "parent-dir versions have PRIORITY"); else
**(2) vendor under this project's own `submodules/<snake_case>/`**. Both copies
MUST stay in sync with main/master. Nested own-org chains FORBIDDEN.

**Current `helix-deps.yaml` (verified `project/helix-deps.yaml`) lists 7 deps:**
`llms_verifier`, `helix_llm`, `helix_agent`, `embeddings`, `helix_qa`,
`challenges`, `docs_chain` — all `layout: grouped` →
`submodules/<snake_case>/`.

**Reconciliation vs R4 targets:**

| R4 target | In helix-deps.yaml? | Correct? | Resolution |
|---|---|---|---|
| helix_llm | ✅ (`git@github.com:HelixDevelopment/HelixLLM.git`) | ✅ vendored dep (provider lib + endpoint convention) | Reuse parent `helix_code/submodules/helix_llm/`; else vendor `submodules/helix_llm/` |
| helix_agent | ✅ (`git@github.com:HelixDevelopment/HelixAgent.git`) | ✅ vendored dep (optional provider) | Reuse parent `helix_code/submodules/helix_agent/`; else vendor `submodules/helix_agent/` |
| llms_verifier | ✅ (`git@github.com:vasic-digital/LLMsVerifier.git`) | ✅ vendored dep (jury/scorer lib) | Reuse parent `helix_code/submodules/llms_verifier/`; else vendor `submodules/llms_verifier/` |
| **HelixTrack** | ❌ **absent** | ✅ **correctly absent** | HelixTrack is an **external peer SERVICE** we push tickets to (§5), **not a build/library dependency**. It is integrated over HTTP, not vendored. **Do NOT add it to helix-deps.yaml** unless a Go client library is actually imported. |

**Verdict:** the 7-entry `helix-deps.yaml` is **consistent with R4** — the three
model-layer components (helix_llm, helix_agent, llms_verifier) are vendored;
HelixTrack is (correctly) integrated as a service, not vendored. The other four
(embeddings, helix_qa, challenges, docs_chain) serve R7/R8/R10, not R4 directly.

**Git URLs (verified from `helix-deps.yaml`):**
- `git@github.com:HelixDevelopment/HelixLLM.git` (ref `main`)
- `git@github.com:HelixDevelopment/HelixAgent.git` (ref `main`)
- `git@github.com:vasic-digital/LLMsVerifier.git` (ref `main`)
- (HelixTrack repo URL not in manifest; the working copy is the sibling
  `../helix_track/` — its remote URL was not read here → **UNCONFIRMED**; if a
  Go client is ever needed, capture it from `helix_track/.git/config` first.)

**Sync obligation (R9):** whichever copy is authoritative, keep the parent
`helix_code/submodules/*` and any local `helix_skills/submodules/*` at the same
`main` tip; `scripts/sync_submodules.sh` (referenced by the manifest header)
drives this dry-run-first.

---

## 9. Port map (verified, for compose / config)

| Service | Port | Transport | Evidence | Note |
|---|---|---|---|---|
| HelixLLM binary (full gateway) | `8443` | HTTPS/TLS | `helix_llm/.env.example:24`, `README.md:37` | binary default |
| HelixLLM coder instance | `18434` | plain HTTP OpenAI-compat | `helix_code/RESUME.md:70`, `tool_calling_concurrent_test.go:84` | **the live default; use this** |
| HelixAgent REST | `7061` | HTTP OpenAI-compat + `/acp` | `helix_agent/internal/config/config.go:315` | README env sample shows `8080` — confirm running `PORT` |
| LLMsVerifier REST | `8080` | HTTP `/api/v1/*` | `llms_verifier/README.md:153` | **plan's `:8081` is wrong** (§4) |
| HelixTrack Core | `8080` | HTTP REST + JWT | `helix_track/compose.helixtrack.yml:44` | **collides with LLMsVerifier + HelixAgent-env-default `:8080`** — assign distinct host ports when co-deploying |
| Skill Graph HTTP | `8080` / `8443` | Gin (HTTP/2 default, HTTP/3 opt) | `SPEC.md §8 [server]` | same `:8080` default — **must remap** to avoid the 4-way `:8080` clash |
| Embeddings (TEI) | `18440` | HTTP | `helix_code/.superpowers/sdd/progress.md:156` | co-resident with coder |

**Port-collision warning (real):** HelixTrack Core, LLMsVerifier REST,
HelixAgent's env-template default, and the Skill Graph's own default all claim
`:8080`. When these run on one host, give each a distinct published port
(§11.4.111 resolve-by-config). The coder (`:18434`) and TEI embeddings
(`:18440`) are already distinct.

---

## 10. Honest gaps / UNCONFIRMED register

| # | Item | Status | Action to close |
|---|---|---|---|
| G1 | LLMsVerifier port `:8081` in `IMPLEMENTATION_PLAN.md:199` | **WRONG/UNCONFIRMED** — README says `:8080` | Fix plan+config to `:8080` or capture a project override |
| G2 | HelixTrack exact ticket-create/read action route + JWT flow | **PARTIALLY CONFIRMED** (envelope pattern + models found; route names not enumerated; no OpenAPI) | Read `helix_track/core/Application/internal/` router registration against a live `htCore` |
| G3 | HelixAgent running bind port | `7061` (code) vs `8080` (README env sample) | Confirm the deployed instance's `PORT` env |
| G4 | Kimi Code MCP config schema | **UNCONFIRMED** — no ecosystem precedent | Fetch Kimi Code live docs (§11.4.99); until then support is DESIGNED not verified |
| G5 | HelixTrack git remote URL | **UNCONFIRMED** (not in helix-deps.yaml; not read) | Read `helix_track/.git/config` only if a Go client is actually needed |
| G6 | HelixLLM binary vs coder in production | both real; env-driven | Default to `:18434` coder; document the `:8443` TLS path as alternate |

**Anti-bluff note (§11.4/§11.4.6):** every "CONFIRMED" above cites a real
`file:line`. Every "UNCONFIRMED"/"PARTIALLY CONFIRMED" says exactly what is not
proven and how to prove it. No interface was invented; where a live doc-fetch or
a router read is required (G2, G4), it is called out as a prerequisite, not
asserted.
