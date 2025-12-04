package repositories

import (
	"fmt"
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
	GetByProjectID(projectID uint) ([]*models.Defect, error)
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

// Delete 软删除缺陷
func (r *defectRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&models.Defect{})
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

// GetMaxDefectSeq 获取项目最大缺陷序号
func (r *defectRepository) GetMaxDefectSeq(projectID uint) (int, error) {
	var maxSeq int
	err := r.db.Model(&models.Defect{}).
		Select("COALESCE(MAX(CAST(SUBSTR(defect_id, 5) AS INTEGER)), 0)").
		Where("project_id = ?", projectID).
		Scan(&maxSeq).Error

	if err != nil {
		return 0, fmt.Errorf("get max defect seq: %w", err)
	}

	return maxSeq, nil
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
