package moduleSystem

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"sky_ISService/services/system/controller"
	"sky_ISService/services/system/repository"
	"sky_ISService/services/system/repository/models"
	"sky_ISService/services/system/service"
	"sky_ISService/utils/database"
)

var SystemModules = fx.Options(
	// 提供依赖项
	fx.Provide(
		// user
		controller.NewUserController,
		service.NewUserService,
		repository.NewAdminsRepository,
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
	fx.Invoke(func(userController *controller.AdminsController, roleController *controller.RoleController, menuController *controller.MenuController, r *gin.Engine) {
		// 注册 user 路由
		userController.UserControllerRoutes(r)
		// 注册 role 路由
		roleController.RoleControllerRoutes(r)
		// 注册 menu 路由
		menuController.MenuControllerRoutes(r)
	}),

	// 调用自动迁移，注册并迁移所有模型
	fx.Invoke(func(db *gorm.DB) {
		// 将所有模型添加到迁移列表
		database.ModelsToMigrate = append(
			database.ModelsToMigrate,
			&models.SkySystemAdmins{},
		)

		// 执行自动迁移
		if err := database.AutoMigrate(db); err != nil {
			panic("迁移失败: " + err.Error())
		}
	}),
)
