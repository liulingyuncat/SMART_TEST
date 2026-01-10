package repositories

import (
	"errors"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// CaseGroupRepository 用例集数据访问层
type CaseGroupRepository struct {
	DB *gorm.DB
}

// NewCaseGroupRepository 创建用例集Repository
func NewCaseGroupRepository(db *gorm.DB) *CaseGroupRepository {
	return &CaseGroupRepository{DB: db}
}

// GetByProjectAndType 获取项目下指定类型的所有用例集（自动过滤软删除记录）
func (r *CaseGroupRepository) GetByProjectAndType(projectID uint, caseType string) ([]models.CaseGroup, error) {
	var groups []models.CaseGroup
	// GORM会自动添加 deleted_at IS NULL 条件过滤软删除记录
	err := r.DB.Where("project_id = ? AND case_type = ?", projectID, caseType).
		Order("display_order ASC, created_at ASC").
		Find(&groups).Error
	return groups, err
}

// GetByID 根据ID获取用例集
func (r *CaseGroupRepository) GetByID(id uint) (*models.CaseGroup, error) {
	var group models.CaseGroup
	err := r.DB.First(&group, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// GetByName 根据名称获取用例集（自动过滤软删除记录）
func (r *CaseGroupRepository) GetByName(projectID uint, caseType, groupName string) (*models.CaseGroup, error) {
	var group models.CaseGroup
	// GORM会自动添加 deleted_at IS NULL 条件过滤软删除记录
	err := r.DB.Where("project_id = ? AND case_type = ? AND group_name = ?", projectID, caseType, groupName).
		First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// Create 创建用例集
func (r *CaseGroupRepository) Create(group *models.CaseGroup) error {
	return r.DB.Create(group).Error
}

// Update 更新用例集
func (r *CaseGroupRepository) Update(group *models.CaseGroup) error {
	// 开启事务
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 先获取原始数据
		var original models.CaseGroup
		if err := tx.First(&original, group.ID).Error; err != nil {
			return err
		}

		// 如果用例集名称发生变化，需要同步更新所有相关用例的 case_group 字段
		if group.GroupName != "" && group.GroupName != original.GroupName {
			// 更新 manual_test_cases 表中的 case_group 字段
			if err := tx.Model(&models.ManualTestCase{}).
				Where("project_id = ? AND case_type = ? AND case_group = ?",
					original.ProjectID, original.CaseType, original.GroupName).
				Update("case_group", group.GroupName).Error; err != nil {
				return err
			}
		}

		// 更新用例集记录
		return tx.Save(group).Error
	})
}

// Delete 删除用例集（软删除）
func (r *CaseGroupRepository) Delete(id uint) error {
	return r.DB.Delete(&models.CaseGroup{}, id).Error
}

// HardDelete 物理删除用例集
func (r *CaseGroupRepository) HardDelete(id uint) error {
	return r.DB.Unscoped().Delete(&models.CaseGroup{}, id).Error
}

// CreateIfNotExists 如果不存在则创建用例集
func (r *CaseGroupRepository) CreateIfNotExists(projectID uint, caseType, groupName string) (*models.CaseGroup, error) {
	// 先查询是否存在
	existing, err := r.GetByName(projectID, caseType, groupName)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return existing, nil
	}

	// 不存在则创建
	group := &models.CaseGroup{
		ProjectID: projectID,
		CaseType:  caseType,
		GroupName: groupName,
	}

	if err := r.Create(group); err != nil {
		return nil, err
	}

	return group, nil
}
