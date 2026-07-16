// Package skill provides CRUD and search operations for skills in the knowledge graph.
package skill

import (
	"context"
	// dbsql aliased (not "sql"): nearly every function in this file
	// declares a LOCAL variable literally named `sql` for its query string
	// (e.g. GetByName's `sql := \`SELECT ...\``), which would shadow an
	// unaliased `database/sql` package import within that exact scope --
	// silently making the import unusable (and unused) anywhere it is
	// actually needed.
	dbsql "database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/helixdevelopment/skill-system/internal/db"
	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	"go.uber.org/zap"
)

// Sentinel errors returned by the skill store and graph operations. Callers
// should compare against these with errors.Is rather than matching strings.
var (
	// ErrSkillNotFound indicates the requested skill does not exist.
	ErrSkillNotFound = errors.New("skill not found")
	// ErrSkillExists indicates a skill with the same unique name already exists.
	ErrSkillExists = errors.New("skill already exists")
	// ErrInvalidSkill indicates a skill failed structural or semantic validation.
	ErrInvalidSkill = errors.New("invalid skill")
	// ErrDependencyNotFound indicates a referenced dependency skill does not exist.
	ErrDependencyNotFound = errors.New("dependency skill not found")
	// ErrCycleDetected indicates an operation would introduce a dependency cycle.
	ErrCycleDetected = errors.New("dependency cycle detected")
	// ErrPartOfUnsupported indicates a TOML import declared a non-empty
	// `part_of` alias, which is NOT yet wired. Unlike depends_on/prerequisite
	// (a SAME-direction fold to `requires` on the skill being imported),
	// part_of is the child declaring its parent — its documented semantics are
	// an INVERTED, CROSS-SKILL edge (parent→child `composes`, per
	// research/skill_granularity_and_composition.md §4.1 alias table). Wiring
	// it correctly means mutating a DIFFERENT skill's edge set (the parent's)
	// inside the child's import transaction and re-running the parent's cycle
	// check — a materially different operation than every other alias, and one
	// that is not yet designed/tested. Per §11.4.6/§11.4.124, rather than ship
	// a half-inversion OR silently drop the edge (the exact class G07 exists to
	// close, research/g06_g07_skilltree_dag_design.md §2.3(4)), ImportFromTOML
	// HARD-ERRORS on a non-empty part_of. The full part_of→inverted-composes
	// inversion is a tracked G07 follow-up (see the deferral note at the
	// part_of guard in import_export.go).
	ErrPartOfUnsupported = errors.New("part_of dependency alias not yet supported")
)

// Store provides data access for skills and related entities.
type Store struct {
	pool *db.Pool
	// embedder is the OPTIONAL query-side embedder that upgrades Search from
	// keyword-only to a genuine hybrid (vector KNN + trigram) search (§G29).
	// It is nil by default (NewStore leaves it unset), so every construction
	// path that does not opt in keeps the historical keyword-only behaviour
	// and never makes an embedding call. It is wired in from configuration at
	// MCP-server construction (internal/mcp/server.go) so the live skill_search,
	// REST /search, and pipeline-dedup paths — all sharing this one Store —
	// become hybrid the moment an embedding provider is configured.
	embedder db.Embedder
	// lastEmbedWarnUnixNano is the UnixNano timestamp of the most recent
	// "query embedding failed" warning (§G29 finding 3), throttled via
	// warnEmbeddingDegraded so a sustained embedding-provider outage logs a
	// steady drumbeat of evidence instead of either flooding the log once per
	// query or going completely silent. Accessed only via the atomic type's
	// own methods -- safe for the concurrent Search callers this Store serves.
	lastEmbedWarnUnixNano atomic.Int64
	// logger is the injected sink for Store's own diagnostics (currently just
	// warnEmbeddingDegraded). Re-review remediation (MAJOR finding, post-G29):
	// warnEmbeddingDegraded originally called the package-level zap.L(), but
	// this codebase's construction path (cmd/server/main.go builds a concrete
	// *zap.Logger via newLogger and threads it explicitly into api.New /
	// mcp.NewMCPServer / validation.NewPipeline -- see internal/api/server.go
	// and internal/mcp/server.go) never calls zap.ReplaceGlobals, so
	// zap.L() resolves to zap's built-in no-op default in every deployed
	// binary. The warning was therefore emitted to a sink nothing reads --
	// dead at the runtime layer despite a green test suite (the pre-fix test
	// only observed the warning by installing its OWN process-global
	// zap.ReplaceGlobals override, which production code never does).
	// Defaults to zap.NewNop() (see NewStore) so every construction path that
	// never calls WithLogger -- every existing test helper, cmd/worker's Store
	// (which never calls Search) -- keeps emitting nothing and never nil-panics,
	// exactly mirroring the historical behaviour for those paths. Wired to a
	// real logger only by internal/mcp.NewMCPServer, mirroring the existing
	// WithEmbedder fluent-option convention immediately below.
	logger *zap.Logger
}

// NewStore creates a new skill store. logger defaults to zap.NewNop() (see
// WithLogger) so a Store used without WithLogger -- every existing test
// helper and any future construction path that does not opt in -- behaves
// exactly as before this field was added: Store's internal diagnostics are
// silently discarded, never a nil-pointer panic.
func NewStore(pool *db.Pool) *Store {
	return &Store{pool: pool, logger: zap.NewNop()}
}

// WithEmbedder configures the query-side embedder that turns Search into a
// genuine hybrid (vector KNN + trigram) search and returns the receiver for
// fluent wiring (§G29). Passing a nil embedder resets Search to keyword-only.
// This is an explicit opt-in: callers that never invoke it (every current test
// helper, and any deployment with no embedding provider configured) keep the
// keyword-only path and never issue an embedding request.
//
// Concurrency contract (Fable code-review remediation, finding 6b): WithEmbedder
// is a plain, unsynchronized field write -- it is SAFE ONLY as a one-time
// construction-time wire-up that happens-before any concurrent Search call, NOT
// as a live runtime reconfiguration switch. The sole production caller,
// internal/mcp.NewMCPServer, relies on exactly this: it calls WithEmbedder on
// the shared *Store it was handed, synchronously, before returning the
// *MCPServer to its caller -- and every transport (stdio/HTTP/ACP) that could
// concurrently invoke Search is started strictly AFTER NewMCPServer returns
// (cmd/server wires the Store, then NewMCPServer, then RegisterTools/
// ListenAndServe/RunStdio/RunACP). There is a SINGLE construction call over the
// Store's lifetime; calling WithEmbedder again after any transport has started
// serving requests is a data race with concurrent Search readers of s.embedder
// and is NOT supported.
func (s *Store) WithEmbedder(e db.Embedder) *Store {
	s.embedder = e
	return s
}

// WithLogger wires the real application logger into Store so its own
// diagnostics (currently just warnEmbeddingDegraded) reach a real sink at
// runtime instead of the package-level zap.L() no-op default (re-review
// remediation, MAJOR finding; see the logger field's doc comment). Mirrors
// the WithEmbedder fluent-option convention above; returns the receiver for
// fluent wiring. A nil logger is a no-op (the field keeps whatever it already
// had -- the zap.NewNop() default from NewStore unless WithLogger was already
// called) rather than falling back to zap.L(), which would silently
// reintroduce the exact dead-sink class this method exists to close.
//
// Concurrency contract: identical to WithEmbedder -- a plain, unsynchronized
// field write, safe ONLY as a one-time construction-time wire-up that
// happens-before any concurrent Search call. internal/mcp.NewMCPServer calls
// WithLogger synchronously, before returning the *MCPServer, alongside its
// existing WithEmbedder call.
func (s *Store) WithLogger(logger *zap.Logger) *Store {
	if logger != nil {
		s.logger = logger
	}
	return s
}

// Pool returns the underlying database pool for operations that need
// direct database access (e.g., audit logging from other packages).
func (s *Store) Pool() *db.Pool {
	return s.pool
}

// rrfK is the Reciprocal Rank Fusion smoothing constant (the field-standard
// default). It damps how quickly a result's contribution decays with rank, so
// no single high-ranked hit in one list can dominate the fused ordering.
const rrfK = 60

// Per-list RRF weights (§G29). Vector recall is weighted marginally above the
// lexical list because the semantic path is the capability G29 restores (R2/R13
// make semantic retrieval core) and this deterministically breaks a rank tie in
// favour of a semantically-near match over an equal-rank purely-lexical one. A
// skill that matches BOTH lexically AND semantically still dominates — it
// accrues both weighted contributions — so exact-name lookups (which land at
// the top of both lists) are unaffected by the slight tilt.
const (
	vectorRRFWeight  = 1.0
	trigramRRFWeight = 0.9
)

// embedDegradeWarnInterval throttles the "query embedding failed, degrading to
// keyword-only" warning (§G29 finding 3) to at most once per this interval per
// Store. A failing embedding provider can degrade EVERY hybrid-search query in
// a hot path; logging every single occurrence would flood the log at exactly
// the moment an operator most needs a readable signal, while suppressing it
// entirely (the pre-fix behaviour) makes an ongoing outage invisible. Once per
// interval keeps a sustained failure OBSERVABLE without spamming.
const embedDegradeWarnInterval = 30 * time.Second

// Search performs a hybrid search over the skill graph.
//
// When a query-side embedder is configured (WithEmbedder) it runs two candidate
// retrievals — a pgvector cosine-KNN over skills.embedding (semantic recall,
// via VectorSearch) and a pg_trgm/ILIKE keyword match (lexical precision) — and
// fuses them with weighted Reciprocal Rank Fusion. Because the fusion is by RANK
// (not by raw score) the incomparable score scales of the two paths cannot
// distort each other, and a semantically-near skill whose text does NOT contain
// the query as a substring can both surface and outrank a purely-lexical match —
// the recall the keyword-only path structurally cannot deliver.
//
// When no embedder is configured, OR the query embedding fails (e.g. an
// embedding provider is temporarily unreachable), Search transparently degrades
// to the keyword-only path rather than returning an error, so keyword search
// keeps working everywhere. A genuine failure of the vector KNN query itself is
// NOT masked — it is returned — because that signals a real misconfiguration
// (e.g. an embedding/column dimension mismatch) that must surface (§11.4.6).
//
// The returned SearchResult.Score is the pg_trgm similarity for the keyword-only
// path, and the fused RRF relevance score for the hybrid path.
func (s *Store) Search(ctx context.Context, query string, limit int) ([]models.SearchResult, error) {
	trigram, err := s.textSearch(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	// No embedder wired ⇒ historical keyword-only behaviour, byte-for-byte.
	if s.embedder == nil {
		return trigram, nil
	}

	vecs, embErr := s.embedder.Embed(ctx, []string{query})
	if embErr != nil {
		// The embedding provider is unavailable for this query: keep search
		// working on the lexical path instead of failing the request -- but,
		// unlike before (§G29 finding 3), make the degradation OBSERVABLE
		// rather than silently swallowing it with zero telemetry.
		s.warnEmbeddingDegraded(embErr)
		return trigram, nil
	}
	if len(vecs) == 0 || len(vecs[0]) == 0 {
		// The provider returned no error but also no usable vector (e.g. an
		// empty batch): degrade quietly, this is not a failure signal worth
		// warning about.
		return trigram, nil
	}

	vector, err := s.VectorSearch(ctx, vecs[0], limit)
	if err != nil {
		// A KNN query error is a real internal fault (e.g. dimension mismatch),
		// not an expected degradation — surface it rather than silently hide it.
		return nil, fmt.Errorf("hybrid search vector leg: %w", err)
	}

	return fuseSearchResults(vector, trigram, limit), nil
}

// warnEmbeddingDegraded logs, at most once per embedDegradeWarnInterval, that
// a query-side embedding call failed and Search degraded to keyword-only
// (§G29 finding 3). Pre-fix, this failure was swallowed with zero telemetry --
// an operator watching logs during a real embedding-provider outage would see
// nothing but degraded search relevance, with no signal pointing at the cause.
// Throttling (rather than logging every occurrence) keeps a sustained outage
// visible without flooding the log once per search request.
//
// Logs via s.logger (injected by WithLogger), NOT the package-level zap.L()
// (re-review remediation, MAJOR finding): this codebase never calls
// zap.ReplaceGlobals, so zap.L() is zap's no-op default in every deployed
// binary -- the pre-fix version of this call was dead at the runtime layer.
// s.logger defaults to zap.NewNop() (see NewStore/WithLogger) so a Store never
// wired with a real logger keeps the historical zero-output behaviour.
func (s *Store) warnEmbeddingDegraded(err error) {
	now := time.Now().UnixNano()
	last := s.lastEmbedWarnUnixNano.Load()
	if now-last < int64(embedDegradeWarnInterval) {
		return
	}
	if !s.lastEmbedWarnUnixNano.CompareAndSwap(last, now) {
		return // another goroutine just logged; avoid a duplicate burst
	}
	s.logger.Warn("hybrid search: query embedding failed, degrading to keyword-only search (§G29)", zap.Error(err))
}

// textSearch is the keyword-only leg of Search: a pg_trgm similarity/ILIKE match
// with a broaden-on-empty ILIKE fallback. Its behaviour is identical to the
// pre-§G29 Search so every keyword-only caller is unchanged.
func (s *Store) textSearch(ctx context.Context, query string, limit int) ([]models.SearchResult, error) {
	// COALESCE(s.description, '') inside the similarity() concatenation
	// (finding 7, extended beyond the reviewer's literal fix text after live
	// testing surfaced the FULL defect, §11.4.194): PostgreSQL's `||` yields
	// NULL when EITHER operand is NULL, so with a nullable description
	// (migrations/001_initial.up.sql) a NULL-description row's ENTIRE
	// concatenation -- and therefore its similarity()/score -- becomes NULL,
	// independent of whatever the raw `s.description` SELECT column scans as.
	// Scanning that NULL score into the non-nullable float64 Score field
	// panics/errors "cannot scan NULL into *float64" BEFORE the description
	// column is ever reached -- proven against a real NULL-description row
	// (§11.4.199 exact reproduction), which is why the NullString fix on
	// Description ALONE (below) does not by itself prevent this query from
	// erroring on such a row.
	sql := `
		SELECT s.id, s.name, s.version, s.title, s.description, s.content,
		       s.metadata, s.status, s.kind, s.created_at, s.updated_at,
		       similarity(s.name || ' ' || s.title || ' ' || COALESCE(s.description, ''), $1) as score
		FROM skills s
		WHERE s.name % $1 OR s.title % $1 OR s.description ILIKE '%' || $1 || '%'
		ORDER BY score DESC, s.name
		LIMIT $2
	`
	rows, err := s.pool.Query(ctx, sql, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search skills: %w", err)
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var r models.SearchResult
		var metadata []byte
		// F2-precedent NullString scan (Fable code-review remediation, finding
		// 7): description is a nullable TEXT column (migrations/001_initial.up.sql)
		// that Store.Create's INSERT always sets to a Go zero-value "" (never
		// SQL NULL) -- but a row written by direct SQL (bypassing Create;
		// e.g. a test fixture, a migration seed, or a future bulk-loader) can
		// leave it genuinely NULL, and scanning a NULL directly into
		// models.Skill's plain (non-nullable) Description string then panics
		// with "cannot scan NULL into *string", exactly as GetByName's
		// fetched_hash/content_cached fix (store.go's GetByName) already
		// established for resources. Scanning through sql.NullString and
		// defaulting to "" on NULL matches that same precedent here.
		var description dbsql.NullString
		err := rows.Scan(
			&r.Skill.ID, &r.Skill.Name, &r.Skill.Version, &r.Skill.Title,
			&description, &r.Skill.Content, &metadata,
			&r.Skill.Status, &r.Skill.Kind, &r.Skill.CreatedAt, &r.Skill.UpdatedAt,
			&r.Score,
		)
		if err != nil {
			return nil, fmt.Errorf("scan search result: %w", err)
		}
		r.Skill.Description = description.String
		r.Skill.Metadata = metadata
		results = append(results, r)
	}
	// F4 (Fable code-review remediation, finding 4): rows.Next() returning
	// false means EITHER the result set is exhausted OR row iteration was cut
	// short by a driver/connection-level error -- the two are indistinguishable
	// from the loop alone. Without checking rows.Err(), a mid-stream failure
	// silently truncates `results` to whatever was scanned so far and this
	// function returns that partial slice with a nil error, masking the
	// failure as "these are all the matches" instead of surfacing it. Mirrors
	// the sibling package's own convention (internal/db/vector.go).
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("search skills rows: %w", err)
	}

	if len(results) == 0 {
		// Fallback: return all skills if no similarity match
		sql = `
			SELECT s.id, s.name, s.version, s.title, s.description, s.content,
			       s.metadata, s.status, s.kind, s.created_at, s.updated_at, 0.0 as score
			FROM skills s
			WHERE s.name ILIKE '%' || $1 || '%' OR s.title ILIKE '%' || $1 || '%'
			ORDER BY s.name
			LIMIT $2
		`
		rows, err := s.pool.Query(ctx, sql, query, limit)
		if err != nil {
			return nil, fmt.Errorf("fallback search: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var r models.SearchResult
			var metadata []byte
			var description dbsql.NullString // F2/finding 7 precedent, see above
			if err := rows.Scan(
				&r.Skill.ID, &r.Skill.Name, &r.Skill.Version, &r.Skill.Title,
				&description, &r.Skill.Content, &metadata,
				&r.Skill.Status, &r.Skill.Kind, &r.Skill.CreatedAt, &r.Skill.UpdatedAt,
				&r.Score,
			); err != nil {
				return nil, fmt.Errorf("scan fallback result: %w", err)
			}
			r.Skill.Description = description.String
			r.Skill.Metadata = metadata
			results = append(results, r)
		}
		// F4/finding 4, fallback leg: see the primary leg's rows.Err() note above.
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("fallback search rows: %w", err)
		}
	}

	return results, nil
}

// fuseSearchResults merges the vector-KNN and keyword candidate lists into one
// ranked list with weighted Reciprocal Rank Fusion (§G29). Each list is assumed
// to already be in descending-relevance order (index 0 = best). A skill present
// in both lists accrues both weighted 1/(rrfK+rank+1) contributions; ties are
// broken deterministically by skill name so the ordering is stable across runs.
// The fused relevance replaces each result's per-leg Score. limit ≤ 0 means no
// cap.
func fuseSearchResults(vector, trigram []models.SearchResult, limit int) []models.SearchResult {
	type agg struct {
		res   models.SearchResult
		score float64
	}
	byID := make(map[uuid.UUID]*agg)
	order := make([]uuid.UUID, 0, len(vector)+len(trigram))

	accumulate := func(list []models.SearchResult, weight float64) {
		for rank, r := range list {
			a, ok := byID[r.Skill.ID]
			if !ok {
				cp := r
				a = &agg{res: cp}
				byID[r.Skill.ID] = a
				order = append(order, r.Skill.ID)
			}
			a.score += weight / float64(rrfK+rank+1)
		}
	}
	accumulate(vector, vectorRRFWeight)
	accumulate(trigram, trigramRRFWeight)

	merged := make([]models.SearchResult, 0, len(order))
	for _, id := range order {
		a := byID[id]
		a.res.Score = a.score
		merged = append(merged, a.res)
	}
	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].Score != merged[j].Score {
			return merged[i].Score > merged[j].Score
		}
		return merged[i].Skill.Name < merged[j].Skill.Name
	})

	if limit > 0 && len(merged) > limit {
		merged = merged[:limit]
	}
	return merged
}

// GetByName retrieves a complete skill by its unique name.
func (s *Store) GetByName(ctx context.Context, name string) (*models.Skill, error) {
	sql := `
		SELECT s.id, s.name, s.version, s.title, s.description, s.content,
		       s.metadata, s.status, s.kind, s.created_at, s.updated_at
		FROM skills s
		WHERE s.name = $1
	`
	var skill models.Skill
	var metadata []byte
	err := s.pool.QueryRow(ctx, sql, name).Scan(
		&skill.ID, &skill.Name, &skill.Version, &skill.Title,
		&skill.Description, &skill.Content, &metadata,
		&skill.Status, &skill.Kind, &skill.CreatedAt, &skill.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Wrap the ErrSkillNotFound sentinel (§11.4.6/§11.4.102 forensic
			// finding, pre-existing defect discovered while implementing the
			// P1.T1 M10 seed-import test): ImportFromTOML's existing-skill
			// guard (import_export.go) checks errors.Is(err, ErrSkillNotFound)
			// to distinguish "not found, OK to create" from a real DB error.
			// The previous plain fmt.Errorf (no %w) never satisfied that
			// check, so EVERY ImportFromTOML call for a brand-new skill name
			// took the "real error" branch and aborted before the INSERT ever
			// ran -- the message text is unchanged, only the wrapping is
			// fixed. Blast-radius audited (N3 correction, Fable code-review
			// remediation): every non-test GetByName call site in the repo,
			// 9 in total across 6 files --
			// import_export.go:40 (ImportFromTOML existing-skill guard,
			// the one that actually needs errors.Is), import_export.go:235
			// (ExportToTOML), graph.go:192 (GetDependencyTree),
			// store.go:201 (GetTree, this same package), pipeline.go:174,246,424
			// (internal/autoexpand, 3 call sites), mcp/tools.go:117
			// (skill_get tool), and main.go:218 (REST skill-by-name route) --
			// all treat a GetByName error either generically (fmt.Errorf-wrap,
			// HTTP 404/500, or a boolean "not found") or via errors.Is on this
			// exact sentinel; none depended on the old unwrapped string form.
			return nil, fmt.Errorf("%w: %s", ErrSkillNotFound, name)
		}
		return nil, fmt.Errorf("get skill: %w", err)
	}
	skill.Metadata = metadata

	// Load dependencies
	// G07 (research/g06_g07_skilltree_dag_design.md §4c): a DETERMINISTIC
	// ORDER BY is required for a byte-stable ExportToTOML round-trip
	// (§2.3(3)). Without it row order is query-plan-dependent, so two exports
	// of the same skill (or two skills carrying the same edge set) could emit
	// the typed-edge lists / [[skill.components]] entries in different orders.
	// Order by relation_type, then sort_order (the umbrella→component ordering,
	// NULLS LAST so unordered edges trail) then the target name as the stable
	// final tiebreak.
	depsSQL := `
		SELECT sd.skill_id, sd.depends_on, sd.relation_type, sd.optional, sd.sort_order,
		       ds.name as depends_on_name, ds.title as depends_on_title
		FROM skill_dependencies sd
		JOIN skills ds ON sd.depends_on = ds.id
		WHERE sd.skill_id = $1
		ORDER BY sd.relation_type, sd.sort_order NULLS LAST, ds.name
	`
	depRows, err := s.pool.Query(ctx, depsSQL, skill.ID)
	if err != nil {
		return nil, fmt.Errorf("get dependencies: %w", err)
	}
	defer depRows.Close()
	for depRows.Next() {
		var d models.SkillDependency
		if err := depRows.Scan(&d.SkillID, &d.DependsOn, &d.RelationType, &d.Optional, &d.SortOrder, &d.DependsOnName, &d.DependsOnTitle); err != nil {
			return nil, fmt.Errorf("scan dependency: %w", err)
		}
		skill.Dependencies = append(skill.Dependencies, d)
	}

	// Load resources
	// G07 §4c: deterministic ordering for byte-stable export (see depsSQL note).
	// url alone is NOT a total order — resources.url has no unique constraint
	// (migrations/001_initial.up.sql), so two resources may share a url while
	// differing in title/resource_type (the only other columns the export emits,
	// see ExportToTOML "Map resources"). id is a v4-random UUID re-minted on every
	// ImportFromTOML, so `ORDER BY url, id` reorders such same-url rows across an
	// export→import→export, breaking the §2.3(3) byte-stability contract. Order by
	// the STABLE emitted columns first (url, title, resource_type); the residual
	// id tiebreak then only ever separates BYTE-IDENTICAL exported rows, so the
	// ordering is a total order that is swap-invariant on the export.
	resSQL := `
		SELECT id, skill_id, url, title, resource_type, fetched_hash, content_cached, last_validated, created_at
		FROM resources WHERE skill_id = $1
		ORDER BY url, title, resource_type, id
	`
	resRows, err := s.pool.Query(ctx, resSQL, skill.ID)
	if err != nil {
		return nil, fmt.Errorf("get resources: %w", err)
	}
	defer resRows.Close()
	for resRows.Next() {
		var r models.Resource
		// B1 fix (Fable code-review remediation): fetched_hash/content_cached
		// are nullable TEXT columns (migrations/001_initial.up.sql) that
		// store.go never sets on INSERT until validation/caching runs
		// (internal/skill/resources.go), so a freshly-imported resource --
		// including every one imported via ImportFromTOML, which the B1 fix
		// makes non-empty for the first time -- has them NULL. Scanning a SQL
		// NULL directly into models.Resource's plain (non-nullable) string
		// fields panics/errors ("cannot scan NULL into *string"); this was
		// never exercised before because ImportFromTOML's resources always
		// decoded empty pre-fix, so GetByName never had a real resource row
		// to load. Scan through sql.NullString and default to "" (the same
		// value NewResource-class helpers already use for an unset hash, see
		// resources.go's own `SET content_cached = '', fetched_hash = ''`
		// reset), so this genuinely resolves the resource-import deadness end
		// to end rather than trading one silent gap for a crash.
		var fetchedHash, contentCached dbsql.NullString
		if err := resRows.Scan(&r.ID, &r.SkillID, &r.URL, &r.Title, &r.ResourceType, &fetchedHash, &contentCached, &r.LastValidated, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan resource: %w", err)
		}
		r.FetchedHash = fetchedHash.String
		r.ContentCached = contentCached.String
		skill.Resources = append(skill.Resources, r)
	}

	return &skill, nil
}

// GetTree returns the dependency tree for a skill up to the specified depth.
func (s *Store) GetTree(ctx context.Context, name string, maxDepth int) (*models.SkillTreeNode, error) {
	skill, err := s.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	root := &models.SkillTreeNode{
		Skill: *skill,
		Depth: 0,
	}

	visited := make(map[uuid.UUID]bool)
	visited[skill.ID] = true

	if err := s.buildTree(ctx, root, 1, maxDepth, visited); err != nil {
		return nil, fmt.Errorf("build tree: %w", err)
	}

	return root, nil
}

func (s *Store) buildTree(ctx context.Context, node *models.SkillTreeNode, depth, maxDepth int, visited map[uuid.UUID]bool) error {
	if depth > maxDepth {
		return nil
	}

	for _, dep := range node.Skill.Dependencies {
		if visited[dep.DependsOn] {
			continue
		}
		visited[dep.DependsOn] = true

		childSQL := `
			SELECT s.id, s.name, s.version, s.title, s.description, s.content,
			       s.metadata, s.status, s.kind, s.created_at, s.updated_at
			FROM skills s WHERE s.id = $1
		`
		var child models.Skill
		var metadata []byte
		err := s.pool.QueryRow(ctx, childSQL, dep.DependsOn).Scan(
			&child.ID, &child.Name, &child.Version, &child.Title,
			&child.Description, &child.Content, &metadata,
			&child.Status, &child.Kind, &child.CreatedAt, &child.UpdatedAt,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				continue // skip missing dependencies
			}
			return err
		}
		child.Metadata = metadata

		// Load child's dependencies
		depsSQL := `
			SELECT sd.skill_id, sd.depends_on, sd.relation_type, sd.optional, sd.sort_order,
			       ds.name, ds.title
			FROM skill_dependencies sd
			JOIN skills ds ON sd.depends_on = ds.id
			WHERE sd.skill_id = $1
		`
		depRows, err := s.pool.Query(ctx, depsSQL, child.ID)
		if err != nil {
			return err
		}
		for depRows.Next() {
			var d models.SkillDependency
			if err := depRows.Scan(&d.SkillID, &d.DependsOn, &d.RelationType, &d.Optional, &d.SortOrder, &d.DependsOnName, &d.DependsOnTitle); err != nil {
				depRows.Close()
				return err
			}
			child.Dependencies = append(child.Dependencies, d)
		}
		depRows.Close()

		childNode := models.SkillTreeNode{
			Skill: child,
			Depth: depth,
		}
		node.Children = append(node.Children, childNode)

		if err := s.buildTree(ctx, &node.Children[len(node.Children)-1], depth+1, maxDepth, visited); err != nil {
			return err
		}
	}

	return nil
}

// Create inserts a new skill into the database.
func (s *Store) Create(ctx context.Context, skill *models.Skill) error {
	if skill.ID == uuid.Nil {
		skill.ID = uuid.New()
	}

	metadataJSON, err := json.Marshal(skill.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	sql := `
		INSERT INTO skills (id, name, version, title, description, content, metadata, status, kind, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		ON CONFLICT (name) DO UPDATE SET
			version = EXCLUDED.version,
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			content = EXCLUDED.content,
			metadata = EXCLUDED.metadata,
			status = EXCLUDED.status,
			kind = EXCLUDED.kind,
			updated_at = NOW()
		RETURNING id
	`
	var returnedID uuid.UUID
	err = s.pool.QueryRow(ctx, sql,
		skill.ID, skill.Name, skill.Version, skill.Title,
		skill.Description, skill.Content, metadataJSON,
		skill.Status, skill.Kind.NormalizeOrAtomic(),
	).Scan(&returnedID)
	if err != nil {
		return fmt.Errorf("create skill: %w", err)
	}

	// Insert dependencies. ON CONFLICT targets the widened
	// (skill_id, depends_on, relation_type) primary key introduced by
	// migrations/002_granularity.up.sql — a pair may now carry more than one
	// typed edge (e.g. both `requires` and `recommends`), so the old
	// (skill_id, depends_on) conflict target no longer matches any unique
	// index (research/p1t1_granularity_schema_migration.md §2 L3).
	for _, dep := range skill.Dependencies {
		depSQL := `
			INSERT INTO skill_dependencies (skill_id, depends_on, relation_type, optional, sort_order)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (skill_id, depends_on, relation_type) DO UPDATE SET
				optional = EXCLUDED.optional,
				sort_order = EXCLUDED.sort_order
		`
		_, err := s.pool.Exec(ctx, depSQL, returnedID, dep.DependsOn, dep.RelationType, dep.Optional, dep.SortOrder)
		if err != nil {
			return fmt.Errorf("create dependency: %w", err)
		}
	}

	// Upsert registry entry
	regSQL := `
		INSERT INTO skill_registry (skill_id, skill_name, missing_deps, stale, auto_expand, coverage)
		VALUES ($1, $2, '{}', false, true, 0.0)
		ON CONFLICT (skill_id) DO NOTHING
	`
	_, err = s.pool.Exec(ctx, regSQL, returnedID, skill.Name)
	if err != nil {
		return fmt.Errorf("create registry entry: %w", err)
	}

	skill.ID = returnedID
	return nil
}

// CreateFromTOML creates a skill from a TOML skill wrapper.
func (s *Store) CreateFromTOML(ctx context.Context, wrapper *models.TOMLSkillWrapper) (*models.Skill, error) {
	metadataJSON, _ := json.Marshal(wrapper.Skill.Metadata)

	skill := &models.Skill{
		ID:          uuid.New(),
		Name:        wrapper.Skill.Name,
		Version:     wrapper.Skill.Version,
		Title:       wrapper.Skill.Title,
		Description: wrapper.Skill.Description,
		Content:     wrapper.Skill.Content,
		Metadata:    metadataJSON,
		Status:      models.SkillStatusDraft,
		Kind:        models.SkillKind(wrapper.Skill.Kind).NormalizeOrAtomic(),
	}

	// Resolve dependencies
	for _, depName := range wrapper.Skill.Dependencies.Requires {
		depID, err := s.resolveSkillID(ctx, depName)
		if err != nil {
			return nil, fmt.Errorf("resolve dependency %q: %w", depName, err)
		}
		skill.Dependencies = append(skill.Dependencies, models.SkillDependency{
			SkillID:      skill.ID,
			DependsOn:    depID,
			RelationType: models.DepTypeRequires,
		})
	}
	for _, depName := range wrapper.Skill.Dependencies.Extends {
		depID, err := s.resolveSkillID(ctx, depName)
		if err != nil {
			return nil, fmt.Errorf("resolve dependency %q: %w", depName, err)
		}
		skill.Dependencies = append(skill.Dependencies, models.SkillDependency{
			SkillID:      skill.ID,
			DependsOn:    depID,
			RelationType: models.DepTypeExtends,
		})
	}
	for _, depName := range wrapper.Skill.Dependencies.Recommends {
		depID, err := s.resolveSkillID(ctx, depName)
		if err != nil {
			return nil, fmt.Errorf("resolve dependency %q: %w", depName, err)
		}
		skill.Dependencies = append(skill.Dependencies, models.SkillDependency{
			SkillID:      skill.ID,
			DependsOn:    depID,
			RelationType: models.DepTypeRecommends,
		})
	}

	// Add resources
	for _, r := range wrapper.Skill.Resources {
		skill.Resources = append(skill.Resources, models.Resource{
			ID:           uuid.New(),
			SkillID:      skill.ID,
			URL:          r.URL,
			Title:        r.Title,
			ResourceType: r.ResourceType,
		})
	}

	if err := s.Create(ctx, skill); err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *Store) resolveSkillID(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, "SELECT id FROM skills WHERE name = $1", name).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, fmt.Errorf("skill %q not found", name)
		}
		return uuid.Nil, err
	}
	return id, nil
}

// GetMissingSkills returns skills with missing dependencies (gaps in the graph).
func (s *Store) GetMissingSkills(ctx context.Context, domain string) ([]models.SkillRegistryEntry, error) {
	sql := `
		SELECT sr.skill_id, sr.skill_name, sr.missing_deps, sr.stale, sr.last_review, sr.auto_expand, sr.coverage
		FROM skill_registry sr
		WHERE array_length(sr.missing_deps, 1) > 0
	`
	args := []interface{}{}
	argIdx := 1

	if domain != "" {
		sql += fmt.Sprintf(` AND EXISTS (
			SELECT 1 FROM skills s
			WHERE s.id = sr.skill_id AND s.metadata->>'domain' = $%d
		)`, argIdx)
		args = append(args, domain)
		argIdx++
	}

	sql += " ORDER BY sr.coverage ASC, sr.skill_name"

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("get missing skills: %w", err)
	}
	defer rows.Close()

	var entries []models.SkillRegistryEntry
	for rows.Next() {
		var e models.SkillRegistryEntry
		if err := rows.Scan(&e.SkillID, &e.SkillName, &e.MissingDeps, &e.Stale, &e.LastReview, &e.AutoExpand, &e.Coverage); err != nil {
			return nil, fmt.Errorf("scan registry entry: %w", err)
		}
		entries = append(entries, e)
	}

	return entries, nil
}

// GetCoverage returns coverage statistics for a domain.
func (s *Store) GetCoverage(ctx context.Context, domain string) (map[string]interface{}, error) {
	// Count total skills
	var total int
	totalSQL := "SELECT COUNT(*) FROM skills"
	var totalArgs []interface{}
	if domain != "" {
		totalSQL += " WHERE metadata->>'domain' = $1"
		totalArgs = append(totalArgs, domain)
	}
	if err := s.pool.QueryRow(ctx, totalSQL, totalArgs...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count skills: %w", err)
	}

	// Count with dependencies
	var withDeps int
	depSQL := `
		SELECT COUNT(DISTINCT s.id) FROM skills s
		WHERE EXISTS (SELECT 1 FROM skill_dependencies sd WHERE sd.skill_id = s.id)
	`
	var depArgs []interface{}
	if domain != "" {
		depSQL += " AND s.metadata->>'domain' = $1"
		depArgs = append(depArgs, domain)
	}
	if err := s.pool.QueryRow(ctx, depSQL, depArgs...).Scan(&withDeps); err != nil {
		return nil, fmt.Errorf("count with deps: %w", err)
	}

	// Count with evidence
	var withEvidence int
	evSQL := `
		SELECT COUNT(DISTINCT s.id) FROM skills s
		WHERE EXISTS (SELECT 1 FROM evidences e WHERE e.skill_id = s.id)
	`
	var evArgs []interface{}
	if domain != "" {
		evSQL += " AND s.metadata->>'domain' = $1"
		evArgs = append(evArgs, domain)
	}
	if err := s.pool.QueryRow(ctx, evSQL, evArgs...).Scan(&withEvidence); err != nil {
		return nil, fmt.Errorf("count with evidence: %w", err)
	}

	// Average coverage from registry
	var avgCoverage float64
	covSQL := `
		SELECT COALESCE(AVG(sr.coverage), 0.0) FROM skill_registry sr
		JOIN skills s ON sr.skill_id = s.id
	`
	var covArgs []interface{}
	if domain != "" {
		covSQL += " WHERE s.metadata->>'domain' = $1"
		covArgs = append(covArgs, domain)
	}
	if err := s.pool.QueryRow(ctx, covSQL, covArgs...).Scan(&avgCoverage); err != nil {
		return nil, fmt.Errorf("avg coverage: %w", err)
	}

	// Count missing dependencies
	var missingCount int
	missSQL := `
		SELECT COUNT(*) FROM skill_registry sr
		JOIN skills s ON sr.skill_id = s.id
		WHERE array_length(sr.missing_deps, 1) > 0
	`
	var missArgs []interface{}
	if domain != "" {
		missSQL += " AND s.metadata->>'domain' = $1"
		missArgs = append(missArgs, domain)
	}
	if err := s.pool.QueryRow(ctx, missSQL, missArgs...).Scan(&missingCount); err != nil {
		return nil, fmt.Errorf("count missing: %w", err)
	}

	coverage := 0.0
	if total > 0 {
		coverage = avgCoverage
	}

	return map[string]interface{}{
		"domain":               domain,
		"total_skills":         total,
		"skills_with_deps":     withDeps,
		"skills_with_evidence": withEvidence,
		"skills_missing_deps":  missingCount,
		"average_coverage":     coverage,
		"coverage_percentage":  fmt.Sprintf("%.1f%%", coverage*100),
	}, nil
}

// SubmitLearningJob creates a new learning job for project analysis.
func (s *Store) SubmitLearningJob(ctx context.Context, projectPath string, languages []string) (*models.LearningJob, error) {
	job := &models.LearningJob{
		ID:          uuid.New(),
		ProjectPath: projectPath,
		Status:      "pending",
		Languages:   languages,
	}

	// For now, just store in audit log. In production, this would insert into a learning_jobs table.
	details, _ := json.Marshal(map[string]interface{}{
		"project_path": projectPath,
		"languages":    languages,
		"job_id":       job.ID,
	})

	_, err := s.pool.Exec(ctx,
		"INSERT INTO audit_log (event, details) VALUES ($1, $2)",
		"learning_job_submitted", details,
	)
	if err != nil {
		return nil, fmt.Errorf("log learning job: %w", err)
	}

	return job, nil
}

// VectorSearch performs vector similarity search using pgvector.
func (s *Store) VectorSearch(ctx context.Context, embedding []float32, limit int) ([]models.SearchResult, error) {
	vec := pgvector.NewVector(embedding)
	// WHERE s.embedding IS NOT NULL (F2, code-review BLOCKING; aligned with the
	// reference sibling internal/db/vector.go's VectorSearch): skills.embedding is
	// a nullable column (migrations/001_initial.up.sql) that store.Create never
	// sets, so it is NULL in the ordinary partially-/un-populated state. On the
	// HNSW index-scan plan NULLs are skipped, but the cost-based planner CAN pick a
	// seqscan/top-N plan (small table, or no usable index) where `ORDER BY
	// s.embedding <=> $1 LIMIT $2` sorts NULL distances LAST and, once LIMIT
	// exceeds the non-NULL row count, RETURNS a NULL-embedding row with a NULL
	// score -- pgx v5 then errors "cannot scan NULL into *float64". Because
	// Store.Search deliberately does not mask a vector-leg error, that single NULL
	// row hard-fails EVERY hybrid Search and discards the trigram results. Filtering
	// NULLs in the vector leg makes correctness independent of the query plan.
	sql := `
		SELECT s.id, s.name, s.version, s.title, s.description, s.content,
		       s.metadata, s.status, s.kind, s.created_at, s.updated_at,
		       1 - (s.embedding <=> $1) as score
		FROM skills s
		WHERE s.embedding IS NOT NULL
		ORDER BY s.embedding <=> $1
		LIMIT $2
	`
	rows, err := s.pool.Query(ctx, sql, vec, limit)
	if err != nil {
		return nil, fmt.Errorf("vector search: %w", err)
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var r models.SearchResult
		var metadata []byte
		var description dbsql.NullString // F2/finding 7 precedent, see textSearch.
		if err := rows.Scan(
			&r.Skill.ID, &r.Skill.Name, &r.Skill.Version, &r.Skill.Title,
			&description, &r.Skill.Content, &metadata,
			&r.Skill.Status, &r.Skill.Kind, &r.Skill.CreatedAt, &r.Skill.UpdatedAt,
			&r.Score,
		); err != nil {
			return nil, fmt.Errorf("scan vector result: %w", err)
		}
		r.Skill.Description = description.String
		r.Skill.Metadata = metadata
		results = append(results, r)
	}
	// F4/finding 4: see textSearch's rows.Err() note above -- same mid-stream
	// truncation risk applies to this KNN query's row iteration.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("vector search rows: %w", err)
	}

	return results, nil
}

// ListSkills returns all skills with optional filtering.
func (s *Store) ListSkills(ctx context.Context, status models.SkillStatus, limit, offset int) ([]models.Skill, error) {
	sql := `
		SELECT id, name, version, title, description, content,
		       metadata, status, kind, created_at, updated_at
		FROM skills
	`
	args := []interface{}{}
	conditions := []string{}

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, status)
	}

	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	sql += fmt.Sprintf(" ORDER BY name LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("list skills: %w", err)
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var sk models.Skill
		var metadata []byte
		if err := rows.Scan(&sk.ID, &sk.Name, &sk.Version, &sk.Title, &sk.Description, &sk.Content, &metadata, &sk.Status, &sk.Kind, &sk.CreatedAt, &sk.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan skill: %w", err)
		}
		sk.Metadata = metadata
		skills = append(skills, sk)
	}

	return skills, nil
}

// logAudit is a helper for audit logging used by other skill package files.
func (s *Store) logAudit(ctx context.Context, tx pgx.Tx, event string, skillID *uuid.UUID, details map[string]interface{}) error {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("marshal audit details: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO audit_log (ts, event, skill_id, details)
		VALUES ($1, $2, $3, $4)
	`, time.Now().UTC(), event, skillID, detailsJSON)

	return err
}
