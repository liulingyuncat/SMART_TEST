package repositories

import (
	"fmt"
	"time"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// RawDocumentRepository 原始文档仓储接口
type RawDocumentRepository interface {
	Create(doc *models.RawDocument) error
	GetByID(id uint) (*models.RawDocument, error)
	FindByID(id uint) (*models.RawDocument, error)
	ListByProjectID(projectID uint) ([]*models.RawDocument, error)
	Update(doc *models.RawDocument) error
	UpdateStatus(id uint, status string, progress int, filename string, filepath string, filesize int64, convertError string) error
	GetConvertStatus(id uint) (*models.ConvertStatusResponse, error)
	Delete(id uint) error
}

type rawDocumentRepository struct {
	db *gorm.DB
}

// NewRawDocumentRepository 创建原始文档仓储实例
func NewRawDocumentRepository(db *gorm.DB) RawDocumentRepository {
	return &rawDocumentRepository{db: db}
}

// Create 创建原始文档记录
func (r *rawDocumentRepository) Create(doc *models.RawDocument) error {
	err := r.db.Create(doc).Error
	if err != nil {
		return fmt.Errorf("create raw document: %w", err)
	}
	return nil
}

// GetByID 根据ID获取原始文档（过滤软删除）
func (r *rawDocumentRepository) GetByID(id uint) (*models.RawDocument, error) {
	var doc models.RawDocument
	err := r.db.Where("id = ?", id).First(&doc).Error
	if err != nil {
		return nil, err // 保留gorm.ErrRecordNotFound
	}
	return &doc, nil
}

// FindByID 根据ID查询完整的RawDocument记录（包括所有转换状态字段）
func (r *rawDocumentRepository) FindByID(id uint) (*models.RawDocument, error) {
	var doc models.RawDocument
	err := r.db.Where("id = ?", id).First(&doc).Error
	if err != nil {
		return nil, fmt.Errorf("find raw document %d: %w", id, err)
	}
	return &doc, nil
}

// ListByProjectID 根据项目ID获取原始文档列表（按创建时间倒序）
func (r *rawDocumentRepository) ListByProjectID(projectID uint) ([]*models.RawDocument, error) {
	var docs []*models.RawDocument
	err := r.db.Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&docs).Error

	if err != nil {
		return nil, fmt.Errorf("list raw documents by project: %w", err)
	}

	return docs, nil
}

// Update 更新原始文档记录
func (r *rawDocumentRepository) Update(doc *models.RawDocument) error {
	err := r.db.Save(doc).Error
	if err != nil {
		return fmt.Errorf("update raw document %d: %w", doc.ID, err)
	}
	return nil
}

// UpdateStatus 更新文档的转换状态（支持原子更新，防止并发问题）
func (r *rawDocumentRepository) UpdateStatus(id uint, status string, progress int, filename string, filepath string, filesize int64, convertError string) error {
	updates := map[string]interface{}{
		"convert_status":      status,
		"convert_progress":    progress,
		"converted_filename":  filename,
		"converted_filepath":  filepath,
		"converted_file_size": filesize,
	}

	// 仅当转换完成或失败时才更新时间和错误信息
	if status == "completed" || status == "failed" {
		now := time.Now()
		if status == "completed" {
			updates["converted_time"] = now
		}
		if convertError != "" {
			updates["convert_error"] = convertError
		}
	}

	// 直接更新状态，不需要额外条件限制
	result := r.db.Model(&models.RawDocument{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("update convert status for document %d: %w", id, result.Error)
	}

	return nil
}

// GetConvertStatus 快速查询文档的转换状态（轻量级查询）
func (r *rawDocumentRepository) GetConvertStatus(id uint) (*models.ConvertStatusResponse, error) {
	var doc struct {
		ConvertStatus     string
		ConvertProgress   int
		ConvertedFilename string
		ConvertError      string
	}

	err := r.db.Model(&models.RawDocument{}).
		Where("id = ?", id).
		Select("convert_status", "convert_progress", "converted_filename", "convert_error").
		Scan(&doc).Error

	if err != nil {
		return nil, fmt.Errorf("get convert status for document %d: %w", id, err)
	}

	return &models.ConvertStatusResponse{
		Status:            doc.ConvertStatus,
		Progress:          doc.ConvertProgress,
		ConvertedFilename: doc.ConvertedFilename,
		ErrorMessage:      doc.ConvertError,
	}, nil
}

// Delete 硬删除原始文档记录
func (r *rawDocumentRepository) Delete(id uint) error {
	result := r.db.Unscoped().Where("id = ?", id).Delete(&models.RawDocument{})
	if result.Error != nil {
		return fmt.Errorf("delete raw document %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
