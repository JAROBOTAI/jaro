package ports

import (
	"context"

	"github.com/JAROBOTAI/jaro/internal/core/domain"
)

// Orchestrator is the primary port (API) for interacting with the JARO system.
// It serves as the main entry point for task management and orchestration operations.
// All external clients (HTTP handlers, CLI, gRPC) should interact through this interface.
type Orchestrator interface {
	// StartTask initializes a new task based on user input and creates an execution plan.
	// Purpose: This is the primary entry point for submitting work to the JARO system.
	// Inputs:
	//   - ctx: Context for cancellation and timeout control
	//   - input: Raw user request in natural language
	//   - userID: Unique identifier of the user submitting the task
	// Outputs:
	//   - *domain.Task: The created task with status NEW or PLANNING
	//   - error: Returns error if input validation fails or system is unavailable
	StartTask(ctx context.Context, input string, userID string) (*domain.Task, error)

	// GetTaskStatus retrieves the current state and progress of a task.
	// Purpose: Allows clients to poll for task status and results.
	// Inputs:
	//   - ctx: Context for cancellation and timeout control
	//   - taskID: Unique identifier of the task to query
	// Outputs:
	//   - *domain.Task: Current task state including status, steps, and artifacts
	//   - error: Returns error if task is not found or access is denied
	GetTaskStatus(ctx context.Context, taskID string) (*domain.Task, error)

	// HandleApproval processes user approval or rejection for high-risk steps.
	// Purpose: Implements the human-in-the-loop pattern for risky operations.
	// Inputs:
	//   - ctx: Context for cancellation and timeout control
	//   - taskID: Unique identifier of the task awaiting approval
	//   - stepID: Unique identifier of the step requiring approval
	//   - approved: User decision (true = approve, false = reject)
	//   - userID: Unique identifier of the user making the decision
	// Outputs:
	//   - error: Returns error if task/step not found, already processed, or unauthorized
	HandleApproval(ctx context.Context, taskID string, stepID string, approved bool, userID string) error
}
