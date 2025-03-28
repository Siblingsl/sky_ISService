package models

import (
	"sky_ISService/utils/database"
)

// SkySecurityUser 继承 CommonBase
type SkySecurityUser struct {
	database.CommonBase `gorm:"embedded"` // 继承公共字段
	ID                  int               `gorm:"primaryKey;autoIncrement" json:"id"`                      // 使用 uint 类型
	Username            string            `gorm:"type:varchar(100);unique;not null;index" json:"username"` // 用户名
	Password            string            `gorm:"type:varchar(255);not null;index" json:"password"`        // 密码（加密存储）
	Email               string            `gorm:"type:varchar(255);index" json:"email"`                    // 邮箱
	Phone               string            `gorm:"type:varchar(20);index" json:"phone"`                     // 电话
}
