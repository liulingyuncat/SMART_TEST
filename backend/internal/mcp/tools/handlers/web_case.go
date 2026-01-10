package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// ListWebGroupsHandler handles listing web case groups.
type ListWebGroupsHandler struct {
	*BaseHandler
}

func NewListWebGroupsHandler(c *client.BackendClient) *ListWebGroupsHandler {
	return &ListWebGroupsHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListWebGroupsHandler) Name() string {
	return "list_web_groups"
}

func (h *ListWebGroupsHandler) Description() string {
	return "获取项目的Web自动化用例集列表"
}

func (h *ListWebGroupsHandler) InputSchema() map[string]interface{} {
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

func (h *ListWebGroupsHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/case-groups", projectID)
	params := map[string]string{
		"case_type": "web",
	}
	data, err := h.client.Get(ctx, path, params)
	if err != nil {
		return tools.NewErrorResult(fmt.Sprintf("请求失败: %v", err)), nil
	}

	// 验证返回数据是否为有效JSON
	var testParse interface{}
	if err := json.Unmarshal(data, &testParse); err != nil {
		return tools.NewErrorResult(fmt.Sprintf("返回数据不是有效JSON: %v, 原始数据: %s", err, string(data[:min(len(data), 200)]))), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetWebGroupMetadataHandler handles getting web case group metadata.
type GetWebGroupMetadataHandler struct {
	*BaseHandler
}

func NewGetWebGroupMetadataHandler(c *client.BackendClient) *GetWebGroupMetadataHandler {
	return &GetWebGroupMetadataHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *GetWebGroupMetadataHandler) Name() string {
	return "get_web_group_metadata"
}

func (h *GetWebGroupMetadataHandler) Description() string {
	return "获取Web用例集的元数据（协议、服务器、端口、用户名、密码），用于自动化执行"
}

func (h *GetWebGroupMetadataHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "Web用例集ID",
			},
		},
		"required": []interface{}{"group_id"},
	}
}

func (h *GetWebGroupMetadataHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	groupID, err := GetInt(args, "group_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 获取用例集详情
	path := fmt.Sprintf("/api/v1/case-groups/%d", groupID)
	data, err := h.client.Get(ctx, path, nil)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 解析响应，提取元数据字段
	var groupData map[string]interface{}
	if err := json.Unmarshal(data, &groupData); err != nil {
		return tools.NewErrorResult(fmt.Sprintf("解析响应失败: %v", err)), nil
	}

	// 构建元数据响应
	metadata := map[string]interface{}{
		"group_id":      groupID,
		"group_name":    groupData["group_name"],
		"meta_protocol": groupData["meta_protocol"],
		"meta_server":   groupData["meta_server"],
		"meta_port":     groupData["meta_port"],
		"meta_user":     groupData["meta_user"],
		"meta_password": groupData["meta_password"],
	}

	return tools.NewJSONResult(tools.MustMarshalJSON(metadata)), nil
}

// ListWebCasesHandler handles listing web test cases.
type ListWebCasesHandler struct {
	*BaseHandler
}

func NewListWebCasesHandler(c *client.BackendClient) *ListWebCasesHandler {
	return &ListWebCasesHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListWebCasesHandler) Name() string {
	return "list_web_cases"
}

func (h *ListWebCasesHandler) Description() string {
	return "获取指定Web用例集的全部测试用例及所有字段（包括script_code脚本代码字段）"
}

func (h *ListWebCasesHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "Web用例集ID",
			},
		},
		"required": []interface{}{"project_id", "group_id"},
	}
}

func (h *ListWebCasesHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	groupID, err := GetInt(args, "group_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/auto-cases", projectID)
	params := map[string]string{
		"case_type":  "web",
		"case_group": fmt.Sprintf("%d", groupID),
		"size":       "99999",
	}

	data, err := h.client.Get(ctx, path, params)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// CreateWebGroupHandler handles creating a web case group.
type CreateWebGroupHandler struct {
	*BaseHandler
}

func NewCreateWebGroupHandler(c *client.BackendClient) *CreateWebGroupHandler {
	return &CreateWebGroupHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *CreateWebGroupHandler) Name() string {
	return "create_web_group"
}

func (h *CreateWebGroupHandler) Description() string {
	return "创建Web自动化用例集"
}

func (h *CreateWebGroupHandler) InputSchema() map[string]interface{} {
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
			"description": map[string]interface{}{
				"type":        "string",
				"description": "用例集描述(可选)",
			},
			"meta_protocol": map[string]interface{}{
				"type":        "string",
				"description": "元数据-协议（如: https, http）",
			},
			"meta_server": map[string]interface{}{
				"type":        "string",
				"description": "元数据-服务器地址",
			},
			"meta_port": map[string]interface{}{
				"type":        "string",
				"description": "元数据-端口号",
			},
			"meta_user": map[string]interface{}{
				"type":        "string",
				"description": "元数据-用户名",
			},
			"meta_password": map[string]interface{}{
				"type":        "string",
				"description": "元数据-密码",
			},
		},
		"required": []interface{}{"project_id", "name"},
	}
}

func (h *CreateWebGroupHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return tools.NewErrorResult("name must be a non-empty string"), nil
	}

	description := ""
	if desc, ok := args["description"].(string); ok {
		description = desc
	}

	// Prepare request body
	requestBody := map[string]interface{}{
		"group_name": name,
		"case_type":  "web",
	}
	if description != "" {
		requestBody["description"] = description
	}

	// 元数据字段
	if metaProtocol := GetOptionalString(args, "meta_protocol", ""); metaProtocol != "" {
		requestBody["meta_protocol"] = metaProtocol
	}
	if metaServer := GetOptionalString(args, "meta_server", ""); metaServer != "" {
		requestBody["meta_server"] = metaServer
	}
	if metaPort := GetOptionalString(args, "meta_port", ""); metaPort != "" {
		requestBody["meta_port"] = metaPort
	}
	if metaUser := GetOptionalString(args, "meta_user", ""); metaUser != "" {
		requestBody["meta_user"] = metaUser
	}
	if metaPassword := GetOptionalString(args, "meta_password", ""); metaPassword != "" {
		requestBody["meta_password"] = metaPassword
	}

	path := fmt.Sprintf("/api/v1/projects/%d/case-groups", projectID)
	resp, err := h.client.Post(ctx, path, requestBody)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(resp)), nil
}

// UpdateWebCasesHandler handles updating multiple web test cases in a group.
type UpdateWebCasesHandler struct {
	*BaseHandler
}

func NewUpdateWebCasesHandler(c *client.BackendClient) *UpdateWebCasesHandler {
	return &UpdateWebCasesHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *UpdateWebCasesHandler) Name() string {
	return "update_web_cases"
}

func (h *UpdateWebCasesHandler) Description() string {
	return "批量更新指定Web用例集内的所有用例的全部字段(除UUID)，支持更新script_code脚本代码字段"
}

func (h *UpdateWebCasesHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "Web用例集ID",
			},
			"cases": map[string]interface{}{
				"type":        "array",
				"description": "要更新的用例数据数组，每个用例对象需包含id和其他字段(不包括UUID)，支持更新script_code脚本代码字段",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "integer",
							"description": "用例ID（必填）",
						},
						"script_code": map[string]interface{}{
							"type":        "string",
							"description": "Playwright脚本代码，格式为async (page) => { ... }",
						},
					},
					"required": []interface{}{"id"},
				},
			},
		},
		"required": []interface{}{"project_id", "group_id", "cases"},
	}
}

func (h *UpdateWebCasesHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
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

	if len(casesInterface) == 0 {
		return tools.NewErrorResult("cases array cannot be empty"), nil
	}

	// 第一步：通过list_web_cases获取该用例集的所有用例，建立 id -> case_id(UUID) 的映射
	listPath := fmt.Sprintf("/api/v1/projects/%d/auto-cases", projectID)
	listParams := map[string]string{
		"case_type":  "web",
		"case_group": fmt.Sprintf("%d", groupID), // 传递group_id，后端会自动转换为group_name
		"page":       "1",
		"size":       "99999", // 获取全部用例以支持更新任意ID的用例
	}
	listData, err := h.client.Get(ctx, listPath, listParams)
	if err != nil {
		return tools.NewErrorResult(fmt.Sprintf("failed to list cases: %v", err)), nil
	}

	// 解析用例列表
	var listResponse struct {
		Data struct {
			Cases []struct {
				ID     int    `json:"id"`
				CaseID string `json:"case_id"`
			} `json:"cases"`
		} `json:"data"`
	}
	if err := json.Unmarshal(listData, &listResponse); err != nil {
		return tools.NewErrorResult(fmt.Sprintf("failed to parse cases list: %v", err)), nil
	}

	// 建立 ID -> UUID 映射
	idToUUID := make(map[int]string)
	for _, c := range listResponse.Data.Cases {
		idToUUID[c.ID] = c.CaseID
	}

	// 第二步：逐个更新用例
	results := []map[string]interface{}{}
	successCount := 0
	failedCount := 0

	for idx, caseItem := range casesInterface {
		caseData, ok := caseItem.(map[string]interface{})
		if !ok {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  "case item must be an object",
			})
			continue
		}

		// 获取用例的数字ID
		var caseNumID int
		if idFloat, ok := caseData["id"].(float64); ok {
			caseNumID = int(idFloat)
		} else if idInt, ok := caseData["id"].(int); ok {
			caseNumID = idInt
		} else {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  "case item must contain 'id' field (integer)",
			})
			continue
		}

		// 查找对应的UUID
		caseUUID, exists := idToUUID[caseNumID]
		if !exists {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":   idx,
				"case_id": caseNumID,
				"status":  "failed",
				"error":   fmt.Sprintf("case with id=%d not found", caseNumID),
			})
			continue
		}

		// 准备更新数据（移除id字段，因为它不是更新的字段）
		updateData := make(map[string]interface{})
		for k, v := range caseData {
			if k != "id" { // 排除id字段
				updateData[k] = v
			}
		}

		// 调用单个更新API
		updatePath := fmt.Sprintf("/api/v1/projects/%d/auto-cases/%s", projectID, caseUUID)
		_, err := h.client.Patch(ctx, updatePath, updateData)
		if err != nil {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":   idx,
				"case_id": caseNumID,
				"status":  "failed",
				"error":   err.Error(),
			})
			continue
		}

		successCount++
		results = append(results, map[string]interface{}{
			"index":   idx,
			"case_id": caseNumID,
			"status":  "success",
		})
	}

	// 返回批量更新结果
	response := map[string]interface{}{
		"success": successCount,
		"failed":  failedCount,
		"results": results,
	}

	responseJSON, _ := json.Marshal(response)
	return tools.NewJSONResult(string(responseJSON)), nil
}

// CreateWebCasesHandler handles batch creating web test cases.
type CreateWebCasesHandler struct {
	*BaseHandler
}

func NewCreateWebCasesHandler(c *client.BackendClient) *CreateWebCasesHandler {
	return &CreateWebCasesHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *CreateWebCasesHandler) Name() string {
	return "create_web_cases"
}

func (h *CreateWebCasesHandler) Description() string {
	return "批量创建Web自动化测试用例，支持script_code脚本代码字段用于自动化执行"
}

func (h *CreateWebCasesHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "Web用例集ID（与group_name二选一）",
			},
			"group_name": map[string]interface{}{
				"type":        "string",
				"description": "Web用例集名称（与group_id二选一），会自动查找对应的group_id",
			},
			"cases": map[string]interface{}{
				"type":        "array",
				"description": "用例数据数组，每个元素包含各语言字段(screen_cn/jp/en, function_cn/jp/en, test_steps_cn/jp/en, expected_result_cn/jp/en等)和script_code脚本代码字段",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"case_number": map[string]interface{}{
							"type":        "string",
							"description": "用例编号，格式如LOGIN-001",
						},
						"script_code": map[string]interface{}{
							"type":        "string",
							"description": "Playwright脚本代码，格式为async (page) => { ... }，用于自动化执行",
						},
					},
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

func (h *CreateWebCasesHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 尝试获取group_id，如果不存在则尝试通过group_name查找
	groupID := 0
	hasGroupID := false
	hasGroupName := false
	groupName := ""

	if _, ok := args["group_id"]; ok {
		hasGroupID = true
	}
	if gn, ok := args["group_name"]; ok {
		hasGroupName = true
		if str, ok := gn.(string); ok {
			groupName = str
		}
	}

	if hasGroupID {
		// 如果提供了group_id，直接使用
		id, err := GetInt(args, "group_id")
		if err != nil {
			return tools.NewErrorResult(err.Error()), nil
		}
		groupID = id
	} else if hasGroupName {
		// 如果提供了group_name，需要先查找对应的group_id
		if groupName == "" {
			return tools.NewErrorResult("group_name must be a non-empty string"), nil
		}

		// 查询Web用例集列表
		groupsPath := fmt.Sprintf("/api/v1/projects/%d/case-groups?case_type=web", projectID)
		groupsData, err := h.client.Get(ctx, groupsPath, nil)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("failed to fetch case groups: %v", err)), nil
		}

		// 解析响应找到matching的group
		var response interface{}
		err = json.Unmarshal(groupsData, &response)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("failed to parse groups response: %v", err)), nil
		}

		// 查找group_name匹配的用例集
		found := false

		// 尝试作为map处理（标准响应格式）
		if respMap, ok := response.(map[string]interface{}); ok {
			if dataVal, ok := respMap["data"].([]interface{}); ok {
				for _, item := range dataVal {
					if itemMap, ok := item.(map[string]interface{}); ok {
						if name, ok := itemMap["group_name"].(string); ok && name == groupName {
							if id, ok := itemMap["id"].(float64); ok {
								groupID = int(id)
								found = true
								break
							}
						}
					}
				}
			}
		} else if dataArray, ok := response.([]interface{}); ok {
			// 直接为数组格式
			for _, item := range dataArray {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if name, ok := itemMap["group_name"].(string); ok && name == groupName {
						if id, ok := itemMap["id"].(float64); ok {
							groupID = int(id)
							found = true
							break
						}
					}
				}
			}
		}

		if !found {
			return tools.NewErrorResult(fmt.Sprintf("web case group '%s' not found in project %d", groupName, projectID)), nil
		}
	} else {
		// 既没有提供group_id也没有提供group_name
		return tools.NewErrorResult("either 'group_id' or 'group_name' must be provided"), nil
	}

	casesInterface, ok := args["cases"].([]interface{})
	if !ok {
		return tools.NewErrorResult("cases must be an array"), nil
	}

	if len(casesInterface) == 0 {
		return tools.NewErrorResult("cases array cannot be empty"), nil
	}

	continueOnError := true
	if val, ok := args["continue_on_error"].(bool); ok {
		continueOnError = val
	}

	results := []map[string]interface{}{}
	successCount := 0
	failedCount := 0

	path := fmt.Sprintf("/api/v1/projects/%d/auto-cases", projectID)

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

		// 添加必填的case_type字段（Web用例类型）
		data["case_type"] = "web"
		// 添加group_id到请求数据（后端会根据group_id查找对应的case_group名称）
		data["group_id"] = groupID

		// 确保script_code字段存在（即使为空）
		if _, exists := data["script_code"]; !exists {
			data["script_code"] = ""
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

		// 解析响应以获取创建的用例ID
		var respData map[string]interface{}
		err = json.Unmarshal(resp, &respData)
		if err == nil {
			successCount++
			// 如果响应包含data字段（如{"code":0,"data":{...}}）
			if dataVal, ok := respData["data"].(map[string]interface{}); ok {
				if id, ok := dataVal["id"].(float64); ok {
					if uuid, ok := dataVal["uuid"].(string); ok {
						results = append(results, map[string]interface{}{
							"index":     idx,
							"status":    "success",
							"case_id":   int(id),
							"case_uuid": uuid,
						})
					} else {
						results = append(results, map[string]interface{}{
							"index":   idx,
							"status":  "success",
							"case_id": int(id),
						})
					}
				} else {
					results = append(results, map[string]interface{}{
						"index":  idx,
						"status": "success",
						"data":   dataVal,
					})
				}
			} else if id, ok := respData["id"].(float64); ok {
				// 如果直接返回用例对象
				if uuid, ok := respData["uuid"].(string); ok {
					results = append(results, map[string]interface{}{
						"index":     idx,
						"status":    "success",
						"case_id":   int(id),
						"case_uuid": uuid,
					})
				} else {
					results = append(results, map[string]interface{}{
						"index":   idx,
						"status":  "success",
						"case_id": int(id),
					})
				}
			} else {
				results = append(results, map[string]interface{}{
					"index":  idx,
					"status": "success",
					"data":   respData,
				})
			}
		} else {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  fmt.Sprintf("failed to parse response: %v", err),
				"data":   string(resp),
			})
			if !continueOnError {
				break
			}
			continue
		}
	}

	resultData := map[string]interface{}{
		"success": successCount,
		"failed":  failedCount,
		"results": results,
	}

	return tools.NewJSONResult(tools.MustMarshalJSON(resultData)), nil
}
