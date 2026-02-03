package repositories

import (
	"fmt"
	"log"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// AutoTestCaseRepository 自动化测试用例仓储接口
type AutoTestCaseRepository interface {
	// 元数据管理
	GetMetadataByProjectID(projectID uint, caseType string) (*models.AutoTestCase, error)
	UpdateMetadata(projectID uint, caseType string, metadata map[string]interface{}) error

	// 用例CRUD方法 - 使用CaseID(UUID)作为主键
	Create(testCase *models.AutoTestCase) error
	GetByCaseID(caseID string) (*models.AutoTestCase, error)
	UpdateByCaseID(caseID string, updates map[string]interface{}) error
	DeleteByCaseID(caseID string) error
	DeleteByCaseGroup(projectID uint, caseType string, caseGroup string) error

	// 分页查询
	GetCasesByType(projectID uint, caseType string, offset int, limit int) ([]*models.AutoTestCase, int64, error)
	GetCasesByTypeAndGroup(projectID uint, caseType string, caseGroup string, offset int, limit int) ([]*models.AutoTestCase, int64, error)

	// 批量操作
	GetByProjectAndType(projectID uint, caseType string) ([]*models.AutoTestCase, error)
	GetByProjectAndTypeOrdered(projectID uint, caseType string) ([]*models.AutoTestCase, error)
	BatchUpdateIDsByCaseID(caseIDOrder []string) error
	GetMaxIDByProjectAndType(projectID uint, caseType string) (uint, error)

	// 新增：插入/删除后的排序辅助方法
	IncrementOrderAfter(projectID uint, caseType string, afterOrder int) error // 将指定位置后的id加1
	DecrementOrderAfter(projectID uint, caseType string, afterOrder int) error // 将指定位置后的id减1
	ReassignDisplayIDs(projectID uint, caseType string) error                  // 重新分配id字段(1,2,3...)
}

type autoTestCaseRepository struct {
	db *gorm.DB
}

// NewAutoTestCaseRepository 创建自动化测试用例仓储实例
func NewAutoTestCaseRepository(db *gorm.DB) AutoTestCaseRepository {
	return &autoTestCaseRepository{db: db}
}

// GetMetadataByProjectID 获取项目元数据(指定类型的首条记录)
func (r *autoTestCaseRepository) GetMetadataByProjectID(projectID uint, caseType string) (*models.AutoTestCase, error) {
	var testCase models.AutoTestCase
	err := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).
		Order("created_at ASC").
		First(&testCase).Error

	if err != nil {
		// 直接返回原始错误,保留gorm.ErrRecordNotFound
		return nil, err
	}
	return &testCase, nil
}

// UpdateMetadata 更新元数据字段
func (r *autoTestCaseRepository) UpdateMetadata(projectID uint, caseType string, metadata map[string]interface{}) error {
	// 查找指定类型的首条记录
	var testCase models.AutoTestCase
	err := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).
		Order("created_at ASC").
		First(&testCase).Error

	if err != nil {
		// 直接返回原始错误,保留gorm.ErrRecordNotFound
		return err
	}

	// 更新指定字段
	err = r.db.Model(&testCase).Updates(metadata).Error
	if err != nil {
		return fmt.Errorf("update metadata for project_id %d: %w", projectID, err)
	}
	return nil
}

// Create 插入单条用例记录
func (r *autoTestCaseRepository) Create(testCase *models.AutoTestCase) error {
	err := r.db.Create(testCase).Error
	if err != nil {
		return fmt.Errorf("create auto test case: %w", err)
	}
	return nil
}

// GetByCaseID 根据CaseID(UUID)查询单条用例记录
func (r *autoTestCaseRepository) GetByCaseID(caseID string) (*models.AutoTestCase, error) {
	var testCase models.AutoTestCase
	err := r.db.Where("case_id = ?", caseID).First(&testCase).Error
	if err != nil {
		return nil, err // 保留gorm.ErrRecordNotFound
	}
	return &testCase, nil
}

// UpdateByCaseID 部分更新用例字段(通过CaseID)
func (r *autoTestCaseRepository) UpdateByCaseID(caseID string, updates map[string]interface{}) error {
	err := r.db.Model(&models.AutoTestCase{}).
		Where("case_id = ?", caseID).
		Updates(updates).Error
	if err != nil {
		return fmt.Errorf("update auto test case %s: %w", caseID, err)
	}
	return nil
}

// DeleteByCaseID 硬删除用例记录(通过CaseID)
func (r *autoTestCaseRepository) DeleteByCaseID(caseID string) error {
	err := r.db.Unscoped().Where("case_id = ?", caseID).Delete(&models.AutoTestCase{}).Error
	if err != nil {
		return fmt.Errorf("delete auto test case %s: %w", caseID, err)
	}
	return nil
}

// GetCasesByType 根据用例类型获取用例列表(分页)
func (r *autoTestCaseRepository) GetCasesByType(projectID uint, caseType string, offset int, limit int) ([]*models.AutoTestCase, int64, error) {
	var cases []*models.AutoTestCase
	var total int64

	// 查询条件
	query := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType)

	// 统计总数
	if err := query.Model(&models.AutoTestCase{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count auto cases by type: %w", err)
	}

	// 查询数据(按id升序排序)
	err := query.Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&cases).Error

	if err != nil {
		return nil, 0, fmt.Errorf("get auto cases by type: %w", err)
	}

	return cases, total, nil
}

// GetCasesByTypeAndGroup 根据用例类型和用例集获取用例列表(分页)
// caseGroup: 可以是用例集名称或用例集ID
func (r *autoTestCaseRepository) GetCasesByTypeAndGroup(projectID uint, caseType string, caseGroup string, offset int, limit int) ([]*models.AutoTestCase, int64, error) {
	var cases []*models.AutoTestCase
	var total int64

	// 查询条件：项目ID、用例类型、用例集名称或ID
	query := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType)

	// 支持按用例集名称或ID查询
	// 先尝试作为ID查询(用于group_id参数)，如果查询结果为空则按名称查询
	var caseGroupID uint
	fmt.Sscanf(caseGroup, "%d", &caseGroupID)

	if caseGroupID > 0 {
		// 按用例集ID查询：先从case_groups表获取名称，再按名称过滤
		log.Printf("[GetCasesByTypeAndGroup] Looking up group_name for id=%d, project_id=%d", caseGroupID, projectID)
		var groupName string
		result := r.db.Where("id = ? AND project_id = ?", caseGroupID, projectID).
			Model(&models.CaseGroup{}).Pluck("group_name", &groupName)
		log.Printf("[GetCasesByTypeAndGroup] Lookup result - Error: %v, groupName: '%s', RowsAffected: %d", result.Error, groupName, result.RowsAffected)

		if result.Error != nil {
			// 如果按ID查询失败，按ID字符串查询
			log.Printf("[GetCasesByTypeAndGroup] Query error, using original caseGroup as filter: '%s'", caseGroup)
			query = query.Where("case_group = ?", caseGroup)
		} else if groupName != "" {
			// 成功获取到名称，使用名称过滤
			log.Printf("[GetCasesByTypeAndGroup] Using groupName as filter: '%s'", groupName)
			query = query.Where("case_group = ?", groupName)
		} else {
			// 用例集不存在（ID对应的group_name为空），不添加额外过滤
			// 这种情况下返回空结果
			log.Printf("[GetCasesByTypeAndGroup] Group not found (groupName empty), returning 0 results")
			return cases, 0, nil
		}
	} else {
		// 按用例集名称查询
		log.Printf("[GetCasesByTypeAndGroup] caseGroup not a number, using as name filter: '%s'", caseGroup)
		query = query.Where("case_group = ?", caseGroup)
	}

	// 统计总数
	if err := query.Model(&models.AutoTestCase{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count auto cases by type and group: %w", err)
	}

	// 查询数据(按id升序排序)
	log.Printf("[GetCasesByTypeAndGroup] Executing final query with offset=%d, limit=%d", offset, limit)
	err := query.Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&cases).Error

	if err != nil {
		return nil, 0, fmt.Errorf("get auto cases by type and group: %w", err)
	}

	log.Printf("[GetCasesByTypeAndGroup] Query completed - Found %d cases, Total: %d", len(cases), total)

	return cases, total, nil
}

// GetByProjectAndType 查询项目下指定类型的所有用例
func (r *autoTestCaseRepository) GetByProjectAndType(projectID uint, caseType string) ([]*models.AutoTestCase, error) {
	var cases []*models.AutoTestCase
	err := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).
		Order("id ASC").
		Find(&cases).Error
	if err != nil {
		return nil, fmt.Errorf("get auto cases by project and type: %w", err)
	}
	return cases, nil
}

// GetByProjectAndTypeOrdered 按id排序查询(用于导出)
func (r *autoTestCaseRepository) GetByProjectAndTypeOrdered(projectID uint, caseType string) ([]*models.AutoTestCase, error) {
	var cases []*models.AutoTestCase
	err := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).
		Order("id ASC").
		Find(&cases).Error
	if err != nil {
		return nil, fmt.Errorf("get auto cases ordered: %w", err)
	}
	return cases, nil
}

// BatchUpdateIDsByCaseID 根据case_id数组的顺序重新分配ID(用于ID重新生成)
// caseIDOrder: 按照期望顺序排列的case_id数组
// 将依次更新为 ID=1, ID=2, ID=3...
func (r *autoTestCaseRepository) BatchUpdateIDsByCaseID(caseIDOrder []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, caseID := range caseIDOrder {
			newID := uint(i + 1) // 从1开始
			if err := tx.Model(&models.AutoTestCase{}).
				Where("case_id = ?", caseID).
				Update("id", newID).Error; err != nil {
				return fmt.Errorf("update id for case %s to %d: %w", caseID, newID, err)
			}
		}
		return nil
	})
}

// GetMaxIDByProjectAndType 获取指定项目和类型的最大显示ID
func (r *autoTestCaseRepository) GetMaxIDByProjectAndType(projectID uint, caseType string) (uint, error) {
	var maxID uint
	err := r.db.Model(&models.AutoTestCase{}).
		Where("project_id = ? AND case_type = ?", projectID, caseType).
		Select("COALESCE(MAX(id), 0)").
		Scan(&maxID).Error
	if err != nil {
		return 0, fmt.Errorf("get max id: %w", err)
	}
	return maxID, nil
}

// IncrementOrderAfter 将指定id之后的所有用例的id加1
// 用于插入用例时调整后续用例的排序
func (r *autoTestCaseRepository) IncrementOrderAfter(projectID uint, caseType string, afterOrder int) error {
	err := r.db.Exec(
		"UPDATE auto_test_cases SET id = id + 1 WHERE project_id = ? AND case_type = ? AND id > ?",
		projectID, caseType, afterOrder,
	).Error
	if err != nil {
		return fmt.Errorf("increment id after %d: %w", afterOrder, err)
	}
	return nil
}

// DecrementOrderAfter 将指定id之后的所有用例的id减1
// 用于删除用例后调整后续用例的排序
func (r *autoTestCaseRepository) DecrementOrderAfter(projectID uint, caseType string, afterOrder int) error {
	err := r.db.Exec(
		"UPDATE auto_test_cases SET id = id - 1 WHERE project_id = ? AND case_type = ? AND id > ?",
		projectID, caseType, afterOrder,
	).Error
	if err != nil {
		return fmt.Errorf("decrement id after %d: %w", afterOrder, err)
	}
	return nil
}

// ReassignDisplayIDs 重新分配所有用例的id字段，从1开始连续递增
// 用于确保No列显示连续的数字
func (r *autoTestCaseRepository) ReassignDisplayIDs(projectID uint, caseType string) error {
	// 1. 查询所有用例，按id排序
	var cases []*models.AutoTestCase
	err := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).
		Order("id ASC").
		Find(&cases).Error
	if err != nil {
		return fmt.Errorf("query cases for reassign: %w", err)
	}

	// 2. 批量更新id字段（在事务中执行）
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, c := range cases {
			newID := uint(i + 1) // id从1开始
			if err := tx.Model(&models.AutoTestCase{}).
				Where("case_id = ?", c.CaseID).
				Update("id", newID).Error; err != nil {
				return fmt.Errorf("reassign id to %d for case %s: %w", newID, c.CaseID, err)
			}
		}
		return nil
	})
}

// DeleteByCaseGroup 删除指定用例集的所有用例(硬删除)
func (r *autoTestCaseRepository) DeleteByCaseGroup(projectID uint, caseType string, caseGroup string) error {
	result := r.db.Unscoped().Where("project_id = ? AND case_type = ? AND case_group = ?", projectID, caseType, caseGroup).
		Delete(&models.AutoTestCase{})

	if result.Error != nil {
		return fmt.Errorf("delete cases by case_group %s: %w", caseGroup, result.Error)
	}

	return nil
}
