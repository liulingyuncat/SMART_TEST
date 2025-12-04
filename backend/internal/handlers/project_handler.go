package handlers

import (
	"errors"
	"log"
	"strconv"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// ProjectHandler 项目处理器接口
type ProjectHandler interface {
	GetProjects(c *gin.Context)
	CreateProject(c *gin.Context)
	UpdateProject(c *gin.Context)
	DeleteProject(c *gin.Context)
	GetProjectByID(c *gin.Context)
}

// projectHandler 项目处理器实现
type projectHandler struct {
	projectService services.ProjectService
}

// NewProjectHandler 创建项目处理器实例
func NewProjectHandler(projectService services.ProjectService) ProjectHandler {
	return &projectHandler{
		projectService: projectService,
	}
}

// GetProjects 获取项目列表
func (h *projectHandler) GetProjects(c *gin.Context) {
	// 从Context获取userID和role(由AuthMiddleware注入)
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}

	roleVal, exists := c.Get("role")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}

	userID := userIDVal.(uint)
	role := roleVal.(string)

	// 调用Service层获取项目列表
	projects, err := h.projectService.GetUserProjects(userID, role)
	if err != nil {
		log.Printf("[Project List Failed] user_id=%d, error=%v", userID, err)
		utils.ResponseError(c, 500, "查询项目失败")
		return
	}

	log.Printf("[Project List] user_id=%d, role=%s, count=%d", userID, role, len(projects))
	utils.ResponseSuccess(c, projects)
}

// CreateProjectRequest 创建项目请求结构体
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// UpdateProjectRequest 更新项目请求结构体
type UpdateProjectRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Status      *string `json:"status" binding:"omitempty,oneof=pending in-progress completed"`
	OwnerID     *int    `json:"owner_id"`
}

// CreateProject 创建项目
func (h *projectHandler) CreateProject(c *gin.Context) {
	// 从Context获取userID(创建者ID)
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 参数绑定和验证
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "参数验证失败")
		return
	}

	// 调用Service层创建项目
	project, err := h.projectService.CreateProject(req.Name, req.Description, userID)
	if err != nil {
		// 业务错误: 项目名已存在
		if errors.Is(err, services.ErrProjectNameExists) {
			log.Printf("[Project Name Conflict] user_id=%d, name=%s", userID, req.Name)
			utils.ResponseError(c, 400, "项目名已存在")
			return
		}

		// 系统错误
		log.Printf("[Project Create Failed] user_id=%d, error=%v", userID, err)
		utils.ResponseError(c, 500, "创建项目失败")
		return
	}

	log.Printf("[Project Created] user_id=%d, project_id=%d, name=%s", userID, project.ID, project.Name)
	c.JSON(201, gin.H{
		"code":    0,
		"message": "项目创建成功",
		"data":    project,
	})
}

// handleUpdateError 处理更新项目错误
func handleUpdateError(c *gin.Context, err error) {
	if errors.Is(err, services.ErrProjectNotFound) {
		utils.ResponseError(c, 404, "项目不存在")
	} else if errors.Is(err, services.ErrProjectNameExists) {
		utils.ResponseError(c, 400, "项目名已存在")
	} else if errors.Is(err, services.ErrPermissionDenied) {
		utils.ResponseError(c, 403, "无权限访问此项目")
	} else {
		log.Printf("[Project Update Failed] error=%v", err)
		utils.ResponseError(c, 500, "更新项目失败")
	}
}

// handleDeleteError 处理删除项目错误
func handleDeleteError(c *gin.Context, err error) {
	if errors.Is(err, services.ErrProjectNotFound) {
		utils.ResponseError(c, 404, "项目不存在")
	} else if errors.Is(err, services.ErrPermissionDenied) {
		utils.ResponseError(c, 403, "无权限访问此项目")
	} else {
		log.Printf("[Project Delete Failed] error=%v", err)
		utils.ResponseError(c, 500, "删除项目失败")
	}
}

// UpdateProject 更新项目
func (h *projectHandler) UpdateProject(c *gin.Context) {
	// 解析项目ID
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的项目ID")
		return
	}

	// 获取用户信息
	userIDVal, _ := c.Get("userID")
	roleVal, _ := c.Get("role")
	userID := userIDVal.(uint)
	role := roleVal.(string)

	// 参数绑定和验证
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "参数验证失败")
		return
	}

	// 构建更新字段映射
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.OwnerID != nil {
		if *req.OwnerID == 0 {
			updates["owner_id"] = nil
		} else {
			updates["owner_id"] = *req.OwnerID
		}
	}

	// 调用Service层更新项目元数据
	project, err := h.projectService.UpdateProjectMetadata(uint(projectID), updates, userID, role)
	if err != nil {
		handleUpdateError(c, err)
		return
	}

	log.Printf("[Project Updated] user_id=%d, project_id=%d, updates=%v", userID, projectID, updates)
	utils.ResponseSuccess(c, project)
}

// DeleteProject 删除项目
func (h *projectHandler) DeleteProject(c *gin.Context) {
	// 解析项目ID
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的项目ID")
		return
	}

	// 获取用户信息
	userIDVal, _ := c.Get("userID")
	roleVal, _ := c.Get("role")
	userID := userIDVal.(uint)
	role := roleVal.(string)

	// 调用Service层删除项目
	err = h.projectService.DeleteProject(uint(projectID), userID, role)
	if err != nil {
		handleDeleteError(c, err)
		return
	}

	log.Printf("[Project Deleted] user_id=%d, project_id=%d", userID, projectID)
	utils.ResponseSuccess(c, gin.H{"message": "项目删除成功"})
}

// GetProjectByID 获取项目详情
func (h *projectHandler) GetProjectByID(c *gin.Context) {
	// 1. 提取并验证项目ID参数
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		log.Printf("[GetProjectByID] Invalid project ID: id=%s, error=%v", projectIDStr, err)
		utils.ResponseError(c, 400, "参数无效")
		return
	}

	// 2. 获取当前用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}
	userID := userIDVal.(uint)

	log.Printf("[GetProjectByID] Request: project_id=%d, user_id=%d", projectID, userID)

	// 3. 调用Service层查询项目
	project, userRole, err := h.projectService.GetByID(uint(projectID), userID)
	if err != nil {
		// 根据错误类型返回不同的HTTP状态码
		if errors.Is(err, services.ErrPermissionDenied) {
			log.Printf("[GetProjectByID] Permission denied: project_id=%d, user_id=%d", projectID, userID)
			utils.ResponseError(c, 403, "无权限访问此项目")
			return
		}
		if errors.Is(err, services.ErrProjectNotFound) {
			log.Printf("[GetProjectByID] Project not found: project_id=%d", projectID)
			utils.ResponseError(c, 404, "项目不存在")
			return
		}
		// 其他错误
		log.Printf("[GetProjectByID] Service error: project_id=%d, error=%v", projectID, err)
		utils.ResponseError(c, 500, "服务器内部错误")
		return
	}

	// 构造响应,包含项目信息和用户角色
	response := gin.H{
		"id":          project.ID,
		"name":        project.Name,
		"description": project.Description,
		"status":      project.Status,
		"owner_id":    project.OwnerID,
		"owner_name":  project.OwnerName,
		"created_at":  project.CreatedAt,
		"updated_at":  project.UpdatedAt,
		"user_role":   userRole,
	}

	log.Printf("[GetProjectByID] Success: project_id=%d, user_id=%d, role=%s", projectID, userID, userRole)
	utils.ResponseSuccess(c, response)
}
