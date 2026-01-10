// Package handlers provides MCP tool handler implementations.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// GetCurrentUserInfoHandler handles getting current user information.
type GetCurrentUserInfoHandler struct {
	*BaseHandler
}

func NewGetCurrentUserInfoHandler(c *client.BackendClient) *GetCurrentUserInfoHandler {
	return &GetCurrentUserInfoHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *GetCurrentUserInfoHandler) Name() string {
	return "get_current_user_info"
}

func (h *GetCurrentUserInfoHandler) Description() string {
	return "获取当前Token对应的用户信息，包括user_id、username、nickname和role"
}

func (h *GetCurrentUserInfoHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
		"required":   []interface{}{},
	}
}

func (h *GetCurrentUserInfoHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	// Call the /api/v1/auth/me endpoint
	data, err := h.client.Get(ctx, "/api/v1/auth/me", nil)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// Parse the response to verify it's valid JSON
	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		return tools.NewErrorResult("failed to parse user info response: " + err.Error()), nil
	}

	// Return the valid JSON string
	return tools.NewJSONResult(string(data)), nil
}

// GetCurrentProjectNameHandler handles getting current user's selected project name.
type GetCurrentProjectNameHandler struct {
	*BaseHandler
}

func NewGetCurrentProjectNameHandler(c *client.BackendClient) *GetCurrentProjectNameHandler {
	return &GetCurrentProjectNameHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *GetCurrentProjectNameHandler) Name() string {
	return "get_current_project_name"
}

func (h *GetCurrentProjectNameHandler) Description() string {
	return "获取当前用户选择的项目名称及详细信息"
}

func (h *GetCurrentProjectNameHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
		"required":   []interface{}{},
	}
}

func (h *GetCurrentProjectNameHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	// Step 1: Get current user info
	userData, err := h.client.Get(ctx, "/api/v1/auth/me", nil)
	if err != nil {
		return tools.NewErrorResult("failed to get current user: " + err.Error()), nil
	}

	var userResponse map[string]interface{}
	if err := json.Unmarshal(userData, &userResponse); err != nil {
		return tools.NewErrorResult("failed to parse user info: " + err.Error()), nil
	}

	// Check user response code
	userCode := int(0)
	if codeVal, ok := userResponse["code"].(float64); ok {
		userCode = int(codeVal)
	}
	if userCode != 0 {
		return tools.NewErrorResult("failed to get user info"), nil
	}

	// Step 2: Get current project ID
	projectIDData, err := h.client.Get(ctx, "/api/v1/profile/current-project", nil)
	if err != nil {
		return tools.NewErrorResult("failed to get current project id: " + err.Error()), nil
	}

	var projectIDResponse map[string]interface{}
	if err := json.Unmarshal(projectIDData, &projectIDResponse); err != nil {
		return tools.NewErrorResult("failed to parse project id response: " + err.Error()), nil
	}

	// Check project ID response code and verify data exists
	projectCode := int(0)
	if codeVal, ok := projectIDResponse["code"].(float64); ok {
		projectCode = int(codeVal)
	}
	if projectCode != 0 {
		// Return the actual error message from the API response if available
		if msg, ok := projectIDResponse["message"].(string); ok {
			return tools.NewErrorResult("failed to get project id: " + msg), nil
		}
		return tools.NewErrorResult("failed to get project id"), nil
	}

	// Extract project ID from data field - handle missing data or null values
	var projectID uint = 0
	if dataObj, ok := projectIDResponse["data"].(map[string]interface{}); ok {
		if pidVal, ok := dataObj["project_id"].(float64); ok {
			projectID = uint(pidVal)
		}
	}

	// If no current project is selected, return empty project info
	if projectID == 0 {
		return tools.NewJSONResult(`{"project_id":0,"name":"","message":"no current project selected"}`), nil
	}

	// Step 3: Get user projects list
	projectsData, err := h.client.Get(ctx, "/api/v1/projects", nil)
	if err != nil {
		return tools.NewErrorResult("failed to get projects list: " + err.Error()), nil
	}

	var projectsResponse map[string]interface{}
	if err := json.Unmarshal(projectsData, &projectsResponse); err != nil {
		return tools.NewErrorResult("failed to parse projects response: " + err.Error()), nil
	}

	// Check projects response code
	projCode := int(0)
	if codeVal, ok := projectsResponse["code"].(float64); ok {
		projCode = int(codeVal)
	}
	if projCode != 0 {
		// Return error with message from API if available
		if msg, ok := projectsResponse["message"].(string); ok {
			return tools.NewErrorResult("failed to get projects list: " + msg), nil
		}
		return tools.NewErrorResult("failed to get projects list"), nil
	}

	// Find the current project in the list
	var currentProject map[string]interface{}
	if dataArr, ok := projectsResponse["data"].([]interface{}); ok {
		for _, item := range dataArr {
			if proj, ok := item.(map[string]interface{}); ok {
				if pid, ok := proj["id"].(float64); ok {
					if uint(pid) == projectID {
						currentProject = proj
						break
					}
				}
			}
		}
	}

	if currentProject == nil {
		// If not found in list, return a special message indicating the project exists but user is not a member
		return tools.NewJSONResult(fmt.Sprintf(`{"project_id":%d,"name":"未知项目","message":"项目存在但当前用户不是该项目的成员"}`, projectID)), nil
	}

	// Return the current project details
	projBytes, _ := json.Marshal(currentProject)
	return tools.NewJSONResult(string(projBytes)), nil
}
