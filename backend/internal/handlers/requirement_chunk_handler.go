package handlers

import (
	"strconv"
	"webtest/internal/repositories"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequirementChunkHandler 需求Chunk处理器接口
type RequirementChunkHandler interface {
	ListChunks(c *gin.Context)
	CreateChunk(c *gin.Context)
	GetChunk(c *gin.Context)
	UpdateChunk(c *gin.Context)
	DeleteChunk(c *gin.Context)
	ReorderChunks(c *gin.Context)
}

// requirementChunkHandler 需求Chunk处理器实现
type requirementChunkHandler struct {
	chunkService services.RequirementChunkService
}

// NewRequirementChunkHandler 创建需求Chunk处理器实例
func NewRequirementChunkHandler(chunkService services.RequirementChunkService) RequirementChunkHandler {
	return &requirementChunkHandler{
		chunkService: chunkService,
	}
}

// ListChunks 获取需求的所有Chunk列表
// GET /api/v1/projects/:id/requirement-items/:itemId/chunks
func (h *requirementChunkHandler) ListChunks(c *gin.Context) {
	idStr := c.Param("itemId")
	if idStr == "" {
		idStr = c.Param("id")
	}
	requirementID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的需求ID")
		return
	}

	chunks, err := h.chunkService.GetChunksByRequirementID(uint(requirementID))
	if err != nil {
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, chunks)
}

// CreateChunk 创建新Chunk
// POST /api/v1/projects/:id/requirement-items/:itemId/chunks
func (h *requirementChunkHandler) CreateChunk(c *gin.Context) {
	idStr := c.Param("itemId")
	if idStr == "" {
		idStr = c.Param("id")
	}
	requirementID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的需求ID")
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	chunk, err := h.chunkService.CreateChunk(uint(requirementID), req.Title, req.Content)
	if err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, chunk)
}

// GetChunk 获取单个Chunk详情
// GET /api/v1/requirement-chunks/:chunkId
func (h *requirementChunkHandler) GetChunk(c *gin.Context) {
	idStr := c.Param("chunkId")
	chunkID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的Chunk ID")
		return
	}

	chunk, err := h.chunkService.GetChunkByID(uint(chunkID))
	if err != nil {
		utils.ResponseError(c, 404, err.Error())
		return
	}

	utils.ResponseSuccess(c, chunk)
}

// UpdateChunk 更新Chunk内容
// PUT /api/v1/requirement-chunks/:chunkId
func (h *requirementChunkHandler) UpdateChunk(c *gin.Context) {
	idStr := c.Param("chunkId")
	chunkID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的Chunk ID")
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	chunk, err := h.chunkService.UpdateChunk(uint(chunkID), req.Title, req.Content)
	if err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, chunk)
}

// DeleteChunk 删除Chunk
// DELETE /api/v1/requirement-chunks/:chunkId
func (h *requirementChunkHandler) DeleteChunk(c *gin.Context) {
	idStr := c.Param("chunkId")
	chunkID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的Chunk ID")
		return
	}

	if err := h.chunkService.DeleteChunk(uint(chunkID)); err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "删除成功"})
}

// ReorderChunks 批量重排序Chunk
// PUT /api/v1/requirement-items/:id/chunks/reorder
func (h *requirementChunkHandler) ReorderChunks(c *gin.Context) {
	var req struct {
		ChunkOrders []repositories.ChunkOrder `json:"chunk_orders" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if err := h.chunkService.ReorderChunks(req.ChunkOrders); err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "排序成功"})
}
