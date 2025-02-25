package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"log"
	"os"
	moduleSystem "sky_ISService/services/system/module"
)

func main() {
	// 创建 Gin 引擎实例
	r := gin.Default()

	// 使用 Fx 创建应用
	app := fx.New(
		// 提供 gin.Engine 实例到容器中
		fx.Provide(
			func() *gin.Engine {
				return r
			},
		),

		// 注册 SystemModules
		moduleSystem.SystemModules,

		// 启动时运行的函数
		fx.Invoke(func() {
			// 启动服务
			port := os.Getenv("PORT")
			if port == "" {
				port = "8082"
			}
			// 启动 Gin 引擎
			if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
				log.Fatalf("System 服务启动失败: %v", err)
			}
		}),
	)

	// 启动 Fx 应用
	app.Run()
}
