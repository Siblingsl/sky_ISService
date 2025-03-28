package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"sync"
)

// Config 定义一个全局变量来存储配置
var Config *InitStructureConfig
var once sync.Once

// ServerConfig 网关总服务
type ServerConfig struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Addr    string `mapstructure:"addr"`
	Weight1 int    `mapstructure:"weight1"`
	Weight2 int    `mapstructure:"weight2"`
}

// SecurityConfig 认证服务配置
type SecurityConfig struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Port1   string `mapstructure:"port1"`
	Addr    string `mapstructure:"addr"`
	Weight1 int    `mapstructure:"weight1"`
	Weight2 int    `mapstructure:"weight2"`
}

// SystemConfig 系统服务配置
type SystemConfig struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Port1   string `mapstructure:"port1"`
	Addr    string `mapstructure:"addr"`
	Weight1 int    `mapstructure:"weight1"`
	Weight2 int    `mapstructure:"weight2"`
}

// 默认服务配置
type defaultConfig struct {
	Addr   string `mapstructure:"addr"`
	Weight int    `mapstructure:"weight"`
}

// PathConfig 配置结构
type PathConfig struct {
	Security string `mapstructure:"security"`
	System   string `mapstructure:"system"`
}

// ElasticsearchConfig Elasticsearch 配置结构
type ElasticsearchConfig struct {
	Host     string   `mapstructure:"host"`
	Port     string   `mapstructure:"port"`
	Indexes  []string `mapstructure:"indexes"`
	User     string   `mapstructure:"user"`
	Password string   `mapstructure:"password"`
}

// LogstashConfig logstash 结构配置
type LogstashConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// PostgresSQLConfig PostgresSQL 配置结构
type PostgresSQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// RedisConfig Redis 配置结构
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// ConsulConfig consul 结构配置
type ConsulConfig struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

// LoggerConfig logger 结构配置
type LoggerConfig struct {
	Level    []string `mapstructure:"level"`
	Filepath string   `mapstructure:"filepath"`
}

// RabbitMQConfig RabbitMQ 配置结构
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"virtual_host"`
}

// JWTSecret JWT
type JWTSecret struct {
	Secret string `mapstructure:"secret"`
}

// AESSecret AES加密
type AESSecret struct {
	Secret string `mapstructure:"secret"`
}

// InitStructureConfig 配置结构
type InitStructureConfig struct {
	// 服务配置
	Server   ServerConfig   `mapstructure:"server"`
	Security SecurityConfig `mapstructure:"security"`
	System   SystemConfig   `mapstructure:"system"`
	Default  defaultConfig  `mapstructure:"default"`

	// 子服务路径
	PathConfig PathConfig `mapstructure:"path_config"`

	// 数据库配置
	Database map[string]struct {
		PostgreSQL PostgresSQLConfig `mapstructure:"postgresql"`
	} `mapstructure:"database"`

	// 缓存 配置
	Cache struct {
		Redis RedisConfig `mapstructure:"redis"`
	} `mapstructure:"cache"`

	// Elasticsearch 配置
	ELK struct {
		Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
	} `mapstructure:"elk"`

	// 服务注册与发现  consul
	RegisterService struct {
		Consul ConsulConfig `mapstructure:"consul"`
	} `mapstructure:"register_service"`

	// 日志配置
	Logger map[string]struct {
		LoggerConfig  LoggerConfig        `mapstructure:"logger"`
		Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
		Logstash      LogstashConfig      `mapstructure:"logstash"`
	} `mapstructure:"logger"`

	// 消息队列配置
	MessageQueue struct {
		RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	} `mapstructure:"message_queue"`

	// JWT
	JWTSecret JWTSecret `mapstructure:"jwt_secret"`

	// AES
	AESSecret AESSecret `mapstructure:"aes_secret"`
}

// InitLoadConfig 加载配置文件
func InitLoadConfig() (*InitStructureConfig, error) {
	var config InitStructureConfig

	// 直接获取当前目录的绝对路径
	absPath, err := filepath.Abs("config/config.yml")

	// 直接指定完整的配置文件路径
	viper.SetConfigFile(absPath)

	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件: %v", err)
	}

	// 将配置文件解析为结构体
	err = viper.Unmarshal(&config)
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
