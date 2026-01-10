package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// ListDefectsHandler handles listing defects.
type ListDefectsHandler struct {
	*BaseHandler
}

func NewListDefectsHandler(c *client.BackendClient) *ListDefectsHandler {
	return &ListDefectsHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListDefectsHandler) Name() string {
	return "list_defects"
}

func (h *ListDefectsHandler) Description() string {
	return "获取项目的缺陷列表"
}

func (h *ListDefectsHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"description": "缺陷状态过滤（可选）",
			},
			"severity": map[string]interface{}{
				"type":        "string",
				"description": "严重程度过滤（可选）",
			},
		},
		"required": []interface{}{"project_id"},
	}
}

func (h *ListDefectsHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/defects", projectID)
	params := map[string]string{
		"size": "99999",
	}

	if status := GetOptionalString(args, "status", ""); status != "" {
		params["status"] = status
	}
	if severity := GetOptionalString(args, "severity", ""); severity != "" {
		params["severity"] = severity
	}

	data, err := h.client.Get(ctx, path, params)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// UpdateDefectHandler handles updating a defect or batch updating defects.
type UpdateDefectHandler struct {
	*BaseHandler
}

func NewUpdateDefectHandler(c *client.BackendClient) *UpdateDefectHandler {
	return &UpdateDefectHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *UpdateDefectHandler) Name() string {
	return "update_defects"
}

func (h *UpdateDefectHandler) Description() string {
	return "更新单个或批量更新缺陷信息，支持指定缺陷ID或批量更新多个缺陷"
}

func (h *UpdateDefectHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "单个缺陷ID（与defects参数二选一），可以是数字（如30）或格式化的字符串（如'000030'）",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"description": "缺陷状态（在使用id参数时可选）",
			},
			"severity": map[string]interface{}{
				"type":        "string",
				"description": "严重程度（在使用id参数时可选）",
			},
			"assignee": map[string]interface{}{
				"type":        "string",
				"description": "指派人（在使用id参数时可选）",
			},
			"comment": map[string]interface{}{
				"type":        "string",
				"description": "备注（在使用id参数时可选）",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "详细描述，支持多行文本，包含实际结果、测试步骤、期望结果等信息（在使用id参数时可选）",
			},
			"defects": map[string]interface{}{
				"type":        "array",
				"description": "批量更新的缺陷数组（与id参数二选一），每个对象需包含id和要更新的字段",
				"items": map[string]interface{}{
					"type": "object",
				},
			},
			"continue_on_error": map[string]interface{}{
				"type":        "boolean",
				"description": "批量更新时，失败是否继续处理（默认: true）",
			},
		},
		"required": []interface{}{"project_id"},
	}
}

func (h *UpdateDefectHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 判断是单个更新还是批量更新
	idVal, hasID := args["id"]
	defectsVal, hasDefects := args["defects"]

	// 如果既有id又有defects，则优先使用defects进行批量更新
	if hasDefects {
		return h.executeBatchUpdate(ctx, projectID, defectsVal, args)
	}

	// 如果只有id，执行单个更新
	if hasID {
		return h.executeSingleUpdate(ctx, projectID, idVal, args)
	}

	return tools.NewErrorResult("必须提供 'id' 参数（单个更新）或 'defects' 参数（批量更新）"), nil
}

// executeSingleUpdate 执行单个缺陷更新
func (h *UpdateDefectHandler) executeSingleUpdate(ctx context.Context, projectID int, idVal interface{}, args map[string]interface{}) (tools.ToolResult, error) {
	idInt, err := toInt(idVal)
	if err != nil {
		return tools.NewErrorResult(fmt.Sprintf("无效的缺陷ID: %v", err)), nil
	}
	idStr := fmt.Sprintf("%06d", idInt)

	body := make(map[string]interface{})

	if status := GetOptionalString(args, "status", ""); status != "" {
		body["status"] = status
	}
	if severity := GetOptionalString(args, "severity", ""); severity != "" {
		body["severity"] = severity
	}
	if assignee := GetOptionalString(args, "assignee", ""); assignee != "" {
		body["assignee"] = assignee
	}
	if comment := GetOptionalString(args, "comment", ""); comment != "" {
		body["comment"] = comment
	}
	if description := GetOptionalString(args, "description", ""); description != "" {
		body["description"] = description
	}

	if len(body) == 0 {
		return tools.NewErrorResult("至少需要提供一个字段来更新"), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/defects/%s", projectID, idStr)
	data, err := h.client.Put(ctx, path, body)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// executeBatchUpdate 执行批量缺陷更新
func (h *UpdateDefectHandler) executeBatchUpdate(ctx context.Context, projectID int, defectsVal interface{}, args map[string]interface{}) (tools.ToolResult, error) {
	defectsInterface, ok := defectsVal.([]interface{})
	if !ok {
		return tools.NewErrorResult("defects 必须是一个数组"), nil
	}

	if len(defectsInterface) == 0 {
		return tools.NewErrorResult("defects 数组不能为空"), nil
	}

	continueOnError := true
	if val, ok := args["continue_on_error"].(bool); ok {
		continueOnError = val
	}

	results := []map[string]interface{}{}
	successCount := 0
	failedCount := 0

	for idx, defectItem := range defectsInterface {
		defectData, ok := defectItem.(map[string]interface{})
		if !ok {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  "缺陷数据必须是对象",
			})
			if !continueOnError {
				break
			}
			continue
		}

		// 获取缺陷ID
		var defectID int
		if idFloat, ok := defectData["id"].(float64); ok {
			defectID = int(idFloat)
		} else if idInt, ok := defectData["id"].(int); ok {
			defectID = idInt
		} else {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  "缺陷对象必须包含 'id' 字段（整数）",
			})
			if !continueOnError {
				break
			}
			continue
		}

		// 准备更新数据（移除id字段）
		updateData := make(map[string]interface{})
		for k, v := range defectData {
			if k != "id" {
				updateData[k] = v
			}
		}

		if len(updateData) == 0 {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":     idx,
				"defect_id": defectID,
				"status":    "failed",
				"error":     "至少需要提供一个要更新的字段",
			})
			if !continueOnError {
				break
			}
			continue
		}

		// 调用单个更新API
		idStr := fmt.Sprintf("%06d", defectID)
		updatePath := fmt.Sprintf("/api/v1/projects/%d/defects/%s", projectID, idStr)
		_, err := h.client.Put(ctx, updatePath, updateData)
		if err != nil {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":     idx,
				"defect_id": defectID,
				"status":    "failed",
				"error":     err.Error(),
			})
			if !continueOnError {
				break
			}
			continue
		}

		successCount++
		results = append(results, map[string]interface{}{
			"index":     idx,
			"defect_id": defectID,
			"status":    "success",
		})
	}

	response := map[string]interface{}{
		"success": successCount,
		"failed":  failedCount,
		"results": results,
	}

	responseJSON, _ := json.Marshal(response)
	return tools.NewJSONResult(string(responseJSON)), nil
}
