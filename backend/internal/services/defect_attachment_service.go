package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"gorm.io/gorm"
)

// DefectAttachmentService 缺陷附件服务接口
type DefectAttachmentService interface {
	Upload(defectID string, projectID uint, userID uint, file *multipart.FileHeader) (*models.AttachmentUploadResponse, error)
	Download(id uint) (*models.DefectAttachment, io.ReadCloser, error)
	Delete(id uint) error
	ListByDefectID(defectID string) ([]*models.DefectAttachment, error)
}

type defectAttachmentService struct {
	repo        repositories.DefectAttachmentRepository
	storagePath string
}

// NewDefectAttachmentService 创建缺陷附件服务实例
func NewDefectAttachmentService(repo repositories.DefectAttachmentRepository, storagePath string) DefectAttachmentService {
	return &defectAttachmentService{
		repo:        repo,
		storagePath: storagePath,
	}
}

// Upload 上传附件
func (s *defectAttachmentService) Upload(defectID string, projectID uint, userID uint, file *multipart.FileHeader) (*models.AttachmentUploadResponse, error) {
	// 验证文件大小
	if file.Size > models.MaxAttachmentSize {
		return nil, fmt.Errorf("file size exceeds limit of %d bytes", models.MaxAttachmentSize)
	}

	// 获取MIME类型
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 验证文件类型（可选，根据需求启用）
	// if !models.IsAllowedMimeType(mimeType) {
	// 	return nil, errors.New("file type not allowed")
	// }

	// 构建存储路径
	timestamp := time.Now().Unix()
	fileName := filepath.Base(file.Filename)
	storageName := fmt.Sprintf("%d_%s", timestamp, fileName)
	relPath := filepath.Join("defects", fmt.Sprintf("%d", projectID), defectID)
	fullDir := filepath.Join(s.storagePath, relPath)
	fullPath := filepath.Join(fullDir, storageName)

	// 创建目录
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("create destination file: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		// 清理已创建的文件
		os.Remove(fullPath)
		return nil, fmt.Errorf("copy file content: %w", err)
	}

	// 创建附件记录
	attachment := &models.DefectAttachment{
		DefectID:   defectID,
		FileName:   fileName,
		FilePath:   filepath.Join(relPath, storageName),
		FileSize:   file.Size,
		MimeType:   mimeType,
		UploadedBy: userID,
	}

	if err := s.repo.Create(attachment); err != nil {
		// 清理已创建的文件
		os.Remove(fullPath)
		return nil, fmt.Errorf("create attachment record: %w", err)
	}

	log.Printf("[Attachment Upload] defect_id=%s, file=%s, size=%d", defectID, fileName, file.Size)

	return &models.AttachmentUploadResponse{
		ID:       attachment.ID,
		FileName: attachment.FileName,
		FileSize: attachment.FileSize,
	}, nil
}

// Download 下载附件
func (s *defectAttachmentService) Download(id uint) (*models.DefectAttachment, io.ReadCloser, error) {
	attachment, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("attachment not found")
		}
		return nil, nil, fmt.Errorf("get attachment: %w", err)
	}

	fullPath := filepath.Join(s.storagePath, attachment.FilePath)
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, errors.New("attachment file not found")
		}
		return nil, nil, fmt.Errorf("open attachment file: %w", err)
	}

	return attachment, file, nil
}

// Delete 删除附件
func (s *defectAttachmentService) Delete(id uint) error {
	// 获取附件信息
	attachment, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("attachment not found")
		}
		return fmt.Errorf("get attachment: %w", err)
	}

	// 删除数据库记录（软删除）
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete attachment record: %w", err)
	}

	// 删除物理文件（可选，根据策略决定是否立即删除）
	fullPath := filepath.Join(s.storagePath, attachment.FilePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		// 记录警告但不返回错误，因为数据库记录已删除
		log.Printf("[Attachment Delete Warning] failed to delete file: %s, error: %v", fullPath, err)
	}

	log.Printf("[Attachment Delete] id=%d, file=%s", id, attachment.FileName)
	return nil
}

// ListByDefectID 获取缺陷的附件列表
func (s *defectAttachmentService) ListByDefectID(defectID string) ([]*models.DefectAttachment, error) {
	attachments, err := s.repo.ListByDefectID(defectID)
	if err != nil {
		return nil, fmt.Errorf("list attachments: %w", err)
	}
	return attachments, nil
}
