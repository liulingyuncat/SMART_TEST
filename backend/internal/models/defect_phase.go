package models

import (
	"time"

	"gorm.io/gorm"
)

// DefectPhase 缺陷测试阶段配置模型
type DefectPhase struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID uint   `gorm:"not null;uniqueIndex:idx_defect_phases_project_name" json:"project_id"`             // 所属项目ID
	Name      string `gorm:"type:varchar(100);not null;uniqueIndex:idx_defect_phases_project_name" json:"name"` // Phase名称
	SortOrder int    `gorm:"default:0" json:"sort_order"`                                                       // 排序顺序

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_defect_phases_deleted_at" json:"-"`
}

// TableName 指定表名
func (DefectPhase) TableName() string {
	return "defect_phases"
}

// DefectPhaseCreateRequest 创建Phase请求
type DefectPhaseCreateRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	SortOrder int    `json:"sort_order"`
}

// DefectPhaseUpdateRequest 更新Phase请求
type DefectPhaseUpdateRequest struct {
	Name      *string `json:"name" binding:"omitempty,max=100"`
	SortOrder *int    `json:"sort_order"`
}
