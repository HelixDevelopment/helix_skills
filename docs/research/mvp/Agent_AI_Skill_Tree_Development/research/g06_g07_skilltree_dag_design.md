# G06 / G07 — Recursive Skill-Dependency DAG Traversal + Round-Trip Integrity Design

**Revision:** 1
**Last modified:** 2026-07-15T15:56:54Z
**Status:** design-research (drives a later Go fix — no code landed by this doc)
**Authority:** remediates `GAPS_AND_RISKS_REGISTER.md` G06 + G07; reconciles with
`research/skill_granularity_and_composition.md` §4/§5 (R16 / P1.T1) and
`research/testing_infrastructure_plan.md` G06/G07 rows.
**Anti-bluff note (§11.4.6/§11.4.123):** every root-cause claim below cites a real
file:line read this session. The G07 serialization root cause **and** its fix are
proven by a captured BurntSushi/toml v1.6.0 (the project's pinned version) decode/
encode experiment — transcript in
`scratchpad/g07_burntsushi_evidence.txt`. Nothing here is inferred from memory.

---

## 0. Method & files read

Read in full: `internal/skill/graph.go`, `internal/skill/store.go`,
`internal/skill/import_export.go`, `internal/skill/resources.go`,
`internal/skill/graph_test.go`, `internal/models/skill.go`,
`migrations/001_initial.up.sql`; the relevant regions of
`internal/api/skills_handler.go`, `internal/api/server.go`,
`internal/mcp/tools.go`; `SPEC.md` §4/§5/§6, `research/skill_granularity_and_composition.md`
§4/§5, `research/testing_infrastructure_plan.md` G06/G07 rows, the 8 `seed/skills/*.toml`.
Environment facts confirmed from source: Postgres image `pgvector/pgvector:pg16`
(`project/docker-compose.yml:15`) → PostgreSQL **16**; deps `jackc/pgx/v5 v5.9.2`,
`BurntSushi/toml v1.6.0` (`project/go.mod`).

### 0.1 The tree-function landscape (load-bearing — the register is imprecise here)

There are **three** rival dependency-tree code paths; the register attributes G06 to
`GetDependencyTree`, but that function is **not called anywhere** (verified:
`grep -rn GetDependencyTree` → only its own definition). Reality:

| Path | Function | Wired to | Correctness (as read) |
|---|---|---|---|
| REST (live server) | `store.GetTree` | `cmd/server/main.go:230` | Go recursion + shared `visited` map — recurses, but N+1 queries and one shared visited-set across siblings ⇒ a diamond's shared node attaches under its **first** parent only (DAG-flatten) |
| MCP `skill_tree` | `store.GetTree` | `internal/mcp/tools.go:215` (`serializeTreeNode` recurses over `node.Children`, `tools.go:226-246`) | same as above — **does recurse**; the register's claim that MCP "emits a 1-level tree" (`GAPS…:100`) is **incorrect**, grounded in `tools.go:215` |
| REST (hardened) | `Pool.GetSkillTree` interface | `internal/api/server.go:54`, handler `skills_handler.go:316` | **no concrete implementation exists** (`grep -rn 'func.*GetSkillTree'` → none); the hardened `internal/api` server is also un-wired to `main` (no `api.New(` caller in `cmd`) — the dead second server of G01/G09 |
| **canonical, intended** | `graph.go GetDependencyTree` (CTE) | **nothing** — dead/rival copy | **broken (depth-1)** — see §1 |

**Consequence for the fix:** G06 is not "one buggy function". It is (a) the broken
canonical CTE function `GetDependencyTree`, plus (b) two rival tree impls (`store.GetTree`
correct-ish, `Pool.GetSkillTree` missing) that must be **consolidated onto one canonical
recursive-CTE function wired to every surface** (§11.4.124 dead-code / §11.4.108
rival-copy). Shipping a fix to `GetDependencyTree` alone would leave it still un-wired.

---

## 1. G06 — recursive tree truncated to depth-1

### 1.1 Root cause (exact)

`GetDependencyTree` (`graph.go:179-310`) has a **correct** recursive CTE and a **broken
Go assembly**:

- **CTE is correct** (`graph.go:194-236`): the `WITH RECURSIVE dep_tree` walks the full
  transitive closure and carries `depth`; it returns every transitive dependency as a flat
  row set ordered by `(depth, name)`. This part is not the bug.
- **Assembly truncates** (`graph.go:271-307`). It builds `childrenMap` keyed by *every*
  parent id from a re-query of `skill_dependencies` (`graph.go:281-303`), i.e. the full
  parent→children map is computed — then attaches **only the root's entry**:

  ```go
  // graph.go:306
  root.Children = childrenMap[rootSkill.ID]
  ```

  Grandchildren are never linked: the child `SkillTreeNode`s built at `graph.go:298-301`
  are created with `Skill`+`Depth` but their `Children` field is left nil, and nothing ever
  walks `childrenMap` for those non-root parent ids. Net result: root + direct children, no
  deeper levels — **depth-1**. (Contrast `GetAllDependencies`, `graph.go:347-371`, which is
  correct because it returns a *flat* closure and needs no assembly.)

- **Two latent secondary defects in the same CTE** that the fix must also close:
  1. **No typed-edge filter.** The CTE joins `skill_dependencies` unconditionally
     (`graph.go:209-211, 228-230`) — it walks `recommends` edges as if they were
     prerequisites. Per the granularity model the hard closure is `{requires, composes,
     extends}` only (`skill_granularity…:126`); `recommends`/`related_to`/`alternative_to`
     are advisory and **not transitive** (§3).
  2. **No visited-set; cycle-safety by depth cap only.** It uses `UNION ALL` with a single
     `WHERE dt.depth + 1 < $2` guard (`graph.go:213, 231`). On a cycle it re-emits the loop
     each level until the depth cap (terminates but produces garbage duplicate rows); on a
     diamond it emits the shared node once per path (row blow-up). A real skill DAG can
     contain a cycle via bad data or the symmetric `related_to`/`alternative_to` edges (§3),
     so a proper visited-set is required, not just a bound.

### 1.2 Decision — one canonical, cycle-safe, depth-bounded, typed-edge-filtered CTE

Replace all three rival paths with a single `Store.GetDependencyTree(ctx, rootName, maxDepth)`
(keep the name; retire `store.GetTree` and implement `Pool.GetSkillTree` as a thin delegate,
or route the hardened handler to the store). The traversal is a PostgreSQL **RECURSIVE CTE**
that (a) filters to the hard-closure edge set, (b) carries a **path array as the visited-set**
for cycle safety, (c) bounds depth, and (d) carries `parent_id`, `relation_type`, and the new
`optional`/`sort_order` attributes so the Go side can assemble and order the tree in **one
pass** (kills the N+1 re-query).

**Cycle-safety mechanism — two equivalent forms (PG16 supports both):**

*Form A — SQL-standard `CYCLE` clause (preferred; PostgreSQL ≥ 14, we run 16).* The engine
tracks a path of a chosen column and stops recursing when a row would revisit it:

```sql
WITH RECURSIVE dep_tree AS (
    -- base: hard-closure edges out of the root
    SELECT sd.skill_id AS parent_id, sd.depends_on AS id,
           sd.relation_type, sd.optional, sd.sort_order, 1 AS depth
    FROM skill_dependencies sd
    WHERE sd.skill_id = $1
      AND sd.relation_type IN ('requires','composes','extends')   -- hard closure only
    UNION ALL
    SELECT sd.skill_id, sd.depends_on,
           sd.relation_type, sd.optional, sd.sort_order, dt.depth + 1
    FROM skill_dependencies sd
    JOIN dep_tree dt ON dt.id = sd.skill_id
    WHERE sd.relation_type IN ('requires','composes','extends')
      AND dt.depth < $2                                            -- depth bound
)
CYCLE id SET is_cycle USING path                                   -- stop on revisit
SELECT dt.parent_id, dt.id, dt.relation_type, dt.optional, dt.sort_order,
       dt.depth, dt.is_cycle,
       s.name, s.version, s.title, s.description, s.content, s.metadata,
       s.kind, s.status, s.created_at, s.updated_at
FROM dep_tree dt
JOIN skills s ON s.id = dt.id
WHERE NOT dt.is_cycle                                              -- drop the closing edge
ORDER BY dt.depth, dt.sort_order NULLS LAST, s.name;
```

*Form B — portable path-array (works on any PG; keep as the fallback / if `CYCLE` is
undesirable):* carry `path uuid[]` and guard `AND NOT sd.depends_on = ANY(dt.path)`, seeding
`ARRAY[sd.skill_id, sd.depends_on]` in the base term and `dt.path || sd.depends_on` in the
recursive term. The `= ANY(path)` predicate is the visited-set; it prevents revisiting a node
on any single path and terminates every cycle.

Decision: **use Form A (`CYCLE`) on PG16** — it is declarative, the engine-native visited-set,
and self-documents intent; keep Form B documented as the portability fallback.

**Reachable-node de-duplication vs. path duplication (the semantics decision):** a DAG rendered
as a tree can either (i) **duplicate** a shared node under every parent (pure tree; can be
exponential on deep diamonds even when cycle-free) or (ii) **emit each reachable node once**
(DAG-as-tree). The register/testing-plan property is **"tree node count == closure size"**
(`GAPS…:103`, `testing_infrastructure_plan.md:292`) — that equality only holds under (ii).
**Decision: emit each reachable node once**, attached under its *first-resolved* (minimum-depth)
parent, matching `GetAllDependencies`'s closure and the existing `store.GetTree` shared-visited
behavior. This makes `count(nodes) == |hard-closure(root)|` a true invariant and bounds work at
O(closure) rather than O(paths). (`sort_order` still orders siblings; the CTE's per-path `CYCLE`
guard is for termination, the Go-side global visited-set is for dedup — both are needed.)

**Go one-pass assembly (replaces `graph.go:271-307`):**

```
rows := <run CTE above>
nodeByID  := map[uuid.UUID]*SkillTreeNode{ root.ID: root }   // root pre-seeded, depth 0
childRows := []row{}                                          // buffer, we need parents first
for each row (already ordered by depth asc):
    if _, seen := nodeByID[row.id]; seen { continue }         // global dedup ⇒ node-count==closure
    n := &SkillTreeNode{ Skill: row.skill, Depth: row.depth } // relation/optional/order also carried
    nodeByID[row.id] = n
    parent := nodeByID[row.parent_id]                         // guaranteed present: parents come first (depth asc)
    parent.Children = append(parent.Children, *n)
// then sort each node.Children by (sort_order NULLS LAST, name) if not already ordered in SQL
```

Because rows arrive depth-ascending, a node's parent is always materialized before the child —
no second pass, no re-query. This structurally cannot truncate (every reachable node is attached
to a real parent), is cycle-safe (CTE `CYCLE` + Go `seen` map), depth-bounded (`dt.depth < $2`),
and typed-edge-filtered.

**Depth bound:** keep the existing clamp (`graph.go:180-185`: default 10, hard cap 50) — it is a
sound guard; the fix does not remove it.

**Sources (§11.4.8):** PostgreSQL 16 manual, *7.8 WITH Queries (CTEs)* — recursive-query form,
`UNION`/`UNION ALL` semantics, and the `SEARCH`/`CYCLE` clauses (`CYCLE col SET flag USING path`)
introduced in PG 14: https://www.postgresql.org/docs/16/queries-with.html . `ORDER BY` on the
outer query for depth-first presentation is documented there.

---

## 2. G07 — TOML/JSON dependency + resource round-trip broken

### 2.1 Root causes — THREE independent defects (one newly proven with captured evidence)

**(A) TOML struct nesting/tag mismatch — silent drop in EVERY import path (new, captured).**
`models.TOMLSkillWrapper` (`models.go:122-126`) models the nested TOML tables
`[skill.dependencies]` and `[[skill.resources]]` as **top-level sibling fields** with dotted
tags:

```go
type TOMLSkillWrapper struct {
    Skill        TOMLSkillDef     `toml:"skill"`
    Dependencies TOMLDependencies `toml:"skill.dependencies"`   // models.go:124
    Resources    []TOMLResource   `toml:"skill.resources"`      // models.go:125
}
```

BurntSushi/toml does **not** interpret a dotted struct tag as a nested-table path — it is a
single literal key. The on-disk seed (`seed/skills/android.toml`) writes
`[skill.dependencies]` / `[[skill.resources]]` as *sub-tables of `[skill]`*, which have no
landing field. **Captured proof** (replica struct + real seed, BurntSushi v1.6.0, transcript in
`scratchpad/g07_burntsushi_evidence.txt`):

```
Dependencies.Requires = []  (len=0)          <- silently empty
Resources             = []  (len=0)          <- silently empty
Undecoded keys        = [skill.dependencies skill.dependencies.requires ... skill.resources ...]
(toml.Decode returned NO error)
Re-encode emits:  ["skill.dependencies"]     <- quoted-literal TOP-LEVEL table (not nested)
```

The loss is **100% silent** because `ImportFromTOML` uses `toml.Unmarshal`
(`import_export.go:24`), which discards the `MetaData`, so `md.Undecoded()` is never checked.
**This defeats the "correct" resolver too:** `ImportFromTOML` (`import_export.go:70-196`) *does*
resolve names→IDs and insert edges/resources properly — but it is handed an already-empty
`wrapper.Dependencies`/`wrapper.Resources`, so it inserts nothing. **This closes the register's
UNCONFIRMED note** (`GAPS…:112`): the MCP `skill_create` path (`tools.go:275` →
`ImportFromTOML`) drops edges *upstream of* its resolver. The encoder's quoted-literal
`["skill.dependencies"]` also means **export is structurally wrong** independent of handler
logic.

**(B) REST `convertTOMLWrapper` throws away the edge target (existing-bug, decidable).**
`skills_handler.go:548-565` builds `SkillDependency{RelationType: …}` but **never sets
`DependsOn` or `DependsOnName`** — literally `_ = depName // placeholder for resolution`
(`skills_handler.go:552,558,564`). Even if (A) were fixed, this path would persist edges with a
zero-UUID target (FK violation) or a nameless edge.

**(C) REST `handleCreateSkill` drops deps/resources entirely (existing-bug, decidable).**
`skills_handler.go:164-176` constructs the `models.Skill` with scalar fields only, then calls
`s.pool.CreateSkill` (`skills_handler.go:179`); `req.Deps` / `req.Resources`
(`CreateDepsRequest`/`CreateResourceReq`, `skills_handler.go:33-45`) are parsed and ignored.

**Export side (works — one direction only):** `ExportToTOML` (`import_export.go:233-301`) and
`exportToTOMLWrapper` (`skills_handler.go:590-629`) *do* read `dep.DependsOnName`
(`import_export.go:261-268`, `skills_handler.go:611-616`), and `GetByName` *does* populate
`skill.Dependencies` + `skill.Resources` (`store.go:143-180`). So names go **out** correctly at
the field level — but come back empty on import (A/B/C), and the encoder emits the wrong nesting
(A) — the exact non-idempotency the register describes (`GAPS…:111`).

### 2.2 Fix — struct restructure (captured proof) + resolve-and-persist

**Restructure the TOML types so nesting is structural** — move `Dependencies` and `Resources`
*inside* `TOMLSkillDef` (the `[skill]` table), tagged `toml:"dependencies"` / `toml:"resources"`.
**Captured proof it round-trips** (`scratchpad/g07_burntsushi_evidence.txt`):

```
Requires=[java.language kotlin.language] Composes=[android.ui.compose] Resources=1 Undecoded=[]
Re-encode -> proper nested [skill.dependencies] + [[skill.resources]] tables
Round-trip re-decode -> Requires/Composes/Resources preserved exactly
```

Then, on the write side:
- **Import** (both REST and `ImportFromTOML`) resolves dependency **names → IDs** and inserts
  edges through `Store.AddDependency` (`graph.go:21-96`), which already validates the relation
  type and runs cycle detection — instead of the current ad-hoc `INSERT` (`import_export.go:178`)
  and the discarded `_ = depName` (`skills_handler.go:552`). Persist resources in the **same
  transaction** as the skill (`ImportFromTOML` already does, `import_export.go:188-196`; the REST
  path must be routed through the store, not `s.pool.CreateSkill`, which has no dep/resource
  persistence for the REST shape).
- **Strict-decode guard (anti-silent-drop):** after decode, fail the import if
  `md.Undecoded()` is non-empty for known-but-unmapped keys (or use `toml.Decode` and inspect
  `MetaData`). A dropped edge/resource must be a **hard error, never a silent skip** — this is
  the mechanical enforcement of the round-trip contract below.

### 2.3 Round-trip integrity contract (the acceptance oracle)

For any skill `x` with dependencies and resources:

1. **Edge preservation:** `export(import(toml)) ⊇` every dependency edge of `toml`, each with its
   `relation_type` **and** the new `optional`/`sort_order` attrs (§3) preserved.
2. **Resource preservation:** every resource (`url`, `title`, `resource_type`) preserved.
3. **Idempotent, byte-equivalent modulo documented normalization:**
   `export(import(export(x))) == export(import(x))` byte-for-byte, where the only permitted
   normalizations are declared explicitly: alias folding (`depends_on`/`prerequisite`→`requires`,
   `part_of`→inverted `composes`, per `skill_granularity…:114-122`), stable key/section ordering,
   and the comment header. Symmetric edges (`related_to`/`alternative_to`) are canonicalized to a
   single stored direction (`skill_granularity…:266`).
4. **No silent loss:** any edge/resource present in the input that cannot be persisted (unknown
   target, unmapped key, FK failure) **aborts the import with a named error** — never a partial,
   "successful" import with fewer edges. This is the direct fix for "edges silently dropped".

JSON path: the same contract holds for the JSON body of `POST /skills`
(`CreateSkillRequest.Deps`/`Resources`, `skills_handler.go:29-30`) — those fields must be
resolved and persisted identically (fixes defect C). (Note: JSON→TOON is a *separate* wire-format
gap, G08 — out of scope here.)

---

## 3. Reconciliation with the granularity model (§4/§5) — and the schema mismatch

The DAG traversal and round-trip must carry the **typed edges + `SkillKind`** from
`skill_granularity_and_composition.md` §4/§5. The current code/schema does **not** — this is a
real mismatch the G06/G07 fix must be coordinated with (it is **P1.T1**'s model:
`IMPLEMENTATION_PLAN.md:170-171`, R16 at `:54`).

| Concern | Granularity model (§4/§5) requires | Current code/schema (as read) | Action |
|---|---|---|---|
| Edge enum | `requires, extends, recommends, composes, related_to, alternative_to` (`…:101-112`) | `models.DependencyType` = only `requires/extends/recommends` (`models.go:24-28`); DB `CHECK` = same three (`migrations/001…:28`, `SPEC.md:230`) | **schema + model change** (§5.3 `…:244-256`) |
| `AddDependency` validation | accept the hard-closure + advisory set | `validTypes` map allows only `requires/extends/recommends` (`graph.go:26-30`) — would **reject** `composes` | widen the map |
| Skill kind | `kind ∈ {atomic, composite, umbrella}` column (`…:197-208, 240-242`) | no `kind` on `skills` (`migrations/001…:6-18`, `models.go:31-48`) | **add column** |
| Edge attrs | `optional BOOLEAN`, `sort_order INT` (`…:224-230, 253-255`) | `skill_dependencies` has neither (`migrations/001…:25-30`) | **add columns** (the CTE §1.2 already selects them) |
| Edge PK | widen to `(skill_id, depends_on, relation_type)` (`…:258`) so a pair may carry >1 typed edge | PK is `(skill_id, depends_on)` (`migrations/001…:29`) — only ONE edge per pair | **widen PK** |
| Hard closure | resolver walks `{requires, composes, extends}` only (`…:126`) | broken CTE walks ALL types (`graph.go:209-211`) | **filter in CTE** (§1.2) |
| Acyclicity scope | cycle check over `{requires, composes, extends}` only; symmetric edges exempt (`…:134-136`) | `hasCycle` (`graph.go:156-175`) and `validate_dag.py` check all three current types jointly | **scope cycle check**; treat `related_to`/`alternative_to` as undirected, excluded |
| Aliases | `depends_on/prerequisite→requires`, `part_of→composes` inverted (`…:114-122`) | none | **normalize at import** |

**Design consequence:** G06's CTE must filter on the **widened** `relation_type` set
`('requires','composes','extends')`, and G07's TOML types must gain the `composes`/`related_to`/
`alternative_to` keys plus per-component `optional`/`sort_order` (the `[[skill.components]]`
ergonomic form, `…:178-191`). Both fixes therefore **depend on the P1.T1 schema migration**
(kind column, widened `relation_type` CHECK, `optional`/`sort_order`, widened PK). Landing G06/G07
against the *current* 3-type schema would ship a design that P1.T1 immediately re-breaks — so the
schema migration is the shared prerequisite (§5, honest-gaps).

---

## 4. Test design (for the later Go fix)

Each PASS cites captured evidence per §11.4.5/§11.4.69; reconciled with
`testing_infrastructure_plan.md:292` (G06) and `:293` (G07). Fixture note: the seed corpus
already contains a **real diamond** — `android.overview` requires both `java.language` and
`kotlin.language`, and `kotlin.language` requires `java.language`
(`seed/skills/{android,kotlin,java}.toml`), so `java.language` is reachable by two paths. Seed
targets `c.language`/`python.language`/`bazel.build_system`/`design_patterns.overview` have **no
seed file** (dangling ⇒ FK-blocked); the multi-level/cycle cases therefore use a **synthetic
A→B→C→D fixture** so they don't depend on unseed­ed skills.

### G06 — traversal (9 cases)

| # | Case | Assertion | Captured evidence a PASS must cite |
|---|---|---|---|
| 1 | Multi-level ≥3 deep (synthetic A→B→C→D) | tree has nodes at depth 1,2,3; D present under C | tree JSON with `max(depth)==3`; per-node id/depth dump |
| 2 | Cycle (A→B→C→A bad data) | traversal **terminates** (bounded time), returns finite tree, no hang | run log with wall-clock < timeout; `is_cycle`-dropped edge count |
| 3 | Transitive edge types (`requires`,`composes`,`extends`) | each pulled into closure | closure list showing a node reached via each type |
| 4 | Non-transitive edges (`recommends`,`related_to`,`alternative_to`) | **excluded** from the hard-closure tree | tree omitting the advisory target; advisory target present in a separate "see-also" set |
| 5 | Diamond (real seed `android.overview`) | shared `java.language` appears **once**; `node_count == \|closure\|` | node-count-equality assertion log; tree JSON showing single `java.language` |
| 6 | Depth bound honored (`maxDepth=2` on the ≥3 chain) | nodes beyond depth 2 absent | tree JSON truncated at depth 2 |
| 7 | Leaf skill (`linux.language`/`java.language`, no hard deps) | tree = root, `Children` empty | tree JSON with empty children |
| 8 | Empty/unknown root | typed `ErrSkillNotFound`, no panic | error string + non-nil-safe return |
| 9 | Wiring parity (REST + MCP + hardened all use the one function) | all three surfaces return the identical tree for a fixed root | side-by-side JSON diff (empty) across surfaces |
| — | **Mutation (§1.1)** | revert assembly to `root.Children = childrenMap[rootID]` (depth-1) ⇒ case 1 FAILS; drop the hard-closure `IN(...)` filter ⇒ case 4 FAILS; remove `CYCLE`/path-guard ⇒ case 2 hangs/FAILS | mutation-run FAIL logs |

### G07 — round-trip (8 cases)

| # | Case | Assertion | Captured evidence |
|---|---|---|---|
| 1 | Decode nested TOML (regression for defect A) | `Dependencies`/`Resources` non-empty; `md.Undecoded()` empty | decoded struct dump; undecoded-keys == `[]` |
| 2 | `convertTOMLWrapper` / import preserves names (defect B) | every edge has resolved `DependsOn`+`DependsOnName` | edge dump with non-nil UUIDs |
| 3 | REST `POST /skills` create-with-deps (defect C) | `skill_dependencies` + `resources` rows exist in DB | `SELECT * FROM skill_dependencies WHERE skill_id=…` dump |
| 4 | export→import→export **byte-stable** | second export == first export byte-for-byte | `diff` output (empty) |
| 5 | All attrs preserved (`relation_type`,`optional`,`sort_order`,resource fields) | attrs survive round-trip | round-trip struct compare showing attrs equal |
| 6 | Alias normalization (`depends_on`,`part_of`) | folded to `requires`/inverted `composes`, one edge (idempotent) | edge dump showing canonical type, no duplicate |
| 7 | Missing-target import | **hard error**, import aborts, **no** partial persist | error return + `SELECT count(*)` == 0 for that skill's edges |
| 8 | MCP `ImportFromTOML` edge fidelity (closes register UNCONFIRMED `GAPS…:112`) | MCP-created skill has all edges/resources in DB | MCP result + DB row dump |
| — | **Mutation (§1.1)** | restore dotted sibling tags ⇒ case 1 FAILS; restore `_ = depName` ⇒ case 2 FAILS; remove `md.Undecoded()` guard ⇒ case 7 (silent-drop) FAILS | mutation-run FAIL logs |

**Counts:** **G06 = 9 cases + 1 mutation triplet; G07 = 8 cases + 1 mutation triplet — 17
positive cases + 6 mutation assertions = 23 test obligations.** This is a superset of the
testing-plan rows (which list unit + integration + property/round-trip + regression/contract +
mutation for each, `:292`/`:293`) — every plan artifact is covered, plus case 9 (wiring parity,
from the rival-copy finding §0.1) and case 1/7 (the newly-proven silent-decode-drop and its
no-partial-persist guard).

---

## 5. Honest gaps (§11.4.6)

**Decidable now (proven this session, no live DB):**
- G06 depth-1 assembly bug (`graph.go:306`), its uncalled/rival status, and the
  missing-typed-filter / no-visited-set secondary defects — read directly from source.
- G07 defects B (`skills_handler.go:552`) and C (`skills_handler.go:164-176,179`) — read directly.
- G07 defect A (TOML nesting/tag mismatch) **and** its fix — **captured** via BurntSushi v1.6.0
  decode/encode experiment (`scratchpad/g07_burntsushi_evidence.txt`); not inferred.
- The schema/model mismatch vs the granularity model — read from `migrations/001…`, `models.go`,
  `graph.go:26-30` against `skill_granularity…` §5.

**Requires the live Postgres (PG16) to prove (integration/property layer, §11.4.108):**
- That the recursive CTE (Form A `CYCLE` / Form B path-array) actually returns the full N-level
  closure, terminates on a real cycle, and yields `node_count == |closure|` on the seed diamond —
  the CTE is DB-executed; `graph_test.go:68-82` already documents (honestly) that these paths need
  a containerized Postgres and cannot be unit-tested in-memory. The G06/G07 tests in §4 that hit
  the DB are **integration** tests gated on a real/containerized PG16 instance.
- End-to-end byte-stable round-trip through the real store transaction (G07 cases 3–8).

**Schema change is REQUIRED (not optional) and must be tracked (§11.4.197):**
- The granularity typed-edge model (`kind` column, widened `relation_type` CHECK,
  `optional`/`sort_order` columns, PK widened to `(skill_id, depends_on, relation_type)`) does not
  exist in `migrations/001_initial.up.sql`, `SPEC.md` §5, or `models.go`. G06's hard-closure filter
  and G07's `composes`/attr round-trip both depend on it. **This migration is owned by P1.T1**
  (`IMPLEMENTATION_PLAN.md:170-171`, R16 `:54`); the G06/G07 Go fix must be sequenced **on/after**
  it, or ship a coordinated migration in the same change. Left un-wired, this is exactly the
  loss-of-requirements failure §11.4.197 forbids — track the G06/G07 fix as a workable item whose
  acceptance includes the P1.T1 schema landing and the §4 integration suite green on PG16.
- **SPEC drift to reconcile:** `SPEC.md:230` and `migrations/001…:28` still show the 3-type
  `CHECK`; both must be updated in lockstep with the migration so the contract, schema, and code
  agree (§11.4.186 anti-divergence).

**Register corrections surfaced (grounded, for the register's next revision):**
- G06 evidence naming `GetDependencyTree` as the live bug is only half-right: that function is
  **uncalled**; the live/MCP tree uses `store.GetTree` (`main.go:230`, `tools.go:215`), and the
  claim that MCP "emits a 1-level tree" (`GAPS…:100`) is **incorrect** — MCP recurses. The fix must
  therefore **consolidate** onto one canonical function and wire it, not merely patch
  `GetDependencyTree`.
- G07's UNCONFIRMED note (`GAPS…:112`) is now **CONFIRMED**: `ImportFromTOML`'s resolver is correct
  but is starved of data by the upstream decode mismatch (defect A, captured) — so the MCP create
  path drops edges too.
