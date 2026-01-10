// Package client provides HTTP client for backend API communication.
package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"webtest/internal/mcp/config"
)

// contextKey is the key type for context values.
type contextKey string

// TokenContextKey is the context key for the API token passed from request headers.
const TokenContextKey contextKey = "api_token"

// AuthManager manages API token authentication for backend API calls.
// Supports both JWT tokens (Authorization: Bearer) and API tokens (X-API-Token).
type AuthManager struct {
	token        string
	isAPIToken   bool // true if using X-API-Token header, false for Bearer token
	dynamicToken bool // true if token should be read from request context
}

// NewAuthManager creates a new AuthManager from the provided config.
// It reads the token from environment variable or file as configured.
func NewAuthManager(cfg config.AuthConfig) (*AuthManager, error) {
	// If dynamic token mode is enabled, don't require a static token
	if cfg.DynamicToken {
		return &AuthManager{
			dynamicToken: true,
			isAPIToken:   true, // Assume API token for dynamic mode
		}, nil
	}

	var token string

	// Try to read from environment variable first
	if cfg.TokenEnv != "" {
		token = os.Getenv(cfg.TokenEnv)
	}

	// If not found in env, try to read from file
	if token == "" && cfg.TokenFile != "" {
		data, err := os.ReadFile(cfg.TokenFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read token file: %w", err)
		}
		token = strings.TrimSpace(string(data))
	}

	if token == "" {
		return nil, fmt.Errorf("no token found: set %s environment variable or provide token_file", cfg.TokenEnv)
	}

	// Detect token type: JWT tokens are typically longer and have 3 parts separated by dots
	// API tokens from the platform are hex strings (64 chars)
	isAPIToken := !strings.Contains(token, ".") && len(token) == 64

	return &AuthManager{token: token, isAPIToken: isAPIToken, dynamicToken: false}, nil
}

// GetAuthHeader returns the Authorization header value (for JWT tokens).
func (a *AuthManager) GetAuthHeader() string {
	if a.isAPIToken {
		return "" // API tokens don't use Authorization header
	}
	return "Bearer " + a.token
}

// GetAPIToken returns the API token value (for X-API-Token header).
func (a *AuthManager) GetAPIToken() string {
	if a.isAPIToken {
		return a.token
	}
	return ""
}

// IsAPIToken returns true if using API token authentication.
func (a *AuthManager) IsAPIToken() bool {
	return a.isAPIToken
}

// IsDynamicToken returns true if token should be read from request context.
func (a *AuthManager) IsDynamicToken() bool {
	return a.dynamicToken
}

// SetAuthHeaders sets the appropriate authentication headers on the request.
// For dynamic token mode, it reads the token from context.
func (a *AuthManager) SetAuthHeaders(req *http.Request) {
	a.SetAuthHeadersWithContext(req.Context(), req)
}

// SetAuthHeadersWithContext sets the appropriate authentication headers on the request using context.
func (a *AuthManager) SetAuthHeadersWithContext(ctx context.Context, req *http.Request) {
	token := a.token

	// In dynamic mode, get token from context
	if a.dynamicToken {
		if ctxToken, ok := ctx.Value(TokenContextKey).(string); ok && ctxToken != "" {
			token = ctxToken
		}
	}

	if token == "" {
		return
	}

	// Detect token type for dynamic tokens
	isAPIToken := a.isAPIToken
	if a.dynamicToken {
		isAPIToken = !strings.Contains(token, ".") && len(token) == 64
	}

	if isAPIToken {
		req.Header.Set("X-API-Token", token)
	} else {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

// ValidateToken validates the token by calling the backend's /api/v1/auth/me endpoint.
func (a *AuthManager) ValidateToken(ctx context.Context, baseURL string) error {
	// In dynamic token mode, we can't validate at startup
	if a.dynamicToken {
		return nil
	}

	url := strings.TrimSuffix(baseURL, "/") + "/api/v1/auth/me"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Use appropriate authentication header
	a.SetAuthHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip certificate verification for self-signed certs
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("token is invalid or expired")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token validation failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// HasToken returns true if a token is configured.
func (a *AuthManager) HasToken() bool {
	return a.token != ""
}
