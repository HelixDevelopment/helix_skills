# P1.T1 — Granularity / Typed-Edge Foundation + NEW-1 Cycle-Scope Fix — QA Evidence

**Revision:** 1
**Last modified:** 2026-07-15T19:55:21Z
**Run-id:** 20260715T195521Z
**Scope:** HelixKnowledge Skill Graph System MVP — `project/` (module `github.com/helixdevelopment/skill-system`)
**Anchors:** §11.4.83 (QA transcript) · §11.4.108 (runtime signature) · §11.4.115 (RED baseline) · §11.4.186 (cross-layer divergence) · §11.4.194 (exhaustive review) · §11.4.209 (Fable-xhigh review) · §11.4.134 (iterate-to-GO)

## What shipped

The P1.T1 granularity/typed-edge foundation: `kind` skill column (atomic/composite/umbrella),
the 6-value `relation_type` typed-edge model, triple primary key `(skill_id, depends_on, relation_type)`,
and migration `002_granularity.{up,down}.sql`. Plus the **NEW-1 fix**: `AddDependency` previously
called `hasCycle` **unconditionally** for all 6 relation types, so an advisory candidate edge with a
hard reverse-path (`parent --composes--> child`, then `child --related_to--> parent`) was falsely
rejected as a cycle. The fix adds `models.IsHardClosure(t)` and gates the `hasCycle` **call** on it,
matching the sibling validator `seed/validate_dag.py` (hard-closure-only cycle check).

## Runtime signature (§11.4.108 definition-of-done)

On a live PostgreSQL 16.14 graph, an **advisory candidate edge over a hard reverse-path is ACCEPTED**
(not falsely rejected), while a **hard/hard cycle is still REJECTED** — the observable that proves the
fix is both active and correct.

## Evidence 1 — 2×2 candidate×reverse-path proven live (PostgreSQL 16.14 @ localhost:55433)

```
--- PASS: TestP1T1W2_ComposesCycleIsRejected (0.26s)         [hard / hard      -> REJECT]
--- PASS: TestP1T1W2_RelatedToBackEdgeIsNotACycle (0.27s)    [advisory / advisory -> ACCEPT]
--- PASS: TestP1T1NEW1_AdvisoryOverHardIsNotACycle (0.28s)   [advisory / hard   -> ACCEPT]  <-- NEW-1 driver
--- PASS: TestP1T1NEW1_HardOverAdvisoryAccepted (0.27s)      [hard / advisory   -> ACCEPT]
--- PASS: TestP1T1W2_SecondTypedEdgePerPairAccepted (0.27s)  [W2(a) triple-scoped dup check]
```

## Evidence 2 — RED-polarity proof (§11.4.115), executed by the reviewer

Backed up `graph.go` (md5 `56dd98b1839073336b3b0dbdedd292cb`), reverted ONLY the `IsHardClosure`
gate to the unconditional call:

```
--- FAIL: TestP1T1NEW1_AdvisoryOverHardIsNotACycle (0.26s)
    graph_granularity_test.go:264: ... wrongly rejected as a cycle:
    dependency cycle detected: adding p1t1.new1.b -> p1t1.new1.a would create a cycle
--- PASS: TestP1T1W2_ComposesCycleIsRejected (0.26s)   [hard/hard rejection intact under mutation]
```

Restored byte-identical (md5 re-verified) → re-ran → `ok`. The test genuinely catches the defect
on the pre-fix artifact and nothing else broke under the mutation.

## Evidence 3 — §11.4.186 cross-layer divergence CLOSED (Go ≡ validate_dag.py)

```
$ python3 seed/validate_dag.py seed/CORPUS.yaml
DAG OK (43 nodes, 56 edges)      # exit 0
```

`HARD_CLOSURE_RELATIONS = (requires, composes, extends)` in Python ≡ `models.HardClosureTypes` in Go;
`hasCycle`'s CTE filters `relation_type = ANY($3)` (hard set) and the call is gated on `IsHardClosure`.
All 4 cells match Python semantics. Honest boundary (§11.4.6): Go's check is **incremental** (assumes
the prior hard subgraph acyclic — holds inductively under API-mediated writes); Python's is
**whole-graph** — equivalent for every state reachable through the store API.

## Evidence 4 — full-suite regression (live DB)

All packages `ok`; `internal/db 4.300s` / `internal/skill 3.332s` (real, non-cached DB timing;
PostgreSQL 16.14 confirmed via psql). All 18 `TestP1T1*` cases PASS live. The single SKIP is the
documented pre-existing `pg_trgm` gap in N6/Search (honest §11.4.3 SKIP-with-reason, tracked follow-up).

## Evidence 5 — commit-boundary re-confirmation (conductor, §11.4.6 — not assumed on the reviewer's word)

```
go version go1.26.4-X:nodwarf5 linux/amd64
go build ./...  -> BUILD_OK
go vet   ./...  -> VET_OK
§11.4.84 residue scan (MUTATED/always-pass/MUTATION/_mutated_/RED_MODE=1) -> CLEAN
gofmt -l: graph.go + all 5 new granularity test files -> clean
```

## Independent code review (§11.4.209 / §11.4.142 / §11.4.194 / §11.4.134)

- **Substrate:** Fable model, **xhigh** effort (§11.4.209).
- **Round 3 verdict:** **GO** — zero blocking findings, zero warnings.
- §11.4.194 exhaustive: both `hasCycle` call sites analysed (the second, `import_export.go:171`, proven
  safe-by-construction — the imported skill's ID is fresh `uuid.New()`, unreachable from the existing DB);
  every graph walker's termination under now-legal advisory cycles verified; the gofmt boundary proven
  pre-existing on all 3 touched files.

## Raw evidence (local, §11.4.128 — `*.log` gitignored, not committed)

`project/qa-results/p1t1_remediation/` — `GREEN_full_p1t1_suite.log`, `W2_red_edge_and_cycle_scope.log`,
`W4_N1_red_validate_dag.log`, `W1_red_down_migration_dataloss.log`, `W3_red_export_omits_kind.log`,
`B1_red_dotted_tag_decode.log`, `INDEX.md`.

## Honest boundaries / tracked follow-ups (§11.4.6, §11.4.197 — not lost)

1. **N6/Search `pg_trgm`**: `CREATE EXTENSION pg_trgm` is not yet issued on a clean deployment → `Store.Search`
   is dead there; the single SKIP above captures it honestly. Tracked for the GAPS register.
2. **Pre-existing gofmt debt** on `models/skill.go`, `store.go`, `api/skills_handler.go` (dirty at HEAD, not
   P1.T1-introduced) — deferred to a dedicated project-wide gofmt commit + gofmt-clean gate (§11.4.114 no-conflation).
3. **`graph_test.go:68` `TestHasCycle_RequiresLiveDatabase` stub** — fully superseded by the 3 live cycle tests;
   removed as its own §11.4.124 commit with git-history citation (follows this commit).
