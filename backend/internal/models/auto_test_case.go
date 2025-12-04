package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AutoTestCase 自动化测试用例模型
type AutoTestCase struct {
	CaseID    string `gorm:"type:varchar(36);primaryKey" json:"case_id"`  // UUID主键
	ID        uint   `gorm:"not null;index:idx_atc_display_id" json:"id"` // 显示序号(可重排,用于前端展示)
	ProjectID uint   `gorm:"not null;index:idx_atc_project" json:"project_id"`
	CaseType  string `gorm:"type:varchar(20);default:'role1';index:idx_atc_type" json:"case_type"` // role1/role2/role3/role4

	// 元数据字段(冗余存储便于导出)
	TestVersion string `gorm:"type:varchar(50)" json:"test_version"` // 测试版本
	TestDate    string `gorm:"type:varchar(20)" json:"test_date"`    // 测试日期(YYYY-MM-DD)

	// 公共字段(不区分语言)
	CaseNumber string `gorm:"type:varchar(50)" json:"case_number"`                                   // 用例编号
	TestResult string `gorm:"type:varchar(10);default:'NR';index:idx_atc_result" json:"test_result"` // OK/NG/NR
	Remark     string `gorm:"type:text" json:"remark"`                                               // 备考

	// ======== 多语言字段 - 画面(三语言) ========
	ScreenCN string `gorm:"type:varchar(100)" json:"screen_cn"`
	ScreenJP string `gorm:"type:varchar(100)" json:"screen_jp"`
	ScreenEN string `gorm:"type:varchar(100)" json:"screen_en"`

	// ======== 多语言字段 - 功能(三语言,简化为单一功能字段) ========
	FunctionCN string `gorm:"type:varchar(200);index:idx_atc_function_cn" json:"function_cn"`
	FunctionJP string `gorm:"type:varchar(200);index:idx_atc_function_jp" json:"function_jp"`
	FunctionEN string `gorm:"type:varchar(200);index:idx_atc_function_en" json:"function_en"`

	// ======== 多语言字段 - 前置条件(三语言) ========
	PreconditionCN string `gorm:"type:text" json:"precondition_cn"`
	PreconditionJP string `gorm:"type:text" json:"precondition_jp"`
	PreconditionEN string `gorm:"type:text" json:"precondition_en"`

	// ======== 多语言字段 - 测试步骤(三语言) ========
	TestStepsCN string `gorm:"type:text" json:"test_steps_cn"`
	TestStepsJP string `gorm:"type:text" json:"test_steps_jp"`
	TestStepsEN string `gorm:"type:text" json:"test_steps_en"`

	// ======== 多语言字段 - 期待值(三语言) ========
	ExpectedResultCN string `gorm:"type:text" json:"expected_result_cn"`
	ExpectedResultJP string `gorm:"type:text" json:"expected_result_jp"`
	ExpectedResultEN string `gorm:"type:text" json:"expected_result_en"`

	// 审计字段
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_atc_deleted_at" json:"-"` // 软删除
}

// TableName 指定表名
func (AutoTestCase) TableName() string {
	return "auto_test_cases"
}

// BeforeCreate GORM钩子:创建前自动生成UUID
func (a *AutoTestCase) BeforeCreate(tx *gorm.DB) error {
	if a.CaseID == "" {
		a.CaseID = uuid.New().String()
	}
	return nil
}
