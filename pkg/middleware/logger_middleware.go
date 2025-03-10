package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sky_ISService/shared/elasticsearch"
	"time"
)

// LoggerMiddleware Gin 中间件：自动记录每个请求的日志
func LoggerMiddleware(indexName string, client *elasticsearch.ElasticsearchClient) gin.HandlerFunc {
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

		// 将日志文档插入到 Elasticsearch
		result, err := client.CreateIndex(indexName, "", logData)
		if err != nil {
			// 打印 Elasticsearch 错误
			fmt.Println("Elasticsearch 错误:", err.Error())
		} else {
			// 成功记录日志
			fmt.Println("请求日志已成功记录到 Elasticsearch")
			// 打印 Elasticsearch 返回的创建结果（此处可以根据实际返回数据进行处理）
			fmt.Println("文档创建成功:", result)
		}
	}
}
