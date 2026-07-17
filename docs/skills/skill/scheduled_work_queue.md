# scheduled-work-queue

> **GENERATED FILE — DO NOT HAND-EDIT.** Regenerated from the live skill
> graph by the `skills-catalog` generator. Edit the skill via CLI/REST/MCP
> (see `docs/scripts/` / `docs/API.md`) — this file will be overwritten.

<!-- skills-catalog:section=header -->
## Header

- **Name:** scheduled-work-queue
- **Title:** Scheduled-Work and Background-Queue Management
- **Version:** 0.1.0
- **Kind:** atomic
- **Status:** active
- **Domain:** agent-infrastructure
- **Complexity:** intermediate
- **Tags:** scheduled-work, background-queue, reminders, async

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Scheduled-work/background-queue management. Manages
`docs/requests/background_queue.md`. Enables agents to track deferred work,
reminders, and background tasks with proper lifecycle management.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

- `action-prefix-system` — for REMINDER/BACKGROUND directive parsing

### Optional

- `workable-item-lifecycle` — for status transitions

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/scheduled-work-queue/` | Directory | Skill source |
| `docs/requests/background_queue.md` | Data file | Queue storage |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Source:** `constitution/skills/scheduled-work-queue/`
- **Consumed via:** Constitution skill (installed via register.sh)
- **Constitution references:** §11.4.140, §11.4.6, §11.4.108
- **Created:** 2026-07-15
- **Last updated:** 2026-07-17

<!-- /skills-catalog:section=metadata -->
