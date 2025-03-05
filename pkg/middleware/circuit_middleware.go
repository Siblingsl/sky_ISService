package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

// 熔断器
type CircuitBreaker struct {
	failureCount int
	state        string
	mu           sync.Mutex
}

var breaker = &CircuitBreaker{state: "CLOSED"}

// CircuitMiddleware (熔断机制中间件)
func CircuitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		breaker.mu.Lock()
		if breaker.state == "OPEN" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "服务不可用"})
			breaker.mu.Unlock()
			c.Abort()
			return
		}
		breaker.mu.Unlock()

		// 处理请求
		c.Next()

		// 监控状态
		breaker.mu.Lock()
		if c.Writer.Status() >= http.StatusInternalServerError {
			breaker.failureCount++
			if breaker.failureCount > 5 { // 连续5次失败，熔断
				breaker.state = "OPEN"
				time.AfterFunc(10*time.Second, func() { // 10秒后恢复
					breaker.mu.Lock()
					breaker.state = "CLOSED"
					breaker.failureCount = 0
					breaker.mu.Unlock()
				})
			}
		}
		breaker.mu.Unlock()
	}
}
