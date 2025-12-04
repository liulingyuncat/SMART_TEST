package main

import (
	"fmt"
	"log"
	"webtest/config"
	"webtest/internal/models"

	"time"
)

func main() {
	// 初始化数据库
	dbConfig := &config.DatabaseConfig{
		Type:            "sqlite",
		DBName:          "webtest.db",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	db, err := config.InitDatabase(dbConfig)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// 查询所有项目
	var projects []models.Project
	db.Find(&projects)
	fmt.Println("=== Projects ===")
	for _, p := range projects {
		fmt.Printf("ID: %d, Name: %s\n", p.ID, p.Name)
	}

	// 查询所有项目成员
	var members []models.ProjectMember
	db.Find(&members)
	fmt.Println("\n=== Project Members ===")
	for _, m := range members {
		fmt.Printf("ProjectID: %d, UserID: %d, Role: %s\n", m.ProjectID, m.UserID, m.Role)
	}

	// 查询所有用户
	var users []models.User
	db.Find(&users)
	fmt.Println("\n=== Users ===")
	for _, u := range users {
		fmt.Printf("ID: %d, Username: %s, Role: %s\n", u.ID, u.Username, u.Role)
	}
}
