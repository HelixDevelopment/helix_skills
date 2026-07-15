# P1.T1 — Granularity Schema Migration Design (R16)

**Revision:** 1
**Last modified:** 2026-07-15T16:31:38Z
**Status:** design-research (no code landed by this doc). Drives migration `002_granularity` + its lockstep Go/SPEC/seed edits.
**Authority:** implements **P1.T1** (`IMPLEMENTATION_PLAN.md:170-177`, R16 at `REQUIREMENTS.md:191-196`); is the **hard prerequisite** for the G06/G07 DAG fix (`research/g06_g07_skilltree_dag_design.md` §3 + §5). Grounds every schema change in `research/skill_granularity_and_composition.md` §4/§5.
**Anti-bluff note (§11.4.6/§11.4.123):** every file:line cited below was read this session. No DDL is inferred from memory. The one runtime-breaking lockstep item (§2, the `ON CONFLICT` collision) is proven from `store.go:313` read against the PK-widening. **Honest gap (§5.3):** every DB-executed property (migrate up/down clean, CHECK/PK enforcement, `composes`-insert-now-succeeds, round-trip) requires a live containerized PostgreSQL 16 — it cannot be run in this design session; `graph_test.go:68-82` already documents this exact boundary.

---

## 0. Method & files read (file:line grounding)

Read in full this session: `project/migrations/001_initial.up.sql`, `project/migrations/001_initial.down.sql`, `project/internal/models/skill.go`, `project/internal/skill/graph.go`, `project/internal/skill/import_export.go`, relevant regions of `project/internal/skill/store.go`, `project/internal/db/migrations.go`, `project/internal/skill/graph_test.go`; `SPEC.md` §4/§5; `REQUIREMENTS.md` R16; `research/skill_granularity_and_composition.md` §3/§4/§5/§8; `research/g06_g07_skilltree_dag_design.md` §3/§5; `research/testing_infrastructure_plan.md` G06/G07 rows + infra section; the 8 `project/seed/skills/*.toml`.

**Environment facts (from source):** Postgres image `pgvector/pgvector:pg16` (`docker-compose.yml:15`; `testing_infrastructure_plan.md:38`) → **PostgreSQL 16**. Migration runner: `internal/db/migrations.go` — files named `NNN_description.up.sql` / `.down.sql`, zero-padded (`:22-23`), executed **in version order** and **each inside one transaction** (`runMigrationSQL`, `:291-326` — `tx.Exec(sql)` then record/delete `schema_migrations`, commit or rollback). Down migration supported (`MigrateDown`, `:84-145`). **Consequence:** the entire up (or down) migration is atomic — partial application is structurally impossible; a failed narrowing rolls the whole thing back.

**Current schema baseline (001):** `skills` has **no `kind`** column (`001_initial.up.sql:6-18`); `skill_dependencies` is `relation_type TEXT DEFAULT 'requires' CHECK (relation_type IN ('requires','extends','recommends'))` with `PRIMARY KEY (skill_id, depends_on)` and **no** `optional`/`sort_order` (`:25-30`). Mirror in `SPEC.md:227-232`. Model `DependencyType` = only the 3 values (`models/skill.go:24-28`; `SPEC.md:101-105`). `graph.go` `validTypes` map allows only those 3 and **rejects `composes`** (`graph.go:26-30`).

---

## 1. Migration DDL (`002_granularity.up.sql` / `.down.sql`)

Design target = the fresh-schema form in `skill_granularity_and_composition.md` §5.3 (`:237-256`) expressed as an **additive `ALTER` migration** on top of 001 (§5.3 `:253-258` mandates exactly this for a non-greenfield DB). All four changes required by the granularity model and the G06/G07 mismatch table (`g06_g07…:296-306`) land here.

### 1.1 Forward — `project/migrations/002_granularity.up.sql`

```sql
-- 002_granularity.up.sql — P1.T1 Granularity & composition schema (R16).
-- Additive, zero-data-loss superset of 001_initial. Runs in ONE transaction
-- (internal/db/migrations.go runMigrationSQL) — all-or-nothing.

-- (1) skills.kind — aggregation axis (default keeps every existing row valid;
--     backward-compat with the 8 seed TOMLs, which omit kind → 'atomic').
ALTER TABLE skills
    ADD COLUMN kind TEXT NOT NULL DEFAULT 'atomic'
    CHECK (kind IN ('atomic', 'composite', 'umbrella'));
CREATE INDEX idx_skills_kind ON skills(kind);

-- (2) widen the relation_type CHECK to the full typed-edge set (§4.1).
ALTER TABLE skill_dependencies
    DROP CONSTRAINT IF EXISTS skill_dependencies_relation_type_check;   -- Postgres default name for the inline 001 CHECK
ALTER TABLE skill_dependencies
    ADD CONSTRAINT skill_dependencies_relation_type_check
    CHECK (relation_type IN (
        'requires', 'extends', 'recommends',        -- existing (001)
        'composes', 'related_to', 'alternative_to'  -- NEW (R16 §4.1)
    ));

-- (3) edge attributes for umbrella→component ordering/optionality (§5.3).
ALTER TABLE skill_dependencies
    ADD COLUMN optional   BOOLEAN NOT NULL DEFAULT FALSE,   -- default keeps 001 rows valid
    ADD COLUMN sort_order INT;                              -- nullable; NULL = unordered

-- (4) widen the edge PK so a (skill_id, depends_on) pair may carry >1 typed edge (§5.3 :258).
--     relation_type must be NON-NULL to join a PK; 001 left it nullable-with-default.
UPDATE skill_dependencies SET relation_type = 'requires' WHERE relation_type IS NULL;  -- defensive backfill (seed has none)
ALTER TABLE skill_dependencies ALTER COLUMN relation_type SET NOT NULL;
ALTER TABLE skill_dependencies DROP CONSTRAINT skill_dependencies_pkey;                 -- Postgres default PK name
ALTER TABLE skill_dependencies
    ADD CONSTRAINT skill_dependencies_pkey
    PRIMARY KEY (skill_id, depends_on, relation_type);
```

Notes, each grounded:
- **Constraint names.** Postgres deterministically names an inline unnamed column CHECK `<table>_<column>_check` and an unnamed table PK `<table>_pkey`, so `skill_dependencies_relation_type_check` / `skill_dependencies_pkey` are correct for the 001 definitions (`001_initial.up.sql:28-29`). `DROP … IF EXISTS` on the CHECK is defensive; the PK drop uses the guaranteed `_pkey` name. **Verify once on live PG16 with `\d skill_dependencies` before merge** (the single name-assumption gap, §5.3).
- **Indexes.** Dropping the old PK removes the implicit unique index on `(skill_id, depends_on)`; the new triple PK supplies the unique index that `store.go`'s `ON CONFLICT (…, relation_type)` (post-fix, §2) resolves against. The explicit `idx_deps_skill(skill_id)` / `idx_deps_depends_on(depends_on)` (`001…:32-33`) are untouched and still serve the pair-existence and reverse lookups (`graph.go:53-55`, `:281-285`).
- **`kind` default = `atomic`** per §3.1 (`skill_granularity…:81`): the 8 seed TOMLs omit `kind`, so the column DEFAULT makes them load valid unchanged — the load-bearing backward-compat lever.

### 1.2 Down — `project/migrations/002_granularity.down.sql` (fail-closed inverse)

```sql
-- 002_granularity.down.sql — inverse of 002. FAIL-CLOSED: if the granularity
-- feature was used (a >1-edge-per-pair, or a composes/related_to/alternative_to
-- row), the narrowing steps ERROR and the whole tx rolls back — never a silent
-- partial drop. Safe pre-use rollback tool; a down AFTER real feature use is a
-- documented destructive rollback (see §3.4).

-- inverse (4): narrow PK back to the pair — ERRORS if any pair now has >1 edge.
ALTER TABLE skill_dependencies DROP CONSTRAINT skill_dependencies_pkey;
ALTER TABLE skill_dependencies
    ADD CONSTRAINT skill_dependencies_pkey PRIMARY KEY (skill_id, depends_on);
ALTER TABLE skill_dependencies ALTER COLUMN relation_type DROP NOT NULL;   -- restore 001 nullable-with-default

-- inverse (3): drop edge attributes.
ALTER TABLE skill_dependencies
    DROP COLUMN IF EXISTS sort_order,
    DROP COLUMN IF EXISTS optional;

-- inverse (2): narrow the CHECK back to the 001 3-value set — ERRORS (rolls back
--   the whole tx) if any composes/related_to/alternative_to row exists.
ALTER TABLE skill_dependencies
    DROP CONSTRAINT IF EXISTS skill_dependencies_relation_type_check;
ALTER TABLE skill_dependencies
    ADD CONSTRAINT skill_dependencies_relation_type_check
    CHECK (relation_type IN ('requires', 'extends', 'recommends'));

-- inverse (1): drop skills.kind (+ index).
DROP INDEX IF EXISTS idx_skills_kind;
ALTER TABLE skills DROP COLUMN IF EXISTS kind;
```

Down ordering is deliberate: the two narrowing statements that can hold real data (PK→pair, CHECK→3-value) run **before** the destructive `DROP COLUMN`s, so a feature-using DB aborts the entire transaction at the first narrowing that would violate — the rollback is atomic (`migrations.go:304-306`). See §3.4 for the exact data-loss boundary.

---

## 2. Lockstep change list (§11.4.186 — SPEC + migration + models + code MUST agree)

Every representation of the edge/kind contract must move together or the system diverges. The table below is **exhaustive** — every file that must change, with file:line and the reason. **Column "Class":** `BREAKING` = the app breaks at runtime/compile if the migration lands without this edit (must ship in the same change); `FEATURE` = required for the granularity feature to actually work / for model↔DB to agree, non-breaking if briefly deferred but then a §11.4.186 divergence; `DOC` = contract document.

| # | File:line | Change | Class | Grounding |
|---|---|---|---|---|
| L1 | `project/migrations/002_granularity.up.sql` (new) | the §1.1 DDL | — | this doc §1 |
| L2 | `project/migrations/002_granularity.down.sql` (new) | the §1.2 DDL | — | this doc §1.2 |
| **L3** | `internal/skill/store.go:311-313` | `ON CONFLICT (skill_id, depends_on)` → `ON CONFLICT (skill_id, depends_on, relation_type)` in the dep upsert | **BREAKING** | PK widening (§1.1 step 4) removes the `(skill_id, depends_on)` unique constraint; the current `ON CONFLICT (skill_id, depends_on)` then has **no matching unique index** → `Create()` errors at runtime on any skill-with-deps. This is the single hard runtime coupling — **it MUST ship with the migration.** Read at `store.go:313`. |
| L4 | `internal/models/skill.go:11-19,21-28,50-59` | add `SkillKind` type + `atomic/composite/umbrella` consts; add `Skill.Kind SkillKind` field (`db:"kind" toml:"kind"`, default `atomic`); add `DepTypeComposes/RelatedTo/Alternative` to the `DependencyType` const block; add `HardClosureTypes` var `{requires,composes,extends}`; add `Optional bool`/`SortOrder *int` to `SkillDependency` | FEATURE | `skill_granularity…:196-230` (§5.2); mismatch table `g06_g07…:298-302` |
| L5 | `internal/models/skill.go:128-147` | `TOMLSkillDef` gains `Kind string toml:"kind"`; `TOMLDependencies` gains `Composes/RelatedTo/Alternative []string` + alias-accepting fields (`DependsOn/Prerequisite/PartOf`); add optional `[[skill.components]]` array-of-tables type (`Name/Order/Optional/Note`) | FEATURE | `skill_granularity…:144-191` (§5.1); TOML authoring form |
| **L6** | `internal/skill/graph.go:26-30` | widen `validTypes` map: add `DepTypeComposes`, `DepTypeRelatedTo`, `DepTypeAlternative` (stop rejecting `composes`) | **BREAKING** (for the feature) | `graph.go:31-33` currently returns `ErrInvalidSkill` for any type outside the 3 — so `AddDependency` would reject every new edge type the migration now permits at the DB layer. Read at `graph.go:26-33`. G07's resolve-and-persist routes through `AddDependency` (`g06_g07…:256-258`), so this must precede it. |
| L7 | `internal/skill/store.go` SELECT+Scan pairs: `GetByName` `:123-124`/`:130-134`; `GetTree`/`buildTree` `:219-…`; `Search` `:55-56`/`:89-90`; `VectorSearch` `:577-578`; `ListSkills` `:612-613` | add `s.kind` (or `kind`) to each SELECT column list **and** `&skill.Kind` to the paired `Scan` | FEATURE | each SELECT that gains a column MUST gain the matching Scan target (else pgx "field count" error) — the classic §11.4.186 SELECT↔Scan↔struct triad. Column lists read at cited lines. |
| L8 | `internal/skill/graph.go` SELECT+Scan: `GetDependencyTree` CTE `:194-236`/`:254-259`; `GetDependents` `:332`; `GetAllDependencies` `:360`; `scanSkills` `:374-388` | add `s.kind` + `&sk.Kind` | FEATURE | same triad. NB the `GetDependencyTree` CTE column list is also rewritten by **G06** (`g06_g07…:104-128` selects `s.kind`, `sd.optional`, `sd.sort_order`) — that rewrite depends on these columns existing, hence P1.T1 first. |
| L9 | `internal/skill/store.go:286-287` (Create INSERT) and `internal/skill/import_export.go:150-153` (ImportFromTOML INSERT) | add `kind` to the column list + value, coercing empty→`atomic` (or omit `kind` to inherit the DB default) | FEATURE | authored `umbrella`/`composite` skills must persist their kind; a Go `Kind=""` passed literally would violate the CHECK — coerce or omit. Both INSERT statements read at cited lines. |
| L10 | `internal/skill/store.go:308-319` (Create dep loop) and `import_export.go:178-181` (dep INSERT) | add `optional,sort_order` columns to the edge INSERTs (default-safe if omitted) | FEATURE | edge attrs are unset today; persisting `[[skill.components]]` order/optional needs them. **The `ON CONFLICT` fix in L3 lives in this same `store.go` block.** |
| L11 | `internal/skill/import_export.go` — `collectDepNames` `:304-328`, `depsToCreate` builder `:101-134`, `ExportToTOML` switch `:270-277` | recognise `composes/related_to/alternative_to` + aliases (`depends_on/prerequisite`→`requires`, `part_of`→inverted `composes`); export the new relation types | FEATURE — **largely G07 scope** | round-trip normalization is owned by G07 (`g06_g07…:240-263`); P1.T1 only needs the *fields to exist* (L5). Flag the boundary so it is not double-owned. |
| L12 | `SPEC.md:99-112` | add the 3 new `DepTypeComposes/RelatedTo/Alternative` consts + `SkillKind` + `Skill.Kind` + `SkillDependency.Optional/SortOrder` to the §4.1 model block | DOC | `SPEC.md:99-112` currently shows only the 3-type enum |
| L13 | `SPEC.md:220,227-232` | add `kind TEXT NOT NULL DEFAULT 'atomic' CHECK(...)` to the `skills` DDL; widen the `relation_type` CHECK to 6 values; add `optional`/`sort_order`; widen PK to the triple | DOC | `SPEC.md:230` still shows the 3-type CHECK + pair PK — the exact drift `g06_g07…:394-396` flags |
| L14 | `SPEC.md:183-186` (TOML example) | add `composes`/`related_to`/`alternative_to` keys + a `kind` line + a `[[skill.components]]` example | DOC | keeps the on-disk contract aligned with L5 |
| L15 | `project/seed/validate_dag.py` | scope the cycle check to `{requires,composes,extends}`; resolve `composes`/`part_of` targets; assert `atomic ⇒ 0 composes` and `umbrella/composite ⇒ ≥1 composes` (§5.4 invariants) | FEATURE | `skill_granularity…:134-136,260-266` (§4.3/§5.4); relaxation-only for the current corpus |
| L16 | `project/seed/skills/*.toml` ×8 | **no edit** — all 8 acquire `kind='atomic'` from the column DEFAULT (§6 table) | FEATURE | backward-compat lever (`skill_granularity…:81`, §8 `:492`) |

**Divergence-guard summary:** the SELECT↔Scan↔`models.Skill.Kind` triad (L4+L7+L8) and the CHECK↔`DependencyType`↔`validTypes`↔TOML-tags chain (L4+L6+L5+L13+L12) are the two axes a `CM-DOC-INTEGRITY-VALIDATION`-class check (§11.4.186) should assert agree after this change.

---

## 3. Migration safety

### 3.1 Forward = zero data loss (proven additive)

Every forward step is additive or a strict widening:
- `ADD COLUMN kind … DEFAULT 'atomic'` — existing `skills` rows materialise `kind='atomic'`; no row rewritten in a lossy way. (PG16 stores the default in catalog metadata for a `NOT NULL DEFAULT` add — fast, no full table rewrite of existing content.)
- CHECK widening — strictly more permissive; every 001-legal `relation_type` stays legal.
- `ADD COLUMN optional NOT NULL DEFAULT FALSE`, `ADD COLUMN sort_order INT` (nullable) — existing edges get `optional=FALSE`, `sort_order=NULL`.
- `relation_type` backfill + `SET NOT NULL` — the backfill only touches NULLs (the seed corpus has none; all `INSERT`s supply `relation_type` — `import_export.go:178-181`, `store.go:311-313`), so it changes nothing in practice and cannot lose a value.

### 3.2 How existing rows acquire defaults

`kind='atomic'`, `optional=FALSE`, `sort_order=NULL` come from the column DEFAULTs at `ADD COLUMN` time — no data-migration step, no per-row UPDATE (except the defensive NULL→`requires` backfill in §1.1 step 4, which is a no-op on real data). This is exactly the "default keeps existing rows valid" guarantee of `skill_granularity…:240,254`.

### 3.3 PK widening effect on existing edges

Under the 001 PK `(skill_id, depends_on)`, **no two rows can share a pair** — so widening to `(skill_id, depends_on, relation_type)` is a **strict superset key**: every existing row is trivially still unique on the triple. **Zero dedup is required** (contrast `skill_granularity…:494`, which flags dedup only for a hypothetical pre-existing exact-duplicate — impossible under the old PK; the seed has none). Post-widening the schema newly *permits* two typed edges on one pair (e.g. a `requires` **and** a `recommends` `A→B`); that is the intended capability (`g06_g07…:302`).

### 3.4 Down = fail-closed, with a documented data-loss boundary

The down migration (§1.2) reverts to the exact 001 shape **when the feature is unused**. When the feature **has been used**, the narrowing statements fail-closed inside the single transaction (`migrations.go:291-326`), rolling back the whole down — no silent partial:
- a `(skill_id, depends_on)` pair carrying >1 edge → the PK re-narrowing raises a duplicate-key error → full rollback;
- any `composes`/`related_to`/`alternative_to` row → the CHECK re-narrowing raises a violation → full rollback.
**Documented boundary:** the only unguarded loss is `DROP COLUMN kind` discarding a non-`atomic` kind on a node that has **no** `composes` edge — a state `validate_dag.py` §5.4 invariant #3 already forbids (an `umbrella`/`composite` must have ≥1 `composes` edge), so it should not exist in a valid corpus. A down after genuine feature use is therefore a deliberate, destructive rollback; the normal path is forward-only.

### 3.5 Atomicity

Both up and down execute inside one `tx` (`migrations.go:298-324`); a failure at any statement rolls back the entire migration and does **not** record/advance `schema_migrations` (`:310-320`). No half-migrated state is reachable.

---

## 4. Test design

All DB-touching cases are **integration** tests gated on a live/containerized `pgvector/pgvector:pg16` (`docker-compose.yml:15`; `testing_infrastructure_plan.md:38,191` type-2 "APPLICABLE-NOW", real pgx, no DB mock). This reconciles with the plan's G06 (`:292`) and G07 (`:293`) rows (which already require "migrate up → insert typed edges → read back" against real pgvector) and with the pre-existing honest boundary at `graph_test.go:68-82` (CTE/edge behaviour is un-unit-testable without a real Postgres). Each PASS cites captured evidence per §11.4.5/§11.4.69.

### 4.1 Migration-specific cases (this migration's own acceptance)

| # | Case | Assertion | Captured evidence a PASS cites |
|---|---|---|---|
| M1 | `migrate up` on fresh `pgvector:pg16` applies 001+002 clean | `schema_migrations` has version 2; no error | migrate-up log; `SELECT version FROM schema_migrations` dump |
| M2 | schema shape post-002 | `skills` has `kind` col + CHECK + `idx_skills_kind`; `skill_dependencies` has 6-value CHECK, `optional`, `sort_order`, triple PK | `\d+ skills` + `\d+ skill_dependencies` dumps |
| M3 | 001-era rows survive 002 | seed loaded pre-002 → post-002 every row has `kind='atomic'`, `optional=FALSE`, `sort_order IS NULL`; row count unchanged | `SELECT kind,count(*)` + before/after row-count dump |
| M4 | **`composes` edge now inserts (was rejected)** | `INSERT … relation_type='composes'` **succeeds** on 002; the identical INSERT **errors** on 001 (RED baseline, §11.4.115) | psql success on 002 + psql CHECK-violation error text on 001 |
| M5 | PK widening — two edges per pair | `requires` **and** `recommends` on the same `(A,B)` both persist | `SELECT * FROM skill_dependencies WHERE skill_id=A AND depends_on=B` dump (2 rows) |
| M6 | Kind/attr round-trip preserves values | insert skill `kind='umbrella'` + edge `optional=TRUE,sort_order=3`; `GetByName`/export reads them back identical | round-trip struct compare (`kind`,`optional`,`sort_order` equal) |
| M7 | CHECK rejects bogus values | `kind='bogus'` and `relation_type='bogus'` both rejected | two CHECK-violation error strings |
| M8 | `migrate down` clean on unused-feature DB | 002 down reverts to exact 001 shape; `\d` == 001 | `\d+` diff (== 001) + `schema_migrations` back to version 1 |
| M9 | **down fail-closed on used-feature DB** | with a `composes` row present, 002 down **errors + rolls back** (schema unchanged, no partial) | psql error text + post-abort `\d+` showing 002 shape intact |
| M10 | 8 seed TOMLs still parse+load (evidence gate) | all 8 import as `kind='atomic'`; `validate_dag.py` green | import log + `validate_dag.py` PASS + `SELECT name,kind` (8× atomic) |
| Mμ1 | **mutation** — strip `'composes'` from the 002 up CHECK | M4 (composes-insert) FAILs | mutation-run FAIL log |
| Mμ2 | **mutation** — leave PK at `(skill_id, depends_on)` | M5 (two-edges-per-pair) FAILs on insert; and L3's `ON CONFLICT` mismatch surfaces | mutation-run FAIL log |
| Mμ3 | **mutation** — revert L3 `ON CONFLICT` to the pair after widening PK | `Create()` with deps errors at runtime | mutation-run FAIL log (no-matching-unique-constraint) |

**Migration-suite count: 10 positive cases (M1–M10) + 3 mutation assertions (Mμ1–Mμ3) = 13 obligations.**

### 4.2 Reconciliation with the G06/G07 suite

The G06/G07 test suite (`g06_g07…:328-362`: **17 positive + 6 mutation = 23 obligations**) is **downstream** — it exercises the *behaviour* this schema enables (hard-closure CTE traversal, `composes`/attr round-trip). P1.T1's 13 obligations are the **schema-layer floor** those 23 build on: G06 case 3 (transitive `requires`/`composes`/`extends`) and G07 case 5 (attrs preserved) are un-runnable until M2/M4/M6 pass. No overlap, no double-counting — P1.T1 = migration acceptance; G06/G07 = consumer acceptance.

---

## 5. Sequencing (§11.4.197) + honest gaps

### 5.1 Order (each step's precondition is the prior step)

1. **P1.T1 (this migration) lands FIRST** — `002_granularity.up/down.sql` **plus the lockstep edits that keep the app compiling+running with zero regression**: the BREAKING items L3 (`ON CONFLICT` triple — must ship with the PK widening) and L6 (`validTypes` widened), the model/SELECT/INSERT triad (L4,L5,L7,L8,L9,L10), the SPEC/validator/seed edits (L12–L16). *Evidence gate (mirrors `IMPLEMENTATION_PLAN.md:175`):* 8 seeds still parse+load as `atomic`; `002` up/down clean on PG16 (M1,M8); `go build ./...` green; existing tests green.
2. **G06 — recursive-CTE hard-closure fix** consumes the new columns: its rewritten `GetDependencyTree` selects `s.kind`, `sd.optional`, `sd.sort_order` and filters `relation_type IN ('requires','composes','extends')` (`g06_g07…:104-128`). **Cannot land before P1.T1** — the columns/types must exist first (`g06_g07…:307-313`).
3. **G07 — TOML/JSON round-trip fix** consumes the new `TOMLDependencies` fields (L5) and resolves-and-persists through the **already-widened** `AddDependency` `validTypes` (L6) (`g06_g07…:240-263`). Depends on P1.T1.
4. **Granularity hard-closure resolver + the §6 `android` reclassification** (LAST): author the 6 component leaves, flip `android.overview`→`umbrella` + `android.aosp.build_system`→`composite`, add their `composes` edges. This needs (a) the schema, (b) `validTypes` widened, (c) G06's `composes`-traversal, **and** (d) satisfying `validate_dag.py` §5.4 invariant #3 (`umbrella/composite ⇒ ≥1 composes edge`). **This is why all 8 seeds stay `atomic` in P1.T1** — reclassifying before the `composes` edges + leaves exist would fail invariant #3. Reclassification is an **additive** later migration (`skill_granularity…:81,496`), never part of the schema migration.

### 5.2 Seed `kind` assignment (P1.T1 scope) — all 8 → `atomic`

Per §6 (`skill_granularity…:81`) and the P1.T1 backward-compat gate: **no seed TOML is edited**; each acquires `kind='atomic'` from the column DEFAULT.

| Seed file | `name` | P1.T1 `kind` | Future reclassification (deferred to step 4) |
|---|---|---|---|
| `java.toml` | `java.language` | `atomic` | stays `atomic` |
| `kotlin.toml` | `kotlin.language` | `atomic` | stays `atomic` |
| `cpp.toml` | `cpp.language` | `atomic` | stays `atomic` |
| `linux.toml` | `linux.os` | `atomic` | stays `atomic` |
| `make.toml` | `make.build_system` | `atomic` | stays `atomic` |
| `cmake.toml` | `cmake.build_system` | `atomic` | stays `atomic` |
| `android.toml` | `android.overview` | `atomic` | → `umbrella` **only when** its 6 `composes` leaves land (§6) |
| `android_aosp.toml` | `android.aosp.build_system` | `atomic` | → `composite` **only when** it gains ≥1 `composes` edge (§6; note it currently uses `extends android.overview` — an `extends` edge does **not** satisfy invariant #3, so it stays `atomic` until it composes real parts) |

The two reclassification candidates are named here explicitly so the requirement is **tracked, not lost** (§11.4.197): they are open work items whose acceptance is "component leaves authored + `composes` edges added + `validate_dag.py` §5.4 green", not a silent P1.T1 default.

### 5.3 Honest gaps (§11.4.6) — needs live PG16 to PROVE

- **Decidable now (from source, no DB):** the exact DDL (§1), the four required changes vs current schema (`g06_g07…:296-306`), the full lockstep file:line list (§2), and the one runtime-breaking coupling — `store.go:313` `ON CONFLICT` vs PK widening (§2 L3).
- **Requires live containerized PG16 to prove (integration layer, §11.4.108):** that `002` up applies clean and yields the intended `\d` shape (M1,M2); that a `composes` INSERT now succeeds while the same INSERT errored on 001 (M4 — the RED→GREEN flip); that the widened PK admits two edges per pair (M5); that kind/optional/sort_order round-trip (M6); that down is fail-closed on a used-feature DB (M9); that all 8 seeds still load (M10). These are DB-executed — unrunnable in this design session, consistent with the pre-existing skip at `graph_test.go:68-82`.
- **One name-assumption to verify on live before merge:** the default constraint names `skill_dependencies_relation_type_check` / `skill_dependencies_pkey` (§1.1) — confirm with `\d skill_dependencies` on a 001 DB; if a prior migration renamed them, adjust the `DROP CONSTRAINT` targets. Deterministic from Postgres naming rules, but not runnable here.
- **SPEC drift is REAL and must land in lockstep** (`g06_g07…:394-396`): `SPEC.md:230` + `migrations/001…:28` still show the 3-type CHECK; L12–L14 close it in the same change so contract↔schema↔code agree (§11.4.186).

---

## 6. Deliverable summary

- **Migration DDL:** `002_granularity.up.sql` (§1.1) — add `skills.kind` (`CHECK atomic|composite|umbrella`, default `atomic`) + index; widen `skill_dependencies.relation_type` CHECK to the 6-value typed set; add `optional BOOLEAN NOT NULL DEFAULT FALSE` + `sort_order INT`; backfill+`SET NOT NULL` on `relation_type` then widen PK to `(skill_id, depends_on, relation_type)`. `002_granularity.down.sql` (§1.2) — fail-closed inverse.
- **Lockstep:** §2 table L1–L16 — the one **BREAKING** runtime coupling is `store.go:313` `ON CONFLICT` (must ship with the PK widening); L6 `validTypes` widening unblocks `AddDependency`; the rest are the SELECT↔Scan↔model triad, INSERT coercion, SPEC §4/§5, `validate_dag.py`, and the (no-edit) seeds.
- **Backward-compat:** all-additive forward, zero data loss; the `kind` DEFAULT makes the 8 seeds load unchanged; PK widening is a strict-superset key → no dedup needed.
- **Test count:** migration suite = **13** (10 positive M1–M10 + 3 mutation); downstream G06/G07 = 23 (unchanged, built on this floor).
- **Honest gap:** every DB-executed property needs a live `pgvector:pg16` container; constraint-name assumption verified on live via `\d` before merge.
- **File written:** `research/p1t1_granularity_schema_migration.md`.
