package service

import (
	"context"
	"sky_ISService/services/auth/repository"
	"sky_ISService/shared/cache"
	"sky_ISService/shared/mq"
	"time"
)

type OrderService struct {
	authRepository *repository.AuthRepository
	rabbitClient   *mq.RabbitMQClient
	redisClient    *cache.RedisClient
	//grpcClient     system.SystemServiceClient
}

func NewOrderService(authRepository *repository.AuthRepository, rabbitClient *mq.RabbitMQClient, redisClient *cache.RedisClient) *OrderService {
	// 初始化 gRPC 客户端
	//grpcClient, _ := grpc.NewSystemClient()
	return &OrderService{
		authRepository: authRepository,
		rabbitClient:   rabbitClient,
		redisClient:    redisClient,
		//grpcClient:     grpcClient,
	}
}

func (s *OrderService) Testxxxx(ctx context.Context, email string) (string, error) {

	startTime := time.Now()

	return startTime.String() + email, nil
}
