package models

import (
	"time"

	"gorm.io/gorm"
)

// ViewpointItem AI观点条目模型(动态观点列表方案)
type ViewpointItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProjectID uint           `gorm:"not null;uniqueIndex:idx_viewpoint_items_project_name" json:"project_id"`
	Name      string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_viewpoint_items_project_name" json:"name"`
	Content   string         `gorm:"type:text" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 外键关联
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定表名
func (ViewpointItem) TableName() string {
	return "viewpoint_items"
}
