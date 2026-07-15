// cmd/server is the main entry point for the HelixKnowledge Skill System API and MCP server.
//
// Usage:
//
//	go run ./cmd/server [--config path/to/config.toml] [--mcp stdio|http|both]
//
// Modes:
//   - API mode (default): Runs the HTTP REST API server on the configured port
//   - MCP stdio mode (--mcp stdio): Runs the MCP server over stdin/stdout for CLI agents
//   - MCP HTTP mode (--mcp http): Runs the MCP server over HTTP/SSE
//   - MCP both mode (--mcp both): Runs both stdio and HTTP transports
//
// Environment variables:
//   HELIX_DB_HOST, HELIX_DB_PORT, HELIX_DB_NAME, HELIX_DB_USER, HELIX_DB_PASSWORD
//   HELIX_LOG_LEVEL, HELIX_MCP_TRANSPORT
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
	mcpTransport := flag.String("mcp", "", "MCP transport: stdio, http, both (overrides config)")
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
		// Pure MCP stdio mode - no HTTP API
		logger.Info("Running in MCP stdio mode (stdout reserved for JSON-RPC)")
		if err := mcpServer.RunStdio(); err != nil {
			logger.Fatal("MCP stdio server failed", zap.Error(err))
		}

	case "http":
		// MCP HTTP mode - serve MCP over HTTP/SSE + REST API
		httpAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)
		if err := mcpServer.RunHTTP(httpAddr); err != nil {
			logger.Fatal("MCP HTTP server failed", zap.Error(err))
		}
		setupAPI(cfg, pool, skillStore, skillRegistry, logger)
		waitForShutdown(ctx, logger, mcpServer)

	case "both":
		// Both stdio (blocking) and HTTP (background)
		httpAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)
		setupAPI(cfg, pool, skillStore, skillRegistry, logger)
		if err := mcpServer.RunBoth(httpAddr); err != nil {
			logger.Fatal("MCP server failed", zap.Error(err))
		}

	default:
		// Standard API mode with MCP over HTTP
		httpAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)
		if err := mcpServer.RunHTTP(httpAddr); err != nil {
			logger.Fatal("MCP HTTP server failed", zap.Error(err))
		}
		setupAPI(cfg, pool, skillStore, skillRegistry, logger)
		waitForShutdown(ctx, logger, mcpServer)
	}
}

// setupAPI configures the Gin REST API server.
func setupAPI(cfg *config.Config, pool *db.Pool, store *skill.Store, reg *registry.Registry, logger *zap.Logger) *gin.Engine {
	if cfg.Server.EnableBrotli {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(apiLoggingMiddleware(logger))
	router.Use(corsMiddleware())

	// Health check
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

	// Skills API
	skills := router.Group("/api/v1/skills")
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
	router.GET("/api/v1/coverage", func(c *gin.Context) {
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
	router.GET("/api/v1/missing", func(c *gin.Context) {
		ctx := c.Request.Context()
		domain := c.Query("domain")
		entries, err := store.GetMissingSkills(ctx, domain)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"missing_skills": entries, "count": len(entries)})
	})

	// Server info
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "HelixKnowledge Skill Graph System",
			"version":     "1.0.0",
			"description": "API and MCP server for AI agent skill management",
			"endpoints": []string{
				"GET  /health",
				"GET  /api/v1/skills",
				"GET  /api/v1/skills/search?q=query",
				"GET  /api/v1/skills/:name",
				"GET  /api/v1/skills/:name/tree",
				"GET  /api/v1/coverage?domain=optional",
				"GET  /api/v1/missing?domain=optional",
				"POST /mcp/v1/messages (JSON-RPC)",
				"GET  /mcp/v1/sse (SSE streaming)",
			},
		})
	})

	// Start the HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)
	logger.Info("API server starting", zap.String("addr", addr))

	go func() {
		srv := &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("API server error", zap.Error(err))
		}
	}()

	return router
}

// waitForShutdown blocks until a shutdown signal is received.
func waitForShutdown(ctx context.Context, logger *zap.Logger, mcpServer *mcp.MCPServer) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	logger.Info("Waiting for shutdown signal (SIGTERM/SIGINT)")
	sig := <-sigCh
	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := mcpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Shutdown error", zap.Error(err))
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

// corsMiddleware adds CORS headers.
func corsMiddleware() gin.HandlerFunc {
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
