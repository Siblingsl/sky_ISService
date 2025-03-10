package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/services/auth/dto"
	"sky_ISService/services/auth/service"
	"sky_ISService/utils"
)

// AuthController 处理认证相关的请求
type AuthController struct {
	service *service.AuthService
}

// NewAuthController 构造函数，依赖注入 AuthService
func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		service: authService,
	}
}

// AuthControllerRoutes 路由类
func (c *AuthController) AuthControllerRoutes(r *gin.Engine) {
	// 创建前缀的路由组
	authGroup := r.Group("/auth")

	// 管理员登陆
	authGroup.POST("/admins/login", func(ctx *gin.Context) {
		var req dto.AdminLoginRequest
		// 绑定请求数据
		if err := ctx.ShouldBindJSON(&req); err != nil {
			// 请求数据错误时返回错误响应
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误")
			return
		}

		// 调用服务层进行登陆验证
		token, err := c.service.AdminLogin(ctx, req)
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// 返回成功响应
		utils.Success(ctx, dto.AdminLoginResponse{
			Token: token,
		})
	})

	// 发送验证码
	authGroup.GET("/admins/code", func(ctx *gin.Context) {
		email := ctx.Query("email")
		if email == "" {
			utils.Error(ctx, http.StatusBadRequest, "邮箱不能为空")
		}
		err := c.service.SendEmailCode(ctx, email)
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}
		utils.Success(ctx, "验证码发送成功")
	})

	// 注册路由
	//authGroup.GET("/register", func(ctx *gin.Context) {
	//	username := ctx.DefaultQuery("username", "shilei")
	//	password := ctx.DefaultQuery("password", "123456")
	//
	//	registerRequest := &pb.RegisterRequest{
	//		Username: username,
	//		Password: password,
	//	}
	//
	//	res, err := c.service.Register(ctx, registerRequest)
	//	if err != nil {
	//		// 记录错误日志
	//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
	//		return
	//	}
	//	utils.Success(ctx, res)
	//})

	// 登陆路由
	//authGroup.GET("/login", func(ctx *gin.Context) {
	//	username := ctx.DefaultQuery("username", "shilei")
	//	password := ctx.DefaultQuery("password", "123456")
	//
	//	// 参数转换为 gRPC 请求
	//	loginRequest := &pb.LoginRequest{
	//		Username: username,
	//		Password: password,
	//	}
	//	// 调用 gRPC 服务
	//	res, err := c.service.Login(ctx, loginRequest)
	//	if err != nil {
	//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//		return
	//	}
	//	utils.Success(ctx, res)
	//})
	//
	//authGroup.GET("/demo", func(ctx *gin.Context) {
	//	ctx.String(200, "12131654513") // 正确的响应方式，返回状态码和内容
	//})
}
