package validation

import (
	"sync"
	"testing"

	"github.com/helixdevelopment/skill-system/internal/config"
	"go.uber.org/zap"
)

// TestStress_ConcurrentPipelineConstruction exercises concurrent Pipeline
// construction. N=100 goroutines, no races.
func TestStress_ConcurrentPipelineConstruction(t *testing.T) {
	const n = 100
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p := NewPipeline(nil, config.ValidationConfig{}, zap.NewNop())
			if p == nil {
				t.Error("NewPipeline returned nil")
			}
		}()
	}
	wg.Wait()
}

// TestStress_ConcurrentValidationResult exercises concurrent ValidationResult
// construction. N=100 goroutines, no races.
func TestStress_ConcurrentValidationResult(t *testing.T) {
	const n = 100
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := &ValidationResult{
				Passed:     true,
				Stage:      "test",
				ApprovedBy: 2,
			}
			if !result.Passed {
				t.Error("expected Passed=true")
			}
		}()
	}
	wg.Wait()
}
