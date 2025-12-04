package repositories

import (
	"fmt"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// RequirementRepository 需求文档仓库接口
type RequirementRepository interface {
	GetByProjectID(projectID uint) (*models.Requirement, error)
	CreateOrUpdate(requirement *models.Requirement) error
	GetFieldByType(projectID uint, docType string) (string, error)
	UpdateFieldByType(projectID uint, docType string, content string) error
}

// requirementRepository 需求文档仓库实现
type requirementRepository struct {
	db *gorm.DB
}

// NewRequirementRepository 创建需求文档仓库实例
func NewRequirementRepository(db *gorm.DB) RequirementRepository {
	return &requirementRepository{db: db}
}

// GetByProjectID 根据项目ID查询需求文档
func (r *requirementRepository) GetByProjectID(projectID uint) (*models.Requirement, error) {
	var requirement models.Requirement
	err := r.db.Where("project_id = ?", projectID).First(&requirement).Error
	if err == gorm.ErrRecordNotFound {
		// 返回空记录而非错误,允许创建新文档
		return &models.Requirement{ProjectID: projectID}, nil
	}
	return &requirement, err
}

// CreateOrUpdate 创建或更新需求文档
func (r *requirementRepository) CreateOrUpdate(requirement *models.Requirement) error {
	var existing models.Requirement
	err := r.db.Where("project_id = ?", requirement.ProjectID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		return r.db.Create(requirement).Error
	}
	if err != nil {
		return fmt.Errorf("查询现有记录失败: %w", err)
	}

	// 更新现有记录
	requirement.ID = existing.ID
	return r.db.Save(requirement).Error
}

// GetFieldByType 根据文档类型获取对应字段内容
func (r *requirementRepository) GetFieldByType(projectID uint, docType string) (string, error) {
	var requirement models.Requirement
	err := r.db.Where("project_id = ?", projectID).First(&requirement).Error
	if err == gorm.ErrRecordNotFound {
		// 文档不存在,返回空字符串
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("查询需求文档失败: %w", err)
	}

	// 根据 docType 映射到对应字段
	content, err := r.mapDocTypeToField(&requirement, docType)
	if err != nil {
		return "", err
	}
	return content, nil
}

// UpdateFieldByType 根据文档类型更新对应字段内容
func (r *requirementRepository) UpdateFieldByType(projectID uint, docType string, content string) error {
	var requirement models.Requirement
	err := r.db.Where("project_id = ?", projectID).First(&requirement).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		requirement.ProjectID = projectID
		if err := r.setFieldByDocType(&requirement, docType, content); err != nil {
			return err
		}
		return r.db.Create(&requirement).Error
	}
	if err != nil {
		return fmt.Errorf("查询需求文档失败: %w", err)
	}

	// 更新字段
	if err := r.setFieldByDocType(&requirement, docType, content); err != nil {
		return err
	}
	return r.db.Save(&requirement).Error
}

// mapDocTypeToField 将 docType 映射到对应的数据库字段值
func (r *requirementRepository) mapDocTypeToField(req *models.Requirement, docType string) (string, error) {
	switch docType {
	case "overall-requirements":
		return req.OverallRequirements, nil
	case "overall-test-viewpoint":
		return req.OverallTestViewpoint, nil
	case "change-requirements":
		return req.ChangeRequirements, nil
	case "change-test-viewpoint":
		return req.ChangeTestViewpoint, nil
	default:
		return "", fmt.Errorf("无效的文档类型: %s", docType)
	}
}

// setFieldByDocType 根据 docType 设置对应字段的值
func (r *requirementRepository) setFieldByDocType(req *models.Requirement, docType string, content string) error {
	switch docType {
	case "overall-requirements":
		req.OverallRequirements = content
	case "overall-test-viewpoint":
		req.OverallTestViewpoint = content
	case "change-requirements":
		req.ChangeRequirements = content
	case "change-test-viewpoint":
		req.ChangeTestViewpoint = content
	default:
		return fmt.Errorf("无效的文档类型: %s", docType)
	}
	return nil
}
