package repository

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"sky_ISService/services/system/repository/models"
	"sky_ISService/utils"
	"sky_ISService/utils/database"
)

type AdminsRepository struct {
	db *gorm.DB
}

func NewAdminsRepository(db *gorm.DB) *AdminsRepository {
	log.Println("UserRepository 实例化")
	return &AdminsRepository{db: db} // ✅ 确保返回指针
}

// IsUsernameExists 通过用户名检查是否存在管理员
func (repo *AdminsRepository) IsUsernameExists(name string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.SkySystemAdmins{}).
		Where("username = ?", name).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// PostCreateAdmin 通过用户ID查询管理员
func (repo *AdminsRepository) PostCreateAdmin(admin *models.SkySystemAdmins) error {
	if err := repo.db.Create(admin).Error; err != nil {
		return err
	}
	return nil
}

// GetAdminByID 查询单个管理员用户
func (repo *AdminsRepository) GetAdminByID(id int) (*models.SkySystemAdmins, error) {
	var admin models.SkySystemAdmins
	if err := repo.db.Where("id = ?", id).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回 nil 表示未找到
		}
		return nil, err
	}
	return &admin, nil
}

func (r *AdminsRepository) GetUsersWithPagination(ctx *gin.Context, page int, size int, conditions map[string]interface{}) (*utils.Pagination, error) {
	// 初始化数据库查询
	db := r.db.Debug().Model(&models.SkySystemAdmins{})

	// 应用动态查询条件
	db = database.ApplyConditions(db, conditions)

	// 执行查询并获取符合条件的数据总数
	var total int64
	var users []models.SkySystemAdmins

	if len(conditions) > 0 {
		// 查询数据总数
		if err := db.Find(&users).Count(&total).Error; err != nil {
			return nil, err
		}
	} else {
		// 执行分页查询
		if err := db.Offset((page - 1) * size).Limit(size).Find(&users).Error; err != nil {
			return nil, err
		}
	}

	// 计算分页信息
	pagination := &utils.Pagination{
		Total: total,
		Data:  users,
	}

	return pagination, nil
}
