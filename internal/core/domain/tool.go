package domain

// ToolMetadata represents metadata information about a tool
type ToolMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	RiskLevel   RiskLevel `json:"risk_level"`
}

// Tool represents an executable tool interface
type Tool interface {
	// Name returns the unique name of the tool
	Name() string
	// Description returns a human-readable description of the tool
	Description() string
	// Execute runs the tool with the given input and returns the result or an error
	Execute(input string) (string, error)
}

// StepResult represents the result of executing a step
type StepResult struct {
	StepID       string `json:"step_id"`
	Success      bool   `json:"success"`
	Output       string `json:"output"`
	ErrorMessage string `json:"error_message,omitempty"`
	DurationMs   int64  `json:"duration_ms"`
}
