package ports

import "time"

// Clock provides time-related operations for testing and time-travel debugging.
// Purpose: Abstracts system time to enable deterministic testing and replay scenarios.
//          Allows mocking current time in tests without system clock manipulation.
// Why: Direct time.Now() calls make tests non-deterministic and impossible to replay.
type Clock interface {
	// Now returns the current time.
	// Purpose: Provides current time for timestamps, expiration checks, etc.
	// Inputs: None
	// Outputs:
	//   - time.Time: Current time (real or mocked)
	Now() time.Time
}

// IDGenerator provides unique identifier generation for entities.
// Purpose: Abstracts ID generation to enable deterministic testing and custom ID strategies.
//          Allows mocking IDs in tests and using different ID formats (UUID, ULID, snowflake).
// Why: Direct uuid.New() calls make tests non-deterministic and couple to specific ID format.
type IDGenerator interface {
	// Generate creates a new unique identifier.
	// Purpose: Generates unique IDs for tasks, plans, steps, events, etc.
	// Inputs: None
	// Outputs:
	//   - string: Unique identifier (UUID format in default implementation)
	Generate() string
}

// Logger provides structured logging capabilities.
// Purpose: Abstracts logging to enable testing, custom log formats, and centralized logging.
//          Separates concerns between business logic and log output format/destination.
// Why: Direct fmt.Println or log.Printf couples code to specific output and makes testing harder.
type Logger interface {
	// Info logs an informational message.
	// Purpose: Records normal operational messages for debugging and monitoring.
	// Inputs:
	//   - msg: Log message
	//   - fields: Key-value pairs for structured logging
	// Outputs: None
	Info(msg string, fields map[string]interface{})

	// Error logs an error message.
	// Purpose: Records error conditions for alerting and troubleshooting.
	// Inputs:
	//   - msg: Error description
	//   - err: Error object
	//   - fields: Additional context as key-value pairs
	// Outputs: None
	Error(msg string, err error, fields map[string]interface{})

	// Warn logs a warning message.
	// Purpose: Records potentially problematic situations that aren't errors.
	// Inputs:
	//   - msg: Warning message
	//   - fields: Key-value pairs for context
	// Outputs: None
	Warn(msg string, fields map[string]interface{})
}
