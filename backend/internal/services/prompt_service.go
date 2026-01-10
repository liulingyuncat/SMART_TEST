package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"webtest/internal/models"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// 错误定义
var (
	ErrPromptNameExists       = errors.New("提示词名称已存在")
	ErrPromptNotFound         = errors.New("提示词不存在")
	ErrPromptPermissionDenied = errors.New("无权限操作此提示词")
	ErrCannotModifySystem     = errors.New("不能修改系统提示词")
)

// PromptService 提示词服务接口
type PromptService interface {
	ListPromptsWithRole(projectID uint, scope string, userID uint, userRole string, page, pageSize int) ([]models.PromptDTO, int64, error)
	ListPrompts(projectID uint, scope string, userID uint, page, pageSize int) ([]models.PromptDTO, int64, error)
	GetPromptByID(id uint, userID uint) (*models.PromptDTO, error)
	GetPromptByName(name string, userID uint) (*models.PromptDTO, error)
	CreatePrompt(projectID uint, name, description, version, content, scope string, arguments []models.PromptArgument, userID uint, userRole string) (*models.PromptDTO, error)
	UpdatePrompt(id uint, description, version, content *string, arguments []models.PromptArgument, userID uint, userRole string) (*models.PromptDTO, error)
	DeletePrompt(id uint, userID uint, userRole string) error
	RefreshSystemPromptsFromDirectory(promptsDir string) error
}

// promptService 提示词服务实现
type promptService struct {
	db *gorm.DB
}

// NewPromptService 创建提示词服务实例
func NewPromptService(db *gorm.DB) PromptService {
	return &promptService{db: db}
}

// ListPromptsWithRole 根据用户角色查询提示词列表
func (s *promptService) ListPromptsWithRole(
	projectID uint,
	scope string,
	userID uint,
	userRole string,
	page, pageSize int,
) ([]models.PromptDTO, int64, error) {
	var prompts []models.Prompt
	var total int64

	query := s.db.Model(&models.Prompt{})
	log.Printf("[PromptService] ListPromptsWithRole: scope=%s, userID=%d, userRole=%s", scope, userID, userRole)
	// 根据scope类型进行查询（提示词与project_id无关）
	if scope == "system" {
		// 系统提示词：全员可见
		query = query.Where("scope = ?", "system")
	} else if scope == "project" {
		// 全员提示词：所有用户可见
		query = query.Where("scope = ?", "project")
	} else if scope == "user" {
		// 个人提示词：仅自己可见
		query = query.Where("scope = ? AND user_id = ?", "user", userID)
	} else {
		// 如果没有指定scope，返回所有可见的提示词（系统+全员+个人）
		query = query.Where("(scope = ?) OR (scope = ?) OR (scope = ? AND user_id = ?)",
			"system", "project", "user", userID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count prompts: %w", err)
	}

	// 分页查询（不加载content字段）
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&prompts).Error; err != nil {
		return nil, 0, fmt.Errorf("query prompts: %w", err)
	}

	// 转换为DTO
	dtos := make([]models.PromptDTO, len(prompts))
	for i, prompt := range prompts {
		dto, err := prompt.ToDTO(false)
		if err != nil {
			return nil, 0, fmt.Errorf("convert to DTO: %w", err)
		}
		dtos[i] = *dto
	}

	return dtos, total, nil
}

// ListPrompts 查询提示词列表
func (s *promptService) ListPrompts(
	projectID uint,
	scope string,
	userID uint,
	page, pageSize int,
) ([]models.PromptDTO, int64, error) {
	var prompts []models.Prompt
	var total int64

	query := s.db.Model(&models.Prompt{})

	log.Printf("[PromptService] ListPrompts called: scope=%s, userID=%d, projectID=%d", scope, userID, projectID)

	// 根据scope类型进行查询（提示词与project_id无关）
	if scope == "system" {
		// 系统提示词：全员可见
		query = query.Where("scope = ?", "system")
	} else if scope == "project" {
		// 全员提示词：所有用户可见
		query = query.Where("scope = ?", "project")
	} else if scope == "user" {
		// 个人提示词：仅自己可见
		query = query.Where("scope = ? AND user_id = ?", "user", userID)
	} else {
		// 如果没有指定scope，返回所有可见的提示词（系统+全员+个人）
		query = query.Where("(scope = ?) OR (scope = ?) OR (scope = ? AND user_id = ?)",
			"system", "project", "user", userID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count prompts: %w", err)
	}

	// 分页查询（不加载content字段）
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&prompts).Error; err != nil {
		return nil, 0, fmt.Errorf("query prompts: %w", err)
	}

	// 转换为DTO
	dtos := make([]models.PromptDTO, len(prompts))
	for i, prompt := range prompts {
		dto, err := prompt.ToDTO(false)
		if err != nil {
			return nil, 0, fmt.Errorf("convert to DTO: %w", err)
		}
		dtos[i] = *dto
	}

	return dtos, total, nil
}

// GetPromptByID 根据ID获取提示词详情
func (s *promptService) GetPromptByID(id uint, userID uint) (*models.PromptDTO, error) {
	var prompt models.Prompt
	if err := s.db.First(&prompt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPromptNotFound
		}
		return nil, fmt.Errorf("query prompt: %w", err)
	}

	// 权限验证：个人Prompt仅所有者可访问
	if prompt.Scope == "user" && prompt.UserID != nil && *prompt.UserID != userID {
		return nil, ErrPromptPermissionDenied
	}

	dto, err := prompt.ToDTO(true)
	if err != nil {
		return nil, fmt.Errorf("convert to DTO: %w", err)
	}

	return dto, nil
}

// GetPromptByName 根据名称获取提示词详情（用于MCP）
func (s *promptService) GetPromptByName(name string, userID uint) (*models.PromptDTO, error) {
	var prompt models.Prompt

	log.Printf("[PromptService] GetPromptByName: name=%s, userID=%d", name, userID)

	// 查询逻辑：按优先级查找 user -> project -> system
	// 1. 先找个人提示词（如果提供了userID）
	if userID != 0 {
		if err := s.db.Where("name = ? AND scope = ? AND user_id = ?", name, "user", userID).First(&prompt).Error; err == nil {
			log.Printf("[PromptService] Found user prompt: id=%d, name=%s", prompt.ID, prompt.Name)
			dto, err := prompt.ToDTO(true)
			if err != nil {
				return nil, fmt.Errorf("convert to DTO: %w", err)
			}
			return dto, nil
		}
	}

	// 2. 查找全员提示词
	if err := s.db.Where("name = ? AND scope = ?", name, "project").First(&prompt).Error; err == nil {
		log.Printf("[PromptService] Found project prompt: id=%d, name=%s", prompt.ID, prompt.Name)
		dto, err := prompt.ToDTO(true)
		if err != nil {
			return nil, fmt.Errorf("convert to DTO: %w", err)
		}
		return dto, nil
	}

	// 3. 查找系统提示词
	if err := s.db.Where("name = ? AND scope = ?", name, "system").First(&prompt).Error; err == nil {
		log.Printf("[PromptService] Found system prompt: id=%d, name=%s", prompt.ID, prompt.Name)
		dto, err := prompt.ToDTO(true)
		if err != nil {
			return nil, fmt.Errorf("convert to DTO: %w", err)
		}
		return dto, nil
	}

	// 找不到
	log.Printf("[PromptService] Prompt not found: name=%s", name)
	return nil, ErrPromptNotFound
}

// CreatePrompt 创建新提示词
func (s *promptService) CreatePrompt(
	projectID uint,
	name, description, version, content, scope string,
	arguments []models.PromptArgument,
	userID uint,
	userRole string,
) (*models.PromptDTO, error) {
	// 权限检查：全员提示词只有系统管理员可以创建
	if scope == "project" && userRole != "system_admin" {
		return nil, ErrPromptPermissionDenied
	}

	// 检查名称唯一性
	// 系统和全员提示词全局唯一，个人提示词仅对用户唯一
	var count int64
	query := s.db.Model(&models.Prompt{}).
		Where("name = ? AND scope = ?", name, scope)

	if scope == "user" {
		// 个人提示词仅对用户唯一
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("check name uniqueness: %w", err)
	}

	if count > 0 {
		return nil, ErrPromptNameExists
	}

	// 构建Prompt对象
	prompt := &models.Prompt{}
	var userIDPtr *uint
	if scope == "user" {
		userIDPtr = &userID
	}

	if err := prompt.FromCreateRequest(
		projectID, name, description, version, content, scope,
		arguments, userIDPtr, userID,
	); err != nil {
		return nil, fmt.Errorf("build prompt: %w", err)
	}

	// 保存到数据库
	if err := s.db.Create(prompt).Error; err != nil {
		return nil, fmt.Errorf("create prompt: %w", err)
	}

	dto, err := prompt.ToDTO(true)
	if err != nil {
		return nil, fmt.Errorf("convert to DTO: %w", err)
	}

	return dto, nil
}

// UpdatePrompt 更新提示词
func (s *promptService) UpdatePrompt(
	id uint,
	description, version, content *string,
	arguments []models.PromptArgument,
	userID uint,
	userRole string,
) (*models.PromptDTO, error) {
	var prompt models.Prompt
	if err := s.db.First(&prompt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPromptNotFound
		}
		return nil, fmt.Errorf("query prompt: %w", err)
	}

	// 权限验证
	// 系统提示词全员只读
	if prompt.Scope == "system" {
		return nil, ErrCannotModifySystem
	}
	// 个人提示词仅自己可编辑
	if prompt.Scope == "user" && prompt.UserID != nil && *prompt.UserID != userID {
		return nil, ErrPromptPermissionDenied
	}
	// 全员提示词仅系统管理员可以编辑
	if prompt.Scope == "project" && userRole != "system_admin" {
		return nil, ErrPromptPermissionDenied
	}

	// 更新字段
	updates := make(map[string]interface{})
	if description != nil {
		updates["description"] = *description
	}
	if version != nil {
		updates["version"] = *version
	}
	if content != nil {
		updates["content"] = *content
	}
	if arguments != nil {
		argsJSON, err := json.Marshal(arguments)
		if err != nil {
			return nil, fmt.Errorf("marshal arguments: %w", err)
		}
		updates["arguments"] = string(argsJSON)
	}

	if err := s.db.Model(&prompt).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update prompt: %w", err)
	}

	// 重新查询
	if err := s.db.First(&prompt, id).Error; err != nil {
		return nil, fmt.Errorf("re-query prompt: %w", err)
	}

	dto, err := prompt.ToDTO(true)
	if err != nil {
		return nil, fmt.Errorf("convert to DTO: %w", err)
	}

	return dto, nil
}

// DeletePrompt 删除提示词(硬删除)
func (s *promptService) DeletePrompt(id uint, userID uint, userRole string) error {
	var prompt models.Prompt
	if err := s.db.First(&prompt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPromptNotFound
		}
		return fmt.Errorf("query prompt: %w", err)
	}

	// 权限验证
	// 系统提示词全员只读
	if prompt.Scope == "system" {
		return ErrCannotModifySystem
	}
	// 个人提示词仅自己可删除
	if prompt.Scope == "user" && prompt.UserID != nil && *prompt.UserID != userID {
		return ErrPromptPermissionDenied
	}
	// 全员提示词仅系统管理员可以删除
	if prompt.Scope == "project" && userRole != "system_admin" {
		return ErrPromptPermissionDenied
	}

	if err := s.db.Unscoped().Delete(&prompt).Error; err != nil {
		return fmt.Errorf("delete prompt: %w", err)
	}

	return nil
}

// PromptFrontMatter 提示词文件的YAML头部结构
type PromptFrontMatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// parsePromptFile 解析提示词文件，提取YAML front matter和内容
func parsePromptFile(filePath string) (*PromptFrontMatter, string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	contentStr := string(content)

	// 使用正则匹配 YAML front matter (---开头和结尾)
	frontMatterRegex := regexp.MustCompile(`(?s)^---\s*\n(.+?)\n---\s*\n?(.*)$`)
	matches := frontMatterRegex.FindStringSubmatch(contentStr)

	if len(matches) != 3 {
		// 没有front matter，使用文件名作为name
		baseName := filepath.Base(filePath)
		name := strings.TrimSuffix(baseName, ".prompt.md")
		return &PromptFrontMatter{
			Name:        name,
			Description: "",
			Version:     "1.0",
		}, contentStr, nil
	}

	// 解析YAML
	var frontMatter PromptFrontMatter
	if err := yaml.Unmarshal([]byte(matches[1]), &frontMatter); err != nil {
		return nil, "", fmt.Errorf("parse yaml front matter: %w", err)
	}

	// 如果name为空，使用文件名
	if frontMatter.Name == "" {
		baseName := filepath.Base(filePath)
		frontMatter.Name = strings.TrimSuffix(baseName, ".prompt.md")
	}

	// 如果version为空，默认1.0
	if frontMatter.Version == "" {
		frontMatter.Version = "1.0"
	}

	return &frontMatter, contentStr, nil
}

// RefreshSystemPromptsFromDirectory 从目录动态扫描并刷新系统提示词
func (s *promptService) RefreshSystemPromptsFromDirectory(promptsDir string) error {
	projectID := uint(1)

	// 扫描目录中所有 .prompt.md 文件
	files, err := filepath.Glob(filepath.Join(promptsDir, "*.prompt.md"))
	if err != nil {
		return fmt.Errorf("scan prompts directory: %w", err)
	}

	fmt.Printf("[Info] Found %d prompt files in %s\n", len(files), promptsDir)

	// 记录当前文件中的所有提示词名称（用于删除检测）
	filePromptNames := make(map[string]bool)

	for _, filePath := range files {
		// 检查文件是否为空
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Printf("[Warn] Failed to stat file %s: %v\n", filePath, err)
			continue
		}
		if fileInfo.Size() == 0 {
			fmt.Printf("[Skip] Empty file: %s\n", filePath)
			continue
		}

		// 解析文件
		frontMatter, content, err := parsePromptFile(filePath)
		if err != nil {
			fmt.Printf("[Warn] Failed to parse prompt file %s: %v\n", filePath, err)
			continue
		}

		filePromptNames[frontMatter.Name] = true

		// 更新或创建数据库记录
		prompt := &models.Prompt{}
		result := s.db.Where("name = ? AND scope = ?", frontMatter.Name, "system").First(prompt)

		if result.Error == nil {
			// 记录存在，更新
			if err := s.db.Model(prompt).Updates(map[string]interface{}{
				"description": frontMatter.Description,
				"version":     frontMatter.Version,
				"content":     content,
			}).Error; err != nil {
				fmt.Printf("[Warn] Failed to update prompt %s: %v\n", frontMatter.Name, err)
				continue
			}
			fmt.Printf("[Updated] System prompt: %s\n", frontMatter.Name)
		} else {
			// 记录不存在，创建
			newPrompt := &models.Prompt{
				ProjectID:   projectID,
				Name:        frontMatter.Name,
				Description: frontMatter.Description,
				Version:     frontMatter.Version,
				Content:     content,
				Scope:       "system",
				UserID:      nil,
			}
			if err := s.db.Create(newPrompt).Error; err != nil {
				fmt.Printf("[Warn] Failed to create prompt %s: %v\n", frontMatter.Name, err)
				continue
			}
			fmt.Printf("[Created] System prompt: %s\n", frontMatter.Name)
		}
	}

	// 删除数据库中存在但文件已不存在的系统提示词
	var existingPrompts []models.Prompt
	if err := s.db.Where("scope = ?", "system").Find(&existingPrompts).Error; err != nil {
		fmt.Printf("[Warn] Failed to query existing system prompts: %v\n", err)
	} else {
		for _, prompt := range existingPrompts {
			if !filePromptNames[prompt.Name] {
				if err := s.db.Unscoped().Delete(&prompt).Error; err != nil {
					fmt.Printf("[Warn] Failed to delete orphan prompt %s: %v\n", prompt.Name, err)
				} else {
					fmt.Printf("[Deleted] Orphan system prompt: %s\n", prompt.Name)
				}
			}
		}
	}

	fmt.Printf("[Info] System prompts refresh completed. Active: %d\n", len(filePromptNames))
	return nil
}
