// Package registry provides scheduled review jobs for skill health monitoring.
package registry

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// ReviewScheduler manages periodic skill review jobs.
type ReviewScheduler struct {
	registry *Registry
	cron     *cron.Cron
	running  bool
	mu       sync.Mutex
	cancel   context.CancelFunc
	ctx      context.Context
	wg       sync.WaitGroup
}

// StartReviewScheduler starts a background cron job that periodically checks
// all skills and marks stale ones. The interval controls how often the review runs.
// Common intervals: 1h (frequent), 24h (daily), 168h (weekly).
func (r *Registry) StartReviewScheduler(interval time.Duration) *ReviewScheduler {
	if interval < time.Minute {
		interval = time.Hour // Minimum 1 hour to avoid excessive DB load
	}

	ctx, cancel := context.WithCancel(context.Background())
	scheduler := &ReviewScheduler{
		registry: r,
		ctx:      ctx,
		cancel:   cancel,
	}

	// Use robfig/cron for more flexible scheduling
	// Map duration to a cron expression (run every N hours)
	hours := int(interval.Hours())
	if hours < 1 {
		hours = 1
	}

	cronExpr := fmt.Sprintf("0 */%d * * *", hours)
	if hours >= 24 {
		days := hours / 24
		if days == 1 {
			cronExpr = "0 0 * * *" // Daily at midnight
		} else {
			cronExpr = fmt.Sprintf("0 0 */%d * *", days)
		}
	}

	scheduler.cron = cron.New()
	_, err := scheduler.cron.AddFunc(cronExpr, func() {
		scheduler.performReview(ctx)
	})
	if err != nil {
		log.Printf("[ReviewScheduler] Failed to add cron job: %v", err)
		// Fall back to simple ticker-based scheduling
		scheduler.startTickerFallback(interval)
		return scheduler
	}

	scheduler.cron.Start()
	scheduler.running = true

	log.Printf("[ReviewScheduler] Started with cron expression: %s (interval: %v)", cronExpr, interval)

	// Run an immediate review on startup
	scheduler.wg.Add(1)
	go func() {
		defer scheduler.wg.Done()
		scheduler.performReview(ctx)
	}()

	return scheduler
}

// startTickerFallback uses a time.Ticker when cron scheduling fails.
func (rs *ReviewScheduler) startTickerFallback(interval time.Duration) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.running = true
	ticker := time.NewTicker(interval)

	rs.wg.Add(1)
	go func() {
		defer rs.wg.Done()
		defer ticker.Stop()
		for {
			select {
			case <-rs.ctx.Done():
				log.Println("[ReviewScheduler] Ticker fallback stopped")
				return
			case <-ticker.C:
				rs.performReview(rs.ctx)
			}
		}
	}()

	// Run immediate review
	rs.wg.Add(1)
	go func() {
		defer rs.wg.Done()
		rs.performReview(rs.ctx)
	}()

	log.Printf("[ReviewScheduler] Fallback ticker started with interval: %v", interval)
}

// performReview checks all skills and marks those needing attention as stale.
// It performs the following checks:
// - Skills not reviewed in the last 30 days
// - Skills with missing dependencies
// - Skills with coverage below 25%
// - Skills with no resources or evidence
func (rs *ReviewScheduler) performReview(ctx context.Context) {
	start := time.Now()
	log.Println("[ReviewScheduler] Starting periodic review...")

	// Check 1: Mark skills not reviewed in 30+ days as stale
	if err := rs.markOldSkillsStale(ctx); err != nil {
		log.Printf("[ReviewScheduler] Error marking old skills stale: %v", err)
	}

	// Check 2: Mark skills with missing dependencies as stale
	if err := rs.markMissingDepSkillsStale(ctx); err != nil {
		log.Printf("[ReviewScheduler] Error marking missing-dep skills stale: %v", err)
	}

	// Check 3: Mark skills with very low coverage as stale
	if err := rs.markLowCoverageSkillsStale(ctx); err != nil {
		log.Printf("[ReviewScheduler] Error marking low-coverage skills stale: %v", err)
	}

	// Check 4: Update coverage scores for all skills
	if err := rs.registry.UpdateCoverage(ctx); err != nil {
		log.Printf("[ReviewScheduler] Error updating coverage: %v", err)
	}

	// Check 5: Recalculate missing dependencies
	if err := rs.registry.CalculateMissingDeps(ctx); err != nil {
		log.Printf("[ReviewScheduler] Error calculating missing deps: %v", err)
	}

	duration := time.Since(start)
	log.Printf("[ReviewScheduler] Review completed in %v", duration)
}

// markOldSkillsStale marks skills that haven't been reviewed in 30+ days as stale.
func (rs *ReviewScheduler) markOldSkillsStale(ctx context.Context) error {
	tag, err := rs.registry.pool.Inner().Exec(ctx, `
		UPDATE skill_registry sr
		SET stale = true
		WHERE sr.last_review < NOW() - INTERVAL '30 days'
		  AND EXISTS (
			  SELECT 1 FROM skills s
			  WHERE s.id = sr.skill_id
			    AND s.status IN ('validated', 'active')
		  )
		  AND sr.stale = false
	`)
	if err != nil {
		return fmt.Errorf("mark old skills stale: %w", err)
	}
	if tag.RowsAffected() > 0 {
		log.Printf("[ReviewScheduler] Marked %d skills stale (not reviewed in 30+ days)", tag.RowsAffected())
	}
	return nil
}

// markMissingDepSkillsStale marks skills with unresolved dependencies as stale.
func (rs *ReviewScheduler) markMissingDepSkillsStale(ctx context.Context) error {
	tag, err := rs.registry.pool.Inner().Exec(ctx, `
		UPDATE skill_registry sr
		SET stale = true
		WHERE sr.missing_deps != '{}'
		  AND sr.stale = false
	`)
	if err != nil {
		return fmt.Errorf("mark missing-dep skills stale: %w", err)
	}
	if tag.RowsAffected() > 0 {
		log.Printf("[ReviewScheduler] Marked %d skills stale (missing dependencies)", tag.RowsAffected())
	}
	return nil
}

// markLowCoverageSkillsStale marks skills with coverage below 25% as stale.
func (rs *ReviewScheduler) markLowCoverageSkillsStale(ctx context.Context) error {
	tag, err := rs.registry.pool.Inner().Exec(ctx, `
		UPDATE skill_registry sr
		SET stale = true
		WHERE sr.coverage < 0.25
		  AND EXISTS (
			  SELECT 1 FROM skills s
			  WHERE s.id = sr.skill_id
			    AND s.status IN ('validated', 'active')
		  )
		  AND sr.stale = false
	`)
	if err != nil {
		return fmt.Errorf("mark low-coverage skills stale: %w", err)
	}
	if tag.RowsAffected() > 0 {
		log.Printf("[ReviewScheduler] Marked %d skills stale (coverage < 25%%)", tag.RowsAffected())
	}
	return nil
}

// Stop gracefully shuts down the review scheduler.
func (rs *ReviewScheduler) Stop() {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if !rs.running {
		return
	}

	log.Println("[ReviewScheduler] Stopping...")

	if rs.cancel != nil {
		rs.cancel()
	}

	if rs.cron != nil {
		stopCtx := rs.cron.Stop()
		<-stopCtx.Done()
	}

	// Wait for background goroutines to finish (with timeout)
	done := make(chan struct{})
	go func() {
		rs.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[ReviewScheduler] Stopped cleanly")
	case <-time.After(10 * time.Second):
		log.Println("[ReviewScheduler] Stop timed out, forcing shutdown")
	}

	rs.running = false
}

// IsRunning returns whether the scheduler is currently active.
func (rs *ReviewScheduler) IsRunning() bool {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.running
}

// ---------------------------------------------------------------------------
// Convenience constructors
// ---------------------------------------------------------------------------

// NewDailyReviewScheduler starts a scheduler that runs reviews once per day.
func (r *Registry) NewDailyReviewScheduler() *ReviewScheduler {
	return r.StartReviewScheduler(24 * time.Hour)
}

// NewHourlyReviewScheduler starts a scheduler that runs reviews every hour.
func (r *Registry) NewHourlyReviewScheduler() *ReviewScheduler {
	return r.StartReviewScheduler(time.Hour)
}

// RunReviewOnce performs a single review cycle synchronously.
func (r *Registry) RunReviewOnce(ctx context.Context) error {
	log.Println("[ReviewScheduler] Running one-off review...")
	start := time.Now()

	// Update coverage
	if err := r.UpdateCoverage(ctx); err != nil {
		return fmt.Errorf("update coverage: %w", err)
	}

	// Calculate missing deps
	if err := r.CalculateMissingDeps(ctx); err != nil {
		return fmt.Errorf("calculate missing deps: %w", err)
	}

	// Mark stale skills based on age
	if err := r.pool.Exec(ctx, `
		UPDATE skill_registry sr
		SET stale = true
		WHERE (sr.last_review < NOW() - INTERVAL '30 days' OR sr.last_review IS NULL)
		  AND EXISTS (
			  SELECT 1 FROM skills s
			  WHERE s.id = sr.skill_id
			    AND s.status IN ('validated', 'active')
		  )
		  AND sr.stale = false
	`); err != nil {
		return fmt.Errorf("mark old skills stale: %w", err)
	}

	// Mark skills with missing deps as stale
	if err := r.pool.Exec(ctx, `
		UPDATE skill_registry
		SET stale = true
		WHERE missing_deps != '{}' AND stale = false
	`); err != nil {
		return fmt.Errorf("mark missing-dep skills stale: %w", err)
	}

	log.Printf("[ReviewScheduler] One-off review completed in %v", time.Since(start))
	return nil
}
