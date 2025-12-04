package repositories

import (
	"fmt"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// DefectSubjectRepository 缺陷主题分类仓储接口
type DefectSubjectRepository interface {
	Create(subject *models.DefectSubject) error
	GetByID(id uint) (*models.DefectSubject, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	ListByProjectID(projectID uint) ([]*models.DefectSubject, error)
	ExistsByName(projectID uint, name string, excludeID uint) (bool, error)
}

type defectSubjectRepository struct {
	db *gorm.DB
}

// NewDefectSubjectRepository 创建缺陷主题分类仓储实例
func NewDefectSubjectRepository(db *gorm.DB) DefectSubjectRepository {
	return &defectSubjectRepository{db: db}
}

// Create 创建Subject
func (r *defectSubjectRepository) Create(subject *models.DefectSubject) error {
	err := r.db.Create(subject).Error
	if err != nil {
		return fmt.Errorf("create subject: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Subject
func (r *defectSubjectRepository) GetByID(id uint) (*models.DefectSubject, error) {
	var subject models.DefectSubject
	err := r.db.Where("id = ?", id).First(&subject).Error
	if err != nil {
		return nil, err // 保留gorm.ErrRecordNotFound
	}
	return &subject, nil
}

// Update 更新Subject
func (r *defectSubjectRepository) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&models.DefectSubject{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("update subject %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete 软删除Subject
func (r *defectSubjectRepository) Delete(id uint) error {
	result := r.db.Where("id = ?", id).Delete(&models.DefectSubject{})
	if result.Error != nil {
		return fmt.Errorf("delete subject %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ListByProjectID 根据项目ID获取Subject列表
func (r *defectSubjectRepository) ListByProjectID(projectID uint) ([]*models.DefectSubject, error) {
	var subjects []*models.DefectSubject
	err := r.db.Where("project_id = ?", projectID).
		Order("sort_order ASC, id ASC").
		Find(&subjects).Error

	if err != nil {
		return nil, fmt.Errorf("list subjects by project: %w", err)
	}

	return subjects, nil
}

// ExistsByName 检查同项目下是否存在同名Subject（排除指定ID）
func (r *defectSubjectRepository) ExistsByName(projectID uint, name string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.DefectSubject{}).
		Where("project_id = ? AND name = ?", projectID, name)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check subject exists: %w", err)
	}

	return count > 0, nil
}
