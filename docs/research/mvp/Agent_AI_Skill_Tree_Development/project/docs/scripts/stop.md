# `scripts/stop.sh` — bring the datastore stack down

**Revision:** 1
**Last modified:** 2026-07-16T00:00:00Z

## Overview

Runs `compose down` against `deploy/docker-compose.yml`. This is the
`ExecStop` command of the `systemctl --user` unit
(`deploy/systemd/helix-skills.service`) as well as the manual entry point.

## Prerequisites

- One of: Docker with `docker compose`, Podman with `podman compose`, or
  `podman-compose`.
- `deploy/docker-compose.yml` present (required).
- `deploy/.env` (optional).

## Usage

```
scripts/stop.sh [--quiet] [-h|--help]
```

| Option | Effect |
|---|---|
| `--quiet`, `-q` | Suppress informational (non-error) output. |
| `-h`, `--help` | Print usage and exit 0. |

### Examples

```bash
scripts/stop.sh            # stop the stack
scripts/stop.sh --quiet    # used by the systemd unit's ExecStop
```

## Inputs

`deploy/docker-compose.yml` (required), `deploy/.env` (optional).

## Outputs

Containers and the compose network are stopped and removed. **Named
volumes are preserved** — the Postgres data directory survives a
stop/start cycle.

## Side-effects

Stops/removes containers + the compose network via the container engine.
Never removes named volumes.

## Edge cases

- **No compose engine found / compose file missing:** dies via `_lib.sh`
  (`hx_detect_engine` / `hx_require_compose_file`) before attempting
  anything.
- **Stack already down:** `compose down` on an already-down stack is a
  no-op / idempotent.
- **Unrecognized flag observed elsewhere in the codebase:** `restore.sh`
  invokes this script as `"$SCRIPT_DIR/stop.sh" --compose 2>/dev/null ||
  true`. `stop.sh`'s argument parser recognizes only `--quiet`/`-q` and
  `-h`/`--help` — any other argument (including `--compose`) falls through
  to the `*` case, which prints `stop.sh: unknown argument: --compose` to
  stderr and exits 2. Because `restore.sh` calls it with `|| true`, that
  non-zero exit is swallowed there rather than surfaced — call sites that
  need the stack actually stopped should not rely on that particular
  invocation.

## Internal behavior

1. Parses `--quiet` / `--help`.
2. `hx_load_env` → `hx_require_compose_file` → `hx_detect_engine`.
3. `hx_compose down`.

## Dependencies

`_lib.sh`, one of `{docker, podman, podman-compose}`.

## Cross-references

`start.sh`, `restart.sh` (composes stop+start), `status.sh`,
`deploy/docker-compose.yml`, `deploy/systemd/helix-skills.service` (its
`ExecStop=`). Also invoked (with an unsupported `--compose` flag, ignored
via `|| true`) from the older `restore.sh`.

## Last verified

2026-07-16, against `project/scripts/stop.sh` (2383 bytes, last modified
2026-07-15).
