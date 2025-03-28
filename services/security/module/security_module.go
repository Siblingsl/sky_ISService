package module

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"sky_ISService/services/security/controller"
	"sky_ISService/services/security/repository"
	"sky_ISService/services/security/repository/models"
	"sky_ISService/services/security/service"
	"sky_ISService/utils/database"
)

var SecurityModule = fx.Options(
	// 提供依赖项
	fx.Provide(
		repository.NewSecurityRepository,
		controller.NewSecurityController,
		service.NewSecurityService,
	),
	// 注册路由
	fx.Invoke(func(securityController *controller.SecurityController, r *gin.Engine) {
		securityController.SecurityControllerRoutes(r)
	}),
	// 调用自动迁移，注册并迁移所有模型
	fx.Invoke(func(db *gorm.DB, r *gin.Engine) {
		// 将所有模型添加到迁移列表
		database.ModelsToMigrate = append(
			database.ModelsToMigrate,
			&models.SkySecurityUser{},
		)
		// 执行自动迁移
		if err := database.AutoMigrate(db); err != nil {
			panic("迁移失败: " + err.Error())
		}
	}),
)
