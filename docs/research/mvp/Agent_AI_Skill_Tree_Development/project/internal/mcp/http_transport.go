package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ============================================================================
// HTTPTransport - MCP over HTTP with SSE streaming
// ============================================================================
//
// This transport wraps the MCP server for HTTP access, supporting:
//   - POST /mcp/v1/messages - Send JSON-RPC requests
//   - GET  /mcp/v1/sse      - Server-Sent Events for streaming
//   - GET  /mcp/v1/tools    - List available tools (REST fallback)
//   - POST /mcp/v1/tools/:name/call - Direct tool call (REST fallback)
//   - GET  /health          - Health check

// HTTPTransport provides HTTP/SSE access to the MCP server.
type HTTPTransport struct {
	server *MCPServer
	engine *gin.Engine
	srv    *http.Server
	logger *zap.Logger
}

// NewHTTPTransport creates a new HTTP transport.
func NewHTTPTransport(server *MCPServer) *HTTPTransport {
	gin.SetMode(gin.ReleaseMode)

	return &HTTPTransport{
		server: server,
		logger: server.logger.With(zap.String("transport", "http")),
	}
}

// Start begins listening on the specified address.
func (t *HTTPTransport) Start(addr string) error {
	t.engine = gin.New()
	t.engine.Use(gin.Recovery())
	t.engine.Use(t.loggingMiddleware())
	t.engine.Use(t.corsMiddleware())

	t.RegisterRoutes(t.engine)

	t.srv = &http.Server{
		Addr:         addr,
		Handler:      t.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	t.logger.Info("HTTP transport starting", zap.String("addr", addr))

	go func() {
		if err := t.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	return nil
}

// Shutdown gracefully stops the HTTP server.
func (t *HTTPTransport) Shutdown(ctx context.Context) error {
	if t.srv == nil {
		return nil
	}

	t.logger.Info("HTTP transport shutting down")
	return t.srv.Shutdown(ctx)
}

// RegisterRoutes sets up all HTTP routes.
func (t *HTTPTransport) RegisterRoutes(router *gin.Engine) {
	mcpGroup := router.Group("/mcp/v1")
	{
		mcpGroup.POST("/messages", t.handleJSONRPC)
		mcpGroup.GET("/sse", t.handleSSE)
		mcpGroup.GET("/tools", t.handleToolsListREST)
		mcpGroup.POST("/tools/:name/call", t.handleToolCallREST)
		mcpGroup.GET("/prompts", t.handlePromptsListREST)
		mcpGroup.GET("/prompts/:name", t.handlePromptsGetREST)
	}

	router.GET("/health", t.handleHealth)
	router.GET("/", t.handleRoot)
}

// handleJSONRPC processes JSON-RPC requests over HTTP.
func (t *HTTPTransport) handleJSONRPC(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &JSONRPCError{
				Code:    ErrCodeParseError,
				Message: "Failed to read request body",
			},
		})
		return
	}

	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &JSONRPCError{
				Code:    ErrCodeParseError,
				Message: "Invalid JSON",
				Data:    err.Error(),
			},
		})
		return
	}

	result, errResp := t.routeJSONRPC(c.Request.Context(), req)
	if errResp != nil {
		c.JSON(http.StatusOK, JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   errResp,
		})
		return
	}

	c.JSON(http.StatusOK, JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	})
}

// routeJSONRPC dispatches JSON-RPC requests to the appropriate handler.
func (t *HTTPTransport) routeJSONRPC(ctx context.Context, req JSONRPCRequest) (interface{}, *JSONRPCError) {
	switch req.Method {
	case "initialize":
		return t.handleInitializeHTTP(req.Params)
	case "notifications/initialized":
		return map[string]interface{}{}, nil
	case "tools/list":
		return t.handleToolsListInternal()
	case "tools/call":
		return t.handleToolsCallInternal(ctx, req.Params)
	case "prompts/list":
		return t.handlePromptsListInternal()
	case "prompts/get":
		return t.handlePromptsGetInternal(req.Params)
	case "ping":
		return map[string]interface{}{}, nil
	default:
		return nil, &JSONRPCError{
			Code:    ErrCodeMethodNotFound,
			Message: fmt.Sprintf("Method not found: %s", req.Method),
		}
	}
}

func (t *HTTPTransport) handleInitializeHTTP(params json.RawMessage) (interface{}, *JSONRPCError) {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{ListChanged: true},
		},
		ServerInfo: ServerInfo{
			Name:    "helix-knowledge-skill-system",
			Version: "1.0.0",
		},
	}
	return result, nil
}

func (t *HTTPTransport) handleToolsListInternal() (interface{}, *JSONRPCError) {
	tools := []map[string]interface{}{
		{
			"name":        "skill_search",
			"description": "Search the skill graph using text or vector similarity",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]string{"type": "string", "description": "Search query"},
					"limit": map[string]interface{}{"type": "integer", "description": "Max results", "default": 5},
				},
				"required": []string{"query"},
			},
		},
		{
			"name":        "skill_get",
			"description": "Retrieve complete skill record by name",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]string{"type": "string", "description": "Exact skill name"},
				},
				"required": []string{"name"},
			},
		},
		{
			"name":        "skill_tree",
			"description": "Get dependency tree for a skill",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":  map[string]string{"type": "string", "description": "Root skill name"},
					"depth": map[string]interface{}{"type": "integer", "description": "Max depth", "default": 5},
				},
				"required": []string{"name"},
			},
		},
		{
			"name":        "skill_create",
			"description": "Create or update a skill using TOML format",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"toml": map[string]string{"type": "string", "description": "TOML skill definition"},
				},
				"required": []string{"toml"},
			},
		},
		{
			"name":        "learn_from_project",
			"description": "Submit a project path for skill extraction",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_path": map[string]string{"type": "string", "description": "Project directory path"},
					"languages":    map[string]interface{}{"type": "array", "items": map[string]string{"type": "string"}},
				},
				"required": []string{"project_path"},
			},
		},
		{
			"name":        "missing_skills",
			"description": "Find gaps in the knowledge graph",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain": map[string]string{"type": "string", "description": "Filter by domain (optional)"},
				},
			},
		},
		{
			"name":        "get_coverage",
			"description": "Get coverage report for a domain",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain": map[string]string{"type": "string", "description": "Domain (optional)"},
				},
			},
		},
	}

	return map[string]interface{}{"tools": tools}, nil
}

func (t *HTTPTransport) handleToolsCallInternal(ctx context.Context, params json.RawMessage) (interface{}, *JSONRPCError) {
	var callParams CallToolParams
	if err := json.Unmarshal(params, &callParams); err != nil {
		return nil, &JSONRPCError{
			Code:    ErrCodeInvalidParams,
			Message: "Invalid tools/call params",
			Data:    err.Error(),
		}
	}

	// Execute the tool via the registered mcp-go tool handler.
	toolResult, err := t.server.dispatchTool(ctx, callParams.Name, callParams.Arguments)
	if err != nil {
		return nil, &JSONRPCError{
			Code:    ErrCodeInternalError,
			Message: fmt.Sprintf("Tool '%s' failed", callParams.Name),
			Data:    err.Error(),
		}
	}

	return map[string]interface{}{
		"content": toolResult.Content,
		"isError": toolResult.IsError,
	}, nil
}

func (t *HTTPTransport) handlePromptsListInternal() (interface{}, *JSONRPCError) {
	prompts := []map[string]interface{}{
		{"name": "system-prompt", "description": "System prompt for HelixKnowledge agents"},
		{"name": "skill-format", "description": "TOML skill format guide"},
	}
	return map[string]interface{}{"prompts": prompts}, nil
}

func (t *HTTPTransport) handlePromptsGetInternal(params json.RawMessage) (interface{}, *JSONRPCError) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, &JSONRPCError{Code: ErrCodeInvalidParams, Message: "Invalid params"}
	}

	var text string
	switch req.Name {
	case "system-prompt":
		text = GetSystemPrompt()
	case "skill-format":
		text = GetSkillFormatPrompt()
	default:
		return nil, &JSONRPCError{Code: ErrCodeInvalidParams, Message: fmt.Sprintf("Prompt not found: %s", req.Name)}
	}

	return map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": map[string]string{"type": "text", "text": text},
			},
		},
	}, nil
}

// handleSSE provides Server-Sent Events for streaming updates.
func (t *HTTPTransport) handleSSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	fmt.Fprintf(c.Writer, "event: connected\ndata: %s\n\n", `{"status":"connected","server":"helix-knowledge-skill-system"}`)
	flusher.Flush()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	clientGone := c.Request.Context().Done()

	for {
		select {
		case <-clientGone:
			t.logger.Debug("SSE client disconnected")
			return
		case <-ticker.C:
			fmt.Fprintf(c.Writer, "event: ping\ndata: %s\n\n", `{"time":"`+time.Now().Format(time.RFC3339)+`"}`)
			flusher.Flush()
		}
	}
}

func (t *HTTPTransport) handleToolsListREST(c *gin.Context) {
	result, errResp := t.handleToolsListInternal()
	if errResp != nil {
		c.JSON(http.StatusOK, gin.H{"error": errResp.Message})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (t *HTTPTransport) handleToolCallREST(c *gin.Context) {
	toolName := c.Param("name")

	var args map[string]interface{}
	if err := c.ShouldBindJSON(&args); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	params, _ := json.Marshal(CallToolParams{
		Name:      toolName,
		Arguments: args,
	})

	result, errResp := t.handleToolsCallInternal(c.Request.Context(), params)
	if errResp != nil {
		c.JSON(http.StatusOK, gin.H{"error": errResp.Message, "details": errResp.Data})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (t *HTTPTransport) handlePromptsListREST(c *gin.Context) {
	result, errResp := t.handlePromptsListInternal()
	if errResp != nil {
		c.JSON(http.StatusOK, gin.H{"error": errResp.Message})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (t *HTTPTransport) handlePromptsGetREST(c *gin.Context) {
	promptName := c.Param("name")
	params, _ := json.Marshal(map[string]string{"name": promptName})
	result, errResp := t.handlePromptsGetInternal(params)
	if errResp != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": errResp.Message})
		return
	}
	c.JSON(http.StatusOK, result)
}

// handleHealth returns health status.
func (t *HTTPTransport) handleHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbStatus := "ok"
	if t.server.pool != nil {
		if err := t.server.pool.Health(ctx); err != nil {
			dbStatus = "error: " + err.Error()
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"server":    "helix-knowledge-skill-system",
		"version":   "1.0.0",
		"database":  dbStatus,
		"transport": "http",
		"tools":     7,
	})
}

// handleRoot returns basic server info.
func (t *HTTPTransport) handleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":        "HelixKnowledge Skill Graph System",
		"version":     "1.0.0",
		"description": "MCP server for AI agent skill management",
		"endpoints": map[string]string{
			"jsonrpc":    "POST /mcp/v1/messages",
			"sse":        "GET /mcp/v1/sse",
			"tools":      "GET /mcp/v1/tools",
			"tool_call":  "POST /mcp/v1/tools/:name/call",
			"prompts":    "GET /mcp/v1/prompts",
			"prompt_get": "GET /mcp/v1/prompts/:name",
			"health":     "GET /health",
		},
	})
}

// loggingMiddleware logs HTTP requests.
func (t *HTTPTransport) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		t.logger.Debug("HTTP request",
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
		)
	}
}

// corsMiddleware adds CORS headers for cross-origin requests.
func (t *HTTPTransport) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
