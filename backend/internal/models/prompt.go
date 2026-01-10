package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Prompt 提示词数据模型
type Prompt struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProjectID   uint           `gorm:"not null;index:idx_project_scope" json:"project_id"`
	Name        string         `gorm:"not null;size:50" json:"name"`
	Description string         `gorm:"size:200" json:"description"`
	Version     string         `gorm:"not null;default:1.0" json:"version"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	Arguments   string         `gorm:"type:json" json:"arguments"` // JSON字符串
	Scope       string         `gorm:"not null;type:varchar(10);check:scope IN ('system','project','user')" json:"scope"`
	UserID      *uint          `gorm:"index:idx_user" json:"user_id"`
	CreatedBy   uint           `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联（使用gorm:"-"忽略这些字段，避免迁移问题）
	Project Project `gorm:"-" json:"-"`
	User    *User   `gorm:"-" json:"-"`
	Creator User    `gorm:"-" json:"-"`
}

// TableName 指定表名
func (Prompt) TableName() string {
	return "prompts"
}

// PromptArgument 提示词参数定义
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// PromptDTO 提示词数据传输对象
type PromptDTO struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Version     string           `json:"version"`
	Content     string           `json:"content,omitempty"` // list时不返回
	Arguments   []PromptArgument `json:"arguments"`
	Scope       string           `json:"scope"`
	UserID      *uint            `json:"user_id,omitempty"`
	CreatedBy   uint             `json:"created_by"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}

// ToDTO 转换为DTO
func (p *Prompt) ToDTO(includeContent bool) (*PromptDTO, error) {
	dto := &PromptDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Version:     p.Version,
		Scope:       p.Scope,
		UserID:      p.UserID,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}

	if includeContent {
		dto.Content = p.Content
	}

	// 解析arguments JSON
	if p.Arguments != "" {
		var args []PromptArgument
		if err := json.Unmarshal([]byte(p.Arguments), &args); err != nil {
			return nil, err
		}
		dto.Arguments = args
	} else {
		dto.Arguments = []PromptArgument{}
	}

	return dto, nil
}

// FromCreateRequest 从创建请求构建模型
func (p *Prompt) FromCreateRequest(
	projectID uint,
	name, description, version, content, scope string,
	arguments []PromptArgument,
	userID *uint,
	createdBy uint,
) error {
	p.ProjectID = projectID
	p.Name = name
	p.Description = description
	p.Version = version
	p.Content = content
	p.Scope = scope
	p.UserID = userID
	p.CreatedBy = createdBy

	// 序列化arguments为JSON
	if len(arguments) > 0 {
		argsJSON, err := json.Marshal(arguments)
		if err != nil {
			return err
		}
		p.Arguments = string(argsJSON)
	}

	return nil
}
