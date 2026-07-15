package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// SkillStatus represents the lifecycle state of a skill
type SkillStatus string

const (
	SkillStatusDraft      SkillStatus = "draft"
	SkillStatusValidated  SkillStatus = "validated"
	SkillStatusActive     SkillStatus = "active"
	SkillStatusDeprecated SkillStatus = "deprecated"
)

// DependencyType defines how skills relate to each other
type DependencyType string

const (
	DepTypeRequires   DependencyType = "requires"
	DepTypeExtends    DependencyType = "extends"
	DepTypeRecommends DependencyType = "recommends"
)

// Skill represents a single knowledge unit in the skill graph
type Skill struct {
	ID          uuid.UUID       `json:"id" db:"id" toml:"-"`
	Name        string          `json:"name" db:"name" toml:"name"`
	Version     string          `json:"version" db:"version" toml:"version"`
	Title       string          `json:"title" db:"title" toml:"title"`
	Description string          `json:"description" db:"description" toml:"description"`
	Content     string          `json:"content" db:"content" toml:"content"`
	Metadata    json.RawMessage `json:"metadata" db:"metadata" toml:"-"`
	Status      SkillStatus     `json:"status" db:"status" toml:"-"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at" toml:"-"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at" toml:"-"`

	// Runtime fields (not persisted directly)
	Dependencies  []SkillDependency `json:"dependencies,omitempty" db:"-" toml:"-"`
	Resources     []Resource        `json:"resources,omitempty" db:"-" toml:"-"`
	Embedding     pgvector.Vector   `json:"-" db:"embedding" toml:"-"`
	TreeDepth     int               `json:"tree_depth,omitempty" db:"-" toml:"-"`
}

// SkillDependency represents a directed edge in the skill DAG
type SkillDependency struct {
	SkillID      uuid.UUID      `json:"skill_id" db:"skill_id"`
	DependsOn    uuid.UUID      `json:"depends_on" db:"depends_on"`
	RelationType DependencyType `json:"relation_type" db:"relation_type"`

	// Join fields
	DependsOnName  string `json:"depends_on_name,omitempty" db:"depends_on_name"`
	DependsOnTitle string `json:"depends_on_title,omitempty" db:"depends_on_title"`
}

// Resource is an external reference attached to a skill
type Resource struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	SkillID       uuid.UUID  `json:"skill_id" db:"skill_id"`
	URL           string     `json:"url" db:"url" toml:"url"`
	Title         string     `json:"title" db:"title" toml:"title"`
	ResourceType  string     `json:"resource_type" db:"resource_type" toml:"resource_type"`
	FetchedHash   string     `json:"fetched_hash" db:"fetched_hash"`
	ContentCached string     `json:"content_cached" db:"content_cached"`
	LastValidated *time.Time `json:"last_validated" db:"last_validated"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// Evidence is a learned pattern from a real codebase
type Evidence struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	SkillID       uuid.UUID       `json:"skill_id" db:"skill_id"`
	SourceProject string          `json:"source_project" db:"source_project"`
	SourceFile    string          `json:"source_file" db:"source_file"`
	CodeSnippet   string          `json:"code_snippet" db:"code_snippet"`
	Pattern       string          `json:"pattern" db:"pattern"`
	Language      string          `json:"language" db:"language"`
	Validated     bool            `json:"validated" db:"validated"`
	Embedding     pgvector.Vector `json:"-" db:"embedding"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}

// SkillRegistryEntry tracks health and completeness
type SkillRegistryEntry struct {
	SkillID     uuid.UUID  `json:"skill_id" db:"skill_id"`
	SkillName   string     `json:"skill_name" db:"skill_name"`
	MissingDeps []string   `json:"missing_deps" db:"missing_deps"`
	Stale       bool       `json:"stale" db:"stale"`
	LastReview  *time.Time `json:"last_review" db:"last_review"`
	AutoExpand  bool       `json:"auto_expand" db:"auto_expand"`
	Coverage    float64    `json:"coverage" db:"coverage"`
}

// AuditLogEntry tracks all system events
type AuditLogEntry struct {
	Timestamp time.Time       `json:"ts" db:"ts"`
	Event     string          `json:"event" db:"event"`
	SkillID   *uuid.UUID      `json:"skill_id,omitempty" db:"skill_id"`
	Details   json.RawMessage `json:"details" db:"details"`
}

// SkillTreeNode wraps a skill with tree traversal info
type SkillTreeNode struct {
	Skill    Skill           `json:"skill"`
	Depth    int             `json:"depth"`
	Children []SkillTreeNode `json:"children,omitempty"`
}

// SkillMetadata is the structured metadata for a skill
type SkillMetadata struct {
	Tags       []string `json:"tags" toml:"tags"`
	Domain     string   `json:"domain" toml:"domain"`
	Complexity string   `json:"complexity" toml:"complexity"`
}

// TOMLSkillWrapper is used for TOML import/export
type TOMLSkillWrapper struct {
	Skill        TOMLSkillDef       `toml:"skill"`
	Dependencies TOMLDependencies   `toml:"skill.dependencies"`
	Resources    []TOMLResource     `toml:"skill.resources"`
}

type TOMLSkillDef struct {
	Name        string          `toml:"name"`
	Version     string          `toml:"version"`
	Title       string          `toml:"title"`
	Description string          `toml:"description"`
	Content     string          `toml:"content"`
	Metadata    SkillMetadata   `toml:"metadata"`
}

type TOMLDependencies struct {
	Requires   []string `toml:"requires"`
	Extends    []string `toml:"extends"`
	Recommends []string `toml:"recommends"`
}

type TOMLResource struct {
	URL          string `toml:"url"`
	Title        string `toml:"title"`
	ResourceType string `toml:"resource_type"`
}

// SearchResult wraps a skill with search relevance score
type SearchResult struct {
	Skill Skill   `json:"skill"`
	Score float64 `json:"score"`
}

// ExpansionJob tracks auto-expansion progress
type ExpansionJob struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SkillName string    `json:"skill_name" db:"skill_name"`
	Status    string    `json:"status" db:"status"` // pending | running | completed | failed
	Depth     int       `json:"depth" db:"depth"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// LearningJob tracks codebase analysis progress
type LearningJob struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectPath string  `json:"project_path" db:"project_path"`
	Status    string    `json:"status" db:"status"`
	Languages []string  `json:"languages" db:"languages"`
	FilesProcessed int  `json:"files_processed" db:"files_processed"`
	PatternsFound  int  `json:"patterns_found" db:"patterns_found"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
