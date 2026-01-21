package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"webtest/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ReviewItemHandler 审阅条目处理器
type ReviewItemHandler struct {
	reviewItemService services.ReviewItemService
	projectService    services.ProjectService
}

// NewReviewItemHandler 创建审阅条目处理器实例
func NewReviewItemHandler(reviewItemService services.ReviewItemService, projectService services.ProjectService) *ReviewItemHandler {
	return &ReviewItemHandler{
		reviewItemService: reviewItemService,
		projectService:    projectService,
	}
}

// ListReviewItems 获取项目所有审阅条目
// GET /api/v1/projects/:id/review-items
func (h *ReviewItemHandler) ListReviewItems(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	// 获取审阅列表
	items, err := h.reviewItemService.ListReviewItems(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取审阅列表失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// CreateReviewItem 创建审阅条目
// POST /api/v1/projects/:id/review-items
func (h *ReviewItemHandler) CreateReviewItem(c *gin.Context) {
	log.Printf("[CreateReviewItem] === Start ===")
	projectIDStr := c.Param("id")
	log.Printf("[CreateReviewItem] Raw project ID param: %s", projectIDStr)

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		log.Printf("[CreateReviewItem] Invalid project ID: %s, error: %v", projectIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}
	log.Printf("[CreateReviewItem] Parsed project ID: %d", projectID)

	// 读取原始请求体用于调试
	bodyBytes, _ := c.GetRawData()
	log.Printf("[CreateReviewItem] Raw request body: %s", string(bodyBytes))

	// 重新设置请求体供后续绑定使用
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateReviewItem] Bind JSON failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "details": err.Error()})
		return
	}

	log.Printf("[CreateReviewItem] ProjectID: %d, Name: '%s'", projectID, req.Name) // 创建审阅条目
	item, err := h.reviewItemService.CreateReviewItem(uint(projectID), req.Name)
	if err != nil {
		if err.Error() == "审阅名称已存在" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建审阅失败: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// GetReviewItem 获取审阅条目详情
// GET /api/v1/projects/:id/review-items/:itemId
func (h *ReviewItemHandler) GetReviewItem(c *gin.Context) {
	itemIDStr := c.Param("itemId")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审阅ID"})
		return
	}

	// 获取审阅详情
	item, err := h.reviewItemService.GetReviewItem(uint(itemID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "审阅记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取审阅失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateReviewItem 更新审阅条目
// PUT /api/v1/projects/:id/review-items/:itemId
func (h *ReviewItemHandler) UpdateReviewItem(c *gin.Context) {
	itemIDStr := c.Param("itemId")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审阅ID"})
		return
	}

	var req struct {
		Name    *string `json:"name"`
		Content *string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 更新审阅条目
	item, err := h.reviewItemService.UpdateReviewItem(uint(itemID), req.Name, req.Content)
	if err != nil {
		if err.Error() == "审阅名称已存在" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "审阅记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("更新审阅失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteReviewItem 删除审阅条目
// DELETE /api/v1/projects/:id/review-items/:itemId
func (h *ReviewItemHandler) DeleteReviewItem(c *gin.Context) {
	itemIDStr := c.Param("itemId")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审阅ID"})
		return
	}

	// 删除审阅条目
	if err := h.reviewItemService.DeleteReviewItem(uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除审阅失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// DownloadReviewItem 下载审阅文档
// GET /api/v1/projects/:id/review-items/:itemId/download
func (h *ReviewItemHandler) DownloadReviewItem(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	itemIDStr := c.Param("itemId")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审阅ID"})
		return
	}

	// 获取用户ID
	userID, _ := c.Get("userID")

	// 获取项目名称
	project, _, err := h.projectService.GetByID(uint(projectID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取项目信息失败"})
		return
	}

	// 生成Markdown文件
	content, filename, err := h.reviewItemService.DownloadReviewItem(uint(itemID), project.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "审阅记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("生成文件失败: %v", err)})
		return
	}

	// 返回文件流
	c.Header("Content-Type", "text/markdown; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.String(http.StatusOK, content)
}
