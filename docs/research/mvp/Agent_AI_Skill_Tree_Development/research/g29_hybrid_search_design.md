# Design: G29 — Wire Hybrid Vector Search

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

**Date:** 2026-07-17
**Status:** DESIGN
**Scope:** Wire `Store.VectorSearch` into `Store.Search` for hybrid vector+keyword search

---

## 1. Current State

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

### What's Implemented
- `internal/skill/store.go:50-118` — `Store.Search` — ILIKE/trigram only (no vector)
- `internal/skill/store.go:574-609` — `Store.VectorSearch` — pgvector KNN (zero callers)
- `internal/db/embedding.go` — Embedding generation (OpenAI/local)
- `internal/skill/hybrid_search_g29_test.go` — Test file exists

### What's NOT Wired
- `Store.VectorSearch` has zero callers across MCP, REST, and pipeline
- `Store.Search` doc-comment claims "hybrid vector search" but is keyword-only
- No query embedding in the search path

---

## 2. Design

### 2.1 Hybrid Search Algorithm

```go
// store.go — Search (hybrid vector + keyword)
func (s *Store) Search(ctx context.Context, query string, limit int) ([]Skill, error) {
    // 1. Embed the query
    queryEmbedding, err := s.embedder.Embed(ctx, query)
    if err != nil {
        // Fallback to keyword-only on embedding failure
        s.logger.Warn("embedding failed, falling back to keyword", zap.Error(err))
        return s.searchKeyword(ctx, query, limit)
    }

    // 2. Vector search (top K*2 for reranking)
    vectorResults, err := s.vectorSearch(ctx, queryEmbedding, limit*2)
    if err != nil {
        s.logger.Warn("vector search failed, falling back to keyword", zap.Error(err))
        return s.searchKeyword(ctx, query, limit)
    }

    // 3. Keyword search (top K*2 for reranking)
    keywordResults, err := s.searchKeyword(ctx, query, limit*2)
    if err != nil {
        return nil, fmt.Errorf("keyword search: %w", err)
    }

    // 4. Reciprocal Rank Fusion (RRF) merge
    merged := s.rrfMerge(vectorResults, keywordResults, limit)

    return merged, nil
}

// rrfMerge implements Reciprocal Rank Fusion
func (s *Store) rrfMerge(vectorResults, keywordResults []Skill, limit int) []Skill {
    const k = 60 // RRF constant

    scores := make(map[uuid.UUID]float64)
    skillMap := make(map[uuid.UUID]Skill)

    // Vector scores
    for rank, skill := range vectorResults {
        scores[skill.ID] += 1.0 / (float64(k) + float64(rank+1))
        skillMap[skill.ID] = skill
    }

    // Keyword scores
    for rank, skill := range keywordResults {
        scores[skill.ID] += 1.0 / (float64(k) + float64(rank+1))
        skillMap[skill.ID] = skill
    }

    // Sort by RRF score
    type scored struct {
        skill Skill
        score float64
    }
    var sorted []scored
    for id, score := range scores {
        sorted = append(sorted, scored{skillMap[id], score})
    }
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].score > sorted[j].score
    })

    // Return top N
    result := make([]Skill, 0, limit)
    for i := 0; i < limit && i < len(sorted); i++ {
        result = append(result, sorted[i].skill)
    }
    return result
}
```

### 2.2 MCP Integration

```go
// mcp/tools.go — skill_search (update to use hybrid)
func (m *MCPServer) skillSearch(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    query := args["query"].(string)
    limit := int(args["limit"].(float64))

    skills, err := m.store.Search(ctx, query, limit)
    if err != nil {
        return nil, fmt.Errorf("search: %w", err)
    }

    return skills, nil
}
```

### 2.3 REST Integration

```go
// api/search_handler.go — SearchSkills (update to use hybrid)
func (h *SearchHandler) SearchSkills(c *gin.Context) {
    query := c.Query("q")
    limit := 10 // default

    skills, err := h.store.Search(c.Request.Context(), query, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, skills)
}
```

---

## 3. Test Plan

| Test Type | Scope | Priority |
|-----------|-------|----------|
| Unit | RRF merge correctness | P0 |
| Unit | Embedding failure → keyword fallback | P0 |
| Integration | Semantic near-match ranks above keyword-only | P0 |
| Integration | Live pgvector KNN | P1 |
| Mutation | Revert to ILIKE-only → hybrid test FAILs | P0 |
| Performance | Search latency < 200ms | P1 |

---

## 4. Dependencies

| Dependency | Status | Blocking |
|------------|--------|----------|
| G10 (embedding dimension) | DESIGN DONE | Yes (embedding must work) |
| G59 (embedding ingestion) | QUEUED | Yes (skills must have embeddings) |
| pgvector extension | ✅ | No |

---

## 5. Honest Gaps

1. **Embedding latency**: Adding embedding to search path adds ~100ms. May need caching for popular queries.
2. **Fallback behavior**: If embedding fails, search degrades to keyword-only. Must be transparent to user.
3. **RRF tuning**: The `k=60` constant may need tuning based on actual data distribution.
