package repositories

import (
	"fmt"
	"log"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// DefectRepository 缺陷仓储接口
type DefectRepository interface {
	// CRUD方法
	Create(defect *models.Defect) error
	GetByID(id string) (*models.Defect, error)
	GetByDefectID(defectID string) (*models.Defect, error)
	Update(id string, updates map[string]interface{}) error
	Delete(id string) error

	// 列表查询
	List(projectID uint, status, keyword string, page, size int) ([]*models.Defect, int64, error)
	GetStatusCounts(projectID uint) (map[string]int64, error)

	// 辅助方法
	GetMaxDefectSeq(projectID uint) (int, error)
	GenerateNextDefectID(projectID uint) (string, error)
	GetByProjectID(projectID uint) ([]*models.Defect, error)
	CleanupDuplicateDefects(projectID uint, defectID string) error
	GetDB() *gorm.DB
}

type defectRepository struct {
	db *gorm.DB
}

// NewDefectRepository 创建缺陷仓储实例
func NewDefectRepository(db *gorm.DB) DefectRepository {
	return &defectRepository{db: db}
}

// Create 创建缺陷
func (r *defectRepository) Create(defect *models.Defect) error {
	err := r.db.Create(defect).Error
	if err != nil {
		return fmt.Errorf("create defect: %w", err)
	}
	return nil
}

// GetByID 根据UUID获取缺陷
func (r *defectRepository) GetByID(id string) (*models.Defect, error) {
	var defect models.Defect
	err := r.db.Preload("Attachments", "deleted_at IS NULL").
		Where("id = ?", id).First(&defect).Error
	if err != nil {
		return nil, err // 保留gorm.ErrRecordNotFound
	}
	return &defect, nil
}

// GetByDefectID 根据显示ID获取缺陷
func (r *defectRepository) GetByDefectID(defectID string) (*models.Defect, error) {
	var defect models.Defect
	err := r.db.Preload("Attachments", "deleted_at IS NULL").
		Preload("CreatedByUser").
		Where("defect_id = ?", defectID).First(&defect).Error
	if err != nil {
		return nil, err // 保留gorm.ErrRecordNotFound
	}
	return &defect, nil
}

// Update 更新缺陷
func (r *defectRepository) Update(id string, updates map[string]interface{}) error {
	result := r.db.Model(&models.Defect{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("update defect %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete 硬删除缺陷
func (r *defectRepository) Delete(id string) error {
	result := r.db.Unscoped().Where("id = ?", id).Delete(&models.Defect{})
	if result.Error != nil {
		return fmt.Errorf("delete defect %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// List 分页查询缺陷列表（支持keyword检索）
func (r *defectRepository) List(projectID uint, status, keyword string, page, size int) ([]*models.Defect, int64, error) {
	var defects []*models.Defect
	var total int64

	query := r.db.Model(&models.Defect{}).Where("project_id = ?", projectID)

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 关键词检索（ID、标题、描述、恢复方法）
	if keyword != "" {
		keywordPattern := "%" + keyword + "%"
		query = query.Where(
			"defect_id LIKE ? OR title LIKE ? OR description LIKE ? OR recovery_method LIKE ?",
			keywordPattern, keywordPattern, keywordPattern, keywordPattern,
		)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count defects: %w", err)
	}

	// 分页查询
	offset := (page - 1) * size
	err := query.Preload("CreatedByUser").
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&defects).Error

	if err != nil {
		return nil, 0, fmt.Errorf("list defects: %w", err)
	}

	return defects, total, nil
}

// GetStatusCounts 获取各状态的缺陷数量
func (r *defectRepository) GetStatusCounts(projectID uint) (map[string]int64, error) {
	type statusCount struct {
		Status string
		Count  int64
	}

	var results []statusCount
	err := r.db.Model(&models.Defect{}).
		Select("status, count(*) as count").
		Where("project_id = ?", projectID).
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("get status counts: %w", err)
	}

	counts := map[string]int64{
		"New":      0,
		"Active":   0,
		"Resolved": 0,
		"Closed":   0,
	}

	for _, r := range results {
		counts[r.Status] = r.Count
	}

	return counts, nil
}

// GetMaxDefectSeq 获取全局最大缺陷序号（包括已删除的记录，用于ID生成的递增）
// 查询整个数据库的最大ID，所有Bug共享一份ID序列
func (r *defectRepository) GetMaxDefectSeq(projectID uint) (int, error) {
	var maxSeq int
	err := r.db.Model(&models.Defect{}).
		Select("COALESCE(MAX(CAST(defect_id AS INTEGER)), 0)").
		Scan(&maxSeq).Error

	if err != nil {
		return 0, fmt.Errorf("get max defect seq: %w", err)
	}

	return maxSeq, nil
}

// GenerateNextDefectID 原子性地生成下一个缺陷ID（用事务保证一致性）
// 查询整个数据库的最大ID（不限于项目），所有Bug共享一份ID序列
// 包括所有记录（包括软删除）的最大值，确保ID永远不会重复
func (r *defectRepository) GenerateNextDefectID(projectID uint) (string, error) {
	var nextSeq int

	err := r.db.Transaction(func(txn *gorm.DB) error {
		var maxSeq int
		// 在事务中读取最大值（包括软删除的记录）
		// 使用 Unscoped() 确保包含软删除的记录
		// 直接用 CAST(defect_id AS INTEGER) 转换，不用 NULLIF
		// 注意：这里不限制 project_id，查询的是整个数据库的最大ID
		result := txn.Unscoped().Model(&models.Defect{}).
			Select("COALESCE(MAX(CAST(defect_id AS INTEGER)), 0)").
			Scan(&maxSeq)

		if result.Error != nil {
			log.Printf("[GenerateNextDefectID] Error querying max defect_id: %v", result.Error)
			return result.Error
		}

		nextSeq = maxSeq + 1
		log.Printf("[GenerateNextDefectID] project_id=%d, global maxSeq=%d, nextSeq=%d", projectID, maxSeq, nextSeq)
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("generate next defect id: %w", err)
	}

	return fmt.Sprintf("%06d", nextSeq), nil
}

// GetByProjectID 获取项目所有缺陷（用于导出）
func (r *defectRepository) GetByProjectID(projectID uint) ([]*models.Defect, error) {
	var defects []*models.Defect
	err := r.db.Where("project_id = ?", projectID).
		Order("defect_id ASC").
		Find(&defects).Error

	if err != nil {
		return nil, fmt.Errorf("get defects by project: %w", err)
	}

	return defects, nil
}

// GetDB 获取数据库实例
func (r *defectRepository) GetDB() *gorm.DB {
	return r.db
}

// CleanupDuplicateDefects 删除指定项目中重复的 defect_id（保留最早的，删除后创建的）
func (r *defectRepository) CleanupDuplicateDefects(projectID uint, defectID string) error {
	// 查找所有相同 defect_id 的记录
	var defects []models.Defect
	err := r.db.Unscoped().
		Where("project_id = ? AND defect_id = ?", projectID, defectID).
		Order("created_at ASC").
		Find(&defects).Error

	if err != nil {
		return fmt.Errorf("query duplicate defects: %w", err)
	}

	// 如果有多条记录，删除除了第一条之外的所有记录
	if len(defects) > 1 {
		log.Printf("[CleanupDuplicateDefects] Found %d duplicate defect_id=%s in project %d, deleting %d records",
			len(defects), defectID, projectID, len(defects)-1)

		// 保留第一条（最早的），删除其他
		for i := 1; i < len(defects); i++ {
			if err := r.db.Unscoped().Delete(&defects[i]).Error; err != nil {
				log.Printf("[CleanupDuplicateDefects] Failed to delete duplicate defect: %v", err)
			}
		}
	}

	return nil
}
