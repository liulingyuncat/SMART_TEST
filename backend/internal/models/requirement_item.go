package models

import (
	"time"

	"gorm.io/gorm"
)

// RequirementItem 需求条目模型(动态需求列表方案)
type RequirementItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProjectID uint           `gorm:"not null;uniqueIndex:idx_requirement_items_project_name" json:"project_id"`
	Name      string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_requirement_items_project_name" json:"name"`
	Content   string         `gorm:"type:text" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 外键关联
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定表名
func (RequirementItem) TableName() string {
	return "requirement_items"
}
