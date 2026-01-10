package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// RequirementItemRepository 需求条目仓库接口
type RequirementItemRepository interface {
	Create(item *models.RequirementItem) error
	Update(item *models.RequirementItem) error
	Delete(id uint) error
	FindByID(id uint) (*models.RequirementItem, error)
	FindByProjectID(projectID uint) ([]*models.RequirementItem, error)
	FindByProjectIDAndName(projectID uint, name string) (*models.RequirementItem, error)
	BulkCreate(items []*models.RequirementItem) error
	BulkUpdate(items []*models.RequirementItem) error
	BulkDelete(ids []uint) error
}

// requirementItemRepository 需求条目仓库实现
type requirementItemRepository struct {
	db *gorm.DB
}

// NewRequirementItemRepository 创建需求条目仓库实例
func NewRequirementItemRepository(db *gorm.DB) RequirementItemRepository {
	return &requirementItemRepository{db: db}
}

// Create 创建需求条目
func (r *requirementItemRepository) Create(item *models.RequirementItem) error {
	return r.db.Create(item).Error
}

// Update 更新需求条目
func (r *requirementItemRepository) Update(item *models.RequirementItem) error {
	return r.db.Save(item).Error
}

// Delete 删除需求条目(硬删除)
func (r *requirementItemRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&models.RequirementItem{}, id).Error
}

// FindByID 根据ID查询需求条目
func (r *requirementItemRepository) FindByID(id uint) (*models.RequirementItem, error) {
	var item models.RequirementItem
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// FindByProjectID 根据项目ID查询所有需求条目
func (r *requirementItemRepository) FindByProjectID(projectID uint) ([]*models.RequirementItem, error) {
	var items []*models.RequirementItem
	if err := r.db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByProjectIDAndName 根据项目ID和名称查询需求条目
func (r *requirementItemRepository) FindByProjectIDAndName(projectID uint, name string) (*models.RequirementItem, error) {
	var item models.RequirementItem
	if err := r.db.Where("project_id = ? AND name = ?", projectID, name).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// BulkCreate 批量创建需求条目
func (r *requirementItemRepository) BulkCreate(items []*models.RequirementItem) error {
	return r.db.Create(items).Error
}

// BulkUpdate 批量更新需求条目
func (r *requirementItemRepository) BulkUpdate(items []*models.RequirementItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Save(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BulkDelete 批量删除需求条目(硬删除)
func (r *requirementItemRepository) BulkDelete(ids []uint) error {
	return r.db.Unscoped().Delete(&models.RequirementItem{}, ids).Error
}
