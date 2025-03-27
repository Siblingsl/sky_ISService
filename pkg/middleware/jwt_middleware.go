package middleware

import (
	"log"
	"net/http"
	"sky_ISService/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 定义不需要 token 验证的路径
		noAuthPaths := []string{
			"/auth/admins/login",
			"/auth/admins/code",
			"/auth/admins/getTest",
			"/swagger/index.html",
		}

		// 检查请求路径是否在不需要验证的路径列表中
		for _, path := range noAuthPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// 从 Header 获取 Token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			// 未提供 Token
			utils.Error(c, http.StatusUnauthorized, "未提供 Token")
			c.Abort()
			return
		}

		// 去除前缀 "Bearer "
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 解析 Token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			// Token 无效
			log.Println("Error parsing token:", err)
			utils.Error(c, http.StatusUnauthorized, "无效的 Token: "+err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims["sub_id"]) // "sub_id" 是用户 ID
		c.Set("role", claims["role"])

		// 继续处理请求
		c.Next()
	}
}
