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
	return "批量创建手工测试用例，必须指定用例集名称（group_name）。如果用例集不存在则自动创建。只生成中文字段的用例"
}

func (h *CreateManualCasesHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"group_name": map[string]interface{}{
				"type":        "string",
				"description": "手工用例集名称（必填），会自动查找对应的group_id，如果不存在则自动创建",
			},
			"cases": map[string]interface{}{
				"type":        "array",
				"description": "用例数据数组，包含中文字段：case_number, major_function_cn, middle_function_cn, minor_function_cn, precondition_cn, test_steps_cn, expected_result_cn。case_type字段可选，默认为overall",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"case_type": map[string]interface{}{
							"type":        "string",
							"enum":        []interface{}{"overall", "change", "ai", "acceptance"},
							"description": "用例类型（可选，默认overall）：overall-整体用例, change-变更用例, ai-AI用例, acceptance-验收用例",
						},
						"case_number": map[string]interface{}{
							"type":        "string",
							"description": "用例编号",
						},
						"major_function_cn": map[string]interface{}{
							"type":        "string",
							"description": "大功能（中文）",
						},
						"middle_function_cn": map[string]interface{}{
							"type":        "string",
							"description": "中功能（中文）",
						},
						"minor_function_cn": map[string]interface{}{
							"type":        "string",
							"description": "小功能（中文）",
						},
						"precondition_cn": map[string]interface{}{
							"type":        "string",
							"description": "前置条件（中文）",
						},
						"test_steps_cn": map[string]interface{}{
							"type":        "string",
							"description": "测试步骤（中文）",
						},
						"expected_result_cn": map[string]interface{}{
							"type":        "string",
							"description": "期待结果（中文）",
						},
					},
				},
			},
			"continue_on_error": map[string]interface{}{
				"type":        "boolean",
				"description": "失败是否继续处理（默认: true）",
			},
		},
		"required": []interface{}{"project_id", "group_name", "cases"},
	}
}

func (h *CreateManualCasesHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// group_name 是必填参数
	groupName, err := GetString(args, "group_name")
	if err != nil {
		return tools.NewErrorResult("必须提供 group_name 参数"), nil
	}

	var groupID int

	// 查找用例集
	caseGroupPath := fmt.Sprintf("/api/v1/projects/%d/case-groups?case_type=overall", projectID)
	groupsData, err := h.client.Get(ctx, caseGroupPath, nil)
	if err != nil {
		return tools.NewErrorResult(fmt.Sprintf("获取用例集列表失败: %v", err)), nil
	}

	// 解析用例集列表
	var groups []map[string]interface{}
	err = json.Unmarshal(groupsData, &groups)
	if err != nil {
		var response map[string]interface{}
		err = json.Unmarshal(groupsData, &response)
		if err == nil {
			if dataVal, ok := response["data"]; ok {
				if dataArray, ok := dataVal.([]interface{}); ok {
					for _, item := range dataArray {
						if itemMap, ok := item.(map[string]interface{}); ok {
							groups = append(groups, itemMap)
						}
					}
				}
			}
		}
	}

	// 查找匹配的用例集
	for _, g := range groups {
		if name, ok := g["group_name"].(string); ok && name == groupName {
			if id, ok := g["id"].(float64); ok {
				groupID = int(id)
				break
			}
		}
	}

	// 如果没找到，则创建新的用例集
	if groupID == 0 {
		createPath := fmt.Sprintf("/api/v1/projects/%d/case-groups", projectID)
		createBody := map[string]interface{}{
			"group_name": groupName,
			"case_type":  "overall",
		}
		createResp, err := h.client.Post(ctx, createPath, createBody)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("创建用例集失败: %v", err)), nil
		}

		var createResult map[string]interface{}
		err = json.Unmarshal(createResp, &createResult)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("解析创建用例集响应失败: %v", err)), nil
		}

		// 尝试从 data 或直接获取 id
		if dataVal, ok := createResult["data"].(map[string]interface{}); ok {
			if id, ok := dataVal["id"].(float64); ok {
				groupID = int(id)
			}
		} else if id, ok := createResult["id"].(float64); ok {
			groupID = int(id)
		}

		if groupID == 0 {
			return tools.NewErrorResult("创建用例集成功但无法获取ID"), nil
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

		// 自动添加 case_type 字段（如果没有提供则默认为 "overall"）
		if _, exists := data["case_type"]; !exists {
			data["case_type"] = "overall"
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
	return "批量更新手工测试用例。支持两种模式：1) 通过cases数组中的id字段直接更新；2) 通过filter条件（大功能/中功能/小功能）查找用例并批量更新。支持的字段：case_number、case_group、major_function_cn/jp/en、middle_function_cn/jp/en、minor_function_cn/jp/en、precondition_cn/jp/en、test_steps_cn/jp/en、expected_result_cn/jp/en、test_result、remark等"
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
				"description": "用例集ID（与group_name二选一，使用filter模式时必填其一）",
			},
			"group_name": map[string]interface{}{
				"type":        "string",
				"description": "用例集名称（与group_id二选一），会自动查找对应的group_id",
			},
			"filter": map[string]interface{}{
				"type":        "object",
				"description": "筛选条件，用于查找要更新的用例。支持：major_function_cn（大功能）、middle_function_cn（中功能）、minor_function_cn（小功能）",
				"properties": map[string]interface{}{
					"major_function_cn": map[string]interface{}{
						"type":        "string",
						"description": "大功能名称（精确匹配）",
					},
					"middle_function_cn": map[string]interface{}{
						"type":        "string",
						"description": "中功能名称（精确匹配）",
					},
					"minor_function_cn": map[string]interface{}{
						"type":        "string",
						"description": "小功能名称（精确匹配）",
					},
				},
			},
			"update_data": map[string]interface{}{
				"type":        "object",
				"description": "使用filter模式时，要更新的字段和值",
			},
			"cases": map[string]interface{}{
				"type":        "array",
				"description": "用例数据数组（直接更新模式），每个元素必须包含id字段（整数型），其他字段可选",
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
		"required": []interface{}{"project_id"},
	}
}

func (h *UpdateManualCasesHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 支持 group_id 或 group_name
	groupID := 0
	groupName := GetOptionalString(args, "group_name", "")

	if gid, ok := args["group_id"]; ok && gid != nil {
		if gidInt, err := GetInt(args, "group_id"); err == nil {
			groupID = gidInt
		} else if gidFloat, ok := gid.(float64); ok {
			groupID = int(gidFloat)
		}
	}

	// 如果没有 group_id 但有 group_name，则查找
	if groupID == 0 && groupName != "" {
		caseGroupPath := fmt.Sprintf("/api/v1/projects/%d/case-groups?case_type=overall", projectID)
		groupsData, err := h.client.Get(ctx, caseGroupPath, nil)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("获取用例集列表失败: %v", err)), nil
		}

		var groups []map[string]interface{}
		err = json.Unmarshal(groupsData, &groups)
		if err != nil {
			var response map[string]interface{}
			err = json.Unmarshal(groupsData, &response)
			if err == nil {
				if dataVal, ok := response["data"]; ok {
					if dataArray, ok := dataVal.([]interface{}); ok {
						for _, item := range dataArray {
							if itemMap, ok := item.(map[string]interface{}); ok {
								groups = append(groups, itemMap)
							}
						}
					}
				}
			}
		}

		for _, g := range groups {
			if name, ok := g["group_name"].(string); ok && name == groupName {
				if id, ok := g["id"].(float64); ok {
					groupID = int(id)
					break
				}
			}
		}

		if groupID == 0 {
			return tools.NewErrorResult(fmt.Sprintf("未找到名为 '%s' 的用例集", groupName)), nil
		}
	}

	continueOnError := true
	if val, ok := args["continue_on_error"].(bool); ok {
		continueOnError = val
	}

	results := []map[string]interface{}{}
	successCount := 0
	failedCount := 0

	// 检查是否使用 filter 模式
	if filterInterface, ok := args["filter"].(map[string]interface{}); ok {
		// Filter 模式：通过大功能/中功能/小功能查找用例并更新
		updateDataInterface, ok := args["update_data"].(map[string]interface{})
		if !ok {
			return tools.NewErrorResult("使用filter模式时必须提供update_data参数"), nil
		}

		if groupID == 0 {
			return tools.NewErrorResult("使用filter模式时必须提供group_id或group_name参数"), nil
		}

		// 获取用例集中的所有用例
		// 首先获取group_name
		caseGroupPath := fmt.Sprintf("/api/v1/projects/%d/case-groups", projectID)
		groupsData, err := h.client.Get(ctx, caseGroupPath, nil)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("获取用例集列表失败: %v", err)), nil
		}

		var groups []map[string]interface{}
		err = json.Unmarshal(groupsData, &groups)
		if err != nil {
			var response map[string]interface{}
			err = json.Unmarshal(groupsData, &response)
			if err == nil {
				if dataVal, ok := response["data"]; ok {
					if dataArray, ok := dataVal.([]interface{}); ok {
						for _, item := range dataArray {
							if itemMap, ok := item.(map[string]interface{}); ok {
								groups = append(groups, itemMap)
							}
						}
					}
				}
			}
		}

		var targetGroupName string
		for _, g := range groups {
			if id, ok := g["id"].(float64); ok && int(id) == groupID {
				if name, ok := g["group_name"].(string); ok {
					targetGroupName = name
					break
				}
			}
		}

		if targetGroupName == "" {
			return tools.NewErrorResult(fmt.Sprintf("未找到ID为 %d 的用例集", groupID)), nil
		}

		// 获取用例列表
		casesPath := fmt.Sprintf("/api/v1/projects/%d/manual-cases", projectID)
		params := map[string]string{
			"case_group":        targetGroupName,
			"size":              "99999",
			"return_all_fields": "true",
		}
		casesData, err := h.client.Get(ctx, casesPath, params)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("获取用例列表失败: %v", err)), nil
		}

		var casesResp map[string]interface{}
		err = json.Unmarshal(casesData, &casesResp)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("解析用例列表失败: %v", err)), nil
		}

		var cases []map[string]interface{}
		if dataVal, ok := casesResp["data"].(map[string]interface{}); ok {
			if casesArray, ok := dataVal["cases"].([]interface{}); ok {
				for _, c := range casesArray {
					if caseMap, ok := c.(map[string]interface{}); ok {
						cases = append(cases, caseMap)
					}
				}
			}
		}

		// 根据 filter 条件筛选用例
		majorFilter := ""
		middleFilter := ""
		minorFilter := ""

		if val, ok := filterInterface["major_function_cn"].(string); ok {
			majorFilter = val
		}
		if val, ok := filterInterface["middle_function_cn"].(string); ok {
			middleFilter = val
		}
		if val, ok := filterInterface["minor_function_cn"].(string); ok {
			minorFilter = val
		}

		matchedCases := []map[string]interface{}{}
		for _, c := range cases {
			match := true

			if majorFilter != "" {
				if val, ok := c["major_function_cn"].(string); !ok || val != majorFilter {
					match = false
				}
			}
			if middleFilter != "" && match {
				if val, ok := c["middle_function_cn"].(string); !ok || val != middleFilter {
					match = false
				}
			}
			if minorFilter != "" && match {
				if val, ok := c["minor_function_cn"].(string); !ok || val != minorFilter {
					match = false
				}
			}

			if match {
				matchedCases = append(matchedCases, c)
			}
		}

		if len(matchedCases) == 0 {
			return tools.NewJSONResult(fmt.Sprintf(`{
				"success": 0,
				"failed": 0,
				"message": "没有找到匹配filter条件的用例",
				"filter": %s
			}`, func() string {
				filterJSON, _ := json.Marshal(filterInterface)
				return string(filterJSON)
			}())), nil
		}

		// 更新匹配的用例
		for idx, c := range matchedCases {
			caseID := 0
			if id, ok := c["id"].(float64); ok {
				caseID = int(id)
			}

			if caseID == 0 {
				failedCount++
				results = append(results, map[string]interface{}{
					"index":  idx,
					"status": "failed",
					"error":  "无法获取用例ID",
				})
				if !continueOnError {
					break
				}
				continue
			}

			path := fmt.Sprintf("/api/v1/projects/%d/case-groups/%d/manual-cases/%d", projectID, groupID, caseID)
			_, err := h.client.Put(ctx, path, updateDataInterface)
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
			"mode": "filter",
			"matched_count": %d,
			"success": %d,
			"failed": %d,
			"filter": %s,
			"results": %s
		}`, len(matchedCases), successCount, failedCount, func() string {
			filterJSON, _ := json.Marshal(filterInterface)
			return string(filterJSON)
		}(), func() string {
			resultJSON, _ := json.Marshal(results)
			return string(resultJSON)
		}())), nil
	}

	// 直接更新模式：通过 cases 数组中的 id 字段更新
	casesInterface, ok := args["cases"].([]interface{})
	if !ok {
		return tools.NewErrorResult("必须提供 cases 数组或 filter 条件"), nil
	}

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

		// 使用不带 group_id 的路径（更可靠）
		path := fmt.Sprintf("/api/v1/projects/%d/manual-cases/%d", projectID, caseID)

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
