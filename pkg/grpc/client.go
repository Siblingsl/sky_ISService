package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"

	pb "sky_ISService/proto/auth"
)

// GRpcClient 结构体，持久化 gRPC 连接
type GRpcClient struct {
	conn   *grpc.ClientConn
	client pb.AuthServiceClient
}

var (
	once     sync.Once
	instance *GRpcClient
)

// NewGRpcClient 创建 gRPC 连接
func NewGRpcClient(host string, port int) *GRpcClient {
	once.Do(func() {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("无法连接 gRPC 服务器: %v", err)
		}

		instance = &GRpcClient{
			conn:   conn,
			client: pb.NewAuthServiceClient(conn),
		}
	})

	return instance
}

// Login 调用 gRPC 认证服务
func (c *GRpcClient) Login(username, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("调用 Login 失败: %v", err)
	}

	return resp.Token, nil
}

// Close 关闭 gRPC 连接
func (c *GRpcClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
