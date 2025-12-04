package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExecutionTask 测试执行任务模型
type ExecutionTask struct {
	TaskUUID        string         `gorm:"type:varchar(36);primaryKey" json:"task_uuid" validate:"omitempty,uuid4"`
	ProjectID       uint           `gorm:"not null;index:idx_tet_project" json:"project_id" validate:"required,min=1"`
	TaskName        string         `gorm:"type:varchar(50);not null" json:"task_name" validate:"required,min=1,max=50"`
	ExecutionType   string         `gorm:"type:varchar(20);not null" json:"execution_type" validate:"required,oneof=manual automation api"`
	TaskStatus      string         `gorm:"type:varchar(20);not null;default:pending;index:idx_tet_status" json:"task_status" validate:"omitempty,oneof=pending in_progress completed"`
	StartDate       *time.Time     `gorm:"type:date" json:"start_date" validate:"omitempty"`
	EndDate         *time.Time     `gorm:"type:date" json:"end_date" validate:"omitempty,gtefield=StartDate"`
	TestVersion     string         `gorm:"type:varchar(50)" json:"test_version" validate:"omitempty,max=50"`
	TestEnv         string         `gorm:"type:varchar(100)" json:"test_env" validate:"omitempty,max=100"`
	TestDate        *time.Time     `gorm:"type:date" json:"test_date" validate:"omitempty"`
	Executor        string         `gorm:"type:varchar(50)" json:"executor" validate:"omitempty,max=50"`
	TaskDescription string         `gorm:"type:text" json:"task_description" validate:"omitempty,max=2000"`
	CreatedBy       uint           `gorm:"not null;index:idx_tet_creator" json:"created_by" validate:"required,min=1"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ExecutionTask) TableName() string {
	return "test_execution_tasks"
}

// BeforeCreate GORM钩子：创建前自动生成UUID
func (t *ExecutionTask) BeforeCreate(tx *gorm.DB) error {
	if t.TaskUUID == "" {
		t.TaskUUID = uuid.New().String()
	}
	// 默认状态为 pending
	if t.TaskStatus == "" {
		t.TaskStatus = "pending"
	}
	return nil
}

// ValidateDateRange 验证日期范围逻辑
func (t *ExecutionTask) ValidateDateRange() error {
	if t.StartDate != nil && t.EndDate != nil {
		if t.EndDate.Before(*t.StartDate) {
			return gorm.ErrInvalidData
		}
	}
	return nil
}
