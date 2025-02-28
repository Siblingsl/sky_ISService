package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DBMiddleware 用于将数据库实例注入到请求上下文
func DBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将数据库实例传递到上下文中
		c.Set("db", db)

		// 在请求处理完后关闭数据库连接
		defer c.Set("db", nil)

		// 执行后续处理
		c.Next()
	}
}

// 获取数据库实例
func GetDB(c *gin.Context) (*gorm.DB, error) {
	db, exists := c.Get("db")
	if !exists {
		return nil, fmt.Errorf("无法获取数据库实例")
	}
	return db.(*gorm.DB), nil
}
