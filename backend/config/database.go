package config

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "modernc.org/sqlite" // 导入纯 Go SQLite 驱动
)

// DatabaseConfig 数据库配置结构
type DatabaseConfig struct {
	Type            string        // "postgres", "sqlite"
	Host            string        // 数据库主机
	Port            int           // 数据库端口
	User            string        // 用户名
	Password        string        // 密码
	DBName          string        // 数据库名
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
	ConnMaxIdleTime time.Duration // 连接最大空闲时间
}

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Type {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
		)
		dialector = postgres.Open(dsn)

	case "sqlite":
		// 使用纯 Go SQLite 驱动 (modernc.org/sqlite)
		dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)", cfg.DBName)
		dialector = sqlite.Dialector{
			DriverName: "sqlite",
			DSN:        dsn,
		}

	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	// GORM 配置
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取通用数据库对象 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池参数(满足 NFR-01 性能要求)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)       // 默认 25
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)       // 默认 10
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 默认 5 分钟
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // 默认 1 分钟

	return db, nil
}
