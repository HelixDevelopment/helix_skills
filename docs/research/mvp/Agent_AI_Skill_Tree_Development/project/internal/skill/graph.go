package skill

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/jackc/pgx/v5"
)

// ---------------------------------------------------------------------------
// Dependency graph operations
// ---------------------------------------------------------------------------

// AddDependency adds a directed edge from skillID to dependsOn with cycle detection.
// The relation type must be one of: requires, extends, recommends.
func (s *Store) AddDependency(ctx context.Context, skillID, dependsOn uuid.UUID, relType models.DependencyType) error {
	if skillID == dependsOn {
		return fmt.Errorf("%w: self-referencing dependency", ErrCycleDetected)
	}

	validTypes := map[models.DependencyType]bool{
		models.DepTypeRequires:   true,
		models.DepTypeExtends:    true,
		models.DepTypeRecommends: true,
	}
	if !validTypes[relType] {
		return fmt.Errorf("%w: invalid relation type %q", ErrInvalidSkill, relType)
	}

	return s.pool.WithTx(ctx, func(tx pgx.Tx) error {
		// Verify both skills exist
		var fromName, toName string
		if err := tx.QueryRow(ctx, `SELECT name FROM skills WHERE id = $1`, skillID).Scan(&fromName); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w: source skill %s", ErrSkillNotFound, skillID)
			}
			return fmt.Errorf("check source skill: %w", err)
		}
		if err := tx.QueryRow(ctx, `SELECT name FROM skills WHERE id = $1`, dependsOn).Scan(&toName); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w: target skill %s", ErrDependencyNotFound, dependsOn)
			}
			return fmt.Errorf("check target skill: %w", err)
		}

		// Check for existing edge
		var exists bool
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM skill_dependencies WHERE skill_id = $1 AND depends_on = $2)
		`, skillID, dependsOn).Scan(&exists); err != nil {
			return fmt.Errorf("check existing dependency: %w", err)
		}
		if exists {
			return fmt.Errorf("dependency already exists: %s -> %s", fromName, toName)
		}

		// Cycle detection: check if dependsOn can already reach skillID (adding edge would create cycle)
		cycle, err := hasCycle(ctx, tx, skillID, dependsOn)
		if err != nil {
			return fmt.Errorf("cycle detection check: %w", err)
		}
		if cycle {
			return fmt.Errorf("%w: adding %s -> %s would create a cycle", ErrCycleDetected, fromName, toName)
		}

		// Insert the edge
		_, err = tx.Exec(ctx, `
			INSERT INTO skill_dependencies (skill_id, depends_on, relation_type)
			VALUES ($1, $2, $3)
		`, skillID, dependsOn, relType)
		if err != nil {
			return fmt.Errorf("insert dependency: %w", err)
		}

		// Update registry: recalculate missing deps for the source skill
		if err := s.recalcMissingDeps(ctx, tx, skillID); err != nil {
			return fmt.Errorf("recalc missing deps: %w", err)
		}

		// Audit log
		if err := s.logAudit(ctx, tx, "dependency.added", &skillID, map[string]interface{}{
			"from":          fromName,
			"to":            toName,
			"relation_type": relType,
		}); err != nil {
			return fmt.Errorf("log audit: %w", err)
		}

		return nil
	})
}

// RemoveDependency removes a directed edge between two skills.
func (s *Store) RemoveDependency(ctx context.Context, skillID, dependsOn uuid.UUID) error {
	return s.pool.WithTx(ctx, func(tx pgx.Tx) error {
		// Get names for audit log
		var fromName, toName string
		_ = tx.QueryRow(ctx, `SELECT name FROM skills WHERE id = $1`, skillID).Scan(&fromName)
		_ = tx.QueryRow(ctx, `SELECT name FROM skills WHERE id = $1`, dependsOn).Scan(&toName)

		cmdTag, err := tx.Exec(ctx, `
			DELETE FROM skill_dependencies WHERE skill_id = $1 AND depends_on = $2
		`, skillID, dependsOn)
		if err != nil {
			return fmt.Errorf("delete dependency: %w", err)
		}
		if cmdTag.RowsAffected() == 0 {
			return fmt.Errorf("dependency not found: %s -> %s", skillID, dependsOn)
		}

		// Recalculate missing deps for the source skill
		if err := s.recalcMissingDeps(ctx, tx, skillID); err != nil {
			return fmt.Errorf("recalc missing deps: %w", err)
		}

		// Audit log
		if err := s.logAudit(ctx, tx, "dependency.removed", &skillID, map[string]interface{}{
			"from": fromName,
			"to":   toName,
		}); err != nil {
			return fmt.Errorf("log audit: %w", err)
		}

		return nil
	})
}

// recalcMissingDeps recalculates and updates the missing_deps array for a skill.
func (s *Store) recalcMissingDeps(ctx context.Context, tx pgx.Tx, skillID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		UPDATE skill_registry
		SET missing_deps = COALESCE((
			SELECT array_agg(s2.name ORDER BY s2.name)
			FROM skill_dependencies sd
			JOIN skills s2 ON s2.id = sd.depends_on
			WHERE sd.skill_id = $1
			  AND NOT EXISTS (
				  SELECT 1 FROM skills s3
				  WHERE s3.id = sd.depends_on
				  AND s3.status IN ('validated', 'active')
			  )
		), '{}')
		WHERE skill_id = $1
	`, skillID)

	return err
}

// hasCycle performs a DFS from toID to see if it can reach fromID.
// If so, adding an edge fromID -> toID would create a cycle.
func hasCycle(ctx context.Context, tx pgx.Tx, fromID, toID uuid.UUID) (bool, error) {
	// Use a PostgreSQL recursive CTE to check reachability from toID to fromID
	var reachable bool
	err := tx.QueryRow(ctx, `
		WITH RECURSIVE reach AS (
			SELECT depends_on AS id
			FROM skill_dependencies
			WHERE skill_id = $1

			UNION

			SELECT sd.depends_on
			FROM skill_dependencies sd
			INNER JOIN reach r ON r.id = sd.skill_id
		)
		SELECT EXISTS(SELECT 1 FROM reach WHERE id = $2)
	`, toID, fromID).Scan(&reachable)

	return reachable, err
}

// GetDependencyTree returns the full dependency tree starting from a root skill name,
// using a PostgreSQL recursive CTE. Results are limited to maxDepth levels.
func (s *Store) GetDependencyTree(ctx context.Context, rootName string, maxDepth int) (*models.SkillTreeNode, error) {
	if maxDepth <= 0 {
		maxDepth = 10
	}
	if maxDepth > 50 {
		maxDepth = 50 // Hard cap to prevent runaway queries
	}

	// First, get the root skill
	rootSkill, err := s.GetByName(ctx, rootName)
	if err != nil {
		return nil, err
	}

	// Use recursive CTE to get all transitive dependencies with depth info
	rows, err := s.pool.Query(ctx, `
		WITH RECURSIVE dep_tree AS (
			SELECT
				s.id,
				s.name,
				s.version,
				s.title,
				s.description,
				s.content,
				s.metadata,
				s.status,
				s.created_at,
				s.updated_at,
				sd.relation_type,
				0 AS depth
			FROM skill_dependencies sd
			JOIN skills s ON s.id = sd.depends_on
			WHERE sd.skill_id = $1

			UNION ALL

			SELECT
				s.id,
				s.name,
				s.version,
				s.title,
				s.description,
				s.content,
				s.metadata,
				s.status,
				s.created_at,
				s.updated_at,
				sd.relation_type,
				dt.depth + 1
			FROM skill_dependencies sd
			JOIN skills s ON s.id = sd.depends_on
			JOIN dep_tree dt ON dt.id = sd.skill_id
			WHERE dt.depth + 1 < $2
		)
		SELECT id, name, version, title, description, content, metadata, status, created_at, updated_at, relation_type, depth
		FROM dep_tree
		ORDER BY depth, name
	`, rootSkill.ID, maxDepth)
	if err != nil {
		return nil, fmt.Errorf("recursive dependency query: %w", err)
	}
	defer rows.Close()

	// Build a flat list of all nodes with their depth and relation type
	type flatNode struct {
		skill        models.Skill
		depth        int
		relationType models.DependencyType
	}

	// Map: skill ID -> flatNode
	nodeMap := make(map[uuid.UUID]*flatNode)
	// Keep insertion order for building tree
	var flatNodes []flatNode

	for rows.Next() {
		var fn flatNode
		var metaJSON []byte
		if err := rows.Scan(
			&fn.skill.ID, &fn.skill.Name, &fn.skill.Version, &fn.skill.Title,
			&fn.skill.Description, &fn.skill.Content, &metaJSON,
			&fn.skill.Status, &fn.skill.CreatedAt, &fn.skill.UpdatedAt,
			&fn.relationType, &fn.depth,
		); err != nil {
			return nil, fmt.Errorf("scan dep tree node: %w", err)
		}
		fn.skill.Metadata = metaJSON
		nodeMap[fn.skill.ID] = &fn
		flatNodes = append(flatNodes, fn)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dep tree: %w", err)
	}

	// Build tree: root -> children
	root := &models.SkillTreeNode{
		Skill:    *rootSkill,
		Depth:    0,
		Children: []models.SkillTreeNode{},
	}

	// For building the tree, we need parent-child relationships.
	// The CTE gives us depth, but not parent. We re-query for parent relationships.
	if len(flatNodes) > 0 {
		parentRows, err := s.pool.Query(ctx, `
			SELECT skill_id, depends_on
			FROM skill_dependencies
			WHERE depends_on = ANY($1)
		`, collectIDs(flatNodes))
		if err != nil {
			return nil, fmt.Errorf("query parent relationships: %w", err)
		}
		defer parentRows.Close()

		childrenMap := make(map[uuid.UUID][]models.SkillTreeNode)
		for parentRows.Next() {
			var parentID, childID uuid.UUID
			if err := parentRows.Scan(&parentID, &childID); err != nil {
				continue
			}
			if node, ok := nodeMap[childID]; ok {
				childrenMap[parentID] = append(childrenMap[parentID], models.SkillTreeNode{
					Skill: node.skill,
					Depth: node.depth,
				})
			}
		}

		// Attach direct children to root
		root.Children = childrenMap[rootSkill.ID]
	}

	return root, nil
}

// collectIDs extracts skill IDs from flat nodes for the parent query.
func collectIDs(nodes []flatNode) []uuid.UUID {
	ids := make([]uuid.UUID, len(nodes))
	for i, n := range nodes {
		ids[i] = n.skill.ID
	}
	return ids
}

// GetDependents returns all skills that directly depend on the given skill.
func (s *Store) GetDependents(ctx context.Context, skillID uuid.UUID) ([]models.Skill, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT s.id, s.name, s.version, s.title, s.description, s.content, s.metadata, s.status, s.created_at, s.updated_at
		FROM skill_dependencies sd
		JOIN skills s ON s.id = sd.skill_id
		WHERE sd.depends_on = $1
		ORDER BY s.name
	`, skillID)
	if err != nil {
		return nil, fmt.Errorf("query dependents: %w", err)
	}
	defer rows.Close()

	return scanSkills(rows)
}

// GetAllDependencies returns a flat list of all transitive dependencies for a skill.
func (s *Store) GetAllDependencies(ctx context.Context, skillID uuid.UUID) ([]models.Skill, error) {
	rows, err := s.pool.Query(ctx, `
		WITH RECURSIVE all_deps AS (
			SELECT depends_on AS id
			FROM skill_dependencies
			WHERE skill_id = $1

			UNION

			SELECT sd.depends_on
			FROM skill_dependencies sd
			INNER JOIN all_deps ad ON ad.id = sd.skill_id
		)
		SELECT s.id, s.name, s.version, s.title, s.description, s.content, s.metadata, s.status, s.created_at, s.updated_at
		FROM all_deps ad
		JOIN skills s ON s.id = ad.id
		ORDER BY s.name
	`, skillID)
	if err != nil {
		return nil, fmt.Errorf("recursive all-deps query: %w", err)
	}
	defer rows.Close()

	return scanSkills(rows)
}

// scanSkills is a helper to scan rows into a slice of Skill.
func scanSkills(rows pgx.Rows) ([]models.Skill, error) {
	var skills []models.Skill
	for rows.Next() {
		var sk models.Skill
		var metaJSON []byte
		if err := rows.Scan(
			&sk.ID, &sk.Name, &sk.Version, &sk.Title,
			&sk.Description, &sk.Content, &metaJSON,
			&sk.Status, &sk.CreatedAt, &sk.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan skill: %w", err)
		}
		sk.Metadata = metaJSON
		skills = append(skills, sk)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate skills: %w", err)
	}

	return skills, nil
}

// ---------------------------------------------------------------------------
// Registry update helpers (called from graph mutations)
// ---------------------------------------------------------------------------

// UpdateRegistryAfterChange recalculates registry state for affected skills.
// Should be called after any skill or dependency mutation.
func (s *Store) UpdateRegistryAfterChange(ctx context.Context, skillID uuid.UUID) error {
	return s.pool.WithTx(ctx, func(tx pgx.Tx) error {
		// Recalculate missing deps for this skill
		if err := s.recalcMissingDeps(ctx, tx, skillID); err != nil {
			return err
		}

		// Recalculate missing deps for all skills that depend on this one
		_, err := tx.Exec(ctx, `
			UPDATE skill_registry sr
			SET missing_deps = COALESCE((
				SELECT array_agg(s2.name ORDER BY s2.name)
				FROM skill_dependencies sd
				JOIN skills s2 ON s2.id = sd.depends_on
				WHERE sd.skill_id = sr.skill_id
				  AND NOT EXISTS (
					  SELECT 1 FROM skills s3
					  WHERE s3.id = sd.depends_on
					  AND s3.status IN ('validated', 'active')
				  )
			), '{}')
			WHERE sr.skill_id IN (
				SELECT skill_id FROM skill_dependencies WHERE depends_on = $1
			)
		`, skillID)

		return err
	})
}

// logAuditGraph is a convenience wrapper for graph-related audit events.
func (s *Store) logAuditGraph(ctx context.Context, tx pgx.Tx, event string, skillID uuid.UUID, details map[string]interface{}) error {
	return s.logAudit(ctx, tx, event, &skillID, details)
}

// logAuditEvent logs a generic audit event without a specific skill ID.
func (s *Store) logAuditEvent(ctx context.Context, tx pgx.Tx, event string, details map[string]interface{}) error {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("marshal audit details: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO audit_log (ts, event, skill_id, details)
		VALUES ($1, $2, NULL, $3)
	`, time.Now().UTC(), event, detailsJSON)

	return err
}
