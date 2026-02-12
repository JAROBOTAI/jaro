package memory

import (
	"time"

	"github.com/JAROBOTAI/jaro/internal/core/ports"
	"github.com/google/uuid"
)

// SystemClock is a real-time clock implementation that uses system time.
// It implements ports.Clock for production use.
type SystemClock struct{}

// NewSystemClock creates a new system clock.
// Purpose: Factory function for creating a real-time clock adapter.
// Inputs: None
// Outputs:
//   - ports.Clock: Clock implementation using system time
func NewSystemClock() ports.Clock {
	return &SystemClock{}
}

// Now returns the current system time.
// Purpose: Provides real system time for production timestamps.
// Inputs: None
// Outputs:
//   - time.Time: Current system time
func (c *SystemClock) Now() time.Time {
	return time.Now()
}

// UUIDGenerator is a UUID-based ID generator implementation.
// It implements ports.IDGenerator using Google's UUID library.
type UUIDGenerator struct{}

// NewUUIDGenerator creates a new UUID generator.
// Purpose: Factory function for creating a UUID-based ID generator.
// Inputs: None
// Outputs:
//   - ports.IDGenerator: ID generator that produces UUID v4 strings
func NewUUIDGenerator() ports.IDGenerator {
	return &UUIDGenerator{}
}

// Generate creates a new UUID v4 identifier.
// Purpose: Generates cryptographically random unique identifiers.
// Inputs: None
// Outputs:
//   - string: UUID v4 string (e.g., "550e8400-e29b-41d4-a716-446655440000")
func (g *UUIDGenerator) Generate() string {
	return uuid.New().String()
}

// ConsoleLogger is a simple console-based logger implementation.
// It implements ports.Logger by printing to stdout/stderr.
type ConsoleLogger struct{}

// NewConsoleLogger creates a new console logger.
// Purpose: Factory function for creating a console-based logger adapter.
// Inputs: None
// Outputs:
//   - ports.Logger: Logger implementation that writes to console
func NewConsoleLogger() ports.Logger {
	return &ConsoleLogger{}
}

// Info logs an informational message to console.
// Purpose: Outputs normal operational messages with optional structured fields.
// Inputs:
//   - msg: Log message
//   - fields: Key-value pairs for context (can be nil)
// Outputs: None (prints to stdout)
func (l *ConsoleLogger) Info(msg string, fields map[string]interface{}) {
	// Simple implementation - can be enhanced with structured logging later
	if len(fields) > 0 {
		println("[INFO]", msg, fields)
	} else {
		println("[INFO]", msg)
	}
}

// Error logs an error message to console.
// Purpose: Outputs error conditions with error object and context.
// Inputs:
//   - msg: Error description
//   - err: Error object
//   - fields: Additional context (can be nil)
// Outputs: None (prints to stderr)
func (l *ConsoleLogger) Error(msg string, err error, fields map[string]interface{}) {
	println("[ERROR]", msg, ":", err.Error())
	if len(fields) > 0 {
		println("  Fields:", fields)
	}
}

// Warn logs a warning message to console.
// Purpose: Outputs warning conditions with optional context.
// Inputs:
//   - msg: Warning message
//   - fields: Key-value pairs for context (can be nil)
// Outputs: None (prints to stdout)
func (l *ConsoleLogger) Warn(msg string, fields map[string]interface{}) {
	if len(fields) > 0 {
		println("[WARN]", msg, fields)
	} else {
		println("[WARN]", msg)
	}
}
