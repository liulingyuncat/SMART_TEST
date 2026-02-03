package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// ViewpointChunkRepository 观点Chunk仓库接口
type ViewpointChunkRepository interface {
	Create(chunk *models.ViewpointChunk) error
	Update(chunk *models.ViewpointChunk) error
	Delete(id uint) error
	FindByID(id uint) (*models.ViewpointChunk, error)
	FindByViewpointID(viewpointID uint) ([]*models.ViewpointChunk, error)
	FindByViewpointIDs(viewpointIDs []uint) ([]*models.ViewpointChunk, error)
	UpdateSortOrders(chunkOrders []ChunkOrder) error
	GetMaxSortOrder(viewpointID uint) (int, error)
}

// viewpointChunkRepository 观点Chunk仓库实现
type viewpointChunkRepository struct {
	db *gorm.DB
}

// NewViewpointChunkRepository 创建观点Chunk仓库实例
func NewViewpointChunkRepository(db *gorm.DB) ViewpointChunkRepository {
	return &viewpointChunkRepository{db: db}
}

// Create 创建Chunk
func (r *viewpointChunkRepository) Create(chunk *models.ViewpointChunk) error {
	return r.db.Create(chunk).Error
}

// Update 更新Chunk
func (r *viewpointChunkRepository) Update(chunk *models.ViewpointChunk) error {
	return r.db.Save(chunk).Error
}

// Delete 软删除Chunk
func (r *viewpointChunkRepository) Delete(id uint) error {
	return r.db.Delete(&models.ViewpointChunk{}, id).Error
}

// FindByID 根据ID查询Chunk
func (r *viewpointChunkRepository) FindByID(id uint) (*models.ViewpointChunk, error) {
	var chunk models.ViewpointChunk
	if err := r.db.First(&chunk, id).Error; err != nil {
		return nil, err
	}
	return &chunk, nil
}

// FindByViewpointID 根据观点ID查询所有Chunk（按sort_order升序）
func (r *viewpointChunkRepository) FindByViewpointID(viewpointID uint) ([]*models.ViewpointChunk, error) {
	var chunks []*models.ViewpointChunk
	if err := r.db.Where("viewpoint_id = ?", viewpointID).Order("sort_order ASC").Find(&chunks).Error; err != nil {
		return nil, err
	}
	return chunks, nil
}

// FindByViewpointIDs 批量根据观点ID查询所有Chunk（按sort_order升序）
func (r *viewpointChunkRepository) FindByViewpointIDs(viewpointIDs []uint) ([]*models.ViewpointChunk, error) {
	if len(viewpointIDs) == 0 {
		return []*models.ViewpointChunk{}, nil
	}
	var chunks []*models.ViewpointChunk
	if err := r.db.Where("viewpoint_id IN ?", viewpointIDs).Order("viewpoint_id, sort_order ASC").Find(&chunks).Error; err != nil {
		return nil, err
	}
	return chunks, nil
}

// UpdateSortOrders 批量更新Chunk排序
func (r *viewpointChunkRepository) UpdateSortOrders(chunkOrders []ChunkOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, order := range chunkOrders {
			if err := tx.Model(&models.ViewpointChunk{}).Where("id = ?", order.ID).Update("sort_order", order.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetMaxSortOrder 获取观点下最大的sort_order
func (r *viewpointChunkRepository) GetMaxSortOrder(viewpointID uint) (int, error) {
	var maxOrder int
	err := r.db.Model(&models.ViewpointChunk{}).
		Where("viewpoint_id = ?", viewpointID).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}
