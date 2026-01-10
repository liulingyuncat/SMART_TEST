package models

import "time"

// CaseReviewItem 用例审阅条目模型
// 支持多文档审阅管理,每个项目可创建多个审阅文档
type CaseReviewItem struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID uint      `json:"project_id" gorm:"not null;index:idx_review_items_project"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (CaseReviewItem) TableName() string {
	return "case_review_items"
}
