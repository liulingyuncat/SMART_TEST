package repositories

import (
	"fmt"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层接口
type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
	FindAll() ([]models.User, error)
	FindByRole(role string) ([]models.User, error)
	Create(user *models.User) error
	InitAdminUsers() error
	// 新增方法 - T18人员管理功能
	FindByID(id uint) (*models.User, error)
	FindByNickname(nickname string) (*models.User, error)
	UpdateNickname(id uint, nickname string) error
	UpdatePassword(id uint, password string) error
	Delete(id uint) error
	// 新增方法 - T21人员分配批量操作
	FindByIDs(ids []uint) ([]models.User, error)
	// 新增方法 - T23 API Token功能
	FindByApiToken(token string) (*models.User, error)
	UpdateApiToken(id uint, token string) error
	// 新增方法 - T50 当前项目管理
	UpdateCurrentProject(id uint, projectID uint) error
}

// userRepository 用户仓库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 未找到返回 nil,nil
		}
		return nil, fmt.Errorf("query user failed: %w", err)
	}
	return &user, nil
}

// FindAll 查询所有用户
func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("deleted_at IS NULL").Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("query all users failed: %w", err)
	}
	return users, nil
}

// FindByRole 根据角色查询用户
func (r *userRepository) FindByRole(role string) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role = ? AND deleted_at IS NULL", role).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("query users by role failed: %w", err)
	}
	return users, nil
}

// Create 创建新用户
func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("create user failed: %w", err)
	}
	return nil
}

// InitAdminUsers 初始化管理员用户
// 密码 hash 需由调用方(Service层)提前生成
func (r *userRepository) InitAdminUsers() error {
	// 检查是否已存在管理员
	var count int64
	if err := r.db.Model(&models.User{}).
		Where("role = ?", "admin").
		Count(&count).Error; err != nil {
		return fmt.Errorf("count admin users failed: %w", err)
	}

	// 已存在管理员则跳过初始化
	if count > 0 {
		return nil
	}

	// 注意: 实际密码 hash 应在 Service 层生成
	// 这里仅作为示例,实际调用时需传入已 hash 的密码
	return nil
}

// FindByIDs 根据ID列表批量查询用户
func (r *userRepository) FindByIDs(ids []uint) ([]models.User, error) {
	var users []models.User
	if len(ids) == 0 {
		return users, nil
	}
	err := r.db.Where("id IN ? AND deleted_at IS NULL", ids).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("query users by IDs failed: %w", err)
	}
	return users, nil
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("deleted_at IS NULL").First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by id failed: %w", err)
	}
	return &user, nil
}

// FindByNickname 根据昵称查找用户
func (r *userRepository) FindByNickname(nickname string) (*models.User, error) {
	var user models.User
	err := r.db.Where("nickname = ? AND deleted_at IS NULL", nickname).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by nickname failed: %w", err)
	}
	return &user, nil
}

// UpdateNickname 更新昵称
func (r *userRepository) UpdateNickname(id uint, nickname string) error {
	result := r.db.Model(&models.User{}).Where("id = ? AND deleted_at IS NULL", id).Update("nickname", nickname)
	if result.Error != nil {
		return fmt.Errorf("update nickname failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// UpdatePassword 更新密码
func (r *userRepository) UpdatePassword(id uint, password string) error {
	result := r.db.Model(&models.User{}).Where("id = ? AND deleted_at IS NULL", id).Update("password", password)
	if result.Error != nil {
		return fmt.Errorf("update password failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// Delete 删除用户(硬删除)
func (r *userRepository) Delete(id uint) error {
	result := r.db.Unscoped().Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete user failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// FindByApiToken 根据API Token查找用户
func (r *userRepository) FindByApiToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("api_token = ? AND deleted_at IS NULL", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by api token failed: %w", err)
	}
	return &user, nil
}

// UpdateApiToken 更新用户的API Token
func (r *userRepository) UpdateApiToken(id uint, token string) error {
	result := r.db.Model(&models.User{}).Where("id = ? AND deleted_at IS NULL", id).Update("api_token", token)
	if result.Error != nil {
		return fmt.Errorf("update api token failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// UpdateCurrentProject 更新用户的当前项目 - T50
func (r *userRepository) UpdateCurrentProject(id uint, projectID uint) error {
	result := r.db.Model(&models.User{}).Where("id = ? AND deleted_at IS NULL", id).Update("current_project_id", projectID)
	if result.Error != nil {
		return fmt.Errorf("update current project failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
