package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"sky_ISService/proto/system"
	"sky_ISService/services/system/service"
	"sync"
)

var grpcServer *grpc.Server
var once sync.Once

// StartSystemGRPCServer 启动 gRPC 服务端
func StartSystemGRPCServer() error {
	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("监听端口失败", err)
		return err
	}

	grpcServer = grpc.NewServer()
	system.RegisterSystemServiceServer(grpcServer, &service.AdminsService{})

	fmt.Println("gRPC 服务器开始监听 9999 端口...")
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Println("gRPC 服务器启动失败", err)
		return err
	}

	fmt.Println("gRPC 服务成功启动")
	return nil
}

// GracefulShutdown 停止 gRPC 服务并清理资源
func GracefulShutdown(ctx context.Context) error {
	// 使用 sync.Once 确保 Stop 只会被调用一次
	once.Do(func() {
		if grpcServer != nil {
			// 停止 gRPC 服务
			grpcServer.GracefulStop()
		}
	})
	return nil
}
