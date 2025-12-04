package services

import (
	"testing"
	"webtest/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository 模拟用户仓库
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByNickname(nickname string) (*models.User, error) {
	args := m.Called(nickname)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) FindByRole(role string) ([]models.User, error) {
	args := m.Called(role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateNickname(id uint, nickname string) error {
	args := m.Called(id, nickname)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(id uint, password string) error {
	args := m.Called(id, password)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) InitAdminUsers() error {
	args := m.Called()
	return args.Error(0)
}

// TestLogin_Success 测试登录成功场景
func TestLogin_Success(t *testing.T) {
	// 准备测试数据
	password := "admin123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUser := &models.User{
		ID:       1,
		Username: "admin",
		Nickname: "管理员",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	// 创建 Mock 仓库
	mockRepo := new(MockUserRepository)
	mockRepo.On("FindByUsername", "admin").Return(mockUser, nil)

	// 创建认证服务
	authService := NewAuthService(mockRepo)

	// 执行登录
	token, err := authService.Login("admin", password)

	// 断言
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

// TestLogin_InvalidPassword 测试密码错误场景
func TestLogin_InvalidPassword(t *testing.T) {
	// 准备测试数据
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	mockUser := &models.User{
		ID:       1,
		Username: "admin",
		Password: string(hashedPassword),
	}

	// 创建 Mock 仓库
	mockRepo := new(MockUserRepository)
	mockRepo.On("FindByUsername", "admin").Return(mockUser, nil)

	// 创建认证服务
	authService := NewAuthService(mockRepo)

	// 执行登录(错误密码)
	token, err := authService.Login("admin", "wrongpassword")

	// 断言
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

// TestLogin_UserNotFound 测试用户不存在场景
func TestLogin_UserNotFound(t *testing.T) {
	// 创建 Mock 仓库
	mockRepo := new(MockUserRepository)
	mockRepo.On("FindByUsername", "nonexistent").Return(nil, nil)

	// 创建认证服务
	authService := NewAuthService(mockRepo)

	// 执行登录
	token, err := authService.Login("nonexistent", "password")

	// 断言
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

// TestValidateToken 测试 Token 验证
func TestValidateToken(t *testing.T) {
	// 准备测试数据
	password := "admin123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUser := &models.User{
		ID:       1,
		Username: "admin",
		Password: string(hashedPassword),
	}

	// 创建 Mock 仓库
	mockRepo := new(MockUserRepository)
	mockRepo.On("FindByUsername", "admin").Return(mockUser, nil)

	// 创建认证服务
	authService := NewAuthService(mockRepo)

	// 生成 Token
	token, _ := authService.Login("admin", password)

	// 验证 Token
	claims, err := authService.ValidateToken(token)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "admin", claims.Subject)
}
