package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// RequirementChunkRepository 需求Chunk仓库接口
type RequirementChunkRepository interface {
	Create(chunk *models.RequirementChunk) error
	Update(chunk *models.RequirementChunk) error
	Delete(id uint) error
	FindByID(id uint) (*models.RequirementChunk, error)
	FindByRequirementID(requirementID uint) ([]*models.RequirementChunk, error)
	FindByRequirementIDs(requirementIDs []uint) ([]*models.RequirementChunk, error)
	UpdateSortOrders(chunkOrders []ChunkOrder) error
	GetMaxSortOrder(requirementID uint) (int, error)
}

// ChunkOrder 用于批量更新排序的结构
type ChunkOrder struct {
	ID        uint `json:"id"`
	SortOrder int  `json:"sort_order"`
}

// requirementChunkRepository 需求Chunk仓库实现
type requirementChunkRepository struct {
	db *gorm.DB
}

// NewRequirementChunkRepository 创建需求Chunk仓库实例
func NewRequirementChunkRepository(db *gorm.DB) RequirementChunkRepository {
	return &requirementChunkRepository{db: db}
}

// Create 创建Chunk
func (r *requirementChunkRepository) Create(chunk *models.RequirementChunk) error {
	return r.db.Create(chunk).Error
}

// Update 更新Chunk
func (r *requirementChunkRepository) Update(chunk *models.RequirementChunk) error {
	return r.db.Save(chunk).Error
}

// Delete 软删除Chunk
func (r *requirementChunkRepository) Delete(id uint) error {
	return r.db.Delete(&models.RequirementChunk{}, id).Error
}

// FindByID 根据ID查询Chunk
func (r *requirementChunkRepository) FindByID(id uint) (*models.RequirementChunk, error) {
	var chunk models.RequirementChunk
	if err := r.db.First(&chunk, id).Error; err != nil {
		return nil, err
	}
	return &chunk, nil
}

// FindByRequirementID 根据需求ID查询所有Chunk（按sort_order升序）
func (r *requirementChunkRepository) FindByRequirementID(requirementID uint) ([]*models.RequirementChunk, error) {
	var chunks []*models.RequirementChunk
	if err := r.db.Where("requirement_id = ?", requirementID).Order("sort_order ASC").Find(&chunks).Error; err != nil {
		return nil, err
	}
	return chunks, nil
}

// FindByRequirementIDs 批量根据需求ID查询所有Chunk（按sort_order升序）
func (r *requirementChunkRepository) FindByRequirementIDs(requirementIDs []uint) ([]*models.RequirementChunk, error) {
	if len(requirementIDs) == 0 {
		return []*models.RequirementChunk{}, nil
	}
	var chunks []*models.RequirementChunk
	if err := r.db.Where("requirement_id IN ?", requirementIDs).Order("requirement_id, sort_order ASC").Find(&chunks).Error; err != nil {
		return nil, err
	}
	return chunks, nil
}

// UpdateSortOrders 批量更新Chunk排序
func (r *requirementChunkRepository) UpdateSortOrders(chunkOrders []ChunkOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, order := range chunkOrders {
			if err := tx.Model(&models.RequirementChunk{}).Where("id = ?", order.ID).Update("sort_order", order.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetMaxSortOrder 获取需求下最大的sort_order
func (r *requirementChunkRepository) GetMaxSortOrder(requirementID uint) (int, error) {
	var maxOrder int
	err := r.db.Model(&models.RequirementChunk{}).
		Where("requirement_id = ?", requirementID).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}
