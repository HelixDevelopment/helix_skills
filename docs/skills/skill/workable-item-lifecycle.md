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
- **Tags:** workable-items, lifecycle, status, type, tracking, closure

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Workable item lifecycle management (§11.4.15, §11.4.16, §11.4.33).
Defines the seven-state Status closed set {Queued | In progress |
Ready for testing | In testing | Reopened | Operator-blocked | Fixed}
and the three-value Type closed set {Bug | Feature | Task}. Enforces
type-aware closure vocabulary (Fixed/Implemented/Completed). Manages
the SQLite SSoT (§11.4.93/§11.4.95) bidirectional sync with markdown
docs.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

_(none — foundational lifecycle definitions)_

### Optional

- `reporting-workable-items` — for intake via reporting directives
- `action-prefix-system` — for grammar-driven item creation

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/workable-item-lifecycle/` | Directory | Skill source |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Created:** 2026-05-14
- **Updated:** 2026-07-16
- **Author:** HelixDevelopment
- **Source:** `constitution/skills/workable-item-lifecycle/`

<!-- /skills-catalog:section=metadata -->
