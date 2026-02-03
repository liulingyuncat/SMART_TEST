package handlers

import (
	"context"
	"encoding/json"
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
	return "获取项目中的AI观点文档列表，包含每个观点的章节(chunks)摘要信息"
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
	return "获取单个AI观点文档的详细内容，包含所有章节(chunks)的完整内容。支持通过ID或名称查询"
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
				"description": "观点文档ID（与name二选一）",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "观点文档名称（与id二选一）",
			},
		},
		"required": []interface{}{"project_id"},
	}
}

func (h *GetViewpointItemHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
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
		path = fmt.Sprintf("/api/v1/projects/%d/viewpoint-items/%d", projectID, idInt)
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

		listPath := fmt.Sprintf("/api/v1/projects/%d/viewpoint-items", projectID)
		listData, err := h.client.Get(ctx, listPath, nil)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("获取观点列表失败: %s", err.Error())), nil
		}

		// 解析列表数据
		var listResp struct {
			Code    int                      `json:"code"`
			Message string                   `json:"message"`
			Data    []map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(listData, &listResp); err != nil {
			return tools.NewErrorResult(fmt.Sprintf("解析观点列表失败: %s", err.Error())), nil
		}

		// 查找匹配的观点
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
			return tools.NewErrorResult(fmt.Sprintf("未找到名称为 '%s' 的观点文档", nameStr)), nil
		}

		// 通过找到的 ID 获取详细内容
		path = fmt.Sprintf("/api/v1/projects/%d/viewpoint-items/%d", projectID, foundID)
		data, err = h.client.Get(ctx, path, nil)
		if err != nil {
			return tools.NewErrorResult(err.Error()), nil
		}
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
	return "创建AI观点文档，可选同时创建多个章节(chunks)"
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

	// 处理 chunks 参数
	if chunks, ok := args["chunks"]; ok && chunks != nil {
		body["chunks"] = chunks
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
	return "更新AI观点文档，可同时对章节(chunks)进行增删改操作"
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
	// 处理 chunks 参数
	if chunks, ok := args["chunks"]; ok && chunks != nil {
		body["chunks"] = chunks
	}

	path := fmt.Sprintf("/api/v1/projects/%d/viewpoint-items/%d", projectID, id)
	data, err := h.client.Put(ctx, path, body)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}
