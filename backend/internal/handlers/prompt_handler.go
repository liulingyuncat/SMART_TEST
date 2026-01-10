package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"webtest/internal/constants"
	"webtest/internal/models"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// PromptHandler 提示词处理器接口
type PromptHandler interface {
	ListPrompts(c *gin.Context)
	GetPromptByID(c *gin.Context)
	GetPromptByName(c *gin.Context)
	CreatePrompt(c *gin.Context)
	UpdatePrompt(c *gin.Context)
	DeletePrompt(c *gin.Context)
}

// promptHandler 提示词处理器实现
type promptHandler struct {
	promptService services.PromptService
}

// NewPromptHandler 创建提示词处理器实例
func NewPromptHandler(promptService services.PromptService) PromptHandler {
	return &promptHandler{
		promptService: promptService,
	}
}

// ListPromptsRequest 查询提示词列表请求
type ListPromptsRequest struct {
	ProjectID uint   `form:"project_id"` // 移除required，提示词与项目无关
	Scope     string `form:"scope"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

// ListPrompts 获取提示词列表
func (h *promptHandler) ListPrompts(c *gin.Context) {
	var userID uint
	var userRole string

	// 获取认证信息（可能不存在，对于公开路由）
	userIDVal, userExists := c.Get("userID")
	if userExists {
		userID = userIDVal.(uint)
		log.Printf("[ListPrompts] Authenticated userID from context: %d", userID)
	} else {
		log.Printf("[ListPrompts] No userID in context (unauthenticated request)")
	}

	// 获取用户角色
	if role, exists := c.Get("role"); exists {
		userRole = role.(string)
		log.Printf("[ListPrompts] User role from context: %s", userRole)
	}

	var req ListPromptsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ResponseError(c, 400, "参数验证失败")
		return
	}

	log.Printf("[ListPrompts] Request params: scope=%s, page=%d, page_size=%d", req.Scope, req.Page, req.PageSize)

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	// 如果project_id为0，使用1作为默认值（提示词与具体项目无关）
	if req.ProjectID == 0 {
		req.ProjectID = 1
	}

	// 权限检查：
	// - 如果scope=project（获取全员提示词），允许无认证访问
	// - 如果scope=user（获取个人提示词），必须有认证
	// - 其他情况下需要特定角色
	if req.Scope == "user" {
		// 个人提示词需要认证
		if !userExists {
			log.Printf("[ListPrompts] Rejecting user scope request: no authentication")
			utils.ResponseError(c, 401, "未授权")
			return
		}
		log.Printf("[ListPrompts] User scope request: userID=%d", userID)
	} else if req.Scope != "project" {
		// 检查角色权限（仅在非全员提示词查询时）
		if userExists && userRole != "" {
			hasPermission := false
			allowedRoles := []string{constants.RoleSystemAdmin, constants.RoleProjectManager, constants.RoleProjectMember}
			for _, role := range allowedRoles {
				if userRole == role {
					hasPermission = true
					break
				}
			}
			if !hasPermission {
				utils.ResponseError(c, 403, "权限不足")
				return
			}
		}
	}

	prompts, total, err := h.promptService.ListPromptsWithRole(
		req.ProjectID, req.Scope, userID, userRole, req.Page, req.PageSize,
	)
	if err != nil {
		log.Printf("[Prompt List Failed] user_id=%d, error=%v", userID, err)
		utils.ResponseError(c, 500, "查询提示词失败")
		return
	}

	log.Printf("[Prompt List] user_id=%d, scope=%s, count=%d, total=%d", userID, req.Scope, len(prompts), total)
	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"items": prompts,
			"total": total,
		},
	})
}

// GetPromptByID 获取提示词详情
func (h *promptHandler) GetPromptByID(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}
	userID := userIDVal.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ResponseError(c, 400, "无效的ID")
		return
	}

	prompt, err := h.promptService.GetPromptByID(uint(id), userID)
	if err != nil {
		if errors.Is(err, services.ErrPromptNotFound) {
			utils.ResponseError(c, 404, "提示词不存在")
			return
		}
		if errors.Is(err, services.ErrPromptPermissionDenied) {
			utils.ResponseError(c, 403, "无权限访问此提示词")
			return
		}
		log.Printf("[Prompt Get Failed] user_id=%d, id=%d, error=%v", userID, id, err)
		utils.ResponseError(c, 500, "获取提示词失败")
		return
	}

	log.Printf("[Prompt Get] user_id=%d, id=%d", userID, id)
	utils.ResponseSuccess(c, prompt)
}

// GetPromptByName 通过名称获取提示词详情（用于MCP）
func (h *promptHandler) GetPromptByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		utils.ResponseError(c, 400, "缺少name参数")
		return
	}

	userIDVal, _ := c.Get("userID")
	userID := uint(0)
	if userIDVal != nil {
		userID = userIDVal.(uint)
	}

	log.Printf("[Prompt Get By Name] name=%s, user_id=%d", name, userID)

	prompt, err := h.promptService.GetPromptByName(name, userID)
	if err != nil {
		if errors.Is(err, services.ErrPromptNotFound) {
			log.Printf("[Prompt Get By Name Failed] name=%s, user_id=%d, error=not found", name, userID)
			utils.ResponseError(c, 404, "提示词不存在")
			return
		}
		if errors.Is(err, services.ErrPromptPermissionDenied) {
			log.Printf("[Prompt Get By Name Failed] name=%s, user_id=%d, error=permission denied", name, userID)
			utils.ResponseError(c, 403, "无权限访问此提示词")
			return
		}
		log.Printf("[Prompt Get By Name Failed] name=%s, user_id=%d, error=%v", name, userID, err)
		utils.ResponseError(c, 500, "获取提示词失败")
		return
	}

	log.Printf("[Prompt Get By Name Success] name=%s, user_id=%d, scope=%s, prompt_id=%d", name, userID, prompt.Scope, prompt.ID)
	utils.ResponseSuccess(c, prompt)
}

// CreatePromptRequest 创建提示词请求
type CreatePromptRequest struct {
	ProjectID   uint   `json:"project_id"` // 全员提示词不需要，可以为0
	Name        string `json:"name" binding:"required,min=3,max=50"`
	Description string `json:"description" binding:"max=200"`
	Version     string `json:"version" binding:"required"`
	Content     string `json:"content" binding:"required,max=10000"`
	Arguments   string `json:"arguments"` // 接收JSON字符串
	Scope       string `json:"scope"`     // 可选，由后端根据用户角色设置
}

// CreatePrompt 创建提示词
func (h *promptHandler) CreatePrompt(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取用户角色
	userRole := ""
	if role, exists := c.Get("role"); exists {
		userRole = role.(string)
	}

	var req CreatePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "参数验证失败")
		return
	}

	// 解析 arguments JSON 字符串
	var arguments []models.PromptArgument
	if req.Arguments != "" {
		if err := json.Unmarshal([]byte(req.Arguments), &arguments); err != nil {
			utils.ResponseError(c, 400, "arguments 格式错误")
			return
		}
	}

	// 根据用户角色设置 scope（忽略前端发送的 scope）
	scope := "user" // 默认为个人
	if userRole == "system_admin" {
		scope = "project" // 系统管理员创建的是全员提示词
	}

	// 如果project_id为0，使用1作为默认值（提示词与具体项目无关）
	projectID := req.ProjectID
	if projectID == 0 {
		projectID = 1
	}

	prompt, err := h.promptService.CreatePrompt(
		projectID, req.Name, req.Description, req.Version,
		req.Content, scope, arguments, userID, userRole,
	)
	if err != nil {
		if errors.Is(err, services.ErrPromptNameExists) {
			utils.ResponseError(c, 409, "提示词名称已存在")
			return
		}
		if errors.Is(err, services.ErrPromptPermissionDenied) {
			utils.ResponseError(c, 403, "无权限创建该类型的提示词")
			return
		}
		log.Printf("[Prompt Create Failed] user_id=%d, error=%v", userID, err)
		utils.ResponseError(c, 500, "创建提示词失败")
		return
	}

	log.Printf("[Prompt Created] user_id=%d, id=%d, name=%s", userID, prompt.ID, prompt.Name)
	c.JSON(201, gin.H{
		"code":    0,
		"message": "提示词创建成功",
		"data":    prompt,
	})
}

// UpdatePromptRequest 更新提示词请求
type UpdatePromptRequest struct {
	Description *string `json:"description" binding:"omitempty,max=200"`
	Version     *string `json:"version"`
	Content     *string `json:"content" binding:"omitempty,max=10000"`
	Arguments   string  `json:"arguments"` // 接收JSON字符串
}

// UpdatePrompt 更新提示词
func (h *promptHandler) UpdatePrompt(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取用户角色
	userRole := ""
	if role, exists := c.Get("role"); exists {
		userRole = role.(string)
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ResponseError(c, 400, "无效的ID")
		return
	}

	var req UpdatePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "参数验证失败")
		return
	}

	// 解析 arguments JSON 字符串
	var arguments []models.PromptArgument
	if req.Arguments != "" {
		if err := json.Unmarshal([]byte(req.Arguments), &arguments); err != nil {
			utils.ResponseError(c, 400, "arguments 格式错误")
			return
		}
	}

	prompt, err := h.promptService.UpdatePrompt(
		uint(id), req.Description, req.Version, req.Content, arguments, userID, userRole,
	)
	if err != nil {
		if errors.Is(err, services.ErrPromptNotFound) {
			utils.ResponseError(c, 404, "提示词不存在")
			return
		}
		if errors.Is(err, services.ErrCannotModifySystem) {
			utils.ResponseError(c, 403, "不能修改系统提示词")
			return
		}
		if errors.Is(err, services.ErrPromptPermissionDenied) {
			utils.ResponseError(c, 403, "无权限操作此提示词")
			return
		}
		log.Printf("[Prompt Update Failed] user_id=%d, id=%d, error=%v", userID, id, err)
		utils.ResponseError(c, 500, "更新提示词失败")
		return
	}

	log.Printf("[Prompt Updated] user_id=%d, id=%d", userID, id)
	utils.ResponseSuccess(c, prompt)
}

// DeletePrompt 删除提示词
func (h *promptHandler) DeletePrompt(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取用户角色
	userRole := ""
	if role, exists := c.Get("role"); exists {
		userRole = role.(string)
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ResponseError(c, 400, "无效的ID")
		return
	}

	err = h.promptService.DeletePrompt(uint(id), userID, userRole)
	if err != nil {
		if errors.Is(err, services.ErrPromptNotFound) {
			utils.ResponseError(c, 404, "提示词不存在")
			return
		}
		if errors.Is(err, services.ErrCannotModifySystem) {
			utils.ResponseError(c, 403, "不能删除系统提示词")
			return
		}
		if errors.Is(err, services.ErrPromptPermissionDenied) {
			utils.ResponseError(c, 403, "无权限操作此提示词")
			return
		}
		log.Printf("[Prompt Delete Failed] user_id=%d, id=%d, error=%v", userID, id, err)
		utils.ResponseError(c, 500, "删除提示词失败")
		return
	}

	log.Printf("[Prompt Deleted] user_id=%d, id=%d", userID, id)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "提示词删除成功",
	})
}

// RefreshPromptsRequest 刷新提示词缓存请求
type RefreshPromptsRequest struct {
	Scope string `json:"scope"` // 可选: system/project/user/all
}

// RefreshPrompts 刷新MCP提示词缓存（热更新）
// 当系统、全员或个人提示词发生变化时调用此接口
// 前端通过此接口通知后端更新MCP层的提示词缓存
func (h *promptHandler) RefreshPrompts(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取用户角色 - 仅允许管理员和项目经理刷新
	userRole := ""
	if role, exists := c.Get("role"); exists {
		userRole = role.(string)
	}

	// 权限检查 - 仅允许系统管理员/项目经理刷新提示词缓存
	if userRole != "SystemAdmin" && userRole != "ProjectManager" {
		utils.ResponseError(c, 403, "无权限刷新提示词缓存")
		return
	}

	var req RefreshPromptsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "参数验证失败")
		return
	}

	// 刷新系统提示词：动态扫描prompts目录并更新数据库
	if req.Scope == "system" || req.Scope == "" || req.Scope == "all" {
		promptsDir := "internal/mcp/prompts"
		if err := h.promptService.RefreshSystemPromptsFromDirectory(promptsDir); err != nil {
			log.Printf("[Prompt Refresh Failed] error=%v", err)
			utils.ResponseError(c, 500, "刷新系统提示词失败")
			return
		}
	}

	log.Printf("[Prompt Refresh] user_id=%d, scope=%s", userID, req.Scope)

	c.JSON(200, gin.H{
		"code":    0,
		"message": "提示词缓存刷新成功",
	})
}
