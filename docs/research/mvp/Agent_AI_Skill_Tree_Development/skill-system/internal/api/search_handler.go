package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SearchRequest is the request body for vector + keyword hybrid search.
type SearchRequest struct {
	Query  string    `json:"query" toml:"query" binding:"required"`
	Vector []float32 `json:"vector,omitempty" toml:"vector"`
	Limit  int       `json:"limit" toml:"limit"`
}

// SimilarSkillsRequest searches for skills similar to provided content.
type SimilarSkillsRequest struct {
	Content string    `json:"content" toml:"content" binding:"required"`
	Vector  []float32 `json:"vector,omitempty" toml:"vector"`
	Limit   int       `json:"limit" toml:"limit"`
}

// handleSearch performs a hybrid vector + keyword search across skills.
//
//	GET /api/v1/search
//
// Query params (for GET):
//	- q:      search query text (required for keyword search)
//	- limit:  max results (default 20, max 100)
//	- vector: comma-separated vector values (optional, for vector search)
func (s *Server) handleSearch(c *gin.Context) {
	var query string
	var vector []float32
	var limit int

	// Support both GET query params and POST JSON body
	if c.Request.Method == "POST" || c.ContentType() != "" {
		var req SearchRequest
		if err := parseRequestBody(c, &req); err == nil {
			query = req.Query
			vector = req.Vector
			limit = req.Limit
		}
	}

	// Fallback to query parameters
	if query == "" {
		query = c.Query("q")
	}
	if len(vector) == 0 {
		if vStr := c.Query("vector"); vStr != "" {
			vector = parseVector(vStr)
		}
	}
	if limit == 0 {
		var err error
		limit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
		if err != nil || limit < 1 {
			limit = 20
		}
	}

	if limit > 100 {
		limit = 100
	}

	// Validate that at least one search method is provided
	if query == "" && len(vector) == 0 {
		RespondErrorWithCode(c, http.StatusBadRequest, "missing_query",
			"Search requires a query parameter 'q' or a vector")
		return
	}

	results, err := s.pool.SearchSkills(c.Request.Context(), query, vector, limit)
	if err != nil {
		zap.L().Error("search failed",
			zap.String("query", query),
			zap.Int("vector_dim", len(vector)),
			zap.Error(err),
		)
		RespondError(c, http.StatusInternalServerError, "Search operation failed")
		return
	}

	NegotiateResponse(c, http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
		"query":   query,
	})
}

// handleSimilarSkills finds skills similar to the provided content or skill ID.
//
//	POST /api/v1/search/similar
//
// Request body (JSON or TOML):
//	- content:  text content to find similar skills for
//	- skill_id: existing skill ID to find similar skills to
//	- vector:   embedding vector (optional)
//	- limit:    max results (default 20, max 100)
func (s *Server) handleSimilarSkills(c *gin.Context) {
	var req SimilarSkillsRequest
	if err := parseRequestBody(c, &req); err != nil {
		// Try alternative format with skill_id
		var altReq struct {
			SkillID string    `json:"skill_id" toml:"skill_id"`
			Vector  []float32 `json:"vector,omitempty" toml:"vector"`
			Limit   int       `json:"limit" toml:"limit"`
		}
		if altErr := parseRequestBody(c, &altReq); altErr != nil {
			RespondErrorWithCode(c, http.StatusBadRequest, "invalid_request",
				fmt.Sprintf("Invalid request body: %s", err.Error()))
			return
		}
		if altReq.SkillID == "" {
			RespondErrorWithCode(c, http.StatusBadRequest, "missing_content",
				"Either 'content' or 'skill_id' is required")
			return
		}

		// Use skill_id-based similarity
		limit := altReq.Limit
		if limit < 1 {
			limit = 20
		}
		if limit > 100 {
			limit = 100
		}

		results, err := s.pool.SimilarSkills(c.Request.Context(), altReq.SkillID, limit)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				RespondErrorWithCode(c, http.StatusNotFound, "skill_not_found",
					fmt.Sprintf("Skill not found: %s", altReq.SkillID))
				return
			}
			zap.L().Error("similar skills search failed",
				zap.String("skill_id", altReq.SkillID),
				zap.Error(err),
			)
			RespondError(c, http.StatusInternalServerError, "Similarity search failed")
			return
		}

		NegotiateResponse(c, http.StatusOK, gin.H{
			"results":  results,
			"count":    len(results),
			"skill_id": altReq.SkillID,
		})
		return
	}

	// Validate content
	if strings.TrimSpace(req.Content) == "" {
		RespondErrorWithCode(c, http.StatusBadRequest, "missing_content",
			"Content is required for similarity search")
		return
	}

	// Parse limit
	limit := req.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Search using content as query text
	results, err := s.pool.SearchSkills(c.Request.Context(), req.Content, req.Vector, limit)
	if err != nil {
		zap.L().Error("similar skills search failed",
			zap.Int("content_len", len(req.Content)),
			zap.Error(err),
		)
		RespondError(c, http.StatusInternalServerError, "Similarity search failed")
		return
	}

	NegotiateResponse(c, http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
		"content": req.Content,
	})
}

// parseVector converts a comma-separated string of floats to a float32 slice.
func parseVector(s string) []float32 {
	parts := strings.Split(s, ",")
	vec := make([]float32, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if f, err := strconv.ParseFloat(p, 32); err == nil {
			vec = append(vec, float32(f))
		}
	}
	return vec
}

