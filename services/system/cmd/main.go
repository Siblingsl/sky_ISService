package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"log"
	"os"
	"sky_ISService/pkg/middleware"

	initconfig "sky_ISService/config"
	loggerconfig "sky_ISService/config"
	moduleSystem "sky_ISService/services/system/module"
	loggerutils "sky_ISService/utils"
)

func main() {
	serviceName := "system"           // 服务名
	configPath := "config/config.yml" // 配置文件路径
	// 初始化配置 Postgres
	db, err := initconfig.InitPostgresConfig(serviceName, configPath)
	if err != nil {
		log.Fatalf("PostgreSQL 初始化失败: %v", err)
	}
	// 初始化配置 Elasticsearch
	elasticClient, err := initconfig.InitElasticsearchConfig(serviceName, configPath)
	if err != nil {
		log.Fatalf("Elasticsearch 初始化失败: %v", err)
	}
	// 初始化日志系统
	logger, err := loggerconfig.InitLogger(serviceName, configPath, elasticClient)
	if err != nil {
		log.Fatalf("11111无法初始化日志系统: %v", err)
	}

	// 打印初始化的日志信息
	loggerutils.LogInfo("日志系统初始化成功")
	// 设置全局 Logger
	loggerconfig.SetLogger(logger)

	// 创建 Gin 引擎实例
	r := gin.Default()

	// 使用数据库中间件
	r.Use(middleware.DBMiddleware(db))

	// 使用日志中间件
	r.Use(middleware.LoggerMiddleware(serviceName, elasticClient))

	// 使用 Fx 创建应用
	app := fx.New(
		// 提供 gin.Engine 实例到容器中
		fx.Provide(
			func() *gin.Engine {
				return r
			},
		),

		// 注册 SystemModules
		moduleSystem.SystemModules,

		// 启动时运行的函数
		fx.Invoke(func() {
			// 启动服务
			port := os.Getenv("PORT")
			if port == "" {
				port = "8082"
			}
			// 启动 Gin 引擎
			if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
				log.Fatalf("System 服务启动失败: %v", err)
			}
		}),
	)

	// 启动 Fx 应用
	app.Run()
}
