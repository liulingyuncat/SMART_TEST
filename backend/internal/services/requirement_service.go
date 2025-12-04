package services

import (
	"fmt"
	"time"
	"webtest/internal/repositories"
)

// RequirementService 需求文档业务逻辑接口
type RequirementService interface {
	GetRequirement(projectID uint, docType string) (string, time.Time, error)
	UpdateRequirement(projectID uint, docType string, content string) error
	UpdateRequirementField(projectID uint, docType string, content string) error
}

// requirementService 需求文档业务逻辑实现
type requirementService struct {
	requirementRepo repositories.RequirementRepository
}

// NewRequirementService 创建需求文档服务实例
func NewRequirementService(requirementRepo repositories.RequirementRepository) RequirementService {
	return &requirementService{
		requirementRepo: requirementRepo,
	}
}

// GetRequirement 获取指定类型的需求文档内容
func (s *requirementService) GetRequirement(projectID uint, docType string) (string, time.Time, error) {
	// 验证 docType 的合法性
	if err := validateDocType(docType); err != nil {
		return "", time.Time{}, err
	}

	// 查询需求文档
	requirement, err := s.requirementRepo.GetByProjectID(projectID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("查询需求文档失败: %w", err)
	}

	// 获取对应字段内容
	content, err := s.requirementRepo.GetFieldByType(projectID, docType)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("获取文档字段失败: %w", err)
	}

	return content, requirement.UpdatedAt, nil
}

// UpdateRequirement 更新指定类型的需求文档内容
func (s *requirementService) UpdateRequirement(projectID uint, docType string, content string) error {
	// 验证 docType 的合法性
	if err := validateDocType(docType); err != nil {
		return err
	}

	// 更新需求文档
	if err := s.requirementRepo.UpdateFieldByType(projectID, docType, content); err != nil {
		return fmt.Errorf("更新需求文档失败: %w", err)
	}

	return nil
}

// UpdateRequirementField 更新需求文档指定字段(用于版本保存)
func (s *requirementService) UpdateRequirementField(projectID uint, docType string, content string) error {
	// 验证 docType 的合法性
	if err := validateDocType(docType); err != nil {
		return err
	}

	// 更新需求文档字段
	if err := s.requirementRepo.UpdateFieldByType(projectID, docType, content); err != nil {
		return fmt.Errorf("更新需求文档字段失败: %w", err)
	}

	return nil
}

// validateDocType 验证文档类型是否合法
func validateDocType(docType string) error {
	validTypes := map[string]bool{
		"overall-requirements":   true,
		"overall-test-viewpoint": true,
		"change-requirements":    true,
		"change-test-viewpoint":  true,
	}

	if !validTypes[docType] {
		return fmt.Errorf("无效的文档类型: %s", docType)
	}
	return nil
}
