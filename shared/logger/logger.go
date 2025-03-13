package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sky_ISService/config"
	"strings"
)

// Logger 是封装后的日志实例
var Logger *logrus.Logger

// LogrusElasticHook 用于将日志发送到 Elasticsearch
type LogrusElasticHook struct {
	Client *elasticsearch.Client
	Index  string // 可通过配置传递索引名称
}

func (hook *LogrusElasticHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *LogrusElasticHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		entry.Context = context.Background()
	}

	// 将文档转为 JSON
	docJSON, err := json.Marshal(entry.Data)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	// 发送日志数据存入 Elasticsearch
	_, err = hook.Client.Index(
		hook.Index, // 使用传递的索引名称
		bytes.NewReader(docJSON),
		hook.Client.Index.WithOpType("_doc"),
	)
	return err
}

// InitLogger 初始化日志系统，接受一个 Elasticsearch 客户端作为参数
func InitLogger(serviceName string, configPath string, client *elasticsearch.Client) (*logrus.Logger, error) {
	fmt.Printf("日志初始化配置: serviceName=%s, configPath=%s\n", serviceName, configPath)
	// 加载配置
	configLogger, err := config.InitLoadConfig()
	if err != nil {
		return nil, fmt.Errorf("无法加载配置: %v", err)
	}

	fmt.Printf("Elasticsearch 配置: %+v\n", configLogger.ELK.Elasticsearch)

	// 确保 serviceName 存在于配置中
	serviceConfig, exists := configLogger.Logger[serviceName]
	if !exists {
		return nil, fmt.Errorf("没有找到服务 %s 的日志配置", serviceName)
	}

	// 初始化 logrus
	Logger := logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{}) // 设置日志格式为 JSON

	// 设置日志级别
	logLevel := logrus.InfoLevel
	if len(serviceConfig.LoggerConfig.Level) > 0 {
		logLevel, err = logrus.ParseLevel(strings.ToLower(serviceConfig.LoggerConfig.Level[0]))
		if err != nil {
			return nil, fmt.Errorf("无效的日志级别配置: %v", err)
		}
	}
	Logger.SetLevel(logLevel)

	// 设置日志输出路径
	logFilePath := "logfile.log"
	if serviceConfig.LoggerConfig.Filepath != "" {
		logFilePath = serviceConfig.LoggerConfig.Filepath
	}

	// 确保日志文件路径存在
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Info("无法打开日志文件，使用标准错误输出")
	} else {
		Logger.SetOutput(file)
	}

	// Elasticsearch Hook 设置
	if len(serviceConfig.Elasticsearch.Indexes) == 0 {

		return nil, fmt.Errorf("没有配置有效的 Elasticsearch 索引")
	}

	// 使用配置中的第一个索引名称
	hook := &LogrusElasticHook{
		Client: client,
		Index:  serviceConfig.Elasticsearch.Indexes[0], // 使用配置中的第一个索引
	}
	Logger.AddHook(hook)

	// 将初始化的 Logger 设置为全局 Logger
	SetLogger(Logger)

	return Logger, nil
}

// SetLogger 用于设置全局 Logger
func SetLogger(logger *logrus.Logger) {
	Logger = logger
}
