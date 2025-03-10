package moduleSystem

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"sky_ISService/services/system/controller"
	"sky_ISService/services/system/repository"
	"sky_ISService/services/system/service"
)

var SystemModules = fx.Options(
	// 提供依赖项
	fx.Provide(
		// user
		controller.NewUserController,
		service.NewUserService,
		repository.NewUserRepository,
		// role
		controller.NewRoleController,
		service.NewRoleService,
		repository.NewRoleRepository,
		// menu
		controller.NewMenuController,
		service.NewMenuService,
		repository.NewMenuRepository,
	),

	// 注册路由
	fx.Invoke(func(userController *controller.UserController, roleController *controller.RoleController, menuController *controller.MenuController, r *gin.Engine) {
		// 注册 user 路由
		userController.UserControllerRoutes(r)
		// 注册 role 路由
		roleController.RoleControllerRoutes(r)
		// 注册 menu 路由
		menuController.MenuControllerRoutes(r)

		// 配置初始化的中间件
		//config := initialize.Config{
		//	InitRouteLoggerMiddleware: true, // 控制是否启用控制台路由日志中间件
		//}
		//// 根据配置初始化功能
		//initialize.InitConfig(r, config)
	}),
)
