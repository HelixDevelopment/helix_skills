# G11 — Panic-Safe, Real-Work Background Worker: Design Decision

**Revision:** 1
**Last modified:** 2026-07-15T16:31:38Z
**Status:** DECIDED (design-research; drives a later Go fix — no code changed by this doc)
**Scope:** Remediation design for gaps-register **G11** (worker does no real work AND can panic the whole process via unchecked type assertions in a recover-less goroutine), composing with **G03** (flagship pipelines are dead code / worker handlers are stubs), and consuming — without re-doing — the interfaces from **G02** (fail-closed validation executor) and **G05** (fail-closed empty jury).
**Authority:** Constitution §11.4.6 (no-guessing / decision-not-maybe), §11.4.108 (SOURCE→ARTIFACT→RUNTIME verification), §11.4.115 (RED-on-broken-artifact + polarity switch), §11.4.119 (single-resource-owner partitioning), §11.4.147 (crash ≠ complete; respawn-until-done), §11.4.194 (exhaustive all-factor review — a multi-factor crash needs every factor proven), §11.4.197 (started work is wired + verified, never left un-wired), §11.4.201 (every guard asserts the REAL condition; a false-negative is a bluff).

> **Deliverable contract:** this is the ONLY file created by this research pass. It produces a *vetted design*, not code. Every claim is grounded in the real tree with `file:line` citations. Paths are relative to `project/`.

---

## 0. Executive summary

G11 is a **compound** gap with two coupled defects that must be fixed together:

1. **No real work.** The `worker.Runner` never holds an `autoexpand.Pipeline` or a `validation.Pipeline` (`internal/worker/runner.go:53-64` — the struct has only `pool/store/cfg/logger`). Its handlers return fabricated success (`runner.go:320-368`) and its cycles only log (`runner.go:440-507`). Auto-growth, validation, and registry-review — the R2 core promise — are inert (G03).
2. **A recover-less panic vector.** `runRegistryReview` performs two **unchecked type assertions** — `coverage["total_skills"].(int)` (`runner.go:518`) and `coverage["coverage_percentage"].(string)` (`runner.go:519`) — inside the `registryReviewWorker` goroutine (`runner.go:414-434`, launched at `runner.go:148-149`), and **no worker goroutine has a `recover()`** (`runner.go:221, 375, 396, 414`). An unrecovered panic in any goroutine terminates the **whole process** (Go runtime `exit(2)`), and the gin `Recovery()` middleware (`cmd/server/main.go:146`, register cites `middleware.go:254-276`) covers only the **server** process's HTTP handlers, never the **separate** worker process (`cmd/worker/main.go`).

These are coupled: **fixing "no real work" without the firewall converts a latent panic into a live one** — the real pipelines (nil derefs, slice OOB, div-by-zero, LLM/DB edge cases) will run inside the same recover-less goroutines. The design therefore lands both halves in one change: a per-goroutine + per-job `recover()`/restart supervisor, comma-ok + typed-struct returns replacing the bare assertions, and the real `autoexpand`/`validation` cycles wired in (composing G03), all single-owner-guarded (§11.4.119).

---

## 1. Root-cause the panic (the exact defect, `file:line`)

### 1.1 The two unchecked assertions

`internal/worker/runner.go:509-527` — `runRegistryReview`:

```go
coverage, err := r.store.GetCoverage(ctx, "")   // runner.go:511  → map[string]interface{}
...
r.logger.Info("registry review completed",
    zap.Int("total_skills", coverage["total_skills"].(int)),          // runner.go:518  UNCHECKED
    zap.String("coverage", coverage["coverage_percentage"].(string)), // runner.go:519  UNCHECKED
)
```

`GetCoverage` (`internal/skill/store.go:456-544`) returns a `map[string]interface{}` populated with `"total_skills": total` (an `int`, `store.go:537`) and `"coverage_percentage": fmt.Sprintf("%.1f%%", …)` (a `string`, `store.go:542`). So on the **current** exact call graph the two assertions happen to succeed. That is the trap: they are **unsafe by construction** and panic under any of —

- **Missing key** — if any future path returns a map without `"total_skills"`, `coverage["total_skills"]` is `nil` and `nil.(int)` panics (`panic: interface conversion: interface {} is nil, not int`).
- **Type drift** — the same value re-hydrated from JSON (the map is already `json.Marshal`-ed into the audit log at `runner.go:523`) comes back as `float64`, not `int`; a different count source could return `int64`. A bare `.(int)` on `float64` panics.
- **Real-work coupling (the load-bearing reason)** — once the real cycles run (§2.3 / G03), a panic from **anywhere** in the cycle (a nil `*Pipeline`, an OOB slice in draft assembly, an LLM/DB nil deref) has the identical blast radius, because there is no firewall around the goroutine.

This is exactly the §11.4.201 defect class: a boundary that does **not** assert the real condition it depends on, taking the non-conservative (crash) path on a bad shape instead of an honest logged error.

### 1.2 The recover-less goroutines

Every worker loop is a bare `go func` with a `defer r.wg.Done()` and **no `defer recover()`**:

- `processJobQueue` — `runner.go:221` (launched `runner.go:133`)
- `autoExpandWorker` — `runner.go:375` (launched `runner.go:138`)
- `validationWorker` — `runner.go:396` (launched `runner.go:143`)
- `registryReviewWorker` — `runner.go:414` (launched `runner.go:149`) ← the one that reaches the assertions

None wraps its body. There is no supervisor, no per-job firewall in `executeJob` (`runner.go:280-314`), and no restart.

### 1.3 Why a stub handler + a bad assertion crashes the *process*

Two Go facts make this a full-process outage rather than a contained error:

1. **An unrecovered panic in *any* goroutine kills the *entire* program.** Go does not isolate a goroutine's panic; it unwinds that goroutine, finds no `recover()`, and the runtime prints the stack and calls `exit(2)` — every other goroutine (job queue, auto-expand, validation) dies with it. So the single `coverage["total_skills"].(int)` on `runner.go:518` can take down the whole worker.
2. **The only `recover()` in the system does not reach here.** `cmd/server/main.go:146` installs `gin.Recovery()`, but that wraps HTTP request goroutines **in the server process**. The worker is a **distinct binary/process** (`cmd/worker/main.go:33-145`) with its own `main`, its own signal loop (`main.go:110-116`), and **no** recovery wrapper. G11's own evidence line makes this explicit (`GAPS_AND_RISKS_REGISTER.md:154`).

The **stub angle** is what makes this urgent rather than merely latent: the handlers today return `Success:true` with no work (`runner.go:334, 348, 362, 367`) and the stub comments assert work that never happens ("Actual expansion is done by the autoExpandWorker polling loop", `runner.go:332`; "Actual validation is done by the validationWorker", `runner.go:347` — bluff-comments, per G03 / `GAPS_AND_RISKS_REGISTER.md:72`). The moment G03 replaces those stubs with real, panic-capable pipeline calls **inside the same recover-less goroutines**, a routine runtime error (an unreachable resource URL, a malformed LLM response, a nil embedder) becomes a process-killing crash loop. **Fixing G03 without G11's firewall ships a worse bug.**

### 1.4 Blast radius (why "outage", not "log line")

On crash: the worker process dies; the ~100-slot in-flight job buffer (`jobChan`, `runner.go:105`) is lost; auto-growth/validation/registry-review all stop; a systemd restart re-enters the same cycle and — if the triggering data shape is persistent — **crash-loops**. Per §11.4.147, a quota/panic-killed worker is `crashed`, never `complete`; the current design has **no** registry to know work is still owed, so the loss is silent.

---

## 2. Decision

### 2.1 Replace the unchecked assertions with a typed struct + comma-ok at any residual boundary

**Primary fix — typed return (compile-time safety).** Introduce a typed coverage struct and have the store expose it, so the worker never touches an `interface{}`:

```go
// internal/skill/store.go (new)
type CoverageStats struct {
    Domain             string  `json:"domain"`
    TotalSkills        int     `json:"total_skills"`
    SkillsWithDeps     int     `json:"skills_with_deps"`
    SkillsWithEvidence int     `json:"skills_with_evidence"`
    SkillsMissingDeps  int     `json:"skills_missing_deps"`
    AverageCoverage    float64 `json:"average_coverage"`
    CoveragePercentage string  `json:"coverage_percentage"`
}

func (s *Store) GetCoverageStats(ctx context.Context, domain string) (*CoverageStats, error)
```

`GetCoverageStats` is the same SQL body as `GetCoverage` (`store.go:456-543`) returning the struct instead of the map. The worker's `runRegistryReview` consumes `stats.TotalSkills` / `stats.CoveragePercentage` directly — **no assertion, no panic surface**. Recommended: migrate the existing `GetCoverage` map-callers (the API `/coverage` path uses `registry.GetCoverageReport`, `cmd/server/main.go:243`, not this map, so the map has few callers) and delete the map form, eliminating the class rather than one instance (§11.4.124 — don't leave the unsafe shape around).

**Secondary fix — comma-ok at any unavoidable `interface{}` boundary.** Where an `interface{}` genuinely cannot be typed away (e.g. an audit-log detail re-decode), use comma-ok with an **honest error**, never a bare assertion — the §11.4.201 conservative-safe default:

```go
v, ok := m["total_skills"].(int)
if !ok {
    r.logger.Error("registry review: coverage shape unexpected",
        zap.String("key", "total_skills"), zap.Any("got", m["total_skills"]))
    return // honest error path, worker survives — NOT a crash
}
```

This also fixes the sibling `pipeline.go` assertion pattern by example: note `autoexpand/pipeline.go:215` (`openaiLLM, ok := p.llm.(*OpenAILLM)`) is *already* comma-ok-safe — the runner's bare `.(int)`/`.(string)` are the outliers, and they align to the safe pattern the codebase already uses elsewhere.

### 2.2 Wrap EVERY worker goroutine (and job) in a `recover()` + restart/backoff supervisor

Add a single supervisor helper; every long-lived loop launches through it, so a panic is **logged + recovered + restarted** and the worker **never dies**:

```go
// supervise runs fn under a panic firewall; on panic it logs the stack,
// records the panic metric, and restarts fn with capped backoff until ctx is done.
func (r *Runner) supervise(ctx context.Context, name string, fn func(context.Context)) {
    defer r.wg.Done()
    backoff := time.Second
    for {
        select { case <-ctx.Done(): return; default: }
        func() {
            defer func() {
                if p := recover(); p != nil {
                    r.metrics.recordPanic()                       // new metric: PanicsRecovered
                    r.logger.Error("worker goroutine panic — recovered, restarting",
                        zap.String("worker", name), zap.Any("panic", p),
                        zap.ByteString("stack", debug.Stack()))
                }
            }()
            fn(ctx)   // normal return = loop exited (ctx cancelled) → supervise returns next tick
        }()
        // panic path falls through to backoff + restart; clean return re-checks ctx and exits
        select { case <-ctx.Done(): return; case <-time.After(backoff): }
        backoff = min(backoff*2, 30*time.Second)
    }
}
```

`Start` (`runner.go:118-155`) launches `go r.supervise(ctx, "registry_review", r.registryReviewLoop)` etc. instead of the bare `go r.registryReviewWorker(ctx)`. The loop bodies (`runner.go:375-434`) lose their own `defer r.wg.Done()` (the supervisor owns it) and keep their `ticker` + `ctx.Done()` select.

**Per-job firewall too.** `executeJob` (`runner.go:280-314`) runs the handler dispatch (`runner.go:290-301`) inside a `func() (res JobResult) { defer func(){ if p:=recover(); p!=nil { res = JobResult{Success:false, Error: fmt.Sprintf("handler panic: %v", p)} } }() ; ... }()`. A panicking job then becomes a **recorded failure** that flows into the existing retry/`recordFailure` path (`runner.go:236-277`), never a process death — the §11.4.147 "crash ≠ done, keep it owed, retry" discipline at the job granularity.

**Backoff bounds (§11.4.6 no-guessing):** restart backoff is capped (1s→30s) so a persistently-panicking cycle does not busy-loop nor hide a real defect — the `PanicsRecovered` metric + logged `debug.Stack()` surface it for systematic-debugging (§11.4.102). A recovered panic is a **tracked signal**, not a swallowed error.

### 2.3 Realize the real cycles (compose G03) — delete the stubs

Add the pipelines to the runner and call them; the cycles do work, not logging:

- **Struct:** add `autoExpand *autoexpand.Pipeline` and `validator *validation.Pipeline` fields to `Runner` (`runner.go:53-64`).
- **Construction:** `NewRunner` (`runner.go:99-107`) builds them — `autoexpand.NewPipeline(store, embedder, cfg.AutoExpand, logger, opts...)` (`autoexpand/pipeline.go:64`) and `validation.NewPipeline(store, cfg.Validation, logger, opts...)` (`validation/pipeline.go:100`). **Wiring dependency:** `autoexpand.NewPipeline` requires a `db.Embedder` that `cmd/worker/main.go` does not currently construct — the worker `main` (`cmd/worker/main.go:86-94`) must build the embedder (from `cfg.Embedding`) and pass it in. Flagged in §5.
- **`runAutoExpandCycle`** (`runner.go:440-481`): for each auto-expandable gap entry (already fetched via `r.store.GetMissingSkills`, `runner.go:442`), call `r.autoExpand.Run(ctx, entry.SkillName, r.cfg.AutoExpand.MaxDepth)` (`autoexpand/pipeline.go:319`) instead of only incrementing `processed` (`runner.go:479`). Honour `MaxNewSkillsPerRun` (`runner.go:460`) and `ctx.Done()` (`runner.go:468-472`).
- **`runValidationCycle`** (`runner.go:483-507`): for each draft skill (already listed via `ListSkills(…, SkillStatusDraft, …)`, `runner.go:485`), call `r.validator.Validate(ctx, &sk)` (`validation/pipeline.go:146`) and persist the `*ValidationResult` (`validation/pipeline.go:31-37`) instead of only logging (`runner.go:502`). A skill only advances to `validated`/`active` on a **recorded** jury verdict with `ApprovedBy ≥ approval_threshold` (composes G03 + G05 fail-closed empty-jury).
- **`runRegistryReview`** (`runner.go:509-527`): keep the coverage read (now `GetCoverageStats`), and additionally do the real work the comment promises ("mark stale skills, recalculate coverage", `runner.go:510`) — update `skill_registry.stale`/`coverage`/`last_review`.
- **Handlers** (`runner.go:320-368`): `handleAutoExpand`/`handleValidate` dispatch to the same pipelines for the on-demand (queued-job) path; **delete the bluff-comments** (`runner.go:317, 332, 347`) per G03.

### 2.4 Alternatives rejected

- **Keep stubs "until P4/P5", fix only the assertion.** Rejected — un-wired research (§11.4.197) and, worse, leaves the recover-less goroutines so the eventual real wiring ships a live crash. G11's own decision line rejects this (`GAPS_AND_RISKS_REGISTER.md:156`).
- **A single top-level `recover()` in `main` only.** Rejected — recover only works in the *same* goroutine's deferred call; a top-level `main` recover cannot catch a panic in a worker goroutine. Per-goroutine + per-job firewalls are mandatory.
- **Swallow recovered panics silently.** Rejected as a §11.4.201/§11.4.6 bluff — a recovered panic MUST be logged with its stack + counted (`PanicsRecovered`) so the underlying defect is investigated (§11.4.102), never hidden.
- **Leave the map + assert everywhere with comma-ok.** Acceptable but weaker than the typed struct; the typed `CoverageStats` removes the boundary entirely (compile-time), which is the stronger guarantee — comma-ok is the fallback only where `interface{}` is unavoidable.

---

## 3. Lifecycle: start / stop / health / single-owner

- **Start/Stop already exist and are sound** — `Start` guards double-start under `r.mu` (`runner.go:118-129`), `Stop` cancels the context and waits on the `WaitGroup` with a shutdown-context timeout (`runner.go:159-183`). Preserve both; only re-route the goroutine launches through `supervise` (§2.2). Graceful shutdown already drains via `wg.Wait()` + `ctx.Done()` selects in each loop (`runner.go:171-182`) — the supervisor's clean-return path re-checks `ctx.Done()` and exits, so cancellation still terminates every loop.
- **Health.** Add `Runner.Health() WorkerHealth` reporting `{running, lastCycleSuccess per worker, panicsRecovered, jobsProcessed/Failed}` sourced from `Metrics` (extend the struct at `runner.go:81-89` with `PanicsRecovered int64` and `LastCycle map[string]time.Time`). Expose it via a small `--health` HTTP endpoint or a liveness file for the `systemd --user` unit (deploy home per G13) so a wedged (but not crashed) worker is detectable — a recovered-panic crash-loop shows as `panicsRecovered` climbing while `lastCycleSuccess` stalls.
- **Single-owner (§11.4.119).** `runAutoExpandCycle` and `runRegistryReview` **write** the shared `skills` / `skill_registry` tables. If more than one worker replica runs, they double-expand and race the registry. Mandate exactly one owner per write-cycle-class via a Postgres **advisory lock** (`SELECT pg_try_advisory_lock($1)` keyed on a per-cycle constant): only the lock holder runs the write cycle; a non-holder logs an honest skip-with-reason (`another worker owns autoexpand`) and moves on. This is the DB analogue of the device-lock — event-driven claim, released at cycle end so the next tick's holder can be any replica. Read-only reporting (health, metrics) needs no lock.

---

## 4. Test design (RED-first, anti-bluff, evidence-cited)

Every test is authored to **reproduce the defect on the pre-fix code first** (§11.4.115 RED-on-broken-artifact), then flip to a GREEN regression guard; every PASS cites captured evidence (§11.4.5/§11.4.69). **9 test cases across 5 layers** (unit / integration / chaos / mutation, plus the reconciliation row):

| # | Layer | Case | RED (pre-fix) | GREEN (post-fix) | Captured evidence |
|---|---|---|---|---|---|
| T1 | unit | `TestRegistryReview_MalformedCoverage_NoPanic` — feed coverage with a missing / wrong-typed key | bare `.(int)` panics → **test process dies** | comma-ok/typed path returns honest error, no panic | `no-panic.log` (recover+continue), `t.Log` verdict |
| T2 | chaos | `TestSupervisor_HandlerPanic_WorkerSurvives` — inject a cycle fn that panics | recover-less `go func` → process `exit(2)`, suite aborts | `supervise` recovers, logs stack, restarts; `runner.IsRunning()==true`, `PanicsRecovered>=1`, next cycle runs | worker log showing recover()+restart; metrics snapshot |
| T3 | chaos | `TestExecuteJob_PanickingJob_RecordedNotFatal` — a handler panics on a job | process dies mid-queue | job → `recordFailure` (`runner.go:592`), queue keeps draining, worker alive | job-failed audit row; worker-alive log |
| T4 | unit | `TestGetCoverageStats_TypedFields` — seed DB → typed struct | n/a (new API) | `TotalSkills` int, `CoveragePercentage` string exact; nil-value boundary → error not panic | struct dump vs expected |
| T5 | integration | `TestAutoExpandCycle_CreatesRealSkill` — seed a gap, run cycle (real DB, §11.4.27 no mocks) | cycle only logs → **0 rows created** (proves "no real work") | `autoExpand.Run` drafts + `store.Create` → new `draft` skill row exists | created-skill row (`SELECT … WHERE name=…`) |
| T6 | integration | `TestValidationCycle_NoValidatedWithoutVerdict` — draft with empty/near-empty jury | (n/a; validation inert pre-fix) | skill stays `draft`/blocked unless `ApprovedBy≥threshold` (composes G05) | jury-verdict log; skill status unchanged |
| T7 | mutation (§1.1) | reintroduce bare `coverage["total_skills"].(int)` | — | **T1 FAILs** (panic returns) — proves the typed/comma-ok guard is load-bearing | mutation diff + failing run |
| T8 | mutation (§1.1) | remove the `recover()` from `supervise` | — | **T2 FAILs** (process dies) — proves the firewall is load-bearing | mutation diff + failing run |
| T9 | mutation (§1.1) | remove the per-job recover in `executeJob` | — | **T3 FAILs** (panicking job kills process) | mutation diff + failing run |

**Reconciliation with `research/testing_infrastructure_plan.md` (§11.4.186 anti-divergence).** The plan's G11 row (`testing_infrastructure_plan.md:297`) specifies exactly: *unit (cycle calls pipeline; coverage type-safety comma-ok), integration (worker creates real skill from seeded gap), chaos (malformed coverage map ⇒ logged error not panic), mutation (reintroduce bare assertion → panic test FAILs)*, with evidence *"worker log showing recover()+continue; created-skill row; chaos malformed-input `no-panic.log`"*. This design's cases map onto it with no divergence: T4/T2 ↔ "coverage type-safety comma-ok" + unit row (`plan:190` "extend to … worker"); T5 ↔ "worker creates real skill from seeded gap"; T1/T2 ↔ the chaos row (`plan:197` "Worker panic-safety (G11: malformed coverage map ⇒ logged error not process death; per-goroutine `recover()`)"); T7 ↔ "reintroduce bare assertion → panic test FAILs". T3/T8/T9 (per-job firewall) **extend** the plan's mutation cell — recommend adding a per-job-panic sub-row to `plan:297` so the two docs stay in lockstep. Cleanup in `t.Cleanup`/`trap` per §11.4.14 (plan:197). Challenge **fabricated-skill-must-fail** and the HelixQA autonomous-pipeline session (plan:289/297) are inherited from G03, not duplicated here.

---

## 5. Honest gaps & cross-gap dependencies (§11.4.6)

- **Requires the live DB.** T4 (`GetCoverageStats`), T5 (real skill creation), T6 (jury verdict), and the §11.4.119 advisory-lock behaviour all need the real Postgres+pgvector (docker-compose per G13 / `plan:197` killable container), §11.4.27 no-mocks. They cannot be proven with a unit fake; they are integration/chaos gated on the test-harness bootstrap (G04/P0.T3).
- **Depends on G03 wiring, not re-done here.** The real cycles (§2.3) consume interfaces this doc does **not** re-implement:
  - `autoexpand.NewPipeline(store, embedder, cfg.AutoExpand, logger, opts...)` (`autoexpand/pipeline.go:64`) → `(*Pipeline).Run(ctx, skillName, maxDepth) (*ExpansionResult, error)` (`pipeline.go:319`).
  - `validation.NewPipeline(store, cfg.Validation, logger, opts...)` (`validation/pipeline.go:100`) → `(*Pipeline).Validate(ctx, *models.Skill) (*ValidationResult, error)` (`pipeline.go:146`), returning `ValidationResult{Passed, Stage, Details, ApprovedBy}` (`pipeline.go:31-37`).
  - **Embedder wiring gap:** `autoexpand.NewPipeline` needs a `db.Embedder` that `cmd/worker/main.go` currently never constructs (`cmd/worker/main.go:86-94` builds only the store). The G11 fix's `NewRunner` change requires the worker `main` to build and inject the embedder — a small but real wiring item, tracked here so it is not lost.
- **Depends on G02 + G05 fail-closed defaults.** The validation cycle's safety rests on G02's `IsolatedExecutor` fail-closed default (SKIP-with-reason `isolation_runtime_absent` until the `vasic-digital/containers` submodule is vendored — `g02_sandbox_faildesign.md:129-132`) and G05's empty-jury→fail-closed (`GAPS_AND_RISKS_REGISTER.md:92-95`). Until the real LLM jury (P3 ModelProvider chain, G05) exists, the wired validation cycle correctly **blocks / requires human review**, never auto-passes — the worker design assumes and must preserve that fail-closed posture, it does not build the jury.
- **Latent-vs-live honesty on the panic.** On today's *exact* call graph the bare `.(int)`/`.(string)` (`runner.go:518-519`) do not panic, because `store.GetCoverage` always returns those keys with those types (`store.go:537, 542`). The precise claim (§11.4.6, no overclaim) is: the assertions are **unsafe by construction** — a §11.4.201 guard-less boundary that panics under any map-shape drift (missing key / JSON-round-trip `float64` / `int64` source) **and** the recover-less goroutine becomes a guaranteed process-crash surface the moment §2.3's real, panic-capable work runs. The fix removes both the boundary and the crash blast radius; it does not claim the pre-fix code panics on every run today.

---

## 6. Runtime signature (definition-of-done, §11.4.108)

The fix is DONE only when, on a clean deployment: (a) a worker goroutine forced to panic logs a recover+stack and the process **stays alive** with `PanicsRecovered≥1` (T2/T8 flip RED→GREEN); (b) a seeded gap produces a real `draft` skill **row** via the wired `autoExpand.Run` (T5); (c) `runRegistryReview` reads coverage through the typed `CoverageStats` with **zero** `interface{}` assertions in the worker package (grep gate: no `coverage[…].(int|string)` in `runner.go`); (d) the §1.1 mutations T7/T8/T9 each make their guard test FAIL. Source-green is not done — the worker is a separate process, so verification runs against the running worker binary on the real DB.
