package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"log"
	moduleAuth "sky_ISService/services/auth/module"
)

func main() {

	// 使用 Fx 创建应用
	app := fx.New(

		// 注册 AuthModule
		moduleAuth.AuthModule,

		// 启动时运行的函数
		fx.Invoke(func(r *gin.Engine,
		) {

			// 启动 Gin 引擎
			fmt.Println(fmt.Sprintf("%s:%s", "0.0.0.0", "8092"))
			go func() {
				if err := r.Run(fmt.Sprintf("%s:%s", "0.0.0.0", "8092")); err != nil {
					log.Fatalf("服务启动失败: %v", err)
				}
			}()

		}),
	)

	// 启动 Fx 应用
	app.Run()
}
