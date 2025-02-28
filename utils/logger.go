package utils

import (
	"github.com/sirupsen/logrus"
	loggerconfig "sky_ISService/config"
)

// LogInfo 记录手动日志
func LogInfo(message string) {
	loggerconfig.Logger.WithFields(logrus.Fields{
		"message": message, // 确保有 message 字段
		"level":   "info",
	}).Info(message)
}

// LogError 记录错误日志
func LogError(message string, err error) {
	loggerconfig.Logger.WithFields(logrus.Fields{
		"message": message,
		"error":   err.Error(),
		"level":   "error",
	}).Error(message)
}
