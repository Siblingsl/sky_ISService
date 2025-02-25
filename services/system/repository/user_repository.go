package repository

import "log"

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	log.Println("UserRepository 实例化")
	return &UserRepository{} // ✅ 确保返回指针
}
