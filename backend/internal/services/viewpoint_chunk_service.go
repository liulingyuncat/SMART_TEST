package services

import (
	"fmt"
	"webtest/internal/models"
	"webtest/internal/repositories"
)

// ViewpointChunkService 观点Chunk服务接口
type ViewpointChunkService interface {
	CreateChunk(viewpointID uint, title, content string) (*models.ViewpointChunk, error)
	UpdateChunk(id uint, title, content string) (*models.ViewpointChunk, error)
	DeleteChunk(id uint) error
	GetChunkByID(id uint) (*models.ViewpointChunk, error)
	GetChunksByViewpointID(viewpointID uint) ([]*models.ViewpointChunk, error)
	ReorderChunks(chunkOrders []repositories.ChunkOrder) error
}

// viewpointChunkService 观点Chunk服务实现
type viewpointChunkService struct {
	chunkRepo repositories.ViewpointChunkRepository
}

// NewViewpointChunkService 创建观点Chunk服务实例
func NewViewpointChunkService(chunkRepo repositories.ViewpointChunkRepository) ViewpointChunkService {
	return &viewpointChunkService{
		chunkRepo: chunkRepo,
	}
}

// CreateChunk 创建Chunk
func (s *viewpointChunkService) CreateChunk(viewpointID uint, title, content string) (*models.ViewpointChunk, error) {
	// 获取当前最大sort_order
	maxOrder, err := s.chunkRepo.GetMaxSortOrder(viewpointID)
	if err != nil {
		return nil, fmt.Errorf("获取排序序号失败: %w", err)
	}

	chunk := &models.ViewpointChunk{
		ViewpointID: viewpointID,
		Title:       title,
		Content:     content,
		SortOrder:   maxOrder + 1,
	}

	if err := s.chunkRepo.Create(chunk); err != nil {
		return nil, fmt.Errorf("创建Chunk失败: %w", err)
	}

	return chunk, nil
}

// UpdateChunk 更新Chunk
func (s *viewpointChunkService) UpdateChunk(id uint, title, content string) (*models.ViewpointChunk, error) {
	chunk, err := s.chunkRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("Chunk不存在: %w", err)
	}

	chunk.Title = title
	chunk.Content = content

	if err := s.chunkRepo.Update(chunk); err != nil {
		return nil, fmt.Errorf("更新Chunk失败: %w", err)
	}

	return chunk, nil
}

// DeleteChunk 软删除Chunk
func (s *viewpointChunkService) DeleteChunk(id uint) error {
	if _, err := s.chunkRepo.FindByID(id); err != nil {
		return fmt.Errorf("Chunk不存在: %w", err)
	}

	if err := s.chunkRepo.Delete(id); err != nil {
		return fmt.Errorf("删除Chunk失败: %w", err)
	}

	return nil
}

// GetChunkByID 根据ID获取Chunk
func (s *viewpointChunkService) GetChunkByID(id uint) (*models.ViewpointChunk, error) {
	chunk, err := s.chunkRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("Chunk不存在: %w", err)
	}
	return chunk, nil
}

// GetChunksByViewpointID 获取观点下所有Chunk
func (s *viewpointChunkService) GetChunksByViewpointID(viewpointID uint) ([]*models.ViewpointChunk, error) {
	chunks, err := s.chunkRepo.FindByViewpointID(viewpointID)
	if err != nil {
		return nil, fmt.Errorf("获取Chunk列表失败: %w", err)
	}
	return chunks, nil
}

// ReorderChunks 批量重排序Chunk
func (s *viewpointChunkService) ReorderChunks(chunkOrders []repositories.ChunkOrder) error {
	if err := s.chunkRepo.UpdateSortOrders(chunkOrders); err != nil {
		return fmt.Errorf("重排序Chunk失败: %w", err)
	}
	return nil
}
