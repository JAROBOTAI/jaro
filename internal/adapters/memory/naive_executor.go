package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/JAROBOTAI/jaro/internal/core/domain"
	"github.com/JAROBOTAI/jaro/internal/core/ports"
)

// NaiveExecutor is a simple mock implementation of the ports.Executor interface.
// It simulates step execution with a short sleep and always returns success.
// This allows testing the orchestrator flow without actual tool or LLM calls.
type NaiveExecutor struct{}

// NewNaiveExecutor creates a new mock executor for testing.
// Purpose: Factory function for creating the mock executor adapter.
// Inputs: None
// Outputs:
//   - ports.Executor: Initialized executor ready for use
func NewNaiveExecutor() ports.Executor {
	return &NaiveExecutor{}
}

// ExecuteStep simulates step execution with a 100ms delay and returns success.
// Purpose: Provides a predictable execution flow for testing without real tools or LLM.
//          Logs step execution to console and simulates processing time.
// Inputs:
//   - ctx: Context for cancellation and timeout control (unused in this implementation)
//   - task: The parent task (provides context information)
//   - step: The step to execute (used for logging and result generation)
// Outputs:
//   - *domain.StepResult: Always returns success with a mock output message
//   - error: Always returns nil (this implementation cannot fail)
func (e *NaiveExecutor) ExecuteStep(ctx context.Context, task *domain.Task, step *domain.Step) (*domain.StepResult, error) {
	// Log execution start
	fmt.Printf("[EXECUTOR] Executing Step: %s (Type: %s)...\n", step.Title, step.Type)

	// Simulate processing time
	startTime := time.Now()
	time.Sleep(100 * time.Millisecond)
	duration := time.Since(startTime)

	// Log execution completion
	fmt.Printf("[EXECUTOR] Step completed: %s (Duration: %dms)\n", step.Title, duration.Milliseconds())

	// Return success result
	result := &domain.StepResult{
		StepID:       step.ID,
		Success:      true,
		Output:       fmt.Sprintf("Step '%s' executed successfully (mock result)", step.Title),
		ErrorMessage: "",
		DurationMs:   duration.Milliseconds(),
	}

	return result, nil
}
