package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"sky_ISService/proto/system"
)

// NewSecurityToSystemClient 创建 gRPC 客户端并向 system 服务发送请求
func NewSecurityToSystemClient() (system.SystemServiceClient, error) {
	// 创建与 system 服务的连接
	conn, err := grpc.Dial("localhost:9999", grpc.WithInsecure()) // 连接到 system 服务
	if err != nil {
		return nil, fmt.Errorf("无法连接到 system 服务: %v", err)
	}
	client := system.NewSystemServiceClient(conn)
	return client, nil
}
