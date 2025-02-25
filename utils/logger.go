// logger/logger.go

package logger

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Logger 是封装后的日志实例
var Logger *logrus.Logger

// LoggerConfig 用于存储从 YAML 配置文件加载的配置信息
type LoggerConfig struct {
	Services map[string]struct {
		Elasticsearch struct {
			Host     string   `yml:"host"`
			Indexes  []string `yml:"indexes"`
			User     string   `yml:"user"`
			Password string   `yml:"password"`
		} `yml:"elasticsearch"`
	} `yml:"services"`
}

// InitLogger 初始化日志系统，接受一个 Elasticsearch 客户端作为参数
func InitLogger(client *elastic.Client, configPath string) error {
	// 加载配置文件
	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("无法加载配置: %v", err)
	}

	// 初始化 logrus
	Logger = logrus.New()

	// 设置日志格式为 JSON
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// 设置日志级别
	Logger.SetLevel(logrus.InfoLevel)

	// 设置输出到文件或标准输出
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Info("无法打开日志文件，使用默认的标准错误输出")
	} else {
		Logger.SetOutput(file)
	}

	// 检查是否有指定的索引配置
	serviceConfig, exists := config.Services["auth"]
	if !exists || len(serviceConfig.Elasticsearch.Indexes) == 0 {
		return fmt.Errorf("没有配置有效的 Elasticsearch 索引")
	}

	// Elasticsearch Hook 设置
	hook := &LogrusElasticHook{
		Client: client,
		Index:  serviceConfig.Elasticsearch.Indexes[0], // 使用第一个索引
	}
	Logger.AddHook(hook)
	return nil
}

// LoadConfig 加载 YAML 配置文件并返回配置结构
func LoadConfig(configPath string) (LoggerConfig, error) {
	var config LoggerConfig
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("无法读取配置文件: %v", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("无法解析 YAML 文件: %v", err)
	}
	return config, nil
}

// LogrusElasticHook Elasticsearch Hook 用于将日志发送到 Elasticsearch
type LogrusElasticHook struct {
	Client *elastic.Client
	Index  string
}

func (hook *LogrusElasticHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *LogrusElasticHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		entry.Context = context.Background()
	}

	// 这里可以添加日志内容调试，查看发送到 Elasticsearch 的内容
	fmt.Printf("日志记录到 Elasticsearch: %v\n", entry.Data)

	// 检查 entry.Data 是否有内容并且转换成正确格式
	if entry.Data == nil {
		return fmt.Errorf("日志数据为空")
	}

	// 将日志内容发送到 Elasticsearch
	_, err := hook.Client.Index().
		Index(hook.Index).
		BodyJson(entry.Data). // 确保 entry.Data 包含数据
		Do(entry.Context)     // 使用 entry.Context
	return err
}

// LogInfo 记录手动日志
func LogInfo(message string) {
	Logger.WithFields(logrus.Fields{
		"message": message, // 确保有 message 字段
		"level":   "info",
	}).Info(message)
}

// LogError 记录错误日志
func LogError(message string, err error) {
	Logger.WithFields(logrus.Fields{
		"message": message,
		"error":   err.Error(),
		"level":   "error",
	}).Error(message)
}
