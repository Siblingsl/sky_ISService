package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sky_ISService/config"
	"time"
)

// InitPostgres 初始化 PostgresSQL 客户端
func initPostgres(serviceName string) (*gorm.DB, error) {
	configSql, err := config.InitLoadConfig()
	if err != nil {
		return nil, err
	}

	// 获取指定服务的 PostgresSQL 配置
	serviceConfig, exists := configSql.Database[serviceName]
	if !exists {
		return nil, fmt.Errorf("服务配置不存在: %v", serviceName)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		serviceConfig.PostgreSQL.Host, serviceConfig.PostgreSQL.Username, serviceConfig.PostgreSQL.Password, serviceConfig.PostgreSQL.Database, serviceConfig.PostgreSQL.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(25)                 // 最大打开连接数
	sqlDB.SetMaxIdleConns(25)                 // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(5 * time.Second) // 连接最大存活时间
	if err != nil {
		return nil, fmt.Errorf("无法连接到 PostgreSQL: %v", err)
	}

	fmt.Println("成功连接到 PostgreSQL")
	return db, nil
}

// InitPostgresConfig 初始化 PostgresSQL
func InitPostgresConfig(serviceName string) (*gorm.DB, error) {
	return initPostgres(serviceName)
}
