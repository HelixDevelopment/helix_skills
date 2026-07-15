# G43 — Docs Chain Context-YAML + Export-Gate Wiring Design (concrete)

**Revision:** 1
**Last modified:** 2026-07-15T18:44:42Z
**Scope:** DESIGN ONLY. No `.docs_chain/` file created, no context YAML written to
disk, no gate script created, no export generated, no submodule added, no git
operation, nothing under `project/` touched by this document's author. The only
artefact this task produced is this one markdown file. The fenced YAML/shell
blocks below are the *proposed content* for the wiring step to author later, not
files that exist today.
**Builds on:** `research/g43_docs_chain_export_wiring_design.md` (Rev 1) — that
doc settled the STRATEGY (adopt the real `vasic-digital/docs_chain` engine; it
builds clean, `go test ./...` green; pandoc 3.10 + weasyprint 69.0 present;
recommend vendoring). This doc does NOT re-derive that; it designs the CONCRETE
per-project wiring: the actual `.docs_chain/contexts/*.yaml` content and the
export gate, grounded field-for-field in the engine's REAL config schema.
**Authority:** Constitution §11.4.106 (Docs Chain mechanical enforcer) / §11.4.65
(universal Markdown export) / §11.4.12 + §11.4.53 (auto-generated docs sync) /
§11.4.86 (fingerprint/content-hash not mtime) / §11.4.44 (revision header) /
§11.4.28(B) + §11.4.177 (decoupling — engine by reference, project literals in
consumer-owned data) / §11.4.50 (deterministic gate) / §11.4.69 (captured
evidence) / §11.4.6 (no-guessing).
**Closes (design half of):** `GAPS_AND_RISKS_REGISTER.md` **G43** (HIGH,
§11.4.106/§11.4.65) — "zero `.html`/`.pdf` exports exist for any tracked doc".

---

## 1. Goal — build on the g43 strategy, not restate it

`g43_docs_chain_export_wiring_design.md` established the WHAT and the WHETHER:
adopt the real Docs Chain engine, vendor it, wire the project's own tracked-doc
export pipeline through it. It deliberately stopped short of the concrete
per-doc context files ("that would be mechanical … exactly the kind of thing the
actual wiring step should generate" — g43 §2 "Honest gap"). This document
supplies exactly that missing layer:

1. the engine's REAL context-YAML schema, cited to the shipped source (§2);
2. the CURRENT MVP tracked-doc inventory with real counts (§3);
3. the concrete `.docs_chain/contexts/*.yaml` content, field-valid against that
   schema (§4);
4. the per-doc `md → html → pdf` chain semantics (§5);
5. the `CM-DOC-EXPORTS-SYNC` gate + its paired §1.1 mutation (§6);
6. the three §11.4.106(F) wiring seams + honest boundary (§7);
7. the vendoring dependency on the G14 decision (§8);
8. the anti-bluff captured-evidence shape (§9).

Two things this doc REFINES (not overrides) in g43, both grounded below:
(a) the gate/helper are **project-owned scripts**, not files added into the
shared `constitution/` submodule (g43 §3.2 was loose here — §11.4.28(B)/§11.4.177
forbid project literals inside the shared engine/governance submodule); (b) the
`.docs_chain/` root is the **MVP project directory itself via `--root`**, not the
`helix_skills` repo root, so the doc set stays self-contained and survives the
G14 extraction into a standalone repo (§4.1).

---

## 2. The real engine context-schema (cited, not invented)

**Schema source (read directly, read-only):** the Docs Chain repo checked out in
the scratchpad during the g43 investigation —
`…/scratchpad/docs_chain_check/` (a clone of `git@github.com:vasic-digital/docs_chain.git`,
outside the `helix_skills` working tree; nothing in this repo was touched). The
two authoritative files:

- **`docs/CONFIG_SCHEMA.md` (Revision 3, Status: IMPLEMENTED)** — the formal YAML
  contract the Phase-4 loader (`internal/config`) parses today (`go test -race
  ./...` green per its own header). This is the field-for-field authority for §4.
- **`.docs_chain/contexts/self-docs.yaml`** — the engine's own dogfood context,
  the exact worked pattern a consumer copies.
- **`cmd/docs_chain/main.go`** — the CLI verb + exit-code contract (§2.3).

If the real schema had not been locatable, this task's mandate is to STOP and
report a BLOCKER rather than invent one (§11.4.6). It WAS located; no schema was
guessed. Every field used in §4 appears in the tables below, each traceable to
`CONFIG_SCHEMA.md`.

### 2.1 Top-level document schema (`CONFIG_SCHEMA.md` §2)

```yaml
context: <string>          # REQUIRED — unique context name; SHOULD match filename stem
description: <string>      # OPTIONAL — human-readable purpose
nodes: <map>               # REQUIRED — node-id -> node-spec (≥1 entry, ids unique)
edges: <list>              # REQUIRED — list of edge-spec (derive-from sub-graph MUST be acyclic)
transforms: <map>          # REQUIRED if any edge names a transform
```

### 2.2 Node / edge / transform specs (the fields §4 uses)

- **Node spec** (`§3`): `{ kind: <node-kind>, path: <project-root-relative path> }`.
  `kind` and `path` are both REQUIRED. `members`/`exclude` are ONLY for
  `kind: fingerprint` (a roster/corpus member-list hash, §11.4.86) — **not used
  here**, because no tracked doc in this set is roster/corpus-backed.
  Node kinds relevant to this design (from the `§3.1` allowed-`kind` table):

  | `kind` | role | direction |
  |---|---|---|
  | `markdown` | canonical `.md` source | input |
  | `html` | pandoc export | derived |
  | `pdf` | weasyprint export | derived |

  (Also available but out of scope: `sqlite`, `summary`, `status`,
  `status_summary`, `fingerprint`. `docx` is a valid derived kind reachable via
  the `pandoc-docx` builtin but is NOT part of the §11.4.65 baseline mandate —
  omitted, see §5.4.)

- **Edge spec — `derive-from`** (one-way; `§4.1`), the only edge shape this
  design uses:

  ```yaml
  - { type: derive-from, from: <node_id>, to: <node_id>, transform: <transform_name> }
  ```

  `from`/`to` MUST exist in `nodes`; `to` MUST NOT introduce a cycle; `transform`
  MUST exist in `transforms`. (The other shape, `type: sync` with `a`/`b`/
  `authority`, is for bidirectional DB↔markdown chains — **not used here**; see
  §5.3 for when it WOULD apply.)

- **Transform spec** (`§5`): exactly one of `{ builtin: <name> }` or
  `{ exec: <cmd>, args: [...] }`. The two builtins this design uses (`§5.1`):

  | builtin | maps | tool |
  |---|---|---|
  | `pandoc-html` | `markdown → html` | pandoc |
  | `weasyprint-pdf` | `html → pdf` (or `markdown → pdf`) | weasyprint |

### 2.3 CLI verbs + exit-code contract (`cmd/docs_chain/main.go`, read directly)

```
docs_chain doctor     [--all | <context>] [--root DIR]   validate contexts (parse+graph+tools, no writes)
docs_chain sync       [--all | <context>] [--root DIR]   propagate atomically, update state.json
docs_chain rebaseline [--all | <context>] [--root DIR]   re-baseline sync edges from AUTHORITY side only
docs_chain verify     [--all | <context>] [--root DIR]   read-only drift check (CI/pre-build gate)
docs_chain graph      <context>           [--root DIR]   print topo order + edges (debug)
docs_chain watch      [--all | <context>] [--root DIR] [--debounce 300ms]   fsnotify daemon
```

Exit codes (documented in `main.go`'s header comment + the `usage()` string):
`0` in-sync/applied/healthy · `1` generic error (bad args, IO, missing contexts
dir) · `2` sync conflict (both sides of a `sync` edge dirty — never silent-merge)
· `3` transform failed, run rolled back, no live changes · `4` cycle/config
error.

**Load-bearing fact for the gate (read from `cmdVerify` source):** `verify`
returns `exit 0` when every selected context reports `in-sync` **OR** an honest
`SKIP (tool absent)`; it escalates to `exit 1` only when at least one node is
`STALE` (`if anyStale { worst = maxExit(worst, exitError) }`). Its per-context
result is a **three-way** print, not two-way — `in-sync` / `STALE: [...]` /
`SKIP (tool absent): <reason>` — driven by a first-class `vr.ToolAbsent` bool, so
the gate never string-matches to tell a real drift from an honest tool-absence
(§2 of §6 below).

### 2.4 Anchors the engine mechanizes (from `docs/CONSTITUTION_INTEGRATION.md` mapping table)

| Anchor | Docs Chain mechanism used here |
|---|---|
| §11.4.65 | derive `.html` + `.pdf` per `.md` (recipe (e),(g)) — **this design** |
| §11.4.12 | derive `summary ← markdown` + early-cutoff — n/a here (no summary docs in set) |
| §11.4.86 | content-hash change detection (NOT mtime) — engine-wide, `state.json` records byte-hashes, `verify` recomputes bytes |
| §11.4.106(E) | typed `ToolAbsentError` — "refusing to fake success" (confirmed in `internal/adapter/adapter.go`) |

---

## 3. MVP doc inventory (real counts, `find` under the MVP dir — not estimated)

Base dir: `docs/research/mvp/Agent_AI_Skill_Tree_Development/` (referred to below
as `<MVP>`). Enumerated directly at design time:

| Group | Path (relative to `<MVP>`) | count | §11.4.44 header? |
|---|---|---|---|
| Top-level tracking | `REQUIREMENTS.md`, `GAPS_AND_RISKS_REGISTER.md`, `IMPLEMENTATION_PLAN.md`, `SPEC.md`, `plan.md` | 5 | **NO** on all 5 (G44) |
| Top-level tracking (has header) | `CONTINUATION.md` | 1 | yes |
| Request ledger | `requests/history.md` | 1 | yes |
| Research design notes | `research/*.md` | **38** | mixed |
| Service product docs | `project/README.md`, `project/docs/{API,ARCHITECTURE,DANGER_ZONES,INSTALL,MCP_INTEGRATION}.md` | 6 | mixed |
| **TOTAL tracked `.md`** | — | **51** | — |

**Honest note on the count (§11.4.6):** `find <MVP> -iname '*.md'` returns **51**
today. This document, once written, becomes the 39th `research/*.md` and the 52nd
overall — so the wiring step will enumerate 52. g43 §2 counted 46 at an earlier
snapshot; the tree grew (research/ is now 38 files). The count is a moving target
by design — which is exactly why the context node-lists must be **generated**
from a `find`, not hand-frozen (§4.4).

**Confirmed zero existing exports:** no `.html`/`.pdf`/`.docx` exists anywhere
under `<MVP>` (the only such files in the whole repo live under `constitution/`,
the governance submodule's own separately-maintained set — not this project's).
So G43's "zero exports exist" is FACT as of this inventory.

**No Status/Issues/Fixed/DB docs in the set:** the only `status`-named artefact
is `project/scripts/status.sh` (a script, not a tracked doc), and there is no
`workable_items.db` present yet (that is the separate G40 adoption effort). So
this design uses **only `derive-from` edges** — no `sync` (DB↔md) edge, no
`summary`/`status_summary` node. §5.3 records what changes when G40's DB lands.

### 3.1 Context grouping (refines g43's 2-context proposal to 3, per the real tree)

g43 proposed a 2-way split (core "living" docs vs research notes). With the
inventory now confirmed, the natural boundary is **3 contexts**, matching the
three physical clusters and their distinct change-cadence / read-frequency:

1. **`skilltree_tracking`** — the 7 top-level "living" tracker docs a resuming
   session reads first: `REQUIREMENTS.md`, `GAPS_AND_RISKS_REGISTER.md`,
   `IMPLEMENTATION_PLAN.md`, `SPEC.md`, `CONTINUATION.md`, `plan.md`,
   `requests/history.md`. High churn → checked by the **per-commit** gate (§6).
2. **`skilltree_research`** — the 38 (→39) `research/*.md` design notes.
   Append-mostly, read less often, larger node count → checked on a **coarser
   cadence** (release-gate sweep, §11.4.40), still `verify`-gated, just not on
   every commit.
3. **`skilltree_service_docs`** — the 6 `project/` product docs (`README.md` +
   `docs/{API,ARCHITECTURE,DANGER_ZONES,INSTALL,MCP_INTEGRATION}.md`). These
   change with the Go service code, live under the `project/` subtree, and are a
   distinct ownership boundary → their own context so a service-code commit that
   touches `project/docs/API.md` re-exports only that context.

The split is a **cadence/ownership** decision, not a correctness one:
`docs_chain sync/verify` operate per-context (`<context>` vs `--all`), so a fast
per-commit gate runs `verify skilltree_tracking` (7 nodes-worth), while `--all`
runs at the release gate. A single combined context would be equally *correct*
(every edge is one-directional `md→html→pdf`, zero cycle risk across groups) —
the 3-way split is purely to keep the per-commit gate cheap. This mirrors g43's
own rationale, refined to the real cluster boundaries.

---

## 4. Concrete context YAMLs (field-valid against §2 schema)

### 4.1 Root anchoring — `.docs_chain/` under `<MVP>`, invoked with `--root`

All node `path`s are **project-root-relative**, and the engine anchors
`.docs_chain/contexts/`, `.docs_chain/state.json`, and `qa-results/docs_chain/`
at whatever `--root` resolves to (confirmed in `loadSelected`: `root =
filepath.Abs(*rootFlag)`, `dir = filepath.Join(root, ".docs_chain/contexts")`;
`state.DefaultPath(root)`; `evidenceDir = filepath.Join(root, "qa-results",
"docs_chain", runID)`).

**Decision: root = `<MVP>` itself**, i.e. the engine is invoked
`docs_chain <verb> --all --root docs/research/mvp/Agent_AI_Skill_Tree_Development`
from the `helix_skills` repo root (or with an absolute path). Consequences:

- context files live at `<MVP>/.docs_chain/contexts/*.yaml`;
- node paths are relative to `<MVP>` (e.g. `REQUIREMENTS.md`, `research/g40_...md`,
  `project/docs/API.md`) — short, and **portable**: when the skill-system is
  eventually extracted into its own standalone repo (the G14 "self-contained and
  clonable in isolation" goal), `<MVP>` becomes that repo's root and the paths
  and `.docs_chain/` move with it **unchanged**;
- `state.json`, `.docs_chain/*.tmp`, and `qa-results/` all land under `<MVP>` →
  add `<MVP>/.docs_chain/state.json`, `<MVP>/.docs_chain/*.docs_chain.tmp`, and
  `<MVP>/qa-results/` to `.gitignore` (mirrors the constitution's own precedent —
  `state.json` is regenerated by `sync`, not authored; the `.html`/`.pdf`
  exports themselves ARE committed, per §11.4.65).

The alternative (root = `helix_skills`, paths like
`docs/research/mvp/Agent_AI_Skill_Tree_Development/REQUIREMENTS.md`) also works
but re-couples every path to the current nesting and does not survive extraction
— rejected for that reason.

### 4.2 `skilltree_tracking.yaml` (fully enumerated — 7 docs, 21 nodes)

```yaml
context: skilltree_tracking
description: Living tracker docs (read first on resume) — md -> html -> pdf per §11.4.65
nodes:
  requirements_md:   { kind: markdown, path: REQUIREMENTS.md }
  requirements_html: { kind: html,     path: REQUIREMENTS.html }
  requirements_pdf:  { kind: pdf,      path: REQUIREMENTS.pdf }
  gaps_md:           { kind: markdown, path: GAPS_AND_RISKS_REGISTER.md }
  gaps_html:         { kind: html,     path: GAPS_AND_RISKS_REGISTER.html }
  gaps_pdf:          { kind: pdf,      path: GAPS_AND_RISKS_REGISTER.pdf }
  implplan_md:       { kind: markdown, path: IMPLEMENTATION_PLAN.md }
  implplan_html:     { kind: html,     path: IMPLEMENTATION_PLAN.html }
  implplan_pdf:      { kind: pdf,      path: IMPLEMENTATION_PLAN.pdf }
  spec_md:           { kind: markdown, path: SPEC.md }
  spec_html:         { kind: html,     path: SPEC.html }
  spec_pdf:          { kind: pdf,      path: SPEC.pdf }
  continuation_md:   { kind: markdown, path: CONTINUATION.md }
  continuation_html: { kind: html,     path: CONTINUATION.html }
  continuation_pdf:  { kind: pdf,      path: CONTINUATION.pdf }
  plan_md:           { kind: markdown, path: plan.md }
  plan_html:         { kind: html,     path: plan.html }
  plan_pdf:          { kind: pdf,      path: plan.pdf }
  history_md:        { kind: markdown, path: requests/history.md }
  history_html:      { kind: html,     path: requests/history.html }
  history_pdf:       { kind: pdf,      path: requests/history.pdf }
edges:
  - { type: derive-from, from: requirements_md,   to: requirements_html,  transform: md2html }
  - { type: derive-from, from: requirements_html, to: requirements_pdf,   transform: html2pdf }
  - { type: derive-from, from: gaps_md,           to: gaps_html,          transform: md2html }
  - { type: derive-from, from: gaps_html,         to: gaps_pdf,           transform: html2pdf }
  - { type: derive-from, from: implplan_md,       to: implplan_html,      transform: md2html }
  - { type: derive-from, from: implplan_html,     to: implplan_pdf,       transform: html2pdf }
  - { type: derive-from, from: spec_md,           to: spec_html,          transform: md2html }
  - { type: derive-from, from: spec_html,         to: spec_pdf,           transform: html2pdf }
  - { type: derive-from, from: continuation_md,   to: continuation_html,  transform: md2html }
  - { type: derive-from, from: continuation_html, to: continuation_pdf,   transform: html2pdf }
  - { type: derive-from, from: plan_md,           to: plan_html,          transform: md2html }
  - { type: derive-from, from: plan_html,         to: plan_pdf,           transform: html2pdf }
  - { type: derive-from, from: history_md,        to: history_html,       transform: md2html }
  - { type: derive-from, from: history_html,      to: history_pdf,        transform: html2pdf }
transforms:
  md2html:  { builtin: pandoc-html }
  html2pdf: { builtin: weasyprint-pdf }
```

### 4.3 `skilltree_service_docs.yaml` (fully enumerated — 6 docs, 18 nodes)

```yaml
context: skilltree_service_docs
description: Skill-graph service product docs under project/ — md -> html -> pdf
nodes:
  readme_md:    { kind: markdown, path: project/README.md }
  readme_html:  { kind: html,     path: project/README.html }
  readme_pdf:   { kind: pdf,      path: project/README.pdf }
  api_md:       { kind: markdown, path: project/docs/API.md }
  api_html:     { kind: html,     path: project/docs/API.html }
  api_pdf:      { kind: pdf,      path: project/docs/API.pdf }
  arch_md:      { kind: markdown, path: project/docs/ARCHITECTURE.md }
  arch_html:    { kind: html,     path: project/docs/ARCHITECTURE.html }
  arch_pdf:     { kind: pdf,      path: project/docs/ARCHITECTURE.pdf }
  danger_md:    { kind: markdown, path: project/docs/DANGER_ZONES.md }
  danger_html:  { kind: html,     path: project/docs/DANGER_ZONES.html }
  danger_pdf:   { kind: pdf,      path: project/docs/DANGER_ZONES.pdf }
  install_md:   { kind: markdown, path: project/docs/INSTALL.md }
  install_html: { kind: html,     path: project/docs/INSTALL.html }
  install_pdf:  { kind: pdf,      path: project/docs/INSTALL.pdf }
  mcp_md:       { kind: markdown, path: project/docs/MCP_INTEGRATION.md }
  mcp_html:     { kind: html,     path: project/docs/MCP_INTEGRATION.html }
  mcp_pdf:      { kind: pdf,      path: project/docs/MCP_INTEGRATION.pdf }
edges:
  - { type: derive-from, from: readme_md,   to: readme_html,  transform: md2html }
  - { type: derive-from, from: readme_html, to: readme_pdf,   transform: html2pdf }
  - { type: derive-from, from: api_md,      to: api_html,     transform: md2html }
  - { type: derive-from, from: api_html,    to: api_pdf,      transform: html2pdf }
  - { type: derive-from, from: arch_md,     to: arch_html,    transform: md2html }
  - { type: derive-from, from: arch_html,   to: arch_pdf,     transform: html2pdf }
  - { type: derive-from, from: danger_md,   to: danger_html,  transform: md2html }
  - { type: derive-from, from: danger_html, to: danger_pdf,   transform: html2pdf }
  - { type: derive-from, from: install_md,  to: install_html, transform: md2html }
  - { type: derive-from, from: install_html,to: install_pdf,  transform: html2pdf }
  - { type: derive-from, from: mcp_md,      to: mcp_html,     transform: md2html }
  - { type: derive-from, from: mcp_html,    to: mcp_pdf,      transform: html2pdf }
transforms:
  md2html:  { builtin: pandoc-html }
  html2pdf: { builtin: weasyprint-pdf }
```

### 4.4 `skilltree_research.yaml` (representative sample + generation rule — 38→39 docs)

Fully enumerating 39 docs × 3 nodes × 2 edges is mechanical, not design-bearing,
and — because `research/` grows every session (§3) — a hand-frozen list would be
STALE the next day. The design decision is: **generate this context file from a
`find`**, not author it by hand. The per-doc triple is identical to §4.2/§4.3.
Representative head (the wiring step emits one such triple per `research/*.md`):

```yaml
context: skilltree_research
description: Research design notes under research/ — md -> html -> pdf (generated from `find research -name '*.md'`)
nodes:
  g40_workable_items_db_adoption_design_md:   { kind: markdown, path: research/g40_workable_items_db_adoption_design.md }
  g40_workable_items_db_adoption_design_html: { kind: html,     path: research/g40_workable_items_db_adoption_design.html }
  g40_workable_items_db_adoption_design_pdf:  { kind: pdf,      path: research/g40_workable_items_db_adoption_design.pdf }
  g43_docs_chain_context_yaml_design_md:      { kind: markdown, path: research/g43_docs_chain_context_yaml_design.md }
  g43_docs_chain_context_yaml_design_html:    { kind: html,     path: research/g43_docs_chain_context_yaml_design.html }
  g43_docs_chain_context_yaml_design_pdf:     { kind: pdf,      path: research/g43_docs_chain_context_yaml_design.pdf }
  # … one { _md, _html, _pdf } triple per research/*.md — 39 docs total …
edges:
  - { type: derive-from, from: g40_workable_items_db_adoption_design_md,   to: g40_workable_items_db_adoption_design_html, transform: md2html }
  - { type: derive-from, from: g40_workable_items_db_adoption_design_html, to: g40_workable_items_db_adoption_design_pdf,  transform: html2pdf }
  - { type: derive-from, from: g43_docs_chain_context_yaml_design_md,      to: g43_docs_chain_context_yaml_design_html,    transform: md2html }
  - { type: derive-from, from: g43_docs_chain_context_yaml_design_html,    to: g43_docs_chain_context_yaml_design_pdf,     transform: html2pdf }
  # … one derive-from pair per research/*.md …
transforms:
  md2html:  { builtin: pandoc-html }
  html2pdf: { builtin: weasyprint-pdf }
```

**Node-id rule (deterministic, so re-generation is stable — §11.4.50):**
`<basename-without-.md, non-`[A-Za-z0-9_]`→`_`>_{md,html,pdf}`. The generator is a
~15-line shell/awk loop the wiring step ships as a project-owned helper
(`<MVP>/scripts/gen_research_context.sh`, §11.4.18-documented); running it after
adding/removing a `research/*.md` re-emits `skilltree_research.yaml`, and the
§11.4.86 content-hash gate (§6) catches any resulting drift. This is the ONLY
generated context — `skilltree_tracking` and `skilltree_service_docs` are stable
enough to hand-maintain and are enumerated in full above.

### 4.5 Decoupling (§11.4.28(B) / §11.4.177)

All three files are **consumer-owned DATA** living under `<MVP>/.docs_chain/`.
They carry every project literal (the doc paths); the engine binary carries none.
They are NEVER written into `constitution/` (that submodule stays project-
agnostic). This is the refinement over g43 §3.2 flagged in §1.

---

## 5. Source→export chains (per-doc semantics)

### 5.1 The chain, per doc (§11.4.65)

`md → html → pdf`, two one-way `derive-from` edges: `md2html` (`pandoc-html`,
markdown→html) then `html2pdf` (`weasyprint-pdf`, html→pdf). PDF derives from the
**html**, not the md, so any HTML-level styling/colorization carries through to
the PDF (matches the engine's own `self-docs.yaml` and `CONFIG_SCHEMA.md §6`).
Propagation: edit a `_md` node → its `_html` regenerates → its `_pdf` regenerates;
edit nothing → early-cutoff → `exit 0`, zero writes (`ARCHITECTURE.md §4`).

### 5.2 Content-hash change detection, NOT mtime (§11.4.86)

The engine records **content hashes** in `state.json` (`internal/state` +
`internal/hash`), and `verify` **recomputes the derived bytes and compares them
to on-disk content** rather than trusting a timestamp (confirmed in `cmdVerify`:
its comment states it "recomputes every derived node and compares the freshly
produced bytes against on-disk content, so it never consults the stored hash
baseline"). This is precisely the §11.4.86 property an ad-hoc `find -newer` sync
script gets wrong (a `touch` with no content change would falsely re-export; a
content change with an unchanged mtime would falsely pass).

### 5.3 Atomic rename + rollback (§9.2 / §11.4.106 — never a partial/fake export)

`sync` is all-or-nothing (`ARCHITECTURE.md §8`): each derived artefact is written
to `<file>.docs_chain.tmp`, `fsync`'d, then the commit phase does an atomic
`rename(2)` of every staged temp over its live file; any error before commit
deletes all temps and aborts, leaving every live file byte-identical to the
pre-run state (`exit 3`). So a failed or tool-absent transform NEVER leaves a
half-written or fake `.html`/`.pdf` on disk.

### 5.4 Honest tool-absence path (§11.4.106(E)) — real, not the expected outcome here

`internal/adapter/adapter.go` defines a typed `ToolAbsentError` whose message is
literally `"adapter: required tool %q not found …; refusing to fake success"`;
callers match it via `errors.As` to SKIP-with-reason. `verify` surfaces it as a
distinct `SKIP (tool absent): <reason>` line (first-class `vr.ToolAbsent` bool),
`doctor` emits a `WARN … not on PATH — runs will SKIP-with-reason`, and `sync`
rolls back rather than writing a partial file.

**On THIS host the SKIP path is a portability safeguard, not the expected result:**
`pandoc` 3.10 (`/home/milos/Factory/software/pandoc/bin/pandoc`) and `weasyprint`
69.0 (`/home/milos/Factory/software/weasyprint/bin/weasyprint`) both resolve
(re-confirmed via `command -v` at design time). So real transforms will run here;
the SKIP path only matters on a future host/CI runner lacking the tools — where
the correct behavior is an honest SKIP + a tracked install follow-up (§11.4.184
pattern), never a fake pass. The §6 mutation set proves the SKIP path fires and
is never mistaken for `STALE` or `in-sync`.

### 5.5 When G40's workable-items DB lands (out of scope, recorded for completeness)

Nothing in the current doc set is DB-backed, so all edges are `derive-from`. When
the separate G40 effort lands `workable_items.db`, a fourth context
(`skilltree_issues`) would add a `type: sync` edge (`markdown ↔ sqlite`,
`authority: <db-node>`) plus `summary`/`html`/`pdf` derivations, exactly per
`CONFIG_SCHEMA.md §7`. That is a G40 concern; this design does not pre-build it.

---

## 6. `CM-DOC-EXPORTS-SYNC` gate + paired §1.1 mutation

### 6.1 Gate design — `docs_chain verify` as the deterministic PASS/FAIL/SKIP oracle

A **project-owned** pre-build/CI gate `CM-DOC-EXPORTS-SYNC` (living under the
project's own `<MVP>/scripts/` or the repo's project gate dir — NOT inside the
shared `constitution/` submodule, §11.4.28(B)/§11.4.177) runs, read-only:

```bash
docs_chain verify --all --root docs/research/mvp/Agent_AI_Skill_Tree_Development
```

Verdict mapping (deterministic, §11.4.50 — read straight off `verify`'s exit code,
NO string-matching needed for the pass/fail decision, per §2.3):

| `verify` outcome | exit | gate verdict |
|---|---|---|
| every context `in-sync` | 0 | **PASS** |
| every context `SKIP (tool absent)` (or a mix of in-sync + skip) | 0 | **PASS (honest SKIP-with-reason)** — §11.4.3, non-blocking |
| any context `STALE: [...]` | 1 | **FAIL** — a source `.md` changed without a `sync`; commit/build refuses |
| conflict / config / IO error | 2 / 4 / 1 | **FAIL** (loud; not applicable to derive-only contexts for exit 2, but handled) |

The gate additionally parses `verify` stdout ONLY to emit a human-readable line
naming the stale nodes (on FAIL) or the skipped context+reason (on SKIP) into its
evidence log — the *decision* is the exit code, the *message* is for the operator.
Because `verify` already separates `SKIP` (exit 0) from `STALE` (exit 1), the gate
never risks counting an honest tool-absence as a failure or a real drift as a pass
(the exact conflation g43 §1.6 warned about — resolved by using `verify`, never
`sync`, as the gate).

Per-cadence invocation (§3.1): the per-commit gate may narrow to
`verify skilltree_tracking` for speed; the release-gate sweep runs `verify --all`.

### 6.2 Paired §1.1 mutation (RED→GREEN polarity, §11.4.115) — proves the gate is load-bearing

Three mutations, each with its golden fixture, wired into the meta-test so the
gate cannot be a tautology:

1. **STALE-drift mutation (the primary §1.1 pair).**
   - *Golden-good:* after `docs_chain sync --all`, `verify --all` prints
     `in-sync` for all three contexts, exit 0 → gate PASS.
   - *Mutation:* append one line to `REQUIREMENTS.md` (a real source edit) WITHOUT
     re-running `sync`. `verify --all` MUST print
     `skilltree_tracking   STALE: [requirements_html requirements_pdf]`, exit 1 →
     gate **FAIL**.
   - *Restore:* revert the edit (or re-`sync`) → `verify` returns to `in-sync`,
     exit 0 → gate PASS. This is the RED-on-drift / GREEN-on-synced flip that
     proves the gate detects real staleness, not a fixed answer.
2. **Backdate/desync mutation (mtime-vs-content proof, §11.4.86).** After a clean
   `sync`, `touch -d '2020-01-01' REQUIREMENTS.html` (older mtime, IDENTICAL
   bytes). `verify` MUST still report `in-sync` (content unchanged) — proving the
   gate keys on content hash, not mtime, and does NOT false-FAIL on a mere
   timestamp change. Conversely, edit `REQUIREMENTS.html` bytes directly (a
   hand-tampered export) → `verify` MUST report `STALE` (recomputed bytes ≠
   on-disk) → gate FAIL. This is the exact "backdate/desync one export → gate
   FAILs" the task calls for, plus its mtime-robustness twin.
3. **Tool-absent SKIP mutation (proves SKIP ≠ FAIL and SKIP ≠ fake pass).** Run
   `verify --all` with `PATH` scoped to exclude
   `/home/milos/Factory/software/{pandoc,weasyprint}/bin`. `verify` MUST print
   `SKIP (tool absent): <reason>` per context, exit 0 → gate PASS-with-honest-SKIP
   (never counted as FAIL, §11.4.3); and `sync` under the same PATH MUST roll back
   (`ROLLED-BACK: SKIP (tool absent)`) leaving zero new bytes — proving no partial
   export and no fake pass.

Meta-test assertion: mutation (1) and (2b) MUST make the gate FAIL; the
golden-good and (2a) and (3) MUST make it PASS. A gate that stays green under
mutation (1) is itself the bluff (§11.4.107(10) analyzer-self-validation applied
to the gate).

---

## 7. Wiring seams (§11.4.106(F)) + honest boundary

Three write-seams enforce eventual-consistency-at-every-write; `sync` is the only
verb that WRITES the exports, `verify` is the read-only gate at each seam.

1. **Commit seam.** A project-owned pre-commit path (helper
   `<MVP>/scripts/sync_docs_chain.sh`, thin wrapper analogous to
   `constitution/scripts/codegraph_sync.sh`) runs `docs_chain sync --all --root
   <MVP>` when a staged change touches any tracked source `.md`, then
   `docs_chain verify --all` to confirm `in-sync` BEFORE the commit is allowed —
   so a committed `.md` never lands with a stale/absent `.html`/`.pdf`. (This is
   the §11.4.12/§11.4.53-class "regenerate-then-verify before commit" discipline
   this Constitution already mandates for the Issues/Fixed family, applied to this
   doc set.)
2. **Build seam.** The `CM-DOC-EXPORTS-SYNC` pre-build gate (§6) runs `verify
   --all` read-only; the build refuses to proceed on any `STALE`.
3. **Constitution-pull seam.** The §11.4.164 `post_update_hook.sh` re-runs
   `verify --all` after a constitution pull (these docs cite constitution anchors;
   an anchor-renumbering pull could invalidate a reference the source author must
   then refresh — the gate surfaces the resulting need to re-sync). Docs Chain
   does NOT author the §11.4.44 header or the anchor text — it only keeps the
   derived exports in lock-step with whatever the source currently says
   (`CONSTITUTION_INTEGRATION.md §4`: "Docs Chain does NOT replace the authoring
   discipline … it replaces the mechanical sync").

**Honest boundary (§11.4.6 / §11.4.106(F) as-written):** these seam hooks give
**eventual-consistency-at-every-write** — the exports are guaranteed in-sync at
each commit, each build, and each constitution pull, NOT continuously in
real-time between those events. Literal real-time would need the `docs_chain
watch` fsnotify daemon (a separate, opt-in long-running process); that is a
tracked §11.4.197 upgrade path, NOT claimed as shipped by this design. The seam
design is exactly the boundary §11.4.106(F) states.

---

## 8. Vendoring dependency (on the G14 decision — cross-referenced, not re-decided)

The context YAMLs (§4) and the gate (§6) are **agnostic to where the `docs_chain`
binary physically lives** — they only invoke `docs_chain <verb> --root <MVP>`.
What the binary's *provenance* depends on is the G14 vendor-fresh submodule
decision, already made by the operator and designed elsewhere; this doc only
states the dependency:

- **Operator decision (`CONTINUATION.md` Rev 8, via `research/g14_vendor_fresh_submodule_layout_design.md`):**
  VENDOR FRESH — each dependency a fresh git submodule under this project's own
  `submodules/<snake_case>/`, self-contained/clonable in isolation.
- **For `docs_chain` specifically (`g14_vendor_fresh_submodule_layout_design.md`
  §2.B row 17, cross-referencing `g43_docs_chain_export_wiring_design.md` §6/§7):**
  the **recommended** landing is `constitution/submodules/docs_chain/` (the same
  §11.4.28(C) depth-1 carve-out already used for `token_optimizer` /
  `session_orchestrator` / `continuum`), **pending the constitution-submodule
  maintainers landing it there** (a constitution-level decision, outside this
  project's unilateral scope). The **interim fallback**, if this project needs
  export wiring before that lands, is to vendor fresh at this project's own
  `submodules/docs_chain/` exactly like any Category-A dependency, then **retire**
  the interim copy the moment `constitution/submodules/docs_chain/` exists (a
  tracked, evidence-backed §11.4.197 migration — never a permanent duplicate).

**This design does NOT re-decide** which of those two lands first. It records the
dependency: the ONLY thing the §4 YAMLs / §6 gate care about is that a built
`docs_chain` binary is resolvable at wiring time (via `command -v docs_chain`, or
a project-declared `DOCS_CHAIN_BIN`/`docs_chain` symlink pointing at whichever
checkout the G14 decision produces). Both landing options satisfy that
identically. The export gap (G43) can therefore be designed to completion now and
wired the moment the G14 binary-provenance step lands — neither blocks the other's
design.

---

## 9. Anti-bluff captured-evidence (self-validated per §11.4.106)

- **Per-`sync`-run evidence (engine-produced, confirmed in source).**
  `runSyncContexts` writes `qa-results/docs_chain/<run-id>/sync.log` (run-id =
  UTC `20060102T150405Z`, one dir per run) containing every context's per-node
  result line (`writeEvidence`, `os.WriteFile(… kind+".log" …)`). This is the
  §11.4.69 captured artefact proving a real transform ran. With `--root <MVP>` it
  lands at `<MVP>/qa-results/docs_chain/<run-id>/sync.log` (gitignored per §4.1).
- **Gate evidence (`verify` — the gate wrapper captures it, NOT the engine).**
  Confirmed: `cmdVerify` does **not** call `writeEvidence` (zero matches in
  source) — `verify` is pure read-only stdout. So the `CM-DOC-EXPORTS-SYNC`
  wrapper is responsible for capturing its own transcript (exit code + the
  per-context `in-sync`/`STALE`/`SKIP` lines) into
  `<MVP>/qa-results/docs_chain/verify_<ts>.log`, and citing that path as the
  gate's captured evidence. This is a design requirement on the wrapper, stated
  explicitly so the wiring step does not assume `verify` self-writes evidence.
- **Self-validation (§11.4.106 / §11.4.107(10)).** Two layers: (a) the engine
  itself is self-validated upstream — `go test -race ./...` green across all
  packages (`g43_docs_chain_export_wiring_design.md` §1.2, proven in the scratch
  clone); (b) the gate's own §1.1 mutation triple (§6.2) is the golden-good /
  golden-bad(STALE) / negative-control(SKIP+mtime-robust) set proving the gate
  distinguishes a real drift from an honest tool-absence from a mere timestamp
  change — a gate that passes its golden-bad STALE fixture is itself the bluff and
  a release blocker.
- **Real-export proof (State-B, once wiring lands).** After the first real `sync
  --all`: each of the 52 tracked docs has a byte-present, non-empty `.html` +
  `.pdf` sibling; an immediately-following `verify --all` reports `in-sync` for
  all three contexts (no false-STALE flapping → proves determinism, §11.4.50).

---

## 10. Honest gaps (§11.4.6)

- No `.docs_chain/` file, no gate script, no export, and no binary was created or
  run inside this repo by this task (read-only mandate). The schema/CLI/exit-code
  facts in §2 were read from the g43 scratch clone (outside `helix_skills`); the
  YAML in §4 is proposed content, field-validated against `CONFIG_SCHEMA.md` but
  not yet parsed by a live `docs_chain doctor` in this repo.
- The `skilltree_research` per-context `verify` latency at ~39×3 nodes was not
  measured here (the engine was never invoked in-repo) — the §3.1 per-commit vs
  release-gate cadence split is a recommendation the wiring step must confirm
  empirically (`PENDING_FORENSICS:` — same open item g43 §7 flagged).
- Whether `docs_chain` lands at `constitution/submodules/docs_chain/` vs this
  project's own `submodules/docs_chain/` is a G14/constitution decision this doc
  cross-references (§8) but does not make.
- DOCX siblings are out of the §11.4.65 baseline mandate; a `pandoc-docx` builtin
  exists (used by the engine's own `self-docs.yaml`) but is not designed in.
- The §11.4.44 revision-header backfill for the 5 header-less top-level docs (G44)
  is a coupled but independent effort (g43 §5) — Docs Chain exports whatever the
  source says, header or not; recommend landing G44 in the same batch so the
  first `sync` produces header-carrying exports.
```