// Package config provides configuration management for the MCP server.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete MCP server configuration.
type Config struct {
	Server  ServerConfig  `yaml:"server" json:"server"`
	Backend BackendConfig `yaml:"backend" json:"backend"`
	Auth    AuthConfig    `yaml:"auth" json:"auth"`
}

// ServerConfig contains server-related settings.
type ServerConfig struct {
	Mode      string `yaml:"mode" json:"mode"`             // stdio | http
	HTTPAddr  string `yaml:"http_addr" json:"http_addr"`   // HTTP server address (e.g., :16410)
	HTTPPath  string `yaml:"http_path" json:"http_path"`   // HTTP endpoint path (default: /mcp)
	LogLevel  string `yaml:"log_level" json:"log_level"`   // debug | info | warn | error
	LogFormat string `yaml:"log_format" json:"log_format"` // json | text
}

// BackendConfig contains backend API connection settings.
type BackendConfig struct {
	BaseURL    string        `yaml:"base_url" json:"base_url"`
	Timeout    time.Duration `yaml:"timeout" json:"timeout"`
	RetryCount int           `yaml:"retry_count" json:"retry_count"`
	RetryDelay time.Duration `yaml:"retry_delay" json:"retry_delay"`
}

// AuthConfig contains authentication settings.
type AuthConfig struct {
	TokenEnv        string `yaml:"token_env" json:"token_env"`
	TokenFile       string `yaml:"token_file" json:"token_file"`
	ValidateOnStart bool   `yaml:"validate_on_start" json:"validate_on_start"`
	DynamicToken    bool   `yaml:"dynamic_token" json:"dynamic_token"` // If true, token is passed from client request headers
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Mode:      "stdio",
			HTTPAddr:  ":16410",
			HTTPPath:  "/mcp",
			LogLevel:  "info",
			LogFormat: "json",
		},
		Backend: BackendConfig{
			BaseURL:    "https://localhost:8443",
			Timeout:    30 * time.Second,
			RetryCount: 2,                      // Reduce retries for faster failure
			RetryDelay: 100 * time.Millisecond, // Faster initial retry
		},
		Auth: AuthConfig{
			TokenEnv:        "MCP_AUTH_TOKEN",
			TokenFile:       "",
			ValidateOnStart: false,
			DynamicToken:    true, // Default to dynamic token from client
		},
	}
}

// LoadConfig loads configuration from a YAML file.
// It applies defaults first, then overrides with file values,
// and finally applies environment variable overrides.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	// Read config file if it exists
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				// Config file not found, use defaults
				return applyEnvOverrides(cfg), nil
			}
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	return applyEnvOverrides(cfg), nil
}

// applyEnvOverrides applies environment variable overrides to the config.
func applyEnvOverrides(cfg *Config) *Config {
	// Server overrides
	if v := os.Getenv("MCP_SERVER_MODE"); v != "" {
		cfg.Server.Mode = v
	}
	if v := os.Getenv("MCP_HTTP_ADDR"); v != "" {
		cfg.Server.HTTPAddr = v
	}
	if v := os.Getenv("MCP_HTTP_PATH"); v != "" {
		cfg.Server.HTTPPath = v
	}
	if v := os.Getenv("MCP_LOG_LEVEL"); v != "" {
		cfg.Server.LogLevel = v
	}
	if v := os.Getenv("MCP_LOG_FORMAT"); v != "" {
		cfg.Server.LogFormat = v
	}

	// Backend overrides
	if v := os.Getenv("MCP_BACKEND_URL"); v != "" {
		cfg.Backend.BaseURL = v
	}
	if v := os.Getenv("MCP_BACKEND_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Backend.Timeout = d
		}
	}

	// Auth overrides
	if v := os.Getenv("MCP_TOKEN_ENV"); v != "" {
		cfg.Auth.TokenEnv = v
	}
	if v := os.Getenv("MCP_TOKEN_FILE"); v != "" {
		cfg.Auth.TokenFile = v
	}

	return cfg
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	// Validate server mode
	if c.Server.Mode != "stdio" && c.Server.Mode != "http" {
		return fmt.Errorf("invalid server mode: %s (must be 'stdio' or 'http')", c.Server.Mode)
	}

	// Validate HTTP address if using HTTP mode
	if c.Server.Mode == "http" && c.Server.HTTPAddr == "" {
		return fmt.Errorf("http_addr is required when using HTTP mode")
	}

	// Validate log level
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.Server.LogLevel] {
		return fmt.Errorf("invalid log level: %s", c.Server.LogLevel)
	}

	// Validate backend URL
	if c.Backend.BaseURL == "" {
		return fmt.Errorf("backend base_url is required")
	}

	// Validate timeout
	if c.Backend.Timeout <= 0 {
		return fmt.Errorf("backend timeout must be positive")
	}

	return nil
}
