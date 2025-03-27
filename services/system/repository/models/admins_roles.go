package models

import (
	"gorm.io/gorm"
	"time"
)

// AdminsRoles 代表 admins_roles 表
type AdminsRoles struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID   int       `gorm:"primaryKey" json:"admin_id"` // 管理员 ID
	RoleID    int       `gorm:"primaryKey" json:"role_id"`  // 角色 ID
	CreatedAt time.Time `json:"created_at"`                 // 创建时间
}

// BeforeCreate 创建前 Hook
func (ar *AdminsRoles) BeforeCreate(tx *gorm.DB) error {
	ar.CreatedAt = time.Now() // 创建时间设置为当前时间
	return nil
}
