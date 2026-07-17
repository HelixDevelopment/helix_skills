# Design: G03 — Wire Validation Pipeline into Worker

**Date:** 2026-07-17
**Status:** DESIGN
**Scope:** Wire `internal/validation.Pipeline` into the worker's `handleValidate` and `runValidationCycle`

---

## 1. Current State

### What's Implemented
- `internal/validation/pipeline.go` — Full validation pipeline with:
  - Stage 1: Resource verification (G21 SSRF guard landed)
  - Stage 2: LLM jury (G05 fail-closed on empty jury)
  - Stage 3: Cross-reference check
  - Stage 4: Sandbox execution (G02 — host execution deleted, StaticValidator only)
- `internal/validation/pipeline.go` — `Validate(ctx, skill) (Verdict, error)` method

### What's NOT Wired
- `internal/worker/runner.go:589-601` — `handleValidate` is a stub (returns `Success:true` with no work)
- `internal/worker/runner.go:735-759` — `runValidationCycle` only logs (no dispatch to pipeline)
- `cmd/server/main.go` — No validation on skill create/update
- `internal/mcp/tools.go` — No validation on `skill_create`

---

## 2. Design

### 2.1 Worker Wiring

```go
// runner.go — handleValidate (replaces stub)
func (r *Runner) handleValidate(ctx context.Context, job *Job) error {
    var payload ValidateJobPayload
    if err := json.Unmarshal(job.Payload, &payload); err != nil {
        return fmt.Errorf("unmarshal validate payload: %w", err)
    }

    skill, err := r.store.GetByName(ctx, payload.SkillName)
    if err != nil {
        return fmt.Errorf("get skill %q: %w", payload.SkillName, err)
    }

    verdict, err := r.validator.Validate(ctx, skill)
    if err != nil {
        return fmt.Errorf("validate skill %q: %w", payload.SkillName, err)
    }

    // Update skill status based on verdict
    if verdict.Consensus {
        skill.Status = models.SkillStatusValidated
    } else {
        skill.Status = models.SkillStatusDraft
        // Record rejection reason in audit log
    }

    if err := r.store.Update(ctx, skill); err != nil {
        return fmt.Errorf("update skill status: %w", err)
    }

    // Record verdict in audit log
    r.auditLog(ctx, "validation_complete", skill.ID, map[string]interface{}{
        "consensus": verdict.Consensus,
        "approvals": verdict.Approvals,
        "rejections": verdict.Rejections,
    })

    return nil
}
```

### 2.2 Cycle Wiring

```go
// runner.go — runValidationCycle (replaces log-only)
func (r *Runner) runValidationCycle(ctx context.Context) error {
    // Get all draft skills
    drafts, err := r.store.ListByStatus(ctx, models.SkillStatusDraft)
    if err != nil {
        return fmt.Errorf("list draft skills: %w", err)
    }

    for _, skill := range drafts {
        // Skip if already has a pending validation job
        hasPending, err := r.hasPendingJob(ctx, JobTypeValidate, skill.Name)
        if err != nil {
            r.logger.Warn("check pending job", zap.Error(err))
            continue
        }
        if hasPending {
            continue
        }

        // Enqueue validation job
        payload := ValidateJobPayload{SkillName: skill.Name}
        if err := r.enqueueJob(ctx, JobTypeValidate, payload); err != nil {
            r.logger.Warn("enqueue validate job",
                zap.String("skill", skill.Name),
                zap.Error(err),
            )
        }
    }

    return nil
}
```

### 2.3 Create-Time Validation

```go
// api/skills_handler.go — CreateSkill (add validation gate)
func (h *SkillsHandler) CreateSkill(c *gin.Context) {
    // ... parse request ...

    // Create skill as draft
    skill.Status = models.SkillStatusDraft
    if err := h.store.Create(c.Request.Context(), skill); err != nil {
        // ...
    }

    // Enqueue validation job (async)
    payload := ValidateJobPayload{SkillName: skill.Name}
    if err := h.worker.EnqueueJob(c.Request.Context(), JobTypeValidate, payload); err != nil {
        h.logger.Warn("enqueue validate job", zap.Error(err))
    }

    c.JSON(http.StatusCreated, skill)
}
```

---

## 3. Invariant

**No skill reaches `validated`/`active` status without a recorded jury verdict.**

```go
// store.go — Status transition guard
func (s *Store) UpdateStatus(ctx context.Context, id uuid.UUID, newStatus models.SkillStatus) error {
    // Verify validation verdict exists for validated/active
    if newStatus == models.SkillStatusValidated || newStatus == models.SkillStatusActive {
        hasVerdict, err := s.hasValidationVerdict(ctx, id)
        if err != nil {
            return err
        }
        if !hasVerdict {
            return fmt.Errorf("cannot transition to %s without validation verdict", newStatus)
        }
    }

    // ... update status ...
}
```

---

## 4. Test Plan

| Test Type | Scope | Priority |
|-----------|-------|----------|
| Unit | handleValidate dispatches to pipeline | P0 |
| Unit | runValidationCycle enqueues jobs for drafts | P0 |
| Unit | Status transition guard blocks unvalidated | P0 |
| Integration | Draft → validate → validated (real DB) | P0 |
| Integration | Draft → validate → rejected (real DB) | P0 |
| E2E | Create skill → auto-validate → status update | P1 |
| Mutation | Remove validation → status test fails | P0 |

---

## 5. Dependencies

| Dependency | Status | Blocking |
|------------|--------|----------|
| G05 (jury fail-closed) | FIXED | No |
| G21 (SSRF guard) | FIXED | No |
| G02 (sandbox deletion) | FIXED | No |
| Worker panic firewall (G11) | FIXED | No |

---

## 6. Honest Gaps

1. **LLM jury**: Requires configured LLM providers. Without providers, jury auto-approves (G05 fixed this to fail-closed).
2. **StaticValidator**: Currently only parses code blocks, doesn't execute them. Sandbox execution is deleted (G02).
3. **Performance**: Validation is async, so skills are created as `draft` and validated later. UI must handle this state.
