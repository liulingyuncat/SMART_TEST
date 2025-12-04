package services

import (
	"fmt"
	"os"
	"path/filepath"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"gorm.io/gorm"
)

// VersionService 版本管理服务接口
type VersionService interface {
	SaveVersion(projectID, userID uint, caseType string) (filename string, err error)
	GetVersionList(projectID uint, caseType string) ([]*models.CaseVersion, error)
	GetVersionByID(versionID uint) (*models.CaseVersion, error)
	DownloadVersion(projectID, versionID uint) (fileBytes []byte, filename string, err error)
	DeleteVersion(projectID, versionID uint) error
	CreateVersion(version *models.CaseVersion) error
	UpdateVersionRemark(projectID, versionID uint, remark string) error
}

type versionService struct {
	db          *gorm.DB
	versionRepo repositories.CaseVersionRepository
	excelSvc    ExcelService
}

// NewVersionService 创建版本管理服务实例
func NewVersionService(db *gorm.DB, versionRepo repositories.CaseVersionRepository, excelSvc ExcelService) VersionService {
	return &versionService{
		db:          db,
		versionRepo: versionRepo,
		excelSvc:    excelSvc,
	}
}

// SaveVersion 保存用例版本(支持overall和change类型)
func (s *versionService) SaveVersion(projectID, userID uint, caseType string) (string, error) {
	// 1. 参数验证
	if caseType != "overall" && caseType != "change" && caseType != "acceptance" {
		return "", fmt.Errorf("unsupported case_type: %s", caseType)
	}

	// 2. 调用ExcelService导出对应类型用例(版本管理不合并执行结果,传空taskUUID)
	fileBytes, filename, err := s.excelSvc.ExportCases(projectID, caseType, "")
	if err != nil {
		return "", fmt.Errorf("export cases failed: %w", err)
	}

	// 3. 构造文件存储路径
	storageDir := filepath.Join("storage", "versions", fmt.Sprintf("%d", projectID))
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", fmt.Errorf("create storage dir: %w", err)
	}

	filePath := filepath.Join(storageDir, filename)

	// 4. 写入文件
	if err := os.WriteFile(filePath, fileBytes, 0644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	// 5. 插入版本记录
	version := &models.CaseVersion{
		ProjectID: projectID,
		DocType:   caseType,
		Filename:  filename,
		FilePath:  filePath,
		FileSize:  int64(len(fileBytes)),
		CreatedBy: &userID,
	}

	if err := s.versionRepo.Create(version); err != nil {
		// 删除已创建的文件
		os.Remove(filePath)
		return "", fmt.Errorf("create version record: %w", err)
	}

	return filename, nil
}

// GetVersionList 获取版本列表(支持按类型过滤)
func (s *versionService) GetVersionList(projectID uint, caseType string) ([]*models.CaseVersion, error) {
	// 如果caseType为空,返回所有版本(向后兼容)
	if caseType == "" {
		return s.versionRepo.GetByProjectID(projectID)
	}
	// 否则按类型过滤查询
	return s.versionRepo.GetByProjectIDAndType(projectID, caseType)
}

// GetVersionByID 根据ID获取版本
func (s *versionService) GetVersionByID(versionID uint) (*models.CaseVersion, error) {
	return s.versionRepo.GetByID(versionID)
}

// DownloadVersion 下载版本文件
func (s *versionService) DownloadVersion(projectID, versionID uint) ([]byte, string, error) {
	// 1. 查询版本记录
	version, err := s.versionRepo.GetByID(versionID)
	if err != nil {
		return nil, "", fmt.Errorf("get version: %w", err)
	}
	if version == nil {
		return nil, "", fmt.Errorf("version not found")
	}

	// 2. 验证项目权限
	if version.ProjectID != projectID {
		return nil, "", fmt.Errorf("unauthorized access")
	}

	// 3. 读取文件
	fileBytes, err := os.ReadFile(version.FilePath)
	if err != nil {
		return nil, "", fmt.Errorf("read file: %w", err)
	}

	return fileBytes, version.Filename, nil
}

// DeleteVersion 删除版本
func (s *versionService) DeleteVersion(projectID, versionID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询版本记录
		version, err := s.versionRepo.GetByID(versionID)
		if err != nil {
			return fmt.Errorf("get version: %w", err)
		}
		if version == nil {
			return fmt.Errorf("version not found")
		}

		// 2. 验证项目权限
		if version.ProjectID != projectID {
			return fmt.Errorf("unauthorized access")
		}

		// 3. 删除文件
		if err := os.Remove(version.FilePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove file: %w", err)
		}

		// 4. 删除数据库记录
		if err := s.versionRepo.Delete(versionID); err != nil {
			return fmt.Errorf("delete record: %w", err)
		}

		return nil
	})
}

// CreateVersion 创建版本记录
func (s *versionService) CreateVersion(version *models.CaseVersion) error {
	return s.versionRepo.Create(version)
}

// UpdateVersionRemark 更新版本备注
func (s *versionService) UpdateVersionRemark(projectID, versionID uint, remark string) error {
	// 1. 获取版本记录
	version, err := s.versionRepo.GetByID(versionID)
	if err != nil {
		return fmt.Errorf("get version: %w", err)
	}
	if version == nil {
		return fmt.Errorf("version not found")
	}

	// 2. 验证项目权限
	if version.ProjectID != projectID {
		return fmt.Errorf("unauthorized access")
	}

	// 3. 更新备注
	version.Remark = remark
	return s.db.Save(version).Error
}
