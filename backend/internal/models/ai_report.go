package models

import (
	"time"

	"gorm.io/gorm"
)

// AIReport AI质量报告模型
// 支持项目级报告管理,包含Markdown内容、版本控制
type AIReport struct {
	ID        string         `gorm:"type:varchar(50);primaryKey" json:"id"` // report_前缀加雪花ID
	ProjectID uint           `gorm:"not null;index:idx_ai_reports_project" json:"project_id"`
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
