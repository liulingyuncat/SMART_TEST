package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// MemberWithUser 成员及用户信息
type MemberWithUser struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

// ProjectMemberRepository 项目成员仓库接口
type ProjectMemberRepository interface {
	IsMember(projectID uint, userID uint) (bool, error)
	IsMemberWithRole(projectID uint, userID uint, role string) (bool, error)
	FindByProjectID(projectID uint) ([]models.ProjectMember, error)
	FindMembersWithUser(projectID uint) ([]MemberWithUser, error)
	AddMember(member *models.ProjectMember) error
	RemoveMember(projectID uint, userID uint) error
	GetMemberRole(projectID uint, userID uint) (string, error)
	// 新增方法 - T21人员分配批量更新
	BatchUpdateMembers(projectID uint, managers []uint, members []uint) error
}

// projectMemberRepository 项目成员仓库实现
type projectMemberRepository struct {
	db *gorm.DB
}

// NewProjectMemberRepository 创建项目成员仓库实例
func NewProjectMemberRepository(db *gorm.DB) ProjectMemberRepository {
	return &projectMemberRepository{db: db}
}

// IsMember 判断用户是否是项目成员
func (r *projectMemberRepository) IsMember(projectID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

// IsMemberWithRole 判断用户是否是项目成员且具有指定角色
func (r *projectMemberRepository) IsMemberWithRole(projectID uint, userID uint, role string) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ? AND role = ?", projectID, userID, role).
		Count(&count).Error
	return count > 0, err
}

// FindByProjectID 查询指定项目的所有成员
func (r *projectMemberRepository) FindByProjectID(projectID uint) ([]models.ProjectMember, error) {
	var members []models.ProjectMember
	err := r.db.Where("project_id = ?", projectID).Find(&members).Error
	return members, err
}

// AddMember 添加项目成员
func (r *projectMemberRepository) AddMember(member *models.ProjectMember) error {
	return r.db.Create(member).Error
}

// RemoveMember 移除项目成员(硬删除)
func (r *projectMemberRepository) RemoveMember(projectID uint, userID uint) error {
	return r.db.Unscoped().Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&models.ProjectMember{}).Error
}

// GetMemberRole 获取用户在项目中的角色
func (r *projectMemberRepository) GetMemberRole(projectID uint, userID uint) (string, error) {
	var member models.ProjectMember
	err := r.db.Where("project_id = ? AND user_id = ?", projectID, userID).
		First(&member).Error
	if err != nil {
		return "", err
	}
	return member.Role, nil
}

// FindMembersWithUser 查询指定项目的所有成员及其用户信息
func (r *projectMemberRepository) FindMembersWithUser(projectID uint) ([]MemberWithUser, error) {
	var members []MemberWithUser
	err := r.db.Table("project_members").
		Select("project_members.user_id, project_members.role, users.username, users.nickname").
		Joins("INNER JOIN users ON project_members.user_id = users.id").
		Where("project_members.project_id = ?", projectID).
		Order("project_members.role, users.username").
		Scan(&members).Error
	return members, err
}

// BatchUpdateMembers 批量更新项目成员（事务）
func (r *projectMemberRepository) BatchUpdateMembers(projectID uint, managers []uint, members []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 步骤1: 删除项目所有现有成员(硬删除)
		if err := tx.Unscoped().Where("project_id = ?", projectID).Delete(&models.ProjectMember{}).Error; err != nil {
			return err
		}

		// 步骤2: 批量插入管理员
		for _, userID := range managers {
			member := &models.ProjectMember{
				ProjectID: projectID,
				UserID:    userID,
				Role:      "project_manager",
			}
			if err := tx.Create(member).Error; err != nil {
				return err
			}
		}

		// 步骤3: 批量插入成员
		for _, userID := range members {
			member := &models.ProjectMember{
				ProjectID: projectID,
				UserID:    userID,
				Role:      "project_member",
			}
			if err := tx.Create(member).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
