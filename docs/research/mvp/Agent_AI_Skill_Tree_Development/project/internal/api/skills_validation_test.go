package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/helixdevelopment/skill-system/internal/validation"
)

// fakePool embeds the Pool interface (nil) and overrides only CreateSkill, which
// is the sole pool method handleCreateSkill invokes. Any other method call panics
// — proving the create path touches nothing else.
type fakePool struct {
	Pool
	created *models.Skill
}

func (f *fakePool) CreateSkill(_ context.Context, s *models.Skill) error {
	f.created = s
	return nil
}

type fakeValidator struct{ passed bool }

func (f fakeValidator) Validate(_ context.Context, _ *models.Skill) (*validation.ValidationResult, error) {
	return &validation.ValidationResult{Passed: f.passed, Stage: "test"}, nil
}

func newCreateServer(t *testing.T, v SkillValidator, enabled bool) (*Server, *fakePool) {
	t.Helper()
	fp := &fakePool{}
	// AuthDisabled=true installs no auth middleware, so the request reaches the
	// handler; the CORS allowlist stays empty (fail-closed).
	cfg := Config{Server: ServerConfig{AuthDisabled: true}}
	srv := New(fp, cfg, zap.NewNop(), WithValidator(v, enabled))
	return srv, fp
}

func postSkill(t *testing.T, srv *Server, body string) int {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/skills", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)
	return w.Code
}

// TestCreateSkill_PromotesOnlyOnPassingVerdict is the paired-mutation-real proof
// of the §G03 request-path fail-closed policy: a client requesting status
// "active" is promoted ONLY when validation passes; on a failing verdict the
// skill is persisted as draft regardless of the requested status.
func TestCreateSkill_PromotesOnlyOnPassingVerdict(t *testing.T) {
	body := `{"name":"my-skill","title":"My Skill","content":"# doc","status":"active"}`

	t.Run("passing verdict -> active honoured", func(t *testing.T) {
		srv, fp := newCreateServer(t, fakeValidator{passed: true}, true)
		if code := postSkill(t, srv, body); code != http.StatusCreated {
			t.Fatalf("status = %d, want 201", code)
		}
		if fp.created == nil {
			t.Fatal("CreateSkill was not called")
		}
		if fp.created.Status != models.SkillStatusActive {
			t.Errorf("status = %q, want %q (active honoured on pass)", fp.created.Status, models.SkillStatusActive)
		}
	})

	t.Run("failing verdict -> forced draft despite requested active", func(t *testing.T) {
		srv, fp := newCreateServer(t, fakeValidator{passed: false}, true)
		if code := postSkill(t, srv, body); code != http.StatusCreated {
			t.Fatalf("status = %d, want 201", code)
		}
		if fp.created == nil {
			t.Fatal("CreateSkill was not called")
		}
		if fp.created.Status != models.SkillStatusDraft {
			t.Errorf("status = %q, want %q (fail-closed: no self-promotion)", fp.created.Status, models.SkillStatusDraft)
		}
	})

	t.Run("validation disabled -> forced draft", func(t *testing.T) {
		srv, fp := newCreateServer(t, fakeValidator{passed: true}, false)
		if code := postSkill(t, srv, body); code != http.StatusCreated {
			t.Fatalf("status = %d, want 201", code)
		}
		if fp.created == nil {
			t.Fatal("CreateSkill was not called")
		}
		if fp.created.Status != models.SkillStatusDraft {
			t.Errorf("status = %q, want %q (disabled ⇒ never auto-promote)", fp.created.Status, models.SkillStatusDraft)
		}
	})
}
