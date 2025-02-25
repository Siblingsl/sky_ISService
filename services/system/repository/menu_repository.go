package repository

import "log"

type MenuRepository struct{}

func NewMenuRepository() *MenuRepository {
	log.Println("MenuRepository 实例化")
	return &MenuRepository{}
}
