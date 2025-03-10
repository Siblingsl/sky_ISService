package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sky_ISService/gateway/proxy"
	"sky_ISService/gateway/router"
	"sky_ISService/gateway/swagger"
	"sky_ISService/pkg/initialize"
	"sky_ISService/pkg/middleware"
	"sky_ISService/pkg/shutdown"
	"sky_ISService/shared/cache"
	"sky_ISService/shared/elasticsearch"
	"sky_ISService/shared/mq"
	consul "sky_ISService/shared/registerservice"
	"sync"
	"syscall"
	"time"
)

// @title SKY
// @version 1.0
// @description 独立站项目

// @contact.name shilei

// @host localhost:8080
func main() {

	// 读取配置路径
	configPath := "config/config.yml"

	// 引入 Elasticsearch、Redis 和 RabbitMQ 客户端
	esClient, redisClient, rmqClient, err := initialize.InitServices(configPath)
	if err != nil {
		log.Fatalf("服务初始化失败: %v", err)
	}

	// 初始化 Consul 客户端
	consulClient, err := consul.InitConsul(configPath)
	if err != nil {
		log.Fatalf("Consul 初始化失败: %v", err)
	}

	// 使用 WaitGroup 来等待所有服务的启动和关闭
	var wg sync.WaitGroup

	// 使用 Fx 创建应用
	app := fx.New(
		// 提供 WaitGroup 到容器中，供服务启动使用
		fx.Provide(func() *sync.WaitGroup { return &wg }),

		// 提供 Proxy 代理实例
		fx.Provide(
			proxy.NewProxy,
		),
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

		// 提供 Gin 引擎
		fx.Provide(
			func(p *proxy.Proxy) *gin.Engine {
				// 创建 Gin 引擎
				r := router.NewRouter(p)

				// 注册中间件
				r.Use(middleware.CircuitMiddleware()) // 熔断中间件
				r.Use(middleware.RecoveryMiddleware())
				r.Use(middleware.ErrorHandlingMiddleware()) // 全局抓错中间件
				r.Use(middleware.JWTAuthMiddleware())       // JWT 验证中间件

				// 初始化 Swagger
				swagger.InitSwagger(r)

				// 注册服务到 Consul
				serviceName := "gateway"
				serviceID := fmt.Sprintf("%s-id", serviceName)
				address := "127.0.0.1" // 服务的 IP 地址
				port := 8080           // 服务的端口
				err := consul.RegisterServiceConsul(consulClient, serviceName, serviceID, address, port)
				if err != nil {
					log.Fatalf("服务注册失败: %v", err)
				}

				return r
			},
		),

		// 注册服务启动和关闭逻辑
		fx.Invoke(func(lc fx.Lifecycle, wg *sync.WaitGroup) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					// 启动服务
					startServiceWithWaitGroup("./services/auth/cmd/main.go", wg)
					startServiceWithWaitGroup("./services/system/cmd/main.go", wg)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					// 这里可以加更多的清理操作
					log.Println("所有服务即将停止")
					return nil
				},
			})
		}),

		// 调用 shutdown 模块的 CloseServices 函数，关闭 RabbitMQ、Redis 连接
		fx.Invoke(func(lc fx.Lifecycle, redisClient *cache.RedisClient, mqClient *mq.RabbitMQClient) {
			shutdown.CloseServices(lc, redisClient, mqClient) // 这里调用我们的关闭服务函数
		}),

		// 启动时运行的函数
		fx.Invoke(func(r *gin.Engine) {
			// 启动 Gin 引擎，监听端口
			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}

			// 启动服务
			go func() {
				if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
					log.Fatalf("服务启动失败: %v", err)
				}
			}()
		}),
	)

	// 启动 Fx 应用
	app.Run()

	// **等待所有子服务退出**
	wg.Wait()
	log.Println("所有子服务已退出")
}

// 启动服务并等待完成
func startServiceWithWaitGroup(servicePath string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startService(servicePath); err != nil {
			log.Printf("启动服务 [%s] 失败: %v", servicePath, err)
		} else {
			log.Printf("服务 [%s] 启动成功", servicePath)
		}
	}()
}

// 启动服务的通用函数
func startService(servicePath string) error {
	cmd := exec.Command("go", "run", servicePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	// 等待服务完成执行
	return cmd.Wait()
}

// 等待系统信号并优雅关闭服务
func waitForShutdown() {
	// 创建一个 channel 来接收系统的中断信号
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 等待信号
	<-signalChannel

	// 超时机制：设置超时时间避免永远等待
	shutdownTimeout := time.After(10 * time.Second)
	select {
	case <-shutdownTimeout:
		log.Println("服务关闭超时，强制退出")
	case <-signalChannel:
		log.Println("收到关闭信号，优雅退出")
	}

	// 这里可以增加更多的清理操作，比如关闭数据库连接等
	log.Println("所有子服务已优雅退出")
}
