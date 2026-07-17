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

## Usage

### Recognizing and expanding action prefixes

When a user prompt starts with a registered action token, the agent replaces the prefix with the action's registered expansion text and executes the remainder under that expansion.

```
# Form 1 — double-colon separator
BACKGROUND :: refactor the parser to support stacked prefixes

# Form 6 — single-colon (most natural for humans)
BUG: subtitles render one frame late on API 35 devices

# Form 3 — slash command
/task update the status doc fleet

# Stacked prefixes (outer-to-inner, left-to-right)
CRITICAL :: BACKGROUND :: fix the crash in session-sync
```

### Registering a new action

Add a row to `constitution/actions/registry.yaml`:

```yaml
actions:
  - name: REVIEW
    summary: "Dispatch an independent code review"
    expansion: "Dispatch an independent code review on the following changes..."
    rules:
      - "Review MUST run on Fable at xhigh effort (§11.4.209)"
    slash_bare: auto
    slash_conflicts: []
```

Then regenerate agent prefix commands:

```bash
bash constitution/scripts/generate_agent_prefix_commands.sh
```

The generator runs automatically on every constitution pull (§11.4.164).

### Escaping a prefix (discuss without invoking)

```
\BACKGROUND :: this is a literal example, not a real directive
```

The leading backslash strips and produces no expansion.

## Constitution References

| Reference | Meaning |
|---|---|
| **§11.4.140** | Universal action-prefix system (`ACTION_NAME ::`). Defines the six equivalent grammar forms, the two-layer architecture (LAYER 1: agent-context-carrier recognition; LAYER 2: mechanical `UserPromptSubmit` hook), the registry-driven extensibility model, and the conflict-resolution rules for bare `/ACTION` vs host slash commands. |
| **§11.4.6** | No-guessing mandate. Unknown tokens matching the grammar shape MUST NOT be silently expanded — either ask or treat literally. Inventing an expansion is a §11.4.6 violation. |
| **§11.4.29** | Lowercase snake_case naming. File/directory names stay lowercase; the action TOKEN itself is UPPERCASE by grammar rule. |
| **§11.4.164** | Constitution auto-propagation hook. The `post_update_hook.sh` re-runs the prefix-command generator on every constitution pull, so new actions are live out of the box. |
| **§11.4.66** | Blocker-resolution interactive clarification. When an unknown token is encountered in forms 1–5, the agent asks which registered action was intended. |
| **§11.4.105** | Natural-language intent recognition. Supports the clarification flow when a token is ambiguous. |
| **§11.4.109** | Anti-forgetting upgrade. LAYER 2 (mechanical hook) ensures expansion holds even if model recall lapses. |

## Cross-links

- **Related skills:** [`workable-item-lifecycle`](workable_item_lifecycle.md) (status transitions after ISSUE/BUG/TASK expansion), [`reporting-workable-items`](reporting_workable_items.md) (ISSUE/BUG/TASK create tracked items), [`scheduled-work-queue`](scheduled_work_queue.md) (REMINDER/BACKGROUND dispatch)
- **Parent domain:** [`agent-infrastructure`](../by-domain/agent-infrastructure.md)
- **Constitution source:** [`constitution/skills/action-prefix-system/`](../../../constitution/skills/action-prefix-system/)
- **Action registry:** [`constitution/actions/registry.yaml`](../../../constitution/actions/registry.yaml)

## Integration

| Surface | How it hooks in |
|---|---|
| **CLI (all agents)** | LAYER 1: the action-prefix recognition instruction is mirrored into every agent context carrier (`CLAUDE.md`, `AGENTS.md`, `QWEN.md`, `GEMINI.md`). The agent self-applies the prefix on every CLI agent — Claude Code, Gemini CLI, Qwen Code, Codex CLI, Cursor, Aider, Cline, etc. |
| **Claude Code hooks** | LAYER 2: a `UserPromptSubmit` / `UserPromptExpansion` hook reads the registry and injects the expansion via `additionalContext`. This is the §11.4.109 anti-forgetting upgrade — the expansion holds even if model recall lapses. |
| **Slash commands** | The generator `generate_agent_prefix_commands.sh` produces per-agent slash-command equivalents from the registry: Claude Code plugin commands, Gemini/Qwen `.toml`, Codex prompt `.md`. |
| **Sub-system shortcuts** | A grammar-shaped UPPERCASE token that is not a behavioral action may name an incorporated sub-system/submodule. It expands to a sub-system context injection (repo + org + checkout + decoupling rules). |
| **Registration lifecycle** | Adding an action = one registry row + re-run generator. The `post_update_hook.sh` (§11.4.164) re-runs the generator automatically on every constitution pull. |
