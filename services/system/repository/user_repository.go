package repository

import (
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	log.Println("UserRepository 实例化")
	return &UserRepository{} // ✅ 确保返回指针
}

// GetUserByID 通过用户ID查询管理员
//func (repo *UserRepository) GetUserByID(userID int) ( error) {
//	//var user models.SkyAuthUser
//	//if err := repo.db.Where("id = ?", userID).First(&user).Error; err != nil {
//	//	return nil, err
//	//}
//	//return &user, nil
//}
