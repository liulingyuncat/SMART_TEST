package repositories

import (
	"fmt"
	"time"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// CaseReviewRepository 评审记录仓储接口
type CaseReviewRepository interface {
	GetByProjectAndType(projectID uint, caseType string) (*models.CaseReview, error)
	Upsert(review *models.CaseReview) error
}

type caseReviewRepository struct {
	db *gorm.DB
}

// NewCaseReviewRepository 创建评审记录仓储实例
func NewCaseReviewRepository(db *gorm.DB) CaseReviewRepository {
	return &caseReviewRepository{db: db}
}

// GetByProjectAndType 根据项目ID和用例类型查询评审记录
func (r *caseReviewRepository) GetByProjectAndType(projectID uint, caseType string) (*models.CaseReview, error) {
	var review models.CaseReview
	err := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).First(&review).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}
	return &review, nil
}

// Upsert 插入或更新评审记录(兼容SQLite和PostgreSQL)
func (r *caseReviewRepository) Upsert(review *models.CaseReview) error {
	// 添加调试日志
	fmt.Printf("[DEBUG] Upsert review - ProjectID: %d, CaseType: %s, ContentLength: %d\n",
		review.ProjectID, review.CaseType, len(review.Content))

	// 先查询是否存在
	existing, err := r.GetByProjectAndType(review.ProjectID, review.CaseType)
	if err != nil {
		fmt.Printf("[ERROR] Query existing review failed - Error: %v\n", err)
		return fmt.Errorf("failed to query existing review: %w", err)
	}

	if existing != nil {
		// 更新现有记录
		fmt.Printf("[DEBUG] Updating existing review ID: %d\n", existing.ID)
		err = r.db.Model(&models.CaseReview{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
			"content":    review.Content,
			"updated_at": time.Now(),
		}).Error
	} else {
		// 插入新记录
		fmt.Printf("[DEBUG] Creating new review\n")
		err = r.db.Create(review).Error
	}

	if err != nil {
		fmt.Printf("[ERROR] Upsert failed - Error: %v\n", err)
		return fmt.Errorf("failed to upsert review: %w", err)
	}

	fmt.Printf("[DEBUG] Upsert successful\n")
	return nil
}
