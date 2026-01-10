package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"gorm.io/gorm"
)

// WebVersionFileInfo Web版本文件信息
type WebVersionFileInfo struct {
	Language  string `json:"language"`
	Filename  string `json:"filename"`
	CaseCount int    `json:"case_count"`
}

// WebVersionInfo Web版本保存返回信息
type WebVersionInfo struct {
	VersionID   string                `json:"version_id"`
	SavedAt     string                `json:"saved_at"`
	Files       []*WebVersionFileInfo `json:"files"`
	ZipFilename string                `json:"zip_filename"`
	TotalCases  int                   `json:"total_cases"`
}

// WebVersionDTO Web用例版本DTO
type WebVersionDTO struct {
	VersionID   string `json:"version_id"`
	ZipFilename string `json:"zip_filename"`
	FileSize    int64  `json:"file_size"`
	CaseCount   int    `json:"case_count"`
	Remark      string `json:"remark"`
	CreatedAt   string `json:"created_at"`
}

// WebVersionListDTO Web版本列表DTO
type WebVersionListDTO struct {
	Versions []*WebVersionDTO `json:"versions"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Size     int              `json:"size"`
}

// WebVersionService Web用例版本服务接口
type WebVersionService interface {
	SaveVersion(projectID uint, userID uint) (*WebVersionInfo, error)
	GetVersionList(projectID uint, userID uint, page int, size int) (*WebVersionListDTO, error)
	DownloadVersion(projectID uint, userID uint, versionID string) ([]byte, string, error)
	DeleteVersion(projectID uint, userID uint, versionID string) error
	UpdateVersionRemark(projectID uint, userID uint, versionID string, remark string) error
}

type webVersionService struct {
	db            *gorm.DB
	projectRepo   repositories.ProjectRepository
	caseGroupRepo *repositories.CaseGroupRepository
	autoCaseRepo  repositories.AutoTestCaseRepository
	excelService  ExcelService
}

// NewWebVersionService 创建Web版本服务实例
func NewWebVersionService(
	db *gorm.DB,
	projectRepo repositories.ProjectRepository,
	caseGroupRepo *repositories.CaseGroupRepository,
	autoCaseRepo repositories.AutoTestCaseRepository,
	excelService ExcelService,
) WebVersionService {
	return &webVersionService{
		db:            db,
		projectRepo:   projectRepo,
		caseGroupRepo: caseGroupRepo,
		autoCaseRepo:  autoCaseRepo,
		excelService:  excelService,
	}
}

// SaveVersion 保存Web用例版本
func (s *webVersionService) SaveVersion(projectID uint, userID uint) (*WebVersionInfo, error) {
	// 验证项目权限
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	if project == nil {
		return nil, errors.New("项目不存在")
	}

	// TODO: 验证用户是否是项目成员（需要ProjectService）

	// 查询所有Web用例数据
	allCases, err := s.autoCaseRepo.GetByProjectAndType(projectID, "web")
	if err != nil {
		return nil, fmt.Errorf("get web cases: %w", err)
	}

	if len(allCases) == 0 {
		return nil, errors.New("没有可导出的Web用例")
	}

	// 生成版本ID：项目名_yyyyMMdd_HHmmss
	now := time.Now()
	versionID := fmt.Sprintf("%s_%s", project.Name, now.Format("20060102_150405"))
	zipFilename := fmt.Sprintf("%s_AIWeb_TestCase_%s.zip", project.Name, now.Format("20060102_150405"))

	// 创建存储目录
	storageDir := filepath.Join("storage", "versions", "web-cases", fmt.Sprintf("%d", projectID))
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}

	// 将 []*models.AutoTestCase 转换为 []models.AutoTestCase
	casesSlice := make([]models.AutoTestCase, len(allCases))
	for i, c := range allCases {
		casesSlice[i] = *c
	}

	// 调用ExcelService生成4个语言版本的Excel文件并打包为zip
	finalZipPath, fileSize, err := s.excelService.GenerateWebCasesZip(projectID, project.Name, casesSlice)
	if err != nil {
		return nil, fmt.Errorf("generate web cases zip: %w", err)
	}

	// 构造文件信息列表
	files := []*WebVersionFileInfo{
		{Language: "All", Filename: fmt.Sprintf("%s_AIWeb_All_TestCase_%s.xlsx", project.Name, now.Format("20060102_150405")), CaseCount: len(allCases)},
		{Language: "CN", Filename: fmt.Sprintf("%s_AIWeb_CN_TestCase_%s.xlsx", project.Name, now.Format("20060102_150405")), CaseCount: len(allCases)},
		{Language: "JP", Filename: fmt.Sprintf("%s_AIWeb_JP_TestCase_%s.xlsx", project.Name, now.Format("20060102_150405")), CaseCount: len(allCases)},
		{Language: "EN", Filename: fmt.Sprintf("%s_AIWeb_EN_TestCase_%s.xlsx", project.Name, now.Format("20060102_150405")), CaseCount: len(allCases)},
	}

	// 插入版本记录
	version := &models.WebCaseVersion{
		VersionID:   versionID,
		ProjectID:   projectID,
		ProjectName: project.Name,
		ZipFilename: zipFilename,
		ZipPath:     finalZipPath,
		FileSize:    fileSize,
		CaseCount:   len(allCases),
		Remark:      "",
		CreatedBy:   &userID,
		CreatedAt:   now,
	}

	if err := s.db.Create(version).Error; err != nil {
		return nil, fmt.Errorf("create version record: %w", err)
	}

	return &WebVersionInfo{
		VersionID:   versionID,
		SavedAt:     now.Format(time.RFC3339),
		Files:       files,
		ZipFilename: zipFilename,
		TotalCases:  len(allCases),
	}, nil
}

// GetVersionList 获取版本列表
func (s *webVersionService) GetVersionList(projectID uint, userID uint, page int, size int) (*WebVersionListDTO, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	} else if size > 100000 {
		size = 100000
	}

	offset := (page - 1) * size

	var versions []*models.WebCaseVersion
	var total int64

	// 查询条件
	query := s.db.Where("project_id = ?", projectID)

	// 统计总数
	if err := query.Model(&models.WebCaseVersion{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count versions: %w", err)
	}

	// 查询数据（按创建时间降序）
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&versions).Error; err != nil {
		return nil, fmt.Errorf("get versions: %w", err)
	}

	// 转换为DTO
	versionDTOs := make([]*WebVersionDTO, 0, len(versions))
	for _, v := range versions {
		versionDTOs = append(versionDTOs, &WebVersionDTO{
			VersionID:   v.VersionID,
			ZipFilename: v.ZipFilename,
			FileSize:    v.FileSize,
			CaseCount:   v.CaseCount,
			Remark:      v.Remark,
			CreatedAt:   v.CreatedAt.Format(time.RFC3339),
		})
	}

	return &WebVersionListDTO{
		Versions: versionDTOs,
		Total:    total,
		Page:     page,
		Size:     size,
	}, nil
} // DownloadVersion 下载版本zip文件
func (s *webVersionService) DownloadVersion(projectID uint, userID uint, versionID string) ([]byte, string, error) {
	// 查询版本记录
	var version models.WebCaseVersion
	if err := s.db.Where("project_id = ? AND version_id = ?", projectID, versionID).
		First(&version).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("版本不存在")
		}
		return nil, "", fmt.Errorf("get version: %w", err)
	}

	// 读取zip文件
	data, err := os.ReadFile(version.ZipPath)
	if err != nil {
		return nil, "", fmt.Errorf("read zip file: %w", err)
	}

	return data, version.ZipFilename, nil
}

// DeleteVersion 删除版本
func (s *webVersionService) DeleteVersion(projectID uint, userID uint, versionID string) error {
	// 查询版本记录
	var version models.WebCaseVersion
	if err := s.db.Where("project_id = ? AND version_id = ?", projectID, versionID).
		First(&version).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("版本不存在")
		}
		return fmt.Errorf("get version: %w", err)
	}

	// 删除文件
	if err := os.Remove(version.ZipPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove zip file: %w", err)
	}

	// 删除数据库记录(硬删除)
	if err := s.db.Unscoped().Delete(&version).Error; err != nil {
		return fmt.Errorf("delete version record: %w", err)
	}

	return nil
}

// UpdateVersionRemark 更新版本备注
func (s *webVersionService) UpdateVersionRemark(projectID uint, userID uint, versionID string, remark string) error {
	// 查询版本记录
	var version models.WebCaseVersion
	if err := s.db.Where("project_id = ? AND version_id = ?", projectID, versionID).
		First(&version).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("版本不存在")
		}
		return fmt.Errorf("get version: %w", err)
	}

	// 更新备注
	if err := s.db.Model(&version).Update("remark", remark).Error; err != nil {
		return fmt.Errorf("update remark: %w", err)
	}

	return nil
}
