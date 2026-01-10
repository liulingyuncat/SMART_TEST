package models

import (
	"time"

	"gorm.io/gorm"
)

// Version 通用版本记录模型
type Version struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProjectID uint           `gorm:"not null;index" json:"project_id"`
	DocType   string         `gorm:"type:varchar(50)" json:"doc_type"`  // 旧字段(兼容)
	ItemType  string         `gorm:"type:varchar(50)" json:"item_type"` // 新字段:版本类型(requirement-batch, viewpoint-batch等)
	Filename  string         `gorm:"type:varchar(255)" json:"filename"`
	FilePath  string         `gorm:"type:varchar(500)" json:"file_path"`
	FileSize  int64          `json:"file_size"`
	FileList  string         `gorm:"type:text" json:"file_list"` // 新字段:JSON数组字符串,存储ZIP内包含的MD文件列表
	Remark    string         `gorm:"type:text" json:"remark"`
	CreatedBy *uint          `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 外键关联
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定表名
func (Version) TableName() string {
	return "versions"
}
