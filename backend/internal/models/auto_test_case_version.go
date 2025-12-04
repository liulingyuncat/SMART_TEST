package models

import "time"

// AutoTestCaseVersion 自动化测试用例版本记录模型
type AutoTestCaseVersion struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VersionID   string    `gorm:"type:varchar(100);not null;index:idx_auto_versions_project_version" json:"version_id"`
	ProjectID   uint      `gorm:"not null;index:idx_auto_versions_project_version" json:"project_id"`
	ProjectName string    `gorm:"type:varchar(100);not null" json:"project_name"`
	RoleType    string    `gorm:"type:varchar(10);not null;check:role_type IN ('role1','role2','role3','role4');index:idx_auto_versions_role" json:"role_type"`
	Filename    string    `gorm:"type:varchar(255);not null" json:"filename"`
	FilePath    string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize    int64     `gorm:"type:bigint" json:"file_size"`
	CaseCount   int       `gorm:"default:0" json:"case_count"`
	Remark      string    `gorm:"type:varchar(200);default:''" json:"remark"`
	CreatedBy   *uint     `gorm:"index" json:"created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index:idx_auto_versions_created,sort:desc" json:"created_at"`
}

// TableName 指定表名
func (AutoTestCaseVersion) TableName() string {
	return "auto_test_case_versions"
}
