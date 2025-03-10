package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// ElasticsearchConfig Elasticsearch 配置结构
type ElasticsearchConfig struct {
	Host     string   `yml:"host"`
	Port     string   `yml:"port"`
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

// RedisConfig Redis 配置结构
type RedisConfig struct {
	Host     string `yml:"host"`
	Port     string `yml:"port"`
	Password string `yml:"password"`
	DB       int    `yml:"db"`
}

// ConsulConfig consul 结构配置
type ConsulConfig struct {
	Address string `yml:"address"`
	Port    string `yml:"port"`
}

// LoggerConfig logger 结构配置
type LoggerConfig struct {
	Level    []string `yml:"level"`
	Filepath string   `yml:"filepath"`
}

// RabbitMQConfig RabbitMQ 配置结构
type RabbitMQConfig struct {
	Host     string `yml:"host"`
	Port     int    `yml:"port"`
	Username string `yml:"username"`
	Password string `yml:"password"`
	VHost    string `yml:"virtual_host"`
}

// JWTSecret JWT
type JWTSecret struct {
	Secret string `yml:"secret"`
}

// AESSecret AES加密
type AESSecret struct {
	Secret string `yml:"secret"`
}

type InitStructureConfig struct {
	// 数据库配置
	Database map[string]struct {
		PostgreSQL PostgresSQLConfig `yaml:"postgresql"`
	} `yaml:"database"`

	// 缓存 配置
	Cache struct {
		Redis RedisConfig `yaml:"redis"`
	} `yaml:"cache"`

	// Elasticsearch 配置
	ELK struct {
		Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
	} `yaml:"elk"`

	// 服务注册与发现  consul
	RegisterService struct {
		Consul ConsulConfig `yaml:"consul"`
	} `yaml:"register_service"`

	// 日志配置
	Logger map[string]struct {
		LoggerConfig  LoggerConfig        `yaml:"logger"`
		Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
		Logstash      LogstashConfig      `yaml:"logstash"`
	} `yaml:"logger"`

	// 消息队列配置
	MessageQueue struct {
		RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
	} `yaml:"message_queue"`

	// JWT
	JWTSecret JWTSecret `yaml:"jwt_secret"`

	// AES
	AESSecret AESSecret `yaml:"aes_secret"`
}

// InitLoadConfig 加载配置文件
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

// GetConfig  全局配置变量或函数  (使用举例： jwtKey := config.GetConfig().JWTSecret.Secret)
func GetConfig() *InitStructureConfig {
	var Config *InitStructureConfig

	return Config
}
