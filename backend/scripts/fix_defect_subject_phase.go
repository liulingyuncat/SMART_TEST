package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Defect struct {
	ID        string `gorm:"primaryKey"`
	ProjectID uint
	Subject   string
	Phase     string
}

type DefectSubject struct {
	ID        uint `gorm:"primaryKey"`
	ProjectID uint
	Name      string
}

type DefectPhase struct {
	ID        uint `gorm:"primaryKey"`
	ProjectID uint
	Name      string
}

func main() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("webtest.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 获取所有Subject和Phase配置
	var subjects []DefectSubject
	if err := db.Find(&subjects).Error; err != nil {
		log.Fatalf("Failed to load subjects: %v", err)
	}

	var phases []DefectPhase
	if err := db.Find(&phases).Error; err != nil {
		log.Fatalf("Failed to load phases: %v", err)
	}

	// 构建名称到ID的映射（按项目分组）
	subjectMap := make(map[uint]map[string]uint) // projectID -> name -> id
	for _, s := range subjects {
		if subjectMap[s.ProjectID] == nil {
			subjectMap[s.ProjectID] = make(map[string]uint)
		}
		subjectMap[s.ProjectID][s.Name] = s.ID
	}

	phaseMap := make(map[uint]map[string]uint) // projectID -> name -> id
	for _, p := range phases {
		if phaseMap[p.ProjectID] == nil {
			phaseMap[p.ProjectID] = make(map[string]uint)
		}
		phaseMap[p.ProjectID][p.Name] = p.ID
	}

	// 获取所有缺陷
	var defects []Defect
	if err := db.Where("subject = '' OR phase = ''").Find(&defects).Error; err != nil {
		log.Fatalf("Failed to load defects: %v", err)
	}

	log.Printf("Found %d defects with empty subject or phase", len(defects))

	// 更新缺陷（这里只是示例，实际需要根据前端发送的ID来填充）
	// 由于历史数据没有保存subject_id和phase_id，我们无法自动修复
	// 需要手动处理或者在前端重新编辑后保存

	log.Println("Note: Historical defects cannot be automatically fixed because subject_id and phase_id were not stored.")
	log.Println("Please re-edit and save these defects in the frontend to populate the Subject and Phase fields.")

	// 显示需要修复的缺陷列表
	for _, d := range defects {
		fmt.Printf("Defect ID: %s, Project: %d, Subject: '%s', Phase: '%s'\n",
			d.ID, d.ProjectID, d.Subject, d.Phase)
	}
}
