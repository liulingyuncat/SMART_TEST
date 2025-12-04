package middleware

import (
	"strings"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ResponseError(c, 401, "authorization header required")
			c.Abort()
			return
		}

		// 提取 Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ResponseError(c, 401, "invalid authorization format")
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 验证 Token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			utils.ResponseError(c, 401, "invalid or expired token")
			c.Abort()
			return
		}

		// 将用户信息设置到上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// 继续处理
		c.Next()
	}
}

// RequireRole 角色权限中间件
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		role, exists := c.Get("role")
		if !exists {
			utils.ResponseError(c, 403, "role not found in context")
			c.Abort()
			return
		}

		// 检查角色是否在允许列表中
		roleStr, ok := role.(string)
		if !ok {
			utils.ResponseError(c, 403, "invalid role type")
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		// 角色不匹配
		utils.ResponseError(c, 403, "insufficient permissions")
		c.Abort()
	}
}
