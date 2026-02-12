package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// NewDefaultConfig creates a Config instance with sensible default values.
// Purpose: Provides safe defaults for all settings to enable running without
//          any configuration. Suitable for development and testing.
// Inputs: None
// Outputs:
//   - *Config: Configuration with all defaults set
func NewDefaultConfig() *Config {
	return &Config{
		// Server defaults
		ServerPort: 8080,
		ServerHost: "0.0.0.0",

		// Timeout defaults
		RequestTimeout: 30 * time.Second,
		IdleTimeout:    60 * time.Second,

		// Limit defaults
		MaxBodySize:       10 * 1024 * 1024, // 10MB
		MaxFileUploadSize: 50 * 1024 * 1024, // 50MB

		// Security defaults
		AllowedMIMETypes: []string{
			"application/json",
			"application/pdf",
			"text/plain",
			"image/png",
			"image/jpeg",
		},

		// Logging defaults
		LogLevel: "info",

		// Feature flags defaults
		EnableMetrics: false,

		// LLM defaults
		OpenAIAPIKey:    "", // Must be set via env var
		AnthropicAPIKey: "", // Must be set via env var
		DefaultLLMModel: "gpt-4o-mini",
		LLMTimeout:      60 * time.Second,
		LLMMaxRetries:   3,
	}
}

// LoadFromEnv creates a Config by loading values from environment variables.
// Purpose: Enables configuration via environment variables for production deployments.
//          Starts with defaults and overrides any values found in environment.
//          This is the recommended way to configure production instances.
// Inputs: None (reads from os.Getenv)
// Outputs:
//   - *Config: Configuration with env overrides applied
//   - error: Returns error if parsing fails (invalid duration, invalid int, etc.)
func LoadFromEnv() (*Config, error) {
	cfg := NewDefaultConfig()

	// Server settings
	if port := os.Getenv("SERVER_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
		}
		cfg.ServerPort = p
	}

	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.ServerHost = host
	}

	// Timeouts
	if timeout := os.Getenv("REQUEST_TIMEOUT"); timeout != "" {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid REQUEST_TIMEOUT: %w", err)
		}
		cfg.RequestTimeout = d
	}

	if timeout := os.Getenv("IDLE_TIMEOUT"); timeout != "" {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid IDLE_TIMEOUT: %w", err)
		}
		cfg.IdleTimeout = d
	}

	// Limits
	if size := os.Getenv("MAX_BODY_SIZE"); size != "" {
		s, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_BODY_SIZE: %w", err)
		}
		cfg.MaxBodySize = s
	}

	if size := os.Getenv("MAX_FILE_SIZE"); size != "" {
		s, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_FILE_SIZE: %w", err)
		}
		cfg.MaxFileUploadSize = s
	}

	// Security
	if mimes := os.Getenv("ALLOWED_MIMES"); mimes != "" {
		cfg.AllowedMIMETypes = strings.Split(mimes, ",")
		// Trim whitespace from each type
		for i := range cfg.AllowedMIMETypes {
			cfg.AllowedMIMETypes[i] = strings.TrimSpace(cfg.AllowedMIMETypes[i])
		}
	}

	// Logging
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.LogLevel = level
	}

	// Feature flags
	if metrics := os.Getenv("ENABLE_METRICS"); metrics != "" {
		cfg.EnableMetrics = metrics == "true" || metrics == "1"
	}

	// LLM configuration
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		cfg.OpenAIAPIKey = key
	}

	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		cfg.AnthropicAPIKey = key
	}

	if model := os.Getenv("DEFAULT_LLM_MODEL"); model != "" {
		cfg.DefaultLLMModel = model
	}

	if timeout := os.Getenv("LLM_TIMEOUT"); timeout != "" {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid LLM_TIMEOUT: %w", err)
		}
		cfg.LLMTimeout = d
	}

	if retries := os.Getenv("LLM_MAX_RETRIES"); retries != "" {
		r, err := strconv.Atoi(retries)
		if err != nil {
			return nil, fmt.Errorf("invalid LLM_MAX_RETRIES: %w", err)
		}
		cfg.LLMMaxRetries = r
	}

	// Validate the loaded configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate checks if the configuration is valid and safe to use.
// Purpose: Ensures configuration values are within acceptable ranges and
//          required fields are set. Prevents runtime errors from bad config.
// Inputs: None (validates receiver)
// Outputs:
//   - error: Returns error if validation fails, nil if config is valid
func (c *Config) Validate() error {
	// Server validation
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d (must be 1-65535)", c.ServerPort)
	}

	if c.ServerHost == "" {
		return fmt.Errorf("server host cannot be empty")
	}

	// Timeout validation
	if c.RequestTimeout < time.Second {
		return fmt.Errorf("request timeout too short: %v (minimum 1s)", c.RequestTimeout)
	}

	if c.IdleTimeout < time.Second {
		return fmt.Errorf("idle timeout too short: %v (minimum 1s)", c.IdleTimeout)
	}

	// Limit validation
	if c.MaxBodySize < 1024 {
		return fmt.Errorf("max body size too small: %d (minimum 1KB)", c.MaxBodySize)
	}

	if c.MaxFileUploadSize < 1024 {
		return fmt.Errorf("max file upload size too small: %d (minimum 1KB)", c.MaxFileUploadSize)
	}

	// Log level validation
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s (must be: debug, info, warn, error)", c.LogLevel)
	}

	// LLM validation
	if c.LLMTimeout < time.Second {
		return fmt.Errorf("LLM timeout too short: %v (minimum 1s)", c.LLMTimeout)
	}

	if c.LLMMaxRetries < 0 {
		return fmt.Errorf("LLM max retries cannot be negative: %d", c.LLMMaxRetries)
	}

	return nil
}

// HasOpenAIKey checks if OpenAI API key is configured.
// Purpose: Allows conditional logic based on available LLM providers.
// Inputs: None
// Outputs:
//   - bool: true if OpenAI API key is set, false otherwise
func (c *Config) HasOpenAIKey() bool {
	return c.OpenAIAPIKey != ""
}

// HasAnthropicKey checks if Anthropic API key is configured.
// Purpose: Allows conditional logic based on available LLM providers.
// Inputs: None
// Outputs:
//   - bool: true if Anthropic API key is set, false otherwise
func (c *Config) HasAnthropicKey() bool {
	return c.AnthropicAPIKey != ""
}
