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

// RequirementItemHandler 需求条目处理器接口
type RequirementItemHandler interface {
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

// requirementItemHandler 需求条目处理器实现
type requirementItemHandler struct {
	itemService    services.RequirementItemService
	projectService services.ProjectService
	storageDir     string
}

// NewRequirementItemHandler 创建需求条目处理器实例
func NewRequirementItemHandler(
	itemService services.RequirementItemService,
	projectService services.ProjectService,
	storageDir string,
) RequirementItemHandler {
	return &requirementItemHandler{
		itemService:    itemService,
		projectService: projectService,
		storageDir:     storageDir,
	}
}

// CreateItem 创建需求条目
// POST /api/projects/:id/requirement-items
func (h *requirementItemHandler) CreateItem(c *gin.Context) {
	projectID, err := h.validateProjectAccess(c)
	if err != nil {
		return
	}

	var req struct {
		Name    string           `json:"name" binding:"required"`
		Content string           `json:"content"`
		Chunks  []dto.ChunkInput `json:"chunks,omitempty"`
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

// UpdateItem 更新需求条目
// PUT /api/requirement-items/:id
// PUT /api/projects/:id/requirement-items/:itemId
func (h *requirementItemHandler) UpdateItem(c *gin.Context) {
	// 优先使用 itemId（项目级别路由）
	idStr := c.Param("itemId")
	if idStr == "" {
		// 如果没有 itemId，使用 id（兼容原来的路由）
		idStr = c.Param("id")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的需求条目ID")
		return
	}

	var req struct {
		Name    string               `json:"name"`
		Content string               `json:"content"`
		Chunks  []dto.ChunkOperation `json:"chunks,omitempty"`
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

// DeleteItem 删除需求条目
// DELETE /api/requirement-items/:id
// DELETE /api/projects/:id/requirement-items/:itemId
func (h *requirementItemHandler) DeleteItem(c *gin.Context) {
	// 优先使用 itemId（项目级别路由）
	idStr := c.Param("itemId")
	if idStr == "" {
		// 如果没有 itemId，使用 id（兼容原来的路由）
		idStr = c.Param("id")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的需求条目ID")
		return
	}

	if err := h.itemService.DeleteItem(uint(id)); err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "删除成功"})
}

// GetItem 获取单个需求条目（包含完整Chunks内容）
// GET /api/v1/projects/:id/requirement-items/:itemId
func (h *requirementItemHandler) GetItem(c *gin.Context) {
	idStr := c.Param("itemId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的需求条目ID")
		return
	}

	item, err := h.itemService.GetItemWithChunks(uint(id))
	if err != nil {
		utils.ResponseError(c, 404, err.Error())
		return
	}

	utils.ResponseSuccess(c, item)
}

// ListItems 获取项目的所有需求条目（包含Chunks摘要）
// GET /api/projects/:id/requirement-items
func (h *requirementItemHandler) ListItems(c *gin.Context) {
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

// BulkCreateItems 批量创建需求条目
// POST /api/projects/:id/requirement-items/bulk
func (h *requirementItemHandler) BulkCreateItems(c *gin.Context) {
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
	items := make([]struct{ Name, Content string }, len(req.Items))
	for i, item := range req.Items {
		items[i] = struct{ Name, Content string }{
			Name:    item.Name,
			Content: item.Content,
		}
	}

	if err := h.itemService.BulkCreateItems(projectID, items); err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "批量创建成功"})
}

// BulkUpdateItems 批量更新需求条目
// PUT /api/requirement-items/bulk
func (h *requirementItemHandler) BulkUpdateItems(c *gin.Context) {
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
		items[i] = struct {
			ID      uint
			Name    string
			Content string
		}{
			ID:      item.ID,
			Name:    item.Name,
			Content: item.Content,
		}
	}

	if err := h.itemService.BulkUpdateItems(items); err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "批量更新成功"})
}

// BulkDeleteItems 批量删除需求条目
// DELETE /api/requirement-items/bulk
func (h *requirementItemHandler) BulkDeleteItems(c *gin.Context) {
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
// POST /api/projects/:id/requirement-items/export
func (h *requirementItemHandler) ExportToZip(c *gin.Context) {
	fmt.Printf("[Handler.ExportToZip] 收到请求\n")

	projectID, err := h.validateProjectAccess(c)
	if err != nil {
		fmt.Printf("[Handler.ExportToZip] 验证项目访问失败\n")
		return
	}
	fmt.Printf("[Handler.ExportToZip] ProjectID=%d\n", projectID)

	userID := h.getUserID(c)
	fmt.Printf("[Handler.ExportToZip] UserID=%d\n", userID)

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)
	fmt.Printf("[Handler.ExportToZip] Remark=%s\n", req.Remark)

	// 获取项目信息
	project, _, err := h.projectService.GetByID(projectID, userID)
	if err != nil {
		utils.ResponseError(c, 400, "项目不存在")
		return
	}

	// 生成ZIP文件路径
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_Requirements_%s.zip", project.Name, timestamp)
	outputPath := filepath.Join(h.storageDir, "versions", filename)
	fmt.Printf("[Handler.ExportToZip] 输出路径=%s\n", outputPath)

	version, err := h.itemService.ExportToZip(projectID, outputPath, req.Remark, userID)
	if err != nil {
		fmt.Printf("[Handler.ExportToZip] 服务层返回错误: %v\n", err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	fmt.Printf("[Handler.ExportToZip] 成功导出版本, ID=%d, Filename=%s\n", version.ID, version.Filename)
	utils.ResponseSuccess(c, version)
}

// ImportFromZip 从ZIP批量版本导入
// POST /api/projects/:id/requirement-items/import
func (h *requirementItemHandler) ImportFromZip(c *gin.Context) {
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
func (h *requirementItemHandler) validateProjectAccess(c *gin.Context) (uint, error) {
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
func (h *requirementItemHandler) getUserID(c *gin.Context) uint {
	userIDVal, _ := c.Get("userID")
	return userIDVal.(uint)
}
