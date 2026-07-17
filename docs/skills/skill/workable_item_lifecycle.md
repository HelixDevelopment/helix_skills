# workable-item-lifecycle

> **GENERATED FILE — DO NOT HAND-EDIT.** Regenerated from the live skill
> graph by the `skills-catalog` generator. Edit the skill via CLI/REST/MCP
> (see `docs/scripts/` / `docs/API.md`) — this file will be overwritten.

<!-- skills-catalog:section=header -->
## Header

- **Name:** workable-item-lifecycle
- **Title:** Workable Item Lifecycle Management
- **Version:** 0.1.0
- **Kind:** atomic
- **Status:** active
- **Domain:** agent-infrastructure
- **Complexity:** intermediate
- **Tags:** workable-items, lifecycle, status, transitions, closure

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Workable item lifecycle management (status transitions, closure vocabulary).
Defines the valid state machine for workable items (FILED → IN_PROGRESS →
QUEUED → CLOSED/FIXED/SUPERSEDED/etc.) and enforces transition rules.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

_None._

### Optional

- `action-prefix-system` — for directive-driven item creation
- `reporting-workable-items` — for ISSUE/BUG/TASK creation

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/workable-item-lifecycle/` | Directory | Skill source |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Source:** `constitution/skills/workable-item-lifecycle/`
- **Consumed via:** Constitution skill (installed via register.sh)
- **Constitution references:** §11.4.93, §11.4.95, §11.4.202
- **Created:** 2026-07-15
- **Last updated:** 2026-07-17

<!-- /skills-catalog:section=metadata -->
