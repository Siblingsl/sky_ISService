package service

import (
	"sky_ISService/pkg/grpc"
	"sky_ISService/services/system/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
	grpcClient     *grpc.GRpcClient
}

func NewUserService(userRepository *repository.UserRepository, grpcClient *grpc.GRpcClient) *UserService {
	return &UserService{
		userRepository: userRepository,
		grpcClient:     grpcClient,
	}
}

func (server *UserService) AddUser() (string, error) {

	// 去管理员库中获取username、pass、code
	//user, err := server.userRepository
	// 在通过grpc请求 身份验证 auth 的 token获取
	//token, err := server.grpcClient.Login()
	// 获取到就可以提供当前模块的 token 使用

	return "用户", nil
}
