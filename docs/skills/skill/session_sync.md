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
- **Tags:** session, sync, persistence, cross-track, state

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Session state persistence and cross-track sync. Wires `session-sync.sh`.
Enables agents to persist session state across invocations and synchronize
state between parallel tracks in multi-track workflows.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

- `continuum` — for content-addressed state storage

### Optional

- `multitrack` — for cross-track synchronization

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/session-sync/` | Directory | Skill source |
| `session-sync.sh` | Script | Sync hook script |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Source:** `constitution/skills/session-sync/`
- **Consumed via:** Constitution skill (installed via register.sh)
- **Constitution references:** §11.4.207
- **Created:** 2026-07-15
- **Last updated:** 2026-07-17

<!-- /skills-catalog:section=metadata -->

## Usage

### Pulling session state from a remote workstation

```bash
# Pull everything from remote (DEFAULT mode)
bash constitution/skills/session-sync/session-sync.sh \
  user@workstation /path/to/project

# Quick mode — recent sessions + memories only (faster)
bash constitution/skills/session-sync/session-sync.sh \
  --quick user@workstation /path/to/project

# Preview without copying
bash constitution/skills/session-sync/session-sync.sh \
  --dry-run user@workstation /path/to/project
```

### Pushing local state to remote

```bash
bash constitution/skills/session-sync/session-sync.sh \
  --push user@workstation /path/to/project
```

### Bidirectional sync

```bash
# Pull first, then push (full two-way)
bash constitution/skills/session-sync/session-sync.sh \
  --bidirectional user@workstation /path/to/project
```

### What gets synced

| Data | Description |
|---|---|
| Memories | Project-specific memory files (`.md` with frontmatter) |
| Session JSONL | Full Claude Code session transcripts |
| Session dirs | Per-session working directories |
| Settings | `.claude/settings.json` + `settings.local.json` |
| Agents | `.claude/agents/` — custom agent definitions |
| Remember logs | `.remember/` — session handoff + daily logs |
| Handoff docs | `docs/SESSION_RESUME.md` + `docs/CONTINUATION.md` |
| Provider config | Global `CLAUDE.md`, `history.jsonl`, `settings.json` |

## Constitution References

| Reference | Meaning |
|---|---|
| **§11.4.207** | Instant multi-stream resume engine. The continuation mechanism is extended by a content-addressed Merkle store so a fresh session resumes ALL work at once, instantly and at token cost O(changed streams). Each stream's state bytes are named by sha256; a snapshot is ONE manifest referencing every stream's content-hash. |
| **§11.4.187** | Multi-track ruler orchestration. Session-sync bridges session state between workstations used by the multi-track conductor. |
| **§11.4.176** | Multi-track work-division. Each track's session state can be synced independently, enabling workstation switching mid-track. |
| **§11.4.131** | Standing session-resumption file. The handoff docs (`SESSION_RESUME.md`, `CONTINUATION.md`) that session-sync transfers are the §11.4.131 resumption artifacts. |
| **§11.4.127** | Session-handoff resumption-prompt mandate. The sync preserves the session handoff data needed for seamless resumption. |
| **§11.4.116** | Real-time conductor↔autonomous-test-framework sync channel. Session-sync is the cross-workstation extension of this discipline. |
| **§11.4.28** | Submodules-as-equal-codebase. The session-sync skill is inherited by reference from the constitution submodule. |

## Cross-links

- **Requires:** `continuum` (content-addressed state storage for cross-track persistence)
- **Optional:** [`multitrack`](multitrack.md) (cross-track synchronization)
- **Related skills:** [`multitrack`](multitrack.md) (parallel track orchestration), `continuum` (Merkle store backend)
- **Parent domain:** [`agent-infrastructure`](../by-domain/agent-infrastructure.md)
- **Constitution source:** [`constitution/skills/session-sync/`](../../../constitution/skills/session-sync/)

## Integration

| Surface | How it hooks in |
|---|---|
| **SSH + rsync** | The script uses SSH key-based access and rsync for efficient file transfer. Derives Claude Code project slugs from absolute paths to find the correct project directory on both hosts. |
| **7-phase sync** | Memories → sessions → session-dirs → settings+agents → remember-logs → handoff-docs → provider-config. Each phase is independently retryable. |
| **Continuum backend** | For cross-track state persistence, session-sync integrates with the continuum Merkle store (§11.4.207). Each track's state is content-addressed; unchanged tracks share a hash → O(changed) resume cost. |
| **Auto-registration** | The skill is auto-registered by `constitution/scripts/post_update_hook.sh` (§11.4.164). On every constitution pull, the hook symlinks the skill into the consuming project's `skills/` directory. |
| **Multi-track bridge** | After syncing, the local project can CONTINUE work exactly where the remote session left off — same memories, same session history, same settings. This is the bridge between multi-track workstations. |
