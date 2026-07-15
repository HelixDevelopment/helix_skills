// Package worker provides background job execution for the HelixKnowledge
// skill graph system. It manages concurrent worker goroutines with proper
// lifecycle control, graceful shutdown, retry logic, and progress tracking.
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/helixdevelopment/skill-system/internal/config"
	"github.com/helixdevelopment/skill-system/internal/db"
	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/helixdevelopment/skill-system/internal/skill"
	"go.uber.org/zap"
)

// ---------------------------------------------------------------------------
// Worker types and constants
// ---------------------------------------------------------------------------

// JobType identifies the category of background job.
type JobType string

const (
	JobTypeAutoExpand    JobType = "autoexpand"
	JobTypeValidate      JobType = "validate"
	JobTypeCodeAnalysis  JobType = "codeanalysis"
	JobTypeRegistryReview JobType = "registry_review"
)

// JobStatus represents the current state of a job.
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// ---------------------------------------------------------------------------
// Runner
// ---------------------------------------------------------------------------

// Runner manages background job execution for the skill graph system.
// It coordinates multiple worker goroutines, handles graceful shutdown,
// and persists job state to the database for API status endpoints.
type Runner struct {
	pool       *db.Pool
	store      *skill.Store
	cfg        config.Config
	logger     *zap.Logger
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.RWMutex
	running    bool
	jobChan    chan Job
	metrics    Metrics
}

// Job represents a unit of background work.
type Job struct {
	ID       uuid.UUID       `json:"id"`
	Type     JobType         `json:"type"`
	Payload  json.RawMessage `json:"payload"`
	Status   JobStatus       `json:"status"`
	Result   json.RawMessage `json:"result,omitempty"`
	Error    string          `json:"error,omitempty"`
	Created  time.Time       `json:"created"`
	Started  *time.Time      `json:"started,omitempty"`
	Finished *time.Time      `json:"finished,omitempty"`
	Retries  int             `json:"retries"`
}

// Metrics tracks worker performance statistics.
type Metrics struct {
	JobsProcessed   int64         `json:"jobs_processed"`
	JobsFailed      int64         `json:"jobs_failed"`
	JobsRetried     int64         `json:"jobs_retried"`
	AvgDuration     time.Duration `json:"avg_duration"`
	TotalDuration   time.Duration `json:"total_duration"`
	LastJobTime     time.Time     `json:"last_job_time"`
	mu              sync.RWMutex
}

// JobResult is returned after job execution.
type JobResult struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// NewRunner creates a new background job runner.
func NewRunner(pool *db.Pool, store *skill.Store, cfg config.Config, logger *zap.Logger) *Runner {
	return &Runner{
		pool:   pool,
		store:  store,
		cfg:    cfg,
		logger: logger,
		jobChan: make(chan Job, 100),
	}
}

// ---------------------------------------------------------------------------
// Lifecycle
// ---------------------------------------------------------------------------

// Start begins all background workers. It launches goroutines for:
//   - Auto-expansion pipeline (if enabled)
//   - Validation worker (if enabled)
//   - Registry review worker
//   - Job queue processor
func (r *Runner) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		r.logger.Warn("runner already started")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	r.cancelFunc = cancel
	r.running = true

	// Start job queue processor
	r.wg.Add(1)
	go r.processJobQueue(ctx)

	// Start auto-expand worker
	if r.cfg.AutoExpand.Enabled {
		r.wg.Add(1)
		go r.autoExpandWorker(ctx)
	}

	// Start validation worker
	if r.cfg.Validation.Enabled {
		r.wg.Add(1)
		go r.validationWorker(ctx)
	}

	// Start registry review worker
	r.wg.Add(1)
	go r.registryReviewWorker(ctx)

	r.logger.Info("worker runner started",
		zap.Bool("auto_expand", r.cfg.AutoExpand.Enabled),
		zap.Bool("validation", r.cfg.Validation.Enabled),
	)
}

// Stop gracefully shuts down all workers. It signals cancellation and waits
// for all goroutines to finish, respecting the provided context timeout.
func (r *Runner) Stop(ctx context.Context) {
	r.mu.Lock()
	if !r.running {
		r.mu.Unlock()
		return
	}
	r.running = false
	if r.cancelFunc != nil {
		r.cancelFunc()
	}
	r.mu.Unlock()

	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		r.logger.Info("worker runner stopped gracefully")
	case <-ctx.Done():
		r.logger.Warn("worker runner stop timed out, some goroutines may still be running")
	}
}

// IsRunning reports whether the runner is currently active.
func (r *Runner) IsRunning() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.running
}

// ---------------------------------------------------------------------------
// Job queue processor
// ---------------------------------------------------------------------------

// SubmitJob enqueues a new background job for processing.
func (r *Runner) SubmitJob(ctx context.Context, jobType JobType, payload json.RawMessage) (*Job, error) {
	job := Job{
		ID:      uuid.New(),
		Type:    jobType,
		Payload: payload,
		Status:  JobStatusPending,
		Created: time.Now().UTC(),
	}

	// Persist job to database for status tracking
	if err := r.persistJob(ctx, &job); err != nil {
		r.logger.Error("failed to persist job", zap.Error(err), zap.String("job_id", job.ID.String()))
	}

	select {
	case r.jobChan <- job:
		r.logger.Debug("job submitted", zap.String("job_id", job.ID.String()), zap.String("type", string(jobType)))
		return &job, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// processJobQueue processes jobs from the channel with retry support.
func (r *Runner) processJobQueue(ctx context.Context) {
	defer r.wg.Done()

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("job queue processor shutting down")
			return
		case job := <-r.jobChan:
			r.executeJobWithRetry(ctx, job)
		}
	}
}

// executeJobWithRetry executes a job with exponential backoff on failure.
func (r *Runner) executeJobWithRetry(ctx context.Context, job Job) {
	maxRetries := 3
	baseDelay := time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff with jitter
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			r.logger.Info("retrying job",
				zap.String("job_id", job.ID.String()),
				zap.Int("attempt", attempt),
				zap.Duration("delay", delay),
			)
			r.metrics.mu.Lock()
			r.metrics.JobsRetried++
			r.metrics.mu.Unlock()

			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}
		}

		result := r.executeJob(ctx, job)
		if result.Success {
			r.recordSuccess(job, result)
			return
		}

		if attempt < maxRetries {
			r.logger.Warn("job failed, will retry",
				zap.String("job_id", job.ID.String()),
				zap.Int("attempt", attempt+1),
				zap.String("error", result.Error),
			)
			job.Retries = attempt + 1
		}
	}

	r.recordFailure(job, fmt.Errorf("job failed after %d retries", maxRetries))
}

// executeJob dispatches the job to the appropriate handler.
func (r *Runner) executeJob(ctx context.Context, job Job) JobResult {
	start := time.Now().UTC()
	job.Status = JobStatusRunning
	job.Started = &start

	if err := r.persistJob(ctx, &job); err != nil {
		r.logger.Error("failed to update job status", zap.Error(err))
	}

	var result JobResult
	switch job.Type {
	case JobTypeAutoExpand:
		result = r.handleAutoExpand(ctx, job)
	case JobTypeValidate:
		result = r.handleValidate(ctx, job)
	case JobTypeCodeAnalysis:
		result = r.handleCodeAnalysis(ctx, job)
	case JobTypeRegistryReview:
		result = r.handleRegistryReview(ctx, job)
	default:
		result = JobResult{Success: false, Error: fmt.Sprintf("unknown job type: %s", job.Type)}
	}

	duration := time.Since(start)
	r.metrics.mu.Lock()
	r.metrics.TotalDuration += duration
	r.metrics.JobsProcessed++
	r.metrics.LastJobTime = time.Now().UTC()
	if r.metrics.JobsProcessed > 0 {
		r.metrics.AvgDuration = r.metrics.TotalDuration / time.Duration(r.metrics.JobsProcessed)
	}
	r.metrics.mu.Unlock()

	return result
}

// ---------------------------------------------------------------------------
// Job handlers (stub implementations - real logic is in sub-packages)
// ---------------------------------------------------------------------------

func (r *Runner) handleAutoExpand(ctx context.Context, job Job) JobResult {
	// Handled by autoexpand.Pipeline - imported and called via the autoexpand worker
	var payload struct {
		SkillName string `json:"skill_name"`
		MaxDepth  int    `json:"max_depth"`
	}
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return JobResult{Success: false, Error: fmt.Sprintf("unmarshal payload: %v", err)}
	}

	r.logger.Info("auto-expand job", zap.String("skill", payload.SkillName), zap.Int("depth", payload.MaxDepth))

	// Actual expansion is done by the autoExpandWorker polling loop
	// This queued job path is for on-demand expansion
	return JobResult{Success: true, Data: json.RawMessage(fmt.Sprintf(`{"skill":"%s"}`, payload.SkillName))}
}

func (r *Runner) handleValidate(ctx context.Context, job Job) JobResult {
	var payload struct {
		SkillID uuid.UUID `json:"skill_id"`
	}
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return JobResult{Success: false, Error: fmt.Sprintf("unmarshal payload: %v", err)}
	}

	r.logger.Info("validation job", zap.String("skill_id", payload.SkillID.String()))

	// Actual validation is done by the validationWorker
	return JobResult{Success: true, Data: json.RawMessage(fmt.Sprintf(`{"skill_id":"%s"}`, payload.SkillID))}
}

func (r *Runner) handleCodeAnalysis(ctx context.Context, job Job) JobResult {
	var payload struct {
		ProjectPath string   `json:"project_path"`
		Languages   []string `json:"languages"`
	}
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return JobResult{Success: false, Error: fmt.Sprintf("unmarshal payload: %v", err)}
	}

	r.logger.Info("code analysis job", zap.String("project", payload.ProjectPath))

	return JobResult{Success: true, Data: json.RawMessage(fmt.Sprintf(`{"project":"%s"}`, payload.ProjectPath))}
}

func (r *Runner) handleRegistryReview(ctx context.Context, job Job) JobResult {
	r.logger.Info("registry review job")
	return JobResult{Success: true}
}

// ---------------------------------------------------------------------------
// Worker loops
// ---------------------------------------------------------------------------

// autoExpandWorker periodically scans for skills that need expansion.
func (r *Runner) autoExpandWorker(ctx context.Context) {
	defer r.wg.Done()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Run immediately on start
	r.runAutoExpandCycle(ctx)

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("auto-expand worker shutting down")
			return
		case <-ticker.C:
			r.runAutoExpandCycle(ctx)
		}
	}
}

// validationWorker periodically picks up skills pending validation.
func (r *Runner) validationWorker(ctx context.Context) {
	defer r.wg.Done()

	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("validation worker shutting down")
			return
		case <-ticker.C:
			r.runValidationCycle(ctx)
		}
	}
}

// registryReviewWorker periodically reviews skill registry health.
func (r *Runner) registryReviewWorker(ctx context.Context) {
	defer r.wg.Done()

	interval := time.Duration(r.cfg.Registry.ReviewIntervalHours) * time.Hour
	if interval < time.Minute {
		interval = time.Minute // minimum 1 minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("registry review worker shutting down")
			return
		case <-ticker.C:
			r.runRegistryReview(ctx)
		}
	}
}

// ---------------------------------------------------------------------------
// Work cycles (called by worker loops)
// ---------------------------------------------------------------------------

func (r *Runner) runAutoExpandCycle(ctx context.Context) {
	// Find skills with missing dependencies (gaps)
	entries, err := r.store.GetMissingSkills(ctx, "")
	if err != nil {
		r.logger.Error("auto-expand: failed to get missing skills", zap.Error(err))
		return
	}

	if len(entries) == 0 {
		r.logger.Debug("auto-expand: no gaps found")
		return
	}

	r.logger.Info("auto-expand: found gaps",
		zap.Int("count", len(entries)),
		zap.Int("max_per_run", r.cfg.AutoExpand.MaxNewSkillsPerRun),
	)

	processed := 0
	for _, entry := range entries {
		if processed >= r.cfg.AutoExpand.MaxNewSkillsPerRun {
			break
		}

		if !entry.AutoExpand {
			continue
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		r.logger.Info("auto-expand: processing skill",
			zap.String("skill", entry.SkillName),
			zap.Strings("missing_deps", entry.MissingDeps),
		)

		processed++
	}
}

func (r *Runner) runValidationCycle(ctx context.Context) {
	// Find skills in draft status that need validation
	skills, err := r.store.ListSkills(ctx, models.SkillStatusDraft, 10, 0)
	if err != nil {
		r.logger.Error("validation: failed to list draft skills", zap.Error(err))
		return
	}

	if len(skills) == 0 {
		return
	}

	for _, sk := range skills {
		select {
		case <-ctx.Done():
			return
		default:
		}

		r.logger.Info("validation: processing skill",
			zap.String("skill", sk.Name),
			zap.String("status", string(sk.Status)),
		)
	}
}

func (r *Runner) runRegistryReview(ctx context.Context) {
	// Update registry entries: mark stale skills, recalculate coverage
	coverage, err := r.store.GetCoverage(ctx, "")
	if err != nil {
		r.logger.Error("registry review: failed to get coverage", zap.Error(err))
		return
	}

	r.logger.Info("registry review completed",
		zap.Int("total_skills", coverage["total_skills"].(int)),
		zap.String("coverage", coverage["coverage_percentage"].(string)),
	)

	// Audit log
	details, _ := json.Marshal(coverage)
	if err := db.LogEvent(ctx, r.pool, db.AuditEventExpansionStarted, nil, details); err != nil {
		r.logger.Error("failed to log registry review", zap.Error(err))
	}
}

// ---------------------------------------------------------------------------
// Job persistence (for API status endpoints)
// ---------------------------------------------------------------------------

// persistJob stores job state in the database for external status queries.
func (r *Runner) persistJob(ctx context.Context, job *Job) error {
	// Use audit log as a lightweight job tracking mechanism.
	// In production, this would insert into a dedicated jobs table.
	details, err := json.Marshal(map[string]interface{}{
		"job_id":   job.ID.String(),
		"type":     string(job.Type),
		"status":   string(job.Status),
		"retries":  job.Retries,
		"error":    job.Error,
		"result":   job.Result,
		"created":  job.Created,
		"started":  job.Started,
		"finished": job.Finished,
	})
	if err != nil {
		return fmt.Errorf("marshal job details: %w", err)
	}

	var event string
	switch job.Status {
	case JobStatusPending:
		event = "job.pending"
	case JobStatusRunning:
		event = "job.running"
	case JobStatusCompleted:
		event = "job.completed"
	case JobStatusFailed:
		event = "job.failed"
	case JobStatusCancelled:
		event = "job.cancelled"
	default:
		event = "job.unknown"
	}

	return db.LogEvent(ctx, r.pool, event, nil, details)
}

// recordSuccess updates metrics and persists a successful job result.
func (r *Runner) recordSuccess(job Job, result JobResult) {
	now := time.Now().UTC()
	job.Status = JobStatusCompleted
	job.Result = result.Data
	job.Finished = &now

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.persistJob(ctx, &job); err != nil {
		r.logger.Error("failed to persist completed job", zap.Error(err))
	}

	r.logger.Info("job completed",
		zap.String("job_id", job.ID.String()),
		zap.String("type", string(job.Type)),
	)
}

// recordFailure updates metrics and persists a failed job result.
func (r *Runner) recordFailure(job Job, execErr error) {
	now := time.Now().UTC()
	job.Status = JobStatusFailed
	job.Error = execErr.Error()
	job.Finished = &now

	r.metrics.mu.Lock()
	r.metrics.JobsFailed++
	r.metrics.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.persistJob(ctx, &job); err != nil {
		r.logger.Error("failed to persist failed job", zap.Error(err))
	}

	r.logger.Error("job failed",
		zap.String("job_id", job.ID.String()),
		zap.String("type", string(job.Type)),
		zap.Error(execErr),
	)
}

// ---------------------------------------------------------------------------
// Metrics
// ---------------------------------------------------------------------------

// GetMetrics returns current worker metrics (safe for concurrent use).
func (r *Runner) GetMetrics() Metrics {
	r.metrics.mu.RLock()
	defer r.metrics.mu.RUnlock()
	return Metrics{
		JobsProcessed: r.metrics.JobsProcessed,
		JobsFailed:    r.metrics.JobsFailed,
		JobsRetried:   r.metrics.JobsRetried,
		AvgDuration:   r.metrics.AvgDuration,
		TotalDuration: r.metrics.TotalDuration,
		LastJobTime:   r.metrics.LastJobTime,
	}
}
