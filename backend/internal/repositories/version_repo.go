package repositories

import (
	"webtest/internal/models"

	"gorm.io/gorm"
)

// VersionRepository 版本仓库接口
type VersionRepository interface {
	Create(version *models.Version) error
	Delete(id uint) error
	FindByID(id uint) (*models.Version, error)
	FindByProjectID(projectID uint) ([]*models.Version, error)
	FindByProjectIDAndDocType(projectID uint, docType string) ([]*models.Version, error)
	FindByProjectIDAndItemType(projectID uint, itemType string) ([]*models.Version, error)
	UpdateRemark(id uint, remark string) error
}

// versionRepository 版本仓库实现
type versionRepository struct {
	db *gorm.DB
}

// NewVersionRepository 创建版本仓库实例
func NewVersionRepository(db *gorm.DB) VersionRepository {
	return &versionRepository{db: db}
}

// Create 创建版本记录
func (r *versionRepository) Create(version *models.Version) error {
	return r.db.Create(version).Error
}

// Delete 删除版本记录(硬删除)
func (r *versionRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&models.Version{}, id).Error
}

// FindByID 根据ID查询版本记录
func (r *versionRepository) FindByID(id uint) (*models.Version, error) {
	var version models.Version
	if err := r.db.First(&version, id).Error; err != nil {
		return nil, err
	}
	return &version, nil
}

// FindByProjectID 根据项目ID查询所有版本记录
func (r *versionRepository) FindByProjectID(projectID uint) ([]*models.Version, error) {
	var versions []*models.Version
	if err := r.db.Where("project_id = ?", projectID).Order("created_at DESC").Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

// FindByProjectIDAndDocType 根据项目ID和文档类型查询版本记录(兼容旧接口)
func (r *versionRepository) FindByProjectIDAndDocType(projectID uint, docType string) ([]*models.Version, error) {
	var versions []*models.Version
	if err := r.db.Where("project_id = ? AND doc_type = ?", projectID, docType).
		Order("created_at DESC").Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

// FindByProjectIDAndItemType 根据项目ID和条目类型查询版本记录(新接口)
func (r *versionRepository) FindByProjectIDAndItemType(projectID uint, itemType string) ([]*models.Version, error) {
	var versions []*models.Version
	if err := r.db.Where("project_id = ? AND item_type = ?", projectID, itemType).
		Order("created_at DESC").Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

// UpdateRemark 更新版本备注
func (r *versionRepository) UpdateRemark(id uint, remark string) error {
	return r.db.Model(&models.Version{}).Where("id = ?", id).Update("remark", remark).Error
}
