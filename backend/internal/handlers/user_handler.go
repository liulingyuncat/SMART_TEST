package handlers

import (
	"strconv"
	"webtest/internal/models"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateUserRequest 创建用户请求（系统管理员）
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
	Role     string `json:"role" binding:"required,oneof=project_manager project_member"`
}

// CreateUserForPMRequest 创建用户请求（项目管理员专用）
type CreateUserForPMRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
	Role     string `json:"role" binding:"required,oneof=project_member"` // 只允许创建项目成员
}

// UpdateNicknameRequest 更新昵称请求
type UpdateNicknameRequest struct {
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
}

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsersAuto 获取用户列表（根据当前用户角色自动分流）
func (h *UserHandler) GetUsersAuto(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}

	roleStr := role.(string)
	var users []models.User
	var err error

	// 根据角色获取不同的用户列表
	if roleStr == "system_admin" {
		users, err = h.userService.GetAllUsers()
	} else if roleStr == "project_manager" {
		users, err = h.userService.GetProjectMembers()
	} else {
		utils.ResponseError(c, 403, "insufficient permissions")
		return
	}

	if err != nil {
		utils.ResponseError(c, 500, "failed to get users: "+err.Error())
		return
	}

	// 转换为 UserInfo 结构
	userInfos := make([]UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Role:     user.Role,
		}
	}

	utils.ResponseSuccess(c, gin.H{
		"users": userInfos,
		"total": len(userInfos),
	})
}

// GetUsersForAdmin 获取用户列表（系统管理员，排除system_admin）
func (h *UserHandler) GetUsersForAdmin(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		utils.ResponseError(c, 500, "failed to get users: "+err.Error())
		return
	}

	// 转换为 UserInfo 结构（复用现有结构）
	userInfos := make([]UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Role:     user.Role,
		}
	}

	utils.ResponseSuccess(c, gin.H{
		"users": userInfos,
		"total": len(userInfos),
	})
}

// GetUsersForPM 获取用户列表（项目管理员，仅显示项目成员）
func (h *UserHandler) GetUsersForPM(c *gin.Context) {
	users, err := h.userService.GetProjectMembers()
	if err != nil {
		utils.ResponseError(c, 500, "failed to get users: "+err.Error())
		return
	}

	// 转换为 UserInfo 结构
	userInfos := make([]UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Role:     user.Role,
		}
	}

	utils.ResponseSuccess(c, gin.H{
		"users": userInfos,
		"total": len(userInfos),
	})
}

// CreateUserAuto 创建用户（根据当前用户角色自动分流）
func (h *UserHandler) CreateUserAuto(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}

	roleStr := role.(string)

	// 根据角色调用不同的创建方法
	if roleStr == "system_admin" {
		h.CreateUserForAdmin(c)
	} else if roleStr == "project_manager" {
		h.CreateUserForPM(c)
	} else {
		utils.ResponseError(c, 403, "insufficient permissions")
	}
}

// CreateUserForAdmin 创建用户（系统管理员）
func (h *UserHandler) CreateUserForAdmin(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "invalid request: "+err.Error())
		return
	}

	user, err := h.userService.CreateUser(req.Username, req.Nickname, req.Role)
	if err != nil {
		if err == services.ErrUserExists {
			utils.ResponseError(c, 400, "user already exists")
			return
		}
		utils.ResponseError(c, 500, "failed to create user: "+err.Error())
		return
	}

	// 返回 UserInfo 格式
	utils.ResponseSuccess(c, UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Role:     user.Role,
	})
}

// CreateUserForPM 创建用户（项目管理员，仅能创建项目成员）
func (h *UserHandler) CreateUserForPM(c *gin.Context) {
	var req CreateUserForPMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "invalid request: "+err.Error())
		return
	}

	user, err := h.userService.CreateUser(req.Username, req.Nickname, req.Role)
	if err != nil {
		if err == services.ErrUserExists {
			utils.ResponseError(c, 400, "user already exists")
			return
		}
		utils.ResponseError(c, 500, "failed to create user: "+err.Error())
		return
	}

	// 返回 UserInfo 格式
	utils.ResponseSuccess(c, UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Role:     user.Role,
	})
}

// UpdateNickname 更新用户昵称
func (h *UserHandler) UpdateNickname(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid user id")
		return
	}

	var req UpdateNicknameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "invalid request: "+err.Error())
		return
	}

	user, err := h.userService.UpdateNickname(uint(userID), req.Nickname)
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

	utils.ResponseSuccess(c, UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Role:     user.Role,
	})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid user id")
		return
	}

	if err := h.userService.DeleteUser(uint(userID)); err != nil {
		if err == services.ErrCannotDeleteAdmin {
			utils.ResponseError(c, 403, "cannot delete system admin")
			return
		}
		if err == services.ErrUserNotFound {
			utils.ResponseError(c, 404, "user not found")
			return
		}
		utils.ResponseError(c, 500, "failed to delete user: "+err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "user deleted"})
}

// ResetPassword 重置用户密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid user id")
		return
	}

	defaultPassword, err := h.userService.ResetPassword(uint(userID))
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ResponseError(c, 404, "user not found")
			return
		}
		utils.ResponseError(c, 500, "failed to reset password: "+err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"message":          "password reset successfully",
		"default_password": defaultPassword,
	})
}

// CheckUnique 检查用户名或昵称唯一性（合并接口）
func (h *UserHandler) CheckUnique(c *gin.Context) {
	username := c.Query("username")
	nickname := c.Query("nickname")

	if username == "" && nickname == "" {
		utils.ResponseError(c, 400, "username or nickname is required")
		return
	}

	var exists bool
	var err error

	if username != "" {
		exists, err = h.userService.CheckUsernameExists(username)
	} else {
		exists, err = h.userService.CheckNicknameExists(nickname)
	}

	if err != nil {
		utils.ResponseError(c, 500, "failed to check uniqueness: "+err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"exists": exists})
}
