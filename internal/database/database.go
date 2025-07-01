package database

import (
	"fmt"
	"log"
	"static-hosting-server/internal/config"
	"static-hosting-server/internal/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Initialize(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// 构建DSN连接字符串，添加超时参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&timeout=30s&readTimeout=30s&writeTimeout=30s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	log.Printf("Attempting to connect to database at %s:%d", cfg.Host, cfg.Port)

	var err error
	var retries = 5

	// 重试连接机制
	for i := 0; i < retries; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			// 连接成功，测试数据库连接
			sqlDB, err := DB.DB()
			if err == nil {
				// 设置连接池参数
				sqlDB.SetMaxOpenConns(10)
				sqlDB.SetMaxIdleConns(5)
				sqlDB.SetConnMaxLifetime(time.Hour)

				// 测试连接
				if err = sqlDB.Ping(); err == nil {
					break
				}
			}
		}

		log.Printf("Database connection attempt %d/%d failed: %v", i+1, retries, err)
		if i < retries-1 {
			log.Printf("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", retries, err)
	}

	// 自动迁移数据库
	if err := autoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return DB, nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&models.Article{},
		&models.User{},
		&models.APIKey{},
		&models.Certificate{},
	)
}

func GetDB() *gorm.DB {
	return DB
}
