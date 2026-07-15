package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TriggerExpandRequest initiates an auto-expansion job.
type TriggerExpandRequest struct {
	SkillName string `json:"skill_name" toml:"skill_name" binding:"required"`
	Depth     int    `json:"depth" toml:"depth"`
}

// handleTriggerExpand starts an auto-expansion job for a skill.
//
//	POST /api/v1/expand
//
// Request body:
//	- skill_name: name of the skill to expand from (required)
//	- depth:      expansion depth (default 3, max 5)
func (s *Server) handleTriggerExpand(c *gin.Context) {
	var req TriggerExpandRequest
	if err := parseRequestBody(c, &req); err != nil {
		RespondErrorWithCode(c, http.StatusBadRequest, "invalid_request",
			fmt.Sprintf("Invalid request body: %s", err.Error()))
		return
	}

	if strings.TrimSpace(req.SkillName) == "" {
		RespondErrorWithCode(c, http.StatusBadRequest, "missing_skill_name",
			"skill_name is required")
		return
	}

	// Parse and validate depth
	depth := req.Depth
	if depth < 1 {
		depth = 3
	}
	if depth > 5 {
		depth = 5
	}

	job, err := s.pool.TriggerExpand(c.Request.Context(), req.SkillName, depth)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondErrorWithCode(c, http.StatusNotFound, "skill_not_found",
				fmt.Sprintf("Skill not found: %s", req.SkillName))
			return
		}
		zap.L().Error("failed to trigger expansion",
			zap.String("skill_name", req.SkillName),
			zap.Int("depth", depth),
			zap.Error(err),
		)
		RespondError(c, http.StatusInternalServerError, "Failed to trigger expansion")
		return
	}

	NegotiateResponse(c, http.StatusAccepted, gin.H{
		"job":     job,
		"message": fmt.Sprintf("Expansion started for skill '%s' with depth %d", req.SkillName, depth),
	})
}

// handleGetExpandStatus returns the status of an expansion job.
//
//	GET /api/v1/expand/status/:id
func (s *Server) handleGetExpandStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		RespondError(c, http.StatusBadRequest, "Job ID is required")
		return
	}

	job, err := s.pool.GetExpandStatus(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondErrorWithCode(c, http.StatusNotFound, "job_not_found",
				fmt.Sprintf("Expansion job not found: %s", id))
			return
		}
		zap.L().Error("failed to get expansion status",
			zap.String("job_id", id),
			zap.Error(err),
		)
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve expansion status")
		return
	}

	NegotiateResponse(c, http.StatusOK, gin.H{
		"job": job,
	})
}

// handleGetGapReport returns a report of gaps in the skill graph.
//
//	GET /api/v1/expand/gaps
func (s *Server) handleGetGapReport(c *gin.Context) {
	report, err := s.pool.GetGapReport(c.Request.Context())
	if err != nil {
		zap.L().Error("failed to get gap report", zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve gap report")
		return
	}

	NegotiateResponse(c, http.StatusOK, report)
}
