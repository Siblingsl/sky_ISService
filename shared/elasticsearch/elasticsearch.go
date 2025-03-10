package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"io"
	"sky_ISService/config"
	"sync"
)

// ElasticsearchClient 封装 Elasticsearch 客户端
type ElasticsearchClient struct {
	Client *elasticsearch.Client // Elasticsearch 客户端实例
	Ctx    context.Context       // 上下文对象，用于管理 Elasticsearch 操作的生命周期
}

// 单例模式: 只有一个 ElasticsearchClient 实例
var (
	once                  sync.Once
	elasticsearchInstance *ElasticsearchClient
)

// initElasticsearch 初始化 Elasticsearch 客户端
// @param configPath string: 配置文件路径
// @return *ElasticsearchClient: Elasticsearch 客户端实例
// @return error: 初始化过程中可能出现的错误
func initElasticsearch(configPath string) (*ElasticsearchClient, error) {
	// 加载配置文件
	configElasticsearch, err := config.InitLoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// 获取指定服务的 Elasticsearch 配置
	cacheElasticConfig := configElasticsearch.ELK.Elasticsearch

	// 构建 Elasticsearch 连接参数
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", cacheElasticConfig.Host, cacheElasticConfig.Port),
		},
		Username: cacheElasticConfig.User,
		Password: cacheElasticConfig.Password,
	})

	if err != nil {
		return nil, fmt.Errorf("无法连接到 Elasticsearch: %v", err)
	}

	fmt.Println("成功连接到 Elasticsearch")

	return &ElasticsearchClient{
		Client: client,
		Ctx:    context.Background(),
	}, nil
}

// InitElasticsearchConfig 获取 Elasticsearch 单例
// @param configPath string: 配置文件路径
// @return *ElasticsearchClient: Elasticsearch 客户端实例
// @return error: 可能的错误，如果 Elasticsearch 无法初始化
func InitElasticsearchConfig(configPath string) (*ElasticsearchClient, error) {
	// 保证 Elasticsearch 客户端只被初始化一次
	once.Do(func() {
		var err error
		elasticsearchInstance, err = initElasticsearch(configPath)
		if err != nil {
			fmt.Printf("初始化 Elasticsearch 失败: %v", err)
		}
	})
	if elasticsearchInstance == nil {
		return nil, fmt.Errorf("elasticsearch 实例为空")
	}
	return elasticsearchInstance, nil
}

// CreateIndex Create 创建文档或索引
// @param index string: 索引名称
// @param docID string: 文档ID， 如果为空将自动生成
// @param document interface{}: 要创建的文档数据，可以是任意结构体或字典
// @return map[string]interface{}: Elasticsearch 返回的结果
// @return error: 如果创建失败，则返回错误
func (es *ElasticsearchClient) CreateIndex(index string, docID string, document interface{}) (map[string]interface{}, error) {
	// 将文档数据编码为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(document); err != nil {
		return nil, fmt.Errorf("编码文档失败: %v", err)
	}

	// 执行创建文档请求
	res, err := es.Client.Index(
		index,                                 // 索引名称
		&buf,                                  // 文档内容
		es.Client.Index.WithDocumentID(docID), // 文档ID，如果为空，Elasticsearch 会自动生成
		es.Client.Index.WithContext(es.Ctx),
	)

	if err != nil {
		return nil, fmt.Errorf("elasticsearch Create 操作失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	// 如果响应错误
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch 创建文档失败: %s", res.String())
	}

	// 解析响应结果
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 Elasticsearch 响应失败: %v", err)
	}

	return result, nil
}
