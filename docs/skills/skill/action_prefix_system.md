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
- **Tags:** action-prefix, grammar, hooks, userpromptsubmit, expansion

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Action-prefix grammar (§11.4.140) recognition and expansion. Provides hooks
for `UserPromptSubmit` events (`action_prefix_expand.sh`) that parse action
directives (BACKGROUND, REMINDER, ISSUE, BUG, TASK, CRITICAL, IMPORTANT,
NOTE, FEATURE) from free-form text and expand them into structured workable
items.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

_None._

### Optional

- `workable-item-lifecycle` — for status-transition validation
- `reporting-workable-items` — for ISSUE/BUG/TASK creation

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/action-prefix-system/` | Directory | Skill source |
| `actions/registry.yaml` | Config | Registered action definitions |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Source:** `constitution/skills/action-prefix-system/`
- **Consumed via:** Constitution skill (installed via register.sh)
- **Constitution references:** §11.4.140
- **Created:** 2026-07-15
- **Last updated:** 2026-07-17

<!-- /skills-catalog:section=metadata -->
