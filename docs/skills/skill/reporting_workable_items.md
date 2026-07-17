# reporting-workable-items

> **GENERATED FILE — DO NOT HAND-EDIT.** Regenerated from the live skill
> graph by the `skills-catalog` generator. Edit the skill via CLI/REST/MCP
> (see `docs/scripts/` / `docs/API.md`) — this file will be overwritten.

<!-- skills-catalog:section=header -->
## Header

- **Name:** reporting-workable-items
- **Title:** Reporting Directives for Workable Items
- **Version:** 0.1.0
- **Kind:** atomic
- **Status:** active
- **Domain:** agent-infrastructure
- **Complexity:** intermediate
- **Tags:** reporting, issue, bug, task, workable-items

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Reporting directives ISSUE/BUG/TASK (§11.4.202) — creates fully-populated,
fully-synced workable items from plain-language reports. Integrates with the
action-prefix system and workable-item lifecycle management.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

- `action-prefix-system` — for directive parsing
- `workable-item-lifecycle` — for status transitions

### Optional

- `scheduled-work-queue` — for background-queue integration

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/reporting-workable-items/` | Directory | Skill source |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Source:** `constitution/skills/reporting-workable-items/`
- **Consumed via:** Constitution skill (installed via register.sh)
- **Constitution references:** §11.4.202
- **Created:** 2026-07-15
- **Last updated:** 2026-07-17

<!-- /skills-catalog:section=metadata -->
