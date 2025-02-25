package service

import (
	"sky_ISService/services/auth/repository"
)

type AuthService struct {
	authRepository *repository.AuthRepository
}

func NewAuthService(authRepository *repository.AuthRepository) *AuthService {
	return &AuthService{authRepository: authRepository}
}

// Register 注册
func (s *AuthService) Register() (string, error) {
	// 注册逻辑
	return "注册成功", nil
}

// Login 登陆
func (s *AuthService) Login() (string, error) {
	// 模拟服务崩溃
	//go func() {
	//	// 这里模拟一个不可恢复的错误，导致服务退出
	//	log.Println("模拟服务崩溃：启动一个 goroutine 执行错误")
	//	time.Sleep(2 * time.Second)
	//	panic("服务崩溃，登录失败")
	//}()
	return "登陆", nil
}
