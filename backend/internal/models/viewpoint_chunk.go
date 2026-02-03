package models

import (
	"time"

	"gorm.io/gorm"
)

// ViewpointChunk 观点Chunk模型（段落/章节单元）
type ViewpointChunk struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ViewpointID uint           `gorm:"not null;index" json:"viewpoint_id"`
	Title       string         `gorm:"type:varchar(255)" json:"title"`
	Content     string         `gorm:"type:text" json:"content"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 外键关联
	ViewpointItem ViewpointItem `gorm:"foreignKey:ViewpointID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定表名
func (ViewpointChunk) TableName() string {
	return "viewpoint_chunks"
}
