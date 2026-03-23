package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"Resume/backend/config"
	"Resume/backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBCharset,
		cfg.DBParseTime,
		cfg.DBLoc,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取通用数据库对象失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("数据库连接成功建立！")

	autoMigrate(DB)
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.FileInfo{},
		&models.Interview{},
		&models.CompanyPosition{},
		&models.RecruitmentSource{},
		&models.Employee{},
	)
	if err != nil {
		log.Fatalf("数据库自动迁移失败: %v", err)
	}
	log.Println("数据库自动迁移完成。")
}
