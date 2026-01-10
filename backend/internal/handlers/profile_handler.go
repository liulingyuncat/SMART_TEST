package handlers

import (
	"webtest/internal/constants"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// ProfileInfo 个人信息响应结构
type ProfileInfo struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

// UpdateProfileNicknameRequest 更新昵称请求
type UpdateProfileNicknameRequest struct {
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
}

// UpdatePasswordRequest 更新密码请求 - T23
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=50"`
}

// TokenResponse Token响应 - T23
type TokenResponse struct {
	Token string `json:"token,omitempty"`
}

// TokenStatusResponse Token状态响应 - T23
type TokenStatusResponse struct {
	HasToken bool `json:"has_token"`
}

// ProfileHandler 个人信息处理器
type ProfileHandler struct {
	userService    services.UserService
	projectService services.ProjectService
}

// NewProfileHandler 创建个人信息处理器实例
func NewProfileHandler(userService services.UserService) *ProfileHandler {
	return &ProfileHandler{userService: userService}
}

// NewProfileHandlerWithProject 创建个人信息处理器实例(包含项目服务)
func NewProfileHandlerWithProject(userService services.UserService, projectService services.ProjectService) *ProfileHandler {
	return &ProfileHandler{
		userService:    userService,
		projectService: projectService,
	}
}

// GetProfile 获取当前用户的个人信息
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// 从 Context 获取 userID（由 AuthMiddleware 注入）
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		utils.ResponseError(c, 500, "invalid user id type")
		return
	}

	// 调用 UserService 查询用户信息
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ResponseError(c, 404, "user not found")
			return
		}
		utils.ResponseError(c, 500, "failed to get user info: "+err.Error())
		return
	}

	// 返回用户信息
	utils.ResponseSuccess(c, ProfileInfo{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Role:     user.Role,
	})
}

// UpdateNickname 更新当前用户的昵称
func (h *ProfileHandler) UpdateNickname(c *gin.Context) {
	// 从 Context 获取 userID（由 AuthMiddleware 注入）
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		utils.ResponseError(c, 500, "invalid user id type")
		return
	}

	// 解析请求体
	var req UpdateProfileNicknameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "invalid request: "+err.Error())
		return
	}

	// 调用 UserService 更新昵称（复用现有逻辑）
	user, err := h.userService.UpdateNickname(userID, req.Nickname)
	if err != nil {
		if err == services.ErrNicknameExists {
			utils.ResponseError(c, 400, "nickname already exists")
			return
		}
		if err == services.ErrUserNotFound {
			utils.ResponseError(c, 404, "user not found")
			return
		}
		utils.ResponseError(c, 500, "failed to update nickname: "+err.Error())
		return
	}

	// 返回更新后的用户信息
	utils.ResponseSuccess(c, ProfileInfo{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Role:     user.Role,
	})
}

// UpdatePassword 更新当前用户的密码 - T23 密码修改功能
func (h *ProfileHandler) UpdatePassword(c *gin.Context) {
	// 从 Context 获取 userID（由 AuthMiddleware 注入）
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		utils.ResponseError(c, 500, "invalid user id type")
		return
	}

	// 解析请求体
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "invalid request: "+err.Error())
		return
	}

	// 调用 UserService 修改密码
	err := h.userService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if err == services.ErrCurrentPasswordIncorrect {
			utils.ResponseError(c, 400, "current password is incorrect")
			return
		}
		if err == services.ErrNewPasswordSameAsCurrent {
			utils.ResponseError(c, 400, "new password cannot be same as current")
			return
		}
		if err == services.ErrUserNotFound {
			utils.ResponseError(c, 404, "user not found")
			return
		}
		utils.ResponseError(c, 500, "failed to update password: "+err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "password updated successfully"})
}

// GenerateToken 生成API Token - T23 Token生成功能
func (h *ProfileHandler) GenerateToken(c *gin.Context) {
	// 从 Context 获取用户信息
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		utils.ResponseError(c, 500, "invalid user id type")
		return
	}

	// 检查用户角色，系统管理员不能生成Token
	roleValue, _ := c.Get("role")
	role, _ := roleValue.(string)
	if role == constants.RoleSystemAdmin {
		utils.ResponseError(c, 403, "system admin cannot generate api token")
		return
	}

	// 调用 UserService 生成Token
	token, err := h.userService.GenerateApiToken(userID)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ResponseError(c, 404, "user not found")
			return
		}
		utils.ResponseError(c, 500, "failed to generate token: "+err.Error())
		return
	}

	utils.ResponseSuccess(c, TokenResponse{Token: token})
}

// GetToken 获取Token状态 - T23 Token状态查询
func (h *ProfileHandler) GetToken(c *gin.Context) {
	// 从 Context 获取 userID
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		utils.ResponseError(c, 500, "invalid user id type")
		return
	}

	// 调用 UserService 检查Token状态
	hasToken, err := h.userService.HasApiToken(userID)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ResponseError(c, 404, "user not found")
			return
		}
		utils.ResponseError(c, 500, "failed to get token status: "+err.Error())
		return
	}

	utils.ResponseSuccess(c, TokenStatusResponse{HasToken: hasToken})
}

// SetCurrentProjectRequest 设置当前项目请求
type SetCurrentProjectRequest struct {
	ProjectID uint `json:"project_id" binding:"required"`
}

// SetCurrentProject 设置当前用户的当前项目
func (h *ProfileHandler) SetCurrentProject(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		utils.ResponseError(c, 500, "invalid user id type")
		return
	}

	var req SetCurrentProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "invalid request: "+err.Error())
		return
	}

	// 调用 UserService 设置当前项目
	if err := h.userService.SetCurrentProject(userID, req.ProjectID); err != nil {
		if err == services.ErrUserNotFound {
			utils.ResponseError(c, 404, "user not found")
			return
		}
		utils.ResponseError(c, 500, "failed to set current project: "+err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"project_id": req.ProjectID})
}

// GetCurrentProject 获取当前用户的当前项目
func (h *ProfileHandler) GetCurrentProject(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		utils.ResponseError(c, 500, "invalid user id type")
		return
	}

	// 调用 UserService 获取当前项目 (返回 *uint)
	projectIDPtr, err := h.userService.GetCurrentProject(userID)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ResponseError(c, 404, "user not found")
			return
		}
		utils.ResponseError(c, 500, "failed to get current project: "+err.Error())
		return
	}

	// 处理指针为nil的情况
	var projectID uint = 0
	if projectIDPtr != nil {
		projectID = *projectIDPtr
	}

	utils.ResponseSuccess(c, gin.H{"project_id": projectID})
}

// GetCurrentProjectInfo 获取当前用户选择的项目详细信息(用于MCP工具)
// 此端点不需要项目成员权限检查，因为项目是用户主动选择的
func (h *ProfileHandler) GetCurrentProjectInfo(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		utils.ResponseError(c, 500, "invalid user id type")
		return
	}

	// 获取当前项目ID (GetCurrentProject返回指针)
	projectIDPtr, err := h.userService.GetCurrentProject(userID)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ResponseError(c, 404, "user not found")
			return
		}
		utils.ResponseError(c, 500, "failed to get current project: "+err.Error())
		return
	}

	// 检查指针是否为nil或项目ID是否为0
	if projectIDPtr == nil || *projectIDPtr == 0 {
		utils.ResponseError(c, 404, "no project selected")
		return
	}

	// 仅返回项目ID，MCP工具可以调用其他端点获取项目详情
	utils.ResponseSuccess(c, gin.H{
		"project_id": *projectIDPtr,
	})
}
