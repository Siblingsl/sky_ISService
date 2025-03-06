package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	pb "sky_ISService/proto/auth"
	"sky_ISService/services/auth/service"
	"sky_ISService/utils"
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
		username := ctx.DefaultQuery("username", "shilei")
		password := ctx.DefaultQuery("password", "123456")

		registerRequest := &pb.RegisterRequest{
			Username: username,
			Password: password,
		}

		res, err := c.authService.Register(ctx, registerRequest)
		if err != nil {
			// 记录错误日志
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		utils.Success(ctx, res)
	})

	// 登陆路由
	authGroup.GET("/login", func(ctx *gin.Context) {
		username := ctx.DefaultQuery("username", "shilei")
		password := ctx.DefaultQuery("password", "123456")

		// 参数转换为 gRPC 请求
		loginRequest := &pb.LoginRequest{
			Username: username,
			Password: password,
		}
		// 调用 gRPC 服务
		res, err := c.authService.Login(ctx, loginRequest)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		utils.Success(ctx, res)
	})

	authGroup.GET("/demo", func(ctx *gin.Context) {
		ctx.String(200, "12131654513") // 正确的响应方式，返回状态码和内容
	})
}
