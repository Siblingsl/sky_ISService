package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"sky_ISService/services/auth/dto"
	"sky_ISService/services/auth/repository"
	"sky_ISService/shared/cache"
	"sky_ISService/shared/mq"
	"sky_ISService/utils"
	"strconv"
	"time"
)

type AuthService struct {
	authRepository *repository.AuthRepository
	rabbitClient   *mq.RabbitMQClient
	redisClient    *cache.RedisClient
	//grpcClient     system.SystemServiceClient
}

func NewAuthService(authRepository *repository.AuthRepository, rabbitClient *mq.RabbitMQClient, redisClient *cache.RedisClient) *AuthService {
	// 初始化 gRPC 客户端
	//grpcClient, _ := grpc.NewSystemClient()
	return &AuthService{
		authRepository: authRepository,
		rabbitClient:   rabbitClient,
		redisClient:    redisClient,
		//grpcClient:     grpcClient,
	}
}

// AdminLoginToken 管理员登陆
func (s *AuthService) AdminLoginToken(ctx context.Context, req dto.AdminLoginRequest) (string, error) {
	// 记录时间
	startTime := time.Now()
	fmt.Println("AdminLogin 执行完成，耗时:", time.Since(startTime))
	fmt.Println(req, "reqreqreqreqreqreqreq")
	token := "12346545222"

	return token, nil
}

//func (s *AuthService) AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (string, error) {
//	startTime := time.Now() // 记录开始时间
//
//	dbCtx, dbCancel := context.WithTimeout(ctx, 2*time.Second)
//	defer dbCancel()
//	user, err := s.authRepository.FindUserByUsername(dbCtx, req.Username)
//	if err != nil {
//		return "", fmt.Errorf("未找到用户")
//	}
//
//	if user.Password != req.Password {
//		return "", fmt.Errorf("密码错误")
//	}
//
//	redisChan := make(chan string, 1)
//	grpcChan := make(chan bool, 1)
//
//	redisCtx, redisCancel := context.WithTimeout(ctx, 1*time.Second)
//	defer redisCancel()
//	go func() {
//		storedCode, err := s.redisClient.Get(redisCtx, req.Email)
//		if err != nil {
//			redisChan <- ""
//		} else {
//			redisChan <- storedCode
//		}
//	}()
//
//	grpcCtx, grpcCancel := context.WithTimeout(ctx, 2*time.Second)
//	defer grpcCancel()
//	go func() {
//		reqVerify := &system.VerifyIsSystemAdminRequest{
//			UserId:   strconv.Itoa(int(user.ID)),
//			UserName: user.Username,
//		}
//		resp, err := s.grpcClient.VerifyIsSystemAdmin(grpcCtx, reqVerify)
//		if err != nil {
//			grpcChan <- false
//		} else {
//			grpcChan <- resp.IsAdmin
//		}
//	}()
//
//	var storedCode string
//	var isAdmin bool
//	for i := 0; i < 2; i++ {
//		select {
//		case storedCode = <-redisChan:
//		case isAdmin = <-grpcChan:
//		case <-time.After(2 * time.Second):
//			return "", fmt.Errorf("超时错误")
//		}
//	}
//
//	if storedCode != req.Code {
//		return "", fmt.Errorf("验证码错误或已过期")
//	}
//	if !isAdmin {
//		return "", fmt.Errorf("该用户不是管理员")
//	}
//
//	token, err := utils.GenerateToken(strconv.Itoa(int(user.ID)), user.Username)
//	if err != nil {
//		return "", fmt.Errorf("生成 Token 失败")
//	}
//
//	fmt.Println("AdminLogin 耗时:", time.Since(startTime))
//	return token, nil
//}

// SendEmailCode 发送邮箱验证码
//func (s *AuthService) SendEmailCode(ctx *gin.Context, email string) error {
//	// 记录时间
//	startTime := time.Now()
//	fmt.Println("SendEmailCode 执行完成，耗时:", time.Since(startTime))
//	return nil
//}

func (s *AuthService) SendEmailCode(ctx *gin.Context, email string) error {
	startTime := time.Now()

	// 1. 校验邮箱格式是否合法
	if !utils.IsValidEmail(email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	// 2. 频率限制 - 1 分钟内同一邮箱只能请求 1 次
	cacheKey := "email_limit:" + email
	if exists, _ := s.redisClient.Client.Exists(ctx, cacheKey).Result(); exists == 1 {
		return fmt.Errorf("请勿频繁请求验证码，稍后再试")
	}

	// 3. **减少 Redis 查询次数，优化请求计数**
	pipe := s.redisClient.Client.Pipeline()
	hourlyLimitKey := "email_hourly:" + email
	ipLimitKey := "ip_limit:" + utils.GetClientIP(ctx)
	blacklistKey := "blacklist:" + email

	// **批量获取 Redis 结果，提高查询效率**
	results, err := pipe.MGet(ctx, hourlyLimitKey, ipLimitKey, blacklistKey).Result()
	if err != nil {
		return fmt.Errorf("Redis 查询失败: %v", err)
	}

	// **防止数组越界**
	if len(results) < 3 {
		return fmt.Errorf("Redis 查询结果异常")
	}

	// **黑名单检查**
	if results[2] != nil {
		return fmt.Errorf("您的邮箱请求过于频繁，请明天再试")
	}

	// **请求次数限制**
	if count, _ := strconv.Atoi(fmt.Sprint(results[0])); count >= 50 {
		// 进入黑名单
		pipe.Set(ctx, blacklistKey, "1", 24*time.Hour)
		return fmt.Errorf("您的邮箱请求过于频繁，请明天再试")
	}

	// **IP 限制**
	if ipCount, _ := strconv.Atoi(fmt.Sprint(results[1])); ipCount >= 100 {
		return fmt.Errorf("您的 IP 请求过于频繁，请稍后再试")
	}

	// 4. **减少 Redis 连接数，优化验证码存储**
	code := utils.GenerateRandomCode(6)
	pipe.Set(ctx, email, code, 30*time.Minute)
	pipe.Incr(ctx, hourlyLimitKey)
	pipe.Expire(ctx, hourlyLimitKey, 1*time.Hour)
	pipe.Incr(ctx, ipLimitKey)
	pipe.Expire(ctx, ipLimitKey, 1*time.Hour)
	pipe.Set(ctx, cacheKey, "1", 1*time.Minute)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("Redis 存储验证码失败: %v", err)
	}

	// **异步发送邮件**
	go func() {
		if err := utils.SendEmail(email, code); err != nil {
			fmt.Println("发送验证码失败:", err)
		}
	}()

	fmt.Println("SendEmailCode 耗时:", time.Since(startTime)) // 打印耗时
	return nil
}

func (s *AuthService) Testxxxx(ctx context.Context, email string) (string, error) {

	startTime := time.Now()

	return startTime.String() + email, nil
}
