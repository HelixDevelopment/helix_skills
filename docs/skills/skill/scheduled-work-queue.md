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
- **Tags:** scheduling, queue, feature, background, reminders

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Scheduled-work and background-queue management. Manages
`docs/requests/background_queue.md` for BACKGROUND/REMINDER action-prefix
entries (§11.4.140) and `docs/requests/feature_queue.md` for FEATURE
directive scheduling (§11.4.213). Backed by the scheduled-work-engine
Go MCP server. Enables the autonomous loop (§11.4.87/§11.4.126) to
claim and drive scheduled items to completion.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

- `reporting-workable-items` — for item creation on FEATURE directive

### Optional

- `session-sync` — for cross-track queue visibility

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/scheduled-work-queue/` | Directory | Skill source |
| `constitution/scripts/scheduled-work-engine/` | Directory | Go MCP server |
| `constitution/mcp/scheduled-work-mcp.json` | File | MCP server definition |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Created:** 2026-07-16
- **Updated:** 2026-07-16
- **Author:** HelixDevelopment
- **Source:** `constitution/skills/scheduled-work-queue/`

<!-- /skills-catalog:section=metadata -->
