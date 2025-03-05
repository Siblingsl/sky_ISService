package utils

import "fmt"

// AppError 定义一个应用错误类型
type AppError struct {
	Code    int    // 错误码
	Message string // 错误信息
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewAppError 工厂方法，方便创建错误, 用于抛出自定义错误
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}
