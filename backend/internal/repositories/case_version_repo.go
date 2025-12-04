package repositories

import (
	"fmt"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// CaseVersionRepository 版本记录仓储接口
type CaseVersionRepository interface {
	Create(version *models.CaseVersion) error
	GetByProjectID(projectID uint) ([]*models.CaseVersion, error)
	GetByProjectIDAndType(projectID uint, caseType string) ([]*models.CaseVersion, error)
	GetByID(id uint) (*models.CaseVersion, error)
	Delete(id uint) error
}

type caseVersionRepository struct {
	db *gorm.DB
}

// NewCaseVersionRepository 创建版本记录仓储实例
func NewCaseVersionRepository(db *gorm.DB) CaseVersionRepository {
	return &caseVersionRepository{db: db}
}

// Create 创建版本记录
func (r *caseVersionRepository) Create(version *models.CaseVersion) error {
	if err := r.db.Create(version).Error; err != nil {
		return fmt.Errorf("failed to create version: %w", err)
	}
	return nil
}

// GetByProjectID 查询项目的所有版本(按创建时间倒序)
func (r *caseVersionRepository) GetByProjectID(projectID uint) ([]*models.CaseVersion, error) {
	var versions []*models.CaseVersion
	err := r.db.Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&versions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}
	return versions, nil
}

// GetByProjectIDAndType 查询项目指定类型的版本(按创建时间倒序)
func (r *caseVersionRepository) GetByProjectIDAndType(projectID uint, caseType string) ([]*models.CaseVersion, error) {
	var versions []*models.CaseVersion
	err := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).
		Order("created_at DESC").
		Find(&versions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get versions by type: %w", err)
	}
	return versions, nil
}

// GetByID 根据ID查询版本记录
func (r *caseVersionRepository) GetByID(id uint) (*models.CaseVersion, error) {
	var version models.CaseVersion
	err := r.db.First(&version, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	return &version, nil
}

// Delete 删除版本记录
func (r *caseVersionRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.CaseVersion{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete version: %w", err)
	}
	return nil
}
