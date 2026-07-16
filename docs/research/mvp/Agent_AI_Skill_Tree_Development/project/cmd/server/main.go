// cmd/server is the main entry point for the HelixKnowledge Skill System API and MCP server.
//
// Usage:
//
//	go run ./cmd/server [--config path/to/config.toml] [--mcp stdio|http|both|acp]
//
// Modes:
//   - API mode (default): Runs the HTTP REST API server on the configured port
//   - MCP stdio mode (--mcp stdio): Runs the MCP server over stdin/stdout for CLI agents
//   - MCP HTTP mode (--mcp http): Runs the MCP server over HTTP/SSE
//   - MCP both mode (--mcp both): Runs both stdio and HTTP transports
//   - MCP acp mode (--mcp acp): Runs the Agent Client Protocol (ACP) adapter
//     over stdin/stdout for ACP-speaking CLI agents (no HTTP listener)
//
// Environment variables:
//
//	HELIX_DB_HOST, HELIX_DB_PORT, HELIX_DB_NAME, HELIX_DB_USER, HELIX_DB_PASSWORD
//	HELIX_LOG_LEVEL, HELIX_MCP_TRANSPORT
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/helixdevelopment/skill-system/internal/api"
	"github.com/helixdevelopment/skill-system/internal/config"
	"github.com/helixdevelopment/skill-system/internal/db"
	"github.com/helixdevelopment/skill-system/internal/mcp"
	"github.com/helixdevelopment/skill-system/internal/registry"
	"github.com/helixdevelopment/skill-system/internal/skill"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "", "Path to TOML config file (optional)")
	mcpTransport := flag.String("mcp", "", "MCP transport: stdio, http, both, acp (overrides config)")
	flag.Parse()

	// 1. Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Override MCP transport from CLI flag
	if *mcpTransport != "" {
		cfg.MCP.Transport = *mcpTransport
	}

	// 2. Initialize logger
	logger, err := newLogger(cfg.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("HelixKnowledge Skill System starting",
		zap.String("version", "1.0.0"),
		zap.String("mcp_transport", cfg.MCP.Transport),
	)

	// 3. Connect to database (db.New takes config value, no context needed)
	pool, err := db.New(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	logger.Info("Database connected",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("database", cfg.Database.Database),
	)

	// 4. Run migrations
	ctx := context.Background()
	if err := db.Migrate(ctx, pool, "./migrations"); err != nil {
		logger.Warn("Migration failed", zap.Error(err))
	} else {
		logger.Info("Migrations completed")
	}

	// 5. Initialize skill store and registry
	skillStore := skill.NewStore(pool)
	skillRegistry := registry.NewRegistry(pool)

	logger.Info("Skill store and registry initialized")

	// 6. Setup MCP server
	mcpServer := mcp.NewMCPServer(pool, skillStore, skillRegistry, cfg, logger)
	mcpServer.RegisterTools()

	// 7. Run based on transport mode
	mode := cfg.MCP.Transport

	switch mode {
	case "stdio":
		// Pure MCP stdio mode - no HTTP listener at all.
		logger.Info("Running in MCP stdio mode (stdout reserved for JSON-RPC)")
		if _, err := runBlockingTransport(mode, mcpServer); err != nil {
			logger.Fatal("MCP stdio server failed", zap.Error(err))
		}

	case "acp":
		// Pure ACP (Agent Client Protocol) mode - no HTTP listener at all.
		// Wires the previously-unwired internal/mcp/acp_adapter.go in as a
		// selectable transport, mirroring the "stdio" case above exactly
		// (same blocking, stdout-reserved-for-JSON-RPC discipline; only the
		// wire protocol dialect differs).
		logger.Info("MCP acp mode (stdout reserved for JSON-RPC)")
		if _, err := runBlockingTransport(mode, mcpServer); err != nil {
			logger.Fatal("MCP acp server failed", zap.Error(err))
		}

	case "both":
		// ONE hardened HTTP listener (background goroutine, MCP routes mounted +
		// auth-guarded) PLUS stdio in the foreground (blocking). There is no
		// second HTTP listener: the MCP HTTP surface lives on the same router as
		// /api/v1 under the same fail-closed policy.
		srv := setupAPI(cfg, pool, skillStore, skillRegistry, mcpServer, logger)
		if err := mcpServer.RunStdio(); err != nil {
			logger.Fatal("MCP stdio server failed", zap.Error(err))
		}
		// stdio ended (EOF/stop): drain the HTTP server and the MCP server.
		gracefulShutdown(ctx, logger, mcpServer, srv)

	default:
		// "http" and the standard API mode both serve a SINGLE hardened HTTP
		// listener with the MCP routes mounted and auth-guarded, then block
		// until a shutdown signal. Exactly one server binds the port.
		srv := setupAPI(cfg, pool, skillStore, skillRegistry, mcpServer, logger)
		waitForShutdown(ctx, logger, mcpServer, srv)
	}
}

// buildRouter assembles the SINGLE hardened Gin router used by every
// HTTP-serving mode. It wires the config-driven CORS allowlist, resolves the
// fail-closed API-key auth ONCE, and guards BOTH the /api/v1 data routes AND the
// mounted MCP /mcp/v1 routes with that same middleware. /health and / are the
// only open routes. No listener is started here (see setupAPI) so this assembly
// is directly unit-testable.
func buildRouter(cfg *config.Config, pool *db.Pool, store *skill.Store, reg *registry.Registry, mcpServer *mcp.MCPServer, logger *zap.Logger) *gin.Engine {
	if cfg.Server.EnableBrotli {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	// Trust NO forwarding proxy (F1). In the R15 no-proxy single-node topology
	// the app is reached directly, so c.ClientIP() MUST resolve to the real TCP
	// socket peer and NEVER to a client-supplied X-Forwarded-For / X-Real-IP. Gin
	// otherwise trusts all proxies by default, which would let an unauthenticated
	// caller spoof its rate-limit identity by rotating a forwarding header.
	// ForwardedByClientIP=false short-circuits header parsing; SetTrustedProxies(nil)
	// additionally empties the trusted set (belt-and-suspenders).
	router.ForwardedByClientIP = false
	if err := router.SetTrustedProxies(nil); err != nil {
		logger.Error("failed to clear trusted proxies (rate-limit identity hardening)", zap.Error(err))
	}
	router.Use(gin.Recovery())
	router.Use(apiLoggingMiddleware(logger))
	// Hardened, config-driven CORS allowlist (internal/api.CORS): an empty
	// allowlist is fail-closed and no wildcard "*" origin is ever emitted with
	// credentials. This replaces the previous wildcard corsMiddleware().
	router.Use(api.CORS(cfg.Server.AllowedOrigins))

	// Per-client token-bucket rate limiting (§G22), installed BEFORE auth so a
	// flood is throttled with 429 ahead of any credential work. Keyed ONLY on the
	// real socket peer (never an attacker-controlled header, F1) with a map that
	// is hard-bounded by MaxClients (F2), so one client cannot starve another and
	// a distinct-IP flood cannot grow the tracking map without bound. Off only
	// when explicitly disabled in config.
	if cfg.Server.RateLimit.Enabled {
		limiter := api.NewRateLimiter(api.RateLimitConfig{
			RequestsPerSecond: cfg.Server.RateLimit.RequestsPerSecond,
			Burst:             cfg.Server.RateLimit.Burst,
			TTL:               cfg.Server.RateLimit.TTL,
			MaxClients:        cfg.Server.RateLimit.MaxClients,
		})
		router.Use(limiter.Middleware())
	}

	// Request-body cap (§G22): reject an oversized body with 413 BEFORE auth so
	// an unauthenticated flood of huge bodies cannot exhaust memory. A
	// non-positive config value falls back to the 100 MiB default.
	maxBody := cfg.Server.MaxRequestBodyBytes
	if maxBody <= 0 {
		maxBody = api.DefaultMaxBodyBytes
	}
	router.Use(api.MaxBodySize(maxBody))

	// Resolve the fail-closed auth middleware ONCE and share the SAME instance
	// across BOTH the /api/v1 data routes and the mounted MCP /mcp/v1 routes so
	// the two surfaces enforce identical authentication (and the startup log
	// fires exactly once). ResolveAPIKeyAuth is fail-CLOSED: with no API keys
	// and auth not explicitly disabled it rejects every request (503) rather
	// than serving these routes open. nil is returned ONLY in the explicit
	// auth-disabled mode.
	authMW := api.ResolveAPIKeyAuth(cfg.Server.APIKeys, cfg.Server.AuthDisabled, logger)

	// Health check (open)
	router.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		dbStatus := "ok"
		if err := pool.Health(ctx); err != nil {
			dbStatus = "error: " + err.Error()
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"server":    "helix-knowledge-skill-system",
			"version":   "1.0.0",
			"database":  dbStatus,
			"transport": cfg.MCP.Transport,
		})
	})

	// All /api/v1 data routes are authenticated under the fail-closed guard.
	v1 := router.Group("/api/v1")
	if authMW != nil {
		v1.Use(authMW)
	}

	// Skills API
	skills := v1.Group("/skills")
	{
		skills.GET("", func(c *gin.Context) {
			ctx := c.Request.Context()
			skills, err := store.ListSkills(ctx, "", 100, 0)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"skills": skills, "count": len(skills)})
		})

		skills.GET("/search", func(c *gin.Context) {
			ctx := c.Request.Context()
			query := c.Query("q")
			if query == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
				return
			}
			results, err := store.Search(ctx, query, 10)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"results": results, "query": query})
		})

		skills.GET("/:name", func(c *gin.Context) {
			ctx := c.Request.Context()
			name := c.Param("name")
			skill, err := store.GetByName(ctx, name)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, skill)
		})

		skills.GET("/:name/tree", func(c *gin.Context) {
			ctx := c.Request.Context()
			name := c.Param("name")
			depth := 5
			tree, err := store.GetTree(ctx, name, depth)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, tree)
		})
	}

	// Coverage API
	v1.GET("/coverage", func(c *gin.Context) {
		ctx := c.Request.Context()
		domain := c.Query("domain")
		report, err := reg.GetCoverageReport(ctx, domain)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, report)
	})

	// Missing skills API
	v1.GET("/missing", func(c *gin.Context) {
		ctx := c.Request.Context()
		domain := c.Query("domain")
		entries, err := store.GetMissingSkills(ctx, domain)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"missing_skills": entries, "count": len(entries)})
	})

	// Mount the MCP HTTP routes (/mcp/v1/*) onto THIS same router, behind the
	// SAME fail-closed auth guard as /api/v1. Previously these ran on a SECOND
	// http.Server bound to the identical host:port with wildcard CORS and NO
	// authentication; whichever server won the bind decided the live security
	// posture, and the hardened /api/v1 routes 404'd when the MCP one won. One
	// listener now serves both surfaces under one hardened policy.
	mcpServer.RegisterHTTPRoutes(router, authMW)

	// Server info (open) — reflects the real, live, auth-guarded route surface.
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "HelixKnowledge Skill Graph System",
			"version":     "1.0.0",
			"description": "API and MCP server for AI agent skill management",
			"endpoints": []string{
				"GET  /health",
				"GET  /api/v1/skills (auth required)",
				"GET  /api/v1/skills/search?q=query (auth required)",
				"GET  /api/v1/skills/:name (auth required)",
				"GET  /api/v1/skills/:name/tree (auth required)",
				"GET  /api/v1/coverage?domain=optional (auth required)",
				"GET  /api/v1/missing?domain=optional (auth required)",
				"POST /mcp/v1/messages (JSON-RPC, auth required)",
				"GET  /mcp/v1/sse (SSE streaming, auth required)",
				"GET  /mcp/v1/tools (auth required)",
				"POST /mcp/v1/tools/:name/call (auth required)",
				"GET  /mcp/v1/prompts (auth required)",
				"GET  /mcp/v1/prompts/:name (auth required)",
			},
		})
	})

	return router
}

// setupAPI builds the single hardened router and starts the ONE HTTP listener,
// returning the *http.Server so the caller can shut it down gracefully. A bind
// failure is FATAL: a server that cannot bind its port is NOT serving, and the
// process must not continue as if it were (that is exactly the silent
// "fixed-but-not-live" failure this remediation closes).
func setupAPI(cfg *config.Config, pool *db.Pool, store *skill.Store, reg *registry.Registry, mcpServer *mcp.MCPServer, logger *zap.Logger) *http.Server {
	router := buildRouter(cfg, pool, store, reg, mcpServer, logger)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("API server starting (single hardened listener; MCP mounted + auth-guarded)",
		zap.String("addr", addr))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// A failed bind (e.g. address already in use) means we are NOT
			// serving — fail hard instead of silently continuing.
			logger.Fatal("API server failed", zap.String("addr", addr), zap.Error(err))
		}
	}()

	return srv
}

// waitForShutdown blocks until a shutdown signal is received, then drains both
// the HTTP server and the MCP server.
func waitForShutdown(ctx context.Context, logger *zap.Logger, mcpServer *mcp.MCPServer, srv *http.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	logger.Info("Waiting for shutdown signal (SIGTERM/SIGINT)")
	sig := <-sigCh
	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	gracefulShutdown(ctx, logger, mcpServer, srv)
}

// gracefulShutdown drains the HTTP listener (triggering ListenAndServe to return
// http.ErrServerClosed) and then the MCP server, within a bounded timeout.
func gracefulShutdown(ctx context.Context, logger *zap.Logger, mcpServer *mcp.MCPServer, srv *http.Server) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if srv != nil {
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("HTTP server shutdown error", zap.Error(err))
		}
	}
	if err := mcpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("MCP shutdown error", zap.Error(err))
	}

	logger.Info("Graceful shutdown complete")
}

// newLogger creates a zap logger based on configuration.
// When running in stdio mode, all logs are written to stderr to avoid
// interfering with stdout which is reserved for JSON-RPC messages.
func newLogger(cfg config.LoggingConfig) (*zap.Logger, error) {
	level := zap.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Always write to stderr to avoid interfering with stdout (JSON-RPC)
	core := zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), level)
	return zap.New(core, zap.AddCaller()), nil
}

// apiLoggingMiddleware logs API requests.
func apiLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Debug("API request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
		)
	}
}
