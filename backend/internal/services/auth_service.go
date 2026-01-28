package services

import (
	"errors"
	"fmt"
	"os"
	"time"
	"webtest/internal/constants"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidCredentials 无效凭据错误
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrTokenGeneration Token 生成失败错误
	ErrTokenGeneration = errors.New("failed to generate token")
)

// AuthService 认证服务接口
type AuthService interface {
	Login(username, password string) (string, error)
	ValidateToken(tokenString string) (*CustomClaims, error)
	GetUserByUsername(username string) (*models.User, error)
	InitAdminUsers() error
}

// authService 认证服务实现
type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret []byte
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repositories.UserRepository) AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_key_change_in_production" // 默认密钥(生产环境必须修改)
	}
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(secret),
	}
}

// Login 用户登录
func (s *authService) Login(username, password string) (string, error) {
	// 查询用户
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", fmt.Errorf("query user failed: %w", err)
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// 生成 JWT Token
	token, err := s.generateToken(user)
	if err != nil {
		return "", ErrTokenGeneration
	}

	return token, nil
}

// generateToken 生成 JWT Token
func (s *authService) generateToken(user *models.User) (string, error) {
	claims := CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "webtest",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign token failed: %w", err)
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT Token
func (s *authService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token failed: %w", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// InitAdminUsers 初始化管理员账号
func (s *authService) InitAdminUsers() error {
	// 预置管理员账号列表
	adminUsers := []struct {
		Username string
		Nickname string
		Password string
		Role     string
	}{
		{Username: "admin", Nickname: "系统管理员", Password: "admin123", Role: constants.RoleSystemAdmin},
		{Username: "root", Nickname: "超级管理员", Password: "root123", Role: constants.RoleSystemAdmin},
		{Username: "Padmin", Nickname: "项目管理员", Password: "123456", Role: constants.RoleProjectManager},
	}

	for _, admin := range adminUsers {
		// 检查是否已存在
		existingUser, err := s.userRepo.FindByUsername(admin.Username)
		if err != nil {
			return fmt.Errorf("query user failed: %w", err)
		}
		if existingUser != nil {
			continue // 已存在则跳过
		}

		// 密码加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password failed: %w", err)
		}

		// 创建用户
		user := &models.User{
			Username: admin.Username,
			Nickname: admin.Nickname,
			Password: string(hashedPassword),
			Role:     admin.Role,
		}
		if err := s.userRepo.Create(user); err != nil {
			return fmt.Errorf("create admin user failed: %w", err)
		}
	}

	return nil
}

// GetUserByUsername 根据用户名获取用户信息
func (s *authService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.FindByUsername(username)
}
