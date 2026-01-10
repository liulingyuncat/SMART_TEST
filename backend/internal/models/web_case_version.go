package models

import (
	"time"
)

// WebCaseVersion Web用例版本模型
type WebCaseVersion struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VersionID   string    `gorm:"type:varchar(100);not null;index" json:"version_id"`
	ProjectID   uint      `gorm:"not null;index" json:"project_id"`
	ProjectName string    `gorm:"type:varchar(100);not null" json:"project_name"`
	ZipFilename string    `gorm:"type:varchar(255);not null" json:"zip_filename"`
	ZipPath     string    `gorm:"type:varchar(500);not null" json:"zip_path"`
	FileSize    int64     `gorm:"type:bigint" json:"file_size"`
	CaseCount   int       `gorm:"default:0" json:"case_count"`
	Remark      string    `gorm:"type:varchar(200);default:''" json:"remark"`
	CreatedBy   *uint     `gorm:"index" json:"created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index:idx_wcv_created,sort:desc" json:"created_at"`
}

// TableName 指定表名
func (WebCaseVersion) TableName() string {
	return "web_case_versions"
}
