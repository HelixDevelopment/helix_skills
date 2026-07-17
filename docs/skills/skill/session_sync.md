# session-sync

> **GENERATED FILE — DO NOT HAND-EDIT.** Regenerated from the live skill
> graph by the `skills-catalog` generator. Edit the skill via CLI/REST/MCP
> (see `docs/scripts/` / `docs/API.md`) — this file will be overwritten.

<!-- skills-catalog:section=header -->
## Header

- **Name:** session-sync
- **Title:** Session State Persistence and Cross-Track Sync
- **Version:** 0.1.0
- **Kind:** atomic
- **Status:** active
- **Domain:** agent-infrastructure
- **Complexity:** advanced
- **Tags:** session, sync, persistence, cross-track, state

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Session state persistence and cross-track sync. Wires `session-sync.sh`.
Enables agents to persist session state across invocations and synchronize
state between parallel tracks in multi-track workflows.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

- `continuum` — for content-addressed state storage

### Optional

- `multitrack` — for cross-track synchronization

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/session-sync/` | Directory | Skill source |
| `session-sync.sh` | Script | Sync hook script |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Source:** `constitution/skills/session-sync/`
- **Consumed via:** Constitution skill (installed via register.sh)
- **Constitution references:** §11.4.207
- **Created:** 2026-07-15
- **Last updated:** 2026-07-17

<!-- /skills-catalog:section=metadata -->
