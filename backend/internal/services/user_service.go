package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"webtest/internal/constants"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists               = errors.New("username or nickname already exists")
	ErrNicknameExists           = errors.New("nickname already exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrCannotDeleteAdmin        = errors.New("cannot delete system admin")
	ErrCurrentPasswordIncorrect = errors.New("current password is incorrect")
	ErrNewPasswordSameAsCurrent = errors.New("new password cannot be same as current")
	ErrTokenGenerationFailed    = errors.New("token generation failed")
	ErrInvalidApiToken          = errors.New("invalid api token")
)

const (
	DefaultPasswordPM     = "admin!123" // 项目管理员默认密码
	DefaultPasswordMember = "user!123"  // 项目成员默认密码
)

// UserService 用户服务接口
type UserService interface {
	GetUsers(currentRole string) ([]models.User, error)
	// T18 人员管理功能
	GetAllUsers() ([]models.User, error)
	GetProjectMembers() ([]models.User, error) // 获取所有项目成员
	CreateUser(username, nickname, role string) (*models.User, error)
	UpdateNickname(userID uint, nickname string) (*models.User, error)
	DeleteUser(userID uint) error
	ResetPassword(userID uint) (string, error)
	CheckUsernameExists(username string) (bool, error)
	CheckNicknameExists(nickname string) (bool, error)
	// T22 个人信息查看
	GetUserByID(userID uint) (*models.User, error)
	// T23 密码修改与Token功能
	ChangePassword(userID uint, currentPwd, newPwd string) error
	GenerateApiToken(userID uint) (string, error)
	ValidateApiToken(token string) (*models.User, error)
	HasApiToken(userID uint) (bool, error)
	// T50 当前项目管理
	SetCurrentProject(userID uint, projectID uint) error
	GetCurrentProject(userID uint) (*uint, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUsers 根据当前用户角色获取用户列表
func (s *userService) GetUsers(currentRole string) ([]models.User, error) {
	// 系统管理员可以看到所有用户
	if currentRole == constants.RoleSystemAdmin {
		return s.userRepo.FindAll()
	}

	// 项目管理员只能看到 project_member 角色用户
	if currentRole == constants.RoleProjectManager {
		return s.userRepo.FindByRole(constants.RoleProjectMember)
	}

	// 其他角色无权限
	return []models.User{}, nil
}

// GetAllUsers 获取所有用户（排除system_admin）- T18人员管理功能
func (s *userService) GetAllUsers() ([]models.User, error) {
	allUsers, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 过滤掉system_admin角色
	users := make([]models.User, 0)
	for _, user := range allUsers {
		if user.Role != constants.RoleSystemAdmin {
			users = append(users, user)
		}
	}

	return users, nil
}

// GetProjectMembers 获取所有项目成员和项目管理员 - 项目管理员专用
// 返回 project_manager 和 project_member 角色的用户，用于人员分配
func (s *userService) GetProjectMembers() ([]models.User, error) {
	// 获取所有用户
	allUsers, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 过滤出 project_manager 和 project_member 角色的用户
	var users []models.User
	for _, user := range allUsers {
		if user.Role == constants.RoleProjectManager || user.Role == constants.RoleProjectMember {
			users = append(users, user)
		}
	}

	return users, nil
}

// CreateUser 创建新用户
func (s *userService) CreateUser(username, nickname, role string) (*models.User, error) {
	// 检查用户名是否存在
	existingUser, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// 检查昵称是否存在
	existingUser, err = s.userRepo.FindByNickname(nickname)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// 根据角色生成默认密码
	defaultPassword := s.getDefaultPassword(role)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		Username: username,
		Nickname: nickname,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateNickname 更新用户昵称
func (s *userService) UpdateNickname(userID uint, nickname string) (*models.User, error) {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 如果昵称未改变，直接返回
	if user.Nickname == nickname {
		return user, nil
	}

	// 检查新昵称是否与其他用户重复
	existingUser, err := s.userRepo.FindByNickname(nickname)
	if err != nil {
		return nil, err
	}
	if existingUser != nil && existingUser.ID != userID {
		return nil, ErrNicknameExists
	}

	// 更新昵称
	if err := s.userRepo.UpdateNickname(userID, nickname); err != nil {
		return nil, err
	}

	// 返回更新后的用户
	user.Nickname = nickname
	return user, nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(userID uint) error {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 防止删除系统管理员
	if user.Role == constants.RoleSystemAdmin {
		return ErrCannotDeleteAdmin
	}

	// 执行删除（软删除）
	return s.userRepo.Delete(userID)
}

// ResetPassword 重置密码为默认值
func (s *userService) ResetPassword(userID uint) (string, error) {
	// 查询用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	// 根据角色获取默认密码
	defaultPassword := s.getDefaultPassword(user.Role)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		return "", err
	}

	return defaultPassword, nil
}

// CheckUsernameExists 检查用户名是否存在
func (s *userService) CheckUsernameExists(username string) (bool, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

// CheckNicknameExists 检查昵称是否存在
func (s *userService) CheckNicknameExists(nickname string) (bool, error) {
	user, err := s.userRepo.FindByNickname(nickname)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

// GetUserByID 根据用户ID获取用户信息 - T22 个人信息查看
func (s *userService) GetUserByID(userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// getDefaultPassword 根据角色获取默认密码
func (s *userService) getDefaultPassword(role string) string {
	if role == constants.RoleProjectManager {
		return DefaultPasswordPM
	}
	return DefaultPasswordMember
}

// ChangePassword 修改用户密码 - T23 密码修改功能
func (s *userService) ChangePassword(userID uint, currentPwd, newPwd string) error {
	// 获取用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 验证当前密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPwd)); err != nil {
		return ErrCurrentPasswordIncorrect
	}

	// 检查新密码是否与旧密码相同
	if currentPwd == newPwd {
		return ErrNewPasswordSameAsCurrent
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	return s.userRepo.UpdatePassword(userID, string(hashedPassword))
}

// GenerateApiToken 生成新的API Token - T23 Token生成功能
func (s *userService) GenerateApiToken(userID uint) (string, error) {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	// 生成32字节随机数据，编码为64字符hex字符串
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", ErrTokenGenerationFailed
	}
	token := hex.EncodeToString(tokenBytes)

	// 保存Token到数据库
	if err := s.userRepo.UpdateApiToken(userID, token); err != nil {
		return "", err
	}

	return token, nil
}

// ValidateApiToken 验证API Token有效性 - T23 Token认证
func (s *userService) ValidateApiToken(token string) (*models.User, error) {
	if token == "" {
		return nil, ErrInvalidApiToken
	}

	user, err := s.userRepo.FindByApiToken(token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidApiToken
	}

	return user, nil
}

// HasApiToken 检查用户是否已有Token - T23 Token状态查询
func (s *userService) HasApiToken(userID uint) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, ErrUserNotFound
	}

	return user.ApiToken != nil && *user.ApiToken != "", nil
}

// SetCurrentProject 设置用户的当前项目 - T50
func (s *userService) SetCurrentProject(userID uint, projectID uint) error {
	return s.userRepo.UpdateCurrentProject(userID, projectID)
}

// GetCurrentProject 获取用户的当前项目 - T50
func (s *userService) GetCurrentProject(userID uint) (*uint, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user.CurrentProjectID, nil
}
