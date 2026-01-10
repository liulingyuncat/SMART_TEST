package services

import (
	"errors"
	"fmt"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MetadataDTO 元数据响应DTO
type MetadataDTO struct {
	TestVersion string `json:"test_version"`
	TestEnv     string `json:"test_env"`
	TestDate    string `json:"test_date"`
	Executor    string `json:"executor"`
}

// UpdateMetadataRequest 更新元数据请求
type UpdateMetadataRequest struct {
	TestVersion string `json:"test_version" binding:"max=50"`
	TestEnv     string `json:"test_env" binding:"max=100"`
	TestDate    string `json:"test_date" binding:"max=20"`
	Executor    string `json:"executor" binding:"max=50"`
}

// CaseDTO 用例DTO
type CaseDTO struct {
	CaseID     string `json:"case_id"`    // UUID主键，用于更新和删除操作
	ID         uint   `json:"id"`         // 显示序号
	DisplayID  uint   `json:"display_id"` // 显示用ID,同一用例的多语言版本共享此ID
	CaseNumber string `json:"case_number"`

	// ======== 单语言字段(AI用例使用) ========
	MajorFunction  string `json:"major_function,omitempty"`  // ai用例
	MiddleFunction string `json:"middle_function,omitempty"` // ai用例
	MinorFunction  string `json:"minor_function,omitempty"`  // ai用例
	Precondition   string `json:"precondition,omitempty"`    // ai用例
	TestSteps      string `json:"test_steps,omitempty"`      // ai用例
	ExpectedResult string `json:"expected_result,omitempty"` // ai用例

	// ======== 多语言字段(整体/变更用例使用) ========
	MajorFunctionCN  string `json:"major_function_cn,omitempty"`  // overall/change用例
	MajorFunctionJP  string `json:"major_function_jp,omitempty"`  // overall/change用例
	MajorFunctionEN  string `json:"major_function_en,omitempty"`  // overall/change用例
	MiddleFunctionCN string `json:"middle_function_cn,omitempty"` // overall/change用例
	MiddleFunctionJP string `json:"middle_function_jp,omitempty"` // overall/change用例
	MiddleFunctionEN string `json:"middle_function_en,omitempty"` // overall/change用例
	MinorFunctionCN  string `json:"minor_function_cn,omitempty"`  // overall/change用例
	MinorFunctionJP  string `json:"minor_function_jp,omitempty"`  // overall/change用例
	MinorFunctionEN  string `json:"minor_function_en,omitempty"`  // overall/change用例
	PreconditionCN   string `json:"precondition_cn,omitempty"`    // overall/change用例
	PreconditionJP   string `json:"precondition_jp,omitempty"`    // overall/change用例
	PreconditionEN   string `json:"precondition_en,omitempty"`    // overall/change用例
	TestStepsCN      string `json:"test_steps_cn,omitempty"`      // overall/change用例
	TestStepsJP      string `json:"test_steps_jp,omitempty"`      // overall/change用例
	TestStepsEN      string `json:"test_steps_en,omitempty"`      // overall/change用例
	ExpectedResultCN string `json:"expected_result_cn,omitempty"` // overall/change用例
	ExpectedResultJP string `json:"expected_result_jp,omitempty"` // overall/change用例
	ExpectedResultEN string `json:"expected_result_en,omitempty"` // overall/change用例

	TestResult string `json:"test_result,omitempty"` // 仅overall/change使用
	Remark     string `json:"remark"`
}

// CaseListDTO 用例列表DTO
type CaseListDTO struct {
	Cases    []*CaseDTO `json:"cases"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	Size     int        `json:"size"`
	Language string     `json:"language"`
}

// CreateCaseRequest 创建用例请求
type CreateCaseRequest struct {
	CaseType   string `json:"case_type" binding:"required,oneof=ai overall change acceptance"`
	Language   string `json:"language,omitempty"` // 保留用于API兼容，不再存储到数据库
	CaseNumber string `json:"case_number" binding:"max=50"`
	CaseGroup  string `json:"case_group,omitempty" binding:"max=100"` // 用例集名称

	// ======== 单语言字段(AI用例使用) ========
	MajorFunction  string `json:"major_function,omitempty" binding:"max=100"`
	MiddleFunction string `json:"middle_function,omitempty" binding:"max=100"`
	MinorFunction  string `json:"minor_function,omitempty" binding:"max=100"`
	Precondition   string `json:"precondition,omitempty"`
	TestSteps      string `json:"test_steps,omitempty"`
	ExpectedResult string `json:"expected_result,omitempty"`

	// ======== 多语言字段(整体/变更用例使用) ========
	MajorFunctionCN  string `json:"major_function_cn,omitempty" binding:"max=100"`
	MajorFunctionJP  string `json:"major_function_jp,omitempty" binding:"max=100"`
	MajorFunctionEN  string `json:"major_function_en,omitempty" binding:"max=100"`
	MiddleFunctionCN string `json:"middle_function_cn,omitempty" binding:"max=100"`
	MiddleFunctionJP string `json:"middle_function_jp,omitempty" binding:"max=100"`
	MiddleFunctionEN string `json:"middle_function_en,omitempty" binding:"max=100"`
	MinorFunctionCN  string `json:"minor_function_cn,omitempty" binding:"max=100"`
	MinorFunctionJP  string `json:"minor_function_jp,omitempty" binding:"max=100"`
	MinorFunctionEN  string `json:"minor_function_en,omitempty" binding:"max=100"`
	PreconditionCN   string `json:"precondition_cn,omitempty"`
	PreconditionJP   string `json:"precondition_jp,omitempty"`
	PreconditionEN   string `json:"precondition_en,omitempty"`
	TestStepsCN      string `json:"test_steps_cn,omitempty"`
	TestStepsJP      string `json:"test_steps_jp,omitempty"`
	TestStepsEN      string `json:"test_steps_en,omitempty"`
	ExpectedResultCN string `json:"expected_result_cn,omitempty"`
	ExpectedResultJP string `json:"expected_result_jp,omitempty"`
	ExpectedResultEN string `json:"expected_result_en,omitempty"`

	TestResult string `json:"test_result,omitempty" binding:"omitempty,oneof=OK NG Block NR"` // 仅overall/change使用
	Remark     string `json:"remark,omitempty"`
}

// UpdateCaseRequest 更新用例请求
type UpdateCaseRequest struct {
	CaseNumber *string `json:"case_number,omitempty" binding:"omitempty,max=50"`

	// ======== 单语言字段(AI用例使用) ========
	MajorFunction  *string `json:"major_function,omitempty" binding:"omitempty,max=100"`
	MiddleFunction *string `json:"middle_function,omitempty"`
	MinorFunction  *string `json:"minor_function,omitempty"`
	Precondition   *string `json:"precondition,omitempty"`
	TestSteps      *string `json:"test_steps,omitempty"`
	ExpectedResult *string `json:"expected_result,omitempty"`

	// ======== 多语言字段(整体/变更用例使用) ========
	MajorFunctionCN  *string `json:"major_function_cn,omitempty" binding:"omitempty,max=100"`
	MajorFunctionJP  *string `json:"major_function_jp,omitempty" binding:"omitempty,max=100"`
	MajorFunctionEN  *string `json:"major_function_en,omitempty" binding:"omitempty,max=100"`
	MiddleFunctionCN *string `json:"middle_function_cn,omitempty"`
	MiddleFunctionJP *string `json:"middle_function_jp,omitempty"`
	MiddleFunctionEN *string `json:"middle_function_en,omitempty"`
	MinorFunctionCN  *string `json:"minor_function_cn,omitempty"`
	MinorFunctionJP  *string `json:"minor_function_jp,omitempty"`
	MinorFunctionEN  *string `json:"minor_function_en,omitempty"`
	PreconditionCN   *string `json:"precondition_cn,omitempty"`
	PreconditionJP   *string `json:"precondition_jp,omitempty"`
	PreconditionEN   *string `json:"precondition_en,omitempty"`
	TestStepsCN      *string `json:"test_steps_cn,omitempty"`
	TestStepsJP      *string `json:"test_steps_jp,omitempty"`
	TestStepsEN      *string `json:"test_steps_en,omitempty"`
	ExpectedResultCN *string `json:"expected_result_cn,omitempty"`
	ExpectedResultJP *string `json:"expected_result_jp,omitempty"`
	ExpectedResultEN *string `json:"expected_result_en,omitempty"`

	TestResult *string `json:"test_result,omitempty" binding:"omitempty,oneof=OK NG Block NR"`
	Remark     *string `json:"remark,omitempty"`
}

// ManualTestCaseService 手工测试用例服务接口
type ManualTestCaseService interface {
	GetMetadata(projectID uint, userID uint, caseType string) (*MetadataDTO, error)
	UpdateMetadata(projectID uint, userID uint, caseType string, req UpdateMetadataRequest) error
	GetCases(projectID uint, userID uint, caseType string, language string, page int, size int, caseGroup string) (*CaseListDTO, error)

	// CRUD方法 - caseID使用UUID字符串
	CreateCase(projectID uint, userID uint, req CreateCaseRequest) (*CaseDTO, error)
	UpdateCase(projectID uint, userID uint, caseID string, req UpdateCaseRequest) error             // 改用UUID
	DeleteCase(projectID uint, userID uint, caseID string) error                                    // 改用UUID
	ReorderCases(projectID uint, userID uint, caseType string, caseIDs []uint) ([]uint, error)      // ID显示序号重排（重新排序按钮）
	ReorderCasesByDrag(projectID uint, userID uint, caseType string, caseIDOrder []string) error    // 拖拽重排：根据case_id顺序重新分配ID
	ReorderAllCasesByID(projectID uint, userID uint, caseType string, language string) (int, error) // 按现有ID顺序重新编号所有用例
	ClearAICases(projectID uint, userID uint) (int, error)

	// 新增：插入和批量删除方法
	InsertCase(projectID uint, userID uint, caseType string, position string, targetCaseID string, language string) (*models.ManualTestCase, error)
	BatchDeleteCases(projectID uint, userID uint, caseType string, caseIDs []string) (deletedCount int, failedCaseIDs []string, err error)

	// 新增：重新分配所有ID
	ReassignAllIDs(projectID uint, userID uint, caseType string) error
}

type manualTestCaseService struct {
	repo           repositories.ManualTestCaseRepository
	projectService ProjectService
}

// NewManualTestCaseService 创建服务实例
func NewManualTestCaseService(repo repositories.ManualTestCaseRepository, projectService ProjectService) ManualTestCaseService {
	return &manualTestCaseService{
		repo:           repo,
		projectService: projectService,
	}
}

// GetMetadata 获取项目元数据
func (s *manualTestCaseService) GetMetadata(projectID uint, userID uint, caseType string) (*MetadataDTO, error) {
	// 验证用户是否是项目成员
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 默认类型为 overall
	if caseType == "" {
		caseType = "overall"
	}

	// 获取元数据
	testCase, err := s.repo.GetMetadataByProjectID(projectID, caseType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未找到记录,返回空值
			return &MetadataDTO{}, nil
		}
		return nil, err
	}

	return &MetadataDTO{
		TestVersion: testCase.TestVersion,
		TestEnv:     testCase.TestEnv,
		TestDate:    testCase.TestDate,
		Executor:    testCase.Executor,
	}, nil
}

// UpdateMetadata 更新元数据
func (s *manualTestCaseService) UpdateMetadata(projectID uint, userID uint, caseType string, req UpdateMetadataRequest) error {
	// 验证用户权限
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 默认类型为 overall
	if caseType == "" {
		caseType = "overall"
	}

	// 尝试更新元数据
	metadata := map[string]interface{}{
		"test_version": req.TestVersion,
		"test_env":     req.TestEnv,
		"test_date":    req.TestDate,
		"executor":     req.Executor,
	}

	err = s.repo.UpdateMetadata(projectID, caseType, metadata)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在,创建默认记录后重试
			if createErr := s.repo.CreateDefaultMetadata(projectID, caseType); createErr != nil {
				return fmt.Errorf("create default metadata: %w", createErr)
			}
			// 重新更新
			if retryErr := s.repo.UpdateMetadata(projectID, caseType, metadata); retryErr != nil {
				return fmt.Errorf("retry update metadata: %w", retryErr)
			}
			return nil
		}
		return err
	}

	return nil
}

// GetCases 获取用例列表
func (s *manualTestCaseService) GetCases(projectID uint, userID uint, caseType string, language string, page int, size int, caseGroup string) (*CaseListDTO, error) {
	// 验证用户权限
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 参数校验
	if caseType == "" {
		caseType = "overall"
	}
	if language == "" {
		language = "中文"
	}
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 50
	}
	// 移除最大值限制，允许获取全部数据
	// if size > 100 {
	// 	size = 50
	// }

	offset := (page - 1) * size
	cases, total, err := s.repo.GetCasesByType(projectID, caseType, offset, size, caseGroup)
	if err != nil {
		return nil, err
	}

	// 转换为DTO并计算display_id
	caseDTOs := make([]*CaseDTO, 0, len(cases))

	// 对于整体/变更/受入用例，返回完整的多语言字段（前端根据language筛选显示）
	if caseType == "overall" || caseType == "change" || caseType == "acceptance" {
		for _, c := range cases {
			caseDTOs = append(caseDTOs, &CaseDTO{
				CaseID:     c.CaseID, // UUID主键
				ID:         c.ID,
				DisplayID:  c.ID, // overall/change/acceptance用例直接使用ID作为display_id（v2.1设计：1条记录存储3语言）
				CaseNumber: c.CaseNumber,

				// 返回所有多语言字段，让前端根据当前语言选择显示哪些列
				MajorFunctionCN:  c.MajorFunctionCN,
				MajorFunctionJP:  c.MajorFunctionJP,
				MajorFunctionEN:  c.MajorFunctionEN,
				MiddleFunctionCN: c.MiddleFunctionCN,
				MiddleFunctionJP: c.MiddleFunctionJP,
				MiddleFunctionEN: c.MiddleFunctionEN,
				MinorFunctionCN:  c.MinorFunctionCN,
				MinorFunctionJP:  c.MinorFunctionJP,
				MinorFunctionEN:  c.MinorFunctionEN,
				PreconditionCN:   c.PreconditionCN,
				PreconditionJP:   c.PreconditionJP,
				PreconditionEN:   c.PreconditionEN,
				TestStepsCN:      c.TestStepsCN,
				TestStepsJP:      c.TestStepsJP,
				TestStepsEN:      c.TestStepsEN,
				ExpectedResultCN: c.ExpectedResultCN,
				ExpectedResultJP: c.ExpectedResultJP,
				ExpectedResultEN: c.ExpectedResultEN,

				TestResult: c.TestResult,
				Remark:     c.Remark,
			})
		}
	} else {
		// AI用例使用单语言字段
		for i, c := range cases {
			caseDTOs = append(caseDTOs, &CaseDTO{
				CaseID:         c.CaseID, // UUID主键
				ID:             c.ID,
				DisplayID:      uint(offset + i + 1), // 基于分页的连续序号
				CaseNumber:     c.CaseNumber,
				MajorFunction:  c.MajorFunction,
				MiddleFunction: c.MiddleFunction,
				MinorFunction:  c.MinorFunction,
				Precondition:   c.Precondition,
				TestSteps:      c.TestSteps,
				ExpectedResult: c.ExpectedResult,
				TestResult:     c.TestResult,
				Remark:         c.Remark,
			})
		}
	}

	return &CaseListDTO{
		Cases:    caseDTOs,
		Total:    total,
		Page:     page,
		Size:     size,
		Language: language,
	}, nil
}

// CreateCase 创建新用例,根据用例类型实现多语言联动创建逻辑(AC-07)
func (s *manualTestCaseService) CreateCase(projectID uint, userID uint, req CreateCaseRequest) (*CaseDTO, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// AI用例强制设置为中文
	if req.CaseType == "ai" {
		req.Language = "中文"
	}

	// 设置默认测试结果(仅overall/change/acceptance用例)
	if req.TestResult == "" && (req.CaseType == "overall" || req.CaseType == "change" || req.CaseType == "acceptance") {
		req.TestResult = "NR"
	}

	// 获取当前类型用例的最大ID，新用例ID = 最大ID + 1
	maxID, err := s.repo.GetMaxIDByProjectAndType(projectID, req.CaseType)
	if err != nil {
		return nil, fmt.Errorf("get max id: %w", err)
	}

	testCase := &models.ManualTestCase{
		ID:         maxID + 1, // 设置新用例的ID
		ProjectID:  projectID,
		CaseType:   req.CaseType,
		CaseNumber: req.CaseNumber,
		CaseGroup:  req.CaseGroup,
		TestResult: req.TestResult,
		Remark:     req.Remark,
	}

	// 根据case_type填充不同的字段集
	if req.CaseType == "ai" {
		// AI用例:使用单语言字段
		testCase.MajorFunction = req.MajorFunction
		testCase.MiddleFunction = req.MiddleFunction
		testCase.MinorFunction = req.MinorFunction
		testCase.Precondition = req.Precondition
		testCase.TestSteps = req.TestSteps
		testCase.ExpectedResult = req.ExpectedResult
	} else {
		// 整体/变更/受入用例:使用多语言字段
		testCase.MajorFunctionCN = req.MajorFunctionCN
		testCase.MajorFunctionJP = req.MajorFunctionJP
		testCase.MajorFunctionEN = req.MajorFunctionEN
		testCase.MiddleFunctionCN = req.MiddleFunctionCN
		testCase.MiddleFunctionJP = req.MiddleFunctionJP
		testCase.MiddleFunctionEN = req.MiddleFunctionEN
		testCase.MinorFunctionCN = req.MinorFunctionCN
		testCase.MinorFunctionJP = req.MinorFunctionJP
		testCase.MinorFunctionEN = req.MinorFunctionEN
		testCase.PreconditionCN = req.PreconditionCN
		testCase.PreconditionJP = req.PreconditionJP
		testCase.PreconditionEN = req.PreconditionEN
		testCase.TestStepsCN = req.TestStepsCN
		testCase.TestStepsJP = req.TestStepsJP
		testCase.TestStepsEN = req.TestStepsEN
		testCase.ExpectedResultCN = req.ExpectedResultCN
		testCase.ExpectedResultJP = req.ExpectedResultJP
		testCase.ExpectedResultEN = req.ExpectedResultEN
	}

	if err := s.repo.Create(testCase); err != nil {
		return nil, fmt.Errorf("create test case: %w", err)
	}

	// 构造DTO返回
	dto := &CaseDTO{
		CaseID:     testCase.CaseID, // UUID主键
		ID:         testCase.ID,
		DisplayID:  testCase.ID,
		CaseNumber: testCase.CaseNumber,
		TestResult: testCase.TestResult,
		Remark:     testCase.Remark,
	}

	if req.CaseType == "ai" {
		// AI用例:返回单语言字段
		dto.MajorFunction = testCase.MajorFunction
		dto.MiddleFunction = testCase.MiddleFunction
		dto.MinorFunction = testCase.MinorFunction
		dto.Precondition = testCase.Precondition
		dto.TestSteps = testCase.TestSteps
		dto.ExpectedResult = testCase.ExpectedResult
	} else {
		// 整体/变更用例:返回多语言字段
		dto.MajorFunctionCN = testCase.MajorFunctionCN
		dto.MajorFunctionJP = testCase.MajorFunctionJP
		dto.MajorFunctionEN = testCase.MajorFunctionEN
		dto.MiddleFunctionCN = testCase.MiddleFunctionCN
		dto.MiddleFunctionJP = testCase.MiddleFunctionJP
		dto.MiddleFunctionEN = testCase.MiddleFunctionEN
		dto.MinorFunctionCN = testCase.MinorFunctionCN
		dto.MinorFunctionJP = testCase.MinorFunctionJP
		dto.MinorFunctionEN = testCase.MinorFunctionEN
		dto.PreconditionCN = testCase.PreconditionCN
		dto.PreconditionJP = testCase.PreconditionJP
		dto.PreconditionEN = testCase.PreconditionEN
		dto.TestStepsCN = testCase.TestStepsCN
		dto.TestStepsJP = testCase.TestStepsJP
		dto.TestStepsEN = testCase.TestStepsEN
		dto.ExpectedResultCN = testCase.ExpectedResultCN
		dto.ExpectedResultJP = testCase.ExpectedResultJP
		dto.ExpectedResultEN = testCase.ExpectedResultEN
	}

	return dto, nil
}

// createWithTranslations overall/change用例使用多语言字段创建单条记录
func (s *manualTestCaseService) createWithTranslations(projectID uint, req CreateCaseRequest) ([]*models.ManualTestCase, error) {
	// 根据请求语言填充对应的多语言字段
	testCase := &models.ManualTestCase{
		ProjectID:  projectID,
		CaseType:   req.CaseType,
		CaseNumber: req.CaseNumber,
		CaseGroup:  req.CaseGroup,
		TestResult: req.TestResult,
		Remark:     req.Remark,
	}

	// 根据请求语言填充对应的CN/JP/EN字段
	switch req.Language {
	case "中文":
		testCase.MajorFunctionCN = req.MajorFunction
		testCase.MiddleFunctionCN = req.MiddleFunction
		testCase.MinorFunctionCN = req.MinorFunction
		testCase.PreconditionCN = req.Precondition
		testCase.TestStepsCN = req.TestSteps
		testCase.ExpectedResultCN = req.ExpectedResult
	case "English":
		testCase.MajorFunctionEN = req.MajorFunction
		testCase.MiddleFunctionEN = req.MiddleFunction
		testCase.MinorFunctionEN = req.MinorFunction
		testCase.PreconditionEN = req.Precondition
		testCase.TestStepsEN = req.TestSteps
		testCase.ExpectedResultEN = req.ExpectedResult
	case "日本語":
		testCase.MajorFunctionJP = req.MajorFunction
		testCase.MiddleFunctionJP = req.MiddleFunction
		testCase.MinorFunctionJP = req.MinorFunction
		testCase.PreconditionJP = req.Precondition
		testCase.TestStepsJP = req.TestSteps
		testCase.ExpectedResultJP = req.ExpectedResult
	}

	if err := s.repo.Create(testCase); err != nil {
		return nil, fmt.Errorf("create case with multilingual fields: %w", err)
	}

	return []*models.ManualTestCase{testCase}, nil
}

// UpdateCase 更新用例指定字段(部分更新) - 使用CaseID(UUID)
func (s *manualTestCaseService) UpdateCase(projectID uint, userID uint, caseID string, req UpdateCaseRequest) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 查询用例是否存在且属于当前项目（使用CaseID）
	testCase, err := s.repo.GetByCaseID(caseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用例不存在")
		}
		return fmt.Errorf("get test case: %w", err)
	}

	if testCase.ProjectID != projectID {
		return errors.New("用例不属于当前项目")
	}

	// 构建updates map (仅包含非nil字段)
	updates := make(map[string]interface{})
	if req.CaseNumber != nil {
		updates["case_number"] = *req.CaseNumber
	}

	// 根据用例类型决定更新哪些字段
	if testCase.CaseType == "overall" || testCase.CaseType == "change" || testCase.CaseType == "acceptance" {
		// overall/change/acceptance用例：支持多语言字段直接更新
		// 优先使用多语言字段(major_function_cn等)，如果不存在则使用单语言字段+language映射
		if req.MajorFunctionCN != nil {
			updates["major_function_cn"] = *req.MajorFunctionCN
		}
		if req.MajorFunctionJP != nil {
			updates["major_function_jp"] = *req.MajorFunctionJP
		}
		if req.MajorFunctionEN != nil {
			updates["major_function_en"] = *req.MajorFunctionEN
		}
		if req.MiddleFunctionCN != nil {
			updates["middle_function_cn"] = *req.MiddleFunctionCN
		}
		if req.MiddleFunctionJP != nil {
			updates["middle_function_jp"] = *req.MiddleFunctionJP
		}
		if req.MiddleFunctionEN != nil {
			updates["middle_function_en"] = *req.MiddleFunctionEN
		}
		if req.MinorFunctionCN != nil {
			updates["minor_function_cn"] = *req.MinorFunctionCN
		}
		if req.MinorFunctionJP != nil {
			updates["minor_function_jp"] = *req.MinorFunctionJP
		}
		if req.MinorFunctionEN != nil {
			updates["minor_function_en"] = *req.MinorFunctionEN
		}
		if req.PreconditionCN != nil {
			updates["precondition_cn"] = *req.PreconditionCN
		}
		if req.PreconditionJP != nil {
			updates["precondition_jp"] = *req.PreconditionJP
		}
		if req.PreconditionEN != nil {
			updates["precondition_en"] = *req.PreconditionEN
		}
		if req.TestStepsCN != nil {
			updates["test_steps_cn"] = *req.TestStepsCN
		}
		if req.TestStepsJP != nil {
			updates["test_steps_jp"] = *req.TestStepsJP
		}
		if req.TestStepsEN != nil {
			updates["test_steps_en"] = *req.TestStepsEN
		}
		if req.ExpectedResultCN != nil {
			updates["expected_result_cn"] = *req.ExpectedResultCN
		}
		if req.ExpectedResultJP != nil {
			updates["expected_result_jp"] = *req.ExpectedResultJP
		}
		if req.ExpectedResultEN != nil {
			updates["expected_result_en"] = *req.ExpectedResultEN
		}

		// 如果前端使用单语言字段，根据language参数映射到对应的多语言字段
		// language参数从查询中获取，默认为"中文"
		if req.MajorFunction != nil {
			switch "中文" { // 默认使用中文
			case "中文":
				updates["major_function_cn"] = *req.MajorFunction
			case "English":
				updates["major_function_en"] = *req.MajorFunction
			case "日本語":
				updates["major_function_jp"] = *req.MajorFunction
			}
		}
		if req.MiddleFunction != nil {
			switch "中文" { // 默认使用中文
			case "中文":
				updates["middle_function_cn"] = *req.MiddleFunction
			case "English":
				updates["middle_function_en"] = *req.MiddleFunction
			case "日本語":
				updates["middle_function_jp"] = *req.MiddleFunction
			}
		}
		if req.MinorFunction != nil {
			switch "中文" {
			case "中文":
				updates["minor_function_cn"] = *req.MinorFunction
			case "English":
				updates["minor_function_en"] = *req.MinorFunction
			case "日本語":
				updates["minor_function_jp"] = *req.MinorFunction
			}
		}
		if req.Precondition != nil {
			switch "中文" {
			case "中文":
				updates["precondition_cn"] = *req.Precondition
			case "English":
				updates["precondition_en"] = *req.Precondition
			case "日本語":
				updates["precondition_jp"] = *req.Precondition
			}
		}
		if req.TestSteps != nil {
			switch "中文" {
			case "中文":
				updates["test_steps_cn"] = *req.TestSteps
			case "English":
				updates["test_steps_en"] = *req.TestSteps
			case "日本語":
				updates["test_steps_jp"] = *req.TestSteps
			}
		}
		if req.ExpectedResult != nil {
			switch "中文" {
			case "中文":
				updates["expected_result_cn"] = *req.ExpectedResult
			case "English":
				updates["expected_result_en"] = *req.ExpectedResult
			case "日本語":
				updates["expected_result_jp"] = *req.ExpectedResult
			}
		}
	} else {
		// AI用例：直接更新单语言字段
		if req.MajorFunction != nil {
			updates["major_function"] = *req.MajorFunction
		}
		if req.MiddleFunction != nil {
			updates["middle_function"] = *req.MiddleFunction
		}
		if req.MinorFunction != nil {
			updates["minor_function"] = *req.MinorFunction
		}
		if req.Precondition != nil {
			updates["precondition"] = *req.Precondition
		}
		if req.TestSteps != nil {
			updates["test_steps"] = *req.TestSteps
		}
		if req.ExpectedResult != nil {
			updates["expected_result"] = *req.ExpectedResult
		}
	}

	// 测试结果和备注字段对所有类型通用
	if req.TestResult != nil {
		updates["test_result"] = *req.TestResult
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	if len(updates) == 0 {
		return nil // 没有字段需要更新
	}

	if err := s.repo.UpdateByCaseID(caseID, updates); err != nil {
		return fmt.Errorf("update test case: %w", err)
	}

	return nil
}

// DeleteCase 软删除用例 - 使用CaseID(UUID)
func (s *manualTestCaseService) DeleteCase(projectID uint, userID uint, caseID string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 查询用例是否存在且属于当前项目（使用CaseID）
	testCase, err := s.repo.GetByCaseID(caseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用例不存在")
		}
		return fmt.Errorf("get test case: %w", err)
	}

	if testCase.ProjectID != projectID {
		return errors.New("用例不属于当前项目")
	}

	// 删除用例并自动重排
	if err := s.repo.DeleteByCaseID(caseID); err != nil {
		return fmt.Errorf("delete test case: %w", err)
	}

	// 删除后自动调整后续用例的display_order(使用ID作为排序字段)
	// 注意: 当前模型使用ID字段作为display_order,后续需要添加独立的display_order字段
	if err := s.repo.DecrementOrderAfter(projectID, testCase.CaseType, int(testCase.ID)); err != nil {
		return fmt.Errorf("decrement order after deletion: %w", err)
	}

	// 重新分配所有用例的id字段（确保No列连续）
	if err := s.repo.ReassignDisplayIDs(projectID, testCase.CaseType); err != nil {
		return fmt.Errorf("reassign display ids: %w", err)
	}

	return nil
}

// deleteWithTranslations 多语言联动删除,通过major_function字段关联查询并批量删除三个语言版本
func (s *manualTestCaseService) deleteWithTranslations(projectID uint, testCase *models.ManualTestCase) error {
	// TODO: GetByCriteria method needs to be implemented in repository
	// 暂时只删除当前用例
	if err := s.repo.DeleteByCaseID(testCase.CaseID); err != nil {
		return fmt.Errorf("delete case: %w", err)
	}

	return nil
}

// ReorderCases 按照指定顺序重新分配ID
func (s *manualTestCaseService) ReorderCases(projectID uint, userID uint, caseType string, caseIDs []uint) ([]uint, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 获取当前类型所有用例
	allCases, err := s.repo.GetByProjectAndType(projectID, caseType)
	if err != nil {
		return nil, fmt.Errorf("get cases by project and type: %w", err)
	}

	// 验证所有caseID归属当前项目且类型匹配
	existingIDs := make(map[uint]bool)
	for _, c := range allCases {
		existingIDs[c.ID] = true
	}

	for _, id := range caseIDs {
		if !existingIDs[id] {
			return nil, fmt.Errorf("case_id %d does not belong to project %d with type %s", id, projectID, caseType)
		}
	}

	// 构建ID映射表: map[oldID]newID (newID从1开始递增)
	caseIDMap := make(map[uint]uint)
	for i, oldID := range caseIDs {
		newID := uint(i + 1)
		caseIDMap[oldID] = newID
	}

	// 批量更新ID
	if err := s.repo.BatchUpdateIDs(caseIDMap); err != nil {
		return nil, fmt.Errorf("batch update ids: %w", err)
	}

	// 返回新的ID序列
	newIDs := make([]uint, len(caseIDs))
	for i := range caseIDs {
		newIDs[i] = uint(i + 1)
	}

	return newIDs, nil
}

// ReorderCasesByDrag 拖拽重排：根据case_id顺序重新分配ID
func (s *manualTestCaseService) ReorderCasesByDrag(projectID uint, userID uint, caseType string, caseIDOrder []string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 调用Repository的新方法，按case_id顺序重新分配ID
	if err := s.repo.BatchUpdateIDsByCaseID(caseIDOrder); err != nil {
		return fmt.Errorf("batch update ids by case id: %w", err)
	}

	return nil
}

// ReorderAllCasesByID 按现有ID顺序重新编号所有用例
func (s *manualTestCaseService) ReorderAllCasesByID(projectID uint, userID uint, caseType string, language string) (int, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return 0, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return 0, errors.New("无项目访问权限")
	}

	// 获取所有用例（按ID排序）- 使用已有的方法
	allCases, err := s.repo.GetByProjectAndTypeOrdered(projectID, caseType)
	if err != nil {
		return 0, fmt.Errorf("get all cases: %w", err)
	}

	if len(allCases) == 0 {
		return 0, nil
	}

	// 按现有ID排序（已经在repo层排序）
	// 提取case_id数组
	caseIDs := make([]string, len(allCases))
	for i, c := range allCases {
		caseIDs[i] = c.CaseID
	}

	// 调用Repository的批量更新方法
	if err := s.repo.BatchUpdateIDsByCaseID(caseIDs); err != nil {
		return 0, fmt.Errorf("batch update ids: %w", err)
	}

	return len(caseIDs), nil
}

// ClearAICases 清空指定项目的AI用例,返回删除数量
func (s *manualTestCaseService) ClearAICases(projectID uint, userID uint) (int, error) {
	// 验证用户权限
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return 0, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return 0, errors.New("无项目访问权限")
	}

	// 先查询要删除的用例数量
	cases, err := s.repo.GetByProjectAndType(projectID, "ai")
	if err != nil {
		return 0, fmt.Errorf("get ai cases: %w", err)
	}
	deletedCount := len(cases)

	// 删除所有AI用例（软删除）
	if err := s.repo.DeleteByCaseType(projectID, "ai"); err != nil {
		return 0, fmt.Errorf("delete ai cases: %w", err)
	}

	return deletedCount, nil
}

// InsertCase 在指定位置插入新用例
func (s *manualTestCaseService) InsertCase(projectID uint, userID uint, caseType string, position string, targetCaseID string, language string) (*models.ManualTestCase, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 查询目标用例
	targetCase, err := s.repo.GetByCaseID(targetCaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目标用例不存在")
		}
		return nil, fmt.Errorf("get target case: %w", err)
	}

	// 计算新用例的display_order(使用ID作为临时排序字段)
	var newOrder int
	if position == "before" {
		newOrder = int(targetCase.ID)
	} else { // after
		newOrder = int(targetCase.ID) + 1
	}

	// 调整现有用例的order
	if position == "before" {
		if err := s.repo.IncrementOrderAfter(projectID, caseType, int(targetCase.ID)-1); err != nil {
			return nil, fmt.Errorf("increment order: %w", err)
		}
	} else {
		if err := s.repo.IncrementOrderAfter(projectID, caseType, int(targetCase.ID)); err != nil {
			return nil, fmt.Errorf("increment order: %w", err)
		}
	}

	// 创建新用例(设置默认值)
	newCase := &models.ManualTestCase{
		CaseID:     uuid.New().String(),
		ProjectID:  projectID,
		CaseType:   caseType,
		ID:         uint(newOrder),
		TestResult: "NR",
	}

	// 保存新用例
	if err := s.repo.Create(newCase); err != nil {
		return nil, fmt.Errorf("create new case: %w", err)
	}

	// 重新分配所有用例的id字段
	if err := s.repo.ReassignDisplayIDs(projectID, caseType); err != nil {
		return nil, fmt.Errorf("reassign display ids: %w", err)
	}

	return newCase, nil
}

// BatchDeleteCases 批量删除用例
func (s *manualTestCaseService) BatchDeleteCases(projectID uint, userID uint, caseType string, caseIDs []string) (int, []string, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return 0, nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return 0, nil, errors.New("无项目访问权限")
	}

	deletedCount := 0
	var failedCaseIDs []string

	// 循环删除每个用例
	for _, caseID := range caseIDs {
		if err := s.repo.DeleteByCaseID(caseID); err != nil {
			failedCaseIDs = append(failedCaseIDs, caseID)
		} else {
			deletedCount++
		}
	}

	// 重新分配所有用例的id字段
	if err := s.repo.ReassignDisplayIDs(projectID, caseType); err != nil {
		return deletedCount, failedCaseIDs, fmt.Errorf("reassign display ids: %w", err)
	}

	return deletedCount, failedCaseIDs, nil
}

// ReassignAllIDs 重新分配所有用例的ID
func (s *manualTestCaseService) ReassignAllIDs(projectID uint, userID uint, caseType string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 调用repository重新分配ID
	if err := s.repo.ReassignDisplayIDs(projectID, caseType); err != nil {
		return fmt.Errorf("reassign display ids: %w", err)
	}

	return nil
}
