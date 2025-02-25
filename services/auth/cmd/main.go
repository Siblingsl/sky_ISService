package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"go.uber.org/fx"
	"log"
	"os"
	"sky_ISService/pkg/middleware"
	moduleAuth "sky_ISService/services/auth/module"
	logger "sky_ISService/utils"
)

// 初始化 Elasticsearch 客户端
func initElasticsearch(serviceName, configPath string) (*elastic.Client, error) {
	// 加载配置
	config, err := logger.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法加载配置: %v", err)
	}

	// 获取指定服务的 Elasticsearch 配置
	serviceConfig, exists := config.Services[serviceName]
	if !exists {
		return nil, fmt.Errorf("服务配置不存在: %v", serviceName)
	}

	// 创建 Elasticsearch 客户端
	client, err := elastic.NewClient(
		elastic.SetURL(serviceConfig.Elasticsearch.Host),                                             // Elasticsearch 地址
		elastic.SetBasicAuth(serviceConfig.Elasticsearch.User, serviceConfig.Elasticsearch.Password), // 用户和密码
	)
	if err != nil {
		return nil, fmt.Errorf("无法创建 Elasticsearch 客户端: %v", err)
	}
	return client, nil
}

func main() {
	// 从配置文件中加载 Elasticsearch 配置
	client, err := initElasticsearch("auth", "config/config.yml")
	if err != nil {
		log.Fatalf("Elasticsearch 初始化失败: %v", err)
	}

	// 初始化日志系统，传递客户端
	err = logger.InitLogger(client, "config/config.yml")
	if err != nil {
		log.Fatalf("日志系统初始化失败: %v", err)
	}

	// 创建 Gin 引擎实例
	r := gin.Default()

	// 使用日志中间件
	r.Use(middleware.LoggerMiddleware(client))

	// 使用 Fx 创建应用
	app := fx.New(
		// 提供 gin.Engine 实例到容器中
		fx.Provide(
			func() *gin.Engine {
				return r
			},
		),

		// 注册 AuthModule
		moduleAuth.AuthModule,

		// 启动时运行的函数
		fx.Invoke(func() {
			// 启动服务
			port := os.Getenv("PORT")
			if port == "" {
				port = "8081"
			}
			// 启动 Gin 引擎
			if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
				log.Fatalf("服务启动失败: %v", err)
			}
		}),
	)

	// 启动 Fx 应用
	app.Run()
}
