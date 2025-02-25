package repository

import "log"

type RoleRepository struct {
}

func NewRoleRepository() *RoleRepository {
	log.Println("RoleRepository 实例化")
	return &RoleRepository{}
}
