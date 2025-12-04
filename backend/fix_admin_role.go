package main

import (
	"log"
	"time"
	"webtest/config"
	"webtest/internal/models"
)

func main() {
	// 数据库配置
	dbConfig := &config.DatabaseConfig{
		Type:            "sqlite",
		DBName:          "webtest.db",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	// 初始化数据库
	db, err := config.InitDatabase(dbConfig)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// 查询所有用户
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("failed to query users: %v", err)
	}

	log.Println("Current users in database:")
	for _, user := range users {
		log.Printf("ID: %d, Username: %s, Nickname: %s, Role: %s\n", user.ID, user.Username, user.Nickname, user.Role)
	}

	// 更新admin用户的角色
	result := db.Model(&models.User{}).Where("username = ?", "admin").Update("role", "system_admin")
	if result.Error != nil {
		log.Fatalf("failed to update admin role: %v", result.Error)
	}
	log.Printf("Updated %d user(s)\n", result.RowsAffected)

	// 更新root用户的角色
	result = db.Model(&models.User{}).Where("username = ?", "root").Update("role", "system_admin")
	if result.Error != nil {
		log.Fatalf("failed to update root role: %v", result.Error)
	}
	log.Printf("Updated %d user(s)\n", result.RowsAffected)

	// 再次查询确认
	users = []models.User{}
	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("failed to query users: %v", err)
	}

	log.Println("\nUsers after update:")
	for _, user := range users {
		log.Printf("ID: %d, Username: %s, Nickname: %s, Role: %s\n", user.ID, user.Username, user.Nickname, user.Role)
	}
}
