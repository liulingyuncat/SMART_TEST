package handlers

import (
	"context"
	"fmt"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// ListViewpointItemsHandler handles listing viewpoint items.
type ListViewpointItemsHandler struct {
	*BaseHandler
}

func NewListViewpointItemsHandler(c *client.BackendClient) *ListViewpointItemsHandler {
	return &ListViewpointItemsHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListViewpointItemsHandler) Name() string {
	return "list_viewpoint_items"
}

func (h *ListViewpointItemsHandler) Description() string {
	return "获取项目中的AI观点文档列表"
}

func (h *ListViewpointItemsHandler) InputSchema() map[string]interface{} {
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

func (h *ListViewpointItemsHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/viewpoint-items", projectID)
	data, err := h.client.Get(ctx, path, nil)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// GetViewpointItemHandler handles getting a single viewpoint item.
type GetViewpointItemHandler struct {
	*BaseHandler
}

func NewGetViewpointItemHandler(c *client.BackendClient) *GetViewpointItemHandler {
	return &GetViewpointItemHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *GetViewpointItemHandler) Name() string {
	return "get_viewpoint_item"
}

func (h *GetViewpointItemHandler) Description() string {
	return "获取单个AI观点文档的详细内容"
}

func (h *GetViewpointItemHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "观点文档ID",
			},
		},
		"required": []interface{}{"project_id", "id"},
	}
}

func (h *GetViewpointItemHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	id, err := GetInt(args, "id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/viewpoint-items/%d", projectID, id)
	data, err := h.client.Get(ctx, path, nil)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// CreateViewpointItemHandler handles creating a viewpoint item.
type CreateViewpointItemHandler struct {
	*BaseHandler
}

func NewCreateViewpointItemHandler(c *client.BackendClient) *CreateViewpointItemHandler {
	return &CreateViewpointItemHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *CreateViewpointItemHandler) Name() string {
	return "create_viewpoint_item"
}

func (h *CreateViewpointItemHandler) Description() string {
	return "创建AI观点文档"
}

func (h *CreateViewpointItemHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "观点名称",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "观点内容",
			},
			"requirement_id": map[string]interface{}{
				"type":        "integer",
				"description": "关联的需求ID（可选）",
			},
		},
		"required": []interface{}{"project_id", "name", "content"},
	}
}

func (h *CreateViewpointItemHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	name, err := GetString(args, "name")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	content, err := GetString(args, "content")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	body := map[string]interface{}{
		"name":    name,
		"content": content,
	}

	if reqID := GetOptionalInt(args, "requirement_id", 0); reqID > 0 {
		body["requirement_id"] = reqID
	}

	path := fmt.Sprintf("/api/v1/projects/%d/viewpoint-items", projectID)
	data, err := h.client.Post(ctx, path, body)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// UpdateViewpointItemHandler handles updating a viewpoint item.
type UpdateViewpointItemHandler struct {
	*BaseHandler
}

func NewUpdateViewpointItemHandler(c *client.BackendClient) *UpdateViewpointItemHandler {
	return &UpdateViewpointItemHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *UpdateViewpointItemHandler) Name() string {
	return "update_viewpoint_item"
}

func (h *UpdateViewpointItemHandler) Description() string {
	return "更新AI观点文档"
}

func (h *UpdateViewpointItemHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "观点文档ID",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "观点名称（可选）",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "观点内容（可选）",
			},
		},
		"required": []interface{}{"project_id", "id"},
	}
}

func (h *UpdateViewpointItemHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	id, err := GetInt(args, "id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	body := make(map[string]interface{})
	if name := GetOptionalString(args, "name", ""); name != "" {
		body["name"] = name
	}
	if content := GetOptionalString(args, "content", ""); content != "" {
		body["content"] = content
	}

	path := fmt.Sprintf("/api/v1/projects/%d/viewpoint-items/%d", projectID, id)
	data, err := h.client.Put(ctx, path, body)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}
