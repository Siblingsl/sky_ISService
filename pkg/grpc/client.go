package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync"
)

// GRpcClient 管理 gRPC 连接
type GRpcClient struct {
	conn *grpc.ClientConn
}

var (
	once     sync.Once
	instance *GRpcClient
)

// NewGRpcClient 创建 gRPC 客户端，连接到指定的服务
func NewGRpcClient(host string, port int) *GRpcClient {
	once.Do(func() {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("无法连接 gRPC 服务器: %v", err)
		}

		instance = &GRpcClient{
			conn: conn,
		}
	})
	return instance
}

// GetConn 获取 gRPC 连接
func (c *GRpcClient) GetConn() *grpc.ClientConn {
	return c.conn
}

// Close 关闭 gRPC 连接
func (c *GRpcClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
