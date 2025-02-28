package service

import (
	"sky_ISService/services/auth/repository"
	"sky_ISService/services/auth/repository/models"
	"sky_ISService/utils"
)

type AuthService struct {
	authRepository *repository.AuthRepository
}

func NewAuthService(authRepository *repository.AuthRepository) *AuthService {
	return &AuthService{authRepository: authRepository}
}

// Register 注册
func (service *AuthService) Register() (*models.SkyAuthUser, error) {
	// 注册逻辑
	user, err := service.authRepository.Demo123456("shilei")
	if err != nil {
		utils.LogError("查询失败：", err)
		return nil, err
	}
	return user, nil
}

// Login 登陆
func (service *AuthService) Login() (string, error) {
	// 模拟服务崩溃
	//go func() {
	//	// 这里模拟一个不可恢复的错误，导致服务退出
	//	log.Println("模拟服务崩溃：启动一个 goroutine 执行错误")
	//	time.Sleep(2 * time.Second)
	//	panic("服务崩溃，登录失败")
	//}()
	return "登陆", nil
}
