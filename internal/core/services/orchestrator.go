package services

import (
	"context"
	"fmt"

	"github.com/JAROBOTAI/jaro/internal/core/domain"
	"github.com/JAROBOTAI/jaro/internal/core/ports"
)

// OrchestratorService is the core implementation of the Orchestrator interface.
// It coordinates task lifecycle management, plan execution, and approval workflows.
// This service follows the Hexagonal Architecture pattern by depending only on ports (interfaces).
type OrchestratorService struct {
	planner   ports.Planner
	executor  ports.Executor
	repo      ports.TaskRepository
	audit     ports.AuditRepository
	clock     ports.Clock
	idGen     ports.IDGenerator
	logger    ports.Logger
}

// NewOrchestrator creates a new OrchestratorService instance with the required dependencies.
// Purpose: Factory function for creating the orchestrator service with dependency injection.
// Inputs:
//   - planner: Implementation of the Planner port for generating execution plans
//   - executor: Implementation of the Executor port for running plan steps
//   - repo: Implementation of the TaskRepository port for task persistence
//   - audit: Implementation of the AuditRepository port for audit logging
//   - clock: Implementation of the Clock port for time operations
//   - idGen: Implementation of the IDGenerator port for ID generation
//   - logger: Implementation of the Logger port for structured logging
// Outputs:
//   - ports.Orchestrator: Fully initialized orchestrator service ready for use
func NewOrchestrator(
	planner ports.Planner,
	executor ports.Executor,
	repo ports.TaskRepository,
	audit ports.AuditRepository,
	clock ports.Clock,
	idGen ports.IDGenerator,
	logger ports.Logger,
) ports.Orchestrator {
	return &OrchestratorService{
		planner:  planner,
		executor: executor,
		repo:     repo,
		audit:    audit,
		clock:    clock,
		idGen:    idGen,
		logger:   logger,
	}
}

// StartTask initializes a new task based on user input and creates an execution plan.
// Purpose: This is the primary entry point for submitting work to the JARO system.
//          Creates a new task, persists it, and logs the creation event.
// Inputs:
//   - ctx: Context for cancellation and timeout control
//   - input: Raw user request in natural language
//   - userID: Unique identifier of the user submitting the task
// Outputs:
//   - *domain.Task: The created task with status NEW
//   - error: Returns error if input validation fails, persistence fails, or system is unavailable
func (s *OrchestratorService) StartTask(ctx context.Context, input string, userID string) (*domain.Task, error) {
	// Validate input
	if input == "" {
		return nil, fmt.Errorf("task input cannot be empty")
	}
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	// Create new task with initial state
	now := s.clock.Now()
	taskID := s.idGen.Generate()
	
	task := &domain.Task{
		ID:               taskID,
		CreatedAt:        now,
		UpdatedAt:        now,
		Status:           domain.TaskStatusNew,
		Input:            input,
		NormalizedIntent: input, // Initial value; will be refined by planner
		UserID:           userID,
		Channel:          "api", // Default channel
		TargetAgent:      "CORE", // Default agent for V1
		Artifacts:        make(map[string]string),
		Metadata:         make(map[string]string),
		UsageTokens:      0,
		CostEstimate:     0.0,
	}

	// Persist the task
	if err := s.repo.SaveTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	// Create audit event for task creation
	auditEvent := &domain.AuditEvent{
		ID:              s.idGen.Generate(),
		TaskID:          taskID,
		CorrelationID:   taskID,
		Timestamp:       now,
		EventType:       "TASK_CREATED",
		Actor:           userID,
		BehaviorVersion: "v1",
		Payload: map[string]interface{}{
			"input":        input,
			"task_id":      taskID,
			"user_id":      userID,
			"target_agent": task.TargetAgent,
		},
	}

	// Log the audit event (non-blocking - don't fail task creation if audit fails)
	if err := s.audit.SaveEvent(ctx, auditEvent); err != nil {
		// Log warning using logger port instead of fmt.Printf
		s.logger.Warn("failed to save audit event", map[string]interface{}{
			"error":   err.Error(),
			"task_id": taskID,
		})
	}

	return task, nil
}

// GetTaskStatus retrieves the current state and progress of a task.
// Purpose: Allows clients to poll for task status and results.
//          Simple pass-through to the repository layer.
// Inputs:
//   - ctx: Context for cancellation and timeout control
//   - taskID: Unique identifier of the task to query
// Outputs:
//   - *domain.Task: Current task state including status, steps, and artifacts
//   - error: Returns error if task is not found or access is denied
func (s *OrchestratorService) GetTaskStatus(ctx context.Context, taskID string) (*domain.Task, error) {
	if taskID == "" {
		return nil, fmt.Errorf("taskID cannot be empty")
	}

	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task status: %w", err)
	}

	return task, nil
}

// HandleApproval processes user approval or rejection for high-risk steps.
// Purpose: Implements the human-in-the-loop pattern for risky operations.
//          Updates task status based on user decision and logs the approval decision.
// Inputs:
//   - ctx: Context for cancellation and timeout control
//   - taskID: Unique identifier of the task awaiting approval
//   - stepID: Unique identifier of the step requiring approval
//   - approved: User decision (true = approve, false = reject)
//   - userID: Unique identifier of the user making the decision
// Outputs:
//   - error: Returns error if task/step not found, not awaiting approval, or persistence fails
func (s *OrchestratorService) HandleApproval(ctx context.Context, taskID string, stepID string, approved bool, userID string) error {
	// Validate inputs
	if taskID == "" {
		return fmt.Errorf("taskID cannot be empty")
	}
	if stepID == "" {
		return fmt.Errorf("stepID cannot be empty")
	}
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	// Load the task
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to load task: %w", err)
	}

	// Verify task is waiting for approval
	if task.Status != domain.TaskStatusWaitingApproval {
		return fmt.Errorf("task %s is not waiting for approval (current status: %s)", taskID, task.Status)
	}

	// Verify the current step matches
	if task.CurrentStepID != stepID {
		return fmt.Errorf("step mismatch: expected %s but got %s", task.CurrentStepID, stepID)
	}

	// Update task status based on approval decision
	now := s.clock.Now()
	if approved {
		task.Status = domain.TaskStatusExecuting
	} else {
		task.Status = domain.TaskStatusCanceled
		task.FinishedAt = now
	}
	task.UpdatedAt = now

	// Persist updated task
	if err := s.repo.SaveTask(ctx, task); err != nil {
		return fmt.Errorf("failed to save task after approval: %w", err)
	}

	// Create audit event for approval decision
	auditEvent := &domain.AuditEvent{
		ID:              s.idGen.Generate(),
		TaskID:          taskID,
		CorrelationID:   taskID,
		Timestamp:       now,
		EventType:       "APPROVAL_DECISION",
		Actor:           userID,
		BehaviorVersion: "v1",
		Payload: map[string]interface{}{
			"task_id":    taskID,
			"step_id":    stepID,
			"approved":   approved,
			"user_id":    userID,
			"new_status": string(task.Status),
		},
	}

	// Log the audit event (non-blocking)
	if err := s.audit.SaveEvent(ctx, auditEvent); err != nil {
		s.logger.Warn("failed to save audit event", map[string]interface{}{
			"error":   err.Error(),
			"task_id": taskID,
			"step_id": stepID,
		})
	}

	return nil
}
