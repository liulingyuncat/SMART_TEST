package models

import (
	"time"

	"gorm.io/gorm"
)

// Project 项目模型
type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"type:varchar(20);default:pending;not null" json:"status"`
	OwnerID     *int           `gorm:"type:int;index" json:"owner_id"`
	OwnerName   string         `gorm:"-" json:"owner_name,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}
