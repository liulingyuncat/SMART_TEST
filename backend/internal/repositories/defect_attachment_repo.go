package repositories

import (
	"fmt"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// DefectAttachmentRepository 缺陷附件仓储接口
type DefectAttachmentRepository interface {
	Create(attachment *models.DefectAttachment) error
	GetByID(id uint) (*models.DefectAttachment, error)
	Delete(id uint) error
	ListByDefectID(defectID string) ([]*models.DefectAttachment, error)
	DeleteByDefectID(defectID string) error
}

type defectAttachmentRepository struct {
	db *gorm.DB
}

// NewDefectAttachmentRepository 创建缺陷附件仓储实例
func NewDefectAttachmentRepository(db *gorm.DB) DefectAttachmentRepository {
	return &defectAttachmentRepository{db: db}
}

// Create 创建附件记录
func (r *defectAttachmentRepository) Create(attachment *models.DefectAttachment) error {
	err := r.db.Create(attachment).Error
	if err != nil {
		return fmt.Errorf("create attachment: %w", err)
	}
	return nil
}

// GetByID 根据ID获取附件
func (r *defectAttachmentRepository) GetByID(id uint) (*models.DefectAttachment, error) {
	var attachment models.DefectAttachment
	err := r.db.Where("id = ?", id).First(&attachment).Error
	if err != nil {
		return nil, err // 保留gorm.ErrRecordNotFound
	}
	return &attachment, nil
}

// Delete 软删除附件记录
func (r *defectAttachmentRepository) Delete(id uint) error {
	result := r.db.Where("id = ?", id).Delete(&models.DefectAttachment{})
	if result.Error != nil {
		return fmt.Errorf("delete attachment %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ListByDefectID 根据缺陷ID获取附件列表
func (r *defectAttachmentRepository) ListByDefectID(defectID string) ([]*models.DefectAttachment, error) {
	var attachments []*models.DefectAttachment
	err := r.db.Where("defect_id = ?", defectID).
		Order("created_at ASC").
		Find(&attachments).Error

	if err != nil {
		return nil, fmt.Errorf("list attachments by defect: %w", err)
	}

	return attachments, nil
}

// DeleteByDefectID 根据缺陷ID删除所有附件（软删除）
func (r *defectAttachmentRepository) DeleteByDefectID(defectID string) error {
	err := r.db.Where("defect_id = ?", defectID).Delete(&models.DefectAttachment{}).Error
	if err != nil {
		return fmt.Errorf("delete attachments by defect %s: %w", defectID, err)
	}
	return nil
}
