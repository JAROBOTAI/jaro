// Package config provides configuration management for the JARO application.
// It supports loading configuration from environment variables and provides
// sensible defaults for all settings. This package follows the "Open Core"
// strategy by keeping all configuration external to code.
package config

import "time"

// Config holds all application configuration settings.
// Purpose: Centralizes all configurable values to eliminate magic numbers
//          and enable different configurations for dev/staging/production.
//          Supports environment-based configuration for secure secret management.
// Fields are organized by category for clarity and maintainability.
type Config struct {
	// Server settings - HTTP server configuration
	ServerPort int    // Port to listen on (default: 8080)
	ServerHost string // Host to bind to (default: "0.0.0.0")

	// Timeouts - Request and connection timeout settings
	RequestTimeout time.Duration // Maximum duration for request processing (default: 30s)
	IdleTimeout    time.Duration // Maximum duration for idle connections (default: 60s)

	// Limits - Size restrictions for request payloads
	MaxBodySize       int64 // Maximum request body size in bytes (default: 10MB)
	MaxFileUploadSize int64 // Maximum file upload size in bytes (default: 50MB)

	// Security - Input validation configuration
	AllowedMIMETypes []string // Whitelist of allowed MIME types for uploads

	// Logging - Log level configuration
	LogLevel string // Log verbosity: debug, info, warn, error (default: "info")

	// Feature Flags - Optional features that can be toggled
	EnableMetrics bool // Enable Prometheus metrics endpoint (default: false)

	// LLM Integration - AI provider configuration (Phase 3)
	OpenAIAPIKey     string // OpenAI API key (load from env: OPENAI_API_KEY)
	AnthropicAPIKey  string // Anthropic API key (load from env: ANTHROPIC_API_KEY)
	DefaultLLMModel  string // Default LLM model to use (default: "gpt-4o-mini")
	LLMTimeout       time.Duration // Maximum duration for LLM requests (default: 60s)
	LLMMaxRetries    int    // Maximum retry attempts for LLM failures (default: 3)
}
