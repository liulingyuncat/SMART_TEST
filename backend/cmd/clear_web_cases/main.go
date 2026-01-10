package main

import (
	"fmt"
	"log"
	"time"
	"webtest/config"
	"webtest/internal/models"

	"gorm.io/gorm"
)

func main() {
	// 初始化数据库配置
	dbConfig := &config.DatabaseConfig{
		Type:            "sqlite",
		DBName:          "webtest.db",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
	}

	// 连接数据库
	db, err := config.InitDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 查询Web用例数量
	var count int64
	if err := db.Model(&models.AutoTestCase{}).Where("case_type = ? AND deleted_at IS NULL", "web").Count(&count).Error; err != nil {
		log.Fatalf("Failed to count web cases: %v", err)
	}

	fmt.Printf("找到 %d 条Web用例\n", count)

	if count == 0 {
		fmt.Println("没有Web用例需要删除")
		return
	}

	// 软删除所有Web用例
	result := db.Model(&models.AutoTestCase{}).Where("case_type = ? AND deleted_at IS NULL", "web").Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP"))
	if result.Error != nil {
		log.Fatalf("Failed to delete web cases: %v", result.Error)
	}

	fmt.Printf("✅ 成功删除 %d 条Web用例\n", result.RowsAffected)
}
