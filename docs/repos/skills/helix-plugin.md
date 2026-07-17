# helix (plugin)

> **Path:** `constitution/plugins/helix/`
> **Type:** Plugin · **Status:** Active

## What it provides

Action directives as native slash commands — every §11.4.140 registered
action (BACKGROUND, REMINDER, ISSUE, BUG, TASK, CRITICAL, IMPORTANT, NOTE,
FEATURE) exposed as `/name` with `/helix:name` as the always-unambiguous
escape. Commands are generated from `actions/registry.yaml`.

## How consumed

Claude Code plugin. Wired via `.claude/settings.local.json` or
`constitution/plugins/helix/.mcp.json`.

## Source paths

- Plugin: `constitution/plugins/helix/`
- Registry: `actions/registry.yaml`

## Dependencies

None (reads from the constitution's action registry).

## Constitution references

§11.4.140
