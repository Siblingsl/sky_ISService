package router

import (
	"github.com/gin-gonic/gin"
	"sky_ISService/gateway/proxy"
)

// SetupRoutes 配置 API Gateway 路由
func SetupRoutes(router *gin.Engine) {
	router.GET("/auth/*path", proxy.ReverseProxy("http://localhost:8081"))   // 身份验证服务
	router.GET("/system/*path", proxy.ReverseProxy("http://localhost:8082")) // 系统管理服务
}
