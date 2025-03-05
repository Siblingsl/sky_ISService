package middlewares

//import (
//	"context"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"sky_ISService/proto/auth"
//	"time"
//)
//
//func AuthMiddleware(client auth.AuthServiceClient) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 从请求头获取 token
//		token := c.GetHeader("Authorization")
//		if token == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{"message": "需要令牌"})
//			c.Abort()
//			return
//		}
//
//		// 调用 auth 服务验证 token
//		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//		defer cancel()
//
//		_, err := client.ValidateToken(ctx, &auth.TokenRequest{Token: token})
//		if err != nil {
//			c.JSON(http.StatusUnauthorized, gin.H{"message": "令牌无效"})
//			c.Abort()
//			return
//		}
//
//		// Token 校验通过，继续执行
//		c.Next()
//	}
//}
