package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sky_ISService/proto/system"
)

// GRpcServer 封装 gRPC 服务端的启动和停止逻辑
type GRpcServer struct {
	server *grpc.Server
	port   int
}

// NewGRpcServer 创建一个新的 gRPC 服务实例
func NewGRpcServer(systemUserService system.SystemServiceServer) *GRpcServer {
	grpcServer := grpc.NewServer()
	system.RegisterSystemServiceServer(grpcServer, systemUserService)
	reflection.Register(grpcServer)

	return &GRpcServer{
		server: grpcServer,
		port:   50051, // 默认端口
	}
}

// Start 启动 gRPC 服务
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

// Stop 关闭 gRPC 服务
func (s *GRpcServer) Stop() {
	fmt.Println("正在关闭 gRPC 服务器...")
	s.server.GracefulStop()
}
