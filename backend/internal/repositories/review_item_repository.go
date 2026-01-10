package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// ReviewItemRepository 审阅条目仓储接口
type ReviewItemRepository interface {
	// Create 创建审阅条目
	Create(item *models.CaseReviewItem) error

	// GetByID 根据ID获取审阅条目
	GetByID(id uint) (*models.CaseReviewItem, error)

	// GetByProjectID 获取项目所有审阅条目
	GetByProjectID(projectID uint) ([]models.CaseReviewItem, error)

	// Update 更新审阅条目
	Update(item *models.CaseReviewItem) error

	// Delete 删除审阅条目
	Delete(id uint) error

	// GetByProjectAndName 根据项目ID和名称获取审阅条目(用于唯一性校验)
	GetByProjectAndName(projectID uint, name string) (*models.CaseReviewItem, error)
}

// reviewItemRepository 审阅条目仓储实现
type reviewItemRepository struct {
	db *gorm.DB
}

// NewReviewItemRepository 创建审阅条目仓储实例
func NewReviewItemRepository(db *gorm.DB) ReviewItemRepository {
	return &reviewItemRepository{db: db}
}

// Create 实现创建审阅条目
func (r *reviewItemRepository) Create(item *models.CaseReviewItem) error {
	return r.db.Create(item).Error
}

// GetByID 实现根据ID获取审阅条目
func (r *reviewItemRepository) GetByID(id uint) (*models.CaseReviewItem, error) {
	var item models.CaseReviewItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetByProjectID 实现获取项目所有审阅条目
func (r *reviewItemRepository) GetByProjectID(projectID uint) ([]models.CaseReviewItem, error) {
	var items []models.CaseReviewItem
	err := r.db.Where("project_id = ?", projectID).
		Order("updated_at DESC").
		Find(&items).Error
	return items, err
}

// Update 实现更新审阅条目
func (r *reviewItemRepository) Update(item *models.CaseReviewItem) error {
	return r.db.Save(item).Error
}

// Delete 实现删除审阅条目(硬删除)
func (r *reviewItemRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&models.CaseReviewItem{}, id).Error
}

// GetByProjectAndName 实现根据项目ID和名称获取审阅条目
func (r *reviewItemRepository) GetByProjectAndName(projectID uint, name string) (*models.CaseReviewItem, error) {
	var item models.CaseReviewItem
	err := r.db.Where("project_id = ? AND name = ?", projectID, name).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
