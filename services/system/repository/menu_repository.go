package repository

import (
	"gorm.io/gorm"
	"log"
	"sky_ISService/services/system/repository/models"
)

type MenuRepository struct {
	db                                     *gorm.DB
	*BaseRepository[models.SkySystemMenus] // 继承 BaseRepository
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	log.Println("MenuRepository 实例化")
	return &MenuRepository{
		db:             db,
		BaseRepository: NewBaseRepository[models.SkySystemMenus](db),
	}
}

// IsMenuExist 查寻菜单是否存在
func (repo *MenuRepository) IsMenuExist(name string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.SkySystemMenus{}).
		Where("menu_name = ? AND is_deleted = false", name).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FetchAllMenus 查询所有菜单
func (repo *MenuRepository) FetchAllMenus() ([]models.SkySystemMenus, error) {
	var menus []models.SkySystemMenus

	// 查询所有未删除的菜单，并按 parent_id 和 menu_sort 排序
	err := repo.db.Where("is_deleted = ?", false).Order("parent_id, menu_sort").Find(&menus).Error
	if err != nil {
		return nil, err // 返回空切片和错误
	}

	// 返回查询到的菜单数据和 nil 错误
	return menus, nil
}

// FetchMenusByRole 获取角色所拥有的菜单
func (repo *MenuRepository) FetchMenusByRole(roleID int, roleMenus *[]models.RolesMenus) error {
	// 查询角色与菜单的关系
	if err := repo.db.Where("role_id = ?", roleID).Find(roleMenus).Error; err != nil {
		return err
	}
	return nil
}
