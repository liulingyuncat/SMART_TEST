package models

import (
	"time"

	"gorm.io/gorm"
)

// DefectComment 缺陷说明模型
type DefectComment struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	DefectID  string         `gorm:"type:varchar(36);not null;index:idx_defect_comments_defect_id" json:"defect_id"` // 关联缺陷UUID
	Content   string         `gorm:"type:text;not null" json:"content"`                                              // 说明内容
	CreatedBy uint           `gorm:"not null" json:"created_by"`                                                     // 创建人ID
	UpdatedBy uint           `json:"updated_by"`                                                                     // 最后编辑人ID
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_defect_comments_deleted_at" json:"-"`

	// 关联 - 不存储在数据库
	CreatedByUser User `gorm:"foreignKey:CreatedBy;references:ID" json:"created_by_user,omitempty"` // 创建人信息
	UpdatedByUser User `gorm:"foreignKey:UpdatedBy;references:ID" json:"updated_by_user,omitempty"` // 编辑人信息
}

// TableName 指定表名
func (DefectComment) TableName() string {
	return "defect_comments"
}

// DefectCommentCreateRequest 创建说明请求
type DefectCommentCreateRequest struct {
	Content string `json:"content" binding:"required,max=2000"`
}

// DefectCommentUpdateRequest 更新说明请求
type DefectCommentUpdateRequest struct {
	Content string `json:"content" binding:"required,max=2000"`
}

// DefectCommentListResponse 说明列表响应
type DefectCommentListResponse struct {
	Comments []DefectComment `json:"comments"`
	Total    int64           `json:"total"`
}
