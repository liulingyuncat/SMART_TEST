package services

import (
	"fmt"
	"webtest/internal/models"
	"webtest/internal/repositories"
)

// RequirementChunkService 需求Chunk服务接口
type RequirementChunkService interface {
	CreateChunk(requirementID uint, title, content string) (*models.RequirementChunk, error)
	UpdateChunk(id uint, title, content string) (*models.RequirementChunk, error)
	DeleteChunk(id uint) error
	GetChunkByID(id uint) (*models.RequirementChunk, error)
	GetChunksByRequirementID(requirementID uint) ([]*models.RequirementChunk, error)
	ReorderChunks(chunkOrders []repositories.ChunkOrder) error
}

// requirementChunkService 需求Chunk服务实现
type requirementChunkService struct {
	chunkRepo repositories.RequirementChunkRepository
}

// NewRequirementChunkService 创建需求Chunk服务实例
func NewRequirementChunkService(chunkRepo repositories.RequirementChunkRepository) RequirementChunkService {
	return &requirementChunkService{
		chunkRepo: chunkRepo,
	}
}

// CreateChunk 创建Chunk
func (s *requirementChunkService) CreateChunk(requirementID uint, title, content string) (*models.RequirementChunk, error) {
	// 获取当前最大sort_order
	maxOrder, err := s.chunkRepo.GetMaxSortOrder(requirementID)
	if err != nil {
		return nil, fmt.Errorf("获取排序序号失败: %w", err)
	}

	chunk := &models.RequirementChunk{
		RequirementID: requirementID,
		Title:         title,
		Content:       content,
		SortOrder:     maxOrder + 1,
	}

	if err := s.chunkRepo.Create(chunk); err != nil {
		return nil, fmt.Errorf("创建Chunk失败: %w", err)
	}

	return chunk, nil
}

// UpdateChunk 更新Chunk
func (s *requirementChunkService) UpdateChunk(id uint, title, content string) (*models.RequirementChunk, error) {
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
func (s *requirementChunkService) DeleteChunk(id uint) error {
	if _, err := s.chunkRepo.FindByID(id); err != nil {
		return fmt.Errorf("Chunk不存在: %w", err)
	}

	if err := s.chunkRepo.Delete(id); err != nil {
		return fmt.Errorf("删除Chunk失败: %w", err)
	}

	return nil
}

// GetChunkByID 根据ID获取Chunk
func (s *requirementChunkService) GetChunkByID(id uint) (*models.RequirementChunk, error) {
	chunk, err := s.chunkRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("Chunk不存在: %w", err)
	}
	return chunk, nil
}

// GetChunksByRequirementID 获取需求下所有Chunk
func (s *requirementChunkService) GetChunksByRequirementID(requirementID uint) ([]*models.RequirementChunk, error) {
	chunks, err := s.chunkRepo.FindByRequirementID(requirementID)
	if err != nil {
		return nil, fmt.Errorf("获取Chunk列表失败: %w", err)
	}
	return chunks, nil
}

// ReorderChunks 批量重排序Chunk
func (s *requirementChunkService) ReorderChunks(chunkOrders []repositories.ChunkOrder) error {
	if err := s.chunkRepo.UpdateSortOrders(chunkOrders); err != nil {
		return fmt.Errorf("重排序Chunk失败: %w", err)
	}
	return nil
}
