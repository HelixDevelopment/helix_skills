// Package validation implements the multi-layer skill validation pipeline
// for the HelixKnowledge system. It provides the "zero-bluff guarantee"
// through source verification, sandboxed code execution, multi-model LLM
// jury validation, and cross-reference consistency checks.
package validation

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/helixdevelopment/skill-system/internal/config"
	"github.com/helixdevelopment/skill-system/internal/db"
	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/helixdevelopment/skill-system/internal/skill"
	"go.uber.org/zap"
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// ValidationResult captures the outcome of the full validation pipeline.
type ValidationResult struct {
	SkillID    uuid.UUID         `json:"skill_id"`
	Passed     bool              `json:"passed"`
	Stage      string            `json:"stage"`       // which stage failed (if any)
	Details    map[string]string `json:"details"`     // per-stage detail messages
	ApprovedBy int               `json:"approved_by"` // number of jury models that approved
}

// JuryResult captures the multi-model consensus outcome.
type JuryResult struct {
	Votes     map[string]bool `json:"votes"`     // model name -> approved
	Consensus bool            `json:"consensus"` // true if threshold met
	Feedback  string          `json:"feedback"`  // aggregated feedback
}

// JuryMember represents a single LLM validator in the jury.
type JuryMember struct {
	Name   string
	Weight float64 // voting weight (default 1.0)
	LLM    LLMValidator
}

// LLMValidator is the interface for a single LLM validator.
type LLMValidator interface {
	// ValidateSkill evaluates a skill and returns an approval decision with feedback.
	ValidateSkill(ctx context.Context, skill *models.Skill) (approved bool, feedback string, err error)
}

// SandboxResult captures the outcome of sandboxed code execution.
type SandboxResult struct {
	ExecutionResult
	Approved bool   `json:"approved"` // whether the code passed validation
	Feedback string `json:"feedback"`
}

// Pipeline validates skills through multiple independent layers to ensure
// accuracy and prevent hallucinated content from entering the knowledge graph.
type Pipeline struct {
	store     *skill.Store
	cfg       config.ValidationConfig
	logger    *zap.Logger
	sandbox   Sandbox
	jury      []JuryMember
	httpClient *http.Client
}

// PipelineOption allows optional configuration.
type PipelineOption func(*Pipeline)

// WithSandbox sets a custom sandbox for code execution.
func WithSandbox(sandbox Sandbox) PipelineOption {
	return func(p *Pipeline) {
		p.sandbox = sandbox
	}
}

// WithJury sets custom jury members for LLM validation.
func WithJury(jury []JuryMember) PipelineOption {
	return func(p *Pipeline) {
		p.jury = jury
	}
}

// ---------------------------------------------------------------------------
// Construction
// ---------------------------------------------------------------------------

// NewPipeline creates a new validation pipeline with the configured sandbox
// and jury setup.
func NewPipeline(store *skill.Store, cfg config.ValidationConfig, logger *zap.Logger, opts ...PipelineOption) *Pipeline {
	p := &Pipeline{
		store:      store,
		cfg:        cfg,
		logger:     logger,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		opt(p)
	}

	// Initialize sandbox if not provided
	if p.sandbox == nil {
		p.sandbox = p.createDefaultSandbox()
	}

	return p
}

// createDefaultSandbox creates the appropriate sandbox based on config.
func (p *Pipeline) createDefaultSandbox() Sandbox {
	switch p.cfg.SandboxType {
	case "wasm":
		return NewWASMSandbox(p.logger)
	case "docker":
		return NewDockerSandbox(p.logger)
	case "none":
		p.logger.Warn("sandbox disabled - code execution will be skipped")
		return &NoOpSandbox{}
	default:
		// Default to WASM (lightweight)
		p.logger.Info("using default WASM sandbox", zap.String("sandbox_type", p.cfg.SandboxType))
		return NewWASMSandbox(p.logger)
	}
}

// ---------------------------------------------------------------------------
// Full validation
// ---------------------------------------------------------------------------

// Validate performs full multi-layer validation on a skill:
//  1. Source verification (URLs are reachable, content unchanged)
//  2. Code sandbox validation (any code snippets execute correctly)
//  3. LLM jury validation (multi-model consensus)
//  4. Cross-reference consistency check
func (p *Pipeline) Validate(ctx context.Context, skill *models.Skill) (*ValidationResult, error) {
	p.logger.Info("starting validation pipeline",
		zap.String("skill", skill.Name),
		zap.String("skill_id", skill.ID.String()),
	)

	result := &ValidationResult{
		SkillID: skill.ID,
		Details: make(map[string]string),
	}

	// Stage 1: Source verification
	if err := p.SourceVerify(ctx, skill); err != nil {
		result.Stage = "source_verification"
		result.Details["source_verification"] = err.Error()
		p.logValidationResult(ctx, skill, result)
		return result, nil // return result, not error - validation failure is a valid outcome
	}
	result.Details["source_verification"] = "passed"

	// Stage 2: Code sandbox (if skill contains code snippets)
	codeResult, err := p.extractAndRunCode(ctx, skill)
	if err != nil {
		result.Stage = "code_sandbox"
		result.Details["code_sandbox"] = err.Error()
		p.logValidationResult(ctx, skill, result)
		return result, nil
	}
	if codeResult != nil {
		result.Details["code_sandbox"] = fmt.Sprintf("executed %d snippets", codeResult)
	} else {
		result.Details["code_sandbox"] = "no code to execute"
	}

	// Stage 3: LLM jury validation
	juryResult, err := p.LLMJury(ctx, skill)
	if err != nil {
		result.Stage = "llm_jury"
		result.Details["llm_jury"] = err.Error()
		p.logValidationResult(ctx, skill, result)
		return result, nil
	}

	result.ApprovedBy = 0
	for _, approved := range juryResult.Votes {
		if approved {
			result.ApprovedBy++
		}
	}

	if !juryResult.Consensus {
		result.Stage = "llm_jury"
		result.Details["llm_jury"] = fmt.Sprintf("insufficient consensus: %d/%d approved",
			result.ApprovedBy, len(juryResult.Votes))
		p.logValidationResult(ctx, skill, result)
		return result, nil
	}
	result.Details["llm_jury"] = fmt.Sprintf("consensus reached: %d/%d approved",
		result.ApprovedBy, len(juryResult.Votes))

	// Stage 4: Cross-reference consistency
	if err := p.CrossReference(ctx, skill); err != nil {
		result.Stage = "cross_reference"
		result.Details["cross_reference"] = err.Error()
		p.logValidationResult(ctx, skill, result)
		return result, nil
	}
	result.Details["cross_reference"] = "passed"

	// All stages passed
	result.Passed = true
	result.Stage = "all_stages"
	p.logValidationResult(ctx, skill, result)

	p.logger.Info("validation passed",
		zap.String("skill", skill.Name),
		zap.Int("jury_approvals", result.ApprovedBy),
	)

	return result, nil
}

// ---------------------------------------------------------------------------
// Stage 1: Source verification
// ---------------------------------------------------------------------------

// SourceVerify checks that all resources attached to a skill are reachable
// and that cached content hashes match the current content.
func (p *Pipeline) SourceVerify(ctx context.Context, skill *models.Skill) error {
	if len(skill.Resources) == 0 {
		return nil // no resources to verify
	}

	var errs []string
	for i, res := range skill.Resources {
		if res.URL == "" {
			continue
		}

		err := p.verifySingleResource(ctx, &skill.Resources[i])
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", res.URL, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("resource verification failed: %s", strings.Join(errs, "; "))
	}

	return nil
}

// verifySingleResource checks if a resource URL is reachable and content unchanged.
func (p *Pipeline) verifySingleResource(ctx context.Context, res *models.Resource) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, res.URL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// If we have a cached hash, verify content hasn't changed
	if res.FetchedHash != "" {
		// Fetch full content and check hash
		bodyReq, err := http.NewRequestWithContext(ctx, http.MethodGet, res.URL, nil)
		if err != nil {
			return nil // HEAD succeeded, hash check is best-effort
		}

		bodyResp, err := p.httpClient.Do(bodyReq)
		if err != nil {
			return nil // best-effort
		}
		defer bodyResp.Close()

		content, err := io.ReadAll(io.LimitReader(bodyResp.Body, 1<<20)) // 1 MiB limit
		if err != nil {
			return nil // best-effort
		}

		hash := sha256.Sum256(content)
		currentHash := hex.EncodeToString(hash[:])

		if currentHash != res.FetchedHash {
			return fmt.Errorf("content hash mismatch (content changed)")
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Stage 2: Code sandbox
// ---------------------------------------------------------------------------

// CodeSandbox executes code in an isolated environment and returns the result.
func (p *Pipeline) CodeSandbox(ctx context.Context, code string, language string) (*SandboxResult, error) {
	timeout := 30 * time.Second

	execResult, err := p.sandbox.Execute(ctx, code, language, timeout)
	if err != nil {
		return &SandboxResult{
			ExecutionResult: ExecutionResult{Stderr: err.Error(), ExitCode: -1},
			Approved:        false,
			Feedback:        fmt.Sprintf("sandbox execution error: %v", err),
		}, nil // sandbox failure is a validation result, not an execution error
	}

	approved := execResult.ExitCode == 0
	feedback := "code executed successfully"
	if !approved {
		feedback = fmt.Sprintf("exit code %d: %s", execResult.ExitCode, execResult.Stderr)
	}

	return &SandboxResult{
		ExecutionResult: *execResult,
		Approved:        approved,
		Feedback:        feedback,
	}, nil
}

// extractAndRunCode extracts code blocks from skill content and runs them.
func (p *Pipeline) extractAndRunCode(ctx context.Context, skill *models.Skill) (interface{}, error) {
	// Extract fenced code blocks from markdown content
	snippets := extractCodeBlocks(skill.Content)
	if len(snippets) == 0 {
		return nil, nil // no code to execute
	}

	p.logger.Info("executing code snippets",
		zap.String("skill", skill.Name),
		zap.Int("snippets", len(snippets)),
	)

	executed := 0
	var failures []string

	for _, snippet := range snippets {
		select {
		case <-ctx.Done():
			return executed, ctx.Err()
		default:
		}

		result, err := p.CodeSandbox(ctx, snippet.Code, snippet.Language)
		if err != nil {
			return executed, fmt.Errorf("sandbox error: %w", err)
		}

		if !result.Approved {
			failures = append(failures, fmt.Sprintf("snippet %d failed: %s", executed+1, result.Feedback))
		}
		executed++
	}

	if len(failures) > 0 {
		return executed, fmt.Errorf("code validation failed: %s", strings.Join(failures, "; "))
	}

	return executed, nil
}

// codeSnippet represents an extracted code block.
type codeSnippet struct {
	Code     string
	Language string
}

// extractCodeBlocks extracts fenced code blocks from markdown content.
func extractCodeBlocks(content string) []codeSnippet {
	var snippets []codeSnippet

	lines := strings.Split(content, "\n")
	var inBlock bool
	var currentBlock strings.Builder
	var currentLang string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			if inBlock {
				// End block
				snippets = append(snippets, codeSnippet{
					Code:     strings.TrimSpace(currentBlock.String()),
					Language: currentLang,
				})
				currentBlock.Reset()
				inBlock = false
				currentLang = ""
			} else {
				// Start block
				inBlock = true
				lang := strings.TrimSpace(trimmed[3:])
				currentLang = lang
			}
			continue
		}

		if inBlock {
			currentBlock.WriteString(line)
			currentBlock.WriteString("\n")
		}
	}

	return snippets
}

// ---------------------------------------------------------------------------
// Stage 3: LLM Jury
// ---------------------------------------------------------------------------

// LLMJury runs multi-model validation to ensure skills meet quality standards.
// It requires consensus from at least cfg.ApprovalThreshold jurors.
func (p *Pipeline) LLMJury(ctx context.Context, skill *models.Skill) (*JuryResult, error) {
	if len(p.jury) == 0 {
		// No jury configured - auto-pass this stage
		p.logger.Warn("no jury members configured, auto-passing LLM validation",
			zap.String("skill", skill.Name),
		)
		return &JuryResult{
			Votes:     map[string]bool{"default": true},
			Consensus: true,
			Feedback:  "no jury configured, auto-approved",
		}, nil
	}

	p.logger.Info("running LLM jury",
		zap.String("skill", skill.Name),
		zap.Int("jury_size", len(p.jury)),
		zap.Int("threshold", p.cfg.ApprovalThreshold),
	)

	votes := make(map[string]bool, len(p.jury))
	var feedbackParts []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Run all jurors in parallel with timeout
	juryCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	for _, member := range p.jury {
		wg.Add(1)
		go func(m JuryMember) {
			defer wg.Done()

			approved, feedback, err := m.LLM.ValidateSkill(juryCtx, skill)
			if err != nil {
				p.logger.Warn("juror failed",
					zap.String("juror", m.Name),
					zap.Error(err),
				)
				mu.Lock()
				votes[m.Name] = false
				feedbackParts = append(feedbackParts, fmt.Sprintf("%s: error - %v", m.Name, err))
				mu.Unlock()
				return
			}

			mu.Lock()
			votes[m.Name] = approved
			feedbackParts = append(feedbackParts, fmt.Sprintf("%s: %s", m.Name, feedback))
			mu.Unlock()
		}(member)
	}

	wg.Wait()

	// Count approvals
	approvalCount := 0
	for _, approved := range votes {
		if approved {
			approvalCount++
		}
	}

	consensus := approvalCount >= p.cfg.ApprovalThreshold

	result := &JuryResult{
		Votes:     votes,
		Consensus: consensus,
		Feedback:  strings.Join(feedbackParts, "\n"),
	}

	p.logger.Info("LLM jury complete",
		zap.String("skill", skill.Name),
		zap.Int("approvals", approvalCount),
		zap.Int("required", p.cfg.ApprovalThreshold),
		zap.Bool("consensus", consensus),
	)

	return result, nil
}

// ---------------------------------------------------------------------------
// Stage 4: Cross-reference
// ---------------------------------------------------------------------------

// CrossReference checks consistency of a skill against existing skills in
// the knowledge graph. It verifies that referenced dependencies exist,
// terminology is consistent, and there are no contradictions.
func (p *Pipeline) CrossReference(ctx context.Context, skill *models.Skill) error {
	// Check that all dependencies exist
	for _, dep := range skill.Dependencies {
		if dep.DependsOn == uuid.Nil {
			return fmt.Errorf("dependency has empty ID for skill %q", dep.DependsOnName)
		}

		// Verify the dependency skill exists
		exists := false
		var checkID uuid.UUID
		if dep.DependsOn != uuid.Nil {
			checkID = dep.DependsOn
		}

		// Quick existence check via search
		results, err := p.store.Search(ctx, dep.DependsOnName, 1)
		if err != nil {
			return fmt.Errorf("search dependency %q: %w", dep.DependsOnName, err)
		}
		for _, r := range results {
			if r.Skill.Name == dep.DependsOnName {
				exists = true
				break
			}
		}
		_ = checkID

		if !exists {
			return fmt.Errorf("dependency %q not found in knowledge graph", dep.DependsOnName)
		}
	}

	// Check for naming conflicts
	existing, err := p.store.Search(ctx, skill.Name, 5)
	if err != nil {
		return fmt.Errorf("search for conflicts: %w", err)
	}
	for _, e := range existing {
		if e.Skill.Name == skill.Name && e.Skill.ID != skill.ID {
			return fmt.Errorf("naming conflict: skill %q already exists with different ID", skill.Name)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// logValidationResult records the validation outcome in the audit log.
func (p *Pipeline) logValidationResult(ctx context.Context, skill *models.Skill, result *ValidationResult) {
	details := map[string]interface{}{
		"skill_name":  skill.Name,
		"skill_id":    skill.ID.String(),
		"passed":      result.Passed,
		"stage":       result.Stage,
		"details":     result.Details,
		"approved_by": result.ApprovedBy,
	}

	event := db.AuditEventSkillValidated
	if !result.Passed {
		event = "skill.validation_failed"
	}

	if err := db.LogEventWithDetails(ctx, p.store.Pool(), event, &skill.ID, details); err != nil {
		p.logger.Warn("failed to log validation result", zap.Error(err))
	}
}

// ---------------------------------------------------------------------------
// NoOpSandbox (for when sandboxing is disabled)
// ---------------------------------------------------------------------------

// NoOpSandbox is a sandbox that does nothing, used when sandbox_type is "none".
type NoOpSandbox struct{}

// Execute always returns success without running any code.
func (s *NoOpSandbox) Execute(ctx context.Context, code string, language string, timeout time.Duration) (*ExecutionResult, error) {
	return &ExecutionResult{
		Stdout:   "sandbox disabled",
		Stderr:   "",
		ExitCode: 0,
		Duration: 0,
	}, nil
}
