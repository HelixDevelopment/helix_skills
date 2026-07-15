package skill

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/jackc/pgx/v5"
)

// ---------------------------------------------------------------------------
// TOML Import / Export
// ---------------------------------------------------------------------------

// ImportFromTOML parses a TOML skill definition and creates the skill along with
// its dependencies and resources in a single transaction.
func (s *Store) ImportFromTOML(ctx context.Context, tomlData []byte) (*models.Skill, error) {
	var wrapper models.TOMLSkillWrapper
	if err := toml.Unmarshal(tomlData, &wrapper); err != nil {
		return nil, fmt.Errorf("parse TOML: %w", err)
	}

	// Validate required fields
	if wrapper.Skill.Name == "" {
		return nil, fmt.Errorf("%w: skill.name is required", ErrInvalidSkill)
	}
	if wrapper.Skill.Title == "" {
		return nil, fmt.Errorf("%w: skill.title is required", ErrInvalidSkill)
	}
	if wrapper.Skill.Content == "" {
		return nil, fmt.Errorf("%w: skill.content is required", ErrInvalidSkill)
	}

	// Check if skill already exists
	existing, err := s.GetByName(ctx, wrapper.Skill.Name)
	if err != nil && !errors.Is(err, ErrSkillNotFound) {
		return nil, fmt.Errorf("check existing skill: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: %s", ErrSkillExists, wrapper.Skill.Name)
	}

	// Marshal metadata
	metadata := models.SkillMetadata{
		Tags:       wrapper.Skill.Metadata.Tags,
		Domain:     wrapper.Skill.Metadata.Domain,
		Complexity: wrapper.Skill.Metadata.Complexity,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

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

	// Resolve dependency names to IDs
	depNames := collectDepNames(wrapper.Dependencies)
	depNameToID := make(map[string]uuid.UUID)

	if len(depNames) > 0 {
		// Query for existing skills that match dependency names
		rows, err := s.pool.Query(ctx, `
			SELECT id, name FROM skills WHERE name = ANY($1)
		`, depNames)
		if err != nil {
			return nil, fmt.Errorf("resolve dependency names: %w", err)
		}
		for rows.Next() {
			var id uuid.UUID
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				rows.Close()
				return nil, fmt.Errorf("scan dep name: %w", err)
			}
			depNameToID[name] = id
		}
		rows.Close()
	}

	// Build dependency records
	var depsToCreate []struct {
		targetID     uuid.UUID
		relationType models.DependencyType
		name         string // for error reporting
	}

	for _, depName := range wrapper.Dependencies.Requires {
		id, ok := depNameToID[depName]
		if !ok {
			return nil, fmt.Errorf("%w: required dependency %q not found", ErrDependencyNotFound, depName)
		}
		depsToCreate = append(depsToCreate, struct {
			targetID     uuid.UUID
			relationType models.DependencyType
			name         string
		}{id, models.DepTypeRequires, depName})
	}
	for _, depName := range wrapper.Dependencies.Extends {
		id, ok := depNameToID[depName]
		if !ok {
			return nil, fmt.Errorf("%w: extends dependency %q not found", ErrDependencyNotFound, depName)
		}
		depsToCreate = append(depsToCreate, struct {
			targetID     uuid.UUID
			relationType models.DependencyType
			name         string
		}{id, models.DepTypeExtends, depName})
	}
	for _, depName := range wrapper.Dependencies.Recommends {
		id, ok := depNameToID[depName]
		if !ok {
			// For recommends, we allow missing deps (soft dependency)
			continue
		}
		depsToCreate = append(depsToCreate, struct {
			targetID     uuid.UUID
			relationType models.DependencyType
			name         string
		}{id, models.DepTypeRecommends, depName})
	}

	// Build resource records
	resources := make([]models.Resource, len(wrapper.Resources))
	for i, r := range wrapper.Resources {
		resources[i] = models.Resource{
			ID:           uuid.New(),
			URL:          r.URL,
			Title:        r.Title,
			ResourceType: r.ResourceType,
		}
	}

	// Execute everything in a transaction
	return skill, s.pool.WithTx(ctx, func(tx pgx.Tx) error {
		// Insert skill
		_, err := tx.Exec(ctx, `
			INSERT INTO skills (id, name, version, title, description, content, metadata, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		`, skill.ID, skill.Name, skill.Version, skill.Title, skill.Description, skill.Content, skill.Metadata, skill.Status)
		if err != nil {
			return fmt.Errorf("insert skill: %w", err)
		}

		// Initialize registry entry
		_, err = tx.Exec(ctx, `
			INSERT INTO skill_registry (skill_id, skill_name, missing_deps, stale, last_review, auto_expand, coverage)
			VALUES ($1, $2, '{}', false, NOW(), true, 0.0)
		`, skill.ID, skill.Name)
		if err != nil {
			return fmt.Errorf("insert registry entry: %w", err)
		}

		// Insert dependencies (with cycle detection per edge)
		for _, dep := range depsToCreate {
			// Cycle check: ensure adding skill -> dep.targetID doesn't create a cycle
			cycle, err := hasCycle(ctx, tx, skill.ID, dep.targetID)
			if err != nil {
				return fmt.Errorf("cycle check for %s: %w", dep.name, err)
			}
			if cycle {
				return fmt.Errorf("%w: adding dependency on %s would create cycle", ErrCycleDetected, dep.name)
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO skill_dependencies (skill_id, depends_on, relation_type)
				VALUES ($1, $2, $3)
			`, skill.ID, dep.targetID, dep.relationType)
			if err != nil {
				return fmt.Errorf("insert dependency %s: %w", dep.name, err)
			}
		}

		// Insert resources
		for _, r := range resources {
			r.SkillID = skill.ID
			_, err := tx.Exec(ctx, `
				INSERT INTO resources (id, skill_id, url, title, resource_type, created_at)
				VALUES ($1, $2, $3, $4, $5, NOW())
			`, r.ID, r.SkillID, r.URL, r.Title, r.ResourceType)
			if err != nil {
				return fmt.Errorf("insert resource %s: %w", r.URL, err)
			}
		}

		// Recalculate registry state
		if err := s.recalcMissingDeps(ctx, tx, skill.ID); err != nil {
			return fmt.Errorf("recalc missing deps: %w", err)
		}
		if err := s.recalcCoverage(ctx, tx, skill.ID); err != nil {
			return fmt.Errorf("recalc coverage: %w", err)
		}

		// Audit log
		if err := s.logAudit(ctx, tx, "skill.imported", &skill.ID, map[string]interface{}{
			"name":            skill.Name,
			"version":         skill.Version,
			"deps_count":      len(depsToCreate),
			"resources_count": len(resources),
		}); err != nil {
			return fmt.Errorf("log audit: %w", err)
		}

		// Set runtime fields on the returned skill
		skill.Dependencies = make([]models.SkillDependency, 0, len(depsToCreate))
		for _, dep := range depsToCreate {
			skill.Dependencies = append(skill.Dependencies, models.SkillDependency{
				SkillID:      skill.ID,
				DependsOn:    dep.targetID,
				RelationType: dep.relationType,
			})
		}
		skill.Resources = resources

		return nil
	})
}

// ExportToTOML exports a skill and its dependencies and resources as TOML.
func (s *Store) ExportToTOML(ctx context.Context, skillName string) ([]byte, error) {
	skill, err := s.GetByName(ctx, skillName)
	if err != nil {
		return nil, err
	}

	// Parse metadata for TOML
	var meta models.SkillMetadata
	if len(skill.Metadata) > 0 {
		if err := json.Unmarshal(skill.Metadata, &meta); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}

	// Build TOML wrapper
	wrapper := models.TOMLSkillWrapper{
		Skill: models.TOMLSkillDef{
			Name:        skill.Name,
			Version:     skill.Version,
			Title:       skill.Title,
			Description: skill.Description,
			Content:     skill.Content,
			Metadata:    meta,
		},
		Resources: make([]models.TOMLResource, len(skill.Resources)),
	}

	// Categorize dependencies by relation type
	for _, dep := range skill.Dependencies {
		depName := dep.DependsOnName
		if depName == "" {
			// Resolve name from ID if not populated
			var name string
			_ = s.pool.QueryRow(ctx, `SELECT name FROM skills WHERE id = $1`, dep.DependsOn).Scan(&name)
			depName = name
		}

		switch dep.RelationType {
		case models.DepTypeRequires:
			wrapper.Dependencies.Requires = append(wrapper.Dependencies.Requires, depName)
		case models.DepTypeExtends:
			wrapper.Dependencies.Extends = append(wrapper.Dependencies.Extends, depName)
		case models.DepTypeRecommends:
			wrapper.Dependencies.Recommends = append(wrapper.Dependencies.Recommends, depName)
		}
	}

	// Map resources
	for i, r := range skill.Resources {
		wrapper.Resources[i] = models.TOMLResource{
			URL:          r.URL,
			Title:        r.Title,
			ResourceType: r.ResourceType,
		}
	}

	// Encode to TOML
	var buf bytes.Buffer
	buf.WriteString("# Skill definition exported from HelixKnowledge\n")
	buf.WriteString(fmt.Sprintf("# Skill: %s (v%s)\n", skill.Name, skill.Version))
	buf.WriteString(fmt.Sprintf("# Exported at: %s\n\n", skill.UpdatedAt.Format("2006-01-02T15:04:05Z")))

	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(wrapper); err != nil {
		return nil, fmt.Errorf("encode TOML: %w", err)
	}

	return buf.Bytes(), nil
}

// collectDepNames gathers all dependency names from the TOML dependency block.
func collectDepNames(deps models.TOMLDependencies) []string {
	seen := make(map[string]bool)
	var names []string

	for _, n := range deps.Requires {
		if !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}
	for _, n := range deps.Extends {
		if !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}
	for _, n := range deps.Recommends {
		if !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}

	return names
}
