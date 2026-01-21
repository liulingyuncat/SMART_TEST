package services

import (
	"errors"
	"fmt"
	"strings"
	"webtest/internal/models"
	"webtest/internal/repositories"
)

// UserDefinedVariableService 用户自定义变量服务接口
type UserDefinedVariableService interface {
	// GetVariablesByGroup 获取用例集的变量列表
	GetVariablesByGroup(groupID uint, groupType string) ([]*models.UserDefinedVariable, error)
	// GetVariablesByTask 获取执行任务的变量列表（优先任务变量，没有则返回用例集变量）
	GetVariablesByTask(taskUUID string, groupID uint, groupType string) ([]*models.UserDefinedVariable, error)
	// SaveVariables 批量保存变量（替换模式）
	SaveVariables(groupID uint, groupType string, projectID uint, variables []*models.UserDefinedVariable) error
	// SaveTaskVariables 批量保存任务变量（替换模式）
	SaveTaskVariables(taskUUID string, groupID uint, groupType string, projectID uint, variables []*models.UserDefinedVariable) error
	// AddVariable 添加单个变量
	AddVariable(variable *models.UserDefinedVariable) error
	// UpdateVariable 更新单个变量
	UpdateVariable(variable *models.UserDefinedVariable) error
	// DeleteVariable 删除单个变量
	DeleteVariable(id uint) error
}

type userDefinedVariableService struct {
	repo repositories.UserDefinedVariableRepository
}

// NewUserDefinedVariableService 创建服务实例
func NewUserDefinedVariableService(repo repositories.UserDefinedVariableRepository) UserDefinedVariableService {
	return &userDefinedVariableService{repo: repo}
}

// GetVariablesByGroup 获取用例集的变量列表
func (s *userDefinedVariableService) GetVariablesByGroup(groupID uint, groupType string) ([]*models.UserDefinedVariable, error) {
	fmt.Printf("[VariableService GetVariablesByGroup] groupID=%d, groupType=%s\n", groupID, groupType)
	if groupID == 0 {
		fmt.Printf("[VariableService GetVariablesByGroup] ERROR: group_id is 0\n")
		return nil, errors.New("group_id is required")
	}
	if groupType == "" {
		fmt.Printf("[VariableService GetVariablesByGroup] ERROR: group_type is empty\n")
		return nil, errors.New("group_type is required")
	}
	variables, err := s.repo.GetByGroupID(groupID, groupType)
	if err != nil {
		fmt.Printf("[VariableService GetVariablesByGroup] ERROR: %v\n", err)
		return nil, err
	}
	fmt.Printf("[VariableService GetVariablesByGroup] Found %d variables\n", len(variables))
	for i, v := range variables {
		fmt.Printf("[VariableService GetVariablesByGroup]   [%d] ID=%d, Key=%s, Value=%s\n", i, v.ID, v.VarKey, v.VarValue)
	}
	return variables, nil
}

// GetVariablesByTask 获取执行任务的变量列表
// 只从任务独立的变量表获取（任务创建时已从用例集继承变量）
func (s *userDefinedVariableService) GetVariablesByTask(taskUUID string, groupID uint, groupType string) ([]*models.UserDefinedVariable, error) {
	if taskUUID == "" {
		return nil, errors.New("task_uuid is required")
	}

	// 只获取任务独立的变量（任务创建时已从用例集继承）
	// 不再 fallback 到用例集，因为：
	// 1. 任务创建时，变量已经从用例集复制到任务变量表
	// 2. 测试者可能在执行画面修改了变量值（如IP地址）
	// 3. 执行时应该只使用任务变量表中的值，确保测试环境隔离
	return s.repo.GetByTaskUUID(taskUUID)
}

// SaveVariables 批量保存变量（替换模式）
func (s *userDefinedVariableService) SaveVariables(groupID uint, groupType string, projectID uint, variables []*models.UserDefinedVariable) error {
	if groupID == 0 {
		return errors.New("group_id is required")
	}
	if groupType == "" {
		return errors.New("group_type is required")
	}

	// 验证并规范化变量
	for _, v := range variables {
		if v.VarKey == "" {
			return errors.New("var_key is required for all variables")
		}
		// 确保var_key为小写
		v.VarKey = strings.ToLower(v.VarKey)
		// 自动生成var_name
		v.VarName = "${" + strings.ToUpper(v.VarKey) + "}"
		// 默认类型
		if v.VarType == "" {
			v.VarType = "custom"
		}
	}

	return s.repo.BatchUpsert(groupID, groupType, projectID, variables)
}

// SaveTaskVariables 批量保存任务变量（替换模式）
func (s *userDefinedVariableService) SaveTaskVariables(taskUUID string, groupID uint, groupType string, projectID uint, variables []*models.UserDefinedVariable) error {
	fmt.Printf("[VariableService SaveTaskVariables] taskUUID=%s, groupID=%d, groupType=%s, projectID=%d\n", taskUUID, groupID, groupType, projectID)
	fmt.Printf("[VariableService SaveTaskVariables] Input variables count: %d\n", len(variables))

	if taskUUID == "" {
		fmt.Printf("[VariableService SaveTaskVariables] ERROR: task_uuid is empty\n")
		return errors.New("task_uuid is required")
	}

	// 验证并规范化变量
	for i, v := range variables {
		fmt.Printf("[VariableService SaveTaskVariables] Processing variable[%d]: Key=%s, Value=%s\n", i, v.VarKey, v.VarValue)
		if v.VarKey == "" {
			fmt.Printf("[VariableService SaveTaskVariables] ERROR: var_key is empty at index %d\n", i)
			return errors.New("var_key is required for all variables")
		}
		// 确保var_key为小写
		v.VarKey = strings.ToLower(v.VarKey)
		// 自动生成var_name
		v.VarName = "${" + strings.ToUpper(v.VarKey) + "}"
		// 默认类型
		if v.VarType == "" {
			v.VarType = "custom"
		}
		fmt.Printf("[VariableService SaveTaskVariables] Normalized variable[%d]: Key=%s, Name=%s, Type=%s\n", i, v.VarKey, v.VarName, v.VarType)
	}

	fmt.Printf("[VariableService SaveTaskVariables] Calling repo.BatchUpsertTaskVariables...\n")
	err := s.repo.BatchUpsertTaskVariables(taskUUID, groupID, groupType, projectID, variables)
	if err != nil {
		fmt.Printf("[VariableService SaveTaskVariables] ERROR: %v\n", err)
		return err
	}
	fmt.Printf("[VariableService SaveTaskVariables] ✅ Successfully saved %d variables\n", len(variables))
	return nil
}

// AddVariable 添加单个变量
func (s *userDefinedVariableService) AddVariable(variable *models.UserDefinedVariable) error {
	if variable.GroupID == 0 {
		return errors.New("group_id is required")
	}
	if variable.VarKey == "" {
		return errors.New("var_key is required")
	}

	// 规范化
	variable.VarKey = strings.ToLower(variable.VarKey)
	variable.VarName = "${" + strings.ToUpper(variable.VarKey) + "}"
	if variable.VarType == "" {
		variable.VarType = "custom"
	}

	return s.repo.Create(variable)
}

// UpdateVariable 更新单个变量
func (s *userDefinedVariableService) UpdateVariable(variable *models.UserDefinedVariable) error {
	if variable.ID == 0 {
		return errors.New("variable id is required")
	}

	// 检查变量是否存在
	existing, err := s.repo.GetByID(variable.ID)
	if err != nil {
		return errors.New("variable not found")
	}

	// 规范化
	if variable.VarKey != "" {
		variable.VarKey = strings.ToLower(variable.VarKey)
		variable.VarName = "${" + strings.ToUpper(variable.VarKey) + "}"
	} else {
		variable.VarKey = existing.VarKey
		variable.VarName = existing.VarName
	}

	// 保留原有的关联信息
	variable.GroupID = existing.GroupID
	variable.GroupType = existing.GroupType
	variable.ProjectID = existing.ProjectID

	return s.repo.Update(variable)
}

// DeleteVariable 删除单个变量
func (s *userDefinedVariableService) DeleteVariable(id uint) error {
	if id == 0 {
		return errors.New("variable id is required")
	}
	return s.repo.Delete(id)
}
