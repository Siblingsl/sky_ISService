package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
)

// Config 配置结构体，用于控制初始化时是否启用特定的功能
type Config struct {
	InitRouteLoggerMiddleware bool // 控制是否启用路由日志中间件
}

// InitConfig 根据配置初始化功能
func InitConfig(r *gin.Engine, config Config) {
	if config.InitRouteLoggerMiddleware {
		InitRouteLoggerMiddleware(r)
	}
}

// InitRouteLoggerMiddleware 在控制台打印所有已注册的路由（初始化中间件，只执行一次）
func InitRouteLoggerMiddleware(r *gin.Engine) {
	// 打印已注册的路由
	for _, route := range r.Routes() {
		log.Println(route.Method, route.Path)
	}
}
