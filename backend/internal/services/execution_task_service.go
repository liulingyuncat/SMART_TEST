package services

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"gorm.io/gorm"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	TaskName      string `json:"task_name" binding:"required,min=1,max=50"`
	ExecutionType string `json:"execution_type" binding:"required,oneof=manual automation api"`
	TaskStatus    string `json:"task_status" binding:"omitempty,oneof=pending in_progress completed"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	TaskName        *string    `json:"task_name" binding:"omitempty,min=1,max=50"`
	ExecutionType   *string    `json:"execution_type" binding:"omitempty,oneof=manual automation api"`
	TaskStatus      *string    `json:"task_status" binding:"omitempty,oneof=pending in_progress completed"`
	CaseGroupID     *uint      `json:"case_group_id"`                               // 关联的用例集ID
	CaseGroupName   *string    `json:"case_group_name" binding:"omitempty,max=100"` // 关联的用例集名称
	DisplayLanguage *string    `json:"display_language" binding:"omitempty,max=10"` // 显示语言(cn/jp/en/all)
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	TestVersion     *string    `json:"test_version" binding:"omitempty,max=50"`
	TestEnv         *string    `json:"test_env" binding:"omitempty,max=100"`
	TestDate        *time.Time `json:"test_date"`
	Executor        *string    `json:"executor" binding:"omitempty,max=50"`
	TaskDescription *string    `json:"task_description" binding:"omitempty,max=2000"`
}

// ExecutionTaskService 测试执行任务服务接口
type ExecutionTaskService interface {
	GetTasksByProject(projectID uint, userID uint) ([]*models.ExecutionTask, error)
	CreateTask(projectID uint, userID uint, req CreateTaskRequest) (*models.ExecutionTask, error)
	UpdateTask(projectID uint, userID uint, taskUUID string, req UpdateTaskRequest) (*models.ExecutionTask, error)
	DeleteTask(projectID uint, userID uint, taskUUID string) error
}

type executionTaskService struct {
	repo        repositories.ExecutionTaskRepository
	projectRepo repositories.ProjectRepository             // 用于权限验证
	ecrRepo     repositories.ExecutionCaseResultRepository // 用于级联删除
}

// NewExecutionTaskService 创建任务服务实例
func NewExecutionTaskService(
	repo repositories.ExecutionTaskRepository,
	projectRepo repositories.ProjectRepository,
	ecrRepo repositories.ExecutionCaseResultRepository,
) ExecutionTaskService {
	return &executionTaskService{
		repo:        repo,
		projectRepo: projectRepo,
		ecrRepo:     ecrRepo,
	}
}

// GetTasksByProject 获取项目的所有任务
func (s *executionTaskService) GetTasksByProject(projectID uint, userID uint) ([]*models.ExecutionTask, error) {
	// 验证用户项目权限(中间件已验证,此处可选)
	// 可选: 调用projectRepo验证项目是否存在

	tasks, err := s.repo.GetByProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("get tasks by project %d: %w", projectID, err)
	}
	return tasks, nil
}

// CreateTask 创建新任务
func (s *executionTaskService) CreateTask(projectID uint, userID uint, req CreateTaskRequest) (*models.ExecutionTask, error) {
	// 1. 验证任务名唯一性(不区分大小写)
	isUnique, err := s.repo.CheckNameUnique(projectID, req.TaskName, "")
	if err != nil {
		return nil, fmt.Errorf("check name unique: %w", err)
	}
	if !isUnique {
		return nil, errors.New("任务名已存在")
	}

	// 2. 构建任务对象
	task := &models.ExecutionTask{
		ProjectID:     projectID,
		TaskName:      req.TaskName,
		ExecutionType: req.ExecutionType,
		TaskStatus:    "pending", // 默认状态
		CreatedBy:     userID,
	}

	// 如果请求中指定了状态,使用指定值
	if req.TaskStatus != "" {
		task.TaskStatus = req.TaskStatus
	}

	// 3. 创建任务(BeforeCreate Hook会生成UUID)
	err = s.repo.Create(task)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

// UpdateTask 更新任务
func (s *executionTaskService) UpdateTask(projectID uint, userID uint, taskUUID string, req UpdateTaskRequest) (*models.ExecutionTask, error) {
	// 1. 验证任务存在且属于该项目
	task, err := s.repo.GetByUUID(taskUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, fmt.Errorf("get task by uuid: %w", err)
	}

	if task.ProjectID != projectID {
		return nil, errors.New("任务不属于该项目")
	}

	// 2. 验证任务名唯一性(如果修改了名称)
	if req.TaskName != nil && *req.TaskName != task.TaskName {
		isUnique, err := s.repo.CheckNameUnique(projectID, *req.TaskName, taskUUID)
		if err != nil {
			return nil, fmt.Errorf("check name unique: %w", err)
		}
		if !isUnique {
			return nil, errors.New("任务名已存在")
		}
	}

	// 3. 验证日期范围逻辑
	startDate := task.StartDate
	endDate := task.EndDate

	if req.StartDate != nil {
		startDate = req.StartDate
	}
	if req.EndDate != nil {
		endDate = req.EndDate
	}

	if startDate != nil && endDate != nil {
		if endDate.Before(*startDate) {
			return nil, errors.New("结束日期不能早于开始日期")
		}
	}

	// 4. 构建更新字段Map
	updates := make(map[string]interface{})

	if req.TaskName != nil {
		updates["task_name"] = *req.TaskName
	}
	if req.ExecutionType != nil {
		updates["execution_type"] = *req.ExecutionType
	}
	if req.TaskStatus != nil {
		updates["task_status"] = *req.TaskStatus
	}
	if req.CaseGroupID != nil {
		updates["case_group_id"] = *req.CaseGroupID
	}
	if req.CaseGroupName != nil {
		updates["case_group_name"] = *req.CaseGroupName
	}
	if req.DisplayLanguage != nil {
		updates["display_language"] = *req.DisplayLanguage
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = *req.EndDate
	}
	if req.TestVersion != nil {
		updates["test_version"] = *req.TestVersion
	}
	if req.TestEnv != nil {
		updates["test_env"] = *req.TestEnv
	}
	if req.TestDate != nil {
		updates["test_date"] = *req.TestDate
	}
	if req.Executor != nil {
		updates["executor"] = *req.Executor
	}
	if req.TaskDescription != nil {
		updates["task_description"] = *req.TaskDescription
	}

	// 5. 执行更新
	err = s.repo.UpdateByUUID(taskUUID, updates)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	// 6. 查询更新后的任务
	updatedTask, err := s.repo.GetByUUID(taskUUID)
	if err != nil {
		return nil, fmt.Errorf("get updated task: %w", err)
	}

	return updatedTask, nil
}

// DeleteTask 删除任务(硬删除，级联删除相关数据)
func (s *executionTaskService) DeleteTask(projectID uint, userID uint, taskUUID string) error {
	// 1. 验证任务存在且属于该项目
	task, err := s.repo.GetByUUID(taskUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("任务不存在")
		}
		return fmt.Errorf("get task by uuid: %w", err)
	}

	if task.ProjectID != projectID {
		return errors.New("任务不属于该项目")
	}

	// 2. 级联删除：先删除执行任务的所有执行用例结果（元数据）
	err = s.ecrRepo.DeleteByTaskUUID(taskUUID)
	if err != nil {
		return fmt.Errorf("delete execution case results: %w", err)
	}

	// 3. 执行硬删除执行任务
	err = s.repo.Delete(taskUUID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	// 可选: 记录审计日志
	// TODO: 添加日志记录

	return nil
}

// 辅助函数: 检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}
