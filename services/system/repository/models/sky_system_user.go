package models

import (
	"sky_ISService/utils/database"
	"time"
)

// SkySystemAdmins 代表 sky_system_admins 表
type SkySystemAdmins struct {
	database.CommonBase
	ID        int       `gorm:"primaryKey" json:"id"`            // 主键
	Username  string    `gorm:"unique;not null" json:"username"` // 用户名
	Password  string    `gorm:"not null" json:"-"`               // 密码，JSON 不返回
	FullName  string    `json:"full_name"`                       // 全名
	Email     string    `gorm:"unique" json:"email"`             // 邮箱
	Phone     string    `json:"phone"`                           // 电话
	Notes     string    `json:"notes"`                           // 备注
	Token     string    `json:"-"`                               // Token，不返回
	Status    bool      `json:"status"`                          // 状态（true: 启用, false: 禁用）
	CreatedBy int       `json:"created_by"`                      // 创建人 ID
	UpdatedBy int       `json:"updated_by"`                      // 更新人 ID
	CreatedAt time.Time `json:"created_at"`                      // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                      // 更新时间
	IsDeleted bool      `gorm:"default:false" json:"is_deleted"` // 删除标志

}
