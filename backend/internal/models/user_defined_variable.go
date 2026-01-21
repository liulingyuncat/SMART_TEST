package models

import (
	"time"

	"gorm.io/gorm"
)

// UserDefinedVariable 用户自定义变量
// 与元数据(case_groups.meta_*)完全独立，用于脚本参数化
// 支持两种模式：
// 1. 用例集变量：TaskUUID 为空，变量属于用例集
// 2. 任务变量：TaskUUID 非空，变量是从用例集复制的任务独立副本
type UserDefinedVariable struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProjectID uint   `gorm:"not null;index:idx_udv_project" json:"project_id"`
	GroupID   uint   `gorm:"not null;index:idx_udv_group" json:"group_id"`
	GroupType string `gorm:"type:varchar(20);not null;index:idx_udv_group" json:"group_type"` // web/api
	TaskUUID  string `gorm:"type:varchar(36);index:idx_udv_task" json:"task_uuid,omitempty"`  // 任务UUID（为空表示用例集变量）

	// 变量定义
	VarName  string `gorm:"type:varchar(100);not null" json:"var_name"`        // 变量名(如: ${BASE_URL})
	VarKey   string `gorm:"type:varchar(100);not null" json:"var_key"`         // 变量键名(如: base_url)
	VarDesc  string `gorm:"type:varchar(500)" json:"var_desc"`                 // 变量描述
	VarValue string `gorm:"type:text" json:"var_value"`                        // 变量值
	VarType  string `gorm:"type:varchar(20);default:'custom'" json:"var_type"` // 变量类型: metadata/path/body/custom

	// 审计字段
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UserDefinedVariable) TableName() string {
	return "user_defined_variables"
}

// BeforeCreate GORM钩子
func (v *UserDefinedVariable) BeforeCreate(tx *gorm.DB) error {
	if v.VarType == "" {
		v.VarType = "custom"
	}
	return nil
}
