package config

import (
	"log"

	"github.com/Zhoudf/blog_backend_by_go/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error
	// 配置 MySQL 连接信息，根据实际情况修改
	dsn := "root:Zdf?0830@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 显示 SQL 日志
	})
	if err != nil {
		log.Printf("数据库连接失败: %v", err)
		return err
	}

	// 自动迁移模型
	err = DB.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
	)
	if err != nil {
		log.Printf("模型迁移失败: %v", err)
		return err
	}

	log.Println("数据库连接和迁移成功")
	return nil
}
