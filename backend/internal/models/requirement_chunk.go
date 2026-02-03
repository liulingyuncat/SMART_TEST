package models

import (
	"time"

	"gorm.io/gorm"
)

// RequirementChunk 需求Chunk模型（段落/章节单元）
type RequirementChunk struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	RequirementID uint           `gorm:"not null;index" json:"requirement_id"`
	Title         string         `gorm:"type:varchar(255)" json:"title"`
	Content       string         `gorm:"type:text" json:"content"`
	SortOrder     int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 外键关联
	RequirementItem RequirementItem `gorm:"foreignKey:RequirementID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定表名
func (RequirementChunk) TableName() string {
	return "requirement_chunks"
}
