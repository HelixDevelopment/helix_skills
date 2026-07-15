// Package skill provides CRUD and search operations for skills in the knowledge graph.
package skill

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/helixdevelopment/skill-system/internal/db"
	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

// Store provides data access for skills and related entities.
type Store struct {
	pool *db.Pool
}

// NewStore creates a new skill store.
func NewStore(pool *db.Pool) *Store {
	return &Store{pool: pool}
}

// Pool returns the underlying database pool for operations that need
// direct database access (e.g., audit logging from other packages).
func (s *Store) Pool() *db.Pool {
	return s.pool
}

// Search performs a hybrid search combining vector similarity and text matching.
func (s *Store) Search(ctx context.Context, query string, limit int) ([]models.SearchResult, error) {
	// For now, use text-based search. In production, this would generate
	// embeddings and use vector similarity + full-text search.
	sql := `
		SELECT s.id, s.name, s.version, s.title, s.description, s.content,
		       s.metadata, s.status, s.created_at, s.updated_at,
		       similarity(s.name || ' ' || s.title || ' ' || s.description, $1) as score
		FROM skills s
		WHERE s.name % $1 OR s.title % $1 OR s.description ILIKE '%' || $1 || '%'
		ORDER BY score DESC, s.name
		LIMIT $2
	`
	rows, err := s.pool.Query(ctx, sql, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search skills: %w", err)
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var r models.SearchResult
		var metadata []byte
		err := rows.Scan(
			&r.Skill.ID, &r.Skill.Name, &r.Skill.Version, &r.Skill.Title,
			&r.Skill.Description, &r.Skill.Content, &metadata,
			&r.Skill.Status, &r.Skill.CreatedAt, &r.Skill.UpdatedAt,
			&r.Score,
		)
		if err != nil {
			return nil, fmt.Errorf("scan search result: %w", err)
		}
		r.Skill.Metadata = metadata
		results = append(results, r)
	}

	if len(results) == 0 {
		// Fallback: return all skills if no similarity match
		sql = `
			SELECT s.id, s.name, s.version, s.title, s.description, s.content,
			       s.metadata, s.status, s.created_at, s.updated_at, 0.0 as score
			FROM skills s
			WHERE s.name ILIKE '%' || $1 || '%' OR s.title ILIKE '%' || $1 || '%'
			ORDER BY s.name
			LIMIT $2
		`
		rows, err := s.pool.Query(ctx, sql, query, limit)
		if err != nil {
			return nil, fmt.Errorf("fallback search: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var r models.SearchResult
			var metadata []byte
			if err := rows.Scan(
				&r.Skill.ID, &r.Skill.Name, &r.Skill.Version, &r.Skill.Title,
				&r.Skill.Description, &r.Skill.Content, &metadata,
				&r.Skill.Status, &r.Skill.CreatedAt, &r.Skill.UpdatedAt,
				&r.Score,
			); err != nil {
				return nil, fmt.Errorf("scan fallback result: %w", err)
			}
			r.Skill.Metadata = metadata
			results = append(results, r)
		}
	}

	return results, nil
}

// GetByName retrieves a complete skill by its unique name.
func (s *Store) GetByName(ctx context.Context, name string) (*models.Skill, error) {
	sql := `
		SELECT s.id, s.name, s.version, s.title, s.description, s.content,
		       s.metadata, s.status, s.created_at, s.updated_at
		FROM skills s
		WHERE s.name = $1
	`
	var skill models.Skill
	var metadata []byte
	err := s.pool.QueryRow(ctx, sql, name).Scan(
		&skill.ID, &skill.Name, &skill.Version, &skill.Title,
		&skill.Description, &skill.Content, &metadata,
		&skill.Status, &skill.CreatedAt, &skill.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("skill not found: %s", name)
		}
		return nil, fmt.Errorf("get skill: %w", err)
	}
	skill.Metadata = metadata

	// Load dependencies
	depsSQL := `
		SELECT sd.skill_id, sd.depends_on, sd.relation_type,
		       ds.name as depends_on_name, ds.title as depends_on_title
		FROM skill_dependencies sd
		JOIN skills ds ON sd.depends_on = ds.id
		WHERE sd.skill_id = $1
	`
	depRows, err := s.pool.Query(ctx, depsSQL, skill.ID)
	if err != nil {
		return nil, fmt.Errorf("get dependencies: %w", err)
	}
	defer depRows.Close()
	for depRows.Next() {
		var d models.SkillDependency
		if err := depRows.Scan(&d.SkillID, &d.DependsOn, &d.RelationType, &d.DependsOnName, &d.DependsOnTitle); err != nil {
			return nil, fmt.Errorf("scan dependency: %w", err)
		}
		skill.Dependencies = append(skill.Dependencies, d)
	}

	// Load resources
	resSQL := `
		SELECT id, skill_id, url, title, resource_type, fetched_hash, content_cached, last_validated, created_at
		FROM resources WHERE skill_id = $1
	`
	resRows, err := s.pool.Query(ctx, resSQL, skill.ID)
	if err != nil {
		return nil, fmt.Errorf("get resources: %w", err)
	}
	defer resRows.Close()
	for resRows.Next() {
		var r models.Resource
		if err := resRows.Scan(&r.ID, &r.SkillID, &r.URL, &r.Title, &r.ResourceType, &r.FetchedHash, &r.ContentCached, &r.LastValidated, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan resource: %w", err)
		}
		skill.Resources = append(skill.Resources, r)
	}

	return &skill, nil
}

// GetTree returns the dependency tree for a skill up to the specified depth.
func (s *Store) GetTree(ctx context.Context, name string, maxDepth int) (*models.SkillTreeNode, error) {
	skill, err := s.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	root := &models.SkillTreeNode{
		Skill: *skill,
		Depth: 0,
	}

	visited := make(map[uuid.UUID]bool)
	visited[skill.ID] = true

	if err := s.buildTree(ctx, root, 1, maxDepth, visited); err != nil {
		return nil, fmt.Errorf("build tree: %w", err)
	}

	return root, nil
}

func (s *Store) buildTree(ctx context.Context, node *models.SkillTreeNode, depth, maxDepth int, visited map[uuid.UUID]bool) error {
	if depth > maxDepth {
		return nil
	}

	for _, dep := range node.Skill.Dependencies {
		if visited[dep.DependsOn] {
			continue
		}
		visited[dep.DependsOn] = true

		childSQL := `
			SELECT s.id, s.name, s.version, s.title, s.description, s.content,
			       s.metadata, s.status, s.created_at, s.updated_at
			FROM skills s WHERE s.id = $1
		`
		var child models.Skill
		var metadata []byte
		err := s.pool.QueryRow(ctx, childSQL, dep.DependsOn).Scan(
			&child.ID, &child.Name, &child.Version, &child.Title,
			&child.Description, &child.Content, &metadata,
			&child.Status, &child.CreatedAt, &child.UpdatedAt,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				continue // skip missing dependencies
			}
			return err
		}
		child.Metadata = metadata

		// Load child's dependencies
		depsSQL := `
			SELECT sd.skill_id, sd.depends_on, sd.relation_type,
			       ds.name, ds.title
			FROM skill_dependencies sd
			JOIN skills ds ON sd.depends_on = ds.id
			WHERE sd.skill_id = $1
		`
		depRows, err := s.pool.Query(ctx, depsSQL, child.ID)
		if err != nil {
			return err
		}
		for depRows.Next() {
			var d models.SkillDependency
			if err := depRows.Scan(&d.SkillID, &d.DependsOn, &d.RelationType, &d.DependsOnName, &d.DependsOnTitle); err != nil {
				depRows.Close()
				return err
			}
			child.Dependencies = append(child.Dependencies, d)
		}
		depRows.Close()

		childNode := models.SkillTreeNode{
			Skill: child,
			Depth: depth,
		}
		node.Children = append(node.Children, childNode)

		if err := s.buildTree(ctx, &node.Children[len(node.Children)-1], depth+1, maxDepth, visited); err != nil {
			return err
		}
	}

	return nil
}

// Create inserts a new skill into the database.
func (s *Store) Create(ctx context.Context, skill *models.Skill) error {
	if skill.ID == uuid.Nil {
		skill.ID = uuid.New()
	}

	metadataJSON, err := json.Marshal(skill.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	sql := `
		INSERT INTO skills (id, name, version, title, description, content, metadata, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		ON CONFLICT (name) DO UPDATE SET
			version = EXCLUDED.version,
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			content = EXCLUDED.content,
			metadata = EXCLUDED.metadata,
			status = EXCLUDED.status,
			updated_at = NOW()
		RETURNING id
	`
	var returnedID uuid.UUID
	err = s.pool.QueryRow(ctx, sql,
		skill.ID, skill.Name, skill.Version, skill.Title,
		skill.Description, skill.Content, metadataJSON,
		skill.Status,
	).Scan(&returnedID)
	if err != nil {
		return fmt.Errorf("create skill: %w", err)
	}

	// Insert dependencies
	for _, dep := range skill.Dependencies {
		depSQL := `
			INSERT INTO skill_dependencies (skill_id, depends_on, relation_type)
			VALUES ($1, $2, $3)
			ON CONFLICT (skill_id, depends_on) DO UPDATE SET relation_type = EXCLUDED.relation_type
		`
		_, err := s.pool.Exec(ctx, depSQL, returnedID, dep.DependsOn, dep.RelationType)
		if err != nil {
			return fmt.Errorf("create dependency: %w", err)
		}
	}

	// Upsert registry entry
	regSQL := `
		INSERT INTO skill_registry (skill_id, skill_name, missing_deps, stale, auto_expand, coverage)
		VALUES ($1, $2, '{}', false, true, 0.0)
		ON CONFLICT (skill_id) DO NOTHING
	`
	_, err = s.pool.Exec(ctx, regSQL, returnedID, skill.Name)
	if err != nil {
		return fmt.Errorf("create registry entry: %w", err)
	}

	skill.ID = returnedID
	return nil
}

// CreateFromTOML creates a skill from a TOML skill wrapper.
func (s *Store) CreateFromTOML(ctx context.Context, wrapper *models.TOMLSkillWrapper) (*models.Skill, error) {
	metadataJSON, _ := json.Marshal(wrapper.Skill.Metadata)

	skill := &models.Skill{
		ID:          uuid.New(),
		Name:        wrapper.Skill.Name,
		Version:     wrapper.Skill.Version,
		Title:       wrapper.Skill.Title,
		Description: wrapper.Skill.Description,
		Content:     wrapper.Skill.Content,
		Metadata:    metadataJSON,
		Status:      models.SkillStatusDraft,
	}

	// Resolve dependencies
	for _, depName := range wrapper.Dependencies.Requires {
		depID, err := s.resolveSkillID(ctx, depName)
		if err != nil {
			return nil, fmt.Errorf("resolve dependency %q: %w", depName, err)
		}
		skill.Dependencies = append(skill.Dependencies, models.SkillDependency{
			SkillID:      skill.ID,
			DependsOn:    depID,
			RelationType: models.DepTypeRequires,
		})
	}
	for _, depName := range wrapper.Dependencies.Extends {
		depID, err := s.resolveSkillID(ctx, depName)
		if err != nil {
			return nil, fmt.Errorf("resolve dependency %q: %w", depName, err)
		}
		skill.Dependencies = append(skill.Dependencies, models.SkillDependency{
			SkillID:      skill.ID,
			DependsOn:    depID,
			RelationType: models.DepTypeExtends,
		})
	}
	for _, depName := range wrapper.Dependencies.Recommends {
		depID, err := s.resolveSkillID(ctx, depName)
		if err != nil {
			return nil, fmt.Errorf("resolve dependency %q: %w", depName, err)
		}
		skill.Dependencies = append(skill.Dependencies, models.SkillDependency{
			SkillID:      skill.ID,
			DependsOn:    depID,
			RelationType: models.DepTypeRecommends,
		})
	}

	// Add resources
	for _, r := range wrapper.Resources {
		skill.Resources = append(skill.Resources, models.Resource{
			ID:           uuid.New(),
			SkillID:      skill.ID,
			URL:          r.URL,
			Title:        r.Title,
			ResourceType: r.ResourceType,
		})
	}

	if err := s.Create(ctx, skill); err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *Store) resolveSkillID(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, "SELECT id FROM skills WHERE name = $1", name).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, fmt.Errorf("skill %q not found", name)
		}
		return uuid.Nil, err
	}
	return id, nil
}

// GetMissingSkills returns skills with missing dependencies (gaps in the graph).
func (s *Store) GetMissingSkills(ctx context.Context, domain string) ([]models.SkillRegistryEntry, error) {
	sql := `
		SELECT sr.skill_id, sr.skill_name, sr.missing_deps, sr.stale, sr.last_review, sr.auto_expand, sr.coverage
		FROM skill_registry sr
		WHERE array_length(sr.missing_deps, 1) > 0
	`
	args := []interface{}{}
	argIdx := 1

	if domain != "" {
		sql += fmt.Sprintf(` AND EXISTS (
			SELECT 1 FROM skills s
			WHERE s.id = sr.skill_id AND s.metadata->>'domain' = $%d
		)`, argIdx)
		args = append(args, domain)
		argIdx++
	}

	sql += " ORDER BY sr.coverage ASC, sr.skill_name"

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("get missing skills: %w", err)
	}
	defer rows.Close()

	var entries []models.SkillRegistryEntry
	for rows.Next() {
		var e models.SkillRegistryEntry
		if err := rows.Scan(&e.SkillID, &e.SkillName, &e.MissingDeps, &e.Stale, &e.LastReview, &e.AutoExpand, &e.Coverage); err != nil {
			return nil, fmt.Errorf("scan registry entry: %w", err)
		}
		entries = append(entries, e)
	}

	return entries, nil
}

// GetCoverage returns coverage statistics for a domain.
func (s *Store) GetCoverage(ctx context.Context, domain string) (map[string]interface{}, error) {
	// Count total skills
	var total int
	totalSQL := "SELECT COUNT(*) FROM skills"
	var totalArgs []interface{}
	if domain != "" {
		totalSQL += " WHERE metadata->>'domain' = $1"
		totalArgs = append(totalArgs, domain)
	}
	if err := s.pool.QueryRow(ctx, totalSQL, totalArgs...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count skills: %w", err)
	}

	// Count with dependencies
	var withDeps int
	depSQL := `
		SELECT COUNT(DISTINCT s.id) FROM skills s
		WHERE EXISTS (SELECT 1 FROM skill_dependencies sd WHERE sd.skill_id = s.id)
	`
	var depArgs []interface{}
	if domain != "" {
		depSQL += " AND s.metadata->>'domain' = $1"
		depArgs = append(depArgs, domain)
	}
	if err := s.pool.QueryRow(ctx, depSQL, depArgs...).Scan(&withDeps); err != nil {
		return nil, fmt.Errorf("count with deps: %w", err)
	}

	// Count with evidence
	var withEvidence int
	evSQL := `
		SELECT COUNT(DISTINCT s.id) FROM skills s
		WHERE EXISTS (SELECT 1 FROM evidences e WHERE e.skill_id = s.id)
	`
	var evArgs []interface{}
	if domain != "" {
		evSQL += " AND s.metadata->>'domain' = $1"
		evArgs = append(evArgs, domain)
	}
	if err := s.pool.QueryRow(ctx, evSQL, evArgs...).Scan(&withEvidence); err != nil {
		return nil, fmt.Errorf("count with evidence: %w", err)
	}

	// Average coverage from registry
	var avgCoverage float64
	covSQL := `
		SELECT COALESCE(AVG(sr.coverage), 0.0) FROM skill_registry sr
		JOIN skills s ON sr.skill_id = s.id
	`
	var covArgs []interface{}
	if domain != "" {
		covSQL += " WHERE s.metadata->>'domain' = $1"
		covArgs = append(covArgs, domain)
	}
	if err := s.pool.QueryRow(ctx, covSQL, covArgs...).Scan(&avgCoverage); err != nil {
		return nil, fmt.Errorf("avg coverage: %w", err)
	}

	// Count missing dependencies
	var missingCount int
	missSQL := `
		SELECT COUNT(*) FROM skill_registry sr
		JOIN skills s ON sr.skill_id = s.id
		WHERE array_length(sr.missing_deps, 1) > 0
	`
	var missArgs []interface{}
	if domain != "" {
		missSQL += " AND s.metadata->>'domain' = $1"
		missArgs = append(missArgs, domain)
	}
	if err := s.pool.QueryRow(ctx, missSQL, missArgs...).Scan(&missingCount); err != nil {
		return nil, fmt.Errorf("count missing: %w", err)
	}

	coverage := 0.0
	if total > 0 {
		coverage = avgCoverage
	}

	return map[string]interface{}{
		"domain":               domain,
		"total_skills":         total,
		"skills_with_deps":     withDeps,
		"skills_with_evidence": withEvidence,
		"skills_missing_deps":  missingCount,
		"average_coverage":     coverage,
		"coverage_percentage":  fmt.Sprintf("%.1f%%", coverage*100),
	}, nil
}

// SubmitLearningJob creates a new learning job for project analysis.
func (s *Store) SubmitLearningJob(ctx context.Context, projectPath string, languages []string) (*models.LearningJob, error) {
	job := &models.LearningJob{
		ID:          uuid.New(),
		ProjectPath: projectPath,
		Status:      "pending",
		Languages:   languages,
	}

	// For now, just store in audit log. In production, this would insert into a learning_jobs table.
	details, _ := json.Marshal(map[string]interface{}{
		"project_path": projectPath,
		"languages":    languages,
		"job_id":       job.ID,
	})

	_, err := s.pool.Exec(ctx,
		"INSERT INTO audit_log (event, details) VALUES ($1, $2)",
		"learning_job_submitted", details,
	)
	if err != nil {
		return nil, fmt.Errorf("log learning job: %w", err)
	}

	return job, nil
}

// VectorSearch performs vector similarity search using pgvector.
func (s *Store) VectorSearch(ctx context.Context, embedding []float32, limit int) ([]models.SearchResult, error) {
	vec := pgvector.NewVector(embedding)
	sql := `
		SELECT s.id, s.name, s.version, s.title, s.description, s.content,
		       s.metadata, s.status, s.created_at, s.updated_at,
		       1 - (s.embedding <=> $1) as score
		FROM skills s
		ORDER BY s.embedding <=> $1
		LIMIT $2
	`
	rows, err := s.pool.Query(ctx, sql, vec, limit)
	if err != nil {
		return nil, fmt.Errorf("vector search: %w", err)
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var r models.SearchResult
		var metadata []byte
		if err := rows.Scan(
			&r.Skill.ID, &r.Skill.Name, &r.Skill.Version, &r.Skill.Title,
			&r.Skill.Description, &r.Skill.Content, &metadata,
			&r.Skill.Status, &r.Skill.CreatedAt, &r.Skill.UpdatedAt,
			&r.Score,
		); err != nil {
			return nil, fmt.Errorf("scan vector result: %w", err)
		}
		r.Skill.Metadata = metadata
		results = append(results, r)
	}

	return results, nil
}

// ListSkills returns all skills with optional filtering.
func (s *Store) ListSkills(ctx context.Context, status models.SkillStatus, limit, offset int) ([]models.Skill, error) {
	sql := `
		SELECT id, name, version, title, description, content,
		       metadata, status, created_at, updated_at
		FROM skills
	`
	args := []interface{}{}
	conditions := []string{}

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, status)
	}

	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	sql += fmt.Sprintf(" ORDER BY name LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("list skills: %w", err)
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var sk models.Skill
		var metadata []byte
		if err := rows.Scan(&sk.ID, &sk.Name, &sk.Version, &sk.Title, &sk.Description, &sk.Content, &metadata, &sk.Status, &sk.CreatedAt, &sk.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan skill: %w", err)
		}
		sk.Metadata = metadata
		skills = append(skills, sk)
	}

	return skills, nil
}

// logAudit is a helper for audit logging used by other skill package files.
func (s *Store) logAudit(ctx context.Context, tx pgx.Tx, event string, skillID *uuid.UUID, details map[string]interface{}) error {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("marshal audit details: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO audit_log (ts, event, skill_id, details)
		VALUES ($1, $2, $3, $4)
	`, time.Now().UTC(), event, skillID, detailsJSON)

	return err
}


