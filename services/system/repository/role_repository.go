package repository

import (
	"gorm.io/gorm"
	"log"
	"sky_ISService/services/system/repository/models"
)

type RoleRepository struct {
	db *gorm.DB
	*BaseRepository[models.SkySystemRoles]
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	log.Println("RoleRepository 实例化")
	return &RoleRepository{
		db:             db,
		BaseRepository: NewBaseRepository[models.SkySystemRoles](db),
	}
}

// IsRoleNameExists 通过角色检查是否存在角色（排除已软删除的数据）
func (repo *RoleRepository) IsRoleNameExists(name string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.SkySystemAdmins{}).
		Where("role_name = ? AND is_deleted = false", name).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// RemoveOldMenusFromRole 删除角色的所有菜单权限
func (repo *RoleRepository) RemoveOldMenusFromRole(roleID int) error {
	// 删除角色现有的所有菜单权限
	if err := repo.db.Where("role_id = ?", roleID).Delete(&models.RolesMenus{}).Error; err != nil {
		return err
	}
	return nil
}

// AssignMenusToRole 批量为角色分配菜单
func (repo *RoleRepository) AssignMenusToRole(roleID int, menuIDs []int) error {
	// 创建菜单权限的切片
	var rolesMenus []models.RolesMenus
	for _, menuID := range menuIDs {
		rolesMenus = append(rolesMenus, models.RolesMenus{
			RoleID: roleID,
			MenuID: menuID,
		})
	}

	// 批量插入菜单权限记录
	if err := repo.db.Create(&rolesMenus).Error; err != nil {
		return err
	}
	return nil
}
