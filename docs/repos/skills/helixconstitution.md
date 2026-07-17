# HelixConstitution

> **Repo:** [HelixDevelopment/HelixConstitution](https://github.com/HelixDevelopment/HelixConstitution.git)
> **Type:** Git submodule · **Status:** Active · **Path:** `constitution/`

## What it provides

Universal agent rules, action-prefix registry, MCP configs, skills, plugins,
scripts, hooks, gates, reporting engine, feature-scheduling engine, guards.

## Skills

- action-prefix-system
- media-validator
- multitrack
- reporting-workable-items
- scheduled-work-queue
- session-sync
- workable-item-lifecycle

## MCP Servers

- media-validator-mcp
- scheduled-work-mcp

## Plugins

- helix (action directives as slash commands)
- scheduled-work (reminder tracking)

## Reusable Engines

- continuum (content-addressed Merkle store)
- session_orchestrator (alias-health registry, flowing-pool)
- token_optimizer (multi-tier prompt caching)

## Upstream Mirrors

| Mirror | URL |
|---|---|
| GitHub (HelixDevelopment, primary) | `git@github.com:HelixDevelopment/HelixConstitution.git` |
| GitLab (helixdevelopment1) | `git@gitlab.com:helixdevelopment1/helixconstitution.git` |
| GitHub (vasic-digital) | `git@github.com:vasic-digital/HelixConstitution.git` |
| GitLab (vasic-digital) | `git@gitlab.com:vasic-digital/HelixConstitution.git` |
| GitFlic | `git@gitflic.ru:helixdevelopment/helixconstitution.git` |
| GitVerse | `git@gitverse.ru:helixdevelopment/HelixConstitution.git` |

## Dependencies

Zero own-org deps for core engines (Go stdlib only). Shared constitution
scripts for session_orchestrator.

## Status

Active. Primary submodule consumed by all Helix Skills projects via
`.gitmodules`.
