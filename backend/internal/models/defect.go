package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DefectStatus 缺陷状态枚举
type DefectStatus string

const (
	DefectStatusNew      DefectStatus = "New"
	DefectStatusActive   DefectStatus = "Active"
	DefectStatusResolved DefectStatus = "Resolved"
	DefectStatusClosed   DefectStatus = "Closed"
)

// ValidDefectStatuses 有效的缺陷状态列表
var ValidDefectStatuses = []DefectStatus{
	DefectStatusNew,
	DefectStatusActive,
	DefectStatusResolved,
	DefectStatusClosed,
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
	DefectSeverityA DefectSeverity = "A"
	DefectSeverityB DefectSeverity = "B"
	DefectSeverityC DefectSeverity = "C"
	DefectSeverityD DefectSeverity = "D"
)

// ValidDefectSeverities 有效的严重程度列表
var ValidDefectSeverities = []DefectSeverity{
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

// Defect 缺陷模型
type Defect struct {
	ID                string `gorm:"type:varchar(36);primaryKey" json:"id"`                                         // UUID主键
	DefectID          string `gorm:"type:varchar(20);uniqueIndex:idx_defects_defect_id;not null" json:"defect_id"`  // 显示ID（XXXXXX）
	ProjectID         uint   `gorm:"not null;index:idx_defects_project_status" json:"project_id"`                   // 所属项目ID
	Title             string `gorm:"type:varchar(200);not null" json:"title"`                                       // 缺陷标题
	Subject           string `gorm:"type:varchar(100)" json:"subject"`                                              // 主题分类
	Description       string `gorm:"type:text" json:"description"`                                                  // 详细描述
	RecoveryMethod    string `gorm:"type:varchar(500)" json:"recovery_method"`                                      // 恢复方法
	Priority          string `gorm:"type:varchar(1);default:'B'" json:"priority"`                                   // 优先级(A/B/C/D)
	Severity          string `gorm:"type:varchar(1);default:'B'" json:"severity"`                                   // 严重程度(A/B/C/D)
	Frequency         string `gorm:"type:varchar(10)" json:"frequency"`                                             // 复现频率
	DetectedInRelease string `gorm:"type:varchar(50)" json:"detected_in_release"`                                   // 发现版本
	Phase             string `gorm:"type:varchar(100)" json:"phase"`                                                // 测试阶段
	CaseID            string `gorm:"type:varchar(50)" json:"case_id"`                                               // 关联的Case ID
	Assignee          string `gorm:"type:varchar(100)" json:"assignee"`                                             // 指派人
	Status            string `gorm:"type:varchar(20);default:'New';index:idx_defects_project_status" json:"status"` // 状态
	CreatedBy         uint   `gorm:"not null" json:"created_by"`                                                    // 创建人ID
	UpdatedBy         uint   `json:"updated_by"`                                                                    // 更新人ID

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
		d.Severity = string(DefectSeverityB)
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
	Title             string `json:"title" binding:"required,max=200"`
	SubjectID         *uint  `json:"subject_id"` // 主题ID
	Subject           string `json:"subject"`    // 兼容直接传名称
	Description       string `json:"description"`
	RecoveryMethod    string `json:"recovery_method"`
	Priority          string `json:"priority"`
	Severity          string `json:"severity"`
	Frequency         string `json:"frequency"`
	DetectedInRelease string `json:"detected_in_release"`
	PhaseID           *uint  `json:"phase_id"`   // 阶段ID
	Phase             string `json:"phase"`      // 兼容直接传名称
	CaseID            string `json:"case_id"`    // 关联的Case ID
	Status            string `json:"status"`     // 状态（导入时使用）
	CreatedAt         string `json:"created_at"` // 创建时间（导入时使用，格式：YYYY-MM-DD）
}

// DefectUpdateRequest 更新缺陷请求
type DefectUpdateRequest struct {
	Title             *string `json:"title"`
	SubjectID         *uint   `json:"subject_id"` // 主题ID
	Subject           *string `json:"subject"`    // 兼容直接传名称
	Description       *string `json:"description"`
	RecoveryMethod    *string `json:"recovery_method"`
	Priority          *string `json:"priority"`
	Severity          *string `json:"severity"`
	Frequency         *string `json:"frequency"`
	DetectedInRelease *string `json:"detected_in_release"`
	PhaseID           *uint   `json:"phase_id"` // 阶段ID
	Phase             *string `json:"phase"`    // 兼容直接传名称
	CaseID            *string `json:"case_id"`  // 关联的Case ID
	Assignee          *string `json:"assignee"`
	Status            *string `json:"status"`
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
