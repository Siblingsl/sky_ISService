package config

import (
	"embed"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/fs"
	"sync"
)

// Config 定义一个全局变量来存储配置
var Config *InitStructureConfig
var once sync.Once // 用来确保配置只加载一次

//go:embed config.yml
var embeddedConfig embed.FS

// PathConfig 配置结构
type PathConfig struct {
	Auth   string `yaml:"auth"`
	System string `yaml:"system"`
}

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
	// 子服务路径
	PathConfig PathConfig `yaml:"path_config"`

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
func InitLoadConfig() (*InitStructureConfig, error) {
	var config InitStructureConfig

	// 从嵌入的文件系统读取配置文件内容
	data, err := fs.ReadFile(embeddedConfig, "config.yml")
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件: %v", err)
	}

	// 解析 YAML 配置文件
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("无法解析配置文件: %v", err)
	}
	return &config, nil
}

// GetConfig 获取全局配置
func GetConfig() *InitStructureConfig {
	once.Do(func() {
		// 只在第一次调用时加载配置
		config, err := InitLoadConfig()
		if err != nil {
			panic(fmt.Sprintf("加载配置失败: %v", err))
		}
		Config = config
	})
	return Config
}

//// InitLoadConfig 加载配置文件
//func InitLoadConfig(configPath string) (*InitStructureConfig, error) {
//	var config InitStructureConfig
//	data, err := os.ReadFile(configPath)
//	if err != nil {
//		return nil, fmt.Errorf("无法读取配置文件: %v", err)
//	}
//	err = yaml.Unmarshal(data, &config)
//	if err != nil {
//		return nil, fmt.Errorf("无法解析配置文件: %v", err)
//	}
//	return &config, nil
//}
//
//// GetConfig  全局配置变量或函数  (使用举例： jwtKey := config.GetConfig().JWTSecret.Secret)
//func GetConfig() *InitStructureConfig {
//	once.Do(func() {
//		// 只在第一次调用时加载配置
//		config, err := InitLoadConfig("config/config.yml") // 确保路径正确
//		if err != nil {
//			panic(fmt.Sprintf("加载配置失败: %v", err))
//		}
//		Config = config
//	})
//	return Config
//}
