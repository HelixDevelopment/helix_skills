# multitrack

> **GENERATED FILE — DO NOT HAND-EDIT.** Regenerated from the live skill
> graph by the `skills-catalog` generator. Edit the skill via CLI/REST/MCP
> (see `docs/scripts/` / `docs/API.md`) — this file will be overwritten.

<!-- skills-catalog:section=header -->
## Header

- **Name:** multitrack
- **Title:** Multi-Track Parallel-Development Orchestration
- **Version:** 0.1.0
- **Kind:** atomic
- **Status:** active
- **Domain:** agent-infrastructure
- **Complexity:** advanced
- **Tags:** multitrack, parallel, orchestration, worktree, concurrency

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Multi-track parallel-development orchestration (§11.4.176, §11.4.187,
§11.4.192). Wires multitrack config + cwd-hook. Enables multiple agents
to work on independent tracks simultaneously with proper isolation via
git worktrees and flowing-pool claim registry.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

- `session_orchestrator` — for flowing-pool claim registry

### Optional

- `continuum` — for cross-track state persistence

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/multitrack/` | Directory | Skill source |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Source:** `constitution/skills/multitrack/`
- **Consumed via:** Constitution skill (installed via register.sh)
- **Constitution references:** §11.4.176, §11.4.187, §11.4.192
- **Created:** 2026-07-15
- **Last updated:** 2026-07-17

<!-- /skills-catalog:section=metadata -->
