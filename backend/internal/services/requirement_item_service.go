package services

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"
)

// RequirementItemService 需求条目服务接口
type RequirementItemService interface {
	// 基础CRUD
	CreateItem(projectID uint, name, content string) (*models.RequirementItem, error)
	UpdateItem(id uint, name, content string) (*models.RequirementItem, error)
	DeleteItem(id uint) error
	GetItemByID(id uint) (*models.RequirementItem, error)
	GetItemsByProjectID(projectID uint) ([]*models.RequirementItem, error)

	// 批量操作
	BulkCreateItems(projectID uint, items []struct{ Name, Content string }) error
	BulkUpdateItems(items []struct {
		ID      uint
		Name    string
		Content string
	}) error
	BulkDeleteItems(ids []uint) error

	// ZIP批量版本打包
	ExportToZip(projectID uint, outputPath, remark string, createdBy uint) (*models.Version, error)
	ImportFromZip(projectID uint, zipPath string, createdBy uint) error
}

// requirementItemService 需求条目服务实现
type requirementItemService struct {
	itemRepo    repositories.RequirementItemRepository
	versionRepo repositories.VersionRepository
	storageDir  string // 存储目录
}

// NewRequirementItemService 创建需求条目服务实例
func NewRequirementItemService(
	itemRepo repositories.RequirementItemRepository,
	versionRepo repositories.VersionRepository,
	storageDir string,
) RequirementItemService {
	return &requirementItemService{
		itemRepo:    itemRepo,
		versionRepo: versionRepo,
		storageDir:  storageDir,
	}
}

// CreateItem 创建需求条目
func (s *requirementItemService) CreateItem(projectID uint, name, content string) (*models.RequirementItem, error) {
	// 检查重名
	existing, _ := s.itemRepo.FindByProjectIDAndName(projectID, name)
	if existing != nil {
		return nil, fmt.Errorf("需求名称 '%s' 已存在", name)
	}

	item := &models.RequirementItem{
		ProjectID: projectID,
		Name:      name,
		Content:   content,
	}

	if err := s.itemRepo.Create(item); err != nil {
		return nil, fmt.Errorf("创建需求条目失败: %w", err)
	}

	return item, nil
}

// UpdateItem 更新需求条目
func (s *requirementItemService) UpdateItem(id uint, name, content string) (*models.RequirementItem, error) {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("需求条目不存在: %w", err)
	}

	// 检查重名(排除自己)
	if name != item.Name {
		existing, _ := s.itemRepo.FindByProjectIDAndName(item.ProjectID, name)
		if existing != nil && existing.ID != id {
			return nil, fmt.Errorf("需求名称 '%s' 已存在", name)
		}
	}

	item.Name = name
	item.Content = content

	if err := s.itemRepo.Update(item); err != nil {
		return nil, fmt.Errorf("更新需求条目失败: %w", err)
	}

	return item, nil
}

// DeleteItem 删除需求条目
func (s *requirementItemService) DeleteItem(id uint) error {
	if err := s.itemRepo.Delete(id); err != nil {
		return fmt.Errorf("删除需求条目失败: %w", err)
	}
	return nil
}

// GetItemByID 获取单个需求条目
func (s *requirementItemService) GetItemByID(id uint) (*models.RequirementItem, error) {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("需求条目不存在: %w", err)
	}
	return item, nil
}

// GetItemsByProjectID 获取项目的所有需求条目
func (s *requirementItemService) GetItemsByProjectID(projectID uint) ([]*models.RequirementItem, error) {
	items, err := s.itemRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("查询需求条目失败: %w", err)
	}
	return items, nil
}

// BulkCreateItems 批量创建需求条目
func (s *requirementItemService) BulkCreateItems(projectID uint, items []struct{ Name, Content string }) error {
	var newItems []*models.RequirementItem

	for _, item := range items {
		// 检查重名
		existing, _ := s.itemRepo.FindByProjectIDAndName(projectID, item.Name)
		if existing != nil {
			return fmt.Errorf("需求名称 '%s' 已存在", item.Name)
		}

		newItems = append(newItems, &models.RequirementItem{
			ProjectID: projectID,
			Name:      item.Name,
			Content:   item.Content,
		})
	}

	if err := s.itemRepo.BulkCreate(newItems); err != nil {
		return fmt.Errorf("批量创建需求条目失败: %w", err)
	}

	return nil
}

// BulkUpdateItems 批量更新需求条目
func (s *requirementItemService) BulkUpdateItems(items []struct {
	ID      uint
	Name    string
	Content string
}) error {
	var updateItems []*models.RequirementItem

	for _, item := range items {
		existing, err := s.itemRepo.FindByID(item.ID)
		if err != nil {
			return fmt.Errorf("需求条目 ID=%d 不存在", item.ID)
		}

		existing.Name = item.Name
		existing.Content = item.Content
		updateItems = append(updateItems, existing)
	}

	if err := s.itemRepo.BulkUpdate(updateItems); err != nil {
		return fmt.Errorf("批量更新需求条目失败: %w", err)
	}

	return nil
}

// BulkDeleteItems 批量删除需求条目
func (s *requirementItemService) BulkDeleteItems(ids []uint) error {
	if err := s.itemRepo.BulkDelete(ids); err != nil {
		return fmt.Errorf("批量删除需求条目失败: %w", err)
	}
	return nil
}

// ExportToZip 导出为ZIP批量版本
func (s *requirementItemService) ExportToZip(projectID uint, outputPath, remark string, createdBy uint) (*models.Version, error) {
	fmt.Printf("[ExportToZip] 开始导出, projectID=%d, outputPath=%s, remark=%s, createdBy=%d\n", projectID, outputPath, remark, createdBy)

	// 查询所有需求条目
	items, err := s.itemRepo.FindByProjectID(projectID)
	if err != nil {
		fmt.Printf("[ExportToZip] 查询需求条目失败: %v\n", err)
		return nil, fmt.Errorf("查询需求条目失败: %w", err)
	}

	fmt.Printf("[ExportToZip] 查询到%d个需求条目\n", len(items))

	if len(items) == 0 {
		return nil, fmt.Errorf("项目无需求条目,无法导出")
	}

	// 创建ZIP文件
	fmt.Printf("[ExportToZip] 准备创建ZIP文件: %s\n", outputPath)
	zipFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("[ExportToZip] 创建ZIP文件失败: %v\n", err)
		return nil, fmt.Errorf("创建ZIP文件失败: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fmt.Printf("[ExportToZip] ZIP文件创建成功\n")

	var fileList []string

	// 将每个需求条目写入ZIP
	for _, item := range items {
		filename := fmt.Sprintf("%s.md", item.Name)
		fileList = append(fileList, filename)

		writer, err := zipWriter.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("创建ZIP条目失败: %w", err)
		}

		if _, err := writer.Write([]byte(item.Content)); err != nil {
			return nil, fmt.Errorf("写入ZIP条目失败: %w", err)
		}
	}

	// 获取ZIP文件信息
	fileInfo, _ := zipFile.Stat()
	fileListJSON, _ := json.Marshal(fileList)

	// 创建版本记录
	version := &models.Version{
		ProjectID: projectID,
		ItemType:  "requirement-batch",
		Filename:  filepath.Base(outputPath),
		FilePath:  outputPath,
		FileSize:  fileInfo.Size(),
		FileList:  string(fileListJSON),
		Remark:    remark,
		CreatedBy: &createdBy,
		CreatedAt: time.Now(),
	}

	fmt.Printf("[ExportToZip] 准备创建版本记录: ProjectID=%d, ItemType=%s, Filename=%s\n", version.ProjectID, version.ItemType, version.Filename)

	if err := s.versionRepo.Create(version); err != nil {
		fmt.Printf("[ExportToZip] 创建版本记录失败: %v\n", err)
		return nil, fmt.Errorf("创建版本记录失败: %w", err)
	}

	fmt.Printf("[ExportToZip] 版本记录创建成功, ID=%d\n", version.ID)
	return version, nil
}

// ImportFromZip 从ZIP批量版本导入
func (s *requirementItemService) ImportFromZip(projectID uint, zipPath string, createdBy uint) error {
	// 打开ZIP文件
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer reader.Close()

	var newItems []*models.RequirementItem

	// 解析ZIP文件
	for _, file := range reader.File {
		if filepath.Ext(file.Name) != ".md" {
			continue
		}

		// 读取文件内容
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("读取ZIP条目失败: %w", err)
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("读取文件内容失败: %w", err)
		}

		// 提取文件名作为需求名称
		name := file.Name[:len(file.Name)-3]

		newItems = append(newItems, &models.RequirementItem{
			ProjectID: projectID,
			Name:      name,
			Content:   string(content),
		})
	}

	if len(newItems) == 0 {
		return fmt.Errorf("ZIP文件中无有效的Markdown文件")
	}

	// 批量导入
	if err := s.itemRepo.BulkCreate(newItems); err != nil {
		return fmt.Errorf("批量导入需求条目失败: %w", err)
	}

	return nil
}
