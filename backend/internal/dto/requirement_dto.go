package dto

import "time"

// ChunkInput 创建Chunk输入
type ChunkInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

// ChunkOperation Chunk操作（用于更新）
type ChunkOperation struct {
	ChunkID *uint  `json:"chunk_id,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Delete  bool   `json:"_delete,omitempty"`
}

// ChunkSummary Chunk摘要（用于列表响应）
type ChunkSummary struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`
}

// ChunkDetail Chunk详情（用于详情响应）
type ChunkDetail struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	SortOrder int    `json:"sort_order"`
}

// CreateRequirementItemRequest 创建需求请求
type CreateRequirementItemRequest struct {
	Name    string       `json:"name" binding:"required"`
	Content string       `json:"content"`
	Chunks  []ChunkInput `json:"chunks,omitempty"`
}

// UpdateRequirementItemRequest 更新需求请求
type UpdateRequirementItemRequest struct {
	Name    string           `json:"name,omitempty"`
	Content string           `json:"content,omitempty"`
	Chunks  []ChunkOperation `json:"chunks,omitempty"`
}

// RequirementItemWithChunks 需求响应（带Chunks摘要，用于列表）
type RequirementItemWithChunks struct {
	ID        uint           `json:"id"`
	ProjectID uint           `json:"project_id"`
	Name      string         `json:"name"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Chunks    []ChunkSummary `json:"chunks"`
}

// RequirementItemWithChunkDetails 需求响应（带Chunks详情，用于详情）
type RequirementItemWithChunkDetails struct {
	ID        uint          `json:"id"`
	ProjectID uint          `json:"project_id"`
	Name      string        `json:"name"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Chunks    []ChunkDetail `json:"chunks"`
}
