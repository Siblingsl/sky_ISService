package service

import "sky_ISService/services/system/repository"

type MenuService struct {
	menuRepository *repository.MenuRepository
}

func NewMenuService(menuRepository *repository.MenuRepository) *MenuService {
	return &MenuService{menuRepository: menuRepository}
}

func (s *MenuService) AddMenu() (string, error) {
	return "菜单", nil
}
