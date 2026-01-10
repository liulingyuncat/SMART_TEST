package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Username         string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Nickname         string         `gorm:"uniqueIndex;not null;size:50" json:"nickname"`
	Password         string         `gorm:"not null;size:255" json:"-"`      // 不序列化到JSON
	Role             string         `gorm:"not null;size:20" json:"role"`    // system_admin, project_manager, project_member
	ApiToken         *string        `gorm:"size:128;index" json:"-"`         // API Token，不序列化到JSON
	CurrentProjectID *uint          `gorm:"index" json:"current_project_id"` // 当前选择的项目ID
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
