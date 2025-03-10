package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync"
)

// GRpcClient 结构体，持久化 gRPC 连接
type GRpcClient struct {
	conn *grpc.ClientConn
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
			conn: conn,
		}
	})

	return instance
}

// GetConn 获取连接
func (c *GRpcClient) GetConn() *grpc.ClientConn {
	return c.conn
}

// Close 关闭 gRPC 连接
func (c *GRpcClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
