package models

import (
	"sky_ISService/utils/database"
)

// SkyAuthUser 继承 CommonBase
type SkyAuthUser struct {
	database.CommonBase        // 嵌套 CommonBase 结构体
	Username            string `gorm:"type:varchar(100);unique;not null" json:"username"` // 用户名
	Password            string `gorm:"type:varchar(255);not null" json:"password"`        // 密码（加密存储）
	Email               string `gorm:"type:varchar(255)" json:"email"`                    // 邮箱
	Code                int    `gorm:"type:int" json:"code"`                              // 验证码
	Phone               string `gorm:"type:varchar(20)" json:"phone"`                     // 电话
	UserType            int    `gorm:"type:int;not null" json:"user_type"`                // 用户类型：1-管理员，2-客户
	Status              bool   `gorm:"type:boolean;default:true" json:"status"`           // 用户状态：是否启用
}
