package initialize

import (
	"log"
	"sky_ISService/shared/cache"
	"sky_ISService/shared/elasticsearch"
	"sky_ISService/shared/mq"
)

// InitServices 初始化 Elasticsearch、Redis 和 RabbitMQ 客户端
// @return *elasticsearch.EsClient: 返回 Elasticsearch 客户端实例
// @return *cache.RedisClient: 返回 Redis 客户端实例
// @return *mq.RabbitMQClient: 返回 RabbitMQ 客户端实例
// @return error: 如果初始化过程中有任何错误，返回错误信息
func InitServices() (*elasticsearch.ElasticsearchClient, *cache.RedisClient, *mq.RabbitMQClient, error) {
	// Elasticsearch 初始化
	esClient, err := elasticsearch.InitElasticsearchConfig()
	if err != nil {
		log.Fatalf("Elasticsearch 初始化失败: %v", err)
		return nil, nil, nil, err
	}

	// Redis 初始化
	redisClient, err := cache.InitRedisConfig()
	if err != nil {
		log.Fatalf("Redis 初始化失败: %v", err)
		return nil, nil, nil, err
	}

	// RabbitMQ 初始化
	rmqClient, err := mq.InitRabbitMQ()
	if err != nil {
		log.Fatalf("RabbitMQ 初始化失败: %v", err)
		return nil, nil, nil, err
	}

	return esClient, redisClient, rmqClient, nil
}
