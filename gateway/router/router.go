package router

import (
	"github.com/gin-gonic/gin"
	"sky_ISService/gateway/proxy"
)

// NewRouter 创建一个新的 Gin 引擎实例，并将代理功能集成进去
func NewRouter(p *proxy.Proxy) *gin.Engine {
	r := gin.Default()
	r.Any("/*path", func(c *gin.Context) {
		p.ServeHTTP(c.Writer, c.Request)
	})
	return r
}
