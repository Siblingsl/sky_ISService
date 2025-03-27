package models

import (
	"gorm.io/gorm"
	"sky_ISService/utils/database"
	"time"
)

// SkySystemAdmins 代表 sky_system_admins 表
type SkySystemAdmins struct {
	database.CommonBase `gorm:"embedded"` // 继承公共字段
	ID                  int               `gorm:"primaryKey;autoIncrement" json:"id"` // 使用 uint 类型
	Username            string            `gorm:"unique;not null" json:"username"`    // 用户名
	Password            string            `gorm:"not null" json:"-"`                  // 密码（JSON 不返回）
	FullName            string            `json:"full_name"`                          // 全名
	UserType            string            `json:"user_type"`                          // 管理员类型（00系统管理员）
	Email               string            `gorm:"unique" json:"email"`                // 邮箱
	Phone               string            `json:"phone"`                              // 电话
	Token               string            `json:"-"`                                  // Token（JSON 不返回）
}

// BeforeCreate 创建前 Hook
func (a *SkySystemAdmins) BeforeCreate(tx *gorm.DB) error {
	currentTime, ok := tx.Statement.Context.Value("current_time").(time.Time)
	if !ok {
		// 如果 Context 里没有这个值，就使用当前时间
		currentTime = time.Now()
	}
	a.CommonBase.CreatedAt = currentTime
	a.CommonBase.UpdatedAt = currentTime
	return nil
}

// BeforeUpdate 更新前 Hook
func (a *SkySystemAdmins) BeforeUpdate(tx *gorm.DB) error {
	a.CommonBase.UpdatedAt = tx.Statement.Context.Value("current_time").(time.Time)
	return nil
}
