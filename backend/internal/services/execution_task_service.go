package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"gorm.io/gorm"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	TaskName      string `json:"task_name" binding:"required,min=1,max=50"`
	ExecutionType string `json:"execution_type" binding:"required,oneof=manual automation api"`
	TaskStatus    string `json:"task_status" binding:"omitempty,oneof=pending in_progress completed"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	TaskName        *string    `json:"task_name" binding:"omitempty,min=1,max=50"`
	ExecutionType   *string    `json:"execution_type" binding:"omitempty,oneof=manual automation api"`
	TaskStatus      *string    `json:"task_status" binding:"omitempty,oneof=pending in_progress completed"`
	CaseGroupID     *uint      `json:"case_group_id"`                               // 关联的用例集ID
	CaseGroupName   *string    `json:"case_group_name" binding:"omitempty,max=100"` // 关联的用例集名称
	DisplayLanguage *string    `json:"display_language" binding:"omitempty,max=10"` // 显示语言(cn/jp/en/all)
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	TestVersion     *string    `json:"test_version" binding:"omitempty,max=50"`
	TestEnv         *string    `json:"test_env" binding:"omitempty,max=100"`
	TestDate        *time.Time `json:"test_date"`
	Executor        *string    `json:"executor" binding:"omitempty,max=50"`
	TaskDescription *string    `json:"task_description" binding:"omitempty,max=2000"`
}

// ExecuteTaskResult 执行结果统计
type ExecuteTaskResult struct {
	Total      int       `json:"total"`
	OKCount    int       `json:"ok_count"`
	NGCount    int       `json:"ng_count"`
	BlockCount int       `json:"block_count"`
	ExecutedAt time.Time `json:"executed_at"`
	ExecutedBy string    `json:"executed_by"`
}

// DockerExecResult Docker执行结果
type DockerExecResult struct {
	Success      bool
	Output       string
	ResponseTime int // 毫秒
}

// ExecutionTaskService 测试执行任务服务接口
type ExecutionTaskService interface {
	GetTasksByProject(projectID uint, userID uint) ([]*models.ExecutionTask, error)
	CreateTask(projectID uint, userID uint, req CreateTaskRequest) (*models.ExecutionTask, error)
	UpdateTask(projectID uint, userID uint, taskUUID string, req UpdateTaskRequest) (*models.ExecutionTask, error)
	DeleteTask(projectID uint, userID uint, taskUUID string) error
	ExecuteTask(projectID uint, userID uint, taskUUID string) (*ExecuteTaskResult, error)
	ExecuteSingleCase(projectID uint, userID uint, taskUUID string, caseResultID uint) (*ExecuteTaskResult, error)
}

type executionTaskService struct {
	repo            repositories.ExecutionTaskRepository
	projectRepo     repositories.ProjectRepository             // 用于权限验证
	ecrRepo         repositories.ExecutionCaseResultRepository // 用于级联删除
	userRepo        repositories.UserRepository                // 用于获取用户名
	pwClient        *PlaywrightExecutorClient                  // Playwright 执行器客户端
	variableService UserDefinedVariableService                 // 用户自定义变量服务
}

// NewExecutionTaskService 创建任务服务实例
func NewExecutionTaskService(
	repo repositories.ExecutionTaskRepository,
	projectRepo repositories.ProjectRepository,
	ecrRepo repositories.ExecutionCaseResultRepository,
	userRepo repositories.UserRepository,
	variableService UserDefinedVariableService,
) ExecutionTaskService {
	// 初始化 Playwright 执行器客户端
	pwClient := NewPlaywrightExecutorClient(DefaultExecutorConfig())

	return &executionTaskService{
		repo:            repo,
		projectRepo:     projectRepo,
		ecrRepo:         ecrRepo,
		userRepo:        userRepo,
		pwClient:        pwClient,
		variableService: variableService,
	}
}

// GetTasksByProject 获取项目的所有任务
func (s *executionTaskService) GetTasksByProject(projectID uint, userID uint) ([]*models.ExecutionTask, error) {
	// 验证用户项目权限(中间件已验证,此处可选)
	// 可选: 调用projectRepo验证项目是否存在

	tasks, err := s.repo.GetByProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("get tasks by project %d: %w", projectID, err)
	}
	return tasks, nil
}

// CreateTask 创建新任务
func (s *executionTaskService) CreateTask(projectID uint, userID uint, req CreateTaskRequest) (*models.ExecutionTask, error) {
	// 1. 验证任务名唯一性(不区分大小写)
	isUnique, err := s.repo.CheckNameUnique(projectID, req.TaskName, "")
	if err != nil {
		return nil, fmt.Errorf("check name unique: %w", err)
	}
	if !isUnique {
		return nil, errors.New("任务名已存在")
	}

	// 2. 构建任务对象
	task := &models.ExecutionTask{
		ProjectID:     projectID,
		TaskName:      req.TaskName,
		ExecutionType: req.ExecutionType,
		TaskStatus:    "pending", // 默认状态
		CreatedBy:     userID,
	}

	// 如果请求中指定了状态,使用指定值
	if req.TaskStatus != "" {
		task.TaskStatus = req.TaskStatus
	}

	// 3. 创建任务(BeforeCreate Hook会生成UUID)
	err = s.repo.Create(task)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

// UpdateTask 更新任务
func (s *executionTaskService) UpdateTask(projectID uint, userID uint, taskUUID string, req UpdateTaskRequest) (*models.ExecutionTask, error) {
	fmt.Printf("\n========== [UpdateTask START] ==========\n")
	fmt.Printf("[UpdateTask] projectID=%d, userID=%d, taskUUID=%s\n", projectID, userID, taskUUID)
	fmt.Printf("[UpdateTask] req.CaseGroupName=%v\n", req.CaseGroupName)
	if req.CaseGroupName != nil {
		fmt.Printf("[UpdateTask] req.CaseGroupName value=%s\n", *req.CaseGroupName)
	}
	fmt.Printf("[UpdateTask] req.CaseGroupID=%v\n", req.CaseGroupID)
	fmt.Printf("[UpdateTask] req.DisplayLanguage=%v\n", req.DisplayLanguage)

	// 1. 验证任务存在且属于该项目
	task, err := s.repo.GetByUUID(taskUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, fmt.Errorf("get task by uuid: %w", err)
	}

	fmt.Printf("[UpdateTask] Found existing task: case_group_id=%d, case_group_name=%s\n", task.CaseGroupID, task.CaseGroupName)

	if task.ProjectID != projectID {
		return nil, errors.New("任务不属于该项目")
	}

	// 2. 验证任务名唯一性(如果修改了名称)
	if req.TaskName != nil && *req.TaskName != task.TaskName {
		isUnique, err := s.repo.CheckNameUnique(projectID, *req.TaskName, taskUUID)
		if err != nil {
			return nil, fmt.Errorf("check name unique: %w", err)
		}
		if !isUnique {
			return nil, errors.New("任务名已存在")
		}
	}

	// 3. 验证日期范围逻辑
	startDate := task.StartDate
	endDate := task.EndDate

	if req.StartDate != nil {
		startDate = req.StartDate
	}
	if req.EndDate != nil {
		endDate = req.EndDate
	}

	if startDate != nil && endDate != nil {
		if endDate.Before(*startDate) {
			return nil, errors.New("结束日期不能早于开始日期")
		}
	}

	// 4. 构建更新字段Map
	updates := make(map[string]interface{})

	if req.TaskName != nil {
		updates["task_name"] = *req.TaskName
		fmt.Printf("[UpdateTask] Adding task_name to updates: %s\n", *req.TaskName)
	}
	if req.ExecutionType != nil {
		updates["execution_type"] = *req.ExecutionType
		fmt.Printf("[UpdateTask] Adding execution_type to updates: %s\n", *req.ExecutionType)
	}
	if req.TaskStatus != nil {
		updates["task_status"] = *req.TaskStatus
		fmt.Printf("[UpdateTask] Adding task_status to updates: %s\n", *req.TaskStatus)
	}
	if req.CaseGroupID != nil {
		updates["case_group_id"] = *req.CaseGroupID
		fmt.Printf("[UpdateTask] Adding case_group_id to updates: %d\n", *req.CaseGroupID)
	}
	if req.CaseGroupName != nil {
		updates["case_group_name"] = *req.CaseGroupName
		fmt.Printf("[UpdateTask] Adding case_group_name to updates: %s\n", *req.CaseGroupName)
	}
	if req.DisplayLanguage != nil {
		updates["display_language"] = *req.DisplayLanguage
		fmt.Printf("[UpdateTask] Adding display_language to updates: %s\n", *req.DisplayLanguage)
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = *req.EndDate
	}
	if req.TestVersion != nil {
		updates["test_version"] = *req.TestVersion
	}
	if req.TestEnv != nil {
		updates["test_env"] = *req.TestEnv
	}
	if req.TestDate != nil {
		updates["test_date"] = *req.TestDate
	}
	if req.Executor != nil {
		updates["executor"] = *req.Executor
	}
	if req.TaskDescription != nil {
		updates["task_description"] = *req.TaskDescription
	}

	fmt.Printf("[UpdateTask] Total updates map size: %d\n", len(updates))
	fmt.Printf("[UpdateTask] Updates map: %+v\n", updates)

	// 5. 执行更新
	fmt.Printf("[UpdateTask] Calling repo.UpdateByUUID...\n")
	err = s.repo.UpdateByUUID(taskUUID, updates)
	if err != nil {
		fmt.Printf("[UpdateTask] ERROR: Failed to update task: %v\n", err)
		return nil, fmt.Errorf("update task: %w", err)
	}
	fmt.Printf("[UpdateTask] Successfully updated task in database\n")

	// 5.5. 如果更新了case_group_name或case_group_id，从用例集继承变量到任务变量表
	// 注意：case_group_id是从前端通过case_group_name查找得到的，所以只要case_group_name变化就触发继承
	fmt.Printf("[UpdateTask] Checking if need to inherit variables...\n")
	fmt.Printf("[UpdateTask] req.CaseGroupName != nil: %v\n", req.CaseGroupName != nil)
	if req.CaseGroupName != nil {
		fmt.Printf("[UpdateTask] CaseGroupName value: '%s', isEmpty: %v\n", *req.CaseGroupName, *req.CaseGroupName == "")
	}

	if req.CaseGroupName != nil && *req.CaseGroupName != "" {
		fmt.Printf("[UpdateTask] ✅ case_group_name changed, inheriting variables from case group: %s\n", *req.CaseGroupName)

		// 6. 查询更新后的任务（包含最新的case_group_id）
		fmt.Printf("[UpdateTask] Fetching updated task to get case_group_id...\n")
		updatedTask, err := s.repo.GetByUUID(taskUUID)
		if err != nil {
			fmt.Printf("[UpdateTask] ERROR: Failed to get updated task: %v\n", err)
			return nil, fmt.Errorf("get updated task: %w", err)
		}

		fmt.Printf("[UpdateTask] Updated task retrieved: case_group_id=%d, case_group_name=%s, execution_type=%s\n",
			updatedTask.CaseGroupID, updatedTask.CaseGroupName, updatedTask.ExecutionType)

		// 如果任务关联了用例集，继承变量
		if updatedTask.CaseGroupID > 0 {
			groupType := "web" // 默认web类型
			if updatedTask.ExecutionType == "api" {
				groupType = "api"
			}

			fmt.Printf("[UpdateTask] ✅ Task has case_group_id, starting variable inheritance...\n")
			fmt.Printf("[UpdateTask] Parameters: group_id=%d, group_type=%s, task_uuid=%s, project_id=%d\n",
				updatedTask.CaseGroupID, groupType, taskUUID, projectID)

			// 先获取用例集的变量
			fmt.Printf("[UpdateTask] Step 1: Getting variables from case group...\n")
			groupVariables, err := s.variableService.GetVariablesByGroup(updatedTask.CaseGroupID, groupType)
			if err != nil {
				fmt.Printf("[UpdateTask] ❌ ERROR: Failed to get case group variables: %v\n", err)
				// 继承变量失败不影响任务更新，只记录日志
			} else {
				fmt.Printf("[UpdateTask] ✅ Successfully retrieved %d variables from case group\n", len(groupVariables))
				if len(groupVariables) > 0 {
					fmt.Printf("[UpdateTask] Variables from case group:\n")
					for i, v := range groupVariables {
						fmt.Printf("[UpdateTask]   [%d] ID=%d, Key=%s, Name=%s, Value=%s\n", i, v.ID, v.VarKey, v.VarName, v.VarValue)
					}

					fmt.Printf("[UpdateTask] Step 2: Copying variables to task variable table...\n")
					// 复制变量到任务变量表（清空任务旧变量，使用用例集的最新变量）
					err = s.variableService.SaveTaskVariables(taskUUID, updatedTask.CaseGroupID, groupType, projectID, groupVariables)
					if err != nil {
						fmt.Printf("[UpdateTask] ❌ ERROR: Failed to save task variables: %v\n", err)
						// 继承变量失败不影响任务更新，只记录日志
					} else {
						fmt.Printf("[UpdateTask] ✅✅✅ Successfully inherited %d variables to task!\n", len(groupVariables))
					}
				} else {
					fmt.Printf("[UpdateTask] ⚠️ WARNING: No variables found in case group (empty array)\n")
				}
			}
		} else {
			fmt.Printf("[UpdateTask] ⚠️ WARNING: case_group_id is 0 or invalid, skipping variable inheritance\n")
		}

		fmt.Printf("[UpdateTask] Returning updated task\n")
		fmt.Printf("========== [UpdateTask END] ==========\n\n")
		return updatedTask, nil
	} else {
		fmt.Printf("[UpdateTask] ℹ️ case_group_name not changed, skipping variable inheritance\n")
	}

	// 6. 查询更新后的任务
	fmt.Printf("[UpdateTask] Fetching final updated task...\n")
	updatedTask, err := s.repo.GetByUUID(taskUUID)
	if err != nil {
		fmt.Printf("[UpdateTask] ERROR: Failed to get updated task: %v\n", err)
		return nil, fmt.Errorf("get updated task: %w", err)
	}

	fmt.Printf("[UpdateTask] Returning updated task\n")
	fmt.Printf("========== [UpdateTask END] ==========\n\n")
	return updatedTask, nil
}

// DeleteTask 删除任务(硬删除，级联删除相关数据)
func (s *executionTaskService) DeleteTask(projectID uint, userID uint, taskUUID string) error {
	// 1. 验证任务存在且属于该项目
	task, err := s.repo.GetByUUID(taskUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("任务不存在")
		}
		return fmt.Errorf("get task by uuid: %w", err)
	}

	if task.ProjectID != projectID {
		return errors.New("任务不属于该项目")
	}

	// 2. 级联删除：先删除执行任务的所有执行用例结果（元数据）
	err = s.ecrRepo.DeleteByTaskUUID(taskUUID)
	if err != nil {
		return fmt.Errorf("delete execution case results: %w", err)
	}

	// 3. 执行硬删除执行任务
	err = s.repo.Delete(taskUUID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	// 可选: 记录审计日志
	// TODO: 添加日志记录

	return nil
}

// 辅助函数: 检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

// ExecuteTask 执行测试任务
func (s *executionTaskService) ExecuteTask(projectID uint, userID uint, taskUUID string) (*ExecuteTaskResult, error) {
	fmt.Printf("[ExecuteTask] 开始执行任务: projectID=%d, userID=%d, taskUUID=%s\n", projectID, userID, taskUUID)

	// 1. 获取任务信息
	task, err := s.repo.GetByUUID(taskUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, fmt.Errorf("get task by uuid: %w", err)
	}
	if task.ProjectID != projectID {
		return nil, errors.New("任务不属于该项目")
	}
	fmt.Printf("[ExecuteTask] 任务信息: task_name=%s, execution_type=%s\n", task.TaskName, task.ExecutionType)

	// 2. 检查执行类型
	if task.ExecutionType == "manual" {
		return nil, errors.New("手工测试类型不支持自动执行")
	}

	// 3. 获取用例列表
	cases, err := s.ecrRepo.GetByTaskUUID(taskUUID)
	if err != nil {
		return nil, fmt.Errorf("get case results: %w", err)
	}
	if len(cases) == 0 {
		return nil, errors.New("没有可执行的用例")
	}
	fmt.Printf("[ExecuteTask] 获取到 %d 个用例\n", len(cases))

	// 3B. 获取任务变量用于脚本替换
	var variables []*models.UserDefinedVariable
	if task.CaseGroupID > 0 {
		groupType := "web"
		if task.ExecutionType == "api" {
			groupType = "api"
		}
		fmt.Printf("[ExecuteTask] 从任务变量表获取变量: taskUUID=%s, groupID=%d, groupType=%s\n", taskUUID, task.CaseGroupID, groupType)
		variables, err = s.variableService.GetVariablesByTask(taskUUID, task.CaseGroupID, groupType)
		if err != nil {
			fmt.Printf("[ExecuteTask] ❌ 警告: 获取变量失败: %v\n", err)
			variables = []*models.UserDefinedVariable{} // 继续执行，不中断
		}
		fmt.Printf("[ExecuteTask] ✅ 获取到 %d 个任务变量\n", len(variables))
		for i, v := range variables {
			fmt.Printf("[ExecuteTask]   变量 %d: var_key=%s, var_value=%s (长度:%d)\n", i+1, v.VarKey, maskValue(v.VarKey, v.VarValue), len(v.VarValue))
		}
	} else {
		fmt.Printf("[ExecuteTask] 跳过变量获取 (case_group_id=%d)\n", task.CaseGroupID)
	}

	// 4. 确定remark语言
	lang := task.DisplayLanguage
	if lang == "" {
		lang = "cn"
	}

	// 5. 逐个执行用例
	var okCount, ngCount, blockCount int
	for i, c := range cases {
		fmt.Printf("[ExecuteTask] 执行用例 %d/%d: case_id=%d, has_script=%v\n",
			i+1, len(cases), c.ID, c.ScriptCode != "")

		if c.ScriptCode == "" {
			blockCount++
			continue
		}

		// 替换脚本中的变量
		fmt.Printf("[ExecuteTask] 用例 %d 脚本替换前长度: %d bytes\n", c.ID, len(c.ScriptCode))
		replacedScript := s.replaceVariables(c.ScriptCode, variables)
		fmt.Printf("[ExecuteTask] 用例 %d 脚本替换后长度: %d bytes\n", c.ID, len(replacedScript))
		if strings.Contains(replacedScript, "${") {
			fmt.Printf("[ExecuteTask] ⚠️ 用例 %d 脚本仍包含变量占位符！\n", c.ID)
		}

		// 调用 Playwright Server 执行脚本
		fmt.Printf("[ExecuteTask] 开始执行 Playwright 脚本...\n")
		execResult, execErr := s.executeViaPlaywright(replacedScript)
		if execErr != nil {
			// 执行失败
			fmt.Printf("[ExecuteTask] 执行失败: %v\n", execErr)
			c.TestResult = "NG"
			c.Remark = s.getRemarkByLang(lang, false, execErr.Error())
			ngCount++
		} else {
			fmt.Printf("[ExecuteTask] 执行成功: response_time=%dms\n", execResult.ResponseTime)
			c.TestResult = "OK"
			c.Remark = s.getRemarkByLang(lang, true, "")
			okCount++
		}
		// 记录执行时间（无论成功失败，只要有结果就记录）
		if execResult != nil && execResult.ResponseTime > 0 {
			c.ResponseTime = fmt.Sprintf("%d", execResult.ResponseTime)
		}

		// 更新数据库
		updates := map[string]interface{}{
			"test_result": c.TestResult,
			"remark":      c.Remark,
		}
		if c.ResponseTime != "" {
			updates["response_time"] = c.ResponseTime
		}
		if updateErr := s.ecrRepo.UpdateResult(c.ID, updates); updateErr != nil {
			return nil, fmt.Errorf("update case result: %w", updateErr)
		}
	}

	fmt.Printf("[ExecuteTask] 执行完成: OK=%d, NG=%d, Block=%d\n", okCount, ngCount, blockCount)

	// 6. 更新任务的测试日期和执行人
	now := time.Now()
	executor := s.getUserName(userID)
	updates := map[string]interface{}{
		"test_date": now,
		"executor":  executor,
	}
	if updateErr := s.repo.UpdateByUUID(taskUUID, updates); updateErr != nil {
		return nil, fmt.Errorf("update task: %w", updateErr)
	}

	return &ExecuteTaskResult{
		Total:      len(cases),
		OKCount:    okCount,
		NGCount:    ngCount,
		BlockCount: blockCount,
		ExecutedAt: now,
		ExecutedBy: executor,
	}, nil
}

// maskValue masks sensitive variable values
func maskValue(key, value string) string {
	lowerKey := strings.ToLower(key)
	if strings.Contains(lowerKey, "password") || strings.Contains(lowerKey, "secret") || strings.Contains(lowerKey, "token") {
		if len(value) <= 3 {
			return "***"
		}
		return value[:2] + "***"
	}
	return value
}

// replaceVariables 替换脚本中的变量占位符
// 支持 ${VAR_NAME} 格式（脚本标准格式）和 {{VAR_NAME}} 格式（兼容旧格式）
func (s *executionTaskService) replaceVariables(script string, variables []*models.UserDefinedVariable) string {
	if len(variables) == 0 {
		return script
	}

	result := script
	for _, v := range variables {
		// 1. 使用 VarKey 替换大写格式 "${BASE_URL}"
		if v.VarKey != "" {
			upperKey := strings.ToUpper(v.VarKey)
			placeholder := fmt.Sprintf("${%s}", upperKey)
			result = strings.ReplaceAll(result, placeholder, v.VarValue)
		}

		// 2. 同时支持小写格式 "${base_url}" (增强兼容性)
		if v.VarKey != "" {
			lowerKey := strings.ToLower(v.VarKey)
			placeholder := fmt.Sprintf("${%s}", lowerKey)
			result = strings.ReplaceAll(result, placeholder, v.VarValue)
		}

		// 3. 兼容：使用 VarName 字段
		if v.VarName != "" {
			result = strings.ReplaceAll(result, v.VarName, v.VarValue)
		}

		// 4. 兼容：{{key}} 格式
		if v.VarKey != "" {
			placeholder := fmt.Sprintf("{{%s}}", v.VarKey)
			result = strings.ReplaceAll(result, placeholder, v.VarValue)
		}
	}

	return result
}

// getRemarkByLang 根据语言生成remark
func (s *executionTaskService) getRemarkByLang(lang string, success bool, errMsg string) string {
	if success {
		switch lang {
		case "jp":
			return "自動実行成功"
		case "en":
			return "Auto execution succeeded"
		default:
			return "自动执行成功"
		}
	}

	switch lang {
	case "jp":
		return fmt.Sprintf("自動実行失敗: %s", errMsg)
	case "en":
		return fmt.Sprintf("Auto execution failed: %s", errMsg)
	default:
		return fmt.Sprintf("自动执行失败: %s", errMsg)
	}
}

// executeViaPlaywright 通过 Playwright Server 执行脚本
// 使用 HTTP 调用 playwright-executor 服务
func (s *executionTaskService) executeViaPlaywright(scriptCode string) (*DockerExecResult, error) {
	fmt.Printf("[executeViaPlaywright] 开始执行脚本，长度: %d bytes\n", len(scriptCode))

	ctx := context.Background()
	result, err := s.pwClient.Execute(ctx, scriptCode)
	if err != nil {
		fmt.Printf("[executeViaPlaywright] 执行失败: %v\n", err)
		return nil, err
	}

	fmt.Printf("[executeViaPlaywright] 执行完成，耗时: %dms\n", result.ResponseTime)
	return result, nil
}

// getUserName 获取用户名
func (s *executionTaskService) getUserName(userID uint) string {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return fmt.Sprintf("user_%d", userID)
	}
	if user.Nickname != "" {
		return user.Nickname
	}
	return user.Username
}

// ExecuteSingleCase 执行单条测试用例
func (s *executionTaskService) ExecuteSingleCase(projectID uint, userID uint, taskUUID string, caseResultID uint) (*ExecuteTaskResult, error) {
	fmt.Printf("[ExecuteSingleCase] 开始执行单条用例: projectID=%d, userID=%d, taskUUID=%s, caseResultID=%d\n", projectID, userID, taskUUID, caseResultID)

	// 1. 获取任务信息
	task, err := s.repo.GetByUUID(taskUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, fmt.Errorf("get task by uuid: %w", err)
	}
	if task.ProjectID != projectID {
		return nil, errors.New("任务不属于该项目")
	}
	fmt.Printf("[ExecuteSingleCase] 任务信息: task_name=%s, execution_type=%s\n", task.TaskName, task.ExecutionType)

	// 2. 检查执行类型
	if task.ExecutionType == "manual" {
		return nil, errors.New("手工测试类型不支持自动执行")
	}

	// 3. 获取指定用例
	c, err := s.ecrRepo.GetByID(caseResultID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用例不存在")
		}
		return nil, fmt.Errorf("get case result: %w", err)
	}
	if c.TaskUUID != taskUUID {
		return nil, errors.New("用例不属于该任务")
	}
	fmt.Printf("[ExecuteSingleCase] 用例信息: case_id=%d, has_script=%v\n", c.ID, c.ScriptCode != "")

	// 4. 检查是否有脚本
	if c.ScriptCode == "" {
		return nil, errors.New("用例没有脚本代码，无法执行")
	}

	// 4B. 获取任务变量用于脚本替换
	var variables []*models.UserDefinedVariable
	if task.CaseGroupID > 0 {
		groupType := "web"
		if task.ExecutionType == "api" {
			groupType = "api"
		}
		fmt.Printf("[ExecuteSingleCase] 从任务变量表获取变量: taskUUID=%s, groupID=%d, groupType=%s\n", taskUUID, task.CaseGroupID, groupType)
		variables, err = s.variableService.GetVariablesByTask(taskUUID, task.CaseGroupID, groupType)
		if err != nil {
			fmt.Printf("[ExecuteSingleCase] ❌ 警告: 获取变量失败: %v\n", err)
			variables = []*models.UserDefinedVariable{} // 继续执行，不中断
		}
		fmt.Printf("[ExecuteSingleCase] ✅ 获取到 %d 个任务变量\n", len(variables))
		for i, v := range variables {
			fmt.Printf("[ExecuteSingleCase]   变量 %d: var_key=%s, var_value=%s (长度:%d)\n", i+1, v.VarKey, maskValue(v.VarKey, v.VarValue), len(v.VarValue))
		}
	} else {
		fmt.Printf("[ExecuteSingleCase] 跳过变量获取 (case_group_id=%d)\n", task.CaseGroupID)
	}

	// 4C. 替换脚本中的变量
	fmt.Printf("[ExecuteSingleCase] 脚本替换前长度: %d bytes\n", len(c.ScriptCode))
	replacedScript := s.replaceVariables(c.ScriptCode, variables)
	fmt.Printf("[ExecuteSingleCase] 脚本替换后长度: %d bytes\n", len(replacedScript))
	if strings.Contains(replacedScript, "${") {
		fmt.Printf("[ExecuteSingleCase] ⚠️ 脚本仍包含变量占位符！\n")
	}

	// 5. 确定remark语言
	lang := task.DisplayLanguage
	if lang == "" {
		lang = "cn"
	}

	// 6. 执行用例
	fmt.Printf("[ExecuteSingleCase] 开始执行 Playwright 脚本...\n")
	execResult, execErr := s.executeViaPlaywright(replacedScript)
	var okCount, ngCount, blockCount int
	if execErr != nil {
		// 执行失败
		fmt.Printf("[ExecuteSingleCase] 执行失败: %v\n", execErr)
		c.TestResult = "NG"
		c.Remark = s.getRemarkByLang(lang, false, execErr.Error())
		ngCount = 1
	} else {
		fmt.Printf("[ExecuteSingleCase] 执行成功: response_time=%dms\n", execResult.ResponseTime)
		c.TestResult = "OK"
		c.Remark = s.getRemarkByLang(lang, true, "")
		okCount = 1
	}
	// 记录执行时间（api和automation类型都记录）
	if execResult != nil && execResult.ResponseTime > 0 {
		c.ResponseTime = fmt.Sprintf("%d", execResult.ResponseTime)
	}

	// 7. 更新数据库
	updates := map[string]interface{}{
		"test_result": c.TestResult,
		"remark":      c.Remark,
	}
	if c.ResponseTime != "" {
		updates["response_time"] = c.ResponseTime
	}
	if updateErr := s.ecrRepo.UpdateResult(c.ID, updates); updateErr != nil {
		return nil, fmt.Errorf("update case result: %w", updateErr)
	}

	// 8. 更新任务的测试日期和执行人
	now := time.Now()
	executor := s.getUserName(userID)
	taskUpdates := map[string]interface{}{
		"test_date": now,
		"executor":  executor,
	}
	if updateErr := s.repo.UpdateByUUID(taskUUID, taskUpdates); updateErr != nil {
		return nil, fmt.Errorf("update task: %w", updateErr)
	}

	fmt.Printf("[ExecuteSingleCase] 执行完成: OK=%d, NG=%d, Block=%d\n", okCount, ngCount, blockCount)
	return &ExecuteTaskResult{
		Total:      1,
		OKCount:    okCount,
		NGCount:    ngCount,
		BlockCount: blockCount,
		ExecutedAt: now,
		ExecutedBy: executor,
	}, nil
}
