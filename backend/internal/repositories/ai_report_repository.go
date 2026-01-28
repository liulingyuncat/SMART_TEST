package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// AIReportRepository AI报告仓库接口
type AIReportRepository interface {
	Create(report *models.AIReport) error
	Update(report *models.AIReport) error
	Delete(id string) error
	FindByID(id string) (*models.AIReport, error)
	FindByProjectID(projectID uint) ([]*models.AIReport, error)
	FindByProjectIDAndType(projectID uint, reportType string) ([]*models.AIReport, error)
	ExistsByProjectAndName(projectID uint, name string, excludeID string) (bool, error)
	ExistsByProjectTypeAndName(projectID uint, reportType, name string, excludeID string) (bool, error)
}

// aiReportRepository AI报告仓库实现
type aiReportRepository struct {
	db *gorm.DB
}

// NewAIReportRepository 创建AI报告仓库实例
func NewAIReportRepository(db *gorm.DB) AIReportRepository {
	return &aiReportRepository{db: db}
}

// Create 创建AI报告
func (r *aiReportRepository) Create(report *models.AIReport) error {
	return r.db.Create(report).Error
}

// Update 更新AI报告
func (r *aiReportRepository) Update(report *models.AIReport) error {
	return r.db.Save(report).Error
}

// Delete 删除AI报告(硬删除)
func (r *aiReportRepository) Delete(id string) error {
	return r.db.Unscoped().Delete(&models.AIReport{}, "id = ?", id).Error
}

// FindByID 根据ID查询AI报告
func (r *aiReportRepository) FindByID(id string) (*models.AIReport, error) {
	var report models.AIReport
	if err := r.db.First(&report, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

// FindByProjectID 根据项目ID查询所有AI报告(按创建时间降序)
func (r *aiReportRepository) FindByProjectID(projectID uint) ([]*models.AIReport, error) {
	var reports []*models.AIReport
	if err := r.db.
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

// ExistsByProjectAndName 检查项目下是否存在同名报告(名称检重,排除指定ID)
func (r *aiReportRepository) ExistsByProjectAndName(projectID uint, name string, excludeID string) (bool, error) {
	var count int64
	query := r.db.Where("project_id = ? AND name = ?", projectID, name)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Model(&models.AIReport{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByProjectIDAndType 根据项目ID和类型查询AI报告(按创建时间降序)
func (r *aiReportRepository) FindByProjectIDAndType(projectID uint, reportType string) ([]*models.AIReport, error) {
	var reports []*models.AIReport
	if err := r.db.
		Where("project_id = ? AND type = ?", projectID, reportType).
		Order("created_at DESC").
		Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

// ExistsByProjectTypeAndName 检查项目下指定类型是否存在同名报告(名称检重,排除指定ID)
func (r *aiReportRepository) ExistsByProjectTypeAndName(projectID uint, reportType, name string, excludeID string) (bool, error) {
	var count int64
	query := r.db.Where("project_id = ? AND type = ? AND name = ?", projectID, reportType, name)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Model(&models.AIReport{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
