package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RetryMiddleware 失败时进行重试 (重试机制中间件)
func RetryMiddleware(maxRetries int, retryDelay time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		for i := 0; i < maxRetries; i++ {
			c.Next()
			if c.Writer.Status() < http.StatusInternalServerError {
				return // 成功，退出
			}
			time.Sleep(retryDelay) // 失败，等待后重试
		}
		// 最终失败
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "服务不可用"})
		c.Abort()
	}
}
