package domain

import "time"

// TaskStatus represents the current status of a task
type TaskStatus string

// Task status constants
const (
	TaskStatusNew              TaskStatus = "NEW"
	TaskStatusPlanning         TaskStatus = "PLANNING"
	TaskStatusExecuting        TaskStatus = "EXECUTING"
	TaskStatusWaitingApproval  TaskStatus = "WAITING_APPROVAL"
	TaskStatusVerifying        TaskStatus = "VERIFYING"
	TaskStatusDone             TaskStatus = "DONE"
	TaskStatusFailed           TaskStatus = "FAILED"
	TaskStatusCanceled         TaskStatus = "CANCELED"
)

// Task represents a user task/request in the system
type Task struct {
	ID                string            `json:"id"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	FinishedAt        time.Time         `json:"finished_at"`
	Status            TaskStatus        `json:"status"`
	Input             string            `json:"input"`
	NormalizedIntent  string            `json:"normalized_intent"`
	UserID            string            `json:"user_id"`
	Channel           string            `json:"channel"`
	Role              string            `json:"role"`
	TargetAgent       string            `json:"target_agent"`
	PlanID            string            `json:"plan_id"`
	CurrentStepID     string            `json:"current_step_id"`
	Artifacts         map[string]string `json:"artifacts"`
	Metadata          map[string]string `json:"metadata"`
	UsageTokens       int               `json:"usage_tokens"`
	CostEstimate      float64           `json:"cost_estimate"`
}
