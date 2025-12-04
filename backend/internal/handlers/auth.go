package handlers

import (
	"errors"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// UserInfo 用户信息结构
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

// AuthHandler 认证处理器
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 处理登录请求
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "invalid request parameters: "+err.Error())
		return
	}

	// 调用认证服务
	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.ResponseError(c, 401, "invalid username or password")
			return
		}
		utils.ResponseError(c, 500, "internal server error")
		return
	}

	// 查询用户信息并返回
	user, err := h.authService.GetUserByUsername(req.Username)
	if err != nil {
		utils.ResponseError(c, 500, "failed to get user info")
		return
	}

	utils.ResponseSuccess(c, LoginResponse{
		Token: token,
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Role:     user.Role,
		},
	})
}
