package domain

// StepType represents the type of a plan step
type StepType string

// Step type constants
const (
	StepTypeThink        StepType = "THINK"
	StepTypeToolCall     StepType = "TOOL_CALL"
	StepTypeDecision     StepType = "DECISION"
	StepTypeApprovalGate StepType = "APPROVAL_GATE"
	StepTypeVerify       StepType = "VERIFY"
)

// StepStatus represents the current status of a step
type StepStatus string

// Step status constants
const (
	StepStatusPending    StepStatus = "PENDING"
	StepStatusInProgress StepStatus = "IN_PROGRESS"
	StepStatusCompleted  StepStatus = "COMPLETED"
	StepStatusFailed     StepStatus = "FAILED"
	StepStatusSkipped    StepStatus = "SKIPPED"
)

// RiskLevel represents the risk level of a step
type RiskLevel string

// Risk level constants
const (
	RiskLevelLow  RiskLevel = "LOW"
	RiskLevelHigh RiskLevel = "HIGH"
)

// Step represents a single step in a plan
type Step struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Type             StepType   `json:"type"`
	Status           StepStatus `json:"status"`
	ToolName         string     `json:"tool_name,omitempty"`
	ToolInput        string     `json:"tool_input"`
	RiskLevel        RiskLevel  `json:"risk_level"`
	RequiresApproval bool       `json:"requires_approval"`
	RetryCount       int        `json:"retry_count"`
	ResultRef        string     `json:"result_ref"`
}

// Plan represents an execution plan for a task
type Plan struct {
	ID          string `json:"id"`
	TaskID      string `json:"task_id"`
	Goal        string `json:"goal"`
	Steps       []Step `json:"steps"`
	RiskSummary string `json:"risk_summary"`
}
