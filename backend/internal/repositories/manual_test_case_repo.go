package repositories

import (
	"fmt"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// ManualTestCaseRepository 手工测试用例仓储接口
type ManualTestCaseRepository interface {
	GetMetadataByProjectID(projectID uint, caseType string) (*models.ManualTestCase, error)
	UpdateMetadata(projectID uint, caseType string, metadata map[string]interface{}) error
	CreateDefaultMetadata(projectID uint, caseType string) error
	GetCasesByType(projectID uint, caseType string, offset int, limit int, caseGroup string) ([]*models.ManualTestCase, int64, error)

	// CRUD方法 - 使用CaseID(UUID)作为主键
	Create(testCase *models.ManualTestCase) error
	GetByCaseID(caseID string) (*models.ManualTestCase, error)                 // 改用CaseID(UUID)
	GetCaseByIntID(projectID uint, intID uint) (*models.ManualTestCase, error) // 通过整数ID查询用例
	UpdateByCaseID(caseID string, updates map[string]interface{}) error        // 改用CaseID(UUID)
	DeleteByCaseID(caseID string) error                                        // 改用CaseID(UUID)
	DeleteByCaseGroup(projectID uint, caseType string, caseGroup string) error // 级联删除用例集的所有用例

	// 批量操作方法
	CreateBatch(testCases []*models.ManualTestCase) error
	DeleteBatch(caseIDs []string) error // 改用CaseID(UUID)
	GetByProjectAndType(projectID uint, caseType string) ([]*models.ManualTestCase, error)
	BatchUpdateIDs(caseIDMap map[uint]uint) error      // ID重排：更新显示序号（用于重新排序按钮）
	BatchUpdateIDsByCaseID(caseIDOrder []string) error // ID重排：根据case_id顺序重新分配ID（用于拖拽排序）
	DeleteByCaseType(projectID uint, caseType string) error
	GetMaxIDByProjectAndType(projectID uint, caseType string) (uint, error) // 获取指定项目和类型的最大显示ID

	// 新增：导出导入支持方法
	GetMaxID(projectID uint, caseType string) (uint, error)                                       // 获取最大ID值(用于导入时生成新No)
	GetByProjectAndTypeOrdered(projectID uint, caseType string) ([]*models.ManualTestCase, error) // 按id排序查询(用于导出)

	// 新增：插入/删除后的排序辅助方法
	IncrementOrderAfter(projectID uint, caseType string, afterOrder int) error // 将指定位置后的display_order加1
	DecrementOrderAfter(projectID uint, caseType string, afterOrder int) error // 将指定位置后的display_order减1
	ReassignDisplayIDs(projectID uint, caseType string) error                  // 重新分配id字段(1,2,3...)

	// 新增：根据groupID获取case_group的名称
	GetCaseGroupName(projectID uint, groupID uint) (string, error) // 根据case_groups表的ID获取group_name
}

type manualTestCaseRepository struct {
	db *gorm.DB
}

// NewManualTestCaseRepository 创建手工测试用例仓储实例
func NewManualTestCaseRepository(db *gorm.DB) ManualTestCaseRepository {
	return &manualTestCaseRepository{db: db}
}

// GetMetadataByProjectID 获取项目元数据(指定类型的首条记录)
func (r *manualTestCaseRepository) GetMetadataByProjectID(projectID uint, caseType string) (*models.ManualTestCase, error) {
	var testCase models.ManualTestCase
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
func (r *manualTestCaseRepository) UpdateMetadata(projectID uint, caseType string, metadata map[string]interface{}) error {
	// 查找指定类型的首条记录
	var testCase models.ManualTestCase
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

// CreateDefaultMetadata 创建默认元数据记录
func (r *manualTestCaseRepository) CreateDefaultMetadata(projectID uint, caseType string) error {
	testCase := models.ManualTestCase{
		ProjectID:   projectID,
		CaseType:    caseType,
		TestVersion: "",
		TestEnv:     "",
		TestDate:    "",
		Executor:    "",
		TestResult:  "NR",
	}

	err := r.db.Create(&testCase).Error
	if err != nil {
		return fmt.Errorf("create default metadata for project_id %d, type %s: %w", projectID, caseType, err)
	}
	return nil
}

// GetCasesByType 根据用例类型获取用例列表(分页)
func (r *manualTestCaseRepository) GetCasesByType(projectID uint, caseType string, offset int, limit int, caseGroup string) ([]*models.ManualTestCase, int64, error) {
	var cases []*models.ManualTestCase
	var total int64

	// 查询条件（不再按语言筛选）
	query := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType)

	// 如果提供了 caseGroup 参数，添加过滤条件
	if caseGroup != "" {
		query = query.Where("case_group = ?", caseGroup)
	}

	// 统计总数
	if err := query.Model(&models.ManualTestCase{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count cases by type: %w", err)
	}

	// 查询数据 - 按display_id排序(已弃用id作为显示编号,改用独立的display_id字段)
	// 注意: display_id在ReassignDisplayIDs中重新分配,确保连续性
	err := query.Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&cases).Error

	if err != nil {
		return nil, 0, fmt.Errorf("get cases by type: %w", err)
	}

	return cases, total, nil
}

// Create 插入单条用例记录
func (r *manualTestCaseRepository) Create(testCase *models.ManualTestCase) error {
	err := r.db.Create(testCase).Error
	if err != nil {
		return fmt.Errorf("create test case: %w", err)
	}
	return nil
}

// GetByCaseID 根据CaseID(UUID)查询单条用例记录
func (r *manualTestCaseRepository) GetByCaseID(caseID string) (*models.ManualTestCase, error) {
	var testCase models.ManualTestCase
	err := r.db.Where("case_id = ?", caseID).First(&testCase).Error
	if err != nil {
		return nil, err // 保留gorm.ErrRecordNotFound
	}
	return &testCase, nil
}

// GetCaseByIntID 通过整数ID查询用例
func (r *manualTestCaseRepository) GetCaseByIntID(projectID uint, intID uint) (*models.ManualTestCase, error) {
	var testCase models.ManualTestCase
	err := r.db.Where("project_id = ? AND id = ?", projectID, intID).First(&testCase).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用例不存在")
		}
		return nil, err
	}
	return &testCase, nil
}

// UpdateByCaseID 部分更新用例字段（通过CaseID）
func (r *manualTestCaseRepository) UpdateByCaseID(caseID string, updates map[string]interface{}) error {
	err := r.db.Model(&models.ManualTestCase{}).
		Where("case_id = ?", caseID).
		Updates(updates).Error
	if err != nil {
		return fmt.Errorf("update test case %s: %w", caseID, err)
	}
	return nil
}

// DeleteByCaseID 硬删除用例记录（通过CaseID）
func (r *manualTestCaseRepository) DeleteByCaseID(caseID string) error {
	err := r.db.Unscoped().Where("case_id = ?", caseID).Delete(&models.ManualTestCase{}).Error
	if err != nil {
		return fmt.Errorf("delete test case %s: %w", caseID, err)
	}
	return nil
}

// CreateBatch 批量插入用例记录,使用事务保证原子性
func (r *manualTestCaseRepository) CreateBatch(testCases []*models.ManualTestCase) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, testCase := range testCases {
			if err := tx.Create(testCase).Error; err != nil {
				return fmt.Errorf("create batch test case: %w", err)
			}
		}
		return nil
	})
}

// DeleteBatch 批量硬删除用例（通过CaseID）,使用事务保证原子性
func (r *manualTestCaseRepository) DeleteBatch(caseIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Unscoped().Where("case_id IN ?", caseIDs).Delete(&models.ManualTestCase{}).Error
		if err != nil {
			return fmt.Errorf("delete batch test cases: %w", err)
		}
		return nil
	})
}

// GetByProjectAndType 查询项目下指定类型的所有用例(用于ID重排前验证)
func (r *manualTestCaseRepository) GetByProjectAndType(projectID uint, caseType string) ([]*models.ManualTestCase, error) {
	var cases []*models.ManualTestCase
	err := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).
		Order("id ASC").
		Find(&cases).Error
	if err != nil {
		return nil, fmt.Errorf("get cases by project and type: %w", err)
	}
	return cases, nil
}

// BatchUpdateIDs 批量更新ID显示序号 - 直接UPDATE（ID不再是主键）
// caseIDMap: map[oldID]newID，例如 {125:1, 126:2, 127:3}
func (r *manualTestCaseRepository) BatchUpdateIDs(caseIDMap map[uint]uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// ID现在是普通字段，可以直接更新，无需临时ID
		for oldID, newID := range caseIDMap {
			if err := tx.Model(&models.ManualTestCase{}).
				Where("id = ?", oldID).
				Update("id", newID).Error; err != nil {
				return fmt.Errorf("update display id from %d to %d: %w", oldID, newID, err)
			}
		}
		return nil
	})
}

// DeleteByCaseType 删除指定项目和用例类型的所有用例
func (r *manualTestCaseRepository) DeleteByCaseType(projectID uint, caseType string) error {
	result := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).Delete(&models.ManualTestCase{})
	if result.Error != nil {
		return fmt.Errorf("delete cases by type: %w", result.Error)
	}
	// 即使没有找到用例(RowsAffected == 0),也返回成功,因为目标状态(没有该类型用例)已达成
	return nil
}

// GetMaxIDByProjectAndType 获取指定项目和类型的最大显示ID
func (r *manualTestCaseRepository) GetMaxIDByProjectAndType(projectID uint, caseType string) (uint, error) {
	var maxID uint
	err := r.db.Model(&models.ManualTestCase{}).
		Where("project_id = ? AND case_type = ?", projectID, caseType).
		Select("COALESCE(MAX(id), 0)").
		Scan(&maxID).Error
	if err != nil {
		return 0, fmt.Errorf("get max id: %w", err)
	}
	return maxID, nil
}

// GetMaxID 获取最大ID值(用于导入时生成新No)
func (r *manualTestCaseRepository) GetMaxID(projectID uint, caseType string) (uint, error) {
	return r.GetMaxIDByProjectAndType(projectID, caseType)
}

// GetByProjectAndTypeOrdered 按id排序查询(用于导出)
func (r *manualTestCaseRepository) GetByProjectAndTypeOrdered(projectID uint, caseType string) ([]*models.ManualTestCase, error) {
	var cases []*models.ManualTestCase
	err := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).
		Order("id ASC").
		Find(&cases).Error
	if err != nil {
		return nil, fmt.Errorf("get cases ordered: %w", err)
	}
	return cases, nil
}

// BatchUpdateIDsByCaseID 根据case_id数组的顺序重新分配ID（用于拖拽排序）
// caseIDOrder: 按照期望顺序排列的case_id数组，例如 [uuid1, uuid2, uuid3]
// 将依次更新为 ID=1, ID=2, ID=3...
func (r *manualTestCaseRepository) BatchUpdateIDsByCaseID(caseIDOrder []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, caseID := range caseIDOrder {
			newID := uint(i + 1) // 从1开始
			if err := tx.Model(&models.ManualTestCase{}).
				Where("case_id = ?", caseID).
				Update("id", newID).Error; err != nil {
				return fmt.Errorf("update id for case %s to %d: %w", caseID, newID, err)
			}
		}
		return nil
	})
}

// IncrementOrderAfter 将指定display_order之后的所有用例的display_order加1
// 用于插入用例时调整后续用例的排序
func (r *manualTestCaseRepository) IncrementOrderAfter(projectID uint, caseType string, afterOrder int) error {
	err := r.db.Exec(
		"UPDATE manual_test_cases SET id = id + 1 WHERE project_id = ? AND case_type = ? AND id > ?",
		projectID, caseType, afterOrder,
	).Error
	if err != nil {
		return fmt.Errorf("increment id after %d: %w", afterOrder, err)
	}
	return nil
}

// DecrementOrderAfter 将指定display_order之后的所有用例的display_order减1
// 用于删除用例后调整后续用例的排序
func (r *manualTestCaseRepository) DecrementOrderAfter(projectID uint, caseType string, afterOrder int) error {
	err := r.db.Exec(
		"UPDATE manual_test_cases SET id = id - 1 WHERE project_id = ? AND case_type = ? AND id > ?",
		projectID, caseType, afterOrder,
	).Error
	if err != nil {
		return fmt.Errorf("decrement id after %d: %w", afterOrder, err)
	}
	return nil
}

// ReassignDisplayIDs 重新分配所有用例的id字段，从1开始连续递增
// 用于确保No列显示连续的数字
func (r *manualTestCaseRepository) ReassignDisplayIDs(projectID uint, caseType string) error {
	// 1. 查询所有用例，按id排序
	var cases []*models.ManualTestCase
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
			if err := tx.Model(&models.ManualTestCase{}).
				Where("case_id = ?", c.CaseID).
				Update("id", newID).Error; err != nil {
				return fmt.Errorf("reassign id to %d for case %s: %w", newID, c.CaseID, err)
			}
		}
		return nil
	})
}

// GetCaseGroupName 根据case_groups表的ID获取group_name
func (r *manualTestCaseRepository) GetCaseGroupName(projectID uint, groupID uint) (string, error) {
	// 需要查询case_groups表，因此需要定义一个简单的结构体
	type CaseGroup struct {
		GroupName string
	}

	var cg CaseGroup
	err := r.db.Table("case_groups").
		Select("group_name").
		Where("id = ? AND project_id = ? AND deleted_at IS NULL", groupID, projectID).
		First(&cg).Error

	if err != nil {
		return "", fmt.Errorf("get case group name: %w", err)
	}

	return cg.GroupName, nil
}

// DeleteByCaseGroup 删除指定用例集的所有用例(硬删除)
func (r *manualTestCaseRepository) DeleteByCaseGroup(projectID uint, caseType string, caseGroup string) error {
	result := r.db.Unscoped().Where("project_id = ? AND case_type = ? AND case_group = ?", projectID, caseType, caseGroup).
		Delete(&models.ManualTestCase{})

	if result.Error != nil {
		return fmt.Errorf("delete manual cases by case_group %s: %w", caseGroup, result.Error)
	}

	return nil
}
