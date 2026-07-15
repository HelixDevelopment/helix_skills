# G20 — Real Auto-Growth: LLMClient Pluggability + Fail-Closed No-Provider Behaviour + Transactional Resource Persistence

**Revision:** 1
**Last modified:** 2026-07-15T16:45:00Z
**Status:** design-research, no code landed
**Scope:** the Go backend `github.com/helixdevelopment/skill-system` under `project/`, package `internal/autoexpand`.
**Authority / mandates served:** G20 (`GAPS_AND_RISKS_REGISTER.md:233-239`) — auto-expand fabricates
placeholder skills without an LLM; couples to concrete `*OpenAILLM`; drafted resources never
persisted. Constitution §11.4.6 (no-guessing), §11.4.8/§11.4.150 (deep-research-before-implementation),
§11.4.108 (four-layer runtime-signature), §11.4.115 (RED-first + polarity switch), §11.4.186
(cross-doc lockstep), §11.4.197 (research fully completed/wired, never left un-wired), §11.4.201
(every guard/gate asserts the REAL condition; fail-closed on an unresolvable signal, never fabricate).
**Read discipline:** every Go fact below was read from the committed baseline ref `255061b` via
`git show 255061b:…project/<path>` (the working tree has an uncommitted G02 change and was NOT read
for `project/` source). No existing file was modified; this is the single new deliverable.
**Composes with (read in full, no re-derivation):** `research/g11_worker_design.md` (G11 — the
panic-safe worker that will actually CALL `Pipeline.Run` on a real cycle; this design's "no
placeholder" contract is what G11's `runAutoExpandCycle` wiring depends on being safe to call
unconditionally), `research/g10_embedding_provider_design.md` (G10 — the `db.Embedder` that
`autoexpand.NewPipeline(store, embedder, cfg.AutoExpand, logger, opts...)` requires as its second
argument; this design does not re-derive the embedder-wiring gap, it inherits G10's fix), G03
(`GAPS_AND_RISKS_REGISTER.md:65-75` — the pipeline is dead code today, never instantiated; this design
makes it *safe* to instantiate, it does not re-litigate the instantiation wiring itself), G05
(`GAPS_AND_RISKS_REGISTER.md:89-95` — fail-closed empty jury; this design's created skills stay
`draft` until a REAL G05 jury verdict advances them, it never marks `validated`/`active` itself).

---

## 0. One-paragraph problem statement

`autoexpand.Pipeline` promises "fully-automatic auto-growth" (R2) but today does the opposite of what
zero-bluff (R11) requires: with no LLM configured — the *default* state, since `Pipeline.llm` is never
set anywhere in the current call graph — it silently writes a **fabricated** skill row containing
boilerplate template text as if it were real knowledge; when an LLM *is* configured, `DraftSkill`
refuses every implementation of its own `LLMClient` interface except the one concrete struct
`*OpenAILLM`, defeating R7's pluggable-provider mandate; and even on the one path that does call a
real LLM, the resources the LLM suggested are computed, given a `SkillID`, and then **silently
dropped** — never written to the database. Three independent defects, one root class: the pipeline
does not fail closed when it cannot do genuine work, it fabricates.

---

## 1. ROOT CAUSE — proven from `git show 255061b:…`, exact file:line

### 1.1 Placeholder fabrication is the DEFAULT path, not a rare fallback

`internal/autoexpand/pipeline.go:209-213` (`DraftSkill`):

```go
209:	// Use LLM to generate skill draft
210:	if p.llm == nil {
211:		// Fallback: create a minimal skill without LLM
212:		draft := p.createMinimalDraft(gap)
213:		return draft, nil, nil
```

`p.llm` is set **only** via the `WithLLMClient` functional option (`pipeline.go:52-57`). Grepping the
committed tree at `255061b` for `NewPipeline(` and `WithLLMClient(` outside `internal/autoexpand`
itself finds **zero call sites** — `autoexpand.NewPipeline` is never invoked anywhere in `cmd/server`
or `cmd/worker` (the same "never instantiated" fact G03 records, `GAPS_AND_RISKS_REGISTER.md:69`, and
that `research/g11_worker_design.md:150` and `research/g10_embedding_provider_design.md:106-115`
independently confirm for the worker's missing `db.Embedder` wiring). **Consequence: on the current
call graph, `p.llm == nil` is not an edge case — it is the ONLY reachable state**, so
`createMinimalDraft` is the pipeline's *entire* current behaviour, not a fallback.

`createMinimalDraft` (`pipeline.go:282-311`) fabricates a real, persisted `models.Skill`:

```go
282:func (p *Pipeline) createMinimalDraft(gap Gap) *models.Skill {
283:	return &models.Skill{
...
288:		Description: fmt.Sprintf("Auto-generated skill to fill gap: %s", gap.Reason),
289:		Content: fmt.Sprintf(`# %s
290:
291:## Overview
292:
293:This skill was auto-generated to fill a gap in the knowledge graph.
...
304:> This skill needs review and enrichment before activation.
305:`, gap.SuggestedTitle, gap.SkillName, gap.MissingDepName, gap.Reason),
306:		Status: models.SkillStatusDraft,
```

`Run` (`pipeline.go:319-420`) then unconditionally persists whatever `DraftSkill` returned —
`p.store.Create(ctx, draft)` at `pipeline.go:393` — with **no branch that distinguishes "genuine LLM
content" from "boilerplate placeholder"**. The register's evidence line is precise
(`GAPS_AND_RISKS_REGISTER.md:236`): this is fake knowledge (R11) written as if it were a real skill,
and once G11 wires `runAutoExpandCycle` to actually call `Run` on a schedule (`g11_worker_design.md
§2.3`), the graph floods with these rows automatically, unattended.

### 1.2 The concrete `*OpenAILLM` coupling defeats the `LLMClient` interface it was handed

`internal/autoexpand/llm.go:26-29` defines the interface the pipeline is supposed to program to:

```go
26:type LLMClient interface {
27:	// Generate creates text from a prompt with a token limit.
28:	Generate(ctx context.Context, prompt string, maxTokens int) (string, error)
29:}
```

But `DraftSkill` immediately discards that abstraction — `pipeline.go:215-218`:

```go
215:	openaiLLM, ok := p.llm.(*OpenAILLM)
216:	if !ok {
217:		return nil, nil, fmt.Errorf("unsupported LLM client type")
218:	}
219:	draft, resources, err := openaiLLM.GenerateSkillDraft(ctx, gap.MissingDepName, context)
```

Any `LLMClient` implementation that is **not literally the struct `*OpenAILLM`** — a HelixLLM client,
an LLMsVerifier-fronted client, a Claude-Toolkit-alias client, or a test double — hits the `!ok` branch
and errors, **which is then swallowed and converted into the same `createMinimalDraft` fabrication**
(`pipeline.go:222-227`, "LLM skill drafting failed, using fallback"). R7 (`REQUIREMENTS.md:68-69`)
explicitly mandates "pluggable ModelProvider, not hardcoded OpenAI" — `DraftSkill` hardcodes OpenAI
via a type assertion on the very interface it declares as its abstraction point.

**Why the assertion exists at all — the real defect is method placement, not the assertion syntax.**
`GenerateSkillDraft(ctx, skillName, existingContext string) (*models.Skill, []models.Resource, error)`
is defined **only** on `*OpenAILLM` (`llm.go`, the "Skill draft generation" section), not on the
`LLMClient` interface. Reading `OpenAILLM.GenerateSkillDraft`'s own body shows it does nothing
`OpenAILLM`-specific: it builds a prompt with the free function `GeneratePrompt(skillName,
existingContext)`, calls `c.Generate(ctx, prompt, 4000)` (the one interface method), and parses the
result with the free function `parseSkillDraft(skillName, response)`. **Both `GeneratePrompt` and
`parseSkillDraft` are already provider-agnostic free functions** — neither touches `c.apiKey`,
`c.baseURL`, or any other `OpenAILLM`-specific field. The type assertion exists purely because this
three-line orchestration (`GeneratePrompt` → `Generate` → `parseSkillDraft`) was written as a method on
the concrete struct instead of as pipeline-level logic driven by the interface it already has.

### 1.3 Drafted resources are computed, tagged, and then dropped

`Run` — `pipeline.go:401-403`:

```go
399:					// Add resources if any
400:					for i := range resources {
401:						resources[i].SkillID = draft.ID
402:					}
```

This loop mutates the **local** `resources` slice's `SkillID` field and then does nothing further with
it — the loop body has no store call. `resources` is never referenced again in `Run` (the next
statement, `pipeline.go:405`, logs `draft.Name` and moves to `nextLayer` bookkeeping). The store
package already provides the exact transactional primitive this needs:
`(*skill.Store).BulkAddResources(ctx, skillID, resources)` (`internal/skill/resources.go:221-254`) —
opens one `pgx.Tx`, inserts every resource row, calls `recalcCoverage`, and writes an audit event
(`resource.bulk_added`) — but `autoexpand.Pipeline.Run` never calls it, nor `AddResource`
(`resources.go:19-67`), nor any other persistence path. Every resource the LLM (when one is wired)
suggests for a newly-drafted skill is silently discarded before `Run` returns.

---

## 2. DECISION (all decisions fail-CLOSED per §11.4.201; never fabricate)

### 2.1 Program `DraftSkill` to `LLMClient.Generate` directly — delete the concrete type assertion

Move the three-line orchestration out of `OpenAILLM.GenerateSkillDraft` and into
`Pipeline.DraftSkill`, calling only the interface method:

```go
// DraftSkill (revised) — no concrete-type assertion; any LLMClient works.
prompt := GeneratePrompt(gap.MissingDepName, context)
response, err := p.llm.Generate(ctx, prompt, 4000)
if err != nil {
    return nil, nil, fmt.Errorf("llm generate: %w", err)
}
draft, resources, err := parseSkillDraft(gap.MissingDepName, response)
if err != nil {
    return nil, nil, fmt.Errorf("parse llm draft: %w", err)
}
```

`GeneratePrompt` and `parseSkillDraft` are already free functions (`llm.go`) — zero change needed to
either. `OpenAILLM.GenerateSkillDraft` MAY remain as a convenience wrapper for direct unit-testing of
`OpenAILLM` in isolation, but `Pipeline.DraftSkill` no longer calls it and no longer references the
concrete `*OpenAILLM` type anywhere. **Any** `LLMClient` implementation — HelixLLM, LLMsVerifier,
Claude-Toolkit alias, a local test double — now drafts skills identically, satisfying R7 with zero new
interface surface. **Rejected alternative:** adding `GenerateSkillDraft` to the `LLMClient` interface
itself — rejected because it forces every future provider implementation to re-implement identical
prompt/parse orchestration; keeping that orchestration at the pipeline layer, driven by the one
`Generate` primitive, is the smaller, more DRY interface (composes with G10's `ModelProvider`
terminology from R7 — the embedding side already treats "provider" as a thin transport, and this
mirrors that shape for the LLM side).

### 2.2 No-provider run is an honest, tracked SKIP — never a fabricated skill

Delete `createMinimalDraft` (`pipeline.go:282-311`) entirely — it is the exact "no genuine LLM content
⇒ fabricate template text as a real skill" pattern G20 forbids by name (register decision line,
`GAPS_AND_RISKS_REGISTER.md:238`: *"Never persist a placeholder as a real skill — either produce
genuine LLM content or mark the gap as unfilled"*). Replace the `p.llm == nil` branch with an honest,
closed-vocabulary skip that creates **nothing**:

```go
if p.llm == nil {
    p.logger.Warn("auto-expand: no LLM/ModelProvider configured — gap left unfilled",
        zap.String("missing_dep", gap.MissingDepName), zap.String("reason", "no_provider_configured"))
    return nil, nil, ErrNoProviderConfigured // sentinel; Run treats it as SkillsSkipped++, never an Errors entry that implies retry-will-help
}
```

Symmetrically, the existing fallback-on-LLM-failure path (today's `pipeline.go:222-227`, "LLM skill
drafting failed, using fallback") is **also** replaced — a genuine `Generate`/parse failure returns the
real error (from §2.1's revised `DraftSkill`) instead of silently degrading into `createMinimalDraft`.
`Run`'s per-gap loop distinguishes the two outcomes:

- `ErrNoProviderConfigured` → increment a new `ExpansionResult.SkillsSkipped` counter with reason
  `no_provider_configured`; do not append to `Errors` (this is a configuration state, not a fault) and
  do not add the gap to `nextLayer` (nothing was created to recurse into).
- any other `DraftSkill` error (LLM unreachable, malformed JSON, etc.) → keep the existing
  `result.Errors = append(...)` behavior (`pipeline.go:387-390`) unchanged — these ARE faults, and
  should be visible for the next §11.4.147 respawn/retry cycle, but still create **zero** skill rows.

**Two distinct nil-LLM causes, both honest, both fail-closed:** (a) `cfg.AutoExpand.Enabled == false` —
the worker/wiring layer (G11/G03) should not even schedule a cycle in this case, so `Run` is never
called; (b) `Enabled == true` but no `LLMClient` was constructed for the pipeline (provider
misconfigured/unresolvable) — `DraftSkill`'s `ErrNoProviderConfigured` path fires per-gap, loud in the
logs, never silent, never a fabricated row. **Rejected alternative:** keeping `createMinimalDraft` as
an opt-in "degraded mode" flag — rejected per the register's own alternatives-rejected line
(`GAPS_AND_RISKS_REGISTER.md:238`): "graceful degradation" here degrades directly into bluff data,
there is no safe opt-in variant of fabricating a skill.

### 2.3 Resources persist in the SAME transaction as the skill they belong to

Add one new transactional store method that the register's decision line names explicitly
(`GAPS_AND_RISKS_REGISTER.md:238`: *"persist resources in the same transaction as the skill"*):

```go
// internal/skill/store.go (new)
func (s *Store) CreateWithResources(ctx context.Context, skill *models.Skill, resources []models.Resource) error {
    return s.pool.WithTx(ctx, func(tx pgx.Tx) error {
        // same skill-insert + dependency-insert + registry-upsert body as Create (store.go:275-333),
        // rewritten to use tx.Exec/tx.QueryRow instead of s.pool.Exec/s.pool.QueryRow
        ...
        for i := range resources {
            resources[i].SkillID = skill.ID
            if resources[i].ID == uuid.Nil { resources[i].ID = uuid.New() }
            if _, err := tx.Exec(ctx, insertResourceSQL, resources[i].ID, resources[i].SkillID,
                resources[i].URL, resources[i].Title, resources[i].ResourceType,
                resources[i].FetchedHash, resources[i].ContentCached, resources[i].LastValidated); err != nil {
                return fmt.Errorf("insert resource %d: %w", i, err)
            }
        }
        if len(resources) > 0 {
            if err := s.recalcCoverage(ctx, tx, skill.ID); err != nil {
                return fmt.Errorf("recalc coverage: %w", err)
            }
        }
        return s.logAudit(ctx, tx, "skill.created_with_resources", &skill.ID,
            map[string]interface{}{"resource_count": len(resources)})
    })
}
```

`Pipeline.Run` (`pipeline.go:393` + the dead `401-403` loop) replaces `p.store.Create(ctx, draft)` +
the tag-only loop with a single `p.store.CreateWithResources(ctx, draft, resources)` call. **Atomicity
guarantee:** if the resource insert fails partway through, the whole transaction — including the skill
row itself — rolls back; `Run` records the failure in `result.Errors` and creates **zero** rows for
that gap, rather than today's possible split state (skill row present, resources silently absent).
This is strictly safer than the pre-fix behaviour, which had no failure path to roll back because it
never attempted the insert at all. Existing non-autoexpand callers of `Store.Create`
(`CreateFromTOML`, the REST/MCP direct-create paths) are **untouched** — `CreateWithResources` is an
additive method, not a signature change to `Create`, so no other call site's behaviour shifts.
**Rejected alternative:** calling `Create` then `BulkAddResources` as two separate transactions (the
existing `resources.go:221-254` method) — rejected because it reintroduces exactly the split-state risk
(skill created, then a resource-insert failure leaves an orphan draft with zero resources, silently
looking "complete" to any reader that only checks skill existence) that the register's "same
transaction" instruction exists to prevent.

### 2.4 What "genuine LLM content" means at the boundary (composes G05 fail-closed jury)

This design does **not** attempt to distinguish "good" LLM content from "bad" LLM content by inspecting
its text — that is explicitly the validation pipeline's job (G05's LLM jury, `internal/validation`),
not auto-expand's. The contract this design enforces is narrower and mechanical: a skill row is
created **only** as the parsed output of a real `LLMClient.Generate` call (§2.1), never as
hand-written template text (§2.2), and it is created with `Status = models.SkillStatusDraft`
unchanged from today (`pipeline.go:232`) — it does **not** advance to `validated`/`active` inside
`autoexpand` at all; that transition is exclusively G05/G11's `runValidationCycle`'s to make, gated on
a real recorded jury verdict (`ApprovedBy ≥ approval_threshold`, per `g11_worker_design.md §2.3`'s
`runValidationCycle` design and G05's fail-closed empty-jury decision,
`GAPS_AND_RISKS_REGISTER.md:94`). G20 removes the *fabrication* at the drafting boundary; G05 (already
decided, not re-derived here) removes the *rubber-stamp* at the jury boundary. Composing the two: a
draft created by this design's fixed `Run` can reach `active` only via a path that (a) came from a real
LLM response and (b) was independently jury-approved — never from a placeholder and never from an
empty-jury auto-pass.

### 2.5 Alternatives rejected (consolidated)

- **Keep `createMinimalDraft` behind a config flag "for demos/offline mode".** Rejected — the register
  itself rejects this class of "graceful degradation" (`GAPS_AND_RISKS_REGISTER.md:238`); a flag that
  can flip a production auto-growth cycle into fabricating content is the exact footgun this design
  closes, and a demo need is served by a fixture skill inserted directly by test/demo tooling, not by
  the production pipeline's code path.
- **Widen the `LLMClient` interface with a `GenerateSkillDraft` method.** Rejected in favour of §2.1's
  approach (orchestration lives at the pipeline layer over the single `Generate` primitive) — smaller
  interface, no duplicated prompt/parse logic per provider.
- **Two-transaction skill-then-resources persistence (reuse existing `BulkAddResources` as-is).**
  Rejected in favour of §2.3's single new `CreateWithResources` — the register is explicit about
  "same transaction," and a two-transaction approach reintroduces a split-state window.

---

## 3. WHY (§11.4.8 external precedent)

**NO external solution found — original work.** This is an internal anti-bluff architecture fix, not a
novel algorithm; the applicable precedent is the codebase's own already-idiomatic Go practice, applied
consistently rather than partially: (a) "accept interfaces, return structs" (Effective Go's standard
guidance) is already the intent of `pipeline.go:52-57`'s `WithLLMClient(client LLMClient)` option — the
bug is that `DraftSkill` re-asserts the concrete type the option was designed to abstract away; §2.1
simply completes that already-adopted pattern. (b) The **local** embedder in this same codebase
(`internal/db/embedding.go:256-262`, cited in `g10_embedding_provider_design.md §3.1`) already
demonstrates the sibling fix for a sibling gap (return-value validation instead of trusting a response
blindly) — no external search needed, the internal precedent is direct. (c) "Fail closed with an
honest, logged skip rather than fabricate a plausible-looking default" is the same discipline G02
(`g02_sandbox_faildesign.md`) and G05 (`GAPS_AND_RISKS_REGISTER.md:89-95`) already apply elsewhere in
this exact codebase — §2.2 is a third instance of one already-established house pattern, not a new one.

---

## 4. RUNTIME SIGNATURE (definition-of-done, §11.4.108)

The fix is DONE only when, on a clean deployment with a real Postgres+pgvector (§11.4.27 no mocks):

| Scenario | Observable on a clean deploy |
|---|---|
| Provider configured, real LLM reachable | Seeding a known gap and running `Pipeline.Run` produces exactly one new `skills` row whose `content` is the LLM's real generated markdown (NOT the string "This skill was auto-generated to fill a gap") **and** every resource the LLM's JSON `resources` array named exists as a row in `resources` with the correct `skill_id`, inserted in the same transaction (a forced mid-transaction resource-insert failure leaves **zero** rows for that gap — skill included). |
| No provider configured (`p.llm == nil`) | Running `Pipeline.Run` against a seeded gap creates **zero** rows in `skills`/`resources`; the log contains `reason="no_provider_configured"`; `ExpansionResult.SkillsSkipped == 1`, `SkillsCreated == 0`. |
| Non-`*OpenAILLM` `LLMClient` configured (a fake/HelixLLM double) | `DraftSkill` succeeds identically to the `*OpenAILLM` case — no `"unsupported LLM client type"` error anywhere, proving R7 pluggability (grep gate: zero `*OpenAILLM` type assertions remain in `pipeline.go`). |
| Full lifecycle | The created `draft` skill advances to `validated`/`active` **only** after `runValidationCycle` (G11) records a real jury verdict with `ApprovedBy ≥ approval_threshold` (composes G05) — never automatically from `autoexpand` alone. |

---

## 5. TEST-CASE COUNT (RED-first per §11.4.115; each PASS cites captured evidence per §11.4.5/§11.4.69)

**Unit (no DB) — 6**
1. `TestDraftSkill_NilLLM_NoPlaceholder` — `p.llm == nil` ⇒ `DraftSkill` returns `(nil, nil,
   ErrNoProviderConfigured)`, never a `*models.Skill`. **[RED on baseline: currently returns a
   fabricated draft via `createMinimalDraft`, `pipeline.go:211-213`]**
2. `TestDraftSkill_AnyLLMClient_NotJustOpenAI` — a table-driven fake `LLMClient` (not `*OpenAILLM`)
   returning a valid JSON draft ⇒ `DraftSkill` succeeds, parses correctly, no type-assertion error.
   **[RED on baseline: hits `"unsupported LLM client type"`, `pipeline.go:217`]**
3. `TestDraftSkill_LLMGenerateError_NoFallbackFabrication` — fake `LLMClient.Generate` returns an error
   ⇒ `DraftSkill` returns that error, never falls back to `createMinimalDraft`. **[RED on baseline:
   falls back and returns a fabricated draft with `err == nil`, `pipeline.go:222-227`]**
4. `TestDraftSkill_MalformedLLMResponse_ErrorsNotFabricates` — fake `LLMClient.Generate` returns
   non-JSON text ⇒ `parseSkillDraft` error propagates, no fallback fabrication.
5. `TestRun_NoProvider_ZeroSkillsCreatedSkippedCounted` — `Pipeline` constructed with `llm=nil`, `Run`
   over a seeded-gap fixture (in-memory store double, unit-level only) ⇒ `SkillsCreated==0`,
   `SkillsSkipped==1`.
6. `createMinimalDraft` **absence** check — a static/compile-level assertion (or a simple `grep`-backed
   test) that the function no longer exists in the package, closing the removal per §11.4.124 (don't
   leave the unsafe shape around after superseding it).

**Integration (real `pgvector:pg16` from `deploy/`, §11.4.27 no mocks) — 4**
7. `TestCreateWithResources_AtomicSuccess` — seed a gap, run the real pipeline against a real (test)
   `LLMClient` fixture returning a valid draft + 2 resources ⇒ both the `skills` row AND both
   `resources` rows exist, correctly linked by `skill_id`, in one committed transaction.
8. `TestCreateWithResources_PartialFailureRollsBackSkill` — inject a resource-insert failure (e.g. a
   resource with a URL exceeding a column constraint, or a forced context-cancel mid-transaction) ⇒
   **zero** rows in `skills` for that draft (the skill row is NOT left orphaned without its
   resources). **[the §11.4.115 RED→GREEN artifact for atomicity]**
9. `TestRun_RealDraft_NoPlaceholderContent` — end-to-end seeded-gap run against a real DB + real
   `LLMClient` fixture ⇒ the created skill's `content` does NOT contain the literal string "This
   skill was auto-generated to fill a gap" (the exact G20 anti-bluff assertion the register + testing
   plan both name, `GAPS_AND_RISKS_REGISTER.md:236`, `testing_infrastructure_plan.md:265`).
10. `TestRun_DraftStaysUnvalidated_UntilRealJuryVerdict` — created draft's `status` remains `draft`
    immediately after `autoExpand.Run` returns (composes G05/G11 — no auto-advance to
    `validated`/`active` from within this pipeline alone).

**Paired §1.1 mutations — 3** (each MUST flip its guard RED then restore GREEN)
- M1: reintroduce `createMinimalDraft`'s call in the `p.llm == nil` branch → test 1 and test 9 (the
  no-placeholder-content assertion) FAIL.
- M2: reintroduce the `openaiLLM, ok := p.llm.(*OpenAILLM)` assertion in `DraftSkill` → test 2 FAILs.
- M3: revert `Run` to the dead `resources[i].SkillID = draft.ID`-only loop (drop the
  `CreateWithResources` call) → test 7 FAILs (no resource rows exist).

**Total: 13 cases** (6 unit + 4 integration + 3 paired mutations).

### 5.1 Reconciliation with `research/testing_infrastructure_plan.md` (§11.4.186)

The plan's G20 row (`testing_infrastructure_plan.md:306`) specifies exactly: *"unit (nil LLM ⇒ no
placeholder persisted; interface pluggability), integration (draft ⇒ resources persisted same tx),
mutation (reintroduce placeholder-persist → anti-bluff test FAILs)"* with evidence *'"no placeholder
row" assertion; persisted-resources row dump'*. This design's cases map onto it with no divergence:
tests 1+9 ↔ "nil LLM ⇒ no placeholder persisted" + the "no placeholder row" evidence; test 2 ↔
"interface pluggability"; tests 7+8 ↔ "draft ⇒ resources persisted same tx" + the "persisted-resources
row dump" evidence; M1/M2/M3 ↔ "reintroduce placeholder-persist → anti-bluff test FAILs" (M1 covers the
named case directly, M2/M3 are the design's necessary extension covering the interface-pluggability and
resource-atomicity halves the plan's one-line mutation cell did not individually enumerate).
**Extension flagged per §11.4.186:** recommend the plan's G20 row (`testing_infrastructure_plan.md:306`)
be extended from one mutation cell to explicitly list all three (M1 placeholder-fabrication, M2
concrete-type-coupling, M3 resource-persistence-atomicity) so the two documents stay in lockstep; no
existing plan content is contradicted, only broadened. The plan's Challenge reference
(`testing_infrastructure_plan.md:245`, CH-G20 "submit a knowingly false skill and assert the
jury/pipeline rejects it") and HelixQA "no placeholder/auto-generated-to-fill-a-gap content"
assertion (`testing_infrastructure_plan.md:265`, `:350`) are inherited unchanged — this design's test 9
is exactly the machine-checkable form of that HelixQA assertion. Harness fit: unit cases use the
sanctioned in-package fake `LLMClient` double (unit-level fakes are permitted per §11.4.27); integration
cases run against the real `deploy/` `pgvector:pg16` with `HELIX_TEST_DATABASE_URL`
SKIP-with-reason when unreachable (`testing_infrastructure_plan.md:88-98`), never a green PASS offline.

---

## 6. HONEST GAPS (§11.4.6) — never guessed

- **Depends on G03 wiring to be exercised at all, not re-done here.** `autoexpand.NewPipeline` is
  still never called from `cmd/server`/`cmd/worker` at `255061b` (§1.1). This design makes the
  pipeline **safe** to wire — the moment G03/G11's `runAutoExpandCycle` calls `p.autoExpand.Run(...)`
  on a real schedule, it will not fabricate. It does not perform that wiring; that is G03/G11's
  scope, cited not duplicated.
- **Depends on a real `LLMClient` provider construction path, which does not yet exist.**
  `config.AutoExpandConfig` (`internal/config/config.go:124-131`) already has `LLMProvider`/`LLMModel`
  string fields, but **no** `NewLLMClientFromConfig(cfg AutoExpandConfig) (LLMClient, error)` factory
  exists anywhere in the tree (unlike the embedding side's `NewEmbedderFromConfig`,
  `g10_embedding_provider_design.md §1.4`/§2.4). Without such a factory, "provider configured" in
  practice still means "an operator manually calls `autoexpand.WithLLMClient(...)`" — there is no
  config-driven path yet. This is a **new, tracked wiring gap** (§11.4.197: flagged so it is not lost),
  symmetric to G10/G11's embedder-construction gap; it is NOT designed in depth here because it is a
  provider-factory concern parallel to (and should reuse the same registry shape as) G10 §2.1/§2.4,
  not a re-derivation of this design's fail-closed/no-fabrication/atomic-persistence contract.
- **`ErrNoProviderConfigured` is a new sentinel; its exact placement (package-level `var` vs a typed
  error) is an implementation detail left open** — either satisfies this design's contract, and the
  choice does not change any test's observable behaviour (`errors.Is` works with either).
- **`Store.Create`'s existing (non-autoexpand) callers remain non-transactional across
  skill+dependency+registry inserts** — this pre-existing characteristic of `Create` (`store.go:275-333`
  uses `s.pool.Exec`/`QueryRow` directly, not `WithTx`) is **not** introduced by this design and is
  **out of scope** for G20; `CreateWithResources` is a new, additive, fully-transactional method used
  only by `autoexpand.Pipeline.Run`, and does not retrofit `Create`'s other call sites. Flagged as an
  honest pre-existing gap, not silently claimed fixed.
- **Live-LLM-required for the "genuine content" half of GREEN.** Tests 7/9 need either a real reachable
  `LLMClient` provider or a high-fidelity fixture double returning realistic JSON; a genuinely live
  cross-provider proof (OpenAI vs a HelixLLM/LLMsVerifier alias, per R7) needs an operator-scheduled
  live acceptance window, exactly the same honest boundary G10 (§5) and G11 (§5) already record for
  their own live-dependency tests — not re-derived here, the same discipline applies.
- **Not in scope / not claimed:** improving the *quality* of LLM-generated content (that is G05's jury
  validation job, not auto-expand's), and choosing which concrete provider becomes the operator's
  default (a product decision, not a code-correctness one). Neither is designed here; both are honest
  deferrals.

*Positive-evidence-only. Every "is/does" claim above is pinned to a file:line read from
`git show 255061b:…`; every forward-looking claim about a not-yet-existing factory/wiring path is
marked as a tracked gap, never asserted as already fixed. No forbidden hedging vocabulary used.*
