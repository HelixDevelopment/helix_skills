# action-prefix-system

> **Path:** `constitution/skills/action-prefix-system/`
> **Type:** Skill · **Status:** Active

## What it provides

Action-prefix grammar (§11.4.140) recognition and expansion. Provides hooks
for `UserPromptSubmit` events (`action_prefix_expand.sh`) that parse action
directives from free-form text.

## How consumed

Constitution skill. Installed via `register.sh`. Hooks into UserPromptSubmit.

## Source paths

- Skill: `constitution/skills/action-prefix-system/`
- Hook: `action_prefix_expand.sh`

## Dependencies

None.

## Constitution references

§11.4.140
