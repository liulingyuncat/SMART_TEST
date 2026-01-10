package models

import (
	"time"

	"gorm.io/gorm"
)

// RawDocument 原始需求文档模型
type RawDocument struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID         uint           `gorm:"not null;index:idx_raw_documents_project_id" json:"project_id"` // 关联项目ID
	OriginalFilename  string         `gorm:"type:varchar(255);not null" json:"original_filename"`           // 原始文件名
	OriginalFilepath  string         `gorm:"type:varchar(500);not null" json:"-"`                           // 原始文件存储路径（不返回前端）
	FileSize          int64          `gorm:"not null" json:"file_size"`                                     // 文件大小(字节)
	MimeType          string         `gorm:"type:varchar(100);not null" json:"mime_type"`                   // MIME类型
	UploadedBy        uint           `gorm:"not null" json:"uploaded_by"`                                   // 上传人ID
	ConvertStatus     string         `gorm:"type:varchar(20);default:'none'" json:"convert_status"`         // 转换状态: none/processing/completed/failed
	ConvertTaskID     string         `gorm:"type:varchar(100)" json:"convert_task_id,omitempty"`            // 转换任务ID
	ConvertProgress   int            `gorm:"default:0" json:"convert_progress"`                             // 转换进度 0-100
	ConvertedFilename string         `gorm:"type:varchar(255)" json:"converted_filename,omitempty"`         // 转换后文件名
	ConvertedFilepath string         `gorm:"type:varchar(500)" json:"-"`                                    // 转换后文件存储路径（不返回前端）
	ConvertedFileSize int64          `gorm:"default:0" json:"converted_file_size,omitempty"`                // 转换后文件大小(字节)
	ConvertedTime     *time.Time     `json:"converted_time,omitempty"`                                      // 转换完成时间
	ConvertError      string         `gorm:"type:text" json:"convert_error,omitempty"`                      // 转换错误信息
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index:idx_raw_documents_deleted_at" json:"-"` // 软删除
}

// TableName 指定表名
func (RawDocument) TableName() string {
	return "raw_documents"
}

// MaxRawDocumentSize 最大原始文档大小 100MB
const MaxRawDocumentSize = 100 * 1024 * 1024

// RawDocumentAllowedMimeTypes 原始需求文档允许的文件类型（15种格式）
var RawDocumentAllowedMimeTypes = map[string]bool{
	// 文档类型
	"application/pdf":    true, // PDF
	"application/msword": true, // DOC
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // DOCX
	"text/plain":      true, // TXT
	"application/rtf": true, // RTF

	// 图片类型
	"image/png":  true, // PNG
	"image/jpeg": true, // JPG/JPEG
	"image/bmp":  true, // BMP
	"image/tiff": true, // TIFF

	// 表格类型
	"application/vnd.ms-excel": true, // XLS
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true, // XLSX
	"text/csv": true, // CSV

	// 演示文稿类型
	"application/vnd.ms-powerpoint":                                             true, // PPT
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true, // PPTX
}

// IsRawDocumentAllowedMimeType 检查原始文档MIME类型是否在白名单中
func IsRawDocumentAllowedMimeType(mimeType string) bool {
	return RawDocumentAllowedMimeTypes[mimeType]
}

// RawDocumentUploadResponse 原始文档上传响应
type RawDocumentUploadResponse struct {
	ID               uint      `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	FileSize         int64     `json:"file_size"`
	MimeType         string    `json:"mime_type"`
	UploadTime       time.Time `json:"upload_time"`
}

// RawDocumentListItem 原始文档列表项
type RawDocumentListItem struct {
	ID                uint       `json:"id"`
	ProjectID         uint       `json:"project_id"`
	OriginalFilename  string     `json:"original_filename"`
	FileSize          int64      `json:"file_size"`
	MimeType          string     `json:"mime_type"`
	UploadedBy        uint       `json:"uploaded_by"`
	UploaderName      string     `json:"uploader_name,omitempty"` // 关联查询上传人姓名
	ConvertStatus     string     `json:"convert_status"`
	ConvertProgress   int        `json:"convert_progress"`
	ConvertedFilename string     `json:"converted_filename,omitempty"`
	ConvertedFileSize int64      `json:"converted_file_size,omitempty"`
	ConvertedTime     *time.Time `json:"converted_time,omitempty"`
	ConvertError      string     `json:"convert_error,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

// ConvertTaskResponse 转换任务响应
type ConvertTaskResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}

// ConvertStatusResponse 转换状态响应
type ConvertStatusResponse struct {
	Status            string `json:"status"`
	Progress          int    `json:"progress"`
	ConvertedFilename string `json:"converted_filename,omitempty"`
	ErrorMessage      string `json:"error_message,omitempty"`
}
