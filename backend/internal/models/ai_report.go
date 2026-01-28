package models

import (
	"time"

	"gorm.io/gorm"
)

// AIReport AI报告模型
// 支持项目级报告管理,包含Markdown内容
// Type字段区分报告类型: R(用例审阅)/A(品质分析)/T(测试结果)/O(其他)
type AIReport struct {
	ID        string         `gorm:"type:varchar(50);primaryKey" json:"id"` // report_前缀加雪花ID
	ProjectID uint           `gorm:"not null;index:idx_ai_reports_project" json:"project_id"`
	Type      string         `gorm:"type:varchar(1);not null;default:'O';index:idx_ai_reports_type" json:"type"` // R/A/T/O
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Content   string         `gorm:"type:text" json:"content"` // Markdown格式
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 外键关联
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定表名
func (AIReport) TableName() string {
	return "ai_reports"
}
