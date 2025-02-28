package repository

import (
	"gorm.io/gorm"
	"log"
	"sky_ISService/services/auth/repository/models"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	log.Println("AuthRepository 实例化")
	return &AuthRepository{db: db}
}

// GetUserByID 通过用户ID查询用户
func (repo *AuthRepository) GetUserByID(userID int) (*models.SkyAuthUser, error) {
	var user models.SkyAuthUser
	if err := repo.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *AuthRepository) Demo123456(username string) (*models.SkyAuthUser, error) {
	var user models.SkyAuthUser
	if err := repo.db.Raw("SELECT * FROM sky_auth_users WHERE username = ?", username).Scan(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
