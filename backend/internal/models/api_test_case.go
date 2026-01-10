package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Method枚举常量
const (
	MethodGET    = "GET"
	MethodPOST   = "POST"
	MethodPUT    = "PUT"
	MethodDELETE = "DELETE"
	MethodPATCH  = "PATCH"
)

// TestResult枚举常量
const (
	TestResultNR = "NR" // Not Run 未执行
	TestResultOK = "OK" // Passed 通过
	TestResultNG = "NG" // Failed 失败
)

// ApiTestCase 接口测试用例模型
type ApiTestCase struct {
	// 主键(UUID)
	ID string `gorm:"type:varchar(36);primaryKey" json:"case_id"` // UUID主键,全局唯一标识符(前端使用case_id字段名)

	// 基本信息
	ProjectID uint   `gorm:"not null;index:idx_api_cases_project" json:"project_id"`
	CaseType  string `gorm:"type:varchar(20);default:'api';index:idx_api_cases_type" json:"case_type"` // 接口用例类型标识
	CaseGroup string `gorm:"type:varchar(100);default:'';index:idx_api_cases_group" json:"case_group"` // 用例集名称

	// 用例基本字段
	CaseNumber string `gorm:"type:varchar(50)" json:"case_number"` // 用例编号(用户自定义)
	Screen     string `gorm:"type:varchar(100)" json:"screen"`     // 画面/接口所属模块

	// API专属字段
	URL      string `gorm:"type:text" json:"url"`                         // 接口地址
	Header   string `gorm:"type:text" json:"header"`                      // 请求头(支持多行)
	Method   string `gorm:"type:varchar(10);default:'GET'" json:"method"` // HTTP方法: GET/POST/PUT/DELETE/PATCH
	Body     string `gorm:"type:text" json:"body"`                        // 请求体(支持多行)
	Response string `gorm:"type:text" json:"response"`                    // 预期响应(支持多行)

	// 可执行脚本(消除AI幻觉)
	ScriptCode string `gorm:"type:text" json:"script_code"` // JS脚本代码，S12直接执行此脚本

	// 测试结果
	TestResult string `gorm:"type:varchar(10);default:'NR';index:idx_api_cases_result" json:"test_result"` // OK/NG/NR
	Remark     string `gorm:"type:text" json:"remark"`                                                     // 备注

	// 显示顺序(用于计算前端No.列序号)
	DisplayOrder int `gorm:"index:idx_api_cases_display_order" json:"display_order"` // 显示顺序,前端根据此字段计算No.序号(1,2,3...)

	// 审计字段
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_api_cases_deleted_at" json:"-"` // 软删除
}

// TableName 指定表名
func (ApiTestCase) TableName() string {
	return "api_test_cases"
}

// BeforeCreate GORM钩子:创建前自动生成UUID
func (a *ApiTestCase) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// Validate 字段验证
func (a *ApiTestCase) Validate() error {
	// 验证Method字段
	validMethods := map[string]bool{
		MethodGET:    true,
		MethodPOST:   true,
		MethodPUT:    true,
		MethodDELETE: true,
		MethodPATCH:  true,
	}
	if a.Method != "" && !validMethods[a.Method] {
		return errors.New("method必须是GET, POST, PUT, DELETE, PATCH之一")
	}

	// 验证TestResult字段
	validResults := map[string]bool{
		TestResultNR: true,
		TestResultOK: true,
		TestResultNG: true,
	}
	if a.TestResult != "" && !validResults[a.TestResult] {
		return errors.New("test_result必须是NR, OK, NG之一")
	}

	return nil
}

// ApiTestCaseVersion 接口测试用例版本管理模型
type ApiTestCaseVersion struct {
	// 主键(UUID)
	ID string `gorm:"type:varchar(36);primaryKey" json:"id"` // 版本UUID

	// 关联信息
	ProjectID uint `gorm:"not null;index:idx_api_versions_project" json:"project_id"`

	// XLSX文件名(包含所有用例集)
	XlsxFilename string `gorm:"type:varchar(255);not null" json:"xlsx_filename"` // XLSX文件名

	// 兼容旧版本的CSV文件名(废弃,保留用于数据迁移)
	FilenameRole1 string `gorm:"type:varchar(255);default:''" json:"filename_role1,omitempty"` // ROLE1文件名(已废弃)
	FilenameRole2 string `gorm:"type:varchar(255);default:''" json:"filename_role2,omitempty"` // ROLE2文件名(已废弃)
	FilenameRole3 string `gorm:"type:varchar(255);default:''" json:"filename_role3,omitempty"` // ROLE3文件名(已废弃)
	FilenameRole4 string `gorm:"type:varchar(255);default:''" json:"filename_role4,omitempty"` // ROLE4文件名(已废弃)

	// 版本信息
	Remark    string `gorm:"type:text" json:"remark"`    // 版本备注(限制500字符,前端校验)
	CreatedBy uint   `gorm:"not null" json:"created_by"` // 创建人ID

	// 审计字段
	CreatedAt time.Time `gorm:"index:idx_api_versions_created_at" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (ApiTestCaseVersion) TableName() string {
	return "api_test_case_versions"
}

// BeforeCreate GORM钩子:创建前自动生成UUID
func (v *ApiTestCaseVersion) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	return nil
}
