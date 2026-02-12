package memory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/JAROBOTAI/jaro/internal/core/domain"
	"github.com/JAROBOTAI/jaro/internal/core/ports"
)

// AuditRepository is an in-memory implementation of the ports.AuditRepository interface.
// It logs audit events to stdout in JSON format for local development and debugging.
// Events are not persisted and are only visible in console output.
type AuditRepository struct{}

// NewAuditRepository creates a new stdout-based audit repository.
// Purpose: Factory function for creating the console-logging audit adapter.
// Inputs: None
// Outputs:
//   - ports.AuditRepository: Initialized repository ready for use
func NewAuditRepository() ports.AuditRepository {
	return &AuditRepository{}
}

// SaveEvent logs an audit event to stdout in JSON format.
// Purpose: Creates an audit trail by printing events to console for debugging.
//          This implementation is non-persistent and suitable only for development.
// Inputs:
//   - ctx: Context for cancellation and timeout control (unused in this implementation)
//   - event: The audit event to log (includes timestamp, actor, event type, payload)
// Outputs:
//   - error: Always returns nil (logging failures are not treated as errors)
func (r *AuditRepository) SaveEvent(ctx context.Context, event *domain.AuditEvent) error {
	if event == nil {
		fmt.Println("[AUDIT] WARNING: attempted to log nil event")
		return nil
	}

	// Serialize payload to JSON for readable output
	payloadJSON, err := json.MarshalIndent(event.Payload, "", "  ")
	if err != nil {
		fmt.Printf("[AUDIT] %s: <failed to serialize payload>\n", event.EventType)
		return nil
	}

	// Log event to console
	fmt.Printf("[AUDIT] %s | TaskID: %s | Actor: %s | Time: %s\nPayload:\n%s\n\n",
		event.EventType,
		event.TaskID,
		event.Actor,
		event.Timestamp.Format("2006-01-02 15:04:05"),
		string(payloadJSON),
	)

	return nil
}
