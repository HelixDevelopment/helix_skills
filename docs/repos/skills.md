# Skills Repository Catalog

> **Last updated:** 2026-07-17

This document catalogs ALL GitHub/GitLab repositories that contain skill definitions, MCP tool servers, plugins, and reusable engines used by or available to the Helix Skills project. The catalog is organized by consumption mechanism.

---

## Table of Contents

- [Direct Git Submodules](#direct-git-submodules)
- [MCP Tool Servers](#mcp-tool-servers)
- [Claude Code Plugins](#claude-code-plugins)
- [Constitution Skills (installed via register.sh)](#constitution-skills-installed-via-registersh)
- [Constitution Sub-System Registry (alias resolutions)](#constitution-sub-system-registry-alias-resolutions)
- [Constitution Depth-1 Reusable Engines](#constitution-depth-1-reusable-engines)
- [Declared Dependencies (via helix-deps.yaml)](#declared-dependencies-via-helix-depssyaml)
- [Project Mirrors / Upstream Remotes](#project-mirrors--upstream-remotes)
- [Constitution Mirrors / Upstream Remotes](#constitution-mirrors--upstream-remotes)
- [Code Intelligence](#code-intelligence)
- [Full vasic-digital Catalogue Reference](#full-vasic-digital-catalogue-reference)
- [Full HelixDevelopment Catalogue Reference](#full-helixdevelopment-catalogue-reference)

---

## Direct Git Submodules

| Repo | URL | Skills/Tools | How consumed | Source path |
|---|---|---|---|---|
| **HelixConstitution** | `git@github.com:HelixDevelopment/HelixConstitution.git` | Universal agent rules, action-prefix registry, MCP configs, skills, plugins, scripts, hooks, gates, reporting engine, feature-scheduling engine, guards | Git submodule (`.gitmodules`) | `constitution/` |

**Upstream mirrors for HelixConstitution:**

| Mirror | URL |
|---|---|
| GitHub (HelixDevelopment, primary) | `git@github.com:HelixDevelopment/HelixConstitution.git` |
| GitLab (helixdevelopment1) | `git@gitlab.com:helixdevelopment1/helixconstitution.git` |
| GitHub (vasic-digital) | `git@github.com:vasic-digital/HelixConstitution.git` |
| GitLab (vasic-digital) | `git@gitlab.com:vasic-digital/HelixConstitution.git` |
| GitFlic | `git@gitflic.ru:helixdevelopment/helixconstitution.git` |
| GitVerse | `git@gitverse.ru:helixdevelopment/HelixConstitution.git` |

---

## MCP Tool Servers

| Server name | Skills/Tools | How consumed | Definition path |
|---|---|---|---|
| **media-validator** | Validate media files (MP4, PNG, TXT) via OCR/metadata/pattern matching for PASS/FAIL with evidence; implements the §11.4.163 Universal Media Validation mandate | MCP server defined in JSON config, executed as bash script | `constitution/mcp/media-validator-mcp.json` → `constitution/skills/media-validator/media-validator.sh` |
| **scheduled-work** | Track scheduled work/reminders (background-queue entries): create, list, status, overdue, needs-verification, mark-done. Decoupled Go MCP server backing the REMINDER/BACKGROUND action-prefix (§11.4.140) | MCP server; compiled Go binary | `constitution/mcp/scheduled-work-mcp.json` → `constitution/scripts/scheduled-work-engine/bin/scheduled-work` |
| **scheduled-work (plugin)** | Same server, wired via plugin path | MCP server via plugin `.mcp.json` | `constitution/plugins/scheduled-work/.mcp.json` |

**scheduled-work-engine Go module:** `github.com/HelixDevelopment/HelixConstitution/scripts/scheduled-work-engine` (Go module in `constitution/scripts/scheduled-work-engine/go.mod`). Compiled binary at `scripts/scheduled-work-engine/bin/scheduled-work`. Dependencies: gin-gonic, quic-go, brotli, uuid, yaml.

---

## Claude Code Plugins

| Plugin name | Skills/Tools | Source path |
|---|---|---|
| **helix** | Action directives as native slash commands — every §11.4.140 registered action (BACKGROUND, REMINDER, ISSUE, BUG, TASK, CRITICAL, IMPORTANT, NOTE, FEATURE) exposed as `/name` with `/helix:name` as the always-unambiguous escape. Commands are generated from `actions/registry.yaml`. | `constitution/plugins/helix/` |
| **scheduled-work** | Scheduled-work/reminder tracking for autonomous agents — records background-queue items and lets a REMINDER re-verify uncertain/blocked/overdue work before reporting done (§11.4.6/§11.4.108). Wraps the scheduled-work-engine Go MCP server. | `constitution/plugins/scheduled-work/` |

Installed from: `constitution/.claude-plugin/marketplace.json` — marketplace entry named `helix-constitution`, owned by HelixDevelopment.

---

## Constitution Skills (installed via register.sh)

Each skill lives under `constitution/skills/<name>/` with a `register.sh` script that installs it as a Claude Code skill. Inherited by reference per §11.4.28/§11.4.177.

| Skill | Directory | Capabilities |
|---|---|---|
| **action-prefix-system** | `constitution/skills/action-prefix-system/` | Action-prefix grammar (§11.4.140) recognition and expansion; hooks for UserPromptSubmit (action_prefix_expand.sh) |
| **media-validator** | `constitution/skills/media-validator/` | Media file validation via OCR/metadata/pattern matching (§11.4.163). Also exposed as MCP server media-validator. |
| **multitrack** | `constitution/skills/multitrack/` | Multi-track parallel-development orchestration (§11.4.176, §11.4.187, §11.4.192). Wires multitrack config + cwd-hook. |
| **reporting-workable-items** | `constitution/skills/reporting-workable-items/` | Reporting directives ISSUE/BUG/TASK (§11.4.202) — creates fully-populated, fully-synced workable items from plain-language reports |
| **scheduled-work-queue** | `constitution/skills/scheduled-work-queue/` | Scheduled-work/background-queue management. Manages `docs/requests/background_queue.md` |
| **session-sync** | `constitution/skills/session-sync/` | Session state persistence and cross-track sync. Wires session-sync.sh. |
| **workable-item-lifecycle** | `constitution/skills/workable-item-lifecycle/` | Workable item lifecycle management (status transitions, closure vocabulary). |

---

## Constitution Sub-System Registry (alias resolutions)

These repositories are registered in `constitution/actions/registry.yaml` under the `subsystems:` section. Each UPPERCASE token (all §11.4.140 grammar forms) expands to a sub-system context injection. They are NOT currently cloned into this project; they are available as resolution targets for the action-prefix resolver.

| Token | Name | Org | URL | Aliases |
|---|---|---|---|---|
| HELIXOTA | HelixOTA | HelixDevelopment | `git@github.com:HelixDevelopment/HelixOTA.git` | HXOTA |
| HELIXTRACK | HelixTrack | Helix-Track | `git@github.com:Helix-Track/HelixTrack.git` | HXTRACK |
| HELIXQA | HelixQA | HelixDevelopment | `git@github.com:HelixDevelopment/HelixQA.git` | HXQA |
| HELIXCODE | HelixCode | HelixDevelopment | `git@github.com:HelixDevelopment/HelixCode.git` | HXCODE |
| HELIXAGENT | HelixAgent | HelixDevelopment | `git@github.com:HelixDevelopment/HelixAgent.git` | HXAGENT |
| HELIXLLM | HelixLLM | HelixDevelopment | `git@github.com:HelixDevelopment/HelixLLM.git` | HXLLM |
| HELIXMEMORY | HelixMemory | HelixDevelopment | `git@github.com:HelixDevelopment/HelixMemory.git` | HXMEM |
| HELIXSPECIFIER | HelixSpecifier | HelixDevelopment | `git@github.com:HelixDevelopment/HelixSpecifier.git` | HXSPEC |
| LLMSVERIFIER | LLMsVerifier | vasic-digital | `git@github.com:vasic-digital/LLMsVerifier.git` | LLMSVERIFIER |
| CLAUDE_TOOLKIT | ClaudeToolkit | vasic-digital | `git@github.com:vasic-digital/claude-toolkit.git` | CTOOLKIT, CLAUDE_TOOLKIT |

---

## Constitution Depth-1 Reusable Engines

These are depth-1 submodules under `constitution/submodules/` per the §11.4.28(C) carve-out. Each carries a `helix-deps.yaml` manifest and nests zero further own-org submodules.

| Engine | In-repo path | Description | Dependencies |
|---|---|---|---|
| **continuum** | `constitution/submodules/continuum/` | Instant multi-stream resume engine — content-addressed Merkle store for whole-fleet continuation snapshots (§11.4.207). Go module, `github.com/vasic-digital/continuum`. Mirrors: GitHub (`vasic-digital/continuum`) + GitLab (`vasic-digital/continuum`). | Zero own-org deps (Go stdlib only) |
| **session_orchestrator** | `constitution/submodules/session_orchestrator/` | Session orchestration engine — alias-health registry, flowing-pool claim registry (§11.4.176), non-failover scheduler. | Zero own-org deps (Go stdlib only; shared constitution scripts) |
| **token_optimizer** | `constitution/submodules/token_optimizer/` | Token efficiency optimization engine — multi-tier prompt caching, context-aware routing (§11.4.141). | TOON, Embeddings, VectorDB, Normalize, LLMProvider, conversation (see Declared Dependencies table) |
| **clickup_sync** | `constitution/submodules/clickup_sync/` | DESIGN ONLY (Phase 0) — ClickUp bidirectional-sync engine design docs and research. No implementation code, no repo URL found. | N/A (design only) |

**Repo URL for continuum:** `git@github.com:vasic-digital/continuum.git` (explicitly declared in `helix-deps.yaml`)

**Note:** `session_orchestrator`, `token_optimizer`, and `clickup_sync` are design/planning directories within the constitution submodule — they do NOT yet declare their own Git repository URLs in published manifests (only `helix-deps.yaml` declaring local deps).

---

## Declared Dependencies (via helix-deps.yaml)

These are repositories declared as dependencies by the constitution's reusable engines (primarily `token_optimizer`). They are NOT currently cloned into this project — they are declared as submodule dependencies to be resolved from the parent project root per §11.4.28(C).

| Repo | URL | Org | Required by | Why (from helix-deps.yaml) |
|---|---|---|---|---|
| **TOON** | `git@github.com:vasic-digital/TOON.git` | vasic-digital | token_optimizer | WS7 wire-encode |
| **Embeddings** | `git@github.com:vasic-digital/Embeddings.git` | vasic-digital | token_optimizer | WS6-L2 embed |
| **VectorDB** | `git@github.com:vasic-digital/VectorDB.git` | vasic-digital | token_optimizer | WS6-L2 ANN store |
| **Normalize** | `git@github.com:vasic-digital/Normalize.git` | vasic-digital | token_optimizer | Request canonicalization |
| **LLMProvider** | `git@github.com:HelixDevelopment/LLMProvider.git` | HelixDevelopment | token_optimizer | Transport-tier provider interface |
| **conversation** | `git@github.com:vasic-digital/conversation.git` | vasic-digital | token_optimizer | Infinite-context compression |

**Agent alias providers (from multitrack config `config/multitrack/the-factory.yaml`):**

The multitrack config references the following toolkits by alias kind:
- `claude1..claude4` — claude-code **native** CLI aliases (the Claude Code tool itself on this host, configured via multiple `CLAUDE_CONFIG_DIR` profiles)
- `opencode` — **provider** alias (the OpenCode AI Coding Agent toolkit)
- `xiaomi` — **provider** alias
- `kimi-for-coding` — **provider** alias (Kimi AI Coding Assistant)

These do not correspond to separate repositories — they are per-alias configurations of the same `claude` or provider CLIs on the development host.

---

## Project Mirrors / Upstream Remotes

The Helix Skills project itself (`helix_skills`) is pushed to all of these upstreams (declared in `upstreams/*.sh` recipe files):

| Host | URL | Role |
|---|---|---|
| GitHub (HelixDevelopment) | `git@github.com:HelixDevelopment/helix_skills.git` | Primary origin |
| GitLab (helixdevelopment1) | `git@gitlab.com:helixdevelopment1/helix_skills.git` | Mirror |
| GitFlic (helixdevelopment) | `git@gitflic.ru:helixdevelopment/helix_skills.git` | Mirror |
| GitVerse (helixdevelopment) | `git@gitverse.ru:helixdevelopment/helix_skills.git` | Mirror |

---

## Constitution Mirrors / Upstream Remotes

The HelixConstitution submodule is pushed to all of these upstreams (declared in `constitution/upstreams/*.sh`):

| Host | URL | Role |
|---|---|---|
| GitHub (HelixDevelopment) | `git@github.com:HelixDevelopment/HelixConstitution.git` | Primary |
| GitLab (helixdevelopment1) | `git@gitlab.com:helixdevelopment1/helixconstitution.git` | Mirror |
| GitHub (vasic-digital) | `git@github.com:vasic-digital/HelixConstitution.git` | Mirror |
| GitLab (vasic-digital) | `git@gitlab.com:vasic-digital/HelixConstitution.git` | Mirror |
| GitFlic (helixdevelopment) | `git@gitflic.ru:helixdevelopment/helixconstitution.git` | Mirror |
| GitVerse (helixdevelopment) | `git@gitverse.ru:helixdevelopment/HelixConstitution.git` | Mirror |

---

## Code Intelligence

| Tool | How consumed | Reference |
|---|---|---|
| **CodeGraph** | Enabled as MCP server in `.claude/settings.local.json` (`enabledMcpjsonServers: ["codegraph"]`) | Used per §11.4.78/§11.4.79/§11.4.80 for code index and sync automation. Sync scripts: `constitution/scripts/codegraph_sync.sh` |

---

## Full vasic-digital Catalogue Reference

The complete catalogue of 122 repositories in the vasic-digital GitHub organization is maintained in `constitution/submodules-catalogue.md` (§4, pages 56--206). Key repositories relevant to agent skills and tools include:

| Repo | URL | Skills/Tools |
|---|---|---|
| **SkillRegistry** | `https://github.com/vasic-digital/SkillRegistry` | CLI agent skill registration and management for AI agent systems |
| **ToolSchema** | `https://github.com/vasic-digital/ToolSchema` | Generic tool schema definition, validation, and execution for AI agent tool systems |
| **Agentic** | `https://github.com/vasic-digital/Agentic` | Graph-based agentic workflow orchestration |
| **AgentWrapper** | `https://github.com/vasic-digital/AgentWrapper` | Wrap AI CLI Coding Agents in Docker containers |
| **LLMOrchestrator** | `https://github.com/vasic-digital/LLMOrchestrator` | Headless CLI agent management for LLM orchestration |
| **MCP_Module** | `https://github.com/vasic-digital/MCP_Module` | Generic reusable Go module: digital.vasic.mcp |
| **HelixQA** | `https://github.com/vasic-digital/HelixQA` | AI-driven QA orchestration for multi-platform testing (also at HelixDevelopment) |
| **Panoptic** | `https://github.com/vasic-digital/Panoptic` | Automated testing, UI recording, and screenshot capture |
| **VisionEngine** | `https://github.com/vasic-digital/VisionEngine` | Computer vision and LLM Vision for UI analysis and navigation |
| **DocProcessor** | `https://github.com/vasic-digital/DocProcessor` | Documentation processing and feature map extraction for QA automation |
| **challenges** | `https://github.com/vasic-digital/challenges` | (challenges/test banks submodule) |
| **Planning** | `https://github.com/vasic-digital/Planning` | AI planning algorithms: HiPlan, MCTS, Tree of Thoughts |
| **RedTeam** | `https://github.com/vasic-digital/RedTeam` | YAML-driven adversarial prompt fixture harness |
| **containers** | `https://github.com/vasic-digital/containers` | Container orchestration layer (used per §11.4.76) |
| **claude-toolkit** | `https://github.com/vasic-digital/claude-toolkit.git` | Claude Code multi-alias toolkit (registered in registry.yaml) |

Full catalogue at `constitution/submodules-catalogue.md`.

---

## Full HelixDevelopment Catalogue Reference

The complete catalogue of 20 repositories in the HelixDevelopment organization is maintained in `constitution/submodules-catalogue.md` (§5, pages 208--230). Key repositories:

| Repo | URL | Skills/Tools |
|---|---|---|
| **HelixConstitution** | `https://github.com/HelixDevelopment/HelixConstitution` | Constitution, AGENTS.md, CLAUDE.md universal agent rules and constraints |
| **HelixCode** | `https://github.com/HelixDevelopment/HelixCode` | AI Coding Agent |
| **HelixAgent** | `https://github.com/HelixDevelopment/HelixAgent` | LLMs Agent |
| **HelixLLM** | `https://github.com/HelixDevelopment/LLMProvider` | Local running super model / LLM provider interface |
| **HelixMemory** | `https://github.com/HelixDevelopment/HelixMemory` | Unified Cognitive Memory Engine |
| **HelixSpecifier** | `https://github.com/HelixDevelopment/HelixSpecifier` | Spec-Driven Development Fusion Engine |
| **DebateOrchestrator** | `https://github.com/HelixDevelopment/DebateOrchestrator` | Multi-agent debate orchestration library |
| **HelixBuilder** | `https://github.com/HelixDevelopment/HelixBuilder` | AI powered application building pipeline |
| **HelixGitpx** | `https://github.com/HelixDevelopment/HelixGitpx` | Helix Git Proxy eXtended |
| **HelixTranslate** | `https://github.com/HelixDevelopment/HelixTranslate` | Universal ebook translation toolkit |

Full catalogue at `constitution/submodules-catalogue.md`.

---

## Summary Statistics

- **Direct git submodules:** 1 (HelixConstitution)
- **MCP tool servers:** 2 (media-validator, scheduled-work)
- **Claude Code plugins:** 2 (helix, scheduled-work)
- **Constitution skills installed via register.sh:** 7
- **Sub-system registry entries (alias-resolvable):** 10 repositories
- **Depth-1 reusable engines (constitution/submodules/):** 4 (continuum, session_orchestrator, token_optimizer, clickup_sync [design only])
- **Declared dependency repos (via helix-deps.yaml):** 6
- **Code intelligence tools:** 1 (CodeGraph MCP)
- **vasic-digital org catalogue:** 122 repos referenced
- **HelixDevelopment org catalogue:** 20 repos referenced
- **Agent provider aliases (multitrack config):** 4 native + 3 provider
