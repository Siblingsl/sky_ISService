package utils

import "github.com/gin-gonic/gin"

// Response 封装统一的响应格式
type Response struct {
	Code    int         `json:"code"`    // 返回状态码
	Message string      `json:"message"` // 返回信息
	Data    interface{} `json:"data"`    // 返回的数据，可以是任意类型
}

// 定义成功和失败的状态码
const (
	SuccessCode = 200
	ErrorCode   = 500
)

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    SuccessCode,
		Message: "success",
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
