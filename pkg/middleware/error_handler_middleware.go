package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/utils"
)

// ErrorHandlingMiddleware 统一处理错误并返回自定义的错误格式
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 进入后续处理
		c.Next()

		// 捕获可能的业务逻辑错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			// 判断错误类型
			var appErr *utils.AppError
			if errors.As(err.Err, &appErr) {
				// 自定义业务错误
				c.JSON(appErr.Code, gin.H{"error": appErr.Message})
				c.Abort()
				return
			}
			// 默认错误处理
			c.JSON(http.StatusInternalServerError, gin.H{"message": "内部服务器错误"})
		}
	}
}
