package models

import (
	"sky_ISService/utils/database"
	"time"
)

// SkyAuthToken 继承 CommonBase
type SkyAuthToken struct {
	database.CommonBase           // 嵌套 CommonBase 结构体
	UserID              uint      `gorm:"type:int;not null" json:"user_id"`        // 关联用户表
	Token               string    `gorm:"type:varchar(512);not null" json:"token"` // 认证 token
	ExpiresAt           time.Time `gorm:"type:timestamptz" json:"expires_at"`      // 过期时间
}
