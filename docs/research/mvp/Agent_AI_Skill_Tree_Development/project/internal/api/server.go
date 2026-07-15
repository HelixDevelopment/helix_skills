package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go/http3"
	"go.uber.org/zap"

	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/helixdevelopment/skill-system/internal/validation"
)

// SkillValidator runs the fail-closed, non-executing validation pipeline against
// a submitted skill before it is persisted (§G03 request-path). It is satisfied
// by *validation.Pipeline.
type SkillValidator interface {
	Validate(ctx context.Context, s *models.Skill) (*validation.ValidationResult, error)
}

// ServerConfig holds the server-specific configuration.
type ServerConfig struct {
	Host         string
	HTTPPort     int
	HTTP3Port    int
	EnableHTTP3  bool
	EnableBrotli bool
	TLSCert      string
	TLSKey       string
	APIKeys      []string
	// AllowedOrigins is the CORS allowlist. Only these origins are echoed
	// back in Access-Control-Allow-Origin. A single "*" entry allows any
	// origin but only without credentials. Empty means no cross-origin access.
	AllowedOrigins []string
	// AuthDisabled, when true, runs the API with NO authentication. It is an
	// explicit, logged mode. When false and APIKeys is empty, the server fails
	// CLOSED (every /api/v1 request is rejected) rather than serving open.
	AuthDisabled bool
}

// Config is the top-level application configuration.
type Config struct {
	Server ServerConfig
}

// Pool defines the interface for database connection pool operations.
// This is satisfied by the internal/db.Pool implementation.
type Pool interface {
	// Skill operations
	ListSkills(ctx context.Context, limit, offset int, status string) ([]models.Skill, int, error)
	GetSkill(ctx context.Context, id string) (*models.Skill, error)
	GetSkillByName(ctx context.Context, name string) (*models.Skill, error)
	CreateSkill(ctx context.Context, skill *models.Skill) error
	UpdateSkill(ctx context.Context, skill *models.Skill) error
	DeleteSkill(ctx context.Context, id string) error
	GetSkillTree(ctx context.Context, rootID string, maxDepth int) (*models.SkillTreeNode, error)
	ImportSkills(ctx context.Context, skills []models.Skill) (int, error)
	ExportSkills(ctx context.Context, id string) ([]models.Skill, error)

	// Search operations
	SearchSkills(ctx context.Context, query string, vector []float32, limit int) ([]models.SearchResult, error)
	SimilarSkills(ctx context.Context, skillID string, limit int) ([]models.SearchResult, error)

	// Registry operations
	GetRegistry(ctx context.Context, limit, offset int) ([]models.SkillRegistryEntry, int, error)
	GetMissingDeps(ctx context.Context, skillID string) ([]string, error)
	GetStaleSkills(ctx context.Context, limit, offset int) ([]models.SkillRegistryEntry, error)
	TriggerReview(ctx context.Context, skillID string) error
	GetCoverage(ctx context.Context) (map[string]interface{}, error)

	// Expansion operations
	TriggerExpand(ctx context.Context, skillName string, depth int) (*models.ExpansionJob, error)
	GetExpandStatus(ctx context.Context, jobID string) (*models.ExpansionJob, error)
	GetGapReport(ctx context.Context) (map[string]interface{}, error)

	// Learning operations
	SubmitProject(ctx context.Context, projectPath string, languages []string) (*models.LearningJob, error)
	GetLearnStatus(ctx context.Context, jobID string) (*models.LearningJob, error)
	GetEvidences(ctx context.Context, skillID string, limit, offset int) ([]models.Evidence, int, error)

	// Health check
	Ping(ctx context.Context) error

	// Close closes the pool
	Close() error
}

// Server is the HTTP API server for the HelixKnowledge Skill Graph System.
type Server struct {
	router      *gin.Engine
	pool        Pool
	cfg         ServerConfig
	logger      *zap.Logger
	httpServer  *http.Server
	http3Server *http3.Server

	// validator runs the fail-closed create-path validation. When nil (or
	// validationEnabled is false) newly-created skills are forced to draft — a
	// client can never self-promote to validated/active without a passing verdict.
	validator         SkillValidator
	validationEnabled bool
}

// Option customizes a Server at construction.
type Option func(*Server)

// WithValidator wires the create-path validation pipeline (§G03). enabled mirrors
// config.Validation.Enabled; validation runs only when both a validator is
// present and enabled is true.
func WithValidator(v SkillValidator, enabled bool) Option {
	return func(s *Server) {
		s.validator = v
		s.validationEnabled = enabled
	}
}

// New creates a new API server with the given database pool and configuration.
func New(pool Pool, cfg Config, logger *zap.Logger, opts ...Option) *Server {
	if logger == nil {
		logger = zap.L()
	}

	// Set Gin mode based on environment
	mode := gin.ReleaseMode
	if os.Getenv("GIN_MODE") != "" {
		mode = os.Getenv("GIN_MODE")
	}
	gin.SetMode(mode)

	router := gin.New()

	s := &Server{
		router: router,
		pool:   pool,
		cfg:    cfg.Server,
		logger: logger,
	}

	for _, opt := range opts {
		opt(s)
	}

	// Register middleware
	s.setupMiddleware()

	// Register routes
	s.SetupRoutes()

	return s
}

// setupMiddleware registers all global middleware on the router.
func (s *Server) setupMiddleware() {
	// Recovery must be first to catch panics from other middleware
	s.router.Use(Recovery())

	// Request ID for tracing
	s.router.Use(RequestID())

	// Structured logging
	s.router.Use(Logger(s.logger))

	// Prometheus metrics
	s.router.Use(MetricsMiddleware())

	// Content negotiation (JSON/TOML)
	s.router.Use(ContentNegotiation())

	// CORS (config-driven allowlist)
	s.router.Use(CORS(s.cfg.AllowedOrigins))

	// Brotli compression (if enabled)
	if s.cfg.EnableBrotli {
		s.router.Use(BrotliMiddleware())
	}

	// Body size limit (100MB for imports)
	s.router.Use(MaxBodySize(100 * 1024 * 1024))
}

// SetupRoutes registers all API routes.
func (s *Server) SetupRoutes() {
	// Health and metrics (no auth required)
	s.router.GET("/health", s.handleHealth)
	s.router.GET("/metrics", s.handleMetrics())
	s.router.GET("/version", s.handleVersion)

	// API v1 group
	v1 := s.router.Group("/api/v1")

	// Apply authentication to all v1 routes. ResolveAPIKeyAuth is fail-CLOSED:
	// with no API keys configured and auth not explicitly disabled it installs
	// a middleware that rejects every request (503), instead of the previous
	// fail-OPEN behaviour that served protected routes with no auth at all.
	if mw := ResolveAPIKeyAuth(s.cfg.APIKeys, s.cfg.AuthDisabled, s.logger); mw != nil {
		v1.Use(mw)
	}

	// Skills CRUD
	skills := v1.Group("/skills")
	{
		skills.GET("", s.handleListSkills)
		skills.POST("", s.handleCreateSkill)
		skills.GET("/:id", s.handleGetSkill)
		skills.PUT("/:id", s.handleUpdateSkill)
		skills.PATCH("/:id", s.handleUpdateSkill)
		skills.DELETE("/:id", s.handleDeleteSkill)
		skills.GET("/:id/tree", s.handleGetSkillTree)
		skills.POST("/import", DetectContentType(), s.handleImportSkills)
		skills.GET("/:id/export", s.handleExportSkill)
	}

	// Search
	v1.GET("/search", s.handleSearch)
	v1.POST("/search/similar", s.handleSimilarSkills)

	// Registry
	registry := v1.Group("/registry")
	{
		registry.GET("", s.handleGetRegistry)
		registry.GET("/missing-deps/:id", s.handleGetMissingDeps)
		registry.GET("/stale", s.handleGetStaleSkills)
		registry.POST("/review/:id", s.handleTriggerReview)
		registry.GET("/coverage", s.handleGetCoverage)
	}

	// Auto-expand
	expand := v1.Group("/expand")
	{
		expand.POST("", s.handleTriggerExpand)
		expand.GET("/status/:id", s.handleGetExpandStatus)
		expand.GET("/gaps", s.handleGetGapReport)
	}

	// Learning
	learn := v1.Group("/learn")
	{
		learn.POST("/projects", s.handleSubmitProject)
		learn.GET("/status/:id", s.handleGetLearnStatus)
		learn.GET("/evidences/:skill_id", s.handleGetEvidences)
	}
}

// Run starts the HTTP/2 server and optionally the HTTP/3 server.
// It blocks until the context is cancelled or a signal is received.
func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.HTTPPort)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start HTTP/3 if enabled and TLS is configured
	if s.cfg.EnableHTTP3 {
		if err := s.setupHTTP3(s.router); err != nil {
			s.logger.Warn("failed to start HTTP/3 server, continuing with HTTP/2 only",
				zap.Error(err),
			)
		}
	}

	// Channel to listen for errors from server goroutines
	serverErr := make(chan error, 1)

	// Start HTTP/2 server in a goroutine
	go func() {
		s.logger.Info("starting HTTP/2 server",
			zap.String("addr", addr),
			zap.Bool("http3_enabled", s.cfg.EnableHTTP3),
		)

		var err error
		if s.cfg.TLSCert != "" && s.cfg.TLSKey != "" {
			err = s.httpServer.ListenAndServeTLS(s.cfg.TLSCert, s.cfg.TLSKey)
		} else {
			s.logger.Warn("running without TLS - HTTP/2 will work but HTTP/3 requires TLS")
			err = s.httpServer.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return fmt.Errorf("HTTP server error: %w", err)
	case sig := <-quit:
		s.logger.Info("shutdown signal received",
			zap.String("signal", sig.String()),
		)
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.Shutdown(shutdownCtx)
}

// Shutdown gracefully shuts down the server, waiting for active connections
// to complete or the context to be cancelled.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("starting graceful shutdown")

	// Shutdown HTTP/3 first if running
	if err := s.shutdownHTTP3(ctx); err != nil {
		s.logger.Warn("HTTP/3 shutdown error", zap.Error(err))
	}

	// Shutdown HTTP/2 server
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("HTTP/2 server shutdown error: %w", err)
		}
	}

	s.logger.Info("server shutdown complete")
	return nil
}

// Router returns the underlying gin.Engine for testing purposes.
func (s *Server) Router() *gin.Engine {
	return s.router
}
