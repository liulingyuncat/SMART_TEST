package services

import (
	"errors"
	"fmt"
	"webtest/internal/models"
	"webtest/internal/repositories"
)

// SaveCaseResultRequest 保存执行结果请求
type SaveCaseResultRequest struct {
	CaseID     string `json:"case_id" binding:"required"`
	DisplayID  uint   `json:"display_id"`
	CaseNum    string `json:"case_num"`  // 用户自定义CaseID
	CaseType   string `json:"case_type"` // 用例类型: overall/acceptance/change/role1-4/api
	TestResult string `json:"test_result" binding:"required,oneof=NR OK NG Block"`
	BugID      string `json:"bug_id" binding:"omitempty,max=50"`
	Remark     string `json:"remark" binding:"omitempty"`

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
}

// ExecutionCaseResultService 测试执行用例结果服务接口
type ExecutionCaseResultService interface {
	GetCaseResults(taskUUID string) ([]*models.ExecutionCaseResult, error)
	SaveCaseResults(taskUUID string, userID uint, requests []SaveCaseResultRequest) error
	GetStatistics(taskUUID string) (map[string]int, error)
	InitTaskResults(taskUUID string, projectID uint, executionType string, userID uint) error
	ClearTaskResults(taskUUID string) error
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
			TaskUUID:   taskUUID,
			CaseID:     req.CaseID,
			DisplayID:  req.DisplayID,
			CaseNum:    req.CaseNum,
			CaseType:   caseType,
			TestResult: req.TestResult,
			BugID:      req.BugID,
			Remark:     req.Remark,
			// AI Web/API 字段
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
			UpdatedBy:        userID,
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
func (s *executionCaseResultService) InitTaskResults(taskUUID string, projectID uint, executionType string, userID uint) error {
	// 验证任务存在
	_, err := s.taskRepo.GetByUUID(taskUUID)
	if err != nil {
		return fmt.Errorf("task %s not found: %w", taskUUID, err)
	}

	// 根据executionType获取用例列表
	var caseIDs []string
	var caseType string

	switch executionType {
	case "manual":
		// 获取整体/受入/变更用例
		caseType = "overall"
		cases, err := s.fetchManualCases(projectID)
		if err != nil {
			return fmt.Errorf("fetch manual cases: %w", err)
		}
		caseIDs = extractCaseIDs(cases)

	case "automation":
		// 获取role1-4用例
		caseType = "role1"
		cases, err := s.fetchAutoCases(projectID)
		if err != nil {
			return fmt.Errorf("fetch auto cases: %w", err)
		}
		caseIDs = extractAutoCaseIDs(cases)

	case "api":
		// 获取API用例
		caseType = "api"
		cases, err := s.fetchApiCases(projectID)
		if err != nil {
			return fmt.Errorf("fetch api cases: %w", err)
		}
		caseIDs = extractApiCaseIDs(cases)

	default:
		return fmt.Errorf("invalid execution_type: %s", executionType)
	}

	// 构建默认执行结果(NR状态)
	results := make([]*models.ExecutionCaseResult, 0, len(caseIDs))
	for _, caseID := range caseIDs {
		result := &models.ExecutionCaseResult{
			TaskUUID:   taskUUID,
			CaseID:     caseID,
			CaseType:   caseType,
			TestResult: "NR",
			UpdatedBy:  userID,
		}
		results = append(results, result)
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

// fetchManualCases 获取手工测试用例
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

// fetchAutoCases 获取自动化测试用例
func (s *executionCaseResultService) fetchAutoCases(projectID uint) ([]*models.AutoTestCase, error) {
	var allCases []*models.AutoTestCase
	caseTypes := []string{"role1", "role2", "role3", "role4"}

	for _, caseType := range caseTypes {
		cases, err := s.autoRepo.GetByProjectAndType(projectID, caseType)
		if err != nil {
			return nil, fmt.Errorf("get %s cases: %w", caseType, err)
		}
		allCases = append(allCases, cases...)
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

// extractCaseIDs 提取手工测试用例ID
func extractCaseIDs(cases []*models.ManualTestCase) []string {
	ids := make([]string, 0, len(cases))
	for _, c := range cases {
		ids = append(ids, c.CaseID)
	}
	return ids
}

// extractAutoCaseIDs 提取自动化测试用例ID
func extractAutoCaseIDs(cases []*models.AutoTestCase) []string {
	ids := make([]string, 0, len(cases))
	for _, c := range cases {
		ids = append(ids, c.CaseID)
	}
	return ids
}

// extractApiCaseIDs 提取API测试用例ID
func extractApiCaseIDs(cases []*models.ApiTestCase) []string {
	ids := make([]string, 0, len(cases))
	for _, c := range cases {
		ids = append(ids, c.ID)
	}
	return ids
}
