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
- **Tags:** reporting, issues, bugs, tasks, workable-items, sync

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Reporting directives ISSUE/BUG/TASK (§11.4.202). A plain-language report
is automatically turned into a fully-populated, fully-synced workable item.
The engine (`constitution/scripts/reporting/report_item.sh`) creates the
item in the SQLite SSoT (§11.4.93), regenerates all derived documents
(§11.4.12/§11.4.53/§11.4.65), and pushes to every configured external
tracker (§11.4.148 D5). Absent trackers skip honestly (§11.4.10).

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

- `workable-item-lifecycle` — for status/type lifecycle
- `action-prefix-system` — for grammar recognition

### Optional

- `scheduled-work-queue` — for FEATURE directive scheduling

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/reporting-workable-items/` | Directory | Skill source |
| `constitution/scripts/reporting/report_item.sh` | Script | Item creation engine |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Created:** 2026-07-15
- **Updated:** 2026-07-16
- **Author:** HelixDevelopment
- **Source:** `constitution/skills/reporting-workable-items/`

<!-- /skills-catalog:section=metadata -->
