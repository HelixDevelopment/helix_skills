package validation

import (
	"testing"

	"github.com/helixdevelopment/skill-system/internal/config"
	"go.uber.org/zap"
)

// TestChaos_NilStore_PipelineConstruction verifies that a Pipeline constructed
// with nil store does not panic.
func TestChaos_NilStore_PipelineConstruction(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("NewPipeline(nil,...) panicked: %v", r)
		}
	}()
	p := NewPipeline(nil, config.ValidationConfig{}, zap.NewNop())
	if p == nil {
		t.Error("expected non-nil pipeline")
	}
}

// TestChaos_ZeroValidationResult verifies that a zero-value ValidationResult
// has sensible defaults.
func TestChaos_ZeroValidationResult(t *testing.T) {
	result := &ValidationResult{}
	if result.Passed {
		t.Error("zero-value Passed should be false")
	}
	if result.Stage != "" {
		t.Error("zero-value Stage should be empty")
	}
}
