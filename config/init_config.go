package config

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
)

// ElasticsearchConfig Elasticsearch 配置结构
type ElasticsearchConfig struct {
	Host     string   `yml:"host"`
	Indexes  []string `yml:"indexes"`
	User     string   `yml:"user"`
	Password string   `yml:"password"`
}

// LogstashConfig logstash 结构配置
type LogstashConfig struct {
	Host string `yml:"host"`
	port string `yml:"port"`
}

// PostgresSQLConfig PostgresSQL 配置结构
type PostgresSQLConfig struct {
	Host     string `yml:"host"`
	Port     string `yml:"port"`
	Database string `yml:"database"`
	Username string `yml:"username"`
	Password string `yml:"password"`
}

// RabbitMQConfig RabbitMQ 配置结构
type RabbitMQConfig struct {
	Host     string `yml:"host"`
	Port     int    `yml:"port"`
	Username string `yml:"username"`
	Password string `yml:"password"`
	VHost    string `yml:"virtual_host"`
}

type InitStructureConfig struct {
	// 数据库
	Database map[string]struct {
		PostgreSQL PostgresSQLConfig `yml:"postgresql"`
	} `yml:"database"`
	// 日志
	Log map[string]struct {
		Logger struct {
			Level         []string            `yml:"level"`
			Filepath      string              `yml:"filepath"`
			Elasticsearch ElasticsearchConfig `yml:"elasticsearch"`
			Logstash      LogstashConfig      `yml:"logstash"`
		}
	} `yml:"log"`
	// 消息队列
	MessageQueue struct {
		RabbitMQ RabbitMQConfig `yml:"rabbitmq"`
	} `yml:"message_queue"`
}

// 加载配置文件
func InitLoadConfig(configPath string) (*InitStructureConfig, error) {
	var config InitStructureConfig
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件: %v", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("无法解析配置文件: %v", err)
	}
	return &config, nil
}

// 初始化 PostgreSQL 客户端
func InitPostgres(serviceName, configPath string) (*gorm.DB, error) {
	config, err := InitLoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// 获取指定服务的 PostgreSQL 配置
	serviceConfig, exists := config.Database[serviceName]
	if !exists {
		return nil, fmt.Errorf("服务配置不存在: %v", serviceName)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		serviceConfig.PostgreSQL.Host, serviceConfig.PostgreSQL.Username, serviceConfig.PostgreSQL.Password, serviceConfig.PostgreSQL.Database, serviceConfig.PostgreSQL.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("无法连接到 PostgreSQL: %v", err)
	}

	fmt.Println("成功连接到 PostgreSQL")
	return db, nil
}

// 初始化 Elasticsearch 客户端
func InitElasticsearch(serviceName, configPath string) (*elasticsearch.Client, error) {
	config, err := InitLoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// 获取指定服务的 Elasticsearch 配置
	serviceConfig, exists := config.Log[serviceName]
	if !exists {
		return nil, fmt.Errorf("服务配置不存在: %v", serviceName)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{serviceConfig.Logger.Elasticsearch.Host},
		Username:  serviceConfig.Logger.Elasticsearch.User,
		Password:  serviceConfig.Logger.Elasticsearch.Password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 Elasticsearch: %v", err)
	}

	fmt.Println("成功连接到 Elasticsearch")
	return client, nil
}

// InitRabbitMQ 初始化 RabbitMQ 连接
func InitRabbitMQ(configPath string) (*amqp.Connection, error) {
	config, err := InitLoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// 获取指定服务的 RabbitMQ 配置
	serviceConfig := config.MessageQueue.RabbitMQ

	// 连接字符串
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		serviceConfig.Username,
		serviceConfig.Password,
		serviceConfig.Host,
		serviceConfig.Port,
		serviceConfig.VHost,
	)

	// 连接 RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 RabbitMQ: %v", err)
	}

	fmt.Println("成功连接到 RabbitMQ")
	return conn, nil
}

// InitPostgresConfig 初始化 PostgreSQL
func InitPostgresConfig(serviceName, configPath string) (*gorm.DB, error) {
	return InitPostgres(serviceName, configPath)
}

// InitElasticsearchConfig 初始化 Elasticsearch
func InitElasticsearchConfig(serviceName, configPath string) (*elasticsearch.Client, error) {
	return InitElasticsearch(serviceName, configPath)
}

// InitRabbitMQConfig 初始化 RabbitMQ 连接
func InitRabbitMQConfig(configPath string) (*amqp.Connection, error) {
	return InitRabbitMQ(configPath)
}
