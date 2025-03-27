package models

import (
	"gorm.io/gorm"
	"time"
)

// RolesMenus 代表 role_menu 表
type RolesMenus struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID    int       `gorm:"primaryKey" json:"role_id"` // 角色 ID
	MenuID    int       `gorm:"primaryKey" json:"menu_id"` // 菜单 ID
	CreatedAt time.Time `json:"created_at"`                // 创建时间
}

// BeforeCreate 创建前 Hook
func (rm *RolesMenus) BeforeCreate(tx *gorm.DB) error {
	rm.CreatedAt = time.Now() // 创建时间设置为当前时间
	return nil
}
