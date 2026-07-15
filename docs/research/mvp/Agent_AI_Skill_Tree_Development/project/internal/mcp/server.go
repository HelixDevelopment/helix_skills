// Package mcp implements the Model Context Protocol server for the HelixKnowledge
// Skill Graph System. It provides 7 tools for AI agents to query, create, and
// manage skills through stdio and HTTP transports.
package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/helixdevelopment/skill-system/internal/config"
	"github.com/helixdevelopment/skill-system/internal/db"
	"github.com/helixdevelopment/skill-system/internal/registry"
	"github.com/helixdevelopment/skill-system/internal/skill"

	mcp_go "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// MCPServer wraps the mcp-go server with application-specific dependencies.
type MCPServer struct {
	server     *server.MCPServer
	skillStore *skill.Store
	registry   *registry.Registry
	pool       *db.Pool
	cfg        *config.Config
	logger     *zap.Logger
	transport  string // "stdio" | "http" | "both"
	stdio      *StdioTransport
	http       *HTTPTransport
}

// NewMCPServer creates a new MCP server with all dependencies.
func NewMCPServer(pool *db.Pool, store *skill.Store, reg *registry.Registry, cfg *config.Config, logger *zap.Logger) *MCPServer {
	mcpServer := server.NewMCPServer(
		"helix-knowledge-skill-system",
		"1.0.0",
		server.WithLogging(),
		server.WithResourceCapabilities(true, true),
	)

	return &MCPServer{
		server:     mcpServer,
		skillStore: store,
		registry:   reg,
		pool:       pool,
		cfg:        cfg,
		logger:     logger,
		transport:  cfg.MCP.Transport,
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

	// The HTTP transport no longer owns a listener — its /mcp/v1 routes are
	// mounted on the shared API server (see RegisterHTTPRoutes), whose lifecycle
	// the caller (cmd/server) manages. Only stdio needs an explicit stop here.

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
