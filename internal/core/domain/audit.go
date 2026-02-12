package domain

import "time"

// AuditEvent represents an audit log entry for tracking system events
type AuditEvent struct {
	ID               string                 `json:"id"`
	TaskID           string                 `json:"task_id"`
	CorrelationID    string                 `json:"correlation_id"`
	Timestamp        time.Time              `json:"timestamp"`
	EventType        string                 `json:"event_type"`
	Payload          map[string]interface{} `json:"payload"`
	Actor            string                 `json:"actor"`
	BehaviorVersion  string                 `json:"behavior_version"`
}
