package main

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"log"
	"os"
	configs "sky_ISService/config"
	"sky_ISService/pkg/middleware"
	loggerutils "sky_ISService/utils"

	moduleAuth "sky_ISService/services/auth/module"
)

func main() {
	// 获取服务名称，可以通过命令行参数或环境变量传入
	serviceName := "auth"             // 服务名
	configPath := "config/config.yml" // 配置文件路径

	// 使用 Fx 创建应用
	app := fx.New(
		// 提供 PostgreSQL 客户端
		fx.Provide(
			func() (*gorm.DB, error) {
				db, err := configs.InitPostgresConfig(serviceName, configPath)
				if err != nil {
					log.Fatalf("PostgreSQL 初始化失败: %v", err)
				}
				return db, nil
			},
		),
		// 提供 Elasticsearch 客户端
		fx.Provide(
			func() (*elasticsearch.Client, error) {
				elasticClient, err := configs.InitElasticsearchConfig(serviceName, configPath)
				if err != nil {
					log.Fatalf("Elasticsearch 初始化失败: %v", err)
				}
				return elasticClient, nil
			},
		),
		// 初始化日志系统
		fx.Provide(
			func(elasticClient *elasticsearch.Client) (*logrus.Logger, error) {
				logger, err := configs.InitLogger(serviceName, configPath, elasticClient)
				if err != nil {
					return nil, fmt.Errorf("日志系统初始化失败: %v", err)
				}
				return logger, nil
			},
		),

		// 提供 gin.Engine 实例到容器中
		fx.Provide(
			func(db *gorm.DB, elasticClient *elasticsearch.Client) *gin.Engine {
				r := gin.Default()
				// 使用中间件
				r.Use(middleware.DBMiddleware(db))
				r.Use(middleware.LoggerMiddleware(serviceName, elasticClient))

				return r
			},
		),

		// 注册 AuthModule
		moduleAuth.AuthModule,

		// 启动时运行的函数
		fx.Invoke(func(r *gin.Engine, logger *logrus.Logger) {
			// 启动服务
			port := os.Getenv("PORT")
			if port == "" {
				port = "8081"
			}

			// 打印初始化的日志信息
			loggerutils.LogInfo("日志系统初始化成功")
			configs.SetLogger(logger)

			// 启动 Gin 引擎
			if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
				log.Fatalf("服务启动失败: %v", err)
			}
		}),
	)

	// 启动 Fx 应用
	app.Run()
}
