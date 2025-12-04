package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// ProjectMemberRepository 项目成员仓库接口
type ProjectMemberRepository interface {
	IsMember(projectID uint, userID uint) (bool, error)
	IsMemberWithRole(projectID uint, userID uint, role string) (bool, error)
	FindByProjectID(projectID uint) ([]models.ProjectMember, error)
	AddMember(member *models.ProjectMember) error
	RemoveMember(projectID uint, userID uint) error
	GetMemberRole(projectID uint, userID uint) (string, error)
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

// RemoveMember 移除项目成员
func (r *projectMemberRepository) RemoveMember(projectID uint, userID uint) error {
	return r.db.Where("project_id = ? AND user_id = ?", projectID, userID).
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
