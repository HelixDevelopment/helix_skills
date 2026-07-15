package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/helixdevelopment/skill-system/internal/models"
)

// CreateSkillRequest is the request body for creating a new skill.
type CreateSkillRequest struct {
	Name        string               `json:"name" toml:"name" binding:"required"`
	Version     string               `json:"version" toml:"version"`
	Title       string               `json:"title" toml:"title" binding:"required"`
	Description string               `json:"description" toml:"description"`
	Content     string               `json:"content" toml:"content" binding:"required"`
	Metadata    json.RawMessage      `json:"metadata" toml:"-"`
	Status      models.SkillStatus   `json:"status" toml:"status"`
	Deps        CreateDepsRequest    `json:"dependencies,omitempty" toml:"dependencies"`
	Resources   []CreateResourceReq  `json:"resources,omitempty" toml:"resources"`
}

// CreateDepsRequest holds dependency definitions.
type CreateDepsRequest struct {
	Requires   []string `json:"requires,omitempty" toml:"requires"`
	Extends    []string `json:"extends,omitempty" toml:"extends"`
	Recommends []string `json:"recommends,omitempty" toml:"recommends"`
}

// CreateResourceReq holds resource definitions for creation.
type CreateResourceReq struct {
	URL          string `json:"url" toml:"url" binding:"required"`
	Title        string `json:"title" toml:"title"`
	ResourceType string `json:"resource_type" toml:"resource_type"`
}

// UpdateSkillRequest is the request body for updating a skill.
type UpdateSkillRequest struct {
	Name        *string              `json:"name,omitempty" toml:"name"`
	Version     *string              `json:"version,omitempty" toml:"version"`
	Title       *string              `json:"title,omitempty" toml:"title"`
	Description *string              `json:"description,omitempty" toml:"description"`
	Content     *string              `json:"content,omitempty" toml:"content"`
	Metadata    *json.RawMessage     `json:"metadata,omitempty" toml:"-"`
	Status      *models.SkillStatus  `json:"status,omitempty" toml:"status"`
}

// ImportSkillsRequest wraps a batch of skills for import.
type ImportSkillsRequest struct {
	Skills []CreateSkillRequest `json:"skills" toml:"skills" binding:"required"`
}

// handleListSkills returns a paginated list of skills.
//
//	GET /api/v1/skills
//
// Query params:
//	- limit:  max items to return (default 20, max 100)
//	- offset: number of items to skip (default 0)
//	- status: filter by status (draft, validated, active, deprecated)
func (s *Server) handleListSkills(c *gin.Context) {
	// Parse limit with default and max
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Parse offset with default
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	// Parse optional status filter
	status := c.Query("status")
	if status != "" {
		validStatuses := map[string]bool{
			"draft": true, "validated": true, "active": true, "deprecated": true,
		}
		if !validStatuses[status] {
			RespondErrorWithCode(c, http.StatusBadRequest, "invalid_status",
				fmt.Sprintf("Invalid status filter: %s. Valid: draft, validated, active, deprecated", status))
			return
		}
	}

	// Fetch from database
	skills, total, err := s.pool.ListSkills(c.Request.Context(), limit, offset, status)
	if err != nil {
		zap.L().Error("failed to list skills", zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve skills")
		return
	}

	RespondPaginated(c, http.StatusOK, skills, total, limit, offset)
}

// handleGetSkill returns a single skill by ID or name.
//
//	GET /api/v1/skills/:id
func (s *Server) handleGetSkill(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		RespondError(c, http.StatusBadRequest, "Skill ID is required")
		return
	}

	// Try UUID first, fallback to name lookup
	var skill *models.Skill
	var err error

	if _, parseErr := uuid.Parse(id); parseErr == nil {
		skill, err = s.pool.GetSkill(c.Request.Context(), id)
	} else {
		skill, err = s.pool.GetSkillByName(c.Request.Context(), id)
	}

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondErrorWithCode(c, http.StatusNotFound, "skill_not_found",
				fmt.Sprintf("Skill not found: %s", id))
			return
		}
		zap.L().Error("failed to get skill", zap.String("id", id), zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve skill")
		return
	}

	NegotiateResponse(c, http.StatusOK, skill)
}

// handleCreateSkill creates a new skill.
//
//	POST /api/v1/skills
func (s *Server) handleCreateSkill(c *gin.Context) {
	var req CreateSkillRequest

	// Parse request body based on content type
	if err := parseRequestBody(c, &req); err != nil {
		RespondErrorWithCode(c, http.StatusBadRequest, "invalid_request",
			fmt.Sprintf("Invalid request body: %s", err.Error()))
		return
	}

	// Validate required fields
	if err := validateCreateRequest(&req); err != nil {
		RespondErrorWithCode(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	// Build skill model
	skill := &models.Skill{
		ID:          uuid.New(),
		Name:        req.Name,
		Version:     defaultString(req.Version, "0.1.0"),
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Metadata:    req.Metadata,
		Status:      defaultStatus(req.Status, models.SkillStatusDraft),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Create in database
	if err := s.pool.CreateSkill(c.Request.Context(), skill); err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "already exists") {
			RespondErrorWithCode(c, http.StatusConflict, "skill_exists",
				fmt.Sprintf("Skill with name '%s' already exists", req.Name))
			return
		}
		zap.L().Error("failed to create skill", zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to create skill")
		return
	}

	NegotiateResponse(c, http.StatusCreated, skill)
}

// handleUpdateSkill updates an existing skill.
//
//	PUT /api/v1/skills/:id
//	PATCH /api/v1/skills/:id
func (s *Server) handleUpdateSkill(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		RespondError(c, http.StatusBadRequest, "Skill ID is required")
		return
	}

	// Verify the skill exists
	existing, err := s.pool.GetSkill(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondErrorWithCode(c, http.StatusNotFound, "skill_not_found",
				fmt.Sprintf("Skill not found: %s", id))
			return
		}
		zap.L().Error("failed to get skill for update", zap.String("id", id), zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve skill")
		return
	}

	var req UpdateSkillRequest
	if err := parseRequestBody(c, &req); err != nil {
		RespondErrorWithCode(c, http.StatusBadRequest, "invalid_request",
			fmt.Sprintf("Invalid request body: %s", err.Error()))
		return
	}

	// Apply updates only for non-nil fields
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Version != nil {
		existing.Version = *req.Version
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Content != nil {
		existing.Content = *req.Content
	}
	if req.Metadata != nil {
		existing.Metadata = *req.Metadata
	}
	if req.Status != nil {
		if !isValidStatus(*req.Status) {
			RespondErrorWithCode(c, http.StatusBadRequest, "invalid_status",
				fmt.Sprintf("Invalid status: %s", *req.Status))
			return
		}
		existing.Status = *req.Status
	}

	existing.UpdatedAt = time.Now().UTC()

	if err := s.pool.UpdateSkill(c.Request.Context(), existing); err != nil {
		zap.L().Error("failed to update skill", zap.String("id", id), zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to update skill")
		return
	}

	NegotiateResponse(c, http.StatusOK, existing)
}

// handleDeleteSkill deletes a skill by ID.
//
//	DELETE /api/v1/skills/:id
func (s *Server) handleDeleteSkill(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		RespondError(c, http.StatusBadRequest, "Skill ID is required")
		return
	}

	// Verify the skill exists first
	if _, err := s.pool.GetSkill(c.Request.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondErrorWithCode(c, http.StatusNotFound, "skill_not_found",
				fmt.Sprintf("Skill not found: %s", id))
			return
		}
		zap.L().Error("failed to check skill existence", zap.String("id", id), zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to delete skill")
		return
	}

	if err := s.pool.DeleteSkill(c.Request.Context(), id); err != nil {
		zap.L().Error("failed to delete skill", zap.String("id", id), zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to delete skill")
		return
	}

	c.Status(http.StatusNoContent)
}

// handleGetSkillTree returns the dependency tree of a skill.
//
//	GET /api/v1/skills/:id/tree
//
// Query params:
//	- max_depth: maximum tree depth (default 5, max 10)
func (s *Server) handleGetSkillTree(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		RespondError(c, http.StatusBadRequest, "Skill ID is required")
		return
	}

	// Parse max depth
	maxDepth, err := strconv.Atoi(c.DefaultQuery("max_depth", "5"))
	if err != nil || maxDepth < 1 {
		maxDepth = 5
	}
	if maxDepth > 10 {
		maxDepth = 10
	}

	tree, err := s.pool.GetSkillTree(c.Request.Context(), id, maxDepth)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondErrorWithCode(c, http.StatusNotFound, "skill_not_found",
				fmt.Sprintf("Skill not found: %s", id))
			return
		}
		zap.L().Error("failed to get skill tree", zap.String("id", id), zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve skill tree")
		return
	}

	NegotiateResponse(c, http.StatusOK, tree)
}

// handleImportSkills imports skills from a batch upload (JSON or TOML).
//
//	POST /api/v1/skills/import
func (s *Server) handleImportSkills(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		RespondErrorWithCode(c, http.StatusBadRequest, "read_error",
			"Failed to read request body")
		return
	}

	// Determine format and parse
	var skills []models.Skill
	bodyStr := string(bodyBytes)
	bodyFormat := c.GetString("body_format")

	switch bodyFormat {
	case "toml":
		var wrapper models.TOMLSkillWrapper
		if err := toml.Unmarshal(bodyBytes, &wrapper); err != nil {
			// Try batch format
			var batch struct {
				Skills []models.TOMLSkillWrapper `toml:"skills"`
			}
			if batchErr := toml.Unmarshal(bodyBytes, &batch); batchErr == nil && len(batch.Skills) > 0 {
				skills = convertTOMLBatch(batch.Skills)
			} else {
				// Single skill fallback
				skill := convertTOMLWrapper(wrapper)
				skills = []models.Skill{skill}
			}
		} else {
			skill := convertTOMLWrapper(wrapper)
			skills = []models.Skill{skill}
		}
	default:
		// Try JSON - single skill or array
		bodyBytes = []byte(bodyStr)
		var single CreateSkillRequest
		if err := json.Unmarshal(bodyBytes, &single); err == nil && single.Name != "" {
			skills = append(skills, createReqToModel(&single))
		} else {
			var batch ImportSkillsRequest
			if err := json.Unmarshal(bodyBytes, &batch); err != nil {
				RespondErrorWithCode(c, http.StatusBadRequest, "parse_error",
					fmt.Sprintf("Failed to parse import body: %s", err.Error()))
				return
			}
			for _, req := range batch.Skills {
				if err := validateCreateRequest(&req); err != nil {
					RespondErrorWithCode(c, http.StatusBadRequest, "validation_error",
						fmt.Sprintf("Invalid skill '%s': %s", req.Name, err.Error()))
					return
				}
				skills = append(skills, createReqToModel(&req))
			}
		}
	}

	if len(skills) == 0 {
		RespondErrorWithCode(c, http.StatusBadRequest, "empty_import",
			"No skills found in import body")
		return
	}

	imported, err := s.pool.ImportSkills(c.Request.Context(), skills)
	if err != nil {
		zap.L().Error("failed to import skills", zap.Int("count", len(skills)), zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to import skills")
		return
	}

	NegotiateResponse(c, http.StatusCreated, gin.H{
		"imported": imported,
		"total":    len(skills),
	})
}

// handleExportSkill exports a single skill in the negotiated format.
//
//	GET /api/v1/skills/:id/export
func (s *Server) handleExportSkill(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		RespondError(c, http.StatusBadRequest, "Skill ID is required")
		return
	}

	skills, err := s.pool.ExportSkills(c.Request.Context(), id)
	if err != nil {
		zap.L().Error("failed to export skill", zap.String("id", id), zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to export skill")
		return
	}

	if len(skills) == 0 {
		RespondErrorWithCode(c, http.StatusNotFound, "skill_not_found",
			fmt.Sprintf("Skill not found: %s", id))
		return
	}

	// For TOML, wrap in the export structure
	if GetResponseFormat(c) == FormatTOML {
		wrapper := exportToTOMLWrapper(&skills[0])
		NegotiateResponse(c, http.StatusOK, wrapper)
		return
	}

	NegotiateResponse(c, http.StatusOK, skills[0])
}

// --- Helper functions ---

// parseRequestBody reads and parses the request body based on Content-Type.
func parseRequestBody(c *gin.Context, dst interface{}) error {
	contentType := c.ContentType()
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	// Re-create body for potential re-reading
	c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	switch {
	case strings.Contains(contentType, "application/toml") || strings.Contains(contentType, "text/x-toml"):
		return toml.Unmarshal(bodyBytes, dst)
	default:
		return json.Unmarshal(bodyBytes, dst)
	}
}

// validateCreateRequest checks required fields on a create request.
func validateCreateRequest(req *CreateSkillRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(req.Content) == "" {
		return fmt.Errorf("content is required")
	}
	// Name must be URL-friendly
	if strings.ContainsAny(req.Name, " /\\?#%") {
		return fmt.Errorf("name contains invalid characters")
	}
	return nil
}

// defaultString returns val if non-empty, otherwise def.
func defaultString(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

// defaultStatus returns val if non-empty, otherwise def.
func defaultStatus(val, def models.SkillStatus) models.SkillStatus {
	if val == "" {
		return def
	}
	return val
}

// isValidStatus checks if a status value is valid.
func isValidStatus(s models.SkillStatus) bool {
	switch s {
	case models.SkillStatusDraft, models.SkillStatusValidated,
		models.SkillStatusActive, models.SkillStatusDeprecated:
		return true
	}
	return false
}

// createReqToModel converts a CreateSkillRequest to a Skill model.
func createReqToModel(req *CreateSkillRequest) models.Skill {
	status := req.Status
	if status == "" {
		status = models.SkillStatusDraft
	}
	return models.Skill{
		ID:          uuid.New(),
		Name:        req.Name,
		Version:     defaultString(req.Version, "0.1.0"),
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Metadata:    req.Metadata,
		Status:      status,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

// convertTOMLWrapper converts a TOML skill wrapper to a Skill model.
func convertTOMLWrapper(w models.TOMLSkillWrapper) models.Skill {
	skill := models.Skill{
		ID:          uuid.New(),
		Name:        w.Skill.Name,
		Version:     defaultString(w.Skill.Version, "0.1.0"),
		Title:       w.Skill.Title,
		Description: w.Skill.Description,
		Content:     w.Skill.Content,
		Status:      models.SkillStatusDraft,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Convert metadata
	if w.Skill.Metadata.Domain != "" || len(w.Skill.Metadata.Tags) > 0 {
		metaBytes, _ := json.Marshal(w.Skill.Metadata)
		skill.Metadata = metaBytes
	}

	// Convert dependencies
	for _, depName := range w.Dependencies.Requires {
		skill.Dependencies = append(skill.Dependencies, models.SkillDependency{
			RelationType: models.DepTypeRequires,
		})
		_ = depName // placeholder for resolution
	}
	for _, depName := range w.Dependencies.Extends {
		skill.Dependencies = append(skill.Dependencies, models.SkillDependency{
			RelationType: models.DepTypeExtends,
		})
		_ = depName
	}
	for _, depName := range w.Dependencies.Recommends {
		skill.Dependencies = append(skill.Dependencies, models.SkillDependency{
			RelationType: models.DepTypeRecommends,
		})
		_ = depName
	}

	// Convert resources
	for _, r := range w.Resources {
		skill.Resources = append(skill.Resources, models.Resource{
			ID:           uuid.New(),
			URL:          r.URL,
			Title:        r.Title,
			ResourceType: r.ResourceType,
		})
	}

	return skill
}

// convertTOMLBatch converts a batch of TOML wrappers to Skill models.
func convertTOMLBatch(wrappers []models.TOMLSkillWrapper) []models.Skill {
	skills := make([]models.Skill, 0, len(wrappers))
	for _, w := range wrappers {
		skills = append(skills, convertTOMLWrapper(w))
	}
	return skills
}

// exportToTOMLWrapper converts a Skill to a TOML export wrapper.
func exportToTOMLWrapper(skill *models.Skill) models.TOMLSkillWrapper {
	var meta models.SkillMetadata
	if skill.Metadata != nil {
		_ = json.Unmarshal(skill.Metadata, &meta)
	}

	wrapper := models.TOMLSkillWrapper{
		Skill: models.TOMLSkillDef{
			Name:        skill.Name,
			Version:     skill.Version,
			Title:       skill.Title,
			Description: skill.Description,
			Content:     skill.Content,
			Metadata:    meta,
		},
	}

	// Group dependencies by type
	for _, dep := range skill.Dependencies {
		switch dep.RelationType {
		case models.DepTypeRequires:
			wrapper.Dependencies.Requires = append(wrapper.Dependencies.Requires, dep.DependsOnName)
		case models.DepTypeExtends:
			wrapper.Dependencies.Extends = append(wrapper.Dependencies.Extends, dep.DependsOnName)
		case models.DepTypeRecommends:
			wrapper.Dependencies.Recommends = append(wrapper.Dependencies.Recommends, dep.DependsOnName)
		}
	}

	// Convert resources
	for _, r := range skill.Resources {
		wrapper.Resources = append(wrapper.Resources, models.TOMLResource{
			URL:          r.URL,
			Title:        r.Title,
			ResourceType: r.ResourceType,
		})
	}

	return wrapper
}
