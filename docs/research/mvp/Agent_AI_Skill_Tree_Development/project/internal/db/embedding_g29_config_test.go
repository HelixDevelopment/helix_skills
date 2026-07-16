package db

// Fable code-review remediation, finding 6a (NIT): NewEmbedderFromConfig's
// "openai" branch used to construct-and-succeed (with only a logged warning)
// even when cfg.APIKey was empty -- the ONLY fail-closed case was "local" with
// a missing endpoint. That asymmetry forced internal/mcp.NewMCPServer to keep
// its OWN, hand-duplicated embeddingConfigured() policy check as a 2nd source
// of truth for "is this provider really usable", which could silently drift
// from this factory's own policy.
//
// The fix makes "openai" fail-closed on a missing api_key too, matching the
// "local" case's existing missing-endpoint check, so a caller can derive
// "configured" SOLELY from this function's error return. This is a pure unit
// test -- no live database required.

import (
	"testing"

	"github.com/helixdevelopment/skill-system/internal/config"
)

// TestG29_NewEmbedderFromConfig_OpenAIRequiresAPIKey is the RED-first
// regression guard for finding 6a: it FAILS on the pre-fix factory (which
// returned a non-nil *OpenAIEmbedder and a nil error for an empty api_key) and
// PASSES post-fix (a §1.1 mutation reverting the "openai" branch to the old
// warn-only behaviour makes it FAIL again).
func TestG29_NewEmbedderFromConfig_OpenAIRequiresAPIKey(t *testing.T) {
	cfg := config.EmbeddingConfig{Provider: "openai", APIKey: "", Model: "text-embedding-3-small", Dimensions: 768}

	emb, err := NewEmbedderFromConfig(cfg)
	if err == nil {
		t.Fatalf("NewEmbedderFromConfig(openai, empty api_key) = (%v, nil), want a non-nil error; "+
			"a caller deriving \"is this provider configured\" from this function's error return "+
			"alone must see a missing api_key rejected, not a silently-unusable embedder", emb)
	}
	if emb != nil {
		t.Errorf("NewEmbedderFromConfig(openai, empty api_key) returned a non-nil Embedder alongside its error: %v", emb)
	}
}

// TestG29_NewEmbedderFromConfig_OpenAIWithAPIKeySucceeds is the accompanying
// positive case: a genuinely-configured openai provider still constructs
// successfully -- the fix narrows the previously-permissive case, it does not
// break the legitimate one.
func TestG29_NewEmbedderFromConfig_OpenAIWithAPIKeySucceeds(t *testing.T) {
	cfg := config.EmbeddingConfig{Provider: "openai", APIKey: "sk-test-not-a-real-key", Model: "text-embedding-3-small", Dimensions: 768}

	emb, err := NewEmbedderFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewEmbedderFromConfig(openai, with api_key) returned an unexpected error: %v", err)
	}
	if emb == nil {
		t.Fatalf("NewEmbedderFromConfig(openai, with api_key) returned a nil Embedder alongside a nil error")
	}
	if got := emb.Dimensions(); got != 768 {
		t.Errorf("Dimensions() = %d, want 768", got)
	}
}

// TestG29_NewEmbedderFromConfig_LocalStillRequiresEndpoint pins the PRE-EXISTING
// "local" fail-closed behaviour (unchanged by this fix) so a future edit
// cannot silently loosen it while tightening the "openai" case.
func TestG29_NewEmbedderFromConfig_LocalStillRequiresEndpoint(t *testing.T) {
	cfg := config.EmbeddingConfig{Provider: "local", LocalEndpoint: ""}
	if _, err := NewEmbedderFromConfig(cfg); err == nil {
		t.Fatal("NewEmbedderFromConfig(local, empty local_endpoint) = nil error, want non-nil")
	}
}
