package repository

import (
	"gorm.io/gorm"
	"log"
	"sky_ISService/services/system/repository/models"
)

type AdminsRepository struct {
	db                                      *gorm.DB
	*BaseRepository[models.SkySystemAdmins] // 继承 BaseRepository
}

func NewAdminsRepository(db *gorm.DB) *AdminsRepository {
	log.Println("UserRepository 实例化")
	return &AdminsRepository{
		db:             db,
		BaseRepository: NewBaseRepository[models.SkySystemAdmins](db),
	}
}

// IsUsernameExists 通过用户名检查是否存在管理员（排除已软删除的数据）
func (repo *AdminsRepository) IsUsernameExists(name string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.SkySystemAdmins{}).
		Where("username = ? AND is_deleted = false", name).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// AdminBindRoles 管理员绑定角色
func (repo *AdminsRepository) AdminBindRoles(adminID int) error {
	result := repo.db.Where("admin_id = ?", adminID).Delete(&models.AdminsRoles{})
	if result.Error != nil {
		return result.Error
	}
	// 允许删除 0 条数据，不返回错误
	return nil
}

func (repo *AdminsRepository) CreateAdminRoles(adminRoles []models.AdminsRoles) error {
	return repo.db.Create(&adminRoles).Error
}

// GetRoleIDsByAdminID 获取管理员角色信息
func (repo *AdminsRepository) GetRoleIDsByAdminID(adminID int) ([]int, error) {
	// 假设你有一个角色与管理员关联的表
	var roleIDs []int
	err := repo.db.Table("admins_roles").Where("admin_id = ?", adminID).Pluck("role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}
	return roleIDs, nil
}

// CheckAdminRoleExist 检查管理员是否已绑定某个角色
func (repo *AdminsRepository) CheckAdminRoleExist(adminID int, roleID int) (bool, error) {
	var count int64
	err := repo.db.Model(&models.AdminsRoles{}).
		Where("admin_id = ? AND role_id = ?", adminID, roleID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// RemoveAdminRole 删除管理员与角色的关联
func (repo *AdminsRepository) RemoveAdminRole(adminID int, roleID int) error {
	err := repo.db.Delete(&models.AdminsRoles{}, "admin_id = ? AND role_id = ?", adminID, roleID).Error
	return err
}

// AddAdminRole 向管理员添加角色
func (repo *AdminsRepository) AddAdminRole(adminID int, roleID int) error {
	// 创建角色绑定
	adminRole := models.AdminsRoles{
		AdminID: adminID,
		RoleID:  roleID,
	}
	// 插入数据库
	err := repo.db.Create(&adminRole).Error
	if err != nil {
		return err
	}
	return nil
}
