package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	logger "sky_ISService/utils"
	"time"
)

// LoggerMiddleware Gin 中间件：自动记录每个请求的日志
func LoggerMiddleware(client *elastic.Client) gin.HandlerFunc {
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

		// 打印日志到控制台，确认中间件被调用
		fmt.Println("正在记录日志：", logData)

		// 使用 Elasticsearch 客户端将日志数据存入 Elasticsearch
		_, err := client.Index().
			Index("logger-auth").
			BodyJson(logData).
			Do(context.Background())

		if err != nil {
			logger.LogError("日志记录失败: ", err)
			fmt.Println("Elasticsearch 错误:", err.Error()) // 打印更详细的错误信息
		} else {
			logger.LogInfo("请求日志已成功记录到 Elasticsearch")
		}
	}
}
