package dto

import "time"

// ViewpointChunkInput 创建观点Chunk输入
type ViewpointChunkInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

// ViewpointChunkOperation 观点Chunk操作（用于更新）
type ViewpointChunkOperation struct {
	ChunkID *uint  `json:"chunk_id,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Delete  bool   `json:"_delete,omitempty"`
}

// ViewpointChunkSummary 观点Chunk摘要（用于列表响应）
type ViewpointChunkSummary struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`
}

// ViewpointChunkDetail 观点Chunk详情（用于详情响应）
type ViewpointChunkDetail struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	SortOrder int    `json:"sort_order"`
}

// CreateViewpointItemRequest 创建观点请求
type CreateViewpointItemRequest struct {
	Name    string                `json:"name" binding:"required"`
	Content string                `json:"content"`
	Chunks  []ViewpointChunkInput `json:"chunks,omitempty"`
}

// UpdateViewpointItemRequest 更新观点请求
type UpdateViewpointItemRequest struct {
	Name    string                    `json:"name,omitempty"`
	Content string                    `json:"content,omitempty"`
	Chunks  []ViewpointChunkOperation `json:"chunks,omitempty"`
}

// ViewpointItemWithChunks 观点响应（带Chunks摘要，用于列表）
type ViewpointItemWithChunks struct {
	ID        uint                    `json:"id"`
	ProjectID uint                    `json:"project_id"`
	Name      string                  `json:"name"`
	Content   string                  `json:"content"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
	Chunks    []ViewpointChunkSummary `json:"chunks"`
}

// ViewpointItemWithChunkDetails 观点响应（带Chunks详情，用于详情）
type ViewpointItemWithChunkDetails struct {
	ID        uint                   `json:"id"`
	ProjectID uint                   `json:"project_id"`
	Name      string                 `json:"name"`
	Content   string                 `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Chunks    []ViewpointChunkDetail `json:"chunks"`
}
