package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// ViewpointItemRepository AI观点条目仓库接口
type ViewpointItemRepository interface {
	Create(item *models.ViewpointItem) error
	Update(item *models.ViewpointItem) error
	Delete(id uint) error
	FindByID(id uint) (*models.ViewpointItem, error)
	FindByProjectID(projectID uint) ([]*models.ViewpointItem, error)
	FindByProjectIDAndName(projectID uint, name string) (*models.ViewpointItem, error)
	BulkCreate(items []*models.ViewpointItem) error
	BulkUpdate(items []*models.ViewpointItem) error
	BulkDelete(ids []uint) error
}

// viewpointItemRepository AI观点条目仓库实现
type viewpointItemRepository struct {
	db *gorm.DB
}

// NewViewpointItemRepository 创建AI观点条目仓库实例
func NewViewpointItemRepository(db *gorm.DB) ViewpointItemRepository {
	return &viewpointItemRepository{db: db}
}

// Create 创建AI观点条目
func (r *viewpointItemRepository) Create(item *models.ViewpointItem) error {
	return r.db.Create(item).Error
}

// Update 更新AI观点条目
func (r *viewpointItemRepository) Update(item *models.ViewpointItem) error {
	return r.db.Save(item).Error
}

// Delete 删除AI观点条目(硬删除)
func (r *viewpointItemRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&models.ViewpointItem{}, id).Error
}

// FindByID 根据ID查询AI观点条目
func (r *viewpointItemRepository) FindByID(id uint) (*models.ViewpointItem, error) {
	var item models.ViewpointItem
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// FindByProjectID 根据项目ID查询所有AI观点条目
func (r *viewpointItemRepository) FindByProjectID(projectID uint) ([]*models.ViewpointItem, error) {
	var items []*models.ViewpointItem
	if err := r.db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByProjectIDAndName 根据项目ID和名称查询AI观点条目
func (r *viewpointItemRepository) FindByProjectIDAndName(projectID uint, name string) (*models.ViewpointItem, error) {
	var item models.ViewpointItem
	if err := r.db.Where("project_id = ? AND name = ?", projectID, name).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// BulkCreate 批量创建AI观点条目
func (r *viewpointItemRepository) BulkCreate(items []*models.ViewpointItem) error {
	return r.db.Create(items).Error
}

// BulkUpdate 批量更新AI观点条目
func (r *viewpointItemRepository) BulkUpdate(items []*models.ViewpointItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Save(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BulkDelete 批量删除AI观点条目(硬删除)
func (r *viewpointItemRepository) BulkDelete(ids []uint) error {
	return r.db.Unscoped().Delete(&models.ViewpointItem{}, ids).Error
}
