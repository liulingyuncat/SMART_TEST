package repositories

import (
	"fmt"

	"webtest/internal/models"

	"gorm.io/gorm"
)

// DefectCommentRepository 缺陷说明仓库接口
type DefectCommentRepository interface {
	// CRUD方法
	Create(comment *models.DefectComment) error
	GetByID(id uint) (*models.DefectComment, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error

	// 列表查询
	ListByDefectID(defectID string) ([]*models.DefectComment, error)

	// 权限检查
	IsCreatedBy(id uint, userID uint) (bool, error)
}

// defectCommentRepository 缺陷说明仓库实现
type defectCommentRepository struct {
	db *gorm.DB
}

// NewDefectCommentRepository 创建缺陷说明仓库实例
func NewDefectCommentRepository(db *gorm.DB) DefectCommentRepository {
	return &defectCommentRepository{db: db}
}

// Create 创建说明
func (r *defectCommentRepository) Create(comment *models.DefectComment) error {
	if err := r.db.Create(comment).Error; err != nil {
		return fmt.Errorf("create defect comment: %w", err)
	}
	return nil
}

// GetByID 根据ID获取说明
func (r *defectCommentRepository) GetByID(id uint) (*models.DefectComment, error) {
	var comment models.DefectComment
	if err := r.db.First(&comment, id).Error; err != nil {
		return nil, fmt.Errorf("get defect comment by id: %w", err)
	}
	return &comment, nil
}

// Update 更新说明
func (r *defectCommentRepository) Update(id uint, updates map[string]interface{}) error {
	if err := r.db.Model(&models.DefectComment{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("update defect comment: %w", err)
	}
	return nil
}

// Delete 软删除说明
func (r *defectCommentRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.DefectComment{}, id).Error; err != nil {
		return fmt.Errorf("delete defect comment: %w", err)
	}
	return nil
}

// ListByDefectID 按缺陷ID查询说明列表（升序，预加载用户信息）
func (r *defectCommentRepository) ListByDefectID(defectID string) ([]*models.DefectComment, error) {
	var comments []*models.DefectComment
	if err := r.db.Where("defect_id = ?", defectID).
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("list defect comments by defect id: %w", err)
	}
	return comments, nil
}

// IsCreatedBy 检查说明是否由指定用户创建
func (r *defectCommentRepository) IsCreatedBy(id uint, userID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.DefectComment{}).
		Where("id = ? AND created_by = ?", id, userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check comment creator: %w", err)
	}
	return count > 0, nil
}
