package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// ============================================================================
// API用例智能排序模块
// 确保批量创建的用例按CRUD生命周期顺序排列，便于自动化执行
// ============================================================================

// apiCaseSortItem 用于排序的用例包装结构
type apiCaseSortItem struct {
	originalIndex int                    // 原始索引
	data          map[string]interface{} // 用例数据
	screen        string                 // 画面名称
	method        string                 // HTTP方法
	responseCode  int                    // 响应码
	url           string                 // URL路径
}

// getMethodWeight 获取HTTP方法的排序权重（CRUD顺序）
// GET(查询) → POST(创建) → PUT(更新) → PATCH(部分更新) → DELETE(删除)
func getMethodWeight(method string) int {
	weights := map[string]int{
		"GET":    1,
		"POST":   2,
		"PUT":    3,
		"PATCH":  4,
		"DELETE": 5,
	}
	if w, ok := weights[strings.ToUpper(method)]; ok {
		return w
	}
	return 99
}

// getResponseCodeWeight 获取响应码的排序权重
// 正常响应 → 客户端错误 → 服务器错误
func getResponseCodeWeight(code int) int {
	switch {
	case code >= 200 && code < 300:
		return 1 // 成功响应优先
	case code >= 400 && code < 500:
		// 细分客户端错误
		switch code {
		case 400:
			return 2 // Bad Request
		case 401:
			return 3 // Unauthorized
		case 403:
			return 4 // Forbidden
		case 404:
			return 5 // Not Found
		case 409:
			return 6 // Conflict
		case 422:
			return 7 // Unprocessable Entity
		default:
			return 8
		}
	case code >= 500:
		return 9 // 服务器错误最后
	default:
		return 10
	}
}

// extractResponseCode 从response字段提取HTTP响应码
func extractResponseCode(response string) int {
	if response == "" {
		return 200 // 默认200
	}

	// 尝试解析JSON格式的response
	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(response), &respData); err == nil {
		// 查找code字段
		if code, ok := respData["code"].(float64); ok {
			return int(code)
		}
	}

	// 尝试匹配 "200:" 或 "401:" 格式
	re := regexp.MustCompile(`^(\d{3}):`)
	if matches := re.FindStringSubmatch(response); len(matches) > 1 {
		if code, err := strconv.Atoi(matches[1]); err == nil {
			return code
		}
	}

	// 尝试匹配 "code": 200 或 "code":401 格式
	re2 := regexp.MustCompile(`"code"\s*:\s*(\d{3})`)
	if matches := re2.FindStringSubmatch(response); len(matches) > 1 {
		if code, err := strconv.Atoi(matches[1]); err == nil {
			return code
		}
	}

	return 200 // 默认200
}

// sortAPICases 对API用例数组进行智能排序
// 排序规则：
// 1. 按screen（画面）分组
// 2. 同一画面内按HTTP方法排序（GET→POST→PUT→DELETE）
// 3. 同一方法内按响应码排序（200→4xx→5xx）
func sortAPICases(cases []interface{}) []interface{} {
	if len(cases) <= 1 {
		return cases
	}

	// 转换为排序结构
	items := make([]apiCaseSortItem, 0, len(cases))
	for i, c := range cases {
		data, ok := c.(map[string]interface{})
		if !ok {
			continue
		}

		item := apiCaseSortItem{
			originalIndex: i,
			data:          data,
		}

		// 提取screen
		if screen, ok := data["screen"].(string); ok {
			item.screen = screen
		}

		// 提取method
		if method, ok := data["method"].(string); ok {
			item.method = strings.ToUpper(method)
		} else {
			item.method = "GET"
		}

		// 提取url
		if url, ok := data["url"].(string); ok {
			item.url = url
		}

		// 提取response并解析响应码
		if response, ok := data["response"].(string); ok {
			item.responseCode = extractResponseCode(response)
		} else {
			item.responseCode = 200
		}

		items = append(items, item)
	}

	// 排序
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]

		// 1. 先按screen排序
		if a.screen != b.screen {
			return a.screen < b.screen
		}

		// 2. 同一screen内，按URL排序（确保同一接口的用例聚合）
		if a.url != b.url {
			return a.url < b.url
		}

		// 3. 同一URL内，按HTTP方法排序
		aMethodWeight := getMethodWeight(a.method)
		bMethodWeight := getMethodWeight(b.method)
		if aMethodWeight != bMethodWeight {
			return aMethodWeight < bMethodWeight
		}

		// 4. 同一方法内，按响应码排序
		aRespWeight := getResponseCodeWeight(a.responseCode)
		bRespWeight := getResponseCodeWeight(b.responseCode)
		return aRespWeight < bRespWeight
	})

	// 转换回原始格式
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = item.data
	}

	return result
}

// ListApiGroupsHandler handles listing API case groups.
type ListApiGroupsHandler struct {
	*BaseHandler
}

func NewListApiGroupsHandler(c *client.BackendClient) *ListApiGroupsHandler {
	return &ListApiGroupsHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListApiGroupsHandler) Name() string {
	return "list_api_groups"
}

func (h *ListApiGroupsHandler) Description() string {
	return "获取项目的接口用例集列表（包含元数据）"
}

func (h *ListApiGroupsHandler) InputSchema() map[string]interface{} {
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

func (h *ListApiGroupsHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/case-groups?case_type=api", projectID)
	data, err := h.client.Get(ctx, path, nil)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// GetApiGroupMetadataHandler handles getting API case group metadata.
type GetApiGroupMetadataHandler struct {
	*BaseHandler
}

func NewGetApiGroupMetadataHandler(c *client.BackendClient) *GetApiGroupMetadataHandler {
	return &GetApiGroupMetadataHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *GetApiGroupMetadataHandler) Name() string {
	return "get_api_group_metadata"
}

func (h *GetApiGroupMetadataHandler) Description() string {
	return "获取接口用例集的元数据（协议、服务器、端口、用户名、密码），用于自动化执行"
}

func (h *GetApiGroupMetadataHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "接口用例集ID",
			},
		},
		"required": []interface{}{"group_id"},
	}
}

func (h *GetApiGroupMetadataHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
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

// ListApiCasesHandler handles listing API test cases.
type ListApiCasesHandler struct {
	*BaseHandler
}

func NewListApiCasesHandler(c *client.BackendClient) *ListApiCasesHandler {
	return &ListApiCasesHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListApiCasesHandler) Name() string {
	return "list_api_cases"
}

func (h *ListApiCasesHandler) Description() string {
	return "获取指定API用例集的全部测试用例及所有字段"
}

func (h *ListApiCasesHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "API用例集ID",
			},
		},
		"required": []interface{}{"project_id", "group_id"},
	}
}

func (h *ListApiCasesHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	groupID, err := GetInt(args, "group_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 首先通过 group_id 获取用例集详情，获取 group_name
	groupPath := fmt.Sprintf("/api/v1/case-groups/%d", groupID)
	groupData, err := h.client.Get(ctx, groupPath, nil)
	if err != nil {
		return tools.NewErrorResult(fmt.Sprintf("获取用例集信息失败: %v", err)), nil
	}

	// 解析用例集信息获取 group_name
	var groupInfo map[string]interface{}
	if err := json.Unmarshal(groupData, &groupInfo); err != nil {
		return tools.NewErrorResult(fmt.Sprintf("解析用例集信息失败: %v", err)), nil
	}

	groupName, ok := groupInfo["group_name"].(string)
	if !ok || groupName == "" {
		return tools.NewErrorResult("用例集名称不存在"), nil
	}

	// 使用 group_name 查询用例列表
	path := fmt.Sprintf("/api/v1/projects/%d/api-cases", projectID)
	params := map[string]string{
		"case_type":  "api",
		"case_group": groupName,
		"size":       "99999",
	}

	data, err := h.client.Get(ctx, path, params)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// CreateApiCaseHandler handles batch creating API test cases.
type CreateApiCaseHandler struct {
	*BaseHandler
}

func NewCreateApiCaseHandler(c *client.BackendClient) *CreateApiCaseHandler {
	return &CreateApiCaseHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *CreateApiCaseHandler) Name() string {
	return "create_api_cases"
}

func (h *CreateApiCaseHandler) Description() string {
	return "批量创建API接口测试用例"
}

func (h *CreateApiCaseHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "API用例集ID（与group_name二选一）",
			},
			"group_name": map[string]interface{}{
				"type":        "string",
				"description": "API用例集名称（与group_id二选一），会自动查找对应的group_id",
			},
			"cases": map[string]interface{}{
				"type":        "array",
				"description": "用例数据数组，每个元素包含用例字段(screen, url, method, header, body, response等)",
				"items": map[string]interface{}{
					"type": "object",
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

func (h *CreateApiCaseHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 尝试获取group_id，如果不存在则尝试通过group_name查找
	groupID := 0
	groupName := ""
	hasGroupID := false
	hasGroupName := false

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
		id, err := GetInt(args, "group_id")
		if err != nil {
			return tools.NewErrorResult(err.Error()), nil
		}
		groupID = id
		// 获取group_name用于创建用例
		groupPath := fmt.Sprintf("/api/v1/case-groups/%d", groupID)
		groupData, err := h.client.Get(ctx, groupPath, nil)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("获取用例集信息失败: %v", err)), nil
		}
		var groupInfo map[string]interface{}
		if err := json.Unmarshal(groupData, &groupInfo); err != nil {
			return tools.NewErrorResult(fmt.Sprintf("解析用例集信息失败: %v", err)), nil
		}

		// 尝试从data字段获取group_name（处理嵌套响应格式）
		if dataField, ok := groupInfo["data"].(map[string]interface{}); ok {
			if name, ok := dataField["group_name"].(string); ok && name != "" {
				groupName = name
			}
		} else if name, ok := groupInfo["group_name"].(string); ok && name != "" {
			// 直接在顶级字段中
			groupName = name
		}

		// 验证groupName是否成功获取
		if groupName == "" {
			return tools.NewErrorResult("无法获取用例集名称，请检查group_id是否正确"), nil
		}
	} else if hasGroupName {
		if groupName == "" {
			return tools.NewErrorResult("group_name must be a non-empty string"), nil
		}
		// 查询API用例集列表找到对应的group_id
		groupsPath := fmt.Sprintf("/api/v1/projects/%d/case-groups?case_type=api", projectID)
		groupsData, err := h.client.Get(ctx, groupsPath, nil)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("failed to fetch case groups: %v", err)), nil
		}

		var response interface{}
		err = json.Unmarshal(groupsData, &response)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("failed to parse groups response: %v", err)), nil
		}

		found := false
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
			return tools.NewErrorResult(fmt.Sprintf("API case group '%s' not found in project %d", groupName, projectID)), nil
		}
	} else {
		return tools.NewErrorResult("either 'group_id' or 'group_name' must be provided"), nil
	}

	casesInterface, ok := args["cases"].([]interface{})
	if !ok {
		return tools.NewErrorResult("cases must be an array"), nil
	}

	if len(casesInterface) == 0 {
		return tools.NewErrorResult("cases array cannot be empty"), nil
	}

	// ========================================
	// 智能排序：确保用例按CRUD生命周期顺序排列
	// 排序规则：screen → url → method(GET→POST→PUT→DELETE) → responseCode(200→4xx→5xx)
	// ========================================
	casesInterface = sortAPICases(casesInterface)

	continueOnError := true
	if val, ok := args["continue_on_error"].(bool); ok {
		continueOnError = val
	}

	results := []map[string]interface{}{}
	successCount := 0
	failedCount := 0

	// 获取该用例集的所有用例，以便获取第一个用例作为插入的目标
	listPath := fmt.Sprintf("/api/v1/projects/%d/api-cases", projectID)
	listParams := map[string]string{
		"case_type":  "api",
		"case_group": groupName,
		"page":       "1",
		"size":       "1",
	}
	listData, err := h.client.Get(ctx, listPath, listParams)
	if err != nil {
		return tools.NewErrorResult(fmt.Sprintf("获取用例列表失败: %v", err)), nil
	}

	// 解析用例列表获取第一个用例
	var listResponse struct {
		Data struct {
			Cases []struct {
				CaseID string `json:"case_id"`
			} `json:"cases"`
		} `json:"data"`
	}
	if err := json.Unmarshal(listData, &listResponse); err != nil {
		return tools.NewErrorResult(fmt.Sprintf("解析用例列表失败: %v", err)), nil
	}

	// 如果用例集为空，需要先创建一个初始用例
	var targetCaseID string
	if len(listResponse.Data.Cases) == 0 {
		// 用例集为空，用第一条用例的数据先创建初始用例，然后更新其字段
		if len(casesInterface) == 0 {
			return tools.NewErrorResult("cases array cannot be empty"), nil
		}

		firstCaseData, ok := casesInterface[0].(map[string]interface{})
		if !ok {
			return tools.NewErrorResult("first case item must be an object"), nil
		}

		// 第1步：用CreateCase创建基础初始用例
		initPath := fmt.Sprintf("/api/v1/projects/%d/api-cases", projectID)
		initRequestData := map[string]interface{}{
			"case_type":   "api",
			"case_group":  groupName,
			"method":      "GET",
			"test_result": "NR",
		}

		initResp, err := h.client.Post(ctx, initPath, initRequestData)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("创建初始用例失败: %v", err)), nil
		}

		// 解析响应获取初始用例ID
		var initRespData map[string]interface{}
		if err := json.Unmarshal(initResp, &initRespData); err == nil {
			if dataVal, ok := initRespData["data"].(map[string]interface{}); ok {
				if id, ok := dataVal["case_id"].(string); ok {
					targetCaseID = id
				}
			}
		}

		if targetCaseID == "" {
			return tools.NewErrorResult("创建初始用例失败，无法获取case_id"), nil
		}

		// 第2步：用UpdateCase更新第一条用例的所有字段
		updatePath := fmt.Sprintf("/api/v1/projects/%d/api-cases/%s", projectID, targetCaseID)
		updateData := make(map[string]interface{})

		// 复制第一条用例的所有字段
		for key, value := range firstCaseData {
			if key != "case_group" && key != "case_id" {
				updateData[key] = value
			}
		}

		_, err = h.client.Patch(ctx, updatePath, updateData)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("更新初始用例字段失败: %v", err)), nil
		}

		successCount++
		results = append(results, map[string]interface{}{
			"index":  0,
			"status": "success",
		})
	} else {
		// 获取用例集中的第一个用例作为目标
		targetCaseID = listResponse.Data.Cases[0].CaseID
	}

	if targetCaseID == "" {
		return tools.NewErrorResult("无法获取目标用例ID，无法插入新用例"), nil
	}

	// 确定从哪个索引开始处理用例（如果用例集为空，第一条已在初始化时处理）
	var startIdx int = 0
	if len(listResponse.Data.Cases) == 0 && targetCaseID != "" {
		startIdx = 1
	}

	path := fmt.Sprintf("/api/v1/projects/%d/api-cases/insert", projectID)

	for i, caseItem := range casesInterface[startIdx:] {
		idx := i + startIdx
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

		// 为insert接口创建请求数据，使用case_data字段包含所有自定义字段
		requestData := map[string]interface{}{
			"case_type":      "api",
			"case_group":     groupName,
			"position":       "after",
			"target_case_id": targetCaseID,
			"case_data":      make(map[string]interface{}),
		}

		// 获取case_data map
		caseDataMap := requestData["case_data"].(map[string]interface{})

		// 设置默认method
		if _, exists := data["method"]; !exists {
			data["method"] = "GET"
		}

		// 复制所有输入字段到case_data中
		for key, value := range data {
			caseDataMap[key] = value
		}

		// 确保必要的字段存在（即使为空）
		ensureFieldExists := func(field string) {
			if _, exists := caseDataMap[field]; !exists {
				caseDataMap[field] = ""
			}
		}
		ensureFieldExists("screen")
		ensureFieldExists("url")
		ensureFieldExists("header")
		ensureFieldExists("body")
		ensureFieldExists("response")
		ensureFieldExists("method")
		ensureFieldExists("script_code")

		resp, err := h.client.Post(ctx, path, requestData)
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

		// 解析响应获取case_id
		var respData map[string]interface{}
		caseID := ""
		if err := json.Unmarshal(resp, &respData); err == nil {
			if dataVal, ok := respData["data"].(map[string]interface{}); ok {
				if id, ok := dataVal["case_id"].(string); ok {
					caseID = id
				}
			}
		}

		// 更新targetCaseID为刚插入的用例ID，确保下一个用例插入在它后面
		// 这样可以保持用例的正确顺序（按数组顺序依次追加）
		if caseID != "" {
			targetCaseID = caseID
		}

		successCount++
		result := map[string]interface{}{
			"index":  idx,
			"status": "success",
		}
		if caseID != "" {
			result["case_id"] = caseID
		}
		results = append(results, result)
	}

	response := map[string]interface{}{
		"success": successCount,
		"failed":  failedCount,
		"results": results,
	}

	responseJSON, _ := json.Marshal(response)
	return tools.NewJSONResult(string(responseJSON)), nil
}

// UpdateApiCaseHandler handles batch updating API test cases.
type UpdateApiCaseHandler struct {
	*BaseHandler
}

func NewUpdateApiCaseHandler(c *client.BackendClient) *UpdateApiCaseHandler {
	return &UpdateApiCaseHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *UpdateApiCaseHandler) Name() string {
	return "update_api_cases"
}

func (h *UpdateApiCaseHandler) Description() string {
	return "批量更新API接口测试用例"
}

func (h *UpdateApiCaseHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"group_id": map[string]interface{}{
				"type":        "integer",
				"description": "API用例集ID",
			},
			"cases": map[string]interface{}{
				"type":        "array",
				"description": "要更新的用例数据数组，每个用例对象需包含case_id(UUID)或id(数字ID)和其他要更新的字段",
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

func (h *UpdateApiCaseHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
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

	continueOnError := true
	if val, ok := args["continue_on_error"].(bool); ok {
		continueOnError = val
	}

	// 首先通过group_id获取用例集名称
	groupPath := fmt.Sprintf("/api/v1/case-groups/%d", groupID)
	groupData, err := h.client.Get(ctx, groupPath, nil)
	if err != nil {
		return tools.NewErrorResult(fmt.Sprintf("获取用例集信息失败: %v", err)), nil
	}

	var groupInfo map[string]interface{}
	if err := json.Unmarshal(groupData, &groupInfo); err != nil {
		return tools.NewErrorResult(fmt.Sprintf("解析用例集信息失败: %v", err)), nil
	}

	// 尝试从data字段获取group_name（处理嵌套响应格式）
	groupName := ""
	if dataField, ok := groupInfo["data"].(map[string]interface{}); ok {
		if name, ok := dataField["group_name"].(string); ok && name != "" {
			groupName = name
		}
	} else if name, ok := groupInfo["group_name"].(string); ok && name != "" {
		// 直接在顶级字段中
		groupName = name
	}

	// 验证groupName是否成功获取
	if groupName == "" {
		return tools.NewErrorResult("无法获取用例集名称，请检查group_id是否正确"), nil
	}

	// 获取该用例集的所有用例，建立 id -> case_id(UUID) 的映射
	listPath := fmt.Sprintf("/api/v1/projects/%d/api-cases", projectID)
	listParams := map[string]string{
		"case_type":  "api",
		"case_group": groupName,
		"page":       "1",
		"size":       "99999", // 获取全部用例以支持更新任意ID的用例
	}
	listData, err := h.client.Get(ctx, listPath, listParams)
	if err != nil {
		return tools.NewErrorResult(fmt.Sprintf("获取用例列表失败: %v", err)), nil
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
		return tools.NewErrorResult(fmt.Sprintf("解析用例列表失败: %v", err)), nil
	}

	// 建立 ID -> UUID 映射
	idToUUID := make(map[int]string)
	for _, c := range listResponse.Data.Cases {
		idToUUID[c.ID] = c.CaseID
	}

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
			if !continueOnError {
				break
			}
			continue
		}

		// 获取用例UUID - 优先使用case_id，其次使用id进行映射
		var caseUUID string
		var caseNumID int

		if cid, ok := caseData["case_id"].(string); ok && cid != "" {
			caseUUID = cid
		} else if idFloat, ok := caseData["id"].(float64); ok {
			caseNumID = int(idFloat)
			if uuid, exists := idToUUID[caseNumID]; exists {
				caseUUID = uuid
			}
		} else if idInt, ok := caseData["id"].(int); ok {
			caseNumID = idInt
			if uuid, exists := idToUUID[caseNumID]; exists {
				caseUUID = uuid
			}
		}

		if caseUUID == "" {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":  idx,
				"status": "failed",
				"error":  "case item must contain 'case_id' (UUID) or 'id' (integer ID)",
			})
			if !continueOnError {
				break
			}
			continue
		}

		// 准备更新数据（移除id和case_id字段）
		updateData := make(map[string]interface{})
		for k, v := range caseData {
			if k != "id" && k != "case_id" {
				updateData[k] = v
			}
		}

		// 调用单个更新API
		updatePath := fmt.Sprintf("/api/v1/projects/%d/api-cases/%s", projectID, caseUUID)
		_, err := h.client.Patch(ctx, updatePath, updateData)
		if err != nil {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":   idx,
				"case_id": caseUUID,
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
			"case_id": caseUUID,
			"status":  "success",
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

// CreateApiGroupHandler handles creating an API case group.
type CreateApiGroupHandler struct {
	*BaseHandler
}

func NewCreateApiGroupHandler(c *client.BackendClient) *CreateApiGroupHandler {
	return &CreateApiGroupHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *CreateApiGroupHandler) Name() string {
	return "create_api_group"
}

func (h *CreateApiGroupHandler) Description() string {
	return "创建接口用例集（支持元数据）"
}

func (h *CreateApiGroupHandler) InputSchema() map[string]interface{} {
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

func (h *CreateApiGroupHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return tools.NewErrorResult("name must be a non-empty string"), nil
	}

	// Prepare request body
	requestBody := map[string]interface{}{
		"group_name": name,
		"case_type":  "api",
	}

	if description := GetOptionalString(args, "description", ""); description != "" {
		requestBody["description"] = description
	}
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
