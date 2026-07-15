# HelixKnowledge Skill Graph System — Execution Plan

## Objective
Build a comprehensive, self-growing Knowledge Skill Graph system for AI CLI agents, fully written in Go, with recursive skill dependencies, auto-growth, validation, and real-time integration via MCP/ACP. Target technologies: Android AOSP, Java, Kotlin, C++, and recursively all dependencies.

## Stage 1: Deep Research & Validation
**Goal**: Validate feasibility, research cutting-edge technologies, identify gaps and danger zones.
**Skills**: deep-research-swarm
**Agents**:
- Agent 1a: Research CodeGraph engines, AST parsers, code analysis tools (tree-sitter, SourceGraph, GraphRAG alternatives)
- Agent 1b: Research memory systems for LLM agents (episodic memory, vector stores, MemGPT, etc.)
- Agent 1c: Research MCP (Model Context Protocol) spec, ACP spec, and CLI agent integration patterns
- Agent 1d: Research vector databases (pgvector, Milvus, Weaviate, Qdrant) — validate pgvector choice
- Agent 1e: Research Go HTTP/3 libraries (quic-go), Brotli compression, TOML serialization performance
- Agent 1f: Research validation techniques for preventing LLM hallucination in knowledge systems (LLM jury, sandboxing)

**Output**: Validated research brief with technology choices confirmed/modified, gaps identified.

## Stage 2: Architecture & Design Documents
**Goal**: Produce comprehensive architecture docs, data models, API specs, system diagrams.
**Agents**:
- Agent 2a: Create detailed architecture document with all components, data flows, state machines
- Agent 2b: Design complete SQL schema with migrations, indexes, constraints
- Agent 2c: Design REST API specification (OpenAPI/Swagger) with all endpoints
- Agent 2d: Design MCP/ACP protocol handlers and tool schemas
- Agent 2e: Create system diagrams (Mermaid: architecture, data flow, skill graph, state machines)

**Output**: architecture.md, schema.sql, api_spec.yaml, mcp_protocol.md, diagrams/

## Stage 3: Core Go Implementation
**Goal**: Build the complete Go codebase with all modules.
**Agents** (parallel where possible):
- Agent 3a: Project bootstrap — go.mod, config, shared packages, database client (pgvector)
- Agent 3b: Skill domain — CRUD, dependency graph, recursive CTE queries, validation
- Agent 3c: API server — Gin with HTTP/3 adapter, Brotli middleware, TOML/JSON handlers
- Agent 3d: MCP server implementation — tools, stdio/SSE transport
- Agent 3e: ACP adapter — gRPC server, protocol translation
- Agent 3f: Auto-growth worker — gap detection, draft generation, validation pipeline
- Agent 3g: Code analysis worker — tree-sitter integration, pattern extraction
- Agent 3h: CLI (Cobra) — all commands for skill management
- Agent 3i: TUI (Bubble Tea) — skill browser, dependency viewer, search
- Agent 3j: Registry & audit — central registrar, health checks, integrity reports

**Output**: Complete Go source code in skill-system/

## Stage 4: Deployment & Operations
**Goal**: Docker/Podman Compose stack, systemd integration, lifecycle scripts.
**Agents**:
- Agent 4a: Dockerfile, docker-compose.yml, .env template
- Agent 4b: Systemd service unit, install/start/stop/restart/status scripts
- Agent 4c: Makefile for build, test, lint, package

**Output**: All deployment artifacts

## Stage 5: Integration & Validation
**Goal**: Integration tests, end-to-end validation, MCP client testing.
**Agents**:
- Agent 5a: Integration tests (API, DB, MCP)
- Agent 5b: End-to-end tests (skill creation, dependency resolution, auto-growth)
- Agent 5c: Seed data — initial Android AOSP skill tree with real dependencies

**Output**: Test suites, seed data, validation reports

## Stage 6: Documentation
**Goal**: Comprehensive documentation with nano-level details.
**Agents**:
- Agent 6a: README.md with quick start, architecture overview
- Agent 6b: INSTALL.md with step-by-step setup guide
- Agent 6c: API documentation (generated + manual)
- Agent 6d: Developer guide — contributing, skill format, validation pipeline
- Agent 6e: Danger zones & mitigation guide (gaps, weak spots, resolved issues)
- Agent 6f: Operations guide — monitoring, troubleshooting, backup/restore

**Output**: Complete documentation suite

## Stage 7: Final Packaging
**Goal**: Generate distributable archives.
**Output**: skill-system.tar.gz, skill-system.zip in /mnt/agents/output/

## Key Design Decisions (from user's blueprint)
- Go 1.22+ as main language
- Gin + quic-go for HTTP/3 API
- PostgreSQL 16 + pgvector for relational + vector needs
- TOML primary serialization, JSON fallback
- tree-sitter for code analysis
- mcp-go for MCP server
- Bubble Tea + Cobra for CLI/TUI
- Docker/Podman Compose + systemd user services
- Sandboxed execution + LLM jury for zero-bluff validation
