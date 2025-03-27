package models

import (
	"gorm.io/gorm"
	"sky_ISService/utils/database"
	"time"
)

// SkySystemMenus 代表 sky_system_menus 表
type SkySystemMenus struct {
	database.CommonBase `gorm:"embedded"` // 继承公共字段
	ID                  int               `gorm:"primaryKey;autoIncrement" json:"id"` // 使用 int 类型
	MenuName            string            `gorm:"not null" json:"menu_name"`          // 菜单名称
	MenuURL             string            `json:"menu_url"`                           // 菜单链接（可以为空）
	ParentID            int               `json:"parent_id"`                          // 父菜单 ID
	MenuSort            int               `json:"menu_sort"`                          // 菜单排序
	MenuType            int               `json:"menu_type"`                          // 菜单类型: 1-目录，2-菜单，3-按钮
	MenuIcon            string            `json:"menu_icon"`                          // 菜单图标
	Description         string            `json:"description"`                        // 描述
}

// BeforeCreate 创建前 Hook
func (m *SkySystemMenus) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前 Hook
func (m *SkySystemMenus) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
