// Package autoexpand provides LLM integration for skill drafting in the
// HelixKnowledge auto-growth pipeline.
package autoexpand

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/helixdevelopment/skill-system/internal/models"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

// ---------------------------------------------------------------------------
// LLMClient interface
// ---------------------------------------------------------------------------

// LLMClient abstracts LLM API calls for skill generation.
type LLMClient interface {
	// Generate creates text from a prompt with a token limit.
	Generate(ctx context.Context, prompt string, maxTokens int) (string, error)
}

// ---------------------------------------------------------------------------
// OpenAI LLM implementation
// ---------------------------------------------------------------------------

// OpenAILLM implements LLMClient using the OpenAI API with rate limiting.
type OpenAILLM struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger

	// Rate limiting
	sem       *semaphore.Weighted
	lastCall  time.Time
	rateMu    sync.Mutex
	rpm       int // requests per minute
}

// NewOpenAILLM creates a new OpenAI LLM client.
func NewOpenAILLM(apiKey, model string, logger *zap.Logger) *OpenAILLM {
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &OpenAILLM{
		apiKey:     apiKey,
		model:      model,
		baseURL:    "https://api.openai.com/v1",
		httpClient: &http.Client{Timeout: 120 * time.Second},
		logger:     logger,
		sem:        semaphore.NewWeighted(10), // max 10 concurrent requests
		rpm:        50,                        // 50 requests per minute
	}
}

// SetBaseURL allows overriding the API base URL (for tests or proxies).
func (c *OpenAILLM) SetBaseURL(url string) {
	c.baseURL = url
}

// SetHTTPClient replaces the default HTTP client.
func (c *OpenAILLM) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// SetRateLimit sets the requests-per-minute rate limit.
func (c *OpenAILLM) SetRateLimit(rpm int) {
	c.rateMu.Lock()
	c.rpm = rpm
	c.rateMu.Unlock()
}

// ---------------------------------------------------------------------------
// Core generation
// ---------------------------------------------------------------------------

// Generate calls the OpenAI chat completions API with rate limiting.
func (c *OpenAILLM) Generate(ctx context.Context, prompt string, maxTokens int) (string, error) {
	// Acquire rate limit semaphore
	if err := c.sem.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("rate limit acquire: %w", err)
	}
	defer c.sem.Release(1)

	// Enforce RPM rate limit
	c.rateMu.Lock()
	minInterval := time.Minute / time.Duration(max(c.rpm, 1))
	elapsed := time.Since(c.lastCall)
	if elapsed < minInterval {
		c.rateMu.Unlock()
		sleepTime := minInterval - elapsed
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(sleepTime):
		}
	} else {
		c.rateMu.Unlock()
	}

	c.rateMu.Lock()
	c.lastCall = time.Now()
	c.rateMu.Unlock()

	reqBody := openAIChatRequest{
		Model: c.model,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens:   maxTokens,
		Temperature: 0.7,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20)) // 10 MiB cap
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result openAIChatResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in API response")
	}

	return result.Choices[0].Message.Content, nil
}

// ---------------------------------------------------------------------------
// Skill draft generation
// ---------------------------------------------------------------------------

// GenerateSkillDraft creates a complete skill draft using the LLM.
// It generates a skill definition with title, description, content,
// metadata, and suggested resources.
func (c *OpenAILLM) GenerateSkillDraft(ctx context.Context, skillName, existingContext string) (*models.Skill, []models.Resource, error) {
	prompt := GeneratePrompt(skillName, existingContext)

	c.logger.Debug("generating skill draft", zap.String("skill", skillName))

	response, err := c.Generate(ctx, prompt, 4000)
	if err != nil {
		return nil, nil, fmt.Errorf("LLM generation: %w", err)
	}

	// Parse the JSON response
	draft, resources, err := parseSkillDraft(skillName, response)
	if err != nil {
		return nil, nil, fmt.Errorf("parse draft: %w", err)
	}

	c.logger.Info("skill draft generated",
		zap.String("skill", draft.Name),
		zap.String("title", draft.Title),
		zap.Int("resources", len(resources)),
	)

	return draft, resources, nil
}

// ---------------------------------------------------------------------------
// Prompt engineering
// ---------------------------------------------------------------------------

// GeneratePrompt creates the system prompt for skill generation.
// The prompt is carefully crafted to produce structured, consistent output
// that can be parsed into skill objects.
func GeneratePrompt(skillName, existingContext string) string {
	return fmt.Sprintf(`You are a technical knowledge graph expert. Generate a skill definition for "%s".

Context from the existing knowledge graph:
%s

Generate a complete skill definition in the following JSON format:

{
  "name": "%s",
  "version": "0.1.0",
  "title": "Human-readable title",
  "description": "Clear, concise description of what this skill covers (2-3 sentences)",
  "content": "# Title\n\n## Overview\n\nDetailed markdown content covering key concepts, patterns, and best practices. Include code examples where relevant.\n\n## Key Concepts\n\n- Concept 1\n- Concept 2\n\n## Best Practices\n\n1. Practice 1\n2. Practice 2\n\n## Common Pitfalls\n\n- Pitfall 1\n- Pitfall 2\n\n## Resources\n\n- [Resource Title](URL)\n",
  "metadata": {
    "tags": ["tag1", "tag2"],
    "domain": "software-engineering",
    "complexity": "intermediate"
  },
  "resources": [
    {
      "url": "https://example.com/resource",
      "title": "Resource Title",
      "resource_type": "documentation"
    }
  ]
}

Requirements:
- Content must be accurate, specific, and actionable
- Include practical code examples where appropriate
- Follow markdown best practices
- Tags should be lowercase, hyphen-separated
- Domain should be a broad category like "software-engineering", "data-science", "devops"
- Complexity should be one of: "beginner", "intermediate", "advanced"
- Resources should be real, verifiable URLs when possible
- Do NOT fabricate URLs - leave resource array empty if you don't know valid URLs
- The content should be comprehensive enough to be genuinely useful

Respond with ONLY the JSON object, no markdown code fences or additional text.`, skillName, existingContext, skillName)
}

// ---------------------------------------------------------------------------
// Parsing
// ---------------------------------------------------------------------------

// parseSkillDraft parses the LLM JSON response into a Skill and Resources.
func parseSkillDraft(skillName, response string) (*models.Skill, []models.Resource, error) {
	// Clean up response: remove markdown code fences if present
	cleaned := response
	if idx := bytes.Index([]byte(cleaned), []byte("```json")); idx != -1 {
		start := idx + 7
		if endIdx := bytes.Index([]byte(cleaned[start:]), []byte("```")); endIdx != -1 {
			cleaned = cleaned[start : start+endIdx]
		}
	} else if idx := bytes.Index([]byte(cleaned), []byte("```")); idx != -1 {
		start := idx + 3
		if endIdx := bytes.Index([]byte(cleaned[start:]), []byte("```")); endIdx != -1 {
			cleaned = cleaned[start : start+endIdx]
		}
	}

	// Trim whitespace
	cleaned = string(bytes.TrimSpace([]byte(cleaned)))

	var raw struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Content     string `json:"content"`
		Metadata    struct {
			Tags       []string `json:"tags"`
			Domain     string   `json:"domain"`
			Complexity string   `json:"complexity"`
		} `json:"metadata"`
		Resources []struct {
			URL          string `json:"url"`
			Title        string `json:"title"`
			ResourceType string `json:"resource_type"`
		} `json:"resources"`
	}

	if err := json.Unmarshal([]byte(cleaned), &raw); err != nil {
		return nil, nil, fmt.Errorf("unmarshal LLM response: %w", err)
	}

	// Build metadata JSON
	metadata := models.SkillMetadata{
		Tags:       raw.Metadata.Tags,
		Domain:     raw.Metadata.Domain,
		Complexity: raw.Metadata.Complexity,
	}
	metadataJSON, _ := json.Marshal(metadata)

	if raw.Version == "" {
		raw.Version = "0.1.0"
	}

	skill := &models.Skill{
		ID:          uuid.New(),
		Name:        skillName,
		Version:     raw.Version,
		Title:       raw.Title,
		Description: raw.Description,
		Content:     raw.Content,
		Metadata:    metadataJSON,
		Status:      models.SkillStatusDraft,
	}

	var resources []models.Resource
	for _, r := range raw.Resources {
		if r.URL == "" {
			continue // skip empty URLs
		}
		resources = append(resources, models.Resource{
			ID:           uuid.New(),
			URL:          r.URL,
			Title:        r.Title,
			ResourceType: r.ResourceType,
		})
	}

	return skill, resources, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ---------------------------------------------------------------------------
// OpenAI API types
// ---------------------------------------------------------------------------

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type openAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int           `json:"index"`
		Message openAIMessage `json:"message"`
		// FinishReason may be absent in some responses.
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
