package repositories

import (
	"fmt"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// DefectPhaseRepository 缺陷测试阶段仓储接口
type DefectPhaseRepository interface {
	Create(phase *models.DefectPhase) error
	GetByID(id uint) (*models.DefectPhase, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	ListByProjectID(projectID uint) ([]*models.DefectPhase, error)
	ExistsByName(projectID uint, name string, excludeID uint) (bool, error)
}

type defectPhaseRepository struct {
	db *gorm.DB
}

// NewDefectPhaseRepository 创建缺陷测试阶段仓储实例
func NewDefectPhaseRepository(db *gorm.DB) DefectPhaseRepository {
	return &defectPhaseRepository{db: db}
}

// Create 创建Phase
func (r *defectPhaseRepository) Create(phase *models.DefectPhase) error {
	err := r.db.Create(phase).Error
	if err != nil {
		return fmt.Errorf("create phase: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Phase
func (r *defectPhaseRepository) GetByID(id uint) (*models.DefectPhase, error) {
	var phase models.DefectPhase
	err := r.db.Where("id = ?", id).First(&phase).Error
	if err != nil {
		return nil, err // 保留gorm.ErrRecordNotFound
	}
	return &phase, nil
}

// Update 更新Phase
func (r *defectPhaseRepository) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&models.DefectPhase{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("update phase %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete 软删除Phase
func (r *defectPhaseRepository) Delete(id uint) error {
	result := r.db.Where("id = ?", id).Delete(&models.DefectPhase{})
	if result.Error != nil {
		return fmt.Errorf("delete phase %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ListByProjectID 根据项目ID获取Phase列表
func (r *defectPhaseRepository) ListByProjectID(projectID uint) ([]*models.DefectPhase, error) {
	var phases []*models.DefectPhase
	err := r.db.Where("project_id = ?", projectID).
		Order("sort_order ASC, id ASC").
		Find(&phases).Error

	if err != nil {
		return nil, fmt.Errorf("list phases by project: %w", err)
	}

	return phases, nil
}

// ExistsByName 检查同项目下是否存在同名Phase（排除指定ID）
func (r *defectPhaseRepository) ExistsByName(projectID uint, name string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.DefectPhase{}).
		Where("project_id = ? AND name = ?", projectID, name)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check phase exists: %w", err)
	}

	return count > 0, nil
}
