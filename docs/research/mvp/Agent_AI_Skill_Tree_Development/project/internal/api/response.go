// Package api provides the HTTP API layer for the HelixKnowledge Skill Graph System.
// It includes handlers, middleware, content negotiation, and response utilities.
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
)

// ResponseFormat indicates the serialization format for API responses.
type ResponseFormat string

const (
	// FormatJSON is the default wire format.
	FormatJSON ResponseFormat = "json"
	// FormatTOML is supported via Accept: application/toml header.
	FormatTOML ResponseFormat = "toml"
)

// contextKey is an unexported type for type-safe context keys.
type contextKey string

const responseFormatKey contextKey = "response_format"

// ErrorResponse is the standardized error envelope returned by the API.
type ErrorResponse struct {
	Error   string `json:"error" toml:"error"`
	Code    string `json:"code" toml:"code"`
	Details string `json:"details,omitempty" toml:"details,omitempty"`
}

// SetResponseFormat stores the negotiated format in the Gin context.
func SetResponseFormat(c *gin.Context, format ResponseFormat) {
	c.Set(string(responseFormatKey), format)
}

// GetResponseFormat retrieves the negotiated format from the Gin context.
// Defaults to JSON when none is set.
func GetResponseFormat(c *gin.Context) ResponseFormat {
	v, exists := c.Get(string(responseFormatKey))
	if !exists {
		return FormatJSON
	}
	if f, ok := v.(ResponseFormat); ok {
		return f
	}
	return FormatJSON
}

// RespondJSON writes a JSON response with the given HTTP status code.
func RespondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// RespondTOML writes a TOML response with the given HTTP status code.
func RespondTOML(c *gin.Context, status int, data interface{}) {
	c.Header("Content-Type", "application/toml; charset=utf-8")
	c.Status(status)
	if err := toml.NewEncoder(c.Writer).Encode(data); err != nil {
		// If TOML encoding fails, fall back to JSON error
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to encode TOML response",
			Code:  "encoding_error",
		})
	}
}

// RespondError writes a structured error response in the negotiated format.
func RespondError(c *gin.Context, status int, message string) {
	resp := ErrorResponse{
		Error: message,
		Code:  http.StatusText(status),
	}
	if format := GetResponseFormat(c); format == FormatTOML {
		RespondTOML(c, status, resp)
		return
	}
	RespondJSON(c, status, resp)
}

// RespondErrorWithCode writes a structured error response with a custom error code.
func RespondErrorWithCode(c *gin.Context, status int, code, message string) {
	resp := ErrorResponse{
		Error: message,
		Code:  code,
	}
	if format := GetResponseFormat(c); format == FormatTOML {
		RespondTOML(c, status, resp)
		return
	}
	RespondJSON(c, status, resp)
}

// NegotiateResponse serializes the response in JSON or TOML based on content negotiation.
func NegotiateResponse(c *gin.Context, status int, data interface{}) {
	if format := GetResponseFormat(c); format == FormatTOML {
		RespondTOML(c, status, data)
		return
	}
	RespondJSON(c, status, data)
}

// PaginatedResponse wraps a list response with pagination metadata.
type PaginatedResponse[T any] struct {
	Data   []T    `json:"data" toml:"data"`
	Total  int    `json:"total" toml:"total"`
	Limit  int    `json:"limit" toml:"limit"`
	Offset int    `json:"offset" toml:"offset"`
}

// RespondPaginated writes a paginated response in the negotiated format.
func RespondPaginated[T any](c *gin.Context, status int, data []T, total, limit, offset int) {
	resp := PaginatedResponse[T]{
		Data:   data,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	NegotiateResponse(c, status, resp)
}
