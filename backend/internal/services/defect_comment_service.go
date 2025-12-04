package services

import (
	"fmt"
	"log"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"
)

// DefectCommentService 缺陷说明服务接口
type DefectCommentService interface {
	// CRUD操作
	Create(defectID string, userID uint, req *models.DefectCommentCreateRequest) (*models.DefectComment, error)
	GetByID(id uint) (*models.DefectComment, error)
	Update(id uint, userID uint, req *models.DefectCommentUpdateRequest) error
	Delete(id uint, userID uint) error

	// 列表查询
	List(defectID string) (*models.DefectCommentListResponse, error)
}

// defectCommentService 缺陷说明服务实现
type defectCommentService struct {
	commentRepo repositories.DefectCommentRepository
	defectRepo  repositories.DefectRepository
}

// NewDefectCommentService 创建缺陷说明服务实例
func NewDefectCommentService(commentRepo repositories.DefectCommentRepository, defectRepo repositories.DefectRepository) DefectCommentService {
	return &defectCommentService{
		commentRepo: commentRepo,
		defectRepo:  defectRepo,
	}
}

// Create 创建说明
func (s *defectCommentService) Create(defectID string, userID uint, req *models.DefectCommentCreateRequest) (*models.DefectComment, error) {
	// 验证缺陷是否存在（使用GetByDefectID支持显示ID）
	_, err := s.defectRepo.GetByDefectID(defectID)
	if err != nil {
		return nil, fmt.Errorf("defect not found: %w", err)
	}

	// 创建说明记录
	comment := &models.DefectComment{
		DefectID:  defectID,
		Content:   req.Content,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		log.Printf("[DefectComment Error] action=create, defect_id=%s, error=%v", defectID, err)
		return nil, fmt.Errorf("create comment: %w", err)
	}

	log.Printf("[DefectComment Create] defect_id=%s, comment_id=%d, user_id=%d, content_length=%d",
		defectID, comment.ID, userID, len(req.Content))

	return comment, nil
}

// GetByID 根据ID获取说明
func (s *defectCommentService) GetByID(id uint) (*models.DefectComment, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get comment: %w", err)
	}
	return comment, nil
}

// Update 更新说明
func (s *defectCommentService) Update(id uint, userID uint, req *models.DefectCommentUpdateRequest) error {
	// 获取说明记录
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	// 权限检查：仅创建人可编辑
	if comment.CreatedBy != userID {
		log.Printf("[DefectComment PermissionDenied] action=update, comment_id=%d, user_id=%d, creator_id=%d",
			id, userID, comment.CreatedBy)
		return fmt.Errorf("permission denied: only creator can edit")
	}

	// 更新说明
	oldLength := len(comment.Content)
	updates := map[string]interface{}{
		"content":    req.Content,
		"updated_by": userID,
		"updated_at": time.Now(),
	}

	if err := s.commentRepo.Update(id, updates); err != nil {
		log.Printf("[DefectComment Error] action=update, comment_id=%d, error=%v", id, err)
		return fmt.Errorf("update comment: %w", err)
	}

	log.Printf("[DefectComment Update] comment_id=%d, user_id=%d, old_length=%d, new_length=%d",
		id, userID, oldLength, len(req.Content))

	return nil
}

// Delete 删除说明
func (s *defectCommentService) Delete(id uint, userID uint) error {
	// 获取说明记录
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	// 权限检查：仅创建人可删除
	if comment.CreatedBy != userID {
		log.Printf("[DefectComment PermissionDenied] action=delete, comment_id=%d, user_id=%d, creator_id=%d",
			id, userID, comment.CreatedBy)
		return fmt.Errorf("permission denied: only creator can delete")
	}

	// 软删除
	if err := s.commentRepo.Delete(id); err != nil {
		log.Printf("[DefectComment Error] action=delete, comment_id=%d, error=%v", id, err)
		return fmt.Errorf("delete comment: %w", err)
	}

	log.Printf("[DefectComment Delete] comment_id=%d, user_id=%d, defect_id=%s",
		id, userID, comment.DefectID)

	return nil
}

// List 获取缺陷的说明列表
func (s *defectCommentService) List(defectID string) (*models.DefectCommentListResponse, error) {
	comments, err := s.commentRepo.ListByDefectID(defectID)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}

	return &models.DefectCommentListResponse{
		Comments: convertToCommentSlice(comments),
		Total:    int64(len(comments)),
	}, nil
}

// convertToCommentSlice 转换指针切片为值切片
func convertToCommentSlice(comments []*models.DefectComment) []models.DefectComment {
	result := make([]models.DefectComment, len(comments))
	for i, c := range comments {
		result[i] = *c
	}
	return result
}
