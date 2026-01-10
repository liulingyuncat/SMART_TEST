package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite" // 纯Go SQLite驱动
)

// getEnv 读取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 读取整数类型的环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// initDatabase 初始化数据库连接（支持 SQLite 和 PostgreSQL）
func initDatabase() (*gorm.DB, string, error) {
	dbType := getEnv("DB_TYPE", "sqlite")

	var dialector gorm.Dialector

	switch dbType {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			getEnv("DB_HOST", "localhost"),
			getEnvInt("DB_PORT", 5432),
			getEnv("DB_USER", "webtest"),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_NAME", "webtest"),
		)
		dialector = postgres.Open(dsn)

	case "sqlite":
		dbName := getEnv("DB_NAME", "webtest.db")
		dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)", dbName)
		dialector = sqlite.Dialector{
			DriverName: "sqlite",
			DSN:        dsn,
		}

	default:
		return nil, "", fmt.Errorf("unsupported database type: %s", dbType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect database: %w", err)
	}

	return db, dbType, nil
}

func main() {
	// 从环境变量读取数据库配置（支持 SQLite 和 PostgreSQL）
	db, dbType, err := initDatabase()
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	fmt.Printf("数据库类型: %s\n", dbType)

	fmt.Println("开始清理软删除的数据...")
	fmt.Println("警告：此操作将永久删除所有软删除的数据！")

	// 使用事务执行清理
	err = db.Transaction(func(tx *gorm.DB) error {
		tables := []struct {
			name  string
			query string
		}{
			{"users", "DELETE FROM users WHERE deleted_at IS NOT NULL"},
			{"project_members", "DELETE FROM project_members WHERE deleted_at IS NOT NULL"},
			{"projects", "DELETE FROM projects WHERE deleted_at IS NOT NULL"},
			{"case_groups", "DELETE FROM case_groups WHERE deleted_at IS NOT NULL"},
			{"manual_test_cases", "DELETE FROM manual_test_cases WHERE deleted_at IS NOT NULL"},
			{"auto_test_cases", "DELETE FROM auto_test_cases WHERE deleted_at IS NOT NULL"},
			{"api_test_cases", "DELETE FROM api_test_cases WHERE deleted_at IS NOT NULL"},
			{"execution_tasks", "DELETE FROM execution_tasks WHERE deleted_at IS NOT NULL"},
			{"execution_case_results", "DELETE FROM execution_case_results WHERE deleted_at IS NOT NULL"},
			{"defects", "DELETE FROM defects WHERE deleted_at IS NOT NULL"},
			{"defect_attachments", "DELETE FROM defect_attachments WHERE deleted_at IS NOT NULL"},
			{"defect_comments", "DELETE FROM defect_comments WHERE deleted_at IS NOT NULL"},
			{"defect_subjects", "DELETE FROM defect_subjects WHERE deleted_at IS NOT NULL"},
			{"defect_phases", "DELETE FROM defect_phases WHERE deleted_at IS NOT NULL"},
			{"versions", "DELETE FROM versions WHERE deleted_at IS NOT NULL"},
			{"case_versions", "DELETE FROM case_versions WHERE deleted_at IS NOT NULL"},
			{"auto_test_case_versions", "DELETE FROM auto_test_case_versions WHERE deleted_at IS NOT NULL"},
			{"api_test_case_versions", "DELETE FROM api_test_case_versions WHERE deleted_at IS NOT NULL"},
			{"web_case_versions", "DELETE FROM web_case_versions WHERE deleted_at IS NOT NULL"},
			{"requirement_items", "DELETE FROM requirement_items WHERE deleted_at IS NOT NULL"},
			{"viewpoint_items", "DELETE FROM viewpoint_items WHERE deleted_at IS NOT NULL"},
			{"raw_documents", "DELETE FROM raw_documents WHERE deleted_at IS NOT NULL"},
			{"ai_reports", "DELETE FROM ai_reports WHERE deleted_at IS NOT NULL"},
			{"prompts", "DELETE FROM prompts WHERE deleted_at IS NOT NULL"},
			{"case_reviews", "DELETE FROM case_reviews WHERE deleted_at IS NOT NULL"},
			{"case_review_items", "DELETE FROM case_review_items WHERE deleted_at IS NOT NULL"},
		}

		totalDeleted := 0
		for _, table := range tables {
			result := tx.Exec(table.query)
			if result.Error != nil {
				errMsg := result.Error.Error()
				// 如果表不存在deleted_at字段或表不存在，跳过该表
				// 兼容 SQLite 和 PostgreSQL 的错误信息格式
				if strings.Contains(errMsg, "no such column") ||
					strings.Contains(errMsg, "does not exist") ||
					strings.Contains(errMsg, "column") && strings.Contains(errMsg, "deleted_at") {
					fmt.Printf("- %s: 无 deleted_at 字段，跳过\n", table.name)
					continue
				}
				if strings.Contains(errMsg, "no such table") ||
					strings.Contains(errMsg, "relation") && strings.Contains(errMsg, "does not exist") {
					fmt.Printf("- %s: 表不存在，跳过\n", table.name)
					continue
				}
				return fmt.Errorf("清理 %s 失败: %v", table.name, result.Error)
			}
			if result.RowsAffected > 0 {
				fmt.Printf("✓ %s: 删除了 %d 条记录\n", table.name, result.RowsAffected)
				totalDeleted += int(result.RowsAffected)
			}
		}

		fmt.Printf("\n总共删除了 %d 条软删除记录\n", totalDeleted)
		return nil
	})

	if err != nil {
		log.Fatalf("清理失败: %v", err)
	}

	fmt.Println("\n清理完成！")

	// 显示剩余记录统计
	fmt.Println("\n各表剩余记录数：")
	tables := []string{
		"users", "projects", "project_members", "case_groups",
		"manual_test_cases", "auto_test_cases", "api_test_cases",
		"execution_tasks", "defects", "requirement_items",
		"viewpoint_items", "raw_documents", "ai_reports", "prompts",
	}

	for _, table := range tables {
		var count int64
		db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		fmt.Printf("  %s: %d\n", table, count)
	}
}
