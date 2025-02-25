package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/services/auth/service"
)

// AuthController 处理认证相关的请求
type AuthController struct {
	authService *service.AuthService
}

// NewAuthController 构造函数，依赖注入 AuthService
func NewAuthController(authService *service.AuthService) *AuthController {
	controller := &AuthController{authService: authService}
	return controller
}

// AuthControllerRoutes 路由类
func (c *AuthController) AuthControllerRoutes(r *gin.Engine) {
	// 创建前缀的路由组
	authGroup := r.Group("/auth")

	// 注册路由
	authGroup.GET("/register", func(ctx *gin.Context) {
		message, err := c.authService.Register()
		if err != nil {
			// 记录错误日志
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": message})
	})

	// 登陆路由
	authGroup.GET("/login", func(ctx *gin.Context) {
		message, err := c.authService.Login()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": message})
	})
}
