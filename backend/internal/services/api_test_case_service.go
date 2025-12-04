package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"webtest/internal/models"
	"webtest/internal/repositories"

	"gorm.io/gorm"
)

// ApiTestCaseService 接口测试用例服务接口
type ApiTestCaseService interface {
	// 用例管理
	GetCases(projectID uint, userID uint, caseType string, page int, size int) ([]*models.ApiTestCase, int64, error)
	CreateCase(projectID uint, userID uint, testCase *models.ApiTestCase) (*models.ApiTestCase, error)
	InsertCase(projectID uint, userID uint, caseType string, position string, targetCaseID string, caseData map[string]interface{}) (*models.ApiTestCase, error)
	DeleteCase(projectID uint, userID uint, caseID string) error
	BatchDeleteCases(projectID uint, userID uint, caseType string, caseIDs []string) (int, []string, error)
	UpdateCase(projectID uint, userID uint, caseID string, updates map[string]interface{}) error

	// 版本管理
	SaveVersion(projectID uint, userID uint, remark string) (*models.ApiTestCaseVersion, error)
	GetVersions(projectID uint, userID uint, page int, size int) ([]*models.ApiTestCaseVersion, int64, error)
	DeleteVersion(projectID uint, userID uint, versionID string) error
	UpdateVersionRemark(projectID uint, userID uint, versionID string, remark string) error
}

type apiTestCaseService struct {
	repo           repositories.ApiTestCaseRepository
	projectService ProjectService
}

// NewApiTestCaseService 创建服务实例
func NewApiTestCaseService(repo repositories.ApiTestCaseRepository, projectService ProjectService) ApiTestCaseService {
	return &apiTestCaseService{
		repo:           repo,
		projectService: projectService,
	}
}

// GetCases 分页查询用例列表
func (s *apiTestCaseService) GetCases(projectID uint, userID uint, caseType string, page int, size int) ([]*models.ApiTestCase, int64, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, 0, errors.New("无项目访问权限")
	}

	// 计算offset
	offset := (page - 1) * size

	// 查询数据(按display_order排序)
	cases, total, err := s.repo.List(projectID, caseType, offset, size)
	if err != nil {
		return nil, 0, fmt.Errorf("list api cases: %w", err)
	}

	return cases, total, nil
}

// CreateCase 创建用例
func (s *apiTestCaseService) CreateCase(projectID uint, userID uint, testCase *models.ApiTestCase) (*models.ApiTestCase, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 字段验证
	if err := testCase.Validate(); err != nil {
		return nil, fmt.Errorf("validate case data: %w", err)
	}

	// 获取当前类型的最大display_order
	cases, _, err := s.repo.List(projectID, testCase.CaseType, 0, 1)
	if err == nil && len(cases) > 0 {
		// 如果有数据,获取最大display_order并+1
		allCases, _ := s.repo.GetByProjectAndType(projectID, testCase.CaseType)
		maxOrder := 0
		for _, c := range allCases {
			if c.DisplayOrder > maxOrder {
				maxOrder = c.DisplayOrder
			}
		}
		testCase.DisplayOrder = maxOrder + 1
	} else {
		// 如果没有数据,设置为1
		testCase.DisplayOrder = 1
	}

	// 保存新用例
	if err := s.repo.Create(testCase); err != nil {
		return nil, fmt.Errorf("create api test case: %w", err)
	}

	return testCase, nil
}

// InsertCase 插入用例(指定位置)
func (s *apiTestCaseService) InsertCase(projectID uint, userID uint, caseType string, position string, targetCaseID string, caseData map[string]interface{}) (*models.ApiTestCase, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 查询目标用例
	targetCase, err := s.repo.GetByID(targetCaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目标用例不存在")
		}
		return nil, fmt.Errorf("get target case: %w", err)
	}

	// 计算新用例的display_order
	var newOrder int
	if position == "before" {
		newOrder = targetCase.DisplayOrder
	} else { // after
		newOrder = targetCase.DisplayOrder + 1
	}

	// 调整现有用例的display_order
	if position == "before" {
		if err := s.repo.IncrementOrderAfter(projectID, caseType, targetCase.DisplayOrder-1); err != nil {
			return nil, fmt.Errorf("increment order: %w", err)
		}
	} else {
		if err := s.repo.IncrementOrderAfter(projectID, caseType, targetCase.DisplayOrder); err != nil {
			return nil, fmt.Errorf("increment order: %w", err)
		}
	}

	// 创建新用例(UUID自动生成)
	newCase := &models.ApiTestCase{
		ProjectID:    projectID,
		CaseType:     caseType,
		DisplayOrder: newOrder,
		TestResult:   "NR",
		Method:       "GET",
	}

	// 填充caseData字段
	if caseNumber, ok := caseData["case_number"].(string); ok {
		newCase.CaseNumber = caseNumber
	}
	if screen, ok := caseData["screen"].(string); ok {
		newCase.Screen = screen
	}
	if url, ok := caseData["url"].(string); ok {
		newCase.URL = url
	}
	if header, ok := caseData["header"].(string); ok {
		newCase.Header = header
	}
	if method, ok := caseData["method"].(string); ok {
		newCase.Method = method
	}
	if body, ok := caseData["body"].(string); ok {
		newCase.Body = body
	}
	if response, ok := caseData["response"].(string); ok {
		newCase.Response = response
	}
	if testResult, ok := caseData["test_result"].(string); ok {
		newCase.TestResult = testResult
	}
	if remark, ok := caseData["remark"].(string); ok {
		newCase.Remark = remark
	}

	// 字段验证
	if err := newCase.Validate(); err != nil {
		return nil, fmt.Errorf("validate case data: %w", err)
	}

	// 保存新用例
	if err := s.repo.Create(newCase); err != nil {
		return nil, fmt.Errorf("create new case: %w", err)
	}

	// 重新分配display_order
	if err := s.repo.ReassignDisplayOrders(projectID, caseType); err != nil {
		return nil, fmt.Errorf("reassign display orders: %w", err)
	}

	return newCase, nil
}

// DeleteCase 删除单个用例
func (s *apiTestCaseService) DeleteCase(projectID uint, userID uint, caseID string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 先检查用例是否存在并属于该项目
	existingCase, err := s.repo.GetByID(caseID)
	if err != nil {
		return errors.New("用例不存在")
	}
	if existingCase.ProjectID != projectID {
		return errors.New("无权限删除该用例")
	}

	// 获取用例类型用于后续重排序
	caseType := existingCase.CaseType

	// 删除用例
	if err := s.repo.Delete(caseID); err != nil {
		return fmt.Errorf("delete case: %w", err)
	}

	// 重新分配display_order
	if err := s.repo.ReassignDisplayOrders(projectID, caseType); err != nil {
		return fmt.Errorf("reassign display orders: %w", err)
	}

	return nil
}

// BatchDeleteCases 批量删除用例
func (s *apiTestCaseService) BatchDeleteCases(projectID uint, userID uint, caseType string, caseIDs []string) (int, []string, error) {
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
		if err := s.repo.Delete(caseID); err != nil {
			failedCaseIDs = append(failedCaseIDs, caseID)
		} else {
			deletedCount++
		}
	}

	// 重新分配display_order
	if err := s.repo.ReassignDisplayOrders(projectID, caseType); err != nil {
		return deletedCount, failedCaseIDs, fmt.Errorf("reassign display orders: %w", err)
	}

	return deletedCount, failedCaseIDs, nil
}

// UpdateCase 更新用例
func (s *apiTestCaseService) UpdateCase(projectID uint, userID uint, caseID string, updates map[string]interface{}) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 更新用例
	if err := s.repo.Update(caseID, updates); err != nil {
		return fmt.Errorf("update case: %w", err)
	}

	return nil
}

// ========== 版本管理 ==========

// SaveVersion 保存版本(生成4个CSV文件)
func (s *apiTestCaseService) SaveVersion(projectID uint, userID uint, remark string) (*models.ApiTestCaseVersion, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 获取项目信息(用于文件名)
	project, _, err := s.projectService.GetByID(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	// 生成时间戳
	timestamp := time.Now().Format("20060102_150405")

	// 创建存储目录
	storageDir := filepath.Join("storage", "versions", "api-cases")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}

	// 生成四个CSV文件
	roleTypes := []string{"role1", "role2", "role3", "role4"}
	filenames := make(map[string]string, 4)

	for _, roleType := range roleTypes {
		// 查询该ROLE的所有用例
		cases, err := s.repo.GetByProjectAndType(projectID, roleType)
		if err != nil {
			return nil, fmt.Errorf("get cases for %s: %w", roleType, err)
		}

		// 生成CSV文件名
		upperRole := fmt.Sprintf("ROLE%s", roleType[len(roleType)-1:])
		filename := fmt.Sprintf("%s_APITestCase_%s_%s.csv", project.Name, upperRole, timestamp)
		filePath := filepath.Join(storageDir, filename)

		// 创建CSV文件
		file, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("create csv file: %w", err)
		}
		defer file.Close()

		// 写入CSV内容
		writer := csv.NewWriter(file)
		defer writer.Flush()

		// 写入表头 (使用英文列名保持与前端一致)
		header := []string{"No.", "CaseID", "Screen", "URL", "Header", "Method", "Body", "Response", "TestResult", "Remark"}
		if err := writer.Write(header); err != nil {
			return nil, fmt.Errorf("write csv header: %w", err)
		}

		// 写入数据行
		for i, c := range cases {
			row := []string{
				fmt.Sprintf("%d", i+1), // No.列基于display_order排序后的序号
				c.CaseNumber,
				c.Screen,
				c.URL,
				c.Header,
				c.Method,
				c.Body,
				c.Response,
				c.TestResult,
				c.Remark,
			}
			if err := writer.Write(row); err != nil {
				return nil, fmt.Errorf("write csv row: %w", err)
			}
		}

		// 保存文件名
		filenames[roleType] = filename
	}

	// 创建版本记录(UUID自动生成)
	version := &models.ApiTestCaseVersion{
		ProjectID:     projectID,
		FilenameRole1: filenames["role1"],
		FilenameRole2: filenames["role2"],
		FilenameRole3: filenames["role3"],
		FilenameRole4: filenames["role4"],
		Remark:        remark,
		CreatedBy:     userID,
	}

	if err := s.repo.CreateVersion(version); err != nil {
		return nil, fmt.Errorf("create version: %w", err)
	}

	return version, nil
}

// GetVersions 获取版本列表
func (s *apiTestCaseService) GetVersions(projectID uint, userID uint, page int, size int) ([]*models.ApiTestCaseVersion, int64, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, 0, errors.New("无项目访问权限")
	}

	// 计算offset
	offset := (page - 1) * size

	// 查询版本列表
	versions, total, err := s.repo.ListVersions(projectID, offset, size)
	if err != nil {
		return nil, 0, fmt.Errorf("list versions: %w", err)
	}

	return versions, total, nil
}

// DeleteVersion 删除版本
func (s *apiTestCaseService) DeleteVersion(projectID uint, userID uint, versionID string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 查询版本记录
	version, err := s.repo.GetVersionByID(versionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("版本不存在")
		}
		return fmt.Errorf("get version: %w", err)
	}

	// 验证版本归属
	if version.ProjectID != projectID {
		return errors.New("版本不属于该项目")
	}

	// 删除CSV文件
	storageDir := filepath.Join("storage", "versions", "api-cases")
	filenames := []string{
		version.FilenameRole1,
		version.FilenameRole2,
		version.FilenameRole3,
		version.FilenameRole4,
	}

	for _, filename := range filenames {
		if filename != "" {
			filePath := filepath.Join(storageDir, filename)
			os.Remove(filePath) // 忽略删除错误
		}
	}

	// 删除版本记录
	if err := s.repo.DeleteVersion(versionID); err != nil {
		return fmt.Errorf("delete version: %w", err)
	}

	return nil
}

// UpdateVersionRemark 更新版本备注
func (s *apiTestCaseService) UpdateVersionRemark(projectID uint, userID uint, versionID string, remark string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 查询版本记录
	version, err := s.repo.GetVersionByID(versionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("版本不存在")
		}
		return fmt.Errorf("get version: %w", err)
	}

	// 验证版本归属
	if version.ProjectID != projectID {
		return errors.New("版本不属于该项目")
	}

	// 更新备注(限制500字符由前端校验)
	if err := s.repo.UpdateVersionRemark(versionID, remark); err != nil {
		return fmt.Errorf("update version remark: %w", err)
	}

	return nil
}
