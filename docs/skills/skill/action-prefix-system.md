# action-prefix-system

> **GENERATED FILE — DO NOT HAND-EDIT.** Regenerated from the live skill
> graph by the `skills-catalog` generator. Edit the skill via CLI/REST/MCP
> (see `docs/scripts/` / `docs/API.md`) — this file will be overwritten.

<!-- skills-catalog:section=header -->
## Header

- **Name:** action-prefix-system
- **Title:** Action-Prefix Grammar Recognition and Expansion
- **Version:** 0.1.0
- **Kind:** atomic
- **Status:** active
- **Domain:** agent-infrastructure
- **Complexity:** intermediate
- **Tags:** actions, prefix, grammar, registry, directives

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Universal action-prefix system (§11.4.140). Recognizes six grammar forms
for registered actions (bare `::`, namespaced `::`, bare slash, namespaced
slash, bare arrow, single-colon). Expands action tokens via the registry
at `constitution/actions/registry.yaml`. Includes the reporting directives
ISSUE/BUG/TASK (§11.4.202), severity markers CRITICAL/IMPORTANT/NOTE,
and the FEATURE scheduling directive (§11.4.213).

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

- `workable-item-lifecycle` — for reporting directive item creation
- `reporting-workable-items` — for ISSUE/BUG/TASK engine

### Optional

- `scheduled-work-queue` — for BACKGROUND/REMINDER queue integration

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/action-prefix-system/` | Directory | Skill source |
| `constitution/actions/registry.yaml` | File | Action registry |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Created:** 2026-06-09
- **Updated:** 2026-07-16
- **Author:** HelixDevelopment
- **Source:** `constitution/skills/action-prefix-system/`

<!-- /skills-catalog:section=metadata -->
