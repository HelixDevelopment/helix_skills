# `scripts/restore.sh` — restore from a backup archive

**Revision:** 3
**Last modified:** 2026-07-16T00:00:00Z

## Overview

Restores the HelixKnowledge Skill Graph datastore from a backup archive
previously created by `backup.sh`: shows the archive's metadata, optionally
confirms with the operator, stops services, starts just the database,
restores the SQL dump (dropping and recreating the target database), and
(unless `--db-only`) restores configuration files and the evidence data
directory, then restarts services.

**Note on which compose stack this targets (G13):** like `backup.sh`, this
script now targets the single canonical compose file
`project/deploy/docker-compose.yml` (via an explicit
`compose -f "$INSTALL_DIR/deploy/docker-compose.yml"`), and its datastore
service is `postgres` — the retired root `project/docker-compose.yml` and its
`db` service name are gone (see `research/ops_hardening_design.md` G13 +
`scripts/check_compose_canonical.sh`). It still reads its environment from
`project/.env` (its own `INSTALL_DIR`), not `project/deploy/.env`, and uses
its own inline `load_env`/`detect_compose` helpers rather than sourcing
`_lib.sh`.

## Prerequisites

- One of: Docker with `docker compose`, `docker-compose`, or
  `podman-compose`.
- A valid backup archive produced by `backup.sh`.
- `jq` (optional — used to pretty-print backup metadata; falls back to
  `cat`-ing the raw JSON if absent).

## Usage

```
./restore.sh <backup-file> [--force] [--db-only] [--help|-h]
```

| Argument/Option | Effect |
|---|---|
| `<backup-file>` | Path to the backup archive (`.tar.gz`), required. |
| `--force` | Skip confirmation prompts (drop-and-recreate database, and the "this will overwrite existing data" warning). |
| `--db-only` | Restore only the database (skip configuration + evidence restore). |
| `--help`, `-h` | Print usage and exit 0. |

### Examples

```bash
./restore.sh data/backups/skill-system_v1.0_20240115.tar.gz
./restore.sh backup.tar.gz --force
./restore.sh backup.tar.gz --db-only
```

## Inputs

A backup archive (`.tar.gz`) as produced by `backup.sh`; `project/.env`
(optional, for DB connection defaults).

## Outputs

The target database dropped and recreated from the archive's SQL dump;
(unless `--db-only`) `project/.env` and `project/config/config.toml`
overwritten from the archive (the previous `.env` is preserved as
`.env.restore-backup.<epoch>` before being overwritten) and
`project/data/` repopulated from the archive's evidence tarball. Services
are stopped before the restore and restarted afterward. On success, prints
a reminder to verify with `curl http://localhost:8080/health`.

## Side-effects

- Drops and recreates the target Postgres database (**destructive** — this
  is the entire point of a restore, but it means any data in the target
  database that is not in the archive is lost).
- Overwrites `project/.env` and `project/config/config.toml` (in non
  `--db-only` mode), after saving a timestamped copy of the previous
  `.env`.
- Extracts the evidence tarball into `project/data/`, potentially
  overwriting existing evidence files with the same names.
- Stops then restarts the compose stack.

## Edge cases

- **Invalid/corrupt archive:** `show_backup_info()` and
  `restore_database()` both attempt `tar -xzf`; on failure the former
  logs "Invalid backup archive" and exits 1 before any destructive action
  is taken.
- **No metadata found in the archive:** logged as a warning; the restore
  proceeds anyway (metadata is informational, not required for the
  restore itself).
- **Target database already exists, no `--force`:** prompts "Drop and
  recreate? [y/N]"; any answer other than `y`/`Y` cancels the restore
  (exit 0) without touching the database.
- **No `--force`, general confirmation:** prompts "Continue with restore?
  [y/N]" before stopping services at all; declining exits 0 immediately.
- **No config / evidence found in the archive:** each restore step
  (`restore_config`, `restore_evidence`) logs a warning and continues
  rather than failing the whole restore. This is distinct from a missing
  database dump — `restore_database()` logs "Database dump not found in
  backup" and exits 1 instead (`scripts/restore.sh:110-114`), since the
  database is the one artifact this script cannot proceed without.
- **`stop.sh --compose`:** this script calls
  `"$SCRIPT_DIR/stop.sh" --compose 2>/dev/null || true` to stop services
  before the restore. `stop.sh`'s own argument parser does not recognize
  `--compose` and would exit 2 on it; because the call is wrapped in
  `|| true`, that failure is swallowed here — see `stop.md`'s "Edge cases"
  for the cross-script detail. In practice this means the "stop services"
  step of a restore may not actually stop anything if `stop.sh` rejects
  the flag; the database is still forcibly dropped/recreated regardless
  since `restore_database()` execs directly against the `postgres` service
  via compose, independent of whether `stop.sh` succeeded.
- **Postgres readiness wait:** after starting just the `postgres` service, the
  script sleeps 5s then polls `pg_isready` up to 30 times (2s apart, ~65s
  total) before proceeding to the restore, regardless of whether readiness
  was actually reached (no hard failure if the loop exhausts retries
  without success — the subsequent `psql` restore command will simply fail
  and surface its own error).

## Internal behavior

1. `load_env()` / `detect_compose()` — same as `backup.sh`.
2. `main()` parses `<backup-file>` (first non-flag argument) plus
   `--force`/`--db-only`/`--help`.
3. `show_backup_info()` — extracts and prints `backup.json` metadata
   (via `jq` if available).
4. Confirmation prompt (skipped with `--force`).
5. Stops services (`stop.sh --compose`, best-effort), starts just
   `postgres`, waits for Postgres readiness (bounded retry loop).
6. `restore_database()` — extracts + decompresses the SQL dump, drops and
   recreates the target database, replays the dump.
7. Unless `--db-only`: `restore_config()` (backs up then overwrites `.env`
   / `config.toml`) and `restore_evidence()` (extracts the evidence
   tarball into `project/data/`).
8. Restarts services via `start.sh`; prints a health-check reminder.

## Dependencies

`bash`, one of `{docker compose, docker-compose, podman-compose}`, `tar`,
`gunzip`, `jq` (optional), `stop.sh`, `start.sh`.

## Cross-references

`backup.sh` (produces the archives this script consumes), `stop.sh`,
`start.sh` (both invoked mid-script).

## Last verified

2026-07-16, against `project/scripts/restore.sh` (10245 bytes, last
modified 2026-07-16) after the G13 canonical-compose change (compose calls
target `-f deploy/docker-compose.yml`; datastore service `postgres`).
