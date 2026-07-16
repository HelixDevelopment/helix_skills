package autoexpand

// G03 remediation (the autoexpand+worker half) -- live-database (no live LLM
// required) proof that draftPersistAndCrossReference (pipeline.go) genuinely
// persists the drafted sub-skill AND cross-references it into the tree by
// adding a `requires` edge from the gap's parent skill, resolved via the
// EXACT Store.GetByName lookup (§G29/§G60), to the newly created skill.
//
// This test deliberately does NOT configure an LLM client (p.llm stays
// nil), so DraftSkill takes its documented no-LLM minimal-fallback path
// (createMinimalDraft) -- proving the PERSIST + CROSS-REFERENCE steps work
// correctly independent of whether a live LLM key is available. The
// LLM-backed drafting half ("does the real Anthropic API get called and
// does its content land in the created skill") is covered separately, with
// a REAL Anthropic call, by pipeline_crossreference_integration_test.go
// (env-gated on ANTHROPIC_API_KEY, `integration` build tag) -- this file's
// job is the part that is verifiable in THIS environment without an
// operator-supplied LLM credential (§11.4.3), and is the file this ticket's
// §1.1 mutation cycle for the cross-reference addition was run against.
//
// CURRENT behavior; the minimal-draft (createMinimalDraft) fallback this
// test asserts gets PERSISTED is slated for removal per G20
// (never-persist-a-placeholder, GAPS_AND_RISKS_REGISTER.md: "Never persist
// a placeholder as a real skill ... degrades into bluff data"). This test
// documents the reconciliation trail: when G20's fallback-removal half
// lands, this test's "no-LLM path persists a placeholder" assertions are
// EXPECTED to flip (the no-LLM path should then mark the gap unfilled
// instead of persisting createMinimalDraft's boilerplate content) -- that
// flip is not a silent regression, it is G20 closing.
//
// Why draftPersistAndCrossReference is called directly with a manually
// constructed Gap, rather than through Pipeline.Run's own gap scan: see
// pipeline_crossreference_integration_test.go's header comment for the
// confirmed, cited FACT (not a guess, §11.4.6) that Run's gap-detection
// cannot fire against any graph the store API constructs (a hand-inserted
// zero-UUID skills row via raw SQL is schema-valid and does make the gate
// reachable -- see G137, GAPS_AND_RISKS_REGISTER.md).
//
// Gated on SKILL_SYSTEM_TEST_DB_HOST (§11.4.3): absent a configured live
// database this honestly t.Skip()s, never a fake PASS (§11.4.27).

import (
	"context"
	"testing"

	"github.com/helixdevelopment/skill-system/internal/config"
	"github.com/helixdevelopment/skill-system/internal/db"
	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/helixdevelopment/skill-system/internal/skill"
	"go.uber.org/zap"
)

func TestDraftPersistAndCrossReference_PersistsAndCrossReferences_RequiresLiveDatabase(t *testing.T) {
	admin, ok := aeSkipIfNoTestDB(t)
	if !ok {
		return
	}
	ctx := context.Background()

	dbCfg, cleanup := aeCreateThrowawayDB(t, admin)
	defer cleanup()

	pool, err := db.New(dbCfg)
	if err != nil {
		t.Fatalf("db.New(dbCfg): %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool, aeRealMigrationsDir); err != nil {
		t.Fatalf("db.Migrate: %v", err)
	}

	store := skill.NewStore(pool)

	parent := &models.Skill{
		Name:    "g03ae.crossref.nollm.parent-skill",
		Title:   "G03 autoexpand cross-reference parent (no-LLM path)",
		Content: "content for the G03 autoexpand draftPersistAndCrossReference no-LLM test",
		Status:  models.SkillStatusActive,
		Kind:    models.SkillKindAtomic,
	}
	if err := store.Create(ctx, parent); err != nil {
		t.Fatalf("create parent skill: %v", err)
	}

	// No WithLLMClient option: p.llm stays nil, DraftSkill takes the
	// documented no-LLM minimal-fallback path (createMinimalDraft) -- no
	// network access anywhere in this test.
	p := NewPipeline(store, nil, config.AutoExpandConfig{MaxDepth: 2, MaxNewSkillsPerRun: 5}, zap.NewNop())

	childName := "g03ae-crossref-nollm-child-skill"
	gap := Gap{
		SkillName:      parent.Name,
		MissingDepName: childName,
		SuggestedTitle: "No-LLM Child Skill",
		Reason:         "test-seeded gap for the no-LLM draftPersistAndCrossReference path",
	}

	result := &ExpansionResult{}
	draft, err := p.draftPersistAndCrossReference(ctx, gap, result)
	if err != nil {
		t.Fatalf("draftPersistAndCrossReference: %v", err)
	}
	if len(result.Errors) != 0 {
		t.Fatalf("draftPersistAndCrossReference reported errors: %v", result.Errors)
	}
	if draft.Name != childName {
		t.Fatalf("draft.Name = %q, want %q", draft.Name, childName)
	}

	// PERSIST: the sub-skill is a real, independently re-readable row.
	reread, err := store.GetByName(ctx, childName)
	if err != nil {
		t.Fatalf("GetByName(%q) after draftPersistAndCrossReference: %v", childName, err)
	}
	if reread.ID != draft.ID {
		t.Fatalf("re-read skill id = %s, want %s", reread.ID, draft.ID)
	}
	if reread.Status != models.SkillStatusDraft {
		t.Errorf("re-read skill status = %q, want %q (a freshly auto-expanded skill starts as a draft "+
			"pending review)", reread.Status, models.SkillStatusDraft)
	}

	// CROSS-REFERENCE: the parent, re-resolved via the EXACT Store.GetByName
	// lookup, now lists the new skill as a `requires` dependency.
	parentReread, err := store.GetByName(ctx, parent.Name)
	if err != nil {
		t.Fatalf("GetByName(%q) (parent) after cross-reference: %v", parent.Name, err)
	}
	found := false
	for _, dep := range parentReread.Dependencies {
		if dep.DependsOn == draft.ID && dep.DependsOnName == childName {
			found = true
			if dep.RelationType != models.DepTypeRequires {
				t.Errorf("cross-referenced edge relation_type = %q, want %q", dep.RelationType, models.DepTypeRequires)
			}
			break
		}
	}
	if !found {
		names := make([]string, 0, len(parentReread.Dependencies))
		for _, d := range parentReread.Dependencies {
			names = append(names, d.DependsOnName)
		}
		t.Fatalf("parent skill %q has no dependency edge to the newly created %q after "+
			"draftPersistAndCrossReference; got dependencies=%v -- the new sub-skill was persisted "+
			"but never cross-referenced into the tree", parent.Name, childName, names)
	}
}
