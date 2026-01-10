package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ManualTestCase 手工测试用例模型
type ManualTestCase struct {
	CaseID    string `gorm:"type:varchar(36);primaryKey" json:"case_id"`  // UUID主键
	ID        uint   `gorm:"not null;index:idx_mtc_display_id" json:"id"` // 显示序号（可修改，用于排序）
	ProjectID uint   `gorm:"not null;index:idx_mtc_project" json:"project_id"`
	CaseType  string `gorm:"type:varchar(20);default:'overall';index:idx_mtc_type" json:"case_type"` // ai(AI用例)/overall(整体用例)/change(变更用例)

	// 元数据字段
	TestVersion string `gorm:"type:varchar(50)" json:"test_version"`
	TestEnv     string `gorm:"type:varchar(100)" json:"test_env"`
	TestDate    string `gorm:"type:varchar(20)" json:"test_date"` // YYYY-MM-DD 格式
	Executor    string `gorm:"type:varchar(50)" json:"executor"`
	CaseNumber  string `gorm:"type:varchar(50)" json:"case_number"`
	CaseGroup   string `gorm:"type:varchar(100);index:idx_mtc_case_group" json:"case_group"` // 用例集名称

	// ======== 单语言字段(AI用例使用) ========
	MajorFunction  string `gorm:"type:varchar(100);index:idx_mtc_major_func" json:"major_function"` // ai用例使用
	MiddleFunction string `gorm:"type:varchar(100)" json:"middle_function"`                         // ai用例使用
	MinorFunction  string `gorm:"type:varchar(100)" json:"minor_function"`                          // ai用例使用
	Precondition   string `gorm:"type:text" json:"precondition"`                                    // ai用例使用
	TestSteps      string `gorm:"type:text" json:"test_steps"`                                      // ai用例使用
	ExpectedResult string `gorm:"type:text" json:"expected_result"`                                 // ai用例使用

	// ======== 多语言字段(整体/变更用例使用) ========
	// 大功能 - 三语言
	MajorFunctionCN string `gorm:"type:varchar(100);index:idx_mtc_major_func_cn" json:"major_function_cn"` // overall/change用例使用
	MajorFunctionJP string `gorm:"type:varchar(100);index:idx_mtc_major_func_jp" json:"major_function_jp"` // overall/change用例使用
	MajorFunctionEN string `gorm:"type:varchar(100);index:idx_mtc_major_func_en" json:"major_function_en"` // overall/change用例使用

	// 中功能 - 三语言
	MiddleFunctionCN string `gorm:"type:varchar(100)" json:"middle_function_cn"` // overall/change用例使用
	MiddleFunctionJP string `gorm:"type:varchar(100)" json:"middle_function_jp"` // overall/change用例使用
	MiddleFunctionEN string `gorm:"type:varchar(100)" json:"middle_function_en"` // overall/change用例使用

	// 小功能 - 三语言
	MinorFunctionCN string `gorm:"type:varchar(100)" json:"minor_function_cn"` // overall/change用例使用
	MinorFunctionJP string `gorm:"type:varchar(100)" json:"minor_function_jp"` // overall/change用例使用
	MinorFunctionEN string `gorm:"type:varchar(100)" json:"minor_function_en"` // overall/change用例使用

	// 前置条件 - 三语言
	PreconditionCN string `gorm:"type:text" json:"precondition_cn"` // overall/change用例使用
	PreconditionJP string `gorm:"type:text" json:"precondition_jp"` // overall/change用例使用
	PreconditionEN string `gorm:"type:text" json:"precondition_en"` // overall/change用例使用

	// 测试步骤 - 三语言
	TestStepsCN string `gorm:"type:text" json:"test_steps_cn"` // overall/change用例使用
	TestStepsJP string `gorm:"type:text" json:"test_steps_jp"` // overall/change用例使用
	TestStepsEN string `gorm:"type:text" json:"test_steps_en"` // overall/change用例使用

	// 期待值 - 三语言
	ExpectedResultCN string `gorm:"type:text" json:"expected_result_cn"` // overall/change用例使用
	ExpectedResultJP string `gorm:"type:text" json:"expected_result_jp"` // overall/change用例使用
	ExpectedResultEN string `gorm:"type:text" json:"expected_result_en"` // overall/change用例使用

	// 共用字段
	TestResult string `gorm:"type:varchar(10);default:'NR'" json:"test_result"` // OK/NG/Block/NR (仅overall/change使用)
	Remark     string `gorm:"type:text" json:"remark"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_mtc_deleted_at" json:"-"`
}

// TableName 指定表名
func (ManualTestCase) TableName() string {
	return "manual_test_cases"
}

// BeforeCreate GORM钩子：创建前自动生成UUID
func (m *ManualTestCase) BeforeCreate(tx *gorm.DB) error {
	if m.CaseID == "" {
		m.CaseID = uuid.New().String()
	}
	return nil
}
