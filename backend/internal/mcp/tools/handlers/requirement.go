package handlers

import (
	"context"
	"encoding/json"
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
	return "获取项目中的AI需求文档列表，包含每个需求的章节(chunks)摘要信息"
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
	return "获取单个AI需求文档的详细内容，包含所有章节(chunks)的完整内容。支持通过ID或名称查询"
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
				"description": "需求文档ID（与name二选一）",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "需求文档名称（与id二选一）",
			},
		},
		"required": []interface{}{"project_id"},
	}
}

func (h *GetRequirementItemHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 优先使用 ID，如果没有 ID 则使用 name
	_, hasID := args["id"]
	_, hasName := args["name"]

	if !hasID && !hasName {
		return tools.NewErrorResult("必须提供 id 或 name 参数之一"), nil
	}

	var path string
	var data []byte

	if hasID {
		// 通过 ID 查询
		idInt, err := GetInt(args, "id")
		if err != nil {
			return tools.NewErrorResult(err.Error()), nil
		}
		path = fmt.Sprintf("/api/v1/projects/%d/requirement-items/%d", projectID, idInt)
		data, err = h.client.Get(ctx, path, nil)
		if err != nil {
			return tools.NewErrorResult(err.Error()), nil
		}
	} else {
		// 通过名称查询：先获取列表，再匹配名称
		nameStr, ok := args["name"].(string)
		if !ok {
			return tools.NewErrorResult("name 参数必须是字符串"), nil
		}

		listPath := fmt.Sprintf("/api/v1/projects/%d/requirement-items", projectID)
		listData, err := h.client.Get(ctx, listPath, nil)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("获取需求列表失败: %s", err.Error())), nil
		}

		// 解析列表数据
		var listResp struct {
			Code    int                      `json:"code"`
			Message string                   `json:"message"`
			Data    []map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(listData, &listResp); err != nil {
			return tools.NewErrorResult(fmt.Sprintf("解析需求列表失败: %s", err.Error())), nil
		}

		// 查找匹配的需求
		var foundID int
		for _, item := range listResp.Data {
			if itemName, ok := item["name"].(string); ok && itemName == nameStr {
				if itemID, ok := item["id"].(float64); ok {
					foundID = int(itemID)
					break
				}
			}
		}

		if foundID == 0 {
			return tools.NewErrorResult(fmt.Sprintf("未找到名称为 '%s' 的需求文档", nameStr)), nil
		}

		// 通过找到的 ID 获取详细内容
		path = fmt.Sprintf("/api/v1/projects/%d/requirement-items/%d", projectID, foundID)
		data, err = h.client.Get(ctx, path, nil)
		if err != nil {
			return tools.NewErrorResult(err.Error()), nil
		}
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
	return "创建AI需求文档，可选同时创建多个章节(chunks)"
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
			"chunks": map[string]interface{}{
				"type":        "array",
				"description": "章节数组（可选），按顺序创建",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type":        "string",
							"description": "章节标题",
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "章节内容",
						},
					},
					"required": []interface{}{"title", "content"},
				},
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

	// 处理 chunks 参数
	if chunks, ok := args["chunks"]; ok && chunks != nil {
		body["chunks"] = chunks
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
	return "更新AI需求文档，可同时对章节(chunks)进行增删改操作"
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
			"chunks": map[string]interface{}{
				"type":        "array",
				"description": "章节操作数组（可选），支持增加/更新/删除",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"chunk_id": map[string]interface{}{
							"type":        "integer",
							"description": "章节ID（更新/删除时必填）",
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "章节标题",
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "章节内容",
						},
						"_delete": map[string]interface{}{
							"type":        "boolean",
							"description": "是否删除此章节",
						},
					},
				},
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
	// 处理 chunks 参数
	if chunks, ok := args["chunks"]; ok && chunks != nil {
		body["chunks"] = chunks
	}

	path := fmt.Sprintf("/api/v1/projects/%d/requirement-items/%d", projectID, id)
	data, err := h.client.Put(ctx, path, body)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}
