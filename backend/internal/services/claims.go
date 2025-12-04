package services

import "github.com/golang-jwt/jwt/v5"

// CustomClaims 自定义JWT Claims
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
