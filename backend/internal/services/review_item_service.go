package services

import (
	"errors"
	"fmt"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"gorm.io/gorm"
)

// ReviewItemService 审阅条目服务接口
type ReviewItemService interface {
	// CreateReviewItem 创建审阅条目
	CreateReviewItem(projectID uint, name string) (*models.CaseReviewItem, error)

	// GetReviewItem 获取审阅条目详情
	GetReviewItem(id uint) (*models.CaseReviewItem, error)

	// ListReviewItems 获取项目所有审阅条目列表
	ListReviewItems(projectID uint) ([]models.CaseReviewItem, error)

	// UpdateReviewItem 更新审阅条目
	UpdateReviewItem(id uint, name, content *string) (*models.CaseReviewItem, error)

	// DeleteReviewItem 删除审阅条目
	DeleteReviewItem(id uint) error

	// DownloadReviewItem 生成审阅文档的Markdown文件内容
	DownloadReviewItem(id uint, projectName string) (string, string, error)
}

// reviewItemService 审阅条目服务实现
type reviewItemService struct {
	repo repositories.ReviewItemRepository
}

// NewReviewItemService 创建审阅条目服务实例
func NewReviewItemService(repo repositories.ReviewItemRepository) ReviewItemService {
	return &reviewItemService{repo: repo}
}

// CreateReviewItem 实现创建审阅条目
func (s *reviewItemService) CreateReviewItem(projectID uint, name string) (*models.CaseReviewItem, error) {
	// 校验名称唯一性
	existingItem, err := s.repo.GetByProjectAndName(projectID, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingItem != nil {
		return nil, errors.New("审阅名称已存在")
	}

	// 创建新审阅条目
	item := &models.CaseReviewItem{
		ProjectID: projectID,
		Name:      name,
		Content:   "",
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetReviewItem 实现获取审阅条目详情
func (s *reviewItemService) GetReviewItem(id uint) (*models.CaseReviewItem, error) {
	return s.repo.GetByID(id)
}

// ListReviewItems 实现获取项目所有审阅条目列表
func (s *reviewItemService) ListReviewItems(projectID uint) ([]models.CaseReviewItem, error) {
	return s.repo.GetByProjectID(projectID)
}

// UpdateReviewItem 实现更新审阅条目
func (s *reviewItemService) UpdateReviewItem(id uint, name, content *string) (*models.CaseReviewItem, error) {
	// 获取现有记录
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 如果更新名称,检查唯一性
	if name != nil && *name != item.Name {
		existingItem, err := s.repo.GetByProjectAndName(item.ProjectID, *name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingItem != nil && existingItem.ID != id {
			return nil, errors.New("审阅名称已存在")
		}
		item.Name = *name
	}

	// 更新内容
	if content != nil {
		item.Content = *content
	}

	// 保存更新
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteReviewItem 实现删除审阅条目
func (s *reviewItemService) DeleteReviewItem(id uint) error {
	return s.repo.Delete(id)
}

// DownloadReviewItem 实现生成Markdown文件内容
func (s *reviewItemService) DownloadReviewItem(id uint, projectName string) (string, string, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return "", "", err
	}

	// 生成文件名: 项目名_审阅名_时间戳.md
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s_%s.md", projectName, item.Name, timestamp)

	return item.Content, filename, nil
}
