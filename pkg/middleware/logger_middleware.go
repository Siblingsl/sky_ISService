package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// LoggerMiddleware Gin 中间件：自动记录每个请求的日志
func LoggerMiddleware(indexName string, client *elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始时间
		start := time.Now()

		// 请求处理前
		c.Next()

		// 请求结束时间
		duration := time.Since(start)

		// 获取客户端的 IP 地址
		clientIP := c.ClientIP()

		// 获取请求的路由和方法
		method := c.Request.Method
		route := c.FullPath()

		// 获取请求的状态码
		status := c.Writer.Status()

		// 记录到 Elasticsearch（或任何你需要记录的内容）
		logData := map[string]interface{}{
			"方法":     method,
			"路由":     route,
			"客户端 IP": clientIP,
			"响应状态":   status,
			"耗时":     duration.String(),
			"时间":     time.Now().Format(time.RFC3339),
		}

		// 将文档转为 JSON
		docJSON, err := json.Marshal(logData)
		if err != nil {
			log.Fatalf("Error marshaling document: %s", err)
		}

		// 发送日志数据存入 Elasticsearch
		_, err = client.Index(
			indexName,                       // 索引名称
			bytes.NewReader(docJSON),        // 请求体内容
			client.Index.WithOpType("_doc"), // 如果使用文档类型，可以指定，但7.x及以上版本通常不需要
		)

		if err != nil {
			fmt.Println("Elasticsearch 错误:", err.Error()) // 打印更详细的错误信息
		} else {
			fmt.Println("请求日志已成功记录到 Elasticsearch")
		}
	}
}
