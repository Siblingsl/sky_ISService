package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sky_ISService/proto/auth"
	"sky_ISService/services/auth/service"
)

// GRpcServer 结构体
type GRpcServer struct {
	server *grpc.Server
	port   int
}

// NewGRpcServer 创建 gRPC 服务器实例
func NewGRpcServer(authService *service.AuthService) *GRpcServer {
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)

	return &GRpcServer{
		server: grpcServer,
		port:   50051, // 默认端口
	}
}

// Start 启动 gRPC 服务器
func (s *GRpcServer) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("无法监听端口 %d: %v", s.port, err)
	}

	fmt.Printf("gRPC 服务器启动，监听端口 %d...\n", s.port)
	if err := s.server.Serve(lis); err != nil {
		log.Fatalf("gRPC 服务器启动失败: %v", err)
	}
}

// Stop 关闭 gRPC 服务器
func (s *GRpcServer) Stop() {
	fmt.Println("正在关闭 gRPC 服务器...")
	s.server.GracefulStop()
}
