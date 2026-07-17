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
- **Tags:** session, sync, continuation, resume, state, merkle

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Session state persistence and cross-track synchronization (§11.4.131,
§11.4.207). Wires session-sync.sh for standing session-resumption file
maintenance and integrates with the continuum engine for content-addressed
Merkle snapshot persistence. Enables instant multi-stream resume in
O(changed) reads via `RestoreAll`. Composes with §12.10 CONTINUATION.md
maintenance.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

_(none — standalone sync mechanism)_

### Optional

- `continuum` — for content-addressed Merkle store persistence

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/session-sync/` | Directory | Skill source |
| `constitution/submodules/continuum/` | Directory | Resume engine |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Created:** 2026-06-07
- **Updated:** 2026-07-16
- **Author:** HelixDevelopment
- **Source:** `constitution/skills/session-sync/`

<!-- /skills-catalog:section=metadata -->
