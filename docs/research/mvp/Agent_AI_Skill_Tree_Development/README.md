# HelixKnowledge Skill Graph System — MVP

A self-growing Knowledge Skill Graph system for AI CLI agents. Each Skill
is a versioned unit of knowledge for a specific technology, with recursive
dependencies forming a DAG. See `SPEC.md` for the technical specification
and `REQUIREMENTS.md` for the consolidated, living requirements doc.

This directory is the MVP project's documentation + planning root. The
Go implementation lives under `project/` (its own `README.md` there
covers build/run/deploy instructions for the software itself).

## Tracked-Items + Status Documents

| Document | Last modified | Revision | Markdown | HTML | PDF |
|---|---|---|---|---|---|
| REQUIREMENTS.md | (no header) | (no header) | [REQUIREMENTS.md](REQUIREMENTS.md) | (pending G43 export) | (pending G43 export) |
| CONTINUATION.md | 2026-07-15T18:24:57Z | 8 | [CONTINUATION.md](CONTINUATION.md) | (pending G43 export) | (pending G43 export) |
| GAPS_AND_RISKS_REGISTER.md | (no header) | (no header) | [GAPS_AND_RISKS_REGISTER.md](GAPS_AND_RISKS_REGISTER.md) | (pending G43 export) | (pending G43 export) |
| requests/history.md | 2026-07-15T18:24:57Z | 2 | [requests/history.md](requests/history.md) | (pending G43 export) | (pending G43 export) |
| SPEC.md | (no header) | (no header) | [SPEC.md](SPEC.md) | (pending G43 export) | (pending G43 export) |

`(no header)` means the document does not yet carry the §11.4.44
`**Revision:**` / `**Last modified:**` header block — this is reported
honestly rather than guessed; adding the header is tracked separately
(not part of this table's own scope). `(pending G43 export)` means no
HTML/PDF export pipeline exists yet for these documents — G43 tracks
building it; these are not broken links, the exports simply do not exist
yet.

## Project layout

- `SPEC.md` — technical specification.
- `REQUIREMENTS.md` — consolidated, living requirements (single source of
  truth as scope evolves).
- `IMPLEMENTATION_PLAN.md`, `plan.md` — planning documents.
- `CONTINUATION.md` — §12.10/§11.4.131 standing session-resumption file;
  read this first when starting a new session.
- `GAPS_AND_RISKS_REGISTER.md` — adversarial gap/risk audit register.
- `requests/history.md` — §11.4.208-class operator request-history ledger.
- `research/` — design research, decision records, and audits.
- `api/` — API contract materials (e.g. OpenAPI spec).
- `project/` — the Go backend implementation (own `README.md`, `Makefile`,
  `docs/`, `scripts/`, etc.).

## Notes on `project/scripts/`

`project/scripts/` currently contains two coexisting script families that
target the datastore compose stack differently:

- **Newer, `deploy/`-based family** (`start.sh`, `stop.sh`, `restart.sh`,
  `status.sh`, `logs.sh`, `install.sh`, `uninstall.sh`, `_lib.sh`): reads
  `project/deploy/docker-compose.yml` and `project/deploy/.env`.
- **Older, project-root-based family** (`backup.sh`, `migrate.sh`,
  `package.sh`, `restore.sh`): reads `project/docker-compose.yml` and
  `project/.env` directly.

Both are documented under `project/docs/scripts/`. This discrepancy is
reported here for operator awareness; reconciling the two conventions is
out of scope for this document.

## Governance

This project inherits the Helix Constitution via the repository root's
`CLAUDE.md` / `AGENTS.md` (see the `constitution/` submodule). This MVP
project directory does not maintain its own `CLAUDE.md` / `AGENTS.md` —
it inherits the repository-root files, so no project-level `QWEN.md` /
`GEMINI.md` mirror applies at this level (§11.4.157 mirroring binds at the
layer that actually owns a `CLAUDE.md`/`AGENTS.md` pair).
