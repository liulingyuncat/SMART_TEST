package models

import "time"

// CaseReview 测试用例评审记录模型
type CaseReview struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"not null;index:idx_reviews_project" json:"project_id"`
	CaseType  string    `gorm:"type:varchar(20);not null" json:"case_type"` // ai/overall/change
	Content   string    `gorm:"type:text" json:"content"`                   // Markdown原文
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (CaseReview) TableName() string {
	return "test_case_reviews"
}
