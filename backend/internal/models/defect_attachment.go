package models

import (
	"time"

	"gorm.io/gorm"
)

// DefectAttachment 缺陷附件模型
type DefectAttachment struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	DefectID   string `gorm:"type:varchar(36);not null;index:idx_defect_attachments_defect_id" json:"defect_id"` // 关联缺陷UUID
	FileName   string `gorm:"type:varchar(255);not null" json:"file_name"`                                       // 原始文件名
	FilePath   string `gorm:"type:varchar(500);not null" json:"file_path"`                                       // 服务器存储路径
	FileSize   int64  `gorm:"not null" json:"file_size"`                                                         // 文件大小(字节)
	MimeType   string `gorm:"type:varchar(100)" json:"mime_type"`                                                // MIME类型
	UploadedBy uint   `gorm:"not null" json:"uploaded_by"`                                                       // 上传人ID

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_defect_attachments_deleted_at" json:"-"`
}

// TableName 指定表名
func (DefectAttachment) TableName() string {
	return "defect_attachments"
}

// MaxAttachmentSize 最大附件大小 100MB
const MaxAttachmentSize = 100 * 1024 * 1024

// AllowedMimeTypes 允许的文件类型
var AllowedMimeTypes = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"image/webp":         true,
	"image/bmp":          true,
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.ms-excel": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
	"application/vnd.ms-powerpoint":                                             true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
	"text/plain":                   true,
	"text/csv":                     true,
	"application/json":             true,
	"application/xml":              true,
	"text/xml":                     true,
	"application/zip":              true,
	"application/x-rar-compressed": true,
	"application/x-7z-compressed":  true,
	"video/mp4":                    true,
	"video/webm":                   true,
	"audio/mpeg":                   true,
	"audio/wav":                    true,
}

// IsAllowedMimeType 检查MIME类型是否允许
func IsAllowedMimeType(mimeType string) bool {
	return AllowedMimeTypes[mimeType]
}

// AttachmentUploadResponse 附件上传响应
type AttachmentUploadResponse struct {
	ID       uint   `json:"id"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}
