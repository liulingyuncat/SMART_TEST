// Package protocol provides JSON-RPC 2.0 protocol handling for MCP.
package protocol

import (
	"context"
	"encoding/json"
	"fmt"
)

// JSON-RPC 2.0 error codes
const (
	// Standard JSON-RPC 2.0 error codes
	ParseError     = -32700 // Invalid JSON
	InvalidRequest = -32600 // Invalid JSON-RPC request
	MethodNotFound = -32601 // Method not found
	InvalidParams  = -32602 // Invalid method parameters
	InternalError  = -32603 // Internal error

	// Custom MCP error codes
	Unauthorized = -32001 // Authentication failed (401)
	Forbidden    = -32002 // Permission denied (403)
	NotFound     = -32003 // Resource not found (404)
	Conflict     = -32004 // Resource conflict (409)
)

// Request represents a JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"` // Can be string, number, or null
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	Context context.Context `json:"-"` // Request context (not serialized)
}

// Response represents a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents a JSON-RPC 2.0 error object.
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("JSON-RPC error %d: %s", e.Code, e.Message)
}

// ParseRequest parses a JSON-RPC 2.0 request from raw bytes.
func ParseRequest(data []byte) (*Request, error) {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, &Error{
			Code:    ParseError,
			Message: "Parse error",
			Data:    err.Error(),
		}
	}

	// Validate JSON-RPC version
	if req.JSONRPC != "2.0" {
		return nil, &Error{
			Code:    InvalidRequest,
			Message: "Invalid Request",
			Data:    "jsonrpc field must be '2.0'",
		}
	}

	// Validate method is present
	if req.Method == "" {
		return nil, &Error{
			Code:    InvalidRequest,
			Message: "Invalid Request",
			Data:    "method field is required",
		}
	}

	return &req, nil
}

// BuildSuccessResponse creates a JSON-RPC 2.0 success response.
func BuildSuccessResponse(id interface{}, result interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// BuildErrorResponse creates a JSON-RPC 2.0 error response.
func BuildErrorResponse(id interface{}, code int, message string, data interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// Serialize converts a Response to JSON bytes.
func (r *Response) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

// IsNotification returns true if the request is a notification (no ID).
func (r *Request) IsNotification() bool {
	return r.ID == nil
}

// MapHTTPStatusToError maps HTTP status codes to JSON-RPC error codes.
func MapHTTPStatusToError(statusCode int) (int, string) {
	switch statusCode {
	case 400:
		return InvalidParams, "Bad Request"
	case 401:
		return Unauthorized, "Unauthorized"
	case 403:
		return Forbidden, "Forbidden"
	case 404:
		return NotFound, "Not Found"
	case 409:
		return Conflict, "Conflict"
	case 500, 502, 503, 504:
		return InternalError, "Internal Server Error"
	default:
		if statusCode >= 400 && statusCode < 500 {
			return InvalidRequest, "Client Error"
		}
		return InternalError, "Server Error"
	}
}
