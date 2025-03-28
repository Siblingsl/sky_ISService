package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"log"
	"sky_ISService/config"
	"sky_ISService/pkg/initialize"
	"sky_ISService/proto/system"
	"sky_ISService/services/security/grpc"
	"sky_ISService/services/security/module"
	"sky_ISService/shared/cache"
	"sky_ISService/shared/elasticsearch"
	"sky_ISService/shared/mq"
	postgres "sky_ISService/shared/postgresql"
)

// StartServer 启动 HTTP 服务器
func StartServer(r *gin.Engine) {
	fmt.Println(fmt.Sprintf("%s:%s", config.GetConfig().Security.Host, config.GetConfig().Security.Port), "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err := r.Run(fmt.Sprintf("%s:%s", config.GetConfig().Security.Host, config.GetConfig().Security.Port)); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func main() {
	serviceName := "security" // 服务名
	// 引入 Elasticsearch、Redis 和 RabbitMQ 客户端
	esClient, redisClient, rmqClient, err := initialize.InitServices()
	if err != nil {
		log.Fatalf("服务初始化失败: %v", err)
	}
	app := fx.New(

		fx.Provide(
			// 提供 Elasticsearch 客户端
			func() *elasticsearch.ElasticsearchClient {
				return esClient
			},
			// 提供 Redis 客户端
			func() *cache.RedisClient {
				return redisClient
			},
			// 提供 Mq 客户端
			func() *mq.RabbitMQClient {
				return rmqClient
			},
			// 提供 PostgreSQL 客户端
			func() (*gorm.DB, error) {
				db, err := postgres.InitPostgresConfig(serviceName)
				if err != nil {
					fmt.Printf("postgreSQL 初始化失败: %v\n", err)
				}
				return db, nil
			},
			// 提供 gRPC 客户端
			func() (system.SystemServiceClient, error) {
				// 使用 grpc.NewSecurityToSystemClient 来创建 gRPC 客户端
				client, err := grpc.NewSecurityToSystemClient()
				if err != nil {
					log.Fatalf("无法创建 gRPC 客户端: %v", err)
				}
				return client, nil
			},
		),

		// 提供 Gin 引擎
		fx.Provide(func() *gin.Engine {
			return gin.Default()
		}),

		// 载入 SecurityModule 模块
		module.SecurityModule,

		// 调用服务器启动函数
		fx.Invoke(StartServer),
	)

	// 启动应用
	app.Run()
}
