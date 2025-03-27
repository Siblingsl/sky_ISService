package models

import (
	"sky_ISService/utils/database"
)

// SkyAuthUser 继承 CommonBase
type SkyAuthUser struct {
	database.CommonBase `gorm:"embedded"` // 继承公共字段
	ID                  uint              `gorm:"primaryKey;autoIncrement" json:"id"`                // 使用 uint 类型
	Username            string            `gorm:"type:varchar(100);unique;not null" json:"username"` // 用户名
	Password            string            `gorm:"type:varchar(255);not null" json:"password"`        // 密码（加密存储）
	Email               string            `gorm:"type:varchar(255)" json:"email"`                    // 邮箱
	Code                int               `gorm:"type:int" json:"code"`                              // 验证码
	Phone               string            `gorm:"type:varchar(20)" json:"phone"`                     // 电话
}
