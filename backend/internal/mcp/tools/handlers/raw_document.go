// Package handlers provides MCP tool handler implementations.
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// BaseHandler provides common functionality for tool handlers.
type BaseHandler struct {
	client *client.BackendClient
}

// NewBaseHandler creates a new BaseHandler.
func NewBaseHandler(c *client.BackendClient) *BaseHandler {
	return &BaseHandler{client: c}
}

// Annotations returns nil by default. Override in specific handlers if needed.
func (h *BaseHandler) Annotations() *tools.ToolAnnotations {
	return nil
}

// GetInt extracts an integer from args, handling float64 from JSON.
func GetInt(args map[string]interface{}, key string) (int, error) {
	val, ok := args[key]
	if !ok {
		return 0, fmt.Errorf("missing required field: %s", key)
	}

	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		// 支持字符串形式的整数
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("field %s must be an integer, got string '%s'", key, v)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("field %s must be an integer, got %T", key, val)
	}
}

// GetString extracts a string from args.
func GetString(args map[string]interface{}, key string) (string, error) {
	val, ok := args[key]
	if !ok {
		return "", fmt.Errorf("missing required field: %s", key)
	}

	switch v := val.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		// 如果是整数值，转换为整数字符串
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10), nil
		}
		return fmt.Sprintf("%v", v), nil
	default:
		return "", fmt.Errorf("field %s must be a string, got %T", key, val)
	}
}

// GetOptionalString extracts an optional string from args.
func GetOptionalString(args map[string]interface{}, key string, defaultVal string) string {
	val, ok := args[key]
	if !ok {
		return defaultVal
	}
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10)
		}
		return fmt.Sprintf("%v", v)
	default:
		return defaultVal
	}
}

// GetOptionalInt extracts an optional integer from args.
func GetOptionalInt(args map[string]interface{}, key string, defaultVal int) int {
	val, ok := args[key]
	if !ok {
		return defaultVal
	}
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return defaultVal
		}
		return i
	default:
		return defaultVal
	}
}

// ListRawDocumentsHandler handles listing raw documents.
type ListRawDocumentsHandler struct {
	*BaseHandler
}

func NewListRawDocumentsHandler(c *client.BackendClient) *ListRawDocumentsHandler {
	return &ListRawDocumentsHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListRawDocumentsHandler) Name() string {
	return "list_raw_documents"
}

func (h *ListRawDocumentsHandler) Description() string {
	return "获取项目中已完成转换的文档列表"
}

func (h *ListRawDocumentsHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
		},
		"required": []interface{}{"project_id"},
	}
}

func (h *ListRawDocumentsHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/raw-documents", projectID)
	data, err := h.client.Get(ctx, path, nil)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 解析响应数据，只返回已完成转换的文档
	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		return tools.NewErrorResult("failed to parse response: " + err.Error()), nil
	}

	// 提取转换完成的文档
	var convertedDocs []map[string]interface{}
	if docs, ok := response["data"].(map[string]interface{}); ok {
		if docList, ok := docs["documents"].([]interface{}); ok {
			for _, doc := range docList {
				if docMap, ok := doc.(map[string]interface{}); ok {
					// 只包含转换状态为completed的文档
					if status, ok := docMap["convert_status"].(string); ok && status == "completed" {
						// 构建转换文档对象，确保文件名是正确的UTF-8编码
						convertedDoc := map[string]interface{}{
							"id":                  docMap["id"],
							"project_id":          docMap["project_id"],
							"converted_filename":  sanitizeUTF8(convertToString(docMap["converted_filename"])),
							"converted_file_size": docMap["converted_file_size"],
							"converted_time":      docMap["converted_time"],
							"original_filename":   sanitizeUTF8(convertToString(docMap["original_filename"])),
						}
						convertedDocs = append(convertedDocs, convertedDoc)
					}
				}
			}
		}
	}

	// 返回转换完成的文档列表
	result := map[string]interface{}{
		"documents": convertedDocs,
		"total":     len(convertedDocs),
	}

	// 使用自定义编码器确保不转义HTML字符，保留原始UTF-8
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(result); err != nil {
		return tools.NewErrorResult("failed to encode result: " + err.Error()), nil
	}

	return tools.NewJSONResult(buf.String()), nil
}

// GetConvertedDocumentHandler handles getting a single converted document's content.
type GetConvertedDocumentHandler struct {
	*BaseHandler
}

func NewGetConvertedDocumentHandler(c *client.BackendClient) *GetConvertedDocumentHandler {
	return &GetConvertedDocumentHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *GetConvertedDocumentHandler) Name() string {
	return "get_converted_document"
}

func (h *GetConvertedDocumentHandler) Description() string {
	return "获取单个转换文档的完整内容"
}

func (h *GetConvertedDocumentHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "转换文档ID",
			},
		},
		"required": []interface{}{"project_id", "id"},
	}
}

func (h *GetConvertedDocumentHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	id, err := GetInt(args, "id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// Call the preview API endpoint which returns filename and content
	// GET /api/v1/raw-documents/:id/converted/preview
	previewPath := fmt.Sprintf("/api/v1/raw-documents/%d/converted/preview", id)
	contentData, err := h.client.Get(ctx, previewPath, nil)
	if err != nil {
		return tools.NewErrorResult("failed to get document content: " + err.Error()), nil
	}

	// Parse the JSON response - API returns {code: 0, message: "success", data: {...}}
	var apiResponse map[string]interface{}
	if err := json.Unmarshal(contentData, &apiResponse); err != nil {
		return tools.NewErrorResult("failed to parse response: " + err.Error()), nil
	}

	// Check if API call was successful
	code := int(0)
	if codeVal, ok := apiResponse["code"].(float64); ok {
		code = int(codeVal)
	}
	if code != 0 {
		msg := "unknown error"
		if msgVal, ok := apiResponse["message"].(string); ok {
			msg = msgVal
		}
		return tools.NewErrorResult(fmt.Sprintf("API error: %s", msg)), nil
	}

	// Extract data from response
	var convertedFilename string
	var content string

	if dataObj, ok := apiResponse["data"].(map[string]interface{}); ok {
		convertedFilename = sanitizeUTF8(convertToString(dataObj["filename"]))
		content = convertToString(dataObj["content"])
	} else {
		return tools.NewErrorResult("invalid response format: missing data field"), nil
	}

	if content == "" {
		return tools.NewErrorResult("document content is empty"), nil
	}

	// Return the complete document content
	result := map[string]interface{}{
		"id":                 id,
		"project_id":         projectID,
		"converted_filename": convertedFilename,
		"content":            content,
		"success":            true,
	}

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(result); err != nil {
		return tools.NewErrorResult("failed to encode result: " + err.Error()), nil
	}

	return tools.NewJSONResult(buf.String()), nil
}

// convertToString converts any value to string safely
func convertToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", val)
	}
}

// sanitizeUTF8 ensures the string is valid UTF-8
// If it contains invalid sequences, it replaces them with replacement character
func sanitizeUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	// Replace invalid UTF-8 sequences
	b := []byte(s)
	var result []byte
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if r == utf8.RuneError && size == 1 {
			// Invalid UTF-8 sequence
			result = append(result, '?')
		} else {
			result = append(result, b[:size]...)
		}
		b = b[size:]
	}
	return string(result)
}
