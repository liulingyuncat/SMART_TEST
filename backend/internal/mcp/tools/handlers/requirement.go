package handlers

import (
	"context"
	"fmt"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// ListRequirementItemsHandler handles listing requirement items.
type ListRequirementItemsHandler struct {
	*BaseHandler
}

func NewListRequirementItemsHandler(c *client.BackendClient) *ListRequirementItemsHandler {
	return &ListRequirementItemsHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListRequirementItemsHandler) Name() string {
	return "list_requirement_items"
}

func (h *ListRequirementItemsHandler) Description() string {
	return "获取项目中的AI需求文档列表"
}

func (h *ListRequirementItemsHandler) InputSchema() map[string]interface{} {
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

func (h *ListRequirementItemsHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/requirement-items", projectID)
	data, err := h.client.Get(ctx, path, nil)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// GetRequirementItemHandler handles getting a single requirement item.
type GetRequirementItemHandler struct {
	*BaseHandler
}

func NewGetRequirementItemHandler(c *client.BackendClient) *GetRequirementItemHandler {
	return &GetRequirementItemHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *GetRequirementItemHandler) Name() string {
	return "get_requirement_item"
}

func (h *GetRequirementItemHandler) Description() string {
	return "获取单个AI需求文档的详细内容"
}

func (h *GetRequirementItemHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "需求文档ID",
			},
		},
		"required": []interface{}{"project_id", "id"},
	}
}

func (h *GetRequirementItemHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	id, err := GetInt(args, "id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/requirement-items/%d", projectID, id)
	data, err := h.client.Get(ctx, path, nil)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// CreateRequirementItemHandler handles creating a requirement item.
type CreateRequirementItemHandler struct {
	*BaseHandler
}

func NewCreateRequirementItemHandler(c *client.BackendClient) *CreateRequirementItemHandler {
	return &CreateRequirementItemHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *CreateRequirementItemHandler) Name() string {
	return "create_requirement_item"
}

func (h *CreateRequirementItemHandler) Description() string {
	return "创建AI需求文档"
}

func (h *CreateRequirementItemHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "需求名称",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "需求内容",
			},
			"parent_id": map[string]interface{}{
				"type":        "integer",
				"description": "父需求ID（可选）",
			},
		},
		"required": []interface{}{"project_id", "name", "content"},
	}
}

func (h *CreateRequirementItemHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
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

	if parentID := GetOptionalInt(args, "parent_id", 0); parentID > 0 {
		body["parent_id"] = parentID
	}

	path := fmt.Sprintf("/api/v1/projects/%d/requirement-items", projectID)
	data, err := h.client.Post(ctx, path, body)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// UpdateRequirementItemHandler handles updating a requirement item.
type UpdateRequirementItemHandler struct {
	*BaseHandler
}

func NewUpdateRequirementItemHandler(c *client.BackendClient) *UpdateRequirementItemHandler {
	return &UpdateRequirementItemHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *UpdateRequirementItemHandler) Name() string {
	return "update_requirement_item"
}

func (h *UpdateRequirementItemHandler) Description() string {
	return "更新AI需求文档"
}

func (h *UpdateRequirementItemHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "需求文档ID",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "需求名称（可选）",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "需求内容（可选）",
			},
		},
		"required": []interface{}{"project_id", "id"},
	}
}

func (h *UpdateRequirementItemHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
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

	path := fmt.Sprintf("/api/v1/projects/%d/requirement-items/%d", projectID, id)
	data, err := h.client.Put(ctx, path, body)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}
