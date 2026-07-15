# Docs Chain Incorporation Plan — HelixKnowledge Skill Graph System

**Revision:** 1
**Last modified:** 2026-07-15T00:00:00Z
**Status:** research/plan — engine FOUND on host (sibling checkout), NOT yet
incorporated into `helix_skills`; contexts DESIGNED, not yet registered.
**Authority:** operator mandate R10 (REQUIREMENTS.md:80) + Constitution §11.4.106
(Docs Chain — mechanical documentation/DB sync engine).
**Scope:** verify the real engine + its real capabilities against actual files,
map this project's derived-doc surface to Docs Chain contexts, reconcile with
§11.4.106 status honesty + §11.4.28C dependency layout, and flag every gap.
**Anti-bluff boundary (§11.4.6):** every capability below is cited to a real file
+ line read on this host; anything not verified is marked UNCONFIRMED / UNKNOWN.

---

## 0. Honest verdict (read this first)

| Question | Answer |
|---|---|
| Is the Docs Chain engine present on this host? | **YES** — but in a *sibling* project, not in `helix_skills`. |
| Real engine path | `/home/milos/Factory/projects/tools_and_research/helix_code/submodules/docs_chain` |
| Git URL | `git@github.com:vasic-digital/docs_chain.git` (+ GitLab mirror `git@gitlab.com:vasic-digital/docs_chain.git`) |
| Is it incorporated into `helix_skills`? | **NO** — no `docs_chain/`, no `submodules/`, no `.docs_chain/`; only the `constitution` submodule is present. **Vendor/incorporation step required.** |
| Can the shipped binary run here? | **NO** — the committed `docs_chain` binary is wrong-arch (`Exec format error`); must `go build` from source on this host. |
| Are the export tools present? | **YES** — `pandoc` and `weasyprint` both on PATH; `go` 1.26.4 present. |
| Chain contexts mapped for this project | **3** (2 buildable now, 1 blocked on a not-yet-created DB) + 1 optional. |

**Bottom line:** the engine is real, its capabilities are as §11.4.106 describes
(verified from source), and this project's doc surface maps cleanly onto it. The
only things between "designed" and "working" are (a) incorporating the engine
into `helix_skills` (operator-gated, §11.4.28C), (b) building the binary on this
host, (c) creating the `.docs_chain/contexts/*.yaml` this doc specifies. The
workable-items DB chain (§11.4.93/.95) is additionally blocked on a DB that does
not exist yet.

---

## 1. Engine location — evidence (Task 1)

Searched, in order, the operator-named candidates:

| Candidate | Result |
|---|---|
| `/home/milos/Projects/docs_chain` | **ABSENT** — `ls: cannot access … No such file or directory`. |
| sibling under `/home/milos/Factory/projects/tools_and_research/` (`docs_chain/`) | **ABSENT** at that level. |
| `helix_code/submodules/docs_chain/` | **FOUND** ✅ |
| `.gitmodules` / `helix-deps.yaml` reference in `helix_skills` | `helix_skills/.gitmodules` declares ONLY `constitution`; **no root `helix-deps.yaml`** in `helix_skills`. The project's own manifest (`.../Agent_AI_Skill_Tree_Development/project/helix-deps.yaml`) DECLARES `docs_chain` as a dep (see §4), but it is not vendored. |

**Real engine path (verified with `ls` + `git remote -v`):**

```
/home/milos/Factory/projects/tools_and_research/helix_code/submodules/docs_chain
  git remote origin  → git@github.com:vasic-digital/docs_chain.git (fetch)
                       git@gitlab.com:vasic-digital/docs_chain.git (push, mirror)
  git HEAD           → 9313c62 "docs: remove dangling cross-project … reference"
```

Top-level layout (real): `cmd/docs_chain/` (CLI), `internal/{adapter,config,graph,hash,orchestrator,runner,state}` (engine), `docs/` (ARCHITECTURE, CONFIG_SCHEMA, CONSTITUTION_INTEGRATION, USE_CASE_CATALOGUE, USER_GUIDE — each md+html+pdf+docx), `.docs_chain/contexts/self-docs.yaml` (dogfood context), `go.mod` (module `digital.vasic.docs_chain`), `helix-deps.yaml` (leaf, `deps: []`), and a committed `docs_chain` binary (**wrong-arch on this host**).

**Honest note on reachability:** the engine lives inside the `helix_code` project's
`submodules/` tree. From `helix_skills`' perspective it is a *foreign checkout*,
not a dependency `helix_skills` resolves. For `helix_skills` the engine is
**present-on-host but not-incorporated** — a vendor step is still owed.

---

## 2. Verified capabilities (Task 2 — from actual source, not from memory)

All of the following were read from `cmd/docs_chain/main.go`,
`docs/CONFIG_SCHEMA.md`, `internal/config/config.go`, `internal/config/validate.go`,
`internal/hash/hash.go`, and the dogfood context `.docs_chain/contexts/self-docs.yaml`.

### 2.1 Subcommands (verified from `main.go` `run()` switch)

| Subcommand | Purpose (from source) |
|---|---|
| `doctor [--all \| <ctx>] [--root DIR]` | validate contexts (parse + graph + per-transform tool availability). Never writes. Tool-absent is **WARN**, not failure. |
| `sync [--all \| <ctx>]` | propagate atomically, update `state.json`, write per-run evidence. |
| `rebaseline [--all \| <ctx>]` | re-baseline sync edges from the AUTHORITY side only (authority side never written). |
| `verify [--all \| <ctx>]` | read-only drift check — the CI / pre-build gate. Non-zero (1) when any node is stale. |
| `graph <ctx>` | print topo order + edges (debug). |
| `watch [--all \| <ctx>] [--debounce 300ms]` | fsnotify daemon: sync on source change. |
| `version` / `help` | banner / usage. |

**Discrepancy flagged (§11.4.6):** the §11.4.106 anchor text says the engine
"registers `sync`/`verify`/`diff`(alias)/`doctor`". The **actual binary at this
HEAD has NO `diff` subcommand** — it exposes `doctor/sync/rebaseline/verify/graph/
watch`. The real CLI is a *superset* of the anchor's list minus the `diff` alias.
Not a defect for our purposes (`sync`/`verify`/`doctor` — the ones we need — all
exist), but the anchor's `diff` alias is not honored by this build. Do not script
against `docs_chain diff`.

### 2.2 Exit-code contract (verified from `main.go` constants + `usage()`)

`0` in-sync/applied/healthy · `1` generic error · `2` sync conflict (both sides
dirty — **never silent-merge**) · `3` transform failed (run rolled back, no live
changes) · `4` cycle / config-validation error. `verify` exits `1` on stale.

### 2.3 Transforms / builtins (verified from `CONFIG_SCHEMA.md §5.1` + `config.go` constants + `main.go builtinTool()`)

| Builtin | Maps | External tool needed | Verified |
|---|---|---|---|
| `pandoc-html` | markdown → html | **pandoc** | `config.go:87`, `builtinTool()` |
| `pandoc-docx` | markdown → docx | **pandoc** | `config.go:89`, self-docs.yaml `md2docx` |
| `weasyprint-pdf` | html → pdf (or md → pdf) | **weasyprint** | `config.go:88`, `builtinTool()` |
| `colorize-html` | html → html post-process (§11.4.23) | internal (x/net/html) | CONFIG_SCHEMA §5.1 |
| `gen-summary` | markdown → summary | configured generator | CONFIG_SCHEMA §5.1 |
| `md-to-sqlite` / `sqlite-to-md` | markdown ↔ sqlite | **internal, pure-Go** `modernc.org/sqlite` (no external `sqlite3`) | CONFIG_SCHEMA §5.1 |
| `members-fingerprint` | members-glob → fingerprint (§11.4.86) | internal sha256-of-sorted-members | CONFIG_SCHEMA §5.1 |
| `exec: <cmd> args: [...]` | any custom transform | the named script/binary | CONFIG_SCHEMA §5.2 |

So it **does** need `pandoc` on PATH (html + docx) and `weasyprint` on PATH (pdf).
It does **not** need an external `sqlite3` — the DB adapters are pure-Go. An
absent tool surfaces as a typed tool-absent rollback → honest SKIP-with-reason
(`main.go` `IsToolAbsent` / `formatSyncResult` "SKIP (tool absent)"), never a fake
PASS or partial write.

### 2.4 Node kinds (verified `CONFIG_SCHEMA.md §3.1` + `validate.go`)

`markdown`, `html`, `pdf`, `docx` (`validate.go:30 graph.KindDOCX` — present in
validator + used by the real dogfood context; note the `docx` row is *missing*
from the CONFIG_SCHEMA §3.1 table — minor doc drift, but the kind is genuinely
accepted), `sqlite` (bidirectional), `summary`, `status` (§11.4.45),
`status_summary` (§11.4.56), `fingerprint` (§11.4.86, requires `members:` glob).

### 2.5 The §11.4.106 core guarantees (verified)

- **Content-hash change detection, NOT mtime (§11.4.86):** `internal/hash/hash.go`
  — "content hash over NORMALIZED content — never by mtime"; `sha256.Sum256`;
  `FingerprintMembers` = sha256 over the SORTED member list. Confirmed.
- **Atomic commit + rollback + conflict:** `main.go` handles
  `StatusCommitted/InSync/Conflict/Cycle/RolledBack`; conflict = exit 2 "no
  writes"; rollback = exit 3 "no live changes". Confirmed at the CLI layer.
  (The temp-stage → atomic-rename mechanics live in `internal/orchestrator` /
  `internal/runner`, referenced by `main.go` but not line-verified here —
  **UNCONFIRMED at byte level**, asserted by ARCHITECTURE.md §5/§8 and the
  exec-transform contract in CONFIG_SCHEMA §5.2 "Docs Chain owns the atomic
  rename".)
- **Both-dirty → conflict, never silent-merge (§11.4.6):** exit 2 path confirmed.
- **Per-run captured evidence (§11.4.69):** `main.go writeEvidence()` →
  `qa-results/docs_chain/<run-id>/<kind>.log`. Confirmed.
- **Consumer-owned data (§11.4.28):** contexts resolve from
  `<root>/.docs_chain/contexts/*.yaml`; state at `<root>/.docs_chain/state.json`.
  Confirmed (`main.go contextsRel`, `state.DefaultPath`).

### 2.6 Runtime evidence attempt (honest)

Running the committed binary FAILED: `./docs_chain version` →
`cannot execute binary file: Exec format error`. The shipped artifact is built
for a different architecture than this host. **Therefore no live `doctor`/`sync`
run could be captured on this host from the shipped binary.** Capability claims
above are from **source inspection**, not a live run. To get a live run, the
engine must be built here (`go build ./cmd/docs_chain` — `go` 1.26.4 is present;
module deps `fsnotify`, `yaml.v3`, `modernc.org/sqlite` are declared in `go.mod`/
`go.sum` and require the Go module cache/network to resolve — **UNCONFIRMED
whether they resolve offline on this host**).

---

## 3. This project's derived-doc surface → Docs Chain contexts (Task 3)

### 3.1 The surface (verified with `find`)

Under `docs/research/mvp/Agent_AI_Skill_Tree_Development/`:

- **Top-level planning docs (5 `.md`):** `REQUIREMENTS.md`, `IMPLEMENTATION_PLAN.md`,
  `GAPS_AND_RISKS_REGISTER.md`, `SPEC.md`, `plan.md`.
- **`research/*.md` (13 files):** the dimension research + incorporation notes
  (this file included).
- **`project/` docs (6 `.md`):** `project/README.md`,
  `project/docs/{ARCHITECTURE,API,INSTALL,MCP_INTEGRATION,DANGER_ZONES}.md`.
- **`api/openapi.yaml`:** an OpenAPI spec — **NOT a Docs Chain node kind**
  (no `openapi`/`yaml` kind). Could be rendered to HTML via an `exec:` transform
  (redoc/swagger) but that is out of scope here → flagged UNKNOWN/optional.
- **Existing exports:** `find` for `*.html`/`*.pdf`/`*.docx` returned **NOTHING** —
  no derived exports exist yet; the chain will *generate* them on first `sync`.
- **§11.4.44 headers:** the planning docs do NOT yet carry a `Revision:` +
  `Last modified:` header (e.g. `REQUIREMENTS.md` opens with a `Status:` line, no
  revision header). Docs Chain does **not** add these (author owns them, §11.4.106
  clause D) — a separate authoring gap, not a chain blocker.

### 3.2 Mapped contexts

| Context | Sources | Derived | Transforms | Buildable now? |
|---|---|---|---|---|
| `mvp_planning_docs` | 5 top-level `.md` | `.html` + `.pdf` + `.docx` per doc | `pandoc-html`, `weasyprint-pdf`, `pandoc-docx` | **YES** (pandoc+weasyprint present; needs built engine) |
| `research_docs` | 13 `research/*.md` | `.html` + `.pdf` per doc | `pandoc-html`, `weasyprint-pdf` | **YES** (same) |
| `workable_items` | future `workable_items.db` ↔ `Issues.md` | `Issues_Summary.md`, `.html`, `.pdf` | `md-to-db`/`db-to-md` (exec), `gen-summary` (exec), `pandoc-html`, `weasyprint-pdf` | **NO — blocked** (DB + generators do not exist yet) |
| `project_docs` (optional) | `project/README.md` + `project/docs/*.md` | `.html` + `.pdf` | `pandoc-html`, `weasyprint-pdf` | YES, once the service repo root is fixed |

### 3.3 Concrete example context — grounded in the REAL schema

`.docs_chain/contexts/mvp_planning_docs.yaml` (paths shown relative to whatever
becomes the consuming repo root; `--root` selects it — see §4.2):

```yaml
context: mvp_planning_docs
description: HelixKnowledge Skill Graph MVP planning docs — md -> html/pdf/docx (§11.4.65)
nodes:
  req_md:    { kind: markdown, path: REQUIREMENTS.md }
  req_html:  { kind: html,     path: REQUIREMENTS.html }
  req_pdf:   { kind: pdf,      path: REQUIREMENTS.pdf }
  req_docx:  { kind: docx,     path: REQUIREMENTS.docx }
  plan_md:   { kind: markdown, path: IMPLEMENTATION_PLAN.md }
  plan_html: { kind: html,     path: IMPLEMENTATION_PLAN.html }
  plan_pdf:  { kind: pdf,      path: IMPLEMENTATION_PLAN.pdf }
  gaps_md:   { kind: markdown, path: GAPS_AND_RISKS_REGISTER.md }
  gaps_html: { kind: html,     path: GAPS_AND_RISKS_REGISTER.html }
  gaps_pdf:  { kind: pdf,      path: GAPS_AND_RISKS_REGISTER.pdf }
  spec_md:   { kind: markdown, path: SPEC.md }
  spec_html: { kind: html,     path: SPEC.html }
  spec_pdf:  { kind: pdf,      path: SPEC.pdf }
edges:
  - { type: derive-from, from: req_md,   to: req_html,  transform: md2html }
  - { type: derive-from, from: req_html, to: req_pdf,   transform: html2pdf }
  - { type: derive-from, from: req_md,   to: req_docx,  transform: md2docx }
  - { type: derive-from, from: plan_md,  to: plan_html, transform: md2html }
  - { type: derive-from, from: plan_html,to: plan_pdf,  transform: html2pdf }
  - { type: derive-from, from: gaps_md,  to: gaps_html, transform: md2html }
  - { type: derive-from, from: gaps_html,to: gaps_pdf,  transform: html2pdf }
  - { type: derive-from, from: spec_md,  to: spec_html, transform: md2html }
  - { type: derive-from, from: spec_html,to: spec_pdf,  transform: html2pdf }
transforms:
  md2html:  { builtin: pandoc-html }
  html2pdf: { builtin: weasyprint-pdf }
  md2docx:  { builtin: pandoc-docx }
```

This is a byte-faithful application of the verified schema (kinds, edge shape,
builtin names) — identical in form to the engine's own dogfood
`.docs_chain/contexts/self-docs.yaml`. `pdf` derives from `html` (not from `md`)
so any HTML styling/colorization carries through — the same pattern the engine
uses on itself.

### 3.4 The future DB-sync context (DESIGNED, flagged pending)

`.docs_chain/contexts/workable_items.yaml` — **DO NOT register until the DB +
generators exist** (§11.4.93/.95). Schema is real; the referenced `exec:`
binaries are **not present in this project** → UNCONFIRMED-pending-DB:

```yaml
context: workable_items
description: §11.4.93 workable-items DB single-source-of-truth <-> Issues.md
nodes:
  items_db:       { kind: sqlite,   path: docs/workable_items.db }   # SSoT, tracked (§11.4.95)
  issues_md:      { kind: markdown, path: docs/Issues.md }
  issues_summary: { kind: summary,  path: docs/Issues_Summary.md }
  issues_html:    { kind: html,     path: docs/Issues.html }
  issues_pdf:     { kind: pdf,      path: docs/Issues.pdf }
edges:
  - { type: sync, a: issues_md, b: items_db, authority: items_db,
      transform_a_to_b: md-to-db, transform_b_to_a: db-to-md }
  - { type: derive-from, from: issues_md,   to: issues_summary, transform: gen-issues-summary }
  - { type: derive-from, from: issues_md,   to: issues_html,    transform: md2html }
  - { type: derive-from, from: issues_html, to: issues_pdf,     transform: html2pdf }
transforms:
  md-to-db:           { exec: "scripts/testing/workable_items", args: ["sync","md-to-db"] }  # MISSING binary
  db-to-md:           { exec: "scripts/testing/workable_items", args: ["sync","db-to-md"] }  # MISSING binary
  gen-issues-summary: { exec: "scripts/testing/generate_issues_summary.sh" }                  # MISSING script
  md2html:            { builtin: pandoc-html }
  html2pdf:           { builtin: weasyprint-pdf }
```

**Blockers:** `docs/workable_items.db` does not exist; the `workable_items`
binary + `generate_issues_summary.sh` do not exist in this project. Until they
do, `docs_chain doctor workable_items` would WARN on the missing `exec:` targets
and `sync` would tool-absent-SKIP the `exec` edges. Do not register this context
as "working" — it is pending the DB deliverable.

---

## 4. §11.4.106 status honesty + §11.4.28C layout + tooling (Task 4)

### 4.1 §11.4.106 status honesty

The engine's own `docs/CONSTITUTION_INTEGRATION.md` (Revision 3) states — and this
matches the constitution's §11.4.106 status line — that **Phase-4 CLI + YAML
loader are IMPLEMENTED** (a first consumer, Herald, wired a verify-green 66-doc
md→html/pdf/docx corpus), while **Phase-6 constitution-submodule distribution is
PLANNED + OPERATOR-GATED** (creating the remote, adding the submodule pointer,
editing governance — all §11.4.66/.35/.26 operator actions, NOT agent actions).
So for `helix_skills`: consume **by reference** to a Projects-root sibling
checkout today (as Herald does); the turnkey `git submodule` distribution is not
yet available and its creation is operator-gated. Claim no unshipped behaviour.

The §11.4.106 anchor→mechanism map (verified in CONSTITUTION_INTEGRATION.md §4)
covers exactly the anchors this project's chains touch: §11.4.12/.53/.45/.56/.57/
.59/.60/.65/.86/.93/.95/.44 + §9.2/§11.4.6/.50/.69/§12.10. Our `mvp_planning_docs`
+ `research_docs` contexts mechanize **§11.4.65** (universal md export); the
future `workable_items` context mechanizes **§11.4.12/.93/.95**.

### 4.2 §11.4.28C dependency layout

§11.4.28C: a dependency resolves to ONE canonical location —
`<repo_root>/<name>/` (ungrouped) or `<repo_root>/submodules/<name>/` (grouped) —
and nested own-org chains are FORBIDDEN. Applied here:

- The project's own `project/helix-deps.yaml` already DECLARES `docs_chain`
  (`git@github.com:vasic-digital/docs_chain.git`, `ref: main`, `layout: grouped`)
  → destination `submodules/docs_chain/` relative to the future service repo root.
- **Parent-copy priority (R9, REQUIREMENTS.md:76-79):** where a parent-dir copy of
  a dep already exists, the parent version has PRIORITY and both copies must stay
  in sync. On this host a parent-level `docs_chain` copy exists inside the SIBLING
  `helix_code/submodules/docs_chain` — but that is a *different project's* tree,
  not `helix_skills`' parent. The GAPS register (X1) records the parent-vs-second-
  copy sync decision as **OPEN** ("no second-copy `--apply` until resolved").
  **Do not vendor a second copy blindly** — resolve X1 first (§11.4.6).
- `docs_chain`'s own `helix-deps.yaml` is a **leaf** (`deps: []`) — incorporating
  it adds zero transitive own-org submodules, so §11.4.28C's no-nested-chain rule
  is satisfied.
- Org note (operator framing for THIS reconciliation): treat the engine as a
  vendored-by-reference dependency; `docs_chain` is under `vasic-digital`. Whether
  it is classified "owned" (§11.4.28 org list) or "third-party" for vendoring
  purposes, the incorporation is **by reference, never copied into logic**
  (§11.4.106 clause A / §11.4.177) — consistent either way.

### 4.3 Tooling presence on THIS host (verified with `command -v`)

| Tool | Needed for | Present? | Path |
|---|---|---|---|
| `pandoc` | `pandoc-html`, `pandoc-docx` | **YES** | `/home/milos/Factory/software/pandoc/bin/pandoc` |
| `weasyprint` | `weasyprint-pdf` | **YES** | `/home/milos/Factory/software/weasyprint/bin/weasyprint` |
| `go` 1.26.4 | build the engine from source | **YES** | `/usr/bin/go` |
| external `sqlite3` | NOT needed (adapters are pure-Go `modernc.org/sqlite`) | (present, android variant) | n/a for chains |

So the two derive contexts (`mvp_planning_docs`, `research_docs`) are
**tooling-ready today** — the only missing piece is a host-built `docs_chain`
binary.

---

## 5. Honest gaps + tracked-item wiring (Task 5, §11.4.6 + §11.4.197)

**Gaps (nothing here is buildable-and-wired until closed):**

1. **Engine not incorporated into `helix_skills`.** No `docs_chain/`,
   `submodules/`, or `.docs_chain/` in the repo; only `constitution` is a
   submodule; `constitution/submodules/` holds `clickup_sync`, `continuum`,
   `session_orchestrator`, `token_optimizer` — **not** `docs_chain`. → operator-
   gated vendor step (§11.4.28C / §11.4.66), pending X1 (parent-copy priority).
2. **Shipped binary wrong-arch** (`Exec format error`) → must `go build
   ./cmd/docs_chain` on this host; **UNCONFIRMED** whether module deps resolve
   offline. No live `sync`/`verify` could be captured here.
3. **No exports exist yet** — first `sync` generates all `.html`/`.pdf`/`.docx`.
4. **`workable_items` chain blocked** — DB (`docs/workable_items.db`), the
   `workable_items` md↔db binary, and `generate_issues_summary.sh` do not exist in
   this project (§11.4.93/.95 not yet delivered).
5. **No `.docs_chain/` scaffold** — `.docs_chain/contexts/` dir, `state.json` +
   `*.docs_chain.tmp` gitignore entries not yet created.
6. **Consuming repo root undecided** — the deliverable sits under
   `docs/research/mvp/…`; the future service's root (where `.docs_chain/` lives +
   what `--root` points at) is UNKNOWN and is an operator/§11.4.35 declaration.
7. **§11.4.44 headers** absent on planning docs (author-owned; not a chain job).
8. **`diff` subcommand** referenced by the §11.4.106 anchor does NOT exist in this
   engine build — script against `sync`/`verify`/`doctor`, never `diff`.
9. **api/openapi.yaml** has no native Docs Chain node kind — optional `exec:`
   render only; out of scope.

**Tracked-item wiring (so this does not become un-wired research — §11.4.197):**
this plan is the design artefact for the ALREADY-TRACKED items —
**REQUIREMENTS.md R10** ("Fully incorporate the Docs Chain submodule (`.docs_chain`)")
and **IMPLEMENTATION_PLAN.md P13.T1** ("Docs Chain submodule fully incorporated
… Evidence gate: submodule resolves; docs build passes through the chain"), with
the parent-copy-priority decision tracked as **GAPS_AND_RISKS_REGISTER.md X1**.
Recommended P13.T1 completion evidence gate (all four): (a) engine resolvable from
the consuming repo root per §11.4.28C; (b) host-built binary runs `doctor --all`
clean; (c) `sync --all` generates the exports with evidence under
`qa-results/docs_chain/<run-id>/`; (d) `verify --all` exits 0. This plan is
**necessary-not-sufficient**: it is DESIGN, and R10/P13.T1 remain OPEN until (a)–(d)
are captured (§11.4.6 — designed ≠ working).

---

## Sources verified (local, 2026-07-15)

- Engine: `/home/milos/Factory/projects/tools_and_research/helix_code/submodules/docs_chain` @ HEAD `9313c62`; `git remote -v` → `git@github.com:vasic-digital/docs_chain.git` (+ GitLab mirror).
- Source read: `cmd/docs_chain/main.go`; `docs/CONFIG_SCHEMA.md`; `docs/CONSTITUTION_INTEGRATION.md` (Rev 3); `internal/config/config.go` (builtin constants); `internal/config/validate.go` (`KindDOCX`); `internal/hash/hash.go` (sha256, not-mtime); `.docs_chain/contexts/self-docs.yaml`; `helix-deps.yaml` (leaf, `deps: []`).
- Host tooling: `command -v pandoc weasyprint go sqlite3`; shipped binary `Exec format error`.
- Project surface: `find` under `docs/research/mvp/Agent_AI_Skill_Tree_Development/` (5 top-level md, 13 research md, 6 project md, `api/openapi.yaml`, zero existing exports); `project/helix-deps.yaml` (docs_chain dep declared, grouped); `REQUIREMENTS.md:80` (R10); `IMPLEMENTATION_PLAN.md:308` (P13.T1); `GAPS_AND_RISKS_REGISTER.md:178,224` (X1 + docs-lint).
- `helix_skills/.gitmodules` (only `constitution`); `constitution/submodules/` (no `docs_chain`).
