package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"sky_ISService/gateway/router"
	"sky_ISService/pkg/middleware"
)

func main() {
	r := gin.Default()

	// 连接 Auth gRPC
	//conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
	//authClient := auth.NewAuthServiceClient(conn)

	// 注册中间件
	r.Use(middleware.RecoveryMiddleware())
	//r.Use(middlewares.AuthMiddleware(authClient))
	r.Use(middleware.CircuitMiddleware())
	r.Use(middleware.ErrorHandlingMiddleware())

	// 注册路由
	router.SetupRoutes(r)

	// 启动网关
	log.Fatal(r.Run(":8080"))
}
