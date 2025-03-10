package service

import (
	"context"
	"sky_ISService/proto/system"
	"sky_ISService/services/system/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
	system.UnimplementedSystemServiceServer
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

// VerifyIsSystemAdmin 方法实现
func (s *UserService) VerifyIsSystemAdmin(ctx context.Context, req *system.VerifyIsSystemAdminRequest) (*system.VerifyIsSystemAdminResponse, error) {
	// 这里实现你的业务逻辑
	return &system.VerifyIsSystemAdminResponse{
		IsAdmin: true, // 示例返回值
	}, nil
}
