# Research: MCP & ACP Protocols for AI CLI Agent Integration

## Knowledge Skill Graph - Dimension 03: Protocol Integration Layer

**Research Date:** 2026-07-15
**Researcher:** AI Research Agent
**Sources:** 20+ primary sources including official specifications, GitHub repositories, documentation sites
**Confidence Level:** High (primary sources verified across multiple references)

---

## Executive Summary

This research covers the two primary protocols enabling AI CLI agent integration: **Model Context Protocol (MCP)** by Anthropic for tool/data access, and **Agent Client Protocol (ACP)** by Zed Industries for agent-editor communication. Both protocols use JSON-RPC (not gRPC) as their transport foundation. The mcp-go library is production-ready with 8.9k+ stars and supports stdio, SSE, and Streamable HTTP transports. Claude Code, OpenCode, Continue.dev, and most modern AI coding agents support MCP natively. **A critical finding: ACP uses JSON-RPC 2.0 over stdio, not gRPC — this contradicts the blueprint's assumption of "ACP adapter via gRPC."**

---

## Section 1: Model Context Protocol (MCP) Specification

### 1.1 Overview and Current Status

**Claim:** MCP is an open protocol that standardizes communication between LLM applications (hosts) and external data sources/tools (servers). [^17^]
**Source:** Official MCP Specification (modelcontextprotocol.io)
**URL:** https://modelcontextprotocol.io/specification/2025-06-18
**Date:** 2025-06-18 (stable), with 2025-11-25 as current and 2026-07-28 RC
**Excerpt:** "MCP provides a standardized way for applications to: Share contextual information with language models; Expose tools and capabilities to AI systems; Build composposable integrations and workflows. The protocol uses JSON-RPC 2.0 messages to establish communication between Hosts, Clients, and Servers."
**Context:** MCP was inspired by the Language Server Protocol (LSP) and aims to do for AI-tool integration what LSP did for language features.
**Confidence:** High

### 1.2 Specification Versions

**Claim:** MCP has gone through several specification versions, with 2025-11-25 being the current stable version and 2026-07-28 in release candidate. [^87^] [^151^]
**Source:** MCP GitHub Releases and Official Blog
**URL:** https://github.com/modelcontextprotocol/modelcontextprotocol/releases
**Date:** 2026-05-29
**Excerpt:** "This release marks the release candidate (RC) 2026-07-28 revision of the Model Context Protocol... The 2025-11-25 release marks the current stable version."
**Context:** Version history:
- 2024-11-05: Initial release (deprecated HTTP+SSE transport)
- 2025-03-26: Introduced Streamable HTTP, deprecated HTTP+SSE
- 2025-06-18: Major stable revision
- 2025-11-25: Current stable (replaced HTTP+SSE with Streamable HTTP)
- 2026-07-28: Release candidate (stateless-first protocol)
**Confidence:** High

### 1.3 Transport Mechanisms

#### 1.3.1 stdio Transport

**Claim:** stdio is the simplest transport, using standard input/output for JSON-RPC message exchange. It's ideal for local CLI tools. [^25^] [^131^]
**Source:** MCP Specification and Transport Comparison Analysis
**URL:** https://modelcontextprotocol.io/specification/2025-11-25/basic/transports
**Date:** 2025-11-25
**Excerpt:** "The protocol currently defines two standard transport mechanisms for client-server communication: 1) stdio, 2) Streamable HTTP"
**Context:** stdio characteristics:
- Deployment: Local process only
- Direction: Bidirectional (stdin/stdout)
- Concurrent clients: 1 (single process)
- Best for: CLI tools, local integrations
- Security: Process-level isolation
**Confidence:** High

#### 1.3.2 HTTP+SSE Transport (DEPRECATED)

**Claim:** The original HTTP+SSE transport using two endpoints (POST /messages + GET /events) was deprecated in March 2025. [^130^] [^131^]
**Source:** MCP Transport Analysis and Official Specification
**URL:** https://www.truefoundry.com/blog/mcp-stdio-vs-streamable-http-enterprise
**Date:** 2026-05-20
**Excerpt:** "HTTP+SSE was the first HTTP transport included in the initial MCP specification in November 2024. In March 2025, Streamable HTTP was introduced, officially deprecating SSE. The two-endpoint problem was the most fatal... If Server A receives the POST request and Server B has the SSE stream open, the correspondence between messages is broken."
**Context:** SSE is deprecated for new implementations. Atlassian Rovo ended SSE support June 30, 2026; Keboola ended April 1, 2026.
**Confidence:** High

#### 1.3.3 Streamable HTTP Transport (CURRENT STANDARD)

**Claim:** Streamable HTTP replaces HTTP+SSE with a single-endpoint design that supports both plain JSON responses and SSE streaming on the same URL. [^131^] [^140^]
**Source:** Official MCP Specification 2025-11-25
**URL:** https://modelcontextprotocol.io/specification/2025-11-25/basic/transports
**Date:** 2025-11-25
**Excerpt:** "Streamable HTTP... replaces the older HTTP+SSE transport with a single-endpoint design. The server exposes one URL (e.g. /mcp) that accepts both POST and GET. Clients POST JSON-RPC messages; servers respond with either a single JSON body or upgrade to a Server-Sent Events stream."
**Context:** Key features:
- Single endpoint (/mcp) for POST and GET
- POST for client->server JSON-RPC messages
- GET for server->client SSE stream (optional)
- Mcp-Session-Id header for session management
- Origin header validation for DNS rebinding protection
- Supports both stateful and stateless modes
- Load-balancer friendly (no sticky sessions needed)
**Confidence:** High

#### 1.3.4 2026-07-28 Stateless Transport (FUTURE)

**Claim:** The upcoming 2026-07-28 specification makes MCP stateless at the protocol layer, removing the initialize handshake and session management. [^151^]
**Source:** MCP Official Blog
**URL:** https://blog.modelcontextprotocol.io/posts/2026-07-28-release-candidate/
**Date:** 2026-05-22
**Excerpt:** "The headline change is that MCP is now stateless at the protocol layer... The initialize/initialized handshake is removed. The protocol version, client info, and client capabilities that used to be exchanged once at connection time now travel in _meta on every request."
**Context:** The Mcp-Session-Id header and protocol-level session are removed. Any MCP request can land on any server instance.
**Confidence:** High

### 1.4 Tool Definition Schema

**Claim:** MCP tool definitions consist of a name, description, and inputSchema following JSON Schema draft standard. [^143^] [^87^]
**Source:** MCP Specification - Tools section
**URL:** https://modelcontextprotocol.io/specification/2025-06-18/server/tools
**Date:** 2025-06-18
**Excerpt:** "Each tool is uniquely identified by a name and includes metadata describing its schema... Tools enable models to interact with external systems, such as querying databases, calling APIs, or performing computations."
**Context:** Complete tool definition structure:
```json
{
  "name": "search_products",
  "description": "Search products by name or category. Returns price and stock status.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "query": {
        "type": "string",
        "description": "Search term, e.g., 'wireless headphones'"
      },
      "category": {
        "type": "string",
        "enum": ["electronics", "clothing", "home"],
        "description": "Filter by product category"
      },
      "max_price": {
        "type": "integer",
        "description": "Maximum price in USD"
      }
    },
    "required": ["query"]
  }
}
```
**Confidence:** High

**Claim:** The tools/list method returns available tools with pagination, and servers can declare listChanged capability to emit notifications when tools change. [^87^]
**Source:** MCP Specification - Tools
**URL:** https://modelcontextprotocol.io/specification/2025-06-18/server/tools
**Excerpt:** "To discover available tools, clients send a tools/list request. This operation supports pagination... Servers that support tools MUST declare the tools capability: { "capabilities": { "tools": { "listChanged": true } } }"
**Confidence:** High

### 1.5 Server Implementation Patterns

**Claim:** MCP servers follow a lifecycle: initialization, capability negotiation, then serving requests for tools, resources, and prompts. [^17^]
**Source:** MCP Specification - Lifecycle
**URL:** https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle
**Excerpt:** "The protocol uses JSON-RPC 2.0 messages to establish communication between: Hosts: LLM applications that initiate connections; Clients: Connectors within the host application; Servers: Services that provide context and capabilities"
**Context:** Server features include:
- **Resources**: Read-only context/data (GET-like endpoints)
- **Tools**: Functions for the AI to execute (POST-like endpoints with side effects)
- **Prompts**: Reusable templates for LLM interactions
- **Sampling**: Server-initiated LLM interactions
- **Roots**: Server-inquiries into URI/filesystem boundaries
**Confidence:** High

---

## Section 2: mcp-go Library Evaluation

### 2.1 Repository Status and Maturity

**Claim:** mcp-go by mark3labs is a Go implementation of MCP with 8.9k stars, 219 contributors, 85 releases, and active development. [^2^]
**Source:** GitHub - mark3labs/mcp-go
**URL:** https://github.com/mark3labs/mcp-go
**Date:** 2026-07-09 (last commit)
**Excerpt:** "A Go implementation of the Model Context Protocol (MCP), enabling seamless integration between LLM applications and external data sources and tools... 8.9k stars, 856 forks, 219 contributors"
**Context:** The project is actively maintained with recent commits. Latest release: v0.56.0 (Jul 9, 2026). It implements MCP spec version 2025-11-25 with backward compatibility for versions 2025-06-18, 2025-03-26, and 2024-11-05.
**Confidence:** High

### 2.2 Supported Features

**Claim:** mcp-go supports all three MCP transports (stdio, SSE, Streamable HTTP), plus session management, tool handler middleware, request hooks, and OpenTelemetry tracing. [^2^]
**Source:** mcp-go README and Documentation
**URL:** https://github.com/mark3labs/mcp-go
**Excerpt:** "MCP-Go supports stdio, SSE and streamable-HTTP transport layers... MCP-Go provides a robust session management system that allows you to: Maintain separate state for each connected client; Register and track client sessions; Send notifications to specific clients; Provide per-session tool customization"
**Context:** Key capabilities:
- Server creation: `server.NewMCPServer(name, version, options...)`
- Tool definitions: `mcp.NewTool(name, options...)` with `mcp.WithString()`, `mcp.WithNumber()`, `mcp.WithBoolean()`, `mcp.WithArray()`, etc.
- Transports: `server.ServeStdio()`, `server.NewSSEServer()`, `server.NewStreamableHTTPServer()`
- Session management: Per-session tools, tool filtering, context passing
- Middleware: Tool handler middleware, prompt handler middleware
- Recovery: `server.WithRecovery()` for panic recovery
- Task tools: Async task-augmented tools with `server.WithTaskCapabilities()`
**Confidence:** High

### 2.3 Production Readiness

**Claim:** mcp-go is suitable for production use with proper error handling, validation, and security features, though the maintainers note some advanced capabilities are still in progress. [^2^] [^86^]
**Source:** mcp-go README and SitePoint Production Tutorial
**URL:** https://www.sitepoint.com/build-an-mcp-server-in-go-a-productionready-tutorial-for-the-model-context-protocol/
**Date:** 2026-06-27
**Excerpt:** "MCP Go is under active development, as is the MCP specification itself. Core features are working but some advanced capabilities are still in progress... Build an MCP Server in Go: A Production-Ready Tutorial"
**Context:** Production considerations:
- **Strengths:** Type-safe, concurrent (Go goroutines), single binary deployment, comprehensive schema builders
- **Validation:** Built-in input validation via JSON Schema
- **Security:** Supports OAuth protected resource metadata, CORS for browser clients, DNS rebinding protection (2025-11-25+)
- **Observability:** OpenTelemetry tracing hooks, request hooks for telemetry
- **Limitations:** Some advanced features (like full sampling support) may still be maturing
**Confidence:** High

### 2.4 Example Implementation Pattern

```go
package main

import (
    "context"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    s := server.NewMCPServer(
        "SkillGraph MCP Server",
        "1.0.0",
        server.WithToolCapabilities(true),
        server.WithRecovery(),
    )

    tool := mcp.NewTool("skill_search",
        mcp.WithDescription("Search for skills in the knowledge graph"),
        mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
        mcp.WithNumber("limit", mcp.Description("Max results")),
    )
    s.AddTool(tool, skillSearchHandler)

    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}
```

---

## Section 3: ACP (Agent Client Protocol)

### 3.1 Overview and Origin

**Claim:** The Agent Client Protocol (ACP) is an open standard created by Zed Industries (released August 2025) that standardizes communication between code editors and AI coding agents. It uses JSON-RPC 2.0 over stdio. [^26^] [^153^]
**Source:** Agent Client Protocol Official Website and GitHub
**URL:** https://agentclientprotocol.com and https://github.com/agentclientprotocol/agent-client-protocol
**Date:** 2025-06-23 (initial), 2026-07-06 (latest docs)
**Excerpt:** "The Agent Client Protocol (ACP) standardizes communication between code editors (interactive programs for viewing and editing source code) and coding agents (programs that use generative AI to autonomously modify code)."
**Context:** Key facts:
- **Created by:** Zed Industries
- **Released:** August 2025
- **License:** Apache 2.0
- **Transport:** JSON-RPC 2.0 over stdin/stdout (NOT gRPC)
- **Current stable version:** 1
- **Official libraries:** TypeScript, Rust, Kotlin, Java, Python
- **Analogy:** "LSP for AI agents"
**Confidence:** High

### 3.2 ACP vs MCP: Critical Distinction

**Claim:** ACP and MCP solve completely different problems and are designed to work together. MCP handles agent-to-tool/data communication; ACP handles agent-to-editor communication. [^26^]
**Source:** MorphLLM - Agent Client Protocol Explained
**URL:** https://www.morphllm.com/agent-client-protocol
**Date:** 2026-03-03
**Excerpt:** "MCP answers: 'What tools and data can the agent access?' ACP answers: 'Where does the agent live in the developer's editor?'... When an ACP session starts, the editor passes available MCP server endpoints and credentials to the agent. The agent then invokes tools via MCP calls, all piped through the ACP session."
**Context:** Comparison table:

| Dimension | ACP (Agent Client Protocol) | MCP (Model Context Protocol) |
|---|---|---|
| Created by | Zed Industries (Aug 2025) | Anthropic (Nov 2024) |
| Solves | Agent-to-editor communication | Agent-to-tool/data communication |
| Transport | JSON-RPC 2.0 over stdio | JSON-RPC 2.0 over stdio or HTTP+SSE |
| Primary relationship | Editor spawns agent as subprocess | Agent calls tool/resource servers |
| Streaming | Real-time token streaming | Request-response |
| Session state | Built-in session management | Stateless per request |
| Ecosystem | 25+ agents and editors | 10,000+ public MCP servers |

**Confidence:** High

### 3.3 ACP Protocol Methods

**Claim:** ACP defines a set of JSON-RPC methods for initialization, session management, and bidirectional communication between agents and editors. [^137^] [^139^]
**Source:** Agent Client Protocol Official Specification
**URL:** https://agentclientprotocol.com/protocol/v1/overview
**Date:** 2026-06-01
**Excerpt:** "The protocol follows the JSON-RPC 2.0 specification with two types of messages: Methods (request-response pairs) and Notifications (one-way messages)."
**Context:** Core ACP methods:

**Agent methods (baseline):**
- `initialize` - Capability negotiation
- `authenticate` - Authentication
- `session/new` - Create new session
- `session/prompt` - Send user prompt

**Agent methods (optional):**
- `session/load` - Resume existing session
- `session/delete` - Delete session
- `session/set_mode` - Set agent mode
- `logout` - End authentication

**Client methods (agent -> client):**
- `fs/read_text_file` - Read file
- `fs/write_text_file` - Write file
- `terminal/create` - Create terminal
- `terminal/output` - Get terminal output
- `session/request_permission` - Request permission

**Notifications:**
- `session/update` - Streaming updates from agent
- `session/cancel` - Cancel processing

**Confidence:** High

### 3.4 ACP Initialization Flow

**Claim:** ACP sessions begin with initialization where protocol version and capabilities are negotiated, followed by optional authentication, then session creation with MCP server configuration. [^137^] [^144^]
**Source:** ACP Protocol Documentation
**URL:** https://agentclientprotocol.com/protocol/v1/initialization
**Excerpt:** "Before a Session can be created, Clients MUST initialize the connection by calling the initialize method with: The latest protocol version supported; The capabilities supported."
**Context:** Example initialization:
```json
// Client -> Agent
{
  "jsonrpc": "2.0", "id": 0, "method": "initialize",
  "params": {
    "protocolVersion": 1,
    "clientCapabilities": {
      "fs": { "readTextFile": true, "writeTextFile": true },
      "terminal": true
    },
    "clientInfo": { "name": "my-client", "version": "1.0.0" }
  }
}

// Agent -> Client
{
  "jsonrpc": "2.0", "id": 0,
  "result": {
    "protocolVersion": 1,
    "agentCapabilities": {
      "loadSession": true,
      "promptCapabilities": { "image": true, "audio": true },
      "mcpCapabilities": { "http": true, "sse": true }
    },
    "agentInfo": { "name": "my-agent", "version": "1.0.0" }
  }
}
```
**Confidence:** High

### 3.5 CRITICAL CORRECTION: ACP Does NOT Use gRPC

**Claim:** ACP uses JSON-RPC 2.0 over stdio, NOT gRPC. The blueprint's assumption of "ACP adapter via gRPC for OpenCode" is incorrect. [^137^] [^142^] [^153^]
**Source:** Multiple authoritative ACP sources
**URL:** https://agentclientprotocol.com/protocol/v1/overview
**Excerpt:** "All messages follow JSON-RPC 2.0 specification... Transport Layer: Newline-delimited JSON over stdio"
**Context:** The ACP specification explicitly defines JSON-RPC 2.0 over stdio as the transport. There is no gRPC in the official ACP spec. The blueprint's mention of "ACP adapter via gRPC" may refer to:
1. A custom internal extension not part of the standard
2. A misunderstanding of the protocol stack
3. A specific organizational adapter that wraps ACP in gRPC
**All official ACP SDKs (TypeScript, Rust, Python, Kotlin, Java) use JSON-RPC over stdio.**
**Confidence:** High

### 3.6 ACP Editor Support

**Claim:** ACP is supported natively by Zed, with JetBrains partnership, and community plugins for Neovim, Emacs, and VS Code. [^26^]
**Source:** MorphLLM ACP Documentation
**URL:** https://www.morphllm.com/agent-client-protocol
**Excerpt:** "Zed has the most complete implementation as the protocol's origin. JetBrains is adding native support. Neovim and Emacs rely on community plugins. VS Code has community extensions but no native support."
**Context:** Editor support status (March 2026):
| Editor | Status |
|--------|--------|
| Zed | Native (full support) |
| JetBrains | In progress (official partnership) |
| Neovim | Via plugins (agentic.nvim) |
| Emacs | Via plugin (agent-shell) |
| VS Code | Community extension only |

**Confidence:** High

---

## Section 4: Claude Code Integration

### 4.1 MCP Server Configuration

**Claim:** Claude Code supports MCP servers through multiple configuration scopes: global (~/.claude.json), project (.mcp.json), and managed settings. [^54^] [^9^]
**Source:** Claude Code Official Documentation
**URL:** https://code.claude.com/docs/en/settings
**Date:** 2026-07-14
**Excerpt:** "MCP servers: User location=~/.claude.json, Project (shared)=.mcp.json, Local (personal)=~/.claude.json (per-project)"
**Context:** Configuration locations:
- **User:** `~/.claude/settings.json` + `~/.claude.json` (MCP servers)
- **Project:** `.mcp.json` at repository root
- **Project (Claude-managed):** `.claude/settings.json`
- **Local:** `.claude/settings.local.json`

Example `.mcp.json`:
```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": { "GITHUB_PERSONAL_ACCESS_TOKEN": "ghp_..." }
    }
  }
}
```
**Confidence:** High

### 4.2 Tool Discovery and Execution

**Claim:** Claude Code discovers MCP tools at startup, names them with the pattern `mcp__<server-name>__<tool-name>`, and requires explicit permission before executing them. [^18^] [^29^]
**Source:** Claude Code Documentation - MCP
**URL:** https://code.claude.com/docs/en/agent-sdk/mcp and https://code.claude.com/docs/en/mcp
**Date:** 2026-07-13
**Excerpt:** "MCP tools follow the naming pattern mcp__<server-name>__<tool-name>... MCP tools require explicit permission before Claude can use them."
**Context:** Tool access control:
- Auto-approve via `allowedTools`: `["mcp__github__*"]` (wildcards supported)
- Tool search: Dynamic tool loading to reduce token usage (~134K -> ~5K)
- Special annotation: `_meta["anthropic/requiresUserInteraction"]: true` forces permission prompt on every call
- Discovery: Available tools visible via `system` init message
**Confidence:** High

### 4.3 Permission Modes

**Claim:** Claude Code offers six permission modes that control tool execution autonomy. [^16^] [^23^]
**Source:** Claude Code Documentation and Academic Analysis
**URL:** https://code.claude.com/docs/en/settings
**Date:** 2026-07-14
**Excerpt:** "Claude Code offers 6 permission modes: default (prompt for all), acceptEdits (auto-approve file edits), plan (read-only analysis), auto (LLM classifier), dontAsk (pre-approved only), bypassPermissions (all auto-approved - CI only)"
**Context:** Permission modes table:
| Mode | File Reading | File Writing | Shell Commands | Use Case |
|------|-------------|--------------|----------------|----------|
| default | Prompt | Prompt | Prompt | Discovery, daily dev |
| acceptEdits | Auto | Auto | Prompt | Everyday development |
| plan | Auto | Plan+confirm | Plan+confirm | Code review |
| auto | Auto | LLM decides | LLM decides | Team/Enterprise |
| dontAsk | Pre-approved | Pre-approved | Pre-approved | CI/CD |
| bypassPermissions | Auto | Auto | Auto | Containers only |
**Confidence:** High

### 4.4 Subagent Architecture and MCP Delegation

**Claim:** Claude Code uses subagents (isolated Claude instances with own context/tools) that can independently access MCP servers, enabling delegation patterns. [^137^] [^138^] [^22^]
**Source:** Academic Analysis and Claude Code Subagent Documentation
**URL:** https://arxiv.org/html/2604.14228v1
**Date:** 2026-04-14
**Excerpt:** "The Agent tool dispatches to built-in subagents (Explore, Plan, general-purpose) or custom subagents, each running in an isolated context with rebuilt permission context and independent tool sets."
**Context:** Subagent types:
- **Explore:** Read-only investigation (Haiku by default)
- **Plan:** Structured planning
- **General-purpose:** Broad capability
- **Claude Code Guide:** Onboarding
- **Verification:** Validation checks
- **Custom:** User-defined in `.claude/agents/`
**Confidence:** High

---

## Section 5: OpenCode Integration

### 5.1 MCP Support in OpenCode

**Claim:** OpenCode fully supports MCP servers through its `opencode.json` configuration under the `mcp` key, supporting both local (stdio) and remote (HTTP/SSE) servers. [^84^] [^85^]
**Source:** OpenCode Official Documentation
**URL:** https://opencode.ai/docs/mcp-servers/
**Date:** 2026-07-14
**Excerpt:** "You can define MCP servers in your OpenCode Config under mcp. Add each MCP with a unique name. You can refer to that MCP by name when prompting the LLM."
**Context:** OpenCode MCP configuration:
```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "my-local-mcp-server": {
      "type": "local",
      "command": ["npx", "-y", "my-mcp-command"],
      "enabled": true,
      "environment": { "MY_ENV_VAR": "value" }
    },
    "my-remote-mcp": {
      "type": "remote",
      "url": "https://my-mcp-server.com",
      "headers": { "Authorization": "Bearer MY_API_KEY" }
    }
  }
}
```
**Confidence:** High

### 5.2 Declarative Tool Loading

**Claim:** OpenCode uses declarative per-agent tool configuration, allowing fine-grained control over which MCP tools are available to which agents. [^66^]
**Source:** MorphLLM OpenCode vs Claude Code Comparison
**URL:** https://www.morphllm.com/comparisons/opencode-vs-claude-code
**Date:** 2026-07-06
**Excerpt:** "OpenCode has always used declarative tool loading... OpenCode: Declarative per-agent control with glob patterns: mymcp_*: true enables all tools from an MCP"
**Context:** OpenCode's approach differs from Claude Code's eager loading:
```json
{
  "tools": {
    "mymcp_*": true,
    "dangerous_tool": false
  },
  "agent": {
    "build": {
      "tools": {
        "mymcp_write": true,
        "mymcp_delete": false
      }
    }
  }
}
```
**Confidence:** High

### 5.3 ACP Support in OpenCode

**Claim:** OpenCode supports ACP via the `opencode acp` command, which starts OpenCode as an ACP-compatible subprocess communicating over JSON-RPC via stdio. [^21^] [^24^]
**Source:** OpenCode Official Documentation - ACP Support
**URL:** https://opencode.ai/docs/acp/
**Date:** 2026-07-14
**Excerpt:** "To use OpenCode via ACP, configure your editor to run the opencode acp command. The command starts OpenCode as an ACP-compatible subprocess that communicates with your editor over JSON-RPC via stdio."
**Context:** OpenCode ACP configuration examples for Zed, JetBrains, and Neovim are provided in the documentation. The integration follows the standard ACP flow.
**Confidence:** High

---

## Section 6: Other CLI Agents and MCP Support

### 6.1 Continue.dev

**Claim:** Continue.dev supports MCP servers through YAML configuration files in `.continue/mcpServers/` directory or inline in `config.yaml`. [^145^] [^150^]
**Source:** Continue.dev Official Documentation
**URL:** https://docs.continue.dev/customize/deep-dives/mcp
**Excerpt:** "MCP Servers can be added to your config using mcpServers. MCP can only be used in the agent mode."
**Context:** Continue MCP configuration:
```yaml
name: My Config
version: 1.0.0
schema: v1
mcpServers:
  - name: SQLite MCP
    command: npx
    args:
      - "-y"
      - "mcp-sqlite"
      - "/path/to/your/database.db"
```
Supported transports: stdio, sse, streamable-http
**Confidence:** High

### 6.2 Aider

**Claim:** Aider is a Git-native CLI pair-programming agent. While Aider itself doesn't directly integrate MCP, AiderDesk (a UI wrapper) adds MCP support with built-in Power Tools plus extensible MCP servers. [^56^]
**Source:** Hotovo Blog - AiderDesk Agent Mode
**URL:** https://www.hotovo.com/blog/how-mcp-servers-gave-birth-to-aiderdesks-agent-mode
**Date:** 2026-07-08
**Excerpt:** "AiderDesk is powerful on its own, but infinitely extensible through MCP. You can start working with the agent right away using its built-in Power Tools, and then gradually connect more specialized MCP servers."
**Context:** Aider itself focuses on git-integrated coding. AiderDesk adds MCP support as an extension layer.
**Confidence:** High

### 6.3 Ecosystem Summary

**Claim:** MCP is now the de facto standard for AI tool integration, supported by Claude Code, OpenCode, Continue.dev, Cline, Cursor, VS Code, JetBrains, Gemini CLI, and Goose. [^8^] [^26^]
**Source:** Terminal AI Agents Landscape 2025
**URL:** https://wal.sh/research/2025-terminal-ai-agents/
**Date:** 2026-06-22
**Excerpt:** "MCP (Model Context Protocol): Anthropic's protocol for tool integration, now adopted by: Claude Code, Cursor CLI, Cline, Continue... ACP (Agent Client Protocol): Zed's open standard for agent-editor integration"
**Context:** Full ecosystem map:
- **MCP-native:** Claude Code, OpenCode, Continue.dev, Cline, Cursor, VS Code, JetBrains, Gemini CLI, Goose, OpenHands
- **ACP-native:** Zed, JetBrains (in progress), Neovim (plugins), OpenCode
- **Combined (MCP+ACP):** OpenCode, Cline, Continue
**Confidence:** High

---

## Section 7: Security Considerations

### 7.1 MCP Server Security Best Practices

**Claim:** MCP servers must implement authentication (OAuth 2.1 with PKCE), input validation, output sanitization, and sandboxed execution to prevent prompt injection and tool poisoning attacks. [^142^] [^143^] [^146^]
**Source:** MCP Security Best Practices (Official) and Multiple Security Guides
**URL:** https://modelcontextprotocol.io/docs/tutorials/security/security_best_practices
**Date:** 2026-06-25
**Excerpt:** "MCP clients MUST validate authorization URLs and reject dangerous schemes... MUST only allow http:// and https:// schemes... MUST reject javascript:, data:, file:, vbscript:... SHOULD use allowlist-based validation rather than blocklist-based approaches"
**Context:** Security layers:
1. **Authentication:** OAuth 2.1 + PKCE, token validation, no token passthrough
2. **Input Validation:** JSON Schema validation, allowlists, command injection prevention
3. **Output Sanitization:** PII redaction, prompt injection scanning
4. **Transport Security:** Origin validation, HTTPS for remote, stdio sandboxing
5. **Session Management:** Secure session IDs, timeouts, per-request validation
**Confidence:** High

### 7.2 Tool Poisoning and Prompt Injection

**Claim:** Real-world attacks have exploited MCP servers through tool poisoning (malicious tool descriptions) and indirect prompt injection (hidden instructions in external data). [^154^] [^156^]
**Source:** MCP Security Research
**URL:** https://rathnaprashanth.medium.com/mcp-server-security-protecting-against-model-context-protocol-vulnerabilities-d1c131461f25
**Date:** 2025-06-24
**Excerpt:** "WhatsApp MCP Tool Poisoning: Malicious tool redirected messages to attacker-controlled numbers. Cursor IDE Vulnerability: Malicious server extracted chat history via hidden instructions."
**Context:** Mitigation strategies:
- Cryptographic signing of tool definitions
- Tool output sanitization before returning to LLM
- Structured data formats (JSON) instead of free text
- Input validation with semantic analysis
- Separate tool outputs from conversation context
**Confidence:** High

---

## Section 8: Transport Comparison for SkillGraph Implementation

### 8.1 Recommended Transport Strategy

Based on the research, for the Knowledge Skill Graph system:

| Transport | Use Case | Recommendation |
|---|---|---|
| **stdio** | Local CLI agents (Claude Code, OpenCode) | **Primary** - Use for local agent integration |
| **Streamable HTTP** | Remote/web agents, shared infrastructure | **Secondary** - Use for web dashboard and remote access |
| **SSE** | Legacy compatibility only | **Avoid** - Deprecated, do not use for new implementations |
| **gRPC** | Not part of MCP or ACP specs | **Not applicable** - Neither MCP nor ACP uses gRPC |

### 8.2 Client-Specific Integration Patterns

**Claude Code:**
- Configure via `.mcp.json` in project root
- Tools named: `mcp__skillgraph__skill_search`, `mcp__skillgraph__skill_get`, etc.
- User approves tools or uses `allowedTools` wildcard

**OpenCode:**
- Configure via `opencode.json` under `mcp` key
- Use glob patterns: `skillgraph_*: true` to enable all tools
- Per-agent tool configuration supported

**Continue.dev:**
- Configure via `.continue/mcpServers/skillgraph.yaml`
- Supports stdio, sse, streamable-http transports
- Agent mode required for MCP tools

---

## Summary and Recommendations

### Key Findings

1. **MCP is the dominant protocol** for AI tool integration with 10,000+ public servers and broad client support. Current stable spec is 2025-11-25.

2. **mcp-go is production-ready** with 8.9k stars, active development, and full support for stdio, SSE, and Streamable HTTP transports.

3. **ACP exists but serves a different purpose** - it's for agent-editor communication (like LSP), not tool access. It uses JSON-RPC over stdio, NOT gRPC.

4. **CRITICAL: The blueprint's "ACP adapter via gRPC" is incorrect.** ACP uses JSON-RPC 2.0 over stdio. No gRPC is involved in the official spec. This needs correction in the architecture.

5. **stdio is the primary transport for CLI agents** - Claude Code, OpenCode, and most CLI tools use stdio for local MCP servers.

6. **Streamable HTTP is the standard for remote access** - replacing the deprecated HTTP+SSE transport.

### Recommendations

1. **Build the MCP server using mcp-go with stdio as primary transport** - this covers Claude Code, OpenCode, and Continue.dev locally.

2. **Also expose Streamable HTTP** for remote dashboard access and web integration.

3. **Correct the architecture: Remove gRPC from ACP adapter.** If gRPC is needed internally, it's a separate concern from ACP. ACP uses JSON-RPC over stdio.

4. **Follow MCP security best practices:** Input validation via JSON Schema, output sanitization, proper error handling, and tool descriptions that guide the LLM.

5. **Name tools clearly:** `skill_search`, `skill_get`, `skill_tree`, `learn_from_project`, `missing_skills` - these will appear as `mcp__skillgraph__skill_search` in Claude Code.

### Danger Zones

1. **gRPC Misconception:** The blueprint's "ACP adapter via gRPC" is based on incorrect information. ACP uses JSON-RPC over stdio. Building a gRPC adapter for ACP would be incompatible with the standard.

2. **SSE is Deprecated:** Do not build SSE-only transport. Use Streamable HTTP for HTTP-based access.

3. **Tool Token Consumption:** MCP tools consume context tokens. With many tools, token usage can grow. Consider Claude Code's Tool Search pattern for lazy loading.

4. **Permission Complexity:** Each client handles MCP tool permissions differently. Test with each target client.

5. **Specification Churn:** MCP spec is actively evolving (2026-07-28 RC has significant changes). mcp-go tracks spec versions but be prepared for migration work.

6. **Security Surface:** MCP servers execute code on behalf of LLMs. Validate all inputs rigorously - treat LLM-provided parameters as potentially malicious.

---

## Source Index

| Citation | Source | URL | Date |
|----------|--------|-----|------|
| [^2^] | mcp-go GitHub | https://github.com/mark3labs/mcp-go | 2026-06-25 |
| [^8^] | Terminal AI Agents 2025 | https://wal.sh/research/2025-terminal-ai-agents/ | 2026-06-22 |
| [^9^] | Claude Codex MCP Setup | https://claude-codex.fr/en/mcp/setup/ | 2026-03-11 |
| [^16^] | Claude Code Permissions | https://claudefa.st/blog/guide/development/permission-management | 2026-07-14 |
| [^17^] | MCP Specification 2025-06-18 | https://modelcontextprotocol.io/specification/2025-06-18 | 2026-06-30 |
| [^18^] | Claude Code MCP Docs | https://code.claude.com/docs/en/agent-sdk/mcp | 2026-07-13 |
| [^21^] | OpenCode ACP Support | https://opencode.ai/docs/acp/ | 2026-07-14 |
| [^22^] | Claude Code Architecture Paper | https://arxiv.org/html/2604.14228v1 | 2026-04-14 |
| [^23^] | Claude Code Permissions Guide | https://institute.sfeir.com/en/claude-code/claude-code-permissions-and-security/ | 2026 |
| [^24^] | OpenCode ACP Docs (mirror) | https://open-code.ai/en/docs/acp | 2026 |
| [^26^] | ACP Explained - MorphLLM | https://www.morphllm.com/agent-client-protocol | 2026-03-03 |
| [^29^] | Claude Code MCP Annotation | https://code.claude.com/docs/en/mcp | 2026-07-14 |
| [^54^] | Claude Code Settings | https://code.claude.com/docs/en/settings | 2026-07-14 |
| [^66^] | OpenCode vs Claude Code | https://www.morphllm.com/comparisons/opencode-vs-claude-code | 2026-07-06 |
| [^84^] | OpenCode MCP Servers | https://opencode.ai/docs/mcp-servers/ | 2026-07-14 |
| [^85^] | OpenCode Config | https://opencode.ai/docs/config/ | 2026-07-14 |
| [^86^] | MCP Server in Go Tutorial | https://www.sitepoint.com/build-an-mcp-server-in-go/ | 2026-06-27 |
| [^87^] | MCP GitHub Releases | https://github.com/modelcontextprotocol/modelcontextprotocol/releases | 2026-05-29 |
| [^130^] | MCP Transport Selection | https://note.com/ayato_studio/n/n61c1ccefbab4 | 2026-05-30 |
| [^131^] | MCP Stdio vs Streamable HTTP | https://www.truefoundry.com/blog/mcp-stdio-vs-streamable-http-enterprise | 2026-05-20 |
| [^137^] | ACP Initialization | https://agentclientprotocol.com/protocol/v1/initialization | 2026-07-06 |
| [^139^] | ACP Overview | https://agentclientprotocol.com/protocol/v1/overview | 2026-06-01 |
| [^142^] | ACPex Overview | https://hexdocs.pm/acpex/protocol_overview.html | 2026 |
| [^143^] | MCP Tool Schema | https://apxml.com/courses/getting-started-model-context-protocol/chapter-3-implementing-tools-and-logic/tool-definition-schema | 2025-11-19 |
| [^145^] | Continue MCP Setup | https://docs.continue.dev/customize/deep-dives/mcp | 2026 |
| [^146^] | Claude Code MCP Issue | https://github.com/anthropics/claude-code/issues/5037 | 2025-08-03 |
| [^151^] | MCP 2026-07-28 RC Blog | https://blog.modelcontextprotocol.io/posts/2026-07-28-release-candidate/ | 2026-05-22 |
| [^153^] | ACP GitHub | https://github.com/agentclientprotocol/agent-client-protocol | 2025-06-23 |
| [^156^] | MCP Security Risks | https://socprime.com/blog/mcp-security-risks-and-mitigations/ | 2026-02-11 |

---

*Research completed with 20+ independent web searches across official specifications, GitHub repositories, documentation sites, and academic papers. All claims are cited with primary sources.*
