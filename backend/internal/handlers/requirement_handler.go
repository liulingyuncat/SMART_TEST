package handlers

import (
	"strconv"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequirementHandler 需求文档处理器接口
type RequirementHandler interface {
	GetRequirement(c *gin.Context)
	UpdateRequirement(c *gin.Context)
}

// requirementHandler 需求文档处理器实现
type requirementHandler struct {
	requirementService services.RequirementService
	projectService     services.ProjectService
}

// NewRequirementHandler 创建需求文档处理器实例
func NewRequirementHandler(requirementService services.RequirementService, projectService services.ProjectService) RequirementHandler {
	return &requirementHandler{
		requirementService: requirementService,
		projectService:     projectService,
	}
}

// GetRequirement 获取需求文档
func (h *requirementHandler) GetRequirement(c *gin.Context) {
	// 提取路径参数
	idStr := c.Param("id")
	docType := c.Param("type")

	// 验证项目ID
	projectID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的项目ID")
		return
	}

	// 验证用户权限(复用T07权限逻辑)
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 检查用户是否为项目成员
	_, role, err := h.projectService.GetByID(uint(projectID), userID)
	if err != nil {
		// 项目不存在或用户无权限
		utils.ResponseError(c, 403, "您没有权限访问此项目的需求文档")
		return
	}
	if role == "" {
		utils.ResponseError(c, 403, "您没有权限访问此项目的需求文档")
		return
	}

	// 获取需求文档内容
	content, updatedAt, err := h.requirementService.GetRequirement(uint(projectID), docType)
	if err != nil {
		utils.ResponseError(c, 400, err.Error())
		return
	}

	// 返回成功响应
	utils.ResponseSuccess(c, gin.H{
		"project_id": projectID,
		"doc_type":   docType,
		"content":    content,
		"updated_at": updatedAt,
	})
}

// UpdateRequirement 更新需求文档
func (h *requirementHandler) UpdateRequirement(c *gin.Context) {
	// 提取路径参数
	idStr := c.Param("id")
	docType := c.Param("type")

	// 验证项目ID
	projectID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的项目ID")
		return
	}

	// 验证用户权限
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 检查用户是否为项目成员
	_, role, err := h.projectService.GetByID(uint(projectID), userID)
	if err != nil {
		utils.ResponseError(c, 403, "您没有权限修改此项目的需求文档")
		return
	}
	if role == "" {
		utils.ResponseError(c, 403, "您没有权限修改此项目的需求文档")
		return
	}

	// 解析请求体
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误")
		return
	}

	// 更新需求文档
	if err := h.requirementService.UpdateRequirement(uint(projectID), docType, req.Content); err != nil {
		utils.ResponseError(c, 500, err.Error())
		return
	}

	// 返回成功响应
	utils.ResponseSuccess(c, gin.H{
		"message": "需求文档已更新",
	})
}
