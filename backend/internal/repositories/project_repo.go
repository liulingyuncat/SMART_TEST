package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// ProjectRepository 项目仓库接口
type ProjectRepository interface {
	FindProjectsByUserID(userID uint) ([]models.Project, error)
	Create(project *models.Project) error
	ExistsByName(name string) (bool, error)
	UpdateName(id uint, name string) (*models.Project, error)
	Update(id uint, updates map[string]interface{}) (*models.Project, error)
	GetByID(id uint) (*models.Project, error)
	DeleteWithCascade(id uint) error
}

// projectRepository 项目仓库实现
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository 创建项目仓库实例
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// FindProjectsByUserID 根据用户ID查询参与的项目
func (r *projectRepository) FindProjectsByUserID(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.
		Joins("JOIN project_members ON project_members.project_id = projects.id").
		Where("project_members.user_id = ?", userID).
		Order("projects.created_at DESC").
		Distinct().
		Find(&projects).Error
	return projects, err
}

// Create 创建新项目
func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

// ExistsByName 检查项目名是否已存在
func (r *projectRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Project{}).
		Where("name = ?", name).
		Count(&count).Error
	return count > 0, err
}

// UpdateName 更新项目名称
func (r *projectRepository) UpdateName(id uint, name string) (*models.Project, error) {
	var project models.Project
	err := r.db.Model(&models.Project{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"name": name}).Error
	if err != nil {
		return nil, err
	}
	err = r.db.First(&project, id).Error
	return &project, err
}

// Update 更新项目字段(支持 name/description/status/owner_id)
func (r *projectRepository) Update(id uint, updates map[string]interface{}) (*models.Project, error) {
	err := r.db.Model(&models.Project{}).
		Where("id = ?", id).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// GetByID 根据ID查询项目
func (r *projectRepository) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.
		Table("projects").
		Select("projects.id, projects.name, projects.description, projects.status, projects.owner_id, projects.created_at, projects.updated_at, projects.deleted_at, users.nickname as owner_name").
		Joins("LEFT JOIN users ON users.id = projects.owner_id").
		Where("projects.id = ?", id).
		Scan(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// DeleteWithCascade 级联删除项目及其关联数据
func (r *projectRepository) DeleteWithCascade(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除项目成员关联
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectMember{}).Error; err != nil {
			return err
		}
		// 删除项目本身
		if err := tx.Delete(&models.Project{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
