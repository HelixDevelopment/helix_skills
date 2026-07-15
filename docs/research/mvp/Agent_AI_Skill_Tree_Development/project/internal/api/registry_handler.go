package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// handleGetRegistry returns the skill registry with health and completeness info.
//
//	GET /api/v1/registry
//
// Query params:
//	- limit:  max items (default 50, max 200)
//	- offset: items to skip (default 0)
func (s *Server) handleGetRegistry(c *gin.Context) {
	limit, offset := parsePagination(c, 50, 200)

	entries, total, err := s.pool.GetRegistry(c.Request.Context(), limit, offset)
	if err != nil {
		zap.L().Error("failed to get registry", zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve registry")
		return
	}

	RespondPaginated(c, http.StatusOK, entries, total, limit, offset)
}

// handleGetMissingDeps returns the missing dependencies for a skill.
//
//	GET /api/v1/registry/missing-deps/:id
func (s *Server) handleGetMissingDeps(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		RespondError(c, http.StatusBadRequest, "Skill ID is required")
		return
	}

	deps, err := s.pool.GetMissingDeps(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondErrorWithCode(c, http.StatusNotFound, "skill_not_found",
				fmt.Sprintf("Skill not found: %s", id))
			return
		}
		zap.L().Error("failed to get missing dependencies",
			zap.String("id", id),
			zap.Error(err),
		)
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve missing dependencies")
		return
	}

	NegotiateResponse(c, http.StatusOK, gin.H{
		"skill_id":      id,
		"missing_deps":  deps,
		"missing_count": len(deps),
	})
}

// handleGetStaleSkills returns skills that are flagged as stale.
//
//	GET /api/v1/registry/stale
//
// Query params:
//	- limit:  max items (default 50, max 200)
//	- offset: items to skip (default 0)
func (s *Server) handleGetStaleSkills(c *gin.Context) {
	limit, offset := parsePagination(c, 50, 200)

	entries, err := s.pool.GetStaleSkills(c.Request.Context(), limit, offset)
	if err != nil {
		zap.L().Error("failed to get stale skills", zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve stale skills")
		return
	}

	NegotiateResponse(c, http.StatusOK, gin.H{
		"stale_skills": entries,
		"count":        len(entries),
	})
}

// handleTriggerReview triggers a manual review for a skill.
//
//	POST /api/v1/registry/review/:id
func (s *Server) handleTriggerReview(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		RespondError(c, http.StatusBadRequest, "Skill ID is required")
		return
	}

	if err := s.pool.TriggerReview(c.Request.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondErrorWithCode(c, http.StatusNotFound, "skill_not_found",
				fmt.Sprintf("Skill not found: %s", id))
			return
		}
		zap.L().Error("failed to trigger review",
			zap.String("id", id),
			zap.Error(err),
		)
		RespondError(c, http.StatusInternalServerError, "Failed to trigger review")
		return
	}

	NegotiateResponse(c, http.StatusAccepted, gin.H{
		"skill_id": id,
		"status":   "review_triggered",
		"message":  "Review has been scheduled for the skill",
	})
}

// handleGetCoverage returns the overall coverage report for the skill graph.
//
//	GET /api/v1/registry/coverage
func (s *Server) handleGetCoverage(c *gin.Context) {
	report, err := s.pool.GetCoverage(c.Request.Context())
	if err != nil {
		zap.L().Error("failed to get coverage report", zap.Error(err))
		RespondError(c, http.StatusInternalServerError, "Failed to retrieve coverage report")
		return
	}

	NegotiateResponse(c, http.StatusOK, report)
}

// parsePagination parses limit and offset query parameters with defaults and maximums.
func parsePagination(c *gin.Context, defaultLimit, maxLimit int) (int, int) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if err != nil || limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}
