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

## Usage

### Recording background work

```
# Via action prefix
BACKGROUND :: run the full test suite overnight and report results

# The agent creates a work item in the scheduled-work queue:
#   title: "Full test suite — overnight run"
#   status: in-progress
#   due_at: tomorrow 09:00
```

### MCP tool calls

```javascript
// Create a work item
create_work_item({
  title: "Overnight test suite",
  description: "Run full test suite and capture evidence",
  status: "in-progress",
  due_at: "2026-07-18T09:00:00Z",
  tags: ["testing", "background"]
})

// On REMINDER: check what needs verification
list_needs_verification({})
// → returns items with status blocked/uncertain/overdue

// After verifying real outcome
mark_work_item_done({
  id: "item-uuid",
  notes: "All 847 tests passed. Evidence: docs/qa/2026-07-18-overnight/"
})
```

### REMINDER flow

```
# When a REMINDER fires, the agent:
# 1. Calls list_needs_verification
# 2. For each item, re-checks the real system state
# 3. Only marks done with captured evidence
# 4. If not done, updates status back to in-progress/blocked
```

## Constitution References

| Reference | Meaning |
|---|---|
| **§11.4.140** | Universal action-prefix system. REMINDER and BACKGROUND are registered actions that dispatch work to this queue. |
| **§11.4.6** | No-guessing. An item left `uncertain` is the honest state. Never mark done on assumption — verify the real outcome first. |
| **§11.4.108** | Four-layer fix-verification. Before marking done, the runtime signature must verify on a clean deployment. Source-committed does not mean working. |
| **§11.4.87** | Endless-loop autonomous work. The autonomous loop must confirm no item is still blocked/uncertain/overdue before declaring its queue empty. |
| **§11.4.126** | Default autonomous-loop working mode. The loop processes every captured request; done-condition must not read satisfied while any item is un-processed. |
| **§11.4.94** | Zero-idle priority-first. If an item did not succeed, update status back to in-progress and keep working. |
| **§11.4.28** | Decoupled, project-agnostic submodule. The scheduled-work engine is inherited by reference, never reimplemented. |
| **§11.4.116** | Real-time conductor↔autonomous-test-framework sync channel. Sibling discipline for agent synchronization. |

## Cross-links

- **Requires:** [`action-prefix-system`](action_prefix_system.md) (REMINDER/BACKGROUND directive parsing)
- **Optional:** [`workable-item-lifecycle`](workable_item_lifecycle.md) (status transitions)
- **Related skills:** [`reporting-workable-items`](reporting_workable_items.md) (ISSUE/BUG/TASK creation), [`action-prefix-system`](action_prefix_system.md) (directive expansion)
- **Parent domain:** [`agent-infrastructure`](../by-domain/agent-infrastructure.md)
- **Constitution source:** [`constitution/skills/scheduled-work-queue/`](../../../constitution/skills/scheduled-work-queue/)

## Integration

| Surface | How it hooks in |
|---|---|
| **MCP server** | `scheduled-work-engine` Go binary exposes MCP tools: `create_work_item`, `list_work_items`, `get_work_item`, `update_work_item_status`, `mark_work_item_done`, `list_overdue_work`, `list_needs_verification`. |
| **REST API** | The engine also exposes a REST surface for non-MCP clients. See `scripts/scheduled-work-engine/README.md`. |
| **Action-prefix hooks** | REMINDER and BACKGROUND are §11.4.140 registered actions. The `UserPromptSubmit` hook expands them; the agent then queries the queue. |
| **Background queue file** | `docs/requests/background_queue.md` — human-readable queue storage alongside the MCP server's database. |
| **Registration** | `mcp/scheduled-work-mcp.json` (project-scoped) or Claude Code plugin `plugins/scheduled-work/`. Auto-registered by §11.4.164 post_update_hook. |
