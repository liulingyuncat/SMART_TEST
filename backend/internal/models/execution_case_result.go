package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// ExecutionCaseResult 测试执行用例结果模型
type ExecutionCaseResult struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskUUID   string `gorm:"type:varchar(36);not null;uniqueIndex:idx_task_case" json:"task_uuid" validate:"required,uuid4"`
	CaseID     string `gorm:"type:varchar(36);not null;uniqueIndex:idx_task_case;index:idx_ecr_case_id" json:"case_id" validate:"required"`
	DisplayID  uint   `gorm:"type:int;not null;default:0" json:"display_id"` // 用例显示ID（序号）
	CaseNum    string `gorm:"type:varchar(100)" json:"case_num"`             // 用户自定义CaseID
	CaseType   string `gorm:"type:varchar(20);not null" json:"case_type" validate:"required,oneof=overall acceptance change ai role1 role2 role3 role4 api"`
	TestResult string `gorm:"type:varchar(10);not null;default:NR;index:idx_ecr_test_result" json:"test_result" validate:"required,oneof=NR OK NG Block"`
	BugID      string `gorm:"type:varchar(50)" json:"bug_id" validate:"omitempty,max=50"`
	Remark     string `gorm:"type:text" json:"remark" validate:"omitempty"`

	// 用例内容快照 - 中文
	ScreenCN         string `gorm:"type:varchar(500)" json:"screen_cn"`
	FunctionCN       string `gorm:"type:varchar(500)" json:"function_cn"`
	MajorFunctionCN  string `gorm:"type:varchar(500)" json:"major_function_cn"`  // 手工测试用例-大功能
	MiddleFunctionCN string `gorm:"type:varchar(500)" json:"middle_function_cn"` // 手工测试用例-中功能
	MinorFunctionCN  string `gorm:"type:varchar(500)" json:"minor_function_cn"`  // 手工测试用例-小功能
	PreconditionCN   string `gorm:"type:text" json:"precondition_cn"`
	TestStepsCN      string `gorm:"type:text" json:"test_steps_cn"`
	ExpectedResultCN string `gorm:"type:text" json:"expected_result_cn"`

	// 用例内容快照 - 日文
	ScreenJP         string `gorm:"type:varchar(500)" json:"screen_jp"`
	FunctionJP       string `gorm:"type:varchar(500)" json:"function_jp"`
	MajorFunctionJP  string `gorm:"type:varchar(500)" json:"major_function_jp"`
	MiddleFunctionJP string `gorm:"type:varchar(500)" json:"middle_function_jp"`
	MinorFunctionJP  string `gorm:"type:varchar(500)" json:"minor_function_jp"`
	PreconditionJP   string `gorm:"type:text" json:"precondition_jp"`
	TestStepsJP      string `gorm:"type:text" json:"test_steps_jp"`
	ExpectedResultJP string `gorm:"type:text" json:"expected_result_jp"`

	// 用例内容快照 - 英文
	ScreenEN         string `gorm:"type:varchar(500)" json:"screen_en"`
	FunctionEN       string `gorm:"type:varchar(500)" json:"function_en"`
	MajorFunctionEN  string `gorm:"type:varchar(500)" json:"major_function_en"`
	MiddleFunctionEN string `gorm:"type:varchar(500)" json:"middle_function_en"`
	MinorFunctionEN  string `gorm:"type:varchar(500)" json:"minor_function_en"`
	PreconditionEN   string `gorm:"type:text" json:"precondition_en"`
	TestStepsEN      string `gorm:"type:text" json:"test_steps_en"`
	ExpectedResultEN string `gorm:"type:text" json:"expected_result_en"`

	UpdatedBy uint           `gorm:"not null" json:"updated_by" validate:"required,min=1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ExecutionCaseResult) TableName() string {
	return "execution_case_results"
}

// Validate 验证字段值
func (e *ExecutionCaseResult) Validate() error {
	// 验证TestResult枚举值
	validResults := map[string]bool{
		"NR":    true,
		"OK":    true,
		"NG":    true,
		"Block": true,
	}
	if !validResults[e.TestResult] {
		return errors.New("test_result must be one of: NR, OK, NG, Block")
	}

	// 验证CaseType枚举值
	validTypes := map[string]bool{
		"overall":    true,
		"acceptance": true,
		"change":     true,
		"ai":         true,
		"role1":      true,
		"role2":      true,
		"role3":      true,
		"role4":      true,
		"api":        true,
	}
	if !validTypes[e.CaseType] {
		return errors.New("case_type must be one of: overall, acceptance, change, ai, role1-4, api")
	}

	return nil
}

// BeforeUpdate GORM钩子：更新前验证
func (e *ExecutionCaseResult) BeforeUpdate(tx *gorm.DB) error {
	return e.Validate()
}

// BeforeCreate GORM钩子：创建前验证
func (e *ExecutionCaseResult) BeforeCreate(tx *gorm.DB) error {
	// 默认TestResult为NR
	if e.TestResult == "" {
		e.TestResult = "NR"
	}
	return e.Validate()
}
