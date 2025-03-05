package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// RecoveryMiddleware 捕获 panic，防止服务崩溃
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				log.Printf("捕获到异常: %v", err)
				// 返回 500 错误
				c.JSON(http.StatusInternalServerError, gin.H{"message": "内部服务器错误"})
				c.Abort()
			}
		}()
		c.Next()
	}
}
