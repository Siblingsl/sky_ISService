package service

import "sky_ISService/services/system/repository"

type RoleService struct {
	roleRepository *repository.RoleRepository
}

func NewRoleService(roleRepository *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepository: roleRepository}
}

func (c *RoleService) AddRole() (string, error) {
	return "角色", nil
}
