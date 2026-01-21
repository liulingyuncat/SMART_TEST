package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// ListExecutionTasksHandler handles listing execution tasks.
type ListExecutionTasksHandler struct {
	*BaseHandler
}

func NewListExecutionTasksHandler(c *client.BackendClient) *ListExecutionTasksHandler {
	return &ListExecutionTasksHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *ListExecutionTasksHandler) Name() string {
	return "list_execution_tasks"
}

func (h *ListExecutionTasksHandler) Description() string {
	return "获取项目的执行任务列表"
}

func (h *ListExecutionTasksHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"description": "任务状态过滤（可选）",
			},
		},
		"required": []interface{}{"project_id"},
	}
}

func (h *ListExecutionTasksHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/execution-tasks", projectID)
	params := make(map[string]string)

	if status := GetOptionalString(args, "status", ""); status != "" {
		params["status"] = status
	}

	data, err := h.client.Get(ctx, path, params)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// GetExecutionTaskMetadataHandler handles getting metadata for an execution task.
type GetExecutionTaskMetadataHandler struct {
	*BaseHandler
}

func NewGetExecutionTaskMetadataHandler(c *client.BackendClient) *GetExecutionTaskMetadataHandler {
	return &GetExecutionTaskMetadataHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *GetExecutionTaskMetadataHandler) Name() string {
	return "get_execution_task_metadata"
}

func (h *GetExecutionTaskMetadataHandler) Description() string {
	return "获取指定执行任务的元数据信息，包括任务名称、状态、执行内容、关联用例集、执行进度(OK/NG/NR/Block数量)、通过率等统计信息。当执行类型为api/web时，会自动获取用例集的连接元数据（协议、服务器、端口、用户名、密码）和用户自定义变量"
}

func (h *GetExecutionTaskMetadataHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"task_id": map[string]interface{}{
				"type":        "string",
				"description": "执行任务标识。支持三种格式：1) 任务名字字符串(如: qweb, 手工, api2)；2) UUID (格式: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)；3) 数字字符串索引（如 '1' 表示第一个任务）",
			},
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID（可选，如不提供则默认为1）",
			},
		},
		"required": []interface{}{"task_id"},
	}
}

func (h *GetExecutionTaskMetadataHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	// 处理task_id参数
	var taskID string
	var isNameLookup bool
	var isIndexLookup bool
	var taskIndex int = -1
	projectID := GetOptionalInt(args, "project_id", 1)

	if val, ok := args["task_id"]; ok {
		switch v := val.(type) {
		case float64:
			taskIndex = int(v)
			taskID = fmt.Sprintf("%d", taskIndex)
			isIndexLookup = true
		case int:
			taskIndex = v
			taskID = fmt.Sprintf("%d", v)
			isIndexLookup = true
		case int64:
			taskIndex = int(v)
			taskID = fmt.Sprintf("%d", v)
			isIndexLookup = true
		case string:
			if idx, err := strconv.Atoi(v); err == nil && idx > 0 {
				taskIndex = idx
				taskID = v
				isIndexLookup = true
			} else if isUUIDFormat(v) {
				taskID = v
			} else {
				taskID = v
				isNameLookup = true
			}
		default:
			taskID = fmt.Sprintf("%v", v)
			isNameLookup = true
		}
	} else {
		return tools.NewErrorResult("task_id parameter is required"), nil
	}

	// 步骤1: 获取项目的执行任务列表
	taskListPath := fmt.Sprintf("/api/v1/projects/%d/execution-tasks", projectID)
	taskListData, err := h.client.Get(ctx, taskListPath, nil)
	if err != nil {
		return tools.NewErrorResult("failed to list execution tasks: " + err.Error()), nil
	}

	// 解析任务列表
	var taskList struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(taskListData, &taskList); err != nil {
		return tools.NewErrorResult("failed to parse task list: " + err.Error()), nil
	}

	// 解析任务数组
	var tasks []struct {
		TaskUUID      string      `json:"task_uuid"`
		TaskName      string      `json:"task_name"`
		ProjectID     int         `json:"project_id"`
		ExecutionType string      `json:"execution_type"`
		TaskStatus    string      `json:"task_status"`
		CaseGroupID   int         `json:"case_group_id"`
		CaseGroupName string      `json:"case_group_name"`
		StartDate     interface{} `json:"start_date"`
		EndDate       interface{} `json:"end_date"`
		TestVersion   string      `json:"test_version"`
		TestEnv       string      `json:"test_env"`
		TestDate      interface{} `json:"test_date"`
		Executor      string      `json:"executor"`
		TaskDesc      string      `json:"task_description"`
		CreatedBy     int         `json:"created_by"`
		CreatedAt     string      `json:"created_at"`
		UpdatedAt     string      `json:"updated_at"`
	}

	if err := json.Unmarshal(taskList.Data, &tasks); err != nil {
		return tools.NewErrorResult("failed to parse tasks array: " + err.Error()), nil
	}

	// 查找匹配的任务
	var foundTask struct {
		TaskUUID      string      `json:"task_uuid"`
		TaskName      string      `json:"task_name"`
		ProjectID     int         `json:"project_id"`
		ExecutionType string      `json:"execution_type"`
		TaskStatus    string      `json:"task_status"`
		CaseGroupID   int         `json:"case_group_id"`
		CaseGroupName string      `json:"case_group_name"`
		StartDate     interface{} `json:"start_date"`
		EndDate       interface{} `json:"end_date"`
		TestVersion   string      `json:"test_version"`
		TestEnv       string      `json:"test_env"`
		TestDate      interface{} `json:"test_date"`
		Executor      string      `json:"executor"`
		TaskDesc      string      `json:"task_description"`
		CreatedBy     int         `json:"created_by"`
		CreatedAt     string      `json:"created_at"`
		UpdatedAt     string      `json:"updated_at"`
	}
	var taskMetadata interface{}
	var taskUUID string
	var found bool

	if isIndexLookup {
		if taskIndex < 1 || taskIndex > len(tasks) {
			return tools.NewErrorResult(fmt.Sprintf("task index %d out of range, project %d has %d tasks", taskIndex, projectID, len(tasks))), nil
		}
		foundTask = tasks[taskIndex-1]
		taskUUID = foundTask.TaskUUID
		taskMetadata = foundTask
		found = true
	} else if isNameLookup {
		for _, task := range tasks {
			if task.TaskName == taskID {
				foundTask = task
				taskUUID = task.TaskUUID
				taskMetadata = task
				found = true
				break
			}
		}
		if !found {
			return tools.NewErrorResult(fmt.Sprintf("task '%s' not found in project %d", taskID, projectID)), nil
		}
	} else {
		for _, task := range tasks {
			if task.TaskUUID == taskID {
				foundTask = task
				taskUUID = task.TaskUUID
				taskMetadata = task
				found = true
				break
			}
		}
		if !found {
			return tools.NewErrorResult(fmt.Sprintf("task with UUID '%s' not found in project %d", taskID, projectID)), nil
		}
	}

	// 步骤2: 获取执行任务的用例结果以计算统计信息
	casesPath := fmt.Sprintf("/api/v1/execution-tasks/%s/case-results", taskUUID)
	casesParams := map[string]string{
		"size": "99999",
	}
	casesData, err := h.client.Get(ctx, casesPath, casesParams)
	if err != nil {
		return tools.NewErrorResult("failed to get case results: " + err.Error()), nil
	}

	// 解析用例结果
	var casesResult struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(casesData, &casesResult); err != nil {
		return tools.NewErrorResult("failed to parse case results: " + err.Error()), nil
	}

	var caseResultsList []map[string]interface{}
	if err := json.Unmarshal(casesResult.Data, &caseResultsList); err != nil {
		return tools.NewErrorResult("failed to parse case results array: " + err.Error()), nil
	}

	// 步骤3: 计算执行统计信息
	totalCases := len(caseResultsList)
	okCount := 0
	ngCount := 0
	nrCount := 0
	blockCount := 0

	for _, caseResult := range caseResultsList {
		if result, ok := caseResult["result"].(string); ok {
			switch result {
			case "OK":
				okCount++
			case "NG":
				ngCount++
			case "NR":
				nrCount++
			case "Block":
				blockCount++
			}
		}
	}

	// 计算执行进度和通过率
	executedCount := okCount + ngCount + blockCount
	var progressRate float64 = 0
	var passRate float64 = 0

	if totalCases > 0 {
		progressRate = float64(executedCount) / float64(totalCases) * 100
	}
	if executedCount > 0 {
		passRate = float64(okCount) / float64(executedCount) * 100
	}

	// 步骤4: 如果是api/web/automation类型，获取用例集的连接元数据和用户自定义变量
	var groupMetadata map[string]interface{}
	var taskVariables []interface{}

	// 确定 groupType
	groupType := "web"
	if foundTask.ExecutionType == "api" {
		groupType = "api"
	}

	// 只对 api/web/automation 类型处理
	if foundTask.ExecutionType == "api" || foundTask.ExecutionType == "web" || foundTask.ExecutionType == "automation" {
		// 确定 groupID：优先使用 CaseGroupID，如果为 0 则通过 CaseGroupName 查找
		groupID := foundTask.CaseGroupID
		if groupID == 0 && foundTask.CaseGroupName != "" {
			// 通过用例集名称查找 group_id
			groupsPath := fmt.Sprintf("/api/v1/projects/%d/case-groups", projectID)
			groupsParams := map[string]string{"case_type": groupType}
			groupsData, err := h.client.Get(ctx, groupsPath, groupsParams)
			if err == nil {
				var groupsList []map[string]interface{}
				if json.Unmarshal(groupsData, &groupsList) == nil {
					for _, g := range groupsList {
						if name, ok := g["group_name"].(string); ok && name == foundTask.CaseGroupName {
							if id, ok := g["id"].(float64); ok {
								groupID = int(id)
								break
							}
						}
					}
				}
			}
		}

		// 获取用例集元数据
		if groupID > 0 {
			groupPath := fmt.Sprintf("/api/v1/case-groups/%d", groupID)
			groupData, err := h.client.Get(ctx, groupPath, nil)
			if err == nil {
				var groupResp map[string]interface{}
				if json.Unmarshal(groupData, &groupResp) == nil {
					groupMetadata = map[string]interface{}{
						"group_id":      groupID,
						"group_name":    groupResp["group_name"],
						"meta_protocol": groupResp["meta_protocol"],
						"meta_server":   groupResp["meta_server"],
						"meta_port":     groupResp["meta_port"],
						"meta_user":     groupResp["meta_user"],
						"meta_password": groupResp["meta_password"],
					}
				}
			}

			// 获取执行任务的用户自定义变量（优先任务独立变量，没有则使用用例集变量）
			varsPath := fmt.Sprintf("/api/v1/projects/%d/execution-tasks/%s/variables", projectID, taskUUID)
			varsParams := map[string]string{
				"group_id":   fmt.Sprintf("%d", groupID),
				"group_type": groupType,
			}
			varsData, err := h.client.Get(ctx, varsPath, varsParams)
			if err == nil {
				var varsResp map[string]interface{}
				if json.Unmarshal(varsData, &varsResp) == nil {
					if vars, ok := varsResp["variables"].([]interface{}); ok {
						taskVariables = vars
					}
				}
			}
		}
	}

	// 步骤5: 构建元数据结果
	resultData := map[string]interface{}{
		"task_uuid":     taskUUID,
		"task_metadata": taskMetadata,
		"statistics": map[string]interface{}{
			"total_cases":   totalCases,
			"ok_count":      okCount,
			"ng_count":      ngCount,
			"nr_count":      nrCount,
			"block_count":   blockCount,
			"executed":      executedCount,
			"progress_rate": fmt.Sprintf("%.1f%%", progressRate),
			"pass_rate":     fmt.Sprintf("%.1f%%", passRate),
		},
	}

	// 如果获取到了用例集元数据，添加到结果中
	if groupMetadata != nil {
		resultData["group_metadata"] = groupMetadata
	}

	// 如果获取到了用户自定义变量，添加到结果中
	if len(taskVariables) > 0 {
		resultData["variables"] = taskVariables
	}

	result := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    resultData,
	}

	resultJSON, _ := json.Marshal(result)
	return tools.NewJSONResult(string(resultJSON)), nil
}

// GetExecutionTaskCasesHandler handles getting cases for an execution task.
type GetExecutionTaskCasesHandler struct {
	*BaseHandler
}

func NewGetExecutionTaskCasesHandler(c *client.BackendClient) *GetExecutionTaskCasesHandler {
	return &GetExecutionTaskCasesHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *GetExecutionTaskCasesHandler) Name() string {
	return "get_execution_task_cases"
}

func (h *GetExecutionTaskCasesHandler) Description() string {
	return "获取指定执行任务的全部用例列表（含所有用例字段：中文、日文、英文描述、前置条件、测试步骤、期望结果、执行状态、备注、script_code脚本代码等）"
}

func (h *GetExecutionTaskCasesHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"task_id": map[string]interface{}{
				"type":        "string",
				"description": "执行任务标识。支持三种格式：1) 任务名字字符串(如: qweb, 手工, api2)；2) UUID (格式: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)；3) 数字字符串索引（如 '1' 表示第一个任务）",
			},
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID（可选，如不提供则默认为1）",
			},
		},
		"required": []interface{}{"task_id"},
	}
}

func (h *GetExecutionTaskCasesHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	// 处理task_id参数
	var taskID string
	var isNameLookup bool
	var isIndexLookup bool
	var taskIndex int = -1
	projectID := GetOptionalInt(args, "project_id", 1)

	if val, ok := args["task_id"]; ok {
		switch v := val.(type) {
		case float64:
			taskIndex = int(v)
			taskID = fmt.Sprintf("%d", taskIndex)
			isIndexLookup = true
		case int:
			taskIndex = v
			taskID = fmt.Sprintf("%d", v)
			isIndexLookup = true
		case int64:
			taskIndex = int(v)
			taskID = fmt.Sprintf("%d", v)
			isIndexLookup = true
		case string:
			if idx, err := strconv.Atoi(v); err == nil && idx > 0 {
				taskIndex = idx
				taskID = v
				isIndexLookup = true
			} else if isUUIDFormat(v) {
				taskID = v
			} else {
				taskID = v
				isNameLookup = true
			}
		default:
			taskID = fmt.Sprintf("%v", v)
			isNameLookup = true
		}
	} else {
		return tools.NewErrorResult("task_id parameter is required"), nil
	}

	// 步骤1: 获取项目的执行任务列表
	taskListPath := fmt.Sprintf("/api/v1/projects/%d/execution-tasks", projectID)
	taskListData, err := h.client.Get(ctx, taskListPath, nil)
	if err != nil {
		return tools.NewErrorResult("failed to list execution tasks: " + err.Error()), nil
	}

	// 解析任务列表
	var taskList struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(taskListData, &taskList); err != nil {
		return tools.NewErrorResult("failed to parse task list: " + err.Error()), nil
	}

	// 解析任务数组
	var tasks []struct {
		TaskUUID string `json:"task_uuid"`
		TaskName string `json:"task_name"`
	}

	if err := json.Unmarshal(taskList.Data, &tasks); err != nil {
		return tools.NewErrorResult("failed to parse tasks array: " + err.Error()), nil
	}

	// 查找匹配的任务
	var taskUUID string

	if isIndexLookup {
		if taskIndex < 1 || taskIndex > len(tasks) {
			return tools.NewErrorResult(fmt.Sprintf("task index %d out of range, project %d has %d tasks", taskIndex, projectID, len(tasks))), nil
		}
		taskUUID = tasks[taskIndex-1].TaskUUID
	} else if isNameLookup {
		found := false
		for _, task := range tasks {
			if task.TaskName == taskID {
				taskUUID = task.TaskUUID
				found = true
				break
			}
		}
		if !found {
			return tools.NewErrorResult(fmt.Sprintf("task '%s' not found in project %d", taskID, projectID)), nil
		}
	} else {
		found := false
		for _, task := range tasks {
			if task.TaskUUID == taskID {
				taskUUID = task.TaskUUID
				found = true
				break
			}
		}
		if !found {
			return tools.NewErrorResult(fmt.Sprintf("task with UUID '%s' not found in project %d", taskID, projectID)), nil
		}
	}

	// 步骤2: 获取执行任务的全部用例
	casesPath := fmt.Sprintf("/api/v1/execution-tasks/%s/case-results", taskUUID)
	casesParams := map[string]string{
		"size": "99999",
	}
	casesData, err := h.client.Get(ctx, casesPath, casesParams)
	if err != nil {
		return tools.NewErrorResult("failed to get case results: " + err.Error()), nil
	}

	return tools.NewJSONResult(string(casesData)), nil
}

// UpdateExecutionCaseResultHandler handles updating execution case result.
// Supports both single update (id + result) and batch update (updates array)
type UpdateExecutionCaseResultHandler struct {
	*BaseHandler
}

func NewUpdateExecutionCaseResultHandler(c *client.BackendClient) *UpdateExecutionCaseResultHandler {
	return &UpdateExecutionCaseResultHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *UpdateExecutionCaseResultHandler) Name() string {
	return "update_execution_case_result"
}

func (h *UpdateExecutionCaseResultHandler) Description() string {
	return "更新执行任务中用例的执行结果，支持单个更新和批量更新"
}

func (h *UpdateExecutionCaseResultHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "执行用例记录ID（单个更新时使用）",
			},
			"result": map[string]interface{}{
				"type":        "string",
				"description": "执行结果（支持：NR-未执行, OK-通过, NG-失败, Block-阻止）",
				"enum":        []interface{}{"NR", "OK", "NG", "Block"},
			},
			"comment": map[string]interface{}{
				"type":        "string",
				"description": "执行备注/备注（可选，与remark字段等价）",
			},
			"remark": map[string]interface{}{
				"type":        "string",
				"description": "执行备注/备注（可选，与comment字段等价）",
			},
			"bug_id": map[string]interface{}{
				"type":        "string",
				"description": "缺陷ID（可选，用于关联测试缺陷）",
			},
			"response_time": map[string]interface{}{
				"type":        "string",
				"description": "响应时间（可选，如：125ms, 1.5s）",
			},
			"updates": map[string]interface{}{
				"type":        "array",
				"description": "批量更新时的用例数组，每个元素包含id、result和可选的comment/remark/bug_id",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "integer",
							"description": "执行用例记录ID",
						},
						"result": map[string]interface{}{
							"type":        "string",
							"description": "执行结果（支持：NR-未执行, OK-通过, NG-失败, Block-阻止）",
							"enum":        []interface{}{"NR", "OK", "NG", "Block"},
						},
						"comment": map[string]interface{}{
							"type":        "string",
							"description": "执行备注/备注（可选，与remark字段等价）",
						},
						"remark": map[string]interface{}{
							"type":        "string",
							"description": "执行备注/备注（可选，与comment字段等价）",
						},
						"bug_id": map[string]interface{}{
							"type":        "string",
							"description": "缺陷ID（可选，用于关联测试缺陷）",
						},
						"response_time": map[string]interface{}{
							"type":        "string",
							"description": "响应时间（可选，如：125ms, 1.5s）",
						},
					},
					"required": []interface{}{"id", "result"},
				},
			},
		},
	}
}

func (h *UpdateExecutionCaseResultHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	// 检查是否为批量更新模式
	if updates, ok := args["updates"]; ok && updates != nil {
		return h.executeBatch(ctx, updates)
	}

	// 单个更新模式
	return h.executeSingle(ctx, args)
}

// executeSingle 执行单个用例更新
func (h *UpdateExecutionCaseResultHandler) executeSingle(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	id, err := GetInt(args, "id")
	if err != nil {
		return tools.NewErrorResult("单个更新模式需要提供id参数: " + err.Error()), nil
	}

	result, err := GetString(args, "result")
	if err != nil {
		return tools.NewErrorResult("单个更新模式需要提供result参数: " + err.Error()), nil
	}

	body := map[string]interface{}{
		"result": result,
	}

	// 支持 comment 和 remark 两个参数名，优先使用 remark
	if remark := GetOptionalString(args, "remark", ""); remark != "" {
		body["comment"] = remark
	} else if comment := GetOptionalString(args, "comment", ""); comment != "" {
		body["comment"] = comment
	}

	// 支持 bug_id 参数
	if bugID := GetOptionalString(args, "bug_id", ""); bugID != "" {
		body["bug_id"] = bugID
	}

	// 支持 response_time 参数
	if responseTime := GetOptionalString(args, "response_time", ""); responseTime != "" {
		body["response_time"] = responseTime
	}

	path := fmt.Sprintf("/api/v1/execution-task-cases/%d", id)
	data, err := h.client.Put(ctx, path, body)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// executeBatch 执行批量更新
func (h *UpdateExecutionCaseResultHandler) executeBatch(ctx context.Context, updates interface{}) (tools.ToolResult, error) {
	updatesSlice, ok := updates.([]interface{})
	if !ok {
		return tools.NewErrorResult("updates参数必须是数组"), nil
	}

	if len(updatesSlice) == 0 {
		return tools.NewErrorResult("updates数组不能为空"), nil
	}

	// 构建批量更新请求体
	batchUpdates := make([]map[string]interface{}, 0, len(updatesSlice))
	for i, item := range updatesSlice {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return tools.NewErrorResult(fmt.Sprintf("updates[%d]格式错误", i)), nil
		}

		// 获取id
		idVal, ok := itemMap["id"]
		if !ok {
			return tools.NewErrorResult(fmt.Sprintf("updates[%d]缺少id字段", i)), nil
		}
		id, err := toInt(idVal)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("updates[%d].id无效: %v", i, err)), nil
		}

		// 获取result
		resultVal, ok := itemMap["result"]
		if !ok {
			return tools.NewErrorResult(fmt.Sprintf("updates[%d]缺少result字段", i)), nil
		}
		result, ok := resultVal.(string)
		if !ok {
			return tools.NewErrorResult(fmt.Sprintf("updates[%d].result必须是字符串", i)), nil
		}

		update := map[string]interface{}{
			"id":     id,
			"result": result,
		}

		// 获取可选的comment/remark字段，优先使用remark
		if remarkVal, ok := itemMap["remark"]; ok && remarkVal != nil {
			if remark, ok := remarkVal.(string); ok && remark != "" {
				update["comment"] = remark
			}
		} else if commentVal, ok := itemMap["comment"]; ok && commentVal != nil {
			if comment, ok := commentVal.(string); ok && comment != "" {
				update["comment"] = comment
			}
		}

		// 获取可选的bug_id字段
		if bugIDVal, ok := itemMap["bug_id"]; ok && bugIDVal != nil {
			if bugID, ok := bugIDVal.(string); ok && bugID != "" {
				update["bug_id"] = bugID
			}
		}

		// 获取可选的response_time字段
		if responseTimeVal, ok := itemMap["response_time"]; ok && responseTimeVal != nil {
			if responseTime, ok := responseTimeVal.(string); ok && responseTime != "" {
				update["response_time"] = responseTime
			}
		}

		batchUpdates = append(batchUpdates, update)
	}

	path := "/api/v1/execution-task-cases/batch"
	data, err := h.client.Put(ctx, path, batchUpdates)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// toInt 将各种类型转换为int
func toInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return i, nil
	default:
		return 0, fmt.Errorf("无法将 %T 转换为int", val)
	}
}

// isUUIDFormat 检查字符串是否为UUID格式（包含4个"-"）
func isUUIDFormat(s string) bool {
	count := 0
	for _, c := range s {
		if c == '-' {
			count++
		}
	}
	return count == 4 && len(s) == 36
}
