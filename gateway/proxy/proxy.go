package proxy

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy 反向代理
func ReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析目标 URL
		targetURL, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// 修改请求头
		c.Request.Host = targetURL.Host
		c.Request.URL.Path = c.Param("path") // 保留原始 Path

		// 执行代理
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// RegisterRoutes 注册 API Gateway 的路由
//func RegisterRoutes(router *gin.Engine) {
//	router.GET("/auth/*path", ReverseProxy("http://localhost:8081"))   // 身份验证服务
//	router.GET("/system/*path", ReverseProxy("http://localhost:8082")) // 系统管理服务
//
//	//router.GET("/order/*path", ReverseProxy("http://localhost:5002"))       // 订单服务
//	//router.GET("/payment/*path", ReverseProxy("http://localhost:5003"))     // 支付服务
//	//router.GET("/logistics/*path", ReverseProxy("http://localhost:5004"))   // 物流服务
//	//router.GET("/inventory/*path", ReverseProxy("http://localhost:5005"))   // 库存服务
//}
