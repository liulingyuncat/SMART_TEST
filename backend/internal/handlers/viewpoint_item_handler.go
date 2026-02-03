package handlers

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"
	"webtest/internal/dto"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// ViewpointItemHandler AI观点条目处理器接口
type ViewpointItemHandler interface {
	CreateItem(c *gin.Context)
	UpdateItem(c *gin.Context)
	DeleteItem(c *gin.Context)
	GetItem(c *gin.Context)
	ListItems(c *gin.Context)
	BulkCreateItems(c *gin.Context)
	BulkUpdateItems(c *gin.Context)
	BulkDeleteItems(c *gin.Context)
	ExportToZip(c *gin.Context)
	ImportFromZip(c *gin.Context)
}

// viewpointItemHandler AI观点条目处理器实现
type viewpointItemHandler struct {
	itemService    services.ViewpointItemService
	projectService services.ProjectService
	storageDir     string
}

// NewViewpointItemHandler 创建AI观点条目处理器实例
func NewViewpointItemHandler(
	itemService services.ViewpointItemService,
	projectService services.ProjectService,
	storageDir string,
) ViewpointItemHandler {
	return &viewpointItemHandler{
		itemService:    itemService,
		projectService: projectService,
		storageDir:     storageDir,
	}
}

// CreateItem 创建AI观点条目
// POST /api/projects/:id/viewpoint-items
func (h *viewpointItemHandler) CreateItem(c *gin.Context) {
	projectID, err := h.validateProjectAccess(c)
	if err != nil {
		return
	}

	var req struct {
		Name          string                    `json:"name" binding:"required"`
		Content       string                    `json:"content"`
		RequirementID *uint                     `json:"requirement_id,omitempty"`
		Chunks        []dto.ViewpointChunkInput `json:"chunks,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 如果有chunks，使用带Chunk的创建方法
	if len(req.Chunks) > 0 {
		result, err := h.itemService.CreateItemWithChunks(projectID, req.Name, req.Content, req.Chunks)
		if err != nil {
			utils.ResponseError(c, 400, err.Error())
			return
		}
		utils.ResponseSuccess(c, result)
		return
	}

	// 向后兼容：无chunks时使用原有方法
	item, err := h.itemService.CreateItem(projectID, req.Name, req.Content)
	if err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, item)
}

// UpdateItem 更新AI观点条目
// PUT /api/viewpoint-items/:id
// PUT /api/projects/:id/viewpoint-items/:itemId
func (h *viewpointItemHandler) UpdateItem(c *gin.Context) {
	// 优先使用 itemId（项目级别路由）
	idStr := c.Param("itemId")
	if idStr == "" {
		// 如果没有 itemId，使用 id（兼容原来的路由）
		idStr = c.Param("id")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的观点条目ID")
		return
	}

	var req struct {
		Name    string                        `json:"name"`
		Content string                        `json:"content"`
		Chunks  []dto.ViewpointChunkOperation `json:"chunks,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 如果有chunks操作，使用带Chunk的更新方法
	if len(req.Chunks) > 0 {
		var namePtr, contentPtr *string
		if req.Name != "" {
			namePtr = &req.Name
		}
		if req.Content != "" {
			contentPtr = &req.Content
		}
		result, err := h.itemService.UpdateItemWithChunks(uint(id), namePtr, contentPtr, req.Chunks)
		if err != nil {
			utils.ResponseError(c, 400, err.Error())
			return
		}
		utils.ResponseSuccess(c, result)
		return
	}

	// 向后兼容：无chunks时使用原有方法
	item, err := h.itemService.UpdateItem(uint(id), req.Name, req.Content)
	if err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, item)
}

// DeleteItem 删除AI观点条目
// DELETE /api/viewpoint-items/:id
// DELETE /api/projects/:id/viewpoint-items/:itemId
func (h *viewpointItemHandler) DeleteItem(c *gin.Context) {
	// 优先使用 itemId（项目级别路由）
	idStr := c.Param("itemId")
	if idStr == "" {
		// 如果没有 itemId，使用 id（兼容原来的路由）
		idStr = c.Param("id")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的观点条目ID")
		return
	}

	if err := h.itemService.DeleteItem(uint(id)); err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "删除成功"})
}

// GetItem 获取单个AI观点条目（包含完整Chunks内容）
// GET /api/viewpoint-items/:id
// GET /api/projects/:id/viewpoint-items/:itemId
func (h *viewpointItemHandler) GetItem(c *gin.Context) {
	// 优先使用 itemId（项目级别路由）
	itemIDStr := c.Param("itemId")
	if itemIDStr == "" {
		// 如果没有 itemId，使用 id（兼容原来的路由）
		itemIDStr = c.Param("id")
	}

	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的观点条目ID")
		return
	}

	item, err := h.itemService.GetItemWithChunks(uint(itemID))
	if err != nil {
		utils.ResponseError(c, 404, err.Error())
		return
	}

	// 如果是项目级别的路由，验证观点项属于该项目
	projectIDStr := c.Param("id")
	if projectIDStr != "" && itemIDStr != projectIDStr {
		projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err == nil && item.ProjectID != uint(projectID) {
			utils.ResponseError(c, 404, "观点条目不属于此项目")
			return
		}
	}

	utils.ResponseSuccess(c, item)
}

// ListItems 获取项目的所有AI观点条目（包含Chunks摘要）
// GET /api/projects/:id/viewpoint-items
func (h *viewpointItemHandler) ListItems(c *gin.Context) {
	projectID, err := h.validateProjectAccess(c)
	if err != nil {
		return
	}

	items, err := h.itemService.GetItemsWithChunksSummary(projectID)
	if err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, items)
}

// BulkCreateItems 批量创建AI观点条目
// POST /api/projects/:id/viewpoint-items/bulk
func (h *viewpointItemHandler) BulkCreateItems(c *gin.Context) {
	projectID, err := h.validateProjectAccess(c)
	if err != nil {
		return
	}

	var req struct {
		Items []struct {
			Name    string `json:"name" binding:"required"`
			Content string `json:"content"`
		} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 转换为Service所需的类型
	items := make([]struct {
		Name    string
		Content string
	}, len(req.Items))
	for i, item := range req.Items {
		items[i].Name = item.Name
		items[i].Content = item.Content
	}

	if err := h.itemService.BulkCreateItems(projectID, items); err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "批量创建成功"})
}

// BulkUpdateItems 批量更新AI观点条目
// PUT /api/viewpoint-items/bulk
func (h *viewpointItemHandler) BulkUpdateItems(c *gin.Context) {
	var req struct {
		Items []struct {
			ID      uint   `json:"id" binding:"required"`
			Name    string `json:"name" binding:"required"`
			Content string `json:"content"`
		} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 转换为Service所需的类型
	items := make([]struct {
		ID      uint
		Name    string
		Content string
	}, len(req.Items))
	for i, item := range req.Items {
		items[i].ID = item.ID
		items[i].Name = item.Name
		items[i].Content = item.Content
	}

	if err := h.itemService.BulkUpdateItems(items); err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "批量更新成功"})
}

// BulkDeleteItems 批量删除AI观点条目
// DELETE /api/viewpoint-items/bulk
func (h *viewpointItemHandler) BulkDeleteItems(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if err := h.itemService.BulkDeleteItems(req.IDs); err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "批量删除成功"})
}

// ExportToZip 导出为ZIP批量版本
// POST /api/projects/:id/viewpoint-items/export
func (h *viewpointItemHandler) ExportToZip(c *gin.Context) {
	projectID, err := h.validateProjectAccess(c)
	if err != nil {
		return
	}

	userID := h.getUserID(c)

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	// 获取项目信息
	project, _, err := h.projectService.GetByID(projectID, userID)
	if err != nil {
		utils.ResponseError(c, 400, "项目不存在")
		return
	}

	// 生成ZIP文件路径
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_Viewpoints_%s.zip", project.Name, timestamp)
	outputPath := filepath.Join(h.storageDir, "versions", filename)

	version, err := h.itemService.ExportToZip(projectID, outputPath, req.Remark, userID)
	if err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, version)
}

// ImportFromZip 从ZIP批量版本导入
// POST /api/projects/:id/viewpoint-items/import
func (h *viewpointItemHandler) ImportFromZip(c *gin.Context) {
	projectID, err := h.validateProjectAccess(c)
	if err != nil {
		return
	}

	userID := h.getUserID(c)

	// 处理文件上传
	file, err := c.FormFile("file")
	if err != nil {
		utils.ResponseError(c, 400, "文件上传失败: "+err.Error())
		return
	}

	// 保存临时文件
	tmpPath := filepath.Join(h.storageDir, "temp", file.Filename)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		utils.ResponseError(c, 500, "保存文件失败: "+err.Error())
		return
	}

	// 导入
	if err := h.itemService.ImportFromZip(projectID, tmpPath, userID); err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "导入成功"})
}

// validateProjectAccess 验证项目访问权限
func (h *viewpointItemHandler) validateProjectAccess(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	projectID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的项目ID")
		return 0, err
	}

	userID := h.getUserID(c)

	// 检查用户是否为项目成员
	_, role, err := h.projectService.GetByID(uint(projectID), userID)
	if err != nil || role == "" {
		utils.ResponseError(c, 403, "您没有权限访问此项目")
		return 0, fmt.Errorf("无权限")
	}

	return uint(projectID), nil
}

// getUserID 获取当前用户ID
func (h *viewpointItemHandler) getUserID(c *gin.Context) uint {
	userIDVal, _ := c.Get("userID")
	return userIDVal.(uint)
}
