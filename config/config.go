package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
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

// LoggerConfig logger 结构配置
type LoggerConfig struct {
	Level    []string `yaml:"level"`
	Filepath string   `yaml:"filepath"`
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
	// 数据库配置
	Database map[string]struct {
		PostgreSQL PostgresSQLConfig `yaml:"postgresql"`
	} `yaml:"database"`

	// Elasticsearch 配置
	ELK struct {
		Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
	} `yaml:"elk"`

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
