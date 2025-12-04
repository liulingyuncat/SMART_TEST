package models

import "time"

// CaseVersion 测试用例版本记录模型
type CaseVersion struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	ProjectID uint       `gorm:"not null;index:idx_versions_project" json:"project_id"`
	DocType   string     `gorm:"type:varchar(50);not null;default:'overall';column:case_type" json:"doc_type"` // 文档类型字段,存储在case_type列
	Filename  string     `gorm:"type:varchar(255);not null" json:"filename"`
	FilePath  string     `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize  int64      `json:"file_size"`
	Remark    string     `gorm:"type:text" json:"remark"` // 备注信息
	CreatedBy *uint      `json:"created_by"`              // 指针类型支持NULL
	CreatedAt time.Time  `gorm:"index:idx_versions_created,sort:desc" json:"created_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"` // 软删除
}

// TableName 指定表名
func (CaseVersion) TableName() string {
	return "test_case_versions"
}
