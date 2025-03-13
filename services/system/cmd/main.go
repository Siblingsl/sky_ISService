package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"log"
	"os"
	"sky_ISService/pkg/initialize"
	"sky_ISService/pkg/middleware"
	"sky_ISService/services/system/grpc"
	"sky_ISService/shared/cache"
	"sky_ISService/shared/elasticsearch"
	"sky_ISService/shared/mq"
	postgres "sky_ISService/shared/postgresql"
	consul "sky_ISService/shared/registerservice"

	moduleSystem "sky_ISService/services/system/module"
)

func main() {
	serviceName := "system"           // 服务名
	configPath := "config/config.yml" // 配置文件路径

	// 引入 Elasticsearch、Redis 和 RabbitMQ 客户端
	esClient, redisClient, rmqClient, err := initialize.InitServices(configPath)
	if err != nil {
		log.Fatalf("服务初始化失败: %v", err)
	}

	// 使用 Fx 创建应用
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
		),

		// 提供 PostgreSQL 客户端
		fx.Provide(
			func() (*gorm.DB, error) {
				db, err := postgres.InitPostgresConfig(serviceName, configPath)
				if err != nil {
					log.Fatalf("PostgreSQL 初始化失败: %v", err)
				}
				return db, nil
			},
		),
		// 初始化日志系统
		//fx.Provide(
		//	func(elasticClient *elasticsearch.Client) (*logrus.Logger, error) {
		//		fmt.Printf("Elasticsearch 客户端: %+v\n", elasticClient)
		//		logs, err := sharedLogger.InitLogger(serviceName, configPath, elasticClient)
		//		if err != nil {
		//			fmt.Printf("日志配置: %+v\n", elasticClient)
		//			return nil, fmt.Errorf("日志系统初始化失败: %v", err)
		//		}
		//		return logs, nil
		//	},
		//),

		// 提供 gin.Engine 实例到容器中
		fx.Provide(
			func(db *gorm.DB, elasticClient *elasticsearch.ElasticsearchClient) *gin.Engine {
				r := gin.Default()
				// 使用中间件
				r.Use(middleware.DBMiddleware(db))
				r.Use(middleware.LoggerMiddleware(serviceName, elasticClient))

				return r
			},
		),

		// 注册 SystemModules
		moduleSystem.SystemModules,

		// 提供 Consul 客户端
		fx.Provide(
			func() (*api.Client, error) {
				client, err := consul.InitConsul(configPath)
				if err != nil {
					log.Fatalf("Consul 初始化失败: %v", err)
				}
				return client, nil
			},
		),

		// 在服务启动时注册服务到 Consul
		fx.Invoke(func(client *api.Client) {
			serviceName := "auth"
			serviceID := fmt.Sprintf("%s-id", serviceName)
			address := "127.0.0.1" // 服务的 IP 地址
			port := 8082           // 服务端口
			// 注册服务到 Consul
			if err := consul.RegisterServiceConsul(client, serviceName, serviceID, address, port); err != nil {
				log.Fatalf("服务注册失败: %v", err)
			}
		}),

		// 启动时运行的函数
		fx.Invoke(func(r *gin.Engine,
			//logger *logrus.Logger,
			mqClient *mq.RabbitMQClient,
			lc fx.Lifecycle,
		) {
			// 启动服务 多端口监听
			port1 := os.Getenv("PORT")
			if port1 == "" {
				port1 = "8083"
			}
			port2 := os.Getenv("PORT2")
			if port2 == "" {
				port2 = "8084" // 默认使用8084端口
			}

			// 打印初始化的日志信息
			//loggerutils.LogInfo("日志系统初始化成功")
			//sharedLogger.SetLogger(logger)

			// 启动 gRPC 服务
			go func() {
				err := grpc.StartSystemGRPCServer()
				fmt.Println("gRPC 服务启动并监听 9999 端口")
				if err != nil {
					log.Fatalf("启动 gRPC 服务失败: %v", err)
				}
			}()

			// 启动 Gin 引擎
			go func() {
				if err := r.Run(fmt.Sprintf(":%s", port1)); err != nil {
					log.Fatalf("服务启动失败: %v", err)
				}
			}()

			go func() {
				if err := r.Run(fmt.Sprintf(":%s", port2)); err != nil {
					log.Fatalf("服务启动失败: %v", err)
				}
			}()

			// 阻塞等待，防止主 goroutine 退出
			select {}
		}),

		// 确保 MQ 连接在应用关闭时正确关闭
		fx.Invoke(func(lc fx.Lifecycle, client *mq.RabbitMQClient) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					log.Println("关闭 RabbitMQ 连接...")
					client.Close()
					return nil
				},
			})
		}),
		// 在应用关闭时释放 Redis 连接
		fx.Invoke(func(lc fx.Lifecycle, client *cache.RedisClient) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					log.Println("关闭 Redis 连接...")
					err := client.Close()
					if err != nil {
						return err
					}
					return nil
				},
			})
		}),
	)

	// 启动 Fx 应用
	app.Run()
}
