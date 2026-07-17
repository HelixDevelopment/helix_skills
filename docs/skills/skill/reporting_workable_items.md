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

## Usage

### Creating a bug report

```
# Via action prefix (form 6 — most natural)
BUG: subtitles render one frame late on API 35 devices

# Via slash command
/helix:bug The settings screen crashes on rotation with 3+ accounts

# The agent creates a fully-populated workable item:
#   Type: Bug
#   Status: Queued
#   ATM-ID: auto-assigned (§11.4.54)
#   Description: WHAT + SCOPE + REPRODUCTION + ACCEPTANCE + EVIDENCE
#   → synced to SQLite DB + regenerated docs + external trackers
```

### Creating a task

```
/task refactor the media-validator to support WAV audio files
```

### Creating an issue (type undecided)

```
/issue the build sometimes fails on CI with no error output
```

The agent classifies into Bug/Feature/Task based on the report content. If ambiguous, it asks (§11.4.66).

### What the engine does

```bash
# The report_item.sh engine drives the full chain:
# 1. Create in SQLite SSoT with stable ATM-ID
# 2. Regenerate Issues.md, Fixed.md, summaries, HTML/PDF/DOCX
# 3. Push to every configured external tracker
bash constitution/scripts/reporting/report_item.sh \
  --type Bug --title "subtitles render one frame late" --description "..."
```

## Constitution References

| Reference | Meaning |
|---|---|
| **§11.4.202** | Reporting directives ISSUE/BUG/TASK. Defines that every report MUST become a fully-populated, fully-synced workable item. A report that is only acknowledged in prose is a lost requirement. |
| **§11.4.148** | Workable-item integrity. The item must carry status+type+id, a comprehensive structured description, and bidirectional external-tracker sync. |
| **§11.4.171** | Comprehensive structured description. A team member not in the room must understand the item without reading code. |
| **§11.4.16** | Type closed set. Bug / Feature / Task. No fourth value. `/issue` is an entry point, not a type. |
| **§11.4.54** | ATM-NNN ticket identifier. Stable, auto-incremented IDs for every workable item. |
| **§11.4.93** | SQLite-backed SSoT. The DB is the authoritative source; all docs are generated output. |
| **§11.4.95** | DB tracked in git. The SQLite DB at `docs/workable_items.db` is NEVER gitignored. |
| **§11.4.197** | No requirements lost. A report accepted but not tracked is a §11.4.197 violation. |
| **§11.4.6** | No-guessing. Classification is stated as FACT from the report's content; if ambiguous, ASK. |
| **§11.4.91** | Summary-doc clarity. Bare fragments like "Composes with" or "Critical" are forbidden as descriptions. |
| **§11.4.10** | Honest skipping. A tracker whose credentials are absent is SKIPPED with a reason — never faked. |

## Cross-links

- **Requires:** [`action-prefix-system`](action_prefix_system.md) (ISSUE/BUG/TASK directive parsing), [`workable-item-lifecycle`](workable_item_lifecycle.md) (status transitions after creation)
- **Optional:** [`scheduled-work-queue`](scheduled_work_queue.md) (background-queue integration)
- **Parent domain:** [`agent-infrastructure`](../by-domain/agent-infrastructure.md)
- **Constitution source:** [`constitution/skills/reporting-workable-items/`](../../../constitution/skills/reporting-workable-items/)

## Integration

| Surface | How it hooks in |
|---|---|
| **Action-prefix system** | ISSUE/BUG/TASK are registered actions in `constitution/actions/registry.yaml`. The `UserPromptSubmit` hook expands them before the agent sees the prompt. |
| **report_item.sh engine** | The constitution script drives: (1) SQLite insert with ATM-ID, (2) doc regeneration from DB, (3) external-tracker push. Never reimplemented — reused by FEATURE and other directives. |
| **SQLite SSoT** | Items are created in `docs/workable_items.db` (§11.4.93). The DB is tracked in git (§11.4.95). All Markdown/HTML/PDF/DOCX surfaces are generator output. |
| **External trackers** | Bidirectional sync to GitHub Issues, GitLab, Jira, etc. (§11.4.148 D5). Absent trackers are honestly skipped (§11.4.10). |
| **Lifecycle handoff** | After creation, the item enters the `workable-item-lifecycle` state machine: Queued → In progress → Ready for testing → terminal closure. |
