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

## Usage

### Moving an item through its lifecycle

```
# Start work on a queued item
# → Status changes: Queued → In progress

# Mark ready for testing
# → Status changes: In progress → Ready for testing

# Close a Bug with evidence
# → Status: Fixed (→ Fixed.md) + captured evidence path

# Close a Feature
# → Status: Implemented (→ Fixed.md) + captured evidence path

# Close a Task
# → Status: Completed (→ Fixed.md) + captured evidence path

# Reopen a closed item
# → Status: Reopened + reopen reason + evidence
```

### Status state machine

```
Queued → In progress → Ready for testing → In testing → [terminal closure]
                  ↕                              ↓
            Operator-blocked                 Reopened
                  ↕
               Obsolete
```

### Closure vocabulary (type-aware, §11.4.33)

| Type | Correct closure | Wrong (violation) |
|---|---|---|
| Bug | `Fixed (→ Fixed.md)` | ~~Implemented~~, ~~Completed~~ |
| Feature | `Implemented (→ Fixed.md)` | ~~Fixed~~, ~~Completed~~ |
| Task | `Completed (→ Fixed.md)` | ~~Fixed~~, ~~Implemented~~ |

### Reopening an item

```
# A reopen requires:
# - By: AI or User
# - On: ISO date
# - Reason: test-failed | manual-testing-detected | captured-evidence-contradicts |
#           end-user-report | cycle-re-discovered | design-reconsidered
# - Evidence: path to the captured artefact
```

## Constitution References

| Reference | Meaning |
|---|---|
| **§11.4.93** | SQLite-backed SSoT. The DB at `docs/workable_items.db` is the authoritative source. Issues/Fixed/summaries are generated output — never hand-edit them. |
| **§11.4.95** | DB tracked in git. The SQLite DB is NEVER gitignored. Committed alongside every state change. |
| **§11.4.202** | Reporting directives. Defines how items are created (ISSUE/BUG/TASK). Lifecycle management starts after creation. |
| **§11.4.15** | Status closed set. Queued, In progress, Ready for testing, In testing, Reopened, Operator-blocked, terminal closure, Obsolete. |
| **§11.4.16** | Type closed set. Bug / Feature / Task. No fourth value. |
| **§11.4.33** | Type-aware closure vocabulary. Each Type has exactly one correct closure status literal. |
| **§11.4.34** | Reopened-source attribution. Every reopen records By (AI/User), On (date), Reason (closed vocabulary), Evidence. |
| **§11.4.55** | Reopens-history tracking. Items with `reopens_count > 0` carry a per-item reopen history doc. High reopen count = strongest fragility signal. |
| **§11.4.90** | Obsolete status + per-item obsolescence audit. Reasons: superseded-by-design-change, superseded-by-later-mandate, feature-removed, duplicate-of, unsupported-topology, not-reproducible. Removing an existing capability requires asking the operator first (§11.4.122). |
| **§11.4.21** | Operator-blocked status. A LAST resort after exhausting every self-resolution path. Must state WHAT / WHY (each exhausted alternative) / UNBLOCK CONDITION / WHO. |
| **§11.4.108** | Four-layer fix-verification. "Done" means the runtime signature verifies on a clean deployment. Source-committed ≠ artifact-contains-it ≠ active-on-clean-target ≠ works-for-the-user. |
| **§11.4.123** | Rock-solid-proof-or-deep-research. If unsure how to validate, that is a research trigger, not a licence to accept a weak PASS. |
| **§11.4.149** | Per-item testing diary. Append-only, one entry per test event. A PASS entry without an evidence path is rejected by the schema. |
| **§11.4.135** | Standing regression-guard suite. Fix ⇒ a permanent regression guard, authored RED on the broken artifact and flipped GREEN on the fixed one. |
| **§11.4.197** | Nothing sits un-wired forever. Every started effort reaches a terminal state: fully COMPLETED-and-wired or explicitly evidence-backed CLOSED. No third state. |
| **§11.4.54** | ATM-NNN ticket identifier. Stable, auto-incremented IDs for every workable item. |
| **§11.4.91** | Summary-doc clarity. Bare fragments forbidden as descriptions. |
| **§11.4.115** | RED-baseline-on-the-broken-artifact + polarity-switch. Regression guards are authored RED first, then flipped GREEN. |

## Cross-links

- **Optional:** [`action-prefix-system`](action_prefix_system.md) (directive-driven item creation), [`reporting-workable-items`](reporting_workable_items.md) (ISSUE/BUG/TASK creation)
- **Related skills:** [`reporting-workable-items`](reporting_workable_items.md) (creates items that enter this lifecycle), [`action-prefix-system`](action_prefix_system.md) (BACKGROUND/REMINDER/ISSUE/BUG/TASK directives)
- **Parent domain:** [`agent-infrastructure`](../by-domain/agent-infrastructure.md)
- **Constitution source:** [`constitution/skills/workable-item-lifecycle/`](../../../constitution/skills/workable-item-lifecycle/)

## Integration

| Surface | How it hooks in |
|---|---|
| **SQLite SSoT** | All lifecycle state lives in `docs/workable_items.db` (§11.4.93). Status changes are DB mutations followed by doc regeneration — never hand-edits to the derived docs. |
| **Doc regeneration** | Every status change flows DB → regenerate docs (Issues.md, Fixed.md, summaries, HTML/PDF/DOCX) in the same commit. The `workable-items` Go binary handles `sync db-to-md`. |
| **report_item.sh** | Item creation (the entry point to the lifecycle) is driven by the §11.4.202 engine. The lifecycle skill picks up from there. |
| **Testing diary** | Each item carries an append-only testing diary (§11.4.149). PASS entries require an evidence path — the anti-bluff mechanism. |
| **Regression guards** | Every fix gets a permanent regression guard (§11.4.135): authored RED on the broken artifact, flipped GREEN on the fixed one. The guard lives beyond the item's closure. |
| **External trackers** | Bidirectional sync to GitHub Issues, GitLab, Jira, etc. (§11.4.148 D5). Status changes propagate to all connected trackers. |
