package models

import (
	"gorm.io/gorm"
	"sky_ISService/utils/database"
	"time"
)

// SkySystemRoles 代表 sky_system_roles 表
type SkySystemRoles struct {
	database.CommonBase `gorm:"embedded"` // 继承公共字段
	ID                  int               `gorm:"primaryKey;autoIncrement" json:"id"` // 使用 uint 类型
	RoleName            string            `gorm:"not null" json:"role_name"`          // 角色名称
	RoleKey             string            `gorm:"not null" json:"role_key"`           // 角色权限字符串
	RoleSort            int               `gorm:"not null" json:"role_sort"`          // 显示顺序
	Description         string            `json:"description"`                        // 角色描述
}

// BeforeCreate 创建前 Hook
func (r *SkySystemRoles) BeforeCreate(tx *gorm.DB) error {
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前 Hook
func (r *SkySystemRoles) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}
