package services

import (
	"errors"
	"fmt"
	"log"
	"webtest/internal/models"
	"webtest/internal/repositories"
)

// truncateString 截断字符串用于日志输出
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// SaveCaseResultRequest 保存执行结果请求
type SaveCaseResultRequest struct {
	CaseID        string `json:"case_id" binding:"required"`
	DisplayID     uint   `json:"display_id"`
	CaseNum       string `json:"case_num"`        // 用户自定义CaseID
	CaseType      string `json:"case_type"`       // 用例类型: overall/acceptance/change/role1-4/api
	CaseGroupID   uint   `json:"case_group_id"`   // 用例集ID
	CaseGroupName string `json:"case_group_name"` // 用例集名字
	TestResult    string `json:"test_result" binding:"required,oneof=NR OK NG Block"`
	BugID         string `json:"bug_id" binding:"omitempty,max=50"`
	Remark        string `json:"remark" binding:"omitempty"`

	// 用例内容快照 - AI Web/API用例
	ScreenCN   string `json:"screen_cn"`
	ScreenJP   string `json:"screen_jp"`
	ScreenEN   string `json:"screen_en"`
	FunctionCN string `json:"function_cn"`
	FunctionJP string `json:"function_jp"`
	FunctionEN string `json:"function_en"`

	// 用例内容快照 - 手工测试用例
	MajorFunctionCN  string `json:"major_function_cn"`
	MajorFunctionJP  string `json:"major_function_jp"`
	MajorFunctionEN  string `json:"major_function_en"`
	MiddleFunctionCN string `json:"middle_function_cn"`
	MiddleFunctionJP string `json:"middle_function_jp"`
	MiddleFunctionEN string `json:"middle_function_en"`
	MinorFunctionCN  string `json:"minor_function_cn"`
	MinorFunctionJP  string `json:"minor_function_jp"`
	MinorFunctionEN  string `json:"minor_function_en"`

	// 通用字段
	PreconditionCN   string `json:"precondition_cn"`
	PreconditionJP   string `json:"precondition_jp"`
	PreconditionEN   string `json:"precondition_en"`
	TestStepsCN      string `json:"test_steps_cn"`
	TestStepsJP      string `json:"test_steps_jp"`
	TestStepsEN      string `json:"test_steps_en"`
	ExpectedResultCN string `json:"expected_result_cn"`
	ExpectedResultJP string `json:"expected_result_jp"`
	ExpectedResultEN string `json:"expected_result_en"`

	// API 用例特有字段
	Screen       string `json:"screen"`
	URL          string `json:"url"`
	Header       string `json:"header"`
	Method       string `json:"method"`
	Body         string `json:"body"`
	Response     string `json:"response"`
	ResponseTime string `json:"response_time"`
	ScriptCode   string `json:"script_code"` // JS脚本代码，用于API测试执行
}

// ExecutionCaseResultService 测试执行用例结果服务接口
type ExecutionCaseResultService interface {
	GetCaseResults(taskUUID string) ([]*models.ExecutionCaseResult, error)
	SaveCaseResults(taskUUID string, userID uint, requests []SaveCaseResultRequest) error
	GetStatistics(taskUUID string) (map[string]int, error)
	InitTaskResults(taskUUID string, projectID uint, executionType string, userID uint, caseGroupID uint, caseGroupName string) error
	ClearTaskResults(taskUUID string) error
	UpdateSingleResult(id uint, result string, comment string, userID uint) error
	UpdateSingleResultWithBugID(id uint, result string, comment string, bugID string, userID uint) error
	UpdateSingleResultWithBugIDAndResponseTime(id uint, result string, comment string, bugID string, responseTime string, userID uint) error
	BatchUpdateResults(updates []UpdateResultRequest, userID uint) ([]UpdateResultResponse, error)
	BatchUpdateResultsWithBugID(updates []UpdateResultRequestWithBugID, userID uint) ([]UpdateResultResponse, error)
}

// UpdateResultRequest 单个用例结果更新请求
type UpdateResultRequest struct {
	ID      uint   `json:"id" binding:"required"`
	Result  string `json:"result" binding:"required,oneof=NR OK NG Block"`
	Comment string `json:"comment"`
}

// UpdateResultRequestWithBugID 单个用例结果更新请求（支持bug_id和response_time）
type UpdateResultRequestWithBugID struct {
	ID           uint   `json:"id" binding:"required"`
	Result       string `json:"result" binding:"required,oneof=NR OK NG Block"`
	Comment      string `json:"comment"`
	BugID        string `json:"bug_id"`
	ResponseTime string `json:"response_time"`
}

// UpdateResultResponse 单个用例结果更新响应
type UpdateResultResponse struct {
	ID      uint   `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type executionCaseResultService struct {
	repo       repositories.ExecutionCaseResultRepository
	taskRepo   repositories.ExecutionTaskRepository
	manualRepo repositories.ManualTestCaseRepository
	autoRepo   repositories.AutoTestCaseRepository
	apiRepo    repositories.ApiTestCaseRepository
}

// NewExecutionCaseResultService 创建服务实例
func NewExecutionCaseResultService(
	repo repositories.ExecutionCaseResultRepository,
	taskRepo repositories.ExecutionTaskRepository,
	manualRepo repositories.ManualTestCaseRepository,
	autoRepo repositories.AutoTestCaseRepository,
	apiRepo repositories.ApiTestCaseRepository,
) ExecutionCaseResultService {
	return &executionCaseResultService{
		repo:       repo,
		taskRepo:   taskRepo,
		manualRepo: manualRepo,
		autoRepo:   autoRepo,
		apiRepo:    apiRepo,
	}
}

// GetCaseResults 获取任务的所有执行结果
func (s *executionCaseResultService) GetCaseResults(taskUUID string) ([]*models.ExecutionCaseResult, error) {
	results, err := s.repo.GetByTaskUUID(taskUUID)
	if err != nil {
		return nil, fmt.Errorf("get case results for task %s: %w", taskUUID, err)
	}
	return results, nil
}

// SaveCaseResults 保存或更新执行结果
func (s *executionCaseResultService) SaveCaseResults(taskUUID string, userID uint, requests []SaveCaseResultRequest) error {
	if len(requests) == 0 {
		return errors.New("requests array is empty")
	}

	// 验证任务存在
	task, err := s.taskRepo.GetByUUID(taskUUID)
	if err != nil {
		return fmt.Errorf("task %s not found: %w", taskUUID, err)
	}

	// 构建执行结果对象数组
	results := make([]*models.ExecutionCaseResult, 0, len(requests))
	for _, req := range requests {
		// 校验枚举值
		if req.TestResult != "NR" && req.TestResult != "OK" &&
			req.TestResult != "NG" && req.TestResult != "Block" {
			return fmt.Errorf("invalid test_result: %s", req.TestResult)
		}

		// 确定用例类型：优先使用请求中的 case_type，否则根据任务类型推断
		caseType := req.CaseType
		if caseType == "" {
			caseType = s.inferCaseType(task.ExecutionType)
		}

		result := &models.ExecutionCaseResult{
			TaskUUID:      taskUUID,
			CaseID:        req.CaseID,
			DisplayID:     req.DisplayID,
			CaseNum:       req.CaseNum,
			CaseType:      caseType,
			CaseGroupID:   req.CaseGroupID,
			CaseGroupName: req.CaseGroupName,
			TestResult:    req.TestResult,
			BugID:         req.BugID,
			Remark:        req.Remark,
			// AI Web 字段
			ScreenCN:   req.ScreenCN,
			ScreenJP:   req.ScreenJP,
			ScreenEN:   req.ScreenEN,
			FunctionCN: req.FunctionCN,
			FunctionJP: req.FunctionJP,
			FunctionEN: req.FunctionEN,
			// 手工测试字段
			MajorFunctionCN:  req.MajorFunctionCN,
			MajorFunctionJP:  req.MajorFunctionJP,
			MajorFunctionEN:  req.MajorFunctionEN,
			MiddleFunctionCN: req.MiddleFunctionCN,
			MiddleFunctionJP: req.MiddleFunctionJP,
			MiddleFunctionEN: req.MiddleFunctionEN,
			MinorFunctionCN:  req.MinorFunctionCN,
			MinorFunctionJP:  req.MinorFunctionJP,
			MinorFunctionEN:  req.MinorFunctionEN,
			// 通用字段
			PreconditionCN:   req.PreconditionCN,
			PreconditionJP:   req.PreconditionJP,
			PreconditionEN:   req.PreconditionEN,
			TestStepsCN:      req.TestStepsCN,
			TestStepsJP:      req.TestStepsJP,
			TestStepsEN:      req.TestStepsEN,
			ExpectedResultCN: req.ExpectedResultCN,
			ExpectedResultJP: req.ExpectedResultJP,
			ExpectedResultEN: req.ExpectedResultEN,
			// API 用例特有字段
			Screen:       req.Screen,
			URL:          req.URL,
			Header:       req.Header,
			Method:       req.Method,
			Body:         req.Body,
			Response:     req.Response,
			ResponseTime: req.ResponseTime,
			ScriptCode:   req.ScriptCode,
			UpdatedBy:    userID,
		}
		results = append(results, result)
	}

	// 批量upsert
	err = s.repo.BatchUpsert(results)
	if err != nil {
		return fmt.Errorf("batch upsert %d results: %w", len(results), err)
	}

	return nil
}

// GetStatistics 获取任务的统计信息
func (s *executionCaseResultService) GetStatistics(taskUUID string) (map[string]int, error) {
	stats, err := s.repo.GetStatistics(taskUUID)
	if err != nil {
		return nil, fmt.Errorf("get statistics for task %s: %w", taskUUID, err)
	}
	return stats, nil
}

// InitTaskResults 初始化任务执行结果
// caseGroupID: 可选，指定用例集ID时仅导入该用例集的用例；为0时导入所有用例
// caseGroupName: 可选，指定用例集名称（如果caseGroupID>0时会查询对应名称）
func (s *executionCaseResultService) InitTaskResults(taskUUID string, projectID uint, executionType string, userID uint, caseGroupID uint, caseGroupName string) error {
	// 验证任务存在
	_, err := s.taskRepo.GetByUUID(taskUUID)
	if err != nil {
		return fmt.Errorf("task %s not found: %w", taskUUID, err)
	}

	// 更新任务的用例集信息（使用UpdateByUUID方法）
	if caseGroupID > 0 || caseGroupName != "" {
		updates := map[string]interface{}{
			"case_group_id":   caseGroupID,
			"case_group_name": caseGroupName,
		}
		if err := s.taskRepo.UpdateByUUID(taskUUID, updates); err != nil {
			return fmt.Errorf("update task case_group info: %w", err)
		}
	}

	// 根据executionType获取用例列表并包含用例集信息
	var results []*models.ExecutionCaseResult

	switch executionType {
	case "manual":
		// 获取手工用例（包含用例集信息）
		cases, err := s.fetchManualCasesByGroup(projectID, caseGroupName)
		if err != nil {
			return fmt.Errorf("fetch manual cases: %w", err)
		}
		results = s.buildExecutionResultsForManual(taskUUID, cases, userID)

	case "automation":
		// 获取自动化用例（包含用例集信息）
		cases, err := s.fetchAutoCasesByGroup(projectID, caseGroupName)
		if err != nil {
			return fmt.Errorf("fetch auto cases: %w", err)
		}
		results = s.buildExecutionResultsForAuto(taskUUID, cases, userID)

	case "api":
		// 获取API用例（包含用例集信息）
		cases, err := s.fetchApiCasesByGroup(projectID, caseGroupName)
		if err != nil {
			return fmt.Errorf("fetch api cases: %w", err)
		}
		results = s.buildExecutionResultsForApi(taskUUID, cases, userID)

	default:
		return fmt.Errorf("invalid execution_type: %s", executionType)
	}

	// 批量创建
	err = s.repo.BatchCreate(results)
	if err != nil {
		return fmt.Errorf("batch create %d default results: %w", len(results), err)
	}

	return nil
}

// ClearTaskResults 清空任务执行结果
func (s *executionCaseResultService) ClearTaskResults(taskUUID string) error {
	err := s.repo.DeleteByTaskUUID(taskUUID)
	if err != nil {
		return fmt.Errorf("clear results for task %s: %w", taskUUID, err)
	}
	return nil
}

// UpdateSingleResult 更新单个用例的执行结果
func (s *executionCaseResultService) UpdateSingleResult(id uint, result string, comment string, userID uint) error {
	updates := map[string]interface{}{
		"test_result": result,
		"updated_by":  userID,
	}
	if comment != "" {
		updates["remark"] = comment
	}

	err := s.repo.UpdateResult(id, updates)
	if err != nil {
		return fmt.Errorf("update result %d: %w", id, err)
	}
	return nil
}

// BatchUpdateResults 批量更新用例的执行结果
func (s *executionCaseResultService) BatchUpdateResults(updates []UpdateResultRequest, userID uint) ([]UpdateResultResponse, error) {
	responses := make([]UpdateResultResponse, len(updates))

	for i, req := range updates {
		err := s.UpdateSingleResult(req.ID, req.Result, req.Comment, userID)
		if err != nil {
			responses[i] = UpdateResultResponse{
				ID:      req.ID,
				Success: false,
				Message: err.Error(),
			}
		} else {
			responses[i] = UpdateResultResponse{
				ID:      req.ID,
				Success: true,
				Message: "updated",
			}
		}
	}

	return responses, nil
}

// UpdateSingleResultWithBugID 更新单个用例的执行结果（支持bug_id）
func (s *executionCaseResultService) UpdateSingleResultWithBugID(id uint, result string, comment string, bugID string, userID uint) error {
	return s.UpdateSingleResultWithBugIDAndResponseTime(id, result, comment, bugID, "", userID)
}

// UpdateSingleResultWithBugIDAndResponseTime 更新单个用例的执行结果（支持bug_id和response_time）
func (s *executionCaseResultService) UpdateSingleResultWithBugIDAndResponseTime(id uint, result string, comment string, bugID string, responseTime string, userID uint) error {
	updates := map[string]interface{}{
		"test_result": result,
		"updated_by":  userID,
	}
	if comment != "" {
		updates["remark"] = comment
	}
	if bugID != "" {
		updates["bug_id"] = bugID
	}
	if responseTime != "" {
		updates["response_time"] = responseTime
	}

	err := s.repo.UpdateResult(id, updates)
	if err != nil {
		return fmt.Errorf("update result %d: %w", id, err)
	}
	return nil
}

// BatchUpdateResultsWithBugID 批量更新用例的执行结果（支持bug_id和response_time）
func (s *executionCaseResultService) BatchUpdateResultsWithBugID(updates []UpdateResultRequestWithBugID, userID uint) ([]UpdateResultResponse, error) {
	responses := make([]UpdateResultResponse, len(updates))

	for i, req := range updates {
		err := s.UpdateSingleResultWithBugIDAndResponseTime(req.ID, req.Result, req.Comment, req.BugID, req.ResponseTime, userID)
		if err != nil {
			responses[i] = UpdateResultResponse{
				ID:      req.ID,
				Success: false,
				Message: err.Error(),
			}
		} else {
			responses[i] = UpdateResultResponse{
				ID:      req.ID,
				Success: true,
				Message: "updated",
			}
		}
	}

	return responses, nil
}

// mapResultStatus 将 passed/failed/blocked/skipped 映射为 OK/NG/Block/NR
func mapResultStatus(result string) string {
	switch result {
	case "passed":
		return "OK"
	case "failed":
		return "NG"
	case "blocked":
		return "Block"
	case "skipped":
		return "NR"
	default:
		return result
	}
}

// ========== 辅助方法 ==========

// inferCaseType 根据executionType推断caseType
func (s *executionCaseResultService) inferCaseType(executionType string) string {
	switch executionType {
	case "manual":
		return "overall"
	case "automation":
		return "role1"
	case "api":
		return "api"
	default:
		return "overall"
	}
}

// buildExecutionResultsForManual 为手工用例构建执行结果，包含用例集信息
func (s *executionCaseResultService) buildExecutionResultsForManual(taskUUID string, cases []*models.ManualTestCase, userID uint) []*models.ExecutionCaseResult {
	results := make([]*models.ExecutionCaseResult, 0, len(cases))
	for _, testCase := range cases {
		result := &models.ExecutionCaseResult{
			TaskUUID:      taskUUID,
			CaseID:        testCase.CaseID,
			CaseGroupID:   0, // 从手工用例中可能无法获取ID，只能获取名字
			CaseGroupName: testCase.CaseGroup,
			CaseType:      testCase.CaseType,
			TestResult:    "NR",
			UpdatedBy:     userID,
		}
		results = append(results, result)
	}
	return results
}

// buildExecutionResultsForAuto 为自动化用例构建执行结果，包含用例集信息和完整用例字段
func (s *executionCaseResultService) buildExecutionResultsForAuto(taskUUID string, cases []*models.AutoTestCase, userID uint) []*models.ExecutionCaseResult {
	results := make([]*models.ExecutionCaseResult, 0, len(cases))
	for i, testCase := range cases {
		result := &models.ExecutionCaseResult{
			TaskUUID:      taskUUID,
			CaseID:        testCase.CaseID,
			DisplayID:     uint(i + 1), // 设置显示序号
			CaseNum:       testCase.CaseNumber,
			CaseGroupID:   0, // 从自动化用例中可能无法获取ID，只能获取名字
			CaseGroupName: testCase.CaseGroup,
			CaseType:      testCase.CaseType,
			TestResult:    "Block", // 默认为Block状态
			Remark:        testCase.Remark,
			// 复制用例内容字段 - 中文
			ScreenCN:         testCase.ScreenCN,
			FunctionCN:       testCase.FunctionCN,
			PreconditionCN:   testCase.PreconditionCN,
			TestStepsCN:      testCase.TestStepsCN,
			ExpectedResultCN: testCase.ExpectedResultCN,
			// 复制用例内容字段 - 日文
			ScreenJP:         testCase.ScreenJP,
			FunctionJP:       testCase.FunctionJP,
			PreconditionJP:   testCase.PreconditionJP,
			TestStepsJP:      testCase.TestStepsJP,
			ExpectedResultJP: testCase.ExpectedResultJP,
			// 复制用例内容字段 - 英文
			ScreenEN:         testCase.ScreenEN,
			FunctionEN:       testCase.FunctionEN,
			PreconditionEN:   testCase.PreconditionEN,
			TestStepsEN:      testCase.TestStepsEN,
			ExpectedResultEN: testCase.ExpectedResultEN,
			// 脚本代码
			ScriptCode: testCase.ScriptCode,
			UpdatedBy:  userID,
		}
		results = append(results, result)
	}
	return results
}

// buildExecutionResultsForApi 为API用例构建执行结果，包含用例集信息和API字段快照
func (s *executionCaseResultService) buildExecutionResultsForApi(taskUUID string, cases []*models.ApiTestCase, userID uint) []*models.ExecutionCaseResult {
	log.Printf("[buildExecutionResultsForApi] Building results for %d cases, taskUUID=%s", len(cases), taskUUID)
	results := make([]*models.ExecutionCaseResult, 0, len(cases))
	for i, testCase := range cases {
		log.Printf("[buildExecutionResultsForApi] Case[%d] ID=%s, ScriptCode length=%d, first50chars=%q",
			i, testCase.ID, len(testCase.ScriptCode), truncateString(testCase.ScriptCode, 50))
		result := &models.ExecutionCaseResult{
			TaskUUID:      taskUUID,
			CaseID:        testCase.ID,
			CaseGroupID:   0, // 从API用例中可能无法获取ID，只能获取名字
			CaseGroupName: testCase.CaseGroup,
			CaseType:      testCase.CaseType,
			TestResult:    "NR",
			UpdatedBy:     userID,
			// API用例字段快照
			Screen:     testCase.Screen,
			URL:        testCase.URL,
			Header:     testCase.Header,
			Method:     testCase.Method,
			Body:       testCase.Body,
			Response:   testCase.Response,
			ScriptCode: testCase.ScriptCode,
		}
		results = append(results, result)
	}
	return results
}

// fetchManualCases 获取手工测试用例（获取所有用例）
func (s *executionCaseResultService) fetchManualCases(projectID uint) ([]*models.ManualTestCase, error) {
	var allCases []*models.ManualTestCase
	caseTypes := []string{"overall", "acceptance", "change"}

	for _, caseType := range caseTypes {
		cases, err := s.manualRepo.GetByProjectAndType(projectID, caseType)
		if err != nil {
			return nil, fmt.Errorf("get %s cases: %w", caseType, err)
		}
		allCases = append(allCases, cases...)
	}

	return allCases, nil
}

// fetchManualCasesByGroup 按用例集名称获取手工测试用例
func (s *executionCaseResultService) fetchManualCasesByGroup(projectID uint, caseGroupName string) ([]*models.ManualTestCase, error) {
	// 如果没有指定用例集，返回所有用例
	if caseGroupName == "" {
		return s.fetchManualCases(projectID)
	}

	// 按用例集名称过滤
	var allCases []*models.ManualTestCase
	caseTypes := []string{"overall", "acceptance", "change"}

	for _, caseType := range caseTypes {
		cases, err := s.manualRepo.GetByProjectAndType(projectID, caseType)
		if err != nil {
			return nil, fmt.Errorf("get %s cases: %w", caseType, err)
		}
		// 过滤出指定用例集的用例
		for _, c := range cases {
			if c.CaseGroup == caseGroupName {
				allCases = append(allCases, c)
			}
		}
	}

	return allCases, nil
}

// fetchAutoCases 获取自动化测试用例（包括role1-4和web类型）
func (s *executionCaseResultService) fetchAutoCases(projectID uint) ([]*models.AutoTestCase, error) {
	var allCases []*models.AutoTestCase
	// 添加 web 类型，支持 AI Web 用例集
	caseTypes := []string{"role1", "role2", "role3", "role4", "web"}

	for _, caseType := range caseTypes {
		cases, err := s.autoRepo.GetByProjectAndType(projectID, caseType)
		if err != nil {
			return nil, fmt.Errorf("get %s cases: %w", caseType, err)
		}
		allCases = append(allCases, cases...)
	}

	return allCases, nil
}

// fetchAutoCasesByGroup 按用例集名称获取自动化测试用例
func (s *executionCaseResultService) fetchAutoCasesByGroup(projectID uint, caseGroupName string) ([]*models.AutoTestCase, error) {
	// 如果没有指定用例集，返回所有用例
	if caseGroupName == "" {
		return s.fetchAutoCases(projectID)
	}

	// 按用例集名称过滤（主要针对 web 类型用例）
	var allCases []*models.AutoTestCase
	caseTypes := []string{"role1", "role2", "role3", "role4", "web"}

	for _, caseType := range caseTypes {
		cases, err := s.autoRepo.GetByProjectAndType(projectID, caseType)
		if err != nil {
			return nil, fmt.Errorf("get %s cases: %w", caseType, err)
		}
		// 过滤出指定用例集的用例
		for _, c := range cases {
			if c.CaseGroup == caseGroupName {
				allCases = append(allCases, c)
			}
		}
	}

	return allCases, nil
}

// fetchApiCases 获取API测试用例
func (s *executionCaseResultService) fetchApiCases(projectID uint) ([]*models.ApiTestCase, error) {
	cases, err := s.apiRepo.GetByProjectAndType(projectID, "api")
	if err != nil {
		return nil, fmt.Errorf("get api cases: %w", err)
	}
	return cases, nil
}

// fetchApiCasesByGroup 按用例集名称获取API测试用例
func (s *executionCaseResultService) fetchApiCasesByGroup(projectID uint, caseGroupName string) ([]*models.ApiTestCase, error) {
	// 如果没有指定用例集，返回所有用例
	if caseGroupName == "" {
		return s.fetchApiCases(projectID)
	}

	// 按用例集名称过滤
	cases, err := s.apiRepo.GetByProjectAndType(projectID, "api")
	if err != nil {
		return nil, fmt.Errorf("get api cases: %w", err)
	}

	// DEBUG: 打印获取到的用例的ScriptCode
	log.Printf("[fetchApiCasesByGroup] Got %d cases from repo for projectID=%d", len(cases), projectID)
	for i, c := range cases {
		log.Printf("[fetchApiCasesByGroup] Case[%d] ID=%s, CaseGroup=%s, ScriptCode length=%d", i, c.ID, c.CaseGroup, len(c.ScriptCode))
	}

	// 过滤出指定用例集的用例
	var filteredCases []*models.ApiTestCase
	for _, c := range cases {
		if c.CaseGroup == caseGroupName {
			filteredCases = append(filteredCases, c)
		}
	}

	log.Printf("[fetchApiCasesByGroup] Filtered to %d cases for caseGroupName=%s", len(filteredCases), caseGroupName)
	for i, c := range filteredCases {
		log.Printf("[fetchApiCasesByGroup] Filtered[%d] ID=%s, ScriptCode length=%d", i, c.ID, len(c.ScriptCode))
	}

	return filteredCases, nil
}
