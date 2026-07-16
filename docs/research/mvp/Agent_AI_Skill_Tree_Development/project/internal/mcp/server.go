// Package mcp implements the Model Context Protocol server for the HelixKnowledge
// Skill Graph System. It provides 7 tools for AI agents to query, create, and
// manage skills through stdio and HTTP transports.
package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/helixdevelopment/skill-system/internal/config"
	"github.com/helixdevelopment/skill-system/internal/db"
	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/helixdevelopment/skill-system/internal/registry"
	"github.com/helixdevelopment/skill-system/internal/skill"
	"github.com/helixdevelopment/skill-system/internal/validation"

	mcp_go "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// skillValidator runs the fail-closed, non-executing validation pipeline against
// a submitted skill before it is persisted (§G03). Satisfied by *validation.Pipeline.
type skillValidator interface {
	Validate(ctx context.Context, s *models.Skill) (*validation.ValidationResult, error)
}

// MCPServer wraps the mcp-go server with application-specific dependencies.
type MCPServer struct {
	server            *server.MCPServer
	skillStore        *skill.Store
	registry          *registry.Registry
	pool              *db.Pool
	cfg               *config.Config
	logger            *zap.Logger
	transport         string // "stdio" | "http" | "both" | "acp"
	stdio             *StdioTransport
	http              *HTTPTransport
	acp               *ACPAdapter
	validator         skillValidator
	validationEnabled bool
}

// NewMCPServer creates a new MCP server with all dependencies.
func NewMCPServer(pool *db.Pool, store *skill.Store, reg *registry.Registry, cfg *config.Config, logger *zap.Logger) *MCPServer {
	mcpServer := server.NewMCPServer(
		"helix-knowledge-skill-system",
		"1.0.0",
		server.WithLogging(),
		server.WithResourceCapabilities(true, true),
	)

	// §G29: wire the query-side embedder onto the shared skill Store so its
	// Search becomes a genuine hybrid (vector KNN + trigram) search. This single
	// Store instance backs the MCP skill_search tool, the REST /search route, and
	// the validation pipeline's dedup lookup, so all three semantic-search paths
	// light up together the moment an embedding provider is configured — and stay
	// keyword-only (no wasted embedding calls) when it is not.
	//
	// "Is the provider configured" is derived SOLELY from
	// db.NewEmbedderFromConfig's own error return (Fable code-review
	// remediation, finding 6a) -- there is no second, hand-maintained
	// per-provider check here that could silently drift from the factory's
	// policy. db.NewEmbedderFromConfig is itself fail-closed (errors on a
	// missing openai api_key / local local_endpoint) rather than construct an
	// embedder guaranteed to fail its first real request.
	if store != nil {
		// Re-review remediation (MAJOR finding, post-G29): wire the real
		// application logger into the Store so its own diagnostics (currently
		// just warnEmbeddingDegraded) reach a real sink at runtime instead of
		// the package-level zap.L() no-op default -- mirrors the WithEmbedder
		// call immediately below and is unconditional (unlike the embedder,
		// which only wires when a provider is configured) because the logger
		// is always available here and Search's degrade-to-keyword-only path
		// can fire regardless of whether an embedder ends up configured.
		store.WithLogger(logger)
		if emb, err := db.NewEmbedderFromConfig(cfg.Embedding); err == nil {
			store.WithEmbedder(emb)
			// NOTE (§11.4.6 honest wording, finding 2): no production
			// ingestion path populates skills.embedding yet -- store.Create
			// never sets the column, and embedding-population is a separate,
			// tracked follow-up (register item pending). Until that lands,
			// this wiring makes the vector leg of Search live, but semantic
			// recall only actually surfaces a skill once SOME out-of-band
			// process has populated ITS embedding column directly; every
			// skill created through the normal write path keeps a NULL
			// embedding and is retrievable via the trigram leg only. This log
			// line must not overclaim full hybrid coverage.
			logger.Info("hybrid skill search wired (§G29): semantic recall is active only for skills with a populated embedding; embedding ingestion is not yet wired and is tracked separately",
				zap.String("embedding_provider", cfg.Embedding.Provider))
		} else {
			logger.Debug("hybrid skill search: no embedding provider configured, using keyword-only search (§G29)", zap.Error(err))
		}
	}

	return &MCPServer{
		server:            mcpServer,
		skillStore:        store,
		registry:          reg,
		pool:              pool,
		cfg:               cfg,
		logger:            logger,
		transport:         cfg.MCP.Transport,
		validator:         validation.NewPipeline(store, cfg.Validation, logger),
		validationEnabled: cfg.Validation.Enabled,
	}
}

// buildSkillFromTOML parses a TOML skill definition into an in-memory model for
// validation ONLY (no persistence). Persistence happens via ImportFromTOML.
//
// buildSkillFromTOML is validation-only (kind is not consumed downstream
// here); a future kind-aware validator rule would need Kind set -- tracked
// for that rule's landing (§11.4.6).
func buildSkillFromTOML(data []byte) (*models.Skill, error) {
	var w models.TOMLSkillWrapper
	if err := toml.Unmarshal(data, &w); err != nil {
		return nil, err
	}
	metaJSON, _ := json.Marshal(w.Skill.Metadata)
	m := &models.Skill{
		ID:          uuid.New(),
		Name:        w.Skill.Name,
		Version:     w.Skill.Version,
		Title:       w.Skill.Title,
		Description: w.Skill.Description,
		Content:     w.Skill.Content,
		Metadata:    metaJSON,
		Status:      models.SkillStatusDraft,
	}
	for _, r := range w.Skill.Resources {
		m.Resources = append(m.Resources, models.Resource{
			ID:           uuid.New(),
			URL:          r.URL,
			Title:        r.Title,
			ResourceType: r.ResourceType,
		})
	}
	return m, nil
}

// validateForCreate runs the fail-closed validation pipeline on a submitted TOML
// skill BEFORE it is persisted (§G03 request-path), returning a machine-readable
// summary for the tool response. It never executes untrusted code.
func (s *MCPServer) validateForCreate(ctx context.Context, tomlData []byte) map[string]interface{} {
	if !s.validationEnabled || s.validator == nil {
		return map[string]interface{}{"ran": false, "reason": "validation disabled"}
	}
	model, err := buildSkillFromTOML(tomlData)
	if err != nil {
		return map[string]interface{}{"ran": false, "reason": "toml parse error"}
	}
	vr, err := s.validator.Validate(ctx, model)
	if err != nil {
		s.logger.Warn("skill_create validation error", zap.String("skill", model.Name), zap.Error(err))
		return map[string]interface{}{"ran": true, "error": err.Error()}
	}
	s.logger.Info("skill_create validated",
		zap.String("skill", model.Name),
		zap.Bool("passed", vr.Passed),
		zap.String("stage", vr.Stage),
	)
	return map[string]interface{}{
		"ran":         true,
		"passed":      vr.Passed,
		"stage":       vr.Stage,
		"approved_by": vr.ApprovedBy,
		"stages":      vr.Stages,
	}
}

// RegisterTools sets up all 7 MCP tool handlers.
func (s *MCPServer) RegisterTools() {
	s.registerSkillSearch()
	s.registerSkillGet()
	s.registerSkillTree()
	s.registerSkillCreate()
	s.registerLearnFromProject()
	s.registerMissingSkills()
	s.registerGetCoverage()

	s.logger.Info("All 7 MCP tools registered",
		zap.String("transport", s.transport),
	)
}

// Server returns the underlying mcp-go server for testing or extension.
func (s *MCPServer) Server() *server.MCPServer {
	return s.server
}

// dispatchTool executes a registered tool by name through the underlying
// mcp-go server's handler and returns the tool result. The custom transports
// (stdio, HTTP, ACP) use this to route their JSON-RPC tool calls into the same
// handlers registered via AddTool, keeping a single execution path.
//
// mcp-go v0.56 has no public CallTool on *server.MCPServer; the supported way
// to invoke a registered tool in-process is to look it up with GetTool and
// call its exported Handler.
func (s *MCPServer) dispatchTool(ctx context.Context, name string, arguments any) (*mcp_go.CallToolResult, error) {
	st := s.server.GetTool(name)
	if st == nil {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	return st.Handler(ctx, mcp_go.CallToolRequest{
		Params: mcp_go.CallToolParams{
			Name:      name,
			Arguments: arguments,
		},
	})
}

// RunStdio starts the stdio transport for CLI agents (blocking).
// All logs are written to stderr only - stdout is reserved for JSON-RPC.
func (s *MCPServer) RunStdio() error {
	s.stdio = NewStdioTransport(s)
	s.logger.Info("Starting MCP stdio transport",
		zap.String("note", "All output on stdout is JSON-RPC; logs go to stderr"),
	)
	return s.stdio.Run()
}

// RunACP starts the ACP (Agent Client Protocol) adapter for CLI agents that
// speak ACP JSON-RPC over stdio instead of raw MCP JSON-RPC (blocking).
// Exactly like RunStdio, ALL logs are written to stderr only -- stdout is
// reserved for ACP JSON-RPC responses (see acp_adapter.go's writeStdout,
// which mirrors stdio.go's own stdout discipline).
func (s *MCPServer) RunACP() error {
	s.acp = NewACPAdapter(s)
	s.logger.Info("Starting MCP acp adapter",
		zap.String("note", "All output on stdout is ACP JSON-RPC; logs go to stderr"),
	)
	return s.acp.Run()
}

// RegisterHTTPRoutes mounts the MCP HTTP routes (/mcp/v1/*) onto the provided
// shared Gin router, guarded by authMW. This replaces the previous standalone
// MCP HTTP listener: the process now serves exactly ONE HTTP listener (the
// hardened API server in cmd/server), eliminating the two-servers-one-port race
// that could expose an unauthenticated, wildcard-CORS MCP surface whenever the
// MCP listener won the bind.
//
// authMW is the SAME fail-closed middleware guarding /api/v1 (from
// api.ResolveAPIKeyAuth). It is nil ONLY in the explicit auth-disabled mode.
func (s *MCPServer) RegisterHTTPRoutes(router *gin.Engine, authMW gin.HandlerFunc) {
	s.http = NewHTTPTransport(s)
	s.http.RegisterRoutes(router, authMW)
	s.logger.Info("MCP HTTP routes mounted on shared router",
		zap.Bool("auth_guarded", authMW != nil),
	)
}

// Shutdown gracefully stops the server and its transports.
func (s *MCPServer) Shutdown(_ context.Context) error {
	s.logger.Info("Shutting down MCP server")

	if s.stdio != nil {
		s.stdio.Stop()
	}
	if s.acp != nil {
		s.acp.Stop()
	}

	// The HTTP transport no longer owns a listener — its /mcp/v1 routes are
	// mounted on the shared API server (see RegisterHTTPRoutes), whose lifecycle
	// the caller (cmd/server) manages. Only stdio/acp need an explicit stop here.

	s.logger.Info("MCP server shutdown complete")
	return nil
}

// newToolResult creates a successful MCP tool result with JSON content.
func (s *MCPServer) newToolResult(data interface{}) *mcp_go.CallToolResult {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		s.logger.Error("Failed to marshal tool result", zap.Error(err))
		return &mcp_go.CallToolResult{
			Content: []mcp_go.Content{
				mcp_go.TextContent{
					Type: "text",
					Text: fmt.Sprintf(`{"error": "failed to serialize result: %v"}`, err),
				},
			},
		}
	}

	return &mcp_go.CallToolResult{
		Content: []mcp_go.Content{
			mcp_go.TextContent{
				Type: "text",
				Text: string(jsonBytes),
			},
		},
	}
}

// newToolResultRaw creates a result from pre-formatted text (for raw string output).
func (s *MCPServer) newToolResultRaw(text string) *mcp_go.CallToolResult {
	return &mcp_go.CallToolResult{
		Content: []mcp_go.Content{
			mcp_go.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}
}

// newToolError creates an error MCP tool result.
func (s *MCPServer) newToolError(errMsg string) *mcp_go.CallToolResult {
	return &mcp_go.CallToolResult{
		Content: []mcp_go.Content{
			mcp_go.TextContent{
				Type: "text",
				Text: fmt.Sprintf(`{"error": %q}`, errMsg),
			},
		},
		IsError: true,
	}
}
