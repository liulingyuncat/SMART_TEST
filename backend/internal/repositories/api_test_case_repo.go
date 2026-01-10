package repositories

import (
	"fmt"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// ApiTestCaseRepository 接口测试用例仓储接口
type ApiTestCaseRepository interface {
	// 用例CRUD方法 - 使用ID(UUID)作为主键
	Create(testCase *models.ApiTestCase) error
	GetByID(id string) (*models.ApiTestCase, error)
	Update(id string, updates map[string]interface{}) error
	Delete(id string) error

	// 分页查询
	List(projectID uint, caseType string, offset int, limit int) ([]*models.ApiTestCase, int64, error)
	// 按用例集筛选的分页查询
	ListByGroup(projectID uint, caseType string, caseGroup string, offset int, limit int) ([]*models.ApiTestCase, int64, error)

	// 批量操作
	GetByProjectAndType(projectID uint, caseType string) ([]*models.ApiTestCase, error)
	GetCaseGroups(projectID uint) ([]string, error)
	GetByProjectAndGroup(projectID uint, caseGroup string) ([]*models.ApiTestCase, error)
	DeleteByCaseGroup(projectID uint, caseGroup string) error // 硬删除指定用例集的所有用例

	// 插入/删除后的排序辅助方法
	IncrementOrderAfter(projectID uint, caseType string, caseGroup string, afterOrder int) error
	ReassignDisplayOrders(projectID uint, caseType string, caseGroup string) error

	// 版本管理
	CreateVersion(version *models.ApiTestCaseVersion) error
	GetVersionByID(versionID string) (*models.ApiTestCaseVersion, error)
	ListVersions(projectID uint, offset int, limit int) ([]*models.ApiTestCaseVersion, int64, error)
	DeleteVersion(versionID string) error
	UpdateVersionRemark(versionID string, remark string) error
}

type apiTestCaseRepository struct {
	db *gorm.DB
}

// NewApiTestCaseRepository 创建接口测试用例仓储实例
func NewApiTestCaseRepository(db *gorm.DB) ApiTestCaseRepository {
	return &apiTestCaseRepository{db: db}
}

// Create 创建用例(UUID自动生成)
func (r *apiTestCaseRepository) Create(testCase *models.ApiTestCase) error {
	err := r.db.Create(testCase).Error
	if err != nil {
		return fmt.Errorf("create api test case: %w", err)
	}
	return nil
}

// GetByID 根据UUID主键查询
func (r *apiTestCaseRepository) GetByID(id string) (*models.ApiTestCase, error) {
	var testCase models.ApiTestCase
	err := r.db.Where("id = ?", id).First(&testCase).Error
	if err != nil {
		return nil, err // 保留gorm.ErrRecordNotFound
	}
	return &testCase, nil
}

// Update 更新用例字段
func (r *apiTestCaseRepository) Update(id string, updates map[string]interface{}) error {
	err := r.db.Model(&models.ApiTestCase{}).
		Where("id = ?", id).
		Updates(updates).Error
	if err != nil {
		return fmt.Errorf("update api test case %s: %w", id, err)
	}
	return nil
}

// Delete 物理删除用例
func (r *apiTestCaseRepository) Delete(id string) error {
	err := r.db.Unscoped().Where("id = ?", id).Delete(&models.ApiTestCase{}).Error
	if err != nil {
		return fmt.Errorf("delete api test case %s: %w", id, err)
	}
	return nil
}

// List 分页查询(按display_order升序)
func (r *apiTestCaseRepository) List(projectID uint, caseType string, offset int, limit int) ([]*models.ApiTestCase, int64, error) {
	var cases []*models.ApiTestCase
	var total int64

	query := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType)

	// 统计总数
	if err := query.Model(&models.ApiTestCase{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count api cases: %w", err)
	}

	// 查询数据(按display_order升序)
	err := query.Order("display_order ASC").
		Offset(offset).
		Limit(limit).
		Find(&cases).Error

	if err != nil {
		return nil, 0, fmt.Errorf("list api cases: %w", err)
	}

	return cases, total, nil
}

// ListByGroup 按用例集筛选的分页查询(按display_order升序)
func (r *apiTestCaseRepository) ListByGroup(projectID uint, caseType string, caseGroup string, offset int, limit int) ([]*models.ApiTestCase, int64, error) {
	var cases []*models.ApiTestCase
	var total int64

	query := r.db.Where("project_id = ? AND case_type = ? AND case_group = ?", projectID, caseType, caseGroup)

	// 统计总数
	if err := query.Model(&models.ApiTestCase{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count api cases by group: %w", err)
	}

	// 查询数据(按display_order升序)
	err := query.Order("display_order ASC").
		Offset(offset).
		Limit(limit).
		Find(&cases).Error

	if err != nil {
		return nil, 0, fmt.Errorf("list api cases by group: %w", err)
	}

	return cases, total, nil
}

// GetByProjectAndType 查询所有用例(按display_order升序)
func (r *apiTestCaseRepository) GetByProjectAndType(projectID uint, caseType string) ([]*models.ApiTestCase, error) {
	var cases []*models.ApiTestCase
	err := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType).
		Order("display_order ASC").
		Find(&cases).Error
	if err != nil {
		return nil, fmt.Errorf("get api cases by type: %w", err)
	}
	return cases, nil
}

// GetCaseGroups 获取项目的所有用例集名称(去重)
func (r *apiTestCaseRepository) GetCaseGroups(projectID uint) ([]string, error) {
	var groups []string
	err := r.db.Model(&models.ApiTestCase{}).
		Where("project_id = ? AND case_group != ''", projectID).
		Distinct("case_group").
		Pluck("case_group", &groups).Error
	if err != nil {
		return nil, fmt.Errorf("get case groups: %w", err)
	}
	return groups, nil
}

// GetByProjectAndGroup 查询指定用例集的所有用例(按display_order升序)
func (r *apiTestCaseRepository) GetByProjectAndGroup(projectID uint, caseGroup string) ([]*models.ApiTestCase, error) {
	var cases []*models.ApiTestCase
	err := r.db.Where("project_id = ? AND case_group = ?", projectID, caseGroup).
		Order("display_order ASC").
		Find(&cases).Error
	if err != nil {
		return nil, fmt.Errorf("get api cases by group: %w", err)
	}
	return cases, nil
}

// IncrementOrderAfter 批量更新display_order+1 (按case_group筛选)
func (r *apiTestCaseRepository) IncrementOrderAfter(projectID uint, caseType string, caseGroup string, afterOrder int) error {
	query := "UPDATE api_test_cases SET display_order = display_order + 1 WHERE project_id = ? AND case_type = ? AND display_order > ?"
	args := []interface{}{projectID, caseType, afterOrder}

	// 如果指定了case_group，添加筛选条件
	if caseGroup != "" {
		query += " AND case_group = ?"
		args = append(args, caseGroup)
	}

	err := r.db.Exec(query, args...).Error
	if err != nil {
		return fmt.Errorf("increment display_order after %d: %w", afterOrder, err)
	}
	return nil
}

// ReassignDisplayOrders 重新分配display_order为1,2,3... (按case_group筛选)
func (r *apiTestCaseRepository) ReassignDisplayOrders(projectID uint, caseType string, caseGroup string) error {
	// 查询用例,按display_order排序
	var cases []*models.ApiTestCase
	query := r.db.Where("project_id = ? AND case_type = ?", projectID, caseType)

	// 如果指定了case_group，添加筛选条件
	if caseGroup != "" {
		query = query.Where("case_group = ?", caseGroup)
	}

	err := query.Order("display_order ASC").Find(&cases).Error
	if err != nil {
		return fmt.Errorf("query cases for reassign: %w", err)
	}

	// 批量更新display_order(事务)
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, c := range cases {
			newOrder := i + 1 // 从1开始
			if err := tx.Model(&models.ApiTestCase{}).
				Where("id = ?", c.ID).
				Update("display_order", newOrder).Error; err != nil {
				return fmt.Errorf("reassign display_order to %d for case %s: %w", newOrder, c.ID, err)
			}
		}
		return nil
	})
}

// ========== 版本管理 ==========

// CreateVersion 创建版本记录
func (r *apiTestCaseRepository) CreateVersion(version *models.ApiTestCaseVersion) error {
	err := r.db.Create(version).Error
	if err != nil {
		return fmt.Errorf("create api version: %w", err)
	}
	return nil
}

// GetVersionByID 根据UUID查询版本
func (r *apiTestCaseRepository) GetVersionByID(versionID string) (*models.ApiTestCaseVersion, error) {
	var version models.ApiTestCaseVersion
	err := r.db.Where("id = ?", versionID).First(&version).Error
	if err != nil {
		return nil, err // 保留gorm.ErrRecordNotFound
	}
	return &version, nil
}

// ListVersions 分页查询版本列表(按创建时间倒序)
func (r *apiTestCaseRepository) ListVersions(projectID uint, offset int, limit int) ([]*models.ApiTestCaseVersion, int64, error) {
	var versions []*models.ApiTestCaseVersion
	var total int64

	query := r.db.Where("project_id = ?", projectID)

	// 统计总数
	if err := query.Model(&models.ApiTestCaseVersion{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count api versions: %w", err)
	}

	// 查询数据(按created_at倒序)
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&versions).Error

	if err != nil {
		return nil, 0, fmt.Errorf("list api versions: %w", err)
	}

	return versions, total, nil
}

// DeleteVersion 删除版本记录(硬删除)
func (r *apiTestCaseRepository) DeleteVersion(versionID string) error {
	err := r.db.Unscoped().Where("id = ?", versionID).Delete(&models.ApiTestCaseVersion{}).Error
	if err != nil {
		return fmt.Errorf("delete api version %s: %w", versionID, err)
	}
	return nil
}

// UpdateVersionRemark 更新版本备注
func (r *apiTestCaseRepository) UpdateVersionRemark(versionID string, remark string) error {
	err := r.db.Model(&models.ApiTestCaseVersion{}).
		Where("id = ?", versionID).
		Update("remark", remark).Error
	if err != nil {
		return fmt.Errorf("update version remark %s: %w", versionID, err)
	}
	return nil
}

// DeleteByCaseGroup 硬删除指定用例集的所有用例
func (r *apiTestCaseRepository) DeleteByCaseGroup(projectID uint, caseGroup string) error {
	result := r.db.Unscoped().
		Where("project_id = ? AND case_group = ?", projectID, caseGroup).
		Delete(&models.ApiTestCase{})
	if result.Error != nil {
		return fmt.Errorf("delete api cases by case_group %s: %w", caseGroup, result.Error)
	}
	return nil
}
