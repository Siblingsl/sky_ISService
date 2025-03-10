package mq

import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
	"sky_ISService/config"
)

// RabbitMQClient 封装 RabbitMQ 客户端
type RabbitMQClient struct {
	Connection  *amqp.Connection
	ChannelPool chan *amqp.Channel // 连接池，维护多个 Channel
}

// 单例模式: 只有一个 RabbitMQClient 实例
var (
	rabbitMQInstance *RabbitMQClient
	once             sync.Once
)

// InitRabbitMQ 初始化 RabbitMQ 客户端并返回单例实例
// @param configPath string: 配置文件路径
// @return *RabbitMQClient: RabbitMQ 客户端实例
// @return error: 初始化过程中可能出现的错误
func InitRabbitMQ(configPath string) (*RabbitMQClient, error) {
	var err error
	once.Do(func() { // 确保只执行一次
		rabbitMQInstance, err = newRabbitMQClient(configPath)
	})
	return rabbitMQInstance, err
}

// newRabbitMQClient 创建 RabbitMQ 连接和通道池
// @param configPath string: 配置文件路径
// @return *RabbitMQClient: 初始化的 RabbitMQ 客户端实例
// @return error: 可能出现的错误
func newRabbitMQClient(configPath string) (*RabbitMQClient, error) {
	// 加载配置
	cfg, err := config.InitLoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	rmqConfig := cfg.MessageQueue.RabbitMQ

	// 创建连接字符串
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		rmqConfig.Username,
		rmqConfig.Password,
		rmqConfig.Host,
		rmqConfig.Port,
		rmqConfig.VHost,
	)

	// 连接到 RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 RabbitMQ: %v", err)
	}

	// 创建通道池
	channelPool := make(chan *amqp.Channel, 5) // 维护 5 个 channel
	for i := 0; i < cap(channelPool); i++ {
		ch, err := conn.Channel()
		if err != nil {
			err := conn.Close()
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("无法创建通道: %v", err)
		}
		channelPool <- ch
	}

	fmt.Println("成功连接到 RabbitMQ")

	return &RabbitMQClient{Connection: conn, ChannelPool: channelPool}, nil
}

// GetChannel 从连接池获取一个可用的 channel
// @return *amqp.Channel: 获取到的 channel
// @return error: 如果无法获取 channel，则返回错误
func (r *RabbitMQClient) GetChannel() (*amqp.Channel, error) {
	if r.ChannelPool == nil {
		// 如果 ChannelPool 是 nil，进行初始化
		r.ChannelPool = make(chan *amqp.Channel, 10) // 假设是一个大小为 10 的缓冲 channel
	}

	// 执行获取操作
	ch := <-r.ChannelPool
	if ch == nil {
		return nil, fmt.Errorf("channel is nil")
	}
	return ch, nil
}

// ReleaseChannel 释放 channel 回到连接池
// @param ch *amqp.Channel: 要释放的 channel
func (r *RabbitMQClient) ReleaseChannel(ch *amqp.Channel) {
	select {
	case r.ChannelPool <- ch:
		// 成功回收通道
	default:
		// 通道池已满，直接关闭
		err := ch.Close()
		if err != nil {
			log.Printf("关闭 channel 失败: %v", err)
		}
	}
}

// SendMessage 发送消息到 RabbitMQ 队列（使用连接池）
// @param queueName string: 队列名称
// @param message string: 要发送的消息
// @return error: 如果发送消息失败，返回错误
func (r *RabbitMQClient) SendMessage(queueName, message string) error {
	ch, err := r.GetChannel()
	if err != nil {
		log.Printf("无法获取通道: %s", err)
		return err
	}
	defer r.ReleaseChannel(ch) // 确保释放通道回到池中

	// 发送消息
	err = ch.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Printf("消息发送失败: %s", err)
		return err
	}

	log.Printf("消息发送成功: %s", message)
	return nil
}

// Close 关闭 RabbitMQ 连接和通道池
// @return error: 如果关闭过程中出现错误，返回错误
func (r *RabbitMQClient) Close() error {
	close(r.ChannelPool) // 关闭连接池
	for ch := range r.ChannelPool {
		err := ch.Close()
		if err != nil {
			log.Printf("关闭通道失败: %v", err)
		}
	}
	err := r.Connection.Close()
	if err != nil {
		return fmt.Errorf("关闭 RabbitMQ 连接失败: %v", err)
	}
	log.Println("RabbitMQ 连接已关闭")
	return nil
}
