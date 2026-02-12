package memory

import (
	"context"

	"github.com/JAROBOTAI/jaro/internal/core/domain"
	"github.com/JAROBOTAI/jaro/internal/core/ports"
	"github.com/google/uuid"
)

// NaivePlanner is a simple hardcoded implementation of the ports.Planner interface.
// It generates a fixed 2-step plan for testing purposes without using an actual LLM.
// This allows the system to run locally without external AI dependencies.
type NaivePlanner struct{}

// NewNaivePlanner creates a new hardcoded planner for testing.
// Purpose: Factory function for creating the mock planner adapter.
// Inputs: None
// Outputs:
//   - ports.Planner: Initialized planner ready for use
func NewNaivePlanner() ports.Planner {
	return &NaivePlanner{}
}

// CreatePlan generates a fixed execution plan with 2 hardcoded steps.
// Purpose: Provides a predictable plan for testing the orchestrator without LLM.
//          Step 1: Analysis (THINK type)
//          Step 2: Execution (TOOL_CALL type with LOW risk)
// Inputs:
//   - ctx: Context for cancellation and timeout control (unused in this implementation)
//   - task: The task requiring a plan (used for TaskID linking)
//   - tools: Available tools (unused in this simple implementation)
// Outputs:
//   - *domain.Plan: Fixed plan with 2 steps and dummy risk summary
//   - error: Always returns nil (this implementation cannot fail)
func (p *NaivePlanner) CreatePlan(ctx context.Context, task *domain.Task, tools []domain.ToolMetadata) (*domain.Plan, error) {
	plan := &domain.Plan{
		ID:          uuid.New().String(),
		TaskID:      task.ID,
		Goal:        task.NormalizedIntent,
		RiskSummary: "Low risk - automated execution with 2 simple steps",
		Steps: []domain.Step{
			{
				ID:               uuid.New().String(),
				Title:            "Analiza zahteva",
				Description:      "Razumevanje korisničkog zahteva i pripreme za izvršenje",
				Type:             domain.StepTypeThink,
				Status:           domain.StepStatusPending,
				ToolName:         "",
				ToolInput:        task.Input,
				RiskLevel:        domain.RiskLevelLow,
				RequiresApproval: false,
				RetryCount:       0,
				ResultRef:        "",
			},
			{
				ID:               uuid.New().String(),
				Title:            "Izvršenje akcije",
				Description:      "Izvršavanje planiranog zadatka prema zahtevima korisnika",
				Type:             domain.StepTypeToolCall,
				Status:           domain.StepStatusPending,
				ToolName:         "mock_executor",
				ToolInput:        task.Input,
				RiskLevel:        domain.RiskLevelLow,
				RequiresApproval: false,
				RetryCount:       0,
				ResultRef:        "",
			},
		},
	}

	return plan, nil
}
