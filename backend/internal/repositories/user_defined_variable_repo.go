package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// UserDefinedVariableRepository 用户自定义变量仓库接口
type UserDefinedVariableRepository interface {
	// GetByGroupID 获取用例集的所有变量（不包含任务变量）
	GetByGroupID(groupID uint, groupType string) ([]*models.UserDefinedVariable, error)
	// GetByTaskUUID 获取执行任务的变量
	GetByTaskUUID(taskUUID string) ([]*models.UserDefinedVariable, error)
	// Create 创建变量
	Create(variable *models.UserDefinedVariable) error
	// Update 更新变量
	Update(variable *models.UserDefinedVariable) error
	// Delete 删除变量
	Delete(id uint) error
	// DeleteByGroupID 删除用例集的所有变量
	DeleteByGroupID(groupID uint, groupType string) error
	// DeleteByTaskUUID 删除任务的所有变量
	DeleteByTaskUUID(taskUUID string) error
	// BatchUpsert 批量更新或插入变量（用例集）
	BatchUpsert(groupID uint, groupType string, projectID uint, variables []*models.UserDefinedVariable) error
	// BatchUpsertTaskVariables 批量更新或插入任务变量
	BatchUpsertTaskVariables(taskUUID string, groupID uint, groupType string, projectID uint, variables []*models.UserDefinedVariable) error
	// GetByID 根据ID获取变量
	GetByID(id uint) (*models.UserDefinedVariable, error)
}

type userDefinedVariableRepository struct {
	db *gorm.DB
}

// NewUserDefinedVariableRepository 创建仓库实例
func NewUserDefinedVariableRepository(db *gorm.DB) UserDefinedVariableRepository {
	return &userDefinedVariableRepository{db: db}
}

// GetByGroupID 获取用例集的所有变量（不包含任务变量，task_uuid为空）
func (r *userDefinedVariableRepository) GetByGroupID(groupID uint, groupType string) ([]*models.UserDefinedVariable, error) {
	var variables []*models.UserDefinedVariable
	err := r.db.Where("group_id = ? AND group_type = ? AND (task_uuid IS NULL OR task_uuid = '')", groupID, groupType).
		Order("id ASC").
		Find(&variables).Error
	return variables, err
}

// GetByTaskUUID 获取执行任务的变量
func (r *userDefinedVariableRepository) GetByTaskUUID(taskUUID string) ([]*models.UserDefinedVariable, error) {
	var variables []*models.UserDefinedVariable
	err := r.db.Where("task_uuid = ?", taskUUID).
		Order("id ASC").
		Find(&variables).Error
	return variables, err
}

// Create 创建变量
func (r *userDefinedVariableRepository) Create(variable *models.UserDefinedVariable) error {
	return r.db.Create(variable).Error
}

// Update 更新变量
func (r *userDefinedVariableRepository) Update(variable *models.UserDefinedVariable) error {
	return r.db.Save(variable).Error
}

// Delete 删除变量
func (r *userDefinedVariableRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserDefinedVariable{}, id).Error
}

// DeleteByGroupID 删除用例集的所有变量（不删除任务变量）
func (r *userDefinedVariableRepository) DeleteByGroupID(groupID uint, groupType string) error {
	return r.db.Where("group_id = ? AND group_type = ? AND (task_uuid IS NULL OR task_uuid = '')", groupID, groupType).
		Delete(&models.UserDefinedVariable{}).Error
}

// DeleteByTaskUUID 删除任务的所有变量
func (r *userDefinedVariableRepository) DeleteByTaskUUID(taskUUID string) error {
	return r.db.Where("task_uuid = ?", taskUUID).
		Delete(&models.UserDefinedVariable{}).Error
}

// BatchUpsert 批量更新或插入变量（用例集变量，task_uuid为空）
// 先删除该用例集的所有变量，再批量插入新变量
func (r *userDefinedVariableRepository) BatchUpsert(groupID uint, groupType string, projectID uint, variables []*models.UserDefinedVariable) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除该用例集的所有变量（不删除任务变量）
		if err := tx.Where("group_id = ? AND group_type = ? AND (task_uuid IS NULL OR task_uuid = '')", groupID, groupType).
			Delete(&models.UserDefinedVariable{}).Error; err != nil {
			return err
		}

		// 2. 批量插入新变量
		if len(variables) > 0 {
			for _, v := range variables {
				v.ID = 0 // 重置ID以便创建新记录
				v.GroupID = groupID
				v.GroupType = groupType
				v.ProjectID = projectID
				v.TaskUUID = "" // 确保是用例集变量
			}
			if err := tx.Create(&variables).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// BatchUpsertTaskVariables 批量更新或插入任务变量
// 先删除该任务的所有变量，再批量插入新变量
func (r *userDefinedVariableRepository) BatchUpsertTaskVariables(taskUUID string, groupID uint, groupType string, projectID uint, variables []*models.UserDefinedVariable) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除该任务的所有变量
		if err := tx.Where("task_uuid = ?", taskUUID).
			Delete(&models.UserDefinedVariable{}).Error; err != nil {
			return err
		}

		// 2. 批量插入新变量
		if len(variables) > 0 {
			for _, v := range variables {
				v.ID = 0 // 重置ID以便创建新记录
				v.GroupID = groupID
				v.GroupType = groupType
				v.ProjectID = projectID
				v.TaskUUID = taskUUID // 设置任务UUID
			}
			if err := tx.Create(&variables).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetByID 根据ID获取变量
func (r *userDefinedVariableRepository) GetByID(id uint) (*models.UserDefinedVariable, error) {
	var variable models.UserDefinedVariable
	err := r.db.First(&variable, id).Error
	if err != nil {
		return nil, err
	}
	return &variable, nil
}
