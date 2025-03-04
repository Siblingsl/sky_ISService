package elasticsearch

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"sky_ISService/config"
)

// InitElasticsearch 初始化 Elasticsearch 客户端
func InitElasticsearch(configPath string) (*elasticsearch.Client, error) {
	configEs, err := config.InitLoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// 获取指定服务的 Elasticsearch 配置
	serviceConfig := configEs.ELK

	cfg := elasticsearch.Config{
		Addresses: []string{serviceConfig.Elasticsearch.Host},
		Username:  serviceConfig.Elasticsearch.User,
		Password:  serviceConfig.Elasticsearch.Password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 Elasticsearch: %v", err)
	}

	fmt.Println("成功连接到 Elasticsearch")
	return client, nil
}

// InitElasticsearchConfig 初始化 Elasticsearch
func InitElasticsearchConfig(configPath string) (*elasticsearch.Client, error) {
	return InitElasticsearch(configPath)
}
