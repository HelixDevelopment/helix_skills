# G40 — Workable-Items SQLite SoT Adoption Design (§11.4.93/§11.4.95)

**Revision:** 1
**Last modified:** 2026-07-15T18:20:00Z
**Scope:** DESIGN ONLY. No code, no schema change, no data migration performed by
this document. The actual migration is a separate, explicitly gated step
(§11.4.101/§11.4.197) — see §6 below for exactly what "gated" means here.
**Authority:** Constitution §11.4.93 (SQLite SSoT) / §11.4.95 (DB tracked in
git) / §11.4.74 (extend-don't-reimplement) / §11.4.28 (decoupling) / §11.4.66
(operator-decision surfacing) / §11.4.186 (anti-divergence) / §11.4.6
(no-guessing).
**Finding this design closes:** `docs/research/mvp/Agent_AI_Skill_Tree_Development/GAPS_AND_RISKS_REGISTER.md`
**G40** (HIGH, §11.4.93/§11.4.95) — "no SQLite `workable_items.db`
single-source-of-truth; 39 findings tracked in this markdown register instead.
The constitution's own Go workable-items engine
(`constitution/scripts/workable-items/`) exists + is unused → adoption-only."
Also closes the coupled **G45** (status/type/reopens closed-vocabulary not
applied to the G0x findings) and resolves the **G47** (§11.4.54, OPERATOR-DECISION
PENDING) tension between the project's existing `G0x`/`R0x`/`P0x` id scheme and
the constitution's `ATM-NNN` convention.

---

## 0. Honest scope correction (§11.4.6)

The dispatching prompt for this design stated "49 findings (G01–G49)". The
register as read on 2026-07-15 confirms this is accurate as a **count** (G01
through G49 inclusive = 49 ids), but the register's own summary table at its
top (`Summary counts`, lines 14–22) is **stale** — it only tallies G01–G27 (27
rows, last updated before G28–G49 were appended). G39–G49 are NOT `### `
H3-heading findings like G01–G38; they are compact bullet-list entries under a
"§11.4.201 gate-methodology note" section (register lines 405–459), each
already carrying a `**STATUS:**` word (`FILED`, `PENDING`, etc.) but not yet
the closed-set §11.4.15 vocabulary. This distinction (rich H3 findings vs.
compact bullet findings) is itself a field-mapping fact captured in §3 below —
it is not a defect, but the migration parser must handle both shapes.

`REQUIREMENTS.md` carries **R01–R24** (24 requirement clusters, confirmed by
`grep -c` — matches the dispatch prompt). `IMPLEMENTATION_PLAN.md` carries
`P0`/`P0.5`/`P1`/… phase headings with per-item lines named e.g. `P0.5.G01`,
`P0.5.G02` — these are NOT independent ids; every `P0.5.Gxx` line is the
**action-plan restatement of the matching Gxx finding** (compare
`IMPLEMENTATION_PLAN.md:127-133` "P0.5.G01 — Wire `internal/api`…" against
`GAPS_AND_RISKS_REGISTER.md` G01's `DECISION:` line — same fix, same file:line
citations). This is the load-bearing fact behind the DEDUP decision in §3.3.

---

## 1. The constitution engine — captured schema + CLI contract (read-only investigation)

### 1.1 Build status

```
cd constitution/scripts/workable-items && go build -o <tmp> ./cmd/workable-items/
```
exits **0** — the binary builds cleanly (confirmed 2026-07-15, produced an
8.9 MB executable). No `go vet`/test run was required by this design task;
the build-clean fact alone is sufficient to answer "does the engine exist and
compile" per the dispatch brief.

### 1.2 CRITICAL pre-existing finding: `schema.sql` and `schema_embed.sql` have DRIFTED

The package ships **two** DDL files:
- `constitution/scripts/workable-items/schema.sql` — a loose reference copy at
  the package root.
- `constitution/scripts/workable-items/cmd/workable-items/schema_embed.sql` —
  the copy actually `//go:embed`-ded into the binary and applied to every DB
  the tool opens (confirmed by `db.go` invoking the embedded schema on open).

A `diff` of the two (performed as part of this investigation) shows they have
**diverged** — each has columns/tables the other lacks:

| Present in `schema.sql` (loose) but **NOT** in `schema_embed.sql` (runtime) | Present in `schema_embed.sql` (runtime) but **NOT** in `schema.sql` (loose) |
|---|---|
| `items.canonical_track` column (§11.4.191 track-pin) | `items.representation` column + `(atm_id, current_location, representation)` composite PK (v6 pipe-table-vs-section dedup) |
| `group_paths` table (§11.4.191 file-scope manifest) | `items.closure_date` / `items.round` / `items.commit_ref` (Fixed.md pipe-row synthesis) |
| — | `items.parent_atm_id` / `items.session_ref` (§11.4.148/§11.4.149 sub-task hierarchy) |
| — | `obsolete_details.reason` extra value `'not-reproducible'` |
| — | `test_diary` table + `test_diary_summary` VIEW (§11.4.149 — entirely absent from `schema.sql`) |

**This is an honest, pre-existing engine-internal inconsistency, out of this
project's scope to fix** (it lives in the constitution submodule, a shared
reusable engine per §11.4.28/§11.4.74 — this project consumes it by
reference and must not silently patch it). It is reported here because it
determines which schema this design targets: **`schema_embed.sql` is
authoritative** for adoption purposes, because it is what the compiled binary
actually applies to `docs/workable_items.db` at runtime. Every field-mapping
decision below is against `schema_embed.sql` (schema_version `'6'`, confirmed
by its `meta` seed row), not the loose `schema.sql`. The `canonical_track` /
`group_paths` absence from the runtime schema means this project's optional
Phase-2 multi-track wiring (§3.4) is **not currently usable out of the box**
without an upstream extension — flagged as Gap E in §4.

### 1.3 CLI contract (from `cmd/workable-items/main.go` usage text + confirmed build)

Relevant subcommands for this adoption (grouped by phase they serve):

| Subcommand | Purpose | Phase it serves |
|---|---|---|
| `sync md-to-db --db <p> [--issues <p>] [--fixed <p>]` | Parse Issues.md + Fixed.md shape, upsert DB | Phase 1 import |
| `sync db-to-md --db <p> [--out-issues <p>] [--out-fixed <p>]` | Regenerate trackers FROM the DB | Phase 3 cutover |
| `diff --db <p> [--issues <p>] [--fixed <p>]` | Show DB-vs-Markdown divergence | Phase 2 round-trip proof |
| `validate --db <p>` | Closed-set + §11.4.91 clarity invariants | every phase (pre-build gate) |
| `add <type> <severity> --db <p> --title <T> --description <D> [--id <id>] [--prefix <P>] [--created-by <h>] [--assigned-to <h>]` | Create one Queued item, with **explicit `--id` override** | Phase 1 (per-item import, or G0x-alias path) |
| `update --id <ID> --db <p> [--status\|--type\|...]` | Mutate an existing item | ongoing use post-adoption |
| `reopen` / `block` / `close` | §11.4.34 / §11.4.21 / §11.4.19 lifecycle transitions | ongoing use post-adoption |
| `report --db <p> [--by-type\|--by-status\|...]` | Read-only grouped tally | validation/audit |
| `group add/list/set/state` + `validate-groups` + `assign` | §11.4.176/§11.4.191 multi-track work-division | optional Phase-2 (see Gap E) |
| `export --db <p> [...]` | Regenerate Issues.md + Fixed.md + Summaries + HTML/PDF/DOCX | Phase 3 cutover (§11.4.12/.53/.65) |

**Load-bearing detail for G47:** `add` accepts an explicit `--id <id>` flag —
the engine does **not** force an `ATM-NNN`-shaped id; `items.atm_id` in
`schema_embed.sql` is declared `TEXT NOT NULL` with **no format `CHECK`**
(only `type`/`status`/`representation` have `CHECK` constraints; `atm_id`
does not). This is the schema fact that resolves G47's tension — see §3.2.

---

## 2. Data to migrate — structure of the two existing streams

### 2.1 `GAPS_AND_RISKS_REGISTER.md` (49 findings, G01–G49)

Two distinct SHAPES coexist in the same file:

**Shape A — rich H3 findings (G01–G38).** Heading `### Gxx — <one-line title>`,
followed by an unordered list of labelled fields in this order (not all
present on every item):
- `**Category:**` (free text: `inconsistency`, `security`, `gap`,
  `existing-bug`, `danger-zone`, `test-coverage`)
- `**Severity:**` (prose sentence, but always opens with
  `critical`/`high`/`medium`/`low`)
- `**Evidence:**` (nested list of `file:line` citations + quoted code)
- `**Why it matters:**` (prose)
- `**DECISION:**` (the accepted fix direction + "Alternatives rejected: …")
- optionally **`- **STATUS (<date>) — <label>:**` or a bare `**STATUS:**`
  line** — free-text narrative, sometimes very long (G01's STATUS block runs
  to ~35 lines covering two remediation attempts, a Fable-xhigh review
  verdict, and two tracked follow-ups O1/O2/O3)
- `**Test coverage:**` (comma-separated test-type list + explicit mutation
  description)
- trailing `**Challenges:**` / `**HelixQA:**` yes/no markers

**Shape B — compact bullet findings (G39–G49).** One bullet per finding,
form: `- **Gxx (SEVERITY, §a.b.c/§d.e.f):** <one-paragraph combined
evidence+decision text>. **STATUS:** <label>.` — no separate Category/Evidence/
DECISION/Test-coverage sub-fields; everything is inline prose.

Both shapes carry, somewhere in text, a severity word from the closed set
`{CRITICAL, HIGH, MEDIUM, LOW}` (register's own summary table, confirmed
consistent with every G39-G49 bullet's parenthetical) and a §-section citation
to the constitution anchor(s) the finding violates.

There is also an **"Adjudication of the 8 mandated open items"** table (lines
463-474) that cross-references G-ids to operator-numbered decision points, and
a **"Resolved / non-findings"** section (476-482) of prose bullets that are
explicitly **not** separate trackable items (verified-clean facts) — these
must **not** become DB rows (they are closure evidence for other findings, not
findings themselves).

### 2.2 `REQUIREMENTS.md` (R01–R24)

A **living narrative document**, not a per-item ledger. The 24 `R`-clusters
appear as bullet points under "## Requirement clusters (from operator, in
order received)" (`REQUIREMENTS.md:51-97`+), each `- **Rn — <short name>.**
<prose defining scope>`. There is **no** per-R Status/Type/Evidence/Decision
field structure — R-items are requirement/scope statements the project is
still building toward, not individually closable defects. The document also
carries a `## Baseline` section explicitly marked **SUPERSEDED** (kept as the
§11.4.7 regression oracle, never treated as current state) and an
`## Addenda` section for new operator mandates folded in later.

### 2.3 `IMPLEMENTATION_PLAN.md` (`P0`, `P0.5`, `P1`, …)

Phase-and-task headings (`## P0.5 — CRITICAL remediation …`) with per-task
bullets named `P0.5.G01`, `P0.5.G02`, etc. Per §0/§2.1 above, these are
**restatements** of the matching `Gxx` finding's DECISION as an action item —
not a third independent id space.

### 2.4 The G47 tension (verbatim from the register)

> **G47 (MED, §11.4.54 — OPERATOR-DECISION):** project uses G0x/R0x/P0x ids
> cross-referenced across 20+ docs, not ATM-NNN; a destructive rename is
> riskier than (a) documenting G0x as an approved alias OR (b) minting
> parallel ATM-NNN without renaming. Surface via §11.4.66. **STATUS:** FILED;
> operator-decision PENDING.

This design does **not** resolve G47 — it is explicitly an operator decision
(§11.4.66) — but §3.2 below shows both paths (a) and (b) are cleanly
expressible in the runtime schema, so neither choice blocks adoption.

---

## 3. Adoption design

### 3.1 Model: DB as SoT, register + REQUIREMENTS as DERIVED

Per §11.4.93/§11.4.106/§11.4.12: `docs/workable_items.db` (this project's own
DB, git-tracked per §11.4.95, WAL-checkpointed before every commit) becomes
the single source of truth for every G-item. `GAPS_AND_RISKS_REGISTER.md`
(and, for the subset of R-items promoted to DB rows — see §3.3 — the relevant
parts of `REQUIREMENTS.md`) become **generated views**: `workable-items sync
db-to-md` + `workable-items export` regenerate them, mirroring the existing
Issues.md/Fixed.md pattern the engine was built for. `IMPLEMENTATION_PLAN.md`
stays **hand-authored narrative** (phase sequencing, cross-stream
coordination) — it is not a generated doc; it *references* G-ids via the
`composes_with` mechanism the schema already provides, it does not duplicate
their content.

### 3.2 Field-mapping table (register field → `schema_embed.sql` v6 column)

| Register / REQUIREMENTS field | Engine column (table.column) | Fit | Notes |
|---|---|---|---|
| `Gxx` / `Rxx` id | `items.atm_id` (TEXT, **no format CHECK**) | **Clean, zero schema change** | G47 Option (a): store the literal string `"G40"`/`"R09"` as `atm_id` — schema already accepts any non-null text. G47 Option (b) (mint parallel `ATM-NNN`, keep `Gxx` as a documented alias) needs a NEW nullable `legacy_id` column or a tiny side-table — **schema gap**, see Gap D below. This design **recommends Option (a)** as the adoption default: zero engine changes, fully non-destructive (the literal string used in 20+ cross-referencing docs is preserved verbatim as the DB's own primary key), and it is still a "stable, unique id" in the §11.4.54 sense even though its prefix is `G`/`R` rather than `ATM`. Final choice stays the operator's per G47. |
| one-line title (after the id, in the `### ` heading or bullet lead) | `items.title` | Clean | |
| `**Severity:**` prose | `items.severity` (free TEXT, no CHECK) | Clean | Extract the leading `critical/high/medium/low` token; keep full sentence too — `severity` is unconstrained TEXT so both fit. |
| `**Category:**` (`inconsistency/security/gap/existing-bug/danger-zone/test-coverage`) | **no matching column** | **Gap A** | See §4. Recommendation: fold into `items.description`'s structured header (e.g. a `Category: security` line) until an upstream `category` column lands, OR request the one-line additive `ALTER TABLE items ADD COLUMN category TEXT` from the constitution maintainers (low-risk, matches the `severity` column's own precedent of being an unconstrained informational field). |
| `**Evidence:**` + `**Why it matters:**` | `items.description` (`NOT NULL`, ≥40 chars/§11.4.91 enforced at insert) | Clean | This is exactly the "what/reproduction" half of the §11.4.148 D2 comprehensive-description contract; concatenate Evidence + Why-it-matters as the description body. |
| `**DECISION:**` (fix direction + alternatives rejected) | `items.closure_criteria` | Clean | Named for exactly this purpose ("Closure criteria (markdown body)"). |
| `**Test coverage:**` / `**Challenges:**` / `**HelixQA:**` lines | append to `items.closure_criteria` | Clean | These are "what proves this is closed" — same field as DECISION, appended, not a separate column. |
| free-text `**STATUS:**` narrative (`FILED`, `DESIGN DONE`, the long G01 attempt-1/attempt-2 block, etc.) | (a) classified into `items.status` (closed 10-value set) **+** (b) the FULL raw narrative preserved via `items.body_md` / `doc_segments` (§11.4.93 byte-identical round-trip mechanism) | **Requires a mapping table, not 1:1** | See Gap B below — free-text status words do not injectively map onto the closed set; a human/engine judgment call table is required per item. The raw prose is never lost regardless: `body_md` captures the ENTIRE `### Gxx … ` block verbatim (or the entire bullet for Shape B) so `db-to-md` reproduces it byte-for-byte even while `items.status` carries a coarser closed-set classification. |
| the whole raw `### Gxx …` block (Shape A) or bullet (Shape B) | `items.body_md` + one `doc_segments` row (`document='Issues'`, `kind='item'`, `atm_id=<id>`) | Clean — this IS the mechanism the engine was designed for | The surrounding narrative (audit preamble, the "Summary counts" table, the "Adjudication" table, "Resolved / non-findings" section) becomes `doc_segments` rows with `kind='raw'` so `db-to-md` reproduces the WHOLE file, not just the item bodies. |
| §11.4.16 `Type` (`Bug\|Feature\|Task`) | `items.type` (CHECK-enforced) | **Requires per-item classification** | See Gap C. |
| §11.4.15 `Status` transitions implied by narrative (e.g. G01's two-attempt cycle) | `item_history` rows (`event_type`, `by`, `on_date`, `reason`, `evidence_path`) | **Backfill is approximate for history predating this design; exact going forward** | See Gap B. |
| `created_by` / `assigned_to` (§11.4.104) | `items.created_by` / `items.assigned_to` (TEXT, default `''`) | Clean, low-stakes default | This project ships no messenger surface (§11.4.104 binds "latently"). Recommended default: `created_by='Claude'` (the audit was AI-authored per the register's own preamble: "Adversarial audit satisfying operator mandate R17"), `assigned_to=''` (unassigned/Operator-default per §11.4.104(C)). |
| `current_location` (`Issues`\|`Fixed`) | `items.current_location` (CHECK-enforced) | Clean, all-`Issues` at import | No G-item is fully closed as of this audit (even G01's attempt-2 leaves O3 open) — every imported row lands in `Issues`. Fine and expected, not a gap. |
| §11.4.176/§11.4.191 track/branch binding (`logic_group`, `destination`) | `items.logic_group` / `items.destination` (nullable) | **Optional, out of core scope** | This project runs a single serialized Go-mutator lane (CONTINUATION.md: "P0.5 critical-remediation spine … single serialized Go-mutator lane"), not the formal multi-track work-division the schema anticipates. Leave `NULL` at import; populate only if/when the project adopts §11.4.176 multi-track dispatch. |
| §11.4.191 `canonical_track` + `group_paths` file-scope manifest | **absent from the runtime `schema_embed.sql`** | **Gap E (engine-level, out of this project's scope)** | See §1.2 / §4. |
| Per-run test results (future, once P0.5 tests actually execute) | `test_diary` table + `test_diary_summary` VIEW (§11.4.149) | Clean, future use | NOT part of the initial import — the register's `Test coverage:` lines are a **plan** (what test types SHOULD prove closure), which maps to `closure_criteria` (above); `test_diary` records **actual executed runs** with PASS-requires-evidence enforced by a `CHECK` constraint. Wire this once P0.5's tests genuinely run and produce evidence paths. |
| `R01`–`R24` requirement clusters | `items.atm_id` = `"R01"`.."R24"`, `items.type` likely `Feature` (defines a capability to build) | **Requires per-item judgment + a DEDUP decision** | See Gap C / §3.3. |
| `P0.5.Gxx` action-plan lines (`IMPLEMENTATION_PLAN.md`) | **NOT separate rows** — referenced via `items.composes_with` (JSON array of §-letter/id refs) pointing back at the owning `Gxx` | Clean — this IS the §11.4.186 anti-divergence-mandated dedup | See §3.3. |

### 3.3 The DEDUP decision (§11.4.186 anti-divergence, forensic precedent cited in the constitution: the SPK-596 "same feature, six rows, two releases" incident)

`P0.5.G01`, `P0.5.G02`, … in `IMPLEMENTATION_PLAN.md` are the **same
underlying work** as `G01`, `G02`, … in the register (verified by comparing
their file:line citations and fix descriptions — identical). Creating a
SECOND DB row per `P0.5.Gxx` would recreate exactly the divergence-risk
pattern §11.4.186 was written to forbid. **Decision: `P0.5.Gxx` items are NOT
imported as separate DB rows.** `IMPLEMENTATION_PLAN.md` stays hand-authored
narrative that cites `Gxx` ids directly; if a machine-checkable cross-link is
wanted later, it can be expressed as a `composes_with` JSON entry on the `Gxx`
item pointing at the plan section — never a duplicate row.

R01–R24 are lower-priority for DB import than the G-items (they are scope
statements, not closable defects) but are explicitly in the task's scope as
"a second workable-item stream." Recommendation: import them as `Type=Feature`
rows by default (each defines a capability the system must gain), with an
explicit per-R override list for the ones that are clearly procedural/
infrastructure (`Type=Task` — e.g. R9 "Submodule resolution + sync", R10 "Docs
Chain incorporation") rather than a new user-facing capability. This
per-R classification is a judgment call (Gap C), not a mechanical rule — it
should be reviewed by the engineer running the Phase-1 import, not
auto-assigned.

### 3.4 Engine reuse vs. extend (§11.4.74 extend-don't-reimplement)

**Verdict: the engine covers this project's core need (G0x findings +
R0x requirements as `items` rows with full-narrative round-trip) with ZERO
required code changes**, using the `add --id <id>` / `sync md-to-db` /
`sync db-to-md` / `export` / `validate` subcommands exactly as they exist
today. The ONLY genuinely-needed extension, if the operator wants it, is:

- **Gap A** (`category` column) — a one-line additive `ALTER TABLE items ADD
  COLUMN category TEXT` (matches the existing precedent of `severity` being
  an unconstrained informational column) — **recommended upstream PR**, low
  risk, purely additive, no `migrateColumns`-style backfill logic needed
  beyond what the engine's existing `ADD COLUMN`-on-open pattern already does
  for `parent_atm_id`/`session_ref`/etc.
- **Gap D** (`legacy_id` alias column, only if G47 resolves to Option (b)) —
  same additive-column shape as Gap A.
- **Gap E** (`canonical_track` + `group_paths` missing from the runtime
  schema) is an **engine-internal drift bug** (schema.sql vs. schema_embed.sql
  disagree with each other) that exists independently of this project and
  should be reported to the constitution submodule maintainers as its own
  finding — it is NOT something this project's adoption should route around
  by re-adding those tables locally (that would violate §11.4.28 decoupling —
  a consumer must never patch a shared engine's schema in its own tree).

No fork is proposed anywhere in this design. Everything above is either
"use as-is" or "request an additive upstream extension."

---

## 4. Honest gaps summary (§11.4.6)

| Gap | Description | Blocking? | Resolution path |
|---|---|---|---|
| **A** | No `category` column in `schema_embed.sql` | No — fold into `description` header line as a workaround | Recommend upstream additive column |
| **B** | Free-text `**STATUS:**` narratives do not injectively map onto the §11.4.15 10-value closed set; historical multi-attempt cycles (G01) can only be **approximately** reconstructed into `item_history` rows (dates/actors for the ORIGINAL narrative were not captured in the §11.4.34 By/On/Reason/Evidence shape at authoring time) | No — the raw narrative is preserved byte-identically via `body_md` regardless of classification accuracy; only the closed-set field needs human judgment | Manual classification table authored during Phase 1 (see §5), reviewed in code-review per §11.4.125/§11.4.209; going FORWARD every new state change is captured with full §11.4.34 fields, so the approximation is a one-time historical-backfill limitation only |
| **C** | `Category`→`Type` and `Rxx`→`Type` mapping is a judgment call, not mechanical (e.g. is "gap" a `Feature` or a `Task`?) | No — defaults proposed in §3.2/§3.3, reviewable, overridable per item | Engineer classification pass during Phase 1 import, captured as an explicit mapping table in the migration's own commit message (§11.4.8/§11.4.92 audit trail) |
| **D** | G47's Option (b) (parallel `ATM-NNN` + alias) has no schema column today | No — Option (a) (use `Gxx`/`Rxx` literally as `atm_id`) needs zero schema change and is the recommended default | Operator decision per G47; only pursue Gap-D column if operator picks (b) |
| **E** | `schema.sql` vs `schema_embed.sql` drift; `canonical_track`/`group_paths` (§11.4.191) absent from the runtime schema actually used | No — this project's core adoption does not need multi-track fields (single serialized lane today) | Report to constitution-submodule maintainers as its own finding; out of this project's fix scope |
| **F** | `reopens_count`/`Reopened-Details` (§11.4.34/§11.4.55) have no source data yet — no G-item has been formally "Reopened" in the register's own vocabulary, though G01's two-attempt cycle is reopen-like in substance | No — reopens_count starts at 0 for every imported item (accurate: none has been reopened under the FORMAL vocabulary); future reopens are captured exactly per §11.4.34 going forward | None needed — honest as-is |

---

## 5. Migration plan (phased, non-destructive)

**Phase 0 (this document).** Design only. No `.db` file created, no register
line touched, no generator wired. Done.

**Phase 1 — Import (additive, reversible, no register edit).**
1. Run `workable-items add <type> <severity> --db docs/workable_items.db --id
   <Gxx-or-Rxx> --title <T> --description <Evidence+Why-it-matters>
   --created-by Claude --assigned-to ''` once per G-item and per R-item
   selected for import, using the field-mapping + classification table from
   §3.2/§3.3 (author the classification table as a reviewable artifact in the
   migration commit, not silently inline).
2. For every item, follow with `workable-items update --id <id> --status
   <classified-status> --closure-criteria <DECISION+test-coverage-text>`.
3. For any item whose narrative shows a genuine attempt/redo cycle (G01), add
   the corresponding `item_history` rows via the engine's mutation path
   (`update`/`reopen` as appropriate) — reconstructing dates from the
   register's own dated STATUS labels (`STATUS (2026-07-15) — …`), never
   inventing a date (§11.4.6).
4. Run `workable-items validate --db docs/workable_items.db` — MUST be clean
   before proceeding (closed-set + §11.4.91 description-clarity invariants).
5. **The `.md` register is NOT touched or replaced in this phase.** This is
   the "additive, reversible" property: the DB is populated as a parallel
   representation while the hand-edited register remains the operative
   document the in-flight P0.5 Go-mutator lane's STATUS-block edits continue
   to land on.

**Phase 2 — Round-trip proof (the gate before the DB may be treated as
authoritative).**
1. Run `workable-items sync db-to-md --db docs/workable_items.db
   --out-issues <scratch-path>` and diff the generated output against the
   live register.
2. Because the register is under ACTIVE EDIT by the concurrent P0.5 lane, a
   single snapshot will not stay byte-identical for long — this phase
   requires **re-running the import (`sync md-to-db`, idempotent upsert) at
   each natural P0.5 pause point** and re-checking `diff` until it reports
   clean, rather than a single one-shot comparison. This is the honest
   consequence of adopting the DB mid-flight against a still-changing source
   document (§11.4.6 — do not claim round-trip proof against a moving
   target).
3. Only when `diff` reports empty (or only whitespace/section-order
   tolerance per the engine's documented round-trip contract) does Phase 2
   close.

**Phase 3 — Cutover (the DB becomes authoritative; generators become the
write path).**
1. Switch the register-regeneration path: future edits happen via
   `workable-items update`/`add`/`close`/`reopen`, then `workable-items
   export` regenerates `GAPS_AND_RISKS_REGISTER.md` (+ any `Summary.md`/HTML/
   PDF/DOCX siblings per §11.4.12/§11.4.53/§11.4.65) — direct hand-edits to
   the register's item bodies are retired (§11.4.93 "text-direct edits
   prohibited" migration-phase discipline), though the narrative
   preamble/adjudication/resolved-non-findings sections (captured as `raw`
   `doc_segments`) may still need occasional hand-touch until Phase 4 tooling
   exists for them.
2. WAL-checkpoint (`PRAGMA wal_checkpoint(TRUNCATE)`) + git-add + commit the
   `.db` file per §11.4.95 (never gitignored — it IS the source, not a build
   artifact).
3. Wire `workable-items validate` into this project's pre-build gate (once
   one exists) per §11.4.93's mandate.

**Phase 4 (optional, deferred).** `logic_groups`/`group_paths`/multi-track
population — only pursue if/when this project actually adopts §11.4.176
multi-track dispatch; not needed for the current single-lane operating mode.
Blocked today anyway by Gap E (the runtime schema lacks `group_paths`).

**Explicit sequencing relative to P0.5 (§11.4.101 — do not disrupt in-flight
work):** Phase 1 (DB population) is itself a **safe-during-active-Go-work**
activity per the constitution's own §11.4.96 SAFE catalogue — it explicitly
lists "workable-items DB ops per §11.4.93+§11.4.95" as SAFE to run in
parallel with a heavy build/mutator lane, because it touches only
`docs/workable_items.db` (a new, currently-nonexistent file) and never
`project/` Go source. It may therefore run as its own background
design/tooling stream concurrently with the P0.5 Go-mutator lane, exactly as
this design task itself was dispatched. **Phase 2's round-trip proof and
Phase 3's cutover, however, are sequenced to land at a natural P0.5 pause
point** (a batch boundary where the register's STATUS-block churn quiesces
briefly) so the round-trip diff is checked against a momentarily-stable
target rather than chasing a continuously-moving one — this is the "separate
gated step" the dispatch brief asked for.

---

## 6. Test plan

- **Unit** — schema round-trip: insert one of each Shape-A/Shape-B item kind,
  read back, assert field equality; `body_md` byte-equality against the
  original captured block.
- **Integration** — real SQLite (§11.4.27 no-fakes-beyond-unit): full Phase-1
  import of all 49 G-items + selected R-items into a real `.db`, followed by
  `sync db-to-md` → `diff` against the live register at a quiesced moment →
  MUST be empty (modulo the engine's documented whitespace/section-order
  tolerance).
- **Paired §1.1 mutation** — corrupt one field mapping (e.g. drop the
  `closure_criteria` population for one item) and assert the round-trip
  diff/`validate` gate FAILs; strip a `body_md` segment and assert `diff`
  detects the loss.
- **Stress** — the §11.4.93-mandated 1000-row insert + 10-concurrent-writer
  suite already ships with the engine itself (constitution-submodule scope,
  not re-authored here); this project's adoption inherits that proof by
  reference, it does not need its own 1000-row stress test for a 49-item
  register — cite the engine's own suite as the evidence source.
- **Chaos** — likewise inherited by reference: the engine's own mid-write
  SIGKILL + corrupt-DB-recovery + disk-full suite (§11.4.93) is the evidence
  for this project's adoption; this project does not re-implement it.
- **Anti-bluff (§11.4.6/§11.4.107(10))** — the migration's own classification
  table (Category→Type, STATUS-narrative→closed-set) ships as a reviewable
  artifact (not silently inline in a script) so a code-reviewer (per
  §11.4.209 Fable-xhigh) can independently verify every judgment call against
  the source register text.

---

## 7. Anti-bluff boundary — what "adopted" means vs. what remains (§11.4.6)

**"Adopted" (the terminal state this design aims at) means:** `docs/workable_items.db`
exists, is git-tracked (§11.4.95), contains every G01–G49 + selected R-item as
a validated row with round-trip-proven `body_md`, `workable-items validate`
runs clean, and `export`/`sync db-to-md` are wired as the register's
regeneration path (§11.4.12/§11.4.106) — i.e., Phases 1–3 above complete and
proven with captured evidence (the `diff` output, the `validate` exit code,
the commit that lands the `.db` file).

**What remains until the migration ACTUALLY runs (this document does not
perform it):** no `.db` file exists yet; no register line has moved; G40/G45
stay `STATUS: FILED` until Phase 3 closes them with cited evidence (the
`validate`-clean run + the round-trip `diff`); G47 stays `STATUS: FILED;
operator-decision PENDING` until the operator is asked (§11.4.66) which alias
strategy (a) or (b) to use — this design recommends (a) but does not decide
for the operator.

---

## Sources verified (§11.4.99)

No external web sources were consulted for this design — it is a pure
internal-codebase investigation (constitution submodule + project's own
docs). No operator-facing installation instructions are produced by this
document, so §11.4.99's cross-reference-against-latest-online-docs mandate
does not apply here.
