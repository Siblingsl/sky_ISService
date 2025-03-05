package main

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"log"
	"os"
	grpc "sky_ISService/pkg/grpc"
	"sky_ISService/pkg/middleware"
	moduleAuth "sky_ISService/services/auth/module"
	"sky_ISService/services/auth/service"
	es "sky_ISService/shared/elasticsearch"
	"sky_ISService/shared/mq"
	postgres "sky_ISService/shared/postgresql"
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
				db, err := postgres.InitPostgresConfig(serviceName, configPath)
				if err != nil {
					log.Fatalf("PostgreSQL 初始化失败: %v", err)
				}
				return db, nil
			},
		),
		// 提供 Elasticsearch 客户端
		fx.Provide(
			func() (*elasticsearch.Client, error) {
				elasticClient, err := es.InitElasticsearchConfig(configPath)
				if err != nil {
					log.Fatalf("Elasticsearch 初始化失败: %v", err)
				}
				return elasticClient, nil
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
		// 初始化 RabbitMQ 客户端
		fx.Provide(
			func() (*mq.RabbitMQClient, error) {
				rmq, err := mq.InitRabbitMQ(configPath)
				if err != nil {
					log.Fatalf("RabbitMQ 初始化失败: %v", err)
				}
				return rmq, nil
			},
		),
		// 提供 gRPC 服务器
		fx.Provide(
			func(authService *service.AuthService) *grpc.GRpcServer {
				return grpc.NewGRpcServer(authService)
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
		fx.Invoke(func(r *gin.Engine,
			//logger *logrus.Logger,
			mqClient *mq.RabbitMQClient,
		) {

			// 启动服务
			port := os.Getenv("PORT")
			if port == "" {
				port = "8081"
			}

			// 打印初始化的日志信息
			//loggerutils.LogInfo("日志系统初始化成功")
			//sharedLogger.SetLogger(logger)

			// 启动 Gin 引擎
			if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
				log.Fatalf("服务启动失败: %v", err)
			}
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
		// 启动 gRPC 服务器
		fx.Invoke(func(lc fx.Lifecycle, grpcServer *grpc.GRpcServer) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go grpcServer.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					grpcServer.Stop()
					return nil
				},
			})
		}),
	)

	// 启动 Fx 应用
	app.Run()
}
