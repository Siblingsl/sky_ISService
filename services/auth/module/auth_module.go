package moduleAuth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"sky_ISService/pkg/middleware"
	"sky_ISService/services/auth/controller"
	"sky_ISService/services/auth/repository"
	"sky_ISService/services/auth/service"
)

var AuthModule = fx.Options(
	// 提供依赖项
	fx.Provide(
		repository.NewAuthRepository, // 提供 repository
		service.NewAuthService,       // 提供 service
		controller.NewAuthController, // 提供 controller
	),

	// 注册路由
	fx.Invoke(func(authController *controller.AuthController, r *gin.Engine) {
		// 通过 controller 注册路由
		authController.AuthControllerRoutes(r)

		// 配置初始化的中间件
		config := middleware.Config{
			InitRouteLoggerMiddleware: true, // 控制是否启用控制台路由日志中间件
		}
		// 根据配置初始化功能
		middleware.InitConfig(r, config)
	}),
)
