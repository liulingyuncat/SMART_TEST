package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// CreateReviewItemHandler handles creating a case review item and optionally writing content.
type CreateReviewItemHandler struct {
	*BaseHandler
}

func NewCreateReviewItemHandler(c *client.BackendClient) *CreateReviewItemHandler {
	return &CreateReviewItemHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *CreateReviewItemHandler) Name() string {
	return "create_review_item"
}

func (h *CreateReviewItemHandler) Description() string {
	return "创建用例评审文档并写入评审内容"
}

func (h *CreateReviewItemHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "审阅条目名称",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "评审内容（可选，Markdown格式）",
			},
		},
		"required": []interface{}{"project_id", "name"},
	}
}

func (h *CreateReviewItemHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	name, err := GetString(args, "name")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// Step 1: 创建审阅文档
	createBody := map[string]interface{}{
		"name": name,
	}

	createPath := fmt.Sprintf("/api/v1/projects/%d/review-items", projectID)
	createData, err := h.client.Post(ctx, createPath, createBody)
	if err != nil {
		return tools.NewErrorResult(fmt.Sprintf("创建审阅文档失败: %s", err.Error())), nil
	}

	// 解析创建响应获取 ID
	var createResp struct {
		ID uint `json:"id"`
	}
	if err := json.Unmarshal(createData, &createResp); err != nil {
		return tools.NewErrorResult(fmt.Sprintf("解析创建响应失败: %s", err.Error())), nil
	}

	// Step 2: 如果提供了 content，更新文档内容
	content := GetOptionalString(args, "content", "")
	if content != "" {
		updateBody := map[string]interface{}{
			"content": content,
		}

		updatePath := fmt.Sprintf("/api/v1/projects/%d/review-items/%d", projectID, createResp.ID)
		updateData, err := h.client.Put(ctx, updatePath, updateBody)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("创建成功但写入内容失败: %s", err.Error())), nil
		}

		return tools.NewJSONResult(string(updateData)), nil
	}

	return tools.NewJSONResult(string(createData)), nil
}
