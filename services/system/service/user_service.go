package service

import "sky_ISService/services/system/repository"

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) AddUser() (string, error) {
	return "用户", nil
}
