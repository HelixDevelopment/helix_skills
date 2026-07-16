package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SubmitProjectRequest initiates a codebase analysis job.
type SubmitProjectRequest struct {
	ProjectPath string   `json:"project_path" toml:"project_path" binding:"required"`
	Languages   []string `json:"languages" toml:"languages"`
}

// handleSubmitProject submits a project for codebase analysis and learning.
//
//	POST /api/v1/learn/projects
//
// Request body:
//   - project_path: absolute or relative path to the project directory (required)
//   - languages:    list of languages to analyze (optional, auto-detected if empty)
func (s *Server) handleSubmitProject(c *gin.Context) {
	var req SubmitProjectRequest
	if err := parseRequestBody(c, &req); err != nil {
		RespondErrorWithCode(c, http.StatusBadRequest, "invalid_request",
			fmt.Sprintf("Invalid request body: %s", err.Error()))
		return
	}

	if strings.TrimSpace(req.ProjectPath) == "" {
		RespondErrorWithCode(c, http.StatusBadRequest, "missing_project_path",
			"project_path is required")
		return
	}

	// Normalize languages
	languages := req.Languages
	if len(languages) == 0 {
		// Auto-detect will be handled by the worker
		languages = nil
	}

	// Validate path doesn't contain dangerous patterns
	if strings.Contains(req.ProjectPath, "..") || strings.Contains(req.ProjectPath, "~") {
		RespondErrorWithCode(c, http.StatusBadRequest, "invalid_path",
			"Project path contains invalid characters")
		return
	}

	job, err := s.pool.SubmitProject(c.Request.Context(), req.ProjectPath, languages)
	if err != nil {
		zap.L().Error("failed to submit project",
			zap.String("project_path", req.ProjectPath),
			zap.Strings("languages", languages),
			zap.Error(err),
		)
		RespondError(c, http.StatusInternalServerError, "Failed to submit project for analysis")
		return
	}

	NegotiateResponse(c, http.StatusAccepted, gin.H{
		"job":     job,
		"message": fmt.Sprintf("Project '%s' submitted for analysis", req.ProjectPath),
	})
}

// handleGetLearnStatus returns the status of a learning job.
//
//	GET /api/v1/learn/status/:id
func (s *Server) handleGetLearnStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		RespondError(c, http.StatusBadRequest, "Job ID is required")
		return
	}

	job, err := s.pool.GetLearnStatus(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondErrorWithCode(c, http.StatusNotFound, "job_not_found",
				fmt.Sprintf("Learning job not found: %s", id))
			return
		}
		zap.L().Error("failed to get learning status",
			zap.String("job_id", id),
			zap.Error(err),
		)
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve learning status")
		return
	}

	NegotiateResponse(c, http.StatusOK, gin.H{
		"job": job,
	})
}

// handleGetEvidences returns learning evidences for a skill.
//
//	GET /api/v1/learn/evidences/:skill_id
//
// Query params:
//   - limit:  max items (default 20, max 100)
//   - offset: items to skip (default 0)
func (s *Server) handleGetEvidences(c *gin.Context) {
	skillID := c.Param("skill_id")
	if skillID == "" {
		RespondError(c, http.StatusBadRequest, "Skill ID is required")
		return
	}

	limit, offset := parsePagination(c, 20, 100)

	evidences, total, err := s.pool.GetEvidences(c.Request.Context(), skillID, limit, offset)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondErrorWithCode(c, http.StatusNotFound, "skill_not_found",
				fmt.Sprintf("Skill not found: %s", skillID))
			return
		}
		zap.L().Error("failed to get evidences",
			zap.String("skill_id", skillID),
			zap.Error(err),
		)
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve evidences")
		return
	}

	RespondPaginated(c, http.StatusOK, evidences, total, limit, offset)
}
