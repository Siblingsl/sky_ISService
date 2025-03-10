package swagger

import (
	"github.com/gin-gonic/gin"
	ginSwaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitSwagger 用于初始化 Swagger 配置
func InitSwagger(r *gin.Engine) {
	// 注册 Swagger UI 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(ginSwaggerFiles.Handler))
}
