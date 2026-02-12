package ports

import (
	"context"

	"github.com/JAROBOTAI/jaro/internal/core/domain"
)

// TaskRepository provides persistence operations for tasks.
// This is a secondary port (infrastructure) that must be implemented by database adapters.
type TaskRepository interface {
	// SaveTask persists a task to the database (insert or update).
	// Purpose: Ensures task state is durably stored and can survive system restarts.
	// Inputs:
	//   - ctx: Context for cancellation and timeout control
	//   - task: The task to save (must have a valid ID)
	// Outputs:
	//   - error: Returns error if database is unavailable or task data is invalid
	SaveTask(ctx context.Context, task *domain.Task) error

	// GetTask retrieves a task by its unique identifier.
	// Purpose: Loads task state from persistent storage for status queries or resumption.
	// Inputs:
	//   - ctx: Context for cancellation and timeout control
	//   - id: Unique identifier of the task to retrieve
	// Outputs:
	//   - *domain.Task: The retrieved task with all fields populated
	//   - error: Returns error if task is not found or database is unavailable
	GetTask(ctx context.Context, id string) (*domain.Task, error)
}

// AuditRepository provides persistence operations for audit events.
// This is a secondary port for logging and compliance tracking.
type AuditRepository interface {
	// SaveEvent persists an audit event to the audit log.
	// Purpose: Creates immutable audit trail for compliance, debugging, and analytics.
	//          All significant system actions should generate audit events.
	// Inputs:
	//   - ctx: Context for cancellation and timeout control
	//   - event: The audit event to log (includes timestamp, actor, event type, payload)
	// Outputs:
	//   - error: Returns error if audit storage is unavailable (should not block task execution)
	SaveEvent(ctx context.Context, event *domain.AuditEvent) error
}

// LLMProvider is the interface to external Large Language Model services.
// This is a secondary port that abstracts LLM provider implementations (OpenAI, Anthropic, etc.).
type LLMProvider interface {
	// GenerateText sends a prompt to the LLM and returns the generated response.
	// Purpose: Provides AI reasoning capabilities for planning, decision-making, and text generation.
	//          Used by Planner for plan generation and Executor for THINK steps.
	// Inputs:
	//   - ctx: Context for cancellation and timeout control
	//   - prompt: The text prompt to send to the LLM (includes instructions and context)
	// Outputs:
	//   - string: The generated text response from the LLM
	//   - error: Returns error if LLM service is unavailable, rate-limited, or prompt is invalid
	GenerateText(ctx context.Context, prompt string) (string, error)
}
