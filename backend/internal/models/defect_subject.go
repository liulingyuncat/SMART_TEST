package models

import (
	"time"

	"gorm.io/gorm"
)

// DefectSubject 缺陷主题分类配置模型
type DefectSubject struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID uint   `gorm:"not null;uniqueIndex:idx_defect_subjects_project_name" json:"project_id"`             // 所属项目ID
	Name      string `gorm:"type:varchar(100);not null;uniqueIndex:idx_defect_subjects_project_name" json:"name"` // Subject名称
	SortOrder int    `gorm:"default:0" json:"sort_order"`                                                         // 排序顺序

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_defect_subjects_deleted_at" json:"-"`
}

// TableName 指定表名
func (DefectSubject) TableName() string {
	return "defect_subjects"
}

// DefectSubjectCreateRequest 创建Subject请求
type DefectSubjectCreateRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	SortOrder int    `json:"sort_order"`
}

// DefectSubjectUpdateRequest 更新Subject请求
type DefectSubjectUpdateRequest struct {
	Name      *string `json:"name" binding:"omitempty,max=100"`
	SortOrder *int    `json:"sort_order"`
}
