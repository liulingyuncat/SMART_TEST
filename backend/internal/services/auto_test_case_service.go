package services

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AutoMetadataDTO 自动化用例元数据DTO
type AutoMetadataDTO struct {
	ScreenCN string `json:"screen_cn"`
	ScreenJP string `json:"screen_jp"`
	ScreenEN string `json:"screen_en"`
}

// UpdateAutoMetadataRequest 更新自动化用例元数据请求
type UpdateAutoMetadataRequest struct {
	ScreenCN string `json:"screen_cn" binding:"max=100"`
	ScreenJP string `json:"screen_jp" binding:"max=100"`
	ScreenEN string `json:"screen_en" binding:"max=100"`
}

// AutoCaseDTO 自动化用例DTO
type AutoCaseDTO struct {
	CaseID  string `json:"case_id"` // UUID主键
	ID      uint   `json:"id"`      // 显示序号
	CaseNum string `json:"case_num"`

	// 多语言字段
	ScreenCN         string `json:"screen_cn"`
	ScreenJP         string `json:"screen_jp"`
	ScreenEN         string `json:"screen_en"`
	FunctionCN       string `json:"function_cn"`
	FunctionJP       string `json:"function_jp"`
	FunctionEN       string `json:"function_en"`
	PreconditionCN   string `json:"precondition_cn"`
	PreconditionJP   string `json:"precondition_jp"`
	PreconditionEN   string `json:"precondition_en"`
	TestStepsCN      string `json:"test_steps_cn"`
	TestStepsJP      string `json:"test_steps_jp"`
	TestStepsEN      string `json:"test_steps_en"`
	ExpectedResultCN string `json:"expected_result_cn"`
	ExpectedResultJP string `json:"expected_result_jp"`
	ExpectedResultEN string `json:"expected_result_en"`

	TestResult string `json:"test_result"`
	Remark     string `json:"remark"`
}

// AutoCaseListDTO 自动化用例列表DTO
type AutoCaseListDTO struct {
	Cases []*AutoCaseDTO `json:"cases"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// CreateAutoCaseRequest 创建自动化用例请求
type CreateAutoCaseRequest struct {
	CaseType string `json:"case_type" binding:"required,oneof=role1 role2 role3 role4"`
	CaseNum  string `json:"case_num" binding:"max=50"`

	ScreenCN         string `json:"screen_cn" binding:"max=100"`
	ScreenJP         string `json:"screen_jp" binding:"max=100"`
	ScreenEN         string `json:"screen_en" binding:"max=100"`
	FunctionCN       string `json:"function_cn" binding:"max=100"`
	FunctionJP       string `json:"function_jp" binding:"max=100"`
	FunctionEN       string `json:"function_en" binding:"max=100"`
	PreconditionCN   string `json:"precondition_cn"`
	PreconditionJP   string `json:"precondition_jp"`
	PreconditionEN   string `json:"precondition_en"`
	TestStepsCN      string `json:"test_steps_cn"`
	TestStepsJP      string `json:"test_steps_jp"`
	TestStepsEN      string `json:"test_steps_en"`
	ExpectedResultCN string `json:"expected_result_cn"`
	ExpectedResultJP string `json:"expected_result_jp"`
	ExpectedResultEN string `json:"expected_result_en"`

	TestResult string `json:"test_result" binding:"omitempty,oneof=OK NG NR"`
	Remark     string `json:"remark"`
}

// UpdateAutoCaseRequest 更新自动化用例请求
type UpdateAutoCaseRequest struct {
	CaseNum *string `json:"case_num,omitempty" binding:"omitempty,max=50"`

	ScreenCN       *string `json:"screen_cn,omitempty" binding:"omitempty,max=100"`
	ScreenJP       *string `json:"screen_jp,omitempty" binding:"omitempty,max=100"`
	ScreenEN       *string `json:"screen_en,omitempty" binding:"omitempty,max=100"`
	FunctionCN     *string `json:"function_cn,omitempty" binding:"omitempty,max=100"`
	FunctionJP     *string `json:"function_jp,omitempty" binding:"omitempty,max=100"`
	FunctionEN     *string `json:"function_en,omitempty" binding:"omitempty,max=100"`
	PreconditionCN *string `json:"precondition_cn,omitempty"`
	PreconditionJP *string `json:"precondition_jp,omitempty"`
	PreconditionEN *string `json:"precondition_en,omitempty"`

	TestStepsCN      *string `json:"test_steps_cn,omitempty"`
	TestStepsJP      *string `json:"test_steps_jp,omitempty"`
	TestStepsEN      *string `json:"test_steps_en,omitempty"`
	ExpectedResultCN *string `json:"expected_result_cn,omitempty"`
	ExpectedResultJP *string `json:"expected_result_jp,omitempty"`
	ExpectedResultEN *string `json:"expected_result_en,omitempty"`

	TestResult *string `json:"test_result,omitempty" binding:"omitempty,oneof=OK NG NR"`
	Remark     *string `json:"remark,omitempty"`
}

// VersionFileInfo 版本文件信息
type VersionFileInfo struct {
	Role     string `json:"role"`
	Filename string `json:"filename"`
	Count    int    `json:"count"`
}

// VersionInfoDTO 版本信息DTO
type VersionInfoDTO struct {
	VersionID  string             `json:"version_id"`
	SavedAt    time.Time          `json:"saved_at"`
	Files      []*VersionFileInfo `json:"files"`
	TotalCases int                `json:"total_cases"`
}

// VersionDTO 版本列表项DTO
type VersionDTO struct {
	VersionID  string             `json:"version_id"`
	SavedAt    time.Time          `json:"saved_at"`
	Files      []*VersionFileInfo `json:"files"`
	TotalCases int                `json:"total_cases"`
	Remark     string             `json:"remark"`
}

// VersionListDTO 版本列表DTO
type VersionListDTO struct {
	Versions []*VersionDTO `json:"versions"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	Size     int           `json:"size"`
}

// AutoTestCaseService 自动化测试用例服务接口
type AutoTestCaseService interface {
	GetMetadata(projectID uint, userID uint, caseType string) (*AutoMetadataDTO, error)
	UpdateMetadata(projectID uint, userID uint, caseType string, req UpdateAutoMetadataRequest) error
	GetCases(projectID uint, userID uint, caseType string, page int, size int) (*AutoCaseListDTO, error)
	CreateCase(projectID uint, userID uint, req CreateAutoCaseRequest) (*AutoCaseDTO, error)
	UpdateCase(projectID uint, userID uint, caseID string, req UpdateAutoCaseRequest) error
	DeleteCase(projectID uint, userID uint, caseID string) error
	ReorderAllCases(projectID uint, userID uint, caseType string) (int, error)
	ReorderByIDs(projectID uint, userID uint, caseType string, caseIDs []string) (int, error)

	// 新增：插入和批量删除方法
	InsertCase(projectID uint, userID uint, caseType string, position string, targetCaseID string) (*models.AutoTestCase, error)
	BatchDeleteCases(projectID uint, userID uint, caseType string, caseIDs []string) (deletedCount int, failedCaseIDs []string, err error)

	// 新增：重新分配所有ID
	ReassignAllIDs(projectID uint, userID uint, caseType string) error

	// 新增：版本管理方法
	BatchSaveVersion(projectID uint, userID uint) (*VersionInfoDTO, error)
	GetVersionList(projectID uint, userID uint, page int, size int) (*VersionListDTO, error)
	DownloadVersion(projectID uint, userID uint, versionID string) ([]byte, string, error)
	DeleteVersion(projectID uint, userID uint, versionID string) error
	UpdateVersionRemark(projectID uint, userID uint, versionID string, remark string) error
}

type autoTestCaseService struct {
	repo           repositories.AutoTestCaseRepository
	projectService ProjectService
	db             *gorm.DB
}

// NewAutoTestCaseService 创建服务实例
func NewAutoTestCaseService(repo repositories.AutoTestCaseRepository, projectService ProjectService, db *gorm.DB) AutoTestCaseService {
	return &autoTestCaseService{
		repo:           repo,
		projectService: projectService,
		db:             db,
	}
}

// GetMetadata 获取元数据
func (s *autoTestCaseService) GetMetadata(projectID uint, userID uint, caseType string) (*AutoMetadataDTO, error) {
	// 验证用户是否是项目成员
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 获取元数据
	testCase, err := s.repo.GetMetadataByProjectID(projectID, caseType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &AutoMetadataDTO{}, nil
		}
		return nil, err
	}

	return &AutoMetadataDTO{
		ScreenCN: testCase.ScreenCN,
		ScreenJP: testCase.ScreenJP,
		ScreenEN: testCase.ScreenEN,
	}, nil
}

// UpdateMetadata 更新元数据
func (s *autoTestCaseService) UpdateMetadata(projectID uint, userID uint, caseType string, req UpdateAutoMetadataRequest) error {
	// 验证用户权限
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 更新元数据
	metadata := map[string]interface{}{
		"screen_cn": req.ScreenCN,
		"screen_jp": req.ScreenJP,
		"screen_en": req.ScreenEN,
	}

	err = s.repo.UpdateMetadata(projectID, caseType, metadata)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未找到元数据记录")
		}
		return err
	}

	return nil
}

// GetCases 获取用例列表
func (s *autoTestCaseService) GetCases(projectID uint, userID uint, caseType string, page int, size int) (*AutoCaseListDTO, error) {
	// 验证用户权限
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 参数校验
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 50
	}

	offset := (page - 1) * size
	cases, total, err := s.repo.GetCasesByType(projectID, caseType, offset, size)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	caseDTOs := make([]*AutoCaseDTO, 0, len(cases))
	for _, c := range cases {
		caseDTOs = append(caseDTOs, &AutoCaseDTO{
			CaseID:           c.CaseID,
			ID:               c.ID,
			CaseNum:          c.CaseNumber,
			ScreenCN:         c.ScreenCN,
			ScreenJP:         c.ScreenJP,
			ScreenEN:         c.ScreenEN,
			FunctionCN:       c.FunctionCN,
			FunctionJP:       c.FunctionJP,
			FunctionEN:       c.FunctionEN,
			PreconditionCN:   c.PreconditionCN,
			PreconditionJP:   c.PreconditionJP,
			PreconditionEN:   c.PreconditionEN,
			TestStepsCN:      c.TestStepsCN,
			TestStepsJP:      c.TestStepsJP,
			TestStepsEN:      c.TestStepsEN,
			ExpectedResultCN: c.ExpectedResultCN,
			ExpectedResultJP: c.ExpectedResultJP,
			ExpectedResultEN: c.ExpectedResultEN,
			TestResult:       c.TestResult,
			Remark:           c.Remark,
		})
	}

	return &AutoCaseListDTO{
		Cases: caseDTOs,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// CreateCase 创建新用例
func (s *autoTestCaseService) CreateCase(projectID uint, userID uint, req CreateAutoCaseRequest) (*AutoCaseDTO, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 设置默认测试结果
	if req.TestResult == "" {
		req.TestResult = "NR"
	}

	// 获取最大ID
	maxID, err := s.repo.GetMaxIDByProjectAndType(projectID, req.CaseType)
	if err != nil {
		return nil, fmt.Errorf("get max id: %w", err)
	}

	testCase := &models.AutoTestCase{
		ID:               maxID + 1,
		ProjectID:        projectID,
		CaseType:         req.CaseType,
		CaseNumber:       req.CaseNum,
		ScreenCN:         req.ScreenCN,
		ScreenJP:         req.ScreenJP,
		ScreenEN:         req.ScreenEN,
		FunctionCN:       req.FunctionCN,
		FunctionJP:       req.FunctionJP,
		FunctionEN:       req.FunctionEN,
		PreconditionCN:   req.PreconditionCN,
		PreconditionJP:   req.PreconditionJP,
		PreconditionEN:   req.PreconditionEN,
		TestStepsCN:      req.TestStepsCN,
		TestStepsJP:      req.TestStepsJP,
		TestStepsEN:      req.TestStepsEN,
		ExpectedResultCN: req.ExpectedResultCN,
		ExpectedResultJP: req.ExpectedResultJP,
		ExpectedResultEN: req.ExpectedResultEN,
		TestResult:       req.TestResult,
		Remark:           req.Remark,
	}

	if err := s.repo.Create(testCase); err != nil {
		return nil, fmt.Errorf("create auto test case: %w", err)
	}

	// 调试日志：检查创建后的CaseID
	log.Printf("[CreateCase] Generated CaseID: %s, ID: %d", testCase.CaseID, testCase.ID)

	return &AutoCaseDTO{
		CaseID:           testCase.CaseID,
		ID:               testCase.ID,
		CaseNum:          testCase.CaseNumber,
		ScreenCN:         testCase.ScreenCN,
		ScreenJP:         testCase.ScreenJP,
		ScreenEN:         testCase.ScreenEN,
		FunctionCN:       testCase.FunctionCN,
		FunctionJP:       testCase.FunctionJP,
		FunctionEN:       testCase.FunctionEN,
		PreconditionCN:   testCase.PreconditionCN,
		PreconditionJP:   testCase.PreconditionJP,
		PreconditionEN:   testCase.PreconditionEN,
		TestStepsCN:      testCase.TestStepsCN,
		TestStepsJP:      testCase.TestStepsJP,
		TestStepsEN:      testCase.TestStepsEN,
		ExpectedResultCN: testCase.ExpectedResultCN,
		ExpectedResultJP: testCase.ExpectedResultJP,
		ExpectedResultEN: testCase.ExpectedResultEN,
		TestResult:       testCase.TestResult,
		Remark:           testCase.Remark,
	}, nil
}

// UpdateCase 更新用例
func (s *autoTestCaseService) UpdateCase(projectID uint, userID uint, caseID string, req UpdateAutoCaseRequest) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 查询用例是否存在且属于当前项目
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

	// 构建updates map
	updates := make(map[string]interface{})
	if req.CaseNum != nil {
		updates["case_number"] = *req.CaseNum // 数据库字段是case_number，不是case_num
	}
	if req.ScreenCN != nil {
		updates["screen_cn"] = *req.ScreenCN
	}
	if req.ScreenJP != nil {
		updates["screen_jp"] = *req.ScreenJP
	}
	if req.ScreenEN != nil {
		updates["screen_en"] = *req.ScreenEN
	}
	if req.FunctionCN != nil {
		updates["function_cn"] = *req.FunctionCN
	}
	if req.FunctionJP != nil {
		updates["function_jp"] = *req.FunctionJP
	}
	if req.FunctionEN != nil {
		updates["function_en"] = *req.FunctionEN
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
	if req.TestResult != nil {
		updates["test_result"] = *req.TestResult
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	if len(updates) == 0 {
		return nil
	}

	return s.repo.UpdateByCaseID(caseID, updates)
}

// DeleteCase 删除用例
func (s *autoTestCaseService) DeleteCase(projectID uint, userID uint, caseID string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 查询用例是否存在且属于当前项目
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

	// 删除用例
	if err := s.repo.DeleteByCaseID(caseID); err != nil {
		return fmt.Errorf("delete test case: %w", err)
	}

	// 删除后自动调整后续用例的display_order
	if err := s.repo.DecrementOrderAfter(projectID, testCase.CaseType, int(testCase.ID)); err != nil {
		return fmt.Errorf("decrement order after deletion: %w", err)
	}

	// 重新分配所有用例的id字段
	if err := s.repo.ReassignDisplayIDs(projectID, testCase.CaseType); err != nil {
		return fmt.Errorf("reassign display ids: %w", err)
	}

	return nil
}

// ReorderAllCases 按现有ID顺序重新编号所有用例
func (s *autoTestCaseService) ReorderAllCases(projectID uint, userID uint, caseType string) (int, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return 0, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return 0, errors.New("无项目访问权限")
	}

	// 获取所有用例(按ID升序)
	cases, err := s.repo.GetByProjectAndTypeOrdered(projectID, caseType)
	if err != nil {
		return 0, fmt.Errorf("get cases ordered: %w", err)
	}

	if len(cases) == 0 {
		return 0, nil
	}

	// 提取case_id顺序
	caseIDOrder := make([]string, len(cases))
	for i, c := range cases {
		caseIDOrder[i] = c.CaseID
	}

	// 批量更新ID
	if err := s.repo.BatchUpdateIDsByCaseID(caseIDOrder); err != nil {
		return 0, fmt.Errorf("batch update ids: %w", err)
	}

	return len(cases), nil
}

// ReorderByIDs 按指定的case_id顺序重新编号用例
func (s *autoTestCaseService) ReorderByIDs(projectID uint, userID uint, caseType string, caseIDs []string) (int, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return 0, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return 0, errors.New("无项目访问权限")
	}

	if len(caseIDs) == 0 {
		return 0, nil
	}

	// 批量更新ID（按caseIDs的顺序重新分配ID: 1, 2, 3...）
	if err := s.repo.BatchUpdateIDsByCaseID(caseIDs); err != nil {
		return 0, fmt.Errorf("batch update ids: %w", err)
	}

	return len(caseIDs), nil
}

// InsertCase 在指定位置插入新用例
func (s *autoTestCaseService) InsertCase(projectID uint, userID uint, caseType string, position string, targetCaseID string) (*models.AutoTestCase, error) {
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

	// 计算新用例的display_order
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

	// 创建新用例
	newCase := &models.AutoTestCase{
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
func (s *autoTestCaseService) BatchDeleteCases(projectID uint, userID uint, caseType string, caseIDs []string) (int, []string, error) {
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
func (s *autoTestCaseService) ReassignAllIDs(projectID uint, userID uint, caseType string) error {
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

// BatchSaveVersion 批量保存版本(ROLE1-4所有用例)
func (s *autoTestCaseService) BatchSaveVersion(projectID uint, userID uint) (*VersionInfoDTO, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 获取项目信息
	project, _, err := s.projectService.GetByID(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	// 并发查询4个ROLE的用例数据
	roleTypes := []string{"role1", "role2", "role3", "role4"}
	type roleData struct {
		roleType string
		cases    []*models.AutoTestCase
		err      error
	}
	resultChan := make(chan roleData, 4)
	var wg sync.WaitGroup

	for _, rt := range roleTypes {
		wg.Add(1)
		go func(roleType string) {
			defer wg.Done()
			cases, err := s.repo.GetByProjectAndType(projectID, roleType)
			resultChan <- roleData{roleType: roleType, cases: cases, err: err}
		}(rt)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	roleDataMap := make(map[string][]*models.AutoTestCase)
	for rd := range resultChan {
		if rd.err != nil {
			return nil, fmt.Errorf("get %s cases: %w", rd.roleType, rd.err)
		}
		roleDataMap[rd.roleType] = rd.cases
	}

	// 生成version_id(格式: {项目名}_{YYYYMMDD_HHMMSS})
	timestamp := time.Now().Format("20060102_150405")
	versionID := fmt.Sprintf("%s_%s", project.Name, timestamp)
	baseDir := filepath.Join("storage", "versions", "auto-cases", fmt.Sprintf("%d", projectID), versionID)

	// 确保目录存在
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}

	// 生成4个Excel文件并保存版本记录
	var files []*VersionFileInfo
	totalCases := 0
	excelService := &excelService{}

	for _, roleType := range roleTypes {
		cases := roleDataMap[roleType]
		// 文件名使用纯英文,避免zip下载时中文乱码
		filename := fmt.Sprintf("%s_AutoTestCase_%s_%s.xlsx", project.Name, strings.ToUpper(roleType), timestamp)
		filePath := filepath.Join(baseDir, filename)

		// 生成Excel
		// 注意：ExportAutoCasesAllLanguages接受[]*models.AutoTestCase
		if err := excelService.ExportAutoCasesAllLanguages(cases, filePath); err != nil {
			return nil, fmt.Errorf("export %s excel: %w", roleType, err)
		}

		// 获取文件大小
		fileInfo, _ := os.Stat(filePath)
		var fileSize int64
		if fileInfo != nil {
			fileSize = fileInfo.Size()
		}

		// 插入版本记录
		version := &models.AutoTestCaseVersion{
			VersionID:   versionID,
			ProjectID:   projectID,
			ProjectName: project.Name,
			RoleType:    roleType,
			Filename:    filename,
			FilePath:    filePath,
			FileSize:    fileSize,
			CaseCount:   len(cases),
			CreatedBy:   &userID,
			CreatedAt:   time.Now(),
		}
		if err := s.db.Create(version).Error; err != nil {
			return nil, fmt.Errorf("save version record: %w", err)
		}

		files = append(files, &VersionFileInfo{
			Role:     roleType,
			Filename: filename,
			Count:    len(cases),
		})
		totalCases += len(cases)
	}

	return &VersionInfoDTO{
		VersionID:  versionID,
		SavedAt:    time.Now(),
		Files:      files,
		TotalCases: totalCases,
	}, nil
}

// GetVersionList 获取版本列表
func (s *autoTestCaseService) GetVersionList(projectID uint, userID uint, page int, size int) (*VersionListDTO, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 查询版本记录(按version_id分组)
	var versions []models.AutoTestCaseVersion
	offset := (page - 1) * size

	if err := s.db.Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&versions).Error; err != nil {
		return nil, fmt.Errorf("query versions: %w", err)
	}

	// 按version_id分组
	versionMap := make(map[string][]*models.AutoTestCaseVersion)
	for i := range versions {
		v := &versions[i]
		versionMap[v.VersionID] = append(versionMap[v.VersionID], v)
	}

	// 转换为DTO并按时间排序
	var versionDTOs []*VersionDTO
	for versionID, records := range versionMap {
		var files []*VersionFileInfo
		totalCases := 0
		remark := ""
		savedAt := records[0].CreatedAt

		for _, record := range records {
			files = append(files, &VersionFileInfo{
				Role:     record.RoleType,
				Filename: record.Filename,
				Count:    record.CaseCount,
			})
			totalCases += record.CaseCount
			if record.Remark != "" {
				remark = record.Remark
			}
		}

		versionDTOs = append(versionDTOs, &VersionDTO{
			VersionID:  versionID,
			SavedAt:    savedAt,
			Files:      files,
			TotalCases: totalCases,
			Remark:     remark,
		})
	}

	// 按创建时间倒序排列(最新的在前)
	sort.Slice(versionDTOs, func(i, j int) bool {
		return versionDTOs[i].SavedAt.After(versionDTOs[j].SavedAt)
	})

	// 分页处理
	total := int64(len(versionDTOs))
	start := offset
	end := offset + size
	if start > len(versionDTOs) {
		start = len(versionDTOs)
	}
	if end > len(versionDTOs) {
		end = len(versionDTOs)
	}

	return &VersionListDTO{
		Versions: versionDTOs[start:end],
		Total:    total,
		Page:     page,
		Size:     size,
	}, nil
}

// DownloadVersion 下载版本(zip打包4个Excel)
func (s *autoTestCaseService) DownloadVersion(projectID uint, userID uint, versionID string) ([]byte, string, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, "", fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, "", errors.New("无项目访问权限")
	}

	// 查询版本记录
	var versions []models.AutoTestCaseVersion
	if err := s.db.Where("project_id = ? AND version_id = ?", projectID, versionID).
		Find(&versions).Error; err != nil {
		return nil, "", fmt.Errorf("query version: %w", err)
	}
	if len(versions) == 0 {
		return nil, "", errors.New("版本不存在")
	}

	// 创建zip buffer
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	defer zipWriter.Close()

	// 添加4个Excel文件到zip
	for _, v := range versions {
		// 验证文件路径安全性
		if err := validateFilePath(v.FilePath); err != nil {
			return nil, "", fmt.Errorf("invalid file path: %w", err)
		}

		// 读取文件
		data, err := os.ReadFile(v.FilePath)
		if err != nil {
			return nil, "", fmt.Errorf("read file %s: %w", v.Filename, err)
		}

		// 写入zip
		w, err := zipWriter.Create(v.Filename)
		if err != nil {
			return nil, "", fmt.Errorf("create zip entry: %w", err)
		}
		if _, err := w.Write(data); err != nil {
			return nil, "", fmt.Errorf("write zip entry: %w", err)
		}
	}

	zipWriter.Close()

	// 生成zip文件名(纯英文,避免下载时中文乱码)
	zipFilename := fmt.Sprintf("%s_AutoTestCase_Version_%s.zip", versions[0].ProjectName, versionID)

	return buf.Bytes(), zipFilename, nil
}

// DeleteVersion 删除版本
func (s *autoTestCaseService) DeleteVersion(projectID uint, userID uint, versionID string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 查询版本记录
	var versions []models.AutoTestCaseVersion
	if err := s.db.Where("project_id = ? AND version_id = ?", projectID, versionID).
		Find(&versions).Error; err != nil {
		return fmt.Errorf("query version: %w", err)
	}
	if len(versions) == 0 {
		return errors.New("版本不存在")
	}

	// 删除物理文件
	for _, v := range versions {
		if err := os.Remove(v.FilePath); err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to delete file %s: %v", v.FilePath, err)
		}
	}

	// 删除数据库记录
	if err := s.db.Where("project_id = ? AND version_id = ?", projectID, versionID).
		Delete(&models.AutoTestCaseVersion{}).Error; err != nil {
		return fmt.Errorf("delete version records: %w", err)
	}

	return nil
}

// UpdateVersionRemark 更新版本备注
func (s *autoTestCaseService) UpdateVersionRemark(projectID uint, userID uint, versionID string, remark string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 更新所有对应记录的备注
	result := s.db.Model(&models.AutoTestCaseVersion{}).
		Where("project_id = ? AND version_id = ?", projectID, versionID).
		Update("remark", remark)

	if result.Error != nil {
		return fmt.Errorf("update remark: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("版本不存在")
	}

	return nil
}

// validateFilePath 验证文件路径安全性(防止路径遍历攻击)
func validateFilePath(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	allowedPrefix := filepath.Join(cwd, "storage", "versions", "auto-cases")
	if !strings.HasPrefix(absPath, allowedPrefix) {
		return errors.New("invalid file path")
	}
	return nil
}
