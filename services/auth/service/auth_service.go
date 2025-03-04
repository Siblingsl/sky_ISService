package service

import (
	"context"
	"fmt"
	pb "sky_ISService/proto/auth"
	"sky_ISService/services/auth/repository"
	"sky_ISService/shared/mq"
	"sky_ISService/utils"
)

type AuthService struct {
	authRepository                    *repository.AuthRepository
	rabbitClient                      *mq.RabbitMQClient
	pb.UnimplementedAuthServiceServer // 默认实现
}

func NewAuthService(authRepository *repository.AuthRepository, rabbitClient *mq.RabbitMQClient) *AuthService {
	return &AuthService{authRepository: authRepository}
}

// Register 注册
func (service *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 假设这里处理注册逻辑，可能会调用 authRepository 来创建新用户
	if req.Username == "" || req.Password == "" {
		return &pb.RegisterResponse{
			Message: "用户名或密码不能为空",
			Status:  2,
		}, nil
	}

	// 注册逻辑
	user, err := service.authRepository.Demo123456("shilei")
	if err != nil {
		utils.LogError("查询失败：", err)
		return &pb.RegisterResponse{
			Message: err.Error(),
			Status:  2,
		}, nil
	}
	result := user.Username + user.Password + "shilei"
	err = service.rabbitClient.SendMessage("Auth_Msg_queue", result)
	if err != nil {
		utils.LogError("成功发送用户注册消息:", err)
	} else {
		utils.LogInfo("发送用户注册消息失败: %v")
	}

	//return user, nil
	return &pb.RegisterResponse{
		Message: "注册成功",
		Status:  2,
	}, nil
}

// Login 登陆
func (service *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	if req.Username == "shilei" && req.Password == "123456" {
		return &pb.LoginResponse{
			Token:  "153424512",
			Status: 2,
		}, nil
	}
	// 如果用户名或密码错误，可以返回相应的错误
	return nil, fmt.Errorf("invalid credentials")
}
