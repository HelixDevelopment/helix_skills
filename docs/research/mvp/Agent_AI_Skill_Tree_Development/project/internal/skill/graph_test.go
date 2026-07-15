package skill

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/helixdevelopment/skill-system/internal/models"
)

// ---------------------------------------------------------------------------
// AddDependency: pure, DB-independent guard clauses.
//
// AddDependency's self-reference check and relation-type validation both
// execute and return BEFORE any *db.Pool method is invoked (see graph.go:
// the `if skillID == dependsOn` and `if !validTypes[relType]` checks precede
// s.pool.WithTx(...)). That makes them safe to exercise against a Store with
// a nil pool -- if either guard clause were ever moved after a pool access,
// or removed, these tests would panic (nil pointer deref) or fail (wrong
// error), so they also serve as a regression trip-wire for that ordering.
// ---------------------------------------------------------------------------

func TestAddDependency_SelfReferenceIsRejectedAsCycle(t *testing.T) {
	s := NewStore(nil) // no DB required: the self-reference check returns first.
	id := uuid.New()

	err := s.AddDependency(context.Background(), id, id, models.DepTypeRequires)
	if err == nil {
		t.Fatal("expected an error for a self-referencing dependency, got nil")
	}
	if !errors.Is(err, ErrCycleDetected) {
		t.Errorf("error = %v, want it to wrap ErrCycleDetected", err)
	}
}

func TestAddDependency_InvalidRelationTypeIsRejected(t *testing.T) {
	s := NewStore(nil) // no DB required: the relation-type validation returns first.
	from := uuid.New()
	to := uuid.New()

	err := s.AddDependency(context.Background(), from, to, models.DependencyType("bogus-relation"))
	if err == nil {
		t.Fatal("expected an error for an invalid relation type, got nil")
	}
	if !errors.Is(err, ErrInvalidSkill) {
		t.Errorf("error = %v, want it to wrap ErrInvalidSkill", err)
	}
}

// Note: we deliberately do NOT exercise AddDependency with a *valid*
// relation type end-to-end here. Past the two guard clauses above, the very
// next line is s.pool.WithTx(...), which immediately calls p.inner.Acquire
// on the (nil in these tests) underlying *pgxpool.Pool and panics with a nil
// pointer dereference. That is the real, DB-bound code path -- see
// TestAddDependency_DuplicateAndPersistence_RequiresLiveDatabase below.

// TestHasCycle_RequiresLiveDatabase documents, rather than bluffs, the
// boundary of what this package can test without infrastructure. The actual
// cycle-detection algorithm (hasCycle in graph.go) is implemented as a
// PostgreSQL recursive CTE executed over a live pgx.Tx -- there is no
// in-memory graph structure in this codebase to walk instead. Faking a
// pgx.Tx/pgx.Rows well enough to prove real reachability semantics would
// either (a) require a real PostgreSQL connection (integration test, out of
// scope here) or (b) reimplement the recursive CTE's semantics in Go and
// compare against a mock -- which would test the mock, not the production
// SQL, and would be exactly the kind of bluff the task forbids.
func TestHasCycle_RequiresLiveDatabase(t *testing.T) {
	t.Skip("hasCycle() executes a recursive CTE against a live pgx.Tx (PostgreSQL); " +
		"no in-memory graph implementation exists to unit test in its place. " +
		"Requires an integration test against a real/containerized Postgres instance.")
}

// TestAddDependency_DuplicateAndPersistence_RequiresLiveDatabase documents
// the same boundary for the rest of AddDependency's behavior (existence
// checks, duplicate-edge detection, the actual INSERT, and the registry
// recalculation) -- all of it runs inside s.pool.WithTx against live SQL.
func TestAddDependency_DuplicateAndPersistence_RequiresLiveDatabase(t *testing.T) {
	t.Skip("AddDependency's existence/duplicate checks and the edge INSERT run inside " +
		"s.pool.WithTx against live PostgreSQL (skills, skill_dependencies, skill_registry " +
		"tables). Requires an integration test against a real/containerized Postgres instance.")
}

// ---------------------------------------------------------------------------
// collectDepNames (import_export.go): pure, in-memory dependency-name
// deduplication used when importing a TOML skill definition.
// ---------------------------------------------------------------------------

func TestCollectDepNames(t *testing.T) {
	tests := []struct {
		name string
		deps models.TOMLDependencies
		want []string
	}{
		{
			name: "empty dependencies yields empty slice",
			deps: models.TOMLDependencies{},
			want: nil,
		},
		{
			name: "single requires entry",
			deps: models.TOMLDependencies{Requires: []string{"foundations"}},
			want: []string{"foundations"},
		},
		{
			name: "requires, extends, recommends combined preserving first-seen order",
			deps: models.TOMLDependencies{
				Requires:   []string{"a", "b"},
				Extends:    []string{"c"},
				Recommends: []string{"d"},
			},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "duplicate name across requires and extends deduplicated, first occurrence kept",
			deps: models.TOMLDependencies{
				Requires: []string{"shared", "only-requires"},
				Extends:  []string{"shared", "only-extends"},
			},
			want: []string{"shared", "only-requires", "only-extends"},
		},
		{
			name: "duplicate name within the same list deduplicated",
			deps: models.TOMLDependencies{
				Requires: []string{"dup", "dup", "unique"},
			},
			want: []string{"dup", "unique"},
		},
		{
			name: "duplicate across all three lists appears exactly once",
			deps: models.TOMLDependencies{
				Requires:   []string{"x"},
				Extends:    []string{"x"},
				Recommends: []string{"x"},
			},
			want: []string{"x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectDepNames(tt.deps)
			if !equalStringSlices(got, tt.want) {
				t.Errorf("collectDepNames(%+v) = %v, want %v", tt.deps, got, tt.want)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
