package models

import "time"

// ProjectMember 项目成员模型
type ProjectMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"not null;index:idx_pm_project_user,priority:1" json:"project_id"`
	UserID    uint      `gorm:"not null;index:idx_pm_project_user,priority:2" json:"user_id"`
	Role      string    `gorm:"not null" json:"role"` // project_manager, project_member
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (ProjectMember) TableName() string {
	return "project_members"
}
