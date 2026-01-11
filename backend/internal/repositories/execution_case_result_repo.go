package repositories

import (
	"fmt"
	"webtest/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ExecutionCaseResultRepository 测试执行用例结果仓储接口
type ExecutionCaseResultRepository interface {
	// 查询方法
	GetByTaskUUID(taskUUID string) ([]*models.ExecutionCaseResult, error)
	GetByCaseID(caseID string) ([]*models.ExecutionCaseResult, error)
	GetByID(id uint) (*models.ExecutionCaseResult, error)

	// 批量操作
	BatchCreate(results []*models.ExecutionCaseResult) error
	BatchUpsert(results []*models.ExecutionCaseResult) error

	// 更新和删除
	UpdateResult(id uint, updates map[string]interface{}) error
	DeleteByTaskUUID(taskUUID string) error

	// 统计方法
	GetStatistics(taskUUID string) (map[string]int, error)
}

type executionCaseResultRepository struct {
	db *gorm.DB
}

// NewExecutionCaseResultRepository 创建仓储实例
func NewExecutionCaseResultRepository(db *gorm.DB) ExecutionCaseResultRepository {
	return &executionCaseResultRepository{db: db}
}

// GetByTaskUUID 获取任务的所有执行结果(按display_id升序)
func (r *executionCaseResultRepository) GetByTaskUUID(taskUUID string) ([]*models.ExecutionCaseResult, error) {
	var results []*models.ExecutionCaseResult
	err := r.db.Where("task_uuid = ?", taskUUID).
		Order("display_id ASC, case_id ASC").
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("get results by task_uuid %s: %w", taskUUID, err)
	}
	return results, nil
}

// GetByCaseID 获取指定用例的所有执行结果
func (r *executionCaseResultRepository) GetByCaseID(caseID string) ([]*models.ExecutionCaseResult, error) {
	var results []*models.ExecutionCaseResult
	err := r.db.Where("case_id = ?", caseID).
		Order("created_at DESC").
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("get results by case_id %s: %w", caseID, err)
	}
	return results, nil
}

// GetByID 根据ID获取单个执行结果
func (r *executionCaseResultRepository) GetByID(id uint) (*models.ExecutionCaseResult, error) {
	var result models.ExecutionCaseResult
	err := r.db.First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchCreate 批量插入执行结果(使用事务)
func (r *executionCaseResultRepository) BatchCreate(results []*models.ExecutionCaseResult) error {
	if len(results) == 0 {
		return nil
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&results).Error
	})

	if err != nil {
		return fmt.Errorf("batch create %d results: %w", len(results), err)
	}
	return nil
}

// BatchUpsert 批量更新或插入(使用ON CONFLICT)
func (r *executionCaseResultRepository) BatchUpsert(results []*models.ExecutionCaseResult) error {
	if len(results) == 0 {
		return nil
	}

	// 使用Clauses配置ON CONFLICT行为
	// idx_task_case是唯一索引(task_uuid, case_id)
	err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "task_uuid"},
			{Name: "case_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"display_id",
			"case_num",
			"case_type",
			"test_result",
			"bug_id",
			"remark",
			// AI Web 字段
			"screen_cn",
			"screen_jp",
			"screen_en",
			"function_cn",
			"function_jp",
			"function_en",
			// 手工测试字段
			"major_function_cn",
			"major_function_jp",
			"major_function_en",
			"middle_function_cn",
			"middle_function_jp",
			"middle_function_en",
			"minor_function_cn",
			"minor_function_jp",
			"minor_function_en",
			// 通用字段
			"precondition_cn",
			"precondition_jp",
			"precondition_en",
			"test_steps_cn",
			"test_steps_jp",
			"test_steps_en",
			"expected_result_cn",
			"expected_result_jp",
			"expected_result_en",
			// API 用例特有字段
			"screen",
			"url",
			"header",
			"method",
			"body",
			"response",
			"updated_by",
			"updated_at",
		}),
	}).Create(&results).Error

	if err != nil {
		return fmt.Errorf("batch upsert %d results: %w", len(results), err)
	}
	return nil
}

// UpdateResult 更新单条记录的字段
// 先查询记录再更新，以确保 GORM 钩子能正确验证
func (r *executionCaseResultRepository) UpdateResult(id uint, updates map[string]interface{}) error {
	// 先查询现有记录
	var existingResult models.ExecutionCaseResult
	if err := r.db.First(&existingResult, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		return fmt.Errorf("find result %d: %w", id, err)
	}

	// 应用更新到现有记录
	if testResult, ok := updates["test_result"].(string); ok {
		existingResult.TestResult = testResult
	}
	if remark, ok := updates["remark"].(string); ok {
		existingResult.Remark = remark
	}
	if bugID, ok := updates["bug_id"].(string); ok {
		existingResult.BugID = bugID
	}
	if responseTime, ok := updates["response_time"].(string); ok {
		existingResult.ResponseTime = responseTime
	}
	if updatedBy, ok := updates["updated_by"].(uint); ok {
		existingResult.UpdatedBy = updatedBy
	}

	// 保存更新（会触发 BeforeUpdate 钩子进行验证）
	result := r.db.Save(&existingResult)
	if result.Error != nil {
		return fmt.Errorf("update result %d: %w", id, result.Error)
	}

	return nil
}

// DeleteByTaskUUID 删除任务的所有执行结果(硬删除)
func (r *executionCaseResultRepository) DeleteByTaskUUID(taskUUID string) error {
	result := r.db.Unscoped().
		Where("task_uuid = ?", taskUUID).
		Delete(&models.ExecutionCaseResult{})

	if result.Error != nil {
		return fmt.Errorf("delete results by task_uuid %s: %w", taskUUID, result.Error)
	}

	return nil
}

// GetStatistics 统计任务的用例总数和各状态数量
func (r *executionCaseResultRepository) GetStatistics(taskUUID string) (map[string]int, error) {
	type StatResult struct {
		Total      int `gorm:"column:total"`
		NRCount    int `gorm:"column:nr_count"`
		OKCount    int `gorm:"column:ok_count"`
		NGCount    int `gorm:"column:ng_count"`
		BlockCount int `gorm:"column:block_count"`
	}

	var stat StatResult
	err := r.db.Model(&models.ExecutionCaseResult{}).
		Select(`
			COUNT(*) as total,
			SUM(CASE WHEN test_result = 'NR' THEN 1 ELSE 0 END) as nr_count,
			SUM(CASE WHEN test_result = 'OK' THEN 1 ELSE 0 END) as ok_count,
			SUM(CASE WHEN test_result = 'NG' THEN 1 ELSE 0 END) as ng_count,
			SUM(CASE WHEN test_result = 'Block' THEN 1 ELSE 0 END) as block_count
		`).
		Where("task_uuid = ?", taskUUID).
		Scan(&stat).Error

	if err != nil {
		return nil, fmt.Errorf("get statistics for task_uuid %s: %w", taskUUID, err)
	}

	return map[string]int{
		"total":       stat.Total,
		"nr_count":    stat.NRCount,
		"ok_count":    stat.OKCount,
		"ng_count":    stat.NGCount,
		"block_count": stat.BlockCount,
	}, nil
}
