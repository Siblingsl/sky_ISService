package repository

import "log"

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	log.Println("AuthRepository 实例化")
	return &AuthRepository{}
}
