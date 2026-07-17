# Design: G59 — Wire Embedding Ingestion

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

**Date:** 2026-07-17
**Status:** DESIGN
**Scope:** Wire `db.StoreSkillEmbedding` into skill create/update paths

---

## 1. Current State

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

### What's Implemented
- `internal/db/vector.go:180` — `StoreSkillEmbedding(ctx, pool, skillID, embedding)` — exists but zero callers
- `internal/db/vector.go:194` — `StoreEvidenceEmbedding(ctx, pool, evidenceID, embedding)` — exists but zero callers
- `internal/db/embedding.go` — Embedding generation (OpenAI/local)
- `internal/skill/store.go:647-699` — `Store.Create` — never writes `skills.embedding`

### What's NOT Wired
- `Store.Create` doesn't call `StoreSkillEmbedding`
- `Store.Update` doesn't call `StoreSkillEmbedding`
- `Store.ImportFromTOML` doesn't call `StoreSkillEmbedding`
- No embedding on evidence creation

---

## 2. Design

### 2.1 Create-Time Embedding

```go
// store.go — Create (add embedding)
func (s *Store) Create(ctx context.Context, skill *models.Skill) error {
    // ... existing create logic ...

    // Generate embedding asynchronously
    go func() {
        embedding, err := s.embedder.Embed(ctx, skill.Content)
        if err != nil {
            s.logger.Warn("embedding failed",
                zap.String("skill", skill.Name),
                zap.Error(err),
            )
            return
        }

        if err := db.StoreSkillEmbedding(ctx, s.pool, skill.ID, embedding); err != nil {
            s.logger.Warn("store embedding failed",
                zap.String("skill", skill.Name),
                zap.Error(err),
            )
        }
    }()

    return nil
}
```

### 2.2 Update-Time Embedding

```go
// store.go — Update (re-embed on content change)
func (s *Store) Update(ctx context.Context, skill *models.Skill) error {
    // ... existing update logic ...

    // Re-embed if content changed
    go func() {
        embedding, err := s.embedder.Embed(ctx, skill.Content)
        if err != nil {
            s.logger.Warn("re-embedding failed",
                zap.String("skill", skill.Name),
                zap.Error(err),
            )
            return
        }

        if err := db.StoreSkillEmbedding(ctx, s.pool, skill.ID, embedding); err != nil {
            s.logger.Warn("store embedding failed",
                zap.String("skill", skill.Name),
                zap.Error(err),
            )
        }
    }()

    return nil
}
```

### 2.3 Evidence Embedding

```go
// skill/evidence.go — CreateEvidence (add embedding)
func (s *Store) CreateEvidence(ctx context.Context, evidence *models.Evidence) error {
    // ... existing create logic ...

    // Generate embedding
    go func() {
        embedding, err := s.embedder.Embed(ctx, evidence.CodeSnippet)
        if err != nil {
            s.logger.Warn("evidence embedding failed",
                zap.String("source", evidence.SourceProject),
                zap.Error(err),
            )
            return
        }

        if err := db.StoreEvidenceEmbedding(ctx, s.pool, evidence.ID, embedding); err != nil {
            s.logger.Warn("store evidence embedding failed",
                zap.String("source", evidence.SourceProject),
                zap.Error(err),
            )
        }
    }()

    return nil
}
```

### 2.4 Batch Embedding (for migration/backfill)

```go
// db/embedding.go — BatchEmbedSkills
func BatchEmbedSkills(ctx context.Context, pool *Pool, store *skill.Store, embedder Embedder) error {
    skills, err := store.ListAll(ctx)
    if err != nil {
        return fmt.Errorf("list skills: %w", err)
    }

    for _, skill := range skills {
        // Skip if already embedded
        if skill.Embedding != nil {
            continue
        }

        embedding, err := embedder.Embed(ctx, skill.Content)
        if err != nil {
            continue
        }

        if err := StoreSkillEmbedding(ctx, pool, skill.ID, embedding); err != nil {
            continue
        }
    }

    return nil
}
```

---

## 3. Test Plan

| Test Type | Scope | Priority |
|-----------|-------|----------|
| Unit | Create writes embedding | P0 |
| Unit | Update re-embeds on content change | P0 |
| Integration | Create → search finds by semantic similarity | P0 |
| Integration | Batch embed existing skills | P1 |
| Mutation | Remove embedding write → search test fails | P0 |

---

## 4. Dependencies

| Dependency | Status | Blocking |
|------------|--------|----------|
| G10 (embedding dimension) | DESIGN DONE | Yes |
| Embedding provider configured | ⚠️ | Yes |

---

## 5. Honest Gaps

1. **Async embedding**: Using goroutines for embedding means failures are logged but not returned to caller. Consider a retry queue.
2. **Rate limits**: OpenAI embedding API has rate limits. Batch operations may hit limits.
3. **Cost**: Embedding generation has API cost. Consider caching and incremental updates.
