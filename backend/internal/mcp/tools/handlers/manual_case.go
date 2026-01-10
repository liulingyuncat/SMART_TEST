package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// ListManualGroupsHandler handles listing manual case groups.
type ListManualGroupsHandler struct {
	*BaseHandler
}

func NewListManualGroupsHandler(c *client.BackendClient) *ListManualGroupsHandler {
	return &ListManualGroupsHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListManualGroupsHandler) Name() string {
	return "list_manual_groups"
}

func (h *ListManualGroupsHandler) Description() string {
	return "获取项目的手工测试用例集列表"
}

func (h *ListManualGroupsHandler) InputSchema() map[string]interface{} {
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

func (h *ListManualGroupsHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/case-groups?case_type=overall", projectID)
	data, err := h.client.Get(ctx, path, nil)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// parseJSON is a helper function to parse JSON bytes
func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// ListManualCasesHandler handles listing manual test cases.
type ListManualCasesHandler struct {
	*BaseHandler
}

func NewListManualCasesHandler(c *client.BackendClient) *ListManualCasesHandler {
	return &ListManualCasesHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListManualCasesHandler) Name() string {
	return "list_manual_cases"
}

func (h *ListManualCasesHandler) Description() string {
	return "获取用例集中的手工测试用例列表，支持返回所有语言字段（CN、JP、EN）"
}

func (h *ListManualCasesHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "用例集ID",
			},
			"return_all_fields": map[string]interface{}{
				"type":        "boolean",
				"description": "是否返回所有字段包括CN、JP、EN (默认false)",
			},
		},
		"required": []interface{}{"project_id", "group_id"},
	}
}

func (h *ListManualCasesHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	groupID, err := GetInt(args, "group_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 获取return_all_fields参数，默认为false
	returnAllFields := false
	if val, ok := args["return_all_fields"].(bool); ok {
		returnAllFields = val
	}

	// 首先获取case_group的名称
	caseGroupPath := fmt.Sprintf("/api/v1/projects/%d/case-groups", projectID)
	groupsData, err := h.client.Get(ctx, caseGroupPath, nil)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 解析JSON获取对应的group_name
	// API 返回直接的数组 [...]，或者包装的响应 {"code": 0, "data": [...]}
	var groups []map[string]interface{}

	// 首先尝试解析为直接数组
	err = parseJSON(groupsData, &groups)
	if err != nil {
		// 如果失败，尝试解析为包装的响应
		var response map[string]interface{}
		err = parseJSON(groupsData, &response)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("failed to parse groups response: %v", err)), nil
		}

		// 获取data字段中的groups列表
		if dataVal, ok := response["data"]; ok {
			// 尝试将data解析为数组
			if dataArray, ok := dataVal.([]interface{}); ok {
				for _, item := range dataArray {
					if itemMap, ok := item.(map[string]interface{}); ok {
						groups = append(groups, itemMap)
					}
				}
			}
		}
	}

	var groupName string
	for _, g := range groups {
		if id, ok := g["id"].(float64); ok && int(id) == groupID {
			if name, ok := g["group_name"].(string); ok {
				groupName = name
				break
			}
		}
	}

	if groupName == "" {
		return tools.NewErrorResult(fmt.Sprintf("case group with id %d not found", groupID)), nil
	}

	// 使用group_name作为query参数过滤指定case-group的用例
	path := fmt.Sprintf("/api/v1/projects/%d/manual-cases", projectID)
	params := map[string]string{
		"case_group": groupName,
		"size":       "99999",
	}

	// 如果需要返回所有字段，添加额外的查询参数
	if returnAllFields {
		params["return_all_fields"] = "true"
	}

	data, err := h.client.Get(ctx, path, params)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// CreateCaseGroupHandler handles creating a case group.
type CreateCaseGroupHandler struct {
	*BaseHandler
}

func NewCreateCaseGroupHandler(c *client.BackendClient) *CreateCaseGroupHandler {
	return &CreateCaseGroupHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *CreateCaseGroupHandler) Name() string {
	return "create_case_group"
}

func (h *CreateCaseGroupHandler) Description() string {
	return "创建测试用例集"
}

func (h *CreateCaseGroupHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "用例集名称",
			},
			"case_type": map[string]interface{}{
				"type":        "string",
				"description": "用例类型（默认: overall）",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "用例集描述（可选）",
			},
		},
		"required": []interface{}{"project_id", "name"},
	}
}

func (h *CreateCaseGroupHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	name, err := GetString(args, "name")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// case_type 有默认值 "overall"
	caseType := GetOptionalString(args, "case_type", "overall")

	body := map[string]interface{}{
		"group_name": name,
		"case_type":  caseType,
	}

	if desc := GetOptionalString(args, "description", ""); desc != "" {
		body["description"] = desc
	}

	path := fmt.Sprintf("/api/v1/projects/%d/case-groups", projectID)
	data, err := h.client.Post(ctx, path, body)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// CreateManualCasesHandler handles batch creating manual test cases.
type CreateManualCasesHandler struct {
	*BaseHandler
}

func NewCreateManualCasesHandler(c *client.BackendClient) *CreateManualCasesHandler {
	return &CreateManualCasesHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *CreateManualCasesHandler) Name() string {
	return "create_manual_cases"
}

func (h *CreateManualCasesHandler) Description() string {
	return "批量创建手工测试用例"
}

func (h *CreateManualCasesHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "用例集ID",
			},
			"cases": map[string]interface{}{
				"type":        "array",
				"description": "用例数据数组，每个元素包含各语言字段",
				"items": map[string]interface{}{
					"type": "object",
				},
			},
			"continue_on_error": map[string]interface{}{
				"type":        "boolean",
				"description": "失败是否继续处理（默认: true）",
			},
		},
		"required": []interface{}{"project_id", "group_id", "cases"},
	}
}

func (h *CreateManualCasesHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	groupID, err := GetInt(args, "group_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	casesInterface, ok := args["cases"].([]interface{})
	if !ok {
		return tools.NewErrorResult("cases must be an array"), nil
	}

	continueOnError := true
	if val, ok := args["continue_on_error"].(bool); ok {
		continueOnError = val
	}

	results := []map[string]interface{}{}
	successCount := 0
	failedCount := 0

	path := fmt.Sprintf("/api/v1/projects/%d/case-groups/%d/manual-cases", projectID, groupID)

	for idx, caseItem := range casesInterface {
		data, ok := caseItem.(map[string]interface{})
		if !ok {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  "case item must be an object",
			})
			if !continueOnError {
				break
			}
			continue
		}

		resp, err := h.client.Post(ctx, path, data)
		if err != nil {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  err.Error(),
			})
			if !continueOnError {
				break
			}
			continue
		}

		// 检查是否接收到HTML错误页面（可能是路由错误）
		respStr := string(resp)
		if strings.Contains(respStr, "<!doctype") || strings.Contains(respStr, "<html") {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  "API返回HTML错误页面，可能是路由配置错误或后端服务问题",
			})
			if !continueOnError {
				break
			}
			continue
		}

		// 解析响应以获取创建的用例ID
		var respData map[string]interface{}
		err = json.Unmarshal(resp, &respData)
		if err == nil {
			successCount++
			// 如果响应包含data字段（如{"code":0,"data":{...}}）
			if dataVal, ok := respData["data"].(map[string]interface{}); ok {
				if id, ok := dataVal["id"].(float64); ok {
					results = append(results, map[string]interface{}{
						"index":   idx,
						"status":  "success",
						"case_id": int(id),
					})
				} else {
					results = append(results, map[string]interface{}{
						"index":  idx,
						"status": "success",
						"data":   dataVal,
					})
				}
			} else if id, ok := respData["id"].(float64); ok {
				// 如果直接返回用例对象
				results = append(results, map[string]interface{}{
					"index":   idx,
					"status":  "success",
					"case_id": int(id),
				})
			} else {
				results = append(results, map[string]interface{}{
					"index":  idx,
					"status": "success",
					"data":   respData,
				})
			}
		} else {
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "success",
				"data":   string(resp),
			})
		}
	}

	return tools.NewJSONResult(fmt.Sprintf(`{
		"success": %d,
		"failed": %d,
		"results": %s
	}`, successCount, failedCount, func() string {
		resultJSON, _ := json.Marshal(results)
		return string(resultJSON)
	}())), nil
}

// UpdateManualCasesHandler handles batch updating manual test cases.
type UpdateManualCasesHandler struct {
	*BaseHandler
}

func NewUpdateManualCasesHandler(c *client.BackendClient) *UpdateManualCasesHandler {
	return &UpdateManualCasesHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *UpdateManualCasesHandler) Name() string {
	return "update_manual_cases"
}

func (h *UpdateManualCasesHandler) Description() string {
	return "批量更新手工测试用例，支持所有字段（除UUID外）。可选参数group_id用于更新特定用例集中的用例。支持的字段：case_number、case_group、major_function_cn/jp/en、middle_function_cn/jp/en、minor_function_cn/jp/en、precondition_cn/jp/en、test_steps_cn/jp/en、expected_result_cn/jp/en、test_result、remark等"
}

func (h *UpdateManualCasesHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "用例集ID（可选但推荐，如果提供则使用用例集API以确保正确性）",
			},
			"cases": map[string]interface{}{
				"type":        "array",
				"description": "用例数据数组，每个元素必须包含id字段（整数型），其他字段可选。支持的字段包括：case_number(字符串)、case_group、major_function_cn/jp/en、middle_function_cn/jp/en、minor_function_cn/jp/en、precondition_cn/jp/en、test_steps_cn/jp/en、expected_result_cn/jp/en、test_result、remark等",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "integer",
							"description": "用例ID（必填）",
						},
						"case_number": map[string]interface{}{
							"type":        "string",
							"description": "用例编号（可选）",
						},
					},
					"required": []interface{}{"id"},
				},
			},
			"continue_on_error": map[string]interface{}{
				"type":        "boolean",
				"description": "失败是否继续处理（默认: true）",
			},
		},
		"required": []interface{}{"project_id", "cases"},
	}
}

func (h *UpdateManualCasesHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 尝试获取group_id（可选）
	groupID := 0
	if gid, ok := args["group_id"]; ok {
		if gidInt, err := GetInt(args, "group_id"); err == nil {
			groupID = gidInt
		} else if gidFloat, ok := gid.(float64); ok {
			groupID = int(gidFloat)
		}
	}

	casesInterface, ok := args["cases"].([]interface{})
	if !ok {
		return tools.NewErrorResult("cases must be an array"), nil
	}

	continueOnError := true
	if val, ok := args["continue_on_error"].(bool); ok {
		continueOnError = val
	}

	results := []map[string]interface{}{}
	successCount := 0
	failedCount := 0

	for idx, caseItem := range casesInterface {
		data, ok := caseItem.(map[string]interface{})
		if !ok {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  "case item must be an object",
			})
			if !continueOnError {
				break
			}
			continue
		}

		// 获取case id
		idInterface, ok := data["id"]
		if !ok {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  "id field is required",
			})
			if !continueOnError {
				break
			}
			continue
		}

		var caseID int
		switch v := idInterface.(type) {
		case float64:
			caseID = int(v)
		case int:
			caseID = v
		default:
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  "id must be an integer",
			})
			if !continueOnError {
				break
			}
			continue
		}

		// 创建更新数据（不包含id，包含所有其他字段）
		updateData := make(map[string]interface{})
		for k, v := range data {
			if k != "id" {
				// 包括所有字段，包括 case_number 和所有其他字段
				// v 可能是 nil、""、0 等，我们需要包括它们以支持清空字段操作
				updateData[k] = v
			}
		}

		// 根据是否有group_id来选择API路径
		var path string
		if groupID > 0 {
			path = fmt.Sprintf("/api/v1/projects/%d/case-groups/%d/manual-cases/%d", projectID, groupID, caseID)
		} else {
			path = fmt.Sprintf("/api/v1/projects/%d/manual-cases/%d", projectID, caseID)
		}

		_, err := h.client.Put(ctx, path, updateData)
		if err != nil {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":   idx,
				"case_id": caseID,
				"status":  "failed",
				"error":   err.Error(),
			})
			if !continueOnError {
				break
			}
			continue
		}

		successCount++
		results = append(results, map[string]interface{}{
			"index":   idx,
			"case_id": caseID,
			"status":  "success",
		})
	}

	return tools.NewJSONResult(fmt.Sprintf(`{
		"success": %d,
		"failed": %d,
		"results": %s
	}`, successCount, failedCount, func() string {
		resultJSON, _ := json.Marshal(results)
		return string(resultJSON)
	}())), nil
}

// UpdateManualCaseHandler handles updating a manual test case.
type UpdateManualCaseHandler struct {
	*BaseHandler
}

func NewUpdateManualCaseHandler(c *client.BackendClient) *UpdateManualCaseHandler {
	return &UpdateManualCaseHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *UpdateManualCaseHandler) Name() string {
	return "update_manual_case"
}

func (h *UpdateManualCaseHandler) Description() string {
	return "更新手工测试用例"
}

func (h *UpdateManualCaseHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "用例ID",
			},
			"data": map[string]interface{}{
				"type":        "object",
				"description": "要更新的用例数据",
			},
		},
		"required": []interface{}{"project_id", "id", "data"},
	}
}

func (h *UpdateManualCaseHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	id, err := GetInt(args, "id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	data, ok := args["data"].(map[string]interface{})
	if !ok {
		return tools.NewErrorResult("data must be an object"), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/manual-cases/%d", projectID, id)
	resp, err := h.client.Put(ctx, path, data)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(resp)), nil
}
