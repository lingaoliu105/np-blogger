package database

import (
	"fmt"
	"np-blogger/internal/config"
	"np-blogger/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// Initialize 初始化数据库连接
func Initialize(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移数据库结构
	err = db.AutoMigrate(
		&model.User{},
		&model.Repository{},
		&model.Blog{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	return nil
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}