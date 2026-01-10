package models

import (
	"time"

	"gorm.io/gorm"
)

// CaseGroup 用例集
type CaseGroup struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	ProjectID    uint   `json:"project_id" gorm:"not null;index:idx_cg_project_type"`
	CaseType     string `json:"case_type" gorm:"type:varchar(20);not null;default:'overall';index:idx_cg_project_type"`
	GroupName    string `json:"group_name" gorm:"type:varchar(100);not null;uniqueIndex:idx_project_type_name"`
	Description  string `json:"description" gorm:"type:text"`
	DisplayOrder int    `json:"display_order" gorm:"default:0;index"`
	// 元数据字段
	MetaProtocol string         `json:"meta_protocol" gorm:"type:varchar(20);default:'https'"`
	MetaServer   string         `json:"meta_server" gorm:"type:varchar(255)"`
	MetaPort     string         `json:"meta_port" gorm:"type:varchar(20)"`
	MetaUser     string         `json:"meta_user" gorm:"type:varchar(100)"`
	MetaPassword string         `json:"meta_password" gorm:"type:varchar(255)"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"` // 不返回给前端
}

// TableName 指定表名
func (CaseGroup) TableName() string {
	return "case_groups"
}
