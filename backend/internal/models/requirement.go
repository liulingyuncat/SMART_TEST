package models

import (
	"time"

	"gorm.io/gorm"
)

// Requirement 需求文档模型(单表多字段方案)
type Requirement struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	ProjectID            uint           `gorm:"not null;uniqueIndex:idx_requirements_project" json:"project_id"`
	OverallRequirements  string         `gorm:"type:text" json:"overall_requirements"`   // 整体需求
	OverallTestViewpoint string         `gorm:"type:text" json:"overall_test_viewpoint"` // 整体测试观点
	ChangeRequirements   string         `gorm:"type:text" json:"change_requirements"`    // 变更需求
	ChangeTestViewpoint  string         `gorm:"type:text" json:"change_test_viewpoint"`  // 变更测试观点
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// 外键关联
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定表名
func (Requirement) TableName() string {
	return "requirements"
}
