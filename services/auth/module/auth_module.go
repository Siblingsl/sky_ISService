package moduleAuth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"sky_ISService/services/auth/controller"
	"sky_ISService/services/auth/repository"
	"sky_ISService/services/auth/repository/models"
	"sky_ISService/services/auth/service"
	"sky_ISService/utils/database"
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
	}),

	// 调用自动迁移，注册并迁移所有模型
	fx.Invoke(func(db *gorm.DB) {
		// 将所有模型添加到迁移列表
		database.ModelsToMigrate = append(
			database.ModelsToMigrate,
			&models.SkyAuthUser{},
			&models.SkyAuthToken{},
		)

		// 执行自动迁移
		if err := database.AutoMigrate(db); err != nil {
			panic("迁移失败: " + err.Error())
		}
	}),
)
