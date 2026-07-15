package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"go.uber.org/zap"
)

// ---------------------------------------------------------------------------
// Vector search
// ---------------------------------------------------------------------------

// VectorSearchResult holds a single result from vector similarity search.
type VectorSearchResult struct {
	ID    uuid.UUID
	Score float64 // cosine similarity (higher = more similar)
}

// VectorSearch performs cosine-similarity search against the embedding
// column of the given table. Returns up to limit results ordered by
// similarity descending.
//
// The table must have an embedding column of type vector(N) and an
// HNSW index for acceptable performance at scale.
func VectorSearch(
	ctx context.Context,
	pool *Pool,
	table string,
	embedding pgvector.Vector,
	limit int,
) ([]VectorSearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	// Validate table name to prevent injection.
	sanitized := sanitizeTableName(table)
	if sanitized == "" {
		return nil, fmt.Errorf("invalid table name: %q", table)
	}

	query := fmt.Sprintf(
		`SELECT id, 1 - (embedding <=> $1) AS score
		 FROM %s
		 WHERE embedding IS NOT NULL
		 ORDER BY embedding <=> $1
		 LIMIT $2`,
		sanitized,
	)

	rows, err := pool.Query(ctx, query, embedding, limit)
	if err != nil {
		return nil, fmt.Errorf("vector search on %s: %w", sanitized, err)
	}
	defer rows.Close()

	var results []VectorSearchResult
	for rows.Next() {
		var r VectorSearchResult
		if err := rows.Scan(&r.ID, &r.Score); err != nil {
			return nil, fmt.Errorf("scan vector search result: %w", err)
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("vector search rows iteration: %w", err)
	}

	return results, nil
}

// ---------------------------------------------------------------------------
// Hybrid search (RRF)
// ---------------------------------------------------------------------------

// HybridSearch combines vector similarity and full-text keyword search
// using Reciprocal Rank Fusion (RRF). The final results are ordered by
// the RRF score descending.
//
// The keywords string is split on whitespace and used in a tsvector
// match against the name, title, and description columns.
func HybridSearch(
	ctx context.Context,
	pool *Pool,
	table string,
	embedding pgvector.Vector,
	keywords string,
	limit int,
) ([]VectorSearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	sanitized := sanitizeTableName(table)
	if sanitized == "" {
		return nil, fmt.Errorf("invalid table name: %q", table)
	}

	// Build tsquery from keywords.
	tsquery := buildTsQuery(keywords)
	if tsquery == "" {
		// No valid keywords — fall back to pure vector search.
		return VectorSearch(ctx, pool, table, embedding, limit)
	}

	const k = 60.0 // RRF constant; higher = smoother ranking

	query := fmt.Sprintf(
		`WITH
		 vector_ranks AS (
			 SELECT id,
					ROW_NUMBER() OVER (ORDER BY embedding <=> $1) AS rank
			 FROM %s
			 WHERE embedding IS NOT NULL
		 ),
		 text_ranks AS (
			 SELECT id,
					ROW_NUMBER() OVER (
						ORDER BY ts_rank_cd(
							setweight(to_tsvector('english', COALESCE(name, '')), 'A') ||
							setweight(to_tsvector('english', COALESCE(title, '')), 'B') ||
							setweight(to_tsvector('english', COALESCE(description, '')), 'C'),
							to_tsquery('english', $3)
						) DESC
					) AS rank
			 FROM %s
			 WHERE to_tsvector('english',
					COALESCE(name, '') || ' ' ||
					COALESCE(title, '') || ' ' ||
					COALESCE(description, '')) @@ to_tsquery('english', $3)
		 )
		 SELECT COALESCE(v.id, t.id) AS id,
				COALESCE(1.0 / (%[3]f + v.rank), 0.0) +
				COALESCE(1.0 / (%[3]f + t.rank), 0.0) AS score
		 FROM vector_ranks v
		 FULL OUTER JOIN text_ranks t ON v.id = t.id
		 ORDER BY score DESC
		 LIMIT $2`,
		sanitized,
		sanitized,
		k,
	)

	rows, err := pool.Query(ctx, query, embedding, limit, tsquery)
	if err != nil {
		return nil, fmt.Errorf("hybrid search on %s: %w", sanitized, err)
	}
	defer rows.Close()

	var results []VectorSearchResult
	for rows.Next() {
		var r VectorSearchResult
		if err := rows.Scan(&r.ID, &r.Score); err != nil {
			return nil, fmt.Errorf("scan hybrid search result: %w", err)
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("hybrid search rows iteration: %w", err)
	}

	return results, nil
}

// ---------------------------------------------------------------------------
// Embedding storage helpers
// ---------------------------------------------------------------------------

// StoreEmbedding updates the embedding vector for a specific skill.
func StoreSkillEmbedding(
	ctx context.Context,
	pool *Pool,
	skillID uuid.UUID,
	embedding pgvector.Vector,
) error {
	const query = `UPDATE skills SET embedding = $1 WHERE id = $2`
	if _, err := pool.Exec(ctx, query, embedding, skillID); err != nil {
		return fmt.Errorf("store embedding for skill %s: %w", skillID, err)
	}
	return nil
}

// StoreEvidenceEmbedding updates the embedding vector for a specific evidence record.
func StoreEvidenceEmbedding(
	ctx context.Context,
	pool *Pool,
	evidenceID uuid.UUID,
	embedding pgvector.Vector,
) error {
	const query = `UPDATE evidences SET embedding = $1 WHERE id = $2`
	if _, err := pool.Exec(ctx, query, embedding, evidenceID); err != nil {
		return fmt.Errorf("store embedding for evidence %s: %w", evidenceID, err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Find similar skills by text (convenience)
// ---------------------------------------------------------------------------

// FindSimilarSkills performs vector search against the skills table and
// returns the matching skill IDs with scores.
func FindSimilarSkills(
	ctx context.Context,
	pool *Pool,
	embedding pgvector.Vector,
	limit int,
) ([]VectorSearchResult, error) {
	return VectorSearch(ctx, pool, "skills", embedding, limit)
}

// FindSimilarEvidences performs vector search against the evidences table.
func FindSimilarEvidences(
	ctx context.Context,
	pool *Pool,
	embedding pgvector.Vector,
	limit int,
) ([]VectorSearchResult, error) {
	return VectorSearch(ctx, pool, "evidences", embedding, limit)
}

// ---------------------------------------------------------------------------
// Vector search with filters
// ---------------------------------------------------------------------------

// VectorSearchFiltered performs vector search with an additional status filter.
func VectorSearchFiltered(
	ctx context.Context,
	pool *Pool,
	table string,
	embedding pgvector.Vector,
	status string,
	limit int,
) ([]VectorSearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	sanitized := sanitizeTableName(table)
	if sanitized == "" {
		return nil, fmt.Errorf("invalid table name: %q", table)
	}

	query := fmt.Sprintf(
		`SELECT id, 1 - (embedding <=> $1) AS score
		 FROM %s
		 WHERE embedding IS NOT NULL AND status = $3
		 ORDER BY embedding <=> $1
		 LIMIT $2`,
		sanitized,
	)

	rows, err := pool.Query(ctx, query, embedding, limit, status)
	if err != nil {
		return nil, fmt.Errorf("filtered vector search on %s: %w", sanitized, err)
	}
	defer rows.Close()

	var results []VectorSearchResult
	for rows.Next() {
		var r VectorSearchResult
		if err := rows.Scan(&r.ID, &r.Score); err != nil {
			return nil, fmt.Errorf("scan filtered vector search result: %w", err)
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("filtered vector search rows iteration: %w", err)
	}

	return results, nil
}

// ---------------------------------------------------------------------------
// Utilities
// ---------------------------------------------------------------------------

// sanitizeTableName ensures the table name contains only alphanumeric
// characters and underscores to prevent SQL injection.
func sanitizeTableName(name string) string {
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// buildTsQuery converts a space-separated keyword string into a PostgreSQL
// tsquery expression with AND semantics.
func buildTsQuery(keywords string) string {
	words := strings.Fields(keywords)
	if len(words) == 0 {
		return ""
	}
	// Escape each word and join with & (AND).
	escaped := make([]string, 0, len(words))
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		// Remove any characters that would break tsquery syntax.
		w = sanitizeTsQueryWord(w)
		if w != "" {
			escaped = append(escaped, w+":*") // prefix matching
		}
	}
	if len(escaped) == 0 {
		return ""
	}
	return strings.Join(escaped, " & ")
}

// sanitizeTsQueryWord removes characters that would break a tsquery.
func sanitizeTsQueryWord(w string) string {
	var b strings.Builder
	for _, r := range w {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '_', r == '-':
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// Vector stats / diagnostics
// ---------------------------------------------------------------------------

// VectorIndexStats returns the number of indexed vectors and index size
// for the given table's embedding column.
func VectorIndexStats(
	ctx context.Context,
	pool *Pool,
	table string,
) (indexedCount int64, indexSizeBytes int64, err error) {
	sanitized := sanitizeTableName(table)
	if sanitized == "" {
		return 0, 0, fmt.Errorf("invalid table name: %q", table)
	}

	// Count non-null embeddings.
	countQuery := fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE embedding IS NOT NULL`, sanitized)
	if err := pool.QueryRow(ctx, countQuery).Scan(&indexedCount); err != nil {
		return 0, 0, fmt.Errorf("count vectors in %s: %w", sanitized, err)
	}

	// Get index size for the HNSW index (best-effort; may not exist).
	indexName := "idx_" + sanitized + "_embedding"
	var sizeMB float64
	err = pool.QueryRow(ctx,
		`SELECT pg_size_pretty(pg_relation_size($1)), pg_relation_size($1)`,
		indexName,
	).Scan(new(string), &indexSizeBytes)
	if err != nil {
		// Index may not exist yet; that's OK.
		indexSizeBytes = 0
	}
	_ = sizeMB

	return indexedCount, indexSizeBytes, nil
}

// WaitForVectorIndexReady polls until the pgvector HNSW index for the
// given table is ready (no invalid index entries), or until the timeout.
func WaitForVectorIndexReady(
	ctx context.Context,
	pool *Pool,
	table string,
	timeout time.Duration,
) error {
	sanitized := sanitizeTableName(table)
	if sanitized == "" {
		return fmt.Errorf("invalid table name: %q", table)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	indexName := "idx_" + sanitized + "_embedding"

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for index %s to be ready: %w", indexName, ctx.Err())
		case <-ticker.C:
			var isValid bool
			err := pool.QueryRow(ctx,
				`SELECT indisvalid FROM pg_index WHERE indexrelname = $1`,
				indexName,
			).Scan(&isValid)
			if err != nil {
				// Index may not exist yet; keep waiting.
				zap.L().Debug("vector index not found yet", zap.String("index", indexName))
				continue
			}
			if isValid {
				return nil
			}
			zap.L().Debug("vector index still building", zap.String("index", indexName))
		}
	}
}
