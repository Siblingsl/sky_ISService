package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sky_ISService/services/security/repository/models"
	"sky_ISService/shared/cache"
	"time"
)

type SecurityRepository struct {
	db          *gorm.DB
	redisClient *cache.RedisClient
}

func NewSecurityRepository(db *gorm.DB, redisClient *cache.RedisClient) *SecurityRepository {
	return &SecurityRepository{
		db:          db,
		redisClient: redisClient,
	}
}

// FindUserByUsername 通过用户名查询用户
func (repo *SecurityRepository) FindUserByUsername(username string) (*models.SkySecurityUser, error) {
	cacheKey := fmt.Sprintf("user:%s", username) // Redis 缓存键
	var user models.SkySecurityUser
	// 1. **先从 Redis 读取缓存**
	cachedData, err := repo.redisClient.Get(repo.redisClient.Ctx, cacheKey)
	if err == nil && cachedData != "" {
		// **如果缓存命中，反序列化 JSON 并返回**
		err = json.Unmarshal([]byte(cachedData), &user)
		if err == nil {
			fmt.Println("从 Redis 缓存中获取用户数据:", user)
			return &user, nil
		}
	}
	// 2. **缓存未命中，查询数据库**
	err = repo.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("用户不存在")
	} else if err != nil {
		return nil, fmt.Errorf("数据库查询出错: %v", err)
	}
	// 3. **将数据库查询结果存入 Redis**
	userData, _ := json.Marshal(user)                              // 序列化为 JSON
	err = repo.redisClient.Set(cacheKey, userData, 10*time.Minute) // 设定 10 分钟缓存
	if err != nil {
		fmt.Println("Redis 缓存存储失败:", err)
	}
	return &user, nil
}
