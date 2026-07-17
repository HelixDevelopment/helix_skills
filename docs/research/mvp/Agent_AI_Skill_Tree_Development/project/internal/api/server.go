package api

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
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
	Close()
}

// Server is the HTTP API server for the HelixKnowledge Skill Graph System.
type Server struct {
	router *gin.Engine
	pool   Pool
	cfg    ServerConfig
	logger *zap.Logger

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

	// System routes (open, no auth)
	s.router.GET("/health", s.handleHealth)
	s.router.GET("/metrics", s.handleMetrics())
	s.router.GET("/version", s.handleVersion)

	// API v1 with auth
	v1 := s.router.Group("/api/v1")
	if mw := ResolveAPIKeyAuth(s.cfg.APIKeys, s.cfg.AuthDisabled, s.logger); mw != nil {
		v1.Use(mw)
	}
	s.RegisterHandlers(v1)

	return s
}

// NewHandler creates a Server suitable for registering handlers onto an
// existing Gin router, without creating its own gin.Engine or registering
// middleware. Use RegisterHandlers to register routes on a provided group.
func NewHandler(pool Pool, logger *zap.Logger, opts ...Option) *Server {
	if logger == nil {
		logger = zap.L()
	}
	s := &Server{
		pool:   pool,
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
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
//
// Deprecated: Use RegisterHandlers instead. This method is retained for
// backward compatibility with legacy callers. It registers system routes
// directly on s.router and data routes on an auth-guarded /api/v1 group.
func (s *Server) SetupRoutes() {
	// System routes (open, no auth)
	s.router.GET("/health", s.handleHealth)
	s.router.GET("/metrics", s.handleMetrics())
	s.router.GET("/version", s.handleVersion)

	// API v1 with auth
	v1 := s.router.Group("/api/v1")
	if mw := ResolveAPIKeyAuth(s.cfg.APIKeys, s.cfg.AuthDisabled, s.logger); mw != nil {
		v1.Use(mw)
	}
	s.RegisterHandlers(v1)
}

// RegisterHandlers registers all API handler routes onto the provided Gin
// router group. The group should already have authentication middleware applied
// where needed. This is the canonical route registration entry point.
func (s *Server) RegisterHandlers(router gin.IRouter) {
	// Skills CRUD
	skills := router.Group("/skills")
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
	router.GET("/search", s.handleSearch)
	router.POST("/search/similar", s.handleSimilarSkills)

	// Registry
	reg := router.Group("/registry")
	{
		reg.GET("", s.handleGetRegistry)
		reg.GET("/missing-deps/:id", s.handleGetMissingDeps)
		reg.GET("/stale", s.handleGetStaleSkills)
		reg.POST("/review/:id", s.handleTriggerReview)
		reg.GET("/coverage", s.handleGetCoverage)
	}

	// Auto-expand
	expand := router.Group("/expand")
	{
		expand.POST("", s.handleTriggerExpand)
		expand.GET("/status/:id", s.handleGetExpandStatus)
		expand.GET("/gaps", s.handleGetGapReport)
	}

	// Learning
	learn := router.Group("/learn")
	{
		learn.POST("/projects", s.handleSubmitProject)
		learn.GET("/status/:id", s.handleGetLearnStatus)
		learn.GET("/evidences/:skill_id", s.handleGetEvidences)
	}
}

// Router returns the underlying gin.Engine for testing purposes.
func (s *Server) Router() *gin.Engine {
	return s.router
}
