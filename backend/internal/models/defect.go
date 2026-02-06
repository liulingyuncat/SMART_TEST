package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DefectStatus 缺陷状态枚举
type DefectStatus string

const (
	DefectStatusNew        DefectStatus = "New"
	DefectStatusInProgress DefectStatus = "InProgress" // 变更：Active → InProgress
	DefectStatusResolved   DefectStatus = "Resolved"
	DefectStatusClosed     DefectStatus = "Closed"
	DefectStatusConfirmed  DefectStatus = "Confirmed" // 新增
	DefectStatusReopened   DefectStatus = "Reopened"  // 新增
	DefectStatusRejected   DefectStatus = "Rejected"  // 新增
	// 向后兼容：保留Active用于旧数据
	DefectStatusActive DefectStatus = "Active" // deprecated: use InProgress instead
)

// ValidDefectStatuses 有效的缺陷状态列表
var ValidDefectStatuses = []DefectStatus{
	DefectStatusNew,
	DefectStatusInProgress,
	DefectStatusResolved,
	DefectStatusClosed,
	DefectStatusConfirmed,
	DefectStatusReopened,
	DefectStatusRejected,
	DefectStatusActive, // 向后兼容
}

// IsValidDefectStatus 检查状态是否有效
func IsValidDefectStatus(status string) bool {
	for _, s := range ValidDefectStatuses {
		if string(s) == status {
			return true
		}
	}
	return false
}

// DefectPriority 缺陷优先级枚举
type DefectPriority string

const (
	DefectPriorityA DefectPriority = "A"
	DefectPriorityB DefectPriority = "B"
	DefectPriorityC DefectPriority = "C"
	DefectPriorityD DefectPriority = "D"
)

// ValidDefectPriorities 有效的优先级列表
var ValidDefectPriorities = []DefectPriority{
	DefectPriorityA,
	DefectPriorityB,
	DefectPriorityC,
	DefectPriorityD,
}

// IsValidDefectPriority 检查优先级是否有效
func IsValidDefectPriority(priority string) bool {
	for _, p := range ValidDefectPriorities {
		if string(p) == priority {
			return true
		}
	}
	return false
}

// DefectSeverity 缺陷严重程度枚举
type DefectSeverity string

const (
	DefectSeverityCritical DefectSeverity = "Critical" // 变更：A → Critical
	DefectSeverityMajor    DefectSeverity = "Major"    // 变更：B → Major
	DefectSeverityMinor    DefectSeverity = "Minor"    // 变更：C → Minor
	DefectSeverityTrivial  DefectSeverity = "Trivial"  // 变更：D → Trivial
	// 向后兼容：保留ABCD用于旧数据
	DefectSeverityA DefectSeverity = "A" // deprecated: use Critical instead
	DefectSeverityB DefectSeverity = "B" // deprecated: use Major instead
	DefectSeverityC DefectSeverity = "C" // deprecated: use Minor instead
	DefectSeverityD DefectSeverity = "D" // deprecated: use Trivial instead
)

// ValidDefectSeverities 有效的严重程度列表
var ValidDefectSeverities = []DefectSeverity{
	DefectSeverityCritical,
	DefectSeverityMajor,
	DefectSeverityMinor,
	DefectSeverityTrivial,
	// 向后兼容
	DefectSeverityA,
	DefectSeverityB,
	DefectSeverityC,
	DefectSeverityD,
}

// IsValidDefectSeverity 检查严重程度是否有效
func IsValidDefectSeverity(severity string) bool {
	for _, s := range ValidDefectSeverities {
		if string(s) == severity {
			return true
		}
	}
	return false
}

// DefectType 缺陷类型枚举（新增）
type DefectType string

const (
	DefectTypeFunctional      DefectType = "Functional"
	DefectTypeUI              DefectType = "UI"
	DefectTypeUIInteraction   DefectType = "UIInteraction"
	DefectTypeCompatibility   DefectType = "Compatibility"
	DefectTypeBrowserSpecific DefectType = "BrowserSpecific"
	DefectTypePerformance     DefectType = "Performance"
	DefectTypeSecurity        DefectType = "Security"
	DefectTypeEnvironment     DefectType = "Environment"
	DefectTypeUserError       DefectType = "UserError"
)

// ValidDefectTypes 有效的缺陷类型列表
var ValidDefectTypes = []DefectType{
	DefectTypeFunctional,
	DefectTypeUI,
	DefectTypeUIInteraction,
	DefectTypeCompatibility,
	DefectTypeBrowserSpecific,
	DefectTypePerformance,
	DefectTypeSecurity,
	DefectTypeEnvironment,
	DefectTypeUserError,
}

// IsValidDefectType 检查缺陷类型是否有效
func IsValidDefectType(defectType string) bool {
	for _, t := range ValidDefectTypes {
		if string(t) == defectType {
			return true
		}
	}
	return false
}

// Defect 缺陷模型
type Defect struct {
	ID              string `gorm:"type:varchar(36);primaryKey" json:"id"`                                         // UUID主键
	DefectID        string `gorm:"type:varchar(20);uniqueIndex:idx_defects_defect_id;not null" json:"defect_id"`  // 显示ID（XXXXXX）
	ProjectID       uint   `gorm:"not null;index:idx_defects_project_status" json:"project_id"`                   // 所属项目ID
	Title           string `gorm:"type:varchar(200);not null" json:"title"`                                       // 缺陷标题
	Subject         string `gorm:"type:varchar(100)" json:"subject"`                                              // 主题分类
	Description     string `gorm:"type:text" json:"description"`                                                  // 详细描述
	RecoveryMethod  string `gorm:"type:varchar(500)" json:"recovery_method"`                                      // 恢复方法
	Priority        string `gorm:"type:varchar(1);default:'B'" json:"priority"`                                   // 优先级(A/B/C/D)
	Severity        string `gorm:"type:varchar(20);default:'Major'" json:"severity"`                              // 严重程度(Critical/Major/Minor/Trivial)
	Type            string `gorm:"type:varchar(30)" json:"type"`                                                  // 缺陷类型（新增）
	Frequency       string `gorm:"type:varchar(10)" json:"frequency"`                                             // 复现频率
	DetectedVersion string `gorm:"type:varchar(50)" json:"detected_version"`                                      // 发现版本
	Phase           string `gorm:"type:varchar(100)" json:"phase"`                                                // 测试阶段
	CaseID          string `gorm:"type:varchar(50)" json:"case_id"`                                               // 关联的Case ID
	Assignee        string `gorm:"type:varchar(100)" json:"assignee"`                                             // 指派人
	RecoveryRank    string `gorm:"type:varchar(50)" json:"recovery_rank"`                                         // 恢复等级（新增）
	DetectionTeam   string `gorm:"type:varchar(100)" json:"detection_team"`                                       // 检测团队（新增）
	Location        string `gorm:"type:varchar(200)" json:"location"`                                             // 位置（新增）
	FixVersion      string `gorm:"type:varchar(50)" json:"fix_version"`                                           // 修复版本（新增）
	SQAMemo         string `gorm:"type:text" json:"sqa_memo"`                                                     // SQA备注（新增）
	Component       string `gorm:"type:varchar(100)" json:"component"`                                            // 组件（新增）
	Resolution      string `gorm:"type:text" json:"resolution"`                                                   // 解决方案（新增）
	Models          string `gorm:"type:varchar(200)" json:"models"`                                               // 机型（新增）
	DetectedBy      string `gorm:"type:varchar(100)" json:"detected_by"`                                          // 提出人名字
	Status          string `gorm:"type:varchar(20);default:'New';index:idx_defects_project_status" json:"status"` // 状态
	CreatedBy       uint   `gorm:"not null" json:"created_by"`                                                    // 创建人ID
	UpdatedBy       uint   `json:"updated_by"`                                                                    // 更新人ID

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_defects_deleted_at" json:"-"`

	// 关联 - 不存储在数据库
	Attachments   []DefectAttachment `gorm:"foreignKey:DefectID;references:ID" json:"attachments,omitempty"`
	CreatedByUser User               `gorm:"foreignKey:CreatedBy;references:ID" json:"created_by_user,omitempty"` // 创建人信息
}

// TableName 指定表名
func (Defect) TableName() string {
	return "defects"
}

// BeforeCreate GORM钩子：创建前自动生成UUID
func (d *Defect) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	if d.Status == "" {
		d.Status = string(DefectStatusNew)
	}
	if d.Priority == "" {
		d.Priority = string(DefectPriorityB)
	}
	if d.Severity == "" {
		d.Severity = string(DefectSeverityMajor)
	}
	return nil
}

// DefectListResponse 缺陷列表响应
type DefectListResponse struct {
	Defects      []Defect         `json:"defects"`
	Total        int64            `json:"total"`
	Page         int              `json:"page"`
	Size         int              `json:"size"`
	StatusCounts map[string]int64 `json:"status_counts"`
}

// DefectCreateRequest 创建缺陷请求
type DefectCreateRequest struct {
	Title           string `json:"title" binding:"required,max=200"`
	SubjectID       *uint  `json:"subject_id"` // 主题ID
	Subject         string `json:"subject"`    // 兼容直接传名称
	Description     string `json:"description"`
	RecoveryMethod  string `json:"recovery_method"`
	Priority        string `json:"priority"`
	Severity        string `json:"severity"`
	Type            string `json:"type"` // 新增：缺陷类型
	Frequency       string `json:"frequency"`
	DetectedVersion string `json:"detected_version"` // 发现版本
	PhaseID         *uint  `json:"phase_id"`         // 阶段ID
	Phase           string `json:"phase"`            // 兼容直接传名称
	CaseID          string `json:"case_id"`          // 关联的Case ID
	RecoveryRank    string `json:"recovery_rank"`    // 新增：恢复等级
	DetectionTeam   string `json:"detection_team"`   // 新增：检测团队
	Location        string `json:"location"`         // 新增：位置
	FixVersion      string `json:"fix_version"`      // 新增：修复版本
	SQAMemo         string `json:"sqa_memo"`         // 新增：SQA备注
	Component       string `json:"component"`        // 新增：组件
	Resolution      string `json:"resolution"`       // 新增：解决方案
	Models          string `json:"models"`           // 新增：机型
	DetectedBy      string `json:"detected_by"`      // 提出人名字（导入时使用）
	Status          string `json:"status"`           // 状态（导入时使用）
	CreatedAt       string `json:"created_at"`       // 创建时间（导入时使用，格式：YYYY-MM-DD）
}

// DefectUpdateRequest 更新缺陷请求
type DefectUpdateRequest struct {
	Title           *string `json:"title"`
	SubjectID       *uint   `json:"subject_id"` // 主题ID
	Subject         *string `json:"subject"`    // 兼容直接传名称
	Description     *string `json:"description"`
	RecoveryMethod  *string `json:"recovery_method"`
	Priority        *string `json:"priority"`
	Severity        *string `json:"severity"`
	Type            *string `json:"type"` // 新增：缺陷类型
	Frequency       *string `json:"frequency"`
	DetectedVersion *string `json:"detected_version"` // 发现版本
	PhaseID         *uint   `json:"phase_id"`         // 阶段ID
	Phase           *string `json:"phase"`            // 兼容直接传名称
	CaseID          *string `json:"case_id"`          // 关联的Case ID
	Assignee        *string `json:"assignee"`
	RecoveryRank    *string `json:"recovery_rank"`  // 新增：恢复等级
	DetectionTeam   *string `json:"detection_team"` // 新增：检测团队
	Location        *string `json:"location"`       // 新增：位置
	FixVersion      *string `json:"fix_version"`    // 新增：修复版本
	SQAMemo         *string `json:"sqa_memo"`       // 新增：SQA备注
	Component       *string `json:"component"`      // 新增：组件
	Resolution      *string `json:"resolution"`     // 新增：解决方案
	Models          *string `json:"models"`         // 新增：机型
	DetectedBy      *string `json:"detected_by"`    // 提出人名字
	Status          *string `json:"status"`
}

// ImportError 导入错误记录
type ImportError struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
}

// ImportResult 导入结果
type ImportResult struct {
	SuccessCount int           `json:"success_count"`
	FailCount    int           `json:"fail_count"`
	Errors       []ImportError `json:"errors"`
}
