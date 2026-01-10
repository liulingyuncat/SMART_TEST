package repositories

import (
	"fmt"
	"strings"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// ExecutionTaskRepository 测试执行任务仓储接口
type ExecutionTaskRepository interface {
	// 查询方法
	GetByProject(projectID uint) ([]*models.ExecutionTask, error)
	GetByUUID(taskUUID string) (*models.ExecutionTask, error)

	// CRUD方法
	Create(task *models.ExecutionTask) error
	UpdateByUUID(taskUUID string, updates map[string]interface{}) error
	Delete(taskUUID string) error

	// 业务逻辑辅助方法
	CheckNameUnique(projectID uint, taskName string, excludeUUID string) (bool, error)
}

type executionTaskRepository struct {
	db *gorm.DB
}

// NewExecutionTaskRepository 创建测试执行任务仓储实例
func NewExecutionTaskRepository(db *gorm.DB) ExecutionTaskRepository {
	return &executionTaskRepository{db: db}
}

// GetByProject 获取项目的所有任务(按状态和创建时间排序)
func (r *executionTaskRepository) GetByProject(projectID uint) ([]*models.ExecutionTask, error) {
	var tasks []*models.ExecutionTask

	// 按状态优先级排序: in_progress(1) > pending(2) > completed(3)
	// 同状态内按创建时间倒序
	// Where会自动过滤deleted_at IS NULL(GORM软删除)
	err := r.db.Where("project_id = ? AND deleted_at IS NULL", projectID).
		Order(`
			CASE task_status 
				WHEN 'in_progress' THEN 1 
				WHEN 'pending' THEN 2 
				WHEN 'completed' THEN 3 
			END ASC, 
			created_at DESC
		`).
		Find(&tasks).Error

	if err != nil {
		return nil, fmt.Errorf("get tasks by project %d: %w", projectID, err)
	}
	return tasks, nil
}

// GetByUUID 根据UUID查询单条任务记录
func (r *executionTaskRepository) GetByUUID(taskUUID string) (*models.ExecutionTask, error) {
	var task models.ExecutionTask
	// Where会自动过滤deleted_at IS NULL(GORM软删除)
	err := r.db.Where("task_uuid = ? AND deleted_at IS NULL", taskUUID).First(&task).Error
	if err != nil {
		// 保留gorm.ErrRecordNotFound
		return nil, err
	}
	return &task, nil
}

// Create 插入单条任务记录
func (r *executionTaskRepository) Create(task *models.ExecutionTask) error {
	err := r.db.Create(task).Error
	if err != nil {
		return fmt.Errorf("create execution task: %w", err)
	}
	return nil
}

// UpdateByUUID 根据UUID更新任务字段
func (r *executionTaskRepository) UpdateByUUID(taskUUID string, updates map[string]interface{}) error {
	result := r.db.Model(&models.ExecutionTask{}).
		Where("task_uuid = ?", taskUUID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("update task %s: %w", taskUUID, result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete 删除任务(硬删除)
func (r *executionTaskRepository) Delete(taskUUID string) error {
	result := r.db.Unscoped().
		Where("task_uuid = ?", taskUUID).
		Delete(&models.ExecutionTask{})

	if result.Error != nil {
		return fmt.Errorf("delete task %s: %w", taskUUID, result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// CheckNameUnique 检查任务名唯一性(不区分大小写,排除已删除和指定UUID)
func (r *executionTaskRepository) CheckNameUnique(projectID uint, taskName string, excludeUUID string) (bool, error) {
	var count int64

	query := r.db.Model(&models.ExecutionTask{}).
		Where("project_id = ?", projectID).
		Where("LOWER(task_name) = ?", strings.ToLower(taskName)).
		Where("deleted_at IS NULL")

	// 排除指定的UUID(用于更新时的唯一性检查)
	if excludeUUID != "" {
		query = query.Where("task_uuid != ?", excludeUUID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check name unique for project %d: %w", projectID, err)
	}

	// count == 0 表示名称唯一
	return count == 0, nil
}
