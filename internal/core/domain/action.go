package domain

// ApprovalStatus represents the status of an approval request
type ApprovalStatus string

// Approval status constants
const (
	ApprovalStatusOpen     ApprovalStatus = "OPEN"
	ApprovalStatusApproved ApprovalStatus = "APPROVED"
	ApprovalStatusRejected ApprovalStatus = "REJECTED"
)

// ToolCall represents a tool execution action
type ToolCall struct {
	ID            string `json:"id"`
	TaskID        string `json:"task_id"`
	StepID        string `json:"step_id"`
	ToolName      string `json:"tool_name"`
	InputPayload  string `json:"input_payload"`
	OutputPayload string `json:"output_payload"`
	Status        string `json:"status"`
	ErrorMessage  string `json:"error_message"`
	DurationMs    int64  `json:"duration_ms"`
}

// ApprovalRequest represents a request for user approval
type ApprovalRequest struct {
	ID            string         `json:"id"`
	TaskID        string         `json:"task_id"`
	StepID        string         `json:"step_id"`
	ActionSummary string         `json:"action_summary"`
	RiskReason    string         `json:"risk_reason"`
	Status        ApprovalStatus `json:"status"`
	ApprovedBy    string         `json:"approved_by"`
}
