package services

import (
	"errors"
	"testing"
	"webtest/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockManualTestCaseRepository 模拟Repository层
type MockManualTestCaseRepository struct {
	mock.Mock
}

func (m *MockManualTestCaseRepository) GetMetadataByProjectID(projectID uint, caseType string) (*models.ManualTestCase, error) {
	args := m.Called(projectID, caseType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ManualTestCase), args.Error(1)
}

func (m *MockManualTestCaseRepository) UpdateMetadata(projectID uint, caseType string, metadata map[string]interface{}) error {
	args := m.Called(projectID, caseType, metadata)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) CreateDefaultMetadata(projectID uint, caseType string) error {
	args := m.Called(projectID, caseType)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) GetCasesByType(projectID uint, caseType string, offset int, limit int) ([]*models.ManualTestCase, int64, error) {
	args := m.Called(projectID, caseType, offset, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.ManualTestCase), int64(args.Int(1)), args.Error(2)
}

func (m *MockManualTestCaseRepository) Create(testCase *models.ManualTestCase) error {
	args := m.Called(testCase)
	if args.Error(0) == nil && testCase.ID == 0 {
		testCase.ID = 1 // 模拟自增ID
	}
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) GetByID(caseID uint) (*models.ManualTestCase, error) {
	args := m.Called(caseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ManualTestCase), args.Error(1)
}

func (m *MockManualTestCaseRepository) UpdateByID(caseID uint, updates map[string]interface{}) error {
	args := m.Called(caseID, updates)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) DeleteByID(caseID uint) error {
	args := m.Called(caseID)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) CreateBatch(testCases []*models.ManualTestCase) error {
	args := m.Called(testCases)
	if args.Error(0) == nil {
		for i, tc := range testCases {
			if tc.ID == 0 {
				tc.ID = uint(i + 1) // 模拟自增ID
			}
		}
	}
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) GetByCriteria(projectID uint, caseType string, majorFunction string, languages []string) ([]*models.ManualTestCase, error) {
	args := m.Called(projectID, caseType, majorFunction, languages)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ManualTestCase), args.Error(1)
}

func (m *MockManualTestCaseRepository) DeleteBatch(caseIDs []uint) error {
	args := m.Called(caseIDs)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) GetByProjectAndType(projectID uint, caseType string) ([]*models.ManualTestCase, error) {
	args := m.Called(projectID, caseType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ManualTestCase), args.Error(1)
}

func (m *MockManualTestCaseRepository) BatchUpdateIDs(caseIDMap map[uint]uint) error {
	args := m.Called(caseIDMap)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) DeleteByCaseType(projectID uint, caseType string) error {
	args := m.Called(projectID, caseType)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) BatchUpdateIDsByCaseID(caseIDs []string) error {
	args := m.Called(caseIDs)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) GetByCaseID(caseID string) (*models.ManualTestCase, error) {
	args := m.Called(caseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ManualTestCase), args.Error(1)
}

func (m *MockManualTestCaseRepository) UpdateByCaseID(caseID string, updates map[string]interface{}) error {
	args := m.Called(caseID, updates)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) DeleteByCaseID(caseID string) error {
	args := m.Called(caseID)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) GetMaxIDByProjectAndType(projectID uint, caseType string) (uint, error) {
	args := m.Called(projectID, caseType)
	return uint(args.Int(0)), args.Error(1)
}

func (m *MockManualTestCaseRepository) GetByProjectAndTypeOrdered(projectID uint, caseType string) ([]*models.ManualTestCase, error) {
	args := m.Called(projectID, caseType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ManualTestCase), args.Error(1)
}

func (m *MockManualTestCaseRepository) ReassignDisplayIDs(projectID uint, caseType string) error {
	args := m.Called(projectID, caseType)
	return args.Error(0)
}

func (m *MockManualTestCaseRepository) DeleteBatchByCaseIDs(caseIDs []string) error {
	args := m.Called(caseIDs)
	return args.Error(0)
}

// MockProjectService 模拟ProjectService
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) GetUserProjects(userID uint, role string) ([]models.Project, error) {
	args := m.Called(userID, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Project), args.Error(1)
}

func (m *MockProjectService) IsProjectMember(projectID uint, userID uint) (bool, error) {
	args := m.Called(projectID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectService) CreateProject(name string, description string, creatorID uint) (*models.Project, error) {
	args := m.Called(name, description, creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectService) UpdateProject(projectID uint, newName string, userID uint, role string) (*models.Project, error) {
	args := m.Called(projectID, newName, userID, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectService) DeleteProject(projectID uint, userID uint, role string) error {
	args := m.Called(projectID, userID, role)
	return args.Error(0)
}

func (m *MockProjectService) GetByID(projectID uint, userID uint) (*models.Project, string, error) {
	args := m.Called(projectID, userID)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.Project), args.String(1), args.Error(2)
}

func (m *MockProjectService) GetProjectMembers(projectID uint, userID uint) (*ProjectMembersResponse, error) {
	args := m.Called(projectID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ProjectMembersResponse), args.Error(1)
}

// 测试用例: CreateCase - AI用例创建(正常流程)
func TestCreateCase_AIType_Success(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock Create方法
	mockRepo.On("Create", mock.AnythingOfType("*models.ManualTestCase")).Return(nil)

	// 执行测试
	result, err := service.CreateCase(1, 123, CreateCaseRequest{
		CaseType:      "ai",
		Language:      "中文",
		MajorFunction: "登录功能",
		TestResult:    "NR",
	})

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotZero(t, result.ID) // 验证ID已生成
	assert.Equal(t, "登录功能", result.MajorFunction)
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: CreateCase - 整体用例多语言联动创建(AC-07)
func TestCreateCase_OverallType_MultiLanguage_Success(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock CreateBatch方法,期望创建1条记录(Language字段已移除)
	mockRepo.On("CreateBatch", mock.MatchedBy(func(cases []*models.ManualTestCase) bool {
		if len(cases) != 1 {
			return false
		}
		// 验证业务字段一致
		for _, c := range cases {
			if c.MajorFunctionCN != "登录功能" || c.ProjectID != 1 {
				return false
			}
		}
		return true
	})).Return(nil)

	// 执行测试
	result, err := service.CreateCase(1, 123, CreateCaseRequest{
		CaseType:      "overall",
		Language:      "中文",
		MajorFunction: "登录功能",
		TestResult:    "NR",
	})

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotZero(t, result.ID)
	assert.Equal(t, "登录功能", result.MajorFunction)
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: CreateCase - 权限校验失败
func TestCreateCase_PermissionDenied(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回false
	mockProjectService.On("IsProjectMember", uint(1), uint(999)).Return(false, nil)

	// 执行测试
	result, err := service.CreateCase(1, 999, CreateCaseRequest{
		CaseType:      "ai",
		Language:      "中文",
		MajorFunction: "登录功能",
	})

	// 断言
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "无项目访问权限")
	mockProjectService.AssertExpectations(t)
}

// 测试用例: UpdateCase - 正常流程
func TestUpdateCase_Success(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock GetByID返回用例
	existingCase := &models.ManualTestCase{
		ID:        10,
		ProjectID: 1,
		CaseType:  "ai",
		Language:  "中文",
	}
	mockRepo.On("GetByID", uint(10)).Return(existingCase, nil)

	// Mock UpdateByID
	testResult := "Pass"
	mockRepo.On("UpdateByID", uint(10), map[string]interface{}{
		"test_result": &testResult,
	}).Return(nil)

	// 执行测试
	err := service.UpdateCase(1, 123, 10, UpdateCaseRequest{
		TestResult: &testResult,
	})

	// 断言
	assert.NoError(t, err)
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: UpdateCase - 用例不存在
func TestUpdateCase_CaseNotFound(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock GetByID返回错误
	mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("record not found"))

	// 执行测试
	testResult := "Pass"
	err := service.UpdateCase(1, 123, 999, UpdateCaseRequest{
		TestResult: &testResult,
	})

	// 断言
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用例不存在")
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: DeleteCase - AI用例删除(正常流程)
func TestDeleteCase_AIType_Success(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock GetByID返回AI用例
	existingCase := &models.ManualTestCase{
		ID:        10,
		ProjectID: 1,
		CaseType:  "ai",
		Language:  "中文",
	}
	mockRepo.On("GetByID", uint(10)).Return(existingCase, nil)

	// Mock DeleteByID
	mockRepo.On("DeleteByID", uint(10)).Return(nil)

	// 执行测试
	err := service.DeleteCase(1, 123, 10)

	// 断言
	assert.NoError(t, err)
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: DeleteCase - 整体用例多语言联动删除(AC-07)
func TestDeleteCase_OverallType_MultiLanguage_Success(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock GetByID返回整体用例
	existingCase := &models.ManualTestCase{
		ID:            10,
		ProjectID:     1,
		CaseType:      "overall",
		Language:      "中文",
		MajorFunction: "登录功能",
	}
	mockRepo.On("GetByID", uint(10)).Return(existingCase, nil)

	// Mock GetByCriteria返回3个语言版本
	relatedCases := []*models.ManualTestCase{
		{ID: 10, Language: "中文", MajorFunction: "登录功能"},
		{ID: 11, Language: "English", MajorFunction: "登录功能"},
		{ID: 12, Language: "日本語", MajorFunction: "登录功能"},
	}
	mockRepo.On("GetByCriteria", uint(1), "overall", "登录功能", []string{"中文", "English", "日本語"}).
		Return(relatedCases, nil)

	// Mock DeleteBatch删除3条记录
	mockRepo.On("DeleteBatch", []uint{10, 11, 12}).Return(nil)

	// 执行测试
	err := service.DeleteCase(1, 123, 10)

	// 断言
	assert.NoError(t, err)
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: ReorderCases - 正常流程
func TestReorderCases_Success(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock GetByProjectAndType返回用例列表
	existingCases := []*models.ManualTestCase{
		{ID: 5, ProjectID: 1, CaseType: "ai"},
		{ID: 3, ProjectID: 1, CaseType: "ai"},
		{ID: 8, ProjectID: 1, CaseType: "ai"},
	}
	mockRepo.On("GetByProjectAndType", uint(1), "ai").Return(existingCases, nil)

	// Mock BatchUpdateIDs
	mockRepo.On("BatchUpdateIDs", mock.MatchedBy(func(idMap map[uint]uint) bool {
		// 验证映射关系: 5->1, 3->2, 8->3
		return len(idMap) == 3 && idMap[5] == 1 && idMap[3] == 2 && idMap[8] == 3
	})).Return(nil)

	// 执行测试
	newIDs, err := service.ReorderCases(1, 123, "ai", []uint{5, 3, 8})

	// 断言
	assert.NoError(t, err)
	assert.Equal(t, []uint{1, 2, 3}, newIDs)
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: ReorderCases - caseID归属验证失败
func TestReorderCases_InvalidCaseID(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock GetByProjectAndType返回用例列表(不包含ID=999)
	existingCases := []*models.ManualTestCase{
		{ID: 5, ProjectID: 1, CaseType: "ai"},
		{ID: 3, ProjectID: 1, CaseType: "ai"},
	}
	mockRepo.On("GetByProjectAndType", uint(1), "ai").Return(existingCases, nil)

	// 执行测试(传入不存在的ID=999)
	newIDs, err := service.ReorderCases(1, 123, "ai", []uint{5, 3, 999})

	// 断言
	assert.Error(t, err)
	assert.Nil(t, newIDs)
	assert.Contains(t, err.Error(), "case_id=999 不属于项目")
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: CreateCase - AI用例使用单语言字段
func TestCreateCase_AIType_SingleLanguageFields(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock Create方法
	mockRepo.On("Create", mock.MatchedBy(func(tc *models.ManualTestCase) bool {
		// 验证AI用例使用单语言字段
		return tc.CaseType == "ai" &&
			tc.MajorFunction == "登录功能" &&
			tc.MiddleFunction == "用户名密码登录" &&
			tc.MinorFunction == "正常登录" &&
			tc.Precondition == "用户已注册" &&
			tc.TestSteps == "1.输入用户名\n2.输入密码\n3.点击登录" &&
			tc.ExpectedResult == "登录成功" &&
			tc.MajorFunctionCN == "" && // AI用例不使用多语言字段
			tc.MajorFunctionJP == "" &&
			tc.MajorFunctionEN == ""
	})).Return(nil)

	// 构造请求
	req := CreateCaseRequest{
		CaseType:       "ai",
		Language:       "中文",
		CaseNumber:     "TC001",
		MajorFunction:  "登录功能",
		MiddleFunction: "用户名密码登录",
		MinorFunction:  "正常登录",
		Precondition:   "用户已注册",
		TestSteps:      "1.输入用户名\n2.输入密码\n3.点击登录",
		ExpectedResult: "登录成功",
		Remark:         "基本功能测试",
	}

	// 执行测试
	result, err := service.CreateCase(1, 123, req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "登录功能", result.MajorFunction)
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: CreateCase - Overall用例使用多语言字段
func TestCreateCase_OverallType_MultiLanguageFields(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock Create方法
	mockRepo.On("Create", mock.MatchedBy(func(tc *models.ManualTestCase) bool {
		// 验证overall用例使用多语言字段
		return tc.CaseType == "overall" &&
			tc.MajorFunctionCN == "登录模块" &&
			tc.MajorFunction == "" && // overall用例不使用单语言字段
			tc.TestResult == "NR" // 默认值
	})).Return(nil)

	// 构造请求(仅填充中文字段)
	req := CreateCaseRequest{
		CaseType:   "overall",
		Language:   "中文",
		CaseNumber: "AC001",
		// 多语言字段
		MajorFunction: "登录模块", // 前端传此字段，后端映射到MajorFunctionCN
		TestResult:    "",     // 空值应设置为NR
		Remark:        "整体功能验证",
	}

	// 执行测试
	result, err := service.CreateCase(1, 123, req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: ClearAICases - 成功清空AI用例
func TestClearAICases_Success(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回true
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(true, nil)

	// Mock DeleteByCaseType返回nil(成功删除)
	mockRepo.On("DeleteByCaseType", uint(1), "ai").Return(nil)

	// 执行测试
	err := service.ClearAICases(1, 123)

	// 断言
	assert.NoError(t, err)
	mockProjectService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试用例: ClearAICases - 权限不足
func TestClearAICases_PermissionDenied(t *testing.T) {
	mockRepo := new(MockManualTestCaseRepository)
	mockProjectService := new(MockProjectService)

	service := &manualTestCaseService{
		repo:           mockRepo,
		projectService: mockProjectService,
	}

	// 权限校验返回false
	mockProjectService.On("IsProjectMember", uint(1), uint(123)).Return(false, nil)

	// 执行测试
	err := service.ClearAICases(1, 123)

	// 断言
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无项目访问权限")
	mockProjectService.AssertExpectations(t)
	// DeleteByCaseType不应该被调用
	mockRepo.AssertNotCalled(t, "DeleteByCaseType")
}
