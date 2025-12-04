package services

import (
	"fmt"
	"webtest/internal/models"
	"webtest/internal/repositories"
)

// ReviewService 评审服务接口
type ReviewService interface {
	GetReview(projectID uint, caseType string) (string, error)
	SaveReview(projectID uint, caseType string, content string) error
}

type reviewService struct {
	reviewRepo repositories.CaseReviewRepository
}

// NewReviewService 创建评审服务实例
func NewReviewService(reviewRepo repositories.CaseReviewRepository) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
	}
}

// GetReview 获取评审记录
func (s *reviewService) GetReview(projectID uint, caseType string) (string, error) {
	review, err := s.reviewRepo.GetByProjectAndType(projectID, caseType)
	if err != nil {
		return "", fmt.Errorf("get review: %w", err)
	}
	if review == nil {
		return "", nil
	}
	return review.Content, nil
}

// SaveReview 保存评审记录(UPSERT)
func (s *reviewService) SaveReview(projectID uint, caseType string, content string) error {
	review := &models.CaseReview{
		ProjectID: projectID,
		CaseType:  caseType,
		Content:   content,
	}
	if err := s.reviewRepo.Upsert(review); err != nil {
		return fmt.Errorf("save review: %w", err)
	}
	return nil
}
