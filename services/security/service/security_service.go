package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"sky_ISService/proto/system"
	"sky_ISService/services/security/dto"
	"sky_ISService/services/security/repository"
	"sky_ISService/shared/cache"
	"sky_ISService/utils"
	"strconv"
	"time"
)

type SecurityService struct {
	securityRepository *repository.SecurityRepository
	redisClient        *cache.RedisClient
	grpcClient         system.SystemServiceClient
}

func NewSecurityService(securityRepository *repository.SecurityRepository, redisClient *cache.RedisClient, grpcClient system.SystemServiceClient) *SecurityService {
	return &SecurityService{
		securityRepository: securityRepository,
		redisClient:        redisClient,
		grpcClient:         grpcClient,
	}
}

func (s *SecurityService) AdminLogin(ctx context.Context, req dto.SecurityAdminLoginRequest) (string, error) {
	user, err := s.securityRepository.FindUserByUsername(req.Username)
	if err != nil {
		return "", fmt.Errorf("未找到用户")
	}
	if user.Password != req.Password {
		return "", fmt.Errorf("密码错误")
	}
	redisChan := make(chan string, 1)
	grpcChan := make(chan bool, 1)

	redisCtx, redisCancel := context.WithTimeout(ctx, 1*time.Second)
	defer redisCancel()
	go func() {
		storedCode, err := s.redisClient.Get(redisCtx, req.Email)
		if err != nil {
			redisChan <- ""
		} else {
			redisChan <- storedCode
		}
	}()

	grpcCtx, grpcCancel := context.WithTimeout(ctx, 2*time.Second)
	defer grpcCancel()
	go func() {
		reqVerify := &system.VerifyIsSystemAdminRequest{
			UserId:   strconv.Itoa(int(user.ID)),
			UserName: user.Username,
		}
		resp, err := s.grpcClient.VerifyIsSystemAdmin(grpcCtx, reqVerify)
		if err != nil {
			grpcChan <- false
		} else {
			grpcChan <- resp.IsAdmin
		}
	}()

	var storedCode string
	var isAdmin bool
	for i := 0; i < 2; i++ {
		select {
		case storedCode = <-redisChan:
		case isAdmin = <-grpcChan:
		case <-time.After(2 * time.Second):
			return "", fmt.Errorf("超时错误")
		}
	}

	if storedCode != req.Code {
		return "", fmt.Errorf("验证码错误或已过期")
	}
	if !isAdmin {
		return "", fmt.Errorf("该用户不是管理员")
	}

	token, err := utils.GenerateToken(strconv.Itoa(int(user.ID)), user.Username)
	if err != nil {
		return "", fmt.Errorf("生成 Token 失败")
	}

	return token, nil
}

// SendEmailCode 发送邮箱验证码
func (s *SecurityService) SendEmailCode(ctx *gin.Context, email string) error {
	startTime := time.Now()

	// 1. 校验邮箱格式是否合法
	if !utils.IsValidEmail(email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	// 2. 频率限制 - 1 分钟内同一邮箱只能请求 1 次
	cacheKey := "email_limit:" + email
	exists, _ := s.redisClient.Client.Exists(ctx, cacheKey).Result()
	if exists == 1 {
		return fmt.Errorf("请勿频繁请求验证码，稍后再试")
	}

	// 3. **单独查询 Redis 中的每个 key**
	hourlyLimitKey := "email_hourly:" + email
	ipLimitKey := "ip_limit:" + utils.GetClientIP(ctx)
	blacklistKey := "blacklist:" + email

	// **黑名单检查**
	blacklistExists, _ := s.redisClient.Client.Exists(ctx, blacklistKey).Result()
	if blacklistExists == 1 {
		return fmt.Errorf("您的邮箱请求过于频繁，请明天再试")
	}

	// **请求次数限制**
	hourlyCountStr, _ := s.redisClient.Client.Get(ctx, hourlyLimitKey).Result()
	hourlyCount, _ := strconv.Atoi(hourlyCountStr)
	if hourlyCount >= 50 {
		// 进入黑名单
		s.redisClient.Client.Set(ctx, blacklistKey, "1", 24*time.Hour)
		return fmt.Errorf("您的邮箱请求过于频繁，请明天再试")
	}

	// **IP 限制**
	ipCountStr, _ := s.redisClient.Client.Get(ctx, ipLimitKey).Result()
	ipCount, _ := strconv.Atoi(ipCountStr)
	if ipCount >= 100 {
		return fmt.Errorf("您的 IP 请求过于频繁，请稍后再试")
	}

	// 4. **存储验证码**
	code := utils.GenerateRandomCode(6)
	s.redisClient.Client.Set(ctx, email, code, 30*time.Minute)
	s.redisClient.Client.Incr(ctx, hourlyLimitKey)
	s.redisClient.Client.Expire(ctx, hourlyLimitKey, 1*time.Hour)
	s.redisClient.Client.Incr(ctx, ipLimitKey)
	s.redisClient.Client.Expire(ctx, ipLimitKey, 1*time.Hour)
	s.redisClient.Client.Set(ctx, cacheKey, "1", 1*time.Minute)

	// **异步发送邮件**
	go func() {
		if err := utils.SendEmail(email, code); err != nil {
			fmt.Println("发送验证码失败:", err)
		}
	}()

	fmt.Println("SendEmailCode 耗时:", time.Since(startTime)) // 打印耗时
	return nil
}

func (s *SecurityService) Testxxxx(ctx context.Context, email string) (string, error) {

	startTime := time.Now()

	return startTime.String() + email, nil
}
