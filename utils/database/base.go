package database

import (
	"gorm.io/gorm"
	"log"
	"time"
)

// CommonBase 实体结构体，映射到数据库的 common_base 表
type CommonBase struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`                           // 主键，自增，使用 int 类型
	Status    bool      `gorm:"default:true" json:"status"`                                   // 状态字段，默认为 true
	IsDeleted bool      `gorm:"default:false" json:"is_deleted"`                              // 删除标志，默认为 false
	CreatedBy int       `gorm:"column:created_by" json:"created_by"`                          // 创建者，使用 int 类型
	UpdatedBy int       `gorm:"size:255" json:"updated_by"`                                   // 更新者，使用 int 类型
	CreatedAt time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
	Notes     string    `gorm:"type:text" json:"notes"`                                       // 备注
}

// ModelsToMigrate 用于保存所有需要迁移的模型
var ModelsToMigrate []interface{}

// AutoMigrate 自动迁移所有模型
func AutoMigrate(db *gorm.DB) error {
	// 遍历并自动迁移所有模型
	for _, model := range ModelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			log.Printf("迁移失败: %v", err)
			return err
		}
	}
	return nil
}
