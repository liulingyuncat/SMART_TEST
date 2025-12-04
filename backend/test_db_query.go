//go:build ignore
// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "file:webtest.db")
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer db.Close()

	// 需要添加的列
	columnsToAdd := []struct {
		name     string
		dataType string
	}{
		{"deleted_at", "DATETIME"},
		{"major_function_cn", "VARCHAR(500)"},
		{"major_function_jp", "VARCHAR(500)"},
		{"major_function_en", "VARCHAR(500)"},
		{"middle_function_cn", "VARCHAR(500)"},
		{"middle_function_jp", "VARCHAR(500)"},
		{"middle_function_en", "VARCHAR(500)"},
		{"minor_function_cn", "VARCHAR(500)"},
		{"minor_function_jp", "VARCHAR(500)"},
		{"minor_function_en", "VARCHAR(500)"},
	}

	for _, col := range columnsToAdd {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('execution_case_results') WHERE name=?", col.name).Scan(&count)
		if err != nil {
			log.Printf("Check column %s error: %v", col.name, err)
			continue
		}

		if count == 0 {
			fmt.Printf("Adding %s column...\n", col.name)
			_, err = db.Exec(fmt.Sprintf("ALTER TABLE execution_case_results ADD COLUMN %s %s", col.name, col.dataType))
			if err != nil {
				log.Printf("Add column %s error: %v", col.name, err)
			} else {
				fmt.Printf("Column %s added successfully\n", col.name)
			}
		} else {
			fmt.Printf("Column %s already exists\n", col.name)
		}
	}

	// 创建索引
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_ecr_deleted_at ON execution_case_results(deleted_at)")
	if err != nil {
		log.Printf("Create index error: %v", err)
	}

	// 获取总数
	var total int
	db.QueryRow("SELECT COUNT(*) FROM execution_case_results").Scan(&total)
	fmt.Printf("Total records: %d\n", total)

	fmt.Println("Migration completed!")
}
