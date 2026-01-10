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
		Table("projects").
		Select("projects.id, projects.name, projects.description, projects.status, projects.owner_id, projects.created_at, projects.updated_at, projects.deleted_at, users.nickname as owner_name").
		Joins("JOIN project_members ON project_members.project_id = projects.id").
		Joins("LEFT JOIN users ON users.id = projects.owner_id").
		Where("project_members.user_id = ?", userID).
		Order("projects.created_at DESC").
		Scan(&projects).Error
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

// DeleteWithCascade 级联删除项目及其关联数据(硬删除)
func (r *projectRepository) DeleteWithCascade(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除项目成员关联(硬删除)
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.ProjectMember{}).Error; err != nil {
			return err
		}

		// 2. 删除用例集
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.CaseGroup{}).Error; err != nil {
			return err
		}

		// 3. 删除手工测试用例
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.ManualTestCase{}).Error; err != nil {
			return err
		}

		// 4. 删除自动化测试用例
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.AutoTestCase{}).Error; err != nil {
			return err
		}

		// 5. 删除API测试用例
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.ApiTestCase{}).Error; err != nil {
			return err
		}

		// 6. 删除执行任务(需要先删除执行用例结果)
		var taskUUIDs []string
		if err := tx.Model(&models.ExecutionTask{}).Where("project_id = ?", id).Pluck("task_uuid", &taskUUIDs).Error; err != nil {
			return err
		}
		if len(taskUUIDs) > 0 {
			// 删除执行用例结果
			if err := tx.Unscoped().Where("task_uuid IN ?", taskUUIDs).Delete(&models.ExecutionCaseResult{}).Error; err != nil {
				return err
			}
		}
		// 删除执行任务
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.ExecutionTask{}).Error; err != nil {
			return err
		}

		// 7. 删除版本管理相关
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.Version{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.CaseVersion{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.AutoTestCaseVersion{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.ApiTestCaseVersion{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.WebCaseVersion{}).Error; err != nil {
			return err
		}

		// 8. 删除缺陷管理相关
		var defectIDs []uint
		if err := tx.Model(&models.Defect{}).Where("project_id = ?", id).Pluck("id", &defectIDs).Error; err != nil {
			return err
		}
		if len(defectIDs) > 0 {
			// 删除缺陷附件
			if err := tx.Unscoped().Where("defect_id IN ?", defectIDs).Delete(&models.DefectAttachment{}).Error; err != nil {
				return err
			}
			// 删除缺陷评论
			if err := tx.Unscoped().Where("defect_id IN ?", defectIDs).Delete(&models.DefectComment{}).Error; err != nil {
				return err
			}
		}
		// 删除缺陷
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.Defect{}).Error; err != nil {
			return err
		}
		// 删除缺陷配置
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.DefectSubject{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.DefectPhase{}).Error; err != nil {
			return err
		}

		// 9. 删除需求管理相关
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.RequirementItem{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.ViewpointItem{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.RawDocument{}).Error; err != nil {
			return err
		}

		// 10. 删除评审相关
		var reviewIDs []uint
		if err := tx.Model(&models.CaseReview{}).Where("project_id = ?", id).Pluck("id", &reviewIDs).Error; err != nil {
			return err
		}
		if len(reviewIDs) > 0 {
			// 删除评审条目
			if err := tx.Unscoped().Where("review_id IN ?", reviewIDs).Delete(&models.CaseReviewItem{}).Error; err != nil {
				return err
			}
		}
		// 删除评审
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.CaseReview{}).Error; err != nil {
			return err
		}

		// 11. 删除AI质量报告
		if err := tx.Unscoped().Where("project_id = ?", id).Delete(&models.AIReport{}).Error; err != nil {
			return err
		}

		// 12. 删除项目相关的提示词(个人提示词和全员提示词,系统提示词不删除)
		if err := tx.Unscoped().Where("project_id = ? AND scope != ?", id, "system").Delete(&models.Prompt{}).Error; err != nil {
			return err
		}

		// 13. 最后删除项目本身(硬删除)
		if err := tx.Unscoped().Delete(&models.Project{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}
