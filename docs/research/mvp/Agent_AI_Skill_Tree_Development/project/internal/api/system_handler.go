package api

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Build-time variables (injected via ldflags).
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

// System metrics.
var (
	uptimeGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "helix",
		Subsystem: "api",
		Name:      "uptime_seconds",
		Help:      "Number of seconds since the API server started",
	})

	goroutineGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "helix",
		Subsystem: "api",
		Name:      "goroutines_count",
		Help:      "Number of active goroutines",
	})

	memoryGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "helix",
		Subsystem: "api",
		Name:      "memory_usage_bytes",
		Help:      "Current memory usage in bytes",
	})
)

// serverStartTime records when the server started for uptime calculation.
var serverStartTime = time.Now()

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status    string            `json:"status" toml:"status"`
	Timestamp time.Time         `json:"timestamp" toml:"timestamp"`
	Uptime    string            `json:"uptime" toml:"uptime"`
	Checks    map[string]string `json:"checks,omitempty" toml:"checks,omitempty"`
	Version   string            `json:"version" toml:"version"`
}

// VersionResponse represents the version information response.
type VersionResponse struct {
	Version   string `json:"version" toml:"version"`
	Commit    string `json:"commit" toml:"commit"`
	BuildTime string `json:"build_time" toml:"build_time"`
	GoVersion string `json:"go_version" toml:"go_version"`
	Platform  string `json:"platform" toml:"platform"`
}

// handleHealth returns the health status of the API server and its dependencies.
//
//	GET /health
func (s *Server) handleHealth(c *gin.Context) {
	ctx := c.Request.Context()
	checks := make(map[string]string)

	// Check database connectivity
	dbStatus := "ok"
	if err := s.pool.Ping(ctx); err != nil {
		dbStatus = "error: " + err.Error()
		zap.L().Warn("health check: database ping failed", zap.Error(err))
	}
	checks["database"] = dbStatus

	// Overall status
	status := "healthy"
	if dbStatus != "ok" {
		status = "degraded"
	}

	// Update system metrics
	updateSystemMetrics()

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC(),
		Uptime:    time.Since(serverStartTime).Round(time.Second).String(),
		Checks:    checks,
		Version:   Version,
	}

	// Return 503 if degraded
	if status != "healthy" {
		NegotiateResponse(c, http.StatusServiceUnavailable, response)
		return
	}

	NegotiateResponse(c, http.StatusOK, response)
}

// handleMetrics returns the Prometheus metrics endpoint.
//
//	GET /metrics
func (s *Server) handleMetrics() gin.HandlerFunc {
	handler := promhttp.Handler()
	return func(c *gin.Context) {
		// Update system metrics before serving
		updateSystemMetrics()
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// handleVersion returns the API server version information.
//
//	GET /version
func (s *Server) handleVersion(c *gin.Context) {
	response := VersionResponse{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}

	NegotiateResponse(c, http.StatusOK, response)
}

// updateSystemMetrics updates Prometheus gauges for system metrics.
func updateSystemMetrics() {
	uptimeGauge.Set(time.Since(serverStartTime).Seconds())
	goroutineGauge.Set(float64(runtime.NumGoroutine()))

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryGauge.Set(float64(m.Alloc))
}
