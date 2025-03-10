package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"sky_ISService/config"
	"sky_ISService/proto/system"
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
	systemClient   system.SystemServiceClient
}

func NewAuthService(authRepository *repository.AuthRepository, rabbitClient *mq.RabbitMQClient, conn *grpc.ClientConn, redisClient *cache.RedisClient) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		rabbitClient:   rabbitClient,
		redisClient:    redisClient,
		systemClient:   system.NewSystemServiceClient(conn),
	}
}

// AdminLogin 管理员登录
func (s *AuthService) AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (string, error) {
	// 1. 查找用户
	user, err := s.authRepository.FindUserByUsername(req.Username)
	if err != nil {
		return "", fmt.Errorf("未找到用户")
	}
	// 2. 验证密码
	if user.Password != req.Password {
		return "", fmt.Errorf("密码错误")
	}
	// 3. 校验邮箱验证码
	// 从 Redis 获取存储的验证码
	storedCode, err := s.redisClient.Get(req.Email)
	if err != nil {
		return "", fmt.Errorf("验证码验证失败")
	}
	// 验证用户提供的验证码是否正确
	if storedCode != req.Code {
		return "", fmt.Errorf("验证码错误或已过期")
	}
	// 4. 调用 system 微服务验证是否为管理员
	reqVerify := &system.VerifyIsSystemAdminRequest{
		UserId:   strconv.Itoa(int(user.ID)),
		UserName: user.Username,
	}
	resp, err := s.systemClient.VerifyIsSystemAdmin(ctx, reqVerify)
	if err != nil {
		return "", fmt.Errorf("验证管理员身份失败: %v", err)
	}
	// 5. 判断是否为管理员
	if !resp.IsAdmin {
		return "", fmt.Errorf("该用户不是管理员")
	}
	// 6. 生成JWT Token
	token, err := utils.GenerateToken(strconv.Itoa(int(user.ID)), user.Username)
	if err != nil {
		return "", fmt.Errorf("生成 Token 失败")
	}
	return token, nil
}

// SendEmailCode 发送邮箱验证码并缓存
func (s *AuthService) SendEmailCode(ctx *gin.Context, email string) error {
	// 1. 校验邮箱格式是否合法
	if !utils.IsValidEmail(email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	// 2. 频率限制 - 1 分钟内同一邮箱只能请求 1 次
	cacheKey := "email_limit:" + email
	if exists, _ := s.redisClient.Client.Exists(ctx, cacheKey).Result(); exists == 1 {
		return fmt.Errorf("请勿频繁请求验证码，稍后再试")
	}

	// 3. 统计 1 小时内请求次数，防止滥用
	hourlyLimitKey := "email_hourly:" + email
	requestCount, _ := s.redisClient.Client.Get(ctx, hourlyLimitKey).Int()
	if requestCount >= 50 {
		// 进入黑名单
		blacklistKey := "blacklist:" + email
		s.redisClient.Client.Set(ctx, blacklistKey, "1", 24*time.Hour)
		return fmt.Errorf("您的邮箱请求过于频繁，请明天再试")
	}

	// 4. 限制同一 IP 1 小时内请求不同邮箱的数量，防止爬虫攻击
	ipLimitKey := "ip_limit:" + utils.GetClientIP(ctx)
	ipCount, _ := s.redisClient.Client.Get(ctx, ipLimitKey).Int()
	if ipCount >= 100 { // 限制同一 IP 1 小时请求 100 个不同邮箱
		return fmt.Errorf("您的 IP 请求过于频繁，请稍后再试")
	}

	// 5. 生成 6 位随机验证码（存储时做加密处理）
	code := utils.GenerateRandomCode(6)
	secretKey := config.GetConfig().AESSecret.Secret
	encryptedCode, _ := utils.EncryptAES(code, secretKey) // 加密存储

	// 6. 存储验证码到 Redis，5 分钟有效
	err := s.redisClient.Set(email, encryptedCode, 1*time.Minute)
	if err != nil {
		return fmt.Errorf("存储验证码到 Redis 失败: %v", err)
	}

	// 7. 更新请求次数
	s.redisClient.Client.Incr(ctx, hourlyLimitKey)                // 邮箱请求次数+1
	s.redisClient.Client.Expire(ctx, hourlyLimitKey, 1*time.Hour) // 1 小时过期
	s.redisClient.Client.Incr(ctx, ipLimitKey)                    // IP 计数 +1
	s.redisClient.Client.Expire(ctx, ipLimitKey, 1*time.Hour)     // 1 小时过期

	// 8. 设置邮箱请求限制，1 分钟只能请求 1 次
	s.redisClient.Set(cacheKey, "1", 1*time.Minute)

	// 9. 发送验证码邮件（异步处理，防止阻塞）
	go func() {
		if err := utils.SendEmail(email, code); err != nil {
			fmt.Println("发送验证码失败:", err)
		}
	}()

	return nil
}

// Register 注册
//func (service *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
//	// 假设这里处理注册逻辑，可能会调用 authRepository 来创建新用户
//	if req.Username == "" || req.Password == "" {
//		return &pb.RegisterResponse{
//			Message: "用户名或密码不能为空",
//			Status:  2,
//		}, nil
//	}
//
//	// 注册逻辑
//	user, err := service.authRepository.Demo123456("shilei")
//	if err != nil {
//		utils.LogError("查询失败：", err)
//		return &pb.RegisterResponse{
//			Message: err.Error(),
//			Status:  2,
//		}, nil
//	}
//
//	err = service.redisClient.Set("user:"+user.Username, user.Username+"|"+user.Password, time.Minute*10)
//	if err != nil {
//		utils.LogError("Redis 存储失败:", err)
//	} else {
//		utils.LogInfo("用户信息已存入 Redis")
//	}
//
//	result := user.Username + user.Password + "shilei"
//	err = service.rabbitClient.SendMessage("Auth_Msg_queue", result)
//	if err != nil {
//		utils.LogError("成功发送用户注册消息:", err)
//	} else {
//		utils.LogInfo("发送用户注册消息失败: %v")
//	}
//
//	//return user, nil
//	return &pb.RegisterResponse{
//		Message: "注册成功",
//		Status:  2,
//	}, nil
//}

// Login 登陆
//func (service *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
//
//	if req.Username == "shilei" && req.Password == "123456" {
//		return &pb.LoginResponse{
//			Token:  "153424512",
//			Status: 2,
//		}, nil
//	}
//	// 如果用户名或密码错误，可以返回相应的错误
//	return nil, fmt.Errorf("invalid credentials")
//}
