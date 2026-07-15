# `scripts/backup.sh` — create a backup archive

**Revision:** 1
**Last modified:** 2026-07-16T00:00:00Z

## Overview

Creates a compressed `tar.gz` archive containing a database dump
(`pg_dump`), configuration files (`.env`, `config/config.toml`,
`docker-compose.yml`), the evidence data directory, and a JSON metadata
file describing the backup. Companion of `restore.sh`, which consumes the
archives this script produces.

**Note on which compose stack this targets:** unlike the newer
`start.sh`/`stop.sh`/`_lib.sh` family, `backup.sh` is part of the older
script family that reads its own `INSTALL_DIR` (= `project/`, i.e. the
parent of `scripts/`) and expects `docker-compose.yml` + `.env` directly
at `project/` root — **not** `project/deploy/docker-compose.yml` /
`project/deploy/.env`. See "Two coexisting script families" in this
project's top-level `README.md` for the observed discrepancy between the
two conventions.

## Prerequisites

- One of: Docker with `docker compose`, `docker-compose` (legacy
  standalone binary), or `podman-compose`.
- A running `db` compose service (the script execs `pg_dump` inside it via
  `compose exec -T db`).
- `project/.env` (optional; falls back to `DB_HOST=db`,
  `DB_PORT=5432`, `DB_NAME=skilldb`, `DB_USER=skilluser`,
  `DB_PASSWORD=skillpassword` if absent or a given variable is unset).

## Usage

```
./backup.sh [--output <dir>] [--name <name>] [--full | --db-only] [--retention <days>] [--list] [--help|-h]
```

| Option | Effect |
|---|---|
| `--output <dir>` | Backup directory (default: `<project>/data/backups`). |
| `--name <name>` | Custom backup name (default: `skill-system_<version>_<timestamp>`). |
| `--full` | Full backup: database + configuration + evidence (default). |
| `--db-only` | Database only. |
| `--retention <days>` | Retention period in days for automatic cleanup of old backups (default: 30; 0 disables cleanup). |
| `--list` | List existing backups and exit (does not create a new backup). |
| `--help`, `-h` | Print usage and exit 0. |

### Examples

```bash
./backup.sh                           # full backup, default retention
./backup.sh --db-only                 # database-only backup
./backup.sh --name pre-upgrade        # named backup before a risky change
./backup.sh --retention 7             # keep only 7 days of backups
./backup.sh --list                    # show existing backups, then exit
```

## Inputs

`project/.env` (optional, for DB connection defaults), a running `db`
compose service, `project/config/config.toml` and `project/docker-compose.yml`
(copied when present, `--full` mode only), `project/data/evidence/`
(archived when present, `--full` mode only).

## Outputs

A single `<name>.tar.gz` archive under the backup directory, containing:
`database/skilldb.sql.gz`, and (in `--full` mode) `config/.env`,
`config/config.toml`, `config/docker-compose.yml`, and
`evidence/evidence.tar.gz`, plus a `metadata/backup.json` describing the
backup's name/type/timestamp/version/hostname/DB connection info/contents
flags. The final archive path is printed (and echoed to stdout).

## Side-effects

Creates a temp directory (cleaned up on both success and the `pg_dump`
failure path), writes the final archive to `$BACKUP_DIR`, and — if
`--retention` is non-zero — deletes archives in that directory older than
the retention window on every run.

## Edge cases

- **No compose command found:** the script's own `detect_compose()`
  checks `docker compose`, then `docker-compose`, then `podman-compose`;
  none found → `log_error` + `exit 1`.
- **`pg_dump` fails on the first attempt:** the script retries once with
  `PGPASSWORD` explicitly injected via `compose exec -T -e
  PGPASSWORD=... db pg_dump ...` before giving up (`exit 1`, temp dir
  removed).
- **No evidence directory:** logged as a warning; the backup proceeds
  without an `evidence.tar.gz` member.
- **`--retention 0`:** disables the automatic cleanup of old backups
  entirely (the `if [ "$RETENTION_DAYS" -gt 0 ]` guard in
  `cleanup_old_backups`).
- **`--list` with no existing backups:** prints "No backups found" and
  returns (does not error).

## Internal behavior

1. `load_env()` sources `project/.env` if present, with fallback defaults.
2. `detect_compose()` picks the first available compose front-end.
3. `main()` parses arguments; `--list` short-circuits to `list_backups()`
   and exits.
4. `generate_name()` builds the backup name from `--name`, or
   `skill-system_<version>_<timestamp>` where `<version>` is read by
   execing the `api` container's `--version` flag (falls back to
   `"unknown"` if that fails).
5. `create_backup()`: creates a temp working tree
   (`database/`, `config/`, `evidence/`, `metadata/`), dumps + gzips the
   database, copies configuration + archives evidence in `--full` mode,
   writes `metadata/backup.json`, tars the whole tree into the final
   archive under `$BACKUP_DIR`, prints the resulting size, then calls
   `cleanup_old_backups()`.

## Dependencies

`bash`, one of `{docker compose, docker-compose, podman-compose}`,
`gzip`, `tar`, `mktemp`, `date`, `du`.

## Cross-references

`restore.sh` (consumes the archives this script produces),
`package.sh` (a separate, unrelated "package the whole project for
distribution" script — not a backup).

## Last verified

2026-07-16, against `project/scripts/backup.sh` (9176 bytes, last modified
2026-07-15).
