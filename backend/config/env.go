package config

import (
	"os"
	"strconv"
	"time"
)

// GetDatabaseConfigFromEnv 从环境变量读取数据库配置
// 支持 SQLite (开发) 和 PostgreSQL (生产) 两种数据库
func GetDatabaseConfigFromEnv() *DatabaseConfig {
	dbType := getEnv("DB_TYPE", "sqlite")

	cfg := &DatabaseConfig{
		Type:            dbType,
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvInt("DB_PORT", 5432),
		User:            getEnv("DB_USER", "webtest"),
		Password:        getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "webtest.db"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 1*time.Minute),
	}

	// 对于 SQLite，如果 DB_NAME 没有指定路径，使用默认文件名
	if dbType == "sqlite" && cfg.DBName == "" {
		cfg.DBName = "webtest.db"
	}

	return cfg
}

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

// getEnvDuration 读取时间间隔类型的环境变量
// 支持格式: "5m", "1h", "30s" 等
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetEnvBool 读取布尔类型的环境变量
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
