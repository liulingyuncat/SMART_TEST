// Package tools provides MCP tool handling functionality.
package tools

import (
	"context"
	"encoding/json"
)

// ToolHandler defines the interface that all MCP tool handlers must implement.
type ToolHandler interface {
	// Name returns the unique name of the tool.
	Name() string

	// Description returns a human-readable description of the tool.
	Description() string

	// InputSchema returns the JSON Schema for the tool's input parameters.
	InputSchema() map[string]interface{}

	// Annotations returns optional annotations for the tool (can return nil).
	Annotations() *ToolAnnotations

	// Execute runs the tool with the provided arguments.
	Execute(ctx context.Context, args map[string]interface{}) (ToolResult, error)
}

// ToolResult represents the result of a tool execution.
type ToolResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ContentItem represents a single content item in a tool result.
type ContentItem struct {
	Type     string `json:"type"`               // "text", "image", "resource"
	Text     string `json:"text,omitempty"`     // For type="text"
	Data     string `json:"data,omitempty"`     // For type="image" (base64)
	MimeType string `json:"mimeType,omitempty"` // MIME type
}

// NewTextResult creates a ToolResult with a single text content item.
func NewTextResult(text string) ToolResult {
	return ToolResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: text,
			},
		},
		IsError: false,
	}
}

// NewErrorResult creates a ToolResult indicating an error.
func NewErrorResult(message string) ToolResult {
	return ToolResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: message,
			},
		},
		IsError: true,
	}
}

// NewJSONResult creates a ToolResult with JSON text content.
func NewJSONResult(jsonText string) ToolResult {
	// 验证输入是否为有效的JSON
	var validateJSON interface{}
	if err := json.Unmarshal([]byte(jsonText), &validateJSON); err != nil {
		// 如果不是有效的JSON，降级为纯文本返回，避免序列化失败
		return ToolResult{
			Content: []ContentItem{
				{
					Type: "text",
					Text: jsonText,
				},
			},
			IsError: false,
		}
	}

	// 有效的JSON，设置正确的MimeType
	return ToolResult{
		Content: []ContentItem{
			{
				Type:     "text",
				Text:     jsonText,
				MimeType: "application/json",
			},
		},
		IsError: false,
	}
}

// MustMarshalJSON marshals an object to JSON string, panics on error.
func MustMarshalJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
