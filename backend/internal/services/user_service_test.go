package services

import (
	"errors"
	"testing"
	"webtest/internal/constants"
	"webtest/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepositoryForUserService 模拟用户仓库(扩展版)
type MockUserRepositoryForUserService struct {
	mock.Mock
}

func (m *MockUserRepositoryForUserService) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepositoryForUserService) FindByNickname(nickname string) (*models.User, error) {
	args := m.Called(nickname)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepositoryForUserService) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepositoryForUserService) FindAll() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepositoryForUserService) FindByRole(role string) ([]models.User, error) {
	args := m.Called(role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepositoryForUserService) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) UpdateNickname(id uint, nickname string) error {
	args := m.Called(id, nickname)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) UpdatePassword(id uint, password string) error {
	args := m.Called(id, password)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) InitAdminUsers() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) FindByApiToken(token string) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepositoryForUserService) UpdateApiToken(id uint, token string) error {
	args := m.Called(id, token)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) FindByIDs(ids []uint) ([]models.User, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepositoryForUserService) UpdateCurrentProject(id uint, projectID uint) error {
	args := m.Called(id, projectID)
	return args.Error(0)
}

// TestGetUsers_SystemAdmin 测试系统管理员获取所有用户
func TestGetUsers_SystemAdmin(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	allUsers := []models.User{
		{ID: 1, Username: "admin", Role: constants.RoleSystemAdmin},
		{ID: 2, Username: "pm001", Role: constants.RoleProjectManager},
	}
	mockRepo.On("FindAll").Return(allUsers, nil)

	users, err := service.GetUsers(constants.RoleSystemAdmin)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	mockRepo.AssertExpectations(t)
}

// TestGetUsers_ProjectManager 测试项目管理员获取项目成员
func TestGetUsers_ProjectManager(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	members := []models.User{
		{ID: 3, Username: "member001", Role: constants.RoleProjectMember},
	}
	mockRepo.On("FindByRole", constants.RoleProjectMember).Return(members, nil)

	users, err := service.GetUsers(constants.RoleProjectManager)

	assert.NoError(t, err)
	assert.Len(t, users, 1)
	mockRepo.AssertExpectations(t)
}

// TestGetUsers_OtherRole 测试其他角色无权限
func TestGetUsers_OtherRole(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	users, err := service.GetUsers(constants.RoleProjectMember)

	assert.NoError(t, err)
	assert.Len(t, users, 0)
	mockRepo.AssertExpectations(t)
}

// TestGetAllUsers_Success 测试获取所有用户成功
func TestGetAllUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	allUsers := []models.User{
		{ID: 1, Username: "admin", Role: constants.RoleSystemAdmin},
		{ID: 2, Username: "pm001", Role: constants.RoleProjectManager},
		{ID: 3, Username: "member001", Role: constants.RoleProjectMember},
	}

	mockRepo.On("FindAll").Return(allUsers, nil)

	users, err := service.GetAllUsers()

	assert.NoError(t, err)
	assert.Len(t, users, 2) // 应该过滤掉system_admin
	assert.Equal(t, "pm001", users[0].Username)
	assert.Equal(t, "member001", users[1].Username)
	mockRepo.AssertExpectations(t)
}

// TestGetAllUsers_RepoError 测试获取用户时仓库错误
func TestGetAllUsers_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindAll").Return(nil, errors.New("database error"))

	users, err := service.GetAllUsers()

	assert.Error(t, err)
	assert.Nil(t, users)
	mockRepo.AssertExpectations(t)
}

// TestCreateUser_Success 测试创建用户成功
func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByUsername", "pm002").Return(nil, nil)
	mockRepo.On("FindByNickname", "项目经理王五").Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := service.CreateUser("pm002", "项目经理王五", constants.RoleProjectManager)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "pm002", user.Username)
	assert.Equal(t, "项目经理王五", user.Nickname)
	assert.Equal(t, constants.RoleProjectManager, user.Role)

	// 验证密码已加密
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(DefaultPasswordPM))
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestCreateUser_UsernameExists 测试创建用户时用户名已存在
func TestCreateUser_UsernameExists(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 1, Username: "pm001"}
	mockRepo.On("FindByUsername", "pm001").Return(existingUser, nil)

	user, err := service.CreateUser("pm001", "新昵称", constants.RoleProjectManager)

	assert.Error(t, err)
	assert.Equal(t, ErrUserExists, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestCreateUser_NicknameExists 测试创建用户时昵称已存在
func TestCreateUser_NicknameExists(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByUsername", "pm002").Return(nil, nil)
	existingUser := &models.User{ID: 1, Nickname: "已存在昵称"}
	mockRepo.On("FindByNickname", "已存在昵称").Return(existingUser, nil)

	user, err := service.CreateUser("pm002", "已存在昵称", constants.RoleProjectManager)

	assert.Error(t, err)
	assert.Equal(t, ErrUserExists, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestUpdateNickname_Success 测试更新昵称成功
func TestUpdateNickname_Success(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 2, Username: "pm001", Nickname: "旧昵称", Role: constants.RoleProjectManager}
	mockRepo.On("FindByID", uint(2)).Return(existingUser, nil)
	mockRepo.On("FindByNickname", "新昵称").Return(nil, nil)
	mockRepo.On("UpdateNickname", uint(2), "新昵称").Return(nil)

	user, err := service.UpdateNickname(2, "新昵称")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "新昵称", user.Nickname)
	mockRepo.AssertExpectations(t)
}

// TestUpdateNickname_UserNotFound 测试更新昵称时用户不存在
func TestUpdateNickname_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	user, err := service.UpdateNickname(999, "新昵称")

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestUpdateNickname_NicknameExists 测试更新昵称时昵称已被占用
func TestUpdateNickname_NicknameExists(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 2, Username: "pm001", Nickname: "旧昵称"}
	anotherUser := &models.User{ID: 3, Username: "pm002", Nickname: "已占用昵称"}
	mockRepo.On("FindByID", uint(2)).Return(existingUser, nil)
	mockRepo.On("FindByNickname", "已占用昵称").Return(anotherUser, nil)

	user, err := service.UpdateNickname(2, "已占用昵称")

	assert.Error(t, err)
	assert.Equal(t, ErrNicknameExists, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestUpdateNickname_NoChange 测试更新昵称时昵称未改变
func TestUpdateNickname_NoChange(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 2, Username: "pm001", Nickname: "相同昵称"}
	mockRepo.On("FindByID", uint(2)).Return(existingUser, nil)

	user, err := service.UpdateNickname(2, "相同昵称")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "相同昵称", user.Nickname)
	mockRepo.AssertExpectations(t)
}

// TestDeleteUser_Success 测试删除用户成功
func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 2, Username: "pm001", Role: constants.RoleProjectManager}
	mockRepo.On("FindByID", uint(2)).Return(existingUser, nil)
	mockRepo.On("Delete", uint(2)).Return(nil)

	err := service.DeleteUser(2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestDeleteUser_UserNotFound 测试删除用户时用户不存在
func TestDeleteUser_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	err := service.DeleteUser(999)

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

// TestDeleteUser_CannotDeleteAdmin 测试无法删除系统管理员
func TestDeleteUser_CannotDeleteAdmin(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	adminUser := &models.User{ID: 1, Username: "admin", Role: constants.RoleSystemAdmin}
	mockRepo.On("FindByID", uint(1)).Return(adminUser, nil)

	err := service.DeleteUser(1)

	assert.Error(t, err)
	assert.Equal(t, ErrCannotDeleteAdmin, err)
	mockRepo.AssertExpectations(t)
}

// TestResetPassword_Success_ProjectManager 测试重置项目管理员密码成功
func TestResetPassword_Success_ProjectManager(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 2, Username: "pm001", Role: constants.RoleProjectManager}
	mockRepo.On("FindByID", uint(2)).Return(existingUser, nil)
	mockRepo.On("UpdatePassword", uint(2), mock.AnythingOfType("string")).Return(nil)

	defaultPassword, err := service.ResetPassword(2)

	assert.NoError(t, err)
	assert.Equal(t, DefaultPasswordPM, defaultPassword)
	mockRepo.AssertExpectations(t)
}

// TestResetPassword_Success_ProjectMember 测试重置项目成员密码成功
func TestResetPassword_Success_ProjectMember(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 3, Username: "member001", Role: constants.RoleProjectMember}
	mockRepo.On("FindByID", uint(3)).Return(existingUser, nil)
	mockRepo.On("UpdatePassword", uint(3), mock.AnythingOfType("string")).Return(nil)

	defaultPassword, err := service.ResetPassword(3)

	assert.NoError(t, err)
	assert.Equal(t, DefaultPasswordMember, defaultPassword)
	mockRepo.AssertExpectations(t)
}

// TestResetPassword_UserNotFound 测试重置密码时用户不存在
func TestResetPassword_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	defaultPassword, err := service.ResetPassword(999)

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	assert.Empty(t, defaultPassword)
	mockRepo.AssertExpectations(t)
}

// TestCheckUsernameExists_True 测试检查用户名存在
func TestCheckUsernameExists_True(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 1, Username: "pm001"}
	mockRepo.On("FindByUsername", "pm001").Return(existingUser, nil)

	exists, err := service.CheckUsernameExists("pm001")

	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

// TestCheckUsernameExists_False 测试检查用户名不存在
func TestCheckUsernameExists_False(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByUsername", "newuser").Return(nil, nil)

	exists, err := service.CheckUsernameExists("newuser")

	assert.NoError(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}

// TestCheckNicknameExists_True 测试检查昵称存在
func TestCheckNicknameExists_True(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 1, Nickname: "已存在昵称"}
	mockRepo.On("FindByNickname", "已存在昵称").Return(existingUser, nil)

	exists, err := service.CheckNicknameExists("已存在昵称")

	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

// TestCheckNicknameExists_False 测试检查昵称不存在
func TestCheckNicknameExists_False(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByNickname", "新昵称").Return(nil, nil)

	exists, err := service.CheckNicknameExists("新昵称")

	assert.NoError(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}

// TestCreateUser_ProjectMemberPassword 测试创建项目成员时使用正确的默认密码
func TestCreateUser_ProjectMemberPassword(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByUsername", "member002").Return(nil, nil)
	mockRepo.On("FindByNickname", "测试员赵六").Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := service.CreateUser("member002", "测试员赵六", constants.RoleProjectMember)

	assert.NoError(t, err)
	assert.NotNil(t, user)

	// 验证密码是项目成员默认密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(DefaultPasswordMember))
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestCheckUsernameExists_RepoError 测试检查用户名时仓库错误
func TestCheckUsernameExists_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(nil, errors.New("database error"))

	exists, err := service.CheckUsernameExists("testuser")

	assert.Error(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}

// TestCheckNicknameExists_RepoError 测试检查昵称时仓库错误
func TestCheckNicknameExists_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByNickname", "测试昵称").Return(nil, errors.New("database error"))

	exists, err := service.CheckNicknameExists("测试昵称")

	assert.Error(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}

// TestCreateUser_RepoError 测试创建用户时仓库错误
func TestCreateUser_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByUsername", "pm003").Return(nil, nil)
	mockRepo.On("FindByNickname", "新用户").Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(errors.New("database error"))

	user, err := service.CreateUser("pm003", "新用户", constants.RoleProjectManager)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestUpdateNickname_RepoError 测试更新昵称时仓库错误
func TestUpdateNickname_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 2, Username: "pm001", Nickname: "旧昵称"}
	mockRepo.On("FindByID", uint(2)).Return(existingUser, nil)
	mockRepo.On("FindByNickname", "新昵称").Return(nil, nil)
	mockRepo.On("UpdateNickname", uint(2), "新昵称").Return(errors.New("database error"))

	user, err := service.UpdateNickname(2, "新昵称")

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestDeleteUser_RepoError 测试删除用户时仓库错误
func TestDeleteUser_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 2, Username: "pm001", Role: constants.RoleProjectManager}
	mockRepo.On("FindByID", uint(2)).Return(existingUser, nil)
	mockRepo.On("Delete", uint(2)).Return(errors.New("database error"))

	err := service.DeleteUser(2)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestResetPassword_RepoError 测试重置密码时仓库错误
func TestResetPassword_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepositoryForUserService)
	service := NewUserService(mockRepo)

	existingUser := &models.User{ID: 2, Username: "pm001", Role: constants.RoleProjectManager}
	mockRepo.On("FindByID", uint(2)).Return(existingUser, nil)
	mockRepo.On("UpdatePassword", uint(2), mock.AnythingOfType("string")).Return(errors.New("database error"))

	defaultPassword, err := service.ResetPassword(2)

	assert.Error(t, err)
	assert.Empty(t, defaultPassword)
	mockRepo.AssertExpectations(t)
}
