// Package config provides configuration loading from TOML files with
// environment variable overrides. It supports ${VAR} interpolation syntax
// in string values, allowing secrets and dynamic values to be injected
// via environment variables.
//
// Usage:
//
//	cfg, err := config.Load("config/config.toml")
//	if err != nil {
//	    log.Fatal(err)
//	}
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// envVarRegex matches ${VAR} or ${VAR:-default} syntax in config values.
var envVarRegex = regexp.MustCompile(`\$\{([^}]+)\}`)

// ---------------------------------------------------------------------------
// Config (root)
// ---------------------------------------------------------------------------

// Config holds all application configuration sections.
type Config struct {
	Server       ServerConfig       `toml:"server"`
	Database     DatabaseConfig     `toml:"database"`
	Embedding    EmbeddingConfig    `toml:"embedding"`
	Validation   ValidationConfig   `toml:"validation"`
	AutoExpand   AutoExpandConfig   `toml:"autoexpand"`
	CodeAnalysis CodeAnalysisConfig `toml:"codeanalysis"`
	MCP          MCPConfig          `toml:"mcp"`
	Registry     RegistryConfig     `toml:"registry"`
	Logging      LoggingConfig      `toml:"logging"`
}

// ---------------------------------------------------------------------------
// Section structs
// ---------------------------------------------------------------------------

// ServerConfig controls the HTTP/HTTPS server behaviour.
type ServerConfig struct {
	Host         string `toml:"host"`
	HTTPPort     int    `toml:"http_port"`
	HTTP3Port    int    `toml:"http3_port"`
	EnableHTTP3  bool   `toml:"enable_http3"`
	EnableBrotli bool   `toml:"enable_brotli"`
	TLSCert      string `toml:"tls_cert"`
	TLSKey       string `toml:"tls_key"`
	// AllowedOrigins is the CORS allowlist of exact origins permitted to make
	// cross-origin requests (e.g. "https://app.example.com"). A single "*"
	// entry allows any origin but only without credentials. Empty (the default)
	// disallows all cross-origin access.
	AllowedOrigins []string `toml:"allowed_origins"`
	// APIKeys is the set of valid keys that authenticate /api/v1 requests via
	// the X-API-Key header. Prefer providing these through the HELIX_API_KEYS
	// environment override (comma-separated) so secrets never live in tracked
	// config (§11.4.10). When APIKeys is empty AND AuthDisabled is false, the
	// server fails CLOSED and refuses every /api/v1 request.
	APIKeys []string `toml:"api_keys"`
	// AuthDisabled explicitly runs the API with NO authentication. It must be
	// set deliberately and is logged loudly at startup. Absent keys without
	// this flag is a fail-closed error, never a silent open server.
	AuthDisabled bool `toml:"auth_disabled"`
}

// DatabaseConfig controls the PostgreSQL connection pool.
type DatabaseConfig struct {
	Host           string        `toml:"host"`
	Port           int           `toml:"port"`
	Database       string        `toml:"database"`
	User           string        `toml:"user"`
	Password       string        `toml:"password"`
	SSLMode        string        `toml:"ssl_mode"`
	MaxConnections int           `toml:"max_connections"`
	ConnectTimeout time.Duration `toml:"connect_timeout"`
}

// DSN returns a PostgreSQL keyword/value connection string for pgx or lib/pq.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		d.Host, d.Port, d.Database, d.User, d.Password, d.SSLMode,
	)
}

// DSNWithTimeout returns a DSN with connect_timeout included.
func (d DatabaseConfig) DSNWithTimeout() string {
	timeoutSec := int(d.ConnectTimeout.Seconds())
	if timeoutSec <= 0 {
		timeoutSec = 10
	}
	return d.DSN() + fmt.Sprintf(" connect_timeout=%d", timeoutSec)
}

// EmbeddingConfig selects the embedding provider and model.
type EmbeddingConfig struct {
	Provider      string `toml:"provider"`       // "openai" | "local"
	Dimensions    int    `toml:"dimensions"`     // e.g. 768
	Model         string `toml:"model"`          // e.g. "text-embedding-3-small"
	APIKey        string `toml:"api_key"`        // OpenAI API key (env override recommended)
	LocalEndpoint string `toml:"local_endpoint"` // URL for local model server
}

// ValidationConfig controls the skill validation pipeline.
type ValidationConfig struct {
	Enabled             bool   `toml:"enabled"`
	SandboxType         string `toml:"sandbox_type"`       // "wasm" | "docker" | "none"
	JurySize            int    `toml:"jury_size"`          // number of validators
	ApprovalThreshold   int    `toml:"approval_threshold"` // votes required
	AutoApproveEvidence bool   `toml:"auto_approve_evidence"`
	RequireHumanReview  bool   `toml:"require_human_review"`
}

// AutoExpandConfig controls the automatic skill-tree expansion.
type AutoExpandConfig struct {
	Enabled            bool   `toml:"enabled"`
	MaxDepth           int    `toml:"max_depth"`
	MaxNewSkillsPerRun int    `toml:"max_new_skills_per_run"`
	LLMProvider        string `toml:"llm_provider"`
	LLMModel           string `toml:"llm_model"`
}

// CodeAnalysisConfig controls the repository-learning subsystem.
type CodeAnalysisConfig struct {
	Enabled         bool     `toml:"enabled"`
	Languages       []string `toml:"languages"`
	MaxFileSizeKB   int      `toml:"max_file_size_kb"`
	ExcludePatterns []string `toml:"exclude_patterns"`
}

// MCPConfig controls the Model Context Protocol integration.
type MCPConfig struct {
	Enabled   bool   `toml:"enabled"`
	Transport string `toml:"transport"` // "stdio" | "http"
}

// RegistryConfig controls skill-registry behaviour.
type RegistryConfig struct {
	ReviewIntervalHours int     `toml:"review_interval_hours"`
	CoverageThreshold   float64 `toml:"coverage_threshold"`
}

// LoggingConfig controls Zap logger output.
type LoggingConfig struct {
	Level  string `toml:"level"`  // "debug" | "info" | "warn" | "error"
	Format string `toml:"format"` // "json" | "console"
}

// ---------------------------------------------------------------------------
// Defaults
// ---------------------------------------------------------------------------

func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			HTTPPort:     8080,
			HTTP3Port:    8443,
			EnableHTTP3:  false,
			EnableBrotli: true,
		},
		Database: DatabaseConfig{
			Host:           "localhost",
			Port:           5432,
			Database:       "skilldb",
			User:           "skill",
			Password:       "secret",
			SSLMode:        "disable",
			MaxConnections: 25,
			ConnectTimeout: 10 * time.Second,
		},
		Embedding: EmbeddingConfig{
			Provider:   "openai",
			Dimensions: 768,
			Model:      "text-embedding-3-small",
		},
		Validation: ValidationConfig{
			Enabled:             true,
			SandboxType:         "wasm",
			JurySize:            3,
			ApprovalThreshold:   2,
			AutoApproveEvidence: false,
			RequireHumanReview:  true,
		},
		AutoExpand: AutoExpandConfig{
			Enabled:            true,
			MaxDepth:           5,
			MaxNewSkillsPerRun: 10,
			LLMProvider:        "openai",
			LLMModel:           "gpt-4o-mini",
		},
		CodeAnalysis: CodeAnalysisConfig{
			Enabled:         true,
			Languages:       []string{"java", "kotlin", "c", "cpp", "python", "go"},
			MaxFileSizeKB:   500,
			ExcludePatterns: []string{"vendor/", "node_modules/", ".git/", "build/"},
		},
		MCP: MCPConfig{
			Enabled:   true,
			Transport: "stdio",
		},
		Registry: RegistryConfig{
			ReviewIntervalHours: 24,
			CoverageThreshold:   0.8,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// Load reads a TOML configuration file, applies environment-variable
// substitution on all string fields, and returns the populated Config.
//
// Environment variables use the ${VAR} syntax. A default value can be
// provided with ${VAR:-default}. If the variable is unset and no default
// is given, the empty string is substituted.
//
// If path is empty, Load searches for config.toml in the current directory,
// then config/config.toml, then /etc/helixskill/config.toml.
func Load(path string) (*Config, error) {
	if path == "" {
		candidates := []string{
			"config.toml",
			"config/config.toml",
			"/etc/helixskill/config.toml",
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				path = c
				break
			}
		}
	}
	if path == "" {
		return nil, fmt.Errorf("no configuration file found")
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	cfg := defaultConfig()

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("decode TOML config from %q: %w", path, err)
	}

	// Apply environment-variable substitution.
	if err := substituteEnv(&cfg); err != nil {
		return nil, fmt.Errorf("environment variable substitution: %w", err)
	}

	// Apply explicit environment overrides.
	applyEnvOverrides(&cfg)

	// Validate critical fields.
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	return &cfg, nil
}

// ---------------------------------------------------------------------------
// Environment variable substitution (${VAR} / ${VAR:-default})
// ---------------------------------------------------------------------------

// substituteEnv walks the Config struct and replaces ${VAR} placeholders
// in every string field with the corresponding environment variable value.
func substituteEnv(cfg *Config) error {
	var errs []string

	sub := func(v string) string {
		replaced, err := interpolate(v)
		if err != nil {
			errs = append(errs, err.Error())
			return v
		}
		return replaced
	}

	// Server
	cfg.Server.Host = sub(cfg.Server.Host)
	cfg.Server.TLSCert = sub(cfg.Server.TLSCert)
	cfg.Server.TLSKey = sub(cfg.Server.TLSKey)
	// Server list fields honor the same ${VAR} promise as scalar fields so a
	// TOML entry like api_keys = ["${PROD_KEY}"] is interpolated from the
	// environment rather than stored as a literal (and dangerously valid) key.
	for i := range cfg.Server.APIKeys {
		cfg.Server.APIKeys[i] = sub(cfg.Server.APIKeys[i])
	}
	for i := range cfg.Server.AllowedOrigins {
		cfg.Server.AllowedOrigins[i] = sub(cfg.Server.AllowedOrigins[i])
	}

	// Database
	cfg.Database.Host = sub(cfg.Database.Host)
	cfg.Database.Database = sub(cfg.Database.Database)
	cfg.Database.User = sub(cfg.Database.User)
	cfg.Database.Password = sub(cfg.Database.Password)
	cfg.Database.SSLMode = sub(cfg.Database.SSLMode)

	// Embedding
	cfg.Embedding.Provider = sub(cfg.Embedding.Provider)
	cfg.Embedding.Model = sub(cfg.Embedding.Model)
	cfg.Embedding.APIKey = sub(cfg.Embedding.APIKey)
	cfg.Embedding.LocalEndpoint = sub(cfg.Embedding.LocalEndpoint)

	// Validation
	cfg.Validation.SandboxType = sub(cfg.Validation.SandboxType)

	// AutoExpand
	cfg.AutoExpand.LLMProvider = sub(cfg.AutoExpand.LLMProvider)
	cfg.AutoExpand.LLMModel = sub(cfg.AutoExpand.LLMModel)

	// MCP
	cfg.MCP.Transport = sub(cfg.MCP.Transport)

	// Logging
	cfg.Logging.Level = sub(cfg.Logging.Level)
	cfg.Logging.Format = sub(cfg.Logging.Format)

	if len(errs) > 0 {
		return fmt.Errorf("errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

// interpolate replaces all ${VAR} occurrences in s with their environment
// variable values. Supports ${VAR:-default} syntax.
func interpolate(s string) (string, error) {
	result := envVarRegex.ReplaceAllStringFunc(s, func(match string) string {
		inner := match[2 : len(match)-1] // strip ${ and }

		// Check for default syntax: VAR:-default
		var envKey, defaultVal string
		if idx := strings.Index(inner, ":-"); idx >= 0 {
			envKey = inner[:idx]
			defaultVal = inner[idx+2:]
		} else {
			envKey = inner
		}

		if v := os.Getenv(envKey); v != "" {
			return v
		}
		if defaultVal != "" {
			return defaultVal
		}
		return ""
	})
	return result, nil
}

// ---------------------------------------------------------------------------
// Explicit environment overrides (HELIX_* prefix)
// ---------------------------------------------------------------------------

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("HELIX_DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("HELIX_DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.Port = port
		}
	}
	if v := os.Getenv("HELIX_DB_NAME"); v != "" {
		cfg.Database.Database = v
	}
	if v := os.Getenv("HELIX_DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("HELIX_DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("HELIX_DB_SSLMODE"); v != "" {
		cfg.Database.SSLMode = v
	}
	if v := os.Getenv("HELIX_LOG_LEVEL"); v != "" {
		cfg.Logging.Level = v
	}
	if v := os.Getenv("HELIX_MCP_TRANSPORT"); v != "" {
		cfg.MCP.Transport = v
	}
	if v := os.Getenv("HELIX_API_KEYS"); v != "" {
		cfg.Server.APIKeys = splitAndTrim(v)
	}
	if v := os.Getenv("HELIX_AUTH_DISABLED"); v != "" {
		cfg.Server.AuthDisabled = v == "1" || strings.EqualFold(v, "true")
	}
}

// splitAndTrim splits a comma-separated string into non-empty, trimmed values.
// Used for list-valued environment overrides such as HELIX_API_KEYS.
func splitAndTrim(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Validation
// ---------------------------------------------------------------------------

func validate(cfg *Config) error {
	var issues []string

	if cfg.Server.HTTPPort <= 0 || cfg.Server.HTTPPort > 65535 {
		issues = append(issues, fmt.Sprintf("invalid server.http_port: %d", cfg.Server.HTTPPort))
	}
	if cfg.Server.HTTP3Port <= 0 || cfg.Server.HTTP3Port > 65535 {
		issues = append(issues, fmt.Sprintf("invalid server.http3_port: %d", cfg.Server.HTTP3Port))
	}

	// Defense-in-depth (§11.4.10): an api_keys/allowed_origins entry that still
	// contains a "${" placeholder AFTER interpolation means a referenced
	// environment variable was never set (or the placeholder was malformed).
	// Fail CLOSED rather than treat the literal placeholder as a valid secret
	// or origin. Values are never echoed — only the field name and index.
	for i, k := range cfg.Server.APIKeys {
		if strings.Contains(k, "${") {
			issues = append(issues, fmt.Sprintf("server.api_keys[%d] contains an uninterpolated ${...} placeholder (set the referenced environment variable or remove the entry)", i))
		}
	}
	for i, o := range cfg.Server.AllowedOrigins {
		if strings.Contains(o, "${") {
			issues = append(issues, fmt.Sprintf("server.allowed_origins[%d] contains an uninterpolated ${...} placeholder", i))
		}
	}

	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		issues = append(issues, fmt.Sprintf("invalid database.port: %d", cfg.Database.Port))
	}
	if cfg.Database.MaxConnections <= 0 {
		cfg.Database.MaxConnections = 25 // apply safe default
	}

	if cfg.Embedding.Dimensions <= 0 {
		issues = append(issues, fmt.Sprintf("invalid embedding.dimensions: %d", cfg.Embedding.Dimensions))
	}

	if cfg.Validation.JurySize <= 0 {
		issues = append(issues, fmt.Sprintf("invalid validation.jury_size: %d", cfg.Validation.JurySize))
	}
	if cfg.Validation.ApprovalThreshold <= 0 {
		issues = append(issues, fmt.Sprintf("invalid validation.approval_threshold: %d", cfg.Validation.ApprovalThreshold))
	}

	if cfg.AutoExpand.MaxDepth <= 0 {
		issues = append(issues, fmt.Sprintf("invalid autoexpand.max_depth: %d", cfg.AutoExpand.MaxDepth))
	}
	if cfg.AutoExpand.MaxNewSkillsPerRun <= 0 {
		issues = append(issues, fmt.Sprintf("invalid autoexpand.max_new_skills_per_run: %d", cfg.AutoExpand.MaxNewSkillsPerRun))
	}

	if cfg.Registry.CoverageThreshold < 0 || cfg.Registry.CoverageThreshold > 1 {
		issues = append(issues, fmt.Sprintf("invalid registry.coverage_threshold: %f (must be 0-1)", cfg.Registry.CoverageThreshold))
	}

	if len(issues) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(issues, "; "))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helper methods on Config
// ---------------------------------------------------------------------------

// ListenAddr returns the HTTP listen address in the form ":port".
func (c *Config) ListenAddr() string {
	return ":" + strconv.Itoa(c.Server.HTTPPort)
}

// HTTP3ListenAddr returns the HTTP/3 listen address.
func (c *Config) HTTP3ListenAddr() string {
	return ":" + strconv.Itoa(c.Server.HTTP3Port)
}
