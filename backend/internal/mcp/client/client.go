package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"webtest/internal/mcp/config"
	"webtest/internal/mcp/protocol"
)

// BackendClient provides HTTP client functionality for backend API calls.
type BackendClient struct {
	httpClient  *http.Client
	baseURL     string
	authManager *AuthManager
	retryCount  int
	retryDelay  time.Duration
}

// sharedTransport is a global HTTP transport with connection pooling for better performance
var sharedTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   10 * time.Second, // Connection timeout
		KeepAlive: 30 * time.Second, // TCP keep-alive
	}).DialContext,
	MaxIdleConns:        100,              // Total max idle connections
	MaxIdleConnsPerHost: 10,               // Max idle connections per host
	MaxConnsPerHost:     20,               // Max total connections per host
	IdleConnTimeout:     90 * time.Second, // How long idle connections stay in pool
	TLSHandshakeTimeout: 10 * time.Second,
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true, // Skip certificate verification for self-signed certs
	},
	ExpectContinueTimeout: 1 * time.Second,
	ForceAttemptHTTP2:     true,  // Enable HTTP/2 when possible
	DisableCompression:    false, // Allow compression
}

// NewBackendClient creates a new BackendClient with the provided configuration.
func NewBackendClient(cfg config.BackendConfig, authManager *AuthManager) *BackendClient {
	return &BackendClient{
		httpClient: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: sharedTransport, // Use shared transport for connection reuse
		},
		baseURL:     strings.TrimSuffix(cfg.BaseURL, "/"),
		authManager: authManager,
		retryCount:  cfg.RetryCount,
		retryDelay:  cfg.RetryDelay,
	}
}

// Get sends a GET request to the specified path with optional query parameters.
func (c *BackendClient) Get(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	fullURL := c.baseURL + path
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		fullURL += "?" + q.Encode()
	}
	return c.doRequest(ctx, http.MethodGet, fullURL, nil)
}

// Post sends a POST request with JSON body.
func (c *BackendClient) Post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequestWithBody(ctx, http.MethodPost, c.baseURL+path, body)
}

// Put sends a PUT request with JSON body.
func (c *BackendClient) Put(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequestWithBody(ctx, http.MethodPut, c.baseURL+path, body)
}

// Patch sends a PATCH request with JSON body.
func (c *BackendClient) Patch(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequestWithBody(ctx, http.MethodPatch, c.baseURL+path, body)
}

// Delete sends a DELETE request.
func (c *BackendClient) Delete(ctx context.Context, path string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodDelete, c.baseURL+path, nil)
}

// doRequestWithBody handles requests with JSON body.
func (c *BackendClient) doRequestWithBody(ctx context.Context, method, url string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}
	return c.doRequest(ctx, method, url, bodyReader)
}

// doRequest performs the actual HTTP request with retry logic.
func (c *BackendClient) doRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= c.retryCount; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := c.retryDelay * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		// Need to re-read body for retry if it's not nil
		var bodyReader io.Reader
		if body != nil {
			if seeker, ok := body.(io.Seeker); ok {
				seeker.Seek(0, io.SeekStart)
				bodyReader = body
			} else {
				// For non-seekable readers, we can't retry with body
				bodyReader = body
			}
		}

		result, err := c.doSingleRequest(ctx, method, url, bodyReader)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Check if error is retryable
		if !c.isRetryable(err) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", c.retryCount, lastErr)
}

// doSingleRequest performs a single HTTP request.
func (c *BackendClient) doSingleRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers - use appropriate authentication method
	if c.authManager != nil {
		c.authManager.SetAuthHeadersWithContext(ctx, req)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &retryableError{err: err}
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, c.handleErrorStatus(resp.StatusCode, respBody)
	}

	return respBody, nil
}

// handleErrorStatus converts HTTP error status to protocol.Error.
func (c *BackendClient) handleErrorStatus(statusCode int, body []byte) error {
	code, message := protocol.MapHTTPStatusToError(statusCode)

	// Try to extract error message from response body
	var apiError struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if json.Unmarshal(body, &apiError) == nil {
		if apiError.Error != "" {
			message = apiError.Error
		} else if apiError.Message != "" {
			message = apiError.Message
		}
	}

	err := &protocol.Error{
		Code:    code,
		Message: message,
		Data: map[string]interface{}{
			"http_status": statusCode,
		},
	}

	// Mark 5xx errors as retryable
	if statusCode >= 500 {
		return &retryableError{err: err}
	}

	return err
}

// isRetryable checks if an error is retryable.
func (c *BackendClient) isRetryable(err error) bool {
	_, ok := err.(*retryableError)
	return ok
}

// retryableError wraps errors that should be retried.
type retryableError struct {
	err error
}

func (e *retryableError) Error() string {
	return e.err.Error()
}

func (e *retryableError) Unwrap() error {
	return e.err
}
