package shutdown

import (
	"context"
	"go.uber.org/fx"
	"log"
	"sky_ISService/shared/cache"
	"sky_ISService/shared/mq"
)

// CloseServices 用于优雅关闭服务的函数
func CloseServices(lc fx.Lifecycle, redisClient *cache.RedisClient, mqClient *mq.RabbitMQClient) {
	// 关闭 RabbitMQ 连接
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("关闭 RabbitMQ 连接...")
			if err := mqClient.Close(); err != nil {
				log.Printf("关闭 RabbitMQ 连接失败: %v", err)
				return err
			}
			return nil
		},
	})

	// 关闭 Redis 连接
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("关闭 Redis 连接...")
			if err := redisClient.Close(); err != nil {
				log.Printf("关闭 Redis 连接失败: %v", err)
				return err
			}
			return nil
		},
	})
}
