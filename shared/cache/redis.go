package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sky_ISService/config"
	"sync"
	"time"
)

// RedisClient 封装 Redis 客户端
type RedisClient struct {
	Client *redis.Client   // Redis 客户端实例
	Ctx    context.Context // 上下文对象，用于管理 Redis 操作的生命周期
}

// 单例模式: 只有一个 RedisClient 实例
var (
	once     sync.Once
	instance *RedisClient
)

// initRedis 初始化 Redis 客户端
// @param configPath string: 配置文件路径
// @return *RedisClient: Redis 客户端实例
// @return error: 初始化过程中可能出现的错误
func initRedis() (*RedisClient, error) {
	// 加载配置文件
	configRedis, err := config.InitLoadConfig()
	if err != nil {
		return nil, err
	}

	// 获取指定服务的 Redis 配置
	cacheRedisConfig := configRedis.Cache.Redis

	// 构建 Redis 连接参数
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cacheRedisConfig.Host, cacheRedisConfig.Port),
		Password: cacheRedisConfig.Password, // 如果没有密码可以不填
		DB:       cacheRedisConfig.DB,       // 默认数据库
	})

	// 测试 Redis 连接
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("无法连接到 Redis: %v", err)
	}

	fmt.Println("成功连接到 Redis")

	return &RedisClient{
		Client: client,
		Ctx:    context.Background(),
	}, nil
}

// InitRedisConfig 获取 Redis 单例
// @return *RedisClient: Redis 客户端实例
// @return error: 可能的错误，如果 Redis 无法初始化
func InitRedisConfig() (*RedisClient, error) {
	// 保证 Redis 客户端只被初始化一次
	once.Do(func() {
		var err error
		instance, err = initRedis()
		if err != nil {
			fmt.Printf("初始化 Redis 失败: %v", err)
		}
	})
	if instance == nil {
		return nil, fmt.Errorf("redis 实例为空")
	}
	return instance, nil
}

// Set 设置 Redis 缓存
// @param key string: 要存储的键
// @param value interface{}: 要存储的值（可以是字符串、整数、JSON 序列化后的数据等）
// @param expiration time.Duration: 过期时间（如果为 0，表示永不过期）
// @return error: 如果存储失败，则返回错误
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(r.Ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis Set 操作失败: %v", err)
	}
	return nil
}

// Get 获取 Redis 缓存的值
// @param key string: 要获取的键
// @return string: 获取到的值，如果 key 不存在则返回空字符串
// @return error: 如果查询失败，则返回错误
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil // key 不存在
		}
		return "", fmt.Errorf("redis Get 操作失败: %v", err)
	}
	return value, nil
}

// Close 关闭 Redis 连接
// @return error: 如果关闭失败，则返回错误
func (r *RedisClient) Close() error {
	err := r.Client.Close()
	if err != nil {
		return fmt.Errorf("redis 关闭连接失败: %v", err)
	}
	fmt.Println("Redis 连接已关闭")
	return nil
}
