package ports

import (
	"context"

	"github.com/JAROBOTAI/jaro/internal/core/domain"
)

// Planner is responsible for generating execution plans from task requirements.
// It uses LLM capabilities to decompose complex tasks into executable steps.
type Planner interface {
	// CreatePlan generates a structured execution plan for the given task.
	// Purpose: Converts natural language task input into a sequence of actionable steps.
	//          Uses LLM to reason about task decomposition and risk assessment.
	// Inputs:
	//   - ctx: Context for cancellation and timeout control
	//   - task: The task requiring a plan (contains user input and normalized intent)
	//   - tools: Available tools that can be used in the plan steps
	// Outputs:
	//   - *domain.Plan: Generated plan with steps, dependencies, and risk summary
	//   - error: Returns error if LLM fails, task is malformed, or no viable plan exists
	CreatePlan(ctx context.Context, task *domain.Task, tools []domain.ToolMetadata) (*domain.Plan, error)
}

// Executor is responsible for executing individual plan steps.
// It handles tool invocations, LLM calls, and step state management.
type Executor interface {
	// ExecuteStep runs a single step from a plan and returns the result.
	// Purpose: Executes one atomic unit of work (tool call, LLM reasoning, decision).
	//          Handles retries, error capture, and duration tracking.
	// Inputs:
	//   - ctx: Context for cancellation and timeout control
	//   - task: The parent task (provides context and metadata)
	//   - step: The step to execute (contains type, tool name, input, etc.)
	// Outputs:
	//   - *domain.StepResult: Execution result including success status, output, and metrics
	//   - error: Returns error if step execution fails critically or context is cancelled
	ExecuteStep(ctx context.Context, task *domain.Task, step *domain.Step) (*domain.StepResult, error)
}

// ToolRegistry manages the collection of available tools in the system.
// It provides discovery and access to registered tools.
type ToolRegistry interface {
	// GetTool retrieves a tool by its unique name.
	// Purpose: Provides access to a specific tool for execution.
	// Inputs:
	//   - name: Unique name of the tool (e.g., "file_write", "http_request")
	// Outputs:
	//   - domain.Tool: The tool instance ready for execution
	//   - error: Returns error if tool is not found or not available
	GetTool(name string) (domain.Tool, error)

	// ListTools returns metadata for all registered tools.
	// Purpose: Enables tool discovery for planning and UI purposes.
	// Inputs: None
	// Outputs:
	//   - []domain.ToolMetadata: List of all available tools with descriptions and risk levels
	ListTools() []domain.ToolMetadata
}
